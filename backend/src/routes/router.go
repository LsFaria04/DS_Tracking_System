package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"app/handlers"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB){
	orderHandler := handlers.OrderHandler{DB : db}

	//ORDER ROUTES
	router.GET("/order/:id", orderHandler.GetOrderByID)

}