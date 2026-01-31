// Package tools 工具测试
package tools

import (
	"context"
	"testing"
)

func TestNewTool(t *testing.T) {
	tool := NewTool(
		"test_tool",
		"A test tool",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"param1": map[string]interface{}{"type": "string"},
			},
		},
		func(ctx context.Context, args map[string]interface{}) (string, error) {
			return "test result", nil
		},
	)

	if tool.Name() != "test_tool" {
		t.Errorf("Expected name 'test_tool', got '%s'", tool.Name())
	}

	if tool.Description() != "A test tool" {
		t.Errorf("Expected description 'A test tool', got '%s'", tool.Description())
	}

	params := tool.Parameters()
	if params == nil {
		t.Error("Parameters should not be nil")
	}
}

func TestToolExecute(t *testing.T) {
	executed := false
	tool := NewTool(
		"execute_test",
		"Test execution",
		map[string]interface{}{},
		func(ctx context.Context, args map[string]interface{}) (string, error) {
			executed = true
			return "executed", nil
		},
	)

	result, err := tool.Execute(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !executed {
		t.Error("Tool should have been executed")
	}

	if result != "executed" {
		t.Errorf("Expected 'executed', got '%s'", result)
	}
}

func TestToolToToolInfo(t *testing.T) {
	tool := NewTool(
		"info_test",
		"Test tool info",
		map[string]interface{}{"type": "object"},
		func(ctx context.Context, args map[string]interface{}) (string, error) {
			return "", nil
		},
	)

	info := tool.ToToolInfo()

	if info.Name != "info_test" {
		t.Errorf("Expected name 'info_test', got '%s'", info.Name)
	}

	if info.Desc != "Test tool info" {
		t.Errorf("Expected desc 'Test tool info', got '%s'", info.Desc)
	}
}

func TestToolRegistry(t *testing.T) {
	registry := NewToolRegistry()

	tool1 := NewTool("tool1", "Tool 1", nil, nil)
	tool2 := NewTool("tool2", "Tool 2", nil, nil)

	registry.Register(tool1)
	registry.Register(tool2)

	// 测试 Get
	got, ok := registry.Get("tool1")
	if !ok {
		t.Error("Should find tool1")
	}
	if got.Name() != "tool1" {
		t.Errorf("Expected 'tool1', got '%s'", got.Name())
	}

	// 测试不存在的工具
	_, ok = registry.Get("nonexistent")
	if ok {
		t.Error("Should not find nonexistent tool")
	}

	// 测试 GetAll
	all := registry.GetAll()
	if len(all) != 2 {
		t.Errorf("Expected 2 tools, got %d", len(all))
	}
}

func TestGetStockDataTool(t *testing.T) {
	if GetStockDataTool == nil {
		t.Fatal("GetStockDataTool should not be nil")
	}

	if GetStockDataTool.Name() != "get_stock_data" {
		t.Errorf("Expected name 'get_stock_data', got '%s'", GetStockDataTool.Name())
	}

	// 测试执行
	result, err := GetStockDataTool.Execute(context.Background(), map[string]interface{}{
		"symbol":     "NVDA",
		"start_date": "2024-05-01",
		"end_date":   "2024-05-10",
	})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == "" {
		t.Error("Result should not be empty")
	}
}

func TestGetIndicatorsTool(t *testing.T) {
	if GetIndicatorsTool == nil {
		t.Fatal("GetIndicatorsTool should not be nil")
	}

	if GetIndicatorsTool.Name() != "get_indicators" {
		t.Errorf("Expected name 'get_indicators', got '%s'", GetIndicatorsTool.Name())
	}

	result, err := GetIndicatorsTool.Execute(context.Background(), map[string]interface{}{
		"symbol": "NVDA",
		"date":   "2024-05-10",
	})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == "" {
		t.Error("Result should not be empty")
	}
}

func TestGetFundamentalsTool(t *testing.T) {
	if GetFundamentalsTool == nil {
		t.Fatal("GetFundamentalsTool should not be nil")
	}

	if GetFundamentalsTool.Name() != "get_fundamentals" {
		t.Errorf("Expected name 'get_fundamentals', got '%s'", GetFundamentalsTool.Name())
	}

	result, err := GetFundamentalsTool.Execute(context.Background(), map[string]interface{}{
		"symbol": "NVDA",
	})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == "" {
		t.Error("Result should not be empty")
	}
}

func TestGetNewsTool(t *testing.T) {
	if GetNewsTool == nil {
		t.Fatal("GetNewsTool should not be nil")
	}

	if GetNewsTool.Name() != "get_news" {
		t.Errorf("Expected name 'get_news', got '%s'", GetNewsTool.Name())
	}

	result, err := GetNewsTool.Execute(context.Background(), map[string]interface{}{
		"query":      "NVDA",
		"start_date": "2024-05-01",
		"end_date":   "2024-05-10",
	})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == "" {
		t.Error("Result should not be empty")
	}
}

func TestGetGlobalNewsTool(t *testing.T) {
	if GetGlobalNewsTool == nil {
		t.Fatal("GetGlobalNewsTool should not be nil")
	}

	if GetGlobalNewsTool.Name() != "get_global_news" {
		t.Errorf("Expected name 'get_global_news', got '%s'", GetGlobalNewsTool.Name())
	}

	result, err := GetGlobalNewsTool.Execute(context.Background(), map[string]interface{}{
		"date":           "2024-05-10",
		"look_back_days": float64(7),
		"limit":          float64(10),
	})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == "" {
		t.Error("Result should not be empty")
	}
}
