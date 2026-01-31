// Package main TradingAgents Golang CLI
package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tradingagents/agent/internal/config"
	"github.com/tradingagents/agent/internal/graph"
)

var (
	cfgFile   string
	ticker    string
	tradeDate string
	debug     bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "tradingagents",
		Short: "TradingAgents - Multi-Agent LLM Financial Trading Framework",
		Long: `TradingAgents is a multi-agent trading framework that mirrors the dynamics 
of real-world trading firms. By deploying specialized LLM-powered agents, 
the platform collaboratively evaluates market conditions and informs trading decisions.`,
	}

	// Run command
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run trading analysis for a stock",
		RunE:  runAnalysis,
	}

	runCmd.Flags().StringVarP(&ticker, "ticker", "t", "NVDA", "Stock ticker symbol")
	runCmd.Flags().StringVarP(&tradeDate, "date", "d", "2024-05-10", "Trade date (YYYY-MM-DD)")
	runCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "Config file path")
	runCmd.Flags().BoolVarP(&debug, "debug", "v", false, "Enable debug mode")

	rootCmd.AddCommand(runCmd)

	// Version command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("TradingAgents v1.0.0 (Golang + Eino)")
		},
	}
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runAnalysis(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// 加载配置
	var cfg *config.Config
	var err error

	if cfgFile != "" {
		cfg, err = config.LoadFromFile(cfgFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
	} else {
		cfg = config.DefaultConfig()
	}

	// 检查 API Key
	if cfg.APIKey == "" {
		return fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║           TradingAgents - Multi-Agent Framework            ║")
	fmt.Println("║                  Golang + Eino Edition                     ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	fmt.Printf("📈 Ticker: %s\n", ticker)
	fmt.Printf("📅 Date: %s\n", tradeDate)
	fmt.Printf("🤖 Deep Think LLM: %s\n", cfg.DeepThinkLLM)
	fmt.Printf("⚡ Quick Think LLM: %s\n", cfg.QuickThinkLLM)
	fmt.Printf("🔄 Debate Rounds: %d\n", cfg.MaxDebateRounds)
	fmt.Printf("⚠️ Risk Discuss Rounds: %d\n", cfg.MaxRiskDiscussRounds)
	fmt.Println()

	// 创建交易图
	tradingGraph, err := graph.NewTradingAgentsGraph(ctx, cfg, nil)
	if err != nil {
		return fmt.Errorf("failed to create trading graph: %w", err)
	}

	// 执行分析
	state, decision, err := tradingGraph.Propagate(ctx, ticker, tradeDate)
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	// 打印结果
	fmt.Println("\n" + strings.Repeat("═", 60))
	fmt.Println("                      ANALYSIS RESULTS")
	fmt.Println(strings.Repeat("═", 60))

	if debug {
		fmt.Println("\n📊 Market Report:")
		fmt.Println(truncate(state.MarketReport, 500))

		fmt.Println("\n💭 Sentiment Report:")
		fmt.Println(truncate(state.SentimentReport, 500))

		fmt.Println("\n📰 News Report:")
		fmt.Println(truncate(state.NewsReport, 500))

		fmt.Println("\n📈 Fundamentals Report:")
		fmt.Println(truncate(state.FundamentalsReport, 500))
	}

	fmt.Println("\n💼 Investment Plan:")
	fmt.Println(truncate(state.InvestmentPlan, 800))

	fmt.Println("\n⚖️ Final Trade Decision:")
	fmt.Println(state.FinalTradeDecision)

	fmt.Println("\n" + strings.Repeat("═", 60))
	fmt.Printf("                 🎯 FINAL DECISION: %s\n", decision)
	fmt.Println(strings.Repeat("═", 60))

	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
