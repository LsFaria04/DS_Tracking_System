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

    "google.golang.org/protobuf/proto"
	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Topic interface for Pub/Sub operations
type Topic interface {
	Publish(ctx context.Context, msg *pubsub.Message) *pubsub.PublishResult
	Exists(ctx context.Context) (bool, error)
	Delete(ctx context.Context) error
}

// Subscription interface for Pub/Sub operations
type Subscription interface {
	Receive(ctx context.Context, f func(context.Context, *pubsub.Message)) error
	Exists(ctx context.Context) (bool, error)
	Delete(ctx context.Context) error
}

// PubSubClient interface defines the interface for Pub/Sub operations
type PubSubClient interface {
	Topic(id string) *pubsub.Topic
	Subscription(id string) *pubsub.Subscription
	CreateTopic(ctx context.Context, id string) (*pubsub.Topic, error)
	CreateSubscription(ctx context.Context, id string, cfg pubsub.SubscriptionConfig) (*pubsub.Subscription, error)
	Close() error
}

// Initializes pubsub client
func StartPubSubClient(ctx context.Context, db *gorm.DB, blockChainClient *blockchain.Client) (*pubsub.Client, error) {

    projectID := os.Getenv("PUBSUB_PROJECT")
    
    if projectID == "" {
        log.Printf("Failed to create Pub/Sub client: PUBSUB_PROJECT env var is empty")
        return nil, fmt.Errorf("PUBSUB_PROJECT environment variable is not set")
    }

    client, err := pubsub.NewClient(ctx, projectID)
    if err != nil {
        log.Printf("Failed to create Pub/Sub client: %v", err)
        return nil, err
    }

	return client, nil
}

func CreateTopicWithID(ctx context.Context, client *pubsub.Client, topicID string) (*pubsub.Topic, error) {
    
    if client == nil {
        log.Printf("Failed to create topic: client is nil")
        return nil, fmt.Errorf("pubsub client is nil")
    }
    
    if topicID == "" {
        log.Printf("Failed to create topic: topicID is empty")
        return nil, fmt.Errorf("topicID is empty")
    }

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

    if client == nil {
        log.Printf("Failed to create subscription: client is nil")
        return nil, fmt.Errorf("pubsub client is nil")
    }
    
    if topicID == "" {
        log.Printf("Failed to create subscription: topicID is empty")
        return nil, fmt.Errorf("topicID is empty")
    }
    
    if subscriptionID == "" {
        log.Printf("Failed to create subscription: subscriptionID is empty")
        return nil, fmt.Errorf("subscriptionID is empty")
    }

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
    
    if client == nil {
        log.Printf("Failed to publish notification: client is nil")
        return fmt.Errorf("pubsub client is nil")
    }
    
    if notification == nil || len(notification) == 0 {
        log.Printf("Failed to publish notification: notification is empty")
        return fmt.Errorf("notification payload is empty")
    }
    
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

    messageNotification := &NotificationRequest{}
    if err := proto.Unmarshal(notification, messageNotification); err != nil {
        log.Printf("Failed to unmarshal notification for logging: %v", err)
        return nil
    }
    log.Printf("Notification payload: %v", messageNotification)
    return nil
}

func buildNotificationPayloadOrder(messageData []byte, db *gorm.DB, blockChainClient *blockchain.Client) []byte {

    if messageData == nil || len(messageData) == 0 {
        log.Printf("Failed to build notification: messageData is empty")
        return nil
    }

    // Parse the update to the order status model
    var order models.Orders
    err := json.Unmarshal(messageData, &order)
    if err != nil {
        log.Printf("Failed to unmarshal order update for notification: %v", err)
        return nil
    }
    
    log.Printf("Order with id %d created for customer %d", order.Id, order.Customer_ID)

    notification := &NotificationRequest{ 
        UserId:     fmt.Sprintf("%d", order.Customer_ID), 
        Type:       "sms",
        Title:      "New Order Created",
        Payload:    fmt.Sprintf("Order with ID %d has been created.", order.Id),
        Hyperlink:  fmt.Sprintf("https://tracking-status-frontend-edneicy3ca-ew.a.run.app/order/%d", order.Id),
        CreatedAt:  time.Now().Format(time.RFC3339),
    }

    // Encrypt to protobuf
    protoData, err := proto.Marshal(notification)
    if err != nil {
        log.Printf("Failed to marshal notification to protobuf: %v", err)
        return nil
    }

    return protoData
}

func buildNotificationPayloadStatus(messageData []byte, db *gorm.DB, blockChainClient *blockchain.Client) []byte {

    if messageData == nil || len(messageData) == 0 {
        log.Printf("Failed to build notification: messageData is empty")
        return nil
    }

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

    notification := &NotificationRequest{ 
        UserId: fmt.Sprintf("%d", userID), 
        Type: "sms", 
        Title: "Order Status Update", 
        Payload: fmt.Sprintf("Your order status has changed to: %s", order_update.Order_Status), 
        Hyperlink: fmt.Sprintf("https://tracking-status-frontend-edneicy3ca-ew.a.run.app/order/%d", order_update.Order_ID),
        CreatedAt: time.Now().Format(time.RFC3339),
    }

    // Encrypt to protobuf
    protoData, err := proto.Marshal(notification)
    if err != nil {
        log.Printf("Failed to marshal notification to protobuf: %v", err)
        return nil
    }

    return protoData
}

func StartListener(ctx context.Context, client *pubsub.Client, sub *pubsub.Subscription, db *gorm.DB, blockChainClient *blockchain.Client) error {
    
    if client == nil {
        return fmt.Errorf("pubsub client is nil")
    }
    
    if sub == nil {
        return fmt.Errorf("subscription is nil")
    }
    
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
    
    if client == nil {
        return fmt.Errorf("pubsub client is nil")
    }
    
    if sub == nil {
        return fmt.Errorf("subscription is nil")
    }
    
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