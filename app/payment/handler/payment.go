package handler

import (
	"context"
	kkerrors "github.com/cloudwego/kitex/pkg/kerrors"

	"TikTokMall/app/payment/biz/service"
	payment "TikTokMall/app/payment/kitex_gen/payment"
)

// PaymentServiceImpl 实现了支付服务接口
type PaymentServiceImpl struct {
	svc *service.ChargeService // 业务逻辑层服务
}

// NewPaymentServiceImpl 创建一个新的支付服务实现实例
func NewPaymentServiceImpl() *PaymentServiceImpl {
	return &PaymentServiceImpl{
		svc: service.NewChargeService(context.Background()), // 初始化业务服务
	}
}

// Charge 处理支付请求
func (s *PaymentServiceImpl) Charge(ctx context.Context, req *payment.ChargeReq) (resp *payment.ChargeResp, err error) {
	// 校验请求参数
	if req == nil {
		// 如果请求为空，则返回业务错误：4004001
		return nil, kkerrors.NewBizStatusError(4004001, "Invalid request")
	}

	// 调用业务逻辑层处理请求
	resp, err = s.svc.Run(req)
	if err != nil {
		// 业务处理发生错误，返回错误
		return nil, err
	}

	// 返回响应数据
	return resp, nil
}
