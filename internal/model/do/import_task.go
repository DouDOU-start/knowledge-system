// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ImportTask is the golang structure of table import_task for DAO operations like Where/Data.
type ImportTask struct {
	g.Meta    `orm:"table:import_task, do:true"`
	Id        interface{} // 任务ID
	Status    interface{} // 任务状态
	Progress  interface{} // 处理进度，0-100
	Total     interface{} // 总条目数
	Processed interface{} // 已处理条目数
	Failed    interface{} // 失败条目数
	Message   interface{} // 任务相关信息
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
