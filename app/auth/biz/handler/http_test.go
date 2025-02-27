package handler

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"TikTokMall/app/auth/biz/service/mock"
	"TikTokMall/app/auth/kitex_gen/auth"
)

func TestAuthHTTPHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name       string
		req        *auth.RegisterRequest
		mockSetup  func() *mock.MockAuthService
		wantStatus int
		wantResp   *auth.RegisterResponse
	}{
		{
			name: "success",
			req: &auth.RegisterRequest{
				Username: "testuser",
				Password: "password123",
				Email:    "test@example.com",
				Phone:    "13800138000",
			},
			mockSetup: func() *mock.MockAuthService {
				mockService := mock.NewMockAuthService(ctrl)
				mockService.EXPECT().Register(
					gomock.Any(),
					gomock.Any(),
				).Return(&auth.RegisterResponse{
					Base: &auth.BaseResp{
						Code:    consts.StatusOK,
						Message: "success",
					},
					Data: &auth.RegisterData{
						UserId: 1,
						Token:  "test-token",
					},
				}, nil)
				return mockService
			},
			wantStatus: consts.StatusOK,
			wantResp: &auth.RegisterResponse{
				Base: &auth.BaseResp{
					Code:    consts.StatusOK,
					Message: "success",
				},
				Data: &auth.RegisterData{
					UserId: 1,
					Token:  "test-token",
				},
			},
		},
		{
			name: "invalid_request",
			req:  &auth.RegisterRequest{},
			mockSetup: func() *mock.MockAuthService {
				return mock.NewMockAuthService(ctrl)
			},
			wantStatus: consts.StatusBadRequest,
			wantResp: &auth.RegisterResponse{
				Base: &auth.BaseResp{
					Code:    consts.StatusBadRequest,
					Message: "username and password are required",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.mockSetup()
			handler := NewAuthHTTPHandler(mockService)

			ctx := context.Background()
			hertzCtx := app.NewContext(16)

			if tt.req != nil {
				reqBody, _ := json.Marshal(tt.req)

				t.Logf("Request body: %s", string(reqBody))

				hertzCtx.Request.SetBody(reqBody)
				hertzCtx.Request.Header.Set("Content-Type", "application/json")
				hertzCtx.Request.Header.SetMethod("POST")
				hertzCtx.Request.SetRequestURI("/v1/auth/register")

				t.Logf("Request method: %s", string(hertzCtx.Request.Method()))
				t.Logf("Request URI: %s", string(hertzCtx.Request.URI().Path()))
				t.Logf("Request Content-Type: %s", string(hertzCtx.Request.Header.ContentType()))
			}

			handler.Register(ctx, hertzCtx)

			t.Logf("Response status: %d", hertzCtx.Response.StatusCode())
			t.Logf("Response body: %s", string(hertzCtx.Response.Body()))

			assert.Equal(t, tt.wantStatus, hertzCtx.Response.StatusCode())
			if tt.wantResp != nil {
				var got auth.RegisterResponse
				err := json.Unmarshal(hertzCtx.Response.Body(), &got)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResp, &got)
			}
		})
	}
}

func TestAuthHTTPHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name       string
		req        *auth.LoginRequest
		mockSetup  func() *mock.MockAuthService
		wantStatus int
		wantResp   interface{}
	}{
		{
			name: "success",
			req: &auth.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			mockSetup: func() *mock.MockAuthService {
				mockService := mock.NewMockAuthService(ctrl)
				mockService.EXPECT().Login(gomock.Any(), gomock.Any()).Return(
					&auth.LoginResponse{
						Base: &auth.BaseResp{
							Code:    consts.StatusOK,
							Message: "success",
						},
						Data: &auth.LoginData{
							Token:        "test-token",
							RefreshToken: "test-refresh-token",
						},
					},
					nil,
				)
				return mockService
			},
			wantStatus: consts.StatusOK,
			wantResp: &auth.LoginResponse{
				Base: &auth.BaseResp{
					Code:    consts.StatusOK,
					Message: "success",
				},
				Data: &auth.LoginData{
					Token:        "test-token",
					RefreshToken: "test-refresh-token",
				},
			},
		},
		{
			name: "invalid_credentials",
			req: &auth.LoginRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			mockSetup: func() *mock.MockAuthService {
				mockService := mock.NewMockAuthService(ctrl)
				mockService.EXPECT().Login(gomock.Any(), gomock.Any()).Return(
					nil,
					auth.ErrInvalidCredentials,
				)
				return mockService
			},
			wantStatus: consts.StatusUnauthorized,
			wantResp: map[string]interface{}{
				"error": "invalid username or password",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mockService *mock.MockAuthService
			if tt.mockSetup != nil {
				mockService = tt.mockSetup()
			} else {
				mockService = mock.NewMockAuthService(ctrl)
			}
			handler := &AuthHTTPHandler{svc: mockService}

			ctx := context.Background()
			hertzCtx := app.NewContext(16)
			if tt.req != nil {
				reqBody, _ := json.Marshal(tt.req)
				hertzCtx.Request.SetBody(reqBody)
			}

			handler.Login(ctx, hertzCtx)

			assert.Equal(t, tt.wantStatus, hertzCtx.Response.StatusCode())
			if tt.wantResp != nil {
				if tt.wantStatus == consts.StatusOK {
					var got auth.LoginResponse
					err := json.Unmarshal(hertzCtx.Response.Body(), &got)
					assert.NoError(t, err)
					assert.Equal(t, tt.wantResp, &got)
				} else {
					var got map[string]interface{}
					err := json.Unmarshal(hertzCtx.Response.Body(), &got)
					assert.NoError(t, err)
					assert.Equal(t, tt.wantResp, got)
				}
			}
		})
	}
}

func TestAuthHTTPHandler_RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name       string
		req        *auth.RefreshTokenRequest
		mockSetup  func() *mock.MockAuthService
		wantStatus int
		wantResp   interface{}
	}{
		{
			name: "success",
			req: &auth.RefreshTokenRequest{
				RefreshToken: "old-refresh-token",
			},
			mockSetup: func() *mock.MockAuthService {
				mockService := mock.NewMockAuthService(ctrl)
				mockService.EXPECT().RefreshToken(gomock.Any(), gomock.Any()).Return(
					&auth.RefreshTokenResponse{
						Base: &auth.BaseResp{
							Code:    consts.StatusOK,
							Message: "success",
						},
						Data: &auth.RefreshTokenData{
							Token:        "new-token",
							RefreshToken: "new-refresh-token",
						},
					},
					nil,
				)
				return mockService
			},
			wantStatus: consts.StatusOK,
			wantResp: &auth.RefreshTokenResponse{
				Base: &auth.BaseResp{
					Code:    consts.StatusOK,
					Message: "success",
				},
				Data: &auth.RefreshTokenData{
					Token:        "new-token",
					RefreshToken: "new-refresh-token",
				},
			},
		},
		{
			name: "invalid_token",
			req: &auth.RefreshTokenRequest{
				RefreshToken: "invalid-token",
			},
			mockSetup: func() *mock.MockAuthService {
				mockService := mock.NewMockAuthService(ctrl)
				mockService.EXPECT().RefreshToken(gomock.Any(), gomock.Any()).Return(
					nil,
					auth.ErrInvalidToken,
				)
				return mockService
			},
			wantStatus: consts.StatusUnauthorized,
			wantResp: map[string]interface{}{
				"error": "invalid token",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mockService *mock.MockAuthService
			if tt.mockSetup != nil {
				mockService = tt.mockSetup()
			} else {
				mockService = mock.NewMockAuthService(ctrl)
			}
			handler := &AuthHTTPHandler{svc: mockService}

			ctx := context.Background()
			hertzCtx := app.NewContext(16)
			if tt.req != nil {
				reqBody, _ := json.Marshal(tt.req)
				hertzCtx.Request.SetBody(reqBody)
			}

			handler.RefreshToken(ctx, hertzCtx)

			assert.Equal(t, tt.wantStatus, hertzCtx.Response.StatusCode())
			if tt.wantResp != nil {
				if tt.wantStatus == consts.StatusOK {
					var got auth.RefreshTokenResponse
					err := json.Unmarshal(hertzCtx.Response.Body(), &got)
					assert.NoError(t, err)
					assert.Equal(t, tt.wantResp, &got)
				} else {
					var got map[string]interface{}
					err := json.Unmarshal(hertzCtx.Response.Body(), &got)
					assert.NoError(t, err)
					assert.Equal(t, tt.wantResp, got)
				}
			}
		})
	}
}

func TestAuthHTTPHandler_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name       string
		req        *auth.ValidateTokenRequest
		mockSetup  func() *mock.MockAuthService
		wantStatus int
		wantResp   interface{}
	}{
		{
			name: "success",
			req: &auth.ValidateTokenRequest{
				Token: "valid-token",
			},
			mockSetup: func() *mock.MockAuthService {
				mockService := mock.NewMockAuthService(ctrl)
				mockService.EXPECT().ValidateToken(gomock.Any(), gomock.Any()).Return(
					&auth.ValidateTokenResponse{
						Base: &auth.BaseResp{
							Code:    consts.StatusOK,
							Message: "success",
						},
						Data: &auth.ValidateTokenData{
							UserId:   1,
							Username: "testuser",
						},
					},
					nil,
				)
				return mockService
			},
			wantStatus: consts.StatusOK,
			wantResp: &auth.ValidateTokenResponse{
				Base: &auth.BaseResp{
					Code:    consts.StatusOK,
					Message: "success",
				},
				Data: &auth.ValidateTokenData{
					UserId:   1,
					Username: "testuser",
				},
			},
		},
		{
			name: "invalid_token",
			req: &auth.ValidateTokenRequest{
				Token: "invalid-token",
			},
			mockSetup: func() *mock.MockAuthService {
				mockService := mock.NewMockAuthService(ctrl)
				mockService.EXPECT().ValidateToken(gomock.Any(), gomock.Any()).Return(
					&auth.ValidateTokenResponse{
						Base: &auth.BaseResp{
							Code:    consts.StatusUnauthorized,
							Message: "invalid token",
						},
					},
					auth.ErrInvalidToken,
				)
				return mockService
			},
			wantStatus: consts.StatusUnauthorized,
			wantResp: map[string]interface{}{
				"error": "invalid token",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mockService *mock.MockAuthService
			if tt.mockSetup != nil {
				mockService = tt.mockSetup()
			} else {
				mockService = mock.NewMockAuthService(ctrl)
			}
			handler := &AuthHTTPHandler{svc: mockService}

			ctx := context.Background()
			hertzCtx := app.NewContext(16)
			if tt.req != nil {
				reqBody, _ := json.Marshal(tt.req)
				hertzCtx.Request.SetBody(reqBody)
			}

			handler.ValidateToken(ctx, hertzCtx)

			assert.Equal(t, tt.wantStatus, hertzCtx.Response.StatusCode())
			if tt.wantResp != nil {
				if tt.wantStatus == consts.StatusOK {
					var got auth.ValidateTokenResponse
					err := json.Unmarshal(hertzCtx.Response.Body(), &got)
					assert.NoError(t, err)
					assert.Equal(t, tt.wantResp, &got)
				} else {
					var got map[string]interface{}
					err := json.Unmarshal(hertzCtx.Response.Body(), &got)
					assert.NoError(t, err)
					assert.Equal(t, tt.wantResp, got)
				}
			}
		})
	}
}

func TestAuthHTTPHandler_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name       string
		req        *auth.LogoutRequest
		mockSetup  func() *mock.MockAuthService
		wantStatus int
		wantResp   interface{}
	}{
		{
			name: "success",
			req: &auth.LogoutRequest{
				Token: "valid-token",
			},
			mockSetup: func() *mock.MockAuthService {
				mockService := mock.NewMockAuthService(ctrl)
				mockService.EXPECT().Logout(gomock.Any(), gomock.Any()).Return(
					&auth.LogoutResponse{
						Base: &auth.BaseResp{
							Code:    consts.StatusOK,
							Message: "success",
						},
					},
					nil,
				)
				return mockService
			},
			wantStatus: consts.StatusOK,
			wantResp: &auth.LogoutResponse{
				Base: &auth.BaseResp{
					Code:    consts.StatusOK,
					Message: "success",
				},
			},
		},
		{
			name: "invalid_token",
			req: &auth.LogoutRequest{
				Token: "invalid-token",
			},
			mockSetup: func() *mock.MockAuthService {
				mockService := mock.NewMockAuthService(ctrl)
				mockService.EXPECT().Logout(gomock.Any(), gomock.Any()).Return(
					nil,
					auth.ErrInvalidToken,
				)
				return mockService
			},
			wantStatus: consts.StatusUnauthorized,
			wantResp: map[string]interface{}{
				"error": "invalid token",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mockService *mock.MockAuthService
			if tt.mockSetup != nil {
				mockService = tt.mockSetup()
			} else {
				mockService = mock.NewMockAuthService(ctrl)
			}
			handler := &AuthHTTPHandler{svc: mockService}

			ctx := context.Background()
			hertzCtx := app.NewContext(16)
			if tt.req != nil {
				reqBody, _ := json.Marshal(tt.req)
				hertzCtx.Request.SetBody(reqBody)
			}

			handler.Logout(ctx, hertzCtx)

			assert.Equal(t, tt.wantStatus, hertzCtx.Response.StatusCode())
			if tt.wantResp != nil {
				if tt.wantStatus == consts.StatusOK {
					var got auth.LogoutResponse
					err := json.Unmarshal(hertzCtx.Response.Body(), &got)
					assert.NoError(t, err)
					assert.Equal(t, tt.wantResp, &got)
				} else {
					var got map[string]interface{}
					err := json.Unmarshal(hertzCtx.Response.Body(), &got)
					assert.NoError(t, err)
					assert.Equal(t, tt.wantResp, got)
				}
			}
		})
	}
}

// 辅助函数用于运行 HTTP 测试
func runHTTPTest(t *testing.T, handlerFunc func(context.Context, *app.RequestContext), tests []struct {
	name       string
	req        interface{}
	mockSetup  func()
	wantStatus int
	wantResp   interface{}
}) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			ctx := context.Background()
			hertzCtx := app.NewContext(16)
			if tt.req != nil {
				reqBody, _ := json.Marshal(tt.req)
				hertzCtx.Request.SetBody(reqBody)
			}

			handlerFunc(ctx, hertzCtx)

			assert.Equal(t, tt.wantStatus, hertzCtx.Response.StatusCode())
			if tt.wantResp != nil {
				if tt.wantStatus == consts.StatusOK {
					got := reflect.New(reflect.TypeOf(tt.wantResp).Elem()).Interface()
					err := json.Unmarshal(hertzCtx.Response.Body(), got)
					assert.NoError(t, err)
					assert.Equal(t, tt.wantResp, got)
				} else {
					var got map[string]interface{}
					err := json.Unmarshal(hertzCtx.Response.Body(), &got)
					assert.NoError(t, err)
					assert.Equal(t, tt.wantResp, got)
				}
			}
		})
	}
}
