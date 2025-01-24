#!/bin/bash

# 设置颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

# MySQL连接参数
MYSQL_USER="root"
MYSQL_PASS="root"
MYSQL_HOST="localhost"
MYSQL_PORT="3306"

# MySQL命令包装器
mysql_cmd() {
    sudo mysql "$@"
}

# 检查依赖
check_dependencies() {
    echo "检查依赖..."
    
    # 检查MySQL
    if ! command -v mysql &> /dev/null; then
        echo -e "${RED}错误: MySQL未安装${NC}"
        exit 1
    fi
    
    # 检查Redis
    if ! command -v redis-cli &> /dev/null; then
        echo -e "${RED}错误: Redis未安装${NC}"
        exit 1
    fi
    
    # 检查Go
    if ! command -v go &> /dev/null; then
        echo -e "${RED}错误: Go未安装${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}所有依赖检查通过${NC}"
}

# 准备测试数据库
prepare_test_db() {
    echo "准备测试数据库..."
    
    # 创建测试数据库
    mysql_cmd -e "CREATE DATABASE IF NOT EXISTS tiktok_mall_test;"
    
    # 导入数据库结构
    mysql_cmd tiktok_mall_test < ../../deploy/docker/mysql/init/init.sql
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}测试数据库准备完成${NC}"
    else
        echo -e "${RED}错误: 测试数据库准备失败${NC}"
        exit 1
    fi
}

# 运行测试
run_tests() {
    echo "运行测试..."
    
    # 设置测试环境变量
    export MYSQL_USER="$MYSQL_USER"
    export MYSQL_PASSWORD="$MYSQL_PASS"
    export MYSQL_HOST="$MYSQL_HOST"
    export MYSQL_PORT="$MYSQL_PORT"
    export MYSQL_DATABASE="tiktok_mall_test"
    export REDIS_ADDR="localhost:6379"
    export REDIS_PASSWORD=""
    
    # 运行所有测试并生成覆盖率报告
    go test -v -race -coverprofile=coverage.out ./...
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}所有测试通过${NC}"
        
        # 显示覆盖率报告
        go tool cover -func=coverage.out
        
        # 生成HTML覆盖率报告
        go tool cover -html=coverage.out -o coverage.html
        echo "覆盖率报告已生成: coverage.html"
    else
        echo -e "${RED}错误: 测试失败${NC}"
        exit 1
    fi
}

# 清理测试环境
cleanup() {
    echo "清理测试环境..."
    
    # 删除测试数据库
    mysql_cmd -e "DROP DATABASE IF EXISTS tiktok_mall_test;"
    
    # 清理Redis测试数据
    redis-cli -n 1 FLUSHDB
    
    # 删除覆盖率文件
    rm -f coverage.out
    
    echo -e "${GREEN}测试环境清理完成${NC}"
}

# 主流程
main() {
    # 检查依赖
    check_dependencies
    
    # 准备测试环境
    prepare_test_db
    
    # 运行测试
    run_tests
    
    # 清理环境
    cleanup
}

# 捕获中断信号，确保清理
trap cleanup EXIT

# 运行主流程
main 