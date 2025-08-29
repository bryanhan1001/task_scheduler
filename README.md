# 任务调度器 (Task Scheduler)

基于Go语言实现的高性能任务调度器，支持多种路由策略，可以有效降低系统成本并提高资源利用率。

## 功能特性

### 🚀 多种路由策略

1. **轮询（Round Robin）**
   - **任务级别轮询**：每个任务维护独立计数器，适合调度时间一致的场景
   - **应用级别轮询**：所有任务共享计数器，保证执行器接收任务次数平均

2. **随机（Random）**
   - 完全随机选择执行器
   - 长期负载均衡，短期可能不均
   - 实现简单，适合对均衡要求不严格的场景

3. **最近最少使用（LFU）**
   - 优先选择使用次数最少的执行器
   - 适合混合多种路由策略的场景
   - 能有效平衡不同策略导致的负载不均

4. **最近最久未使用（LRU）**
   - 优先选择最久未使用的执行器
   - 基于时间的负载均衡
   - 适合需要考虑时间因素的调度场景

### 📋 核心功能

- ✅ 支持Cron表达式的定时任务调度
- ✅ 多执行器负载均衡
- ✅ 执行器健康检查
- ✅ 任务状态管理
- ✅ 统计信息收集
- ✅ 并发安全
- ✅ 优雅启停

## 项目结构

```
task_scheduler/
├── go.mod              # Go模块定义
├── types.go            # 核心类型定义
├── executor.go         # 执行器实现
├── router.go           # 路由策略实现
├── scheduler.go        # 调度器核心逻辑
├── main.go             # 示例程序
├── router_test.go      # 路由策略测试
└── README.md           # 项目文档
```

## 快速开始

### 环境要求

- Go 1.21+
- 依赖：`github.com/robfig/cron/v3`

### 安装依赖

```bash
go mod tidy
```

### 运行示例

```bash
go run .
```

### 运行测试

```bash
# 运行所有测试
go test -v

# 运行性能测试
go test -bench=.
```

## 使用示例

### 基本使用

```go
package main

import (
    "log"
    "time"
)

func main() {
    // 1. 创建调度器配置
    config := &SchedulerConfig{
        MaxConcurrentTasks:  5,
        HealthCheckInterval: 30 * time.Second,
        DefaultStrategy:     RoundRobinApp,
    }

    // 2. 创建任务调度器
    scheduler := NewTaskScheduler(config)

    // 3. 添加执行器
    executors := []*SimpleExecutor{
        NewSimpleExecutor("executor-1", "http://localhost:8001"),
        NewSimpleExecutor("executor-2", "http://localhost:8002"),
        NewSimpleExecutor("executor-3", "http://localhost:8003"),
    }

    for _, executor := range executors {
        scheduler.AddExecutor(executor)
    }

    // 4. 创建任务
    task := &Task{
        ID:       "order-timeout-check",
        Name:     "订单超时检查",
        Cron:     "0 */1 * * * *", // 每分钟执行一次
        Handler:  "orderTimeoutHandler",
        Strategy: RoundRobinTask,
    }

    // 5. 添加任务并启动调度器
    scheduler.AddTask(task)
    scheduler.Start()

    // 6. 运行一段时间后停止
    time.Sleep(5 * time.Minute)
    scheduler.Stop()
}
```

### 路由策略选择建议

根据文章建议和实际场景：

```go
// 场景1：所有任务负载相近
task := &Task{
    Strategy: RoundRobinApp, // 使用应用级别轮询
}

// 场景2：存在大任务和小任务
bigTask := &Task{
    Strategy: RoundRobinTask, // 大任务使用任务级轮询
}
smallTask := &Task{
    Strategy: Random, // 小任务使用随机或其他策略
}

// 场景3：混合策略环境中的新任务
newTask := &Task{
    Strategy: LFU, // 使用LFU平衡负载
}

// 场景4：需要考虑时间因素
timeBasedTask := &Task{
    Strategy: LRU, // 使用LRU基于时间调度
}
```

## 架构设计

### 核心组件

1. **Scheduler（调度器）**
   - 任务管理和调度
   - 执行器管理
   - 健康检查

2. **Router（路由器）**
   - 实现多种路由策略
   - 支持策略动态切换
   - 负载均衡算法

3. **Executor（执行器）**
   - 任务执行
   - 状态管理
   - 统计信息收集

4. **Task（任务）**
   - 任务定义
   - Cron调度配置
   - 路由策略配置

### 设计原则

- **单一职责**：每个组件职责明确
- **开闭原则**：易于扩展新的路由策略
- **依赖倒置**：面向接口编程
- **并发安全**：使用适当的同步机制

## 性能优化

### 路由策略性能对比

根据基准测试结果：

- **轮询策略**：性能最高，O(1)时间复杂度
- **随机策略**：性能良好，O(1)时间复杂度
- **LFU策略**：中等性能，O(n)时间复杂度
- **LRU策略**：中等性能，O(n)时间复杂度

### 优化建议

1. **高频任务**：优先使用轮询或随机策略
2. **负载均衡要求高**：使用LFU或LRU策略
3. **混合场景**：根据任务特性选择不同策略

## 监控和统计

### 任务统计

```go
// 获取任务统计信息
stats := scheduler.GetTaskStats()
fmt.Printf("总任务数: %v\n", stats["total_tasks"])
fmt.Printf("状态分布: %v\n", stats["status_distribution"])
fmt.Printf("策略分布: %v\n", stats["strategy_distribution"])
```

### 执行器统计

```go
// 获取执行器统计信息
executorStats := scheduler.GetExecutorStats()
for _, stat := range executorStats {
    fmt.Printf("执行器 %s: 使用次数=%d, 健康状态=%v\n",
        stat.ID, stat.UsageCount, stat.IsHealthy)
}
```

## 扩展开发

### 添加新的路由策略

1. 实现`Router`接口：

```go
type CustomRouter struct {
    BaseRouter
    // 自定义字段
}

func (r *CustomRouter) Route(task *Task, executors []Executor) (Executor, error) {
    // 实现自定义路由逻辑
    return selectedExecutor, nil
}
```

2. 在`RouterFactory`中注册：

```go
func (rf *RouterFactory) CreateRouter(strategy RouteStrategy) Router {
    switch strategy {
    case CustomStrategy:
        return NewCustomRouter()
    // ... 其他策略
    }
}
```

### 添加新的执行器类型

实现`Executor`接口：

```go
type CustomExecutor struct {
    // 自定义字段
}

func (e *CustomExecutor) Execute(task *Task) error {
    // 实现自定义执行逻辑
    return nil
}

// 实现其他接口方法...
```

## 最佳实践

### 1. 路由策略选择

- **任务负载相近**：使用应用级别轮询
- **大小任务混合**：大任务用任务级轮询，小任务用其他策略
- **负载不均场景**：使用LFU策略
- **时间敏感场景**：使用LRU策略

### 2. 执行器配置

- 合理设置执行器数量，避免资源浪费
- 定期进行健康检查，及时发现问题
- 监控执行器负载，动态调整

### 3. 任务设计

- 任务应该是幂等的
- 避免长时间运行的任务
- 合理设置任务超时时间

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License

## 参考资料

- [合理选择任务调度的路由策略，可以帮助降本 50%](https://mp.weixin.qq.com/s/c6IEiXsgtggeBINdtyzUTw)
- [Cron表达式文档](https://pkg.go.dev/github.com/robfig/cron/v3)
- [Go并发编程](https://golang.org/doc/effective_go.html#concurrency)
