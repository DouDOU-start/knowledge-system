package helper

import (
	"context"
	"knowledge-system-api/internal/model"
)

// LLMClassifyFunc 大模型分类函数类型
type LLMClassifyFunc func(ctx context.Context, content string) ([]model.LabelScore, string, error)

// FilterLabelsFunc 标签过滤函数类型
type FilterLabelsFunc func(labels []model.LabelScore, threshold int) []model.LabelScore

// QdrantUpsertFunc Qdrant 向量库插入函数类型
type QdrantUpsertFunc func(id string, vector []float32, content string, repoName string, labels []model.LabelScore, summary string) error

// KnowledgeServiceFunc 获取知识库服务接口实例函数类型
type KnowledgeServiceFunc func() interface{}

// 全局函数变量
var (
	// LLMClassify 大模型分类函数
	LLMClassify LLMClassifyFunc

	// FilterLabels 标签过滤函数
	FilterLabels FilterLabelsFunc

	// QdrantUpsert Qdrant 向量库插入函数
	QdrantUpsert QdrantUpsertFunc

	// GetKnowledgeService 获取知识库服务接口实例函数
	GetKnowledgeService KnowledgeServiceFunc
)

// SetLLMClassify 设置大模型分类函数
func SetLLMClassify(fn LLMClassifyFunc) {
	LLMClassify = fn
}

// SetFilterLabels 设置标签过滤函数
func SetFilterLabels(fn FilterLabelsFunc) {
	FilterLabels = fn
}

// SetQdrantUpsert 设置 Qdrant 向量库插入函数
func SetQdrantUpsert(fn QdrantUpsertFunc) {
	QdrantUpsert = fn
}

// SetKnowledgeService 设置获取知识库服务接口实例函数
func SetKnowledgeService(fn KnowledgeServiceFunc) {
	GetKnowledgeService = fn
}
