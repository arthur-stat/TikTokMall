package service

import (
	checkout "TikTokMall/app/checkout/kitex_gen/checkout"
	"context"
)

// CheckoutService 定义结账服务接口
type CheckoutService interface {
	Run(ctx context.Context, req *checkout.CheckoutReq) (*checkout.CheckoutResp, error)
}
