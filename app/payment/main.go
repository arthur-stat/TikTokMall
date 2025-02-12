package main

import (
	"log"

	"TikTokMall/app/payment/handler"
	"TikTokMall/app/payment/kitex_gen/payment/paymentservice"
)

func main() {
	// 创建支付服务的 handler 实例
	svr := paymentservice.NewServer(handler.NewPaymentServiceImpl())

	// 启动服务
	err := svr.Run()
	if err != nil {
		log.Println("Payment service run error:", err.Error())
	}
}
