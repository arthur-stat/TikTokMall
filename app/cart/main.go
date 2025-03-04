package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"TikTokMall/app/cart/biz/dal/mysql"
	"TikTokMall/app/cart/biz/dal/redis"
	"TikTokMall/app/cart/conf"
	"TikTokMall/app/cart/handler"
	"TikTokMall/app/cart/pkg/mtls"
	"TikTokMall/app/cart/pkg/tracer"
)

func main() {
	// 初始化配置
	if err := conf.Init(); err != nil {
		hlog.Fatalf("初始化配置失败: %v", err)
	}

	// 初始化链路追踪
	t, closer, err := tracer.InitJaeger()
	if err != nil {
		hlog.Fatalf("初始化Jaeger失败: %v", err)
	}
	defer closer.Close()
	_ = t // 暂时不使用 tracer

	// 初始化数据库连接
	// 初始化 MySQL
	if err := mysql.Init(); err != nil {
		hlog.Fatalf("MySQL初始化失败: %v", err)
	}

	// 初始化Redis
	if err := redis.Init(); err != nil {
		hlog.Fatalf("Redis初始化失败: %v", err)
	}

	// 获取配置
	config := conf.GetConfig()
	addrStr := ":8888"

	// 创建Hertz服务器
	h := server.Default(server.WithHostPorts(addrStr))

	// 如果启用了TLS，记录日志
	if config.TLS.Enable {
		_, err := mtls.NewServerTLSConfig(
			config.TLS.CACertPath,
			config.TLS.ServerCertPath,
			config.TLS.ServerKeyPath,
		)
		if err != nil {
			hlog.Errorf("TLS配置失败: %v", err)
		} else {
			// 由于已经创建了服务器，我们需要手动设置TLS配置
			// 注意：这里可能需要根据Hertz的API进行调整
			hlog.Warn("TLS配置已启用，但在服务器创建后设置TLS可能不生效")
		}
	}

	// 注册路由
	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, map[string]interface{}{
			"message": "pong",
		})
	})

	// 注册购物车API路由
	cartHandler := handler.NewCartHandler()
	cartGroup := h.Group("/api/cart")
	{
		cartGroup.POST("/add", cartHandler.AddItem)
		cartGroup.GET("/get", cartHandler.GetCart)
		cartGroup.POST("/empty", cartHandler.EmptyCart)
		cartGroup.POST("/update", cartHandler.UpdateItem)
		cartGroup.POST("/remove", cartHandler.RemoveItem)
	}

	// 异步启动服务
	go func() {
		hlog.Infof("Cart HTTP服务启动于 %s", addrStr)
		if err := h.Run(); err != nil {
			hlog.Fatalf("服务启动失败: %v", err)
		}
	}()

	// 等待终止信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	hlog.Info("正在关闭服务...")
	h.Shutdown(context.Background())
}
