// handlers/order_handler_test.go
package handlers

import (
	"app/requestModels"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// --- Helpers ---

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm with sqlmock: %v", err)
	}
	return gdb, mock
}

func performRequest(r *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// --- Tests ---

func TestGetOrderByID_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	h := &OrderHandler{DB: db}
	r := gin.Default()
	r.GET("/order/:id", h.GetOrderByID)

	mock.ExpectQuery(`SELECT \* FROM "orders" WHERE "orders"."id" = \$1 ORDER BY "orders"."id" LIMIT \$2`).
		WithArgs("999", 1). // two args: id and limit
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	// Expect preload query on products (empty result)
	mock.ExpectQuery(`SELECT \* FROM "products" WHERE "products"."order_id" IN \(\$1\)`).
		WithArgs("999").
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id"}))

	req := httptest.NewRequest(http.MethodGet, "/order/999", nil)
	w := performRequest(r, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetOrderByID_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	h := &OrderHandler{DB: db}
	r := gin.Default()
	r.GET("/order/:id", h.GetOrderByID)

	// Mock row
	mock.ExpectQuery(`SELECT \* FROM "orders" WHERE "orders"."id" = \$1 ORDER BY "orders"."id" LIMIT \$2`).
		WithArgs("1", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "delivery_address"}).
			AddRow(1, "Rua Nova"))

		// Expect preload query on products (empty result)
	mock.ExpectQuery(`SELECT \* FROM "order_products" WHERE "order_products"."order_id" = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id"}))

	req := httptest.NewRequest(http.MethodGet, "/order/1", nil)
	w := performRequest(r, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAddOrder_BadInput(t *testing.T) {
	db, _ := setupMockDB(t)
	h := &OrderHandler{DB: db}
	r := gin.Default()
	r.POST("/order/add", h.AddOrder)

	// Send malformed JSON to trigger 400
	req := httptest.NewRequest(http.MethodPost, "/order/add", bytes.NewBufferString(`not-a-json`))
	req.Header.Set("Content-Type", "application/json")
	w := performRequest(r, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpdateOrder_BadInput(t *testing.T) {
	db, _ := setupMockDB(t)
	h := &OrderHandler{DB: db}
	r := gin.Default()
	r.POST("/order/update", h.UpdateOrder)

	req := httptest.NewRequest(http.MethodPost, "/order/update", bytes.NewBufferString(`not-a-json`))
	req.Header.Set("Content-Type", "application/json")
	w := performRequest(r, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpdateOrder_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	h := &OrderHandler{DB: db}
	r := gin.Default()
	r.POST("/order/update", h.UpdateOrder)

	// Expect SELECT returning no rows
	mock.ExpectQuery(`SELECT \* FROM "orders" WHERE "orders"."id" = \$1 ORDER BY "orders"."id" LIMIT \$2`).
		WithArgs(999, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	payload := requestModels.UpdateOrderRequest{
		OrderID:           999,
		DeliveryAddress:   "New Address",
		DeliveryLatitude:  41.15,
		DeliveryLongitude: -8.61,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/order/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := performRequest(r, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateOrder_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	h := &OrderHandler{DB: db}
	r := gin.Default()
	r.POST("/order/update", h.UpdateOrder)

	// Expect SELECT for order
	mock.ExpectQuery(`SELECT \* FROM "orders" WHERE "orders"."id" = \$1 ORDER BY "orders"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "delivery_address", "delivery_latitude", "delivery_longitude", "seller_latitude", "seller_longitude"}).
			AddRow(1, "Old Address", 41.1, -8.6, 41.2, -8.5))

	// Expect SELECT for order_status_history
	mock.ExpectQuery(`SELECT \* FROM "order_status_history" WHERE order_id = \$1 ORDER BY timestamp_history desc,"order_status_history"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "order_status", "timestamp_history"}).
			AddRow(1, 1, "PROCESSING", time.Now()))

		// Expect transaction begin
	mock.ExpectBegin()

	// Expect UPDATE statement
	mock.ExpectExec(`UPDATE "orders"`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect commit
	mock.ExpectCommit()

	payload := requestModels.UpdateOrderRequest{
		OrderID:           1,
		DeliveryAddress:   "Updated Address",
		DeliveryLatitude:  41.1496,
		DeliveryLongitude: -8.6109,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/order/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := performRequest(r, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCancelOrder_BadInput(t *testing.T) {
	db, _ := setupMockDB(t)
	h := &OrderHandler{DB: db}
	r := gin.Default()
	r.POST("/order/cancel", h.CancelOrder)

	req := httptest.NewRequest(http.MethodPost, "/order/cancel", bytes.NewBufferString(`not-a-json`))
	req.Header.Set("Content-Type", "application/json")
	w := performRequest(r, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCancelOrder_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	h := &OrderHandler{DB: db}
	r := gin.Default()
	r.POST("/order/cancel", h.CancelOrder)

	// Expect SELECT returning no rows
	mock.ExpectQuery(`SELECT \* FROM "orders" WHERE "orders"."id" = \$1 ORDER BY "orders"."id" LIMIT \$2`).
		WithArgs(999, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	payload := requestModels.CancelOrderRequest{
		OrderID: 999,
		Reason:  "Customer request",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/order/cancel", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := performRequest(r, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestCancelOrder_NotProcessing(t *testing.T) {
	db, mock := setupMockDB(t)
	h := &OrderHandler{DB: db}
	r := gin.Default()
	r.POST("/order/cancel", h.CancelOrder)

	// Expect SELECT for order
	mock.ExpectQuery(`SELECT \* FROM "orders" WHERE "orders"."id" = \$1 ORDER BY "orders"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "tracking_code", "delivery_address", "delivery_latitude", "delivery_longitude", "seller_latitude", "seller_longitude", "price"}).
			AddRow(1, "TRACK123", "Rua Nova", 41.1, -8.6, 41.2, -8.5, 100.00))

	// Expect SELECT for order_status_history (SHIPPED status)
	mock.ExpectQuery(`SELECT \* FROM "order_status_history" WHERE order_id = \$1 ORDER BY timestamp_history desc,"order_status_history"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "order_status", "timestamp_history"}).
			AddRow(1, 1, "SHIPPED", time.Now()))

	payload := requestModels.CancelOrderRequest{
		OrderID: 1,
		Reason:  "Customer request",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/order/cancel", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := performRequest(r, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["error"] != "Cannot cancel order with status: SHIPPED" {
		t.Fatalf("expected error message about SHIPPED status, got: %s", response["error"])
	}
}

func TestCancelOrder_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	h := &OrderHandler{DB: db}
	r := gin.Default()
	r.POST("/order/cancel", h.CancelOrder)

	// Expect SELECT for order
	mock.ExpectQuery(`SELECT \* FROM "orders" WHERE "orders"."id" = \$1 ORDER BY "orders"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "tracking_code", "delivery_address", "delivery_latitude", "delivery_longitude", "seller_latitude", "seller_longitude", "price"}).
			AddRow(1, "TRACK123", "Rua Nova", 41.1, -8.6, 41.2, -8.5, 100.00))

	// Expect SELECT for order_status_history (PROCESSING status)
	mock.ExpectQuery(`SELECT \* FROM "order_status_history" WHERE order_id = \$1 ORDER BY timestamp_history desc,"order_status_history"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "order_status", "timestamp_history"}).
			AddRow(1, 1, "PROCESSING", time.Now()))

	// Expect transaction begin
	mock.ExpectBegin()

	// Expect INSERT for order_status_history with RETURNING clause
	mock.ExpectQuery(`INSERT INTO "order_status_history"`).
		WithArgs(
			sqlmock.AnyArg(),                  // order_id
			sqlmock.AnyArg(),                  // timestamp_history
			"CANCELLED",                       // order_status
			"Customer requested cancellation", // note
			"",                                // blockchain_transaction
			"SYSTEM",                          // order_location
		).
		WillReturnRows(sqlmock.NewRows([]string{"storage_id", "id"}).
			AddRow(nil, 2))

	// Expect commit
	mock.ExpectCommit()

	payload := requestModels.CancelOrderRequest{
		OrderID: 1,
		Reason:  "Customer requested cancellation",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/order/cancel", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := performRequest(r, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["message"] != "Order cancelled successfully" {
		t.Fatalf("expected success message, got: %v", response["message"])
	}

	if response["order_id"] != float64(1) {
		t.Fatalf("expected order_id 1, got: %v", response["order_id"])
	}

	if response["status"] != "CANCELLED" {
		t.Fatalf("expected status CANCELLED, got: %v", response["status"])
	}
}
