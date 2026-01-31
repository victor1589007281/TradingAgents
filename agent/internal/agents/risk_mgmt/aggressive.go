// Package risk_mgmt 风险管理 Agents
package risk_mgmt

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
	"github.com/tradingagents/agent/internal/llm"
	"github.com/tradingagents/agent/internal/states"
)

// AggressiveDebator 激进派风险分析师
type AggressiveDebator struct {
	llm *llm.ChatModelWrapper
}

// NewAggressiveDebator 创建激进派
func NewAggressiveDebator(chatModel *llm.ChatModelWrapper) *AggressiveDebator {
	return &AggressiveDebator{
		llm: chatModel,
	}
}

// Run 运行激进派分析师
func (d *AggressiveDebator) Run(ctx context.Context, state *states.AgentState) (*states.AgentState, error) {
	riskState := state.RiskDebateState
	history := riskState.History
	riskyHistory := riskState.RiskyHistory
	currentSafeResponse := riskState.CurrentSafeResponse
	currentNeutralResponse := riskState.CurrentNeutralResponse

	// 获取分析报告和交易员决策
	marketReport := state.MarketReport
	sentimentReport := state.SentimentReport
	newsReport := state.NewsReport
	fundamentalsReport := state.FundamentalsReport
	traderDecision := state.TraderInvestmentPlan

	// 构建提示词
	prompt := fmt.Sprintf(`As the Risky Risk Analyst, your role is to actively champion high-reward, high-risk opportunities, emphasizing bold strategies and competitive advantages. When evaluating the trader's decision or plan, focus intently on the potential upside, growth potential, and innovative benefits—even when these come with elevated risk. Use the provided market data and sentiment analysis to strengthen your arguments and challenge the opposing views. Specifically, respond directly to each point made by the conservative and neutral analysts, countering with data-driven rebuttals and persuasive reasoning. Highlight where their caution might miss critical opportunities or where their assumptions may be overly conservative. Here is the trader's decision:

%s

Your task is to create a compelling case for the trader's decision by questioning and critiquing the conservative and neutral stances to demonstrate why your high-reward perspective offers the best path forward. Incorporate insights from the following sources into your arguments:

Market Research Report: %s
Social Media Sentiment Report: %s
Latest World Affairs Report: %s
Company Fundamentals Report: %s
Here is the current conversation history: %s
Here are the last arguments from the conservative analyst: %s
Here are the last arguments from the neutral analyst: %s

If there are no responses from the other viewpoints, do not hallucinate and just present your point.

Engage actively by addressing any specific concerns raised, refuting the weaknesses in their logic, and asserting the benefits of risk-taking to outpace market norms. Maintain a focus on debating and persuading, not just presenting data. Challenge each counterpoint to underscore why a high-risk approach is optimal. Output conversationally as if you are speaking without any special formatting.`,
		traderDecision, marketReport, sentimentReport, newsReport, fundamentalsReport, history, currentSafeResponse, currentNeutralResponse)

	// 调用 LLM
	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: prompt,
		},
	}

	response, err := d.llm.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("aggressive debator LLM call failed: %w", err)
	}

	argument := fmt.Sprintf("Risky Analyst: %s", response.Content)

	// 更新风险辩论状态
	state.RiskDebateState = states.RiskDebateState{
		History:                history + "\n" + argument,
		RiskyHistory:           riskyHistory + "\n" + argument,
		SafeHistory:            riskState.SafeHistory,
		NeutralHistory:         riskState.NeutralHistory,
		LatestSpeaker:          "Risky",
		CurrentRiskyResponse:   argument,
		CurrentSafeResponse:    riskState.CurrentSafeResponse,
		CurrentNeutralResponse: riskState.CurrentNeutralResponse,
		Count:                  riskState.Count + 1,
		JudgeDecision:          riskState.JudgeDecision,
	}

	return state, nil
}

// GetName 获取名称
func (d *AggressiveDebator) GetName() string {
	return "Risky Analyst"
}
