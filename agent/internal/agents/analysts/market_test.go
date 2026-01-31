// Package analysts 市场分析师测试
package analysts

import (
	"context"
	"testing"

	"github.com/tradingagents/agent/internal/config"
	"github.com/tradingagents/agent/internal/states"
)

func TestNewMarketAnalyst(t *testing.T) {
	cfg := config.DefaultConfig()
	analyst := NewMarketAnalyst(cfg)

	if analyst == nil {
		t.Fatal("MarketAnalyst should not be nil")
	}

	if analyst.Name() != "market_analyst" {
		t.Errorf("Expected name 'market_analyst', got '%s'", analyst.Name())
	}
}

func TestMarketAnalystDescription(t *testing.T) {
	cfg := config.DefaultConfig()
	analyst := NewMarketAnalyst(cfg)

	desc := analyst.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}

func TestMarketAnalystGetTools(t *testing.T) {
	cfg := config.DefaultConfig()
	analyst := NewMarketAnalyst(cfg)

	tools := analyst.GetTools()
	if len(tools) == 0 {
		t.Error("Market analyst should have tools")
	}

	// 验证必要的工具存在
	expectedTools := []string{"get_stock_data", "get_indicators"}
	toolNames := make(map[string]bool)
	for _, tool := range tools {
		toolNames[tool.Name()] = true
	}

	for _, expected := range expectedTools {
		if !toolNames[expected] {
			t.Errorf("Expected tool '%s' not found", expected)
		}
	}
}

func TestMarketAnalystProcess(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.APIKey = "" // 不使用实际 API
	analyst := NewMarketAnalyst(cfg)

	state := states.NewAgentState("NVDA", "2024-05-10")
	ctx := context.Background()

	// 测试处理（不调用 LLM）
	newState, err := analyst.Process(ctx, state)
	if err != nil {
		// 没有 API key 应该返回错误或空结果
		t.Logf("Expected error without API key: %v", err)
		return
	}

	if newState == nil {
		t.Error("Returned state should not be nil")
	}
}

func TestMarketAnalystSystemPrompt(t *testing.T) {
	cfg := config.DefaultConfig()
	analyst := NewMarketAnalyst(cfg)

	prompt := analyst.GetSystemPrompt()
	if prompt == "" {
		t.Error("System prompt should not be empty")
	}

	// 验证提示词包含关键内容
	expectedContent := []string{"market", "analyst", "stock"}
	for _, content := range expectedContent {
		if !containsIgnoreCase(prompt, content) {
			t.Errorf("System prompt should contain '%s'", content)
		}
	}
}

func containsIgnoreCase(s, substr string) bool {
	s = toLowerCase(s)
	substr = toLowerCase(substr)
	return contains(s, substr)
}

func toLowerCase(s string) string {
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

func contains(s, substr string) bool {
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
