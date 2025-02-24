package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"TikTokMall/app/order/kitex_gen/order"
)

func TestOrderService_PlaceOrder(t *testing.T) {
	svc := NewOrderService()

	tests := []struct {
		name    string
		req     *order.PlaceOrderReq
		wantErr bool
	}{
		{
			name: "valid order",
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
						Item: &order.CartItem{
							ProductId: 1,
							Quantity:  2,
						},
						Cost: 100.0,
					},
				},
			},
			wantErr: false,
		},
		// 添加更多测试用例
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := svc.PlaceOrder(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, resp.Order.OrderId)
		})
	}
}
