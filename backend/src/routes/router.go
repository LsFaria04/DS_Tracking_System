package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"app/handlers"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB){
	orderHandler := handlers.OrderHandler{DB : db}
	orderStatusHistory := handlers.OrderStatusHistoryHandler{DB : db}
	storageHandler := handlers.StorageHandler{DB : db}
	blockchainHandler := handlers.BlockchainHandler{}

	//routes for the orders
	router.GET("/order/:id", orderHandler.GetOrderByID)

	//routes for the order history
	router.GET("/order/history/:order_id", orderStatusHistory.GetOrderStatusByOrderID)

	router.GET("/order/history/set/:order_id", orderStatusHistory.AddOrderUpdate)

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