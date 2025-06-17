package knowledge

import (
	"context"
	"encoding/json"
	"fmt"
	"knowledge-system-api/internal/dao"
	"knowledge-system-api/internal/helper"
	"knowledge-system-api/internal/model"
	"knowledge-system-api/internal/model/do"
	"knowledge-system-api/internal/model/entity"
	"knowledge-system-api/internal/service/interfaces"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

// 活跃任务管理器
var (
	// 全局任务映射表，用于跟踪正在处理的任务
	activeTasks = sync.Map{}

	// 任务队列，用于处理任务
	taskChan = make(chan string, 100)

	// 任务执行并发数
	concurrentWorkers = 3

	// 任务处理是否已初始化
	taskProcessorInitialized = false

	// 初始化锁
	initLock = sync.Mutex{}
)

// InitTaskProcessor 初始化任务处理器
func (s *Knowledge) InitTaskProcessor() {
	initLock.Lock()
	defer initLock.Unlock()

	if taskProcessorInitialized {
		return
	}

	// 启动指定数量的工作协程来处理任务
	for i := 0; i < concurrentWorkers; i++ {
		go s.taskWorker()
	}

	g.Log().Debug(gctx.New(), "任务处理器已初始化，工作协程数量:", concurrentWorkers)
	taskProcessorInitialized = true
}

// 工作协程，从任务队列中获取任务并处理
func (s *Knowledge) taskWorker() {
	ctx := gctx.New()
	for taskID := range taskChan {
		g.Log().Debug(ctx, "开始处理任务:", taskID)
		s.processTask(ctx, taskID)
	}
}

// CreateImportTask 创建导入任务
func (s *Knowledge) CreateImportTask(ctx context.Context, items []model.TaskItem) (string, error) {
	// 确保任务处理器已初始化
	s.InitTaskProcessor()

	// 生成任务ID
	taskID := uuid.NewString()
	now := gtime.Now()

	// 序列化任务项
	itemsJSON, err := json.Marshal(items)
	if err != nil {
		return "", err
	}

	// 创建任务记录
	task := do.ImportTask{
		Id:        taskID,
		Status:    "pending",
		Progress:  0,
		Total:     len(items),
		Processed: 0,
		Failed:    0,
		Message:   "任务已创建，等待处理",
		Items:     string(itemsJSON),
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 保存到数据库
	_, err = dao.ImportTask.Ctx(ctx).Data(task).Insert()
	if err != nil {
		return "", err
	}

	// 将任务ID放入任务队列
	select {
	case taskChan <- taskID:
		g.Log().Debug(ctx, "任务已加入队列:", taskID)
	default:
		// 队列已满，更新任务状态为失败
		dao.ImportTask.Ctx(ctx).Data(do.ImportTask{
			Status:    "failed",
			Message:   "任务队列已满，无法处理",
			UpdatedAt: gtime.Now(),
		}).Where(do.ImportTask{Id: taskID}).Update()
		return taskID, fmt.Errorf("任务队列已满，请稍后再试")
	}

	return taskID, nil
}

// GetTaskStatus 获取任务状态
func (s *Knowledge) GetTaskStatus(ctx context.Context, taskID string) (*model.ImportTask, error) {
	var entity entity.ImportTask
	err := dao.ImportTask.Ctx(ctx).Where(do.ImportTask{Id: taskID}).Scan(&entity)
	if err != nil {
		return nil, err
	}

	if entity.Id == "" {
		return nil, fmt.Errorf("任务不存在")
	}

	var items []model.TaskItem
	if entity.Items != "" {
		if err := json.Unmarshal([]byte(entity.Items), &items); err != nil {
			g.Log().Warning(ctx, "解析任务项JSON失败", err)
			items = []model.TaskItem{}
		}
	}

	return &model.ImportTask{
		TaskID:    entity.Id,
		Status:    entity.Status,
		Progress:  entity.Progress,
		Total:     entity.Total,
		Processed: entity.Processed,
		Failed:    entity.Failed,
		Message:   entity.Message,
		Items:     items,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}, nil
}

// UpdateTaskStatus 更新任务状态
func (s *Knowledge) UpdateTaskStatus(ctx context.Context, taskID string, status string, progress int, processed int, failed int, message string) error {
	_, err := dao.ImportTask.Ctx(ctx).Data(do.ImportTask{
		Status:    status,
		Progress:  progress,
		Processed: processed,
		Failed:    failed,
		Message:   message,
		UpdatedAt: gtime.Now(),
	}).Where(do.ImportTask{Id: taskID}).Update()
	return err
}

// processTask 处理导入任务
func (s *Knowledge) processTask(ctx context.Context, taskID string) {
	// 标记任务为处理中
	s.UpdateTaskStatus(ctx, taskID, "processing", 0, 0, 0, "任务处理中")

	// 获取任务详情
	task, err := s.GetTaskStatus(ctx, taskID)
	if err != nil {
		g.Log().Error(ctx, "获取任务详情失败:", err)
		s.UpdateTaskStatus(ctx, taskID, "failed", 0, 0, 0, "获取任务详情失败: "+err.Error())
		return
	}

	// 处理每个任务项
	processed := 0
	failed := 0
	items := task.Items
	totalItems := len(items)

	if totalItems == 0 {
		s.UpdateTaskStatus(ctx, taskID, "completed", 100, 0, 0, "任务中没有可处理的条目")
		return
	}

	// 更新任务项状态
	for i, item := range items {
		// 更新任务状态
		progress := int(float64(i) / float64(totalItems) * 100)
		s.UpdateTaskStatus(ctx, taskID, "processing", progress, processed, failed,
			fmt.Sprintf("正在处理第 %d/%d 条", i+1, totalItems))

		g.Log().Debug(ctx, fmt.Sprintf("处理任务 %s 的第 %d/%d 条: %s", taskID, i+1, totalItems, item.ID))

		// 处理单个条目
		err := s.processTaskItem(ctx, item)
		if err != nil {
			g.Log().Error(ctx, "处理任务项失败:", err, "item:", item)
			items[i].Status = "failed"
			items[i].Message = err.Error()
			failed++
		} else {
			items[i].Status = "completed"
			processed++
		}

		// 每10个条目或最后一个条目时更新任务项状态
		if (i+1)%10 == 0 || i == totalItems-1 {
			itemsJSON, _ := json.Marshal(items)
			dao.ImportTask.Ctx(ctx).Data(do.ImportTask{
				Items: string(itemsJSON),
			}).Where(do.ImportTask{Id: taskID}).Update()
		}

		// 添加短暂延迟，避免处理过快导致资源占用过高
		time.Sleep(100 * time.Millisecond)
	}

	// 更新最终状态
	finalStatus := "completed"
	finalMessage := fmt.Sprintf("任务完成: 共 %d 条, 成功 %d 条, 失败 %d 条", totalItems, processed, failed)
	if failed == totalItems {
		finalStatus = "failed"
		finalMessage = "所有条目处理失败"
	} else if failed > 0 {
		finalStatus = "completed_with_errors"
		finalMessage = fmt.Sprintf("任务部分完成: 共 %d 条, 成功 %d 条, 失败 %d 条", totalItems, processed, failed)
	}

	s.UpdateTaskStatus(ctx, taskID, finalStatus, 100, processed, failed, finalMessage)
	g.Log().Debug(ctx, "任务处理完成:", taskID, finalMessage)
}

// processTaskItem 处理单个任务项
func (s *Knowledge) processTaskItem(ctx context.Context, item model.TaskItem) error {
	// 检查服务是否已初始化
	if helper.LLMClassify == nil {
		return fmt.Errorf("LLM分类服务未初始化")
	}
	if helper.Vectorize == nil {
		return fmt.Errorf("向量化服务未初始化")
	}
	if helper.QdrantUpsert == nil {
		return fmt.Errorf("Qdrant服务未初始化")
	}
	if helper.GetKnowledgeService == nil {
		return fmt.Errorf("知识库服务未初始化")
	}

	// 调用LLM进行标签分类和摘要生成
	labels, summary, err := helper.LLMClassify(ctx, item.Content)
	if err != nil {
		return fmt.Errorf("LLM推理失败: %v", err)
	}

	// 过滤标签
	filtered := helper.FilterLabels(labels, 3)

	// 向量化
	vector, err := helper.Vectorize(ctx, item.Content)
	if err != nil {
		return fmt.Errorf("向量化失败: %v", err)
	}

	// 生成ID
	id := item.ID
	if id == "" {
		id = uuid.NewString()
	}

	// 存入向量数据库
	if err := helper.QdrantUpsert(id, vector, item.Content, filtered, summary); err != nil {
		return fmt.Errorf("Qdrant入库失败: %v", err)
	}

	// 获取知识库服务接口
	knowledgeService, ok := helper.GetKnowledgeService().(interfaces.KnowledgeService)
	if !ok {
		return fmt.Errorf("知识库服务类型转换失败")
	}

	// 存入MySQL
	if err := knowledgeService.CreateKnowledge(ctx, id, item.Content, filtered, summary); err != nil {
		return fmt.Errorf("MySQL入库失败: %v", err)
	}

	return nil
}
