package auth

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserBanned         = errors.New("user is banned")
	ErrTooManyAttempts    = errors.New("too many login attempts")
	// ... 其他错误定义
)
