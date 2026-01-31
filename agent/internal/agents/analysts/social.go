// Package analysts 社交媒体分析师
package analysts

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/schema"
	"github.com/tradingagents/agent/internal/llm"
	"github.com/tradingagents/agent/internal/states"
	"github.com/tradingagents/agent/internal/tools"
)

// SocialMediaAnalyst 社交媒体/情绪分析师
type SocialMediaAnalyst struct {
	llm      *llm.ChatModelWrapper
	registry *tools.ToolRegistry
}

// NewSocialMediaAnalyst 创建社交媒体分析师
func NewSocialMediaAnalyst(chatModel *llm.ChatModelWrapper) *SocialMediaAnalyst {
	registry := tools.NewToolRegistry()
	registry.Register(tools.GetNewsTool)

	return &SocialMediaAnalyst{
		llm:      chatModel,
		registry: registry,
	}
}

// SystemPrompt 系统提示词
const SocialMediaAnalystSystemPrompt = `You are a social media and sentiment analyst tasked with analyzing public sentiment and social media trends about a company.

Your responsibilities:
1. Analyze social media sentiment (positive, negative, neutral)
2. Identify trending topics and discussions
3. Evaluate public perception and market mood
4. Track sentiment changes over time

Please write a comprehensive sentiment report including:
- Overall sentiment score and trend
- Key themes in public discussion
- Notable mentions or viral content
- Sentiment comparison with historical patterns
- Potential impact on stock price

Make sure to append a Markdown table at the end of the report to organize key sentiment metrics.

Use the available tools:
- get_news: Get news and social media mentions for sentiment analysis`

// Run 运行分析师
func (a *SocialMediaAnalyst) Run(ctx context.Context, state *states.AgentState) (*states.AgentState, error) {
	ticker := state.CompanyOfInterest
	tradeDate := state.TradeDate

	// 构建消息
	messages := []*schema.Message{
		{
			Role:    schema.System,
			Content: SocialMediaAnalystSystemPrompt,
		},
		{
			Role:    schema.User,
			Content: fmt.Sprintf("Analyze the social media sentiment for %s as of %s. Use the tools to gather data and provide a comprehensive sentiment analysis report.", ticker, tradeDate),
		},
	}

	// 获取工具信息
	toolInfos := a.registry.GetToolInfos()

	// 调用 LLM
	boundLLM := a.llm.BindTools(toolInfos)
	response, err := boundLLM.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("social media analyst LLM call failed: %w", err)
	}

	// 处理工具调用
	for len(response.ToolCalls) > 0 {
		for _, tc := range response.ToolCalls {
			var args map[string]interface{}
			if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
				continue
			}

			result, err := a.registry.ExecuteTool(ctx, tc.Function.Name, args)
			if err != nil {
				result = fmt.Sprintf("Error executing tool: %v", err)
			}

			messages = append(messages, &schema.Message{
				Role:    schema.Assistant,
				Content: response.Content,
			})
			messages = append(messages, &schema.Message{
				Role:       schema.Tool,
				Content:    result,
				ToolCallID: tc.ID,
			})
		}

		response, err = boundLLM.Generate(ctx, messages)
		if err != nil {
			return nil, fmt.Errorf("social media analyst follow-up LLM call failed: %w", err)
		}
	}

	// 更新状态
	state.SentimentReport = response.Content
	state.AddMessage(response)

	return state, nil
}

// GetName 获取名称
func (a *SocialMediaAnalyst) GetName() string {
	return "Social Media Analyst"
}
