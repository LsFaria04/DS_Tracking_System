package models

type OrderProduct struct {
    ID        uint    `gorm:"primaryKey"`
    OrderID   uint    `gorm:"not null"`
    ProductID uint    `gorm:"not null"`
    Quantity  uint     `gorm:"not null"`
    ProductNameAtPurchase string `gorm:"not null"`
    ProductPriceAtPurchase float64 `gorm:"not null"`
}