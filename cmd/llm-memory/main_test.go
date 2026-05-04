package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/salemarsm/llm-memory/config"
)

func TestTokenCreateHelpDoesNotMutateConfig(t *testing.T) {
	home := t.TempDir()
	if err := tokenCommand(home, []string{"create", "--help"}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(configPath(home)); !os.IsNotExist(err) {
		t.Fatalf("token create --help must not create config, stat err=%v", err)
	}
}

func TestTokenTopLevelHelpDoesNotMutateConfig(t *testing.T) {
	home := t.TempDir()
	if err := tokenCommand(home, []string{"--help"}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(configPath(home)); !os.IsNotExist(err) {
		t.Fatalf("token --help must not create config, stat err=%v", err)
	}
}

func TestSetupClaudeCodeHelpDoesNotMutateConfig(t *testing.T) {
	home := t.TempDir()
	if err := setupCommand(home, []string{"claude-code", "--help"}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(configPath(home)); !os.IsNotExist(err) {
		t.Fatalf("setup claude-code --help must not create config, stat err=%v", err)
	}
}

func TestSetupTopLevelHelpDoesNotMutateConfig(t *testing.T) {
	home := t.TempDir()
	if err := setupCommand(home, []string{"--help"}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(configPath(home)); !os.IsNotExist(err) {
		t.Fatalf("setup --help must not create config, stat err=%v", err)
	}
}

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

func TestMergeClaudeSettings(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	if err := os.WriteFile(path, []byte(`{"permissions":{"allow":["Bash(git status)"]},"mcpServers":{"other":{"command":"other"}}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	original, merged, err := mergeClaudeSettings(path, map[string]any{"command": "/bin/memmcp", "args": []string{"-db", "/tmp/memory.db"}})
	if err != nil {
		t.Fatal(err)
	}
	if original == nil {
		t.Fatal("expected original bytes for existing settings")
	}
	var out map[string]any
	if err := json.Unmarshal(merged, &out); err != nil {
		t.Fatal(err)
	}
	servers := out["mcpServers"].(map[string]any)
	if _, ok := servers["other"]; !ok {
		t.Fatal("existing MCP server was not preserved")
	}
	ginko := servers["ginko"].(map[string]any)
	if ginko["command"] != "/bin/memmcp" {
		t.Fatalf("unexpected ginko command: %#v", ginko)
	}
	if _, ok := out["permissions"]; !ok {
		t.Fatal("unrelated settings were not preserved")
	}
}
