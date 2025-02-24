package hertz

import (
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hertz-contrib/registry/consul"
)

// NewConsulRegister 创建 Consul 注册器
func NewConsulRegister(addr string) (registry.Registry, error) {
	// 创建 Consul 客户端配置
	cfg := consulapi.DefaultConfig()
	cfg.Address = addr

	// 创建 Consul 客户端
	cli, err := consulapi.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	// 创建并返回注册器
	r := consul.NewConsulRegister(cli)
	return r, nil
}
