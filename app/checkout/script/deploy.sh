#!/usr/bin/env bash

CURDIR=$(cd $(dirname $0); pwd)
cd $CURDIR

# 检查环境
if [ -z "$KUBE_CONFIG" ]; then
    echo "未设置KUBE_CONFIG环境变量"
    exit 1
fi

# 构建镜像
docker build -t tiktok-mall/checkout:latest ..

# 推送镜像到仓库
docker push tiktok-mall/checkout:latest

# 创建配置
kubectl create configmap checkout-config --from-file=../conf/prod/conf.yaml -n tiktok-mall

# 部署服务
kubectl apply -f ../deploy/checkout.yaml

# 等待部署完成
kubectl rollout status deployment/checkout -n tiktok-mall

# 检查服务状态
kubectl get pods -l app=checkout -n tiktok-mall 