package model

import (
	"time"
)

// Payments 模型
type Payments struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement"`
	OrderID       string    `gorm:"column:order_id;uniqueIndex;not null"`
	UserID        int64     `gorm:"column:user_id;index;not null"`
	Amount        float32   `gorm:"column:amount;type:decimal(10,2);not null"`
	Status        int8      `gorm:"column:status;default:1;not null"`
	PaymentMethod string    `gorm:"column:payment_method;size:32;not null"`
	TransactionID string    `gorm:"column:transaction_id;uniqueIndex;size:128"`
	RefundID      string    `gorm:"column:transaction_id;uniqueIndex;size:128"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Payments) TableName() string {
	return "payments"
}
