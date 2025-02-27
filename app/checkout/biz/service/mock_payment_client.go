package service

import (
	"context"
)

// MockPaymentClient 模拟支付客户端
type MockPaymentClient struct{}

// ProcessPayment 处理支付 - 使用我们的自定义接口
func (m *MockPaymentClient) ProcessPayment(ctx context.Context, orderID string, amount float64, userID uint32) (string, bool, string, error) {
	// 一个简单的模拟实现
	if orderID != "" && amount > 0 && userID > 0 {
		return "mock-transaction-123", true, "", nil
	}
	return "", false, "支付处理失败", nil
}
