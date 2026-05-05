package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultDataDir(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("USERPROFILE", tmp)
	want := filepath.Join(tmp, ".ginko")
	if got := DefaultDataDir(); got != want {
		t.Fatalf("DefaultDataDir()=%q want %q", got, want)
	}
}

func TestMaybeMigrateLegacyDataDirNoop(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("USERPROFILE", tmp)
	if err := MaybeMigrateLegacyDataDir(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(tmp, ".ginko")); err == nil {
		t.Fatal("expected MaybeMigrateLegacyDataDir to avoid creating data dir")
	}
}
