package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"TikTokMall/app/order/biz/service/mock"
	"TikTokMall/app/order/kitex_gen/cart"
	"TikTokMall/app/order/kitex_gen/order"
)

func TestOrderHTTPHandler_PlaceOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockOrderService(ctrl)
	handler := &OrderHTTPHandler{svc: mockService}

	tests := []struct {
		name       string
		req        *order.PlaceOrderReq
		mockSetup  func()
		wantStatus int
		wantResp   interface{}
	}{
		{
			name: "success",
			req: &order.PlaceOrderReq{
				UserId:       1,
				UserCurrency: "CNY",
				Email:        "test@example.com",
				Address: &order.Address{
					StreetAddress: "123 Main St",
					City:          "Beijing",
					State:         "Beijing",
					Country:       "China",
					ZipCode:       100000,
				},
				OrderItems: []*order.OrderItem{
					{
						Item: &cart.CartItem{
							ProductId: 1,
							Quantity:  2,
						},
						Cost: 100.0,
					},
				},
			},
			mockSetup: func() {
				mockService.EXPECT().PlaceOrder(gomock.Any(), gomock.Any()).Return(
					&order.PlaceOrderResp{
						Order: &order.OrderResult{
							OrderId: "test-order-id",
						},
					},
					nil,
				)
			},
			wantStatus: consts.StatusOK,
			wantResp: &order.PlaceOrderResp{
				Order: &order.OrderResult{
					OrderId: "test-order-id",
				},
			},
		},
		// 添加更多测试用例
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置 mock 期望
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			// 创建请求上下文
			ctx := context.Background()
			hertzCtx := app.NewContext(16)
			reqBody, _ := json.Marshal(tt.req)
			hertzCtx.Request.SetBody(reqBody)

			// 调用处理器
			handler.PlaceOrder(ctx, hertzCtx)

			// 验证响应
			assert.Equal(t, tt.wantStatus, hertzCtx.Response.StatusCode())
			if tt.wantResp != nil {
				var got order.PlaceOrderResp
				err := json.Unmarshal(hertzCtx.Response.Body(), &got)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResp, &got)
			}
		})
	}
}

func TestOrderHTTPHandler_ListOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockOrderService(ctrl)
	handler := &OrderHTTPHandler{svc: mockService}

	tests := []struct {
		name       string
		setupReq   func(*app.RequestContext)
		mockSetup  func()
		wantStatus int
		wantResp   interface{}
	}{
		{
			name: "invalid_request",
			setupReq: func(ctx *app.RequestContext) {
				// 设置为 GET 请求，但不设置任何参数
				ctx.Request.Header.SetMethod("GET")
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			mockSetup:  func() {},
			wantStatus: consts.StatusBadRequest,
			wantResp: map[string]interface{}{
				"error": "invalid user_id",
			},
		},
		{
			name: "success",
			setupReq: func(ctx *app.RequestContext) {
				// 设置为 GET 请求
				ctx.Request.Header.SetMethod("GET")
				ctx.Request.Header.Set("Content-Type", "application/json")

				// 设置查询字符串
				ctx.Request.SetRequestURI("/v1/order/list?user_id=1")
			},
			mockSetup: func() {
				// 验证传入的参数
				mockService.EXPECT().ListOrder(gomock.Any(), &order.ListOrderReq{
					UserId: 1,
				}).Return(
					&order.ListOrderResp{
						Orders: []*order.Order{
							{
								OrderId:      "test-order-1",
								UserId:       1,
								UserCurrency: "CNY",
								Email:        "test@example.com",
								Address: &order.Address{
									StreetAddress: "123 Main St",
									City:          "Beijing",
									State:         "Beijing",
									Country:       "China",
									ZipCode:       100000,
								},
								CreatedAt: 1234567890,
							},
						},
					},
					nil,
				).AnyTimes()
			},
			wantStatus: consts.StatusOK,
			wantResp: &order.ListOrderResp{
				Orders: []*order.Order{
					{
						OrderId:      "test-order-1",
						UserId:       1,
						UserCurrency: "CNY",
						Email:        "test@example.com",
						Address: &order.Address{
							StreetAddress: "123 Main St",
							City:          "Beijing",
							State:         "Beijing",
							Country:       "China",
							ZipCode:       100000,
						},
						CreatedAt: 1234567890,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			ctx := context.Background()
			hertzCtx := app.NewContext(16)
			tt.setupReq(hertzCtx)

			handler.ListOrder(ctx, hertzCtx)

			assert.Equal(t, tt.wantStatus, hertzCtx.Response.StatusCode())
			if tt.wantResp != nil {
				if tt.wantStatus == consts.StatusOK {
					var got order.ListOrderResp
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

func TestOrderHTTPHandler_MarkOrderPaid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockOrderService(ctrl)
	handler := &OrderHTTPHandler{svc: mockService}

	tests := []struct {
		name       string
		req        *order.MarkOrderPaidReq
		mockSetup  func()
		wantStatus int
		wantResp   interface{}
	}{
		{
			name: "success",
			req: &order.MarkOrderPaidReq{
				UserId:  1,
				OrderId: "test-order-1",
			},
			mockSetup: func() {
				mockService.EXPECT().MarkOrderPaid(gomock.Any(), gomock.Any()).Return(
					&order.MarkOrderPaidResp{},
					nil,
				)
			},
			wantStatus: consts.StatusOK,
			wantResp:   &order.MarkOrderPaidResp{},
		},
		{
			name: "order_not_found",
			req: &order.MarkOrderPaidReq{
				UserId:  1,
				OrderId: "non-existent",
			},
			mockSetup: func() {
				mockService.EXPECT().MarkOrderPaid(gomock.Any(), gomock.Any()).Return(
					nil,
					fmt.Errorf("order not found"),
				)
			},
			wantStatus: consts.StatusInternalServerError,
			wantResp: map[string]interface{}{
				"error": "order not found",
			},
		},
	}

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

			handler.MarkOrderPaid(ctx, hertzCtx)

			assert.Equal(t, tt.wantStatus, hertzCtx.Response.StatusCode())
			if tt.wantResp != nil {
				if tt.wantStatus == consts.StatusOK {
					var got order.MarkOrderPaidResp
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
