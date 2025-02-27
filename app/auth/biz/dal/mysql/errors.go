package mysql

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrInvalidInput   = errors.New("invalid input")
)
