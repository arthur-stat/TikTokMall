package handler

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"TikTokMall/app/cart/biz/service"
	"TikTokMall/app/cart/kitex_gen/cart"
)

// CartHandler 购物车HTTP处理器
type CartHandler struct {
	cartService service.CartService
}

// NewCartHandler 创建购物车处理器
func NewCartHandler() *CartHandler {
	return &CartHandler{
		cartService: service.NewCartService(),
	}
}

// AddItem 添加商品到购物车
func (h *CartHandler) AddItem(ctx context.Context, c *app.RequestContext) {
	var req struct {
		UserID    int64  `json:"user_id"`
		ProductID int64  `json:"product_id"`
		Quantity  int32  `json:"quantity"`
		Name      string `json:"name"`
		Price     int64  `json:"price"`
		Image     string `json:"image"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 创建符合接口的请求对象
	addItemReq := &cart.AddItemReq{
		UserId: uint32(req.UserID),
		Item: &cart.CartItem{
			ProductId: uint32(req.ProductID),
			Quantity:  req.Quantity,
		},
	}

	resp, err := h.cartService.AddItem(ctx, addItemReq)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "添加商品失败: " + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "添加商品成功",
		"data":    resp,
	})
}

// GetCart 获取购物车
func (h *CartHandler) GetCart(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Query("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "用户ID参数错误",
		})
		return
	}

	// 创建符合接口的请求对象
	getCartReq := &cart.GetCartReq{
		UserId: uint32(userID),
	}

	resp, err := h.cartService.GetCart(ctx, getCartReq)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取购物车失败: " + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "获取购物车成功",
		"data":    resp,
	})
}

// EmptyCart 清空购物车
func (h *CartHandler) EmptyCart(ctx context.Context, c *app.RequestContext) {
	var req struct {
		UserID int64 `json:"user_id"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 创建符合接口的请求对象
	emptyCartReq := &cart.EmptyCartReq{
		UserId: uint32(req.UserID),
	}

	resp, err := h.cartService.EmptyCart(ctx, emptyCartReq)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "清空购物车失败: " + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "清空购物车成功",
		"data":    resp,
	})
}

// UpdateItem 更新购物车商品
func (h *CartHandler) UpdateItem(ctx context.Context, c *app.RequestContext) {
	var req struct {
		UserID    int64 `json:"user_id"`
		ProductID int64 `json:"product_id"`
		Quantity  int32 `json:"quantity"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 由于proto中没有定义UpdateItem，我们使用AddItem来实现更新功能
	addItemReq := &cart.AddItemReq{
		UserId: uint32(req.UserID),
		Item: &cart.CartItem{
			ProductId: uint32(req.ProductID),
			Quantity:  req.Quantity,
		},
	}

	resp, err := h.cartService.AddItem(ctx, addItemReq)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "更新商品失败: " + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "更新商品成功",
		"data":    resp,
	})
}

// RemoveItem 从购物车移除商品
func (h *CartHandler) RemoveItem(ctx context.Context, c *app.RequestContext) {
	var req struct {
		UserID    int64 `json:"user_id"`
		ProductID int64 `json:"product_id"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 由于proto中没有定义RemoveItem，我们使用AddItem来实现移除功能
	// 设置数量为0表示移除
	addItemReq := &cart.AddItemReq{
		UserId: uint32(req.UserID),
		Item: &cart.CartItem{
			ProductId: uint32(req.ProductID),
			Quantity:  0, // 数量为0表示移除
		},
	}

	resp, err := h.cartService.AddItem(ctx, addItemReq)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "移除商品失败: " + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "移除商品成功",
		"data":    resp,
	})
}
