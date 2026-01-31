# FinRobot Golang + Eino 迁移开发 ToDo

> 使用 CloudWeGo Eino 框架，将 FinRobot 金融 AI Agent 平台完整复制到 Golang

## 📊 项目进度总览

```
总进度: [░░░░░░░░░░░░░░░░░░░░] 0%

Phase 1 基础框架: [░░░░░░░░░░] 0%
Phase 2 数据源层: [░░░░░░░░░░] 0%
Phase 3 工作流层: [░░░░░░░░░░] 0%
Phase 4 功能模块: [░░░░░░░░░░] 0%
Phase 5 集成测试: [░░░░░░░░░░] 0%
```

---

## 🎯 Phase 1: 基础框架搭建

### 1.1 项目初始化
- [ ] 创建 `agent/finrobot` 目录结构
- [ ] 初始化 `go.mod`
- [ ] 配置基础依赖
  - [ ] `github.com/cloudwego/eino`
  - [ ] `github.com/cloudwego/eino-ext`
  - [ ] `github.com/spf13/cobra` (CLI)
  - [ ] `github.com/spf13/viper` (配置)

### 1.2 配置管理
- [ ] 实现 `internal/config/config.go`
  - [ ] LLM 配置 (model, api_key, base_url)
  - [ ] API Keys 管理 (FinnHub, FMP, SEC)
  - [ ] 工作目录配置

### 1.3 LLM 客户端封装
- [ ] 实现 `internal/llm/client.go`
  - [ ] 统一 LLM 接口定义
  - [ ] OpenAI 适配器
  - [ ] 错误处理和重试机制

### 1.4 核心 Agent 基类
- [ ] 实现 `internal/agents/finrobot.go`
  - [ ] FinRobot Agent 结构体
  - [ ] 系统消息配置
  - [ ] 工具注册机制
  - [ ] Agent Library 集成

### 1.5 Agent 库定义
- [ ] 实现 `internal/agents/library.go`
  - [ ] Market_Analyst 配置
  - [ ] Expert_Investor 配置
  - [ ] Financial_Analyst 配置
  - [ ] Data_Analyst 配置
  - [ ] Software_Developer 配置

### 1.6 提示词模板
- [ ] 实现 `internal/agents/prompts.go`
  - [ ] Role System Message 模板
  - [ ] Leader System Message 模板
  - [ ] 各角色专业提示词

---

## 🔌 Phase 2: 数据源层实现

### 2.1 数据源接口定义
- [ ] 实现 `internal/datasource/interface.go`
  - [ ] CompanyDataSource 接口
  - [ ] FinancialDataSource 接口
  - [ ] NewsDataSource 接口

### 2.2 FinnHub API 集成
- [ ] 实现 `internal/datasource/finnhub.go`
  - [ ] `GetCompanyProfile(symbol string) (*CompanyProfile, error)`
  - [ ] `GetCompanyNews(symbol, startDate, endDate string, maxNum int) ([]News, error)`
  - [ ] `GetBasicFinancials(symbol string) (*Financials, error)`
  - [ ] `GetBasicFinancialsHistory(symbol, freq, startDate, endDate string) (*FinancialsHistory, error)`
  - [ ] 单元测试

### 2.3 YFinance API 集成
- [ ] 实现 `internal/datasource/yfinance.go`
  - [ ] `GetStockData(symbol, startDate, endDate string) (*StockData, error)`
  - [ ] `GetStockInfo(symbol string) (*StockInfo, error)`
  - [ ] `GetCompanyInfo(symbol string) (*CompanyInfo, error)`
  - [ ] `GetIncomeStmt(symbol string) (*IncomeStatement, error)`
  - [ ] `GetBalanceSheet(symbol string) (*BalanceSheet, error)`
  - [ ] `GetCashFlow(symbol string) (*CashFlow, error)`
  - [ ] `GetAnalystRecommendations(symbol string) (*Recommendations, error)`
  - [ ] HTTP 客户端封装 (调用 Yahoo Finance API)
  - [ ] 单元测试

### 2.4 FMP API 集成
- [ ] 实现 `internal/datasource/fmp.go`
  - [ ] `GetSECReport(symbol, year string) (*SECReport, error)`
  - [ ] `GetFinancialMetrics(symbol string, years int) (*FinancialMetrics, error)`
  - [ ] `GetCompetitorMetrics(symbol string, competitors []string) (*CompetitorMetrics, error)`
  - [ ] `GetHistoricalMarketCap(symbol, date string) (float64, error)`
  - [ ] `GetHistoricalBVPS(symbol, date string) (float64, error)`
  - [ ] `GetTargetPrice(symbol, date string) (float64, error)`
  - [ ] 单元测试

### 2.5 SEC Filings 集成
- [ ] 实现 `internal/datasource/sec.go`
  - [ ] `Get10KSection(symbol, year string, section interface{}) (string, error)`
  - [ ] Section 映射 (1, 1A, 7, 7A 等)
  - [ ] HTML 解析和清洗
  - [ ] 单元测试

### 2.6 Reddit 数据 (可选)
- [ ] 实现 `internal/datasource/reddit.go`
  - [ ] 子版块搜索
  - [ ] 情绪分析数据

---

## 🔄 Phase 3: 工作流层实现

### 3.1 工具注册机制
- [ ] 实现 `internal/tools/registry.go`
  - [ ] 工具注册中心
  - [ ] InferTool 自动推断
  - [ ] 工具分组管理

### 3.2 数据源工具注册
- [ ] 实现 `internal/tools/datasource_tools.go`
  - [ ] FinnHub 工具组
  - [ ] YFinance 工具组
  - [ ] FMP 工具组
  - [ ] SEC 工具组

### 3.3 SingleAssistant 工作流
- [ ] 实现 `internal/agents/workflow/single.go`
  - [ ] 基础对话循环
  - [ ] 工具调用处理
  - [ ] 终止条件判断
  - [ ] 历史记录管理

### 3.4 SingleAssistantRAG 工作流
- [ ] 实现 `internal/agents/workflow/single_rag.go`
  - [ ] 向量数据库集成 (Qdrant/Milvus)
  - [ ] 文档检索功能
  - [ ] 上下文注入

### 3.5 SingleAssistantShadow 工作流
- [ ] 实现 `internal/agents/workflow/single_shadow.go`
  - [ ] 影子 Agent 创建
  - [ ] 嵌套对话处理
  - [ ] 指令解析和转发

### 3.6 MultiAssistant 工作流
- [ ] 实现 `internal/agents/workflow/multi.go`
  - [ ] 群聊管理器
  - [ ] Speaker 选择策略
  - [ ] 消息路由

### 3.7 MultiAssistantWithLeader 工作流
- [ ] 实现 `internal/agents/workflow/multi_leader.go`
  - [ ] Leader Agent 定义
  - [ ] 任务分配机制
  - [ ] 嵌套聊天触发

---

## 📊 Phase 4: 功能模块实现

### 4.1 报表分析工具
- [ ] 实现 `internal/functional/analyzer.go`
  - [ ] `AnalyzeIncomeStmt(symbol, fyear, savePath string) error`
  - [ ] `AnalyzeBalanceSheet(symbol, fyear, savePath string) error`
  - [ ] `AnalyzeCashFlow(symbol, fyear, savePath string) error`
  - [ ] `AnalyzeSegmentStmt(symbol, fyear, savePath string) error`
  - [ ] `IncomeSummarization(symbol, fyear, incomeAnalysis, segmentAnalysis, savePath string) error`
  - [ ] `GetRiskAssessment(symbol, fyear, savePath string) error`
  - [ ] `GetCompetitorsAnalysis(symbol string, competitors []string, fyear, savePath string) error`
  - [ ] `AnalyzeBusinessHighlights(symbol, fyear, savePath string) error`
  - [ ] `AnalyzeCompanyDescription(symbol, fyear, savePath string) error`
  - [ ] `GetKeyData(symbol, filingDate string) (*KeyData, error)`

### 4.2 图表生成
- [ ] 实现 `internal/functional/charting.go`
  - [ ] 使用 go-echarts 库
  - [ ] `PlotSharePerformance(symbol, startDate, endDate, savePath string) error`
  - [ ] `PlotPEEPS(symbol, savePath string) error`
  - [ ] `PlotFinancialMetrics(metrics *FinancialMetrics, savePath string) error`

### 4.3 PDF 报告生成
- [ ] 实现 `internal/functional/reportlab.go`
  - [ ] 使用 gofpdf 或 maroto 库
  - [ ] `BuildAnnualReport(params *ReportParams) error`
    - [ ] 封面设计
    - [ ] 双栏布局
    - [ ] 图表嵌入
    - [ ] 表格格式化

### 4.4 文本处理
- [ ] 实现 `internal/functional/text.go`
  - [ ] `CheckTextLength(text string, minLen, maxLen int) bool`
  - [ ] `SummarizeText(text string, maxLen int) string`
  - [ ] `ExtractKeyPoints(text string) []string`

### 4.5 RAG 功能
- [ ] 实现 `internal/functional/rag.go`
  - [ ] 文档加载器
  - [ ] 文本分块器
  - [ ] 向量存储
  - [ ] 检索查询

---

## 🧪 Phase 5: 集成测试与 CLI

### 5.1 CLI 实现
- [ ] 实现 `cmd/main.go`
  - [ ] `forecast` 命令 - 市场预测
  - [ ] `report` 命令 - 研究报告生成
  - [ ] `analyze` 命令 - 财务分析
  - [ ] 全局参数配置

### 5.2 Demo 1: Market Forecaster
- [ ] 复现 `agent_fingpt_forecaster` 功能
  - [ ] 输入: ticker, date
  - [ ] 输出: 股价预测和分析

### 5.3 Demo 2: Equity Research Report
- [ ] 复现 `agent_annual_report` 功能
  - [ ] 输入: company, fyear
  - [ ] 输出: PDF 研究报告

### 5.4 Demo 3: Trade Strategist (可选)
- [ ] 复现交易策略功能
  - [ ] 多模态分析
  - [ ] 策略建议

### 5.5 文档编写
- [ ] README.md - 快速开始
- [ ] API 文档
- [ ] 示例代码

---

## 📋 详细 API 对照表

### FinnHubUtils

| Python 方法 | Golang 方法 | 状态 |
|-------------|-------------|------|
| `get_company_profile(symbol)` | `GetCompanyProfile(symbol)` | [ ] |
| `get_company_news(symbol, start, end, max)` | `GetCompanyNews(symbol, start, end, max)` | [ ] |
| `get_basic_financials(symbol)` | `GetBasicFinancials(symbol)` | [ ] |
| `get_basic_financials_history(symbol, freq, start, end)` | `GetBasicFinancialsHistory(symbol, freq, start, end)` | [ ] |

### YFinanceUtils

| Python 方法 | Golang 方法 | 状态 |
|-------------|-------------|------|
| `get_stock_data(symbol, start, end)` | `GetStockData(symbol, start, end)` | [ ] |
| `get_stock_info(symbol)` | `GetStockInfo(symbol)` | [ ] |
| `get_company_info(symbol)` | `GetCompanyInfo(symbol)` | [ ] |
| `get_income_stmt(symbol)` | `GetIncomeStmt(symbol)` | [ ] |
| `get_balance_sheet(symbol)` | `GetBalanceSheet(symbol)` | [ ] |
| `get_cash_flow(symbol)` | `GetCashFlow(symbol)` | [ ] |
| `get_analyst_recommendations(symbol)` | `GetAnalystRecommendations(symbol)` | [ ] |

### FMPUtils

| Python 方法 | Golang 方法 | 状态 |
|-------------|-------------|------|
| `get_sec_report(symbol, year)` | `GetSECReport(symbol, year)` | [ ] |
| `get_financial_metrics(symbol, years)` | `GetFinancialMetrics(symbol, years)` | [ ] |
| `get_competitor_financial_metrics(symbol, competitors, years)` | `GetCompetitorMetrics(symbol, competitors, years)` | [ ] |
| `get_historical_market_cap(symbol, date)` | `GetHistoricalMarketCap(symbol, date)` | [ ] |
| `get_historical_bvps(symbol, date)` | `GetHistoricalBVPS(symbol, date)` | [ ] |
| `get_target_price(symbol, date)` | `GetTargetPrice(symbol, date)` | [ ] |

### SECUtils

| Python 方法 | Golang 方法 | 状态 |
|-------------|-------------|------|
| `get_10k_section(symbol, year, section)` | `Get10KSection(symbol, year, section)` | [ ] |

### ReportAnalysisUtils

| Python 方法 | Golang 方法 | 状态 |
|-------------|-------------|------|
| `analyze_income_stmt(symbol, fyear, path)` | `AnalyzeIncomeStmt(symbol, fyear, path)` | [ ] |
| `analyze_balance_sheet(symbol, fyear, path)` | `AnalyzeBalanceSheet(symbol, fyear, path)` | [ ] |
| `analyze_cash_flow(symbol, fyear, path)` | `AnalyzeCashFlow(symbol, fyear, path)` | [ ] |
| `analyze_segment_stmt(symbol, fyear, path)` | `AnalyzeSegmentStmt(symbol, fyear, path)` | [ ] |
| `income_summarization(...)` | `IncomeSummarization(...)` | [ ] |
| `get_risk_assessment(symbol, fyear, path)` | `GetRiskAssessment(symbol, fyear, path)` | [ ] |
| `get_competitors_analysis(...)` | `GetCompetitorsAnalysis(...)` | [ ] |
| `analyze_business_highlights(symbol, fyear, path)` | `AnalyzeBusinessHighlights(symbol, fyear, path)` | [ ] |
| `analyze_company_description(symbol, fyear, path)` | `AnalyzeCompanyDescription(symbol, fyear, path)` | [ ] |
| `get_key_data(symbol, date)` | `GetKeyData(symbol, date)` | [ ] |

### ReportLabUtils

| Python 方法 | Golang 方法 | 状态 |
|-------------|-------------|------|
| `build_annual_report(...)` | `BuildAnnualReport(...)` | [ ] |

---

## 🎯 下一步行动

1. **立即开始**: 创建目录结构和 go.mod
2. **第一个里程碑**: 完成 FinnHub API 集成 + SingleAssistant 工作流
3. **第二个里程碑**: 完成 Market Forecaster Demo
4. **最终里程碑**: 完成 Equity Research Report Demo

---

## 📝 开发笔记

### 技术决策记录
| 日期 | 决策 | 原因 |
|------|------|------|
| 2026-01-31 | 使用 Eino ChatModelAgent | 官方推荐，功能完整 |
| 2026-01-31 | YFinance 使用 HTTP 封装 | 无原生 Golang 库 |
| 2026-01-31 | PDF 选用 gofpdf | 功能强大，社区活跃 |

### 遇到的问题
_待记录_

### 解决方案
_待记录_

---

*文档创建: 2026-01-31*
*最后更新: 2026-01-31*
