package service

import (
	"context"
)

// PaymentClient 支付客户端接口
type PaymentClient interface {
	ProcessPayment(ctx context.Context, orderID string, amount float64, userID uint32) (string, bool, string, error)
}
