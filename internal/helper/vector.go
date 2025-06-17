package helper

import (
	"context"
	"knowledge-system-api/internal/model"
)

// VectorizeFunc 向量化函数类型
type VectorizeFunc func(ctx context.Context, text string) ([]float32, error)

// VectorSearchFunc 向量搜索函数类型
type VectorSearchFunc func(query string, vector []float32, repoName string, limit int) ([]model.VectorSearchResult, error)

// 全局函数变量
var (
	// Vectorize 向量化函数
	Vectorize VectorizeFunc

	// VectorSearch 向量搜索函数
	VectorSearch VectorSearchFunc
)

// SetVectorize 设置向量化函数
func SetVectorize(fn VectorizeFunc) {
	Vectorize = fn
}

// SetVectorSearch 设置向量搜索函数
func SetVectorSearch(fn VectorSearchFunc) {
	VectorSearch = fn
}
