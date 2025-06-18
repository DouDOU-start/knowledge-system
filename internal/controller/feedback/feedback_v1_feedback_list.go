package feedback

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"knowledge-system-api/api/feedback/v1"
)

func (c *ControllerV1) FeedbackList(ctx context.Context, req *v1.FeedbackListReq) (res *v1.FeedbackListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
