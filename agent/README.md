# TradingAgents - Golang + Eino Edition

> 使用 CloudWeGo Eino 框架和 Golang 实现的多智能体金融交易框架

## 📖 简介

这是 [TradingAgents](https://github.com/TauricResearch/TradingAgents) 项目的 Golang 版本，使用字节跳动开源的 [Eino](https://github.com/cloudwego/eino) 框架实现多智能体协作。

## 🏗️ 架构

```
┌─────────────────────────────────────────────────────────────┐
│                     TradingAgentsGraph                       │
├─────────────────────────────────────────────────────────────┤
│  Phase 1: Analyst Team                                       │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────────┐│
│  │ Market   │ │ Social   │ │ News     │ │ Fundamentals     ││
│  │ Analyst  │ │ Analyst  │ │ Analyst  │ │ Analyst          ││
│  └──────────┘ └──────────┘ └──────────┘ └──────────────────┘│
├─────────────────────────────────────────────────────────────┤
│  Phase 2: Investment Debate                                  │
│  ┌──────────────┐          ┌──────────────┐                 │
│  │ Bull         │◄────────►│ Bear         │                 │
│  │ Researcher   │  Debate  │ Researcher   │                 │
│  └──────────────┘          └──────────────┘                 │
│              │                    │                          │
│              └────────┬───────────┘                          │
│                       ▼                                      │
│              ┌──────────────┐                                │
│              │ Research     │                                │
│              │ Manager      │                                │
│              └──────────────┘                                │
├─────────────────────────────────────────────────────────────┤
│  Phase 3: Trader Decision                                    │
│              ┌──────────────┐                                │
│              │   Trader     │                                │
│              └──────────────┘                                │
├─────────────────────────────────────────────────────────────┤
│  Phase 4: Risk Management                                    │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐                     │
│  │ Risky    │ │ Neutral  │ │ Safe     │                     │
│  │ Analyst  │ │ Analyst  │ │ Analyst  │                     │
│  └──────────┘ └──────────┘ └──────────┘                     │
│              │       │       │                               │
│              └───────┼───────┘                               │
│                      ▼                                       │
│              ┌──────────────┐                                │
│              │ Risk Manager │ ──► FINAL DECISION             │
│              └──────────────┘                                │
└─────────────────────────────────────────────────────────────┘
```

## 🚀 快速开始

### 环境要求

- Go 1.21+
- OpenAI API Key

### 安装

```bash
cd agent
go mod tidy
```

### 设置环境变量

```bash
export OPENAI_API_KEY=your_openai_api_key
export ALPHA_VANTAGE_API_KEY=your_alpha_vantage_api_key  # 可选
```

### 运行

```bash
# 运行交易分析
go run cmd/main.go run --ticker NVDA --date 2024-05-10

# 启用调试模式
go run cmd/main.go run --ticker NVDA --date 2024-05-10 --debug

# 查看帮助
go run cmd/main.go --help
```

## 📁 项目结构

```
agent/
├── cmd/
│   └── main.go                    # CLI 入口
├── internal/
│   ├── agents/                    # Agent 实现
│   │   ├── analysts/              # 分析师 (Market, Social, News, Fundamentals)
│   │   ├── researchers/           # 研究员 (Bull, Bear)
│   │   ├── managers/              # 管理者 (Research, Risk)
│   │   ├── risk_mgmt/             # 风险管理 (Aggressive, Neutral, Conservative)
│   │   └── trader/                # 交易员
│   ├── config/                    # 配置管理
│   ├── dataflows/                 # 数据源适配
│   ├── graph/                     # 图编排
│   ├── llm/                       # LLM 封装
│   ├── memory/                    # 向量记忆
│   ├── states/                    # 状态定义
│   └── tools/                     # 数据工具
├── go.mod
└── README.md
```

## ⚙️ 配置

### 默认配置

```yaml
llm_provider: openai
deep_think_llm: o4-mini
quick_think_llm: gpt-4o-mini
backend_url: https://api.openai.com/v1

max_debate_rounds: 1
max_risk_discuss_rounds: 1

data_vendors:
  core_stock_apis: yfinance
  technical_indicators: yfinance
  fundamental_data: alpha_vantage
  news_data: alpha_vantage
```

### 自定义配置

创建 `config.yaml` 文件：

```yaml
llm_provider: openai
deep_think_llm: gpt-4o
quick_think_llm: gpt-4o-mini
max_debate_rounds: 2
```

运行时指定配置文件：

```bash
go run cmd/main.go run --config config.yaml --ticker AAPL
```

## 🔧 API 使用

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/tradingagents/agent/internal/config"
    "github.com/tradingagents/agent/internal/graph"
)

func main() {
    ctx := context.Background()
    cfg := config.DefaultConfig()
    
    // 创建交易图
    tradingGraph, err := graph.NewTradingAgentsGraph(ctx, cfg, nil)
    if err != nil {
        panic(err)
    }
    
    // 执行分析
    state, decision, err := tradingGraph.Propagate(ctx, "NVDA", "2024-05-10")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Decision: %s\n", decision)
    
    // 反思学习（可选）
    tradingGraph.ReflectAndRemember(ctx, state, 10.5) // 10.5% 收益
}
```

## 📊 智能体说明

### 分析师团队
- **Market Analyst**: 技术分析，OHLCV 数据，技术指标
- **Social Media Analyst**: 社交媒体情绪分析
- **News Analyst**: 新闻和宏观经济分析
- **Fundamentals Analyst**: 公司基本面分析

### 研究员团队
- **Bull Researcher**: 看涨论点，强调增长潜力
- **Bear Researcher**: 看跌论点，强调风险因素
- **Research Manager**: 综合评估，制定投资计划

### 风险管理团队
- **Risky Analyst**: 激进派，强调高风险高回报
- **Neutral Analyst**: 中立派，平衡分析
- **Safe Analyst**: 保守派，强调风险规避
- **Risk Manager**: 最终裁决，输出 BUY/SELL/HOLD

## 📚 参考

- [原版 TradingAgents (Python)](https://github.com/TauricResearch/TradingAgents)
- [Eino Framework](https://github.com/cloudwego/eino)
- [论文: TradingAgents](https://arxiv.org/abs/2412.20138)

## 📄 License

Apache License 2.0
