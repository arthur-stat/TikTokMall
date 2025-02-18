package service

import (
	"context"
	"errors"

	"TikTokMall/app/user/biz/dal/mysql"
	user "TikTokMall/app/user/kitex_gen/user"
	"golang.org/x/crypto/bcrypt"
)

type RegisterService struct {
	ctx context.Context
}

func NewRegisterService(ctx context.Context) *RegisterService {
	return &RegisterService{ctx: ctx}
}

func (s *RegisterService) Run(req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	// Validate passwords match
	if req.Password != req.ConfirmPassword {
		return nil, errors.New("passwords do not match")
	}

	// Check if email already exists
	exists, err := mysql.CheckEmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &mysql.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Status:   1, // Active
	}

	if err := mysql.CreateUser(user); err != nil {
		return nil, err
	}

	return &user.RegisterResp{
		UserId: int32(user.ID),
	}, nil
}