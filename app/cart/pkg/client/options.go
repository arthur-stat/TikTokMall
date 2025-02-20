package client

import (
	"github.com/cloudwego/kitex/client"
)

// WithHostPorts 设置服务地址
func WithHostPorts(hostPorts string) client.Option {
	return client.WithHostPorts(hostPorts)
}

// WithMuxConnection 设置连接复用
func WithMuxConnection(muxConnection int) client.Option {
	return client.WithMuxConnection(muxConnection)
}
