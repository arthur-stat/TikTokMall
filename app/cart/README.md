# 购物车服务 (Cart Service)

## 功能概述

购物车服务提供了基本的购物车管理功能，包括：

- 添加商品到购物车
- 获取购物车内容
- 清空购物车

## 技术架构

### 数据存储

- MySQL: 持久化存储购物车数据
- Redis: 缓存购物车数据，提高读取性能

### 目录结构

```
app/cart/
├── biz/                    # 业务逻辑层
│   ├── dal/               # 数据访问层
│   │   ├── mysql/        # MySQL 相关操作
│   │   └── redis/        # Redis 相关操作
│   ├── handler/          # 业务处理器
│   ├── model/            # 数据模型
│   └── service/          # 服务实现
├── conf/                  # 配置文件
├── handler/              # RPC 处理器
├── kitex_gen/            # Kitex 生成的代码
└── scripts/              # 脚本文件
    └── test.sh          # 测试脚本
```

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
