package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"TikTokMall/app/order/biz/dal/mysql"
	"TikTokMall/app/order/biz/dal/redis"
	"TikTokMall/app/order/kitex_gen/cart"
	"TikTokMall/app/order/kitex_gen/order"
)

type mockOrderRepo struct {
	mock.Mock
}

func (m *mockOrderRepo) CreateOrder(ctx context.Context, order *mysql.Order, items []*mysql.OrderItem) error {
	args := m.Called(ctx, order, items)
	return args.Error(0)
}

func (m *mockOrderRepo) ListOrdersByUserID(ctx context.Context, userID uint32) ([]*mysql.Order, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*mysql.Order), args.Error(1)
}

func (m *mockOrderRepo) UpdateOrderStatus(ctx context.Context, orderID int64, status int8) error {
	args := m.Called(ctx, orderID, status)
	return args.Error(0)
}

func TestOrderService_PlaceOrder(t *testing.T) {
	// 初始化 Redis 客户端
	if err := redis.Init(); err != nil {
		t.Fatalf("初始化 Redis 失败: %v", err)
	}

	repo := new(mockOrderRepo)
	svc := NewOrderService(repo)

	tests := []struct {
		name    string
		req     *order.PlaceOrderReq
		setup   func()
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
						Item: &cart.CartItem{
							ProductId: 1,
							Quantity:  2,
						},
						Cost: 100.0,
					},
				},
			},
			setup: func() {
				repo.On("CreateOrder", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			resp, err := svc.PlaceOrder(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, resp.Order.OrderId)
			repo.AssertExpectations(t)
		})
	}
}
