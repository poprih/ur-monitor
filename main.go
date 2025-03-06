package main

import (
	"log"

	"github.com/poprih/ur-monitor/cmd/server"
)

func main() {
	log.Println("Starting UR Monitor Service...")

	// 加载配置
	// server.LoadConfig()

	// 初始化数据库
	// server.InitDB()

	// 启动 Web 服务器
	go server.StartServer()

	// 启动定时任务
	// danchiList := os.Getenv("DANCHI_LIST")
	// lineToken := os.Getenv("LINE_NOTIFY_TOKEN")
	// if danchiList == "" || lineToken == "" {
	// 	log.Fatal("DANCHI_LIST and LINE_NOTIFY_TOKEN environment variables are required")
	// }
	// danchis := strings.Split(danchiList, ",")
	// go jobs.StartMonitorJob(danchis, lineToken)

	select {} // Keep the main goroutine running
}

// 本地部署步骤：
// 1. 设置环境变量 DANCHI_LIST 和 LINE_NOTIFY_TOKEN。
// 2. 确保数据库已初始化并可用。
// 3. 运行 `go run main.go` 启动服务。
