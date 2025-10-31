package handlers

import (
	"app/blockchain"
	"app/models"
	"app/requestModels"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderHandler struct{
	DB *gorm.DB
	Client *blockchain.Client
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
		fmt.Printf("Delivery: %s", order.Delivery_Estimate)
		c.JSON(http.StatusOK,gin.H{"order" : order})
	}
	
}

func (h *OrderHandler) AddOrder(c *gin.Context){

	//get the order request
    var input requestModels.AddOrderRequest
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Input"})
        return
    }

	//assign a unique tracking code
	trackingCode := uuid.New().String()

	//assign the values to the order model
	var order models.Orders
	order.Tracking_Code = trackingCode
	order.Customer_ID = input.CustomerId
	order.Delivery_Address = input.DeliveryAddress
	order.Delivery_Estimate = time.Now().Add(48 * time.Hour) //just a mock estimate for now
	order.Delivery_Latitude = input.DeliveryLatitude
	order.Delivery_Longitude = input.DeliveryLongitude
	order.Seller_Address = input.SellerAddress
	order.Seller_ID = input.SellerId
	order.Seller_Latitude = input.SellerLatitude
	order.Seller_Longitude = input.SellerLongitude
	order.Created_At = time.Now()

	
	//the next operations are made inside a transaction to ensure atomicity
	transaction := h.DB.Begin()

	result := transaction.Create(&order)
	if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "Order not Created"})
			transaction.Rollback()
			return
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			transaction.Rollback()
			return
        }
	}

	//create the order products associated to the order
	for _, productRequest := range input.Products{
		var orderProduct models.OrderProduct
		orderProduct.OrderID = order.Id
		orderProduct.ProductID = productRequest.ProductID
		orderProduct.Quantity = productRequest.Quantity

		result := transaction.Create(&orderProduct)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Order was created with errors in the products"})
				transaction.Rollback()
				return
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				transaction.Rollback()
				return
			}
		}

	}

	//insert a first update (processing)
	var statusHistory models.OrderStatusHistory
	
	//get an instance of the contract
    auth := h.Client.Auth
    ethClient := h.Client.EthClient
    addr := os.Getenv("BLOCKCHAIN_CONTRACT_ADDRESS")
    contract, err := blockchain.GetContractInstance(ethClient, addr)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		transaction.Rollback()
        return
    }

    statusHistory.Note = "Processing the Order"
	statusHistory.Order_ID = order.Id
	statusHistory.Order_Location = order.Seller_Address
	statusHistory.Timestamp_History = time.Now()
	statusHistory.Order_Status = "PROCESSING"

    //hash the main components of the request
    data := fmt.Sprintf("%d|%s|%s|%s",
        statusHistory.Order_ID,
        statusHistory.Order_Status,
        statusHistory.Timestamp_History.Format(time.RFC3339),
        statusHistory.Order_Location,
    )

    hash := sha256.Sum256([]byte(data))

    //store the hash in the blockchain
    if err := StoreUpdateHash(auth, contract, uint64(statusHistory.Order_ID), hash); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save update"})
		transaction.Rollback()
        return
    }

    //store the update into the database
    if err := transaction.Create(&statusHistory).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save update"})
		transaction.Rollback()
        return
    }

	//commits the transaction
	transaction.Commit()

    c.JSON(http.StatusOK, gin.H{"message": "Update stored successfully"})

}