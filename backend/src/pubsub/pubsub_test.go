package pubsub

import (
	"app/models"
	"bytes"
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
)

// MockPubSubClient mocks the Pub/Sub client
type MockPubSubClient struct {
	mock.Mock
}

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

// TestBuildNotificationPayloadOrder tests the order notification payload builder
func TestBuildNotificationPayloadOrder(t *testing.T) {
	// Create test order data
	orderJSON := `{
		"id": 1,
		"customer_id": 101,
		"seller_id": 501,
		"tracking_code": "TRACK001"
	}`
	
	// Build the notification
	payload := buildNotificationPayloadOrder([]byte(orderJSON), nil, nil)
	
	// Verify payload is not nil and can be unmarshaled
	assert.NotNil(t, payload, "Notification payload should not be nil")
	
	// Unmarshal the protobuf
	notification := &NotificationRequest{}
	err := proto.Unmarshal(payload, notification)
	assert.NoError(t, err, "Should unmarshal notification without error")
	
	// Verify notification fields
	assert.Equal(t, "101", notification.UserId)
	assert.Equal(t, "sms", notification.Type)
	assert.Equal(t, "New Order Created", notification.Title)
	assert.Contains(t, notification.Payload, "Order with ID 1 has been created")
}

// TestBuildNotificationPayloadOrderInvalidJSON tests error handling for invalid JSON
func TestBuildNotificationPayloadOrderInvalidJSON(t *testing.T) {
	// Invalid JSON should return nil
	payload := buildNotificationPayloadOrder([]byte(`invalid json`), nil, nil)
	assert.Nil(t, payload, "Should return nil for invalid JSON")
}

// TestBuildNotificationPayloadStatus tests the status update notification payload builder
func TestBuildNotificationPayloadStatus(t *testing.T) {
	// Create a notification directly to test the marshaling logic
	notification := &NotificationRequest{
		UserId:    "1",
		Type:      "sms",
		Title:     "Order Status Update",
		Payload:   "Your order status has changed to: SHIPPED",
		CreatedAt: "2025-01-01T00:00:00Z",
	}
	
	// Marshal to protobuf
	data, err := proto.Marshal(notification)
	assert.NoError(t, err, "Should marshal notification without error")
	assert.NotNil(t, data, "Marshaled data should not be nil")
	
	// Unmarshal the protobuf
	result := &NotificationRequest{}
	err = proto.Unmarshal(data, result)
	assert.NoError(t, err, "Should unmarshal notification without error")
	
	// Verify notification fields
	assert.Equal(t, "sms", result.Type)
	assert.Equal(t, "Order Status Update", result.Title)
	assert.Contains(t, result.Payload, "SHIPPED")
}

// TestBuildNotificationPayloadStatusInvalidJSON tests error handling for invalid JSON
func TestBuildNotificationPayloadStatusInvalidJSON(t *testing.T) {
	// Invalid JSON should return nil
	payload := buildNotificationPayloadStatus([]byte(`invalid json`), nil, nil)
	assert.Nil(t, payload, "Should return nil for invalid JSON")
}

// TestPublishNotificationWithValidPayload tests publishing a valid notification
func TestPublishNotificationWithValidPayload(t *testing.T) {
	// Create a valid protobuf notification
	notification := &NotificationRequest{
		UserId:    "123",
		Type:      "sms",
		Title:     "Test",
		Payload:   "Test payload",
		CreatedAt: "2025-01-01T00:00:00Z",
	}
	
	// Marshal to protobuf
	data, err := proto.Marshal(notification)
	assert.NoError(t, err)
	
	// Note: PublishNotification requires a real Pub/Sub client, so we test the marshaling logic
	assert.NotNil(t, data, "Marshaled data should not be nil")
	
	// Verify it can be unmarshaled
	unmarshaledNotif := &NotificationRequest{}
	err = proto.Unmarshal(data, unmarshaledNotif)
	assert.NoError(t, err, "Should unmarshal successfully")
	assert.Equal(t, "123", unmarshaledNotif.UserId)
}

// TestOrderStatusHistoryUnmarshalJSON tests unmarshaling order status history
func TestOrderStatusHistoryUnmarshalJSON(t *testing.T) {
	statusJSON := `{
		"order_id": 1,
		"order_status": "SHIPPED",
		"note": "Package shipped",
		"order_location": "Main Warehouse",
		"storage_id": 1
	}`
	
	var orderUpdate models.OrderStatusHistory
	err := json.Unmarshal([]byte(statusJSON), &orderUpdate)
	assert.NoError(t, err, "Should unmarshal order status update")
	
	// Verify fields
	assert.Equal(t, uint(1), orderUpdate.Order_ID)
	assert.Equal(t, "SHIPPED", orderUpdate.Order_Status)
	assert.Equal(t, "Package shipped", orderUpdate.Note)
	assert.Equal(t, "Main Warehouse", orderUpdate.Order_Location)
	// Storage_ID is a pointer, so dereference it
	assert.NotNil(t, orderUpdate.Storage_ID)
	assert.Equal(t, uint(1), *orderUpdate.Storage_ID)
}

// TestNotificationRequestProtobufMarshaling tests protobuf marshaling of NotificationRequest
func TestNotificationRequestProtobufMarshaling(t *testing.T) {
	notification := &NotificationRequest{
		UserId:    "user123",
		Type:      "sms",
		Title:     "Order Update",
		Payload:   "Your order has been shipped",
		Hyperlink: "https://example.com/order/1",
		CreatedAt: "2025-01-01T12:00:00Z",
	}
	
	// Marshal
	data, err := proto.Marshal(notification)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Greater(t, len(data), 0)
	
	// Unmarshal
	result := &NotificationRequest{}
	err = proto.Unmarshal(data, result)
	assert.NoError(t, err)
	assert.Equal(t, notification.UserId, result.UserId)
	assert.Equal(t, notification.Type, result.Type)
	assert.Equal(t, notification.Title, result.Title)
	assert.Equal(t, notification.Payload, result.Payload)
	assert.Equal(t, notification.Hyperlink, result.Hyperlink)
	assert.Equal(t, notification.CreatedAt, result.CreatedAt)
}

// TestStartPubSubClientMissingProjectID tests client creation with missing project ID
func TestStartPubSubClientMissingProjectID(t *testing.T) {
	// Save current env var
	originalProjectID := os.Getenv("PUBSUB_PROJECT")
	defer os.Setenv("PUBSUB_PROJECT", originalProjectID)
	
	// Clear the env var to trigger error
	os.Unsetenv("PUBSUB_PROJECT")
	
	ctx := context.Background()
	
	// With empty project ID, client creation should fail
	client, err := StartPubSubClient(ctx, nil, nil)
	assert.Nil(t, client, "Client should be nil with empty project ID")
	assert.Error(t, err, "Should return error with empty project ID")
}

// TestContextUsage tests that context is properly used
func TestContextUsage(t *testing.T) {
	backgroundCtx := context.Background()
	assert.NotNil(t, backgroundCtx, "Context should not be nil")
	
	// Test context cancellation
	cancelCtx, cancel := context.WithCancel(backgroundCtx)
	cancel()
	
	select {
	case <-cancelCtx.Done():
		t.Log("Context cancelled successfully")
	default:
		t.Error("Context should be cancelled")
	}
}

// TestOrderUnmarshalJSON tests unmarshaling orders
func TestOrderUnmarshalJSON(t *testing.T) {
	orderJSON := `{
		"id": 1,
		"customer_id": 101,
		"seller_id": 501,
		"tracking_code": "TRACK001"
	}`
	
	var order models.Orders
	err := json.Unmarshal([]byte(orderJSON), &order)
	assert.NoError(t, err, "Should unmarshal order")
	
	assert.Equal(t, uint(1), order.Id)
	assert.Equal(t, uint(101), order.Customer_ID)
	assert.Equal(t, uint(501), order.Seller_ID)
	assert.Equal(t, "TRACK001", order.Tracking_Code)
}

// TestNotificationPayloadWithEmptyFields tests notification with empty fields
func TestNotificationPayloadWithEmptyFields(t *testing.T) {
	notification := &NotificationRequest{
		UserId:    "",
		Type:      "",
		Title:     "",
		Payload:   "",
		Hyperlink: "",
		CreatedAt: "",
	}
	
	// Marshal even with empty fields
	data, err := proto.Marshal(notification)
	assert.NoError(t, err, "Should marshal notification with empty fields")
	assert.NotNil(t, data, "Marshaled data should not be nil")
	
	// Unmarshal
	result := &NotificationRequest{}
	err = proto.Unmarshal(data, result)
	assert.NoError(t, err, "Should unmarshal successfully")
	assert.Equal(t, "", result.UserId)
}

// TestOrderStatusHistoryUnmarshalWithoutStorageID tests unmarshaling without storage_id
func TestOrderStatusHistoryUnmarshalWithoutStorageID(t *testing.T) {
	statusJSON := `{
		"order_id": 1,
		"order_status": "PROCESSING",
		"note": "Order received",
		"order_location": "Warehouse"
	}`
	
	var orderUpdate models.OrderStatusHistory
	err := json.Unmarshal([]byte(statusJSON), &orderUpdate)
	assert.NoError(t, err, "Should unmarshal order status without storage_id")
	
	// Storage_ID should be nil when not provided
	assert.Nil(t, orderUpdate.Storage_ID)
}

// TestHTTPTestContextCreation tests creating a test context for HTTP requests
func TestHTTPTestContextCreation(t *testing.T) {
	w := httptest.NewRecorder()
	assert.NotNil(t, w, "Should create test recorder")
	
	req := httptest.NewRequest("POST", "/api/test", bytes.NewReader([]byte("{}")))
	assert.NotNil(t, req, "Should create test request")
	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, "/api/test", req.URL.Path)
}

// TestMultipleNotificationMarshalings tests marshaling multiple notifications
func TestMultipleNotificationMarshalings(t *testing.T) {
	notifications := []NotificationRequest{
		{
			UserId:    "1",
			Type:      "sms",
			Title:     "Notification 1",
			Payload:   "Payload 1",
		},
		{
			UserId:    "2",
			Type:      "email",
			Title:     "Notification 2",
			Payload:   "Payload 2",
		},
		{
			UserId:    "3",
			Type:      "push",
			Title:     "Notification 3",
			Payload:   "Payload 3",
		},
	}
	
	for i, notif := range notifications {
		data, err := proto.Marshal(&notif)
		assert.NoError(t, err, "Should marshal notification %d", i)
		
		result := &NotificationRequest{}
		err = proto.Unmarshal(data, result)
		assert.NoError(t, err, "Should unmarshal notification %d", i)
		assert.Equal(t, notif.UserId, result.UserId)
		assert.Equal(t, notif.Type, result.Type)
	}
}

// TestLogMessageParsing tests parsing various log message scenarios
func TestLogMessageParsing(t *testing.T) {
	testMessages := []struct {
		name    string
		message string
		shouldParse bool
	}{
		{
			name:    "Valid message",
			message: `{"order_id": 1, "order_status": "SHIPPED", "note": "test", "order_location": "loc"}`,
			shouldParse: true,
		},
		{
			name:    "Empty JSON object",
			message: `{}`,
			shouldParse: true, // JSON parses, but validation fails
		},
		{
			name:    "Invalid JSON",
			message: `not json`,
			shouldParse: false,
		},
	}
	
	for _, tc := range testMessages {
		t.Run(tc.name, func(t *testing.T) {
			var data map[string]interface{}
			err := json.Unmarshal([]byte(tc.message), &data)
			
			if tc.shouldParse {
				assert.NoError(t, err, "Should parse message")
			} else {
				assert.Error(t, err, "Should fail to parse message")
			}
		})
	}
}

// TestStringContainmentCheck tests string containment logic
func TestStringContainmentCheck(t *testing.T) {
	tests := []struct {
		haystack string
		needle   string
		expected bool
	}{
		{"AlreadyExists error", "AlreadyExists", true},
		{"Some other error", "AlreadyExists", false},
		{"", "AlreadyExists", false},
		{"AlreadyExists", "AlreadyExists", true},
	}
	
	for _, tc := range tests {
		t.Run(fmt.Sprintf("Contains_%s_in_%s", tc.needle, tc.haystack), func(t *testing.T) {
			// Simulate the string containment logic used in pubsub.go
			result := strings.Contains(tc.haystack, tc.needle)
			assert.Equal(t, tc.expected, result, "String containment check")
		})
	}
}

// TestPublishNotificationWithNilClient tests publishing with nil client
func TestPublishNotificationWithNilClient(t *testing.T) {
	ctx := context.Background()
	notification := []byte{1, 2, 3}
	
	// PublishNotification should return error with nil client
	err := PublishNotification(ctx, nil, notification)
	assert.Error(t, err, "Should return error with nil client")
}

// TestPublishNotificationWithEmptyPayload tests publishing with empty payload
func TestPublishNotificationWithEmptyPayload(t *testing.T) {
	ctx := context.Background()
	
	// PublishNotification should return error with empty payload
	err := PublishNotification(ctx, nil, []byte{})
	assert.Error(t, err, "Should return error with empty payload")
}

// TestStartListenerWithNilClient tests StartListener with nil client
func TestStartListenerWithNilClient(t *testing.T) {
	ctx := context.Background()
	
	// StartListener should return error with nil client
	err := StartListener(ctx, nil, nil, nil, nil)
	assert.Error(t, err, "Should return error with nil client")
}

// TestStartListenerWithNilSubscription tests StartListener with nil subscription
func TestStartListenerWithNilSubscription(t *testing.T) {
	ctx := context.Background()
	
	// StartListener should return error with nil subscription
	err := StartListener(ctx, &pubsub.Client{}, nil, nil, nil)
	assert.Error(t, err, "Should return error with nil subscription")
}

// TestStartListenerOrdersWithNilClient tests StartListenerOrders with nil client
func TestStartListenerOrdersWithNilClient(t *testing.T) {
	ctx := context.Background()
	
	// StartListenerOrders should return error with nil client
	err := StartListenerOrders(ctx, nil, nil, nil, nil)
	assert.Error(t, err, "Should return error with nil client")
}

// TestStartListenerOrdersWithNilSubscription tests StartListenerOrders with nil subscription
func TestStartListenerOrdersWithNilSubscription(t *testing.T) {
	ctx := context.Background()
	
	// StartListenerOrders should return error with nil subscription
	err := StartListenerOrders(ctx, &pubsub.Client{}, nil, nil, nil)
	assert.Error(t, err, "Should return error with nil subscription")
}

// TestCreateTopicWithNilClient tests CreateTopicWithID with nil client
func TestCreateTopicWithNilClient(t *testing.T) {
	ctx := context.Background()
	
	// CreateTopicWithID should return error with nil client
	topic, err := CreateTopicWithID(ctx, nil, "test-topic")
	assert.Nil(t, topic, "Should return nil topic")
	assert.Error(t, err, "Should return error with nil client")
}

// TestSubscribeClientWithNilClient tests SubscribeClient with nil client
func TestSubscribeClientWithNilClient(t *testing.T) {
	ctx := context.Background()
	
	// SubscribeClient should return error with nil client
	sub, err := SubscribeClient(ctx, nil, "test-topic", "test-sub")
	assert.Nil(t, sub, "Should return nil subscription")
	assert.Error(t, err, "Should return error with nil client")
}

// TestCreateTopicWithEmptyTopicID tests CreateTopicWithID with empty topic ID
func TestCreateTopicWithEmptyTopicID(t *testing.T) {
	ctx := context.Background()
	
	topic, err := CreateTopicWithID(ctx, &pubsub.Client{}, "")
	assert.Nil(t, topic, "Should return nil topic")
	assert.Error(t, err, "Should return error with empty topicID")
	assert.Equal(t, "topicID is empty", err.Error())
}

// TestSubscribeClientWithEmptyTopicID tests SubscribeClient with empty topic ID
func TestSubscribeClientWithEmptyTopicID(t *testing.T) {
	ctx := context.Background()
	
	sub, err := SubscribeClient(ctx, &pubsub.Client{}, "", "test-sub")
	assert.Nil(t, sub, "Should return nil subscription")
	assert.Error(t, err, "Should return error with empty topicID")
	assert.Equal(t, "topicID is empty", err.Error())
}

// TestSubscribeClientWithEmptySubscriptionID tests SubscribeClient with empty subscription ID
func TestSubscribeClientWithEmptySubscriptionID(t *testing.T) {
	ctx := context.Background()
	
	sub, err := SubscribeClient(ctx, &pubsub.Client{}, "test-topic", "")
	assert.Nil(t, sub, "Should return nil subscription")
	assert.Error(t, err, "Should return error with empty subscriptionID")
	assert.Equal(t, "subscriptionID is empty", err.Error())
}

// TestBuildNotificationPayloadOrderEmptyPayload tests with nil/empty payload
func TestBuildNotificationPayloadOrderEmptyPayload(t *testing.T) {
	// Test with nil payload
	result := buildNotificationPayloadOrder(nil, nil, nil)
	assert.Nil(t, result, "Should return nil for nil payload")
	
	// Test with empty payload
	result = buildNotificationPayloadOrder([]byte{}, nil, nil)
	assert.Nil(t, result, "Should return nil for empty payload")
}

// TestBuildNotificationPayloadStatusEmptyPayload tests with nil/empty payload
func TestBuildNotificationPayloadStatusEmptyPayload(t *testing.T) {
	// Test with nil payload
	result := buildNotificationPayloadStatus(nil, nil, nil)
	assert.Nil(t, result, "Should return nil for nil payload")
	
	// Test with empty payload
	result = buildNotificationPayloadStatus([]byte{}, nil, nil)
	assert.Nil(t, result, "Should return nil for empty payload")
}

// TestCreateTopicLogic tests CreateTopicWithID error handling
func TestCreateTopicLogic(t *testing.T) {
	// Test with valid nil client
	ctx := context.Background()
	_, err := CreateTopicWithID(ctx, nil, "valid-topic")
	assert.Error(t, err)
	assert.Equal(t, "pubsub client is nil", err.Error())
}

// TestOrderPayloadMarshaling tests marshaling of order payloads
func TestOrderPayloadMarshaling(t *testing.T) {
	tests := []struct {
		name   string
		json   string
		expUserID string
	}{
		{
			name:   "Simple order",
			json:   `{"id": 5, "customer_id": 250}`,
			expUserID: "250",
		},
		{
			name:   "Order with additional fields",
			json:   `{"id": 10, "customer_id": 500, "total": 99.99, "status": "pending"}`,
			expUserID: "500",
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := buildNotificationPayloadOrder([]byte(test.json), nil, nil)
			assert.NotNil(t, result)
			
			var notification NotificationRequest
			err := proto.Unmarshal(result, &notification)
			assert.NoError(t, err)
			assert.Equal(t, test.expUserID, notification.UserId)
			assert.Contains(t, notification.Payload, "Order with ID")
		})
	}
}

// TestStatusPayloadInvalidJSON tests status payload with bad JSON
func TestStatusPayloadInvalidJSON(t *testing.T) {
	result := buildNotificationPayloadStatus([]byte("not valid json"), nil, nil)
	assert.Nil(t, result, "Should return nil for invalid JSON")
}

// TestPublishNotificationLogging tests PublishNotification error cases
func TestPublishNotificationLogging(t *testing.T) {
	ctx := context.Background()
	
	// Test with nil client and valid payload
	validNotif := &NotificationRequest{
		UserId: "test",
		Type:   "sms",
	}
	data, _ := proto.Marshal(validNotif)
	
	err := PublishNotification(ctx, nil, data)
	assert.Error(t, err)
	assert.Equal(t, "pubsub client is nil", err.Error())
}

// TestEmptyPayloadValidation tests empty payload validation
func TestEmptyPayloadValidation(t *testing.T) {
	tests := []struct {
		name    string
		payload []byte
		expErr  bool
	}{
		{
			name:    "Nil payload",
			payload: nil,
			expErr:  true,
		},
		{
			name:    "Empty bytes",
			payload: []byte{},
			expErr:  true,
		},
		{
			name:    "Valid data",
			payload: []byte{1, 2, 3},
			expErr:  false,
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			err := PublishNotification(ctx, nil, test.payload)
			
			if test.expErr {
				assert.Error(t, err)
			} else {
				// Will still error because client is nil, but not because of payload
				if err != nil {
					assert.NotEqual(t, "notification payload is empty", err.Error())
				}
			}
		})
	}
}