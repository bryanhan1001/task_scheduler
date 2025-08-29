package main

import (
	"log"
	"time"

	"task_scheduler/pkg/executor"
	"task_scheduler/pkg/scheduler"
	"task_scheduler/pkg/types"
)

// 简单示例：演示基本的任务调度功能
func main() {
	// 创建调度器
	config := &types.SchedulerConfig{
		MaxConcurrentTasks:  3,
		HealthCheckInterval: 10 * time.Second,
		DefaultStrategy:     types.RoundRobinApp,
	}

	ts := scheduler.New(config)

	// 添加执行器
	exec1 := executor.NewSimpleExecutor("worker-1", "http://localhost:8001")
	exec2 := executor.NewSimpleExecutor("worker-2", "http://localhost:8002")

	ts.AddExecutor(exec1)
	ts.AddExecutor(exec2)

	// 创建简单任务
	task := &types.Task{
		ID:       "hello-world",
		Name:     "Hello World Task",
		Cron:     "0 */10 * * * *", // 每10秒执行一次
		Handler:  "helloWorldHandler",
		Params:   map[string]interface{}{"message": "Hello, World!"},
		Strategy: types.RoundRobinApp,
	}

	// 添加任务
	err := ts.AddTask(task)
	if err != nil {
		log.Fatalf("Failed to add task: %v", err)
	}

	// 启动调度器
	err = ts.Start()
	if err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}

	log.Println("Simple scheduler started. Running for 1 minute...")

	// 运行1分钟
	time.Sleep(1 * time.Minute)

	// 停止调度器
	ts.Stop()
	log.Println("Scheduler stopped.")
}
