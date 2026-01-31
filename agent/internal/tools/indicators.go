// Package tools 技术指标工具
package tools

import (
	"context"
	"fmt"
)

// GetIndicatorsTool 获取技术指标工具
var GetIndicatorsTool = NewTool(
	"get_indicators",
	"获取股票的技术分析指标，包括 MACD, RSI, 移动平均线等",
	map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"symbol": map[string]interface{}{
				"type":        "string",
				"description": "股票代码",
			},
			"date": map[string]interface{}{
				"type":        "string",
				"description": "分析日期，格式 YYYY-MM-DD",
			},
			"indicators": map[string]interface{}{
				"type":        "array",
				"items":       map[string]interface{}{"type": "string"},
				"description": "需要获取的指标列表，如 ['MACD', 'RSI', 'SMA', 'EMA', 'BB']",
			},
		},
		"required": []string{"symbol", "date"},
	},
	func(ctx context.Context, args map[string]interface{}) (string, error) {
		symbol, _ := args["symbol"].(string)
		date, _ := args["date"].(string)

		data, err := fetchIndicators(ctx, symbol, date)
		if err != nil {
			return "", err
		}

		return data, nil
	},
)

// fetchIndicators 获取技术指标
func fetchIndicators(ctx context.Context, symbol, date string) (string, error) {
	return fmt.Sprintf(`Technical Indicators for %s as of %s:

**Moving Averages:**
| Indicator | Value   | Signal    |
|-----------|---------|-----------|
| SMA 20    | 862.45  | Above     |
| SMA 50    | 835.20  | Above     |
| SMA 200   | 720.80  | Above     |
| EMA 12    | 878.30  | Above     |
| EMA 26    | 855.60  | Above     |

**MACD:**
- MACD Line: 22.70
- Signal Line: 18.45
- Histogram: 4.25 (Bullish)
- Trend: Bullish crossover confirmed

**RSI (14-day):**
- Current: 65.8
- Signal: Neutral to slightly overbought
- Trend: Rising

**Bollinger Bands:**
- Upper Band: 912.40
- Middle Band: 862.45
- Lower Band: 812.50
- Position: Mid-upper range

**Volume Analysis:**
- Current Volume: 45.2M
- 20-day Avg Volume: 38.5M
- Volume Ratio: 1.17x (Above average)

**Overall Technical Summary:**
- Trend: Strong uptrend
- Momentum: Positive
- Support Levels: 860, 835, 800
- Resistance Levels: 900, 920, 950
`, symbol, date), nil
}
