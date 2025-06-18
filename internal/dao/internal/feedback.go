// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// FeedbackDao is the data access object for the table feedback.
type FeedbackDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  FeedbackColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// FeedbackColumns defines and stores column names for the table feedback.
type FeedbackColumns struct {
	Id                   string // 主键ID
	SessionId            string // 用户会话ID
	UserQuery            string // 用户查询内容
	RetrievedKnowledgeId string // 被检索到的知识ID
	Action               string // 反馈操作类型
	Timestamp            string // 反馈时间
}

// feedbackColumns holds the columns for the table feedback.
var feedbackColumns = FeedbackColumns{
	Id:                   "id",
	SessionId:            "session_id",
	UserQuery:            "user_query",
	RetrievedKnowledgeId: "retrieved_knowledge_id",
	Action:               "action",
	Timestamp:            "timestamp",
}

// NewFeedbackDao creates and returns a new DAO object for table data access.
func NewFeedbackDao(handlers ...gdb.ModelHandler) *FeedbackDao {
	return &FeedbackDao{
		group:    "default",
		table:    "feedback",
		columns:  feedbackColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *FeedbackDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *FeedbackDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *FeedbackDao) Columns() FeedbackColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *FeedbackDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *FeedbackDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *FeedbackDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
