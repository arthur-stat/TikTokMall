package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"TikTokMall/app/cart/biz/service"
	"TikTokMall/app/cart/kitex_gen/cart"
)

type CartHTTPHandler struct {
	Svc service.CartService
}

func NewCartHTTPHandler() *CartHTTPHandler {
	return &CartHTTPHandler{
		Svc: service.NewCartService(),
	}
}

// AddItem handles HTTP request for adding items to cart
func (h *CartHTTPHandler) AddItem(c context.Context, ctx *app.RequestContext) {
	var req cart.AddItemReq
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	resp, err := h.Svc.AddItem(c, &req)
	if err != nil {
		statusCode := consts.StatusInternalServerError
		if err == service.ErrInvalidQuantity {
			statusCode = consts.StatusBadRequest
		}
		ctx.JSON(statusCode, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, resp)
}

// GetCart handles HTTP request for getting cart
func (h *CartHTTPHandler) GetCart(c context.Context, ctx *app.RequestContext) {
	var req cart.GetCartReq
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	resp, err := h.Svc.GetCart(c, &req)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, resp)
}

// EmptyCart handles HTTP request for emptying cart
func (h *CartHTTPHandler) EmptyCart(c context.Context, ctx *app.RequestContext) {
	var req cart.EmptyCartReq
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	resp, err := h.Svc.EmptyCart(c, &req)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, resp)
}
