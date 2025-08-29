package scheduler

import (
	"task_scheduler/pkg/executor"
	"task_scheduler/pkg/router"
	"task_scheduler/pkg/types"
)

// GetTaskStats 获取任务统计信息
func (ts *TaskScheduler) GetTaskStats() map[string]interface{} {
	ts.taskMutex.RLock()
	defer ts.taskMutex.RUnlock()

	stats := make(map[string]interface{})
	statusCount := make(map[types.TaskStatus]int)
	strategyCount := make(map[types.RouteStrategy]int)

	for _, task := range ts.tasks {
		statusCount[task.Status]++
		strategyCount[task.Strategy]++
	}

	stats["total_tasks"] = len(ts.tasks)
	stats["status_distribution"] = statusCount
	stats["strategy_distribution"] = strategyCount
	stats["total_executors"] = len(ts.executorManager.GetExecutors())

	return stats
}

// GetExecutorStats 获取执行器统计信息
func (ts *TaskScheduler) GetExecutorStats() []*types.ExecutorStats {
	executors := ts.executorManager.GetExecutors()
	stats := make([]*types.ExecutorStats, 0, len(executors))

	for _, exec := range executors {
		if simpleExec, ok := exec.(*executor.SimpleExecutor); ok {
			stats = append(stats, simpleExec.GetStats())
		}
	}

	return stats
}

// GetRouter 获取路由器（用于手动路由）
func (ts *TaskScheduler) GetRouter() *router.MultiStrategyRouter {
	return ts.router
}
