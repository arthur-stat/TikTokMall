# TikTokMall

TikTokMall 是一个基于字节跳动开源的 CloudWeGo 中间件集合的高性能微服务架构电商后端系统，使用 Go 语言开发，采用 Protocol Buffers 进行服务间通信。

## 开发指南
```bash
/docs/development_guide.md
```
## 技术栈

### 基础框架
- Go 1.23.4
- CloudWeGo 微服务框架集
  - Kitex: 高性能 RPC 框架
  - Hertz: HTTP 框架
- Protocol Buffers v29.3

### 存储
- MySQL: 主要数据存储
- Redis: 缓存、会话管理
- MongoDB: 日志存储（可选）

### 中间件
- Consul: 服务注册与发现
- Jaeger: 分布式链路追踪
- Prometheus: 监控系统
- Grafana: 可视化监控
- ELK Stack: 日志收集与分析
  - Elasticsearch
  - Logstash
  - Kibana

### 开发工具
- Docker: 容器化部署
- Docker Compose: 本地开发环境编排
- Make: 项目构建工具
- Swagger: API 文档生成
- GoMock: 单元测试 mock 工具

### 项目特性
- 微服务架构
- RPC 通信
- 服务注册与发现
- 负载均衡
- 熔断降级
- 链路追踪
- 监控告警
- 日志管理
- 缓存机制
- 分布式事务
- 统一认证

## 相关链接

- 仓库地址: [TikTok Mall](https://github.com/arthur-stat/TikTokMall)
- 文档地址: [TikTok Mall](https://uestc.feishu.cn/docx/T6HfdUzLqorZqaxpUfschLf2nKj)

## 项目结构

```bash
TikTokMall/
├── app/                    # 各微服务的服务端代码目录
│   ├── auth/               # 用户认证服务
│   ├── cart/               # 购物车服务
│   ├── checkout/           # 结算服务
│   ├── order/              # 订单服务
│   ├── payment/            # 支付服务
│   ├── product/            # 商品服务
│   └── user/               # 用户服务
├── deploy/                # 部署和基础设施配置
│   └── docker/           # Docker 相关配置
│       ├── docker-compose.yaml  # 基础服务编排配置
│       ├── mysql/        # MySQL 配置和初始化脚本
│       ├── prometheus/   # Prometheus 监控配置
│       └── logstash/     # ELK 日志收集配置
├── idl/                    # Protocol Buffers 定义目录
│   ├── api.proto          # API 通用注解文件
│   ├── auth.proto         # 用户认证服务的 .proto 文件
│   ├── cart.proto         # 购物车服务的 .proto 文件
│   ├── checkout.proto     # 结算服务的 .proto 文件
│   ├── order.proto        # 订单服务的 .proto 文件
│   ├── payment.proto      # 支付服务的 .proto 文件
│   ├── product.proto      # 商品服务的 .proto 文件
│   └── user.proto         # 用户服务的 .proto 文件
├── rpc_gen/                # 生成的客户端代码目录
│   ├── kitex_gen/         # Kitex 生成的客户端代码
│   └── rpc/               # 自定义的 RPC 客户端代码
├── README.md              # 项目说明文件
├── clean_generated_code.sh # 清理生成代码的脚本
├── generate_code.sh       # 生成代码的脚本
└── tidy_all.sh           # 整理和拉取依赖的脚本
```

## 微服务说明

1. **认证服务 (auth)**
   - 处理用户登录、注册、token管理等认证相关功能

2. **用户服务 (user)**
   - 管理用户信息、个人资料等

3. **商品服务 (product)**
   - 处理商品信息、库存管理等

4. **购物车服务 (cart)**
   - 管理用户购物车

5. **结算服务 (checkout)**
   - 处理订单结算流程

6. **订单服务 (order)**
   - 管理订单生命周期

7. **支付服务 (payment)**
   - 处理支付相关功能

## 微服务通用结构

每个微服务都遵循以下目录结构：

```bash
service/                   # 服务根目录
├── biz/                  # 业务逻辑目录
│   ├── dal/             # 数据访问层(Data Access Layer)
│   │   ├── init.go     # 数据库连接初始化
│   │   ├── mysql/      # MySQL 相关代码
│   │   │   ├── init.go        # MySQL 连接初始化
│   │   │   ├── model.go       # 数据库模型定义
│   │   │   └── crud.go        # 基础的 CRUD 操作
│   │   └── redis/      # Redis 相关代码
│   │       ├── init.go        # Redis 连接初始化
│   │       └── cache.go       # 缓存操作方法
│   └── service/        # 业务逻辑实现层
│       ├── service.go         # 服务接口定义
│       ├── service_impl.go    # 接口具体实现
│       └── service_test.go    # 单元测试
├── conf/                # 配置文件目录
│   ├── conf.go         # 配置结构定义和加载逻辑
│   ├── dev/           # 开发环境配置
│   │   └── conf.yaml   # 开发环境配置文件
│   ├── test/          # 测试环境配置
│   │   └── conf.yaml   # 测试环境配置文件
│   └── online/        # 生产环境配置
│       └── conf.yaml   # 生产环境配置文件
├── handler.go          # RPC 请求处理入口
├── main.go            # 服务启动入口
├── build.sh           # 服务构建脚本
├── Dockerfile         # Docker 构建文件
└── kitex.yaml         # Kitex RPC 框架配置
```

### 目录结构说明

1. **biz/dal (数据访问层)**
   - `init.go`: 负责初始化所有数据存储连接
   - **mysql/**: MySQL 数据库操作相关代码
     - `init.go`: MySQL 连接池初始化
     - `model.go`: 数据库模型定义，包含表结构和字段映射
     - `crud.go`: 基础的增删改查操作实现
   - **redis/**: Redis 缓存操作相关代码
     - `init.go`: Redis 连接池初始化
     - `cache.go`: 缓存的存取、更新、删除等操作

2. **biz/service (业务逻辑层)**
   - `service.go`: 定义服务接口和数据结构
   - `service_impl.go`: 实现具体的业务逻辑
   - `service_test.go`: 业务逻辑的单元测试

3. **conf (配置管理)**
   - `conf.go`: 定义配置结构和配置加载逻辑
   - 支持多环境配置：
     - `dev/`: 开发环境
     - `test/`: 测试环境
     - `online/`: 生产环境
   - 配置文件包含：
     - 服务基本信息（名称、端口等）
     - 数据库连接信息
     - 缓存配置
     - 日志配置
     - 服务发现配置

4. **根目录文件**
   - `handler.go`: 处理 RPC 请求，调用对应的业务逻辑
   - `main.go`: 服务启动入口，负责初始化和启动服务
   - `build.sh`: 服务构建脚本
   - `Dockerfile`: 定义 Docker 镜像构建步骤
   - `kitex.yaml`: Kitex RPC 框架的配置文件，包含服务端配置

## 快速开始

### 环境要求

- Go 1.23.4 或更高版本
- Protocol Buffers v29.3
- MySQL
- Redis
- Consul
- Docker & Docker Compose


### 常用命令

1. 生成 RPC 代码：
```bash
./generate_code.sh
```

2. 清理生成的代码：
```bash
./clean_generated_code.sh
```

3. 整理和拉取依赖：
```bash
./tidy_all.sh
```

## 贡献指南

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开一个 Pull Request

## 许可证

[MIT License](LICENSE)

## 系统架构

### 整体架构
```
┌─────────────────────────────────────────────────────────────┐
│                     TikTokMall 微服务架构                    │
├─────────────┬─────────────┬─────────────┬─────────────┬────┤
│  服务注册中心  │   消息队列   │    缓存     │   数据存储   │监控│
│   Consul    │   Kafka    │   Redis    │   MySQL    │    │
└─────────────┴─────────────┴─────────────┴─────────────┴────┘
         │            │            │            │         │
         └────────────┴────────────┴────────────┴─────────┘
                              │
                    ┌─────────┴─────────┐
                    │    微服务集群      │
        ┌──────────┬┴────────┬────────┬─┴──────────┐
        │          │         │        │            │
    ┌───────┐ ┌───────┐ ┌───────┐┌───────┐   ┌───────┐
    │ Auth  │ │ Cart  │ │Product││ Order │   │Payment│
    │Service│ │Service│ │Service││Service│   │Service│
    └───────┘ └───────┘ └───────┘└───────┘   └───────┘

```

### 核心组件

1. **服务注册与发现 (Consul)**
   - 服务自动注册
   - 服务健康检查
   - 服务发现与负载均衡
   - 配置中心

2. **数据存储**
   - MySQL: 持久化存储
   - Redis: 缓存层
   - MongoDB: 日志存储

3. **监控系统**
   - Prometheus: 指标收集
   - Grafana: 可视化面板
   - Jaeger: 分布式追踪
   - ELK Stack: 日志管理

### 微服务通信

1. **同步通信 (RPC)**
   - 使用 Kitex RPC 框架
   - Protobuf 序列化
   - 服务发现与负载均衡

2. **异步通信**
   - Kafka 消息队列
   - 事件驱动架构

### 服务说明

1. **认证服务 (Auth Service)**
   - 用户认证与授权
   - Token 管理
   - 端口: 8000
2. **用户服务（User Service）**
   - 与认证服务共同组成用户模块
   - 用户数据管理
   - 端口: 8001
3. **购物车服务 (Cart Service)**
   - 购物车管理
   - 商品缓存
   - 端口: 8002
4. **商品服务 (Product Service)**
   - 商品管理
   - 库存控制
   - 端口: 8003
5. **订单服务 (Order Service)**
   - 订单处理
   - 支付集成
   - 端口: 8004
6. **支付服务 (Payment Service)**
   - 支付处理
   - 退款管理
   - 端口: 8005

### 技术特性

1. **高可用**
   - 服务注册与发现
   - 负载均衡
   - 熔断降级
   - 失败重试

2. **可扩展**
   - 水平扩展
   - 模块化设计
   - 插件化架构

3. **可观测**
   - 分布式追踪
   - 性能监控
   - 日志聚合
   - 告警系统

## 本地部署

由于本项目的不同服务组件由不同的成员编写，虽热我们可以保证所有服务均在我们的开发环境下通过了测试，但受制于有限的时间与精力，我们没有对项目异地部署环境做充分的检测与兼容性。因此，若在不同的环境下尝试直接启动服务，则有启动失败或不能正常启动的可能。不过，如果严格按照说明文档配置中间件，则服务理应能够正常运行。

例如，要正常启动auth认证服务，则必须保证`consul`、`redis`、`mysql`、`jaeger`与`prometheus`这几个中间件全部事先正常启动并运行，才可以通过`go run app/auth/main.go`命令直接启动服务，否则启动失败。该服务的中间件配置与启动流程已在其核心组件路径下的说明文档中给出。