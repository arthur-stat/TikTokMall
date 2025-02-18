# Checkout Service

结账服务，负责处理用户订单结算和支付流程。

## 功能特性

- 订单创建和管理
- 支付处理
- 地址管理
- 用户信息管理

## 技术栈

- Kitex (RPC框架)
- GORM (ORM框架)
- MySQL (数据存储)
- Jaeger (链路追踪)
- Prometheus (监控)
- Consul (服务发现)

## 快速开始

1. 安装依赖
```bash
go mod tidy
```

2. 配置环境
```bash
cp conf/dev/conf.yaml.example conf/dev/conf.yaml
# 修改配置文件
```

3. 启动服务
```bash
sh build.sh
./output/bin/checkout
```

## API文档

### Checkout

```protobuf
service CheckoutService {
    rpc Run(CheckoutReq) returns (CheckoutResp)
}
```

请求参数：
- user_id: 用户ID
- firstname: 名字
- lastname: 姓氏
- email: 邮箱
- address: 地址信息
- credit_card: 信用卡信息

响应参数：
- order_id: 订单ID
- transaction_id: 交易ID

## 监控指标

- checkout_total: 结账请求总数
- checkout_duration_seconds: 结账处理时间
- payment_total: 支付处理总数

## 开发团队

- 开发人员1
- 开发人员2

## 许可证

MIT License 