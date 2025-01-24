# 认证服务实现文档

## 1. 概述

认证服务是TikTokMall项目的核心基础服务之一，负责处理用户认证相关的所有功能，包括用户注册、登录、令牌刷新、登出和令牌验证等操作。

## 2. 目录结构

```
app/auth/
├── biz/
│   ├── dal/
│   │   ├── mysql/
│   │   │   ├── init.go    # MySQL初始化
│   │   │   ├── user.go    # 用户表操作
│   │   │   └── token.go   # Token表操作
│   │   └── redis/
│   │       ├── init.go    # Redis初始化
│   │       └── token.go   # Token缓存操作
│   ├── handler/
│   │   └── auth.go        # 请求处理层
│   └── service/
│       └── auth.go        # 业务逻辑层
└── README.md              # 本文档
```

## 3. 数据模型

### 3.1 User表
```sql
CREATE TABLE users (
    id         bigint       NOT NULL AUTO_INCREMENT,
    username   varchar(64)  NOT NULL,
    password   varchar(128) NOT NULL,
    email      varchar(128),
    phone      varchar(20),
    status     tinyint     NOT NULL DEFAULT 1,
    created_at timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY idx_username (username),
    KEY idx_email (email),
    KEY idx_phone (phone)
);
```

### 3.2 Token表
```sql
CREATE TABLE tokens (
    id            bigint       NOT NULL AUTO_INCREMENT,
    user_id       bigint       NOT NULL,
    token         varchar(512) NOT NULL,
    refresh_token varchar(512),
    expired_at    timestamp    NOT NULL,
    created_at    timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_user_id (user_id),
    KEY idx_token (token(191)),
    KEY idx_refresh_token (refresh_token(191))
);
```

## 4. 接口实现

### 4.1 用户注册 (Register)
- **接口**: `/v1/auth/register`
- **方法**: POST
- **功能**:
  - 验证用户名、邮箱、手机号唯一性
  - 密码加密存储
  - 创建用户记录
  - 生成访问令牌
- **安全措施**:
  - 密码使用bcrypt加密
  - 用户名长度限制：3-32字符
  - 密码长度限制：6-32字符
  - 邮箱格式验证
  - 手机号格式验证（中国大陆手机号）

### 4.2 用户登录 (Login)
- **接口**: `/v1/auth/login`
- **方法**: POST
- **功能**:
  - 验证用户名和密码
  - 生成访问令牌和刷新令牌
  - 记录登录重试次数
- **安全措施**:
  - 密码错误次数限制（5次/小时）
  - 返回统一的错误信息（不区分用户名或密码错误）
  - 用户状态检查

### 4.3 刷新令牌 (RefreshToken)
- **接口**: `/v1/auth/refresh`
- **方法**: POST
- **功能**:
  - 验证刷新令牌
  - 生成新的访问令牌和刷新令牌
  - 使旧令牌失效
- **安全措施**:
  - 刷新令牌有效期检查
  - 旧令牌自动加入黑名单

### 4.4 用户登出 (Logout)
- **接口**: `/v1/auth/logout`
- **方法**: POST
- **功能**:
  - 使当前令牌失效
  - 清理缓存
- **安全措施**:
  - 令牌加入黑名单
  - 同时清理数据库和缓存

### 4.5 验证令牌 (ValidateToken)
- **接口**: `/v1/auth/validate`
- **方法**: POST
- **功能**:
  - 验证令牌有效性
  - 返回用户信息
- **安全措施**:
  - 黑名单检查
  - 令牌过期检查
  - 用户状态验证

## 5. 缓存策略

### 5.1 Redis键设计
- 令牌缓存: `auth:token:{token}`
- 令牌黑名单: `auth:blacklist:{token}`
- 登录重试计数: `auth:retry:{username}`

### 5.2 缓存时间
- 访问令牌: 24小时
- 刷新令牌: 7天
- 黑名单记录: 与令牌有效期相同
- 登录重试记录: 1小时

## 6. 安全特性

1. **密码安全**
   - 使用bcrypt加密存储
   - 密码强度要求
   - 统一的错误消息

2. **令牌管理**
   - 双令牌机制（访问令牌+刷新令牌）
   - 令牌自动过期
   - 黑名单机制

3. **访问控制**
   - 登录重试限制
   - 用户状态检查
   - 令牌有效性验证

4. **数据安全**
   - 参数验证
   - SQL注入防护
   - XSS防护

## 7. 性能优化

1. **缓存优化**
   - 令牌信息缓存
   - 双写一致性保证
   - 缓存穿透防护

2. **数据库优化**
   - 合适的索引设计
   - 连接池管理
   - 超时控制

3. **并发处理**
   - 连接池配置
   - 超时设置
   - 错误处理

## 8. 错误处理

所有错误响应统一格式：
```json
{
    "base": {
        "code": <错误码>,
        "message": "<错误信息>"
    }
}
```

主要错误码：
- 400: 请求参数错误
- 401: 未授权
- 403: 禁止访问
- 429: 请求过于频繁
- 500: 服务器内部错误

## 9. 监控指标

1. **业务指标**
   - 注册成功率
   - 登录成功率
   - 令牌刷新率
   - 活跃用户数

2. **性能指标**
   - 接口响应时间
   - 缓存命中率
   - 错误率统计
   - 并发请求数

3. **安全指标**
   - 登录失败次数
   - 异常令牌使用
   - 黑名单命中率

## 10. 部署说明

1. **环境要求**
   - Go 1.20+
   - MySQL 8.0+
   - Redis 6.0+

2. **配置项**
   - 数据库连接
   - Redis连接
   - 令牌配置
   - 安全参数

3. **启动步骤**
   - 初始化数据库
   - 配置环境变量
   - 启动服务

## 11. 后续优化计划

1. **功能优化**
   - 添加OAuth2.0支持
   - 实现手机号验证码登录
   - 添加用户角色权限

2. **性能优化**
   - 引入分布式缓存
   - 优化数据库查询
   - 添加请求限流

3. **安全加强**
   - 实现IP限制
   - 添加设备指纹
   - 风控系统集成 