package errno

import "errors"

var (
	ErrInvalidRequest         = errors.New("invalid request")
	ErrEmailRequired         = errors.New("email is required")
	ErrPasswordRequired      = errors.New("password is required")
	ErrConfirmPasswordRequired = errors.New("confirm password is required")
	ErrPasswordTooShort      = errors.New("password must be at least 8 characters")
	ErrPasswordsDoNotMatch   = errors.New("passwords do not match")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidPassword      = errors.New("invalid password")
)