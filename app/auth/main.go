package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/cors"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"TikTokMall/app/auth/biz/dal/mysql"
	"TikTokMall/app/auth/biz/dal/redis"
	"TikTokMall/app/auth/biz/handler"
	"TikTokMall/app/auth/biz/registry"
	"TikTokMall/app/auth/biz/utils"
	"TikTokMall/app/auth/conf"
	"TikTokMall/app/auth/tracer"
)

func main() {
	// 1. 初始化配置
	conf.Init()
	config := conf.GetConfig()

	// 2. 初始化追踪
	tracer, closer, err := tracer.InitJaeger()
	if err != nil {
		hlog.Fatalf("初始化Jaeger失败: %v", err)
	}
	defer closer.Close()

	// 3. 初始化Prometheus
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Prometheus.Port), nil); err != nil {
			hlog.Fatalf("启动Prometheus metrics服务失败: %v", err)
		}
	}()

	// 初始化数据库连接
	if err := initDeps(); err != nil {
		hlog.Fatalf("init dependencies failed: %v", err)
	}

	// 创建 Consul 注册器
	r, err := consul.NewConsulRegister("localhost:8500")
	if err != nil {
		hlog.Fatalf("create consul register failed: %v", err)
	}

	// 创建HTTP服务器
	h := server.Default(
		server.WithHostPorts(":8000"),
		server.WithRegistry(r, &registry.Info{
			ServiceName: "auth",
			Addr:        utils.NewNetAddr("tcp", "localhost:8000"),
			Weight:      10,
			Tags:        []string{"auth", "v1"},
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
	authHandler := handler.NewAuthHandler()

	// 注册路由
	v1 := h.Group("/v1/auth")
	{
		v1.POST("/register", authHandler.Register)
		v1.POST("/login", authHandler.Login)
		v1.POST("/refresh", authHandler.RefreshToken)
		v1.POST("/logout", authHandler.Logout)
		v1.POST("/validate", authHandler.ValidateToken)
	}

	// 启动服务器
	if err := h.Run(); err != nil {
		hlog.Fatalf("start server failed: %v", err)
	}
}

// initDeps 初始化依赖
func initDeps() error {
	// 初始化MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		getEnvOrDefault("MYSQL_USER", "root"),
		getEnvOrDefault("MYSQL_PASSWORD", "123456"),
		getEnvOrDefault("MYSQL_HOST", "localhost"),
		getEnvOrDefault("MYSQL_PORT", "3306"),
		getEnvOrDefault("MYSQL_DATABASE", "tiktok_mall"),
	)
	if err := mysql.Init(dsn); err != nil {
		return fmt.Errorf("init mysql failed: %v", err)
	}

	// 初始化Redis
	if err := redis.Init(
		getEnvOrDefault("REDIS_ADDR", "localhost:6379"),
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
