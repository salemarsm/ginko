package memory

import (
	"context"
	"testing"
)

func TestMemoryTimelineIncludesDirectAndContextEvents(t *testing.T) {
	ctx := context.Background()
	s, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	m, err := s.UpsertMemory(ctx, Memory{Type: TypeDecision, Subject: "project", Content: "Use SQLite as canonical truth.", Source: Source{Kind: "test", Ref: "timeline"}, Scope: ScopeProject, Confidence: 0.9})
	if err != nil {
		t.Fatal(err)
	}
	if err := s.AppendEvent(ctx, Event{MemoryID: &m.ID, Kind: "memory.upserted", Payload: m.ID, Source: Source{Kind: "test", Ref: "timeline"}}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.BuildContext(ctx, ContextRequest{Query: "SQLite canonical", Subject: "project", Scopes: []Scope{ScopeProject}, MaxTokens: 200}); err != nil {
		t.Fatal(err)
	}
	out, err := s.MemoryTimeline(ctx, m.ID, 20)
	if err != nil {
		t.Fatal(err)
	}
	if out.Memory == nil || out.Memory.ID != m.ID {
		t.Fatalf("missing memory: %#v", out.Memory)
	}
	kinds := map[string]bool{}
	for _, e := range out.Events {
		kinds[e.Kind] = true
	}
	if !kinds["memory.upserted"] || !kinds["context.built"] {
		t.Fatalf("expected upsert and context events, got %#v", out.Events)
	}
}
