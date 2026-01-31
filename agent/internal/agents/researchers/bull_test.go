// Package researchers Bull 研究员测试
package researchers

import (
	"context"
	"testing"

	"github.com/tradingagents/agent/internal/config"
	"github.com/tradingagents/agent/internal/memory"
	"github.com/tradingagents/agent/internal/states"
)

func TestNewBullResearcher(t *testing.T) {
	cfg := config.DefaultConfig()
	mem := memory.NewFinancialSituationMemory("test_bull", cfg)
	researcher := NewBullResearcher(cfg, mem)

	if researcher == nil {
		t.Fatal("BullResearcher should not be nil")
	}

	if researcher.Name() != "bull_researcher" {
		t.Errorf("Expected name 'bull_researcher', got '%s'", researcher.Name())
	}
}

func TestBullResearcherDescription(t *testing.T) {
	cfg := config.DefaultConfig()
	mem := memory.NewFinancialSituationMemory("test_bull", cfg)
	researcher := NewBullResearcher(cfg, mem)

	desc := researcher.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}

func TestBullResearcherSystemPrompt(t *testing.T) {
	cfg := config.DefaultConfig()
	mem := memory.NewFinancialSituationMemory("test_bull", cfg)
	researcher := NewBullResearcher(cfg, mem)

	prompt := researcher.GetSystemPrompt()
	if prompt == "" {
		t.Error("System prompt should not be empty")
	}

	// 验证看多立场
	bullishKeywords := []string{"bull", "invest", "opportunity", "growth"}
	found := false
	for _, keyword := range bullishKeywords {
		if containsIgnoreCaseResearcher(prompt, keyword) {
			found = true
			break
		}
	}

	if !found {
		t.Error("Bull researcher prompt should contain bullish keywords")
	}
}

func TestBullResearcherProcess(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.APIKey = ""
	mem := memory.NewFinancialSituationMemory("test_bull", cfg)
	researcher := NewBullResearcher(cfg, mem)

	state := states.NewAgentState("NVDA", "2024-05-10")
	state.MarketReport = "Strong market momentum"
	state.SentimentReport = "Positive sentiment"
	state.NewsReport = "Good earnings"
	state.FundamentalsReport = "Strong balance sheet"

	ctx := context.Background()

	// 测试处理
	newState, err := researcher.Process(ctx, state)
	if err != nil {
		t.Logf("Expected error without API key: %v", err)
		return
	}

	if newState == nil {
		t.Error("Returned state should not be nil")
	}
}

func TestBullResearcherMemoryIntegration(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.APIKey = ""
	mem := memory.NewFinancialSituationMemory("test_bull", cfg)

	// 添加一些记忆
	ctx := context.Background()
	mem.AddSituations(ctx, []memory.SituationAdvice{
		{
			Situation:      "Strong earnings, positive momentum",
			Recommendation: "Buy was correct, gained 20%",
		},
	})

	researcher := NewBullResearcher(cfg, mem)

	// 验证记忆被正确关联
	if researcher.GetMemory() != mem {
		t.Error("Memory should be properly linked")
	}
}

func containsIgnoreCaseResearcher(s, substr string) bool {
	s = toLowerCaseResearcher(s)
	substr = toLowerCaseResearcher(substr)
	return containsResearcher(s, substr)
}

func toLowerCaseResearcher(s string) string {
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

func containsResearcher(s, substr string) bool {
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
