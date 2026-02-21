# Ollie

A Go CLI that connects Ollama LLMs to MCP (Model Context Protocol) servers for tool use. Built because I couldn't find what I wanted with existing tools.

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

Chat history displays in a scrollable viewport. Type in the multi-line text area below, submit with Enter (Ctrl+J for new lines).

Tools depend on which MCP servers you configure. Example with filesystem server:

```
> What files are in the current directory?
[Calling tool from filesystem server]
Found: main.go, config/, mcp/, ollama/, tools/
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
