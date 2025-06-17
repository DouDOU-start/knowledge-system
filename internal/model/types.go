package model

import "github.com/gogf/gf/v2/os/gtime"

// LabelScore 标签分数
type LabelScore struct {
	LabelID string `json:"label_id"`
	Score   int    `json:"score"`
}

// KnowledgeItem 知识条目
type KnowledgeItem struct {
	ID        string       `json:"id"`
	Content   string       `json:"content"`
	Labels    []LabelScore `json:"labels"`
	Summary   string       `json:"summary"`
	Vector    []float32    `json:"vector,omitempty"`
	CreatedAt *gtime.Time  `json:"created_at"`
	UpdatedAt *gtime.Time  `json:"updated_at"`
}

// ImportTask 导入任务
type ImportTask struct {
	TaskID    string      `json:"task_id"`
	Status    string      `json:"status"`   // pending, processing, completed, failed
	Progress  int         `json:"progress"` // 0-100
	Total     int         `json:"total"`
	Processed int         `json:"processed"`
	Failed    int         `json:"failed"`
	Message   string      `json:"message,omitempty"`
	Items     []TaskItem  `json:"items,omitempty"`
	CreatedAt *gtime.Time `json:"created_at"`
	UpdatedAt *gtime.Time `json:"updated_at"`
}

// TaskItem 任务项
type TaskItem struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Status  string `json:"status"` // pending, processing, completed, failed
	Message string `json:"message,omitempty"`
}

// SearchResult 搜索结果
type SearchResult struct {
	ID      string       `json:"id"`
	Content string       `json:"content"`
	Labels  []LabelScore `json:"labels"`
	Summary string       `json:"summary"`
	Score   float64      `json:"score"`
}

// VectorSearchResult 向量搜索结果
type VectorSearchResult struct {
	ID      string                 `json:"id"`
	Score   float64                `json:"score"`
	Vector  []float32              `json:"vector,omitempty"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}
