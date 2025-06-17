package service

import (
	"context"
	"sync"
)

var (
	embeddingClientInstance EmbeddingClient
	embeddingOnce           sync.Once
)

const defaultEmbeddingConfigPath = "hack/config.yaml"

type EmbeddingConfig struct {
	Backend string `yaml:"backend"`
	Ollama  struct {
		BaseURL string `yaml:"base_url"`
		Model   string `yaml:"model"`
	} `yaml:"ollama"`
}

// LoadEmbeddingConfig 读取embedding相关配置
func LoadEmbeddingConfig() (*EmbeddingConfig, error) {
	var cfg struct {
		Embedding EmbeddingConfig `yaml:"embedding"`
	}
	err := LoadYAMLConfig(defaultEmbeddingConfigPath, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg.Embedding, nil
}

// GetEmbeddingClient 工厂方法，根据配置返回对应实现
func GetEmbeddingClient() EmbeddingClient {
	embeddingOnce.Do(func() {
		cfg, err := LoadEmbeddingConfig()
		if err != nil {
			panic("加载embedding配置失败: " + err.Error())
		}
		switch cfg.Backend {
		case "ollama":
			embeddingClientInstance = NewOllamaEmbeddingClient(OllamaEmbeddingConfig{
				BaseURL: cfg.Ollama.BaseURL,
				Model:   cfg.Ollama.Model,
			})
		// 预留其他后端
		default:
			panic("不支持的embedding后端: " + cfg.Backend)
		}
	})
	return embeddingClientInstance
}

// EmbeddingClient 向量化统一接口
type EmbeddingClient interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}
