package main

import (
	"app/blockchain"
	"app/routes"
	"context"
	"fmt"
	"log"
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
func configRouter(db *gorm.DB) (*gin.Engine, error) {
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
		return nil, err
	}

	//registers the routes
	routes.RegisterRoutes(router, db, blockChainClient)
	return router, nil
}

// Test PubSub emulator connection
func testPubSub() error {

	pubsubEmulatorHost := os.Getenv("PUBSUB_EMULATOR_HOST")
	if pubsubEmulatorHost == "" {
		log.Println("PUBSUB_EMULATOR_HOST not set, skipping PubSub test")
		return nil
	}

	log.Printf("Testing PubSub emulator at: %s", pubsubEmulatorHost)

	ctx := context.Background()

	// Use the project ID from environment or default to "madeinportugal"
	projectID := os.Getenv("PUBSUB_PROJECT")
	if projectID == "" {
		projectID = "madeinportugal"
	}

	// Create PubSub client
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to create pubsub client: %v", err)
	}
	defer client.Close()

	// Test publishing to "orders" topic
	topic := client.Topic("orders")
	defer topic.Stop()

	// Check if topic exists, if not create it
	exists, err := topic.Exists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if topic exists: %v", err)
	}

	if !exists {
		log.Println("Topic 'orders' does not exist, creating...")
		topic, err = client.CreateTopic(ctx, "orders")
		if err != nil {
			return fmt.Errorf("failed to create topic: %v", err)
		}
	}

	// Publish a test message
	testMessage := fmt.Sprintf("Test message from backend - %s", time.Now().Format(time.RFC3339))
	result := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(testMessage),
		Attributes: map[string]string{
			"source": "tracking-status",
			"type":   "health-check",
		},
	})

	// Block until the message is published
	msgID, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	log.Printf("PubSub test successful! Message ID: %s", msgID)
	log.Printf("Published message: %s", testMessage)

	return nil
}

func main() {
	db, err := configDB()

	if err != nil {
		log.Printf("Error while conecting to the database: %v", err)
		return
	}

	// Test PubSub emulator if in development mode
	if err := testPubSub(); err != nil {
		log.Printf("⚠️  PubSub test failed: %v", err)
		log.Println("Continuing without PubSub...")
	}

	router, err := configRouter(db)

	if err != nil {
		log.Printf("Error while configuring the routing: %v", err)
		return
	}

	router.Run(":8080") // listens on 0.0.0.0:8080 by default
}
