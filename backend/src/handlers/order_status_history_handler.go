package handlers

import (
	"app/models"
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderStatusHistoryHandler struct{
	DB *gorm.DB
}
	
func (h *OrderStatusHistoryHandler) GetOrderStatusByOrderID(c *gin.Context){
	id := c.Param("order_id")
	var orderStatus []models.OrderStatusHistory
	result := h.DB.Where("Order_ID = ?", id).Order("Timestamp_History desc").Find(&orderStatus)

	//check if there was an error with the database request
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "Order Status not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        }
        
    } else{
		c.JSON(http.StatusOK,gin.H{"order_status_history" : orderStatus})
	}
	
}