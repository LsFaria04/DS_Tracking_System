package main

import (
    "app/blockchain"
    "app/routes"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
)


func TestConfigBlockChainClient_NoRPCURL(t *testing.T) {
    os.Setenv("BLOCKCHAIN_RPC_URL", "")
    client, err := configBlockChainClient()
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if client != nil {
        t.Errorf("expected nil client when BLOCKCHAIN_RPC_URL is empty")
    }
}

func TestConfigRouter_PingRoute(t *testing.T) {
    r := gin.Default()
    // Use nil DB and dummy blockchain client
    routes.RegisterRoutes(r, &gorm.DB{}, &blockchain.Client{})

    req := httptest.NewRequest(http.MethodGet, "/ping", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", w.Code)
    }
    if w.Body.String() != `{"message":"pong"}` {
        t.Errorf("unexpected body: %s", w.Body.String())
    }
}

func TestConfigRouter_CORS(t *testing.T) {
    r, _, err := configRouter(&gorm.DB{})
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    // Simulate request with allowed origin
    req := httptest.NewRequest(http.MethodGet, "/ping", nil)
    req.Header.Set("Origin", "http://localhost:3000")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
        t.Errorf("CORS not applied correctly, got %s", w.Header().Get("Access-Control-Allow-Origin"))
    }
}
