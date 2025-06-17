// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ImportTaskItemDao is the data access object for the table import_task_item.
type ImportTaskItemDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  ImportTaskItemColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// ImportTaskItemColumns defines and stores column names for the table import_task_item.
type ImportTaskItemColumns struct {
	Id           string // 条目自增ID
	TaskId       string // 所属任务ID
	Status       string // 条目处理状态
	SourceData   string // 原始数据 (如单条知识的JSON)
	ErrorMessage string // 处理失败时的错误信息
	CreatedAt    string // 创建时间
	UpdatedAt    string // 更新时间
}

// importTaskItemColumns holds the columns for the table import_task_item.
var importTaskItemColumns = ImportTaskItemColumns{
	Id:           "id",
	TaskId:       "task_id",
	Status:       "status",
	SourceData:   "source_data",
	ErrorMessage: "error_message",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewImportTaskItemDao creates and returns a new DAO object for table data access.
func NewImportTaskItemDao(handlers ...gdb.ModelHandler) *ImportTaskItemDao {
	return &ImportTaskItemDao{
		group:    "default",
		table:    "import_task_item",
		columns:  importTaskItemColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ImportTaskItemDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ImportTaskItemDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ImportTaskItemDao) Columns() ImportTaskItemColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ImportTaskItemDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ImportTaskItemDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ImportTaskItemDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
