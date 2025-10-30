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
func configDB() (*gorm.DB, error) {
  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT") )
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  return db, err 
}

func configBlockChainClient() (*blockchain.Client, error){
  //connect to the block chain client
    client, err := blockchain.NewClient()
    if err != nil {
        
        return nil, err
    }
    return client, nil
}

//Configure the router that will be used for the API
func configRouter(db *gorm.DB) (*gin.Engine, error){
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

  blockChainClient, err := configBlockChainClient()

  if err != nil{
    return nil, err
  }

  //registers the routes
  routes.RegisterRoutes(router, db, blockChainClient)
  return router, nil
}

func main() {
  db, err := configDB()

  if err != nil {
    log.Printf("Error while conecting to the database: %v", err)
    return
  }

  router,err := configRouter(db)

  if err != nil{
    log.Printf("Error while configuring the routing: %v", err)
    return
  }

  router.Run(":8080") // listens on 0.0.0.0:8080 by default
}