package handlers

import (
	"app/models"
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandler struct{
	DB *gorm.DB
}

// GetAllProducts retrieves all available products
func (h *ProductHandler) GetAllProducts(c *gin.Context){
	var products []models.Product
	result := h.DB.Find(&products)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

// GetProductByID retrieves a specific product by ID
func (h *ProductHandler) GetProductByID(c *gin.Context){
	id := c.Param("id")

	var product models.Product
	result := h.DB.First(&product, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}
