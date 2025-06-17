// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TaskQueueDao is the data access object for the table task_queue.
type TaskQueueDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  TaskQueueColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// TaskQueueColumns defines and stores column names for the table task_queue.
type TaskQueueColumns struct {
	Id        string // 队列项ID
	TaskId    string // 任务ID
	Priority  string // 优先级
	Status    string // 状态
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
	StartedAt string // 开始处理时间
	EndedAt   string // 处理结束时间
}

// taskQueueColumns holds the columns for the table task_queue.
var taskQueueColumns = TaskQueueColumns{
	Id:        "id",
	TaskId:    "task_id",
	Priority:  "priority",
	Status:    "status",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	StartedAt: "started_at",
	EndedAt:   "ended_at",
}

// NewTaskQueueDao creates and returns a new DAO object for table data access.
func NewTaskQueueDao(handlers ...gdb.ModelHandler) *TaskQueueDao {
	return &TaskQueueDao{
		group:    "default",
		table:    "task_queue",
		columns:  taskQueueColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TaskQueueDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TaskQueueDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TaskQueueDao) Columns() TaskQueueColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TaskQueueDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TaskQueueDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *TaskQueueDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
