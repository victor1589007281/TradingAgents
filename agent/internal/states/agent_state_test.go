// Package states Agent 状态测试
package states

import (
	"testing"

	"github.com/cloudwego/eino/schema"
)

func TestNewAgentState(t *testing.T) {
	state := NewAgentState("NVDA", "2024-05-10")

	if state.CompanyOfInterest != "NVDA" {
		t.Errorf("Expected CompanyOfInterest to be 'NVDA', got '%s'", state.CompanyOfInterest)
	}

	if state.TradeDate != "2024-05-10" {
		t.Errorf("Expected TradeDate to be '2024-05-10', got '%s'", state.TradeDate)
	}

	if len(state.Messages) != 0 {
		t.Errorf("Expected empty Messages, got %d", len(state.Messages))
	}

	if state.InvestmentDebateState.Count != 0 {
		t.Errorf("Expected InvestmentDebateState.Count to be 0, got %d", state.InvestmentDebateState.Count)
	}

	if state.RiskDebateState.Count != 0 {
		t.Errorf("Expected RiskDebateState.Count to be 0, got %d", state.RiskDebateState.Count)
	}
}

func TestAgentStateClone(t *testing.T) {
	state := NewAgentState("AAPL", "2024-01-15")
	state.MarketReport = "Test market report"
	state.AddMessage(&schema.Message{
		Role:    schema.User,
		Content: "Test message",
	})

	cloned := state.Clone()

	// 验证克隆的值相同
	if cloned.CompanyOfInterest != state.CompanyOfInterest {
		t.Error("Clone should have same CompanyOfInterest")
	}

	if cloned.MarketReport != state.MarketReport {
		t.Error("Clone should have same MarketReport")
	}

	// 验证是独立的副本
	cloned.MarketReport = "Modified report"
	if state.MarketReport == cloned.MarketReport {
		t.Error("Clone should be independent copy")
	}
}

func TestAgentStateAddMessage(t *testing.T) {
	state := NewAgentState("GOOGL", "2024-02-20")

	msg1 := &schema.Message{Role: schema.User, Content: "Message 1"}
	msg2 := &schema.Message{Role: schema.Assistant, Content: "Message 2"}

	state.AddMessage(msg1)
	if len(state.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(state.Messages))
	}

	state.AddMessage(msg2)
	if len(state.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(state.Messages))
	}
}

func TestAgentStateClearMessages(t *testing.T) {
	state := NewAgentState("MSFT", "2024-03-10")
	state.AddMessage(&schema.Message{Role: schema.User, Content: "Message 1"})
	state.AddMessage(&schema.Message{Role: schema.Assistant, Content: "Message 2"})

	if len(state.Messages) != 2 {
		t.Errorf("Expected 2 messages before clear, got %d", len(state.Messages))
	}

	state.ClearMessages()

	if len(state.Messages) != 0 {
		t.Errorf("Expected 0 messages after clear, got %d", len(state.Messages))
	}
}

func TestAgentStateGetAllReports(t *testing.T) {
	state := NewAgentState("TSLA", "2024-04-05")
	state.MarketReport = "Market Report"
	state.SentimentReport = "Sentiment Report"
	state.NewsReport = "News Report"
	state.FundamentalsReport = "Fundamentals Report"

	reports := state.GetAllReports()

	if reports == "" {
		t.Error("GetAllReports should not return empty string")
	}

	// 验证所有报告都在结果中
	expectedContents := []string{"Market Report", "Sentiment Report", "News Report", "Fundamentals Report"}
	for _, content := range expectedContents {
		if !contains(reports, content) {
			t.Errorf("GetAllReports should contain '%s'", content)
		}
	}
}

func TestInvestDebateState(t *testing.T) {
	state := InvestDebateState{
		BullHistory:     "Bull argument 1",
		BearHistory:     "Bear argument 1",
		History:         "Full history",
		CurrentResponse: "Current response",
		JudgeDecision:   "Buy",
		Count:           2,
	}

	if state.Count != 2 {
		t.Errorf("Expected Count to be 2, got %d", state.Count)
	}

	if state.JudgeDecision != "Buy" {
		t.Errorf("Expected JudgeDecision to be 'Buy', got '%s'", state.JudgeDecision)
	}
}

func TestRiskDebateState(t *testing.T) {
	state := RiskDebateState{
		RiskyHistory:           "Risky argument",
		SafeHistory:            "Safe argument",
		NeutralHistory:         "Neutral argument",
		History:                "Full history",
		LatestSpeaker:          "Risky",
		CurrentRiskyResponse:   "Current risky",
		CurrentSafeResponse:    "Current safe",
		CurrentNeutralResponse: "Current neutral",
		JudgeDecision:          "Hold",
		Count:                  3,
	}

	if state.Count != 3 {
		t.Errorf("Expected Count to be 3, got %d", state.Count)
	}

	if state.LatestSpeaker != "Risky" {
		t.Errorf("Expected LatestSpeaker to be 'Risky', got '%s'", state.LatestSpeaker)
	}
}

// 辅助函数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
