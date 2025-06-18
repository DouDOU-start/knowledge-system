package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 添加反馈请求
type FeedbackAddReq struct {
	g.Meta      `path:"/feedback" method:"post" tags:"反馈管理" summary:"添加用户反馈"`
	SessionID   string `json:"session_id" dc:"用户会话ID，可选，不传则使用默认测试ID"`
	UserQuery   string `v:"required#用户查询不能为空" json:"user_query" dc:"用户查询内容"`
	KnowledgeID string `v:"required#知识ID不能为空" json:"knowledge_id" dc:"被检索到的知识ID"`
	Action      string `v:"required|in:like,dislike#操作类型必须是like或dislike" json:"action" dc:"反馈操作类型：like或dislike"`
}

// 添加反馈响应
type FeedbackAddRes struct {
	ID string `json:"id" dc:"反馈ID"`
}

// 查询反馈请求
type FeedbackListReq struct {
	g.Meta      `path:"/feedback/list" method:"get" tags:"反馈管理" summary:"查询用户反馈列表"`
	SessionID   string `json:"session_id" in:"query" dc:"按会话ID过滤"`
	KnowledgeID string `json:"knowledge_id" in:"query" dc:"按知识ID过滤"`
	Action      string `json:"action" in:"query" v:"in:,like,dislike#操作类型必须是like或dislike" dc:"按操作类型过滤"`
	StartTime   string `json:"start_time" in:"query" dc:"开始时间，格式:YYYY-MM-DD HH:MM:SS"`
	EndTime     string `json:"end_time" in:"query" dc:"结束时间，格式:YYYY-MM-DD HH:MM:SS"`
	Page        int    `json:"page" in:"query" d:"1" dc:"页码"`
	PageSize    int    `json:"page_size" in:"query" d:"10" v:"max:100#每页最多100条" dc:"每页数量"`
}

// 查询反馈响应
type FeedbackListRes struct {
	List  []FeedbackItem `json:"list" dc:"反馈列表"`
	Total int            `json:"total" dc:"总条数"`
	Page  int            `json:"page" dc:"当前页码"`
}

// 反馈项
type FeedbackItem struct {
	ID          string `json:"id" dc:"反馈ID"`
	SessionID   string `json:"session_id" dc:"会话ID"`
	UserQuery   string `json:"user_query" dc:"用户查询"`
	KnowledgeID string `json:"knowledge_id" dc:"知识ID"`
	Action      string `json:"action" dc:"操作类型"`
	Timestamp   string `json:"timestamp" dc:"反馈时间"`
}

// 触发更新请求
type FeedbackProcessReq struct {
	g.Meta `path:"/feedback/process" method:"post" tags:"反馈管理" summary:"手动触发反馈处理"`
}

// 触发更新响应
type FeedbackProcessRes struct {
	Success bool   `json:"success" dc:"是否成功"`
	Message string `json:"message" dc:"处理消息"`
}
