// Package dataflows Alpha Vantage 数据源适配
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

// AlphaVantageVendor Alpha Vantage 数据源
type AlphaVantageVendor struct {
	cfg    *config.Config
	apiKey string
}

// NewAlphaVantageVendor 创建 Alpha Vantage 数据源
func NewAlphaVantageVendor(cfg *config.Config) *AlphaVantageVendor {
	return &AlphaVantageVendor{
		cfg:    cfg,
		apiKey: cfg.AlphaVantageAPIKey,
	}
}

func (a *AlphaVantageVendor) Name() string { return "alpha_vantage" }

const alphaVantageBaseURL = "https://www.alphavantage.co/query"

// makeRequest 发送 Alpha Vantage API 请求
func (a *AlphaVantageVendor) makeRequest(ctx context.Context, params url.Values) (map[string]interface{}, error) {
	params.Set("apikey", a.apiKey)
	fullURL := alphaVantageBaseURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// 检查是否有错误或限流信息
	if note, ok := result["Note"].(string); ok {
		return nil, fmt.Errorf("API limit reached: %s", note)
	}
	if errMsg, ok := result["Error Message"].(string); ok {
		return nil, fmt.Errorf("API error: %s", errMsg)
	}

	return result, nil
}

// GetStockData 获取股票数据
func (a *AlphaVantageVendor) GetStockData(ctx context.Context, symbol, startDate, endDate string) (string, error) {
	params := url.Values{}
	params.Set("function", "TIME_SERIES_DAILY")
	params.Set("symbol", symbol)
	params.Set("outputsize", "full")

	result, err := a.makeRequest(ctx, params)
	if err != nil {
		return "", err
	}

	timeSeries, ok := result["Time Series (Daily)"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}

	output := fmt.Sprintf("Stock Data for %s (from Alpha Vantage):\n", symbol)
	output += "Date       | Open    | High    | Low     | Close   | Volume\n"
	output += "-----------|---------|---------|---------|---------|----------\n"

	for date, data := range timeSeries {
		if date < startDate || date > endDate {
			continue
		}
		dayData := data.(map[string]interface{})
		output += fmt.Sprintf("%s | %s | %s | %s | %s | %s\n",
			date,
			dayData["1. open"],
			dayData["2. high"],
			dayData["3. low"],
			dayData["4. close"],
			dayData["5. volume"],
		)
	}

	return output, nil
}

// GetIndicators 获取技术指标
func (a *AlphaVantageVendor) GetIndicators(ctx context.Context, symbol, date string) (string, error) {
	output := fmt.Sprintf("Technical Indicators for %s (from Alpha Vantage):\n\n", symbol)

	// 获取 RSI
	rsiParams := url.Values{}
	rsiParams.Set("function", "RSI")
	rsiParams.Set("symbol", symbol)
	rsiParams.Set("interval", "daily")
	rsiParams.Set("time_period", "14")
	rsiParams.Set("series_type", "close")

	rsiResult, err := a.makeRequest(ctx, rsiParams)
	if err == nil {
		output += "RSI (14):\n"
		if rsiData, ok := rsiResult["Technical Analysis: RSI"].(map[string]interface{}); ok {
			count := 0
			for rsiDate, value := range rsiData {
				if count >= 5 {
					break
				}
				v := value.(map[string]interface{})
				output += fmt.Sprintf("  %s: %s\n", rsiDate, v["RSI"])
				count++
			}
		}
	}

	// 获取 MACD
	macdParams := url.Values{}
	macdParams.Set("function", "MACD")
	macdParams.Set("symbol", symbol)
	macdParams.Set("interval", "daily")
	macdParams.Set("series_type", "close")

	macdResult, err := a.makeRequest(ctx, macdParams)
	if err == nil {
		output += "\nMACD:\n"
		if macdData, ok := macdResult["Technical Analysis: MACD"].(map[string]interface{}); ok {
			count := 0
			for macdDate, value := range macdData {
				if count >= 5 {
					break
				}
				v := value.(map[string]interface{})
				output += fmt.Sprintf("  %s: MACD=%s, Signal=%s, Hist=%s\n",
					macdDate, v["MACD"], v["MACD_Signal"], v["MACD_Hist"])
				count++
			}
		}
	}

	return output, nil
}

// GetFundamentals 获取基本面数据
func (a *AlphaVantageVendor) GetFundamentals(ctx context.Context, symbol string) (string, error) {
	params := url.Values{}
	params.Set("function", "OVERVIEW")
	params.Set("symbol", symbol)

	result, err := a.makeRequest(ctx, params)
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("Company Fundamentals for %s:\n\n", symbol)
	output += fmt.Sprintf("Name: %v\n", result["Name"])
	output += fmt.Sprintf("Sector: %v\n", result["Sector"])
	output += fmt.Sprintf("Industry: %v\n", result["Industry"])
	output += fmt.Sprintf("Market Cap: %v\n", result["MarketCapitalization"])
	output += fmt.Sprintf("P/E Ratio: %v\n", result["PERatio"])
	output += fmt.Sprintf("P/B Ratio: %v\n", result["PriceToBookRatio"])
	output += fmt.Sprintf("EPS: %v\n", result["EPS"])
	output += fmt.Sprintf("Dividend Yield: %v\n", result["DividendYield"])
	output += fmt.Sprintf("Profit Margin: %v\n", result["ProfitMargin"])
	output += fmt.Sprintf("Operating Margin: %v\n", result["OperatingMarginTTM"])
	output += fmt.Sprintf("ROE: %v\n", result["ReturnOnEquityTTM"])
	output += fmt.Sprintf("ROA: %v\n", result["ReturnOnAssetsTTM"])
	output += fmt.Sprintf("Revenue TTM: %v\n", result["RevenueTTM"])
	output += fmt.Sprintf("Gross Profit TTM: %v\n", result["GrossProfitTTM"])
	output += fmt.Sprintf("52 Week High: %v\n", result["52WeekHigh"])
	output += fmt.Sprintf("52 Week Low: %v\n", result["52WeekLow"])
	output += fmt.Sprintf("50 Day MA: %v\n", result["50DayMovingAverage"])
	output += fmt.Sprintf("200 Day MA: %v\n", result["200DayMovingAverage"])

	return output, nil
}

// GetBalanceSheet 获取资产负债表
func (a *AlphaVantageVendor) GetBalanceSheet(ctx context.Context, symbol, period string) (string, error) {
	params := url.Values{}
	params.Set("function", "BALANCE_SHEET")
	params.Set("symbol", symbol)

	result, err := a.makeRequest(ctx, params)
	if err != nil {
		return "", err
	}

	reportKey := "annualReports"
	if period == "quarterly" {
		reportKey = "quarterlyReports"
	}

	reports, ok := result[reportKey].([]interface{})
	if !ok || len(reports) == 0 {
		return "", fmt.Errorf("no balance sheet data available")
	}

	output := fmt.Sprintf("Balance Sheet for %s (%s):\n\n", symbol, period)
	latestReport := reports[0].(map[string]interface{})

	output += fmt.Sprintf("Fiscal Date: %v\n", latestReport["fiscalDateEnding"])
	output += fmt.Sprintf("Total Assets: %v\n", latestReport["totalAssets"])
	output += fmt.Sprintf("Total Current Assets: %v\n", latestReport["totalCurrentAssets"])
	output += fmt.Sprintf("Cash: %v\n", latestReport["cashAndCashEquivalentsAtCarryingValue"])
	output += fmt.Sprintf("Total Liabilities: %v\n", latestReport["totalLiabilities"])
	output += fmt.Sprintf("Total Current Liabilities: %v\n", latestReport["totalCurrentLiabilities"])
	output += fmt.Sprintf("Total Shareholder Equity: %v\n", latestReport["totalShareholderEquity"])

	return output, nil
}

// GetCashflow 获取现金流量表
func (a *AlphaVantageVendor) GetCashflow(ctx context.Context, symbol, period string) (string, error) {
	params := url.Values{}
	params.Set("function", "CASH_FLOW")
	params.Set("symbol", symbol)

	result, err := a.makeRequest(ctx, params)
	if err != nil {
		return "", err
	}

	reportKey := "annualReports"
	if period == "quarterly" {
		reportKey = "quarterlyReports"
	}

	reports, ok := result[reportKey].([]interface{})
	if !ok || len(reports) == 0 {
		return "", fmt.Errorf("no cash flow data available")
	}

	output := fmt.Sprintf("Cash Flow Statement for %s (%s):\n\n", symbol, period)
	latestReport := reports[0].(map[string]interface{})

	output += fmt.Sprintf("Fiscal Date: %v\n", latestReport["fiscalDateEnding"])
	output += fmt.Sprintf("Operating Cash Flow: %v\n", latestReport["operatingCashflow"])
	output += fmt.Sprintf("Capital Expenditures: %v\n", latestReport["capitalExpenditures"])
	output += fmt.Sprintf("Dividend Payout: %v\n", latestReport["dividendPayout"])
	output += fmt.Sprintf("Net Income: %v\n", latestReport["netIncome"])

	return output, nil
}

// GetIncomeStatement 获取利润表
func (a *AlphaVantageVendor) GetIncomeStatement(ctx context.Context, symbol, period string) (string, error) {
	params := url.Values{}
	params.Set("function", "INCOME_STATEMENT")
	params.Set("symbol", symbol)

	result, err := a.makeRequest(ctx, params)
	if err != nil {
		return "", err
	}

	reportKey := "annualReports"
	if period == "quarterly" {
		reportKey = "quarterlyReports"
	}

	reports, ok := result[reportKey].([]interface{})
	if !ok || len(reports) == 0 {
		return "", fmt.Errorf("no income statement data available")
	}

	output := fmt.Sprintf("Income Statement for %s (%s):\n\n", symbol, period)
	latestReport := reports[0].(map[string]interface{})

	output += fmt.Sprintf("Fiscal Date: %v\n", latestReport["fiscalDateEnding"])
	output += fmt.Sprintf("Total Revenue: %v\n", latestReport["totalRevenue"])
	output += fmt.Sprintf("Gross Profit: %v\n", latestReport["grossProfit"])
	output += fmt.Sprintf("Operating Income: %v\n", latestReport["operatingIncome"])
	output += fmt.Sprintf("Net Income: %v\n", latestReport["netIncome"])
	output += fmt.Sprintf("EPS: %v\n", latestReport["reportedEPS"])

	return output, nil
}

// GetNews 获取新闻
func (a *AlphaVantageVendor) GetNews(ctx context.Context, query, startDate, endDate string) (string, error) {
	params := url.Values{}
	params.Set("function", "NEWS_SENTIMENT")
	params.Set("tickers", query)

	result, err := a.makeRequest(ctx, params)
	if err != nil {
		return "", err
	}

	feed, ok := result["feed"].([]interface{})
	if !ok {
		return "", fmt.Errorf("no news data available")
	}

	output := fmt.Sprintf("News for %s:\n\n", query)
	for i, item := range feed {
		if i >= 10 {
			break
		}
		news := item.(map[string]interface{})
		output += fmt.Sprintf("**%s**\n", news["title"])
		output += fmt.Sprintf("Source: %v | Time: %v\n", news["source"], news["time_published"])
		output += fmt.Sprintf("Sentiment: %v (Score: %v)\n", news["overall_sentiment_label"], news["overall_sentiment_score"])
		output += fmt.Sprintf("Summary: %v\n\n", news["summary"])
	}

	return output, nil
}

// GetGlobalNews 获取全球新闻
func (a *AlphaVantageVendor) GetGlobalNews(ctx context.Context, date string, lookBackDays, limit int) (string, error) {
	params := url.Values{}
	params.Set("function", "NEWS_SENTIMENT")
	params.Set("topics", "economy_macro,finance,technology")

	result, err := a.makeRequest(ctx, params)
	if err != nil {
		return "", err
	}

	feed, ok := result["feed"].([]interface{})
	if !ok {
		return "", fmt.Errorf("no global news data available")
	}

	output := "Global Market News:\n\n"
	for i, item := range feed {
		if i >= limit {
			break
		}
		news := item.(map[string]interface{})
		output += fmt.Sprintf("**%s**\n", news["title"])
		output += fmt.Sprintf("Source: %v | Time: %v\n", news["source"], news["time_published"])
		output += fmt.Sprintf("Summary: %v\n\n", news["summary"])
	}

	return output, nil
}

// GetInsiderSentiment 获取内部人士情绪
func (a *AlphaVantageVendor) GetInsiderSentiment(ctx context.Context, symbol string) (string, error) {
	// Alpha Vantage 不直接提供 insider sentiment，使用新闻情绪替代
	return a.GetNews(ctx, symbol, "", "")
}

// GetInsiderTransactions 获取内部人士交易
func (a *AlphaVantageVendor) GetInsiderTransactions(ctx context.Context, symbol string) (string, error) {
	// Alpha Vantage 不直接提供 insider transactions
	return fmt.Sprintf("Insider transactions for %s not available from Alpha Vantage", symbol), nil
}
