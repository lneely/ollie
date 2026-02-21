package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"ollie/config"
	"ollie/mcp"
	"ollie/ollama"
	"ollie/tools"
)

type model struct {
	textarea textarea.Model
	viewport viewport.Model
	messages []ollama.Message
	client   *ollama.Client
	tools    []tools.ToolInfo
	executor *tools.Executor
	model    string
	history  []string
	ready    bool
	hooks    map[string]string
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: ollama-mcp-client <config.json> [model]")
	}

	modelName := "qwen3:8b"
	if len(os.Args) > 2 {
		modelName = os.Args[2]
	}

	cfg, err := config.Load(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	hooks := make(map[string]string)
	if cfg.Hooks != nil {
		hooks = cfg.Hooks
	}

	executor := tools.NewExecutor()
	connectedServers := 0
	for name, serverCfg := range cfg.MCPServers {
		if serverCfg.Disabled {
			log.Printf("Skipping disabled server: %s", name)
			continue
		}
		if serverCfg.Command != "" {
			transport := mcp.NewSTDIOTransport(serverCfg.Command, serverCfg.Args, serverCfg.Env)
			client, err := transport.Connect()
			if err != nil {
				log.Printf("Failed to connect to %s: %v", name, err)
				continue
			}
			executor.AddServer(name, client)
			connectedServers++
			log.Printf("Connected to server: %s", name)
		}
	}

	if connectedServers == 0 {
		log.Fatal("No MCP servers connected")
	}

	toolsList, err := executor.ListTools()
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}
	log.Printf("Loaded %d tools", len(toolsList))

	ta := textarea.New()
	ta.Placeholder = "Type your message..."
	ta.Prompt = ""
	ta.ShowLineNumbers = false
	ta.CharLimit = 0
	ta.SetHeight(5)
	ta.KeyMap.InsertNewline = key.NewBinding(key.WithKeys("ctrl+j"))
	ta.Focus()

	p := tea.NewProgram(model{
		textarea: ta,
		client:   ollama.NewClient("http://localhost:11434"),
		tools:    toolsList,
		executor: executor,
		model:    modelName,
		hooks:    hooks,
	})

	if hook := hooks["agentSpawn"]; hook != "" {
		exec.Command("sh", "-c", hook).Run()
	}

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

type responseMsg struct {
	history []string
	messages []ollama.Message
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-5)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 5
		}
		m.textarea.SetWidth(msg.Width)
	case responseMsg:
		m.history = msg.history
		m.messages = msg.messages
		m.viewport.SetContent(strings.Join(m.history, "\n"))
		m.viewport.GotoBottom()
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "esc" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			input := strings.TrimSpace(m.textarea.Value())
			if input == "" {
				return m, nil
			}
			m.history = append(m.history, "You: "+input)
			m.viewport.SetContent(strings.Join(m.history, "\n"))
			m.viewport.GotoBottom()
			m.textarea.Reset()
			
			if hook := m.hooks["userPromptSubmit"]; hook != "" {
				exec.Command("sh", "-c", hook).Run()
			}
			
			return m, m.handleInputCmd(input)
		}
	}

	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m model) handleInputCmd(input string) tea.Cmd {
	return func() tea.Msg {
		history := append([]string{}, m.history...)
		messages := append([]ollama.Message{}, m.messages...)
		
		messages = append(messages, ollama.Message{Role: "user", Content: input})

		ollamaTools := make([]ollama.Tool, len(m.tools))
		for i, t := range m.tools {
			ollamaTools[i] = ollama.Tool{
				Type: "function",
				Function: ollama.Function{
					Name:        t.Name,
					Description: t.Description,
					Parameters:  t.InputSchema,
				},
			}
		}

		resp, err := m.client.Chat(ollama.ChatRequest{
			Model:    m.model,
			Messages: messages,
			Tools:    ollamaTools,
			Stream:   false,
		})
		if err != nil {
			history = append(history, fmt.Sprintf("Error: %v", err))
			return responseMsg{history: history, messages: messages}
		}

		messages = append(messages, resp.Message)

		seenCalls := make(map[string]bool)
		for len(resp.Message.ToolCalls) > 0 {
			for _, tc := range resp.Message.ToolCalls {
				callKey := fmt.Sprintf("%s:%v", tc.Function.Name, tc.Function.Arguments)
				if seenCalls[callKey] {
					history = append(history, fmt.Sprintf("Skipping duplicate tool call: %s", tc.Function.Name))
					continue
				}
				seenCalls[callKey] = true
				
				history = append(history, fmt.Sprintf("Running tool: %s", tc.Function.Name))
				var toolInfo *tools.ToolInfo
				for _, t := range m.tools {
					if t.Name == tc.Function.Name {
						toolInfo = &t
						break
					}
				}
				if toolInfo == nil {
					history = append(history, fmt.Sprintf("→ %s: error - tool not found", tc.Function.Name))
					continue
				}
				result, err := m.executor.Execute(toolInfo.Server, tc.Function.Name, tc.Function.Arguments)
				if err != nil {
					history = append(history, fmt.Sprintf("→ %s: error - %v", tc.Function.Name, err))
					continue
				}
				resultStr := string(result)
				history = append(history, fmt.Sprintf("→ %s: success - %s", tc.Function.Name, resultStr))
				messages = append(messages, ollama.Message{
					Role:    "tool",
					Content: resultStr,
				})
			}
			
			resp, err = m.client.Chat(ollama.ChatRequest{
				Model:    m.model,
				Messages: messages,
				Tools:    ollamaTools,
				Stream:   false,
			})
			if err != nil {
				history = append(history, fmt.Sprintf("Error: %v", err))
				return responseMsg{history: history, messages: messages}
			}
			messages = append(messages, resp.Message)
		}

		if resp.Message.Content != "" {
			history = append(history, "Bot: "+resp.Message.Content)
		}
		
		if hook := m.hooks["stop"]; hook != "" {
			exec.Command("sh", "-c", hook).Run()
		}
		
		return responseMsg{history: history, messages: messages}
	}
}

func (m model) View() string {
	if !m.ready {
		return "Loading..."
	}
	return m.viewport.View() + "\n" + m.textarea.View()
}


