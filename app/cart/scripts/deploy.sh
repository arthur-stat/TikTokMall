#!/bin/bash

# Color definitions
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Service configuration
SERVICE_NAME="cart"
SERVICE_PORT=8002

# Database configuration
MYSQL_HOST=${MYSQL_HOST:-"localhost"}
MYSQL_PORT=${MYSQL_PORT:-"3306"}
MYSQL_USER=${MYSQL_USER:-"root"}
MYSQL_PASS=${MYSQL_PASS:-"root"}
MYSQL_DB=${MYSQL_DB:-"tiktok_mall"}

# Redis configuration
REDIS_ADDR=${REDIS_ADDR:-"localhost:6379"}
REDIS_PASS=${REDIS_PASS:-""}

# Check dependencies
echo -e "${YELLOW}Checking dependencies...${NC}"

# Check MySQL client
if ! command -v mysql &> /dev/null; then
    echo -e "${RED}Error: MySQL client is not installed${NC}"
    exit 1
fi

# Check Redis
if ! command -v redis-cli &> /dev/null; then
    echo -e "${RED}Error: Redis is not installed${NC}"
    exit 1
fi

# Check Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    exit 1
fi

echo -e "${GREEN}All dependencies are installed${NC}"

# Check service status
echo -e "${YELLOW}Checking service status...${NC}"

# Test MySQL connection
if ! mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASS" -e "SELECT 1" &>/dev/null; then
    echo -e "${RED}Error: Could not connect to MySQL${NC}"
    exit 1
fi

# Test Redis connection
if ! redis-cli -h "${REDIS_ADDR%:*}" -p "${REDIS_ADDR#*:}" ping &>/dev/null; then
    echo -e "${RED}Error: Could not connect to Redis${NC}"
    exit 1
fi

echo -e "${GREEN}All services are running${NC}"

# Update dependencies
echo -e "${YELLOW}Updating dependencies...${NC}"
go mod tidy
if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Failed to update dependencies${NC}"
    exit 1
fi

# Build service
echo -e "${YELLOW}Building service...${NC}"
go build -o cart_service
if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Build failed${NC}"
    exit 1
fi

# Check if service is already running
PID=$(pgrep -f "cart_service")
if [ ! -z "$PID" ]; then
    echo -e "${YELLOW}Stopping existing service (PID: $PID)...${NC}"
    kill $PID
    sleep 2
fi

# Start service
echo -e "${YELLOW}Starting service...${NC}"
./cart_service > cart_service.log 2>&1 &
PID=$!

# Verify service is running
sleep 2
if ps -p $PID > /dev/null; then
    echo -e "${GREEN}Service started successfully (PID: $PID)${NC}"
else
    echo -e "${RED}Error: Service failed to start${NC}"
    exit 1
fi

echo -e "${GREEN}Deployment completed successfully${NC}" 
