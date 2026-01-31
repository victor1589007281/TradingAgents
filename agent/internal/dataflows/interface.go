// Package dataflows 数据源接口和路由
package dataflows

import (
	"context"
	"fmt"

	"github.com/tradingagents/agent/internal/config"
)

// DataVendor 数据源接口
type DataVendor interface {
	// Name 数据源名称
	Name() string
	// GetStockData 获取股票数据
	GetStockData(ctx context.Context, symbol, startDate, endDate string) (string, error)
	// GetIndicators 获取技术指标
	GetIndicators(ctx context.Context, symbol, date string) (string, error)
	// GetFundamentals 获取基本面数据
	GetFundamentals(ctx context.Context, symbol string) (string, error)
	// GetBalanceSheet 获取资产负债表
	GetBalanceSheet(ctx context.Context, symbol, period string) (string, error)
	// GetCashflow 获取现金流量表
	GetCashflow(ctx context.Context, symbol, period string) (string, error)
	// GetIncomeStatement 获取利润表
	GetIncomeStatement(ctx context.Context, symbol, period string) (string, error)
	// GetNews 获取新闻
	GetNews(ctx context.Context, query, startDate, endDate string) (string, error)
	// GetGlobalNews 获取全球新闻
	GetGlobalNews(ctx context.Context, date string, lookBackDays, limit int) (string, error)
	// GetInsiderSentiment 获取内部人士情绪
	GetInsiderSentiment(ctx context.Context, symbol string) (string, error)
	// GetInsiderTransactions 获取内部人士交易
	GetInsiderTransactions(ctx context.Context, symbol string) (string, error)
}

// DataRouter 数据路由器
type DataRouter struct {
	cfg     *config.Config
	vendors map[string]DataVendor
}

// NewDataRouter 创建数据路由器
func NewDataRouter(cfg *config.Config) *DataRouter {
	router := &DataRouter{
		cfg:     cfg,
		vendors: make(map[string]DataVendor),
	}

	// 注册数据源
	router.vendors["yfinance"] = NewYFinanceVendor(cfg)
	router.vendors["alpha_vantage"] = NewAlphaVantageVendor(cfg)
	router.vendors["mock"] = NewMockVendor()

	return router
}

// GetVendor 获取数据源
func (r *DataRouter) GetVendor(category string) DataVendor {
	var vendorName string
	switch category {
	case "core_stock_apis":
		vendorName = r.cfg.DataVendors.CoreStockAPIs
	case "technical_indicators":
		vendorName = r.cfg.DataVendors.TechnicalIndicators
	case "fundamental_data":
		vendorName = r.cfg.DataVendors.FundamentalData
	case "news_data":
		vendorName = r.cfg.DataVendors.NewsData
	default:
		vendorName = "mock"
	}

	if vendor, ok := r.vendors[vendorName]; ok {
		return vendor
	}

	// 回退到 mock
	return r.vendors["mock"]
}

// RouteStockData 路由获取股票数据
func (r *DataRouter) RouteStockData(ctx context.Context, symbol, startDate, endDate string) (string, error) {
	vendor := r.GetVendor("core_stock_apis")
	data, err := vendor.GetStockData(ctx, symbol, startDate, endDate)
	if err != nil {
		// 尝试回退到其他数据源
		if fallback, ok := r.vendors["mock"]; ok {
			return fallback.GetStockData(ctx, symbol, startDate, endDate)
		}
		return "", err
	}
	return data, nil
}

// RouteIndicators 路由获取技术指标
func (r *DataRouter) RouteIndicators(ctx context.Context, symbol, date string) (string, error) {
	vendor := r.GetVendor("technical_indicators")
	data, err := vendor.GetIndicators(ctx, symbol, date)
	if err != nil {
		if fallback, ok := r.vendors["mock"]; ok {
			return fallback.GetIndicators(ctx, symbol, date)
		}
		return "", err
	}
	return data, nil
}

// RouteFundamentals 路由获取基本面数据
func (r *DataRouter) RouteFundamentals(ctx context.Context, symbol string) (string, error) {
	vendor := r.GetVendor("fundamental_data")
	data, err := vendor.GetFundamentals(ctx, symbol)
	if err != nil {
		if fallback, ok := r.vendors["mock"]; ok {
			return fallback.GetFundamentals(ctx, symbol)
		}
		return "", err
	}
	return data, nil
}

// RouteNews 路由获取新闻
func (r *DataRouter) RouteNews(ctx context.Context, query, startDate, endDate string) (string, error) {
	vendor := r.GetVendor("news_data")
	data, err := vendor.GetNews(ctx, query, startDate, endDate)
	if err != nil {
		if fallback, ok := r.vendors["mock"]; ok {
			return fallback.GetNews(ctx, query, startDate, endDate)
		}
		return "", err
	}
	return data, nil
}

// MockVendor Mock 数据源（用于测试和开发）
type MockVendor struct{}

// NewMockVendor 创建 Mock 数据源
func NewMockVendor() *MockVendor {
	return &MockVendor{}
}

func (m *MockVendor) Name() string { return "mock" }

func (m *MockVendor) GetStockData(ctx context.Context, symbol, startDate, endDate string) (string, error) {
	return fmt.Sprintf("Mock stock data for %s from %s to %s", symbol, startDate, endDate), nil
}

func (m *MockVendor) GetIndicators(ctx context.Context, symbol, date string) (string, error) {
	return fmt.Sprintf("Mock indicators for %s on %s", symbol, date), nil
}

func (m *MockVendor) GetFundamentals(ctx context.Context, symbol string) (string, error) {
	return fmt.Sprintf("Mock fundamentals for %s", symbol), nil
}

func (m *MockVendor) GetBalanceSheet(ctx context.Context, symbol, period string) (string, error) {
	return fmt.Sprintf("Mock balance sheet for %s (%s)", symbol, period), nil
}

func (m *MockVendor) GetCashflow(ctx context.Context, symbol, period string) (string, error) {
	return fmt.Sprintf("Mock cashflow for %s (%s)", symbol, period), nil
}

func (m *MockVendor) GetIncomeStatement(ctx context.Context, symbol, period string) (string, error) {
	return fmt.Sprintf("Mock income statement for %s (%s)", symbol, period), nil
}

func (m *MockVendor) GetNews(ctx context.Context, query, startDate, endDate string) (string, error) {
	return fmt.Sprintf("Mock news for %s from %s to %s", query, startDate, endDate), nil
}

func (m *MockVendor) GetGlobalNews(ctx context.Context, date string, lookBackDays, limit int) (string, error) {
	return fmt.Sprintf("Mock global news for %s", date), nil
}

func (m *MockVendor) GetInsiderSentiment(ctx context.Context, symbol string) (string, error) {
	return fmt.Sprintf("Mock insider sentiment for %s", symbol), nil
}

func (m *MockVendor) GetInsiderTransactions(ctx context.Context, symbol string) (string, error) {
	return fmt.Sprintf("Mock insider transactions for %s", symbol), nil
}
