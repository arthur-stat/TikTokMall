// Code generated by MockGen. DO NOT EDIT.
package mock

import (
	"context"
	"reflect"

	"github.com/golang/mock/gomock"

	"TikTokMall/app/auth/kitex_gen/auth"
)

// MockAuthService is a mock of AuthService interface
type MockAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceMockRecorder
}

// MockAuthServiceMockRecorder is the mock recorder for MockAuthService
type MockAuthServiceMockRecorder struct {
	mock *MockAuthService
}

// NewMockAuthService creates a new mock instance
func NewMockAuthService(ctrl *gomock.Controller) *MockAuthService {
	mock := &MockAuthService{ctrl: ctrl}
	mock.recorder = &MockAuthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAuthService) EXPECT() *MockAuthServiceMockRecorder {
	return m.recorder
}

// Register mocks base method
func (m *MockAuthService) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, req)
	ret0, _ := ret[0].(*auth.RegisterResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register
func (mr *MockAuthServiceMockRecorder) Register(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockAuthService)(nil).Register), ctx, req)
}

// Login mocks base method
func (m *MockAuthService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, req)
	ret0, _ := ret[0].(*auth.LoginResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login
func (mr *MockAuthServiceMockRecorder) Login(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthService)(nil).Login), ctx, req)
}

// RefreshToken mocks base method
func (m *MockAuthService) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshToken", ctx, req)
	ret0, _ := ret[0].(*auth.RefreshTokenResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshToken indicates an expected call of RefreshToken
func (mr *MockAuthServiceMockRecorder) RefreshToken(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshToken", reflect.TypeOf((*MockAuthService)(nil).RefreshToken), ctx, req)
}

// ValidateToken mocks base method
func (m *MockAuthService) ValidateToken(ctx context.Context, req *auth.ValidateTokenRequest) (*auth.ValidateTokenResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateToken", ctx, req)
	ret0, _ := ret[0].(*auth.ValidateTokenResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateToken indicates an expected call of ValidateToken
func (mr *MockAuthServiceMockRecorder) ValidateToken(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateToken", reflect.TypeOf((*MockAuthService)(nil).ValidateToken), ctx, req)
}

// Logout mocks base method
func (m *MockAuthService) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout", ctx, req)
	ret0, _ := ret[0].(*auth.LogoutResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Logout indicates an expected call of Logout
func (mr *MockAuthServiceMockRecorder) Logout(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockAuthService)(nil).Logout), ctx, req)
}

// ... 其他方法的 mock 实现 