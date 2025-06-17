package service

import (
	"context"
	"knowledge-system-api/internal/helper"
	"knowledge-system-api/internal/logic/knowledge"
	"knowledge-system-api/internal/model"
	"knowledge-system-api/internal/service/interfaces"
)

func init() {
	// 初始化向量化函数
	helper.SetVectorize(Vectorize)

	// 初始化向量搜索函数
	helper.SetVectorSearch(QdrantSearch)
}

// LLMClassifyByConfig 调用配置指定的大模型推理后端
func LLMClassifyByConfig(ctx context.Context, content string) (labels []model.LabelScore, summary string, err error) {
	return GetLLMClient().Classify(ctx, content)
}

// FilterLabels 过滤低分标签
func FilterLabels(labels []model.LabelScore, threshold int) []model.LabelScore {
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

// QdrantSearch 调用 Qdrant 搜索
func QdrantSearch(query string, vector []float32, limit int) ([]model.VectorSearchResult, error) {
	// 这里实现向量搜索逻辑
	// 临时实现，实际应该调用 Qdrant 客户端
	return []model.VectorSearchResult{}, nil
}

// 知识库服务实现
type knowledgeServiceImpl struct{}

// 初始化时创建实例
var knowledgeService = &knowledgeServiceImpl{}

// KnowledgeService 获取知识库服务实例
func KnowledgeService() interfaces.KnowledgeService {
	return knowledgeService
}

// CreateKnowledge 创建知识条目
func (s *knowledgeServiceImpl) CreateKnowledge(ctx context.Context, id, content string, labels []model.LabelScore, summary string) error {
	return knowledge.New().CreateKnowledge(ctx, id, content, labels, summary)
}

// GetKnowledgeById 根据ID获取知识条目
func (s *knowledgeServiceImpl) GetKnowledgeById(ctx context.Context, id string) (*model.KnowledgeItem, error) {
	return knowledge.New().GetKnowledgeById(ctx, id)
}

// SearchKnowledgeByKeyword 关键词搜索知识条目
func (s *knowledgeServiceImpl) SearchKnowledgeByKeyword(ctx context.Context, keyword string, limit int) ([]model.SearchResult, error) {
	return knowledge.New().SearchKnowledgeByKeyword(ctx, keyword, limit)
}

// SearchKnowledgeBySemantic 语义搜索知识条目
func (s *knowledgeServiceImpl) SearchKnowledgeBySemantic(ctx context.Context, query string, limit int) ([]model.SearchResult, error) {
	return knowledge.New().SearchKnowledgeBySemantic(ctx, query, limit)
}

// SearchKnowledgeByHybrid 混合搜索知识条目（关键词+语义）
func (s *knowledgeServiceImpl) SearchKnowledgeByHybrid(ctx context.Context, query string, limit int) ([]model.SearchResult, error) {
	return knowledge.New().SearchKnowledgeByHybrid(ctx, query, limit)
}
