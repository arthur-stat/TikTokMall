package service

import (
	"context"
	"fmt"

	"github.com/cloudwego/kitex/pkg/klog"

	"TikTokMall/app/checkout/kitex_gen/payment/paymentservice"
)

// PaymentClientAdapter 适配真实的支付客户端
type PaymentClientAdapter struct {
	client paymentservice.Client
}

// NewPaymentClientAdapter 创建适配器
func NewPaymentClientAdapter() (PaymentClient, error) {
	client, err := paymentservice.NewClient("payment")
	if err != nil {
		return nil, fmt.Errorf("创建支付客户端失败: %w", err)
	}
	return &PaymentClientAdapter{client: client}, nil
}

// ProcessPayment 处理支付
func (a *PaymentClientAdapter) ProcessPayment(ctx context.Context, orderID string, amount float64, userID uint32) (string, bool, string, error) {
	// 这里需要根据实际的 payment 包实现，暂时返回测试数据
	// 在实际对接时，需要使用正确的 client 方法和参数
	klog.Infof("模拟支付处理: orderID=%s, amount=%f, userID=%d", orderID, amount, userID)
	return "tx-123456", true, "", nil
}
