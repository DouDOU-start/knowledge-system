// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ImportTaskDao is the data access object for the table import_task.
type ImportTaskDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ImportTaskColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ImportTaskColumns defines and stores column names for the table import_task.
type ImportTaskColumns struct {
	Id        string // 任务ID
	Status    string // 任务状态
	Progress  string // 处理进度，0-100
	Total     string // 总条目数
	Processed string // 已处理条目数
	Failed    string // 失败条目数
	Message   string // 任务相关信息
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// importTaskColumns holds the columns for the table import_task.
var importTaskColumns = ImportTaskColumns{
	Id:        "id",
	Status:    "status",
	Progress:  "progress",
	Total:     "total",
	Processed: "processed",
	Failed:    "failed",
	Message:   "message",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewImportTaskDao creates and returns a new DAO object for table data access.
func NewImportTaskDao(handlers ...gdb.ModelHandler) *ImportTaskDao {
	return &ImportTaskDao{
		group:    "default",
		table:    "import_task",
		columns:  importTaskColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ImportTaskDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ImportTaskDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ImportTaskDao) Columns() ImportTaskColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ImportTaskDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ImportTaskDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ImportTaskDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
