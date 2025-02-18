package main

import (
	"context"
	"net"
	"time"

	"TikTokMall/app/user/biz/dal"
	"TikTokMall/app/user/conf"
	"TikTokMall/app/user/kitex_gen/user/userservice"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
	prometheus "github.com/kitex-contrib/monitor-prometheus"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	// Initialize configuration
	conf.Init()

	// Initialize logger
	initLogger()

	// Initialize data access layer (MySQL and Redis)
	dal.Init()

	// Initialize service options
	opts := initServiceOptions()

	// Create and run server
	svr := userservice.NewServer(NewUserServiceImpl(), opts...)

	if err := svr.Run(); err != nil {
		klog.Fatal("server stopped with error:", err)
	}
}

func initLogger() {
	// Configure rotating log file
	logger := &lumberjack.Logger{
		Filename:   conf.GetConf().Kitex.LogFileName,
		MaxSize:    conf.GetConf().Kitex.LogMaxSize,
		MaxBackups: conf.GetConf().Kitex.LogMaxBackups,
		MaxAge:     conf.GetConf().Kitex.LogMaxAge,
		Compress:   true,
	}

	klog.SetLevel(conf.LogLevel())
	klog.SetOutput(logger)
}

func initServiceOptions() []server.Option {
	var opts []server.Option

	// Service address
	addr, err := net.ResolveTCPAddr("tcp", conf.GetConf().Kitex.Address)
	if err != nil {
		klog.Fatal("failed to resolve address:", err)
	}
	opts = append(opts, server.WithServiceAddr(addr))

	// Service info
	opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: conf.GetConf().Kitex.Service,
	}))

	// Rate limiting
	opts = append(opts, server.WithLimit(&limit.Option{
		MaxConnections: 2000,
		MaxQPS:        500,
	}))

	// Consul service registration
	r, err := consul.NewConsulRegister(conf.GetConf().Consul.Addr)
	if err != nil {
		klog.Fatal("failed to create consul register:", err)
	}
	opts = append(opts, server.WithRegistry(r))

	// Prometheus metrics
	opts = append(opts, server.WithTracer(prometheus.NewServerTracer(":9091", "/metrics")))

	// OpenTelemetry tracing
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(conf.GetConf().Kitex.Service),
		provider.WithExportEndpoint(conf.GetConf().Jaeger.Endpoint),
	)
	defer p.Shutdown(context.Background())
	opts = append(opts, server.WithSuite(tracing.NewServerSuite()))

	// Graceful shutdown
	opts = append(opts, server.WithExitWaitTime(time.Second*3))

	return opts
}