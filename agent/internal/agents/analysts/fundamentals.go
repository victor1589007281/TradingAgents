// Package analysts 基本面分析师
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

// FundamentalsAnalyst 基本面分析师
type FundamentalsAnalyst struct {
	llm      *llm.ChatModelWrapper
	registry *tools.ToolRegistry
}

// NewFundamentalsAnalyst 创建基本面分析师
func NewFundamentalsAnalyst(chatModel *llm.ChatModelWrapper) *FundamentalsAnalyst {
	registry := tools.NewToolRegistry()
	registry.Register(tools.GetFundamentalsTool)
	registry.Register(tools.GetBalanceSheetTool)
	registry.Register(tools.GetCashflowTool)
	registry.Register(tools.GetIncomeStatementTool)

	return &FundamentalsAnalyst{
		llm:      chatModel,
		registry: registry,
	}
}

// SystemPrompt 系统提示词
const FundamentalsAnalystSystemPrompt = `You are a fundamentals analyst tasked with analyzing fundamental information about a company.

Your responsibilities:
1. Analyze company financial documents and statements
2. Evaluate valuation metrics (P/E, P/B, P/S, EV/EBITDA)
3. Assess profitability, growth, and financial health
4. Compare metrics to industry peers

Please write a comprehensive fundamental analysis report including:
- Company overview and business model
- Financial statement analysis (Balance Sheet, Income Statement, Cash Flow)
- Key valuation metrics and comparison
- Growth metrics and trends
- Financial health indicators
- Investment thesis based on fundamentals

Make sure to append a Markdown table at the end of the report to organize key financial metrics.

Use the available tools:
- get_fundamentals: Get company overview and key metrics
- get_balance_sheet: Get balance sheet data
- get_cashflow: Get cash flow statement
- get_income_statement: Get income statement`

// Run 运行分析师
func (a *FundamentalsAnalyst) Run(ctx context.Context, state *states.AgentState) (*states.AgentState, error) {
	ticker := state.CompanyOfInterest
	tradeDate := state.TradeDate

	// 构建消息
	messages := []*schema.Message{
		{
			Role:    schema.System,
			Content: FundamentalsAnalystSystemPrompt,
		},
		{
			Role:    schema.User,
			Content: fmt.Sprintf("Analyze the fundamentals for %s as of %s. Use the tools to gather financial data and provide a comprehensive fundamental analysis report.", ticker, tradeDate),
		},
	}

	// 获取工具信息
	toolInfos := a.registry.GetToolInfos()

	// 调用 LLM
	boundLLM := a.llm.BindTools(toolInfos)
	response, err := boundLLM.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("fundamentals analyst LLM call failed: %w", err)
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
			return nil, fmt.Errorf("fundamentals analyst follow-up LLM call failed: %w", err)
		}
	}

	// 更新状态
	state.FundamentalsReport = response.Content
	state.AddMessage(response)

	return state, nil
}

// GetName 获取名称
func (a *FundamentalsAnalyst) GetName() string {
	return "Fundamentals Analyst"
}
