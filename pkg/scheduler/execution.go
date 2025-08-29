package scheduler

import (
	"fmt"
	"log"
	"time"

	"task_scheduler/pkg/types"
)

// Start 启动调度器
func (ts *TaskScheduler) Start() error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	if ts.running {
		return fmt.Errorf("scheduler is already running")
	}

	// 启动cron调度器
	ts.cron.Start()

	// 启动健康检查
	go ts.healthCheckLoop()

	ts.running = true
	log.Println("Task scheduler started successfully")
	return nil
}

// Stop 停止调度器
func (ts *TaskScheduler) Stop() error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	if !ts.running {
		return fmt.Errorf("scheduler is not running")
	}

	// 停止cron调度器
	ts.cron.Stop()

	// 取消上下文，停止所有goroutine
	ts.cancel()

	// 停止所有正在运行的任务
	ts.taskMutex.Lock()
	for _, task := range ts.tasks {
		if task.Status == types.TaskStatusRunning {
			task.Status = types.TaskStatusStopped
		}
	}
	ts.taskMutex.Unlock()

	ts.running = false
	log.Println("Task scheduler stopped successfully")
	return nil
}

// executeTask 执行任务
func (ts *TaskScheduler) executeTask(task *types.Task) {
	// 检查任务状态
	if task.Status == types.TaskStatusRunning {
		log.Printf("Task %s is already running, skipping", task.ID)
		return
	}

	if task.Status == types.TaskStatusStopped {
		log.Printf("Task %s is stopped, skipping", task.ID)
		return
	}

	// 获取可用执行器
	executors := ts.executorManager.GetExecutors()
	if len(executors) == 0 {
		log.Printf("No available executors for task %s", task.ID)
		task.Status = types.TaskStatusFailed
		return
	}

	// 使用路由策略选择执行器
	executor, err := ts.router.Route(task, executors)
	if err != nil {
		log.Printf("Failed to route task %s: %v", task.ID, err)
		task.Status = types.TaskStatusFailed
		return
	}

	// 异步执行任务
	go func() {
		// 更新任务状态
		task.Status = types.TaskStatusRunning
		task.LastRunTime = time.Now()

		log.Printf("Executing task %s on executor %s (strategy: %v)",
			task.ID, executor.GetID(), task.Strategy)

		// 执行任务
		err := executor.Execute(task)
		if err != nil {
			log.Printf("Task %s execution failed: %v", task.ID, err)
			task.Status = types.TaskStatusFailed
		} else {
			log.Printf("Task %s executed successfully", task.ID)
			task.Status = types.TaskStatusCompleted
		}

		// 计算下次运行时间（如果是周期性任务）
		if task.Cron != "" {
			// 这里可以添加计算下次运行时间的逻辑
			task.Status = types.TaskStatusPending // 重置为待执行状态
		}
	}()
}

// healthCheckLoop 健康检查循环
func (ts *TaskScheduler) healthCheckLoop() {
	ticker := time.NewTicker(ts.config.HealthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ts.ctx.Done():
			return
		case <-ticker.C:
			ts.performHealthCheck()
		}
	}
}

// performHealthCheck 执行健康检查
func (ts *TaskScheduler) performHealthCheck() {
	executors := ts.executorManager.GetExecutors()
	for _, exec := range executors {
		go func(executor types.Executor) {
			// 这里可以添加实际的健康检查逻辑
			// 比如发送ping请求、检查响应时间等
			// 模拟健康检查
			// 实际应用中可以发送HTTP请求或RPC调用
			// 这里可以通过接口方法来设置健康状态
			// healthy := true
		}(exec)
	}
}
