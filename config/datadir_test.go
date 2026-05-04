package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMaybeMigrateLegacyDataDir(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("USERPROFILE", tmp)

	legacy := filepath.Join(tmp, ".llm-memory")
	if err := os.MkdirAll(legacy, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(legacy, "memory.db"), []byte("dbcontent"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := MaybeMigrateLegacyDataDir(); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	newDB := filepath.Join(tmp, ".ginko", "ginko.db")
	b, err := os.ReadFile(newDB)
	if err != nil || string(b) != "dbcontent" {
		t.Fatalf("expected ginko.db with original content; got err=%v content=%q", err, string(b))
	}
	if _, err := os.Stat(filepath.Join(tmp, ".ginko", "MIGRATED_FROM_LLM_MEMORY")); err != nil {
		t.Fatalf("expected migration marker: %v", err)
	}
	if err := MaybeMigrateLegacyDataDir(); err != nil {
		t.Fatalf("second call: %v", err)
	}
}

func TestMaybeMigrateNoLegacy(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("USERPROFILE", tmp)
	if err := MaybeMigrateLegacyDataDir(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(tmp, ".ginko")); err == nil {
		t.Fatal("expected ~/.ginko to NOT be created when no legacy dir exists")
	}
}
