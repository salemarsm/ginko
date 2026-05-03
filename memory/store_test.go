package memory

import (
	"context"
	"testing"
)

func TestStoreLifecycle(t *testing.T) {
	ctx := context.Background()
	s, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	m, err := s.UpsertMemory(ctx, Memory{
		Type:       TypePreference,
		Subject:    "user",
		Content:    "Prefere respostas diretas e técnicas.",
		Source:     Source{Kind: "test", Ref: "msg-1"},
		Scope:      ScopeGlobal,
		Confidence: 0.9,
		Tags:       []string{"style"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if m.ID == "" {
		t.Fatal("expected generated id")
	}

	got, err := s.GetMemory(ctx, m.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Content != m.Content {
		t.Fatalf("content mismatch: %q", got.Content)
	}

	results, err := s.Search(ctx, Query{Text: "diretas", Subject: "user", Limit: 5})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	newer, err := s.Supersede(ctx, m.ID, Memory{
		Type:       TypePreference,
		Subject:    "user",
		Content:    "Prefere respostas extremamente diretas.",
		Source:     Source{Kind: "test", Ref: "msg-2"},
		Scope:      ScopeGlobal,
		Confidence: 0.95,
		Tags:       []string{"style"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if newer.SupersedesID == nil || *newer.SupersedesID != m.ID {
		t.Fatal("expected supersedes link")
	}

	results, err = s.Search(ctx, Query{Subject: "user", Limit: 5})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].ID != newer.ID {
		t.Fatalf("expected only newer active memory, got %#v", results)
	}

	if err := s.Forget(ctx, newer.ID); err != nil {
		t.Fatal(err)
	}
}

func TestSearchTagsAreExact(t *testing.T) {
	ctx := context.Background()
	s, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	_, err = s.UpsertMemory(ctx, Memory{Type: TypeNote, Subject: "user", Content: "short tag", Source: Source{Kind: "test", Ref: "tag-1"}, Scope: ScopeGlobal, Confidence: 0.8, Tags: []string{"proj"}})
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.UpsertMemory(ctx, Memory{Type: TypeNote, Subject: "user", Content: "long tag", Source: Source{Kind: "test", Ref: "tag-2"}, Scope: ScopeGlobal, Confidence: 0.8, Tags: []string{"project"}})
	if err != nil {
		t.Fatal(err)
	}
	items, err := s.Search(ctx, Query{Tags: []string{"proj"}, Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Content != "short tag" {
		t.Fatalf("expected exact tag match, got %#v", items)
	}
}

func TestAppendEventMemoryID(t *testing.T) {
	ctx := context.Background()
	s, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	memID := "mem_test"
	if err := s.AppendEvent(ctx, Event{Kind: "memory.test", MemoryID: &memID, Payload: "ok", Source: Source{Kind: "test", Ref: "event"}}); err != nil {
		t.Fatal(err)
	}
	events, err := s.ListEvents(ctx, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 || events[0].MemoryID == nil || *events[0].MemoryID != memID {
		t.Fatalf("missing event memory id: %#v", events)
	}
}
