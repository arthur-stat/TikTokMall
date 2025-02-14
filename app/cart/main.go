package main

import (
	"log"
	"net"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	consul "github.com/kitex-contrib/registry-consul"

	"TikTokMall/app/cart/handler"
	"TikTokMall/app/cart/kitex_gen/cart/cartservice"
)

func main() {
	// 创建 Consul 注册器
	r, err := consul.NewConsulRegister("localhost:8500")
	if err != nil {
		log.Fatalf("create consul register failed: %v", err)
	}

	// 创建服务器
	addr, _ := net.ResolveTCPAddr("tcp", ":8002")
	svr := cartservice.NewServer(
		handler.NewCartServiceImpl(),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "cart",
			Tags:        []string{"cart", "v1"},
		}),
	)

	// 启动服务器
	if err := svr.Run(); err != nil {
		log.Fatalf("server start failed: %v", err)
	}
}
