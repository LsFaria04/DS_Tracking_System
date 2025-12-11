package handlers

import (
	"app/models"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// --- GetOrderStatusByOrderID Tests ---

func TestGetOrderStatusByOrderID_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	h := &OrderStatusHistoryHandler{DB: db}
	r := gin.Default()
	r.GET("/order/:order_id/status", h.GetOrderStatusByOrderID)

	ts := time.Now()
	mock.ExpectQuery(`SELECT \* FROM "order_status_history" WHERE Order_ID = \$1 ORDER BY Timestamp_History desc`).
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "timestamp_history", "order_status", "note", "order_location"}).
			AddRow(2, 1, ts, "SHIPPED", "Shipped to customer", "Warehouse A").
			AddRow(1, 1, ts.Add(-time.Hour), "PROCESSING", "Order received", "Origin"))

	// Expect preload for Storage
	mock.ExpectQuery(`SELECT \* FROM "storage" WHERE "storage"."id" IN \(\$1,\$2\)`).
		WithArgs(nil, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	req := httptest.NewRequest(http.MethodGet, "/order/1/status", nil)
	w := performRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "order_status_history")
}

func TestGetOrderStatusByOrderID_EmptyHistory(t *testing.T) {
	db, mock := setupMockDB(t)
	h := &OrderStatusHistoryHandler{DB: db}
	r := gin.Default()
	r.GET("/order/:order_id/status", h.GetOrderStatusByOrderID)

	mock.ExpectQuery(`SELECT \* FROM "order_status_history" WHERE Order_ID = \$1 ORDER BY Timestamp_History desc`).
		WithArgs("999").
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "timestamp_history", "order_status", "note", "order_location"}))

	// Empty preload
	mock.ExpectQuery(`SELECT \* FROM "storage"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	req := httptest.NewRequest(http.MethodGet, "/order/999/status", nil)
	w := performRequest(r, req)

	// Should return 200 with empty array
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "order_status_history")
}

func TestGetOrderStatusByOrderID_DBError(t *testing.T) {
	db, mock := setupMockDB(t)
	h := &OrderStatusHistoryHandler{DB: db}
	r := gin.Default()
	r.GET("/order/:order_id/status", h.GetOrderStatusByOrderID)

	mock.ExpectQuery(`SELECT \* FROM "order_status_history" WHERE Order_ID = \$1`).
		WithArgs("1").
		WillReturnError(errors.New("db failure"))

	req := httptest.NewRequest(http.MethodGet, "/order/1/status", nil)
	w := performRequest(r, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// --- AddOrderUpdate Tests ---

func TestAddOrderUpdate_BadInput(t *testing.T) {
	db, _ := setupMockDB(t)
	h := &OrderStatusHistoryHandler{DB: db}
	r := gin.Default()
	r.POST("/order/update", h.AddOrderUpdate)

	req := httptest.NewRequest(http.MethodPost, "/order/update", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	w := performRequest(r, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddOrderUpdate_Success_NoBlockchain(t *testing.T) {
	db, mock := setupMockDB(t)
	h := &OrderStatusHistoryHandler{DB: db, Client: nil} // No blockchain
	r := gin.Default()
	r.POST("/order/update", h.AddOrderUpdate)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "order_status_history"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	payload := models.OrderStatusHistory{
		Order_ID:       1,
		Order_Status:   "SHIPPED",
		Note:           "Package shipped",
		Order_Location: "Warehouse B",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/order/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := performRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Update stored successfully")
}

func TestAddOrderUpdate_DBCreateError(t *testing.T) {
	db, mock := setupMockDB(t)
	h := &OrderStatusHistoryHandler{DB: db, Client: nil}
	r := gin.Default()
	r.POST("/order/update", h.AddOrderUpdate)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "order_status_history"`).
		WillReturnError(errors.New("db insert failed"))
	mock.ExpectRollback()

	payload := models.OrderStatusHistory{
		Order_ID:       1,
		Order_Status:   "DELIVERED",
		Note:           "Package delivered",
		Order_Location: "Customer Address",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/order/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := performRequest(r, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAddOrderUpdate_TimestampDefault(t *testing.T) {
	db, mock := setupMockDB(t)
	h := &OrderStatusHistoryHandler{DB: db, Client: nil}
	r := gin.Default()
	r.POST("/order/update", h.AddOrderUpdate)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "order_status_history"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// Payload without timestamp - should use current time
	payload := models.OrderStatusHistory{
		Order_ID:       1,
		Order_Status:   "IN_TRANSIT",
		Note:           "Package in transit",
		Order_Location: "Hub C",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/order/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := performRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// --- Tracking State Tests ---

func TestAddOrderUpdate_AllTrackingStates(t *testing.T) {
	trackingStates := []string{"PROCESSING", "SHIPPED", "IN_TRANSIT", "DELIVERED"}

	for _, state := range trackingStates {
		t.Run("State_"+state, func(t *testing.T) {
			db, mock := setupMockDB(t)
			h := &OrderStatusHistoryHandler{DB: db, Client: nil}
			r := gin.Default()
			r.POST("/order/update", h.AddOrderUpdate)

			mock.ExpectBegin()
			mock.ExpectQuery(`INSERT INTO "order_status_history"`).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			mock.ExpectCommit()

			payload := models.OrderStatusHistory{
				Order_ID:       1,
				Order_Status:   state,
				Note:           "Status update to " + state,
				Order_Location: "Location",
			}
			body, _ := json.Marshal(payload)

			req := httptest.NewRequest(http.MethodPost, "/order/update", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := performRequest(r, req)

			assert.Equal(t, http.StatusOK, w.Code, "Failed for state: %s", state)
			assert.Contains(t, w.Body.String(), "Update stored successfully")
		})
	}
}

func TestAddOrderUpdate_WithStorageID(t *testing.T) {
	db, mock := setupMockDB(t)
	h := &OrderStatusHistoryHandler{DB: db, Client: nil}
	r := gin.Default()
	r.POST("/order/update", h.AddOrderUpdate)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "order_status_history"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	storageID := uint(5)
	payload := models.OrderStatusHistory{
		Order_ID:       1,
		Order_Status:   "SHIPPED",
		Note:           "Stored in warehouse",
		Order_Location: "Warehouse D",
		Storage_ID:     &storageID,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/order/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := performRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// --- Edge Case Tests ---

func TestGetOrderStatusByOrderID_RecordNotFoundError(t *testing.T) {
	db, mock := setupMockDB(t)
	h := &OrderStatusHistoryHandler{DB: db}
	r := gin.Default()
	r.GET("/order/:order_id/status", h.GetOrderStatusByOrderID)

	mock.ExpectQuery(`SELECT \* FROM "order_status_history" WHERE Order_ID = \$1`).
		WithArgs("1").
		WillReturnError(gorm.ErrRecordNotFound)

	req := httptest.NewRequest(http.MethodGet, "/order/1/status", nil)
	w := performRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
