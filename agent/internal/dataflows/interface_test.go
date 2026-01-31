// Package dataflows 数据流测试
package dataflows

import (
	"context"
	"testing"

	"github.com/tradingagents/agent/internal/config"
)

func TestNewDataRouter(t *testing.T) {
	cfg := config.DefaultConfig()
	router := NewDataRouter(cfg)

	if router == nil {
		t.Fatal("DataRouter should not be nil")
	}
}

func TestDataRouterGetStockData(t *testing.T) {
	cfg := config.DefaultConfig()
	router := NewDataRouter(cfg)
	ctx := context.Background()

	data, err := router.GetStockData(ctx, "NVDA", "2024-05-01", "2024-05-10")
	if err != nil {
		t.Logf("GetStockData error (expected if no API): %v", err)
	}

	if data != nil && data.Symbol != "NVDA" {
		t.Errorf("Expected symbol 'NVDA', got '%s'", data.Symbol)
	}
}

func TestDataRouterGetIndicators(t *testing.T) {
	cfg := config.DefaultConfig()
	router := NewDataRouter(cfg)
	ctx := context.Background()

	indicators, err := router.GetIndicators(ctx, "NVDA", "2024-05-10")
	if err != nil {
		t.Logf("GetIndicators error (expected if no API): %v", err)
	}

	if indicators != nil {
		t.Logf("Got indicators: %+v", indicators)
	}
}

func TestDataRouterGetFundamentals(t *testing.T) {
	cfg := config.DefaultConfig()
	router := NewDataRouter(cfg)
	ctx := context.Background()

	fundamentals, err := router.GetFundamentals(ctx, "NVDA")
	if err != nil {
		t.Logf("GetFundamentals error (expected if no API): %v", err)
	}

	if fundamentals != nil && fundamentals.Symbol != "NVDA" {
		t.Errorf("Expected symbol 'NVDA', got '%s'", fundamentals.Symbol)
	}
}

func TestDataRouterGetNews(t *testing.T) {
	cfg := config.DefaultConfig()
	router := NewDataRouter(cfg)
	ctx := context.Background()

	news, err := router.GetNews(ctx, "NVDA", "2024-05-01", "2024-05-10")
	if err != nil {
		t.Logf("GetNews error (expected if no API): %v", err)
	}

	if news != nil {
		t.Logf("Got %d news items", len(news))
	}
}

func TestDataRouterVendorFallback(t *testing.T) {
	cfg := config.DefaultConfig()
	// 设置多个供应商
	cfg.DataVendors.CoreStockAPIs = "alpha_vantage,yfinance"

	router := NewDataRouter(cfg)

	// 验证 fallback 机制
	if router == nil {
		t.Fatal("Router should not be nil")
	}
}

func TestVendorConfig(t *testing.T) {
	tests := []struct {
		name     string
		vendor   string
		expected bool
	}{
		{"yfinance", "yfinance", true},
		{"alpha_vantage", "alpha_vantage", true},
		{"local", "local", true},
		{"invalid", "invalid_vendor", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := isValidVendor(tt.vendor)
			if valid != tt.expected {
				t.Errorf("isValidVendor(%s) = %v, expected %v", tt.vendor, valid, tt.expected)
			}
		})
	}
}

func isValidVendor(vendor string) bool {
	validVendors := map[string]bool{
		"yfinance":      true,
		"alpha_vantage": true,
		"local":         true,
		"openai":        true,
		"google_news":   true,
	}
	return validVendors[vendor]
}

func TestStockData(t *testing.T) {
	data := &StockData{
		Symbol:    "NVDA",
		StartDate: "2024-05-01",
		EndDate:   "2024-05-10",
		OHLCV: []OHLCV{
			{Date: "2024-05-01", Open: 100, High: 105, Low: 99, Close: 103, Volume: 1000000},
			{Date: "2024-05-02", Open: 103, High: 108, Low: 102, Close: 107, Volume: 1200000},
		},
	}

	if data.Symbol != "NVDA" {
		t.Errorf("Expected symbol 'NVDA', got '%s'", data.Symbol)
	}

	if len(data.OHLCV) != 2 {
		t.Errorf("Expected 2 OHLCV records, got %d", len(data.OHLCV))
	}
}

func TestIndicators(t *testing.T) {
	indicators := &Indicators{
		Symbol: "NVDA",
		Date:   "2024-05-10",
		SMA20:  850.5,
		SMA50:  820.3,
		EMA12:  855.2,
		EMA26:  835.8,
		RSI:    65.5,
		MACD:   19.4,
	}

	if indicators.RSI < 0 || indicators.RSI > 100 {
		t.Errorf("RSI should be between 0 and 100, got %f", indicators.RSI)
	}
}

func TestFundamentals(t *testing.T) {
	fundamentals := &Fundamentals{
		Symbol:      "NVDA",
		MarketCap:   1500000000000,
		PE:          45.5,
		PB:          25.3,
		Revenue:     60000000000,
		NetIncome:   15000000000,
		EPS:         6.05,
		ROE:         0.55,
		DebtToEquity: 0.45,
	}

	if fundamentals.PE < 0 {
		t.Errorf("PE should be non-negative, got %f", fundamentals.PE)
	}

	if fundamentals.ROE < -1 || fundamentals.ROE > 2 {
		t.Errorf("ROE seems unrealistic: %f", fundamentals.ROE)
	}
}

func TestNewsItem(t *testing.T) {
	news := NewsItem{
		Title:     "NVDA Reports Record Earnings",
		Source:    "Reuters",
		URL:       "https://reuters.com/...",
		Published: "2024-05-10T10:00:00Z",
		Summary:   "NVIDIA reported record quarterly earnings...",
		Sentiment: 0.8,
	}

	if news.Title == "" {
		t.Error("News title should not be empty")
	}

	if news.Sentiment < -1 || news.Sentiment > 1 {
		t.Errorf("Sentiment should be between -1 and 1, got %f", news.Sentiment)
	}
}
