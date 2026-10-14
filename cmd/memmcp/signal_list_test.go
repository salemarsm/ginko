package main

import (
	"strings"
	"testing"
)

// O handler HTTP garante lista vazia como [] e nao nil
// (server/server.go: "if sigs == nil { sigs = []memory.AgentSignal{} }").
// A tool MCP precisa fazer o mesmo: json.Unmarshal trata null e [] de forma
// diferente, e um monitor que compara comprimento ou itera o resultado quebra
// com null.
func TestSignalListEmptyReturnsArrayNotNull(t *testing.T) {
	srv := &mcpServer{store: openTestStore(t), project: "ginko"}
	out := strings.TrimSpace(callToolJSON(t, srv, "signal_list", `{"project":"ginko"}`))
	if out == "null" {
		t.Fatal("signal_list sem sinais devolveu null; esperado [] (o handler HTTP ja garante isso)")
	}
	if !strings.HasPrefix(out, "[") {
		t.Fatalf("signal_list deveria devolver um array JSON, devolveu %q", out)
	}
}
