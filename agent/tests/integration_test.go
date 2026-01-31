// Package tests 集成测试
package tests

import (
	"context"
	"os"
	"testing"

	"github.com/tradingagents/agent/internal/config"
	"github.com/tradingagents/agent/internal/graph"
	"github.com/tradingagents/agent/internal/memory"
	"github.com/tradingagents/agent/internal/states"
)

// TestIntegration_StateFlow 测试状态流转
func TestIntegration_StateFlow(t *testing.T) {
	// 创建初始状态
	state := states.NewAgentState("NVDA", "2024-05-10")

	// 验证初始状态
	if state.CompanyOfInterest != "NVDA" {
		t.Errorf("Expected company 'NVDA', got '%s'", state.CompanyOfInterest)
	}

	// 模拟分析报告
	state.MarketReport = "Market shows bullish trend with strong momentum"
	state.SentimentReport = "Social media sentiment is positive"
	state.NewsReport = "Recent earnings beat expectations"
	state.FundamentalsReport = "Strong balance sheet and growth metrics"

	// 验证报告已设置
	if state.MarketReport == "" {
		t.Error("MarketReport should not be empty")
	}

	// 模拟辩论状态
	state.InvestmentDebateState = states.InvestDebateState{
		BullHistory:     "Bull: Strong growth potential...",
		BearHistory:     "Bear: Valuation concerns...",
		History:         "Full debate history...",
		CurrentResponse: "Bull: Rebuttal...",
		Count:           2,
	}

	if state.InvestmentDebateState.Count != 2 {
		t.Errorf("Expected debate count 2, got %d", state.InvestmentDebateState.Count)
	}

	// 模拟交易决策
	state.InvestmentPlan = "Buy with 5% position size"
	state.TraderInvestmentPlan = "Confirmed buy recommendation"

	// 模拟风险辩论
	state.RiskDebateState = states.RiskDebateState{
		RiskyHistory:   "High reward potential...",
		SafeHistory:    "Consider position sizing...",
		NeutralHistory: "Balanced approach recommended...",
		Count:          3,
	}

	// 最终决策
	state.FinalTradeDecision = "FINAL TRANSACTION PROPOSAL: **BUY**"

	// 验证完整流程
	if state.FinalTradeDecision == "" {
		t.Error("FinalTradeDecision should not be empty")
	}
}

// TestIntegration_MemorySystem 测试记忆系统集成
func TestIntegration_MemorySystem(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.APIKey = "" // 不使用实际 API

	ctx := context.Background()

	// 创建多个记忆实例
	bullMemory := memory.NewFinancialSituationMemory("bull_memory", cfg)
	bearMemory := memory.NewFinancialSituationMemory("bear_memory", cfg)
	traderMemory := memory.NewFinancialSituationMemory("trader_memory", cfg)

	// 添加经验
	bullMemory.AddSituations(ctx, []memory.SituationAdvice{
		{
			Situation:      "Strong earnings, positive sentiment, upward trend",
			Recommendation: "Buy recommendation was correct, gained 15%",
		},
	})

	bearMemory.AddSituations(ctx, []memory.SituationAdvice{
		{
			Situation:      "Overvalued metrics, negative news cycle",
			Recommendation: "Sell recommendation was correct, avoided 10% loss",
		},
	})

	traderMemory.AddSituations(ctx, []memory.SituationAdvice{
		{
			Situation:      "Mixed signals, high volatility",
			Recommendation: "Hold recommendation was appropriate",
		},
	})

	// 验证记忆数量
	if bullMemory.Count() != 1 {
		t.Errorf("Expected bull memory count 1, got %d", bullMemory.Count())
	}

	if bearMemory.Count() != 1 {
		t.Errorf("Expected bear memory count 1, got %d", bearMemory.Count())
	}

	if traderMemory.Count() != 1 {
		t.Errorf("Expected trader memory count 1, got %d", traderMemory.Count())
	}
}

// TestIntegration_ConfigLoading 测试配置加载集成
func TestIntegration_ConfigLoading(t *testing.T) {
	// 设置环境变量
	os.Setenv("OPENAI_API_KEY", "test-key")
	os.Setenv("ALPHA_VANTAGE_API_KEY", "test-av-key")
	defer func() {
		os.Unsetenv("OPENAI_API_KEY")
		os.Unsetenv("ALPHA_VANTAGE_API_KEY")
	}()

	cfg := config.DefaultConfig()

	// 验证环境变量被正确读取
	if cfg.APIKey != "test-key" {
		t.Errorf("Expected API key 'test-key', got '%s'", cfg.APIKey)
	}

	if cfg.AlphaVantageAPIKey != "test-av-key" {
		t.Errorf("Expected Alpha Vantage key 'test-av-key', got '%s'", cfg.AlphaVantageAPIKey)
	}

	// 验证默认配置
	if cfg.MaxDebateRounds != 1 {
		t.Errorf("Expected MaxDebateRounds 1, got %d", cfg.MaxDebateRounds)
	}
}

// TestIntegration_GraphCreation 测试图创建（需要 API key）
func TestIntegration_GraphCreation(t *testing.T) {
	// 检查是否有 API key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" || apiKey == "test-key" {
		t.Skip("Skipping integration test: OPENAI_API_KEY not set")
	}

	ctx := context.Background()
	cfg := config.DefaultConfig()

	// 尝试创建图
	tradingGraph, err := graph.NewTradingAgentsGraph(ctx, cfg, []string{"market"})
	if err != nil {
		t.Fatalf("Failed to create trading graph: %v", err)
	}

	if tradingGraph == nil {
		t.Error("Trading graph should not be nil")
	}

	// 验证配置
	if tradingGraph.GetConfig() != cfg {
		t.Error("Config should match")
	}
}

// TestIntegration_EndToEnd 端到端测试（需要 API key）
func TestIntegration_EndToEnd(t *testing.T) {
	// 检查是否有 API key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" || apiKey == "test-key" {
		t.Skip("Skipping end-to-end test: OPENAI_API_KEY not set")
	}

	ctx := context.Background()
	cfg := config.DefaultConfig()
	cfg.MaxDebateRounds = 1
	cfg.MaxRiskDiscussRounds = 1

	// 创建图
	tradingGraph, err := graph.NewTradingAgentsGraph(ctx, cfg, nil)
	if err != nil {
		t.Fatalf("Failed to create trading graph: %v", err)
	}

	// 执行分析
	state, decision, err := tradingGraph.Propagate(ctx, "NVDA", "2024-05-10")
	if err != nil {
		t.Fatalf("Propagate failed: %v", err)
	}

	// 验证结果
	if state == nil {
		t.Error("State should not be nil")
	}

	if decision == "" {
		t.Error("Decision should not be empty")
	}

	// 决策应该是 BUY, SELL, 或 HOLD
	validDecisions := map[string]bool{"BUY": true, "SELL": true, "HOLD": true}
	if !validDecisions[decision] {
		t.Errorf("Invalid decision: %s", decision)
	}

	// 验证报告已生成
	if state.MarketReport == "" {
		t.Error("MarketReport should not be empty")
	}

	if state.FinalTradeDecision == "" {
		t.Error("FinalTradeDecision should not be empty")
	}
}

// TestIntegration_ReflectAndRemember 测试反思学习
func TestIntegration_ReflectAndRemember(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" || apiKey == "test-key" {
		t.Skip("Skipping reflection test: OPENAI_API_KEY not set")
	}

	ctx := context.Background()
	cfg := config.DefaultConfig()

	tradingGraph, err := graph.NewTradingAgentsGraph(ctx, cfg, nil)
	if err != nil {
		t.Fatalf("Failed to create trading graph: %v", err)
	}

	// 创建模拟状态
	state := states.NewAgentState("NVDA", "2024-05-10")
	state.MarketReport = "Market analysis report"
	state.SentimentReport = "Sentiment report"
	state.NewsReport = "News report"
	state.FundamentalsReport = "Fundamentals report"

	// 测试反思学习
	err = tradingGraph.ReflectAndRemember(ctx, state, 10.5)
	if err != nil {
		t.Errorf("ReflectAndRemember failed: %v", err)
	}
}

// BenchmarkIntegration_StateCreation 状态创建性能测试
func BenchmarkIntegration_StateCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		states.NewAgentState("NVDA", "2024-05-10")
	}
}

// BenchmarkIntegration_MemoryOperations 记忆操作性能测试
func BenchmarkIntegration_MemoryOperations(b *testing.B) {
	cfg := config.DefaultConfig()
	cfg.APIKey = ""
	mem := memory.NewFinancialSituationMemory("bench_memory", cfg)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mem.AddSituations(ctx, []memory.SituationAdvice{
			{Situation: "Test", Recommendation: "Test rec"},
		})
	}
}
