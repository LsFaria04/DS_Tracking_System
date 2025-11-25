package main

import (
	"app/blockchain"
	"app/routes"
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
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// configure the database connection using gorm
func configDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}

func configBlockChainClient() (*blockchain.Client, error) {
	rpcURL := os.Getenv("BLOCKCHAIN_RPC_URL")
	if rpcURL == "" {
		return nil, nil
	}
	//connect to the block chain client
	client, err := blockchain.NewClient()
	if err != nil {

		return nil, err
	}
	return client, nil
}

// Configure the router that will be used for the API
func configRouter(db *gorm.DB) (*gin.Engine, *blockchain.Client, error) {
	router := gin.Default()

	// Configure CORS middleware (Allow frontend and localhost)

	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// Allow localhost for development

			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			// Allow any .run.app domain (Cloud Run)
			if strings.HasSuffix(origin, ".run.app") {
				return true
			}
			return false
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	blockChainClient, err := configBlockChainClient()

	if err != nil {
		return nil,nil, err
	}

	//registers the routes
	routes.RegisterRoutes(router, db, blockChainClient)
	return router, blockChainClient, nil
}

func testPubSub(db *gorm.DB, blockChainClient *blockchain.Client) error {
    ctx := context.Background()

    projectID := os.Getenv("PUBSUB_PROJECT")
    topicID := "orders_status"
    subscriptionID := "orders_status-sub" 

    client, err := pubsub.NewClient(ctx, projectID)
    if err != nil {
        log.Printf("Failed to create Pub/Sub client: %v", err)
        return err
    }
    defer client.Close()

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
                return err
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
                return err
            }
        }
    }

    // Handler
    orderStatusHistory := handlers.OrderStatusHistoryHandler{DB: db, Client: blockChainClient}

    fmt.Println("Listening for messages...")

    err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
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

    return nil
}

func main() {
	db, err := configDB()

	if err != nil {
		log.Printf("Error while conecting to the database: %v", err)
		return
	}

	router, blockChainClient, err := configRouter(db)

	if err != nil {
		log.Printf("Error while configuring the routing: %v", err)
		return
	}

	// Test PubSub emulator in another goroutine
	go func(){
		if err := testPubSub(db, blockChainClient); err != nil {
			log.Printf(" PubSub test failed: %v", err)
		}
	}()

	router.Run(":8080") // listens on 0.0.0.0:8080 by default
}
