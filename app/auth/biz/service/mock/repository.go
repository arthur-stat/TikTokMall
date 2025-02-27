package mock

import (
	"github.com/stretchr/testify/mock"

	"TikTokMall/app/auth/biz/dal/mysql"
)

// MockAuthRepository 是 AuthRepository 的 mock 实现
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) CreateUser(user *mysql.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockAuthRepository) GetUserByUsername(username string) (*mysql.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mysql.User), args.Error(1)
}

func (m *MockAuthRepository) GetUserByID(id int64) (*mysql.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mysql.User), args.Error(1)
}

func (m *MockAuthRepository) CreateToken(token *mysql.Token) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockAuthRepository) GetTokenByToken(token string) (*mysql.Token, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mysql.Token), args.Error(1)
}

func (m *MockAuthRepository) DeleteToken(token string) error {
	args := m.Called(token)
	return args.Error(0)
}
