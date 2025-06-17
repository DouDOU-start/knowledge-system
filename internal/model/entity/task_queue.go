package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TaskQueue 任务队列实体模型
type TaskQueue struct {
	Id        string      `json:"id"`         // 队列项ID
	TaskId    string      `json:"task_id"`    // 任务ID
	Priority  int         `json:"priority"`   // 优先级
	Status    string      `json:"status"`     // 状态：waiting, processing, completed, failed
	CreatedAt *gtime.Time `json:"created_at"` // 创建时间
	UpdatedAt *gtime.Time `json:"updated_at"` // 更新时间
	StartedAt *gtime.Time `json:"started_at"` // 开始处理时间
	EndedAt   *gtime.Time `json:"ended_at"`   // 处理结束时间
}
