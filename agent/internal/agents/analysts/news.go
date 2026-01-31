// Package analysts 新闻分析师
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

// NewsAnalyst 新闻分析师
type NewsAnalyst struct {
	llm      *llm.ChatModelWrapper
	registry *tools.ToolRegistry
}

// NewNewsAnalyst 创建新闻分析师
func NewNewsAnalyst(chatModel *llm.ChatModelWrapper) *NewsAnalyst {
	registry := tools.NewToolRegistry()
	registry.Register(tools.GetNewsTool)
	registry.Register(tools.GetGlobalNewsTool)
	registry.Register(tools.GetInsiderSentimentTool)
	registry.Register(tools.GetInsiderTransactionsTool)

	return &NewsAnalyst{
		llm:      chatModel,
		registry: registry,
	}
}

// SystemPrompt 系统提示词
const NewsAnalystSystemPrompt = `You are a news analyst tasked with analyzing recent news and trends over the past week.

Your responsibilities:
1. Analyze company-specific news and announcements
2. Evaluate macroeconomic news and global events
3. Assess insider trading activities and sentiment
4. Identify potential market-moving events

Please write a comprehensive news report including:
- Key company news and developments
- Relevant industry news
- Macroeconomic factors and global events
- Insider activity analysis
- Potential impact on investment decisions

Make sure to append a Markdown table at the end of the report to organize key news items and their potential impact.

Use the available tools:
- get_news: Get company-specific news
- get_global_news: Get macroeconomic and market news
- get_insider_sentiment: Get insider sentiment data
- get_insider_transactions: Get insider trading records`

// Run 运行分析师
func (a *NewsAnalyst) Run(ctx context.Context, state *states.AgentState) (*states.AgentState, error) {
	ticker := state.CompanyOfInterest
	tradeDate := state.TradeDate

	// 构建消息
	messages := []*schema.Message{
		{
			Role:    schema.System,
			Content: NewsAnalystSystemPrompt,
		},
		{
			Role:    schema.User,
			Content: fmt.Sprintf("Analyze the news for %s as of %s. Use the tools to gather data and provide a comprehensive news analysis report covering company news, global events, and insider activities.", ticker, tradeDate),
		},
	}

	// 获取工具信息
	toolInfos := a.registry.GetToolInfos()

	// 调用 LLM
	boundLLM := a.llm.BindTools(toolInfos)
	response, err := boundLLM.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("news analyst LLM call failed: %w", err)
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
			return nil, fmt.Errorf("news analyst follow-up LLM call failed: %w", err)
		}
	}

	// 更新状态
	state.NewsReport = response.Content
	state.AddMessage(response)

	return state, nil
}

// GetName 获取名称
func (a *NewsAnalyst) GetName() string {
	return "News Analyst"
}
