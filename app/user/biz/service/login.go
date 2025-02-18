package service

import (
	"context"
	"errors"

	"TikTokMall/app/user/biz/dal/mysql"
	"TikTokMall/app/user/biz/dal/redis"
	user "TikTokMall/app/user/kitex_gen/user"
	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	ctx context.Context
}

func NewLoginService(ctx context.Context) *LoginService {
	return &LoginService{ctx: ctx}
}

func (s *LoginService) Run(req *user.LoginReq) (resp *user.LoginResp, err error) {
	// Get user by email
	u, err := mysql.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user not found")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid password")
	}

	// Update last login time
	if err := mysql.UpdateUserLoginTime(u.ID); err != nil {
		return nil, err
	}

	// Update cache
	if err := redis.SetUser(s.ctx, u); err != nil {
		// Log error but don't fail the login
		// TODO: Add proper logging
	}

	return &user.LoginResp{
		UserId: int32(u.ID),
	}, nil
}