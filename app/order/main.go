package main

import (
	"fmt"
	"os"

	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/network/standard"
	"github.com/hertz-contrib/cors"

	"TikTokMall/app/order/biz/dal/mysql"
	"TikTokMall/app/order/biz/dal/redis"
	"TikTokMall/app/order/biz/handler"
	"TikTokMall/app/order/biz/utils"
	"TikTokMall/app/order/conf"
	"TikTokMall/app/order/pkg/hertz"
	"TikTokMall/app/order/pkg/mtls"
	"TikTokMall/app/order/pkg/tracer"
)

func main() {
	// 1. 初始化配置
	conf.Init()
	config := conf.GetConf()

	// 2. 初始化追踪
	tracer, closer, err := tracer.InitJaeger()
	if err != nil {
		hlog.Fatalf("初始化Jaeger失败: %v", err)
	}
	defer closer.Close()
	_ = tracer // 暂时不使用 tracer

	// 初始化数据库连接
	if err := initDeps(); err != nil {
		hlog.Fatalf("init dependencies failed: %v", err)
	}

	// 创建 Consul 注册器
	r, err := hertz.NewConsulRegister("localhost:8501")
	if err != nil {
		hlog.Fatalf("create consul register failed: %v", err)
	}

	// 创建服务器
	var h *server.Hertz

	// 检查是否启用TLS
	if config.TLS.Enabled {
		// 创建TLS配置
		tlsConfig, err := mtls.NewKitexServerTLSConfig(
			config.TLS.CACertPath,
			config.TLS.CertPath,
			config.TLS.KeyPath,
		)
		if err != nil {
			hlog.Fatalf("创建TLS配置失败: %v", err)
		}

		// 使用TLS配置创建服务器
		h = server.Default(
			server.WithHostPorts(":8000"),
			server.WithRegistry(r, &registry.Info{
				ServiceName: "order",
				Addr:        utils.NewNetAddr("tcp", "localhost:8000"),
				Weight:      10,
				Tags: map[string]string{
					"version": "v1",
					"service": "order",
				},
			}),
			server.WithTLS(tlsConfig),
			server.WithTransport(standard.NewTransporter),
		)

		hlog.Infof("TLS已启用，服务器将使用HTTPS")
	} else {
		// 不使用TLS配置创建服务器
		h = server.Default(
			server.WithHostPorts(":8000"),
			server.WithRegistry(r, &registry.Info{
				ServiceName: "order",
				Addr:        utils.NewNetAddr("tcp", "localhost:8000"),
				Weight:      10,
				Tags: map[string]string{
					"version": "v1",
					"service": "order",
				},
			}),
		)
	}

	// 添加CORS中间件
	h.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		MaxAge:       3600,
	}))

	// 添加恢复中间件
	h.Use(recovery.Recovery())

	// 创建处理器
	orderHandler := handler.NewOrderHTTPHandler()

	// 注册路由
	v1 := h.Group("/v1/order")
	{
		v1.POST("/create", orderHandler.PlaceOrder)
		v1.GET("/list", orderHandler.ListOrder)
		v1.POST("/mark_paid", orderHandler.MarkOrderPaid)
	}

	// 启动服务器
	if err := h.Run(); err != nil {
		hlog.Fatalf("start server failed: %v", err)
	}
}

// initDeps 初始化依赖
func initDeps() error {
	// 初始化MySQL
	if err := mysql.Init(); err != nil {
		return fmt.Errorf("init mysql failed: %v", err)
	}

	// 初始化Redis
	if err := redis.Init(); err != nil {
		return fmt.Errorf("init redis failed: %v", err)
	}

	return nil
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
