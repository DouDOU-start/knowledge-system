package feedback

import (
	"context"
	"fmt"

	v1 "knowledge-system-api/api/feedback/v1"
	"knowledge-system-api/internal/service"
)

func (c *ControllerV1) FeedbackAdd(ctx context.Context, req *v1.FeedbackAddReq) (res *v1.FeedbackAddRes, err error) {
	// 调用服务
	feedbackID, err := service.Feedback().Add(ctx, req.SessionID, req.UserQuery, req.KnowledgeID, req.Action)
	if err != nil {
		return nil, err
	}

	// 返回结果
	res = &v1.FeedbackAddRes{
		ID: fmt.Sprintf("%d", feedbackID),
	}
	return
}
