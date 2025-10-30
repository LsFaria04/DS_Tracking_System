package handlers

import (
    "app/models"
    "net/http"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type StorageHandler struct {
    DB *gorm.DB
}

func (h *StorageHandler) GetAllStorages(c *gin.Context) {
    var storages []models.Storage
    result := h.DB.Order("id ASC").Find(&storages)

    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"storages": storages})
}