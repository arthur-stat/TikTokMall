package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	payment "TikTokMall/app/payment/kitex_gen/payment"
)

// validRefundReq 返回一个合法的 RefundReq 对象，用于成功流程测试
func validRefundReq() *payment.RefundReq {
	return &payment.RefundReq{
		TransactionId: "valid-transaction-id",
		UserId:        1,
		OrderId:       1,
		Amount:        100.0,
	}
}

// invalidTransactionIdRefundReq 返回一个交易 ID 无效的 RefundReq
func invalidTransactionIdRefundReq() *payment.RefundReq {
	return &payment.RefundReq{
		TransactionId: "invalid-transaction-id",
		UserId:        1,
		OrderId:       1,
		Amount:        100.0,
	}
}

func TestRefundService_Run_Success(t *testing.T) {
	ctx := context.Background()

	// 先创建支付记录
	chargeSvc := NewChargeService(ctx)
	chargeReq := validChargeReq()
	chargeReq.OrderId = time.Now().UnixNano() // 生成唯一订单号
	chargeResp, err := chargeSvc.Run(chargeReq)
	assert.NoError(t, err)

	// 使用真实生成的 transactionId
	req := &payment.RefundReq{
		TransactionId: chargeResp.TransactionId,
		UserId:        chargeReq.UserId,
		OrderId:       chargeReq.OrderId,
		Amount:        chargeReq.Amount,
	}

	svc := NewRefundService(ctx)
	resp, err := svc.Run(req)
	assert.NoError(t, err, "合法请求不应返回错误")
	assert.NotNil(t, resp, "返回结果不应为 nil")
	if resp != nil {
		assert.NotEmpty(t, resp.RefundId, "RefundId 不应为空")
	}
}

// TestRefundService_Run_InvalidTransactionId 测试交易 ID 无效的情况
func TestRefundService_Run_InvalidTransactionId(t *testing.T) {
	ctx := context.Background()
	svc := NewRefundService(ctx)

	req := invalidTransactionIdRefundReq()

	resp, err := svc.Run(req)
	assert.Error(t, err, "交易 ID 无效应返回错误")
	assert.Nil(t, resp, "错误情况下返回结果应为 nil")
	assert.Contains(t, err.Error(), "4004003", "错误码应为 4004003")
}
