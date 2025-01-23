# Deploy 目录说明

本目录包含了 TikTokMall 项目的所有基础设施和部署相关的配置文件。

## 目录结构

```bash
deploy/
├── docker/                     # Docker 相关配置
│   ├── docker-compose.yaml    # 基础服务编排配置
│   ├── mysql/                # MySQL 相关配置
│   │   └── init/            # MySQL 初始化脚本
│   │       └── init.sql     # 数据库初始化 SQL
│   ├── prometheus/          # Prometheus 监控配置
│   │   └── prometheus.yml   # Prometheus 配置文件
│   └── logstash/           # Logstash 日志收集配置
│       ├── config/         # Logstash 主配置
│       │   └── logstash.yml
│       └── pipeline/       # Logstash 管道配置
│           └── logstash.conf
```

## 基础服务说明

### 1. MySQL (端口: 3306)
- 主数据库，存储所有服务的业务数据
- 默认创建数据库: `tiktok_mall`
- 默认用户: `tiktok`
- 默认密码: `tiktok123`
- 包含的表:
  - users: 用户信息
  - products: 商品信息
  - cart_items: 购物车项
  - orders: 订单信息
  - order_items: 订单项
  - payments: 支付信息

### 2. Redis (端口: 6379)
- 缓存服务
- 用途：
  - 会话管理
  - 数据缓存
  - 限流计数

### 3. Consul (端口: 8500)
- 服务注册与发现中心
- 功能：
  - 服务注册
  - 健康检查
  - 配置中心
- 访问 UI: http://localhost:8500

### 4. Jaeger (端口: 16686)
- 分布式链路追踪系统
- 功能：
  - 请求链路追踪
  - 性能分析
  - 问题定位
- 访问 UI: http://localhost:16686

### 5. Prometheus & Grafana (端口: 9090, 3000)
- 监控系统
- Prometheus (9090):
  - 指标收集
  - 时序数据存储
  - 告警规则配置
- Grafana (3000):
  - 数据可视化
  - 默认账号: admin
  - 默认密码: admin123
- 访问地址:
  - Prometheus: http://localhost:9090
  - Grafana: http://localhost:3000

### 6. ELK Stack
- 日志收集与分析系统
- Elasticsearch (端口: 9200, 9300):
  - 日志存储
  - 全文检索
- Kibana (端口: 5601):
  - 日志可视化
  - 数据分析
- Logstash (端口: 5044, 5000):
  - 日志收集
  - 日志解析
  - 日志转发
- 访问地址:
  - Kibana: http://localhost:5601

## 使用说明

### 启动所有服务
```bash
cd deploy/docker
docker-compose up -d
```

### 查看服务状态
```bash
docker-compose ps
```

### 查看服务日志
```bash
# 查看所有服务日志
docker-compose logs

# 查看特定服务日志
docker-compose logs [service_name]
```

### 停止服务
```bash
docker-compose down
```

### 重启服务
```bash
docker-compose restart [service_name]
```

## 配置说明

### MySQL 配置
- 初始化脚本位置: `mysql/init/init.sql`
- 包含所有数据库表结构
- 添加新表需要在此文件中添加相应的 SQL

### Prometheus 配置
- 配置文件位置: `prometheus/prometheus.yml`
- 已配置所有微服务的监控目标
- 支持 Consul 服务发现

### Logstash 配置
- 主配置: `logstash/config/logstash.yml`
- 管道配置: `logstash/pipeline/logstash.conf`
- 支持多种日志输入源
- 统一输出到 Elasticsearch

## 注意事项

1. 数据持久化
- 所有服务的数据都通过 Docker volumes 持久化
- 位置: `/var/lib/docker/volumes/`

2. 网络配置
- 所有服务都在 `tiktok_mall_net` 网络中
- 服务间可以通过服务名互相访问

3. 安全性
- 所有服务都配置了默认密码
- 生产环境部署时需要修改默认密码
- 建议配置防火墙规则

4. 资源需求
- 建议最小配置:
  - CPU: 4 核
  - 内存: 8GB
  - 磁盘: 50GB

## 常见问题

1. 服务无法启动
- 检查端口占用
- 检查资源使用情况
- 查看具体服务日志

2. 数据持久化问题
- 确保 volumes 目录有足够权限
- 定期备份重要数据

3. 性能问题
- 适当调整各服务的资源限制
- 监控服务资源使用情况
- 根据需要扩展配置 