package models

import (
	"time"
)

type OrderStatusHistory struct{
    Id                uint      `gorm:"primaryKey"`
    Order_ID          uint      `gorm:"not null"`
    Timestamp_History time.Time `gorm:"not null"`
    Order_Status      string    `gorm:"not null"`
    Note              string
    Blockchain_Transaction string  `gorm:"not null"`   
    Order_Location    string    `gorm:"not null"`
    Storage_ID        *uint     `gorm:"default:null"`
    Order             *Orders   `gorm:"foreignKey:Order_ID;references:Id"`
    Storage           *Storage  `gorm:"foreignKey:Storage_ID;references:Id"`
}

func (OrderStatusHistory) TableName() string {
    return "order_status_history" 
}