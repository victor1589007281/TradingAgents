// Package config 配置管理测试
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	// 测试默认值
	if cfg.LLMProvider != "openai" {
		t.Errorf("Expected LLMProvider to be 'openai', got '%s'", cfg.LLMProvider)
	}

	if cfg.DeepThinkLLM != "o4-mini" {
		t.Errorf("Expected DeepThinkLLM to be 'o4-mini', got '%s'", cfg.DeepThinkLLM)
	}

	if cfg.QuickThinkLLM != "gpt-4o-mini" {
		t.Errorf("Expected QuickThinkLLM to be 'gpt-4o-mini', got '%s'", cfg.QuickThinkLLM)
	}

	if cfg.MaxDebateRounds != 1 {
		t.Errorf("Expected MaxDebateRounds to be 1, got %d", cfg.MaxDebateRounds)
	}

	if cfg.MaxRiskDiscussRounds != 1 {
		t.Errorf("Expected MaxRiskDiscussRounds to be 1, got %d", cfg.MaxRiskDiscussRounds)
	}

	// 测试数据源配置
	if cfg.DataVendors.CoreStockAPIs != "yfinance" {
		t.Errorf("Expected CoreStockAPIs to be 'yfinance', got '%s'", cfg.DataVendors.CoreStockAPIs)
	}
}

func TestLoadFromFile(t *testing.T) {
	// 测试加载不存在的文件（应该返回默认配置）
	cfg, err := LoadFromFile("nonexistent.yaml")
	if err != nil {
		t.Errorf("Expected no error for nonexistent file, got %v", err)
	}

	if cfg == nil {
		t.Error("Expected config to be non-nil")
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test_config.yaml")

	// 创建并保存配置
	cfg := DefaultConfig()
	cfg.DeepThinkLLM = "test-model"
	cfg.MaxDebateRounds = 5

	err := cfg.SaveToFile(configPath)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// 加载配置
	loadedCfg, err := LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 验证加载的配置
	if loadedCfg.DeepThinkLLM != "test-model" {
		t.Errorf("Expected DeepThinkLLM to be 'test-model', got '%s'", loadedCfg.DeepThinkLLM)
	}

	if loadedCfg.MaxDebateRounds != 5 {
		t.Errorf("Expected MaxDebateRounds to be 5, got %d", loadedCfg.MaxDebateRounds)
	}
}

func TestEnvVarOverride(t *testing.T) {
	// 设置环境变量
	os.Setenv("OPENAI_API_KEY", "test-api-key")
	os.Setenv("ALPHA_VANTAGE_API_KEY", "test-av-key")
	defer func() {
		os.Unsetenv("OPENAI_API_KEY")
		os.Unsetenv("ALPHA_VANTAGE_API_KEY")
	}()

	cfg := DefaultConfig()

	if cfg.APIKey != "test-api-key" {
		t.Errorf("Expected APIKey to be 'test-api-key', got '%s'", cfg.APIKey)
	}

	if cfg.AlphaVantageAPIKey != "test-av-key" {
		t.Errorf("Expected AlphaVantageAPIKey to be 'test-av-key', got '%s'", cfg.AlphaVantageAPIKey)
	}
}

func TestDataVendorConfig(t *testing.T) {
	cfg := DefaultConfig()

	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{"CoreStockAPIs", cfg.DataVendors.CoreStockAPIs, "yfinance"},
		{"TechnicalIndicators", cfg.DataVendors.TechnicalIndicators, "yfinance"},
		{"FundamentalData", cfg.DataVendors.FundamentalData, "alpha_vantage"},
		{"NewsData", cfg.DataVendors.NewsData, "alpha_vantage"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("Expected %s to be '%s', got '%s'", tt.name, tt.expected, tt.got)
			}
		})
	}
}
