// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Feedback is the golang structure for table feedback.
type Feedback struct {
	Id                   uint64      `json:"id"                   orm:"id"                     description:"主键ID"`      // 主键ID
	SessionId            string      `json:"sessionId"            orm:"session_id"             description:"用户会话ID"`    // 用户会话ID
	UserQuery            string      `json:"userQuery"            orm:"user_query"             description:"用户查询内容"`    // 用户查询内容
	RetrievedKnowledgeId string      `json:"retrievedKnowledgeId" orm:"retrieved_knowledge_id" description:"被检索到的知识ID"` // 被检索到的知识ID
	Action               string      `json:"action"               orm:"action"                 description:"反馈操作类型"`    // 反馈操作类型
	Timestamp            *gtime.Time `json:"timestamp"            orm:"timestamp"              description:"反馈时间"`      // 反馈时间
}
