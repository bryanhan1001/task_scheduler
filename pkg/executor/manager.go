package executor

import (
	"fmt"
	"sync"

	"task_scheduler/pkg/types"
)

// Manager 执行器管理器
type Manager struct {
	executors map[string]types.Executor
	mutex     sync.RWMutex
}

// NewManager 创建执行器管理器
func NewManager() *Manager {
	return &Manager{
		executors: make(map[string]types.Executor),
	}
}

// AddExecutor 添加执行器
func (em *Manager) AddExecutor(executor types.Executor) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	if _, exists := em.executors[executor.GetID()]; exists {
		return fmt.Errorf("executor %s already exists", executor.GetID())
	}

	em.executors[executor.GetID()] = executor
	return nil
}

// RemoveExecutor 移除执行器
func (em *Manager) RemoveExecutor(executorID string) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	if _, exists := em.executors[executorID]; !exists {
		return fmt.Errorf("executor %s not found", executorID)
	}

	delete(em.executors, executorID)
	return nil
}

// GetExecutors 获取所有健康的执行器
func (em *Manager) GetExecutors() []types.Executor {
	em.mutex.RLock()
	defer em.mutex.RUnlock()

	var healthyExecutors []types.Executor
	for _, executor := range em.executors {
		if executor.IsHealthy() {
			healthyExecutors = append(healthyExecutors, executor)
		}
	}

	return healthyExecutors
}

// GetExecutor 根据ID获取执行器
func (em *Manager) GetExecutor(executorID string) (types.Executor, error) {
	em.mutex.RLock()
	defer em.mutex.RUnlock()

	executor, exists := em.executors[executorID]
	if !exists {
		return nil, fmt.Errorf("executor %s not found", executorID)
	}

	return executor, nil
}
