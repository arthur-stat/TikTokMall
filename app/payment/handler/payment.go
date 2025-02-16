package handler

import (
	"context"
	"fmt"
	kkerrors "github.com/cloudwego/kitex/pkg/kerrors"
	"net/http"

	"TikTokMall/app/payment/biz/service"
	payment "TikTokMall/app/payment/kitex_gen/payment"
	"github.com/cloudwego/hertz/pkg/app"
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

// ChargeHandler处理支付请求的 HTTP 路由
func ChargeHandler(c context.Context, ctx *app.RequestContext) {
	// 从请求中解析支付请求参数
	var req payment.ChargeReq
	if err := ctx.Bind(&req); err != nil {
		// 如果绑定失败，返回 400 错误
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
		return
	}

	// 调用业务层的 Charge 方法处理请求
	paymentService := NewPaymentServiceImpl()
	resp, err := paymentService.Charge(c, &req)
	if err != nil {
		// 如果发生错误，返回 500 错误
		ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to process payment: %v", err),
		})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, resp)
}
