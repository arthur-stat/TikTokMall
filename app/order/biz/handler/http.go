package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"TikTokMall/app/order/biz/service"
	"TikTokMall/app/order/kitex_gen/order"
)

type OrderHTTPHandler struct {
	svc service.OrderService
}

func NewOrderHTTPHandler() *OrderHTTPHandler {
	return &OrderHTTPHandler{
		svc: service.NewOrderService(),
	}
}

// PlaceOrder handles HTTP request for placing order
func (h *OrderHTTPHandler) PlaceOrder(c context.Context, ctx *app.RequestContext) {
	var req order.PlaceOrderReq
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	resp, err := h.svc.PlaceOrder(c, &req)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, resp)
}

// ListOrder handles HTTP request for listing orders
func (h *OrderHTTPHandler) ListOrder(c context.Context, ctx *app.RequestContext) {
	var req order.ListOrderReq
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	resp, err := h.svc.ListOrder(c, &req)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, resp)
}

// MarkOrderPaid handles HTTP request for marking order as paid
func (h *OrderHTTPHandler) MarkOrderPaid(c context.Context, ctx *app.RequestContext) {
	var req order.MarkOrderPaidReq
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	resp, err := h.svc.MarkOrderPaid(c, &req)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, resp)
}
