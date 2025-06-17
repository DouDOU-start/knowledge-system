package service

import (
	"context"

	"github.com/tmc/langchaingo/llms/ollama"
)

type LangchainOllamaClient struct {
	llm *ollama.LLM
}

// NewLangchainOllamaClient 构造函数
func NewLangchainOllamaClient(baseURL, model string) (*LangchainOllamaClient, error) {
	llm, err := ollama.New(
		ollama.WithModel(model),
		ollama.WithServerURL(baseURL),
	)
	if err != nil {
		return nil, err
	}
	return &LangchainOllamaClient{llm: llm}, nil
}

// Generate 用于通用内容生成
func (c *LangchainOllamaClient) Generate(ctx context.Context, prompt string) (string, error) {
	return c.llm.Call(ctx, prompt, nil)
}
