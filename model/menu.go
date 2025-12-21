package model

import "time"

type Menu struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"size:150" json:"name"`
    Description string    `gorm:"size:500" json:"description"`
    Price       float64   `json:"price"`
    CategoryID  *uint     `json:"category_id"`
    Category    Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
    ImageURL    string    `gorm:"size:512" json:"image_url"`
    IsAvailable bool      `gorm:"default:true" json:"is_available"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
