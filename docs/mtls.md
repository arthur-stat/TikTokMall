# TikTokMall微服务mTLS实现指南

本文档介绍了如何在TikTokMall微服务架构中实现mTLS（双向TLS）认证，以确保服务间通信的安全性。

## 什么是mTLS？

mTLS（mutual TLS，双向TLS）是一种安全通信协议，它要求通信双方都提供证书进行身份验证。与传统的单向TLS不同，mTLS不仅要求服务器向客户端提供证书，还要求客户端向服务器提供证书。这样可以确保：

1. 服务器的身份是可信的（客户端验证服务器）
2. 客户端的身份是可信的（服务器验证客户端）
3. 通信内容是加密的，防止中间人攻击

## 实现步骤

### 1. 生成证书

首先，需要生成CA证书和各服务的证书：

```bash
# 执行证书生成脚本
chmod +x scripts/generate_certs.sh
./scripts/generate_certs.sh
```

这将在`certs`目录下生成以下文件：
- `ca-cert.pem`：CA证书
- `ca-key.pem`：CA私钥
- `auth-cert.pem`：Auth服务证书
- `auth-key.pem`：Auth服务私钥
- 其他服务的证书和私钥（如果已生成）

### 2. 配置服务

每个服务需要在配置文件中添加TLS配置：

```yaml
tls:
  enabled: true
  ca_cert: "certs/ca-cert.pem"
  server_cert: "certs/服务名-cert.pem"
  server_key: "certs/服务名-key.pem"
  client_cert: "certs/服务名-cert.pem"
  client_key: "certs/服务名-key.pem"
```

### 3. 服务端配置

在服务的main.go文件中，添加TLS配置：

```go
// 如果启用了TLS，添加TLS配置
if conf.Config.TLS.Enabled {
    // 创建TLS配置
    tlsConfig, err := mtls.NewServerTLSConfig(mtls.CertConfig{
        CACertPath:     conf.Config.TLS.CACert,
        ServerCertPath: conf.Config.TLS.ServerCert,
        ServerKeyPath:  conf.Config.TLS.ServerKey,
    })
    if err != nil {
        log.Fatalf("创建TLS配置失败: %v", err)
    }
    
    // 添加TLS配置到服务器选项
    opts = append(opts, server.WithTLS(tlsConfig))
}
```

### 4. 客户端配置

在客户端代码中，添加TLS配置：

```go
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
```

## 证书管理

在生产环境中，证书管理是一个重要的问题。建议：

1. 使用安全的方式存储证书和私钥
2. 定期轮换证书（例如每年）
3. 考虑使用证书管理工具（如Vault）来管理证书
4. 在Kubernetes环境中，可以使用Secret来存储证书

## 故障排除

如果遇到TLS相关的问题，可以尝试：

1. 检查证书路径是否正确
2. 确保证书和私钥匹配
3. 验证CA证书是否正确导入
4. 检查证书的有效期
5. 使用OpenSSL工具验证证书：
   ```bash
   openssl verify -CAfile certs/ca-cert.pem certs/auth-cert.pem
   ```

## 下一步

完成Auth服务的mTLS配置后，可以按照类似的方式为其他服务添加mTLS支持。确保所有服务都使用相同的CA证书，这样它们才能相互验证身份。 