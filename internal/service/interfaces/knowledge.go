package interfaces

import (
	"context"
	"knowledge-system-api/internal/model"
)

// KnowledgeService 知识库业务接口
type KnowledgeService interface {
	// CreateKnowledge 创建知识条目
	CreateKnowledge(ctx context.Context, id, repoName, content string, labels []model.LabelScore, summary string) error

	// GetKnowledgeById 根据ID获取知识条目
	GetKnowledgeById(ctx context.Context, id string) (*model.KnowledgeItem, error)

	// SearchKnowledgeByHybrid 混合搜索知识条目（关键词+语义）
	SearchKnowledgeByHybrid(ctx context.Context, query string, repoName string, limit uint64) ([]model.SearchResult, error)

	// CreateImportTask 创建导入任务
	CreateImportTask(ctx context.Context, items []model.TaskItem) (string, error)

	// GetTaskStatus 获取任务状态
	GetTaskStatus(ctx context.Context, taskId string) (*model.ImportTask, error)

	// UpdateTaskStatus 更新任务状态
	UpdateTaskStatus(ctx context.Context, taskId string, status string, progress uint, processed uint, failed uint, message string) error

	// GetAllRepos 获取所有知识库名称
	GetAllRepos(ctx context.Context) ([]string, error)
}

// EmbeddingService 向量嵌入服务接口
type EmbeddingService interface {
	// Embed 将文本转换为向量
	Embed(ctx context.Context, text string) ([]float32, error)
}

// LLMService 大语言模型服务接口
type LLMService interface {
	// Classify 对文本进行分类并生成摘要
	Classify(ctx context.Context, text string) ([]model.LabelScore, string, error)
}

// VectorDBService 向量数据库服务接口
type VectorDBService interface {
	// Upsert 插入或更新向量
	Upsert(id string, vector []float32, content string, repoName string, labels []model.LabelScore, summary string) error

	// Search 向量搜索
	Search(query string, vector []float32, repoName string, limit int) ([]model.VectorSearchResult, error)
}
