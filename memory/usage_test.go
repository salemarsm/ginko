package memory

import (
	"context"
	"encoding/json"
	"testing"
)

func TestMemoryUsageStatsFromContextEvents(t *testing.T) {
	ctx := context.Background()
	s, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	m, err := s.UpsertMemory(ctx, Memory{Type: TypeDecision, Subject: "proj", Content: "Use SQLite as canonical memory.", Source: Source{Kind: "test", Ref: "usage"}, Scope: ScopeProject, Confidence: 0.9})
	if err != nil {
		t.Fatal(err)
	}
	payload, _ := json.Marshal(map[string]any{"context_id": "ctx_usage", "memory_ids": []string{m.ID}})
	if err := s.AppendEvent(ctx, Event{Kind: "context.built", Payload: string(payload), Source: Source{Kind: "test", Ref: "usage"}}); err != nil {
		t.Fatal(err)
	}
	if err := s.RecordContextFeedback(ctx, ContextFeedback{ContextID: "ctx_usage", Useful: true, MemoryIDsUsed: []string{m.ID}, Source: Source{Kind: "test", Ref: "usage"}}); err != nil {
		t.Fatal(err)
	}

	rows, err := s.ListMemoryUsage(ctx, Query{Subject: "proj", Limit: 10}, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected one usage row, got %#v", rows)
	}
	if rows[0].Usage.ContextUses != 1 || rows[0].Usage.UsefulVotes != 1 || rows[0].Usage.MemoryID != m.ID {
		t.Fatalf("bad usage stats: %#v", rows[0].Usage)
	}
}
