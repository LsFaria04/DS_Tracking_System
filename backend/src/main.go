package main

import "github.com/gin-gonic/gin"
import (
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
)

import "os"
import "fmt"



func main() {
  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT") )
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

  // Simple raw query to test connection
  sqlDB, err := db.DB()

  err = sqlDB.Ping()

  router := gin.Default()
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