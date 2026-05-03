package memory

import (
	"context"
	"strings"
	"testing"
)

func TestBuildContextBudget(t *testing.T) {
	s, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	ctx := context.Background()

	_, err = s.UpsertMemory(ctx, Memory{Type: TypeFact, Subject: "botmaster", Content: "OpenClaw should receive compact memory context, not raw documents.", Source: Source{Kind: "test", Ref: "ctx"}, Scope: ScopeProject, Confidence: 0.9})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := s.BuildContext(ctx, ContextRequest{Query: "compact memory context", Subject: "botmaster", Scopes: []Scope{ScopeProject}, MaxTokens: 200})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resp.Context, "compact memory context") {
		t.Fatalf("unexpected context: %q", resp.Context)
	}
	if resp.EstimatedTokens <= 0 || resp.EstimatedTokens > resp.BudgetTokens {
		t.Fatalf("bad token budget: %#v", resp)
	}
}
