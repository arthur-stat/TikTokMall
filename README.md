# TikTokMall

依赖 / 版本：

- Go, 1.23.4
- Protocol Buffers, v29.3

仓库地址：[TikTok Mall](https://github.com/arthur-stat/TikTokMall)

# 文件结构

```powershell
│   main.go
│   README.md
│
├───.idea
│       .gitignore
│       .name
│       modules.xml
│       TikTokMall.iml
│       vcs.xml
│       workspace.xml
│
├───api
│   │   original_api.zip
│   │
│   ├───auth
│   │       auth.proto
│   │
│   ├───cart
│   │       cart.proto
│   │
│   ├───checkout
│   │       checkout.proto
│   │
│   ├───order
│   │       order.proto
│   │
│   ├───payment
│   │       payment.proto
│   │
│   ├───product
│   │       product.proto
│   │
│   └───user
│           user.proto
│
├───auth
│       auth.pb.go
│       auth_grpc.pb.go
│
├───cart
│       cart.pb.go
│       cart_grpc.pb.go
│
├───checkout
│       checkout.pb.go
│       checkout_grpc.pb.go
│
├───order
│       order.pb.go
│       order_grpc.pb.go
│
├───payment
│       payment.pb.go
│       payment_grpc.pb.go
│
├───product
│       product.pb.go
│       product_grpc.pb.go
│
└───user
        user.pb.go
        user_grpc.pb.go
```

注释：

- auth：认证服务
- cart：购物车服务
- checkout：结算服务
- order：订单服务
- payment：支付服务
- product：商品服务
- user：用户服务

# Protobuf接口编译

项目提供了`.proto`接口文件，需要对其进行编译为go文件

---

Linux安装：

```bash
sudo apt install protobuf-compiler  # or 'brew install protobuf'
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Windows下载：[Protocol Buffers](https://github.com/protocolbuffers/protobuf/releases)，安装后配置环境变量

---

示例

```bash
protoc --go_out=. --go-grpc_out=. proto/cart/cart.proto
```

