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
