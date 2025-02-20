#!/usr/bin/env bash

CURDIR=$(cd $(dirname $0); pwd)
cd $CURDIR

# 设置环境变量
export GO_ENV=prod

# 检查配置文件
if [ ! -f "../conf/prod/conf.yaml" ]; then
    echo "配置文件不存在"
    exit 1
fi

# 启动服务
../bin/checkout -conf ../conf/prod/conf.yaml

# 检查启动结果
if [ $? -ne 0 ]; then
    echo "服务启动失败"
    exit 1
fi
