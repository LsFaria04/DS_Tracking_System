package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"app/handlers"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB){
	orderHandler := handlers.OrderHandler{DB : db}
	orderStatusHistory := handlers.OrderStatusHistoryHandler{DB : db}
	blockchainHandler := handlers.BlockchainHandler{}

	//routes for the orders
	router.GET("/order/:id", orderHandler.GetOrderByID)

	//routes for the order history
	router.GET("/order/history/:order_id", orderStatusHistory.GetOrderStatusByOrderID)

	// Blockchain status endpoint
	router.GET("/blockchain/status", blockchainHandler.GetBlockchainStatus)

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