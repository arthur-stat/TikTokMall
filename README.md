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
- Nginx: 负载均衡与API网关
  - 请求路由与分发
  - SSL终止
  - 静态资源缓存
  - 限流与熔断
- Kafka/RabbitMQ: 消息队列
  - 异步通信
  - 事件驱动架构
  - 流量削峰
- Etcd: 分布式配置中心
  - 动态配置管理
  - 服务发现备选
- Redis Sentinel/Cluster: 高可用缓存
  - 主从复制
  - 故障转移

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
- 双向 TLS 认证 (mTLS)
  - 服务间安全通信
  - 证书管理
  - 身份验证
  - 防止中间人攻击

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
│   │   ├── conf/          # 配置文件目录
│   │   ├── pkg/           # 工具包目录
│   │   │   ├── client/    # 客户端工具
│   │   │   ├── hertz/     # Hertz框架工具
│   │   │   ├── metrics/   # 监控指标工具
│   │   │   ├── middleware/# 中间件工具
│   │   │   ├── mtls/      # mTLS 配置工具
│   │   │   ├── tracer/    # 链路追踪工具
│   │   │   └── utils/     # 通用工具
│   ├── payment/            # 支付服务
│   ├── product/            # 商品服务
│   └── user/               # 用户服务
├── certs/                  # 证书目录
│   ├── ca-cert.pem        # CA 证书
│   ├── ca-key.pem         # CA 私钥
│   ├── ca-cert.srl        # CA 证书序列号
│   ├── auth-cert.pem      # 认证服务证书
│   └── auth-key.pem       # 认证服务私钥
├── deploy/                # 部署和基础设施配置
│   └── docker/           # Docker 相关配置
│       ├── docker-compose.yaml  # 基础服务编排配置
│       ├── mysql/        # MySQL 配置和初始化脚本
│       ├── nginx/        # Nginx 配置
│       │   ├── conf.d/   # Nginx 配置文件
│       │   └── ssl/      # SSL 证书
│       ├── kafka/        # Kafka 配置
│       ├── rabbitmq/     # RabbitMQ 配置
│       ├── prometheus/   # Prometheus 监控配置
│       └── logstash/     # ELK 日志收集配置
├── docs/                   # 文档目录
│   ├── development_guide.md # 开发指南
│   └── mtls.md            # mTLS 配置文档
├── idl/                    # Protocol Buffers 定义目录
│   ├── api.proto          # API 通用注解文件
│   ├── auth.proto         # 用户认证服务的 .proto 文件
│   ├── cart.proto         # 购物车服务的 .proto 文件
│   ├── checkout.proto     # 结算服务的 .proto 文件
│   ├── order.proto        # 订单服务的 .proto 文件
│   ├── payment.proto      # 支付服务的 .proto 文件
│   ├── product.proto      # 商品服务的 .proto 文件
│   └── user.proto         # 用户服务的 .proto 文件
├── pkg/                    # 共享工具包目录
│   └── security/          # 安全相关工具
│       └── mtls/          # 共享 mTLS 工具
│           └── config.go  # TLS 配置实现
├── rpc_gen/                # 生成的客户端代码目录
│   ├── kitex_gen/         # Kitex 生成的客户端代码
│   └── rpc/               # 自定义的 RPC 客户端代码
├── scripts/                # 脚本目录
│   ├── clean_generated_code.sh # 清理生成代码的脚本
│   ├── generate_certs.sh  # 生成证书的脚本
│   ├── generate_code.sh   # 生成RPC代码的脚本
│   └── tidy_all.sh        # 整理和拉取依赖的脚本
├── README.md              # 项目说明文件
├── go.mod                 # Go模块定义
└── go.sum                 # Go依赖校验和
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
├── pkg/                # 工具包目录
│   ├── client/        # 客户端工具
│   │   └── options.go # 客户端选项配置
│   ├── mtls/          # mTLS 配置工具
│   │   └── config.go  # TLS 配置实现
│   └── utils/         # 通用工具
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
     - TLS 配置（证书路径、是否启用等）

4. **pkg (工具包)**
   - `client/`: 客户端工具
     - `options.go`: 客户端选项配置，包括 TLS 配置
   - `mtls/`: mTLS 配置工具
     - `config.go`: 实现 TLS 配置创建，包括服务端和客户端配置
   - `utils/`: 通用工具

5. **根目录文件**
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
- Consul/Etcd
- Kafka/RabbitMQ
- Nginx
- Docker & Docker Compose

### 中间件配置

#### Nginx 配置

TikTokMall 使用 Nginx 作为 API 网关和负载均衡器。基本配置示例：

```nginx
# TikTokMall API 网关配置
http {
    upstream auth_service {
        server localhost:8000;
        # 可添加更多实例实现负载均衡
    }
    
    upstream cart_service {
        server localhost:8002;
    }
    
    upstream order_service {
        server localhost:8004;
    }
    
    # 其他服务...
    
    server {
        listen 80;
        server_name api.tikmall.local;
        
        # SSL 配置
        # listen 443 ssl;
        # ssl_certificate /path/to/cert.pem;
        # ssl_certificate_key /path/to/key.pem;
        
        # 认证服务
        location /auth/ {
            proxy_pass http://auth_service/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }
        
        # 购物车服务
        location /cart/ {
            proxy_pass http://cart_service/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }
        
        # 订单服务
        location /order/ {
            proxy_pass http://order_service/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }
        
        # 限流配置
        limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;
        location /api/ {
            limit_req zone=api_limit burst=20;
            proxy_pass http://backend_service;
        }
    }
}
```

#### 消息队列配置

TikTokMall 支持 Kafka 或 RabbitMQ 作为消息队列。在服务配置中添加相应配置：

```yaml
# Kafka 配置示例
kafka:
  brokers:
    - "localhost:9092"
  topics:
    order_created: "order-created"
    payment_processed: "payment-processed"
  consumer_group: "tikmall-service"

# RabbitMQ 配置示例
rabbitmq:
  host: "localhost"
  port: 5672
  username: "guest"
  password: "guest"
  exchange: "tikmall"
  queues:
    order_created: "order-created"
    payment_processed: "payment-processed"
```

### mTLS 配置

TikTokMall 使用双向 TLS (mTLS) 进行服务间安全通信。要启用 mTLS，请按照以下步骤操作：

1. 证书准备：
   - 项目根目录的 `certs/` 文件夹包含了预生成的证书
   - 包括 CA 证书、服务证书和密钥

2. 服务配置：
   - 每个服务的配置文件中都有 TLS 相关配置项
   - 默认情况下，TLS 是禁用的 (`enabled: false`)
   - 要启用 TLS，将 `enabled` 设置为 `true`

3. 示例配置：
```yaml
tls:
  enabled: true
  ca_cert_path: "../../certs/ca-cert.pem"
  cert_path: "../../certs/service-cert.pem"
  key_path: "../../certs/service-key.pem"
  server_name: "service.name"
  client_verify: true
```

4. 客户端配置：
   - 每个服务都有对应的客户端包，用于其他服务调用
   - 客户端会自动读取 TLS 配置并建立安全连接

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
┌─────────────────────────────────────────────────────────────────────────────┐
│                          TikTokMall 微服务架构                               │
├─────────────┬─────────────┬─────────────┬─────────────┬────────────┬────────┤
│  服务注册中心  │   消息队列   │    缓存     │   数据存储   │  API网关   │  监控  │
│ Consul/Etcd  │Kafka/RabbitMQ│Redis Cluster│  MySQL    │   Nginx    │        │
└─────────────┴─────────────┴─────────────┴─────────────┴────────────┴────────┘
         │            │            │            │            │           │
         └────────────┴────────────┴────────────┴────────────┴───────────┘
                                   │
                        ┌──────────┴──────────┐
                        │      微服务集群      │
        ┌──────────┬────┴────────┬────────┬───┴──────────┐
        │          │             │        │              │
    ┌───────┐ ┌───────┐     ┌───────┐┌───────┐     ┌───────┐
    │ Auth  │ │ Cart  │     │Product││ Order │     │Payment│
    │Service│ │Service│     │Service││Service│     │Service│
    └───────┘ └───────┘     └───────┘└───────┘     └───────┘

```

### 核心组件

1. **服务注册与发现 (Consul/Etcd)**
   - 服务自动注册
   - 服务健康检查
   - 服务发现与负载均衡
   - 配置中心

2. **API网关 (Nginx)**
   - 外部请求入口
   - 路由分发到各微服务
   - 负载均衡
   - 请求限流与熔断
   - SSL终止
   - 静态资源缓存
   - 请求合并与转发

3. **消息队列 (Kafka/RabbitMQ)**
   - 异步通信
   - 事件驱动架构
   - 发布/订阅模式
   - 流量削峰填谷
   - 系统解耦

4. **数据存储**
   - MySQL: 持久化存储
   - Redis Cluster: 高可用缓存层
     - 主从复制
     - 故障转移
   - MongoDB: 日志存储

5. **监控系统**
   - Prometheus: 指标收集
   - Grafana: 可视化面板
   - Jaeger: 分布式链路追踪
   - ELK Stack: 日志管理

### 微服务通信

1. **API网关 (Nginx)**
   - 外部请求入口
   - 路由分发到各微服务
   - 负载均衡
   - 请求限流与熔断
   - SSL终止
   - 静态资源缓存
   - 请求合并与转发

2. **同步通信 (RPC)**
   - 使用 Kitex RPC 框架
   - Protobuf 序列化
   - 服务发现与负载均衡
   - 双向 TLS 认证 (mTLS)
     - 基于 X.509 证书的身份验证
     - 服务间通信加密
     - 防止中间人攻击和未授权访问

3. **异步通信**
   - Kafka/RabbitMQ 消息队列
   - 事件驱动架构
   - 发布/订阅模式
   - 流量削峰填谷
   - 系统解耦

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

4. **安全性**
   - 双向 TLS 认证 (mTLS)
     - 服务间通信加密
     - 基于证书的身份验证
     - 防止未授权服务接入
     - 支持证书轮换
   - 统一认证授权
   - 数据加密存储
   - 安全审计日志

## 本地部署

由于本项目的不同服务组件由不同的成员编写，虽热我们可以保证所有服务均在我们的开发环境下通过了测试，但受制于有限的时间与精力，我们没有对项目异地部署环境做充分的检测与兼容性。因此，若在不同的环境下尝试直接启动服务，则有启动失败或不能正常启动的可能。不过，如果严格按照说明文档配置中间件，则服务理应能够正常运行。

例如，要正常启动auth认证服务，则必须保证`consul`、`redis`、`mysql`、`jaeger`与`prometheus`这几个中间件全部事先正常启动并运行，才可以通过`go run app/auth/main.go`命令直接启动服务，否则启动失败。该服务的中间件配置与启动流程已在其核心组件路径下的说明文档中给出。