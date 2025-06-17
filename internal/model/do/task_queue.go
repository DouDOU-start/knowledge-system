// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TaskQueue is the golang structure of table task_queue for DAO operations like Where/Data.
type TaskQueue struct {
	g.Meta    `orm:"table:task_queue, do:true"`
	Id        interface{} // 队列项ID
	TaskId    interface{} // 任务ID
	Priority  interface{} // 优先级
	Status    interface{} // 状态
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
	StartedAt *gtime.Time // 开始处理时间
	EndedAt   *gtime.Time // 处理结束时间
}
