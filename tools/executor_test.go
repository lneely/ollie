package tools

import (
	"testing"
)

func TestExecutorAddServer(t *testing.T) {
	executor := NewExecutor()
	
	if len(executor.servers) != 0 {
		t.Errorf("Expected 0 servers, got %d", len(executor.servers))
	}
	
	executor.AddServer("test", nil)
	
	if len(executor.servers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(executor.servers))
	}
}

