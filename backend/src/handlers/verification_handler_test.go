package handlers

import (
	"app/requestModels"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"app/blockchain"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupMockDB creates a gorm DB backed by sqlmock

// Helpers
func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
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

// Tests
func TestVerifyOrder_NoBlockchain(t *testing.T) {
	db, _ := setupMockDB(t)
	h := &VerificationHandler{DB: db}
	r := gin.Default()
	r.GET("/order/verify/:order_id", h.VerifyOrder)

	req := httptest.NewRequest(http.MethodGet, "/order/verify/1", nil)
	w := performRequest(r, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp requestModels.VerificationResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Status != "BLOCKCHAIN_NOT_AVAILABLE" {
		t.Fatalf("expected BLOCKCHAIN_NOT_AVAILABLE, got %s", resp.Status)
	}
}

func TestVerifyOrder_ReturnsTransactionHashes(t *testing.T) {
	db, mock := setupMockDB(t)
	h := &VerificationHandler{DB: db}
	r := gin.Default()
	r.GET("/order/verify/:order_id", h.VerifyOrder)

	// Setup mock DB to return one order history row
	orderID := 1
	ts := time.Now().UTC().Truncate(time.Second)
	rows := sqlmock.NewRows([]string{
		"id",
		"order_id",
		"timestamp_history",
		"order_status",
		"order_location",
	}).AddRow(1, orderID, ts, "DELIVERED", "POINT(1 1)")
	mock.ExpectQuery(`SELECT .* FROM "order_status_history" WHERE .*`).
		WithArgs("1").WillReturnRows(rows)

	// Compute expected hash for the single update
	data := fmt.Sprintf(
		"%d|%s|%s|%s",
		orderID,
		"DELIVERED",
		ts.Format(time.RFC3339),
		"POINT(1 1)",
	)
	computed := sha256.Sum256([]byte(data))

	// Setup handler with fake blockchain client; non-nil client so it attempts to use contract
	h.Client = &blockchain.Client{EthClient: &ethclient.Client{}}

	// Inject GetUpdateHashesFunc to return our computed hash
	h.GetUpdateHashesFunc = func(
		contract *blockchain.Blockchain,
		orderIDBig *big.Int,
	) ([][32]byte, error) {
		return [][32]byte{computed}, nil
	}


	req := httptest.NewRequest(http.MethodGet, "/order/verify/1", nil)
	w := performRequest(r, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp requestModels.VerificationResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if !resp.Verified {
		t.Fatalf("expected verified true, got false: %+v", resp)
	}
}
