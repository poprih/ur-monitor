ur-monitor/
├── cmd/
│ ├── server/ # 主要应用入口
│ │ ├── main.go # 初始化 web 服务器和定时任务
│ │ ├── config.go # 配置加载
│ │ ├── routes.go # 注册 API 路由
│ │ ├── jobs.go # 定时任务调度
│ │ ├── line.go # 处理 LINE 事件
│ │ ├── logger.go # 日志初始化
│ │ ├── database.go # 数据库连接
│ ├── worker/ # 独立 Worker（可选，如果需要更好的并发）
│ │ ├── main.go # 监听队列任务
│
├── internal/
│ ├── controllers/ # 控制器，处理 HTTP 请求
│ │ ├── line_controller.go
│ │ ├── monitor_controller.go
│ ├── services/ # 业务逻辑
│ │ ├── line_service.go
│ │ ├── monitor_service.go
│ │ ├── notification_service.go
│ ├── repositories/ # 数据访问层（数据库操作）
│ │ ├── user_repository.go
│ │ ├── subscription_repository.go
│ ├── models/ # 数据模型
│ │ ├── user.go
│ │ ├── subscription.go
│ ├── jobs/ # 定时任务逻辑
│ │ ├── monitor_job.go
│ ├── clients/ # 外部 API 客户端
│ │ ├── line_client.go
│ │ ├── ur_client.go
│
├── pkg/ # 公共工具库
│ ├── utils.go
│ ├── line_webhook.go
│ ├── scheduler.go
│ ├── http_client.go
│
├── configs/ # 配置文件
│ ├── config.yaml
│
├── migrations/ # 数据库迁移文件
│ ├── 001_create_users.sql
│ ├── 002_create_subscriptions.sql
│
├── test/ # 测试代码
│ ├── line_service_test.go
│ ├── monitor_service_test.go
│
├── go.mod # Go 依赖管理
├── go.sum
└── README.md
