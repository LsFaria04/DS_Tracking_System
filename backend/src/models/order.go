package models

import "time"

type Orders struct{
    Id               uint      `gorm:"primaryKey"`
    Customer_ID       int       `gorm:"not null"`
    Created_At        time.Time
    Tracking_Code     string    `gorm:"unique;not null"`
    Delivery_Estimates time.Time
    Delivery_Address  string    `gorm:"not null"`
    Price             float64   `gorm:"not null;default:0"`

    Products []OrderProduct `gorm:"foreignKey:OrderID"`
}

type OrderProduct struct {
    ID        uint    `gorm:"primaryKey"`
    OrderID   uint    `gorm:"not null"`
    ProductID uint    `gorm:"not null"`
    Quantity  int     `gorm:"not null"`
    Product   *Product `gorm:"foreignKey:ProductID;references:ID"`
}

type Product struct {
    ID    uint    `gorm:"primaryKey"`
    Name  string  `gorm:"not null"`
    Price float64 `gorm:"not null"`
}