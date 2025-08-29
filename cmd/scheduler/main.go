package main

import (
	"fmt"
	"log"
	"time"

	"task_scheduler/pkg/executor"
	"task_scheduler/pkg/scheduler"
	"task_scheduler/pkg/types"
)

func main() {
	// 创建调度器配置
	config := &types.SchedulerConfig{
		MaxConcurrentTasks:  5,
		HealthCheckInterval: 30 * time.Second,
		DefaultStrategy:     types.RoundRobinApp,
	}

	// 创建任务调度器
	ts := scheduler.New(config)

	// 添加执行器
	executors := []*executor.SimpleExecutor{
		executor.NewSimpleExecutor("executor-1", "http://localhost:8001"),
		executor.NewSimpleExecutor("executor-2", "http://localhost:8002"),
		executor.NewSimpleExecutor("executor-3", "http://localhost:8003"),
	}

	for _, exec := range executors {
		err := ts.AddExecutor(exec)
		if err != nil {
			log.Fatalf("Failed to add executor %s: %v", exec.GetID(), err)
		}
		log.Printf("Added executor: %s", exec.GetID())
	}

	// 创建不同策略的任务
	tasks := []*types.Task{
		{
			ID:       "order-timeout-check",
			Name:     "订单超时检查",
			Cron:     "0 */1 * * * *", // 每分钟执行一次
			Handler:  "orderTimeoutHandler",
			Params:   map[string]interface{}{"timeout": 30},
			Strategy: types.RoundRobinTask, // 任务级别轮询
		},
		{
			ID:       "risk-monitoring",
			Name:     "风险监控",
			Cron:     "0 */1 * * * *", // 每分钟执行一次
			Handler:  "riskMonitoringHandler",
			Params:   map[string]interface{}{"threshold": 100},
			Strategy: types.Random, // 随机路由
		},
		{
			ID:       "data-sync",
			Name:     "数据同步",
			Cron:     "0 0 2 * * *", // 每天凌晨2点执行
			Handler:  "dataSyncHandler",
			Params:   map[string]interface{}{"tables": []string{"inventory", "stores"}},
			Strategy: types.LFU, // 最少使用优先
		},
		{
			ID:       "cache-cleanup",
			Name:     "缓存清理",
			Cron:     "0 */30 * * * *", // 每30分钟执行一次
			Handler:  "cacheCleanupHandler",
			Params:   map[string]interface{}{"max_age": 3600},
			Strategy: types.LRU, // 最久未使用优先
		},
		{
			ID:       "health-check",
			Name:     "健康检查",
			Cron:     "0 */5 * * * *", // 每5分钟执行一次
			Handler:  "healthCheckHandler",
			Params:   map[string]interface{}{"services": []string{"api", "db", "cache"}},
			Strategy: types.RoundRobinApp, // 应用级别轮询
		},
	}

	// 添加任务到调度器
	for _, task := range tasks {
		err := ts.AddTask(task)
		if err != nil {
			log.Fatalf("Failed to add task %s: %v", task.ID, err)
		}
		log.Printf("Added task: %s with strategy %v", task.Name, task.Strategy)
	}

	// 启动调度器
	err := ts.Start()
	if err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}

	log.Println("=== 任务调度器启动成功 ===")
	log.Println("演示不同的路由策略效果...")

	// 运行一段时间来观察调度效果
	time.Sleep(2 * time.Minute)

	// 打印统计信息
	printStats(ts)

	// 演示手动执行任务
	log.Println("\n=== 演示手动任务执行 ===")
	manualTask := &types.Task{
		ID:       "manual-task",
		Name:     "手动任务",
		Handler:  "manualHandler",
		Params:   map[string]interface{}{"type": "manual"},
		Strategy: types.Random,
	}

	// 手动执行任务几次，观察不同策略的效果
	for i := 0; i < 10; i++ {
		executors := ts.GetExecutors()
		if len(executors) > 0 {
			executor, err := ts.GetRouter().Route(manualTask, executors)
			if err == nil {
				log.Printf("Manual task %d routed to executor: %s", i+1, executor.GetID())
				executor.Execute(manualTask)
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	// 再次打印统计信息
	log.Println("\n=== 最终统计信息 ===")
	printStats(ts)

	// 停止调度器
	err = ts.Stop()
	if err != nil {
		log.Printf("Error stopping scheduler: %v", err)
	} else {
		log.Println("任务调度器已停止")
	}
}

// printStats 打印统计信息
func printStats(ts *scheduler.TaskScheduler) {
	log.Println("\n=== 任务统计信息 ===")
	taskStats := ts.GetTaskStats()
	for key, value := range taskStats {
		log.Printf("%s: %v", key, value)
	}

	log.Println("\n=== 执行器统计信息 ===")
	executorStats := ts.GetExecutorStats()
	for _, stat := range executorStats {
		log.Printf("Executor %s: 使用次数=%d, 最后使用时间=%v, 健康状态=%v",
			stat.ID, stat.UsageCount, stat.LastUsedTime.Format("15:04:05"), stat.IsHealthy)
	}

	printStrategyDescription()
}

// printStrategyDescription 打印策略说明
func printStrategyDescription() {
	log.Println("\n=== 路由策略说明 ===")
	fmt.Println(`
根据文章描述，不同路由策略的特点：

1. 任务级别轮询 (RoundRobinTask):
   - 每个任务维护独立计数器
   - 适合任务调度时间一致的场景
   - 初始化时随机，避免首次压力集中

2. 应用级别轮询 (RoundRobinApp):
   - 所有任务共享一个计数器
   - 保证所有执行器接收任务次数平均
   - 适合任务负载和执行时间相近的场景

3. 随机路由 (Random):
   - 完全随机选择执行器
   - 长期看负载均衡，短期可能不均
   - 实现简单，适合对均衡要求不严格的场景

4. 最近最少使用 (LFU):
   - 优先选择使用次数最少的执行器
   - 适合混合多种路由策略的场景
   - 能有效平衡不同策略导致的负载不均

5. 最近最久未使用 (LRU):
   - 优先选择最久未使用的执行器
   - 基于时间的负载均衡
   - 适合需要考虑时间因素的调度场景

建议：
- 任务负载相近：使用应用级别轮询
- 存在大小任务：大任务用任务级轮询，小任务用其他策略
- 混合策略场景：新任务使用LFU策略平衡负载`)
}
