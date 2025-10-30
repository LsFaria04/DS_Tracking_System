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
	Updates []OrderStatusHistory `gorm:"foreignKey:Order_ID"`
}



