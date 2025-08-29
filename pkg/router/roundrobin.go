package router

import (
	"errors"
	"math/rand"
	"sync"
	"sync/atomic"

	"task_scheduler/pkg/types"
)

// RoundRobinTaskRouter 任务级别轮询路由器
type RoundRobinTaskRouter struct {
	BaseRouter
	taskCounters map[string]*int64 // 每个任务的计数器
	mutex        sync.RWMutex
}

// NewRoundRobinTaskRouter 创建任务级别轮询路由器
func NewRoundRobinTaskRouter() *RoundRobinTaskRouter {
	return &RoundRobinTaskRouter{
		BaseRouter:   BaseRouter{strategy: types.RoundRobinTask},
		taskCounters: make(map[string]*int64),
	}
}

// Route 任务级别轮询路由
func (r *RoundRobinTaskRouter) Route(task *types.Task, executors []types.Executor) (types.Executor, error) {
	if len(executors) == 0 {
		return nil, errors.New("no available executors")
	}

	r.mutex.Lock()
	counter, exists := r.taskCounters[task.ID]
	if !exists {
		// 初始化时随机一次，缓解首次压力
		initialValue := int64(rand.Intn(100))
		counter = &initialValue
		r.taskCounters[task.ID] = counter
	}
	r.mutex.Unlock()

	// 原子操作增加计数器
	count := atomic.AddInt64(counter, 1)
	if count > 1000000 {
		// 重置计数器，避免溢出
		atomic.StoreInt64(counter, int64(rand.Intn(100)))
		count = atomic.LoadInt64(counter)
	}

	index := int(count) % len(executors)
	return executors[index], nil
}

// RoundRobinAppRouter 应用级别轮询路由器
type RoundRobinAppRouter struct {
	BaseRouter
	globalCounter int64
}

// NewRoundRobinAppRouter 创建应用级别轮询路由器
func NewRoundRobinAppRouter() *RoundRobinAppRouter {
	return &RoundRobinAppRouter{
		BaseRouter: BaseRouter{strategy: types.RoundRobinApp},
	}
}

// Route 应用级别轮询路由
func (r *RoundRobinAppRouter) Route(task *types.Task, executors []types.Executor) (types.Executor, error) {
	if len(executors) == 0 {
		return nil, errors.New("no available executors")
	}

	count := atomic.AddInt64(&r.globalCounter, 1)
	if count > 1000000 {
		// 重置计数器，避免溢出
		atomic.StoreInt64(&r.globalCounter, 0)
		count = 0
	}

	index := int(count) % len(executors)
	return executors[index], nil
}
