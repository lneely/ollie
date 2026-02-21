package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	content := `{
		"mcpServers": {
			"test": {
				"command": "test-cmd",
				"args": ["arg1"],
				"env": {"KEY": "value"}
			}
		}
	}`
	
	tmpfile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()
	
	cfg, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	
	if len(cfg.MCPServers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(cfg.MCPServers))
	}
	
	server := cfg.MCPServers["test"]
	if server.Command != "test-cmd" {
		t.Errorf("Expected command 'test-cmd', got '%s'", server.Command)
	}
	if len(server.Args) != 1 || server.Args[0] != "arg1" {
		t.Errorf("Expected args ['arg1'], got %v", server.Args)
	}
	if server.Env["KEY"] != "value" {
		t.Errorf("Expected env KEY=value, got %v", server.Env)
	}
}
