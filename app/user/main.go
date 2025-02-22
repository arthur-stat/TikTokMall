package main

import (
	"fmt"
	"os"

	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/cors"

	"TikTokMall/app/auth/biz/dal/mysql"
	"TikTokMall/app/auth/biz/dal/redis"
	"TikTokMall/app/auth/biz/handler"
	"TikTokMall/app/auth/biz/utils"
	"TikTokMall/app/auth/conf"
	"TikTokMall/app/auth/pkg/hertz"
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
	h := server.Default(
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
		getEnvOrDefault("MYSQL_USER", "tiktok"),
		getEnvOrDefault("MYSQL_PASSWORD", "tiktok123"),
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

// // 旧入口代码
//func main() {
//	opts := kitexInit()
//
//	svr := userservice.NewServer(new(UserServiceImpl), opts...)
//
//	err := svr.Run()
//	if err != nil {
//		klog.Error(err.Error())
//	}
//}
//
//func kitexInit() (opts []server.Option) {
//	// address
//	addr, err := net.ResolveTCPAddr("tcp", conf.GetConf().Kitex.Address)
//	if err != nil {
//		panic(err)
//	}
//	opts = append(opts, server.WithServiceAddr(addr))
//
//	// service info
//	opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
//		ServiceName: conf.GetConf().Kitex.Service,
//	}))
//
//	// klog
//	logger := kitexlogrus.NewLogger()
//	klog.SetLogger(logger)
//	klog.SetLevel(conf.LogLevel())
//	asyncWriter := &zapcore.BufferedWriteSyncer{
//		WS: zapcore.AddSync(&lumberjack.Logger{
//			Filename:   conf.GetConf().Kitex.LogFileName,
//			MaxSize:    conf.GetConf().Kitex.LogMaxSize,
//			MaxBackups: conf.GetConf().Kitex.LogMaxBackups,
//			MaxAge:     conf.GetConf().Kitex.LogMaxAge,
//		}),
//		FlushInterval: time.Minute,
//	}
//	klog.SetOutput(asyncWriter)
//	server.RegisterShutdownHook(func() {
//		asyncWriter.Sync()
//	})
//	return
//}
