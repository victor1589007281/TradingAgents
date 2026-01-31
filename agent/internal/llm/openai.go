// Package llm LLM 模型封装
package llm

import (
	"context"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/tradingagents/agent/internal/config"
)

// ChatModelWrapper LLM 模型封装
type ChatModelWrapper struct {
	model     model.ChatModel
	modelName string
}

// NewChatModel 创建新的 ChatModel
func NewChatModel(ctx context.Context, cfg *config.Config, modelName string) (*ChatModelWrapper, error) {
	// 使用 OpenAI 兼容的模型配置
	chatModel, err := NewOpenAIChatModel(ctx, &OpenAIConfig{
		APIKey:  cfg.APIKey,
		BaseURL: cfg.BackendURL,
		Model:   modelName,
	})
	if err != nil {
		return nil, err
	}

	return &ChatModelWrapper{
		model:     chatModel,
		modelName: modelName,
	}, nil
}

// Generate 生成回复
func (c *ChatModelWrapper) Generate(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	return c.model.Generate(ctx, messages, opts...)
}

// Stream 流式生成
func (c *ChatModelWrapper) Stream(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	return c.model.Stream(ctx, messages, opts...)
}

// BindTools 绑定工具
func (c *ChatModelWrapper) BindTools(tools []schema.ToolInfo) model.ChatModel {
	return c.model.BindTools(tools)
}

// GetModelName 获取模型名称
func (c *ChatModelWrapper) GetModelName() string {
	return c.modelName
}

// OpenAIConfig OpenAI 配置
type OpenAIConfig struct {
	APIKey  string
	BaseURL string
	Model   string
}

// NewOpenAIChatModel 创建 OpenAI ChatModel
func NewOpenAIChatModel(ctx context.Context, cfg *OpenAIConfig) (model.ChatModel, error) {
	// 这里使用简化的实现，实际项目中应该使用 eino-ext/components/model/openai
	return &SimpleChatModel{
		apiKey:  cfg.APIKey,
		baseURL: cfg.BaseURL,
		model:   cfg.Model,
	}, nil
}

// SimpleChatModel 简化的 ChatModel 实现
type SimpleChatModel struct {
	apiKey  string
	baseURL string
	model   string
}

// Generate 实现 ChatModel 接口
func (s *SimpleChatModel) Generate(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	// 调用 OpenAI API
	response, err := callOpenAI(ctx, s.apiKey, s.baseURL, s.model, messages)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Stream 实现 ChatModel 接口
func (s *SimpleChatModel) Stream(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	// 简化实现，不支持流式
	return nil, nil
}

// BindTools 实现 ChatModel 接口
func (s *SimpleChatModel) BindTools(tools []schema.ToolInfo) model.ChatModel {
	return &ToolBoundChatModel{
		SimpleChatModel: s,
		tools:           tools,
	}
}

// ToolBoundChatModel 绑定工具的 ChatModel
type ToolBoundChatModel struct {
	*SimpleChatModel
	tools []schema.ToolInfo
}

// Generate 实现 ChatModel 接口
func (t *ToolBoundChatModel) Generate(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	// 调用带工具的 OpenAI API
	response, err := callOpenAIWithTools(ctx, t.apiKey, t.baseURL, t.model, messages, t.tools)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Stream 实现 ChatModel 接口
func (t *ToolBoundChatModel) Stream(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	return nil, nil
}

// BindTools 实现 ChatModel 接口
func (t *ToolBoundChatModel) BindTools(tools []schema.ToolInfo) model.ChatModel {
	return &ToolBoundChatModel{
		SimpleChatModel: t.SimpleChatModel,
		tools:           tools,
	}
}
