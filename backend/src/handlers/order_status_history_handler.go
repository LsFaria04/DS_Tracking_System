package handlers

import (
	"app/blockchain"
	"app/models"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderStatusHistoryHandler struct {
	DB     *gorm.DB
	Client *blockchain.Client
}

func (h *OrderStatusHistoryHandler) GetOrderStatusByOrderID(c *gin.Context) {
	id := c.Param("order_id")
	var orderStatus []models.OrderStatusHistory
	result := h.DB.Preload("Storage").Where("Order_ID = ?", id).Order("Timestamp_History desc").Find(&orderStatus)

	//check if there was an error with the database request
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order Status not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}

	} else {
		c.JSON(http.StatusOK, gin.H{"order_status_history": orderStatus})
	}
}

func (h *OrderStatusHistoryHandler) AddOrderUpdate(c *gin.Context) {

	//get the order status from the post request
	var input models.OrderStatusHistory
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Input"})
		return
	}

	//assign a value to timestamp if there is none
	if input.Timestamp_History.IsZero() {
		input.Timestamp_History = time.Now()
	}

	if h.Client != nil {
		//get an instance of the contract
		auth := h.Client.Auth
		ethClient := h.Client.EthClient
		addr := os.Getenv("BLOCKCHAIN_CONTRACT_ADDRESS")
		contract, err := blockchain.GetContractInstance(ethClient, addr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		//hash the main components of the request
		data := fmt.Sprintf("%d|%s|%s|%s",
			input.Order_ID,
			input.Order_Status,
			input.Timestamp_History.Format(time.RFC3339),
			input.Order_Location,
		)

		hash := sha256.Sum256([]byte(data))

		//store the hash in the blockchain
		if err := StoreUpdateHash(auth, contract, uint64(input.Order_ID), hash); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save update"})
			return
		}
	}

	//store the update into the database
	if err := h.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Update stored successfully"})
}
