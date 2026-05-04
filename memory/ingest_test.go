package memory

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestIngestPathRecursive(t *testing.T) {
	ctx := context.Background()
	s, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "root.md"), []byte("# Root\n\nMemory is conclusion."), 0o644); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(dir, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nested, "child.txt"), []byte("RAG is evidence."), 0o644); err != nil {
		t.Fatal(err)
	}

	resp, err := s.IngestPath(ctx, IngestRequest{Path: dir, Recursive: true})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Run.Status != "ok" || len(resp.Documents) != 2 || len(resp.Chunks) != 2 {
		t.Fatalf("unexpected ingest response: status=%s docs=%d chunks=%d skipped=%d", resp.Run.Status, len(resp.Documents), len(resp.Chunks), len(resp.Skipped))
	}
	docs, err := s.ListDocuments(ctx, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(docs) != 2 {
		t.Fatalf("expected 2 docs, got %d", len(docs))
	}
}
