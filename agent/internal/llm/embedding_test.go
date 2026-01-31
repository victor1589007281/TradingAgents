// Package llm Embedding 服务测试
package llm

import (
	"testing"

	"github.com/tradingagents/agent/internal/config"
)

func TestNewEmbeddingService(t *testing.T) {
	cfg := config.DefaultConfig()
	svc := NewEmbeddingService(cfg)

	if svc == nil {
		t.Fatal("EmbeddingService should not be nil")
	}
}

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a        []float64
		b        []float64
		expected float64
	}{
		{
			name:     "identical vectors",
			a:        []float64{1, 0, 0},
			b:        []float64{1, 0, 0},
			expected: 1.0,
		},
		{
			name:     "orthogonal vectors",
			a:        []float64{1, 0, 0},
			b:        []float64{0, 1, 0},
			expected: 0.0,
		},
		{
			name:     "opposite vectors",
			a:        []float64{1, 0, 0},
			b:        []float64{-1, 0, 0},
			expected: -1.0,
		},
		{
			name:     "similar vectors",
			a:        []float64{1, 1, 0},
			b:        []float64{1, 0, 0},
			expected: 0.7071067811865475, // 1/sqrt(2)
		},
		{
			name:     "empty vectors",
			a:        []float64{},
			b:        []float64{},
			expected: 0.0,
		},
		{
			name:     "different lengths",
			a:        []float64{1, 2, 3},
			b:        []float64{1, 2},
			expected: 0.0,
		},
		{
			name:     "zero vector",
			a:        []float64{0, 0, 0},
			b:        []float64{1, 2, 3},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CosineSimilarity(tt.a, tt.b)
			
			// 允许浮点数误差
			tolerance := 0.0001
			if diff := result - tt.expected; diff < -tolerance || diff > tolerance {
				t.Errorf("CosineSimilarity(%v, %v) = %f, expected %f", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestSqrt(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{0, 0},
		{1, 1},
		{4, 2},
		{9, 3},
		{2, 1.4142135623730951},
		{-1, 0}, // 负数返回 0
	}

	for _, tt := range tests {
		result := sqrt(tt.input)
		tolerance := 0.0001
		if diff := result - tt.expected; diff < -tolerance || diff > tolerance {
			t.Errorf("sqrt(%f) = %f, expected %f", tt.input, result, tt.expected)
		}
	}
}

// Benchmark 测试
func BenchmarkCosineSimilarity(b *testing.B) {
	// 创建 1536 维向量（OpenAI embedding 维度）
	a := make([]float64, 1536)
	c := make([]float64, 1536)
	for i := range a {
		a[i] = float64(i) / 1536.0
		c[i] = float64(1536-i) / 1536.0
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CosineSimilarity(a, c)
	}
}
