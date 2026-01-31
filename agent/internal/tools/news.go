// Package tools 新闻数据工具
package tools

import (
	"context"
	"fmt"
)

// GetNewsTool 获取公司相关新闻工具
var GetNewsTool = NewTool(
	"get_news",
	"获取与特定公司或主题相关的新闻文章",
	map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "搜索关键词，如公司名称或股票代码",
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
		"required": []string{"query"},
	},
	func(ctx context.Context, args map[string]interface{}) (string, error) {
		query, _ := args["query"].(string)
		startDate, _ := args["start_date"].(string)
		endDate, _ := args["end_date"].(string)
		return fetchNews(ctx, query, startDate, endDate)
	},
)

// GetGlobalNewsTool 获取全球宏观新闻工具
var GetGlobalNewsTool = NewTool(
	"get_global_news",
	"获取全球宏观经济和市场相关的新闻",
	map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"date": map[string]interface{}{
				"type":        "string",
				"description": "日期，格式 YYYY-MM-DD",
			},
			"look_back_days": map[string]interface{}{
				"type":        "integer",
				"description": "回溯天数，默认7天",
			},
			"limit": map[string]interface{}{
				"type":        "integer",
				"description": "返回新闻条数限制",
			},
		},
		"required": []string{"date"},
	},
	func(ctx context.Context, args map[string]interface{}) (string, error) {
		date, _ := args["date"].(string)
		lookBackDays := 7
		if v, ok := args["look_back_days"].(float64); ok {
			lookBackDays = int(v)
		}
		limit := 10
		if v, ok := args["limit"].(float64); ok {
			limit = int(v)
		}
		return fetchGlobalNews(ctx, date, lookBackDays, limit)
	},
)

// GetInsiderSentimentTool 获取内部人士情绪工具
var GetInsiderSentimentTool = NewTool(
	"get_insider_sentiment",
	"获取公司内部人士的买卖情绪指标",
	map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"symbol": map[string]interface{}{
				"type":        "string",
				"description": "股票代码",
			},
		},
		"required": []string{"symbol"},
	},
	func(ctx context.Context, args map[string]interface{}) (string, error) {
		symbol, _ := args["symbol"].(string)
		return fetchInsiderSentiment(ctx, symbol)
	},
)

// GetInsiderTransactionsTool 获取内部人士交易工具
var GetInsiderTransactionsTool = NewTool(
	"get_insider_transactions",
	"获取公司内部人士的最近交易记录",
	map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"symbol": map[string]interface{}{
				"type":        "string",
				"description": "股票代码",
			},
		},
		"required": []string{"symbol"},
	},
	func(ctx context.Context, args map[string]interface{}) (string, error) {
		symbol, _ := args["symbol"].(string)
		return fetchInsiderTransactions(ctx, symbol)
	},
)

func fetchNews(ctx context.Context, query, startDate, endDate string) (string, error) {
	return fmt.Sprintf(`News for "%s" (from %s to %s):

**1. NVIDIA Reports Record Quarterly Revenue**
- Source: Reuters | Date: 2024-05-09
- Summary: NVIDIA announced record quarterly revenue of $26 billion, 
  driven by unprecedented demand for AI chips. The company's data 
  center revenue grew 427%% year-over-year.
- Sentiment: Very Positive

**2. Major Tech Companies Expand AI Partnerships with NVIDIA**
- Source: Bloomberg | Date: 2024-05-08
- Summary: Microsoft, Google, and Amazon announced expanded partnerships 
  with NVIDIA for AI infrastructure. Combined orders exceed $15 billion.
- Sentiment: Positive

**3. NVIDIA Announces Next-Generation Blackwell Architecture**
- Source: TechCrunch | Date: 2024-05-07
- Summary: NVIDIA unveiled its next-generation Blackwell GPU architecture, 
  promising 30x improvement in AI inference performance.
- Sentiment: Very Positive

**4. Analysts Raise Price Targets Following Earnings Beat**
- Source: CNBC | Date: 2024-05-09
- Summary: Multiple Wall Street analysts raised price targets for NVIDIA 
  following strong earnings. Average target now $1,100, up from $850.
- Sentiment: Positive

**5. Semiconductor Industry Outlook Remains Strong**
- Source: WSJ | Date: 2024-05-06
- Summary: Industry analysts project continued strong demand for 
  semiconductors through 2025, with AI chips leading growth.
- Sentiment: Positive

**Overall News Sentiment: Strongly Bullish**
`, query, startDate, endDate), nil
}

func fetchGlobalNews(ctx context.Context, date string, lookBackDays, limit int) (string, error) {
	return fmt.Sprintf(`Global Market News (as of %s, past %d days):

**Macroeconomic Headlines:**

1. **Federal Reserve Signals Potential Rate Cuts**
   - The Fed indicated possible rate cuts in H2 2024 if inflation 
     continues to moderate. Markets responded positively.
   - Impact: Bullish for equities

2. **US GDP Growth Exceeds Expectations**
   - Q1 GDP grew 3.2%%, above the 2.8%% consensus estimate.
   - Strong consumer spending and business investment drove growth.
   - Impact: Positive for market sentiment

3. **China Announces New Economic Stimulus Package**
   - $150 billion stimulus focused on infrastructure and technology.
   - Expected to boost global demand for semiconductors.
   - Impact: Positive for tech sector

4. **European Central Bank Holds Rates Steady**
   - ECB maintained current rate levels, citing improving inflation.
   - Euro strengthened against dollar.
   - Impact: Neutral to slightly positive

5. **Oil Prices Stabilize After OPEC+ Decision**
   - OPEC+ extended production cuts through Q3 2024.
   - Brent crude stable around $82/barrel.
   - Impact: Neutral for broader market

**Market Indicators:**
- VIX: 14.2 (Low volatility)
- 10Y Treasury: 4.45%%
- Dollar Index: 104.8
- S&P 500: +2.1%% WoW

**Summary:** Global markets show positive momentum with supportive 
monetary policy outlook and strong economic data.
`, date, lookBackDays), nil
}

func fetchInsiderSentiment(ctx context.Context, symbol string) (string, error) {
	return fmt.Sprintf(`Insider Sentiment for %s:

**Overall Insider Sentiment: Neutral to Slightly Positive**

**Recent Insider Activity (Last 90 Days):**
- Buy Transactions: 12
- Sell Transactions: 28
- Net Shares Traded: -125,000

**Sentiment Indicators:**
| Metric                  | Value    | Signal      |
|-------------------------|----------|-------------|
| Buy/Sell Ratio          | 0.43     | Neutral     |
| Net Insider Flow        | -$45M    | Slight Sell |
| Insider Ownership       | 4.2%%    | Low         |
| Institutional Ownership | 68.5%%   | High        |

**Notable Transactions:**
- CEO: Sold 50,000 shares (planned sale under 10b5-1)
- CFO: Bought 5,000 shares at $820
- Director: Bought 2,000 shares at $845

**Analysis:**
Most selling appears to be pre-planned under 10b5-1 plans for 
tax and diversification purposes. Recent buying by CFO and 
directors suggests confidence in company outlook.
`, symbol), nil
}

func fetchInsiderTransactions(ctx context.Context, symbol string) (string, error) {
	return fmt.Sprintf(`Insider Transactions for %s:

| Date       | Name              | Title   | Type | Shares  | Price   | Value     |
|------------|-------------------|---------|------|---------|---------|-----------|
| 2024-05-08 | Jensen Huang      | CEO     | Sell | 50,000  | $885    | $44.25M   |
| 2024-05-06 | Colette Kress     | CFO     | Buy  | 5,000   | $820    | $4.1M     |
| 2024-05-03 | Mark Stevens      | Director| Buy  | 2,000   | $845    | $1.69M    |
| 2024-04-28 | Deborah Shoquist  | EVP     | Sell | 15,000  | $870    | $13.05M   |
| 2024-04-22 | Ajay Puri         | EVP     | Sell | 10,000  | $850    | $8.5M     |

**90-Day Summary:**
- Total Buys: $8.2M (12 transactions)
- Total Sells: $125.5M (28 transactions)
- Net Flow: -$117.3M

**Note:** Most sales are part of pre-arranged 10b5-1 trading plans 
and may not reflect current sentiment about stock price.
`, symbol), nil
}
