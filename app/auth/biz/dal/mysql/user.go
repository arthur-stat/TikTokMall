package mysql

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Username  string    `gorm:"column:username;unique;not null"`
	Password  string    `gorm:"column:password;not null"`
	Email     string    `gorm:"column:email"`
	Phone     string    `gorm:"column:phone"`
	Status    int       `gorm:"column:status;default:1"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// CreateUser 创建用户
func CreateUser(user *User) error {
	return DB.Create(user).Error
}

// GetUserByUsername 通过用户名获取用户
func GetUserByUsername(username string) (*User, error) {
	var user User
	err := DB.Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

// GetUserByID 通过ID获取用户
func GetUserByID(id int64) (*User, error) {
	var user User
	err := DB.First(&user, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

// UpdateUser 更新用户信息
func UpdateUser(user *User) error {
	return DB.Save(user).Error
}

// UpdateUserPassword 更新用户密码
func UpdateUserPassword(userID int64, password string) error {
	return DB.Model(&User{}).Where("id = ?", userID).
		Update("password", password).Error
}

// UpdateUserStatus 更新用户状态
func UpdateUserStatus(userID int64, status int8) error {
	return DB.Model(&User{}).Where("id = ?", userID).
		Update("status", status).Error
}

// CheckUserExists 检查用户是否存在
func CheckUserExists(username string) (bool, error) {
	var count int64
	err := DB.Model(&User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// CheckEmailExists 检查邮箱是否已被使用
func CheckEmailExists(email string) (bool, error) {
	var count int64
	err := DB.Model(&User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// CheckPhoneExists 检查手机号是否已被使用
func CheckPhoneExists(phone string) (bool, error) {
	var count int64
	err := DB.Model(&User{}).Where("phone = ?", phone).Count(&count).Error
	return count > 0, err
}
