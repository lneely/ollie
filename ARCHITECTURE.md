# Ollama MCP Client - Go Architecture

## Overview
Simple REPL-based CLI that connects Ollama LLMs to MCP servers. Minimal implementation focused on core functionality.

## Directory Structure
```
.
├── main.go                 # Entry point, REPL loop
├── config/
│   └── config.go          # JSON config parser (mcpServers format)
├── mcp/
│   ├── client.go          # MCP protocol client
│   └── transport.go       # STDIO/SSE/HTTP transports
├── ollama/
│   └── client.go          # Ollama API client
└── tools/
    └── executor.go        # Tool execution handler
```

## Core Components

### 1. Main (main.go)
- Simple REPL with `>` prompt
- Multi-line input support (submit on ENTER)
- Read user input → Send to Ollama → Handle tool calls → Display response

### 2. Config (config/config.go)
```go
type Config struct {
    MCPServers map[string]ServerConfig `json:"mcpServers"`
}

type ServerConfig struct {
    Command  string            `json:"command,omitempty"`  // STDIO
    Args     []string          `json:"args,omitempty"`
    Env      map[string]string `json:"env,omitempty"`
    Type     string            `json:"type,omitempty"`     // "sse" or "streamable_http"
    URL      string            `json:"url,omitempty"`      // SSE/HTTP
    Headers  map[string]string `json:"headers,omitempty"`
    Disabled bool              `json:"disabled,omitempty"`
}
```

### 3. MCP Client (mcp/client.go)
- Connect to MCP servers (STDIO/SSE/HTTP)
- List available tools
- Execute tool calls
- Handle JSON-RPC protocol

### 4. Ollama Client (ollama/client.go)
- Send chat requests with tool definitions
- Stream responses
- Parse tool call requests from model

### 5. Tool Executor (tools/executor.go)
- Route tool calls to appropriate MCP server
- Execute and return results
- Handle errors

## Data Flow
```
User Input → Ollama (with tools) → Tool Call? → MCP Server → Tool Result → Ollama → Response
```

## Dependencies
- Standard library (encoding/json, bufio, os/exec)
- Ollama Go SDK (github.com/ollama/ollama/api)
- MCP Go SDK (if available) or custom JSON-RPC implementation

## Implementation Notes
- Keep it minimal - no TUI, just basic REPL
- Focus on STDIO transport first (simplest)
- Add SSE/HTTP support after STDIO works
- Error handling: log and continue
- No conversation history persistence initially
