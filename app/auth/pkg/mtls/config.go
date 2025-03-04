package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

// CertConfig 包含证书配置信息
type CertConfig struct {
	CACertPath     string
	ServerCertPath string
	ServerKeyPath  string
	ClientCertPath string
	ClientKeyPath  string
}

// NewServerTLSConfig 创建服务器TLS配置
func NewServerTLSConfig(config CertConfig) (*tls.Config, error) {
	// 加载服务器证书
	serverCert, err := tls.LoadX509KeyPair(config.ServerCertPath, config.ServerKeyPath)
	if err != nil {
		return nil, fmt.Errorf("加载服务器证书失败: %w", err)
	}

	// 加载CA证书
	caCert, err := ioutil.ReadFile(config.CACertPath)
	if err != nil {
		return nil, fmt.Errorf("加载CA证书失败: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("解析CA证书失败")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS12,
	}, nil
}

// NewClientTLSConfig 创建客户端TLS配置
func NewClientTLSConfig(config CertConfig) (*tls.Config, error) {
	// 加载客户端证书
	clientCert, err := tls.LoadX509KeyPair(config.ClientCertPath, config.ClientKeyPath)
	if err != nil {
		return nil, fmt.Errorf("加载客户端证书失败: %w", err)
	}

	// 加载CA证书
	caCert, err := ioutil.ReadFile(config.CACertPath)
	if err != nil {
		return nil, fmt.Errorf("加载CA证书失败: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("解析CA证书失败")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
		MinVersion:   tls.VersionTLS12,
	}, nil
}

// NewHertzServerTLSConfig 创建Hertz服务器TLS配置
func NewHertzServerTLSConfig(config CertConfig) (*tls.Config, error) {
	// 复用服务器TLS配置
	return NewServerTLSConfig(config)
}
