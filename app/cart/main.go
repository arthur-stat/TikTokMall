package main

import (
	"fmt"
	"os"

	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/cors"

	"TikTokMall/app/cart/biz/dal/mysql"
	"TikTokMall/app/cart/biz/dal/redis"
	"TikTokMall/app/cart/biz/handler"
	"TikTokMall/app/cart/biz/utils"
	"TikTokMall/app/cart/conf"
	"TikTokMall/app/cart/pkg/hertz"
	"TikTokMall/app/cart/pkg/tracer"
)

func main() {
	// 1. 初始化配置
	if err := conf.Init(); err != nil {
		hlog.Fatalf("初始化配置失败: %v", err)
	}

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
	h := server.Default(
		server.WithHostPorts(":8000"),
		server.WithRegistry(r, &registry.Info{
			ServiceName: "cart",
			Addr:        utils.NewNetAddr("tcp", "localhost:8000"),
			Weight:      10,
			Tags: map[string]string{
				"version": "v1",
				"service": "cart",
			},
		}),
	)

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
	cartHandler := handler.NewCartHandler()

	// 注册路由
	v1 := h.Group("/v1/cart")
	{
		v1.POST("/add", cartHandler.AddItem)
		v1.POST("/update", cartHandler.UpdateItem)
		v1.POST("/delete", cartHandler.DeleteItem)
		v1.GET("/list", cartHandler.GetCart)
		v1.POST("/clear", cartHandler.ClearCart)
	}

	// 注释掉 Prometheus 相关代码
	/*
		// 初始化 Prometheus
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Prometheus.Port), nil); err != nil {
				hlog.Fatalf("启动Prometheus metrics服务失败: %v", err)
			}
		}()
	*/

	// 启动服务器
	if err := h.Run(); err != nil {
		hlog.Fatalf("start server failed: %v", err)
	}
}

// initDeps 初始化依赖
func initDeps() error {
	// 初始化MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		getEnvOrDefault("MYSQL_USER", "tiktok"),
		getEnvOrDefault("MYSQL_PASSWORD", "tiktok123"),
		getEnvOrDefault("MYSQL_HOST", "localhost"),
		getEnvOrDefault("MYSQL_PORT", "3307"),
		getEnvOrDefault("MYSQL_DATABASE", "tiktok_mall"),
	)
	if err := mysql.Init(dsn); err != nil {
		return fmt.Errorf("init mysql failed: %v", err)
	}

	// 初始化Redis
	if err := redis.Init(
		getEnvOrDefault("REDIS_ADDR", "localhost:6380"),
		getEnvOrDefault("REDIS_PASSWORD", ""),
		0,
	); err != nil {
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
