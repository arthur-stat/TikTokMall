package discovery

import (
	"fmt"

	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
)

// NewConsulResolver 创建一个基于 Consul 的服务发现解析器
func NewConsulResolver(consulAddr string) (client.Resolver, error) {
	r, err := consul.NewConsulResolver(consulAddr)
	if err != nil {
		return nil, fmt.Errorf("create consul resolver failed: %v", err)
	}
	return r, nil
}

// GetConsulClient 获取带有服务发现功能的客户端选项
func GetConsulClient(consulAddr string) ([]client.Option, error) {
	r, err := NewConsulResolver(consulAddr)
	if err != nil {
		return nil, err
	}

	return []client.Option{
		client.WithResolver(r),
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()), // 使用加权负载均衡
	}, nil
}
