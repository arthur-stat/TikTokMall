# 支付服务 (Payment Service)

## 介绍

支付服务是 TikTokMall 电商平台的核心服务之一，负责处理用户的支付请求、支付状态管理以及退款等相关功能。该服务基于 [Kitex](https://github.com/cloudwego/kitex/) 框架开发，提供了高性能的 RPC 接口，支持分布式部署和服务发现。

## 主要功能

- **支付处理**：处理用户的支付请求，支持多种支付方式（如信用卡支付）。
- **支付状态管理**：记录支付状态，提供支付结果的查询接口。（未完成）
- **退款管理**：处理用户的退款请求，支持部分退款和全额退款。（未完成）

## 目录结构

| 目录         | 介绍                                                                 |
|--------------|--------------------------------------------------------------------|
| `conf`       | 配置文件目录，包含服务的配置信息。                                       |
| `main.go`    | 服务启动文件，初始化服务并启动 RPC 和 HTTP 服务。                          |
| `handler`    | 请求处理层，负责接收和处理 RPC 请求，并返回响应。                           |
| `kitex_gen`  | Kitex 框架自动生成的代码，包含 RPC 接口定义和客户端代码。                    |
| `biz/service`| 业务逻辑层，实现支付、退款等核心业务逻辑。                                  |
| `biz/dal`    | 数据访问层，负责与数据库（如 MySQL）和缓存（如 Redis）进行交互。              |

## 依赖项

支付服务依赖以下外部服务和技术栈：

- **MySQL**：用于持久化存储支付记录和交易状态。
- **Redis**：用于缓存支付结果，提高查询性能。
- **Consul**：用于服务注册与发现，支持动态服务发现和负载均衡。
- **Kitex**：高性能 RPC 框架，用于服务间通信。
- **Hertz**：用于提供 HTTP 接口，支持 RESTful API。

## 如何运行

### 1. 环境配置

确保你已经安装了以下工具和环境：

- Go 1.16 或更高版本
- MySQL
- Redis
- Consul
- Kitex

### 2. 启动服务

由根目录进入项目根目录
``` bash
cd app/payment
```
启动前，先确保依赖安装成功
``` bash
go mod tidy
```
在项目根目录下执行以下命令来启动支付服务
``` bash
go run .
```
服务启动后，将会在 8885 端口上监听 RPC 请求，
并在 8005 端口上监听 HTTP 请求。

### 3.API 接口

3.1 支付请求 (Charge)
接口描述： 处理用户的支付请求，支持多种支付方式（如信用卡支付）。
请求方法：
```Proto
    rpc Charge(ChargeReq) returns (ChargeResp)
```
请求参数：

| 字段名            | 类型             | 必填 | 描述                     |
|------------------|----------------|---|------------------------|
| `amount`         | float          | 是 | 支付金额                   |
| `credit_card`    | CreditCardInfo | 是 | 信用卡信息                  |
| `payment_method` | string         | 否 | 支付方式，默认为 "credit_card" |
| `order_id`       | int64          | 是 | 订单ID                   |
| `user_id`        | int64          | 是 | 用户ID                   |

CreditCardInfo 结构：

| 字段名                            | 类型     | 必填 | 描述     |
|--------------------------------|--------|-|--------|
| `credit_card_number`           | string | 是 | 信用卡号   |
| `credit_card_cvv`              | nt32   | 是 | CVV  |
| `credit_card_expiration_year`  | int32  | 是 | 信用卡过期年份   |
| `credit_card_expiration_month` | int32  | 是 | 信用卡过期月份   |

响应参数：


| 字段名                        | 类型     | 描述   |
|----------------------------|--------|------|
| `transaction_id`           | string | 交易ID |

示例请求：
```json
{
  "amount": 100.0,
  "order_id": 123456,
  "user_id": 789,
  "payment_method": "credit_card",
  "credit_card": {
    "credit_card_number": "4111111111111111",
    "credit_card_cvv": 123,
    "credit_card_expiration_year": 2025,
    "credit_card_expiration_month": 12
  }
}
```
示例响应：
```json
{
    "transaction_id": "b24e05ec-3aff-4bee-b8fc-46345f50ea02"
}
```


### 4.错误码

| 错误码       | 描述      |
|-----------|---------|
| `4004001` | 请求参数无效  |
| `5005001` | 内部服务器错误 |
| `4004002` | 支付失败    |
| `4004003` | 订单不存在   |
| `4004004` | 用户不存在   |


