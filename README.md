# Ollama MCP Client

A minimal Go CLI application that connects Ollama LLMs to MCP (Model Context Protocol) servers, enabling tool use.

## Features

- Simple REPL interface with `>` prompt
- Multi-line input support (submit on ENTER)
- JSON configuration for MCP servers
- STDIO transport support for MCP servers
- Tool discovery and execution
- Conversation history

## Requirements

- Go 1.21+
- Ollama running locally (default: http://localhost:11434)
- MCP servers configured in JSON file

## Installation

```bash
go build -o ollama-mcp-client
```

## Configuration

Create a JSON configuration file with your MCP servers:

```json
{
  "mcpServers": {
    "anvilmcp": {
      "command": "anvilmcp",
      "args": []
    },
    "filesystem": {
      "command": "mcp-server-filesystem",
      "args": ["/path/to/allowed/directory"],
      "env": {
        "DEBUG": "true"
      }
    },
    "disabled-server": {
      "command": "some-server",
      "disabled": true
    }
  }
}
```

### Configuration Options

- `command`: Executable command for STDIO MCP server
- `args`: Command-line arguments (optional)
- `env`: Environment variables (optional)
- `disabled`: Skip this server (optional)

## Usage

```bash
./ollama-mcp-client config.json [model]
```

Default model: `qwen3:8b`

Example with custom model:
```bash
./ollama-mcp-client config.json llama3.2
```

### Example Session

```
> What files are in the current directory?
[Tool execution: list_directory]
The current directory contains: main.go, config/, mcp/, ollama/, tools/

> Create a new file called test.txt
[Tool execution: write_file]
Created test.txt successfully.
```

## Architecture

- `main.go`: REPL loop and orchestration
- `config/`: JSON configuration parser
- `mcp/`: MCP protocol client and STDIO transport
- `ollama/`: Ollama API client
- `tools/`: Tool executor and manager

## Troubleshooting

### No servers connected
- Check that MCP server commands are in your PATH
- Verify command names and arguments in config
- Check logs for connection errors

### Tools not working
- Ensure MCP servers support the `tools/list` and `tools/call` methods
- Check tool execution logs for errors
- Verify tool arguments match expected schema

### Ollama errors
- Ensure Ollama is running: `ollama serve`
- Check model is available: `ollama list`
- Verify Ollama URL (default: http://localhost:11434)

## License

MIT
