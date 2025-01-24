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

### 首次使用
1. 进入 docker 目录：
```bash
cd deploy/docker
```

2. 启动所有基础服务：
```bash
docker-compose up -d
```

3. 验证服务状态：
```bash
docker-compose ps
```

### 服务访问

1. **Web 界面访问**
- Consul UI: http://localhost:8500
- Jaeger UI: http://localhost:16686
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000
  - 默认账号: admin
  - 默认密码: admin123
- Kibana: http://localhost:5601

2. **数据库连接信息**
```
MySQL:
- 主机: localhost:3306
- 数据库: tiktok_mall
- 用户名: tiktok
- 密码: tiktok123

Redis:
- 主机: localhost:6379
- 无密码
```

### 服务管理命令

1. **查看服务状态**
```bash
# 查看所有服务状态
docker-compose ps

# 查看特定服务状态
docker-compose ps [service_name]
```

2. **查看服务日志**
```bash
# 查看所有服务日志
docker-compose logs

# 查看特定服务日志
docker-compose logs [service_name]

# 实时查看日志
docker-compose logs -f [service_name]
```

3. **服务生命周期管理**
```bash
# 启动所有服务
docker-compose up -d

# 重启特定服务
docker-compose restart [service_name]

# 停止所有服务
docker-compose down

# 停止并删除所有数据（包括数据卷）
docker-compose down -v
```

### 配置管理

1. **MySQL 配置**
```bash
# 编辑初始化脚本
vim mysql/init/init.sql

# 应用更改（需要重新创建容器）
docker-compose down
docker-compose up -d
```

2. **Prometheus 配置**
```bash
# 编辑配置文件
vim prometheus/prometheus.yml

# 重启 Prometheus 服务
docker-compose restart prometheus
```

3. **Logstash 配置**
```bash
# 编辑管道配置
vim logstash/pipeline/logstash.conf

# 重启 Logstash 服务
docker-compose restart logstash
```

### 数据持久化

所有服务数据都通过 Docker volumes 持久化，包括：
- mysql_data: MySQL 数据文件
- redis_data: Redis 数据文件
- consul_data: Consul 数据
- prometheus_data: 监控数据
- grafana_data: 仪表盘配置
- elasticsearch_data: 日志数据

数据卷位置：`/var/lib/docker/volumes/`

### 故障排查指南

1. **服务无法启动**
```bash
# 1. 检查服务状态
docker-compose ps

# 2. 查看详细日志
docker-compose logs -f [service_name]

# 3. 检查端口占用
netstat -tunlp | grep [port]

# 4. 重新创建服务
docker-compose up -d --force-recreate [service_name]
```

2. **网络问题**
```bash
# 检查网络连接
docker network inspect tiktok_mall_net

# 重新创建网络
docker-compose down
docker-compose up -d
```

3. **数据持久化问题**
```bash
# 检查数据卷
docker volume ls

# 检查数据卷权限
ls -la /var/lib/docker/volumes/

# 备份数据
docker run --rm -v [volume_name]:/source -v $(pwd):/backup alpine tar -czf /backup/backup.tar.gz -C /source .
```

## 性能优化建议

1. **系统要求**
- 最小配置：
  - CPU: 4 核
  - 内存: 8GB
  - 磁盘: 50GB

2. **性能调优**
- 适当调整各服务的资源限制
- 监控服务资源使用情况
- 根据需要扩展配置

3. **注意事项**
- 生产环境部署前修改默认密码
- 配置适当的防火墙规则
- 定期备份重要数据
- 监控系统资源使用情况

## 安全建议

1. **密码管理**
- 修改所有默认密码
- 使用强密码策略
- 定期轮换密码

2. **网络安全**
- 限制服务访问范围
- 配置防火墙规则
- 使用 HTTPS 进行通信

3. **数据安全**
- 定期备份数据
- 加密敏感信息
- 实施访问控制

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