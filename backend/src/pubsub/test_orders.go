package pubsub

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "cloud.google.com/go/pubsub"
)

func TestOrdersPubSub() {
    ctx := context.Background()

    pubsubEmulatorHost := os.Getenv("PUBSUB_EMULATOR_HOST")
    if pubsubEmulatorHost != "" {
        log.Printf("Using Pub/Sub emulator at %s", pubsubEmulatorHost)
    }

    projectID := os.Getenv("GCP_PROJECT_ID")
    if projectID == "" {
        projectID = "madeinportugal"
    }

    client, err := pubsub.NewClient(ctx, projectID)
    if err != nil {
        log.Fatalf("Failed to create Pub/Sub client: %v", err)
    }
    defer client.Close()

    topicName := "checkout_orders"
    topic := client.Topic(topicName)
    exists, err := topic.Exists(ctx)
    if err != nil {
        log.Fatalf("Failed to check if topic exists: %v", err)
    }
    if !exists {
        if _, err := client.CreateTopic(ctx, topicName); err != nil {
            log.Fatalf("Failed to create topic: %v", err)
        }
        topic = client.Topic(topicName)
    }
    defer topic.Stop()

	seller_address, seller_latitude, seller_longitude, delivery_address, delivery_latitude, delivery_longitude := "Dona Lurdes, Rua Dom Afonso Henriques 12, 2800-012 Almada, Portugal", 38.6780, -9.1580, "Rua Padre Joaquim Alves Correia 5, 1990-152 Lisboa, Portugal", 38.7680, -9.1000
    time.Sleep(5 * time.Second)

	msgData := fmt.Sprintf(`{
		"customer_id": 101,
		"seller_id": 501,
		"seller_address": "%s",
		"seller_latitude": %f,
		"seller_longitude": %f,
		"delivery_address": "%s",
		"delivery_latitude": %f,
		"delivery_longitude": %f,
		"products": [
			{	
				"product_id": 32865210,
				"quantity": 2
			},
			{
				"product_id": 32865209,
				"quantity": 1
			}
		]
	}`, seller_address, seller_latitude, seller_longitude, delivery_address, delivery_latitude, delivery_longitude)       

	res := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(msgData),
		Attributes: map[string]string{
			"source": "pubsub-test",
			"type":   "new-order",

		},
	})

	id, err := res.Get(ctx)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	} else {
		log.Printf("Published message ID: %s, Data: %s", id, msgData)
	}
}

