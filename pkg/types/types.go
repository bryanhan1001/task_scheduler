package types

import (
	"sync"
	"time"
)

// RouteStrategy 路由策略类型
type RouteStrategy int

const (
	// RoundRobinTask 任务级别轮询
	RoundRobinTask RouteStrategy = iota
	// RoundRobinApp 应用级别轮询
	RoundRobinApp
	// Random 随机路由
	Random
	// LFU 最近最少使用
	LFU
	// LRU 最近最久未使用
	LRU
)

// Executor 执行器接口
type Executor interface {
	GetID() string
	GetAddress() string
	IsHealthy() bool
	Execute(task *Task) error
	GetLastUsedTime() time.Time
	GetUsageCount() int64
	IncrementUsage()
}

// Task 任务定义
type Task struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Cron        string        `json:"cron"`
	Handler     string        `json:"handler"`
	Params      interface{}   `json:"params"`
	Strategy    RouteStrategy `json:"strategy"`
	CreatedAt   time.Time     `json:"created_at"`
	LastRunTime time.Time     `json:"last_run_time"`
	NextRunTime time.Time     `json:"next_run_time"`
	Status      TaskStatus    `json:"status"`
}

// TaskStatus 任务状态
type TaskStatus int

const (
	TaskStatusPending TaskStatus = iota
	TaskStatusRunning
	TaskStatusCompleted
	TaskStatusFailed
	TaskStatusStopped
)

// Router 路由器接口
type Router interface {
	Route(task *Task, executors []Executor) (Executor, error)
	GetStrategy() RouteStrategy
}

// Scheduler 调度器接口
type Scheduler interface {
	AddTask(task *Task) error
	RemoveTask(taskID string) error
	Start() error
	Stop() error
	GetTasks() []*Task
	AddExecutor(executor Executor) error
	RemoveExecutor(executorID string) error
	GetExecutors() []Executor
}

// ExecutorStats 执行器统计信息
type ExecutorStats struct {
	ID           string    `json:"id"`
	Address      string    `json:"address"`
	UsageCount   int64     `json:"usage_count"`
	LastUsedTime time.Time `json:"last_used_time"`
	IsHealthy    bool      `json:"is_healthy"`
	mutex        sync.RWMutex
}

// SchedulerConfig 调度器配置
type SchedulerConfig struct {
	MaxConcurrentTasks  int           `json:"max_concurrent_tasks"`
	HealthCheckInterval time.Duration `json:"health_check_interval"`
	DefaultStrategy     RouteStrategy `json:"default_strategy"`
}
