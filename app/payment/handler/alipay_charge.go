package handler

import (
	"TikTokMall/app/payment/biz/service"
	payment "TikTokMall/app/payment/kitex_gen/payment"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	kkerrors "github.com/cloudwego/kitex/pkg/kerrors"
	"net/http"
)

func (s *PaymentServiceImpl) AlipayCharge(ctx context.Context, req *payment.AlipayChargeReq) (resp *payment.AlipayChargeResp, err error) {
	if req == nil {
		return nil, kkerrors.NewBizStatusError(4004001, "Invalid request")
	}
	resp, err = service.NewAlipayChargeService(ctx).Run(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// AlipayChargeHandler 处理支付请求的 HTTP 路由
func AlipayChargeHandler(c context.Context, ctx *app.RequestContext) {
	// 从请求中解析支付请求参数
	var req payment.AlipayChargeReq
	if err := ctx.Bind(&req); err != nil {
		handleError(ctx, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	// 调用业务层的 AlipayCharge 方法处理请求
	paymentService := &PaymentServiceImpl{}
	resp, err := paymentService.AlipayCharge(c, &req)
	if err != nil {
		handleError(ctx, http.StatusInternalServerError, "Failed to process payment", err)
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, resp)
}
