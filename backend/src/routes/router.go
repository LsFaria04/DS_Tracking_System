package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"app/handlers"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB){
	orderHandler := handlers.OrderHandler{DB : db}

	//routes for the orders
	router.GET("/order/:id", orderHandler.GetOrderByID)

	//old routes
	router.GET("/pong", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "pong",
    })
  })
  router.GET("/", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "DS is awes!",
    })
  })
}