// =================================================================================
// This file is manually written. Do not edit by tools.
// =================================================================================

package entity

import "github.com/gogf/gf/v2/os/gtime"

// ImportTask is the golang structure for table import_task.
type ImportTask struct {
	Id        string      `json:"id"        orm:"id"        description:"任务ID"`   // 任务ID
	Status    string      `json:"status"    orm:"status"    description:"任务状态"`   // 任务状态：pending, processing, completed, failed, completed_with_errors
	Progress  int         `json:"progress"  orm:"progress"  description:"处理进度"`   // 处理进度，0-100
	Total     int         `json:"total"     orm:"total"     description:"总条目数"`   // 总条目数
	Processed int         `json:"processed" orm:"processed" description:"已处理条目数"` // 已处理条目数
	Failed    int         `json:"failed"    orm:"failed"    description:"失败条目数"`  // 失败条目数
	Message   string      `json:"message"   orm:"message"   description:"任务相关信息"` // 任务相关信息
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`  // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`  // 更新时间
}
