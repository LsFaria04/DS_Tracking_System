package handlers

import (
	"app/models"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestGetOrderProducts_MissingParam(t *testing.T) {
    db, _ := setupMockDB(t)
    h := &OrderProductHandler{DB: db}
    r := gin.Default()
    r.GET("/order_products", h.GetOrderProducts)

    req := httptest.NewRequest(http.MethodGet, "/order_products", nil)
    w := performRequest(r, req)

    if w.Code != http.StatusBadRequest {
        t.Fatalf("expected 400, got %d", w.Code)
    }
}

func TestGetOrderProducts_Success(t *testing.T) {
    db, mock := setupMockDB(t)
    h := &OrderProductHandler{DB: db}
    r := gin.Default()
    r.GET("/order_products", h.GetOrderProducts)

    mock.ExpectQuery(`SELECT \* FROM "order_products" WHERE order_id = \$1`).
        WithArgs("1").
        WillReturnRows(sqlmock.NewRows([]string{"id","order_id","product_id","quantity"}).
            AddRow(1,1,10,2))

    req := httptest.NewRequest(http.MethodGet, "/order_products?order_id=1", nil)
    w := performRequest(r, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", w.Code)
    }
}

func TestGetOrderProducts_DBError(t *testing.T) {
    db, mock := setupMockDB(t)
    h := &OrderProductHandler{DB: db}
    r := gin.Default()
    r.GET("/order_products", h.GetOrderProducts)

    mock.ExpectQuery(`SELECT \* FROM "order_products" WHERE order_id = \$1`).
        WithArgs("1").
        WillReturnError(errors.New("db fail"))

    req := httptest.NewRequest(http.MethodGet, "/order_products?order_id=1", nil)
    w := performRequest(r, req)

    if w.Code != http.StatusInternalServerError {
        t.Fatalf("expected 500, got %d", w.Code)
    }
}


func TestGetOrderProductByID_NotFound(t *testing.T) {
    db, mock := setupMockDB(t)
    h := &OrderProductHandler{DB: db}
    r := gin.Default()
    r.GET("/order_product/:id", h.GetOrderProductByID)

    mock.ExpectQuery(`SELECT \* FROM "order_products" WHERE "order_products"."id" = \$1 ORDER BY "order_products"."id" LIMIT \$2`).
        WithArgs("999",1).
        WillReturnError(gorm.ErrRecordNotFound)

    req := httptest.NewRequest(http.MethodGet, "/order_product/999", nil)
    w := performRequest(r, req)

    if w.Code != http.StatusNotFound {
        t.Fatalf("expected 404, got %d", w.Code)
    }
}

func TestGetOrderProductByID_Success(t *testing.T) {
    db, mock := setupMockDB(t)
    h := &OrderProductHandler{DB: db}
    r := gin.Default()
    r.GET("/order_product/:id", h.GetOrderProductByID)

    mock.ExpectQuery(`SELECT \* FROM "order_products" WHERE "order_products"."id" = \$1 ORDER BY "order_products"."id" LIMIT \$2`).
        WithArgs("1",1).
        WillReturnRows(sqlmock.NewRows([]string{"id","order_id","product_id","quantity"}).
            AddRow(1,1,10,2))

    req := httptest.NewRequest(http.MethodGet, "/order_product/1", nil)
    w := performRequest(r, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", w.Code)
    }
}

func TestGetOrderProductByID_DBError(t *testing.T) {
    db, mock := setupMockDB(t)
    h := &OrderProductHandler{DB: db}
    r := gin.Default()
    r.GET("/order_product/:id", h.GetOrderProductByID)

    mock.ExpectQuery(`SELECT \* FROM "order_products" WHERE "order_products"."id" = \$1 ORDER BY "order_products"."id" LIMIT \$2`).
        WithArgs("1",1).
        WillReturnError(errors.New("db fail"))

    req := httptest.NewRequest(http.MethodGet, "/order_product/1", nil)
    w := performRequest(r, req)

    if w.Code != http.StatusInternalServerError {
        t.Fatalf("expected 500, got %d", w.Code)
    }
}


func TestAddOrderProduct_BadJSON(t *testing.T) {
    db, _ := setupMockDB(t)
    h := &OrderProductHandler{DB: db}
    r := gin.Default()
    r.POST("/order_product/add", h.AddOrderProduct)

    req := httptest.NewRequest(http.MethodPost, "/order_product/add", bytes.NewBufferString("not-json"))
    req.Header.Set("Content-Type","application/json")
    w := performRequest(r, req)

    if w.Code != http.StatusBadRequest {
        t.Fatalf("expected 400, got %d", w.Code)
    }
}

func TestAddOrderProduct_InvalidQuantity(t *testing.T) {
    db, _ := setupMockDB(t)
    h := &OrderProductHandler{DB: db}
    r := gin.Default()
    r.POST("/order_product/add", h.AddOrderProduct)

    payload := models.OrderProduct{Order_ID: 1, Product_ID: 10, Quantity: 0}
    body,_ := json.Marshal(payload)

    req := httptest.NewRequest(http.MethodPost,"/order_product/add",bytes.NewBuffer(body))
    req.Header.Set("Content-Type","application/json")
    w := performRequest(r,req)

    if w.Code != http.StatusBadRequest {
        t.Fatalf("expected 400, got %d", w.Code)
    }
}

func TestAddOrderProduct_OrderNotFound(t *testing.T) {
    db, mock := setupMockDB(t)
    h := &OrderProductHandler{DB: db}
    r := gin.Default()
    r.POST("/order_product/add", h.AddOrderProduct)

    payload := models.OrderProduct{Order_ID: 999, Product_ID: 10, Quantity: 2}
    body,_ := json.Marshal(payload)

    mock.ExpectQuery(`SELECT \* FROM "orders" WHERE "orders"."id" = \$1 ORDER BY "orders"."id" LIMIT \$2`).
        WithArgs(999,1).
        WillReturnError(gorm.ErrRecordNotFound)

    req := httptest.NewRequest(http.MethodPost,"/order_product/add",bytes.NewBuffer(body))
    req.Header.Set("Content-Type","application/json")
    w := performRequest(r,req)

    if w.Code != http.StatusNotFound {
        t.Fatalf("expected 404, got %d", w.Code)
    }
}

func TestAddOrderProduct_DBErrorOnOrderLookup(t *testing.T) {
    db, mock := setupMockDB(t)
    h := &OrderProductHandler{DB: db}
    r := gin.Default()
    r.POST("/order_product/add", h.AddOrderProduct)

    payload := models.OrderProduct{Order_ID: 1, Product_ID: 10, Quantity: 2}
    body,_ := json.Marshal(payload)

    mock.ExpectQuery(`SELECT \* FROM "orders" WHERE "orders"."id" = \$1 ORDER BY "orders"."id" LIMIT \$2`).
        WithArgs(1,1).
        WillReturnError(errors.New("db fail"))

    req := httptest.NewRequest(http.MethodPost,"/order_product/add",bytes.NewBuffer(body))
    req.Header.Set("Content-Type","application/json")
    w := performRequest(r,req)

    if w.Code != http.StatusInternalServerError {
        t.Fatalf("expected 500, got %d", w.Code)
    }
}

func TestAddOrderProduct_DBErrorOnCreate(t *testing.T) {
    db, mock := setupMockDB(t)
    h := &OrderProductHandler{DB: db}
    r := gin.Default()
    r.POST("/order_product/add", h.AddOrderProduct)

    payload := models.OrderProduct{Order_ID: 1, Product_ID: 10, Quantity: 2}
    body,_ := json.Marshal(payload)

    // Order exists
    mock.ExpectQuery(`SELECT \* FROM "orders" WHERE "orders"."id" = \$1 ORDER BY "orders"."id" LIMIT \$2`).
        WithArgs(1,1).
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

    // Fail on create
    mock.ExpectBegin()
    mock.ExpectExec(`INSERT INTO "order_products"`).WillReturnError(errors.New("insert fail"))
    mock.ExpectRollback()

    req := httptest.NewRequest(http.MethodPost,"/order_product/add",bytes.NewBuffer(body))
    req.Header.Set("Content-Type","application/json")
    w := performRequest(r,req)

    if w.Code != http.StatusInternalServerError {
        t.Fatalf("expected 500, got %d", w.Code)
    }
}

func TestAddOrderProduct_Success(t *testing.T) {
    db, mock := setupMockDB(t)
    h := &OrderProductHandler{DB: db}
    r := gin.Default()
    r.POST("/order_product/add", h.AddOrderProduct)

    payload := models.OrderProduct{Order_ID: 1, Product_ID: 10, Quantity: 2}
    body,_ := json.Marshal(payload)

    // Order exists
    mock.ExpectQuery(`SELECT \* FROM "orders" WHERE "orders"."id" = \$1 ORDER BY "orders"."id" LIMIT \$2`).
        WithArgs(1,1).
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

    // Create succeeds
    mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "order_products" ("order_id","product_id","quantity","product_name_at_purchase","product_price_at_purchase") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(1, 10, 2, "", float64(0)). 
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()



    // Reload product
    mock.ExpectQuery(`SELECT \* FROM "order_products" WHERE "order_products"."id" = \$1 ORDER BY "order_products"."id" LIMIT \$2`).
        WithArgs(1,1).
        WillReturnRows(sqlmock.NewRows([]string{"id","order_id","product_id","quantity"}).
            AddRow(1,1,10,2))

    req := httptest.NewRequest(http.MethodPost,"/order_product/add",bytes.NewBuffer(body))
    req.Header.Set("Content-Type","application/json")
    w := performRequest(r,req)

    if w.Code != http.StatusCreated {
        t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
    }
}



func TestUpdateOrderProduct_NotFound(t *testing.T) {
    db, mock := setupMockDB(t)
    h := &OrderProductHandler{DB: db}
    r := gin.Default()
    r.PUT("/order_product/:id", h.UpdateOrderProduct)

    mock.ExpectQuery(`SELECT \* FROM "order_products" WHERE "order_products"."id" = \$1 ORDER BY "order_products"."id" LIMIT \$2`).
        WithArgs("999",1).
        WillReturnError(gorm.ErrRecordNotFound)

    req := httptest.NewRequest(http.MethodPut,"/order_product/999",bytes.NewBufferString(`{"quantity":5}`))
    req.Header.Set("Content-Type","application/json")
    w := performRequest(r,req)

    if w.Code != http.StatusNotFound {
        t.Fatalf("expected 404, got %d", w.Code)
    }
}

func TestUpdateOrderProduct_InvalidQuantity(t *testing.T) {
    db, mock := setupMockDB(t)
    h := &OrderProductHandler{DB: db}
    r := gin.Default()
    r.PUT("/order_product/:id", h.UpdateOrderProduct)

    mock.ExpectQuery(`SELECT \* FROM "order_products" WHERE "order_products"."id" = \$1 ORDER BY "order_products"."id" LIMIT \$2`).
        WithArgs("1",1).
        WillReturnRows(sqlmock.NewRows([]string{"id","order_id","product_id","quantity"}).
            AddRow(1,1,10,2))

    req := httptest.NewRequest(http.MethodPut,"/order_product/1",bytes.NewBufferString(`{"quantity":0}`))
    req.Header.Set("Content-Type","application/json")
    w := performRequest(r,req)

    if w.Code != http.StatusBadRequest {
        t.Fatalf("expected 400, got %d", w.Code)
    }
}

func TestUpdateOrderProduct_BadJSON(t *testing.T) {
    db, mock := setupMockDB(t)
    h := &OrderProductHandler{DB: db}
    r := gin.Default()
    r.PUT("/order_product/:id", h.UpdateOrderProduct)

    mock.ExpectQuery(`SELECT \* FROM "order_products" WHERE "order_products"."id" = \$1 ORDER BY "order_products"."id" LIMIT \$2`).
        WithArgs("1",1).
        WillReturnRows(sqlmock.NewRows([]string{"id","order_id","product_id","quantity"}).
            AddRow(1,1,10,2))

    req := httptest.NewRequest(http.MethodPut,"/order_product/1",bytes.NewBufferString("not-json"))
    req.Header.Set("Content-Type","application/json")
    w := performRequest(r,req)

    if w.Code != http.StatusBadRequest {
        t.Fatalf("expected 400, got %d", w.Code)
    }
}

func TestUpdateOrderProduct_SaveError(t *testing.T) {
    db, mock := setupMockDB(t)
    h := &OrderProductHandler{DB: db}
    r := gin.Default()
    r.PUT("/order_product/:id", h.UpdateOrderProduct)

    mock.ExpectQuery(`SELECT \* FROM "order_products" WHERE "order_products"."id" = \$1 ORDER BY "order_products"."id" LIMIT \$2`).
        WithArgs("1",1).
        WillReturnRows(sqlmock.NewRows([]string{"id","order_id","product_id","quantity"}).
            AddRow(1,1,10,2))

    mock.ExpectBegin()
    mock.ExpectExec(`UPDATE "order_products"`).WillReturnError(errors.New("update fail"))
    mock.ExpectRollback()

    req := httptest.NewRequest(http.MethodPut,"/order_product/1",bytes.NewBufferString(`{"quantity":5}`))
    req.Header.Set("Content-Type","application/json")
    w := performRequest(r,req)

    if w.Code != http.StatusInternalServerError {
        t.Fatalf("expected 500, got %d", w.Code)
    }
}

func TestUpdateOrderProduct_Success(t *testing.T) {
    db, mock := setupMockDB(t)
    h := &OrderProductHandler{DB: db}
    r := gin.Default()
    r.PUT("/order_product/:id", h.UpdateOrderProduct)

    // Existing product
    mock.ExpectQuery(`SELECT \* FROM "order_products" WHERE "order_products"."id" = \$1 ORDER BY "order_products"."id" LIMIT \$2`).
        WithArgs("1",1).
        WillReturnRows(sqlmock.NewRows([]string{"id","order_id","product_id","quantity"}).
            AddRow(1,1,10,2))

    // Update succeeds
    mock.ExpectBegin()
    mock.ExpectExec(`UPDATE "order_products"`).
        WillReturnResult(sqlmock.NewResult(1,1))
    mock.ExpectCommit()

    // Reload product
    mock.ExpectQuery(`SELECT \* FROM "order_products" WHERE "order_products"."id" = \$1 ORDER BY "order_products"."id" LIMIT \$2`).
        WithArgs(1,1).
        WillReturnRows(sqlmock.NewRows([]string{"id","order_id","product_id","quantity"}).
            AddRow(1,1,10,5))

    req := httptest.NewRequest(http.MethodPut,"/order_product/1",bytes.NewBufferString(`{"quantity":5}`))
    req.Header.Set("Content-Type","application/json")
    w := performRequest(r,req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
    }
}



func TestDeleteOrderProduct_NotFound(t *testing.T) {
    db, mock := setupMockDB(t)
    h := &OrderProductHandler{DB: db}
    r := gin.Default()
    r.DELETE("/order_product/:id", h.DeleteOrderProduct)

    // Simulate record not found
    mock.ExpectQuery(`SELECT \* FROM "order_products" WHERE "order_products"."id" = \$1 ORDER BY "order_products"."id" LIMIT \$2`).
        WithArgs("999", 1).
        WillReturnError(gorm.ErrRecordNotFound)

    req := httptest.NewRequest(http.MethodDelete, "/order_product/999", nil)
    w := performRequest(r, req)

    if w.Code != http.StatusNotFound {
        t.Fatalf("expected 404, got %d", w.Code)
    }
}

func TestDeleteOrderProduct_Success(t *testing.T) {
    db, mock := setupMockDB(t)
    h := &OrderProductHandler{DB: db}
    r := gin.Default()
    r.DELETE("/order_product/:id", h.DeleteOrderProduct)

    // Record exists
    mock.ExpectQuery(`SELECT \* FROM "order_products" WHERE "order_products"."id" = \$1 ORDER BY "order_products"."id" LIMIT \$2`).
        WithArgs("1", 1).
        WillReturnRows(sqlmock.NewRows([]string{"id","order_id","product_id","quantity"}).
            AddRow(1, 1, 10, 2))

    // Transaction + delete
    mock.ExpectBegin()
    mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "order_products" WHERE "order_products"."id" = $1`)).
        WithArgs(1).
        WillReturnResult(sqlmock.NewResult(1, 1))
    mock.ExpectCommit()

    req := httptest.NewRequest(http.MethodDelete, "/order_product/1", nil)
    w := performRequest(r, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatalf("unmet expectations: %v", err)
    }
}

func TestDeleteOrderProduct_DBErrorOnLookup(t *testing.T) {
    db, mock := setupMockDB(t)
    h := &OrderProductHandler{DB: db}
    r := gin.Default()
    r.DELETE("/order_product/:id", h.DeleteOrderProduct)

    mock.ExpectQuery(`SELECT \* FROM "order_products" WHERE "order_products"."id" = \$1 ORDER BY "order_products"."id" LIMIT \$2`).
        WithArgs("1",1).
        WillReturnError(errors.New("db fail"))

    req := httptest.NewRequest(http.MethodDelete,"/order_product/1",nil)
    w := performRequest(r,req)

    if w.Code != http.StatusInternalServerError {
        t.Fatalf("expected 500, got %d", w.Code)
    }
}

func TestDeleteOrderProduct_DeleteError(t *testing.T) {
    db, mock := setupMockDB(t)
    h := &OrderProductHandler{DB: db}
    r := gin.Default()
    r.DELETE("/order_product/:id", h.DeleteOrderProduct)

    mock.ExpectQuery(`SELECT \* FROM "order_products" WHERE "order_products"."id" = \$1 ORDER BY "order_products"."id" LIMIT \$2`).
        WithArgs("1",1).
        WillReturnRows(sqlmock.NewRows([]string{"id","order_id","product_id","quantity"}).
            AddRow(1,1,10,2))

    mock.ExpectBegin()
    mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "order_products" WHERE "order_products"."id" = $1`)).
        WithArgs(1).
        WillReturnError(errors.New("delete fail"))
    mock.ExpectRollback()

    req := httptest.NewRequest(http.MethodDelete,"/order_product/1",nil)
    w := performRequest(r,req)

    if w.Code != http.StatusInternalServerError {
        t.Fatalf("expected 500, got %d", w.Code)
    }
}






