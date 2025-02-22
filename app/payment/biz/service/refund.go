package service

import (
	"TikTokMall/app/payment/biz/dal/mysql"
	payment "TikTokMall/app/payment/kitex_gen/payment"
	"context"
	kkerrors "github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/google/uuid"
)

type RefundService struct {
	ctx context.Context
}

// NewRefundService new RefundService
func NewRefundService(ctx context.Context) *RefundService {
	return &RefundService{ctx: ctx}
}

// Run 处理退款请求
func (s *RefundService) Run(req *payment.RefundReq) (resp *payment.RefundResp, err error) {
	// 校验请求参数
	if req == nil {
		return nil, kkerrors.NewBizStatusError(4004002, "Invalid refund request")
	}

	thePayment, err := mysql.GetPaymentByTransactionID(mysql.DB, s.ctx, req.TransactionId)
	if err != nil {
		return nil, kkerrors.NewBizStatusError(4004003, err.Error())
	}
	if thePayment == nil {
		return nil, kkerrors.NewBizStatusError(4004004, "Payment not found")
	}
	if thePayment.Status != 1 {
		return nil, kkerrors.NewBizStatusError(4004005, "Payment status is not valid")
	}
	if thePayment.UserID != req.UserId {
		return nil, kkerrors.NewBizStatusError(4004006, "User ID does not match")
	}
	if thePayment.OrderID != req.OrderId {
		return nil, kkerrors.NewBizStatusError(4004007, "Order ID does not match")
	}
	if thePayment.Amount != req.Amount {
		return nil, kkerrors.NewBizStatusError(4004008, "Amount does not match")
	}

	// 生成退款 ID（暂未使用支付接口）
	refundId, err := uuid.NewRandom()
	if err != nil {
		return nil, kkerrors.NewBizStatusError(4005005, err.Error())
	}

	err = mysql.StatusToRefundStatus(mysql.DB, s.ctx, req.TransactionId, refundId.String())
	if err != nil {
		return nil, kkerrors.NewBizStatusError(4005006, err.Error())
	}

	return &payment.RefundResp{RefundId: refundId.String()}, nil
}
