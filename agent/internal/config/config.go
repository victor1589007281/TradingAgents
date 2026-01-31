// Package config 配置管理
package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// DataVendorConfig 数据源配置
type DataVendorConfig struct {
	CoreStockAPIs       string `yaml:"core_stock_apis" json:"core_stock_apis"`
	TechnicalIndicators string `yaml:"technical_indicators" json:"technical_indicators"`
	FundamentalData     string `yaml:"fundamental_data" json:"fundamental_data"`
	NewsData            string `yaml:"news_data" json:"news_data"`
}

// Config 应用配置
type Config struct {
	// 项目路径
	ProjectDir   string `yaml:"project_dir" json:"project_dir"`
	ResultsDir   string `yaml:"results_dir" json:"results_dir"`
	DataCacheDir string `yaml:"data_cache_dir" json:"data_cache_dir"`

	// LLM 配置
	LLMProvider   string `yaml:"llm_provider" json:"llm_provider"`
	DeepThinkLLM  string `yaml:"deep_think_llm" json:"deep_think_llm"`
	QuickThinkLLM string `yaml:"quick_think_llm" json:"quick_think_llm"`
	BackendURL    string `yaml:"backend_url" json:"backend_url"`
	APIKey        string `yaml:"api_key" json:"api_key"`

	// Embedding 配置
	EmbeddingModel string `yaml:"embedding_model" json:"embedding_model"`

	// 辩论配置
	MaxDebateRounds      int `yaml:"max_debate_rounds" json:"max_debate_rounds"`
	MaxRiskDiscussRounds int `yaml:"max_risk_discuss_rounds" json:"max_risk_discuss_rounds"`
	MaxRecurLimit        int `yaml:"max_recur_limit" json:"max_recur_limit"`

	// 数据源配置
	DataVendors DataVendorConfig `yaml:"data_vendors" json:"data_vendors"`

	// API Keys
	AlphaVantageAPIKey string `yaml:"alpha_vantage_api_key" json:"alpha_vantage_api_key"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	projectDir, _ := os.Getwd()
	return &Config{
		ProjectDir:   projectDir,
		ResultsDir:   filepath.Join(projectDir, "results"),
		DataCacheDir: filepath.Join(projectDir, "data_cache"),

		LLMProvider:   "openai",
		DeepThinkLLM:  "o4-mini",
		QuickThinkLLM: "gpt-4o-mini",
		BackendURL:    "https://api.openai.com/v1",
		APIKey:        os.Getenv("OPENAI_API_KEY"),

		EmbeddingModel: "text-embedding-3-small",

		MaxDebateRounds:      1,
		MaxRiskDiscussRounds: 1,
		MaxRecurLimit:        100,

		DataVendors: DataVendorConfig{
			CoreStockAPIs:       "yfinance",
			TechnicalIndicators: "yfinance",
			FundamentalData:     "alpha_vantage",
			NewsData:            "alpha_vantage",
		},

		AlphaVantageAPIKey: os.Getenv("ALPHA_VANTAGE_API_KEY"),
	}
}

// LoadFromFile 从文件加载配置
func LoadFromFile(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	// 环境变量覆盖
	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		cfg.APIKey = apiKey
	}
	if avKey := os.Getenv("ALPHA_VANTAGE_API_KEY"); avKey != "" {
		cfg.AlphaVantageAPIKey = avKey
	}

	return cfg, nil
}

// SaveToFile 保存配置到文件
func (c *Config) SaveToFile(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
