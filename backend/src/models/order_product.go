package models

type OrderProduct struct {
    ID        uint    `gorm:"primaryKey"`
    OrderID   uint    `gorm:"not null"`
    ProductID uint    `gorm:"not null"`
    Quantity  uint     `gorm:"not null"`
    Product   *Product `gorm:"foreignKey:ProductID;references:ID"`
}