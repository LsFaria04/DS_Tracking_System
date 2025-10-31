package handlers

import (
	"app/blockchain"
	"app/models"
	"crypto/sha256"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VerificationHandler struct {
	DB     *gorm.DB
	Client *blockchain.Client
}

type VerificationResponse struct {
	Verified         bool     `json:"verified"`
	TotalUpdates     int      `json:"total_updates"`
	VerifiedUpdates  int      `json:"verified_updates"`
	BlockchainHashes int      `json:"blockchain_hashes"`
	Status           string   `json:"status"`
	Message          string   `json:"message"`
	Mismatches       []string `json:"mismatches,omitempty"`
}

// VerifyOrder verifies all updates for an order against blockchain
func (h *VerificationHandler) VerifyOrder(c *gin.Context) {
	orderID := c.Param("order_id")

	// Fetch all order status updates from database
	var orderHistory []models.OrderStatusHistory
	if err := h.DB.Where("Order_ID = ?", orderID).Order("Timestamp_History asc").Find(&orderHistory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order history"})
		return
	}

	if len(orderHistory) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No order history found"})
		return
	}

	// Get blockchain contract instance
	ethClient := h.Client.EthClient
	addr := os.Getenv("BLOCKCHAIN_CONTRACT_ADDRESS")
	contract, err := blockchain.GetContractInstance(ethClient, addr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to blockchain"})
		return
	}

	// Get all hashes from blockchain for this order
	orderIDBigInt := new(big.Int)
	orderIDBigInt.SetString(orderID, 10)
	blockchainHashes, err := contract.GetUpdateHash(nil, orderIDBigInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blockchain hashes"})
		return
	}

	response := VerificationResponse{
		TotalUpdates:     len(orderHistory),
		BlockchainHashes: len(blockchainHashes),
		Mismatches:       []string{},
	}

	// Verify each update
	verifiedCount := 0
	for i, update := range orderHistory {
		// Compute hash for this update (same logic as when storing)
		data := fmt.Sprintf("%d|%s|%s|%s|%d",
			update.Order_ID,
			update.Order_Status,
			update.Timestamp_History.Format(time.RFC3339),
			update.Order_Location,
			update.Storage_ID,
		)
		computedHash := sha256.Sum256([]byte(data))

		// Check if this hash exists in blockchain
		found := false
		for j, blockchainHash := range blockchainHashes {
			if computedHash == blockchainHash {
				found = true
				verifiedCount++
				break
			}
			// For debugging - check if it's just in wrong order
			_ = j
		}

		if !found {
			response.Mismatches = append(response.Mismatches, fmt.Sprintf("Update #%d (%s) not found in blockchain", i+1, update.Order_Status))
		}
	}

	response.VerifiedUpdates = verifiedCount
	response.Verified = (verifiedCount == len(orderHistory)) && (len(orderHistory) == len(blockchainHashes))

	// Determine status message
	if response.Verified {
		response.Status = "VERIFIED"
		response.Message = "All order updates are verified on the blockchain"
	} else if verifiedCount == 0 {
		response.Status = "NOT_VERIFIED"
		response.Message = "No updates found on blockchain"
	} else if verifiedCount < len(orderHistory) {
		response.Status = "PARTIALLY_VERIFIED"
		response.Message = fmt.Sprintf("Only %d out of %d updates are verified", verifiedCount, len(orderHistory))
	} else if len(blockchainHashes) > len(orderHistory) {
		response.Status = "EXTRA_HASHES"
		response.Message = "More hashes on blockchain than in database"
	} else {
		response.Status = "MISMATCH"
		response.Message = "Database and blockchain data mismatch"
	}

	c.JSON(http.StatusOK, response)
}
