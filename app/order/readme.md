# 订单服务 (Order Service)

## 功能介绍

订单服务是 TikTokMall 电商平台的核心服务之一，负责处理订单的创建、查询、支付状态更新等功能。该服务基于 [Kitex](https://github.com/cloudwego/kitex/) 框架开发，提供了高性能的 RPC 接口。

## 主要功能

- 创建订单
- 查询订单列表
- 更新订单支付状态
- 订单状态管理
- 订单缓存处理

## 技术栈

- Go 1.23.4
- Kitex (RPC 框架)
- MySQL (数据存储)
- Redis (缓存)
- Jaeger (链路追踪)
- Prometheus (监控)

## 快速开始

### 环境要求

- Go 1.23.4 或更高版本
- Docker 和 Docker Compose
- MySQL 8.0
- Redis 6.0

### 本地开发

1. **启动依赖服务**
```bash
docker-compose up -d
```

2. **初始化数据库**
```bash
mysql -h 127.0.0.1 -P 3307 -u tiktok -ptiktok123 tiktok_mall < biz/dal/mysql/migrations/001_init_schema.sql
```

3. **运行服务**
```bash
go run main.go
```

### 配置说明

配置文件位于 `conf/dev/conf.yaml`，主要包含：

- 服务基本配置
- MySQL 配置
- Redis 配置
- 注册中心配置
- 日志配置
- 链路追踪配置
- 监控配置

### 接口说明

#### PlaceOrder
- 功能：创建订单
- 请求：`PlaceOrderReq`
- 响应：`PlaceOrderResp`

#### ListOrder
- 功能：获取订单列表
- 请求：`ListOrderReq`
- 响应：`ListOrderResp`

#### MarkOrderPaid
- 功能：标记订单已支付
- 请求：`MarkOrderPaidReq`
- 响应：`MarkOrderPaidResp`

## 部署

### Docker 部署
```bash
# 构建镜像
docker build -t tiktok-mall/order .

# 运行容器
docker run -d -p 8083:8080 tiktok-mall/order
```

### Kubernetes 部署
```bash
# 创建配置
kubectl create configmap order-config --from-file=conf/prod/conf.yaml -n tiktok-mall

# 部署服务
kubectl apply -f deploy/order.yaml
```

## 监控指标

- `order_total`: 订单操作总数
- `order_duration_seconds`: 订单操作处理时间
- `order_amount`: 订单金额分布

## 开发规范

- 代码格式化：`go fmt ./...`
- 代码检查：`go vet ./...`
- 单元测试：`go test ./...`
- 基准测试：`go test -bench=. ./...`

## 目录结构

```
app/order/
├── biz/                # 业务逻辑
│   ├── dal/           # 数据访问层
│   ├── handler/       # 请求处理器
│   ├── model/         # 数据模型
│   └── service/       # 服务实现
├── conf/              # 配置文件
├── pkg/               # 公共包
├── deploy/            # 部署配置
└── main.go           # 入口文件
```

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交代码
4. 创建 Pull Request

## 许可证

MIT License
