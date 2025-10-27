package main

import "github.com/gin-gonic/gin"
import "github.com/gin-contrib/cors"
import (
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
)

import "os"
import "fmt"
import "time"
import "strings"


func main() {
  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT") )
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

  // Simple raw query to test connection
  sqlDB, err := db.DB()

  err = sqlDB.Ping()

  router := gin.Default()

  // Configure CORS middleware (Allow frontend and localhost)

  router.Use(cors.New(cors.Config{
    AllowOriginFunc: func(origin string) bool {
      // Allow localhost for development 

      fmt.Println("CORS Origin:", origin)
      
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

  router.GET("/pong", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "pong",
    })
  })
  router.GET("/", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "DS is awesome!",
      "connected": err,
    })
  })

  router.Run() // listens on 0.0.0.0:8080 by default
}