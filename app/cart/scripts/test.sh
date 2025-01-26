#!/bin/bash

# 设置错误时退出
set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

# 日志文件
LOG_FILE="test.log"

# 清理之前的日志
> $LOG_FILE

echo "检查依赖..."
echo "所有依赖检查通过"

# 设置MySQL环境变量
export MYSQL_USER="root"
export MYSQL_PASSWORD="root"
export MYSQL_HOST="127.0.0.1"
export MYSQL_PORT="3306"
export MYSQL_DATABASE="tiktok_mall_test"

# 设置Redis环境变量
export REDIS_HOST=127.0.0.1
export REDIS_PORT=6379
export REDIS_DB=0

echo "准备测试数据库..."
MYSQL_PWD=root mysql -uroot -h 127.0.0.1 -P 3306 -e "DROP DATABASE IF EXISTS tiktok_mall_test;"
MYSQL_PWD=root mysql -uroot -h 127.0.0.1 -P 3306 -e "CREATE DATABASE tiktok_mall_test;"
echo "测试数据库准备完成"

echo "运行测试..."
echo "检查Redis连接..."
if timeout 5 redis-cli -h 127.0.0.1 -p 6379 ping > /dev/null 2>&1; then
    echo "Redis连接成功"
    REDIS_AVAILABLE=true
else
    echo "Redis未运行或无法连接，跳过Redis相关测试"
    export SKIP_REDIS_TESTS=true
    REDIS_AVAILABLE=false
fi

echo "检查MySQL连接..."
MYSQL_PWD=root mysql -uroot -h 127.0.0.1 -P 3306 -e "SELECT 1;" > /dev/null 2>&1 || { 
    echo "MySQL未运行，请启动MySQL服务"
    exit 1
}

echo "开始运行单元测试..."

# 定义要测试的包
declare -A test_packages=(
    ["mysql"]="./biz/dal/mysql/..."
    ["model"]="./biz/model/..."
    ["biz_handler"]="./biz/handler/..."
    ["handler"]="./handler/..."
    ["conf"]="./conf/..."
)

# 测试结果
failed=0
coverage_data=""

# 逐个测试包
for key in "${!test_packages[@]}"; do
    pkg="${test_packages[$key]}"
    echo "测试包: $pkg"
    
    # 如果是Redis相关的包且Redis不可用，则跳过
    if [[ "$key" == "redis" ]] && [[ "$REDIS_AVAILABLE" == "false" ]]; then
        echo "跳过Redis测试"
        continue
    fi
    
    # 运行测试并捕获覆盖率
    if ! go test -v -cover "$pkg" 2>&1 | tee -a $LOG_FILE; then
        failed=1
        echo -e "${RED}包 $pkg 测试失败${NC}"
    else
        # 获取覆盖率数据
        coverage=$(go test -cover "$pkg" | grep -o '[0-9]*\.[0-9]*%')
        if [ ! -z "$coverage" ]; then
            coverage_data="$coverage_data\n$pkg: $coverage"
        fi
    fi
done

# 显示覆盖率报告
if [ ! -z "$coverage_data" ]; then
    echo -e "\n覆盖率报告:"
    echo -e "$coverage_data"
fi

if [ $failed -eq 0 ]; then
    echo -e "${GREEN}所有测试通过${NC}"
else
    echo -e "${RED}部分测试失败${NC}"
fi

echo "详细日志已保存到: $LOG_FILE"

echo "清理测试环境..."
MYSQL_PWD=root mysql -uroot -h 127.0.0.1 -P 3306 -e "DROP DATABASE IF EXISTS tiktok_mall_test;"
echo "测试环境清理完成"

# 如果测试失败，返回非零状态码
exit $failed
