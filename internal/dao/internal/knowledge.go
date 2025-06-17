// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// KnowledgeDao is the data access object for the table knowledge.
type KnowledgeDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  KnowledgeColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// KnowledgeColumns defines and stores column names for the table knowledge.
type KnowledgeColumns struct {
	Id      string // 唯一ID，服务端生成uuid
	Content string // 知识内容
	Labels  string // 标签分数数组，存储为JSON字符串
	Summary string // 内容摘要
}

// knowledgeColumns holds the columns for the table knowledge.
var knowledgeColumns = KnowledgeColumns{
	Id:      "id",
	Content: "content",
	Labels:  "labels",
	Summary: "summary",
}

// NewKnowledgeDao creates and returns a new DAO object for table data access.
func NewKnowledgeDao(handlers ...gdb.ModelHandler) *KnowledgeDao {
	return &KnowledgeDao{
		group:    "default",
		table:    "knowledge",
		columns:  knowledgeColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *KnowledgeDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *KnowledgeDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *KnowledgeDao) Columns() KnowledgeColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *KnowledgeDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *KnowledgeDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *KnowledgeDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
