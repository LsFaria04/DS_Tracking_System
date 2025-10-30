package routes

import (
	"app/blockchain"
	"app/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, blockChainClient *blockchain.Client){
	orderHandler := handlers.OrderHandler{DB : db}
	orderStatusHistory := handlers.OrderStatusHistoryHandler{DB : db, Client: blockChainClient}
	storageHandler := handlers.StorageHandler{DB : db}
	blockchainHandler := handlers.BlockchainHandler{}

	//routes for the orders
	router.GET("/order/:id", orderHandler.GetOrderByID)

	//routes for the order history
	router.GET("/order/history/:order_id", orderStatusHistory.GetOrderStatusByOrderID)
	router.POST("/order/history/add", orderStatusHistory.AddOrderUpdate)

	//routes for the storages
    router.GET("/storages", storageHandler.GetAllStorages)

	// Blockchain endpoints (should not be public in the production)
	router.GET("/blockchain/status", blockchainHandler.GetBlockchainStatus)
	router.GET("/blockchain/deploy", blockchainHandler.DeployContract)

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