// Package knowledge 知识库模块接口定义
package knowledge

import (
	"context"

	v1 "knowledge-system-api/api/knowledge/v1"
)

// IKnowledgeV1 知识库V1版本接口定义
type IKnowledgeV1 interface {
	// BatchImport 批量导入知识条目
	BatchImport(ctx context.Context, req *v1.BatchImportReq) (res *v1.BatchImportRes, err error)

	// Classify 单条内容标签打分
	Classify(ctx context.Context, req *v1.ClassifyReq) (res *v1.ClassifyRes, err error)

	// Search 知识检索
	Search(ctx context.Context, req *v1.SearchReq) (res *v1.SearchRes, err error)
}
