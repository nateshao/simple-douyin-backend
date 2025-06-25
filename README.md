# simple-douyin-backend

[![build](https://img.shields.io/badge/build-0.1.0-brightgreen)](https://github.com/StellarisW/douyin)
[![go-version](https://img.shields.io/badge/go-~%3D1.18-30dff3?logo=go)](https://github.com/StellarisW/douyin)
[![OpenTracing Badge](https://img.shields.io/badge/OpenTracing-disabled-blue.svg)](http://opentracing.io)

> 极简版抖音后端，微服务架构，gRPC 通信，MySQL/Redis/Kafka 支撑，支持高并发与高可用，适合学习与团队协作。

---

## 🚀 快速开始

### 1. 环境准备

- Go 1.18+
- Docker & Docker Compose

### 2. 拉取依赖

```sh
go mod tidy
```

### 3. 启动依赖服务

```sh
docker-compose up -d
```
> 等待 MySQL、Redis、Kafka 等服务 healthy

### 4. 编译与启动服务

```sh
go build -o simple_douyin_backend ./cmd/simple_douyin_backend
go build -o service_level ./cmd/service_level
go build -o dao_level ./cmd/dao_level
```

分别启动（建议分终端窗口）：

```sh
./dao_level
./service_level
./simple_douyin_backend
```

### 5. 功能验证

- 访问健康检查接口: [http://127.0.0.1:8888/health](http://127.0.0.1:8888/health)
- 注册/登录示例：

```sh
curl -X POST "http://127.0.0.1:8888/douyin/user/register/" -d "username=testuser&password=testpass"
curl -X POST "http://127.0.0.1:8888/douyin/user/login/" -d "username=testuser&password=testpass"
```

---

## 🏗️ 项目架构与特色

### 架构亮点

- **微服务分层**：基于 DDD，分为 Controller / Service / Dao 层，领域服务独立，gRPC 通信。
- **高可用**：限流、熔断、降级、消息队列削峰，支持服务注册与发现。
- **高性能**：Redis 缓存、Kafka 异步、协程并发、算法优化。
- **安全性**：JWT 鉴权、密码加密、边界校验。
- **易维护**：全局错误码、日志、自动化测试、代码生成工具。

### 技术栈

| 层级      | 技术/工具         | 说明                   |
|---------|----------------|----------------------|
| HTTP    | Hertz          | 高性能 HTTP 框架         |
| RPC     | gRPC           | 高性能 RPC 框架          |
| ORM     | Gorm           | 主流 ORM，防 SQL 注入     |
| 缓存     | Redis          | 高速缓存                |
| 消息队列   | Kafka          | 异步削峰                |
| 配置中心   | Apollo         | 配置管理（可选）          |
| 日志     | Zlog           | 日志全覆盖               |
| 服务发现   | Etcd           | 服务注册与发现            |
| 追踪     | Jaeger         | 分布式链路追踪（可选）      |
| 代理     | Nginx          | 反向代理与限流            |
| 其他     | bcrypt, JWT, Docker, Kubernetes, GoConvey, SnowFlake, CI/CD |

---

## 📁 目录结构

```
.
├── api/           # gRPC/HTTP 接口定义
├── cmd/           # 各服务启动入口
├── configs/       # 配置文件
├── deploy/        # 部署脚本与说明
├── docs/          # 设计文档与总结
├── init/          # 初始化相关
├── internal/      # 业务核心代码（Controller/Service/Dao）
├── test/          # 测试用例
├── third_party/   # 第三方依赖
├── docker-compose.yml
├── Dockerfile
└── README.md
```

- 详细目录说明见各目录下 `package_info` 文件。

---

## 📚 设计与开发规范

- 遵循 [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
- Controller/Service/Dao 三层分明，领域服务独立
- 详见 `/docs` 目录下设计文档

---

## 📝 数据库与安全规范

- 建表规范、字段注释、软删除字段等见 `/docs`
- 密码加密存储，JWT 鉴权，防 SQL 注入
- 参考 [腾讯 Go 安全指南](https://github.com/Tencent/secguide/blob/main/Go%E5%AE%89%E5%85%A8%E6%8C%87%E5%8D%97.md)

---

## 🧩 贡献与协作

- 欢迎 issue、PR、代码 review
- 建议先阅读 `/docs` 下的设计与开发规范
- 单元测试、集成测试、压力测试已覆盖主要链路

---

## 🔗 更多信息

- 详细设计、表结构、缓存、消息队列、安全性等见 `/docs`
- 优化与调试历史见 `PROJECT_DEBUG_HISTORY.md`
- TODO 任务追踪见 `TODO.md`

---

## 📷 架构图

![架构图](/docs/images/structure.png)

---

## License

MIT
