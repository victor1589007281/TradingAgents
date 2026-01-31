// Package risk_mgmt 中立派风险分析师
package risk_mgmt

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
	"github.com/tradingagents/agent/internal/llm"
	"github.com/tradingagents/agent/internal/states"
)

// NeutralDebator 中立派风险分析师
type NeutralDebator struct {
	llm *llm.ChatModelWrapper
}

// NewNeutralDebator 创建中立派
func NewNeutralDebator(chatModel *llm.ChatModelWrapper) *NeutralDebator {
	return &NeutralDebator{
		llm: chatModel,
	}
}

// Run 运行中立派分析师
func (d *NeutralDebator) Run(ctx context.Context, state *states.AgentState) (*states.AgentState, error) {
	riskState := state.RiskDebateState
	history := riskState.History
	neutralHistory := riskState.NeutralHistory
	currentRiskyResponse := riskState.CurrentRiskyResponse
	currentSafeResponse := riskState.CurrentSafeResponse

	// 获取分析报告和交易员决策
	marketReport := state.MarketReport
	sentimentReport := state.SentimentReport
	newsReport := state.NewsReport
	fundamentalsReport := state.FundamentalsReport
	traderDecision := state.TraderInvestmentPlan

	// 构建提示词
	prompt := fmt.Sprintf(`As the Neutral Risk Analyst, your role is to provide a balanced, objective evaluation of the trader's decision. Consider both the potential rewards and the risks involved. Your analysis should:

1. Acknowledge valid points from both the aggressive and conservative analysts
2. Identify areas where their arguments may be biased or incomplete
3. Provide a balanced perspective that weighs opportunity against risk
4. Suggest modifications or conditions that could improve the risk-reward profile

Here is the trader's decision:
%s

Use the following data to support your balanced analysis:

Market Research Report: %s
Social Media Sentiment Report: %s
Latest World Affairs Report: %s
Company Fundamentals Report: %s
Current conversation history: %s
Risky analyst's arguments: %s
Conservative analyst's arguments: %s

If there are no responses from the other viewpoints, present your independent balanced assessment.

Your goal is to find the middle ground that maximizes potential while managing downside risk. Present your analysis conversationally without special formatting.`,
		traderDecision, marketReport, sentimentReport, newsReport, fundamentalsReport, history, currentRiskyResponse, currentSafeResponse)

	// 调用 LLM
	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: prompt,
		},
	}

	response, err := d.llm.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("neutral debator LLM call failed: %w", err)
	}

	argument := fmt.Sprintf("Neutral Analyst: %s", response.Content)

	// 更新风险辩论状态
	state.RiskDebateState = states.RiskDebateState{
		History:                history + "\n" + argument,
		RiskyHistory:           riskState.RiskyHistory,
		SafeHistory:            riskState.SafeHistory,
		NeutralHistory:         neutralHistory + "\n" + argument,
		LatestSpeaker:          "Neutral",
		CurrentRiskyResponse:   riskState.CurrentRiskyResponse,
		CurrentSafeResponse:    riskState.CurrentSafeResponse,
		CurrentNeutralResponse: argument,
		Count:                  riskState.Count + 1,
		JudgeDecision:          riskState.JudgeDecision,
	}

	return state, nil
}

// GetName 获取名称
func (d *NeutralDebator) GetName() string {
	return "Neutral Analyst"
}
