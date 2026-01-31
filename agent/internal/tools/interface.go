// Package tools 数据工具定义
package tools

import (
	"context"

	"github.com/cloudwego/eino/schema"
)

// Tool 工具接口
type Tool interface {
	// Name 工具名称
	Name() string
	// Description 工具描述
	Description() string
	// Parameters 参数 schema
	Parameters() map[string]interface{}
	// Execute 执行工具
	Execute(ctx context.Context, args map[string]interface{}) (string, error)
}

// BaseTool 基础工具
type BaseTool struct {
	name        string
	description string
	parameters  map[string]interface{}
	executeFunc func(ctx context.Context, args map[string]interface{}) (string, error)
}

// NewTool 创建新工具
func NewTool(name, description string, parameters map[string]interface{}, executeFunc func(ctx context.Context, args map[string]interface{}) (string, error)) *BaseTool {
	return &BaseTool{
		name:        name,
		description: description,
		parameters:  parameters,
		executeFunc: executeFunc,
	}
}

// Name 实现 Tool 接口
func (t *BaseTool) Name() string {
	return t.name
}

// Description 实现 Tool 接口
func (t *BaseTool) Description() string {
	return t.description
}

// Parameters 实现 Tool 接口
func (t *BaseTool) Parameters() map[string]interface{} {
	return t.parameters
}

// Execute 实现 Tool 接口
func (t *BaseTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	return t.executeFunc(ctx, args)
}

// ToToolInfo 转换为 schema.ToolInfo
func (t *BaseTool) ToToolInfo() schema.ToolInfo {
	return schema.ToolInfo{
		Name:        t.name,
		Desc:        t.description,
		ParamsOneOf: t.parameters,
	}
}

// ToolRegistry 工具注册表
type ToolRegistry struct {
	tools map[string]Tool
}

// NewToolRegistry 创建工具注册表
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]Tool),
	}
}

// Register 注册工具
func (r *ToolRegistry) Register(tool Tool) {
	r.tools[tool.Name()] = tool
}

// Get 获取工具
func (r *ToolRegistry) Get(name string) (Tool, bool) {
	tool, ok := r.tools[name]
	return tool, ok
}

// GetAll 获取所有工具
func (r *ToolRegistry) GetAll() []Tool {
	tools := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	return tools
}

// GetToolInfos 获取所有工具信息
func (r *ToolRegistry) GetToolInfos() []schema.ToolInfo {
	infos := make([]schema.ToolInfo, 0, len(r.tools))
	for _, tool := range r.tools {
		if bt, ok := tool.(*BaseTool); ok {
			infos = append(infos, bt.ToToolInfo())
		}
	}
	return infos
}

// ExecuteTool 执行工具
func (r *ToolRegistry) ExecuteTool(ctx context.Context, name string, args map[string]interface{}) (string, error) {
	tool, ok := r.Get(name)
	if !ok {
		return "", nil
	}
	return tool.Execute(ctx, args)
}
