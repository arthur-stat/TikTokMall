package main

import (
	"TikTokMall/app/payment/biz/dal"
	"TikTokMall/app/payment/conf"
	"TikTokMall/app/payment/handler"
	"TikTokMall/app/payment/kitex_gen/payment/paymentservice"
	"context"
	hserver "github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	kserver "github.com/cloudwego/kitex/server"
	"github.com/hashicorp/consul/api"
	"github.com/joho/godotenv"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	// 初始化 Kitex 服务
	opts := kitexInit()
	// 获取环境变量
	_ = godotenv.Load()
	// 初始化 Redis 和 MySQL 客户端
	err := dal.Init()
	if err != nil {
		klog.Error("Failed to initialize Redis and MySQL: %v", err)
	}

	wg.Add(1)
	// 创建支付服务的 handler 实例
	svr := paymentservice.NewServer(new(handler.PaymentServiceImpl), opts...)
	// 启动服务
	go func() {
		defer wg.Done()
		err := svr.Run()
		if err != nil {
			klog.Error("Payment service run error:", err.Error())
			panic(err)
		}
	}()

	// 启动 HTTP 服务并注册服务
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := startHTTPServer()
		if err != nil {
			hlog.Error("Failed to start HTTP server: %v", err)
			panic(err)
		}
	}()

	wg.Wait()
}

// startHTTPServer 启动 HTTP 服务并注册服务
func startHTTPServer() error {
	h := hserver.Default(
		hserver.WithHostPorts(":8005"),
		hserver.WithKeepAlive(true),
	)

	v1 := h.Group("/payment")
	{
		v1.GET("/health", handler.HealthHandler)
		v1.POST("/charge", handler.ChargeHandler)
		v1.POST("/refund", handler.RefundHandler)
		v1.POST("/alipay/charge", handler.AlipayChargeHandler)
	}

	// 添加 OnRun 钩子，在服务启动后注册 Consul
	h.OnRun = append(h.OnRun, func(ctx context.Context) error {
		err := registerServiceToConsul("payment", "localhost", 8005, "http")
		if err != nil {
			hlog.Error("Failed to register service to Consul: %v", err)
			return err
		}
		hlog.Info("Service registered to Consul after server started")
		return nil
	})

	return h.Run()
}

func kitexInit() (opts []kserver.Option) {
	// address
	addr, err := net.ResolveTCPAddr("tcp", conf.GetConf().Kitex.Address)
	if err != nil {
		panic(err)
	}
	opts = append(opts, kserver.WithServiceAddr(addr))

	// service info
	opts = append(opts, kserver.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: conf.GetConf().Kitex.Service,
	}))

	// klog
	logger := kitexlogrus.NewLogger()
	klog.SetLogger(logger)
	klog.SetLevel(conf.LogLevel())
	asyncWriter := &zapcore.BufferedWriteSyncer{
		WS: zapcore.AddSync(&lumberjack.Logger{
			Filename:   conf.GetConf().Kitex.LogFileName,
			MaxSize:    conf.GetConf().Kitex.LogMaxSize,
			MaxBackups: conf.GetConf().Kitex.LogMaxBackups,
			MaxAge:     conf.GetConf().Kitex.LogMaxAge,
		}),
		FlushInterval: time.Minute,
	}
	klog.SetOutput(asyncWriter)
	kserver.RegisterShutdownHook(func() {
		asyncWriter.Sync()
	})
	return
}

// registerServiceToConsul 将服务注册到 Consul
func registerServiceToConsul(serviceName, host string, port int, protocol string) error {
	// 创建Consul客户端
	config := &api.Config{
		Address: "localhost:8500", // Consul地址
	}
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}

	// 注册服务到 Consul
	registration := &api.AgentServiceRegistration{
		ID:      serviceName,        // 服务ID
		Name:    serviceName,        // 服务名称
		Address: host,               // 服务地址
		Port:    port,               // 服务端口
		Tags:    []string{protocol}, // 服务协议类型
		Check: &api.AgentServiceCheck{ // 健康检查配置
			HTTP:     "http://localhost:8005/payment/health", // 健康检查路径
			Interval: "10s",                                  // 健康检查时间间隔
		},
	}

	// 注册服务到Consul
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	}

	hlog.Info("Service %s registered successfully ", serviceName)
	return nil
}
