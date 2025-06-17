// =================================================================================
// This file is manually written. Do not edit by tools.
// =================================================================================

package entity

import "github.com/gogf/gf/v2/os/gtime"

// ImportTaskItem is the golang structure for table import_task_item.
type ImportTaskItem struct {
	Id           int64       `json:"id"            orm:"id"            description:"条目自增ID"` // 条目自增ID
	TaskId       string      `json:"task_id"       orm:"task_id"       description:"所属任务ID"` // 所属任务ID
	Status       string      `json:"status"        orm:"status"        description:"条目处理状态"` // 条目处理状态：pending, processing, completed, failed
	SourceData   string      `json:"source_data"   orm:"source_data"   description:"原始数据"`   // 原始数据 (如单条知识的JSON)
	ErrorMessage string      `json:"error_message" orm:"error_message" description:"错误信息"`   // 处理失败时的错误信息
	CreatedAt    *gtime.Time `json:"created_at"    orm:"created_at"    description:"创建时间"`   // 创建时间
	UpdatedAt    *gtime.Time `json:"updated_at"    orm:"updated_at"    description:"更新时间"`   // 更新时间
}
