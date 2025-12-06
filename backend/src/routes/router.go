package routes

import (
	"app/blockchain"
	"app/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, blockChainClient *blockchain.Client) {
	orderHandler := handlers.OrderHandler{DB: db, Client: blockChainClient}
	orderStatusHistory := handlers.OrderStatusHistoryHandler{DB: db, Client: blockChainClient}
	storageHandler := handlers.StorageHandler{DB: db}
	orderProductHandler := handlers.OrderProductHandler{DB: db}
	productHandler := handlers.ProductHandler{DB: db}
	blockchainHandler := handlers.BlockchainHandler{}
	verificationHandler := handlers.VerificationHandler{DB: db, Client: blockChainClient}

	apiRoutes := router.Group("/api")

	//routes for the order history
	apiRoutes.GET("/order/history/:order_id", orderStatusHistory.GetOrderStatusByOrderID)
	apiRoutes.POST("/order/history/add", orderStatusHistory.AddOrderUpdate)

	//routes for the orders
	apiRoutes.GET("/orders", orderHandler.GetAllOrders)
	apiRoutes.GET("/order/:id", orderHandler.GetOrderByID)
	apiRoutes.GET("/order/verify/:order_id", verificationHandler.VerifyOrder)
	apiRoutes.POST("/order/add", orderHandler.AddOrder)
	apiRoutes.POST("/order/update", orderHandler.UpdateOrder)
	apiRoutes.POST("/order/cancel", orderHandler.CancelOrder)

	//routes for order products (using order-products path to avoid conflicts)
	apiRoutes.GET("/order-products", orderProductHandler.GetOrderProducts) // Query param: ?order_id=X
	apiRoutes.POST("/order-products", orderProductHandler.AddOrderProduct)
	apiRoutes.GET("/order-products/:id", orderProductHandler.GetOrderProductByID)
	apiRoutes.PUT("/order-products/:id", orderProductHandler.UpdateOrderProduct)
	apiRoutes.DELETE("/order-products/:id", orderProductHandler.DeleteOrderProduct)

	//routes for products
	apiRoutes.GET("/products", productHandler.GetAllProducts)
	apiRoutes.GET("/products/:id", productHandler.GetProductByID)

	//routes for the storages
	apiRoutes.GET("/storages", storageHandler.GetAllStorages)

	// Blockchain endpoints (should not be public in the production)
	apiRoutes.GET("/blockchain/status", blockchainHandler.GetBlockchainStatus)
	apiRoutes.GET("/blockchain/deploy", blockchainHandler.DeployContract)

	//old routes for testing
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Online", // "DS is awesome!"
		})
	})
}
