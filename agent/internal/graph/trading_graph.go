// Package graph 交易智能体图编排
package graph

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/tradingagents/agent/internal/agents/analysts"
	"github.com/tradingagents/agent/internal/agents/managers"
	"github.com/tradingagents/agent/internal/agents/researchers"
	"github.com/tradingagents/agent/internal/agents/risk_mgmt"
	"github.com/tradingagents/agent/internal/agents/trader"
	"github.com/tradingagents/agent/internal/config"
	"github.com/tradingagents/agent/internal/llm"
	"github.com/tradingagents/agent/internal/memory"
	"github.com/tradingagents/agent/internal/states"
)

// TradingAgentsGraph 交易智能体图
type TradingAgentsGraph struct {
	cfg *config.Config

	// LLM
	deepThinkingLLM  *llm.ChatModelWrapper
	quickThinkingLLM *llm.ChatModelWrapper

	// 记忆
	bullMemory        *memory.FinancialSituationMemory
	bearMemory        *memory.FinancialSituationMemory
	traderMemory      *memory.FinancialSituationMemory
	investJudgeMemory *memory.FinancialSituationMemory
	riskManagerMemory *memory.FinancialSituationMemory

	// 分析师
	marketAnalyst       *analysts.MarketAnalyst
	socialMediaAnalyst  *analysts.SocialMediaAnalyst
	newsAnalyst         *analysts.NewsAnalyst
	fundamentalsAnalyst *analysts.FundamentalsAnalyst

	// 研究员
	bullResearcher  *researchers.BullResearcher
	bearResearcher  *researchers.BearResearcher
	researchManager *managers.ResearchManager

	// 交易员
	traderAgent *trader.Trader

	// 风险管理
	riskyAnalyst   *risk_mgmt.AggressiveDebator
	neutralAnalyst *risk_mgmt.NeutralDebator
	safeAnalyst    *risk_mgmt.ConservativeDebator
	riskManager    *managers.RiskManager

	// 选择的分析师
	selectedAnalysts []string
}

// NewTradingAgentsGraph 创建交易智能体图
func NewTradingAgentsGraph(ctx context.Context, cfg *config.Config, selectedAnalysts []string) (*TradingAgentsGraph, error) {
	if len(selectedAnalysts) == 0 {
		selectedAnalysts = []string{"market", "social", "news", "fundamentals"}
	}

	tag := &TradingAgentsGraph{
		cfg:              cfg,
		selectedAnalysts: selectedAnalysts,
	}

	// 初始化 LLM
	var err error
	tag.deepThinkingLLM, err = llm.NewChatModel(ctx, cfg, cfg.DeepThinkLLM)
	if err != nil {
		return nil, fmt.Errorf("failed to create deep thinking LLM: %w", err)
	}

	tag.quickThinkingLLM, err = llm.NewChatModel(ctx, cfg, cfg.QuickThinkLLM)
	if err != nil {
		return nil, fmt.Errorf("failed to create quick thinking LLM: %w", err)
	}

	// 初始化记忆
	tag.bullMemory = memory.NewFinancialSituationMemory("bull_memory", cfg)
	tag.bearMemory = memory.NewFinancialSituationMemory("bear_memory", cfg)
	tag.traderMemory = memory.NewFinancialSituationMemory("trader_memory", cfg)
	tag.investJudgeMemory = memory.NewFinancialSituationMemory("invest_judge_memory", cfg)
	tag.riskManagerMemory = memory.NewFinancialSituationMemory("risk_manager_memory", cfg)

	// 初始化分析师
	tag.marketAnalyst = analysts.NewMarketAnalyst(tag.quickThinkingLLM)
	tag.socialMediaAnalyst = analysts.NewSocialMediaAnalyst(tag.quickThinkingLLM)
	tag.newsAnalyst = analysts.NewNewsAnalyst(tag.quickThinkingLLM)
	tag.fundamentalsAnalyst = analysts.NewFundamentalsAnalyst(tag.quickThinkingLLM)

	// 初始化研究员
	tag.bullResearcher = researchers.NewBullResearcher(tag.quickThinkingLLM, tag.bullMemory)
	tag.bearResearcher = researchers.NewBearResearcher(tag.quickThinkingLLM, tag.bearMemory)
	tag.researchManager = managers.NewResearchManager(tag.deepThinkingLLM, tag.investJudgeMemory)

	// 初始化交易员
	tag.traderAgent = trader.NewTrader(tag.quickThinkingLLM, tag.traderMemory)

	// 初始化风险管理
	tag.riskyAnalyst = risk_mgmt.NewAggressiveDebator(tag.quickThinkingLLM)
	tag.neutralAnalyst = risk_mgmt.NewNeutralDebator(tag.quickThinkingLLM)
	tag.safeAnalyst = risk_mgmt.NewConservativeDebator(tag.quickThinkingLLM)
	tag.riskManager = managers.NewRiskManager(tag.deepThinkingLLM, tag.riskManagerMemory)

	return tag, nil
}

// Propagate 执行交易决策
func (g *TradingAgentsGraph) Propagate(ctx context.Context, company, tradeDate string) (*states.AgentState, string, error) {
	// 初始化状态
	state := states.NewAgentState(company, tradeDate)

	fmt.Printf("🚀 Starting analysis for %s on %s\n", company, tradeDate)

	// 阶段 1: 分析师团队
	fmt.Println("\n📊 Phase 1: Analyst Team")

	for _, analystType := range g.selectedAnalysts {
		var err error
		switch analystType {
		case "market":
			fmt.Println("  → Running Market Analyst...")
			state, err = g.marketAnalyst.Run(ctx, state)
		case "social":
			fmt.Println("  → Running Social Media Analyst...")
			state, err = g.socialMediaAnalyst.Run(ctx, state)
		case "news":
			fmt.Println("  → Running News Analyst...")
			state, err = g.newsAnalyst.Run(ctx, state)
		case "fundamentals":
			fmt.Println("  → Running Fundamentals Analyst...")
			state, err = g.fundamentalsAnalyst.Run(ctx, state)
		}
		if err != nil {
			return nil, "", fmt.Errorf("analyst %s failed: %w", analystType, err)
		}
		state.ClearMessages() // 清理消息以保持上下文简洁
	}

	// 阶段 2: 投资辩论
	fmt.Println("\n💬 Phase 2: Investment Debate")
	maxDebateRounds := g.cfg.MaxDebateRounds

	for round := 0; round < maxDebateRounds; round++ {
		fmt.Printf("  → Debate Round %d\n", round+1)

		// 看涨研究员
		fmt.Println("    • Bull Researcher arguing...")
		state, err := g.bullResearcher.Run(ctx, state)
		if err != nil {
			return nil, "", fmt.Errorf("bull researcher failed: %w", err)
		}

		// 看跌研究员
		fmt.Println("    • Bear Researcher countering...")
		state, err = g.bearResearcher.Run(ctx, state)
		if err != nil {
			return nil, "", fmt.Errorf("bear researcher failed: %w", err)
		}
	}

	// 研究经理做出判断
	fmt.Println("  → Research Manager making decision...")
	state, err := g.researchManager.Run(ctx, state)
	if err != nil {
		return nil, "", fmt.Errorf("research manager failed: %w", err)
	}

	// 阶段 3: 交易员决策
	fmt.Println("\n💰 Phase 3: Trader Decision")
	fmt.Println("  → Trader evaluating investment plan...")
	state, err = g.traderAgent.Run(ctx, state)
	if err != nil {
		return nil, "", fmt.Errorf("trader failed: %w", err)
	}

	// 阶段 4: 风险管理辩论
	fmt.Println("\n⚠️ Phase 4: Risk Management Debate")
	maxRiskRounds := g.cfg.MaxRiskDiscussRounds

	for round := 0; round < maxRiskRounds; round++ {
		fmt.Printf("  → Risk Debate Round %d\n", round+1)

		// 激进派
		fmt.Println("    • Risky Analyst advocating...")
		state, err = g.riskyAnalyst.Run(ctx, state)
		if err != nil {
			return nil, "", fmt.Errorf("risky analyst failed: %w", err)
		}

		// 保守派
		fmt.Println("    • Safe Analyst warning...")
		state, err = g.safeAnalyst.Run(ctx, state)
		if err != nil {
			return nil, "", fmt.Errorf("safe analyst failed: %w", err)
		}

		// 中立派
		fmt.Println("    • Neutral Analyst balancing...")
		state, err = g.neutralAnalyst.Run(ctx, state)
		if err != nil {
			return nil, "", fmt.Errorf("neutral analyst failed: %w", err)
		}
	}

	// 风险管理经理做出最终决策
	fmt.Println("  → Risk Manager making final decision...")
	state, err = g.riskManager.Run(ctx, state)
	if err != nil {
		return nil, "", fmt.Errorf("risk manager failed: %w", err)
	}

	// 处理信号
	decision := g.processSignal(state.FinalTradeDecision)

	fmt.Printf("\n✅ Final Decision: %s\n", decision)

	return state, decision, nil
}

// processSignal 从完整响应中提取交易信号
func (g *TradingAgentsGraph) processSignal(fullSignal string) string {
	// 尝试匹配 "FINAL TRANSACTION PROPOSAL: **BUY**" 格式
	re := regexp.MustCompile(`FINAL TRANSACTION PROPOSAL:\s*\*{0,2}(BUY|SELL|HOLD)\*{0,2}`)
	matches := re.FindStringSubmatch(strings.ToUpper(fullSignal))
	if len(matches) > 1 {
		return matches[1]
	}

	// 尝试简单匹配
	upper := strings.ToUpper(fullSignal)
	if strings.Contains(upper, "BUY") && !strings.Contains(upper, "SELL") {
		return "BUY"
	}
	if strings.Contains(upper, "SELL") && !strings.Contains(upper, "BUY") {
		return "SELL"
	}

	return "HOLD"
}

// ReflectAndRemember 反思并记忆
func (g *TradingAgentsGraph) ReflectAndRemember(ctx context.Context, state *states.AgentState, returnsLosses float64) error {
	// 构建情境描述
	situation := state.GetAllReports()

	// 根据收益/损失生成反思
	var recommendation string
	if returnsLosses > 0 {
		recommendation = fmt.Sprintf("Positive outcome (+%.2f%%). The analysis and decision were effective. Key factors: accurate market assessment, proper risk evaluation.", returnsLosses)
	} else {
		recommendation = fmt.Sprintf("Negative outcome (%.2f%%). Review needed: market conditions may have changed unexpectedly, risk assessment may have been insufficient.", returnsLosses)
	}

	advice := memory.SituationAdvice{
		Situation:      situation,
		Recommendation: recommendation,
	}

	// 更新各个记忆
	g.bullMemory.AddSituations(ctx, []memory.SituationAdvice{advice})
	g.bearMemory.AddSituations(ctx, []memory.SituationAdvice{advice})
	g.traderMemory.AddSituations(ctx, []memory.SituationAdvice{advice})
	g.riskManagerMemory.AddSituations(ctx, []memory.SituationAdvice{advice})

	return nil
}

// GetConfig 获取配置
func (g *TradingAgentsGraph) GetConfig() *config.Config {
	return g.cfg
}
