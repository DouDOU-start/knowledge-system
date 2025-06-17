package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"

	"knowledge-system-api/internal/model"

	"github.com/gogf/gf/v2/frame/g"
	"gopkg.in/yaml.v3"
)

// ===================== 检索主流程 =====================

// SearchKnowledge 检索主入口，支持 keyword/semantic/hybrid 三种模式
func SearchKnowledge(ctx context.Context, query, mode string, topK int) ([]model.KnowledgeItem, error) {
	if topK <= 0 {
		topK = 5
	}
	labels, _, err := LLMClassifyByConfig(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("LLM标签分析失败: %w", err)
	}
	vector, err := Vectorize(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("向量化失败: %w", err)
	}

	switch mode {
	case "keyword":
		return keywordSearch(ctx, query, topK)
	case "semantic":
		return semanticSearch(ctx, vector, topK)
	case "hybrid":
		return hybridSearch(ctx, vector, labels, topK)
	default:
		return hybridSearch(ctx, vector, labels, topK)
	}
}

// ===================== Qdrant相关 =====================

// qdrantSearch 向量检索，返回id和分数列表
func qdrantSearch(ctx context.Context, vector []float32, topK int) ([]struct {
	ID    string
	Score float32
}, error) {
	cfg, err := loadSearchConfig()
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}
	body := map[string]interface{}{
		"vector":       vector,
		"top":          topK,
		"with_payload": false,
		"with_vector":  false,
	}
	jsonBody, _ := json.Marshal(body)
	url := fmt.Sprintf("%s/collections/%s/points/search", cfg.Qdrant.BaseURL, cfg.Qdrant.Collection)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("Qdrant请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Qdrant请求失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Qdrant返回错误状态码 %d: %s", resp.StatusCode, string(body))
	}
	var result struct {
		Result []struct {
			ID    string  `json:"id"`
			Score float32 `json:"score"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析Qdrant响应失败: %w", err)
	}
	var out []struct {
		ID    string
		Score float32
	}
	for _, p := range result.Result {
		out = append(out, struct {
			ID    string
			Score float32
		}{ID: p.ID, Score: p.Score})
	}
	return out, nil
}

// ===================== MySQL相关 =====================

// getKnowledgeByIDs 批量查MySQL详情
func getKnowledgeByIDs(ctx context.Context, ids []string) ([]model.KnowledgeItem, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var items []struct {
		ID      string `json:"id"`
		Content string `json:"content"`
		Labels  string `json:"labels"`
		Summary string `json:"summary"`
	}
	err := g.Model("knowledge").
		Where("id IN(?)", ids).
		Scan(&items)
	if err != nil {
		return nil, fmt.Errorf("查询MySQL失败: %w", err)
	}
	result := make([]model.KnowledgeItem, 0, len(items))
	for _, item := range items {
		var labels []model.LabelScore
		if err := json.Unmarshal([]byte(item.Labels), &labels); err != nil {
			return nil, fmt.Errorf("解析标签JSON失败: %w", err)
		}
		result = append(result, model.KnowledgeItem{
			ID:      item.ID,
			Content: item.Content,
			Labels:  labels,
			Summary: item.Summary,
		})
	}
	// 按照传入的ids顺序排序
	idMap := make(map[string]int)
	for i, id := range ids {
		idMap[id] = i
	}
	sort.Slice(result, func(i, j int) bool {
		return idMap[result[i].ID] < idMap[result[j].ID]
	})
	return result, nil
}

// ===================== 检索模式实现 =====================

// semanticSearch 只用Qdrant分数
func semanticSearch(ctx context.Context, vector []float32, topK int) ([]model.KnowledgeItem, error) {
	points, err := qdrantSearch(ctx, vector, topK)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, p := range points {
		ids = append(ids, p.ID)
	}
	items, err := getKnowledgeByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	// 按Qdrant分数排序
	id2score := map[string]float32{}
	for _, p := range points {
		id2score[p.ID] = p.Score
	}
	for i := range items {
		items[i].Vector = []float32{id2score[items[i].ID]} // 用Vector字段临时存分数
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Vector[0] > items[j].Vector[0] })
	if len(items) > topK {
		items = items[:topK]
	}
	return items, nil
}

// hybridSearch 语义分数+标签分数融合
func hybridSearch(ctx context.Context, vector []float32, queryLabels []model.LabelScore, topK int) ([]model.KnowledgeItem, error) {
	cfg, err := loadSearchConfig()
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}
	points, err := qdrantSearch(ctx, vector, topK*2)
	if err != nil {
		return nil, fmt.Errorf("向量检索失败: %w", err)
	}
	if len(points) == 0 {
		return nil, nil
	}
	var ids []string
	for _, p := range points {
		ids = append(ids, p.ID)
	}
	items, err := getKnowledgeByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("获取知识条目详情失败: %w", err)
	}
	id2score := map[string]float32{}
	for _, p := range points {
		id2score[p.ID] = p.Score
	}
	for i := range items {
		labelScore := CalcLabelScore(queryLabels, items[i].Labels)
		finalScore := id2score[items[i].ID]*cfg.Hybrid.Alpha + labelScore*cfg.Hybrid.Beta
		items[i].Vector = []float32{finalScore}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Vector[0] > items[j].Vector[0] })
	if len(items) > topK {
		items = items[:topK]
	}
	return items, nil
}

// keywordSearch 关键词检索（MySQL全文索引）
func keywordSearch(ctx context.Context, query string, topK int) ([]model.KnowledgeItem, error) {
	var items []struct {
		ID      string  `json:"id"`
		Content string  `json:"content"`
		Labels  string  `json:"labels"`
		Summary string  `json:"summary"`
		Score   float64 `json:"score"`
	}
	err := g.Model("knowledge").
		Fields("*, MATCH(content) AGAINST(? IN NATURAL LANGUAGE MODE) as score", query).
		Where("MATCH(content) AGAINST(? IN NATURAL LANGUAGE MODE)", query).
		Order("score DESC").
		Limit(topK).
		Scan(&items)
	if err != nil {
		return nil, fmt.Errorf("MySQL全文检索失败: %w", err)
	}
	result := make([]model.KnowledgeItem, 0, len(items))
	for _, item := range items {
		var labels []model.LabelScore
		if err := json.Unmarshal([]byte(item.Labels), &labels); err != nil {
			return nil, fmt.Errorf("解析标签JSON失败: %w", err)
		}
		result = append(result, model.KnowledgeItem{
			ID:      item.ID,
			Content: item.Content,
			Labels:  labels,
			Summary: item.Summary,
			Vector:  []float32{float32(item.Score)},
		})
	}
	return result, nil
}

// ===================== 配置与工具 =====================

type searchConfig struct {
	Hybrid struct {
		Alpha float32 `yaml:"alpha"`
		Beta  float32 `yaml:"beta"`
	} `yaml:"hybrid"`
	Qdrant struct {
		BaseURL    string `yaml:"base_url"`
		Collection string `yaml:"collection"`
	} `yaml:"qdrant"`
}

var defaultSearchConfig = searchConfig{
	Hybrid: struct {
		Alpha float32 `yaml:"alpha"`
		Beta  float32 `yaml:"beta"`
	}{
		Alpha: 0.7,
		Beta:  0.3,
	},
	Qdrant: struct {
		BaseURL    string `yaml:"base_url"`
		Collection string `yaml:"collection"`
	}{
		BaseURL:    "http://localhost:6333",
		Collection: "knowledge",
	},
}

// loadSearchConfig 读取检索相关配置
func loadSearchConfig() (*searchConfig, error) {
	cfg := defaultSearchConfig
	b, err := os.ReadFile("hack/config.yaml")
	if err != nil {
		g.Log().Warningf(nil, "加载配置文件失败，使用默认配置: %v", err)
		return &cfg, nil
	}
	var fullCfg struct {
		Search searchConfig `yaml:"search"`
	}
	if err := yaml.Unmarshal(b, &fullCfg); err != nil {
		g.Log().Warningf(nil, "解析配置文件失败，使用默认配置: %v", err)
		return &cfg, nil
	}
	if fullCfg.Search.Hybrid.Alpha != 0 {
		cfg.Hybrid.Alpha = fullCfg.Search.Hybrid.Alpha
	}
	if fullCfg.Search.Hybrid.Beta != 0 {
		cfg.Hybrid.Beta = fullCfg.Search.Hybrid.Beta
	}
	if fullCfg.Search.Qdrant.BaseURL != "" {
		cfg.Qdrant.BaseURL = fullCfg.Search.Qdrant.BaseURL
	}
	if fullCfg.Search.Qdrant.Collection != "" {
		cfg.Qdrant.Collection = fullCfg.Search.Qdrant.Collection
	}
	return &cfg, nil
}
