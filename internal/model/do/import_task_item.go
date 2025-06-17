// =================================================================================
// This file is manually written. Do not edit by tools.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ImportTaskItem is the golang structure of table import_task_item for DAO operations like Where/Data.
type ImportTaskItem struct {
	g.Meta       `orm:"table:import_task_item, do:true"`
	Id           interface{} // 条目自增ID
	TaskId       interface{} // 所属任务ID
	Status       interface{} // 条目处理状态
	SourceData   interface{} // 原始数据
	ErrorMessage interface{} // 错误信息
	CreatedAt    interface{} // 创建时间
	UpdatedAt    interface{} // 更新时间
}
