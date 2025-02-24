# 结算服务 (Checkout Service)

## 介绍

结算服务是 TikTokMall 电商平台的核心服务之一，负责处理订单结算、支付流程、订单创建等功能。该服务基于 [Kitex](https://github.com/cloudwego/kitex/) 框架开发，提供了高性能的 RPC 接口，支持分布式部署和服务发现。

## 启动说明

### 1. 环境要求
- Go 1.23.4 或更高版本
- Docker 和 Docker Compose
- Git

### 2. 启动所有依赖服务
```bash
# 停止并删除所有已存在的容器和卷
docker compose down -v

# 启动所有依赖服务（MySQL、Redis、Consul、Jaeger、Prometheus）
docker compose up -d
```

### 3. 验证服务状态
```bash
# 检查 MySQL 连接
mysql -h 127.0.0.1 -P 3307 -u tiktok -ptiktok123 tiktok_mall

# 检查 Redis 连接
redis-cli -p 6380 ping

# 检查 Consul 状态
curl localhost:8501/v1/status/leader

# 检查 Jaeger UI
访问 http://localhost:16687

# 检查 Prometheus
访问 http://localhost:9091
```

### 4. 启动服务
```bash
# 编译并运行服务
go build -o checkout_service
./checkout_service

# 或直接运行
go run main.go
```

### 5. 服务端口说明
- RPC 服务端口：8888
- HTTP 服务端口：8000
- MySQL 端口：3307
- Redis 端口：6380
- Consul 端口：8501
- Jaeger 端口：6832(UDP), 16687(UI)
- Prometheus 端口：9091

## 配置说明

配置文件位于 `conf/test/conf.yaml`，包含：
- MySQL 连接信息
- Redis 连接信息
- Consul 服务发现配置
- Jaeger 链路追踪配置
- Prometheus 监控配置
- 服务端口设置
- 日志级别设置

## 主要功能

- **订单创建**：从购物车创建订单
- **支付处理**：集成支付服务，处理支付流程
- **订单状态管理**：管理订单生命周期
- **库存锁定**：确保订单商品库存
- **价格计算**：计算订单总价、优惠等

[其他内容与 auth 服务的 readme 类似，但针对结算服务特性进行相应修改]

## 目录结构

| 目录            | 介绍                                                    |
|----------------|--------------------------------------------------------|
| `conf`         | 配置文件目录，包含服务的配置信息                        |
| `main.go`      | 服务启动文件，初始化服务并启动 RPC 服务                |
| `handler`      | 请求处理层，负责接收和处理 RPC 请求                    |
| `kitex_gen`    | Kitex 框架自动生成的代码，包含 RPC 接口定义            |
| `biz/service`  | 业务逻辑层，实现订单处理、支付等核心业务逻辑           |
| `biz/dal`      | 数据访问层，负责与数据库和缓存交互                     |
| `pkg`          | 公共工具包，包含监控、追踪等组件                       |
| `deploy`       | 部署相关配置和脚本                                     |
| `scripts`      | 构建、测试和部署脚本                                   |

## 依赖项

结算服务依赖以下外部服务和技术栈：

- **MySQL**：用于持久化存储订单信息和交易记录
- **Redis**：用于缓存购物车数据和会话管理
- **Consul**：用于服务注册与发现
- **Jaeger**：用于分布式追踪
- **Prometheus**：用于监控指标收集
- **Payment Service**：用于处理实际支付流程

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
cd app/checkout
```

安装依赖：
```bash
go mod tidy
```

启动依赖服务：
```bash
docker-compose up -d
```

启动结算服务：
```bash
sh scripts/build.sh
./output/bin/checkout
```

### 3. API 接口

#### 结算请求 (Run)
接口描述：处理用户的结算请求，包括订单创建和支付处理。

请求方法：
```protobuf
rpc Run(CheckoutReq) returns (CheckoutResp)
```

请求参数：

| 字段名      | 类型           | 必填 | 描述        |
|------------|----------------|------|------------|
| user_id    | uint32        | 是   | 用户ID      |
| firstname  | string        | 是   | 名字        |
| lastname   | string        | 是   | 姓氏        |
| email      | string        | 是   | 邮箱        |
| address    | Address       | 是   | 地址信息     |
| credit_card| CreditCardInfo| 是   | 信用卡信息   |

响应参数：

| 字段名          | 类型    | 描述    |
|----------------|---------|---------|
| order_id       | string  | 订单ID   |
| transaction_id | string  | 交易ID   |

### 4. 错误码

| 错误码    | 描述          |
|----------|---------------|
| 4001001  | 请求参数无效   |
| 5001001  | 内部服务器错误 |
| 4001002  | 订单创建失败   |
| 4001003  | 支付处理失败   |
| 4001004  | 地址验证失败   |

## 许可证

MIT License 