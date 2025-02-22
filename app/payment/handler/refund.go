package handler

import (
	"TikTokMall/app/payment/biz/service"
	payment "TikTokMall/app/payment/kitex_gen/payment"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	kkerrors "github.com/cloudwego/kitex/pkg/kerrors"
	"net/http"
)

// Refund 处理支付请求
func (s *PaymentServiceImpl) Refund(ctx context.Context, req *payment.RefundReq) (resp *payment.RefundResp, err error) {
	if req == nil {
		return nil, kkerrors.NewBizStatusError(4004001, "Invalid request")
	}
	resp, err = service.NewRefundService(ctx).Run(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RefundHandler 处理支付请求的 HTTP 路由
func RefundHandler(c context.Context, ctx *app.RequestContext) {
	// 从请求中解析支付请求参数
	var req payment.RefundReq
	if err := ctx.Bind(&req); err != nil {
		handleError(ctx, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	paymentService := &PaymentServiceImpl{}
	resp, err := paymentService.Refund(c, &req)
	if err != nil {
		handleError(ctx, http.StatusInternalServerError, "Failed to process refund", err)
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, resp)
}
