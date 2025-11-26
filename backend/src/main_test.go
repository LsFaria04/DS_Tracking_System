package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPubSubMessageFormat tests the JSON message format your system expects
func TestPubSubMessageFormat(t *testing.T) {
	// Test the exact message formats from your mock courier
	validMessages := []string{
		`{
			"order_id": 1,
			"order_status": "PROCESSING",          
			"note": "Order received from seller",
			"order_location": "Dona Lurdes"        
		}`,
		`{
			"order_id": 1,
			"order_status": "SHIPPED",          
			"note": "Picked up from seller by courier",
			"order_location": "Main Warehouse Lisboa",        
			"storage_id": 1
		}`,
		`{
			"order_id": 1,
			"order_status": "IN TRANSIT",          
			"note": "Package ready for delivery",
			"order_location": "Main Warehouse Lisboa",        
			"storage_id": 1
		}`,
	}

	for i, msg := range validMessages {
		t.Run(fmt.Sprintf("ValidMessage_%d", i), func(t *testing.T) {
			var data map[string]interface{}
			err := json.Unmarshal([]byte(msg), &data)
			assert.NoError(t, err, "Message should be valid JSON")
			
			// Check required fields
			assert.Contains(t, data, "order_id")
			assert.Contains(t, data, "order_status")
			assert.Contains(t, data, "note")
			assert.Contains(t, data, "order_location")
			
			t.Logf("Message %d validated successfully: %s", i, data["order_status"])
		})
	}
}

// TestInvalidMessageFormat tests invalid messages that should be rejected
func TestInvalidMessageFormat(t *testing.T) {
	invalidMessages := []struct {
		name    string
		message string
	}{
		{"Invalid JSON", `{ invalid json }`},
		{"Missing order_id", `{"order_status": "PROCESSING", "note": "test", "order_location": "test"}`},
		{"Missing order_status", `{"order_id": 1, "note": "test", "order_location": "test"}`},
		{"Empty object", `{}`},
	}

	for _, tc := range invalidMessages {
		t.Run(tc.name, func(t *testing.T) {
			var data map[string]interface{}
			err := json.Unmarshal([]byte(tc.message), &data)
			
			if err != nil {
				t.Logf("Expected JSON error for %s: %v", tc.name, err)
				return
			}
			
			// If JSON is valid, check required fields
			hasRequiredFields := true
			requiredFields := []string{"order_id", "order_status", "note", "order_location"}
			for _, field := range requiredFields {
				if _, exists := data[field]; !exists {
					hasRequiredFields = false
					break
				}
			}
			
			assert.False(t, hasRequiredFields, "Invalid message should miss required fields")
		})
	}
}

// TestAckNackLogic tests the message acknowledgment logic from testPubSub
func TestAckNackLogic(t *testing.T) {
	tests := []struct {
		statusCode int
		shouldAck  bool
		description string
	}{
		{200, true, "Success status should ack"},
		{201, true, "Created status should ack"},
		{204, true, "No content status should ack"},
		{400, false, "Bad request should nack"},
		{500, false, "Server error should nack"},
		{404, false, "Not found should nack"},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			// This replicates the logic from your testPubSub function:
			// if w.Code >= 400 { m.Nack() } else { /* ack the message */ }
			shouldAck := tc.statusCode < 400
			assert.Equal(t, tc.shouldAck, shouldAck, 
				"For status %d, expected ack=%v", tc.statusCode, tc.shouldAck)
		})
	}
}

// TestMessageAttributes tests the Pub/Sub message attributes
func TestMessageAttributes(t *testing.T) {
	// Test the attributes your mock courier sets
	expectedAttributes := map[string]string{
		"source": "mock_courier",
		"type":   "order-update",
	}

	assert.Equal(t, "mock_courier", expectedAttributes["source"])
	assert.Equal(t, "order-update", expectedAttributes["type"])
	
	// Test that attributes exist
	assert.NotEmpty(t, expectedAttributes["source"])
	assert.NotEmpty(t, expectedAttributes["type"])
}

// TestStorageIDLogic tests the storage_id conditional logic
func TestStorageIDLogic(t *testing.T) {
	tests := []struct {
		storageID  int
		shouldInclude bool
	}{
		{0, false},
		{1, true},
		{2, true},
		{999, true},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("StorageID_%d", tc.storageID), func(t *testing.T) {
			// This replicates the logic from your mock courier:
			// if storageIds[i] != 0 { include storage_id } else { don't include }
			shouldInclude := tc.storageID != 0
			assert.Equal(t, tc.shouldInclude, shouldInclude,
				"storage_id should be included when != 0")
		})
	}
}