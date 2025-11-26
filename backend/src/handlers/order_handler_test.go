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
        WithArgs(999,1).
        WillReturnRows(sqlmock.NewRows([]string{"id"}))

    payload := requestModels.UpdateOrderRequest{
        OrderID:          999,
        DeliveryAddress:  "New Address",
        DeliveryLatitude: 41.15,
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

