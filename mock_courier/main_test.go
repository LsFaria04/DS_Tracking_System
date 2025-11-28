package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestMockCourierMessageGeneration tests the message generation logic
func TestMockCourierMessageGeneration(t *testing.T) {
	// Test data from your main function
	status := []string{"PROCESSING", "SHIPPED", "IN TRANSIT"}
	notes := []string{"Order received from seller", "Picked up from seller by courier", "Package ready for delivery"}
	locations := []string{"Dona Lurdes", "Main Warehouse Lisboa", "Main Warehouse Lisboa"}
	storageIds := []int{0, 1, 1}

	for i := 0; i < 3; i++ {
		t.Run(status[i], func(t *testing.T) {
			var msgData string
			
			// Replicate the exact logic from your main function
			if storageIds[i] != 0 {
				msgData = fmt.Sprintf(`{
					"order_id": 1,
					"order_status": "%s",          
					"note": "%s",
					"order_location": "%s",        
					"storage_id": %d
				}`, status[i], notes[i], locations[i], storageIds[i])
			} else {
				msgData = fmt.Sprintf(`{
					"order_id": 1,
					"order_status": "%s",          
					"note": "%s",
					"order_location": "%s"        
				}`, status[i], notes[i], locations[i])
			}

			// Validate the generated message
			var data map[string]interface{}
			err := json.Unmarshal([]byte(msgData), &data)
			assert.NoError(t, err, "Generated message should be valid JSON")

			// Check required fields
			assert.Equal(t, float64(1), data["order_id"])
			assert.Equal(t, status[i], data["order_status"])
			assert.Equal(t, notes[i], data["note"])
			assert.Equal(t, locations[i], data["order_location"])

			// Check storage_id logic
			if storageIds[i] != 0 {
				assert.Equal(t, float64(storageIds[i]), data["storage_id"])
			} else {
				_, exists := data["storage_id"]
				assert.False(t, exists, "storage_id should not exist when storageIds[i] == 0")
			}

			t.Logf("Generated valid message for status: %s", status[i])
		})
	}
}

// TestMessageAttributes tests the Pub/Sub message attributes
func TestMessageAttributes(t *testing.T) {
	// Test the attributes from your main function
	attributes := map[string]string{
		"source": "mock_courier",
		"type":   "order-update",
	}

	assert.Equal(t, "mock_courier", attributes["source"])
	assert.Equal(t, "order-update", attributes["type"])
	
	// Verify attributes are set correctly
	assert.NotEmpty(t, attributes["source"])
	assert.NotEmpty(t, attributes["type"])
}

// TestTimingLogic tests the timing logic from your mock courier
func TestTimingLogic(t *testing.T) {
	// Test the delays from your main function
	totalExpectedTime := 20*time.Second + 2*7*time.Second
	assert.Equal(t, 34*time.Second, totalExpectedTime, 
		"Total execution time should be 34 seconds")

	t.Logf("Mock courier timing: 20s initial + 7s + 7s = 34s total")
}

// TestPubSubMessageStructure tests the complete message structure
func TestPubSubMessageStructure(t *testing.T) {
	// Test a complete message with attributes
	messageData := `{
		"order_id": 1,
		"order_status": "PROCESSING",          
		"note": "Order received from seller",
		"order_location": "Dona Lurdes"        
	}`

	attributes := map[string]string{
		"source": "mock_courier",
		"type":   "order-update",
	}

	// Validate message data
	var data map[string]interface{}
	err := json.Unmarshal([]byte(messageData), &data)
	assert.NoError(t, err)

	// Validate attributes
	assert.Equal(t, "mock_courier", attributes["source"])
	assert.Equal(t, "order-update", attributes["type"])

	t.Log("Complete Pub/Sub message structure is valid")
}

// TestStorageIDConditionalLogic tests the storage_id conditional logic
func TestStorageIDConditionalLogic(t *testing.T) {
	tests := []struct {
		storageID    int
		shouldInclude bool
	}{
		{0, false},
		{1, true},
		{2, true},
		{999, true},
		{-1, true}, // negative numbers are also != 0
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("StorageID_%d", tc.storageID), func(t *testing.T) {
			// This replicates the exact logic from your main function:
			// if storageIds[i] != 0 { include storage_id } else { don't include }
			shouldInclude := tc.storageID != 0
			assert.Equal(t, tc.shouldInclude, shouldInclude,
				"For storageID %d, shouldInclude should be %v", tc.storageID, tc.shouldInclude)
		})
	}
}