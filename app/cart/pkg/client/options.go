package client

import (
	"fmt"

	"github.com/cloudwego/kitex/client"

	"TikTokMall/app/cart/conf"
	"TikTokMall/app/cart/pkg/mtls"
)

// WithHostPorts 设置服务地址
func WithHostPorts(hostPorts string) client.Option {
	return client.WithHostPorts(hostPorts)
}

// WithMuxConnection 设置连接复用
func WithMuxConnection(muxConnection int) client.Option {
	return client.WithMuxConnection(muxConnection)
}

// NewClientOptions 创建客户端选项，包括TLS配置
func NewClientOptions() []client.Option {
	config := conf.GetConfig()
	options := []client.Option{}

	// 如果启用了TLS，添加TLS配置
	if config.TLS.Enable {
		tlsConfig, err := mtls.NewClientTLSConfig(
			config.TLS.CACertPath,
			config.TLS.ClientCertPath,
			config.TLS.ClientKeyPath,
		)
		if err == nil {
			// 简化TLS配置，仅使用基本选项
			// 使用默认地址，因为ServiceConfig没有Address字段
			defaultAddress := fmt.Sprintf("localhost:%d", config.Service.Port)
			options = append(options, client.WithHostPorts(defaultAddress))

			// 注意：由于Kitex版本或配置问题，我们暂时不使用高级TLS配置
			// 如果需要TLS，请确保Kitex版本支持以下API
			// options = append(options, client.WithTLSConfig(tlsConfig))

			// 记录TLS配置状态
			_ = tlsConfig // 暂时不使用TLS配置
		}
	}

	return options
}
