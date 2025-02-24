# 购物车服务 (Cart Service)

## 介绍

购物车服务是 TikTokMall 电商平台的核心服务之一，负责处理购物车管理、商品添加删除、数量修改等功能。该服务基于 [Kitex](https://github.com/cloudwego/kitex/) 框架开发，提供了高性能的 RPC 接口，支持分布式部署和服务发现。

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
go build -o cart_service
./cart_service

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

- **购物车管理**：添加、删除、修改购物车商品
- **商品数量控制**：修改商品数量，库存检查
- **商品选择**：选择/取消选择购物车商品
- **价格计算**：实时计算选中商品总价
- **库存同步**：与商品服务同步库存信息

## 目录结构

| 目录            | 介绍                                                    |
|----------------|--------------------------------------------------------|
| `conf`         | 配置文件目录，包含服务的配置信息                            |
| `main.go`      | 服务启动文件，初始化服务并启动 RPC 服务                    |
| `handler`      | 请求处理层，负责接收和处理 RPC 请求                        |
| `kitex_gen`    | Kitex 框架自动生成的代码，包含 RPC 接口定义                |
| `biz/service`  | 业务逻辑层，实现购物车管理等核心业务逻辑                    |
| `biz/dal`      | 数据访问层，负责与数据库和缓存交互                         |
| `pkg`          | 公共工具包，包含监控、追踪等组件                           |
| `deploy`       | 部署相关配置和脚本                                       |
| `scripts`      | 构建、测试和部署脚本                                     |

## 依赖项

购物车服务依赖以下外部服务和技术栈：

- **MySQL**：用于持久化存储购物车数据
- **Redis**：用于缓存购物车数据，提高读取性能
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
cd app/cart
```

安装依赖：
```bash
go mod tidy
```

启动依赖服务：
```bash
docker-compose up -d
```

启动购物车服务：
```bash
sh scripts/build.sh
./output/bin/cart
```

### 3. API 接口

#### 添加商品 (AddItem)
接口描述：添加商品到购物车。

请求方法：
```protobuf
rpc AddItem(AddItemReq) returns (AddItemResp)
```

请求参数：

| 字段名      | 类型           | 必填 | 描述        |
|------------|----------------|------|------------|
| user_id    | uint32        | 是   | 用户ID      |
| product_id | uint32        | 是   | 商品ID      |
| quantity   | uint32        | 是   | 商品数量     |
| selected   | bool          | 否   | 是否选中     |

#### 获取购物车 (GetCart)
接口描述：获取用户购物车内容。

请求方法：
```protobuf
rpc GetCart(GetCartReq) returns (GetCartResp)
```

请求参数：

| 字段名      | 类型    | 必填 | 描述        |
|------------|---------|------|------------|
| user_id    | uint32  | 是   | 用户ID      |

响应参数：

| 字段名     | 类型    | 描述         |
|-----------|---------|--------------|
| items     | []Item  | 购物车商品列表 |
| total     | float64 | 总价         |

### 4. 错误码

| 错误码    | 描述              |
|----------|-------------------|
| 4002001  | 请求参数无效       |
| 5002001  | 内部服务器错误     |
| 4002002  | 商品不存在         |
| 4002003  | 库存不足          |
| 4002004  | 购物车为空         |

## 许可证

MIT License

## 数据模型

### CartItem (购物车项)

```go
type CartItem struct {
    ID        uint32    // 唯一标识
    UserID    uint32    // 用户ID
    ProductID uint32    // 商品ID
    Quantity  uint32    // 数量
    Selected  bool      // 是否选中
    CreatedAt time.Time // 创建时间
    UpdatedAt time.Time // 更新时间
}
```

## API 接口

### AddItem

添加商品到购物车

```protobuf
rpc AddItem(AddItemReq) returns (AddItemResp)

message AddItemReq {
    uint32 user_id = 1;
    CartItem item = 2;
}
```

### GetCart

获取购物车内容

```protobuf
rpc GetCart(GetCartReq) returns (GetCartResp)

message GetCartReq {
    uint32 user_id = 1;
}

message GetCartResp {
    Cart cart = 1;
}
```

### EmptyCart

清空购物车

```protobuf
rpc EmptyCart(EmptyCartReq) returns (EmptyCartResp)

message EmptyCartReq {
    uint32 user_id = 1;
}
```

## 运行测试

### 环境要求

- MySQL 服务
- Redis 服务
- Go 1.23.4 或以上版本

### 配置数据库

1. MySQL 配置

```bash
export MYSQL_USER="root"
export MYSQL_PASSWORD="root"
export MYSQL_HOST="127.0.0.1"
export MYSQL_PORT="3306"
export MYSQL_DATABASE="tiktok_mall_test"
```

2. Redis 配置

```bash
export REDIS_HOST="127.0.0.1"
export REDIS_PORT="6379"
export REDIS_DB="0"
```

### 运行测试

```bash
cd app/cart
./scripts/test.sh
```

测试脚本会自动：

1. 检查依赖是否满足
2. 准备测试数据库
3. 运行单元测试
4. 生成覆盖率报告
5. 清理测试环境

### 测试用例说明

#### handler 测试

- `TestNewCartServiceImpl`: 测试服务实例创建
- `TestCartServiceImpl_AddItem`: 测试添加商品功能
  - 测试成功添加商品
  - 测试添加无效数量商品
- `TestCartServiceImpl_GetCart`: 测试获取购物车功能
  - 测试获取已存在的购物车
  - 测试获取空购物车
- `TestCartServiceImpl_EmptyCart`: 测试清空购物车功能
  - 测试清空非空购物车
  - 验证清空后的状态

#### dal/mysql 测试

- `TestCartItem_CRUD`: 测试基本的增删改查操作
- `TestCartItem_BatchOperations`: 测试批量操作
- `TestCartItem_Timestamps`: 测试时间戳自动更新

## 性能优化

1. Redis 缓存策略
   - 读取时优先从 Redis 获取
   - 写入时同步更新 MySQL 并失效缓存
   - 使用 Redis Pipeline 减少网络往返

2. MySQL 索引优化
   - 用户ID索引 (idx_user_id)
   - 自动递增主键

## 错误处理

服务定义了以下错误类型：

- `ErrInvalidQuantity`: 商品数量无效
- `ErrUserNotFound`: 用户不存在

## 注意事项

1. 测试数据隔离
   - 每个测试用例使用独立的用户ID
   - 测试完成后自动清理数据

2. 并发处理
   - 使用事务确保数据一致性
   - Redis 操作保证原子性

3. 环境配置
   - 确保 MySQL 和 Redis 服务正常运行
   - 正确设置环境变量
