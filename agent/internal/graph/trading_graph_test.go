// Package graph Trading Graph 测试
package graph

import (
	"context"
	"os"
	"testing"

	"github.com/tradingagents/agent/internal/config"
	"github.com/tradingagents/agent/internal/states"
)

func TestNewTradingAgentsGraph(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.APIKey = "test-key" // 使用测试 key

	ctx := context.Background()
	graph, err := NewTradingAgentsGraph(ctx, cfg, nil)

	// 没有真实 API key 应该返回错误或跳过
	if err != nil {
		t.Logf("Expected error without valid API key: %v", err)
		return
	}

	if graph == nil {
		t.Error("Graph should not be nil when creation succeeds")
	}
}

func TestTradingAgentsGraphWithAnalysts(t *testing.T) {
	cfg := config.DefaultConfig()
	ctx := context.Background()

	// 测试选择特定分析师
	selectedAnalysts := []string{"market", "fundamentals"}
	graph, err := NewTradingAgentsGraph(ctx, cfg, selectedAnalysts)

	if err != nil {
		t.Logf("Expected error without valid API key: %v", err)
		return
	}

	if graph == nil {
		t.Error("Graph should not be nil")
	}

	// 验证配置
	if graph.GetConfig() != cfg {
		t.Error("Config should match")
	}
}

func TestTradingAgentsGraphGetConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.MaxDebateRounds = 5
	ctx := context.Background()

	graph, _ := NewTradingAgentsGraph(ctx, cfg, nil)
	if graph != nil {
		gotCfg := graph.GetConfig()
		if gotCfg.MaxDebateRounds != 5 {
			t.Errorf("Expected MaxDebateRounds 5, got %d", gotCfg.MaxDebateRounds)
		}
	}
}

func TestTradingAgentsGraphPropagate(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping test: OPENAI_API_KEY not set")
	}

	cfg := config.DefaultConfig()
	ctx := context.Background()

	graph, err := NewTradingAgentsGraph(ctx, cfg, nil)
	if err != nil {
		t.Fatalf("Failed to create graph: %v", err)
	}

	state, decision, err := graph.Propagate(ctx, "NVDA", "2024-05-10")
	if err != nil {
		t.Fatalf("Propagate failed: %v", err)
	}

	if state == nil {
		t.Error("State should not be nil")
	}

	if decision == "" {
		t.Error("Decision should not be empty")
	}

	validDecisions := map[string]bool{"BUY": true, "SELL": true, "HOLD": true}
	if !validDecisions[decision] {
		t.Errorf("Invalid decision: %s", decision)
	}
}

func TestTradingAgentsGraphReflect(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping test: OPENAI_API_KEY not set")
	}

	cfg := config.DefaultConfig()
	ctx := context.Background()

	graph, err := NewTradingAgentsGraph(ctx, cfg, nil)
	if err != nil {
		t.Fatalf("Failed to create graph: %v", err)
	}

	state := states.NewAgentState("NVDA", "2024-05-10")
	state.MarketReport = "Test market report"
	state.FinalTradeDecision = "BUY"

	err = graph.ReflectAndRemember(ctx, state, 10.5)
	if err != nil {
		t.Errorf("ReflectAndRemember failed: %v", err)
	}
}

func TestGraphOptions(t *testing.T) {
	tests := []struct {
		name      string
		analysts  []string
		wantError bool
	}{
		{
			name:      "default analysts",
			analysts:  nil,
			wantError: false,
		},
		{
			name:      "single analyst",
			analysts:  []string{"market"},
			wantError: false,
		},
		{
			name:      "multiple analysts",
			analysts:  []string{"market", "news", "fundamentals"},
			wantError: false,
		},
		{
			name:      "all analysts",
			analysts:  []string{"market", "social", "news", "fundamentals"},
			wantError: false,
		},
	}

	cfg := config.DefaultConfig()
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTradingAgentsGraph(ctx, cfg, tt.analysts)
			if (err != nil) != tt.wantError {
				t.Errorf("NewTradingAgentsGraph() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// Benchmark 测试
func BenchmarkGraphCreation(b *testing.B) {
	cfg := config.DefaultConfig()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewTradingAgentsGraph(ctx, cfg, nil)
	}
}
