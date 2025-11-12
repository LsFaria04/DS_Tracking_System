package handlers

import (
	"app/models"
	"app/requestModels"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandler struct {
	DB *gorm.DB
}

// GetAllProducts retrieves all available products from the jumpseller API
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/products.json", os.Getenv("JUMPSELLER_BASE_URL"))
	req, err := http.NewRequest("GET",url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	req.SetBasicAuth(os.Getenv("LOGIN_JUMPSELLER_API"), os.Getenv("TOKEN_JUMPSELLER_API"))

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var result []requestModels.ProductResponse
	json.NewDecoder(resp.Body).Decode(&result)
	err = resp.Body.Close()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":  err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"products": result})
}

// GetProductByID retrieves a specific product by ID from the jump seller API
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")

	var product *models.Product
	product, err := GetProductByIDAPI(id)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, gin.H{"product": &product})
}

//gets a product from the Jumpseller API
func GetProductByIDAPI(id string) (*models.Product, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/products/%s.json", os.Getenv("JUMPSELLER_BASE_URL"), id)
	req, err := http.NewRequest("GET",url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(os.Getenv("LOGIN_JUMPSELLER_API"), os.Getenv("TOKEN_JUMPSELLER_API"))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var result requestModels.ProductResponse
	json.NewDecoder(resp.Body).Decode(&result)
	err = resp.Body.Close()

	if err != nil {
		return nil, err
	}
	
	var product models.Product
	product.ID = uint(result.Product.ID)
	product.Name = result.Product.Name
	product.Price = result.Product.Price

	return &product, nil
}
