// internal/service/feedback.go
package service

import (
	"context"
	"time"
)

// FeedbackItem 反馈项
type FeedbackItem struct {
	ID          string
	SessionID   string
	UserQuery   string
	KnowledgeID string
	Action      string
	Timestamp   time.Time
}

// IFeedback 反馈服务接口
type IFeedback interface {
	// Add 添加反馈
	Add(ctx context.Context, sessionID, userQuery, knowledgeID, action string) (int64, error)

	// List 查询反馈列表
	List(ctx context.Context, sessionID, knowledgeID, action, startTime, endTime string, page, pageSize int) ([]FeedbackItem, int, error)

	// ProcessFeedbacks 处理反馈数据，更新标签分数
	ProcessFeedbacks(ctx context.Context) error
}

var (
	localFeedback IFeedback
)

// Feedback 获取反馈服务
func Feedback() IFeedback {
	if localFeedback == nil {
		panic("implement not found for interface IFeedback, forgot register?")
	}
	return localFeedback
}

// RegisterFeedback 注册反馈服务
func RegisterFeedback(i IFeedback) {
	localFeedback = i
}
