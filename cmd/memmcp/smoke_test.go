package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMCPInitialize(t *testing.T) {
	store := openTestStore(t)
	srv := &mcpServer{store: store, project: "test"}

	req := rpcRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params:  json.RawMessage(`{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"0.0.1"}}`),
	}
	resp := srv.handle(req)

	if resp.Error != nil {
		t.Fatalf("initialize error: %+v", resp.Error)
	}
	m, ok := resp.Result.(map[string]any)
	if !ok {
		t.Fatal("expected map result")
	}
	if m["protocolVersion"] != "2024-11-05" {
		t.Fatalf("unexpected protocolVersion: %v", m["protocolVersion"])
	}
	if _, ok := m["serverInfo"]; !ok {
		t.Fatal("missing serverInfo")
	}
}

func TestMCPToolsList(t *testing.T) {
	store := openTestStore(t)
	srv := &mcpServer{store: store, project: "test"}

	resp := srv.handle(rpcRequest{JSONRPC: "2.0", ID: 2, Method: "tools/list"})
	if resp.Error != nil {
		t.Fatalf("tools/list error: %+v", resp.Error)
	}
	m, _ := resp.Result.(map[string]any)
	list, _ := m["tools"].([]map[string]any)
	if len(list) == 0 {
		// result may be []any
		raw, _ := json.Marshal(m["tools"])
		var items []any
		_ = json.Unmarshal(raw, &items)
		if len(items) == 0 {
			t.Fatal("expected at least one tool")
		}
	}
}

func TestMCPUnknownMethod(t *testing.T) {
	store := openTestStore(t)
	srv := &mcpServer{store: store}

	resp := srv.handle(rpcRequest{JSONRPC: "2.0", ID: 3, Method: "unknown/method"})
	if resp.Error == nil {
		t.Fatal("expected error for unknown method")
	}
	if resp.Error.Code != -32601 {
		t.Fatalf("expected code -32601, got %d", resp.Error.Code)
	}
}

func TestMCPRememberAndContext(t *testing.T) {
	store := openTestStore(t)
	srv := &mcpServer{store: store, project: "smoke-test"}

	// save a memory
	remResp := srv.handle(rpcRequest{
		JSONRPC: "2.0", ID: 4, Method: "tools/call",
		Params: json.RawMessage(`{"name":"memory_remember","arguments":{"type":"fact","subject":"smoke-test","content":"MCP smoke test memory","scope":"project","confidence":0.9}}`),
	})
	if remResp.Error != nil {
		t.Fatalf("memory_remember error: %+v", remResp.Error)
	}

	// retrieve context
	ctxResp := srv.handle(rpcRequest{
		JSONRPC: "2.0", ID: 5, Method: "tools/call",
		Params: json.RawMessage(`{"name":"memory_context","arguments":{"query":"MCP smoke","subject":"smoke-test"}}`),
	})
	if ctxResp.Error != nil {
		t.Fatalf("memory_context error: %+v", ctxResp.Error)
	}
	raw, _ := json.Marshal(ctxResp.Result)
	if !strings.Contains(string(raw), "smoke") {
		t.Fatalf("expected context to contain saved memory, got: %s", string(raw))
	}
}
