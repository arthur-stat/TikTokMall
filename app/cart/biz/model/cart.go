package model

import (
	"time"
)

// CartItem 购物车项
type CartItem struct {
	ID        uint32    `gorm:"primaryKey;autoIncrement"`
	UserID    uint32    `gorm:"not null;index"`
	ProductID uint32    `gorm:"not null"`
	Quantity  uint32    `gorm:"not null;default:1"`
	Selected  bool      `gorm:"not null;default:true"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;ON UPDATE CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (CartItem) TableName() string {
	return "cart_items"
}
