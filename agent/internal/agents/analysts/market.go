// Package analysts 分析师 Agents
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

// MarketAnalyst 技术分析师
type MarketAnalyst struct {
	llm      *llm.ChatModelWrapper
	registry *tools.ToolRegistry
}

// NewMarketAnalyst 创建技术分析师
func NewMarketAnalyst(chatModel *llm.ChatModelWrapper) *MarketAnalyst {
	registry := tools.NewToolRegistry()
	registry.Register(tools.GetStockDataTool)
	registry.Register(tools.GetIndicatorsTool)

	return &MarketAnalyst{
		llm:      chatModel,
		registry: registry,
	}
}

// SystemPrompt 系统提示词
const MarketAnalystSystemPrompt = `You are a market analyst tasked with analyzing market trends and technical indicators over the past week. 

Your responsibilities:
1. Analyze OHLCV (Open, High, Low, Close, Volume) data
2. Evaluate technical indicators (MACD, RSI, Moving Averages, Bollinger Bands)
3. Identify trends, support/resistance levels, and trading patterns
4. Provide detailed insights that may help traders make decisions

Please write a comprehensive report including:
- Price trend analysis
- Volume analysis
- Technical indicator signals
- Key support and resistance levels
- Overall market outlook

Make sure to append a Markdown table at the end of the report to organize key points.

Use the available tools:
- get_stock_data: Get OHLCV historical data
- get_indicators: Get technical indicators`

// Run 运行分析师
func (a *MarketAnalyst) Run(ctx context.Context, state *states.AgentState) (*states.AgentState, error) {
	ticker := state.CompanyOfInterest
	tradeDate := state.TradeDate

	// 构建消息
	messages := []*schema.Message{
		{
			Role:    schema.System,
			Content: MarketAnalystSystemPrompt,
		},
		{
			Role:    schema.User,
			Content: fmt.Sprintf("Analyze the market data for %s as of %s. Use the tools to gather data and provide a comprehensive technical analysis report.", ticker, tradeDate),
		},
	}

	// 获取工具信息
	toolInfos := a.registry.GetToolInfos()

	// 调用 LLM（绑定工具）
	boundLLM := a.llm.BindTools(toolInfos)
	response, err := boundLLM.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("market analyst LLM call failed: %w", err)
	}

	// 处理工具调用
	for len(response.ToolCalls) > 0 {
		// 执行工具调用
		for _, tc := range response.ToolCalls {
			var args map[string]interface{}
			if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
				continue
			}

			result, err := a.registry.ExecuteTool(ctx, tc.Function.Name, args)
			if err != nil {
				result = fmt.Sprintf("Error executing tool: %v", err)
			}

			// 添加工具结果消息
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

		// 再次调用 LLM
		response, err = boundLLM.Generate(ctx, messages)
		if err != nil {
			return nil, fmt.Errorf("market analyst follow-up LLM call failed: %w", err)
		}
	}

	// 更新状态
	state.MarketReport = response.Content
	state.AddMessage(response)

	return state, nil
}

// GetName 获取名称
func (a *MarketAnalyst) GetName() string {
	return "Market Analyst"
}
