package handlers

import (
	"app/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderProductHandler struct {
	DB *gorm.DB
}

// GetOrderProducts retrieves all products for a specific order
func (h *OrderProductHandler) GetOrderProducts(c *gin.Context) {
	orderID := c.Query("order_id")

	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_id query parameter is required"})
		return
	}

	var orderProducts []models.OrderProduct
	result := h.DB.Preload("Product").Where("order_id = ?", orderID).Find(&orderProducts)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order_products": orderProducts})
}

// GetOrderProductByID retrieves a specific order product by ID
func (h *OrderProductHandler) GetOrderProductByID(c *gin.Context) {
	id := c.Param("id")

	var orderProduct models.OrderProduct
	result := h.DB.First(&orderProduct, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"order_product": orderProduct})
}

// AddOrderProduct adds a new product to an order
func (h *OrderProductHandler) AddOrderProduct(c *gin.Context) {
	var orderProduct models.OrderProduct

	if err := c.ShouldBindJSON(&orderProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if quantity is valid
	if orderProduct.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quantity must be greater than 0"})
		return
	}

	// Check if order exists
	var order models.Orders
	if err := h.DB.First(&order, orderProduct.Order_ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	// Create the order product
	result := h.DB.Create(&orderProduct)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order product"})
		return
	}

	// Load the product relation for the response
	h.DB.First(&orderProduct, orderProduct.ID)

	c.JSON(http.StatusCreated, gin.H{"order_product": orderProduct})
}

// UpdateOrderProduct updates the quantity of an order product
func (h *OrderProductHandler) UpdateOrderProduct(c *gin.Context) {
	id := c.Param("id")

	var orderProduct models.OrderProduct
	if err := h.DB.First(&orderProduct, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	var updateData struct {
		Quantity uint `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updateData.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quantity must be greater than 0"})
		return
	}

	orderProduct.Quantity = updateData.Quantity
	if err := h.DB.Save(&orderProduct).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order product"})
		return
	}

	// Load the product relation for the response
	h.DB.Preload("Product").First(&orderProduct, orderProduct.ID)

	c.JSON(http.StatusOK, gin.H{"order_product": orderProduct})
}

// DeleteOrderProduct removes a product from an order
func (h *OrderProductHandler) DeleteOrderProduct(c *gin.Context) {
	id := c.Param("id")

	var orderProduct models.OrderProduct
	if err := h.DB.First(&orderProduct, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	if err := h.DB.Delete(&orderProduct).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order product deleted successfully"})
}
