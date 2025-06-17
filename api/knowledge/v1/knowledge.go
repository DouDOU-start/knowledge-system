package v1

import "github.com/gogf/gf/v2/frame/g"

// KnowledgeItem 知识条目
type KnowledgeItem struct {
	ID      string `json:"id" v:""`
	Content string `json:"content" v:"required#内容不能为空"`
}

// LabelScore 标签分数
type LabelScore struct {
	LabelID string `json:"label_id" v:"required#标签ID不能为空"`
	Score   int    `json:"score" v:"min:0#分数不能为负数"`
}

// KnowledgeResult 知识检索结果
type KnowledgeResult struct {
	ID      string       `json:"id"`
	Content string       `json:"content"`
	Labels  []LabelScore `json:"labels"`
	Summary string       `json:"summary"`
}

// 批量导入
//
type BatchImportReq struct {
	g.Meta `path:"/batch_import" method:"post" tags:"Knowledge" summary:"批量导入知识条目"`
	Items  []KnowledgeItem `json:"items" v:"required|array#导入条目不能为空|导入条目必须为数组"`
}

type BatchImportRes struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// 单条内容标签打分
//
type ClassifyReq struct {
	g.Meta  `path:"/classify" method:"post" tags:"Knowledge" summary:"单条内容标签打分"`
	Content string `json:"content" v:"required#内容不能为空"`
}

type ClassifyRes struct {
	Labels  []LabelScore `json:"labels"`
	Summary string       `json:"summary"`
}

// 知识检索
//
type SearchReq struct {
	g.Meta `path:"/search" method:"post" tags:"Knowledge" summary:"知识检索"`
	Query  string `json:"query" v:"required#检索关键词不能为空"`
	Mode   string `json:"mode" v:"in:keyword,semantic,hybrid#检索模式必须是 keyword/semantic/hybrid 之一"`
	TopK   int    `json:"top_k" v:"min:1#返回结果数量必须大于0"`
}

type SearchRes struct {
	Items []KnowledgeResult `json:"items"`
}
