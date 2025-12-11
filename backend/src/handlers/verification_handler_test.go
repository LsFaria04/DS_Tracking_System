package handlers

import (
	"app/requestModels"
	"crypto/sha256"
	"encoding/json"
	"errors"
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
	"github.com/stretchr/testify/assert"
)

// setupMockDB creates a gorm DB backed by sqlmock


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

func TestVerifyOrder_DBError(t *testing.T) {
	db, mock := setupMockDB(t)
	h := &VerificationHandler{
		DB:     db,
		Client: &blockchain.Client{EthClient: &ethclient.Client{}},
	}
	r := gin.Default()
	r.GET("/order/verify/:order_id", h.VerifyOrder)

	mock.ExpectQuery(`SELECT .* FROM "order_status_history" WHERE .*`).
		WithArgs("1").
		WillReturnError(errors.New("db failure"))

	req := httptest.NewRequest(http.MethodGet, "/order/verify/1", nil)
	w := performRequest(r, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestVerifyOrder_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	h := &VerificationHandler{
		DB:     db,
		Client: &blockchain.Client{EthClient: &ethclient.Client{}},
	}
	r := gin.Default()
	r.GET("/order/verify/:order_id", h.VerifyOrder)

	// Return empty result set
	mock.ExpectQuery(`SELECT .* FROM "order_status_history" WHERE .*`).
		WithArgs("999").
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "timestamp_history", "order_status", "order_location"}))

	req := httptest.NewRequest(http.MethodGet, "/order/verify/999", nil)
	w := performRequest(r, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestVerifyOrder_BlockchainHashError(t *testing.T) {
	db, mock := setupMockDB(t)
	h := &VerificationHandler{
		DB:     db,
		Client: &blockchain.Client{EthClient: &ethclient.Client{}},
	}
	r := gin.Default()
	r.GET("/order/verify/:order_id", h.VerifyOrder)

	ts := time.Now().UTC().Truncate(time.Second)
	rows := sqlmock.NewRows([]string{
		"id", "order_id", "timestamp_history", "order_status", "order_location",
	}).AddRow(1, 1, ts, "PROCESSING", "Origin")
	mock.ExpectQuery(`SELECT .* FROM "order_status_history" WHERE .*`).
		WithArgs("1").
		WillReturnRows(rows)

	// Inject GetUpdateHashesFunc to return error
	h.GetUpdateHashesFunc = func(
		contract *blockchain.Blockchain,
		orderIDBig *big.Int,
	) ([][32]byte, error) {
		return nil, errors.New("blockchain hash retrieval failed")
	}

	req := httptest.NewRequest(http.MethodGet, "/order/verify/1", nil)
	w := performRequest(r, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestVerifyOrder_PartiallyVerified(t *testing.T) {
	db, mock := setupMockDB(t)
	h := &VerificationHandler{
		DB:     db,
		Client: &blockchain.Client{EthClient: &ethclient.Client{}},
	}
	r := gin.Default()
	r.GET("/order/verify/:order_id", h.VerifyOrder)

	ts := time.Now().UTC().Truncate(time.Second)
	// Two updates in DB
	rows := sqlmock.NewRows([]string{
		"id", "order_id", "timestamp_history", "order_status", "order_location",
	}).AddRow(1, 1, ts, "PROCESSING", "Origin").
		AddRow(2, 1, ts.Add(time.Hour), "SHIPPED", "Warehouse")
	mock.ExpectQuery(`SELECT .* FROM "order_status_history" WHERE .*`).
		WithArgs("1").
		WillReturnRows(rows)

	// Only first hash matches
	data1 := fmt.Sprintf("%d|%s|%s|%s", 1, "PROCESSING", ts.Format(time.RFC3339), "Origin")
	hash1 := sha256.Sum256([]byte(data1))

	h.GetUpdateHashesFunc = func(
		contract *blockchain.Blockchain,
		orderIDBig *big.Int,
	) ([][32]byte, error) {
		return [][32]byte{hash1}, nil // Only one hash, missing second
	}

	req := httptest.NewRequest(http.MethodGet, "/order/verify/1", nil)
	w := performRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp requestModels.VerificationResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "PARTIALLY_VERIFIED", resp.Status)
	assert.Equal(t, 1, resp.VerifiedUpdates)
	assert.False(t, resp.Verified)
}

func TestVerifyOrder_NotVerified(t *testing.T) {
	db, mock := setupMockDB(t)
	h := &VerificationHandler{
		DB:     db,
		Client: &blockchain.Client{EthClient: &ethclient.Client{}},
	}
	r := gin.Default()
	r.GET("/order/verify/:order_id", h.VerifyOrder)

	ts := time.Now().UTC().Truncate(time.Second)
	rows := sqlmock.NewRows([]string{
		"id", "order_id", "timestamp_history", "order_status", "order_location",
	}).AddRow(1, 1, ts, "PROCESSING", "Origin")
	mock.ExpectQuery(`SELECT .* FROM "order_status_history" WHERE .*`).
		WithArgs("1").
		WillReturnRows(rows)

	// Return empty hashes - nothing verified
	h.GetUpdateHashesFunc = func(
		contract *blockchain.Blockchain,
		orderIDBig *big.Int,
	) ([][32]byte, error) {
		return [][32]byte{}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/order/verify/1", nil)
	w := performRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp requestModels.VerificationResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "NOT_VERIFIED", resp.Status)
	assert.Equal(t, 0, resp.VerifiedUpdates)
}

func TestVerifyOrder_ExtraHashes(t *testing.T) {
	db, mock := setupMockDB(t)
	h := &VerificationHandler{
		DB:     db,
		Client: &blockchain.Client{EthClient: &ethclient.Client{}},
	}
	r := gin.Default()
	r.GET("/order/verify/:order_id", h.VerifyOrder)

	ts := time.Now().UTC().Truncate(time.Second)
	rows := sqlmock.NewRows([]string{
		"id", "order_id", "timestamp_history", "order_status", "order_location",
	}).AddRow(1, 1, ts, "PROCESSING", "Origin")
	mock.ExpectQuery(`SELECT .* FROM "order_status_history" WHERE .*`).
		WithArgs("1").
		WillReturnRows(rows)

	// Compute matching hash plus an extra one
	data := fmt.Sprintf("%d|%s|%s|%s", 1, "PROCESSING", ts.Format(time.RFC3339), "Origin")
	hash := sha256.Sum256([]byte(data))
	extraHash := sha256.Sum256([]byte("extra"))

	h.GetUpdateHashesFunc = func(
		contract *blockchain.Blockchain,
		orderIDBig *big.Int,
	) ([][32]byte, error) {
		return [][32]byte{hash, extraHash}, nil // More hashes than DB entries
	}

	req := httptest.NewRequest(http.MethodGet, "/order/verify/1", nil)
	w := performRequest(r, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp requestModels.VerificationResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "EXTRA_HASHES", resp.Status)
}

