#!/usr/bin/env bash

CURDIR=$(cd $(dirname $0); pwd)
cd $CURDIR

# 检查环境
if [ -z "$KUBE_CONFIG" ]; then
    echo "未设置KUBE_CONFIG环境变量"
    exit 1
fi

# 构建镜像
docker build -t tiktok-mall/cart:latest ..

# 推送镜像到仓库
docker push tiktok-mall/cart:latest

# 创建配置
kubectl create configmap cart-config --from-file=../conf/prod/conf.yaml -n tiktok-mall

# 部署服务
kubectl apply -f ../deploy/cart.yaml

# 等待部署完成
kubectl rollout status deployment/cart -n tiktok-mall

# 检查服务状态
kubectl get pods -l app=cart -n tiktok-mall 
