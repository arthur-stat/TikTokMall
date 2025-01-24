package mysql

import (
	"time"

	"gorm.io/gorm"
)

// Token 令牌模型
type Token struct {
	ID           int64     `gorm:"primaryKey"`
	UserID       int64     `gorm:"index;not null"`
	Token        string    `gorm:"size:512;not null;index"`
	RefreshToken string    `gorm:"size:512;index"`
	ExpiredAt    time.Time `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

// CreateToken 创建Token记录
func CreateToken(token *Token) error {
	return DB.Create(token).Error
}

// GetTokenByToken 通过Token获取记录
func GetTokenByToken(token string) (*Token, error) {
	var t Token
	err := DB.Where("token = ?", token).First(&t).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &t, err
}

// GetTokenByRefreshToken 通过RefreshToken获取记录
func GetTokenByRefreshToken(refreshToken string) (*Token, error) {
	var t Token
	err := DB.Where("refresh_token = ?", refreshToken).First(&t).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &t, err
}

// DeleteToken 删除Token记录
func DeleteToken(token string) error {
	return DB.Where("token = ?", token).Delete(&Token{}).Error
}

// DeleteTokenByUserID 删除用户的所有Token
func DeleteTokenByUserID(userID int64) error {
	return DB.Where("user_id = ?", userID).Delete(&Token{}).Error
}

// UpdateToken 更新Token记录
func UpdateToken(token *Token) error {
	return DB.Save(token).Error
}

// CleanExpiredTokens 清理过期的Token
func CleanExpiredTokens() error {
	return DB.Where("expired_at < ?", time.Now()).Delete(&Token{}).Error
}

// GetValidTokenByUserID 获取用户的有效Token
func GetValidTokenByUserID(userID int64) (*Token, error) {
	var t Token
	err := DB.Where("user_id = ? AND expired_at > ?", userID, time.Now()).
		First(&t).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &t, err
}

// IsTokenValid 检查Token是否有效
func IsTokenValid(token string) (bool, error) {
	var count int64
	err := DB.Model(&Token{}).
		Where("token = ? AND expired_at > ?", token, time.Now()).
		Count(&count).Error
	return count > 0, err
}
