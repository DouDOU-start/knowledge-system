package service

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
)

var (
	embeddingClientInstance EmbeddingClient
	embeddingOnce           sync.Once
)

type EmbeddingConfig struct {
	Backend string `yaml:"backend" json:"backend"`
	Ollama  struct {
		BaseURL string `yaml:"base_url" json:"base_url"`
		Model   string `yaml:"model" json:"model"`
	} `yaml:"ollama" json:"ollama"`
}

// LoadEmbeddingConfig 读取embedding相关配置
func LoadEmbeddingConfig(ctx context.Context) (*EmbeddingConfig, error) {
	var cfg EmbeddingConfig
	if err := g.Cfg().MustGet(ctx, "embedding").Scan(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// GetEmbeddingClient 工厂方法，根据配置返回对应实现
func GetEmbeddingClient() EmbeddingClient {
	embeddingOnce.Do(func() {
		ctx := context.Background()
		cfg, err := LoadEmbeddingConfig(ctx)
		if err != nil {
			g.Log().Errorf(ctx, "加载embedding配置失败: %v", err)
			// 使用默认配置
			cfg = &EmbeddingConfig{
				Backend: "ollama",
				Ollama: struct {
					BaseURL string `yaml:"base_url" json:"base_url"`
					Model   string `yaml:"model" json:"model"`
				}{
					BaseURL: "http://localhost:11434",
					Model:   "nomic-embed-text",
				},
			}
		}

		switch cfg.Backend {
		case "ollama":
			embeddingClientInstance = NewOllamaEmbeddingClient(OllamaEmbeddingConfig{
				BaseURL: cfg.Ollama.BaseURL,
				Model:   cfg.Ollama.Model,
			})
		// 预留其他后端
		default:
			g.Log().Errorf(ctx, "不支持的embedding后端: %s，使用默认ollama后端", cfg.Backend)
			embeddingClientInstance = NewOllamaEmbeddingClient(OllamaEmbeddingConfig{
				BaseURL: "http://localhost:11434",
				Model:   "nomic-embed-text",
			})
		}
	})
	return embeddingClientInstance
}

// EmbeddingClient 向量化统一接口
type EmbeddingClient interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}
