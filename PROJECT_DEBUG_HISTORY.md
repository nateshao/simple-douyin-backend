# simple-douyin-backend 项目调试与修复全记录

## 1. 项目简介

simple-douyin-backend 是一个类抖音后端，采用微服务架构，gRPC、MySQL、Redis、Kafka 等多种组件。

---

## 2. 启动与调试全流程

### 2.1 依赖环境准备
- Go 1.18+
- Docker & Docker Compose

### 2.2 拉取依赖
```sh
go mod tidy
```

### 2.3 启动依赖服务
```sh
docker-compose up -d
```
- 等待 MySQL、Redis、Kafka 等服务 healthy

### 2.4 检查配置
- configs/ 下数据库、Redis、Kafka 配置需与 docker-compose 保持一致
- MySQL 默认 root/123456，端口 3306
- Redis 默认 6379，无密码
- Kafka 默认 9092

### 2.5 编译与启动服务
```sh
go build -o simple_douyin_backend ./cmd/simple_douyin_backend
go build -o service_level ./cmd/service_level
go build -o dao_level ./cmd/dao_level
```
分别启动：
```sh
./dao_level
./service_level
./simple_douyin_backend
```

---

## 3. 主要问题与修复过程

### 3.1 数据库连接失败
**日志：**
```
[error] failed to initialize database, got error Error 1045 (28000): Access denied for user ''@'localhost' (using password: NO)
```
**分析：** 数据库配置未正确读取，用户名/密码为空。
**修复：** 检查 configs/ 下数据库配置，确保与 docker-compose 保持一致。

---

### 3.2 Redis/Kafka 未初始化导致 panic
**日志：**
```
panic: runtime error: invalid memory address or nil pointer dereference
.../internal/service.(*userService).userRegisterInfo
```
**分析：** service_level 的 redisClient 未初始化。
**修复：**
- 在 service_level 启动时，确保调用 internal/service/rdb.go 的 initRedis()，初始化 redisClient。
- 在 userRegisterInfo 等函数开头加 redisClient nil 检查和日志，便于排查。

---

### 3.3 注册/登录 userId=0
**分析：**
- 注册成功后，JWT 登录流程通过 requestContext.Query("username") 获取参数，但注册接口为 POST，参数在 body，导致获取为空。
- gRPC 查询不到用户，userId=0。
**修复：**
- 注册成功后，将用户名和密码写入 query，确保 JWT 登录流程能正确获取参数。

---

### 3.4 Kafka 消费协程 panic
**日志：**
```
panic: runtime error: invalid memory address or nil pointer dereference
.../internal/dao.(*favoriteDao).getFromMessageQueue
```
**分析：** Kafka 消费协程依赖的 client/dao 未初始化。
**修复：**
- 检查 initKafkaClient 的初始化顺序，确保所有依赖初始化后再启动消费协程。

---

## 4. 验证主链路

### 4.1 服务健康检查
- 访问 http://127.0.0.1:8888/health
- 日志无 panic、无连接错误

### 4.2 注册/登录功能
```sh
curl -X POST "http://127.0.0.1:8888/douyin/user/register/" -d "username=testuser&password=testpass"
curl -X POST "http://127.0.0.1:8888/douyin/user/login/" -d "username=testuser&password=testpass"
```
- 返回应包含 user_id > 0 和 token

---

## 5. 关键日志片段

### 5.1 数据库初始化失败
```
[error] failed to initialize database, got error Error 1045 (28000): Access denied for user ''@'localhost' (using password: NO)
```

### 5.2 Redis 未初始化
```
redisClient is nil in userRegisterInfo!
```

### 5.3 Kafka 消费协程 panic
```
panic: runtime error: invalid memory address or nil pointer dereference
.../internal/dao.(*favoriteDao).getFromMessageQueue
```

### 5.4 注册/登录 userId=0
```
注册后 userId=0，原因：JWT 登录流程未获取到正确参数
```

---

## 6. 总结与建议

- 所有依赖服务建议用 docker-compose 启动
- 配置文件与实际服务端口/密码一致
- 所有服务都要编译并分别启动
- service_level 的 redisClient/kafkaClient 必须在服务启动时初始化
- 注册/登录链路已修复，userId 不应为 0

如需进一步排查，请贴出接口响应和最新日志。 