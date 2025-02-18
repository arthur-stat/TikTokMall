# Auth Service Documentation

## 目录
- [服务概述](#服务概述)
- [API 接口](#api-接口)
- [数据模型](#数据模型)
- [测试说明](#测试说明)
- [部署指南](#部署指南)
- [常见问题](#常见问题)

## 服务概述

Auth Service 是 TikTokMall 的认证服务，提供用户注册、登录、令牌管理等功能。主要特性包括：

- 用户注册和登录
- 访问令牌和刷新令牌管理
- 登录重试限制
- 令牌验证和黑名单
- Redis 缓存支持
- MySQL 持久化存储

## API 接口

### 1. 用户注册 (Register)

```protobuf
rpc Register(RegisterRequest) returns (RegisterResponse)
```

**请求参数：**
- username: 用户名（必填，长度3-32）
- password: 密码（必填，长度6-32）
- email: 邮箱（可选）
- phone: 手机号（可选）

**响应：**
- base: 基础响应信息
- data: 包含用户ID和访问令牌

**示例：**
```bash
curl -X POST http://localhost:8888/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "email": "test@example.com",
    "phone": "13800138000"
  }'
```

### 2. 用户登录 (Login)

```protobuf
rpc Login(LoginRequest) returns (LoginResponse)
```

**请求参数：**
- username: 用户名（必填）
- password: 密码（必填）

**响应：**
- base: 基础响应信息
- data: 包含访问令牌和刷新令牌

**特性：**
- 登录重试限制（5次/小时）
- 密码加密存储（bcrypt）

### 3. 刷新令牌 (RefreshToken)

```protobuf
rpc RefreshToken(RefreshTokenRequest) returns (RefreshTokenResponse)
```

**请求参数：**
- refresh_token: 刷新令牌（必填）

**响应：**
- base: 基础响应信息
- data: 包含新的访问令牌和刷新令牌

### 4. 登出 (Logout)

```protobuf
rpc Logout(LogoutRequest) returns (LogoutResponse)
```

**请求参数：**
- token: 当前访问令牌（必填）

**响应：**
- base: 基础响应信息

### 5. 验证令牌 (ValidateToken)

```protobuf
rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse)
```

**请求参数：**
- token: 访问令牌（必填）

**响应：**
- base: 基础响应信息
- data: 包含用户ID和用户名

## 数据模型

### User 表
```sql
CREATE TABLE users (
    id BIGINT NOT NULL AUTO_INCREMENT,
    username VARCHAR(32) NOT NULL,
    password VARCHAR(60) NOT NULL,
    email VARCHAR(64),
    phone VARCHAR(16),
    status TINYINT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_username (username),
    UNIQUE KEY uk_email (email),
    UNIQUE KEY uk_phone (phone)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### Token 表
```sql
CREATE TABLE tokens (
    id BIGINT NOT NULL AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    token VARCHAR(512) NOT NULL,
    refresh_token VARCHAR(512),
    expired_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_user_id (user_id),
    KEY idx_token (token(191)),
    KEY idx_refresh_token (refresh_token(191))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## 测试说明

### 运行测试
```bash
cd app/auth
./scripts/test.sh
```

测试脚本会：
1. 检查必要的依赖（MySQL、Redis、Go）
2. 准备测试数据库环境
3. 运行单元测试和集成测试
4. 生成测试覆盖率报告
5. 清理测试环境

### 测试用例说明

1. **注册测试**
   - 正常注册
   - 用户名已存在
   - 邮箱已存在
   - 手机号已存在

2. **登录测试**
   - 正常登录
   - 密码错误
   - 用户不存在
   - 登录重试限制

3. **令牌测试**
   - 刷新令牌
   - 令牌验证
   - 令牌过期
   - 令牌黑名单

### 测试覆盖率
当前测试覆盖率：
- 业务逻辑 (service): 71.1%
- 数据访问层 (dal): 70%+
- HTTP处理器 (handler): 需要改进
- 配置管理 (conf): 需要改进

## 部署指南

### 环境要求
- Go 1.20+
- MySQL 8.0+
- Redis 6.0+

### 配置说明
环境变量：
- `MYSQL_HOST`: MySQL 主机地址
- `MYSQL_PORT`: MySQL 端口
- `MYSQL_USER`: MySQL 用户名
- `MYSQL_PASSWORD`: MySQL 密码
- `MYSQL_DATABASE`: 数据库名称
- `REDIS_ADDR`: Redis 地址
- `REDIS_PASSWORD`: Redis 密码
- `REDIS_DB`: Redis 数据库编号

### 部署步骤
1. 编译服务
```bash
go build -o auth_service
```

2. 运行服务
```bash
./auth_service
```

## 常见问题

### 1. 登录重试限制
Q: 如何重置登录重试次数？
A: 成功登录后会自动重置，或等待1小时后自动重置。

### 2. 令牌过期
Q: 令牌有效期是多久？
A: 访问令牌24小时，刷新令牌7天。

### 3. 性能优化
- 使用 Redis 缓存令牌信息
- 定期清理过期令牌
- 使用令牌黑名单机制 

# Auth Service

认证服务，负责处理用户认证、授权和令牌管理。

## 功能特性

- 用户注册和登录
- 令牌管理（生成、刷新、验证）
- 登录重试限制
- 令牌黑名单
- 分布式会话管理

## 技术栈

- Kitex (RPC框架)
- GORM (ORM框架)
- MySQL (数据存储)
- Redis (缓存和会话管理)
- Jaeger (链路追踪)
- Prometheus (监控)
- Consul (服务发现)

## 快速开始

1. 安装依赖
```bash
go mod tidy
```

2. 配置环境
```bash
cp conf/dev/conf.yaml.example conf/dev/conf.yaml
# 修改配置文件
```

3. 启动服务
```bash
sh scripts/build.sh
./output/bin/auth
```

## API文档

### 注册
```protobuf
rpc Register(RegisterRequest) returns (RegisterResponse)
```

### 登录
```protobuf
rpc Login(LoginRequest) returns (LoginResponse)
```

### 刷新令牌
```protobuf
rpc RefreshToken(RefreshTokenRequest) returns (RefreshTokenResponse)
```

### 登出
```protobuf
rpc Logout(LogoutRequest) returns (LogoutResponse)
```

### 验证令牌
```protobuf
rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse)
```

## 监控指标

- auth_total: 认证请求总数
- auth_duration_seconds: 认证处理时间
- token_total: 令牌操作总数

## 配置说明

```yaml
env: "dev"  # 环境：dev, test, prod

kitex:
  service: "auth"  # 服务名
  address: ":8888" # 服务地址
  log_level: "info" # 日志级别

mysql:
  dsn: "用户名:密码@tcp(主机:端口)/数据库名"

redis:
  address: "localhost:6379"
  password: ""
  db: 0

registry:
  registry_address:
    - "localhost:8500"  # Consul地址

jaeger:
  host: "localhost"
  port: 6831

prometheus:
  port: 9090
  path: "/metrics"
```

## 部署

支持以下部署方式：
1. 直接部署
2. Docker部署
3. Kubernetes部署

详细部署步骤请参考[部署文档](deploy/README.md)

## 开发团队

- 开发人员1
- 开发人员2

## 许可证

MIT License 