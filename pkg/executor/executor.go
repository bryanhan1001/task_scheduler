package executor

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"task_scheduler/pkg/types"
)

// SimpleExecutor 简单执行器实现
type SimpleExecutor struct {
	id         string
	address    string
	usageCount int64
	lastUsed   time.Time
	isHealthy  bool
	mutex      sync.RWMutex
}

// NewSimpleExecutor 创建新的简单执行器
func NewSimpleExecutor(id, address string) *SimpleExecutor {
	return &SimpleExecutor{
		id:        id,
		address:   address,
		isHealthy: true,
		lastUsed:  time.Now(),
	}
}

// GetID 获取执行器ID
func (e *SimpleExecutor) GetID() string {
	return e.id
}

// GetAddress 获取执行器地址
func (e *SimpleExecutor) GetAddress() string {
	return e.address
}

// IsHealthy 检查执行器是否健康
func (e *SimpleExecutor) IsHealthy() bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.isHealthy
}

// SetHealthy 设置执行器健康状态
func (e *SimpleExecutor) SetHealthy(healthy bool) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.isHealthy = healthy
}

// Execute 执行任务
func (e *SimpleExecutor) Execute(task *types.Task) error {
	if !e.IsHealthy() {
		return fmt.Errorf("executor %s is not healthy", e.id)
	}

	// 更新使用统计
	e.IncrementUsage()
	e.updateLastUsedTime()

	// 模拟任务执行
	fmt.Printf("Executor %s executing task %s (handler: %s)\n", e.id, task.ID, task.Handler)

	// 这里可以添加实际的任务执行逻辑
	// 比如HTTP调用、RPC调用等

	return nil
}

// GetLastUsedTime 获取最后使用时间
func (e *SimpleExecutor) GetLastUsedTime() time.Time {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.lastUsed
}

// GetUsageCount 获取使用次数
func (e *SimpleExecutor) GetUsageCount() int64 {
	return atomic.LoadInt64(&e.usageCount)
}

// IncrementUsage 增加使用次数
func (e *SimpleExecutor) IncrementUsage() {
	atomic.AddInt64(&e.usageCount, 1)
}

// updateLastUsedTime 更新最后使用时间
func (e *SimpleExecutor) updateLastUsedTime() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.lastUsed = time.Now()
}

// GetStats 获取执行器统计信息
func (e *SimpleExecutor) GetStats() *types.ExecutorStats {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	return &types.ExecutorStats{
		ID:           e.id,
		Address:      e.address,
		UsageCount:   atomic.LoadInt64(&e.usageCount),
		LastUsedTime: e.lastUsed,
		IsHealthy:    e.isHealthy,
	}
}
