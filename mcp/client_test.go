package mcp

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestClientCall(t *testing.T) {
	var buf bytes.Buffer
	client := NewClient(&buf, &buf)
	
	go func() {
		resp := Message{
			JSONRPC: "2.0",
			ID:      1,
			Result:  json.RawMessage(`{"status":"ok"}`),
		}
		data, _ := json.Marshal(resp)
		buf.Write(append(data, '\n'))
	}()
	
	result, err := client.Call("test_method", map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("Call failed: %v", err)
	}
	
	var resultMap map[string]string
	if err := json.Unmarshal(result, &resultMap); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	
	if resultMap["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", resultMap["status"])
	}
}
