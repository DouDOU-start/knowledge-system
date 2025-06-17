package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"knowledge-system-api/internal/model"
	"net/http"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

var (
	qdrantConfigInstance *QdrantConfig
	qdrantOnce           sync.Once
)

// QdrantConfig Qdrant配置
type QdrantConfig struct {
	URL        string `yaml:"url" json:"url"`
	Collection string `yaml:"collection" json:"collection"`
	Dimension  int    `yaml:"dimension" json:"dimension"`
}

// LoadQdrantConfig 加载Qdrant配置
func LoadQdrantConfig(ctx context.Context) (*QdrantConfig, error) {
	var cfg QdrantConfig
	if err := g.Cfg().MustGet(ctx, "qdrant").Scan(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// GetQdrantConfig 获取Qdrant配置单例
func GetQdrantConfig() *QdrantConfig {
	qdrantOnce.Do(func() {
		ctx := context.Background()
		cfg, err := LoadQdrantConfig(ctx)
		if err != nil {
			g.Log().Errorf(ctx, "加载qdrant配置失败: %v", err)
			// 使用默认配置
			cfg = &QdrantConfig{
				URL:        "http://localhost:6333",
				Collection: "knowledge",
				Dimension:  768,
			}
		}
		qdrantConfigInstance = cfg
	})
	return qdrantConfigInstance
}

// QdrantPoint Qdrant单条数据结构
// 参考Qdrant API文档
// https://qdrant.tech/documentation/concepts/points/
type QdrantPoint struct {
	ID      string                 `json:"id"`
	Vector  []float32              `json:"vector"`
	Payload map[string]interface{} `json:"payload"`
}

type QdrantUpsertRequest struct {
	Points []QdrantPoint `json:"points"`
}

type QdrantUpsertResponse struct {
	Status string `json:"status"`
}

// QdrantUpsert 将知识条目写入Qdrant向量库
func QdrantUpsert(id string, vector []float32, content string, repoName string, labels []model.LabelScore, summary string) error {
	cfg := GetQdrantConfig()

	// 1. 检查向量
	if len(vector) == 0 {
		return fmt.Errorf("QdrantUpsert: 向量为空")
	}
	// 2. 转换 labels 为[]map[string]interface{}
	labelArr := make([]map[string]interface{}, 0, len(labels))
	for _, l := range labels {
		labelArr = append(labelArr, map[string]interface{}{
			"label_id": l.LabelID,
			"score":    l.Score,
		})
	}
	point := QdrantPoint{
		ID:     id,
		Vector: vector,
		Payload: map[string]interface{}{
			"content":   content,
			"repo_name": repoName,
			"labels":    labelArr,
			"summary":   summary,
		},
	}
	upsertReq := QdrantUpsertRequest{
		Points: []QdrantPoint{point},
	}
	body, err := json.Marshal(upsertReq)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/collections/%s/points?wait=true", cfg.URL, cfg.Collection)
	// 新增：自动创建collection并重试
	retry := false
RETRY_UPSERT:
	req, err := http.NewRequest("PUT", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 && !retry {
		// collection不存在，自动创建
		if err := createQdrantCollection(len(vector)); err != nil {
			return fmt.Errorf("自动创建Qdrant collection失败: %w", err)
		}
		retry = true
		goto RETRY_UPSERT
	}
	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("qdrant upsert failed: %s, detail: %s", resp.Status, string(bodyBytes))
	}
	return nil
}

// QdrantSearch 向量搜索
func QdrantSearch(ctx context.Context, query string, vector []float32, repoName string, limit int) ([]model.VectorSearchResult, error) {
	cfg := GetQdrantConfig()

	if len(vector) == 0 {
		return nil, fmt.Errorf("QdrantSearch: 向量为空")
	}

	// 构建查询
	searchRequest := map[string]interface{}{
		"vector":       vector,
		"top":          limit,
		"with_payload": true,
	}

	// 如果指定了知识库名称，添加过滤条件
	if repoName != "" {
		searchRequest["filter"] = map[string]interface{}{
			"must": []map[string]interface{}{
				{
					"key": "repo_name",
					"match": map[string]interface{}{
						"value": repoName,
					},
				},
			},
		}
	}

	jsonBody, _ := json.Marshal(searchRequest)
	url := fmt.Sprintf("%s/collections/%s/points/search", cfg.URL, cfg.Collection)
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("Qdrant请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Qdrant请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("qdrant search failed: %s, detail: %s", resp.Status, string(bodyBytes))
	}

	var result struct {
		Result []struct {
			ID      string                 `json:"id"`
			Score   float64                `json:"score"`
			Payload map[string]interface{} `json:"payload"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析Qdrant响应失败: %w", err)
	}

	var out []model.VectorSearchResult
	for _, p := range result.Result {
		out = append(out, model.VectorSearchResult{
			ID:      p.ID,
			Score:   p.Score,
			Payload: p.Payload,
		})
	}

	return out, nil
}

// 新增：自动创建collection
func createQdrantCollection(vectorSize int) error {
	cfg := GetQdrantConfig()

	url := fmt.Sprintf("%s/collections/%s", cfg.URL, cfg.Collection)
	body := fmt.Sprintf(`{"vectors":{"size":%d,"distance":"Cosine"}}`, vectorSize)
	req, err := http.NewRequest("PUT", url, bytes.NewReader([]byte(body)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create collection failed: %s, detail: %s", resp.Status, string(bodyBytes))
	}
	return nil
}
