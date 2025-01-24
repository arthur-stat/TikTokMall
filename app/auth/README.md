# Auth Service

认证服务负责处理用户认证相关的所有功能，包括注册、登录、token管理等。

## 功能列表

### 1. 用户认证
- 用户注册
- 用户登录
- 密码重置
- 登出操作

### 2. Token 管理
- JWT Token 生成
- Token 验证
- Token 刷新
- Token 撤销

### 3. 权限管理
- 角色管理
- 权限验证
- 访问控制

## API 接口

### 1. 注册接口
```
POST /v1/auth/register
```
请求体：
```json
{
    "username": "string",
    "password": "string",
    "email": "string",
    "phone": "string"
}
```
响应：
```json
{
    "code": 0,
    "message": "string",
    "data": {
        "user_id": "string",
        "token": "string"
    }
}
```

### 2. 登录接口
```
POST /v1/auth/login
```
请求体：
```json
{
    "username": "string",
    "password": "string"
}
```
响应：
```json
{
    "code": 0,
    "message": "string",
    "data": {
        "token": "string",
        "refresh_token": "string"
    }
}
```

### 3. 刷新Token接口
```
POST /v1/auth/refresh
```
请求头：
```
Authorization: Bearer {refresh_token}
```
响应：
```json
{
    "code": 0,
    "message": "string",
    "data": {
        "token": "string",
        "refresh_token": "string"
    }
}
```

### 4. 登出接口
```
POST /v1/auth/logout
```
请求头：
```
Authorization: Bearer {token}
```
响应：
```json
{
    "code": 0,
    "message": "string"
}
```

## 数据库设计

### users 表
```sql
CREATE TABLE IF NOT EXISTS `users` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `username` varchar(64) NOT NULL,
    `password` varchar(128) NOT NULL,
    `email` varchar(128),
    `phone` varchar(20),
    `status` tinyint NOT NULL DEFAULT 1,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`username`),
    KEY `idx_email` (`email`),
    KEY `idx_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### tokens 表
```sql
CREATE TABLE IF NOT EXISTS `tokens` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `user_id` bigint NOT NULL,
    `token` varchar(512) NOT NULL,
    `refresh_token` varchar(512),
    `expired_at` timestamp NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_token` (`token`(191)),
    KEY `idx_refresh_token` (`refresh_token`(191))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

## 目录结构

```bash
auth/
├── biz/                     # 业务逻辑目录
│   ├── dal/                # 数据访问层
│   │   ├── mysql/         # MySQL 相关代码
│   │   │   ├── init.go    # 数据库初始化
│   │   │   ├── user.go    # 用户表操作
│   │   │   └── token.go   # Token表操作
│   │   └── redis/         # Redis 相关代码
│   │       ├── init.go    # Redis 初始化
│   │       └── token.go   # Token 缓存操作
│   └── service/           # 业务逻辑实现
│       ├── auth.go        # 认证相关业务逻辑
│       └── token.go       # Token 相关业务逻辑
├── conf/                   # 配置文件目录
│   ├── dev/              # 开发环境配置
│   ├── test/             # 测试环境配置
│   └── online/           # 生产环境配置
├── handler.go             # 请求处理入口
└── main.go               # 服务入口文件
```

## 依赖服务

1. **MySQL**
- 存储用户信息
- 存储 Token 信息

2. **Redis**
- 缓存 Token
- 黑名单管理
- 限流计数

3. **Consul**
- 服务注册
- 服务发现

## 中间件集成

1. **JWT**
- 用于生成和验证 Token
- 使用 RS256 算法
- Token 有效期设置

2. **Jaeger**
- 请求链路追踪
- 性能监控

3. **Prometheus**
- 业务指标监控
- 性能指标采集

4. **ELK**
- 日志收集
- 日志分析

## 安全措施

1. **密码安全**
- 使用 bcrypt 加密存储
- 密码强度验证
- 密码重试限制

2. **Token 安全**
- JWT 签名验证
- Token 过期机制
- 刷新 Token 机制

3. **接口安全**
- 参数验证
- 频率限制
- 防重放攻击

## 监控指标

1. **业务指标**
- 注册用户数
- 登录次数
- Token 刷新次数
- 登录失败率

2. **性能指标**
- 接口响应时间
- 并发请求数
- 错误率统计

## 错误码设计

```go
const (
    SuccessCode           = 0
    InvalidParamCode     = 1001
    UserNotFoundCode     = 1002
    PasswordErrorCode    = 1003
    TokenExpiredCode     = 1004
    TokenInvalidCode     = 1005
    UserExistsCode       = 1006
    SystemErrorCode      = 5000
)
```

## 测试计划

1. **单元测试**
- 业务逻辑测试
- 数据访问测试
- 工具函数测试

2. **集成测试**
- API 接口测试
- 中间件集成测试
- 数据库操作测试

3. **性能测试**
- 并发性能测试
- 压力测试
- 稳定性测试 