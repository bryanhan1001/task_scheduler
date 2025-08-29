# Task Scheduler 架构设计

## 概述

Task Scheduler 是一个高性能、可扩展的 Go 任务调度系统，支持多种路由策略和分布式执行。

## 核心组件

### 1. Types (pkg/types)

定义了系统的核心接口和数据结构：

- `RouteStrategy`: 路由策略枚举
- `Executor`: 执行器接口
- `Task`: 任务结构体
- `Router`: 路由器接口
- `Scheduler`: 调度器接口

### 2. Executor (pkg/executor)

执行器组件负责任务的实际执行：

- `SimpleExecutor`: 基础执行器实现
- `Manager`: 执行器管理器，负责执行器的生命周期管理

### 3. Router (pkg/router)

路由组件实现了多种任务分发策略：

- `RoundRobinTaskRouter`: 任务级轮询路由
- `RoundRobinAppRouter`: 应用级轮询路由
- `RandomRouter`: 随机路由
- `LFURouter`: 最少使用频率路由
- `LRURouter`: 最近最少使用路由
- `MultiStrategyRouter`: 多策略路由器

### 4. Scheduler (pkg/scheduler)

调度器是系统的核心组件：

- `TaskScheduler`: 主调度器实现
- 支持 Cron 表达式定时任务
- 健康检查机制
- 统计信息收集

## 架构图

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client API    │    │   Scheduler     │    │   Executors     │
│                 │    │                 │    │                 │
│ - Add Task      │───▶│ - Task Queue    │───▶│ - Worker 1      │
│ - Remove Task   │    │ - Cron Jobs     │    │ - Worker 2      │
│ - Get Stats     │    │ - Health Check  │    │ - Worker N      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │     Router      │
                       │                 │
                       │ - Round Robin   │
                       │ - Random        │
                       │ - LFU/LRU       │
                       └─────────────────┘
```

## 设计原则

### 1. 单一职责原则
每个组件都有明确的职责：
- Types: 定义接口和数据结构
- Executor: 任务执行
- Router: 任务路由
- Scheduler: 任务调度

### 2. 开闭原则
系统对扩展开放，对修改关闭：
- 新的路由策略可以通过实现 Router 接口添加
- 新的执行器类型可以通过实现 Executor 接口添加

### 3. 依赖倒置原则
高层模块不依赖低层模块，都依赖于抽象：
- Scheduler 依赖 Router 接口，而不是具体实现
- Router 依赖 Executor 接口，而不是具体实现

## 并发安全

系统在多个层面保证并发安全：

1. **Scheduler 层面**: 使用 sync.RWMutex 保护任务和执行器列表
2. **Router 层面**: 各路由器使用 sync.Mutex 保护内部状态
3. **Executor 层面**: 执行器状态更新使用原子操作

## 扩展性

### 水平扩展
- 支持动态添加/移除执行器
- 执行器可以分布在不同的机器上

### 垂直扩展
- 支持配置最大并发任务数
- 可调整健康检查间隔

## 监控和观测

### 统计信息
- 任务执行统计（成功/失败次数、平均执行时间）
- 执行器统计（使用次数、健康状态）

### 健康检查
- 定期检查执行器健康状态
- 自动移除不健康的执行器

## 性能优化

1. **路由算法优化**: 不同策略针对不同场景优化
2. **内存管理**: 合理使用对象池减少 GC 压力
3. **并发控制**: 精细化锁粒度，减少锁竞争
4. **批处理**: 支持批量操作减少系统调用

## 容错机制

1. **执行器故障**: 自动检测并移除故障执行器
2. **任务失败**: 支持任务重试机制
3. **网络分区**: 优雅处理网络异常
4. **资源限制**: 防止资源耗尽的保护机制