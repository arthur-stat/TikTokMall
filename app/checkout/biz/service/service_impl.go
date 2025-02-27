package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"

	"TikTokMall/app/checkout/biz/dal/mysql"
	"TikTokMall/app/checkout/kitex_gen/checkout"
	"TikTokMall/app/checkout/kitex_gen/payment"
	"TikTokMall/app/checkout/pkg/metrics"
	"TikTokMall/app/checkout/pkg/opentracing"
	"TikTokMall/app/checkout/pkg/prometheus"
)

type checkoutServiceImpl struct {
	paymentClient PaymentClient
}

func NewCheckoutService() CheckoutService {
	// 创建支付客户端适配器
	paymentClient, err := NewPaymentClientAdapter()
	if err != nil {
		klog.Fatalf("创建支付客户端失败: %v", err)
	}

	return &checkoutServiceImpl{
		paymentClient: paymentClient,
	}
}

// Run 实现结账流程
func (s *checkoutServiceImpl) Run(ctx context.Context, req *checkout.CheckoutReq) (*checkout.CheckoutResp, error) {
	// 开始计时
	timer := prometheus.NewTimer(metrics.CheckoutDuration.WithLabelValues("total"))
	defer timer.ObserveDuration()

	// 创建span
	span, ctx := opentracing.StartSpanFromContext(ctx, "Checkout.Run")
	defer span.Finish()

	// 记录请求信息
	span.SetTag("user_id", req.UserId)
	span.SetTag("email", req.Email)

	klog.Infof("收到结账请求: user_id=%d, name=%s %s, email=%s",
		req.UserId, req.Firstname, req.Lastname, req.Email)

	// 1. 验证输入参数
	if err := s.validateRequest(req); err != nil {
		metrics.CheckoutTotal.WithLabelValues("failed").Inc()
		span.SetTag("error", true)
		span.LogKV("error.message", err.Error())
		return nil, fmt.Errorf("参数验证失败: %w", err)
	}

	// 2. 创建订单
	order := &mysql.Order{
		OrderNo:   generateOrderNo(req.UserId),
		UserID:    int64(req.UserId),
		Status:    mysql.OrderStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ShippingAddress: sql.NullString{
			String: formatAddress(req.Address),
			Valid:  true,
		},
		UserInfo: sql.NullString{
			String: formatUserInfo(req.Firstname, req.Lastname, req.Email),
			Valid:  true,
		},
	}

	// 3. 处理支付
	transactionID, err := s.processPayment(ctx, order, req.CreditCard)
	if err != nil {
		klog.Errorf("支付处理失败: %v", err)
		return nil, fmt.Errorf("支付处理失败: %w", err)
	}

	// 更新订单支付信息
	order.Status = mysql.OrderStatusPaid
	order.PaymentTime = sql.NullTime{Time: time.Now(), Valid: true}
	order.TransactionID = sql.NullString{String: transactionID, Valid: true}

	// 4. 保存订单
	if err := mysql.DB.Create(order).Error; err != nil {
		klog.Errorf("保存订单失败: %v", err)
		return nil, fmt.Errorf("保存订单失败: %w", err)
	}

	klog.Infof("结账完成: order_id=%s, transaction_id=%s", order.OrderNo, transactionID)

	// 5. 返回结果
	resp := &checkout.CheckoutResp{
		OrderId:       order.OrderNo,
		TransactionId: transactionID,
	}

	// 成功计数
	metrics.CheckoutTotal.WithLabelValues("success").Inc()
	return resp, nil
}

// 辅助方法

func (s *checkoutServiceImpl) validateRequest(req *checkout.CheckoutReq) error {
	if req.UserId == 0 {
		return ErrInvalidInput
	}
	if req.Address == nil {
		return ErrAddressInvalid
	}
	if req.CreditCard == nil {
		return ErrCreditCardInvalid
	}
	if req.Firstname == "" || req.Lastname == "" || req.Email == "" {
		return ErrInvalidInput
	}
	return nil
}

func (s *checkoutServiceImpl) processPayment(ctx context.Context, order *mysql.Order, creditCard *payment.CreditCardInfo) (string, error) {
	// 在测试模式下直接返回模拟数据
	if os.Getenv("TESTING") == "1" {
		return "mock-transaction-id", nil
	}

	// 调用支付服务
	transactionID, success, errMsg, err := s.paymentClient.ProcessPayment(
		ctx,
		order.OrderNo,
		order.TotalAmount,
		uint32(order.UserID),
	)

	if err != nil {
		metrics.PaymentTotal.WithLabelValues("failed").Inc()
		return "", fmt.Errorf("调用支付服务失败: %w", err)
	}

	if !success {
		metrics.PaymentTotal.WithLabelValues("failed").Inc()
		return "", fmt.Errorf("支付失败: %s", errMsg)
	}

	// 记录支付成功
	metrics.PaymentTotal.WithLabelValues("success").Inc()
	return transactionID, nil
}

func formatAddress(addr *checkout.Address) string {
	if addr == nil {
		return ""
	}
	return fmt.Sprintf("%s, %s, %s, %s %s",
		addr.StreetAddress,
		addr.City,
		addr.State,
		addr.Country,
		addr.ZipCode,
	)
}

func formatUserInfo(firstname, lastname, email string) string {
	return fmt.Sprintf("%s %s <%s>", firstname, lastname, email)
}

func generateOrderNo(userID uint32) string {
	return fmt.Sprintf("ORD-%d-%d", time.Now().UnixNano(), userID)
}
