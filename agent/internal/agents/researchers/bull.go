// Package researchers 研究员 Agents
package researchers

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
	"github.com/tradingagents/agent/internal/llm"
	"github.com/tradingagents/agent/internal/memory"
	"github.com/tradingagents/agent/internal/states"
)

// BullResearcher 看涨研究员
type BullResearcher struct {
	llm    *llm.ChatModelWrapper
	memory *memory.FinancialSituationMemory
}

// NewBullResearcher 创建看涨研究员
func NewBullResearcher(chatModel *llm.ChatModelWrapper, mem *memory.FinancialSituationMemory) *BullResearcher {
	return &BullResearcher{
		llm:    chatModel,
		memory: mem,
	}
}

// Run 运行研究员
func (r *BullResearcher) Run(ctx context.Context, state *states.AgentState) (*states.AgentState, error) {
	debateState := state.InvestmentDebateState
	history := debateState.History
	bullHistory := debateState.BullHistory
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
	prompt := fmt.Sprintf(`You are a Bull Analyst advocating for investing in the stock. Your task is to build a strong, evidence-based case emphasizing growth potential, competitive advantages, and positive market indicators. Leverage the provided research and data to address concerns and counter bearish arguments effectively.

Key points to focus on:
- Growth Potential: Highlight the company's market opportunities, revenue projections, and scalability.
- Competitive Advantages: Emphasize factors like unique products, strong branding, or dominant market positioning.
- Positive Indicators: Use financial health, industry trends, and recent positive news as evidence.
- Bear Counterpoints: Critically analyze the bear argument with specific data and sound reasoning, addressing concerns thoroughly and showing why the bull perspective holds stronger merit.
- Engagement: Present your argument in a conversational style, engaging directly with the bear analyst's points and debating effectively rather than just listing data.

Resources available:
Market research report: %s
Social media sentiment report: %s
Latest world affairs news: %s
Company fundamentals report: %s
Conversation history of the debate: %s
Last bear argument: %s
Reflections from similar situations and lessons learned: %s

Use this information to deliver a compelling bull argument, refute the bear's concerns, and engage in a dynamic debate that demonstrates the strengths of the bull position. You must also address reflections and learn from lessons and mistakes you made in the past.`,
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
		return nil, fmt.Errorf("bull researcher LLM call failed: %w", err)
	}

	argument := fmt.Sprintf("Bull Analyst: %s", response.Content)

	// 更新辩论状态
	state.InvestmentDebateState = states.InvestDebateState{
		History:         history + "\n" + argument,
		BullHistory:     bullHistory + "\n" + argument,
		BearHistory:     debateState.BearHistory,
		CurrentResponse: argument,
		Count:           debateState.Count + 1,
		JudgeDecision:   debateState.JudgeDecision,
	}

	return state, nil
}

// GetName 获取名称
func (r *BullResearcher) GetName() string {
	return "Bull Researcher"
}
