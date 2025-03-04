package gonet

import (
	"crypto/tls"

	"github.com/cloudwego/kitex/pkg/remote"
	"github.com/cloudwego/kitex/pkg/remote/trans/gonet"
)

// NewSvrTransHandlerFactory 创建一个基于gonet的服务端传输处理工厂
// 这允许Kitex服务器使用标准库的网络功能，从而支持TLS
func NewSvrTransHandlerFactory() remote.ServerTransHandlerFactory {
	return gonet.NewSvrTransHandlerFactory()
}

// NewCliTransHandlerFactory 创建一个基于gonet的客户端传输处理工厂
// 这允许Kitex客户端使用标准库的网络功能，从而支持TLS
func NewCliTransHandlerFactory(tlsConfig *tls.Config) remote.ClientTransHandlerFactory {
	return gonet.NewCliTransHandlerFactory()
}
