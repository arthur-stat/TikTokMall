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
	// 1. 初始化配置
	if err := conf.Init(); err != nil {
		hlog.Fatalf("initialize configuration failed: %v", err)
	}

	// 2. 初始化 tracer：Jaeger 追踪
	tracer, closer, err := tracer.InitJaeger()
	if err != nil {
		hlog.Fatalf("initialize Jaeger failed: %v", err)
	}
	defer closer.Close()
	_ = tracer // 暂时不使用 tracer

	// Comment out
	/*
		// 3. 初始化Prometheus
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Prometheus.Port), nil); err != nil {
				hlog.Fatalf("start Prometheus metrics server failed: %v", err)
			}
		}()
	*/

	// 初始化数据库连接，具体实现见后文
	if err := initDeps(); err != nil {
		hlog.Fatalf("init dependencies failed: %v", err)
	}

	// 创建 Consul 注册器
	r, err := hertz.NewConsulRegister("localhost:8500")
	if err != nil {
		hlog.Fatalf("create consul register failed: %v", err)
	}

	// 创建服务器（server）
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

	// 添加CORS中间件
	h.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		MaxAge:       3600,
	}))

	// 添加恢复中间件
	h.Use(recovery.Recovery())

	// 创建处理器（handler）
	userHandler := handler.NewUserHandler()

	// Register routes 注册路由
	v1 := h.Group("/v1/user")
	{
		v1.POST("/register", userHandler.Register) // 注册账户，同auth服务
		v1.POST("/login", userHandler.Login)       // 账户登录，同auth服务
		v1.POST("/logout", userHandler.Logout)     // 账户登出，同auth服务
		v1.POST("/delete", userHandler.Delete)     // 删除账户
		v1.POST("/update", userHandler.Update)     // 更新账户
		v1.POST("/info", userHandler.Info)         // 获取账户身份信息
	}

	// 启动服务器
	if err := h.Run(); err != nil {
		hlog.Fatalf("start server failed: %v", err)
	}
}

// initDeps 初始化依赖：初始化数据库
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

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
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
