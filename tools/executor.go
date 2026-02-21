package tools

import (
	"encoding/json"
	"fmt"
	"ollie/mcp"
)

type Executor struct {
	servers map[string]*mcp.Client
}

func NewExecutor() *Executor {
	return &Executor{
		servers: make(map[string]*mcp.Client),
	}
}

func (e *Executor) AddServer(name string, client *mcp.Client) {
	e.servers[name] = client
}

func (e *Executor) ListTools() ([]ToolInfo, error) {
	var allTools []ToolInfo
	for serverName, client := range e.servers {
		result, err := client.Call("tools/list", nil)
		if err != nil {
			return nil, fmt.Errorf("server %s: %w", serverName, err)
		}
		var resp struct {
			Tools []struct {
				Name        string          `json:"name"`
				Description string          `json:"description"`
				InputSchema json.RawMessage `json:"inputSchema"`
			} `json:"tools"`
		}
		if err := json.Unmarshal(result, &resp); err != nil {
			return nil, err
		}
		for _, t := range resp.Tools {
			allTools = append(allTools, ToolInfo{
				Server:      serverName,
				Name:        t.Name,
				Description: t.Description,
				InputSchema: t.InputSchema,
			})
		}
	}
	return allTools, nil
}

func (e *Executor) Execute(serverName, toolName string, args json.RawMessage) (json.RawMessage, error) {
	client, ok := e.servers[serverName]
	if !ok {
		return nil, fmt.Errorf("server not found: %s", serverName)
	}
	params := map[string]interface{}{
		"name":      toolName,
		"arguments": args,
	}
	return client.Call("tools/call", params)
}

type ToolInfo struct {
	Server      string
	Name        string
	Description string
	InputSchema json.RawMessage
}
