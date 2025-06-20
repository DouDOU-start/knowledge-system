package service

import (
	"context"
	"knowledge-system-api/internal/helper"
	"knowledge-system-api/internal/model"
	"knowledge-system-api/internal/service/interfaces"

	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	// 初始化向量化函数
	helper.SetVectorize(Vectorize)

	// 初始化向量搜索函数
	helper.SetVectorSearch(func(repoName string, content string, labels []model.LabelScore, limit uint64) ([]model.VectorSearchResult, error) {
		ctx := context.Background()
		return QdrantSearch(ctx, repoName, content, labels, limit)
	})

	// 初始化 LLM 分类函数
	helper.SetLLMClassify(LLMClassifyByConfig)

	// 初始化标签过滤函数
	helper.SetFilterLabels(FilterLabels)

	// 初始化 Qdrant 向量库插入函数
	helper.SetQdrantUpsert(QdrantUpsert)

	// 初始化知识库服务获取函数
	helper.SetKnowledgeService(func() interface{} {
		return KnowledgeService()
	})
}

// LLMClassifyByConfig 调用配置指定的大模型推理后端
func LLMClassifyByConfig(ctx context.Context, content string) (labels []model.LabelScore, summary string, err error) {
	return GetLLMClient().Classify(ctx, content)
}

// FilterLabels 过滤低分标签
func FilterLabels(labels []model.LabelScore, threshold float32) []model.LabelScore {
	var filtered []model.LabelScore
	for _, l := range labels {
		if l.Score >= threshold {
			filtered = append(filtered, l)
		}
	}
	return filtered
}

// Vectorize 调用配置指定的向量化后端
func Vectorize(ctx context.Context, content string) ([]float32, error) {
	return GetEmbeddingClient().Embed(ctx, content)
}

// 知识库服务接口实现
type knowledgeServiceImpl struct{}

// 初始化时创建实例
var knowledgeService = &knowledgeServiceImpl{}

// KnowledgeService 获取知识库服务实例
func KnowledgeService() interfaces.KnowledgeService {
	return knowledgeService
}

// 以下方法将通过 logic 层注入实现
var (
	// CreateKnowledgeLogic 创建知识条目逻辑
	CreateKnowledgeLogic func(ctx context.Context, id, repoName, content string, labels []model.LabelScore, summary string) error

	// GetKnowledgeByIdLogic 根据ID获取知识条目逻辑
	GetKnowledgeByIdLogic func(ctx context.Context, id string) (*model.KnowledgeItem, error)

	// SearchKnowledgeByHybridLogic 混合搜索知识条目逻辑
	SearchKnowledgeByHybridLogic func(ctx context.Context, query string, repoName string, limit uint64) ([]model.SearchResult, error)

	// CreateImportTaskLogic 创建导入任务逻辑
	CreateImportTaskLogic func(ctx context.Context, items []model.TaskItem) (string, error)

	// GetTaskStatusLogic 获取任务状态逻辑
	GetTaskStatusLogic func(ctx context.Context, taskId string) (*model.ImportTask, error)

	// UpdateTaskStatusLogic 更新任务状态逻辑
	UpdateTaskStatusLogic func(ctx context.Context, taskId string, status string, progress uint, processed uint, failed uint, message string) error

	// GetAllReposLogic 获取所有知识库名称逻辑
	GetAllReposLogic func(ctx context.Context) ([]string, error)

	// RecoverTasksLogic 恢复未完成任务逻辑
	RecoverTasksLogic func()
)

// RegisterKnowledgeLogic 注册知识库业务逻辑实现
func RegisterKnowledgeLogic(
	createKnowledge func(ctx context.Context, id, repoName, content string, labels []model.LabelScore, summary string) error,
	getKnowledgeById func(ctx context.Context, id string) (*model.KnowledgeItem, error),
	searchByHybrid func(ctx context.Context, query string, repoName string, limit uint64) ([]model.SearchResult, error),
	createImportTask func(ctx context.Context, items []model.TaskItem) (string, error),
	getTaskStatus func(ctx context.Context, taskId string) (*model.ImportTask, error),
	updateTaskStatus func(ctx context.Context, taskId string, status string, progress uint, processed uint, failed uint, message string) error,
	getAllRepos func(ctx context.Context) ([]string, error),
	recoverTasks func(),
) {
	CreateKnowledgeLogic = createKnowledge
	GetKnowledgeByIdLogic = getKnowledgeById
	SearchKnowledgeByHybridLogic = searchByHybrid
	CreateImportTaskLogic = createImportTask
	GetTaskStatusLogic = getTaskStatus
	UpdateTaskStatusLogic = updateTaskStatus
	GetAllReposLogic = getAllRepos
	RecoverTasksLogic = recoverTasks
}

// CreateKnowledge 创建知识条目
func (s *knowledgeServiceImpl) CreateKnowledge(ctx context.Context, id, repoName, content string, labels []model.LabelScore, summary string) error {
	if CreateKnowledgeLogic == nil {
		return context.Canceled
	}
	return CreateKnowledgeLogic(ctx, id, repoName, content, labels, summary)
}

// GetKnowledgeById 根据ID获取知识条目
func (s *knowledgeServiceImpl) GetKnowledgeById(ctx context.Context, id string) (*model.KnowledgeItem, error) {
	if GetKnowledgeByIdLogic == nil {
		return nil, context.Canceled
	}
	return GetKnowledgeByIdLogic(ctx, id)
}

// SearchKnowledgeByHybrid 混合搜索知识条目（关键词+语义）
func (s *knowledgeServiceImpl) SearchKnowledgeByHybrid(ctx context.Context, query string, repoName string, limit uint64) ([]model.SearchResult, error) {
	if SearchKnowledgeByHybridLogic == nil {
		return nil, context.Canceled
	}
	return SearchKnowledgeByHybridLogic(ctx, query, repoName, limit)
}

// CreateImportTask 创建导入任务
func (s *knowledgeServiceImpl) CreateImportTask(ctx context.Context, items []model.TaskItem) (string, error) {
	if CreateImportTaskLogic == nil {
		return "", context.Canceled
	}
	return CreateImportTaskLogic(ctx, items)
}

// GetTaskStatus 获取任务状态
func (s *knowledgeServiceImpl) GetTaskStatus(ctx context.Context, taskId string) (*model.ImportTask, error) {
	if GetTaskStatusLogic == nil {
		return nil, context.Canceled
	}
	return GetTaskStatusLogic(ctx, taskId)
}

// UpdateTaskStatus 更新任务状态
func (s *knowledgeServiceImpl) UpdateTaskStatus(ctx context.Context, taskId string, status string, progress uint, processed uint, failed uint, message string) error {
	if UpdateTaskStatusLogic == nil {
		return context.Canceled
	}
	return UpdateTaskStatusLogic(ctx, taskId, status, progress, processed, failed, message)
}

// GetAllRepos 获取所有知识库名称
func (s *knowledgeServiceImpl) GetAllRepos(ctx context.Context) ([]string, error) {
	if GetAllReposLogic == nil {
		return nil, context.Canceled
	}
	return GetAllReposLogic(ctx)
}

// RecoverUnfinishedTasks 恢复未完成的任务
// 在服务启动时调用
func RecoverUnfinishedTasks(ctx context.Context) {
	// 直接调用 logic 层的知识库实例的恢复任务方法
	g.Log().Info(ctx, "开始恢复未完成任务...")

	// 注意：这里我们不能直接引用 knowledge.Knowledge{}，因为会导致导入循环
	// 我们将在 init.go 中注册 RecoverTasksLogic 函数
	if RecoverTasksLogic != nil {
		RecoverTasksLogic()
	} else {
		g.Log().Warning(ctx, "任务恢复函数未注册，无法恢复未完成任务")
	}
}
