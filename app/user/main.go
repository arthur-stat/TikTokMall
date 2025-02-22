package main

import (
	"fmt"
	"os"

	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/cors"

	"TikTokMall/app/user/biz/dal/mysql"
	"TikTokMall/app/user/biz/dal/redis"
	"TikTokMall/app/user/biz/handler"
	"TikTokMall/app/user/biz/utils"
	"TikTokMall/app/user/conf"
	"TikTokMall/app/user/pkg/hertz"
	"TikTokMall/app/user/pkg/tracer"
)

func main() {
	// 1. Initialize configuration
	if err := conf.Init(); err != nil {
		hlog.Fatalf("initialize configuration failed: %v", err)
	}

	// 2. Initialize tracing
	tracer, closer, err := tracer.InitJaeger()
	if err != nil {
		hlog.Fatalf("initialize Jaeger failed: %v", err)
	}
	defer closer.Close()
	_ = tracer // temporarily not using tracer

	// Comment out this code
	/*
		// 3. Initialize Prometheus
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Prometheus.Port), nil); err != nil {
				hlog.Fatalf("start Prometheus metrics server failed: %v", err)
			}
		}()
	*/

	// Initialize database connections
	if err := initDeps(); err != nil {
		hlog.Fatalf("init dependencies failed: %v", err)
	}

	// Create Consul registry
	r, err := hertz.NewConsulRegister("localhost:8500")
	if err != nil {
		hlog.Fatalf("create consul register failed: %v", err)
	}

	// Create server
	h := server.Default(
		server.WithHostPorts(":8001"),
		server.WithRegistry(r, &registry.Info{
			ServiceName: "user",
			Addr:        utils.NewNetAddr("tcp", "localhost:8001"),
			Weight:      10,
			Tags: map[string]string{
				"version": "v1",
				"service": "user",
			},
		}),
	)

	// Add CORS middleware
	h.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		MaxAge:       3600,
	}))

	// Add recovery middleware
	h.Use(recovery.Recovery())

	// Create handler
	userHandler := handler.NewUserHandler()

	// Register routes
	v1 := h.Group("/v1/user")
	{
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)
		v1.POST("/refresh", userHandler.RefreshToken)
		v1.POST("/logout", userHandler.Logout)
		v1.POST("/validate", userHandler.ValidateToken)
	}

	// Start server
	if err := h.Run(); err != nil {
		hlog.Fatalf("start server failed: %v", err)
	}
}

// initDeps initializes dependencies
func initDeps() error {
	// Initialize MySQL
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

	// Initialize Redis
	if err := redis.Init(
		getEnvOrDefault("REDIS_ADDR", "localhost:6379"),
		getEnvOrDefault("REDIS_PASSWORD", ""),
		0,
	); err != nil {
		return fmt.Errorf("init redis failed: %v", err)
	}

	return nil
}

// getEnvOrDefault gets environment variable, returns default value if not exist
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// // 旧服务入口代码
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
