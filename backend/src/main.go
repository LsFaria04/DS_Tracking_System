package main

import (
	"app/routes"
	"fmt"
	"os"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func main() {
  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT") )
  db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

  router := gin.Default()

  // Configure CORS middleware
  router.Use(cors.New(cors.Config{
    AllowOrigins: []string{
      "http://localhost:3000",
      "https://production-url.com", // Add here the production url when deployed to GCP Cloud Run
    },
    AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
    ExposeHeaders: []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge: 12 * time.Hour,
  }))

  //registers the routes
  routes.RegisterRoutes(router, db)

  router.Run(":8080") // listens on 0.0.0.0:8080 by default
}