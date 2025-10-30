package handlers

import (
	"app/models"
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderHandler struct{
	DB *gorm.DB
}
	
func (h *OrderHandler) GetOrderByID(c *gin.Context){
	id := c.Param("id")

	var order models.Orders
	result := h.DB.Preload("Products.Product").First(&order, id)

	//check if there was an error with the database request
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        }
        
    } else{
		c.JSON(http.StatusOK,gin.H{"order" : order})
	}
	
}