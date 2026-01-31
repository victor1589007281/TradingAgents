// Package risk_mgmt 保守派风险分析师
package risk_mgmt

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
	"github.com/tradingagents/agent/internal/llm"
	"github.com/tradingagents/agent/internal/states"
)

// ConservativeDebator 保守派风险分析师
type ConservativeDebator struct {
	llm *llm.ChatModelWrapper
}

// NewConservativeDebator 创建保守派
func NewConservativeDebator(chatModel *llm.ChatModelWrapper) *ConservativeDebator {
	return &ConservativeDebator{
		llm: chatModel,
	}
}

// Run 运行保守派分析师
func (d *ConservativeDebator) Run(ctx context.Context, state *states.AgentState) (*states.AgentState, error) {
	riskState := state.RiskDebateState
	history := riskState.History
	safeHistory := riskState.SafeHistory
	currentRiskyResponse := riskState.CurrentRiskyResponse
	currentNeutralResponse := riskState.CurrentNeutralResponse

	// 获取分析报告和交易员决策
	marketReport := state.MarketReport
	sentimentReport := state.SentimentReport
	newsReport := state.NewsReport
	fundamentalsReport := state.FundamentalsReport
	traderDecision := state.TraderInvestmentPlan

	// 构建提示词
	prompt := fmt.Sprintf(`As the Conservative Risk Analyst, your role is to prioritize capital preservation and risk mitigation. When evaluating the trader's decision, focus on:

1. Identifying potential downside risks and worst-case scenarios
2. Highlighting market uncertainties and volatility concerns
3. Questioning aggressive assumptions in the investment thesis
4. Recommending risk mitigation strategies or position sizing adjustments

Here is the trader's decision:
%s

Counter the risky analyst's arguments by pointing out:
- Overlooked risks or blind spots
- Historical examples of similar situations that went wrong
- Current market conditions that warrant caution
- The importance of protecting capital

Use the following data to support your conservative stance:

Market Research Report: %s
Social Media Sentiment Report: %s
Latest World Affairs Report: %s
Company Fundamentals Report: %s
Current conversation history: %s
Risky analyst's arguments: %s
Neutral analyst's arguments: %s

If there are no responses from the other viewpoints, present your independent risk assessment.

Your goal is to ensure that potential downsides are fully considered before any investment decision. Present your analysis conversationally without special formatting.`,
		traderDecision, marketReport, sentimentReport, newsReport, fundamentalsReport, history, currentRiskyResponse, currentNeutralResponse)

	// 调用 LLM
	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: prompt,
		},
	}

	response, err := d.llm.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("conservative debator LLM call failed: %w", err)
	}

	argument := fmt.Sprintf("Safe Analyst: %s", response.Content)

	// 更新风险辩论状态
	state.RiskDebateState = states.RiskDebateState{
		History:                history + "\n" + argument,
		RiskyHistory:           riskState.RiskyHistory,
		SafeHistory:            safeHistory + "\n" + argument,
		NeutralHistory:         riskState.NeutralHistory,
		LatestSpeaker:          "Safe",
		CurrentRiskyResponse:   riskState.CurrentRiskyResponse,
		CurrentSafeResponse:    argument,
		CurrentNeutralResponse: riskState.CurrentNeutralResponse,
		Count:                  riskState.Count + 1,
		JudgeDecision:          riskState.JudgeDecision,
	}

	return state, nil
}

// GetName 获取名称
func (d *ConservativeDebator) GetName() string {
	return "Safe Analyst"
}
