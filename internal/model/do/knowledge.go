// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Knowledge is the golang structure of table knowledge for DAO operations like Where/Data.
type Knowledge struct {
	g.Meta    `orm:"table:knowledge, do:true"`
	Id        interface{} // 唯一ID，服务端生成UUID
	RepoName  interface{} // 知识库名称
	Content   interface{} // 知识内容
	Labels    interface{} // 标签分数数组
	Summary   interface{} // 内容摘要
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
