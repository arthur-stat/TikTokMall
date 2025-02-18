package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"

	"TikTokMall/app/checkout/biz/dal/mysql"
	"TikTokMall/app/checkout/kitex_gen/checkout"
	"TikTokMall/app/checkout/kitex_gen/payment"
	"TikTokMall/app/checkout/pkg/client"
	"TikTokMall/app/checkout/pkg/metrics"
	"TikTokMall/app/checkout/pkg/opentracing"
	"TikTokMall/app/checkout/pkg/prometheus"
)

type checkoutService struct {
	paymentClient payment.Client
}

func NewCheckoutService() CheckoutService {
	paymentClient, err := payment.NewClient(
		"payment",
		client.WithHostPorts("127.0.0.1:8081"),
		client.WithMuxConnection(1),
	)
	if err != nil {
		klog.Fatalf("创建支付服务客户端失败: %v", err)
	}

	return &checkoutService{
		paymentClient: paymentClient,
	}
}

// Run 实现结账流程
func (s *checkoutService) Run(ctx context.Context, req *checkout.CheckoutReq) (*checkout.CheckoutResp, error) {
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

func (s *checkoutService) validateRequest(req *checkout.CheckoutReq) error {
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

func (s *checkoutService) processPayment(ctx context.Context, order *mysql.Order, creditCard *payment.CreditCardInfo) (string, error) {
	// 创建支付请求
	paymentReq := &payment.ProcessPaymentRequest{
		OrderId:    order.OrderNo,
		Amount:     order.TotalAmount,
		UserId:     uint32(order.UserID),
		CreditCard: creditCard,
	}

	// 记录支付开始
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("payment.order_id", order.OrderNo)
	span.SetTag("payment.amount", order.TotalAmount)

	// 调用支付服务
	resp, err := s.paymentClient.ProcessPayment(ctx, paymentReq)
	if err != nil {
		metrics.PaymentTotal.WithLabelValues("failed").Inc()
		span.SetTag("error", true)
		span.LogKV("error.message", err.Error())
		return "", fmt.Errorf("调用支付服务失败: %w", err)
	}

	// 检查支付结果
	if !resp.Success {
		metrics.PaymentTotal.WithLabelValues("failed").Inc()
		return "", ErrPaymentFailed
	}

	// 记录支付成功
	metrics.PaymentTotal.WithLabelValues("success").Inc()
	span.SetTag("payment.transaction_id", resp.TransactionId)

	return resp.TransactionId, nil
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
