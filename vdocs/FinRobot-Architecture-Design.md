# FinRobot Golang 架构设计文档

> 基于 Eino 框架的金融 AI Agent 平台架构设计

---

## 1. 系统架构总览

```mermaid
graph TB
    subgraph "**用户层**"
        CLI[**CLI 命令行**]
        API[**API 接口**]
    end
    
    subgraph "**Agent 编排层**"
        SA[**SingleAssistant**]
        SAR[**SingleAssistantRAG**]
        SAS[**SingleAssistantShadow**]
        MA[**MultiAssistant**]
        MAL[**MultiAssistantWithLeader**]
    end
    
    subgraph "**核心 Agent**"
        FR[**FinRobot Agent**]
        LIB[**Agent Library**]
    end
    
    subgraph "**工具层**"
        TR[**Tool Registry**]
        DT[**DataSource Tools**]
        AT[**Analysis Tools**]
    end
    
    subgraph "**数据源层**"
        FH[**FinnHub**]
        YF[**YFinance**]
        FMP[**FMP API**]
        SEC[**SEC Filings**]
    end
    
    subgraph "**功能模块**"
        RA[**报表分析**]
        CH[**图表生成**]
        PDF[**PDF 报告**]
        RAG[**RAG 检索**]
    end
    
    subgraph "**LLM 层**"
        LLM[**LLM Client**]
        OAI[**OpenAI**]
    end
    
    CLI --> SA
    CLI --> MA
    API --> SA
    API --> MA
    
    SA --> FR
    SAR --> FR
    SAS --> FR
    MA --> FR
    MAL --> FR
    
    FR --> LIB
    FR --> TR
    
    TR --> DT
    TR --> AT
    
    DT --> FH
    DT --> YF
    DT --> FMP
    DT --> SEC
    
    AT --> RA
    AT --> CH
    AT --> PDF
    AT --> RAG
    
    FR --> LLM
    LLM --> OAI
    
    style CLI fill:#ffe1e1,stroke:#333,stroke-width:2px,color:#000
    style API fill:#ffe1e1,stroke:#333,stroke-width:2px,color:#000
    style FR fill:#e1ffe1,stroke:#333,stroke-width:2px,color:#000
    style LLM fill:#e1f5ff,stroke:#333,stroke-width:2px,color:#000
    style TR fill:#fff3e1,stroke:#333,stroke-width:2px,color:#000
```

---

## 2. Agent 工作流详解

### 2.1 SingleAssistant 工作流

```mermaid
sequenceDiagram
    participant U as **用户**
    participant SA as **SingleAssistant**
    participant FR as **FinRobot Agent**
    participant LLM as **LLM**
    participant T as **Tools**
    
    U->>SA: **发送消息**
    SA->>FR: **初始化 Agent**
    
    rect rgb(240, 248, 255)
    Note over SA,T: **对话循环**
    
    loop 直到终止条件
        FR->>LLM: **请求响应**
        LLM-->>FR: **返回响应/工具调用**
        
        alt **需要工具调用**
            FR->>T: **执行工具**
            T-->>FR: **工具结果**
            FR->>LLM: **提交工具结果**
        else **直接回复**
            FR-->>SA: **生成回复**
        end
    end
    end
    
    SA-->>U: **最终结果**
```

### 2.2 MultiAssistant 群聊工作流

```mermaid
sequenceDiagram
    participant U as **用户**
    participant GM as **GroupChat Manager**
    participant A1 as **Agent 1**
    participant A2 as **Agent 2**
    participant A3 as **Agent 3**
    
    U->>GM: **发起群聊**
    GM->>GM: **选择发言者**
    
    rect rgb(255, 250, 205)
    Note over GM,A3: **群聊轮次**
    
    loop 直到任务完成
        GM->>A1: **轮到你发言**
        A1-->>GM: **回复消息**
        GM->>GM: **选择下一个发言者**
        GM->>A2: **轮到你发言**
        A2-->>GM: **回复消息**
        
        alt **需要工具执行**
            GM->>U: **请求执行工具**
            U-->>GM: **执行结果**
        end
    end
    end
    
    GM-->>U: **群聊总结**
```

### 2.3 MultiAssistantWithLeader 工作流

```mermaid
sequenceDiagram
    participant U as **用户**
    participant L as **Leader Agent**
    participant E1 as **Employee 1**
    participant E2 as **Employee 2**
    
    U->>L: **分配任务**
    L->>L: **分析任务**
    
    rect rgb(230, 255, 230)
    Note over L,E2: **任务分配与执行**
    
    L->>E1: **[Employee1] 子任务 1**
    activate E1
    E1->>E1: **执行任务**
    E1-->>L: **任务 1 结果**
    deactivate E1
    
    L->>E2: **[Employee2] 子任务 2**
    activate E2
    E2->>E2: **执行任务**
    E2-->>L: **任务 2 结果**
    deactivate E2
    end
    
    L->>L: **综合分析结果**
    L-->>U: **最终报告**
```

---

## 3. 数据流架构

```mermaid
graph LR
    subgraph "**输入**"
        I1[**Ticker Symbol**]
        I2[**Date Range**]
        I3[**Fiscal Year**]
    end
    
    subgraph "**数据采集**"
        D1[**公司基本信息**]
        D2[**股价数据**]
        D3[**财务报表**]
        D4[**新闻舆情**]
        D5[**10-K 报告**]
    end
    
    subgraph "**数据分析**"
        A1[**财务分析**]
        A2[**风险评估**]
        A3[**竞对分析**]
        A4[**市场预测**]
    end
    
    subgraph "**输出**"
        O1[**分析报告**]
        O2[**PDF 文档**]
        O3[**预测结果**]
    end
    
    I1 --> D1
    I1 --> D2
    I1 --> D3
    I2 --> D2
    I2 --> D4
    I3 --> D5
    
    D1 --> A1
    D2 --> A4
    D3 --> A1
    D4 --> A4
    D5 --> A2
    
    A1 --> O1
    A2 --> O1
    A3 --> O1
    A1 --> O2
    A2 --> O2
    A3 --> O2
    A4 --> O3
    
    style I1 fill:#ffe1e1,stroke:#333,stroke-width:2px,color:#000
    style I2 fill:#ffe1e1,stroke:#333,stroke-width:2px,color:#000
    style I3 fill:#ffe1e1,stroke:#333,stroke-width:2px,color:#000
    style O1 fill:#e1ffe1,stroke:#333,stroke-width:2px,color:#000
    style O2 fill:#e1ffe1,stroke:#333,stroke-width:2px,color:#000
    style O3 fill:#e1ffe1,stroke:#333,stroke-width:2px,color:#000
```

---

## 4. 类图设计

### 4.1 Agent 层类图

```mermaid
classDiagram
    class IAgent {
        <<interface>>
        +Chat(message string) string
        +Reset()
        +RegisterTools(tools []Tool)
    }
    
    class FinRobotAgent {
        -name string
        -systemMessage string
        -llmConfig LLMConfig
        -toolkits []Tool
        +NewFinRobotAgent(config AgentConfig) *FinRobotAgent
        +Chat(message string) string
        +Reset()
        +RegisterTools(tools []Tool)
    }
    
    class AgentLibrary {
        -agents map[string]AgentConfig
        +GetAgent(name string) AgentConfig
        +RegisterAgent(config AgentConfig)
    }
    
    class AgentConfig {
        +Name string
        +Profile string
        +Description string
        +Toolkits []string
    }
    
    IAgent <|.. FinRobotAgent
    FinRobotAgent --> AgentLibrary
    AgentLibrary --> AgentConfig
```

### 4.2 数据源层类图

```mermaid
classDiagram
    class IDataSource {
        <<interface>>
        +GetCompanyProfile(symbol string) CompanyProfile
        +GetStockData(symbol, start, end string) StockData
    }
    
    class FinnHubClient {
        -apiKey string
        -httpClient *http.Client
        +GetCompanyProfile(symbol string) CompanyProfile
        +GetCompanyNews(symbol, start, end string) []News
        +GetBasicFinancials(symbol string) Financials
    }
    
    class YFinanceClient {
        -httpClient *http.Client
        +GetStockData(symbol, start, end string) StockData
        +GetIncomeStmt(symbol string) IncomeStatement
        +GetBalanceSheet(symbol string) BalanceSheet
        +GetCashFlow(symbol string) CashFlow
    }
    
    class FMPClient {
        -apiKey string
        -httpClient *http.Client
        +GetSECReport(symbol, year string) SECReport
        +GetFinancialMetrics(symbol string, years int) FinancialMetrics
    }
    
    class SECClient {
        -httpClient *http.Client
        +Get10KSection(symbol, year string, section int) string
    }
    
    IDataSource <|.. FinnHubClient
    IDataSource <|.. YFinanceClient
    IDataSource <|.. FMPClient
    IDataSource <|.. SECClient
```

---

## 5. 工具注册机制

```mermaid
graph TB
    subgraph "**Eino Tool 定义**"
        IT[**InferTool**<br/>自动推断参数]
        NT[**NewTool**<br/>手动定义参数]
    end
    
    subgraph "**工具注册中心**"
        TR[**ToolRegistry**]
        TG1[**DataSource Group**]
        TG2[**Analysis Group**]
        TG3[**Reporting Group**]
    end
    
    subgraph "**具体工具**"
        T1[**GetCompanyProfile**]
        T2[**GetStockData**]
        T3[**AnalyzeIncomeStmt**]
        T4[**BuildReport**]
    end
    
    IT --> TR
    NT --> TR
    
    TR --> TG1
    TR --> TG2
    TR --> TG3
    
    TG1 --> T1
    TG1 --> T2
    TG2 --> T3
    TG3 --> T4
    
    style TR fill:#fff3e1,stroke:#333,stroke-width:2px,color:#000
    style IT fill:#e1f5ff,stroke:#333,stroke-width:2px,color:#000
    style NT fill:#e1f5ff,stroke:#333,stroke-width:2px,color:#000
```

### 5.1 InferTool 示例

```go
// 使用 struct tag 定义工具参数
type GetCompanyProfileParams struct {
    Symbol string `json:"symbol" jsonschema:"description=股票代码,example=AAPL"`
}

// 注册工具
tool, _ := utils.InferTool(
    "get_company_profile",
    "获取公司基本信息",
    func(ctx context.Context, params *GetCompanyProfileParams) (string, error) {
        return finnhub.GetCompanyProfile(params.Symbol)
    },
)
```

---

## 6. Equity Research Report 生成流程

```mermaid
graph TB
    subgraph "**Phase 1: 数据收集**"
        P1A[**获取 10-K 报告**]
        P1B[**获取财务报表**]
        P1C[**获取市场数据**]
        P1D[**获取竞对数据**]
    end
    
    subgraph "**Phase 2: 分析**"
        P2A[**收入分析**]
        P2B[**资产负债分析**]
        P2C[**现金流分析**]
        P2D[**分部分析**]
        P2E[**风险评估**]
        P2F[**竞对分析**]
    end
    
    subgraph "**Phase 3: 综合**"
        P3A[**业务概述**]
        P3B[**市场定位**]
        P3C[**运营结果**]
    end
    
    subgraph "**Phase 4: 可视化**"
        P4A[**股价表现图**]
        P4B[**PE/EPS 图**]
        P4C[**财务指标表**]
    end
    
    subgraph "**Phase 5: 输出**"
        P5[**PDF 研究报告**]
    end
    
    P1A --> P2A
    P1A --> P2B
    P1A --> P2C
    P1A --> P2E
    P1B --> P2A
    P1B --> P2B
    P1B --> P2C
    P1B --> P2D
    P1C --> P4A
    P1C --> P4B
    P1D --> P2F
    
    P2A --> P3C
    P2B --> P3C
    P2C --> P3C
    P2D --> P3C
    P2E --> P5
    P2F --> P5
    
    P3A --> P5
    P3B --> P5
    P3C --> P5
    P4A --> P5
    P4B --> P5
    P4C --> P5
    
    style P5 fill:#e1ffe1,stroke:#333,stroke-width:3px,color:#000
```

---

## 7. 目录结构

```
TradingAgents/agent/finrobot/
│
├── cmd/
│   └── main.go                        # CLI 入口
│
├── internal/
│   ├── agents/
│   │   ├── finrobot.go                # 核心 Agent
│   │   ├── library.go                 # Agent 库
│   │   ├── prompts.go                 # 提示词模板
│   │   └── workflow/
│   │       ├── single.go              # SingleAssistant
│   │       ├── single_rag.go          # SingleAssistantRAG
│   │       ├── single_shadow.go       # SingleAssistantShadow
│   │       ├── multi.go               # MultiAssistant
│   │       └── multi_leader.go        # MultiAssistantWithLeader
│   │
│   ├── datasource/
│   │   ├── interface.go               # 数据源接口
│   │   ├── finnhub.go                 # FinnHub API
│   │   ├── yfinance.go                # YFinance API
│   │   ├── fmp.go                     # FMP API
│   │   └── sec.go                     # SEC Filings
│   │
│   ├── functional/
│   │   ├── analyzer.go                # 报表分析
│   │   ├── charting.go                # 图表生成
│   │   ├── reportlab.go               # PDF 生成
│   │   ├── text.go                    # 文本处理
│   │   └── rag.go                     # RAG 功能
│   │
│   ├── tools/
│   │   ├── registry.go                # 工具注册中心
│   │   ├── datasource_tools.go        # 数据源工具
│   │   └── analysis_tools.go          # 分析工具
│   │
│   ├── llm/
│   │   ├── client.go                  # LLM 客户端接口
│   │   └── openai.go                  # OpenAI 适配器
│   │
│   └── config/
│       └── config.go                  # 配置管理
│
├── pkg/
│   └── utils/
│       └── utils.go                   # 通用工具函数
│
├── go.mod
├── go.sum
└── README.md
```

---

## 8. 配置示例

```yaml
# config.yaml
llm:
  provider: openai
  model: gpt-4o
  api_key: ${OPENAI_API_KEY}
  base_url: https://api.openai.com/v1
  timeout: 120
  temperature: 0.5

api_keys:
  finnhub: ${FINNHUB_API_KEY}
  fmp: ${FMP_API_KEY}
  sec: ${SEC_API_KEY}

agent:
  default: Expert_Investor
  max_turns: 50
  human_input_mode: NEVER

output:
  work_dir: ./output
  report_dir: ./reports
```

---

## 9. 技术栈

| 组件 | 技术选型 | 说明 |
|------|----------|------|
| **Agent 框架** | CloudWeGo Eino | 字节开源 AI Agent 框架 |
| **LLM** | OpenAI GPT-4 | 主要 LLM 提供者 |
| **HTTP 客户端** | net/http + resty | API 调用 |
| **JSON 处理** | encoding/json | 数据序列化 |
| **图表生成** | go-echarts | 可视化图表 |
| **PDF 生成** | gofpdf / maroto | 报告输出 |
| **CLI** | cobra | 命令行界面 |
| **配置** | viper | 配置管理 |
| **日志** | zap / zerolog | 结构化日志 |

---

## 10. 关键设计决策

| 决策点 | 选择 | 理由 |
|--------|------|------|
| Agent 框架 | Eino | 官方支持，与 Golang 生态深度集成 |
| 工具定义方式 | InferTool | 减少样板代码，自动推断参数 |
| 数据源封装 | HTTP 封装 | 灵活性高，无需依赖第三方库 |
| 工作流编排 | Eino Workflow | 原生支持顺序/并行执行 |
| PDF 库 | gofpdf | 功能完整，社区活跃 |

---

*文档版本: v1.0*
*创建日期: 2026-01-31*
