package mysql

import (
	"database/sql"
	"time"
)

// Order 订单模型
type Order struct {
	ID              int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo         string         `gorm:"type:varchar(32);not null;uniqueIndex" json:"order_no"`
	UserID          uint32         `gorm:"not null;index" json:"user_id"`
	UserCurrency    string         `gorm:"type:varchar(3);not null" json:"user_currency"`
	TotalAmount     float64        `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status          int8           `gorm:"type:tinyint;not null;default:1" json:"status"`
	Email           string         `gorm:"type:varchar(128)" json:"email"`
	ShippingAddress sql.NullString `gorm:"type:json" json:"shipping_address"`
	CreatedAt       time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP;ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

// OrderItem 订单项模型
type OrderItem struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID   int64     `gorm:"not null;index" json:"order_id"`
	ProductID uint32    `gorm:"not null" json:"product_id"`
	Quantity  int32     `gorm:"not null" json:"quantity"`
	Cost      float64   `gorm:"type:decimal(10,2);not null" json:"cost"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName 指定Order模型的表名
func (Order) TableName() string {
	return "orders"
}

// TableName 指定OrderItem模型的表名
func (OrderItem) TableName() string {
	return "order_items"
}
