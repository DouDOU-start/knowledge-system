// =================================================================================
// This file is manually written. Do not edit by tools.
// =================================================================================

package dao

import (
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// importTaskItemDao is the data access object for the table import_task_item.
type importTaskItemDao struct {
	table   string                // 表名
	group   string                // 分组
	columns importTaskItemColumns // 列名
}

// importTaskItemColumns defines and stores column names for table import_task_item.
type importTaskItemColumns struct {
	Id           string // 条目自增ID
	TaskId       string // 所属任务ID
	Status       string // 条目处理状态
	SourceData   string // 原始数据
	ErrorMessage string // 错误信息
	CreatedAt    string // 创建时间
	UpdatedAt    string // 更新时间
}

// importTaskItemDao is a globally accessible object for table import_task_item operations.
var (
	ImportTaskItem = &importTaskItemDao{
		table: "import_task_item",
		group: "default",
		columns: importTaskItemColumns{
			Id:           "id",
			TaskId:       "task_id",
			Status:       "status",
			SourceData:   "source_data",
			ErrorMessage: "error_message",
			CreatedAt:    "created_at",
			UpdatedAt:    "updated_at",
		},
	}
)

// DB returns the underlying database connection of current dao.
func (dao *importTaskItemDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *importTaskItemDao) Ctx(ctx g.Ctx) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error if it returns non-nil error.
// It commits the transaction and returns nil if it returns nil.
func (dao *importTaskItemDao) Transaction(ctx g.Ctx, f func(ctx g.Ctx, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
