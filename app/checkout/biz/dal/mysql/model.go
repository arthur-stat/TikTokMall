package mysql

import (
	"database/sql"
	"time"
)

// Cart 购物车模型
type Cart struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"not null;index" json:"user_id"`
	ProductID int64     `gorm:"not null;index" json:"product_id"`
	Quantity  int       `gorm:"not null;default:1" json:"quantity"`
	Selected  bool      `gorm:"not null;default:true" json:"selected"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

// Order 订单模型
type Order struct {
	ID              int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo         string         `gorm:"type:varchar(32);not null;uniqueIndex" json:"order_no"`
	UserID          int64          `gorm:"not null;index" json:"user_id"`
	UserInfo        sql.NullString `gorm:"type:json" json:"user_info"`
	TotalAmount     float64        `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status          int8           `gorm:"type:tinyint;not null;default:1" json:"status"`
	PaymentMethod   sql.NullInt8   `gorm:"type:tinyint" json:"payment_method"`
	PaymentTime     sql.NullTime   `json:"payment_time"`
	TransactionID   sql.NullString `gorm:"type:varchar(64)" json:"transaction_id"`
	ShippingAddress sql.NullString `gorm:"type:json" json:"shipping_address"`
	CreatedAt       time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP;ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

// OrderItem 订单项模型
type OrderItem struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID      int64     `gorm:"not null;index" json:"order_id"`
	ProductID    int64     `gorm:"not null" json:"product_id"`
	ProductName  string    `gorm:"type:varchar(128);not null" json:"product_name"`
	ProductImage string    `gorm:"type:varchar(256)" json:"product_image"`
	Quantity     int       `gorm:"not null" json:"quantity"`
	Price        float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	CreatedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// 订单状态常量
const (
	OrderStatusPending  int8 = 1 // 待支付
	OrderStatusPaid     int8 = 2 // 已支付
	OrderStatusCanceled int8 = 3 // 已取消
	OrderStatusComplete int8 = 4 // 已完成
)

// TableName 指定Cart模型的表名
func (Cart) TableName() string {
	return "carts"
}

// TableName 指定Order模型的表名
func (Order) TableName() string {
	return "orders"
}

// TableName 指定OrderItem模型的表名
func (OrderItem) TableName() string {
	return "order_items"
}
