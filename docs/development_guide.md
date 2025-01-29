# TikTokMall 微服务开发指南

## 目录
1. [开发流程概述](#开发流程概述)
2. [环境准备](#环境准备)
3. [服务开发步骤](#服务开发步骤)
4. [测试规范](#测试规范)
5. [部署指南](#部署指南)
6. [最佳实践](#最佳实践)
7. [服务注册与发现](#服务注册与发现)

## 开发流程概述

### 1. 整体开发流程
1. 定义服务接口（Proto）
2. 生成基础代码
3. 实现业务逻辑
4. 编写单元测试
5. 集成测试
6. 部署服务

### 2. 代码结构规范
```
service/
├── biz/                    # 业务逻辑目录
│   ├── dal/               # 数据访问层
│   │   ├── mysql/        # MySQL 相关代码
│   │   └── redis/        # Redis 相关代码
│   ├── handler/          # 业务处理器
│   ├── model/            # 数据模型
│   └── service/          # 服务实现
├── conf/                  # 配置文件
├── handler/              # RPC 处理器
├── kitex_gen/            # 生成的代码
└── scripts/              # 脚本文件
```

## 环境准备

### 1. 开发环境要求
- Go 1.23.4 或更高版本
- Protocol Buffers v29.3
- MySQL 8.0+
- Redis 7.0+
- Docker & Docker Compose
- Kitex 代码生成工具
- Protoc 及其插件

### 2. 工具安装

#### 2.1 安装 CloudWeGo 工具集
```bash
# 安装 cwgo 工具
go install github.com/cloudwego/cwgo@latest

# 验证安装
cwgo -v
```

#### 2.2 环境变量配置
```bash
# 将 GOPATH/bin 添加到 PATH
# Linux/macOS (添加到 ~/.bashrc 或 ~/.zshrc)
export PATH="$PATH:$(go env GOPATH)/bin"

# Windows (添加到系统环境变量)
# %GOPATH%\bin
```

### 3. 手动代码生成

#### 3.1 生成服务代码
```bash
# 1. 进入项目根目录
cd /path/to/TikTokMall

# 2. 生成 RPC 代码（以 auth 服务为例）
cwgo server -idl idl/auth.proto -module TikTokMall -service auth

# 3. 为其他服务生成（替换参数为对应服务）
cwgo server -idl idl/cart.proto -module TikTokMall -service cart

# 参数说明：
# -idl: proto 文件路径
# -module: Go 模块名
# -service: 服务名
```

#### 3.2 生成客户端代码
```bash
# 1. 生成客户端代码（以 auth 服务为例）
cwgo client -idl idl/auth.proto -module TikTokMall -service auth

# 2. 为其他服务生成（替换参数为对应服务）
cwgo client -idl idl/cart.proto -module TikTokMall -service cart
```

#### 3.3 代码生成后的检查
```bash
# 1. 检查生成的目录结构
tree kitex_gen/

# 2. 检查 go.mod 依赖
go mod tidy

# 3. 验证导入路径
go build ./...
```

### 4. 环境配置
```bash
# 1. 安装依赖
go mod tidy

# 2. 配置环境变量
export MYSQL_USER="root"
export MYSQL_PASSWORD="root"
export MYSQL_HOST="127.0.0.1"
export MYSQL_PORT="3306"
export MYSQL_DATABASE="tiktok_mall"

export REDIS_HOST="127.0.0.1"
export REDIS_PORT="6379"
```

## 服务开发步骤

### 0. 开发前的准备工作

1. 检查 Go 语言环境：
```bash
# 在终端中输入：
go version

# 应该看到类似下面这样的输出
# go version go1.23.4 windows/amd64

```

2. 检查 MySQL 数据库：
```bash
# 在终端中输入：
mysql --version

# 如果看到类似这样的输出，就说明 MySQL 安装好了：
# mysql  Ver 8.0.xx

```

3. 检查 Redis 缓存：
```bash
# 在终端中输入：
redis-cli --version

# 如果看到类似这样的输出，说明 Redis 已经就位：
# redis-cli 7.0.xx

```

#### 0.2 安装缺少的工具

如果上面的检查发现有工具没安装：

1. 安装 MySQL：
   ```bash
   # 验证安装
   mysql --version
   ```

2. 安装 Redis：
   ```bash
   # Windows 用户：
   # 1. 访问：https://github.com/microsoftarchive/redis/releases
   # 2. 下载最新的 .msi 文件
   # 3. 双击安装，全部默认选项即可
   
   # Mac 用户：
   # 打开终端，输入：
   brew install redis
   # 如果提示没有 brew，先安装 brew：
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
   
   # Linux 用户：
   # Ubuntu:
   sudo apt-get update
   sudo apt-get install redis-server
   
   # CentOS:
   sudo yum install redis
   ```

#### 0.3 配置开发环境

```bash
# 第一步：设置 MySQL 的信息
# Windows 用户在终端中输入：
set MYSQL_USER=root
set MYSQL_PASSWORD=你的密码  # 替换成你之前设置的密码
set MYSQL_HOST=127.0.0.1
set MYSQL_PORT=3306
set MYSQL_DATABASE=tiktok_mall

# Mac/Linux 用户在终端中输入：
export MYSQL_USER=root
export MYSQL_PASSWORD=你的密码  # 替换成你之前设置的密码
export MYSQL_HOST=127.0.0.1
export MYSQL_PORT=3306
export MYSQL_DATABASE=tiktok_mall

# 第二步：设置 Redis 的信息
# Windows:
set REDIS_HOST=127.0.0.1
set REDIS_PORT=6379

# Mac/Linux:
export REDIS_HOST=127.0.0.1
export REDIS_PORT=6379
```

### 1. 服务注册与发现配置

#### 1.1 Consul 简介
Consul 是我们使用的服务注册与发现中间件,提供以下功能:
- 服务注册: 服务启动时自动注册到 Consul
- 健康检查: 定期检查服务是否健康
- 服务发现: 服务间调用时自动发现目标服务地址
- 负载均衡: 支持多实例间的负载均衡

#### 1.2 启动 Consul
```bash
# 使用 docker-compose 启动
docker-compose up -d consul

# 检查 Consul 状态
curl localhost:8500/v1/status/leader

# 访问 Consul UI
open http://localhost:8500/ui
```

#### 1.3 Consul 配置
```yaml
# conf/conf.yaml
registry:
  registry_address: ["localhost:8500"]  # Consul 地址
  username: ""                          # 认证用户名(可选)
  password: ""                          # 认证密码(可选)
```

### 2. 实现服务注册

#### 2.1 HTTP 服务注册 (以 auth 服务为例)
```go
// 创建 Consul 注册器
r, err := consul.NewConsulRegister("localhost:8500")
if err != nil {
    hlog.Fatalf("create consul register failed: %v", err)
}

// 创建 HTTP 服务器并注册到 Consul
h := server.Default(
    server.WithHostPorts(":8000"),
    server.WithRegistry(r, &registry.Info{
        ServiceName: "auth",
        Addr:        utils.NewNetAddr("tcp", "localhost:8000"),
        Weight:      10,
        Tags:        []string{"auth", "v1"},
    }),
)
```

#### 2.2 RPC 服务注册 (以 cart 服务为例)
```go
// 创建 Consul 注册器
r, err := consul.NewConsulRegister("localhost:8500")
if err != nil {
    log.Fatalf("create consul register failed: %v", err)
}

// 创建 RPC 服务器并注册到 Consul
addr, _ := net.ResolveTCPAddr("tcp", ":8002")
svr := cartservice.NewServer(
    handler.NewCartServiceImpl(),
    server.WithServiceAddr(addr),
    server.WithRegistry(r),
    server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
        ServiceName: "cart",
        Tags:        []string{"cart", "v1"},
    }),
)
```

### 3. 实现服务发现

#### 3.1 创建服务发现工具包
```go
// pkg/discovery/consul.go
func GetConsulClient(consulAddr string) ([]client.Option, error) {
    r, err := consul.NewConsulResolver(consulAddr)
    if err != nil {
        return nil, err
    }

    return []client.Option{
        client.WithResolver(r),
        client.WithLoadBalancer(loadbalance.NewWeightedBalancer()),
    }, nil
}
```

#### 3.2 使用服务发现调用其他服务
```go
// 创建带服务发现的客户端
consulOpts, err := discovery.GetConsulClient("localhost:8500")
if err != nil {
    return err
}

// 创建服务客户端 (以 cart 服务为例)
client, err := cartservice.NewClient("cart", consulOpts...)
if err != nil {
    return err
}

// 调用服务方法
resp, err := client.GetCart(ctx, req)
```

### 4. 验证服务注册与发现

#### 4.1 检查服务注册状态
```bash
# 通过 API 查看所有服务
curl localhost:8500/v1/catalog/services

# 查看特定服务的健康状态
curl localhost:8500/v1/health/service/[service-name]
```

#### 4.2 常见问题排查
- 服务注册失败: 检查 Consul 地址和连接
- 服务发现失败: 确认服务名称和标签是否正确
- 负载均衡异常: 检查服务权重配置

### 5. 实现数据访问层

在开始开发之前，我们需要了解两个主要的数据存储组件：

#### 1.1 MySQL 与 Redis 
- MySQL：用作主数据库，存储持久化数据
  - 适用场景：用户信息、订单数据等需要持久化的结构化数据
  - 特点：ACID 特性，适合事务操作
  
- Redis：用作缓存和临时存储
  - 适用场景：会话信息、临时令牌、热点数据缓存
  - 特点：高性能的键值存储，支持多种数据结构

### 2. 实现数据库操作（详细步骤）

#### 2.1 准备必要的工具包

首先，我们需要安装一些好用的工具包，就像是在工具箱里添加新工具一样：

```bash
# 装需要用到的包
# 帮我们更容易地操作数据库
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
go get -u github.com/redis/go-redis/v9

# 安装完成后，go.mod 文件会自动更新
# 我们可以检查一下：
cat go.mod
```

#### 2.2 创建数据库连接

接下来，我们要写代码来连接数据库。这就像是在我们的程序和数据库之间搭建一座桥：

1. 首先，创建 MySQL 连接文件：
```bash
# 打开这个文件：
app/user/biz/dal/mysql/init.go
```

2. 在文件中写入以下内容：
```go
package mysql

import (
    "fmt"
    "time"
    
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

// DB 是我们的数据库连接
// 就像是一个特殊的电话，我们可以通过它和数据库通话
var DB *gorm.DB

// Init 初始化数据库连接
// 这个函数就建立与数据库的连接
func Init() error {
    // 1. 准备连接信息
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        os.Getenv("MYSQL_USER"),        // 数据库用户名
        os.Getenv("MYSQL_PASSWORD"),    // 数据库密码
        os.Getenv("MYSQL_HOST"),        // 数据库地址
        os.Getenv("MYSQL_PORT"),        // 数据库端口
        os.Getenv("MYSQL_DATABASE"),    // 数据库名称
    )
    
    // 2. 尝试连接（拨打电话）
    var err error
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return fmt.Errorf("连接数据库失败: %w", err)
    }
    
    // 3. 设置连接池（就像是设置可以同时打多少个电话）
    sqlDB, err := DB.DB()
    if err != nil {
        return fmt.Errorf("获取数据库实例失败: %w", err)
    }
    
    // 最多同时连接10个
    sqlDB.SetMaxIdleConns(10)
    // 最多保持100个连接
    sqlDB.SetMaxOpenConns(100)
    // 一个连接最多重用1小时
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    // 4. 测试连接
    if err := sqlDB.Ping(); err != nil {
        return fmt.Errorf("无法连接数据库: %w", err)
    }
    
    return nil
}
```

3. 创建 Redis 连接文件：
```bash
# 创建并打开这个文件：
app/user/biz/dal/redis/init.go
```

4. 写入 Redis 连接代码：
```go
package redis

import (
    "context"
    "fmt"
    "time"
    
    "github.com/redis/go-redis/v9"
)

// 我们的 Redis 客户端
var RDB *redis.Client

// 初始化 Redis 连接
func Init() error {
    // 1. 创建 Redis 客户端（建立连接）
    RDB = redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
        Password: os.Getenv("REDIS_PASSWORD"), // 如果没有密码，就是空字符串
        DB:       0,                           // 使用默认数据库
    })
    
    // 2. 测试连接
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // 发送一个 PING 命令测试连接
    if err := RDB.Ping(ctx).Err(); err != nil {
        return fmt.Errorf("无法连接到 Redis: %w", err)
    }
    
    return nil
}
```

#### 2.3 创建用户模型

现在我们要定义用户数据的结构，就像是设计一张表格，规定每个用户都有哪些信息：

1. 创建用户模型文件：
```bash
# 创建并打开这个文件：
app/user/biz/model/user.go
```

2. 写入用户模型代码：
```go
package model

import (
    "time"
)

// User 代表一个用户
// 就像是一张表格，记录用户的各种信息
type User struct {
    ID        int64     `json:"id" gorm:"primaryKey"`  // 用户的唯一编号
    Username  string    `json:"username"`              // 用户名
    Password  string    `json:"-"`                     // 密码（json:"-" 表示不要在 JSON 中显示密码）
    Email     string    `json:"email"`                 // 邮箱
    Phone     string    `json:"phone"`                 // 手机号
    CreatedAt time.Time `json:"created_at"`           // 创建时间
    UpdatedAt time.Time `json:"updated_at"`           // 更新时间
}

// TableName 指定表名
// 告诉数据库这个结构对应哪张表
func (User) TableName() string {
    return "users"
}
```

好了！现在我们已经：
1. ✅ 安装了需要的工具包
2. ✅ 创建了数据库连接
3. ✅ 定义了用户数据结构

接下来，我们就可以开始实现具体的业务逻辑了

### 6. 实现业务逻辑

#### 3.1 创建用户操作接口

首先，我们要定义用户可以进行哪些操作。这就像是写一份菜单，列出我们的服务能提供什么功能：

1. 创建服务接口文件：
```bash
# 创建并打开这个文件：
app/user/biz/service/user.go
```

2. 写入服务接口代码：
```go
package service

import (
    "context"
    "TikTokMall/app/user/biz/model"
)

// UserService 定义了用户服务可以做的所有事情
// 就像一份菜单，列出所有可用的功能
type UserService interface {
    // 创建新用户
    CreateUser(ctx context.Context, username, password, email, phone string) (*model.User, error)
    
    // 根据用户名获取用户信息
    GetUserByUsername(ctx context.Context, username string) (*model.User, error)
    
    // 验证用户密码
    ValidatePassword(ctx context.Context, username, password string) (bool, error)
    
    // 更新用户信息
    UpdateUser(ctx context.Context, user *model.User) error
}
```

#### 3.2 实现具体的业务逻辑

现在我们知道要做什么了，接下来就要实现这些功能。

1. 创建服务实现文件：
```bash
# 创建并打开这个文件：
app/user/biz/service/user_impl.go
```

2. 写入具体实现代码：
```go
package service

import (
    "context"
    "fmt"
    "time"
    "golang.org/x/crypto/bcrypt"
    
    "TikTokMall/app/user/biz/dal/mysql"
    "TikTokMall/app/user/biz/model"
)

// userService 是我们的厨师，负责实现所有的功能
type userService struct{}

// NewUserService 创建一个新的用户服务
func NewUserService() UserService {
    return &userService{}
}

// CreateUser 创建新用户
func (s *userService) CreateUser(ctx context.Context, username, password, email, phone string) (*model.User, error) {
    // 1. 检查用户名是否已经存在
    existingUser, err := mysql.GetUserByUsername(ctx, username)
    if err != nil {
        return nil, fmt.Errorf("检查用户名时出错: %w", err)
    }
    if existingUser != nil {
        return nil, fmt.Errorf("用户名 %s 已经被使用了", username)
    }
    
    // 2. 对密码进行加密
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, fmt.Errorf("密码加密失败: %w", err)
    }
    
    // 3. 创建新用户
    user := &model.User{
        Username:  username,
        Password:  string(hashedPassword),
        Email:     email,
        Phone:     phone,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    // 4. 保存到数据库
    if err := mysql.CreateUser(ctx, user); err != nil {
        return nil, fmt.Errorf("创建用户失败: %w", err)
    }
    
    return user, nil
}

// GetUserByUsername 根据用户名查找用户
func (s *userService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
    user, err := mysql.GetUserByUsername(ctx, username)
    if err != nil {
        return nil, fmt.Errorf("查找用户失败: %w", err)
    }
    return user, nil
}

// ValidatePassword 验证用户密码
func (s *userService) ValidatePassword(ctx context.Context, username, password string) (bool, error) {
    // 1. 先找到用户
    user, err := s.GetUserByUsername(ctx, username)
    if err != nil {
        return false, err
    }
    if user == nil {
        return false, nil
    }
    
    // 2. 验证密码
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    return err == nil, nil
}

// UpdateUser 更新用户信息
func (s *userService) UpdateUser(ctx context.Context, user *model.User) error {
    user.UpdatedAt = time.Now()
    if err := mysql.UpdateUser(ctx, user); err != nil {
        return fmt.Errorf("更新用户信息失败: %w", err)
    }
    return nil
}
```

#### 3.3 添加错误处理

在实际使用中，难免会遇到一些问题，我们需要优雅地处理这些错误：

1. 创建错误定义文件：
```bash
# 创建并打开这个文件：
app/user/biz/service/errors.go
```

2. 写入错误定义代码：
```go
package service

import "errors"

// 定义可能遇到的错误
var (
    // 当用户名已经存在时的错误
    ErrUserExists = errors.New("用户名已存在")
    
    // 当用户不存在时的错误
    ErrUserNotFound = errors.New("用户不存在")
    
    // 当密码错误时的错误
    ErrInvalidPassword = errors.New("密码不正确")
    
    // 当输入的信息不完整时的错误
    ErrInvalidInput = errors.New("请填写完整的信息")
)
```

好了，现在我们已经：
1. ✅ 定义了用户服务的功能清单
2. ✅ 实现了每个功能的具体逻辑
3. ✅ 添加了错误处理机制

接下来，我们就要开始实现 RPC 接口，让其他服务能够调用我们的功能了

### 7. 实现 RPC 接口

RPC接口是服务对外提供的入口，它接收请求并返回响应。

```go
// handler/handler.go
package handler

import (
    "context"

    "TikTokMall/app/user/biz/service"
    "TikTokMall/app/user/kitex_gen/user"
)

// UserServiceImpl 实现了 RPC 服务接口
type UserServiceImpl struct {
    svc service.UserService  // 业务服务实例
}

// NewUserServiceImpl 创建一个新的服务实现实例
func NewUserServiceImpl() *UserServiceImpl {
    return &UserServiceImpl{
        svc: service.NewUserService(),
    }
}

// CreateUser 处理创建用户的 RPC 请求
func (s *UserServiceImpl) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
    // 1. 检查请求参数是否合法
    if req == nil || req.Username == "" {
        return &user.CreateUserResponse{
            BaseResp: &user.BaseResp{
                Code:    400,
                Message: "用户名不能为空",
            },
        }, nil
    }

    // 2. 调用业务层创建用户
    userModel, err := s.svc.CreateUser(ctx, req.Username, req.Password, req.Email, req.Phone)
    if err != nil {
        return &user.CreateUserResponse{
            BaseResp: &user.BaseResp{
                Code:    500,
                Message: err.Error(),
            },
        }, nil
    }

    // 3. 返回成功响应
    return &user.CreateUserResponse{
        BaseResp: &user.BaseResp{
            Code:    200,
            Message: "创建成功",
        },
        Data: &user.User{
            Id:       userModel.ID,
            Username: userModel.Username,
            Email:    userModel.Email,
            Phone:    userModel.Phone,
        },
    }, nil
}

// ... 实现其他 RPC 方法 ...
```

### 8. 编写启动程序

最后，我们需要编写主程序来启动服务。

```go
// main.go
package main

import (
    "log"
    "net"

    "TikTokMall/app/user/biz/dal/mysql"
    "TikTokMall/app/user/biz/dal/redis"
    "TikTokMall/app/user/conf"
    "TikTokMall/app/user/handler"
    "TikTokMall/app/user/kitex_gen/user/userservice"
)

func main() {
    // 1. 加载配置文件
    config := conf.LoadConfig()

    // 2. 初始化 MySQL 连接
    if err := mysql.Init(config); err != nil {
        log.Fatalf("初始化 MySQL 失败: %v", err)
    }

    // 3. 初始化 Redis 连接
    if err := redis.Init(config); err != nil {
        log.Fatalf("初始化 Redis 失败: %v", err)
    }

    // 4. 创建 RPC 服务器
    addr, _ := net.ResolveTCPAddr("tcp", ":8888")
    svr := userservice.NewServer(
        handler.NewUserServiceImpl(),
        server.WithServiceAddr(addr),
    )

    // 5. 启动服务
    log.Println("服务开始启动...")
    if err := svr.Run(); err != nil {
        log.Fatalf("服务运行失败: %v", err)
    }
}
```

### 9. 如何运行和测试

1. 确保环境准备好了：
   - MySQL 已经启动
   - Redis 已经启动
   - 配置文件内容正确

2. 编译和运行：
```bash
# 在项目根目录下
go mod tidy     # 更新依赖
go build        # 编译程序
./user          # 运行程序
```

3. 测试服务是否正常：
```bash
# 使用测试工具（如 BloomRPC）发送请求：
# 创建用户请求示例：
{
    "username": "test_user",
    "password": "test_password",
    "email": "test@example.com",
    "phone": "1234567890"
}
```

## 测试规范

### 1. 单元测试
```go
// service_test.go
func TestService_Method(t *testing.T) {
    tests := []struct {
        name    string
        input   *pb.Request
        want    *pb.Response
        wantErr bool
    }{
        // 定义测试用例
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 执行测试
        })
    }
}
```

### 2. 运行测试
```bash
# 在服务目录下执行
./scripts/test.sh
```

测试脚本会：
1. 检查依赖
2. 准备测试数据库
3. 运行单元测试
4. 生成覆盖率报告
5. 清理测试环境

## 部署指南

### 1. 本地部署
```bash
# 在服务目录下执行
./scripts/deploy.sh
```

### 2. Docker 部署
```bash
# 在项目根目录下执行
docker-compose up -d service_name
```

## 最佳实践

### 1. 代码规范
- 使用 gofmt 格式化代码
- 遵循 Go 官方代码规范
- 添加适当的注释

### 2. 错误处理
```go
// 定义错误类型
var (
    ErrNotFound = errors.New("resource not found")
    ErrInvalid  = errors.New("invalid input")
)

// 错误处理
if err != nil {
    return nil, fmt.Errorf("failed to process: %w", err)
}
```

### 3. 日志规范
```go
// 使用结构化日志
klog.InfoF("processing request: %v", req)
klog.ErrorF("failed to process: %v", err)
```

### 4. 配置管理
```yaml
# conf/conf.yaml
service:
  name: "your_service"
  port: 8080

mysql:
  host: "localhost"
  port: 3306
```

### 5. 性能优化
- 使用连接池
- 实现缓存策略
- 添加适当的索引
- 使用批量操作

### 6. 监控指标
- 请求延迟
- 错误率
- 并发连接数
- 资源使用率

## 示例服务

### Auth 服务
参考 `app/auth/` 目录下的实现：
- 用户认证
- Token 管理
- 登录重试限制

### Cart 服务
参考 `app/cart/` 目录下的实现：
- 购物车管理
- 缓存策略
- 数据一致性

## 常见问题

### 1. 代码生成
Q: 修改 proto 文件后需要做什么？
A: 需要重新生成代码：
1. 使用 cwgo server 重新生成服务端代码
2. 使用 cwgo client 重新生成客户端代码
3. 运行 `go mod tidy` 更新依赖
4. 更新相关业务实现

Q: 生成代码时遇到错误怎么办？
A: 常见问题检查：
1. 确保 cwgo 已正确安装且在 PATH 中
2. 检查 proto 文件语法是否正确
3. 确保在正确的目录下执行命令
4. 检查模块名是否与 go.mod 一致

Q: proto 文件之间有依赖关系，应该如何处理？
A: 按依赖顺序生成代码：
1. 先生成被依赖的 proto 文件代码
2. 再生成依赖其他 proto 的文件代码
3. 确保 proto 文件中的 import 路径正确

### 2. 测试数据库
Q: 如何准备测试数据库？
A: 测试脚本会自动创建和初始化测试数据库。

### 3. 依赖管理
Q: 如何处理依赖冲突？
A: 使用 `go mod tidy` 和 `go mod vendor` 管理依赖。

## 开发工具推荐

### 1. Proto 文件编辑
- VSCode + vscode-proto3: 提供语法高亮和自动完成
- GoLand: 内置 Proto 支持
- BloomRPC: 用于测试 gRPC 接口

### 2. 数据库工具
- MySQL Workbench: MySQL 官方工具
- Navicat: 支持多种数据库
- Redis Desktop Manager: Redis 可视化管理

## 联系方式

如有问题，请联系：
- DENGYI：[QQ:212294929, 微信:18375657303] 