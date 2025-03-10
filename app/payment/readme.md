# 支付服务 (Payment Service)

## 1.介绍

支付服务是 TikTokMall 电商平台的核心服务之一，负责处理用户的支付请求、支付状态管理以及退款等相关功能。该服务基于 [Kitex](https://github.com/cloudwego/kitex/) 和 [Hertz](https://github.com/cloudwego/hertz/)框架开发，提供了高性能的 RPC 和 HTTP 接口，支持分布式部署和服务发现。

## 2.主要功能

- **支付处理**：处理用户的支付请求，支持用卡支付（暂未使用支付接口）
- **退款管理**：处理用户的退款请求，支持部分退款和全额退款（暂未使用支付接口）（未完成）
- **支付宝支付**: 支持支付宝支付，通过支付接口进行支付和退款操作
- **支付宝退款**: 支持支付宝退款，通过支付接口进行退款操作（未完成）
- **支付状态管理**：记录支付状态，提供支付结果的查询接口（未完成）


## 3.目录结构

| 目录         | 介绍                                                                 |
|--------------|--------------------------------------------------------------------|
| `conf`       | 配置文件目录，包含服务的配置信息。                                       |
| `main.go`    | 服务启动文件，初始化服务并启动 RPC 和 HTTP 服务。                          |
| `handler`    | 请求处理层，负责接收和处理 RPC 请求，并返回响应。                           |
| `kitex_gen`  | Kitex 框架自动生成的代码，包含 RPC 接口定义和客户端代码。                    |
| `biz/service`| 业务逻辑层，实现支付、退款等核心业务逻辑。                                  |
| `biz/dal`    | 数据访问层，负责与数据库（如 MySQL）和缓存（如 Redis）进行交互。              |

## 4.依赖项

支付服务依赖以下外部服务和技术栈：

- **MySQL**：用于持久化存储支付记录和交易状态。
- **Redis**：用于缓存支付结果，提高查询性能。
- **Consul**：用于服务注册与发现，支持动态服务发现和负载均衡。
- **Kitex**：高性能 RPC 框架，用于服务间通信。
- **Hertz**：用于提供 HTTP 接口，支持 RESTful API。

## 5.如何运行

### 5.1. 环境配置

确保你已经安装了以下工具和环境：

- Go 1.16 或更高版本
- MySQL
- Redis
- Consul
- Kitex

### 5.2. 启动服务

若要使用支付宝相关服务，请修改conf文件

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

## 6. API 接口

### 6.1. 支付请求 (Charge)
接口描述： 处理用户的信用卡支付请求。（暂未使用支付接口）
#### 请求方法：
**RPC 请求：**
```Proto
    rpc Charge(ChargeReq) returns (ChargeResp)
```

**HTTP 请求：**
- URL: /payment/charge
- Method: POST
- Content-Type: application/json

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
| `credit_card_cvv`              | int32  | 是 | CVV  |
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

### 6.2. 退款请求 (Refund)
接口描述： 处理用户的信用卡退款请求。（暂未使用支付接口）
#### 请求方法：
**RPC 请求：**
```Proto
    rpc Charge(ChargeReq) returns (ChargeResp)
```

**HTTP 请求：**
- URL: /payment/charge
- Method: POST
- Content-Type: application/json

请求参数：

| 字段名              | 类型     | 必填 | 描述   |
|------------------|--------|----|------|
| `transaction_id` | string | 是  | 支付金额 |
| `order_id`       | int64  | 是  | 订单ID |
| `amount`         | string | 是  | 退款金额 |
| `user_id`        | int64  | 是  | 用户ID |

响应参数：

| 字段名 | 类型     | 描述   |
|-----|--------|------|
| `refund_id` | string | 退款单号 |

示例请求：
```json
{
  "transaction_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "order_id": 123456,
  "amount": 100.0,
  "user_id": 1
}
```
示例响应：
```json
{
  "refund_id": "a3f8d7e0-4b5c-11ee-be56-0242ac120002"
}
```

## 7. 错误码

### 7.1. charge 接口

#### 7.1.1. 错误码：**4004001**
- **描述**：信用卡信息验证失败，可能是卡号、CVV 或到期日期格式不正确
- **解决方案**：请检查信用卡信息并确保所有字段都符合标准格式

#### 7.1.2. 错误码：**4005001**
- **描述**：在尝试生成交易 ID 时遇到异常
- **解决方案**：请检查系统是否存在 UUID 生成或系统资源相关的问题/接口是否正常

#### 7.1.3. 错误码：**4005002**
- **描述**：订单已支付或已经存在相同订单，不能进行重复支付
- **解决方案**：请确认订单状态，确保未重复支付该订单

#### 7.1.4. 错误码：**4005003**
- **描述**：无法从 Redis 缓存中获取支付记录
- **解决方案**：检查 Redis 连接和配置，确保缓存服务可用

#### 7.1.5. 错误码：**4005004**
- **描述**：成功创建支付记录后，无法将其缓存到 Redis
- **解决方案**：检查 Redis 配置和服务状态，确认写入操作是否被正常处理

### 7.2. refund 接口

#### 7.2.1. 错误码：**4004002**
- **描述**：无效的退款请求
- **解决方案**：检查请求体格式

#### 7.2.2. 错误码：**4004003**
- **描述**：数据库查询失败
- **解决方案**：检查数据库连接及日志

#### 7.2.3. 错误码：**4004004**
- **描述**：支付记录不存在
- **解决方案**：验证transaction_id是否正确

#### 7.2.4. 错误码：**4004005**
- **描述**：支付状态不可退款

#### 7.2.5. 错误码：**4004006**
- **描述**：用户ID不匹配

#### 7.2.6. 错误码：**4004007**
- **描述**：订单ID不匹配

#### 7.2.7. 错误码：**4004008**
- **描述**：退款金额不符

#### 7.2.8. 错误码：**4005005**
- **描述**：退款单号生成失败
- - **解决方案**：重试或检查UUID生成器/接口是否正常

#### 7.2.9. 错误码：**4005006**
- **描述**：状态更新失败
- - **解决方案**：检查数据库事务及连接
