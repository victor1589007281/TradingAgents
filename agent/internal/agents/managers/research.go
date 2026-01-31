// Package managers 管理者 Agents
package managers

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
	"github.com/tradingagents/agent/internal/llm"
	"github.com/tradingagents/agent/internal/memory"
	"github.com/tradingagents/agent/internal/states"
)

// ResearchManager 研究经理
type ResearchManager struct {
	llm    *llm.ChatModelWrapper
	memory *memory.FinancialSituationMemory
}

// NewResearchManager 创建研究经理
func NewResearchManager(chatModel *llm.ChatModelWrapper, mem *memory.FinancialSituationMemory) *ResearchManager {
	return &ResearchManager{
		llm:    chatModel,
		memory: mem,
	}
}

// Run 运行研究经理
func (m *ResearchManager) Run(ctx context.Context, state *states.AgentState) (*states.AgentState, error) {
	debateState := state.InvestmentDebateState
	history := debateState.History

	// 获取所有分析报告
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
	prompt := fmt.Sprintf(`As the portfolio manager and debate facilitator, your role is to critically evaluate this round of debate and make a definitive decision: align with the bear analyst, the bull analyst, or choose Hold only if it is strongly justified based on the arguments presented.

Summarize the key points from both sides concisely, focusing on the most compelling evidence or reasoning. Your recommendation—Buy, Sell, or Hold—must be clear and actionable. Avoid defaulting to Hold simply because both sides have valid points; commit to a stance grounded in the debate's strongest arguments.

Additionally, develop a detailed investment plan for the trader. This should include:

Your Recommendation: A decisive stance supported by the most convincing arguments.
Rationale: An explanation of why these arguments lead to your conclusion.
Strategic Actions: Concrete steps for implementing the recommendation.

Take into account your past mistakes on similar situations. Use these insights to refine your decision-making and ensure you are learning and improving. Present your analysis conversationally, as if speaking naturally, without special formatting.

Here are your past reflections on mistakes:
"%s"

Here is the debate:
Debate History:
%s`, pastMemoryStr, history)

	// 调用 LLM
	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: prompt,
		},
	}

	response, err := m.llm.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("research manager LLM call failed: %w", err)
	}

	// 更新状态
	state.InvestmentDebateState = states.InvestDebateState{
		JudgeDecision:   response.Content,
		History:         debateState.History,
		BearHistory:     debateState.BearHistory,
		BullHistory:     debateState.BullHistory,
		CurrentResponse: response.Content,
		Count:           debateState.Count,
	}
	state.InvestmentPlan = response.Content

	return state, nil
}

// GetName 获取名称
func (m *ResearchManager) GetName() string {
	return "Research Manager"
}
