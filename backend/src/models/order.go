package models

import "time"

type Orders struct{
    Id                  uint      `gorm:"primaryKey"`
    Customer_ID         uint       `gorm:"not null"`
    Seller_ID           uint       `gorm:"not null"`
    Seller_Address      string    `gorm:"not null"`
    Seller_Latitude     float64  `gorm:"type:decimal(10,8)"`
    Seller_Longitude    float64  `gorm:"type:decimal(11,8)"`
    Created_At          time.Time
    Tracking_Code       string    `gorm:"unique;not null"`
    Delivery_Estimate  time.Time
    Delivery_Address    string    `gorm:"not null"`
    Delivery_Latitude   float64  `gorm:"type:decimal(10,8)"`
    Delivery_Longitude  float64  `gorm:"type:decimal(11,8)"`
    Price               float64   `gorm:"not null;default:0"`

    Products []OrderProduct        `gorm:"foreignKey:OrderID"`
    Updates  []OrderStatusHistory  `gorm:"foreignKey:Order_ID"`
}