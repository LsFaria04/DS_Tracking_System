package main

import (
	"app/blockchain"
	"app/routes"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//configure the database connection using gorm
func configDB() (*gorm.DB) {
  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT") )
  db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  return db  
}

//Configure the router that will be used for the API
func configRouter(db *gorm.DB) (*gin.Engine){
  router := gin.Default()

  // Configure CORS middleware (Allow frontend and localhost)

  router.Use(cors.New(cors.Config{
    AllowOriginFunc: func(origin string) bool {
      // Allow localhost for development 
      
      if strings.HasPrefix(origin, "http://localhost") {
        return true
      }
      // Allow any .run.app domain (Cloud Run)
      if strings.HasSuffix(origin, ".run.app") {
        return true
      }
      return false
    },
    AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
    ExposeHeaders: []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge: 12 * time.Hour,
  }))

  //registers the routes
  routes.RegisterRoutes(router, db)
  return router
}

func main() {
  db := configDB()

  // Initialize blockchain client (Sepolia testnet)
  blockchainClient, err := blockchain.NewClient()
  if err != nil {
    log.Printf("⚠️  Warning: Failed to initialize blockchain client: %v", err)
    log.Println("⚠️  Backend will continue without blockchain functionality")
  } else {
    defer blockchainClient.Close()
    
    // Test blockchain connection
    balance, err := blockchainClient.GetWalletBalance()
    if err != nil {
      log.Printf("⚠️  Warning: Failed to get wallet balance: %v", err)
    } else {
      log.Printf("💰 Wallet balance: %s", blockchain.FormatBalance(balance))
    }
    
    blockNumber, err := blockchainClient.GetBlockNumber()
    if err != nil {
      log.Printf("⚠️  Warning: Failed to get block number: %v", err)
    } else {
      log.Printf("🔗 Current Sepolia block: %d", blockNumber)
    }
  }

  router := configRouter(db)

  router.Run(":8080") // listens on 0.0.0.0:8080 by default
}