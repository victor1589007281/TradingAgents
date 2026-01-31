// Package trader Trader 测试
package trader

import (
	"context"
	"testing"

	"github.com/tradingagents/agent/internal/config"
	"github.com/tradingagents/agent/internal/memory"
	"github.com/tradingagents/agent/internal/states"
)

func TestNewTrader(t *testing.T) {
	cfg := config.DefaultConfig()
	mem := memory.NewFinancialSituationMemory("test_trader", cfg)
	trader := NewTrader(cfg, mem)

	if trader == nil {
		t.Fatal("Trader should not be nil")
	}

	if trader.Name() != "trader" {
		t.Errorf("Expected name 'trader', got '%s'", trader.Name())
	}
}

func TestTraderDescription(t *testing.T) {
	cfg := config.DefaultConfig()
	mem := memory.NewFinancialSituationMemory("test_trader", cfg)
	trader := NewTrader(cfg, mem)

	desc := trader.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}

func TestTraderSystemPrompt(t *testing.T) {
	cfg := config.DefaultConfig()
	mem := memory.NewFinancialSituationMemory("test_trader", cfg)
	trader := NewTrader(cfg, mem)

	prompt := trader.GetSystemPrompt()
	if prompt == "" {
		t.Error("System prompt should not be empty")
	}

	// 验证交易相关关键词
	tradingKeywords := []string{"trade", "transaction", "BUY", "SELL", "HOLD"}
	found := false
	for _, keyword := range tradingKeywords {
		if containsIgnoreCaseTrader(prompt, keyword) {
			found = true
			break
		}
	}

	if !found {
		t.Error("Trader prompt should contain trading keywords")
	}
}

func TestTraderProcess(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.APIKey = ""
	mem := memory.NewFinancialSituationMemory("test_trader", cfg)
	trader := NewTrader(cfg, mem)

	state := states.NewAgentState("NVDA", "2024-05-10")
	state.MarketReport = "Strong market"
	state.InvestmentPlan = "Buy with 5% position"
	state.InvestmentDebateState = states.InvestDebateState{
		JudgeDecision: "Buy",
		Count:         2,
	}

	ctx := context.Background()

	newState, err := trader.Process(ctx, state)
	if err != nil {
		t.Logf("Expected error without API key: %v", err)
		return
	}

	if newState == nil {
		t.Error("Returned state should not be nil")
	}
}

func TestTraderGetTools(t *testing.T) {
	cfg := config.DefaultConfig()
	mem := memory.NewFinancialSituationMemory("test_trader", cfg)
	trader := NewTrader(cfg, mem)

	tools := trader.GetTools()
	// Trader 可能没有工具，或者有一些辅助工具
	t.Logf("Trader has %d tools", len(tools))
}

func TestTraderMemoryIntegration(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.APIKey = ""
	mem := memory.NewFinancialSituationMemory("test_trader", cfg)

	// 添加交易记忆
	ctx := context.Background()
	mem.AddSituations(ctx, []memory.SituationAdvice{
		{
			Situation:      "Strong buy signal with positive momentum",
			Recommendation: "Executed buy, gained 15%",
		},
		{
			Situation:      "Weak fundamentals, negative sentiment",
			Recommendation: "Executed sell, avoided 10% loss",
		},
	})

	trader := NewTrader(cfg, mem)

	if trader.GetMemory() != mem {
		t.Error("Memory should be properly linked")
	}

	if mem.Count() != 2 {
		t.Errorf("Expected 2 memories, got %d", mem.Count())
	}
}

func containsIgnoreCaseTrader(s, substr string) bool {
	s = toLowerCaseTrader(s)
	substr = toLowerCaseTrader(substr)
	return containsTrader(s, substr)
}

func toLowerCaseTrader(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}

func containsTrader(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
