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

	"TikTokMall/app/auth/biz/dal/mysql"
	"TikTokMall/app/auth/biz/dal/redis"
	"TikTokMall/app/auth/biz/handler"
	"TikTokMall/app/auth/biz/utils"
	"TikTokMall/app/auth/conf"
	"TikTokMall/app/auth/pkg/hertz"
	"TikTokMall/app/auth/pkg/mtls"
	"TikTokMall/app/auth/pkg/tracer"
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

	// 注释掉这段代码
	// TODO: 存在版本兼容问题
	/*
		// 3. 初始化Prometheus
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Prometheus.Port), nil); err != nil {
				hlog.Fatalf("启动Prometheus metrics服务失败: %v", err)
			}
		}()
	*/

	// 初始化数据库连接
	if err := initDeps(); err != nil {
		hlog.Fatalf("init dependencies failed: %v", err)
	}

	// 创建 Consul 注册器
	r, err := hertz.NewConsulRegister("localhost:8500")
	if err != nil {
		hlog.Fatalf("create consul register failed: %v", err)
	}

	// 创建服务器
	var h *server.Hertz

	// 如果启用了TLS，使用TLS配置创建服务器
	if conf.GetConf().TLS.Enabled {
		// 创建TLS配置
		tlsConfig, err := mtls.NewHertzServerTLSConfig(mtls.CertConfig{
			CACertPath:     conf.GetConf().TLS.CACert,
			ServerCertPath: conf.GetConf().TLS.ServerCert,
			ServerKeyPath:  conf.GetConf().TLS.ServerKey,
		})
		if err != nil {
			hlog.Fatalf("创建TLS配置失败: %v", err)
		}

		// 使用TLS和标准网络库创建服务器
		h = server.New(
			server.WithHostPorts(":8000"),
			server.WithTLS(tlsConfig),
			server.WithTransport(standard.NewTransporter),
			server.WithRegistry(r, &registry.Info{
				ServiceName: "auth",
				Addr:        utils.NewNetAddr("tcp", "localhost:8000"),
				Weight:      10,
				Tags: map[string]string{
					"version": "v1",
					"service": "auth",
				},
			}),
		)

		hlog.Info("已启用mTLS安全通信")
	} else {
		// 不使用TLS创建服务器
		h = server.New(
			server.WithHostPorts(":8000"),
			server.WithRegistry(r, &registry.Info{
				ServiceName: "auth",
				Addr:        utils.NewNetAddr("tcp", "localhost:8000"),
				Weight:      10,
				Tags: map[string]string{
					"version": "v1",
					"service": "auth",
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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		getEnvOrDefault("MYSQL_USER", conf.GetConf().MySQL.User),
		getEnvOrDefault("MYSQL_PASSWORD", conf.GetConf().MySQL.Password),
		getEnvOrDefault("MYSQL_HOST", conf.GetConf().MySQL.Host),
		conf.GetConf().MySQL.Port,
		getEnvOrDefault("MYSQL_DATABASE", conf.GetConf().MySQL.DBName),
	)
	if err := mysql.Init(dsn); err != nil {
		return fmt.Errorf("init mysql failed: %v", err)
	}

	// 初始化Redis
	if err := redis.Init(
		getEnvOrDefault("REDIS_ADDR", conf.GetConf().Redis.Addr),
		getEnvOrDefault("REDIS_PASSWORD", conf.GetConf().Redis.Password),
		conf.GetConf().Redis.DB,
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
