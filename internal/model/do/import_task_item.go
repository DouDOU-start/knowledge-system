// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ImportTaskItem is the golang structure of table import_task_item for DAO operations like Where/Data.
type ImportTaskItem struct {
	g.Meta       `orm:"table:import_task_item, do:true"`
	Id           interface{} // 条目自增ID
	TaskId       interface{} // 所属任务ID
	Status       interface{} // 条目处理状态
	SourceData   interface{} // 原始数据 (如单条知识的JSON)
	ErrorMessage interface{} // 处理失败时的错误信息
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
}
