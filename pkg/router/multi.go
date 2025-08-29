package router

import (
	"sync"

	"task_scheduler/pkg/types"
)

// MultiStrategyRouter 多策略路由器管理
type MultiStrategyRouter struct {
	routers map[types.RouteStrategy]types.Router
	factory *Factory
	mutex   sync.RWMutex
}

// NewMultiStrategyRouter 创建多策略路由器
func NewMultiStrategyRouter() *MultiStrategyRouter {
	return &MultiStrategyRouter{
		routers: make(map[types.RouteStrategy]types.Router),
		factory: &Factory{},
	}
}

// Route 根据任务策略进行路由
func (msr *MultiStrategyRouter) Route(task *types.Task, executors []types.Executor) (types.Executor, error) {
	msr.mutex.RLock()
	router, exists := msr.routers[task.Strategy]
	msr.mutex.RUnlock()

	if !exists {
		msr.mutex.Lock()
		// 双重检查
		if router, exists = msr.routers[task.Strategy]; !exists {
			router = msr.factory.CreateRouter(task.Strategy)
			msr.routers[task.Strategy] = router
		}
		msr.mutex.Unlock()
	}

	return router.Route(task, executors)
}

// GetStrategy 获取路由策略（实现Router接口）
func (msr *MultiStrategyRouter) GetStrategy() types.RouteStrategy {
	return types.RoundRobinApp // 默认策略
}
