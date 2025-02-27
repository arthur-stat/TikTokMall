package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"

	"TikTokMall/app/cart/biz/dal/mysql"
	"TikTokMall/app/cart/biz/dal/redis"
	"TikTokMall/app/cart/conf"
	"TikTokMall/app/cart/handler"
	"TikTokMall/app/cart/kitex_gen/cart/cartservice"
	"TikTokMall/app/cart/pkg/tracer"
)

func main() {
	// 初始化配置
	if err := conf.Init(); err != nil {
		klog.Fatalf("初始化配置失败: %v", err)
	}

	// 初始化链路追踪
	t, closer, err := tracer.InitJaeger()
	if err != nil {
		klog.Fatalf("初始化Jaeger失败: %v", err)
	}
	defer closer.Close()
	_ = t // 暂时不使用 tracer

	// 初始化数据库连接
	// 初始化 MySQL
	if err := mysql.Init(); err != nil {
		klog.Fatalf("MySQL初始化失败: %v", err)
	}

	// 初始化Redis
	if err := redis.Init(); err != nil {
		klog.Fatalf("Redis初始化失败: %v", err)
	}

	// 启动RPC服务
	addrStr := ":8888"
	config := conf.GetConfig()

	// 创建网络地址
	addr, err := net.ResolveTCPAddr("tcp", addrStr)
	if err != nil {
		klog.Fatalf("解析地址失败: %v", err)
	}

	// 使用handler中已有的CartServiceImpl，避免重复定义
	impl := handler.NewCartServiceImpl()

	svr := cartservice.NewServer(
		impl,
		server.WithServiceAddr(addr), // 使用正确的地址对象
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: config.Service.Name,
		}),
	)

	// 异步启动服务
	go func() {
		klog.Infof("Cart RPC服务启动于 %s", addrStr)
		if err := svr.Run(); err != nil {
			klog.Fatalf("服务启动失败: %v", err)
		}
	}()

	// 等待终止信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	klog.Info("正在关闭服务...")
}
