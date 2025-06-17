// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Knowledge is the golang structure for table knowledge.
type Knowledge struct {
	Id        string      `json:"id"        orm:"id"         description:"唯一ID，服务端生成UUID"` // 唯一ID，服务端生成UUID
	RepoName  string      `json:"repoName"  orm:"repo_name"  description:"知识库名称"`          // 知识库名称
	Content   string      `json:"content"   orm:"content"    description:"知识内容"`           // 知识内容
	Labels    string      `json:"labels"    orm:"labels"     description:"标签分数数组"`         // 标签分数数组
	Summary   string      `json:"summary"   orm:"summary"    description:"内容摘要"`           // 内容摘要
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`           // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`           // 更新时间
}
