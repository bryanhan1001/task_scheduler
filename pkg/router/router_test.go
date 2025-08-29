package router

import (
	"testing"
	"time"

	"task_scheduler/pkg/executor"
	"task_scheduler/pkg/types"
)

// 创建测试用的执行器
func createTestExecutors() []types.Executor {
	return []types.Executor{
		executor.NewSimpleExecutor("exec-1", "http://localhost:8001"),
		executor.NewSimpleExecutor("exec-2", "http://localhost:8002"),
		executor.NewSimpleExecutor("exec-3", "http://localhost:8003"),
	}
}

// 创建测试任务
func createTestTask(id string, strategy types.RouteStrategy) *types.Task {
	return &types.Task{
		ID:       id,
		Name:     "Test Task",
		Handler:  "testHandler",
		Strategy: strategy,
	}
}

func TestRoundRobinTaskRouter(t *testing.T) {
	router := NewRoundRobinTaskRouter()
	executors := createTestExecutors()
	task := createTestTask("test-task", types.RoundRobinTask)

	// 测试多次路由，应该轮询分配
	results := make(map[string]int)
	for i := 0; i < 9; i++ {
		exec, err := router.Route(task, executors)
		if err != nil {
			t.Fatalf("Route failed: %v", err)
		}
		results[exec.GetID()]++
	}

	// 验证每个执行器都被分配了任务
	for _, exec := range executors {
		if results[exec.GetID()] == 0 {
			t.Errorf("Executor %s was not assigned any tasks", exec.GetID())
		}
	}
}

func TestRoundRobinAppRouter(t *testing.T) {
	router := NewRoundRobinAppRouter()
	executors := createTestExecutors()
	task1 := createTestTask("task-1", types.RoundRobinApp)
	task2 := createTestTask("task-2", types.RoundRobinApp)

	// 测试不同任务使用同一个计数器
	exec1, _ := router.Route(task1, executors)
	exec2, _ := router.Route(task2, executors)
	exec3, _ := router.Route(task1, executors)

	// 应该按顺序分配
	if exec1.GetID() == exec2.GetID() {
		t.Error("App-level round robin should use different executors for consecutive calls")
	}
	if exec2.GetID() == exec3.GetID() {
		t.Error("App-level round robin should use different executors for consecutive calls")
	}
}

func TestRandomRouter(t *testing.T) {
	router := NewRandomRouter()
	executors := createTestExecutors()
	task := createTestTask("random-task", types.Random)

	// 测试随机路由
	results := make(map[string]int)
	for i := 0; i < 100; i++ {
		exec, err := router.Route(task, executors)
		if err != nil {
			t.Fatalf("Route failed: %v", err)
		}
		results[exec.GetID()]++
	}

	// 验证所有执行器都被使用过（概率性检查）
	for _, exec := range executors {
		if results[exec.GetID()] == 0 {
			t.Errorf("Executor %s was never selected in 100 random selections", exec.GetID())
		}
	}
}

func TestLFURouter(t *testing.T) {
	router := NewLFURouter()
	executors := createTestExecutors()
	task := createTestTask("lfu-task", types.LFU)

	// 第一次调用应该选择第一个执行器
	exec1, err := router.Route(task, executors)
	if err != nil {
		t.Fatalf("Route failed: %v", err)
	}

	// 第二次调用应该选择不同的执行器（因为第一个已经被使用过）
	exec2, err := router.Route(task, executors)
	if err != nil {
		t.Fatalf("Route failed: %v", err)
	}

	if exec1.GetID() == exec2.GetID() {
		t.Error("LFU router should select different executors when usage counts are different")
	}
}

func TestLRURouter(t *testing.T) {
	router := NewLRURouter()
	executors := createTestExecutors()
	task := createTestTask("lru-task", types.LRU)

	// 第一次调用
	exec1, err := router.Route(task, executors)
	if err != nil {
		t.Fatalf("Route failed: %v", err)
	}

	// 等待一小段时间
	time.Sleep(10 * time.Millisecond)

	// 第二次调用应该选择不同的执行器
	exec2, err := router.Route(task, executors)
	if err != nil {
		t.Fatalf("Route failed: %v", err)
	}

	if exec1.GetID() == exec2.GetID() {
		t.Error("LRU router should select different executors when last used times are different")
	}
}

func TestMultiStrategyRouter(t *testing.T) {
	router := NewMultiStrategyRouter()
	executors := createTestExecutors()

	// 测试不同策略
	strategies := []types.RouteStrategy{
		types.RoundRobinTask,
		types.RoundRobinApp,
		types.Random,
		types.LFU,
		types.LRU,
	}

	for _, strategy := range strategies {
		task := createTestTask("multi-task", strategy)
		exec, err := router.Route(task, executors)
		if err != nil {
			t.Fatalf("Route failed for strategy %v: %v", strategy, err)
		}
		if exec == nil {
			t.Fatalf("Route returned nil executor for strategy %v", strategy)
		}
	}
}

func TestEmptyExecutors(t *testing.T) {
	router := NewRoundRobinAppRouter()
	task := createTestTask("empty-test", types.RoundRobinApp)

	// 测试空执行器列表
	_, err := router.Route(task, []types.Executor{})
	if err == nil {
		t.Error("Route should return error when no executors available")
	}
}

// 基准测试
func BenchmarkRoundRobinTaskRouter(b *testing.B) {
	router := NewRoundRobinTaskRouter()
	executors := createTestExecutors()
	task := createTestTask("bench-task", types.RoundRobinTask)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = router.Route(task, executors)
	}
}

func BenchmarkRandomRouter(b *testing.B) {
	router := NewRandomRouter()
	executors := createTestExecutors()
	task := createTestTask("bench-task", types.Random)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = router.Route(task, executors)
	}
}
