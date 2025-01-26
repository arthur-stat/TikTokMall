package model

import (
    "time"
)

// CartItem represents a single item in a user's shopping cart
type CartItem struct {
    ID        uint32    `json:"id" gorm:"primaryKey;autoIncrement"`
    UserID    uint32    `json:"user_id" gorm:"not null;index:idx_user_id"`
    ProductID uint32    `json:"product_id" gorm:"not null"`
    Quantity  uint32    `json:"quantity" gorm:"not null;default:1"`
    Selected  bool      `json:"selected" gorm:"not null;default:true"`
    CreatedAt time.Time `json:"created_at" gorm:"not null;autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"not null;autoUpdateTime"`
}

// TableName specifies the table name for CartItem
func (CartItem) TableName() string {
    return "cart_items"
} 
