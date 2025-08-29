package router

import (
	"errors"
	"sync"
	"time"

	"task_scheduler/pkg/types"
)

// LFURouter 最近最少使用路由器
type LFURouter struct {
	BaseRouter
	taskCounters map[string]map[string]int64 // taskID -> executorID -> count
	mutex        sync.RWMutex
}

// NewLFURouter 创建LFU路由器
func NewLFURouter() *LFURouter {
	return &LFURouter{
		BaseRouter:   BaseRouter{strategy: types.LFU},
		taskCounters: make(map[string]map[string]int64),
	}
}

// Route LFU路由
func (r *LFURouter) Route(task *types.Task, executors []types.Executor) (types.Executor, error) {
	if len(executors) == 0 {
		return nil, errors.New("no available executors")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 初始化任务计数器
	if _, exists := r.taskCounters[task.ID]; !exists {
		r.taskCounters[task.ID] = make(map[string]int64)
	}

	taskCounter := r.taskCounters[task.ID]

	// 找到使用次数最少的执行器
	var selectedExecutor types.Executor
	minCount := int64(-1)

	for _, executor := range executors {
		executorID := executor.GetID()
		count, exists := taskCounter[executorID]
		if !exists {
			count = 0
		}

		if minCount == -1 || count < minCount {
			minCount = count
			selectedExecutor = executor
		}
	}

	// 增加选中执行器的计数
	if selectedExecutor != nil {
		taskCounter[selectedExecutor.GetID()]++
	}

	return selectedExecutor, nil
}

// LRURouter 最近最久未使用路由器
type LRURouter struct {
	BaseRouter
	taskLastUsed map[string]map[string]time.Time // taskID -> executorID -> lastUsedTime
	mutex        sync.RWMutex
}

// NewLRURouter 创建LRU路由器
func NewLRURouter() *LRURouter {
	return &LRURouter{
		BaseRouter:   BaseRouter{strategy: types.LRU},
		taskLastUsed: make(map[string]map[string]time.Time),
	}
}

// Route LRU路由
func (r *LRURouter) Route(task *types.Task, executors []types.Executor) (types.Executor, error) {
	if len(executors) == 0 {
		return nil, errors.New("no available executors")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 初始化任务最后使用时间记录
	if _, exists := r.taskLastUsed[task.ID]; !exists {
		r.taskLastUsed[task.ID] = make(map[string]time.Time)
	}

	taskUsage := r.taskLastUsed[task.ID]
	now := time.Now()

	// 找到最久未使用的执行器
	var selectedExecutor types.Executor
	oldestTime := now

	for _, executor := range executors {
		executorID := executor.GetID()
		lastUsed, exists := taskUsage[executorID]
		if !exists {
			// 如果从未使用过，直接选择
			selectedExecutor = executor
			break
		}

		if lastUsed.Before(oldestTime) {
			oldestTime = lastUsed
			selectedExecutor = executor
		}
	}

	// 更新选中执行器的最后使用时间
	if selectedExecutor != nil {
		taskUsage[selectedExecutor.GetID()] = now
	}

	return selectedExecutor, nil
}
