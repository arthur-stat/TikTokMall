#!/usr/bin/env bash

CURDIR=$(cd $(dirname $0); pwd)
PROJECT_ROOT=$(dirname "$CURDIR")

export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct

echo "开始运行测试..."

# 切换到项目根目录
cd $PROJECT_ROOT

# 运行单元测试
echo "运行单元测试..."

# 测试 service 包
go test -v -cover ./biz/service/...
if [ $? -ne 0 ]; then
    echo "service 包测试失败"
    exit 1
fi

# 测试 dal 包
go test -v -cover ./biz/dal/...
if [ $? -ne 0 ]; then
    echo "dal 包测试失败"
    exit 1
fi

# 测试 pkg 包
go test -v -cover ./pkg/...
if [ $? -ne 0 ]; then
    echo "pkg 包测试失败"
    exit 1
fi

# 运行集成测试
echo "运行集成测试..."
go test -v -tags=integration ./test/...
if [ $? -ne 0 ]; then
    echo "集成测试失败"
    exit 1
fi

# 运行基准测试
echo "运行基准测试..."
go test -v -bench=. -benchmem ./biz/service/...
if [ $? -ne 0 ]; then
    echo "基准测试失败"
    exit 1
fi

# 检查代码覆盖率
echo "检查代码覆盖率..."
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# 清理测试文件
rm coverage.out

echo "所有测试完成!"
exit 0 