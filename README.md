# Ollie

A Go CLI that connects Ollama LLMs to MCP (Model Context Protocol) servers for tool use.

## Quick Start

**Requirements:**
- Go 1.21+
- Ollama running locally (http://localhost:11434)

**Build:**
```bash
go build
```

**Run:**
```bash
./ollie config.json [model]
```

Default model is `qwen3:8b`. Override with any Ollama model name.

## Configuration

Create a JSON file defining your MCP servers:

```json
{
  "mcpServers": {
    "filesystem": {
      "command": "mcp-server-filesystem",
      "args": ["/path/to/allowed/directory"],
      "env": {
        "DEBUG": "true"
      }
    },
    "disabled-example": {
      "command": "some-server",
      "disabled": true
    }
  }
}
```

**Options:**
- `command` - MCP server executable (required)
- `args` - Command arguments (optional)
- `env` - Environment variables (optional)
- `disabled` - Skip this server (optional)

MCP server commands must be in your PATH.

## Usage

The REPL accepts multi-line input and submits on ENTER. Tools are discovered automatically from connected MCP servers.

```
> What files are in the current directory?
[Tool: list_directory]
Found: main.go, config/, mcp/, ollama/, tools/

> Create a file called test.txt
[Tool: write_file]
Created test.txt
```

## Troubleshooting

**No servers connected**
- Verify MCP server commands are in PATH
- Check command names and arguments in config

**Tools not working**
- Ensure servers support `tools/list` and `tools/call`
- Check tool arguments match expected schema

**Ollama errors**
- Start Ollama: `ollama serve`
- Verify model exists: `ollama list`

## License

GPLv3
