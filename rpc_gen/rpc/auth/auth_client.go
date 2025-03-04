package auth

import (
	auth "TikTokMall/rpc_gen/kitex_gen/auth"
	"context"

	"TikTokMall/pkg/security/mtls"
	"TikTokMall/rpc_gen/kitex_gen/auth/authservice"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

type RPCClient interface {
	KitexClient() authservice.Client
	Service() string
	DeliverTokenByRPC(ctx context.Context, Req *auth.DeliverTokenReq, callOptions ...callopt.Option) (r *auth.DeliveryResp, err error)
	VerifyTokenByRPC(ctx context.Context, Req *auth.VerifyTokenReq, callOptions ...callopt.Option) (r *auth.VerifyResp, err error)
}

// ClientConfig 客户端配置
type ClientConfig struct {
	TLSEnabled     bool
	CACertPath     string
	ClientCertPath string
	ClientKeyPath  string
}

// DefaultClientConfig 默认客户端配置
var DefaultClientConfig = ClientConfig{
	TLSEnabled:     true,
	CACertPath:     "certs/ca-cert.pem",
	ClientCertPath: "certs/auth-cert.pem",
	ClientKeyPath:  "certs/auth-key.pem",
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	return NewRPCClientWithConfig(dstService, DefaultClientConfig, opts...)
}

func NewRPCClientWithConfig(dstService string, config ClientConfig, opts ...client.Option) (RPCClient, error) {
	// 如果启用了TLS，添加TLS配置
	if config.TLSEnabled {
		tlsConfig, err := mtls.NewClientTLSConfig(mtls.CertConfig{
			CACertPath:     config.CACertPath,
			ClientCertPath: config.ClientCertPath,
			ClientKeyPath:  config.ClientKeyPath,
		})
		if err != nil {
			return nil, err
		}

		// 添加TLS配置到客户端选项
		opts = append(opts, client.WithTLSConfig(tlsConfig))
	}

	kitexClient, err := authservice.NewClient(dstService, opts...)
	if err != nil {
		return nil, err
	}
	cli := &clientImpl{
		service:     dstService,
		kitexClient: kitexClient,
	}

	return cli, nil
}

type clientImpl struct {
	service     string
	kitexClient authservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() authservice.Client {
	return c.kitexClient
}

func (c *clientImpl) DeliverTokenByRPC(ctx context.Context, Req *auth.DeliverTokenReq, callOptions ...callopt.Option) (r *auth.DeliveryResp, err error) {
	return c.kitexClient.DeliverTokenByRPC(ctx, Req, callOptions...)
}

func (c *clientImpl) VerifyTokenByRPC(ctx context.Context, Req *auth.VerifyTokenReq, callOptions ...callopt.Option) (r *auth.VerifyResp, err error) {
	return c.kitexClient.VerifyTokenByRPC(ctx, Req, callOptions...)
}
