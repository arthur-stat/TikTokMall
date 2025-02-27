package handler

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"TikTokMall/app/cart/biz/service/mock"
	"TikTokMall/app/cart/kitex_gen/cart"
)

func TestCartHTTPHandler_AddItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name       string
		req        *cart.AddItemReq
		mockSetup  func() *mock.MockCartService
		wantStatus int
		wantResp   interface{}
	}{
		{
			name: "success",
			req: &cart.AddItemReq{
				UserId: 1,
				Item: &cart.CartItem{
					ProductId: 101,
					Quantity:  2,
				},
			},
			mockSetup: func() *mock.MockCartService {
				mockService := mock.NewMockCartService(ctrl)
				mockService.EXPECT().AddItem(
					gomock.Any(),
					gomock.Any(),
				).Return(&cart.AddItemResp{}, nil)
				return mockService
			},
			wantStatus: consts.StatusOK,
			wantResp:   &cart.AddItemResp{},
		},
		{
			name: "invalid_quantity",
			req: &cart.AddItemReq{
				UserId: 1,
				Item: &cart.CartItem{
					ProductId: 101,
					Quantity:  0,
				},
			},
			mockSetup: func() *mock.MockCartService {
				mockService := mock.NewMockCartService(ctrl)
				mockService.EXPECT().AddItem(
					gomock.Any(),
					gomock.Any(),
				).Return(nil, service.ErrInvalidQuantity)
				return mockService
			},
			wantStatus: consts.StatusBadRequest,
			wantResp: map[string]interface{}{
				"error": "quantity must be greater than 0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.mockSetup()
			handler := &CartHTTPHandler{svc: mockService}

			ctx := context.Background()
			hertzCtx := app.NewContext(16)

			// 设置请求
			if tt.req != nil {
				reqBody, _ := json.Marshal(tt.req)

				t.Logf("Request body: %s", string(reqBody))

				hertzCtx.Request.SetBody(reqBody)
				hertzCtx.Request.Header.Set("Content-Type", "application/json")
				hertzCtx.Request.Header.SetMethod("POST")
				hertzCtx.Request.SetRequestURI("/v1/cart/add_item")

				t.Logf("Request method: %s", string(hertzCtx.Request.Method()))
				t.Logf("Request URI: %s", string(hertzCtx.Request.URI().Path()))
			}

			// 调用处理器
			handler.AddItem(ctx, hertzCtx)

			t.Logf("Response status: %d", hertzCtx.Response.StatusCode())
			t.Logf("Response body: %s", string(hertzCtx.Response.Body()))

			// 验证响应
			assert.Equal(t, tt.wantStatus, hertzCtx.Response.StatusCode())
			if tt.wantResp != nil {
				if tt.wantStatus == consts.StatusOK {
					var got cart.AddItemResp
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

func TestCartHTTPHandler_GetCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name       string
		req        *cart.GetCartReq
		mockSetup  func() *mock.MockCartService
		wantStatus int
		wantResp   interface{}
	}{
		{
			name: "success",
			req: &cart.GetCartReq{
				UserId: 1,
			},
			mockSetup: func() *mock.MockCartService {
				mockService := mock.NewMockCartService(ctrl)
				mockService.EXPECT().GetCart(
					gomock.Any(),
					gomock.Any(),
				).Return(&cart.GetCartResp{
					Cart: &cart.Cart{
						UserId: 1,
						Items: []*cart.CartItem{
							{
								ProductId: 101,
								Quantity:  2,
							},
						},
					},
				}, nil)
				return mockService
			},
			wantStatus: consts.StatusOK,
			wantResp: &cart.GetCartResp{
				Cart: &cart.Cart{
					UserId: 1,
					Items: []*cart.CartItem{
						{
							ProductId: 101,
							Quantity:  2,
						},
					},
				},
			},
		},
		{
			name: "empty_cart",
			req: &cart.GetCartReq{
				UserId: 2,
			},
			mockSetup: func() *mock.MockCartService {
				mockService := mock.NewMockCartService(ctrl)
				mockService.EXPECT().GetCart(
					gomock.Any(),
					gomock.Any(),
				).Return(&cart.GetCartResp{
					Cart: &cart.Cart{
						UserId: 2,
						Items:  []*cart.CartItem{},
					},
				}, nil)
				return mockService
			},
			wantStatus: consts.StatusOK,
			wantResp: &cart.GetCartResp{
				Cart: &cart.Cart{
					UserId: 2,
					Items:  []*cart.CartItem{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.mockSetup()
			handler := &CartHTTPHandler{svc: mockService}

			ctx := context.Background()
			hertzCtx := app.NewContext(16)

			// 设置请求
			if tt.req != nil {
				reqBody, _ := json.Marshal(tt.req)
				hertzCtx.Request.SetBody(reqBody)
				hertzCtx.Request.Header.Set("Content-Type", "application/json")
				hertzCtx.Request.Header.SetMethod("GET")
				hertzCtx.Request.SetRequestURI("/v1/cart/get")
			}

			// 调用处理器
			handler.GetCart(ctx, hertzCtx)

			// 验证响应
			assert.Equal(t, tt.wantStatus, hertzCtx.Response.StatusCode())
			if tt.wantResp != nil {
				var got cart.GetCartResp
				err := json.Unmarshal(hertzCtx.Response.Body(), &got)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResp, &got)
			}
		})
	}
}

func TestCartHTTPHandler_EmptyCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name       string
		req        *cart.EmptyCartReq
		mockSetup  func() *mock.MockCartService
		wantStatus int
		wantResp   interface{}
	}{
		{
			name: "success",
			req: &cart.EmptyCartReq{
				UserId: 1,
			},
			mockSetup: func() *mock.MockCartService {
				mockService := mock.NewMockCartService(ctrl)
				mockService.EXPECT().EmptyCart(
					gomock.Any(),
					gomock.Any(),
				).Return(&cart.EmptyCartResp{}, nil)
				return mockService
			},
			wantStatus: consts.StatusOK,
			wantResp:   &cart.EmptyCartResp{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.mockSetup()
			handler := &CartHTTPHandler{svc: mockService}

			ctx := context.Background()
			hertzCtx := app.NewContext(16)

			// 设置请求
			if tt.req != nil {
				reqBody, _ := json.Marshal(tt.req)
				hertzCtx.Request.SetBody(reqBody)
				hertzCtx.Request.Header.Set("Content-Type", "application/json")
				hertzCtx.Request.Header.SetMethod("POST")
				hertzCtx.Request.SetRequestURI("/v1/cart/empty")
			}

			// 调用处理器
			handler.EmptyCart(ctx, hertzCtx)

			// 验证响应
			assert.Equal(t, tt.wantStatus, hertzCtx.Response.StatusCode())
			if tt.wantResp != nil {
				var got cart.EmptyCartResp
				err := json.Unmarshal(hertzCtx.Response.Body(), &got)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResp, &got)
			}
		})
	}
}
