// Package trader 交易员 Agent
package trader

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
	"github.com/tradingagents/agent/internal/llm"
	"github.com/tradingagents/agent/internal/memory"
	"github.com/tradingagents/agent/internal/states"
)

// Trader 交易员
type Trader struct {
	llm    *llm.ChatModelWrapper
	memory *memory.FinancialSituationMemory
}

// NewTrader 创建交易员
func NewTrader(chatModel *llm.ChatModelWrapper, mem *memory.FinancialSituationMemory) *Trader {
	return &Trader{
		llm:    chatModel,
		memory: mem,
	}
}

// Run 运行交易员
func (t *Trader) Run(ctx context.Context, state *states.AgentState) (*states.AgentState, error) {
	companyName := state.CompanyOfInterest
	investmentPlan := state.InvestmentPlan

	// 获取所有分析报告
	marketReport := state.MarketReport
	sentimentReport := state.SentimentReport
	newsReport := state.NewsReport
	fundamentalsReport := state.FundamentalsReport

	// 组合当前情境
	currSituation := fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s", marketReport, sentimentReport, newsReport, fundamentalsReport)

	// 从记忆中获取相似情境的经验
	pastMemoryStr := "No past memories found."
	if t.memory != nil {
		pastMemories, err := t.memory.GetMemories(ctx, currSituation, 2)
		if err == nil && len(pastMemories) > 0 {
			pastMemoryStr = ""
			for _, rec := range pastMemories {
				pastMemoryStr += rec.Recommendation + "\n\n"
			}
		}
	}

	// 构建系统消息
	systemPrompt := fmt.Sprintf(`You are a trading agent analyzing market data to make investment decisions. Based on your analysis, provide a specific recommendation to buy, sell, or hold. End with a firm decision and always conclude your response with 'FINAL TRANSACTION PROPOSAL: **BUY/HOLD/SELL**' to confirm your recommendation. Do not forget to utilize lessons from past decisions to learn from your mistakes. Here is some reflections from similar situations you traded in and the lessons learned: %s`, pastMemoryStr)

	// 构建用户消息
	userPrompt := fmt.Sprintf(`Based on a comprehensive analysis by a team of analysts, here is an investment plan tailored for %s. This plan incorporates insights from current technical market trends, macroeconomic indicators, and social media sentiment. Use this plan as a foundation for evaluating your next trading decision.

Proposed Investment Plan: %s

Leverage these insights to make an informed and strategic decision.`, companyName, investmentPlan)

	// 调用 LLM
	messages := []*schema.Message{
		{
			Role:    schema.System,
			Content: systemPrompt,
		},
		{
			Role:    schema.User,
			Content: userPrompt,
		},
	}

	response, err := t.llm.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("trader LLM call failed: %w", err)
	}

	// 更新状态
	state.TraderInvestmentPlan = response.Content
	state.AddMessage(response)
	state.Sender = "Trader"

	return state, nil
}

// GetName 获取名称
func (t *Trader) GetName() string {
	return "Trader"
}
