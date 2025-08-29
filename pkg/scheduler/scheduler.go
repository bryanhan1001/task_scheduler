package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	"task_scheduler/pkg/executor"
	"task_scheduler/pkg/router"
	"task_scheduler/pkg/types"
)

// TaskScheduler 任务调度器实现
type TaskScheduler struct {
	tasks           map[string]*types.Task
	executorManager *executor.Manager
	router          *router.MultiStrategyRouter
	cron            *cron.Cron
	config          *types.SchedulerConfig
	running         bool
	ctx             context.Context
	cancel          context.CancelFunc
	mutex           sync.RWMutex
	taskMutex       sync.RWMutex
}

// New 创建新的任务调度器
func New(config *types.SchedulerConfig) *TaskScheduler {
	if config == nil {
		config = &types.SchedulerConfig{
			MaxConcurrentTasks:  10,
			HealthCheckInterval: 30 * time.Second,
			DefaultStrategy:     types.RoundRobinApp,
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &TaskScheduler{
		tasks:           make(map[string]*types.Task),
		executorManager: executor.NewManager(),
		router:          router.NewMultiStrategyRouter(),
		cron:            cron.New(cron.WithSeconds()),
		config:          config,
		ctx:             ctx,
		cancel:          cancel,
	}
}

// AddTask 添加任务
func (ts *TaskScheduler) AddTask(task *types.Task) error {
	ts.taskMutex.Lock()
	defer ts.taskMutex.Unlock()

	if _, exists := ts.tasks[task.ID]; exists {
		return fmt.Errorf("task %s already exists", task.ID)
	}

	// 设置默认策略
	if task.Strategy == 0 {
		task.Strategy = ts.config.DefaultStrategy
	}

	// 设置任务状态
	task.Status = types.TaskStatusPending
	task.CreatedAt = time.Now()

	// 添加到cron调度器
	if task.Cron != "" {
		entryID, err := ts.cron.AddFunc(task.Cron, func() {
			ts.executeTask(task)
		})
		if err != nil {
			return fmt.Errorf("failed to add cron job for task %s: %v", task.ID, err)
		}

		// 可以将entryID存储到task中，用于后续管理
		_ = entryID
	}

	ts.tasks[task.ID] = task
	log.Printf("Task %s added successfully with strategy %v", task.ID, task.Strategy)
	return nil
}

// RemoveTask 移除任务
func (ts *TaskScheduler) RemoveTask(taskID string) error {
	ts.taskMutex.Lock()
	defer ts.taskMutex.Unlock()

	task, exists := ts.tasks[taskID]
	if !exists {
		return fmt.Errorf("task %s not found", taskID)
	}

	// 停止正在运行的任务
	if task.Status == types.TaskStatusRunning {
		task.Status = types.TaskStatusStopped
	}

	delete(ts.tasks, taskID)
	log.Printf("Task %s removed successfully", taskID)
	return nil
}

// GetTasks 获取所有任务
func (ts *TaskScheduler) GetTasks() []*types.Task {
	ts.taskMutex.RLock()
	defer ts.taskMutex.RUnlock()

	tasks := make([]*types.Task, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// AddExecutor 添加执行器
func (ts *TaskScheduler) AddExecutor(executor types.Executor) error {
	return ts.executorManager.AddExecutor(executor)
}

// RemoveExecutor 移除执行器
func (ts *TaskScheduler) RemoveExecutor(executorID string) error {
	return ts.executorManager.RemoveExecutor(executorID)
}

// GetExecutors 获取所有执行器
func (ts *TaskScheduler) GetExecutors() []types.Executor {
	return ts.executorManager.GetExecutors()
}
