package mysql

import (
	"gorm.io/gorm"
)

// User 用户表结构（与auth服务共享表结构）
type User struct {
	gorm.Model
	ID       int64  `gorm:"primaryKey;autoIncrement"`
	Username string `gorm:"uniqueIndex;size:32;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"size:64"`
	Phone    string `gorm:"size:16"`
	Status   int    `gorm:"default:1"`
}

// CreateUser 创建用户
func CreateUser(user *User) error {
	return DB.Create(user).Error
}

// GetUserByUsername 根据用户名查询用户
func GetUserByUsername(username string) (*User, error) {
	var user User
	err := DB.Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

// DeleteUser 删除用户
func DeleteUser(userID int64) error {
	return DB.Delete(&User{}, userID).Error
}

// UpdateUser 更新用户信息
func UpdateUser(userID int64, updates map[string]interface{}) error {
	return DB.Model(&User{}).Where("id = ?", userID).Updates(updates).Error
}

// GetUserByID 根据ID查询用户
func GetUserByID(userID int64) (*User, error) {
	var user User
	err := DB.First(&user, userID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

// CheckUserExists 检查用户名唯一性
func CheckUserExists(username string) (bool, error) {
	var count int64
	err := DB.Model(&User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}
