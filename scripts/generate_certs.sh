#!/bin/bash
# 证书生成脚本 - 为TikTokMall微服务创建mTLS证书

# 设置变量
CERT_DIR="certs"
CA_KEY="$CERT_DIR/ca-key.pem"
CA_CERT="$CERT_DIR/ca-cert.pem"

# 创建证书目录
mkdir -p $CERT_DIR

echo "生成CA证书..."
# 生成CA私钥和证书
openssl genrsa -out $CA_KEY 4096
openssl req -new -x509 -days 365 -key $CA_KEY -out $CA_CERT -subj "/CN=TikTokMall CA"

# 为服务生成证书的函数
generate_service_cert() {
    SERVICE_NAME=$1
    echo "生成 $SERVICE_NAME 服务证书..."
    
    # 生成私钥
    openssl genrsa -out "$CERT_DIR/$SERVICE_NAME-key.pem" 2048
    
    # 生成CSR (Certificate Signing Request)
    openssl req -new -key "$CERT_DIR/$SERVICE_NAME-key.pem" \
        -out "$CERT_DIR/$SERVICE_NAME.csr" \
        -subj "/CN=$SERVICE_NAME.tiktok-mall.service"
    
    # 使用CA签名证书
    openssl x509 -req -days 365 \
        -in "$CERT_DIR/$SERVICE_NAME.csr" \
        -CA $CA_CERT -CAkey $CA_KEY -CAcreateserial \
        -out "$CERT_DIR/$SERVICE_NAME-cert.pem"
    
    # 清理CSR文件
    rm "$CERT_DIR/$SERVICE_NAME.csr"
    
    echo "$SERVICE_NAME 服务证书生成完成"
}

# 生成Auth服务的证书
generate_service_cert "auth"

# 可以为其他服务生成证书
# generate_service_cert "user"
# generate_service_cert "product"
# generate_service_cert "order"
# generate_service_cert "payment"
# generate_service_cert "cart"
# generate_service_cert "checkout"

echo "所有证书生成完成，保存在 $CERT_DIR 目录"
echo "请确保将证书复制到相应的服务目录" 