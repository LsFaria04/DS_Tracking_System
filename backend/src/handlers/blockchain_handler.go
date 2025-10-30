package handlers

import (
	"app/blockchain"
	"app/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BlockchainHandler struct{}

// GetBlockchainStatus returns the current blockchain connection status
func (h *BlockchainHandler) GetBlockchainStatus(c *gin.Context) {
	client, err := blockchain.NewClient()
	if err != nil {
		c.JSON(http.StatusOK, models.BlockchainStatusResponse{
			Connected: false,
			Network:   "sepolia",
			Error:     err.Error(),
		})
		return
	}
	defer client.Close()

	response := models.BlockchainStatusResponse{
		Connected:     true,
		Network:       "sepolia",
		WalletAddress: client.WalletAddress.Hex(),
	}

	// Get wallet balance
	balance, err := client.GetWalletBalance()
	if err == nil {
		response.WalletBalance = blockchain.FormatBalance(balance)
	}

	// Get current block number
	blockNumber, err := client.GetBlockNumber()
	if err == nil {
		response.BlockNumber = blockNumber
	}

	// Add contract address if configured
	if client.ContractAddress.Hex() != "0x0000000000000000000000000000000000000000" {
		response.ContractAddress = client.ContractAddress.Hex()
	}

	c.JSON(http.StatusOK, response)
}

func (h *BlockchainHandler) DeployContract(c *gin.Context){
	address, err := blockchain.DeployContract()

	if err != nil {
		c.JSON(
		http.StatusInternalServerError,
		gin.H{
			"message": err.Error(),
	})
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "Contract deployed",
			"address": address,
	})

}


