package models

import "time"

type Storage struct {
    Id        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"not null"`
    Address   string    `gorm:"type:text;not null"`
    Latitude  float64   `gorm:"type:decimal(10,8);not null"`
    Longitude float64   `gorm:"type:decimal(11,8);not null"`
    Created_At time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (Storage) TableName() string {
    return "storages"
}