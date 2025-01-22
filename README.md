# TikTokMall

依赖 / 版本：

- Go, 1.23.4
- Protocol Buffers, v29.3

仓库地址:[TikTok Mall](https://github.com/arthur-stat/TikTokMall)

文档地址:[TikTok Mall](https://uestc.feishu.cn/docx/T6HfdUzLqorZqaxpUfschLf2nKj)

# 项目结构
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
├── idl/                    # Protocol Buffers 定义目录
│   ├── auth.proto          # 用户认证服务的 .proto 文件
│   ├── cart.proto          # 购物车服务的 .proto 文件
│   ├── checkout.proto      # 结算服务的 .proto 文件
│   ├── order.proto         # 订单服务的 .proto 文件
│   ├── payment.proto       # 支付服务的 .proto 文件
│   ├── product.proto       # 商品服务的 .proto 文件
│   └── user.proto          # 用户服务的 .proto 文件
├── rpc_gen/                # 生成的客户端代码目录
├── README.md               # 项目简介文件
├── clean_generated_code.sh # 清理生成代码的脚本
├── generate_code.sh        # 生成代码的脚本
└── tidy_all.sh             # 整理和拉取依赖的脚本

```

# 每个微服务的文件结构
```bash

├── biz // 业务逻辑目录
│   ├── dal // 数据访问层 - 用来连接外部数据库进行database初始化、table创建等等
│   │   ├── init.go
│   │   ├── mysql
│   │   │   └── init.go
│   │   └── redis
│   │       └── init.go
│   └── service // service 层，业务逻辑存放的地方。更新时，新的方法会追加文件。
│       ├── HelloMethod.go
│       └── HelloMethod_test.go
├── build.sh
├── conf // 存放不同环境下的配置文件 - online/test/dev，通过环境变量设置
│     └── ...
├── docker-compose.yaml - docker启动mysql,consul等服务
├── go.mod // go.mod 文件，如不在命令行指定，则默认使用相对于 GOPATH 的相对路径作为 module 名
├── handler.go // 业务逻辑入口，更新时会全量覆盖
├── idl - 不一定在这里，我放在了项目根目录下
│   └── hello.thrift
├── kitex.yaml
├── kitex_gen // IDL 内容相关的生成代码，勿动 - 我生成的统一放在项目根目录的rpc_gen/kitex_gen下面了
│     └── ...
├── main.go // 程序入口 - 该服务的程序入口，比如auth服务逻辑从这里运行
├── readme.md
└── script // 启动脚本
└── bootstrap.sh
```

注释：

- auth：认证服务
- cart：购物车服务
- checkout：结算服务
- order：订单服务
- payment：支付服务
- product：商品服务
- user：用户服务

# 脚本
## 代码生成
```bash
./generate_code.sh
```
## 清理代码
```bash
./clean_generated_code.sh
```
## 整理和拉取依赖
```bash
./tidy_all.sh
```
