package mysql

import "errors"

var (
	// ErrRecordNotFound 记录未找到
	ErrRecordNotFound = errors.New("record not found")
	// ErrInvalidInput 无效的输入
	ErrInvalidInput = errors.New("invalid input")
	// ErrDuplicateKey 重复的键值
	ErrDuplicateKey = errors.New("duplicate key")
)
