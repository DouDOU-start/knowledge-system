// =================================================================================
// This file is manually written. Do not edit by tools.
// =================================================================================

package dao

import (
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// importTaskDao is the data access object for the table import_task.
type importTaskDao struct {
	table   string            // 表名
	group   string            // 分组
	columns importTaskColumns // 列名
}

// importTaskColumns defines and stores column names for table import_task.
type importTaskColumns struct {
	Id        string // 任务ID
	Status    string // 任务状态
	Progress  string // 处理进度
	Total     string // 总条目数
	Processed string // 已处理条目数
	Failed    string // 失败条目数
	Message   string // 任务相关信息
	Items     string // 任务项JSON数据
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// importTaskDao is a globally accessible object for table import_task operations.
var (
	ImportTask = &importTaskDao{
		table: "import_task",
		group: "default",
		columns: importTaskColumns{
			Id:        "id",
			Status:    "status",
			Progress:  "progress",
			Total:     "total",
			Processed: "processed",
			Failed:    "failed",
			Message:   "message",
			Items:     "items",
			CreatedAt: "created_at",
			UpdatedAt: "updated_at",
		},
	}
)

// DB returns the underlying database connection of current dao.
func (dao *importTaskDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *importTaskDao) Ctx(ctx g.Ctx) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error if it returns non-nil error.
// It commits the transaction and returns nil if it returns nil.
func (dao *importTaskDao) Transaction(ctx g.Ctx, f func(ctx g.Ctx, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
