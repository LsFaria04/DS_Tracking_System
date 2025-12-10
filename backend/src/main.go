package main

import (
	"app/blockchain"
	"app/routes"
    "app/pubsub"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"


    googlepubsub "cloud.google.com/go/pubsub"
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

func configPubSubClient(db *gorm.DB, blockChainClient *blockchain.Client, topicID []string, subscriptionID []string) (*googlepubsub.Client, []*googlepubsub.Subscription, error) {
    ctx := context.Background()
    pubsubClient, err := pubsub.StartPubSubClient(ctx, db, blockChainClient)

    if err != nil {
        return nil, nil, err
    }

    subs := make([]*googlepubsub.Subscription, len(subscriptionID))
    for i, subID := range subscriptionID {
        subs[i], err = pubsub.SubscribeClient(ctx, pubsubClient, topicID[i], subID)
        if err != nil {
            return nil, nil, err
        } 
    }
    return pubsubClient, subs, nil
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
			if strings.HasSuffix(origin, ".run.app") || strings.HasPrefix(origin, "https://frontend.madeinportugal.store") {
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

	// Create and start the Pub/Sub client
    ctx := context.Background()

	// Configure Pub/Sub client and subscriptions 
    client, subs, err := configPubSubClient(db, blockChainClient, []string{}, []string{}) // Empty topic and subscription IDs since they are not created yet by the checkout team

	if err != nil {
        log.Printf("Error configuring PubSub: %v", err)
        // Continue without PubSub
        log.Printf("Continuing without PubSub functionality")
    } else if client != nil {
		// List all topics and subscriptions (for debugging)
		pubsub.ListAllTopics(ctx, client)
		pubsub.ListAllSubscriptions(ctx, client)

		// Create notifications topic to send notifications whenever status updates occur
		//_, err = pubsub.CreateTopicWithID(ctx, client, "tracking-notifications")  -  Notifications are already created
        defer client.Close() // Close client when main exits
        
        // Only start listeners if we have subscriptions
        if len(subs) > 0 {
            err = pubsub.StartListener(ctx, client, subs[0], db, blockChainClient)
            if err != nil {
                log.Printf("Error starting PubSub listener for order status: %v", err)
            }
        } else {
            log.Printf("No subscriptions configured, skipping listeners")
        }
        
        if len(subs) > 1 {
            err = pubsub.StartListenerOrders(ctx, client, subs[1], db, blockChainClient)
            if err != nil {
                log.Printf("Error starting PubSub listener for checkout orders: %v", err)
            }
        }

        //pubsub.TestOrdersPubSub()  // Uncomment to test order publishing
    }
	

	router.Run(":8080") // listens on 0.0.0.0:8080 by default
}
