package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "cloud.google.com/go/pubsub"
)

func main() {
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

    topicName := "orders_status"
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

    status := []string{"PROCESSING", "SHIPPED", "IN TRANSIT"}
    notes := []string{"Order received from seller", "Picked up from seller by courier", "Package ready for delivery"}
    locations := []string{"Dona Lurdes", "Main Warehouse Lisboa", "Main Warehouse Lisboa"}
    storageIds := []int{0, 1, 1}

    time.Sleep(20 * time.Second)
    for i := 0; i < 3; i++ {
        time.Sleep(7 * time.Second)

        var res *pubsub.PublishResult

        if storageIds[i] != 0 {

            msgData := fmt.Sprintf(`{
                "order_id": 1,
                "order_status": "%s",          
                "note": "%s",
                "order_location": "%s",        
                "storage_id": %d
            }`, status[i], notes[i], locations[i], storageIds[i])

            res = topic.Publish(ctx, &pubsub.Message{
                Data: []byte(msgData),
                Attributes: map[string]string{
                    "source": "mock_courier",
                    "type":   "order-update",
                },
            })
        } else {
            msgData := fmt.Sprintf(`{
                "order_id": 1,
                "order_status": "%s",          
                "note": "%s",
                "order_location": "%s"        
            }`, status[i], notes[i], locations[i])

            res = topic.Publish(ctx, &pubsub.Message{
                Data: []byte(msgData),
                Attributes: map[string]string{
                    "source": "mock_courier",
                    "type":   "order-update",
                },
            })
        }

        id, err := res.Get(ctx)
        if err != nil {
            log.Printf("Failed to publish message: %v", err)
        } else {
            log.Printf("Published message ID: %s, status: %s", id, status[i])
        }
    }
}
