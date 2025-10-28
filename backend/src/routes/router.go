package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"app/handlers"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB){
	orderHandler := handlers.OrderHandler{DB : db}
	blockchainHandler := handlers.BlockchainHandler{}

	//routes for the orders
	router.GET("/order/:id", orderHandler.GetOrderByID)

	// Blockchain status endpoint
	router.GET("/blockchain/status", blockchainHandler.GetBlockchainStatus)

	//old routes
	router.GET("/pong", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "pong",
    })
  })
  router.GET("/", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "DS is awesome!",
    })
  })
}