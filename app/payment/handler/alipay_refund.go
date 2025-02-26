package handler

import (
	"TikTokMall/app/payment/biz/service"
	"TikTokMall/app/payment/kitex_gen/payment"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	kkerrors "github.com/cloudwego/kitex/pkg/kerrors"
	"net/http"
)

func (s *PaymentServiceImpl) AlipayRefund(ctx context.Context, req *payment.AlipayRefundReq) (resp *payment.AlipayRefundResp, err error) {
	if req == nil {
		return nil, kkerrors.NewBizStatusError(40040012, "Invalid request")
	}
	resp, err = service.NewAlipayRefundService(ctx).Run(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func AlipayRefundHandler(c context.Context, ctx *app.RequestContext) {
	// 从请求中解析支付请求参数
	var req payment.AlipayRefundReq
	if err := ctx.Bind(&req); err != nil {
		handleError(ctx, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	// 调用业务层的 AlipayCharge 方法处理请求
	paymentService := &PaymentServiceImpl{}
	resp, err := paymentService.AlipayRefund(c, &req)
	if err != nil {
		handleError(ctx, http.StatusInternalServerError, "Failed to process payment", err)
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, resp)
}
