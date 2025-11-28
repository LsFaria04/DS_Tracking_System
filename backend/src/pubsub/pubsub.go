package pubsub

import (
	"app/blockchain"
	"app/handlers"
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"time"

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

func SubscribeClient(ctx context.Context, client *pubsub.Client, topicID string, subscriptionID string) (*pubsub.Subscription, error) {

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

    // Get or create subscription
    sub := client.Subscription(subscriptionID)
    exists, err = sub.Exists(ctx)
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

func StartListener(ctx context.Context, sub *pubsub.Subscription, db *gorm.DB, blockChainClient *blockchain.Client) error {
	// Handler
    orderStatusHistory := handlers.OrderStatusHistoryHandler{DB: db, Client: blockChainClient}
	fmt.Println("Listening for messages...")
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
			}
    	})

		if err != nil {
			log.Printf("PubSub listener stopped: %v", err)
		}

	}()
	return nil
}