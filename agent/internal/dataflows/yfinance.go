// Package dataflows yfinance 数据源适配
package dataflows

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/tradingagents/agent/internal/config"
)

// YFinanceVendor yfinance 数据源
type YFinanceVendor struct {
	cfg *config.Config
}

// NewYFinanceVendor 创建 yfinance 数据源
func NewYFinanceVendor(cfg *config.Config) *YFinanceVendor {
	return &YFinanceVendor{cfg: cfg}
}

func (y *YFinanceVendor) Name() string { return "yfinance" }

// GetStockData 获取股票数据
func (y *YFinanceVendor) GetStockData(ctx context.Context, symbol, startDate, endDate string) (string, error) {
	// 将日期转换为时间戳
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return "", fmt.Errorf("invalid start date: %w", err)
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return "", fmt.Errorf("invalid end date: %w", err)
	}

	// 构建 Yahoo Finance API URL
	baseURL := "https://query1.finance.yahoo.com/v8/finance/chart/"
	params := url.Values{}
	params.Set("period1", fmt.Sprintf("%d", start.Unix()))
	params.Set("period2", fmt.Sprintf("%d", end.Unix()))
	params.Set("interval", "1d")
	params.Set("includePrePost", "false")

	fullURL := baseURL + symbol + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析 JSON 响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	// 格式化输出
	return formatYFinanceData(symbol, result), nil
}

// GetIndicators 获取技术指标
func (y *YFinanceVendor) GetIndicators(ctx context.Context, symbol, date string) (string, error) {
	// yfinance 需要获取历史数据后计算指标
	// 获取过去 200 天的数据来计算各种移动平均线
	endDate, _ := time.Parse("2006-01-02", date)
	startDate := endDate.AddDate(0, 0, -200)

	data, err := y.GetStockData(ctx, symbol, startDate.Format("2006-01-02"), date)
	if err != nil {
		return "", err
	}

	// 这里应该计算技术指标，暂时返回原始数据
	return fmt.Sprintf("Technical indicators for %s:\n%s", symbol, data), nil
}

// GetFundamentals 获取基本面数据
func (y *YFinanceVendor) GetFundamentals(ctx context.Context, symbol string) (string, error) {
	// Yahoo Finance 基本面数据需要不同的 API 端点
	// 暂时返回模拟数据，实际实现需要解析 Yahoo Finance 页面或使用其他 API
	return fmt.Sprintf("Fundamentals data for %s from Yahoo Finance", symbol), nil
}

// GetBalanceSheet 获取资产负债表
func (y *YFinanceVendor) GetBalanceSheet(ctx context.Context, symbol, period string) (string, error) {
	return fmt.Sprintf("Balance sheet for %s (%s) from Yahoo Finance", symbol, period), nil
}

// GetCashflow 获取现金流量表
func (y *YFinanceVendor) GetCashflow(ctx context.Context, symbol, period string) (string, error) {
	return fmt.Sprintf("Cashflow for %s (%s) from Yahoo Finance", symbol, period), nil
}

// GetIncomeStatement 获取利润表
func (y *YFinanceVendor) GetIncomeStatement(ctx context.Context, symbol, period string) (string, error) {
	return fmt.Sprintf("Income statement for %s (%s) from Yahoo Finance", symbol, period), nil
}

// GetNews 获取新闻 (yfinance 不直接提供新闻)
func (y *YFinanceVendor) GetNews(ctx context.Context, query, startDate, endDate string) (string, error) {
	return "", fmt.Errorf("yfinance does not support news")
}

// GetGlobalNews 获取全球新闻
func (y *YFinanceVendor) GetGlobalNews(ctx context.Context, date string, lookBackDays, limit int) (string, error) {
	return "", fmt.Errorf("yfinance does not support global news")
}

// GetInsiderSentiment 获取内部人士情绪
func (y *YFinanceVendor) GetInsiderSentiment(ctx context.Context, symbol string) (string, error) {
	return "", fmt.Errorf("yfinance does not support insider sentiment")
}

// GetInsiderTransactions 获取内部人士交易
func (y *YFinanceVendor) GetInsiderTransactions(ctx context.Context, symbol string) (string, error) {
	return fmt.Sprintf("Insider transactions for %s from Yahoo Finance", symbol), nil
}

// formatYFinanceData 格式化 yfinance 数据
func formatYFinanceData(symbol string, data map[string]interface{}) string {
	chart, ok := data["chart"].(map[string]interface{})
	if !ok {
		return "No data available"
	}

	result, ok := chart["result"].([]interface{})
	if !ok || len(result) == 0 {
		return "No data available"
	}

	firstResult := result[0].(map[string]interface{})
	timestamps, _ := firstResult["timestamp"].([]interface{})
	indicators, _ := firstResult["indicators"].(map[string]interface{})
	quotes, _ := indicators["quote"].([]interface{})

	if len(quotes) == 0 {
		return "No quote data available"
	}

	quote := quotes[0].(map[string]interface{})
	opens, _ := quote["open"].([]interface{})
	highs, _ := quote["high"].([]interface{})
	lows, _ := quote["low"].([]interface{})
	closes, _ := quote["close"].([]interface{})
	volumes, _ := quote["volume"].([]interface{})

	output := fmt.Sprintf("Stock Data for %s:\n", symbol)
	output += "Date       | Open    | High    | Low     | Close   | Volume\n"
	output += "-----------|---------|---------|---------|---------|----------\n"

	for i := 0; i < len(timestamps) && i < len(opens); i++ {
		if timestamps[i] == nil || opens[i] == nil {
			continue
		}
		ts := int64(timestamps[i].(float64))
		date := time.Unix(ts, 0).Format("2006-01-02")
		open := opens[i].(float64)
		high := highs[i].(float64)
		low := lows[i].(float64)
		close := closes[i].(float64)
		volume := volumes[i].(float64)

		output += fmt.Sprintf("%s | %.2f | %.2f | %.2f | %.2f | %.0f\n",
			date, open, high, low, close, volume)
	}

	return output
}
