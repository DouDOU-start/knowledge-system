package knowledge

import (
	"context"
	v1 "knowledge-system-api/api/knowledge/v1"
	"knowledge-system-api/internal/model"
	"knowledge-system-api/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
)

// ControllerV1 知识库V1接口控制器
// 只做参数校验、调用service、返回响应
// 错误全部用gerror/gcode统一处理
type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

// BatchImport 批量导入知识条目
func (c *ControllerV1) BatchImport(ctx context.Context, req *v1.BatchImportReq) (res *v1.BatchImportRes, err error) {
	// 参数校验由框架自动完成，这里只需处理业务逻辑
	for _, item := range req.Items {
		// 调用LLM进行标签分类和摘要生成
		labels, summary, err := service.LLMClassifyByConfig(ctx, item.Content)
		if err != nil {
			g.Log().Errorf(ctx, "LLM推理失败: %v", err)
			return nil, gerror.NewCodef(gcode.CodeInternalError, "LLM推理失败: %s", err.Error())
		}

		// 过滤标签
		filtered := service.FilterLabels(labels, 3)

		// 始终生成新的 ID，不使用用户提供的 ID
		id := uuid.New().String()

		// 存入向量数据库
		if err := service.QdrantUpsert(ctx, req.RepoName, id, item.Content, summary, filtered); err != nil {
			g.Log().Errorf(ctx, "Qdrant入库失败: %v", err)
			return nil, gerror.NewCodef(gcode.CodeInternalError, "Qdrant入库失败: %s", err.Error())
		}

		// 存入MySQL
		if err := service.KnowledgeService().CreateKnowledge(ctx, id, req.RepoName, item.Content, filtered, summary); err != nil {
			g.Log().Errorf(ctx, "MySQL入库失败: %v", err)
			return nil, gerror.NewCodef(gcode.CodeInternalError, "MySQL入库失败: %s", err.Error())
		}
	}

	return &v1.BatchImportRes{Success: true}, nil
}

// BatchImportAsync 批量异步导入知识条目
func (c *ControllerV1) BatchImportAsync(ctx context.Context, req *v1.BatchImportAsyncReq) (res *v1.BatchImportAsyncRes, err error) {
	// 参数校验由框架自动完成，这里只需处理业务逻辑

	// 创建任务项
	var taskItems []model.TaskItem
	for _, item := range req.Items {
		taskItems = append(taskItems, model.TaskItem{
			// 不使用用户提供的 ID，任务项 ID 由系统在处理时生成
			RepoName: req.RepoName,
			Content:  item.Content,
			Status:   "pending",
		})
	}

	// 创建导入任务
	taskID, err := service.KnowledgeService().CreateImportTask(ctx, taskItems)
	if err != nil {
		g.Log().Errorf(ctx, "创建导入任务失败: %v", err)
		return nil, gerror.NewCodef(gcode.CodeInternalError, "创建导入任务失败: %s", err.Error())
	}

	return &v1.BatchImportAsyncRes{
		TaskID:  taskID,
		Message: "任务已创建，正在后台处理",
	}, nil
}

// TaskStatus 查询任务状态
func (c *ControllerV1) TaskStatus(ctx context.Context, req *v1.TaskStatusReq) (res *v1.TaskStatusRes, err error) {
	// 参数校验由框架自动完成，这里只需处理业务逻辑

	// 获取任务状态
	task, err := service.KnowledgeService().GetTaskStatus(ctx, req.TaskID)
	if err != nil {
		g.Log().Errorf(ctx, "获取任务状态失败: %v", err)
		return nil, gerror.NewCodef(gcode.CodeInternalError, "获取任务状态失败: %s", err.Error())
	}

	if task == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "任务不存在")
	}

	return &v1.TaskStatusRes{
		TaskID:    task.TaskID,
		Status:    task.Status,
		Progress:  task.Progress,
		Total:     task.Total,
		Processed: task.Processed,
		Failed:    task.Failed,
		Message:   task.Message,
	}, nil
}

// Classify 单条内容标签打分
func (c *ControllerV1) Classify(ctx context.Context, req *v1.ClassifyReq) (res *v1.ClassifyRes, err error) {
	// 参数校验由框架自动完成

	// 调用LLM进行标签分类和摘要生成
	// labels, summary, err := service.LLMClassifyByConfig(ctx, req.Content)
	// if err != nil {
	// 	g.Log().Errorf(ctx, "LLM推理失败: %v", err)
	// 	return nil, gerror.NewCodef(gcode.CodeInternalError, "LLM推理失败: %s", err.Error())
	// }

	// 转换为API响应格式
	// var outLabels []v1.LabelScore
	// for _, l := range labels {
	// 	outLabels = append(outLabels, v1.LabelScore{
	// 		LabelID: l.LabelID,
	// 		Score:   l.Score,
	// 	})
	// }

	// return &v1.ClassifyRes{
	// 	Labels:  outLabels,
	// 	Summary: summary,
	// }, nil

	return nil, gerror.NewCode(gcode.CodeNotImplemented, "单条内容标签打分功能尚未实现，请稍后再试")
}

// Search 知识检索
func (c *ControllerV1) Search(ctx context.Context, req *v1.SearchReq) (res *v1.SearchRes, err error) {
	// 参数校验由框架自动完成

	// 设置默认值
	if req.TopK <= 0 {
		req.TopK = 5
	}

	// 根据模式选择不同的搜索方式
	var items []model.SearchResult
	switch req.Mode {
	case "hybrid", "":
		items, err = service.KnowledgeService().SearchKnowledgeByHybrid(ctx, req.Query, req.RepoName, uint64(req.TopK))
	default:
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "不支持的搜索模式")
	}

	if err != nil {
		g.Log().Errorf(ctx, "知识检索失败: %v", err)
		return nil, gerror.NewCodef(gcode.CodeInternalError, "知识检索失败: %s", err.Error())
	}

	// 转换为API响应格式
	var outItems []v1.KnowledgeResult
	for _, item := range items {
		var outLabels []v1.LabelScore
		for _, l := range item.Labels {
			outLabels = append(outLabels, v1.LabelScore{
				Name:  l.Name,
				Score: l.Score,
			})
		}

		outItems = append(outItems, v1.KnowledgeResult{
			ID:       item.ID,
			RepoName: item.RepoName,
			Content:  item.Content,
			Labels:   outLabels,
			Summary:  item.Summary,
			Score:    item.Score,
		})
	}

	return &v1.SearchRes{Items: outItems}, nil
}

// GetRepos 获取所有知识库
func (c *ControllerV1) GetRepos(ctx context.Context, req *v1.GetReposReq) (res *v1.GetReposRes, err error) {
	// 获取所有知识库名称
	repos, err := service.KnowledgeService().GetAllRepos(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "获取知识库列表失败: %v", err)
		return nil, gerror.NewCodef(gcode.CodeInternalError, "获取知识库列表失败: %s", err.Error())
	}

	return &v1.GetReposRes{Repos: repos}, nil
}
