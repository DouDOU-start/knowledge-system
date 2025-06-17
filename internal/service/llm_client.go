package service

import (
	"context"
	"encoding/json"
	"fmt"
	"knowledge-system-api/internal/model"
	"os"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

// LLMClient 大模型推理统一接口
// 所有推理后端（如Ollama、阿里云等）都需实现该接口
// 这里Classify接口兼容原有标签打分和摘要
type LLMClient interface {
	Classify(ctx context.Context, content string) (labels []model.LabelScore, summary string, err error)
}

var (
	llmClientInstance LLMClient
	llmOnce           sync.Once
)

type LLMConfig struct {
	Backend string `yaml:"backend" json:"backend"`
	Ollama  struct {
		BaseURL    string `yaml:"base_url" json:"base_url"`
		Model      string `yaml:"model" json:"model"`
		PromptPath string `yaml:"prompt_path" json:"prompt_path"`
	} `yaml:"ollama" json:"ollama"`
}

// LoadLLMConfig 读取llm相关配置
func LoadLLMConfig(ctx context.Context) (*LLMConfig, error) {
	var cfg LLMConfig
	if err := g.Cfg().MustGet(ctx, "llm").Scan(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// GetLLMClient 工厂方法，返回LLMClient实例
func GetLLMClient() LLMClient {
	llmOnce.Do(func() {
		ctx := context.Background()
		cfg, err := LoadLLMConfig(ctx)
		if err != nil {
			g.Log().Errorf(ctx, "加载llm配置失败: %v", err)
			// 使用默认配置
			cfg = &LLMConfig{
				Backend: "ollama",
				Ollama: struct {
					BaseURL    string `yaml:"base_url" json:"base_url"`
					Model      string `yaml:"model" json:"model"`
					PromptPath string `yaml:"prompt_path" json:"prompt_path"`
				}{
					BaseURL:    "http://localhost:11434",
					Model:      "llama3",
					PromptPath: "resource/prompts/classify.txt",
				},
			}
		}

		switch cfg.Backend {
		case "ollama":
			llmClientInstance = &LangchainOllamaLLMAdapter{
				BaseURL:    cfg.Ollama.BaseURL,
				Model:      cfg.Ollama.Model,
				PromptPath: cfg.Ollama.PromptPath,
			}
		// 预留其他后端
		default:
			g.Log().Errorf(ctx, "不支持的llm后端: %s，使用默认ollama后端", cfg.Backend)
			llmClientInstance = &LangchainOllamaLLMAdapter{
				BaseURL:    "http://localhost:11434",
				Model:      "llama3",
				PromptPath: "resource/prompts/classify.txt",
			}
		}
	})
	return llmClientInstance
}

// LangchainOllamaLLMAdapter 适配器，兼容原有Classify接口
// prompt模板读取和拼接逻辑与原有一致
type LangchainOllamaLLMAdapter struct {
	BaseURL    string
	Model      string
	PromptPath string
}

func (a *LangchainOllamaLLMAdapter) Classify(ctx context.Context, content string) (labels []model.LabelScore, summary string, err error) {
	// 日志记录当前配置
	glog.Debugf(ctx, "Classify: BaseURL=%s, Model=%s, PromptPath=%s", a.BaseURL, a.Model, a.PromptPath)

	promptTmpl, err := LoadPromptTemplate(a.PromptPath)
	if err != nil {
		glog.Errorf(ctx, "加载Prompt模板失败: %v", err)
		return nil, "", err
	}
	prompt := promptTmpl + content
	llm, err := ollama.New(
		ollama.WithModel(a.Model),
		ollama.WithServerURL(a.BaseURL),
	)
	if err != nil {
		glog.Errorf(ctx, "ollama.New error: %v", err)
		return nil, "", err
	}
	if llm == nil {
		glog.Errorf(ctx, "ollama.New returned nil LLM")
		return nil, "", fmt.Errorf("LLM 初始化失败: BaseURL=%s, Model=%s", a.BaseURL, a.Model)
	}
	resp, err := llm.Call(ctx, prompt,
		llms.WithTemperature(0.8),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			// fmt.Print(string(chunk))
			return nil
		}),
	)
	if err != nil {
		glog.Errorf(ctx, "llm.Call error: %v", err)
		return nil, "", err
	}

	// 新增：提取resp中的JSON部分
	jsonStr, extractErr := ExtractJSONFromLLMResponse(resp)
	if extractErr != nil {
		glog.Errorf(ctx, "提取大模型JSON失败: %v, resp=%s", extractErr, resp)
		return nil, "", fmt.Errorf("提取大模型JSON失败: %w", extractErr)
	}

	// 解析大模型返回的JSON
	var parsed struct {
		C1TopicScores map[string]int `json:"C1_Topic_Scores"`
		C2TypeScores  map[string]int `json:"C2_Type_Scores"`
		Summary       string         `json:"summary"`
	}
	if err = json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		glog.Errorf(ctx, "解析大模型JSON失败: %v, resp=%s", err, jsonStr)
		err = fmt.Errorf("解析大模型JSON失败: %w", err)
		return
	}
	for k, v := range parsed.C1TopicScores {
		labels = append(labels, model.LabelScore{LabelID: k, Score: v})
	}
	for k, v := range parsed.C2TypeScores {
		labels = append(labels, model.LabelScore{LabelID: k, Score: v})
	}
	summary = parsed.Summary
	return
}

// LoadPromptTemplate 读取prompt模板内容
func LoadPromptTemplate(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
