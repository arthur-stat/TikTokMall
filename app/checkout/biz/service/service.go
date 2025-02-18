package service

import (
	"context"
	"fmt"

	"TikTokMall/app/checkout/kitex_gen/checkout"
)

// CheckoutService 定义结账服务的接口
type CheckoutService interface {
	// Checkout 处理结账请求
	Run(ctx context.Context, req *checkout.CheckoutReq) (*checkout.CheckoutResp, error)
}

// 错误定义
var (
	ErrInvalidInput      = fmt.Errorf("无效的输入参数")
	ErrPaymentFailed     = fmt.Errorf("支付处理失败")
	ErrAddressInvalid    = fmt.Errorf("地址信息无效")
	ErrCreditCardInvalid = fmt.Errorf("信用卡信息无效")
	ErrOrderCreateFailed = fmt.Errorf("订单创建失败")
)
