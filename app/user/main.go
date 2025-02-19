package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/hertz-contrib/cors"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
	prometheus "github.com/kitex-contrib/monitor-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"

	"TikTokMall/app/user/biz/dal"
	"TikTokMall/app/user/biz/handler"
	"TikTokMall/app/user/conf"
	"TikTokMall/app/user/kitex_gen/user/userservice"
)

func main() {
	// Initialize configuration
	conf.Init()
	config := conf.GetConf()

	// Initialize logger
	initLogger()

	// Initialize dependencies
	dal.Init()

	// Initialize tracer
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.Kitex.Service),
		provider.WithExportEndpoint(config.Jaeger.Endpoint),
	)
	defer p.Shutdown(context.Background())

	// Start Prometheus metrics server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Prometheus.Port), nil); err != nil {
			klog.Fatal("start prometheus server failed:", err)
		}
	}()

	// Create Consul registry
	r, err := consul.NewConsulRegister(config.Consul.Addr)
	if err != nil {
		klog.Fatal("create consul register failed:", err)
	}

	// Start HTTP server
	go func() {
		h := server.Default(
			server.WithHostPorts(config.HTTP.Address),
			server.WithRegistry(r, &registry.Info{
				ServiceName: config.HTTP.ServiceName,
				Addr:       utils.NewNetAddr("tcp", config.HTTP.Address),
				Weight:     10,
				Tags:       []string{"user", "http", "v1"},
			}),
		)

		// Add CORS middleware
		h.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))

		// Add recovery middleware
		h.Use(recovery.Recovery())

		// Initialize user handler
		userHandler := handler.NewUserHandler()

		// Register HTTP routes
		v1 := h.Group("/v1/user")
		{
			v1.POST("/register", userHandler.Register)
			v1.POST("/login", userHandler.Login)
		}

		klog.Info("HTTP server is starting on ", config.HTTP.Address)
		if err := h.Run(); err != nil {
			klog.Fatal("start http server failed:", err)
		}
	}()

	// Initialize RPC server options
	opts := []kitex.server.Option{
		server.WithServiceAddr(&net.TCPAddr{Port: config.RPC.Port}),
		server.WithRegistry(r),
		server.WithLimit(&limit.Option{
			MaxConnections: 2000,
			MaxQPS:        500,
		}),
		server.WithTracer(prometheus.NewServerTracer(":9091", "/metrics")),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: config.RPC.ServiceName,
		}),
	}

	// Create and start RPC server
	svr := userservice.NewServer(NewUserServiceImpl(), opts...)
	klog.Info("RPC server is starting on port ", config.RPC.Port)
	if err := svr.Run(); err != nil {
		klog.Fatal("start rpc server failed:", err)
	}
}