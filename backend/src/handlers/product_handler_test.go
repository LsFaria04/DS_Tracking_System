package handlers

import (
	"app/models"
	"app/requestModels"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)



func setupTestEnv() {
	os.Setenv("JUMPSELLER_BASE_URL", "http://mock-api")
	os.Setenv("LOGIN_JUMPSELLER_API", "user")
	os.Setenv("TOKEN_JUMPSELLER_API", "token")
}

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGetAllProducts(t *testing.T) {
	setupTestEnv()

	// Mock external Jumpseller API
	mockProducts := []requestModels.ProductResponse{
		{Product: requestModels.Product{
			ID:    1,
			Name:  "Test Product",
			Price: 10.5,
		}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mockProducts)
	}))
	defer server.Close()

	// Override env to point to mock server
	os.Setenv("JUMPSELLER_BASE_URL", server.URL)

	// Prepare Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler := ProductHandler{}
	handler.GetAllProducts(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var response map[string][]requestModels.ProductResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	if len(response["products"]) != 1 {
		t.Fatalf("expected 1 product, got %d", len(response["products"]))
	}

	if response["products"][0].Product.Name != "Test Product" {
		t.Fatalf("unexpected product name")
	}
}

func TestGetAllProducts_Success(t *testing.T) {
	setupTestEnv()

	mockProducts := []requestModels.ProductResponse{
		{Product: requestModels.Product{ID: 1, Name: "Test", Price: 10}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mockProducts)
	}))
	defer server.Close()

	os.Setenv("JUMPSELLER_BASE_URL", server.URL)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler := ProductHandler{}
	handler.GetAllProducts(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetAllProducts_HTTPClientError(t *testing.T) {
	setupTestEnv()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hj, _ := w.(http.Hijacker)
		conn, _, _ := hj.Hijack()
		conn.Close()
	}))
	defer server.Close()

	os.Setenv("JUMPSELLER_BASE_URL", server.URL)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler := ProductHandler{}
	handler.GetAllProducts(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}


func TestGetProductByID(t *testing.T) {
	setupTestEnv()

	// Mock the API function
	GetProductByIDAPIFunc = func(id string) (*models.Product, error) {
		return &models.Product{
			ID:    1,
			Name:  "Mocked Product",
			Price: 20.0,
		}, nil
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	handler := ProductHandler{}
	handler.GetProductByID(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var response map[string]*models.Product
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["product"].Name != "Mocked Product" {
		t.Fatalf("unexpected product name")
	}
}

func TestGetProductByID_Success(t *testing.T) {
	GetProductByIDAPIFunc = func(id string) (*models.Product, error) {
		return &models.Product{ID: 1, Name: "OK"}, nil
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	handler := ProductHandler{}
	handler.GetProductByID(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetProductByID_Error(t *testing.T) {
	GetProductByIDAPIFunc = func(id string) (*models.Product, error) {
		return nil, fmt.Errorf("not found")
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	handler := ProductHandler{}
	handler.GetProductByID(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}


func TestGetProductByIDAPI(t *testing.T) {
	setupTestEnv()

	mockResponse := requestModels.ProductResponse{
		Product: requestModels.Product{
			ID:    5,
			Name:  "API Product",
			Price: 99.99,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	os.Setenv("JUMPSELLER_BASE_URL", server.URL)

	product, err := GetProductByIDAPI("5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if product.Name != "API Product" {
		t.Fatalf("expected API Product, got %s", product.Name)
	}

	if product.Price != 99.99 {
		t.Fatalf("unexpected price")
	}
}

func TestGetProductByIDAPI_Success(t *testing.T) {
	setupTestEnv()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(requestModels.ProductResponse{
			Product: requestModels.Product{
				ID:    10,
				Name:  "Test",
				Price: 99.9,
			},
		})
	}))
	defer server.Close()

	os.Setenv("JUMPSELLER_BASE_URL", server.URL)

	product, err := GetProductByIDAPI("10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if product.Name != "Test" {
		t.Fatalf("wrong product name")
	}
}

func TestGetProductByIDAPI_HTTPError(t *testing.T) {
	setupTestEnv()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hj, _ := w.(http.Hijacker)
		conn, _, _ := hj.Hijack()
		conn.Close()
	}))
	defer server.Close()

	os.Setenv("JUMPSELLER_BASE_URL", server.URL)

	_, err := GetProductByIDAPI("10")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

