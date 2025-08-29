package router

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"task_scheduler/pkg/types"
)

// RandomRouter 随机路由器
type RandomRouter struct {
	BaseRouter
	rnd   *rand.Rand
	mutex sync.Mutex
}

// NewRandomRouter 创建随机路由器
func NewRandomRouter() *RandomRouter {
	return &RandomRouter{
		BaseRouter: BaseRouter{strategy: types.Random},
		rnd:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Route 随机路由
func (r *RandomRouter) Route(task *types.Task, executors []types.Executor) (types.Executor, error) {
	if len(executors) == 0 {
		return nil, errors.New("no available executors")
	}

	r.mutex.Lock()
	index := r.rnd.Intn(len(executors))
	r.mutex.Unlock()

	return executors[index], nil
}
