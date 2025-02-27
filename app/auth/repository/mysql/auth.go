package mysql

import (
	"TikTokMall/app/auth/biz/dal/mysql"
	"TikTokMall/app/auth/biz/service"
)

type AuthRepository struct{}

func NewAuthRepository() service.AuthRepository {
	return &AuthRepository{}
}

func (r *AuthRepository) CreateUser(user *mysql.User) error {
	return mysql.DB.Create(user).Error
}

func (r *AuthRepository) GetUserByUsername(username string) (*mysql.User, error) {
	var user mysql.User
	err := mysql.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) GetUserByID(id int64) (*mysql.User, error) {
	var user mysql.User
	err := mysql.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) CreateToken(token *mysql.Token) error {
	return mysql.DB.Create(token).Error
}

func (r *AuthRepository) GetTokenByToken(token string) (*mysql.Token, error) {
	var t mysql.Token
	err := mysql.DB.Where("token = ?", token).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *AuthRepository) DeleteToken(token string) error {
	return mysql.DB.Where("token = ?", token).Delete(&mysql.Token{}).Error
}
