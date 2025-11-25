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

// Test PubSub emulator connection to mock_courier
func testPubSub() error {

	ctx := context.Background()

    projectID := os.Getenv("PUBSUB_PROJECT")
    subscriptionID := "orders_status-sub" 

    client, err := pubsub.NewClient(ctx, projectID)
    if err != nil {
        log.Fatalf("Failed to create Pub/Sub client: %v", err)
    }
    defer client.Close()

    sub := client.Subscription(subscriptionID)

    fmt.Println("Listening for messages...")

	err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
        fmt.Printf("Received message: %s\n", string(m.Data))
        m.Ack() // acknowledge the message so it's not redelivered
    })
    if err != nil {
        log.Fatalf("Failed to receive messages: %v", err)
    }

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
		log.Printf(" PubSub test failed: %v", err)
		log.Println("Continuing without PubSub...")
	}

	router, err := configRouter(db)

	if err != nil {
		log.Printf("Error while configuring the routing: %v", err)
		return
	}

	router.Run(":8080") // listens on 0.0.0.0:8080 by default
}
