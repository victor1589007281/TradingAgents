// Package memory 向量记忆测试
package memory

import (
	"context"
	"testing"

	"github.com/tradingagents/agent/internal/config"
)

func TestNewFinancialSituationMemory(t *testing.T) {
	cfg := config.DefaultConfig()
	memory := NewFinancialSituationMemory("test_memory", cfg)

	if memory == nil {
		t.Fatal("Memory should not be nil")
	}

	if memory.GetName() != "test_memory" {
		t.Errorf("Expected name 'test_memory', got '%s'", memory.GetName())
	}

	if memory.Count() != 0 {
		t.Errorf("Expected count 0, got %d", memory.Count())
	}
}

func TestMemoryAddAndGetWithoutEmbedding(t *testing.T) {
	cfg := config.DefaultConfig()
	// 不设置 API key，测试没有实际 embedding 的情况
	cfg.APIKey = ""
	
	memory := NewFinancialSituationMemory("test_memory", cfg)
	ctx := context.Background()

	// 添加情境
	situations := []SituationAdvice{
		{
			Situation:      "Market is bullish with strong momentum",
			Recommendation: "Consider buying growth stocks",
		},
		{
			Situation:      "Market shows bearish signals",
			Recommendation: "Consider defensive positions",
		},
	}

	err := memory.AddSituations(ctx, situations)
	if err != nil {
		t.Errorf("Unexpected error adding situations: %v", err)
	}

	if memory.Count() != 2 {
		t.Errorf("Expected count 2, got %d", memory.Count())
	}
}

func TestMemoryClear(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.APIKey = ""
	
	memory := NewFinancialSituationMemory("test_memory", cfg)
	ctx := context.Background()

	// 添加数据
	memory.AddSituations(ctx, []SituationAdvice{
		{Situation: "Test", Recommendation: "Test rec"},
	})

	if memory.Count() != 1 {
		t.Errorf("Expected count 1, got %d", memory.Count())
	}

	// 清空
	memory.Clear()

	if memory.Count() != 0 {
		t.Errorf("Expected count 0 after clear, got %d", memory.Count())
	}
}

func TestSituationAdvice(t *testing.T) {
	advice := SituationAdvice{
		Situation:      "Test situation",
		Recommendation: "Test recommendation",
	}

	if advice.Situation != "Test situation" {
		t.Errorf("Expected 'Test situation', got '%s'", advice.Situation)
	}

	if advice.Recommendation != "Test recommendation" {
		t.Errorf("Expected 'Test recommendation', got '%s'", advice.Recommendation)
	}
}

func TestMemoryRecord(t *testing.T) {
	record := MemoryRecord{
		ID:              "1",
		Situation:       "Test situation",
		Recommendation:  "Test recommendation",
		Embedding:       []float64{0.1, 0.2, 0.3},
		SimilarityScore: 0.95,
	}

	if record.ID != "1" {
		t.Errorf("Expected ID '1', got '%s'", record.ID)
	}

	if record.SimilarityScore != 0.95 {
		t.Errorf("Expected SimilarityScore 0.95, got %f", record.SimilarityScore)
	}

	if len(record.Embedding) != 3 {
		t.Errorf("Expected Embedding length 3, got %d", len(record.Embedding))
	}
}

func TestGetMemoriesEmpty(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.APIKey = ""
	
	memory := NewFinancialSituationMemory("test_memory", cfg)
	ctx := context.Background()

	// 空记忆应返回 nil
	results, err := memory.GetMemories(ctx, "test query", 5)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if results != nil && len(results) != 0 {
		t.Errorf("Expected empty results, got %d", len(results))
	}
}
