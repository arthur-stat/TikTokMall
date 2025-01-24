#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 配置参数
SERVICE_NAME="auth_service"
SERVICE_PORT=8888
MYSQL_HOST=${MYSQL_HOST:-"localhost"}
MYSQL_PORT=${MYSQL_PORT:-"3306"}
MYSQL_USER=${MYSQL_USER:-"root"}
MYSQL_PASSWORD=${MYSQL_PASSWORD:-"123456"}
MYSQL_DATABASE=${MYSQL_DATABASE:-"tiktok_mall"}
REDIS_ADDR=${REDIS_ADDR:-"localhost:6379"}
REDIS_PASSWORD=${REDIS_PASSWORD:-""}
REDIS_DB=${REDIS_DB:-"0"}

# 检查依赖
check_dependencies() {
    echo -e "${YELLOW}正在检查依赖...${NC}"
    
    # 检查 Go
    if ! command -v go &> /dev/null; then
        echo -e "${RED}错误: Go 未安装${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Go 已安装${NC}"
    
    # 检查 MySQL
    if ! command -v mysql &> /dev/null; then
        echo -e "${RED}错误: MySQL 未安装${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ MySQL 已安装${NC}"
    
    # 检查 Redis
    if ! command -v redis-cli &> /dev/null; then
        echo -e "${RED}错误: Redis 未安装${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Redis 已安装${NC}"
}

# 检查服务状态
check_services() {
    echo -e "${YELLOW}正在检查服务状态...${NC}"
    
    # 检查 MySQL 连接
    if ! mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" -e "SELECT 1" &> /dev/null; then
        echo -e "${RED}错误: 无法连接到 MySQL${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ MySQL 连接正常${NC}"
    
    # 检查 Redis 连接
    if ! redis-cli -h $(echo $REDIS_ADDR | cut -d: -f1) -p $(echo $REDIS_ADDR | cut -d: -f2) ping &> /dev/null; then
        echo -e "${RED}错误: 无法连接到 Redis${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Redis 连接正常${NC}"
}

# 初始化数据库
init_database() {
    echo -e "${YELLOW}正在初始化数据库...${NC}"
    
    # 创建数据库（如果不存在）
    mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" -e "CREATE DATABASE IF NOT EXISTS $MYSQL_DATABASE DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
    
    # 导入数据库结构
    mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" < ../../deploy/docker/mysql/init/init.sql
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ 数据库初始化完成${NC}"
    else
        echo -e "${RED}错误: 数据库初始化失败${NC}"
        exit 1
    fi
}

# 编译服务
build_service() {
    echo -e "${YELLOW}正在编译服务...${NC}"
    
    # 更新依赖
    go mod tidy
    if [ $? -ne 0 ]; then
        echo -e "${RED}错误: 依赖更新失败${NC}"
        exit 1
    fi
    
    # 编译
    go build -o $SERVICE_NAME
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ 服务编译完成${NC}"
    else
        echo -e "${RED}错误: 服务编译失败${NC}"
        exit 1
    fi
}

# 启动服务
start_service() {
    echo -e "${YELLOW}正在启动服务...${NC}"
    
    # 检查服务是否已经运行
    if pgrep -f "$SERVICE_NAME" > /dev/null; then
        echo -e "${YELLOW}服务已在运行，正在重启...${NC}"
        pkill -f "$SERVICE_NAME"
        sleep 2
    fi
    
    # 设置环境变量
    export MYSQL_HOST="$MYSQL_HOST"
    export MYSQL_PORT="$MYSQL_PORT"
    export MYSQL_USER="$MYSQL_USER"
    export MYSQL_PASSWORD="$MYSQL_PASSWORD"
    export MYSQL_DATABASE="$MYSQL_DATABASE"
    export REDIS_ADDR="$REDIS_ADDR"
    export REDIS_PASSWORD="$REDIS_PASSWORD"
    export REDIS_DB="$REDIS_DB"
    
    # 后台启动服务
    nohup ./$SERVICE_NAME > auth.log 2>&1 &
    
    # 检查服务是否成功启动
    sleep 2
    if pgrep -f "$SERVICE_NAME" > /dev/null; then
        echo -e "${GREEN}✓ 服务启动成功${NC}"
        echo -e "${GREEN}服务日志: auth.log${NC}"
    else
        echo -e "${RED}错误: 服务启动失败${NC}"
        exit 1
    fi
}

# 主函数
main() {
    echo -e "${YELLOW}开始部署 Auth Service...${NC}"
    
    # 执行部署步骤
    check_dependencies
    check_services
    init_database
    build_service
    start_service
    
    echo -e "${GREEN}✓ Auth Service 部署完成${NC}"
    echo -e "${GREEN}服务运行在端口: $SERVICE_PORT${NC}"
}

# 执行主函数
main 