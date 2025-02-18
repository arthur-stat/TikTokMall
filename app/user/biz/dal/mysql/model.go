package mysql

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Email       string    `gorm:"column:email;unique;not null"`
	Password    string    `gorm:"column:password;not null"`
	Nickname    string    `gorm:"column:nickname"`
	Avatar      string    `gorm:"column:avatar"`
	Status      int       `gorm:"column:status;default:1"`
	LastLoginAt time.Time `gorm:"column:last_login_at"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}