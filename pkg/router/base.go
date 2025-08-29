package router

import (
	"task_scheduler/pkg/types"
)

// BaseRouter 基础路由器
type BaseRouter struct {
	strategy types.RouteStrategy
}

// GetStrategy 获取路由策略
func (br *BaseRouter) GetStrategy() types.RouteStrategy {
	return br.strategy
}

// Factory 路由器工厂
type Factory struct{}

// CreateRouter 根据策略创建路由器
func (rf *Factory) CreateRouter(strategy types.RouteStrategy) types.Router {
	switch strategy {
	case types.RoundRobinTask:
		return NewRoundRobinTaskRouter()
	case types.RoundRobinApp:
		return NewRoundRobinAppRouter()
	case types.Random:
		return NewRandomRouter()
	case types.LFU:
		return NewLFURouter()
	case types.LRU:
		return NewLRURouter()
	default:
		return NewRoundRobinAppRouter() // 默认使用应用级别轮询
	}
}
