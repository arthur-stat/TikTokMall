package handler

import (
	"TikTokMall/app/payment/biz/service"
	payment "TikTokMall/app/payment/kitex_gen/payment"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	kkerrors "github.com/cloudwego/kitex/pkg/kerrors"
	"net/http"
)

func (s *PaymentServiceImpl) Charge(ctx context.Context, req *payment.ChargeReq) (resp *payment.ChargeResp, err error) {
	if req == nil {
		return nil, kkerrors.NewBizStatusError(4004001, "Invalid request")
	}
	resp, err = service.NewChargeService(ctx).Run(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ChargeHandler 处理支付请求的 HTTP 路由
func ChargeHandler(c context.Context, ctx *app.RequestContext) {
	// 从请求中解析支付请求参数
	var req payment.ChargeReq
	if err := ctx.Bind(&req); err != nil {
		handleError(ctx, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	// 调用业务层的 Charge 方法处理请求
	paymentService := &PaymentServiceImpl{}
	resp, err := paymentService.Charge(c, &req)
	if err != nil {
		handleError(ctx, http.StatusInternalServerError, "Failed to process payment", err)
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, resp)
}
