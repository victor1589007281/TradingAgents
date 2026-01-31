// Package tools 基本面数据工具
package tools

import (
	"context"
	"fmt"
)

// GetFundamentalsTool 获取公司基本面工具
var GetFundamentalsTool = NewTool(
	"get_fundamentals",
	"获取公司的基本面信息，包括公司概况、财务摘要、估值指标等",
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
		return fetchFundamentals(ctx, symbol)
	},
)

// GetBalanceSheetTool 获取资产负债表工具
var GetBalanceSheetTool = NewTool(
	"get_balance_sheet",
	"获取公司的资产负债表数据",
	map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"symbol": map[string]interface{}{
				"type":        "string",
				"description": "股票代码",
			},
			"period": map[string]interface{}{
				"type":        "string",
				"description": "时间周期: annual 或 quarterly",
				"enum":        []string{"annual", "quarterly"},
			},
		},
		"required": []string{"symbol"},
	},
	func(ctx context.Context, args map[string]interface{}) (string, error) {
		symbol, _ := args["symbol"].(string)
		period, _ := args["period"].(string)
		if period == "" {
			period = "annual"
		}
		return fetchBalanceSheet(ctx, symbol, period)
	},
)

// GetCashflowTool 获取现金流量表工具
var GetCashflowTool = NewTool(
	"get_cashflow",
	"获取公司的现金流量表数据",
	map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"symbol": map[string]interface{}{
				"type":        "string",
				"description": "股票代码",
			},
			"period": map[string]interface{}{
				"type":        "string",
				"description": "时间周期: annual 或 quarterly",
			},
		},
		"required": []string{"symbol"},
	},
	func(ctx context.Context, args map[string]interface{}) (string, error) {
		symbol, _ := args["symbol"].(string)
		period, _ := args["period"].(string)
		if period == "" {
			period = "annual"
		}
		return fetchCashflow(ctx, symbol, period)
	},
)

// GetIncomeStatementTool 获取利润表工具
var GetIncomeStatementTool = NewTool(
	"get_income_statement",
	"获取公司的利润表数据",
	map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"symbol": map[string]interface{}{
				"type":        "string",
				"description": "股票代码",
			},
			"period": map[string]interface{}{
				"type":        "string",
				"description": "时间周期: annual 或 quarterly",
			},
		},
		"required": []string{"symbol"},
	},
	func(ctx context.Context, args map[string]interface{}) (string, error) {
		symbol, _ := args["symbol"].(string)
		period, _ := args["period"].(string)
		if period == "" {
			period = "annual"
		}
		return fetchIncomeStatement(ctx, symbol, period)
	},
)

func fetchFundamentals(ctx context.Context, symbol string) (string, error) {
	return fmt.Sprintf(`Company Fundamentals for %s:

**Company Profile:**
- Name: NVIDIA Corporation
- Sector: Technology
- Industry: Semiconductors
- Market Cap: $2.2 Trillion
- Employees: 29,600

**Valuation Metrics:**
| Metric           | Value   | Industry Avg |
|------------------|---------|--------------|
| P/E Ratio        | 65.2    | 28.5         |
| P/S Ratio        | 35.8    | 8.2          |
| P/B Ratio        | 42.5    | 5.8          |
| EV/EBITDA        | 55.3    | 15.2         |
| PEG Ratio        | 1.25    | 1.8          |

**Growth Metrics:**
- Revenue Growth (YoY): +122%%
- EPS Growth (YoY): +168%%
- Free Cash Flow Growth: +145%%

**Profitability:**
- Gross Margin: 76.2%%
- Operating Margin: 62.1%%
- Net Margin: 55.8%%
- ROE: 85.3%%
- ROA: 45.2%%

**Financial Health:**
- Current Ratio: 4.2
- Debt/Equity: 0.41
- Interest Coverage: 58.5x

**Summary:**
Extremely strong fundamentals with exceptional growth rates and 
profitability metrics. Premium valuation justified by market 
leadership in AI/GPU space.
`, symbol), nil
}

func fetchBalanceSheet(ctx context.Context, symbol, period string) (string, error) {
	return fmt.Sprintf(`Balance Sheet for %s (%s):

**Assets:**
| Item                    | Amount (B)  |
|-------------------------|-------------|
| Cash & Equivalents      | $26.0       |
| Short-term Investments  | $18.5       |
| Accounts Receivable     | $12.3       |
| Inventory               | $5.8        |
| Total Current Assets    | $65.2       |
| Property & Equipment    | $4.8        |
| Goodwill & Intangibles  | $5.2        |
| **Total Assets**        | **$82.5**   |

**Liabilities:**
| Item                    | Amount (B)  |
|-------------------------|-------------|
| Accounts Payable        | $2.8        |
| Short-term Debt         | $1.2        |
| Total Current Liab.     | $12.5       |
| Long-term Debt          | $8.5        |
| **Total Liabilities**   | **$24.8**   |

**Equity:**
| Item                    | Amount (B)  |
|-------------------------|-------------|
| Common Stock            | $0.05       |
| Retained Earnings       | $52.2       |
| **Total Equity**        | **$57.7**   |
`, symbol, period), nil
}

func fetchCashflow(ctx context.Context, symbol, period string) (string, error) {
	return fmt.Sprintf(`Cash Flow Statement for %s (%s):

**Operating Activities:**
| Item                        | Amount (B) |
|-----------------------------|------------|
| Net Income                  | $29.8      |
| Depreciation & Amortization | $1.5       |
| Stock-based Compensation    | $2.8       |
| Changes in Working Capital  | -$3.2      |
| **Operating Cash Flow**     | **$32.5**  |

**Investing Activities:**
| Item                        | Amount (B) |
|-----------------------------|------------|
| Capital Expenditures        | -$2.1      |
| Acquisitions                | -$0.8      |
| Investment Purchases        | -$15.2     |
| **Investing Cash Flow**     | **-$18.1** |

**Financing Activities:**
| Item                        | Amount (B) |
|-----------------------------|------------|
| Dividends Paid              | -$0.4      |
| Share Repurchases           | -$10.5     |
| Debt Repayment              | -$1.2      |
| **Financing Cash Flow**     | **-$12.1** |

**Free Cash Flow: $30.4B**
`, symbol, period), nil
}

func fetchIncomeStatement(ctx context.Context, symbol, period string) (string, error) {
	return fmt.Sprintf(`Income Statement for %s (%s):

| Item                    | Amount (B) | YoY Change |
|-------------------------|------------|------------|
| Revenue                 | $60.9      | +122%%     |
| Cost of Revenue         | $14.5      | +65%%      |
| **Gross Profit**        | **$46.4**  | +147%%     |
| Operating Expenses      | $8.5       | +42%%      |
| **Operating Income**    | **$37.9**  | +185%%     |
| Interest Income         | $1.8       | +95%%      |
| Interest Expense        | -$0.3      | -15%%      |
| **Pre-tax Income**      | **$39.4**  | +180%%     |
| Income Tax              | $9.6       | +165%%     |
| **Net Income**          | **$29.8**  | +185%%     |

**Per Share Data:**
- EPS (Basic): $12.05
- EPS (Diluted): $11.93
- Shares Outstanding: 2.47B
`, symbol, period), nil
}
