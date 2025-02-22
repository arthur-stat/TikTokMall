> to co-workers:
> 
> 为了开发方便，目前跟根目录所有项目完全解耦合，基于product文件下的docker启动服务，所有测试数据都在测试文件中生成。后期联调部署时再做调整

# Quick Start

## How to run?

```shell
# 1. run docker compose
cd app/product
docker compose up -d

# 2. run service
cd app/product
go run .
```

## How to run tests?


```shell
# 1. service tests
cd app/product/biz/service
go test -v -cover

# 2. model tests
cd app/product/model
go test -v -cover
```

# Introduce - product

## Features

✅查询商品信息
- 根据ID查询单个商品
- 根据搜索查询商品列表
- 根据分类查询商品列表

❌创建商品（可选）

❌删除商品（可选）

❌修改商品信息（可选）

## Directory structure

| 目录/文件 | 说明 |
|----------|------|
| `biz/` | 业务逻辑层 |
| ├── `dal/` | 数据访问层 |
| ├── `service/` | 服务层实现及测试 |
| `conf/` | 配置文件 |
| ├── `dev/` | 开发环境配置 |
| ├── `online/` | 线上环境配置 |
| ├── `test/` | 测试环境配置 |
| `factory/` | 工厂模式 |
| `kitex_gen/` | Kitex 框架自动生成的API |
| `model/` | 数据模型定义、数据库交互层及测试 |
| `scripts/` | 构建、测试和部署脚本 |
| `handler.go` | 服务初始化及依赖注入组装 |
| `main.go` | 程序入口文件 |



