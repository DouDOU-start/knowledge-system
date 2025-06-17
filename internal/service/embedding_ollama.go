package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type OllamaEmbeddingConfig struct {
	BaseURL string
	Model   string
}

type OllamaEmbeddingClient struct {
	cfg OllamaEmbeddingConfig
}

func NewOllamaEmbeddingClient(cfg OllamaEmbeddingConfig) *OllamaEmbeddingClient {
	return &OllamaEmbeddingClient{cfg: cfg}
}

// Embed 调用Ollama生成向量
func (c *OllamaEmbeddingClient) Embed(ctx context.Context, text string) ([]float32, error) {
	body := map[string]interface{}{
		"model":  c.cfg.Model,
		"prompt": text,
	}
	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST", c.cfg.BaseURL+"/api/embeddings", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("ollama embedding返回错误: %s", resp.Status)
	}
	respBytes, _ := ioutil.ReadAll(resp.Body)
	var result struct {
		Embedding []float32 `json:"embedding"`
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, errors.New("解析Ollama embedding响应失败")
	}
	return result.Embedding, nil
}
