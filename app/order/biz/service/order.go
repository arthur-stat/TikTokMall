package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"TikTokMall/app/order/biz/dal/mysql"
	"TikTokMall/app/order/biz/dal/redis"
	"TikTokMall/app/order/kitex_gen/order"
)

type orderService struct{}

func NewOrderService() OrderService {
	return &orderService{}
}

func (s *orderService) PlaceOrder(ctx context.Context, req *order.PlaceOrderReq) (*order.PlaceOrderResp, error) {
	// 1. 创建订单记录
	addr, err := json.Marshal(req.Address)
	if err != nil {
		return nil, fmt.Errorf("marshal address failed: %w", err)
	}

	newOrder := &mysql.Order{
		OrderNo:      generateOrderNo(req.UserId),
		UserID:       req.UserId,
		UserCurrency: req.UserCurrency,
		Email:        req.Email,
		ShippingAddress: sql.NullString{
			String: string(addr),
			Valid:  true,
		},
		Status:    mysql.OrderStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 2. 创建订单项
	var orderItems []*mysql.OrderItem
	var totalAmount float64
	for _, item := range req.OrderItems {
		orderItems = append(orderItems, &mysql.OrderItem{
			ProductID: item.Item.ProductId,
			Quantity:  item.Item.Quantity,
			Cost:      float64(item.Cost),
		})
		totalAmount += float64(item.Cost)
	}
	newOrder.TotalAmount = totalAmount

	// 3. 保存到数据库
	if err := mysql.CreateOrder(ctx, newOrder, orderItems); err != nil {
		return nil, fmt.Errorf("create order failed: %w", err)
	}

	// 4. 缓存订单信息
	if err := redis.CacheOrder(ctx, newOrder); err != nil {
		// 记录错误但不影响主流程
		fmt.Printf("cache order failed: %v\n", err)
	}

	return &order.PlaceOrderResp{
		Order: &order.OrderResult{
			OrderId: newOrder.OrderNo,
		},
	}, nil
}

func (s *orderService) ListOrder(ctx context.Context, req *order.ListOrderReq) (*order.ListOrderResp, error) {
	// 获取用户订单列表
	orders, err := mysql.ListOrdersByUserID(ctx, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("list orders failed: %w", err)
	}

	// 转换为响应格式
	var respOrders []*order.Order
	for _, o := range orders {
		var addr order.Address
		if o.ShippingAddress.Valid {
			if err := json.Unmarshal([]byte(o.ShippingAddress.String), &addr); err != nil {
				return nil, fmt.Errorf("unmarshal address failed: %w", err)
			}
		}

		respOrders = append(respOrders, &order.Order{
			OrderId:      o.OrderNo,
			UserId:       o.UserID,
			UserCurrency: o.UserCurrency,
			Address:      &addr,
			Email:        o.Email,
			CreatedAt:    int32(o.CreatedAt.Unix()),
		})
	}

	return &order.ListOrderResp{
		Orders: respOrders,
	}, nil
}

func (s *orderService) MarkOrderPaid(ctx context.Context, req *order.MarkOrderPaidReq) (*order.MarkOrderPaidResp, error) {
	// 获取订单ID
	orderID, err := strconv.ParseInt(req.OrderId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid order id: %w", err)
	}

	// 更新订单状态
	if err := mysql.UpdateOrderStatus(ctx, orderID, mysql.OrderStatusPaid); err != nil {
		return nil, fmt.Errorf("update order status failed: %w", err)
	}

	// 使缓存失效
	if err := redis.InvalidateOrderCache(ctx, orderID); err != nil {
		// 记录错误但不影响主流程
		fmt.Printf("invalidate order cache failed: %v\n", err)
	}

	return &order.MarkOrderPaidResp{}, nil
}

// 生成订单号
func generateOrderNo(userID uint32) string {
	return fmt.Sprintf("ORD-%d-%d", time.Now().UnixNano(), userID)
}
