package pubsub

import (
	"app/blockchain"
	"app/handlers"
	"app/models"
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"time"
    "encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Initializes pubsub client
func StartPubSubClient(ctx context.Context, db *gorm.DB, blockChainClient *blockchain.Client) (*pubsub.Client, error) {

    projectID := os.Getenv("PUBSUB_PROJECT")

    client, err := pubsub.NewClient(ctx, projectID)
    if err != nil {
        log.Printf("Failed to create Pub/Sub client: %v", err)
        return nil, err
    }

	return client, nil
}

func CreateTopicWithID(ctx context.Context, client *pubsub.Client, topicID string) (*pubsub.Topic, error) {

    // Get or create topic
    topic := client.Topic(topicID)
    exists, err := topic.Exists(ctx)
    if err != nil {
        log.Printf("Failed to check if topic exists: %v", err)
        // Continue anyway
    }
    
    if !exists || err != nil {
        log.Printf("Creating topic: %s", topicID)
        topic, err = client.CreateTopic(ctx, topicID)
        if err != nil {
            // If topic already exists, just use the existing one
            if strings.Contains(err.Error(), "AlreadyExists") {
                log.Printf("Topic already exists, using existing topic: %s", topicID)
                topic = client.Topic(topicID)
            } else {
                log.Printf("Failed to create topic: %v", err)
                return nil, err
            }
        }
    }
    return topic, nil
}


func SubscribeClient(ctx context.Context, client *pubsub.Client, topicID string, subscriptionID string) (*pubsub.Subscription, error) {

    // Ensure topic exists
    topic, err := CreateTopicWithID(ctx, client, topicID)
    if err != nil {
        return nil, err
    }

    // Get or create subscription
    sub := client.Subscription(subscriptionID)
    exists, err := sub.Exists(ctx)
    if err != nil {
        log.Printf("Failed to check if subscription exists: %v", err)
        // Continue anyway
    }

    if !exists || err != nil {
        log.Printf("Creating subscription: %s", subscriptionID)
        sub, err = client.CreateSubscription(ctx, subscriptionID, pubsub.SubscriptionConfig{
            Topic:       topic,
            AckDeadline: 20 * time.Second,
        })
        if err != nil {
            // If subscription already exists, just use the existing one
            if strings.Contains(err.Error(), "AlreadyExists") {
                log.Printf("Subscription already exists, using existing subscription: %s", subscriptionID)
                sub = client.Subscription(subscriptionID)
            } else {
                log.Printf("Failed to create subscription: %v", err)
                return nil, err
            }
        }
    }

    return sub, nil
}

func PublishNotification(ctx context.Context, client *pubsub.Client, notification []byte) error {
    topic := client.Topic("notifications")

    result := topic.Publish(ctx, &pubsub.Message{
        Data: notification,
    })

    id, err := result.Get(ctx)
    if err != nil {
        log.Printf("Failed to publish notification: %v", err)
        return err
    }

    log.Printf("Published notification with message ID: %s", id)
    log.Printf("Notification payload: %s", string(notification))
    return nil
}

func buildNotificationPayloadOrder(messageData []byte, db *gorm.DB, blockChainClient *blockchain.Client) []byte {

    // Parse the update to the order status model
    var order models.Orders
    err := json.Unmarshal(messageData, &order)
    if err != nil {
        log.Printf("Failed to unmarshal order update for notification: %v", err)
        return nil
    }

    notification := fmt.Sprintf(`{
       "user_id": %d, 
        "type": "sms", 
        "title": "New Order Created", 
        "payload": "Order with ID %d has been created.", 
        "hyperlink": "https://tracking-status-frontend-edneicy3ca-ew.a.run.app/order/%d", 
        "created_at": "` + time.Now().Format(time.RFC3339) + `" 
    }`, order.Customer_ID, order.Id)

    return []byte(notification)
}

func buildNotificationPayloadStatus(messageData []byte, db *gorm.DB, blockChainClient *blockchain.Client) []byte {

    // Parse the update to the order status model
    var order_update models.OrderStatusHistory
    err := json.Unmarshal(messageData, &order_update)
    if err != nil {
        log.Printf("Failed to unmarshal order update for notification: %v", err)
        return nil
    }

    // Create handler to get user ID
    handler := handlers.OrderHandler{DB: db, Client: blockChainClient}
    userID, err := handlers.GetUserIDByOrderID(handler.DB, order_update.Order_ID)
    if err != nil {
        log.Printf("Failed to get user ID for order ID %d: %v", order_update.Order_ID, err)
        return nil
    }

    notification := fmt.Sprintf(`{
       "user_id": %d, 
        "type": "sms", 
        "title": "Order Status Update", 
        "payload": %s, 
        "hyperlink": "https://tracking-status-frontend-edneicy3ca-ew.a.run.app/order/%d", 
        "created_at": "` + time.Now().Format(time.RFC3339) + `" 
    }`, userID, order_update.Order_Status, order_update.Order_ID)

    return []byte(notification)
}

func StartListener(ctx context.Context, client *pubsub.Client, sub *pubsub.Subscription, db *gorm.DB, blockChainClient *blockchain.Client) error {
	// Handler
    orderStatusHistory := handlers.OrderStatusHistoryHandler{DB: db, Client: blockChainClient}
	fmt.Println("Listening for order status update messages...")
	go func() {
		err := sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
			fmt.Printf("RAW MESSAGE %s\n", string(m.Data))
			m.Ack()

			// Create a mock Gin context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			
			// Create a request with the PubSub data as JSON body
			c.Request = httptest.NewRequest("POST", "/order/history/add", bytes.NewReader(m.Data))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call the existing handler
			orderStatusHistory.AddOrderUpdate(c)

			// Check the response
			fmt.Printf("Response Status: %d, Body: %s\n", w.Code, w.Body.String())
			
			if w.Code >= 400 {
				fmt.Printf("Failed to save update: Status %d, Response: %s\n", w.Code, w.Body.String())
				m.Nack()
			} else {
				fmt.Printf("Successfully saved update via HTTP handler\n")

                // Send notification for status update
                notificationPayload := buildNotificationPayloadStatus(m.Data, db, blockChainClient)
                if err:= PublishNotification(ctx, client, notificationPayload); err != nil {
                    log.Printf("Failed to publish notification: %v", err)
                }
			}
    	})

		if err != nil {
			log.Printf("PubSub listener stopped: %v", err)
		}

	}()
	return nil
}

func StartListenerOrders(ctx context.Context, client *pubsub.Client, sub *pubsub.Subscription, db *gorm.DB, blockChainClient *blockchain.Client) error {
	// Handler
    orderHandler := handlers.OrderHandler{DB: db, Client: blockChainClient}
	fmt.Println("Listening for new order messages...")
	go func() {
		err := sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
			fmt.Printf("RAW MESSAGE %s\n", string(m.Data))
			m.Ack()

			// Create a mock Gin context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			
			// Create a request with the PubSub data as JSON body
			c.Request = httptest.NewRequest("POST", "/order/add", bytes.NewReader(m.Data))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call the existing handler
			orderHandler.AddOrder(c)

			// Check the response
			fmt.Printf("Response Status: %d, Body: %s\n", w.Code, w.Body.String())
			
			if w.Code >= 400 {
				fmt.Printf("Failed to save update: Status %d, Response: %s\n", w.Code, w.Body.String())
				m.Nack()
			} else {
				fmt.Printf("Successfully saved update via HTTP handler\n")

                // Send notification for status update
                notificationPayload := buildNotificationPayloadOrder(m.Data, db, blockChainClient)
                if err:= PublishNotification(ctx, client, notificationPayload); err != nil {
                    log.Printf("Failed to publish notification: %v", err)
                }
			}
    	})

		if err != nil {
			log.Printf("PubSub listener stopped: %v", err)
		}

	}()
	return nil
}