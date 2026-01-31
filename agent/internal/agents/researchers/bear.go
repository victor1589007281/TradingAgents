// Package researchers 看跌研究员
package researchers

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
	"github.com/tradingagents/agent/internal/llm"
	"github.com/tradingagents/agent/internal/memory"
	"github.com/tradingagents/agent/internal/states"
)

// BearResearcher 看跌研究员
type BearResearcher struct {
	llm    *llm.ChatModelWrapper
	memory *memory.FinancialSituationMemory
}

// NewBearResearcher 创建看跌研究员
func NewBearResearcher(chatModel *llm.ChatModelWrapper, mem *memory.FinancialSituationMemory) *BearResearcher {
	return &BearResearcher{
		llm:    chatModel,
		memory: mem,
	}
}

// Run 运行研究员
func (r *BearResearcher) Run(ctx context.Context, state *states.AgentState) (*states.AgentState, error) {
	debateState := state.InvestmentDebateState
	history := debateState.History
	bearHistory := debateState.BearHistory
	currentResponse := debateState.CurrentResponse

	// 获取所有分析报告
	marketReport := state.MarketReport
	sentimentReport := state.SentimentReport
	newsReport := state.NewsReport
	fundamentalsReport := state.FundamentalsReport

	// 组合当前情境
	currSituation := fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s", marketReport, sentimentReport, newsReport, fundamentalsReport)

	// 从记忆中获取相似情境的经验
	pastMemoryStr := ""
	if r.memory != nil {
		pastMemories, err := r.memory.GetMemories(ctx, currSituation, 2)
		if err == nil && len(pastMemories) > 0 {
			for _, rec := range pastMemories {
				pastMemoryStr += rec.Recommendation + "\n\n"
			}
		}
	}

	// 构建提示词
	prompt := fmt.Sprintf(`You are a Bear Analyst tasked with critically evaluating the stock investment decision. Your role is to identify potential risks, weaknesses, and warning signs that could negatively impact the investment.

Key points to focus on:
- Risk Identification: Highlight potential market, operational, and financial risks.
- Weakness Analysis: Point out competitive disadvantages, market threats, or operational inefficiencies.
- Negative Indicators: Cite concerning trends in financial data, market sentiment, or industry dynamics.
- Bull Counterpoints: Critically examine the bull's arguments, identify flaws in reasoning or overly optimistic assumptions.
- Engagement: Present your argument conversationally, directly challenging the bull's points with evidence-based rebuttals.

Resources available:
Market research report: %s
Social media sentiment report: %s
Latest world affairs news: %s
Company fundamentals report: %s
Conversation history of the debate: %s
Last bull argument: %s
Reflections from similar situations and lessons learned: %s

Deliver a compelling bearish case that challenges optimistic assumptions and presents a balanced view of the investment risks. Learn from past mistakes and apply those lessons.`,
		marketReport, sentimentReport, newsReport, fundamentalsReport, history, currentResponse, pastMemoryStr)

	// 调用 LLM
	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: prompt,
		},
	}

	response, err := r.llm.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("bear researcher LLM call failed: %w", err)
	}

	argument := fmt.Sprintf("Bear Analyst: %s", response.Content)

	// 更新辩论状态
	state.InvestmentDebateState = states.InvestDebateState{
		History:         history + "\n" + argument,
		BullHistory:     debateState.BullHistory,
		BearHistory:     bearHistory + "\n" + argument,
		CurrentResponse: argument,
		Count:           debateState.Count + 1,
		JudgeDecision:   debateState.JudgeDecision,
	}

	return state, nil
}

// GetName 获取名称
func (r *BearResearcher) GetName() string {
	return "Bear Researcher"
}
