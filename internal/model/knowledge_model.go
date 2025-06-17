package model

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KnowledgeItem 知识条目业务模型
type KnowledgeItem struct {
	ID        string       `json:"id"`               // 唯一ID
	RepoName  string       `json:"repo_name"`        // 知识库名称
	Content   string       `json:"content"`          // 知识内容
	Labels    []LabelScore `json:"labels"`           // 标签分数数组
	Summary   string       `json:"summary"`          // 内容摘要
	Vector    []float32    `json:"vector,omitempty"` // 向量，用于临时存储分数
	CreatedAt *gtime.Time  `json:"created_at"`       // 创建时间
	UpdatedAt *gtime.Time  `json:"updated_at"`       // 更新时间
}

// LabelScore 标签分数
type LabelScore struct {
	LabelID string `json:"label_id"` // 标签ID
	Score   int    `json:"score"`    // 分数
}

// SearchResult 搜索结果
type SearchResult struct {
	ID       string       `json:"id"`        // 条目ID
	RepoName string       `json:"repo_name"` // 知识库名称
	Content  string       `json:"content"`   // 知识内容
	Labels   []LabelScore `json:"labels"`    // 标签分数数组
	Summary  string       `json:"summary"`   // 内容摘要
	Score    float64      `json:"score"`     // 搜索匹配分数
}

// VectorSearchResult 向量搜索结果
type VectorSearchResult struct {
	ID      string                 `json:"id"`      // 条目ID
	Score   float64                `json:"score"`   // 搜索分数
	Payload map[string]interface{} `json:"payload"` // 负载数据
}

// ImportTask 导入任务
type ImportTask struct {
	TaskID    string      `json:"task_id"`    // 任务ID
	Status    string      `json:"status"`     // 任务状态：pending, processing, completed, failed, completed_with_errors
	Progress  uint        `json:"progress"`   // 处理进度，0-100
	Total     uint        `json:"total"`      // 总条目数
	Processed uint        `json:"processed"`  // 已处理条目数
	Failed    uint        `json:"failed"`     // 失败条目数
	Message   string      `json:"message"`    // 任务相关信息
	Items     []TaskItem  `json:"items"`      // 任务条目（用于展示，非持久化）
	CreatedAt *gtime.Time `json:"created_at"` // 创建时间
	UpdatedAt *gtime.Time `json:"updated_at"` // 更新时间
}

// TaskItem 任务条目
type TaskItem struct {
	ID           int64  `json:"id,omitempty"`      // 条目ID，可选，数据库自增
	TaskID       string `json:"task_id,omitempty"` // 所属任务ID
	RepoName     string `json:"repo_name"`         // 知识库名称
	Content      string `json:"content"`           // 知识内容
	Status       string `json:"status"`            // 处理状态：pending, processing, completed, failed
	ErrorMessage string `json:"error_message"`     // 处理失败时的错误信息
}
