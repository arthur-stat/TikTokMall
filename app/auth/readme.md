# 认证服务 (Auth Service)

## 介绍

认证服务是 TikTokMall 电商平台的核心服务之一，负责处理用户认证、授权和令牌管理等功能。该服务基于 [Kitex](https://github.com/cloudwego/kitex/) 框架开发，提供了高性能的 RPC 接口，支持分布式部署和服务发现。

## 快速启动

### 1. 环境要求
- Go 1.21 或更高版本
- Docker 和 Docker Compose
- Git

### 2. 启动基础服务
```bash
# 停止并删除所有已存在的容器和卷
docker compose down -v

# 启动 MySQL 和 Redis
docker compose up -d mysql redis
```

### 3. 配置数据库
等待 MySQL 完全启动后（约15秒），执行以下命令：
```bash
# 连接到 MySQL
mysql -h 127.0.0.1 -P 3307 -u root -proot123

# 在 MySQL 中执行：
CREATE DATABASE IF NOT EXISTS tiktok_mall;
DROP USER IF EXISTS 'tiktok'@'%';
CREATE USER 'tiktok'@'%' IDENTIFIED WITH mysql_native_password BY 'tiktok123';
GRANT ALL PRIVILEGES ON tiktok_mall.* TO 'tiktok'@'%';
FLUSH PRIVILEGES;
```

### 4. 验证服务状态
```bash
# 检查 MySQL 连接
mysql -h 127.0.0.1 -P 3307 -u tiktok -ptiktok123 tiktok_mall

# 检查 Redis 连接
redis-cli -p 6380 ping
```

### 5. 启动服务
```bash
go run main.go
```

### 6. 服务端口
- RPC 服务端口：8888
- HTTP 服务端口：8000
- Prometheus 监控端口：9091
- MySQL 端口：3307
- Redis 端口：6380

## 配置说明

配置文件位于 `conf/test/conf.yaml`，包含：
- MySQL 连接信息
- Redis 连接信息
- 服务端口设置
- 日志级别设置

## 常见问题

1. MySQL 连接失败
   - 检查 MySQL 容器是否正常运行
   - 验证用户名和密码是否正确
   - 确认端口映射是否正确

2. Redis 连接失败
   - 检查 Redis 容器是否正常运行
   - 确认端口映射是否正确

3. 服务启动失败
   - 检查配置文件是否正确
   - 确认所有依赖服务是否正常运行
   - 查看日志获取详细错误信息

## 注意事项
- 首次启动时需要等待 MySQL 完全初始化
- 确保所需端口未被其他服务占用
- 开发环境下使用 test 配置文件

## 主要功能

- **用户认证**：处理用户注册和登录请求，支持多种认证方式。
- **令牌管理**：生成、刷新和验证访问令牌，管理令牌生命周期。
- **会话管理**：维护用户会话状态，支持分布式会话存储。
- **权限控制**：实现基于角色的访问控制，管理用户权限。
- **安全防护**：登录重试限制，令牌黑名单，防止恶意访问。

## 目录结构

| 目录            | 介绍                                                    |
|----------------|--------------------------------------------------------|
| `conf`         | 配置文件目录，包含服务的配置信息                            |
| `main.go`      | 服务启动文件，初始化服务并启动 RPC 和 HTTP 服务              |
| `handler`      | 请求处理层，负责接收和处理 RPC 请求                         |
| `kitex_gen`    | Kitex 框架自动生成的代码，包含 RPC 接口定义                 |
| `biz/service`  | 业务逻辑层，实现认证、授权等核心业务逻辑                     |
| `biz/dal`      | 数据访问层，负责与数据库和缓存交互                          |
| `pkg`          | 公共工具包，包含监控、追踪等组件                            |
| `deploy`       | 部署相关配置和脚本                                        |
| `scripts`      | 构建、测试和部署脚本                                      |

## 依赖项

认证服务依赖以下外部服务和技术栈：

- **MySQL**：用于持久化存储用户信息和令牌数据
- **Redis**：用于缓存会话信息和令牌黑名单
- **Consul**：用于服务注册与发现
- **Jaeger**：用于分布式追踪
- **Prometheus**：用于监控指标收集

## 如何运行

### 1. 环境配置

确保已安装以下工具和环境：

- Go 1.16 或更高版本
- MySQL
- Redis
- Consul
- Docker & Docker Compose

### 2. 启动服务

进入项目目录：
```bash
cd app/auth
```

安装依赖：
```bash
go mod tidy
```

启动依赖服务：
```bash
docker-compose up -d
```

启动认证服务：
```bash
sh scripts/build.sh
./output/bin/auth
```

### 3. API 接口

#### 用户注册 (Register)
接口描述：处理新用户注册请求。

请求方法：
```protobuf
rpc Register(RegisterRequest) returns (RegisterResponse)
```

请求参数：

| 字段名      | 类型    | 必填 | 描述        |
|------------|---------|------|------------|
| username   | string  | 是   | 用户名      |
| password   | string  | 是   | 密码        |
| email      | string  | 否   | 邮箱        |
| phone      | string  | 否   | 手机号      |

响应参数：

| 字段名     | 类型    | 描述         |
|-----------|---------|--------------|
| user_id   | int64   | 用户ID       |
| token     | string  | 访问令牌     |

#### 用户登录 (Login)
接口描述：处理用户登录请求。

请求方法：
```protobuf
rpc Login(LoginRequest) returns (LoginResponse)
```

请求参数：

| 字段名      | 类型    | 必填 | 描述        |
|------------|---------|------|------------|
| username   | string  | 是   | 用户名      |
| password   | string  | 是   | 密码        |

响应参数：

| 字段名          | 类型    | 描述         |
|----------------|---------|--------------|
| token          | string  | 访问令牌     |
| refresh_token  | string  | 刷新令牌     |

### 4. 错误码

| 错误码    | 描述              |
|----------|-------------------|
| 4001001  | 请求参数无效       |
| 5001001  | 内部服务器错误     |
| 4001002  | 用户名已存在       |
| 4001003  | 用户名或密码错误   |
| 4001004  | 令牌无效或已过期   |
| 4001005  | 登录重试次数超限   |

## 许可证

MIT License 