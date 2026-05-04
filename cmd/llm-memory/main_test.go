package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/salemarsm/llm-memory/config"
)

func TestTokenCreateListRevoke(t *testing.T) {
	home := t.TempDir()
	if err := tokenCommand(home, []string{"create"}); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(configPath(home))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.AuthToken == "" || len(cfg.Server.AuthToken) < 32 {
		t.Fatalf("expected generated auth token, got %q", cfg.Server.AuthToken)
	}
	if err := tokenCommand(home, []string{"list"}); err != nil {
		t.Fatal(err)
	}
	if err := tokenCommand(home, []string{"revoke"}); err != nil {
		t.Fatal(err)
	}
	cfg, err = config.Load(configPath(home))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.AuthToken != "" || cfg.Server.AuthTokenEnv != "" {
		t.Fatalf("expected auth config cleared: %#v", cfg.Server)
	}
}

func TestDoctorHandlesBadConfig(t *testing.T) {
	home := t.TempDir()
	if err := os.WriteFile(filepath.Join(home, "config.json"), []byte(`{"server":`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := doctor(home); err != nil {
		t.Fatal(err)
	}
}

func TestWriteConfigModeAndRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	cfg := config.Default()
	cfg.Server.AuthToken = strings.Repeat("x", 32)
	if err := writeConfig(path, cfg); err != nil {
		t.Fatal(err)
	}
	st, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if st.Mode().Perm() != 0o600 {
		t.Fatalf("expected 0600 config permissions, got %o", st.Mode().Perm())
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var decoded config.Config
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Server.AuthToken != cfg.Server.AuthToken {
		t.Fatal("auth token did not round trip")
	}
}
