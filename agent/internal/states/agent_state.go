// Package states 定义 Agent 状态
package states

import (
	"github.com/cloudwego/eino/schema"
)

// InvestDebateState 投资辩论状态
type InvestDebateState struct {
	BullHistory     string `json:"bull_history"`      // 看涨研究员历史
	BearHistory     string `json:"bear_history"`      // 看跌研究员历史
	History         string `json:"history"`           // 完整辩论历史
	CurrentResponse string `json:"current_response"`  // 最新回复
	JudgeDecision   string `json:"judge_decision"`    // 裁判决定
	Count           int    `json:"count"`             // 辩论轮数
}

// RiskDebateState 风险辩论状态
type RiskDebateState struct {
	RiskyHistory           string `json:"risky_history"`            // 激进派历史
	SafeHistory            string `json:"safe_history"`             // 保守派历史
	NeutralHistory         string `json:"neutral_history"`          // 中立派历史
	History                string `json:"history"`                  // 完整辩论历史
	LatestSpeaker          string `json:"latest_speaker"`           // 最近发言者
	CurrentRiskyResponse   string `json:"current_risky_response"`   // 激进派最新回复
	CurrentSafeResponse    string `json:"current_safe_response"`    // 保守派最新回复
	CurrentNeutralResponse string `json:"current_neutral_response"` // 中立派最新回复
	JudgeDecision          string `json:"judge_decision"`           // 裁判决定
	Count                  int    `json:"count"`                    // 辩论轮数
}

// AgentState 主状态结构
type AgentState struct {
	// 消息历史
	Messages []*schema.Message `json:"messages"`

	// 交易信息
	CompanyOfInterest string `json:"company_of_interest"` // 交易标的
	TradeDate         string `json:"trade_date"`          // 交易日期
	Sender            string `json:"sender"`              // 发送者

	// 分析报告
	MarketReport       string `json:"market_report"`       // 技术分析报告
	SentimentReport    string `json:"sentiment_report"`    // 情绪分析报告
	NewsReport         string `json:"news_report"`         // 新闻分析报告
	FundamentalsReport string `json:"fundamentals_report"` // 基本面分析报告

	// 辩论状态
	InvestmentDebateState InvestDebateState `json:"investment_debate_state"` // 投资辩论状态
	RiskDebateState       RiskDebateState   `json:"risk_debate_state"`       // 风险辩论状态

	// 决策结果
	InvestmentPlan       string `json:"investment_plan"`        // 投资计划
	TraderInvestmentPlan string `json:"trader_investment_plan"` // 交易员投资计划
	FinalTradeDecision   string `json:"final_trade_decision"`   // 最终交易决策
}

// NewAgentState 创建新的 Agent 状态
func NewAgentState(company, tradeDate string) *AgentState {
	return &AgentState{
		CompanyOfInterest: company,
		TradeDate:         tradeDate,
		Messages:          make([]*schema.Message, 0),
		InvestmentDebateState: InvestDebateState{
			Count: 0,
		},
		RiskDebateState: RiskDebateState{
			Count: 0,
		},
	}
}

// Clone 克隆状态
func (s *AgentState) Clone() *AgentState {
	newState := *s
	newState.Messages = make([]*schema.Message, len(s.Messages))
	copy(newState.Messages, s.Messages)
	return &newState
}

// AddMessage 添加消息
func (s *AgentState) AddMessage(msg *schema.Message) {
	s.Messages = append(s.Messages, msg)
}

// ClearMessages 清除消息
func (s *AgentState) ClearMessages() {
	s.Messages = make([]*schema.Message, 0)
}

// GetAllReports 获取所有分析报告的合并文本
func (s *AgentState) GetAllReports() string {
	return s.MarketReport + "\n\n" + s.SentimentReport + "\n\n" + s.NewsReport + "\n\n" + s.FundamentalsReport
}
