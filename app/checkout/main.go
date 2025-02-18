package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/server"
	consul "github.com/kitex-contrib/registry-consul"

	"TikTokMall/app/checkout/biz/dal/mysql"
	"TikTokMall/app/checkout/conf"
	"TikTokMall/app/checkout/kitex_gen/checkout/checkoutservice"
	"TikTokMall/app/checkout/pkg/tracer"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// 1. 初始化配置
	conf.Init()
	config := conf.GetConfig()

	// 2. 初始化追踪
	tracer, closer, err := tracer.InitJaeger()
	if err != nil {
		log.Fatalf("初始化Jaeger失败: %v", err)
	}
	defer closer.Close()

	// 3. 初始化Prometheus
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Prometheus.Port), nil); err != nil {
			log.Fatalf("启动Prometheus metrics服务失败: %v", err)
		}
	}()

	// 4. 初始化数据库连接
	if err := mysql.Init(); err != nil {
		log.Fatalf("初始化MySQL失败: %v", err)
	}

	// 5. 创建服务注册器
	r, err := consul.NewConsulRegister(config.Registry.RegistryAddress[0])
	if err != nil {
		log.Fatalf("创建服务注册器失败: %v", err)
	}

	// 6. 创建 RPC 服务器
	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", config.Service.Port))
	opts := []server.Option{
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithRegistryInfo(&registry.Info{
			ServiceName: config.Service.Name,
			Tags:        []string{"v1"},
		}),
		server.WithHealthCheck(true),
	}

	svr := checkoutservice.NewServer(NewCheckoutServiceImpl(), opts...)

	// 7. 启动服务
	log.Printf("checkout service starting on port %d...", config.Service.Port)
	if err := svr.Run(); err != nil {
		log.Fatalf("服务运行失败: %v", err)
	}
}
