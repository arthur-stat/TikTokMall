package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

// NewServerTLSConfig 创建服务端TLS配置
func NewServerTLSConfig(caCertPath, certPath, keyPath string) (*tls.Config, error) {
	// 加载CA证书
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, fmt.Errorf("读取CA证书失败: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("添加CA证书到证书池失败")
	}

	// 加载服务器证书和密钥
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("加载服务器证书和密钥失败: %w", err)
	}

	// 创建TLS配置
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS12,
	}

	return tlsConfig, nil
}

// NewClientTLSConfig 创建客户端TLS配置
func NewClientTLSConfig(caCertPath, certPath, keyPath string) (*tls.Config, error) {
	// 加载CA证书
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, fmt.Errorf("读取CA证书失败: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("添加CA证书到证书池失败")
	}

	// 加载客户端证书和密钥
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("加载客户端证书和密钥失败: %w", err)
	}

	// 创建TLS配置
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		MinVersion:   tls.VersionTLS12,
	}

	return tlsConfig, nil
}
