// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ImportTaskItem is the golang structure for table import_task_item.
type ImportTaskItem struct {
	Id           uint64      `json:"id"           orm:"id"            description:"条目自增ID"`            // 条目自增ID
	TaskId       string      `json:"taskId"       orm:"task_id"       description:"所属任务ID"`            // 所属任务ID
	Status       string      `json:"status"       orm:"status"        description:"条目处理状态"`            // 条目处理状态
	SourceData   string      `json:"sourceData"   orm:"source_data"   description:"原始数据 (如单条知识的JSON)"` // 原始数据 (如单条知识的JSON)
	ErrorMessage string      `json:"errorMessage" orm:"error_message" description:"处理失败时的错误信息"`        // 处理失败时的错误信息
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:"创建时间"`              // 创建时间
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:"更新时间"`              // 更新时间
}
