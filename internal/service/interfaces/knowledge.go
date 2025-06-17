package interfaces

import (
	"context"
	"knowledge-system-api/internal/model"
)

// KnowledgeService 知识库业务接口
type KnowledgeService interface {
	// CreateKnowledge 创建知识条目
	CreateKnowledge(ctx context.Context, id, content string, labels []model.LabelScore, summary string) error

	// GetKnowledgeById 根据ID获取知识条目
	GetKnowledgeById(ctx context.Context, id string) (*model.KnowledgeItem, error)

	// SearchKnowledgeByKeyword 关键词搜索知识条目
	SearchKnowledgeByKeyword(ctx context.Context, keyword string, limit int) ([]model.SearchResult, error)

	// SearchKnowledgeBySemantic 语义搜索知识条目
	SearchKnowledgeBySemantic(ctx context.Context, query string, limit int) ([]model.SearchResult, error)

	// SearchKnowledgeByHybrid 混合搜索知识条目（关键词+语义）
	SearchKnowledgeByHybrid(ctx context.Context, query string, limit int) ([]model.SearchResult, error)
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
	Upsert(id string, vector []float32, content string, labels []model.LabelScore, summary string) error

	// Search 向量搜索
	Search(query string, vector []float32, limit int) ([]model.VectorSearchResult, error)
}
