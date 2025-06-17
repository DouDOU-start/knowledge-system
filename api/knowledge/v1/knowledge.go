package v1

import "github.com/gogf/gf/v2/frame/g"

// KnowledgeItem 知识条目
type KnowledgeItem struct {
	ID      string `json:"id,omitempty" v:""` // ID 字段设为可选，系统会自动生成
	Content string `json:"content" v:"required#内容不能为空"`
}

// LabelScore 标签分数
type LabelScore struct {
	LabelID string `json:"label_id" v:"required#标签ID不能为空"`
	Score   int    `json:"score" v:"min:0#分数不能为负数"`
}

// KnowledgeResult 知识检索结果
type KnowledgeResult struct {
	ID       string       `json:"id"`
	RepoName string       `json:"repo_name"`
	Content  string       `json:"content"`
	Labels   []LabelScore `json:"labels"`
	Summary  string       `json:"summary"`
}

// 批量导入
//
type BatchImportReq struct {
	g.Meta   `path:"/batch_import" method:"post" tags:"Knowledge" summary:"批量导入知识条目"`
	RepoName string          `json:"repo_name" v:"required#知识库名称不能为空"`
	Items    []KnowledgeItem `json:"items" v:"required|array#导入条目不能为空|导入条目必须为数组"`
}

type BatchImportRes struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// 批量异步导入
//
type BatchImportAsyncReq struct {
	g.Meta   `path:"/batch_import_async" method:"post" tags:"Knowledge" summary:"批量异步导入知识条目"`
	RepoName string          `json:"repo_name" v:"required#知识库名称不能为空"`
	Items    []KnowledgeItem `json:"items" v:"required|array#导入条目不能为空|导入条目必须为数组"`
}

type BatchImportAsyncRes struct {
	TaskID  string `json:"task_id"` // 任务ID，用于查询进度
	Message string `json:"message"` // 提示信息
}

// 任务状态
//
type TaskStatusReq struct {
	g.Meta `path:"/task/:task_id" method:"get" tags:"Knowledge" summary:"查询任务状态"`
	TaskID string `json:"task_id" in:"path" v:"required#任务ID不能为空"`
}

type TaskStatusRes struct {
	TaskID    string `json:"task_id"`           // 任务ID
	Status    string `json:"status"`            // 任务状态：pending, processing, completed, failed
	Progress  int    `json:"progress"`          // 处理进度，0-100
	Total     int    `json:"total"`             // 总条目数
	Processed int    `json:"processed"`         // 已处理条目数
	Failed    int    `json:"failed"`            // 失败条目数
	Message   string `json:"message,omitempty"` // 任务相关信息
}

// 单条内容标签打分
//
type ClassifyReq struct {
	g.Meta   `path:"/classify" method:"post" tags:"Knowledge" summary:"单条内容标签打分"`
	Content  string `json:"content" v:"required#内容不能为空"`
	RepoName string `json:"repo_name" v:"required#知识库名称不能为空"`
}

type ClassifyRes struct {
	Labels  []LabelScore `json:"labels"`
	Summary string       `json:"summary"`
}

// 知识检索
//
type SearchReq struct {
	g.Meta   `path:"/search" method:"post" tags:"Knowledge" summary:"知识检索"`
	Query    string `json:"query" v:"required#检索关键词不能为空"`
	RepoName string `json:"repo_name" v:"#知识库名称，不填则搜索所有知识库"`
	Mode     string `json:"mode" v:"in:keyword,semantic,hybrid#检索模式必须是 keyword/semantic/hybrid 之一"`
	TopK     int    `json:"top_k" v:"min:1#返回结果数量必须大于0"`
}

type SearchRes struct {
	Items []KnowledgeResult `json:"items"`
}

// 获取所有知识库
//
type GetReposReq struct {
	g.Meta `path:"/repos" method:"get" tags:"Knowledge" summary:"获取所有知识库"`
}

type GetReposRes struct {
	Repos []string `json:"repos"`
}
