package routes

import (
    "app/blockchain"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "testing"
)

func TestRegisterRoutes_AllEndpointsExist(t *testing.T) {
    r := gin.Default()
    RegisterRoutes(r, &gorm.DB{}, &blockchain.Client{})

    // Collect all registered routes
    routes := r.Routes()

    // Expected endpoints (method + path)
    expected := map[string]bool{
        "GET-/api/order/history/:order_id": true,
        "POST-/api/order/history/add":      true,
        "GET-/api/orders":                  true,
        "GET-/api/order/:id":               true,
        "GET-/api/order/verify/:order_id":  true,
        "POST-/api/order/add":              true,
        "POST-/api/order/update":           true,
        "GET-/api/order-products":          true,
        "POST-/api/order-products":         true,
        "GET-/api/order-products/:id":      true,
        "PUT-/api/order-products/:id":      true,
        "DELETE-/api/order-products/:id":   true,
        "GET-/api/products":                true,
        "GET-/api/products/:id":            true,
        "GET-/api/storages":                true,
        "GET-/api/blockchain/status":       true,
        "GET-/api/blockchain/deploy":       true,
        "GET-/ping":                        true,
        "GET-/":                            true,
    }

    // Check that all expected routes are present
    for _, rt := range routes {
        key := rt.Method + "-" + rt.Path
        if expected[key] {
            delete(expected, key)
        }
    }

    if len(expected) > 0 {
        t.Errorf("missing routes: %v", expected)
    }
}
