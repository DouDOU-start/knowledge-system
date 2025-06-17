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
	"knowledge-system-api/internal/service"
	"sync"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
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

	// 初始化持久化队列
	InitPersistentQueue()

	// 启动指定数量的工作协程来处理任务
	for i := 0; i < concurrentWorkers; i++ {
		go s.taskWorker()
	}

	// 启动任务恢复协程
	go s.recoverTasks()

	g.Log().Debug(gctx.New(), "任务处理器已初始化，工作协程数量:", concurrentWorkers)
	taskProcessorInitialized = true
}

// 工作协程，从任务队列中获取任务并处理
func (s *Knowledge) taskWorker() {
	ctx := gctx.New()
	for {
		// 从持久化队列获取任务
		taskID := DequeueTask()
		if taskID != "" {
			g.Log().Debug(ctx, "开始处理任务:", taskID)
			success := true

			// 处理任务
			err := s.processTask(ctx, taskID)
			if err != nil {
				g.Log().Error(ctx, "处理任务失败:", err)
				success = false
			}

			// 标记任务完成
			CompleteTask(ctx, taskID, success)
		}
	}
}

// CreateImportTask 创建导入任务
func (s *Knowledge) CreateImportTask(ctx context.Context, items []model.TaskItem) (string, error) {
	// 确保任务处理器已初始化
	s.InitTaskProcessor()

	// 生成任务ID
	taskID := uuid.NewString()
	now := gtime.Now()

	// 创建任务记录
	task := do.ImportTask{
		Id:        taskID,
		Status:    "pending",
		Progress:  0,
		Total:     len(items),
		Processed: 0,
		Failed:    0,
		Message:   "任务已创建，等待处理",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 开启事务
	err := dao.ImportTask.Transaction(ctx, func(ctx g.Ctx, tx gdb.TX) error {
		// 1. 保存主任务记录
		_, err := dao.ImportTask.Ctx(ctx).Data(task).Insert()
		if err != nil {
			return err
		}

		// 2. 保存任务条目记录
		for _, item := range items {
			// 将每个任务条目序列化为 JSON
			sourceData, err := json.Marshal(map[string]interface{}{
				"repo_name": item.RepoName,
				"content":   item.Content,
			})
			if err != nil {
				return err
			}

			// 创建任务条目
			_, err = dao.ImportTaskItem.Ctx(ctx).Data(do.ImportTaskItem{
				TaskId:       taskID,
				Status:       "pending",
				SourceData:   string(sourceData),
				ErrorMessage: "",
				CreatedAt:    now,
				UpdatedAt:    now,
			}).Insert()

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("创建导入任务失败: %w", err)
	}

	// 将任务加入持久化队列
	err = EnqueueTask(ctx, taskID, 0) // 默认优先级为0
	if err != nil {
		// 更新任务状态为失败
		dao.ImportTask.Ctx(ctx).Data(do.ImportTask{
			Status:    "failed",
			Message:   "加入任务队列失败: " + err.Error(),
			UpdatedAt: gtime.Now(),
		}).Where(do.ImportTask{Id: taskID}).Update()
		return taskID, fmt.Errorf("加入任务队列失败: %w", err)
	}

	g.Log().Debug(ctx, "任务已加入队列:", taskID)
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

	// 从 import_task_item 表中获取任务条目
	var taskItems []model.TaskItem
	var items []struct {
		Id           int64  `json:"id"`
		TaskId       string `json:"task_id"`
		Status       string `json:"status"`
		SourceData   string `json:"source_data"`
		ErrorMessage string `json:"error_message"`
	}

	err = dao.ImportTaskItem.Ctx(ctx).Where("task_id=?", taskID).
		OrderAsc("id").
		Scan(&items)

	if err != nil {
		g.Log().Warning(ctx, "获取任务条目失败", err)
	} else {
		for _, item := range items {
			var sourceData map[string]interface{}
			if err := json.Unmarshal([]byte(item.SourceData), &sourceData); err != nil {
				g.Log().Warning(ctx, "解析任务条目数据失败", err)
				continue
			}

			taskItems = append(taskItems, model.TaskItem{
				ID:           item.Id,
				TaskID:       item.TaskId,
				RepoName:     sourceData["repo_name"].(string),
				Content:      sourceData["content"].(string),
				Status:       item.Status,
				ErrorMessage: item.ErrorMessage,
			})
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
		Items:     taskItems,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}, nil
}

// UpdateTaskStatus 更新任务状态
func (s *Knowledge) UpdateTaskStatus(ctx context.Context, taskID string, status string, progress uint, processed uint, failed uint, message string) error {
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
func (s *Knowledge) processTask(ctx context.Context, taskID string) error {
	// 标记任务为处理中
	s.UpdateTaskStatus(ctx, taskID, "processing", 0, 0, 0, "任务处理中")

	// 获取任务详情
	task, err := s.GetTaskStatus(ctx, taskID)
	if err != nil {
		g.Log().Error(ctx, "获取任务详情失败:", err)
		s.UpdateTaskStatus(ctx, taskID, "failed", 0, 0, 0, "获取任务详情失败: "+err.Error())
		return err
	}

	// 使用通用结构查询待处理的任务条目
	var items []struct {
		Id           int64  `json:"id"`
		TaskId       string `json:"task_id"`
		Status       string `json:"status"`
		SourceData   string `json:"source_data"`
		ErrorMessage string `json:"error_message"`
	}

	err = dao.ImportTaskItem.Ctx(ctx).Where("task_id=? AND status=?", taskID, "pending").
		OrderAsc("id").
		Scan(&items)

	if err != nil {
		errMsg := fmt.Sprintf("获取任务条目失败: %s", err.Error())
		g.Log().Error(ctx, errMsg)
		s.UpdateTaskStatus(ctx, taskID, "failed", 0, task.Processed, task.Failed, errMsg)
		return err
	}

	totalItems := len(items)
	if totalItems == 0 {
		// 检查是否所有条目都已处理完成
		if task.Processed+task.Failed == task.Total {
			finalStatus := "completed"
			finalMessage := fmt.Sprintf("任务完成: 共 %d 条, 成功 %d 条, 失败 %d 条", task.Total, task.Processed, task.Failed)

			if task.Failed == task.Total {
				finalStatus = "failed"
				finalMessage = "所有条目处理失败"
			} else if task.Failed > 0 {
				finalStatus = "completed_with_errors"
				finalMessage = fmt.Sprintf("任务部分完成: 共 %d 条, 成功 %d 条, 失败 %d 条", task.Total, task.Processed, task.Failed)
			}

			s.UpdateTaskStatus(ctx, taskID, finalStatus, 100, task.Processed, task.Failed, finalMessage)
		} else {
			s.UpdateTaskStatus(ctx, taskID, "processing",
				uint(float64(task.Processed+task.Failed)/float64(task.Total)*100),
				task.Processed, task.Failed, "待处理队列为空，但仍有条目未完成")
		}
		return nil
	}

	// 更新任务进度
	processed := task.Processed
	failed := task.Failed

	// 日志记录恢复信息
	g.Log().Info(ctx, fmt.Sprintf("开始处理任务 %s：已处理 %d 项，失败 %d 项，待处理 %d 项",
		taskID, processed, failed, totalItems))

	// 更新任务项状态
	for i, item := range items {
		// 解析任务条目数据
		var itemData map[string]interface{}
		if err := json.Unmarshal([]byte(item.SourceData), &itemData); err != nil {
			g.Log().Warning(ctx, "解析任务条目数据失败", err)
			continue
		}

		// 准备任务项数据
		content := itemData["content"].(string)
		repoName := itemData["repo_name"].(string)

		// 更新任务状态
		progress := uint(float64(processed+failed+uint(i)) / float64(task.Total) * 100)
		s.UpdateTaskStatus(ctx, taskID, "processing", progress, processed, failed,
			fmt.Sprintf("正在处理第 %d/%d 条", processed+failed+uint(i)+1, task.Total))

		g.Log().Debug(ctx, fmt.Sprintf("处理任务 %s 的第 %d/%d 条: %d",
			taskID, processed+failed+uint(i)+1, task.Total, item.Id))

		// 更新任务条目状态为处理中
		dao.ImportTaskItem.Ctx(ctx).Data(do.ImportTaskItem{
			Status:    "processing",
			UpdatedAt: gtime.Now(),
		}).Where(do.ImportTaskItem{Id: item.Id}).Update()

		// 处理单个条目
		err := s.processTaskItemContent(ctx, content, repoName)
		if err != nil {
			g.Log().Error(ctx, "处理任务项失败:", err, "item:", item)
			// 更新任务条目状态为失败
			dao.ImportTaskItem.Ctx(ctx).Data(do.ImportTaskItem{
				Status:       "failed",
				ErrorMessage: err.Error(),
				UpdatedAt:    gtime.Now(),
			}).Where(do.ImportTaskItem{Id: item.Id}).Update()
			failed++
		} else {
			// 更新任务条目状态为完成
			dao.ImportTaskItem.Ctx(ctx).Data(do.ImportTaskItem{
				Status:    "completed",
				UpdatedAt: gtime.Now(),
			}).Where(do.ImportTaskItem{Id: item.Id}).Update()
			processed++
		}

		// 每处理一项就立即更新主任务状态，确保重启后可以恢复
		dao.ImportTask.Ctx(ctx).Data(do.ImportTask{
			Processed: processed,
			Failed:    failed,
			UpdatedAt: gtime.Now(),
		}).Where(do.ImportTask{Id: taskID}).Update()

		// 添加短暂延迟，避免处理过快导致资源占用过高
		time.Sleep(100 * time.Millisecond)
	}

	// 更新最终状态
	finalStatus := "completed"
	finalMessage := fmt.Sprintf("任务完成: 共 %d 条, 成功 %d 条, 失败 %d 条", task.Total, processed, failed)
	if failed == task.Total {
		finalStatus = "failed"
		finalMessage = "所有条目处理失败"
	} else if failed > 0 {
		finalStatus = "completed_with_errors"
		finalMessage = fmt.Sprintf("任务部分完成: 共 %d 条, 成功 %d 条, 失败 %d 条", task.Total, processed, failed)
	}

	s.UpdateTaskStatus(ctx, taskID, finalStatus, 100, processed, failed, finalMessage)
	g.Log().Debug(ctx, "任务处理完成:", taskID, finalMessage)
	return nil
}

// processTaskItemContent 处理单个任务条目内容
func (s *Knowledge) processTaskItemContent(ctx context.Context, content string, repoName string) error {
	// 检查服务是否已初始化
	if helper.LLMClassify == nil {
		return fmt.Errorf("LLM分类服务未初始化")
	}

	if helper.Vectorize == nil {
		return fmt.Errorf("向量化服务未初始化")
	}

	// 1. 生成唯一ID
	id := uuid.NewString()

	// 2. 调用LLM进行分类，获取标签和摘要
	labels, summary, err := helper.LLMClassify(ctx, content)
	if err != nil {
		return fmt.Errorf("LLM分类失败: %w", err)
	}

	// 3. 过滤低分标签
	labels = helper.FilterLabels(labels, 70)

	// 4. 向量化文本
	vector, err := helper.Vectorize(ctx, content)
	if err != nil {
		return fmt.Errorf("向量化失败: %w", err)
	}

	// 5. 存入向量数据库
	if err := helper.QdrantUpsert(id, vector, content, repoName, labels, summary); err != nil {
		return fmt.Errorf("保存到向量库失败: %w", err)
	}

	// 6. 存入MySQL
	knowledgeService := service.KnowledgeService()
	if err := knowledgeService.CreateKnowledge(ctx, id, repoName, content, labels, summary); err != nil {
		return fmt.Errorf("保存到MySQL失败: %w", err)
	}

	return nil
}

// recoverTasks 恢复未完成的任务
func (s *Knowledge) recoverTasks() {
	ctx := context.Background()
	g.Log().Debug(ctx, "开始恢复未完成的任务...")

	// 查找所有未完成的任务
	var entities []entity.ImportTask
	err := dao.ImportTask.Ctx(ctx).
		Where("status IN(?)", g.Slice{"pending", "processing"}).
		OrderAsc("created_at").
		Scan(&entities)

	if err != nil {
		g.Log().Error(ctx, "查询未完成任务失败:", err)
		return
	}

	if len(entities) == 0 {
		g.Log().Debug(ctx, "没有需要恢复的任务")
		return
	}

	g.Log().Infof(ctx, "找到 %d 个未完成的任务", len(entities))

	// 恢复任务
	for _, e := range entities {
		// 更新任务状态
		g.Log().Infof(ctx, "恢复任务: %s, 原状态: %s", e.Id, e.Status)

		// 标记为正在处理
		dao.ImportTask.Ctx(ctx).Data(do.ImportTask{
			Status:    "processing",
			Message:   "任务恢复处理中",
			UpdatedAt: gtime.Now(),
		}).Where(do.ImportTask{Id: e.Id}).Update()

		// 将任务加入队列
		err := EnqueueTask(ctx, e.Id, 0)
		if err != nil {
			g.Log().Error(ctx, "将任务加入队列失败:", err)
			dao.ImportTask.Ctx(ctx).Data(do.ImportTask{
				Status:    "failed",
				Message:   "恢复任务失败: " + err.Error(),
				UpdatedAt: gtime.Now(),
			}).Where(do.ImportTask{Id: e.Id}).Update()
		}
	}

	g.Log().Info(ctx, "任务恢复完成")
}
