// Package managers 风险管理经理
package managers

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
	"github.com/tradingagents/agent/internal/llm"
	"github.com/tradingagents/agent/internal/memory"
	"github.com/tradingagents/agent/internal/states"
)

// RiskManager 风险管理经理
type RiskManager struct {
	llm    *llm.ChatModelWrapper
	memory *memory.FinancialSituationMemory
}

// NewRiskManager 创建风险管理经理
func NewRiskManager(chatModel *llm.ChatModelWrapper, mem *memory.FinancialSituationMemory) *RiskManager {
	return &RiskManager{
		llm:    chatModel,
		memory: mem,
	}
}

// Run 运行风险管理经理
func (m *RiskManager) Run(ctx context.Context, state *states.AgentState) (*states.AgentState, error) {
	riskState := state.RiskDebateState
	history := riskState.History
	traderPlan := state.InvestmentPlan

	// 获取分析报告
	marketReport := state.MarketReport
	sentimentReport := state.SentimentReport
	newsReport := state.NewsReport
	fundamentalsReport := state.FundamentalsReport

	// 组合当前情境
	currSituation := fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s", marketReport, sentimentReport, newsReport, fundamentalsReport)

	// 从记忆中获取相似情境的经验
	pastMemoryStr := ""
	if m.memory != nil {
		pastMemories, err := m.memory.GetMemories(ctx, currSituation, 2)
		if err == nil && len(pastMemories) > 0 {
			for _, rec := range pastMemories {
				pastMemoryStr += rec.Recommendation + "\n\n"
			}
		}
	}

	// 构建提示词
	prompt := fmt.Sprintf(`As the Risk Management Judge and Debate Facilitator, your goal is to evaluate the debate between three risk analysts—Risky, Neutral, and Safe/Conservative—and determine the best course of action for the trader. Your decision must result in a clear recommendation: Buy, Sell, or Hold. Choose Hold only if strongly justified by specific arguments, not as a fallback when all sides seem valid. Strive for clarity and decisiveness.

Guidelines for Decision-Making:
1. **Summarize Key Arguments**: Extract the strongest points from each analyst, focusing on relevance to the context.
2. **Provide Rationale**: Support your recommendation with direct quotes and counterarguments from the debate.
3. **Refine the Trader's Plan**: Start with the trader's original plan, **%s**, and adjust it based on the analysts' insights.
4. **Learn from Past Mistakes**: Use lessons from **%s** to address prior misjudgments and improve the decision you are making now to make sure you don't make a wrong BUY/SELL/HOLD call that loses money.

Deliverables:
- A clear and actionable recommendation: Buy, Sell, or Hold.
- Detailed reasoning anchored in the debate and past reflections.
- End your response with: FINAL TRANSACTION PROPOSAL: **BUY** or **SELL** or **HOLD**

---

**Analysts Debate History:**  
%s

---

Focus on actionable insights and continuous improvement. Build on past lessons, critically evaluate all perspectives, and ensure each decision advances better outcomes.`, traderPlan, pastMemoryStr, history)

	// 调用 LLM
	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: prompt,
		},
	}

	response, err := m.llm.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("risk manager LLM call failed: %w", err)
	}

	// 更新风险辩论状态
	state.RiskDebateState = states.RiskDebateState{
		JudgeDecision:          response.Content,
		History:                riskState.History,
		RiskyHistory:           riskState.RiskyHistory,
		SafeHistory:            riskState.SafeHistory,
		NeutralHistory:         riskState.NeutralHistory,
		LatestSpeaker:          "Judge",
		CurrentRiskyResponse:   riskState.CurrentRiskyResponse,
		CurrentSafeResponse:    riskState.CurrentSafeResponse,
		CurrentNeutralResponse: riskState.CurrentNeutralResponse,
		Count:                  riskState.Count,
	}

	// 设置最终交易决策
	state.FinalTradeDecision = response.Content

	return state, nil
}

// GetName 获取名称
func (m *RiskManager) GetName() string {
	return "Risk Manager"
}
