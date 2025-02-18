package mysql

import (
	"time"

	"gorm.io/gorm"
)

// CreateUser creates a new user
func CreateUser(user *User) error {
	return DB.Create(user).Error
}

// GetUserByEmail retrieves user by email
func GetUserByEmail(email string) (*User, error) {
	var user User
	err := DB.Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

// GetUserByID retrieves user by ID
func GetUserByID(id int64) (*User, error) {
	var user User
	err := DB.First(&user, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

// UpdateUser updates user information
func UpdateUser(user *User) error {
	return DB.Save(user).Error
}

// UpdateUserLoginTime updates user's last login time
func UpdateUserLoginTime(userID int64) error {
	return DB.Model(&User{}).Where("id = ?", userID).
		Update("last_login_at", time.Now()).Error
}

// CheckEmailExists checks if email is already in use
func CheckEmailExists(email string) (bool, error) {
	var count int64
	err := DB.Model(&User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}