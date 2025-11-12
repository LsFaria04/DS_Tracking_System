package models

type OrderProduct struct {
    ID        uint    `gorm:"primaryKey"`
    Order_ID   uint    `gorm:"not null"`
    Product_ID uint    `gorm:"not null"`
    Quantity  uint     `gorm:"not null"`
    Product_Name_At_Purchase string `gorm:"not null"`
    Product_Price_At_Purchase float64 `gorm:"not null"`
}