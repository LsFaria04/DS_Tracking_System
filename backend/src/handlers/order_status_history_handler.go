package handlers

import (
    "app/blockchain"
    "app/models"
    "errors"
    "encoding/hex"
    "log"
    "math/big"
    "net/http"
    "os"
    "strconv"

    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type OrderStatusHistoryHandler struct{
    DB *gorm.DB
}
    
func (h *OrderStatusHistoryHandler) GetOrderStatusByOrderID(c *gin.Context){
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
        
    } else{
        c.JSON(http.StatusOK,gin.H{"order_status_history" : orderStatus})
    }
}

//TODO: Is only to test the blockchain for now and still needs to be improved
func (h *OrderStatusHistoryHandler) AddOrderUpdate(c *gin.Context){
    id := c.Param("order_id")

    // parse order id to uint64
    orderID, err := strconv.ParseUint(id, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
        return
    }

    client, err := blockchain.NewClient()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }
    auth := client.Auth
    ethClient := client.EthClient
    add := os.Getenv("BLOCKCHAIN_CONTRACT_ADDRESS")
    log.Printf("address : %s", add)
    contract, err := blockchain.GetContractInstance(ethClient, "0x472b477d30c45cfbd89e76d9f9700ad1f90cc370")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get contract instance"})
        return
    }

    hexStr := "012345" //mock str to hash
    if len(hexStr) >= 2 && hexStr[:2] == "0x" {
        hexStr = hexStr[2:]
    }
    b, err := hex.DecodeString(hexStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hash hex"})
        return
    }
    if len(b) > 32 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "hash too long"})
        return
    }
    var hash [32]byte
    copy(hash[:], b)

    if err := StoreUpdateHash(ethClient, auth, contract, orderID, hash); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }


    c.JSON(http.StatusOK, gin.H{"status": "hash stored"})
}

//stores the hash in the block chain
func StoreUpdateHash(client *ethclient.Client, auth *bind.TransactOpts, contract *blockchain.Blockchain, orderID uint64, hash [32]byte) error {
    tx, err := contract.StoreUpdateHash(auth, big.NewInt(int64(orderID)), hash)
    if err != nil {
        return err
    }
    log.Printf("Stored hash on-chain. TX: %s", tx.Hash().Hex())
    return nil
}