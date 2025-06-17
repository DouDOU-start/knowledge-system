// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package knowledge

import (
	"context"

	"knowledge-system-api/api/knowledge/v1"
)

type IKnowledgeV1 interface {
	BatchImport(ctx context.Context, req *v1.BatchImportReq) (res *v1.BatchImportRes, err error)
	BatchImportAsync(ctx context.Context, req *v1.BatchImportAsyncReq) (res *v1.BatchImportAsyncRes, err error)
	TaskStatus(ctx context.Context, req *v1.TaskStatusReq) (res *v1.TaskStatusRes, err error)
	Classify(ctx context.Context, req *v1.ClassifyReq) (res *v1.ClassifyRes, err error)
	Search(ctx context.Context, req *v1.SearchReq) (res *v1.SearchRes, err error)
}
