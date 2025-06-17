// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TaskQueue is the golang structure for table task_queue.
type TaskQueue struct {
	Id        string      `json:"id"        orm:"id"         description:"队列项ID"`  // 队列项ID
	TaskId    string      `json:"taskId"    orm:"task_id"    description:"任务ID"`   // 任务ID
	Priority  int         `json:"priority"  orm:"priority"   description:"优先级"`    // 优先级
	Status    string      `json:"status"    orm:"status"     description:"状态"`     // 状态
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`   // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`   // 更新时间
	StartedAt *gtime.Time `json:"startedAt" orm:"started_at" description:"开始处理时间"` // 开始处理时间
	EndedAt   *gtime.Time `json:"endedAt"   orm:"ended_at"   description:"处理结束时间"` // 处理结束时间
}
