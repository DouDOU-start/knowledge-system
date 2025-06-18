package knowledge

import (
	"context"
	"fmt"
	"knowledge-system-api/internal/dao"
	"knowledge-system-api/internal/model/do"
	"knowledge-system-api/internal/model/entity"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

// 持久化任务队列管理器
type PersistentTaskQueue struct {
	// 内存队列，用于提高效率
	memoryQueue chan string
	// 互斥锁
	lock sync.Mutex
	// 轮询间隔
	pollInterval time.Duration
	// 是否已初始化
	initialized bool
}

// 全局持久化任务队列实例
var persistentQueue = &PersistentTaskQueue{
	memoryQueue:  make(chan string, 100),
	pollInterval: 5 * time.Second,
}

// InitPersistentQueue 初始化持久化任务队列
func InitPersistentQueue() {
	if persistentQueue.initialized {
		return
	}

	persistentQueue.lock.Lock()
	defer persistentQueue.lock.Unlock()

	if persistentQueue.initialized {
		return
	}

	// 启动任务消费者
	go persistentQueue.taskConsumer()

	// 启动数据库轮询器
	go persistentQueue.databasePoller()

	persistentQueue.initialized = true
	g.Log().Info(gctx.New(), "持久化任务队列已初始化")
}

// EnqueueTask 将任务加入队列
func EnqueueTask(ctx context.Context, taskId string, priority int) error {
	// 确保队列已初始化
	InitPersistentQueue()

	// 创建队列项
	queueId := uuid.NewString()
	now := gtime.Now()
	queueItem := do.TaskQueue{
		Id:        queueId,
		TaskId:    taskId,
		Priority:  priority,
		Status:    "waiting",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 保存到数据库
	_, err := dao.TaskQueue.Ctx(ctx).Insert(queueItem)
	if err != nil {
		return fmt.Errorf("保存任务队列项失败: %w", err)
	}

	// 尝试放入内存队列
	select {
	case persistentQueue.memoryQueue <- taskId:
		g.Log().Debug(ctx, "任务已加入内存队列:", taskId)
	default:
		// 内存队列已满，仅保存到数据库即可，稍后会由轮询器处理
		g.Log().Debug(ctx, "内存队列已满，任务已保存到数据库:", taskId)
	}

	return nil
}

// DequeueTask 从队列获取下一个任务
func DequeueTask() string {
	// 确保队列已初始化
	InitPersistentQueue()

	// 从内存队列获取任务
	taskId := <-persistentQueue.memoryQueue
	// 更新数据库中的任务状态
	ctx := gctx.New()
	dao.TaskQueue.Ctx(ctx).
		Fields("*").
		Data(do.TaskQueue{
			Status:    "processing",
			StartedAt: gtime.Now(),
			UpdatedAt: gtime.Now(),
		}).
		Where("task_id", taskId).
		Where("status", "waiting").
		Update()

	return taskId
}

// taskConsumer 任务消费者，处理任务
func (q *PersistentTaskQueue) taskConsumer() {
	ctx := gctx.New()
	for taskId := range q.memoryQueue {
		// 将任务状态更新为处理中
		dao.TaskQueue.Ctx(ctx).
			Fields("*").
			Data(do.TaskQueue{
				Status:    "processing",
				StartedAt: gtime.Now(),
				UpdatedAt: gtime.Now(),
			}).
			Where("task_id", taskId).
			Where("status", "waiting").
			Update()

		// 加入实际的处理队列
		taskChan <- taskId
	}
}

// databasePoller 数据库轮询器，定期从数据库拉取等待中的任务
func (q *PersistentTaskQueue) databasePoller() {
	ctx := gctx.New()
	ticker := time.NewTicker(q.pollInterval)
	defer ticker.Stop()

	for {
		<-ticker.C

		// 查询等待中的任务
		var waitingTasks []entity.TaskQueue
		// 使用标准链式操作查询
		err := dao.TaskQueue.Ctx(ctx).
			Fields("*").
			Where("status", "waiting").
			OrderDesc("priority").
			OrderAsc("created_at").
			Limit(20).
			Scan(&waitingTasks)

		if err != nil {
			g.Log().Error(ctx, "查询等待中的任务失败:", err)
			continue
		}

		for _, task := range waitingTasks {
			// 尝试放入内存队列
			select {
			case q.memoryQueue <- task.TaskId:
				g.Log().Debug(ctx, "从数据库轮询的任务已加入内存队列:", task.TaskId)
			default:
				// 内存队列已满，稍后再尝试
				g.Log().Debug(ctx, "内存队列已满，任务将在下次轮询时处理:", task.TaskId)
			}
		}
	}
}

// CompleteTask 标记任务完成
func CompleteTask(ctx context.Context, taskId string, success bool) {
	status := "completed"
	if !success {
		status = "failed"
	}

	dao.TaskQueue.Ctx(ctx).
		Fields("*").
		Data(do.TaskQueue{
			Status:    status,
			EndedAt:   gtime.Now(),
			UpdatedAt: gtime.Now(),
		}).
		Where("task_id", taskId).
		Update()
}
