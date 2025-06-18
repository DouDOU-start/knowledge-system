// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package feedback

import (
	"context"

	"knowledge-system-api/api/feedback/v1"
)

type IFeedbackV1 interface {
	FeedbackAdd(ctx context.Context, req *v1.FeedbackAddReq) (res *v1.FeedbackAddRes, err error)
	FeedbackList(ctx context.Context, req *v1.FeedbackListReq) (res *v1.FeedbackListRes, err error)
	FeedbackProcess(ctx context.Context, req *v1.FeedbackProcessReq) (res *v1.FeedbackProcessRes, err error)
}
