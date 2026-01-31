// Package llm OpenAI API 客户端
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cloudwego/eino/schema"
)

// OpenAI API 请求/响应结构

type openAIMessage struct {
	Role       string          `json:"role"`
	Content    string          `json:"content,omitempty"`
	ToolCalls  []openAIToolCall `json:"tool_calls,omitempty"`
	ToolCallID string          `json:"tool_call_id,omitempty"`
}

type openAIToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type openAITool struct {
	Type     string `json:"type"`
	Function struct {
		Name        string      `json:"name"`
		Description string      `json:"description"`
		Parameters  interface{} `json:"parameters"`
	} `json:"function"`
}

type openAIRequest struct {
	Model    string          `json:"model"`
	Messages []openAIMessage `json:"messages"`
	Tools    []openAITool    `json:"tools,omitempty"`
}

type openAIResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Role       string           `json:"role"`
			Content    string           `json:"content"`
			ToolCalls  []openAIToolCall `json:"tool_calls,omitempty"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// callOpenAI 调用 OpenAI API
func callOpenAI(ctx context.Context, apiKey, baseURL, model string, messages []*schema.Message) (*schema.Message, error) {
	return callOpenAIWithTools(ctx, apiKey, baseURL, model, messages, nil)
}

// callOpenAIWithTools 调用带工具的 OpenAI API
func callOpenAIWithTools(ctx context.Context, apiKey, baseURL, model string, messages []*schema.Message, tools []schema.ToolInfo) (*schema.Message, error) {
	// 转换消息格式
	openAIMessages := make([]openAIMessage, 0, len(messages))
	for _, msg := range messages {
		openAIMsg := openAIMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
		if msg.ToolCallID != "" {
			openAIMsg.ToolCallID = msg.ToolCallID
		}
		openAIMessages = append(openAIMessages, openAIMsg)
	}

	// 构建请求
	req := openAIRequest{
		Model:    model,
		Messages: openAIMessages,
	}

	// 添加工具
	if len(tools) > 0 {
		openAITools := make([]openAITool, 0, len(tools))
		for _, tool := range tools {
			openAITool := openAITool{
				Type: "function",
			}
			openAITool.Function.Name = tool.Name
			openAITool.Function.Description = tool.Desc
			openAITool.Function.Parameters = tool.ParamsOneOf
			openAITools = append(openAITools, openAITool)
		}
		req.Tools = openAITools
	}

	// 序列化请求
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 构建 HTTP 请求
	url := baseURL + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(respBody))
	}

	// 解析响应
	var openAIResp openAIResponse
	if err := json.Unmarshal(respBody, &openAIResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	choice := openAIResp.Choices[0]

	// 构建返回消息
	resultMsg := &schema.Message{
		Role:    schema.Assistant,
		Content: choice.Message.Content,
	}

	// 处理工具调用
	if len(choice.Message.ToolCalls) > 0 {
		toolCalls := make([]schema.ToolCall, 0, len(choice.Message.ToolCalls))
		for _, tc := range choice.Message.ToolCalls {
			toolCalls = append(toolCalls, schema.ToolCall{
				ID:   tc.ID,
				Type: tc.Type,
				Function: schema.FunctionCall{
					Name:      tc.Function.Name,
					Arguments: tc.Function.Arguments,
				},
			})
		}
		resultMsg.ToolCalls = toolCalls
	}

	return resultMsg, nil
}
