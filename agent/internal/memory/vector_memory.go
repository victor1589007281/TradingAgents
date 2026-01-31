// Package memory 向量记忆系统
package memory

import (
	"context"
	"sort"
	"sync"

	"github.com/tradingagents/agent/internal/config"
	"github.com/tradingagents/agent/internal/llm"
)

// MemoryRecord 记忆记录
type MemoryRecord struct {
	ID              string    `json:"id"`
	Situation       string    `json:"situation"`
	Recommendation  string    `json:"recommendation"`
	Embedding       []float64 `json:"embedding"`
	SimilarityScore float64   `json:"similarity_score"`
}

// SituationAdvice 情境和建议
type SituationAdvice struct {
	Situation      string
	Recommendation string
}

// FinancialSituationMemory 金融情境记忆
type FinancialSituationMemory struct {
	name      string
	embedder  *llm.EmbeddingService
	records   []MemoryRecord
	mu        sync.RWMutex
	idCounter int
}

// NewFinancialSituationMemory 创建金融情境记忆
func NewFinancialSituationMemory(name string, cfg *config.Config) *FinancialSituationMemory {
	return &FinancialSituationMemory{
		name:     name,
		embedder: llm.NewEmbeddingService(cfg),
		records:  make([]MemoryRecord, 0),
	}
}

// AddSituations 添加情境和建议
func (m *FinancialSituationMemory) AddSituations(ctx context.Context, situations []SituationAdvice) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, sa := range situations {
		// 获取 embedding
		embedding, err := m.embedder.Embed(ctx, sa.Situation)
		if err != nil {
			// 如果获取 embedding 失败，使用空向量
			embedding = make([]float64, 1536)
		}

		m.idCounter++
		record := MemoryRecord{
			ID:             string(rune(m.idCounter)),
			Situation:      sa.Situation,
			Recommendation: sa.Recommendation,
			Embedding:      embedding,
		}
		m.records = append(m.records, record)
	}

	return nil
}

// GetMemories 获取相似记忆
func (m *FinancialSituationMemory) GetMemories(ctx context.Context, currentSituation string, nMatches int) ([]MemoryRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.records) == 0 {
		return nil, nil
	}

	// 获取当前情境的 embedding
	queryEmbedding, err := m.embedder.Embed(ctx, currentSituation)
	if err != nil {
		// 如果获取 embedding 失败，返回空结果
		return nil, nil
	}

	// 计算相似度
	results := make([]MemoryRecord, len(m.records))
	for i, record := range m.records {
		similarity := llm.CosineSimilarity(queryEmbedding, record.Embedding)
		results[i] = MemoryRecord{
			ID:              record.ID,
			Situation:       record.Situation,
			Recommendation:  record.Recommendation,
			SimilarityScore: similarity,
		}
	}

	// 按相似度排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].SimilarityScore > results[j].SimilarityScore
	})

	// 返回前 n 个
	if nMatches > len(results) {
		nMatches = len(results)
	}

	return results[:nMatches], nil
}

// GetName 获取记忆名称
func (m *FinancialSituationMemory) GetName() string {
	return m.name
}

// Count 获取记录数量
func (m *FinancialSituationMemory) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.records)
}

// Clear 清空记忆
func (m *FinancialSituationMemory) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.records = make([]MemoryRecord, 0)
	m.idCounter = 0
}
