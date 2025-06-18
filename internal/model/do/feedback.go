// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Feedback is the golang structure of table feedback for DAO operations like Where/Data.
type Feedback struct {
	g.Meta               `orm:"table:feedback, do:true"`
	Id                   interface{} // 主键ID
	SessionId            interface{} // 用户会话ID
	UserQuery            interface{} // 用户查询内容
	RetrievedKnowledgeId interface{} // 被检索到的知识ID
	Action               interface{} // 反馈操作类型
	Timestamp            *gtime.Time // 反馈时间
}
