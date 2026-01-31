// Package tools 股票数据工具
package tools

import (
	"context"
	"fmt"
)

// GetStockDataTool 获取股票数据工具
var GetStockDataTool = NewTool(
	"get_stock_data",
	"获取股票的 OHLCV (开盘价、最高价、最低价、收盘价、成交量) 历史数据",
	map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"symbol": map[string]interface{}{
				"type":        "string",
				"description": "股票代码，如 NVDA, AAPL, GOOGL",
			},
			"start_date": map[string]interface{}{
				"type":        "string",
				"description": "开始日期，格式 YYYY-MM-DD",
			},
			"end_date": map[string]interface{}{
				"type":        "string",
				"description": "结束日期，格式 YYYY-MM-DD",
			},
		},
		"required": []string{"symbol", "start_date", "end_date"},
	},
	func(ctx context.Context, args map[string]interface{}) (string, error) {
		symbol, _ := args["symbol"].(string)
		startDate, _ := args["start_date"].(string)
		endDate, _ := args["end_date"].(string)

		// 调用数据源获取数据
		data, err := fetchStockData(ctx, symbol, startDate, endDate)
		if err != nil {
			return "", err
		}

		return data, nil
	},
)

// fetchStockData 获取股票数据 (实际实现需要对接数据源)
func fetchStockData(ctx context.Context, symbol, startDate, endDate string) (string, error) {
	// 这里应该调用实际的数据源 API
	// 暂时返回模拟数据
	return fmt.Sprintf(`Stock Data for %s from %s to %s:
Date       | Open    | High    | Low     | Close   | Volume
-----------|---------|---------|---------|---------|----------
2024-05-10 | 875.50  | 892.30  | 871.20  | 887.90  | 45.2M
2024-05-09 | 860.25  | 878.40  | 858.10  | 875.50  | 42.8M
2024-05-08 | 845.80  | 865.20  | 842.50  | 860.25  | 38.5M
2024-05-07 | 832.40  | 848.90  | 830.10  | 845.80  | 35.2M
2024-05-06 | 820.15  | 835.60  | 818.30  | 832.40  | 32.1M

Summary:
- 5-day trend: Upward (+8.2%%)
- Average volume: 38.76M
- Volatility: Moderate`, symbol, startDate, endDate), nil
}
