package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// DefaultDataDir returns ~/.ginko. Falls back to .ginko if HOME is unavailable.
func DefaultDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ".ginko"
	}
	return filepath.Join(home, ".ginko")
}

// DefaultDBPath returns ~/.ginko/ginko.db.
func DefaultDBPath() string {
	return filepath.Join(DefaultDataDir(), "ginko.db")
}

// DefaultConfigPath returns ~/.ginko/config.json.
func DefaultConfigPath() string {
	return filepath.Join(DefaultDataDir(), "config.json")
}

// MaybeMigrateLegacyDataDir performs a one-time migration from ~/.llm-memory/
// to ~/.ginko/ when the new dir does not exist and the legacy dir does.
// It is idempotent and silent when not needed.
func MaybeMigrateLegacyDataDir() error {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return nil
	}
	newDir := filepath.Join(home, ".ginko")
	legacyDir := filepath.Join(home, ".llm-memory")

	if _, err := os.Stat(newDir); err == nil {
		return nil
	}
	if _, err := os.Stat(legacyDir); err != nil {
		return nil
	}

	if err := os.MkdirAll(newDir, 0o755); err != nil {
		return fmt.Errorf("create %s: %w", newDir, err)
	}

	pairs := []struct{ from, to string }{
		{filepath.Join(legacyDir, "memory.db"), filepath.Join(newDir, "ginko.db")},
		{filepath.Join(legacyDir, "memory.db-wal"), filepath.Join(newDir, "ginko.db-wal")},
		{filepath.Join(legacyDir, "memory.db-shm"), filepath.Join(newDir, "ginko.db-shm")},
		{filepath.Join(legacyDir, "config.json"), filepath.Join(newDir, "config.json")},
	}
	migrated := 0
	for _, p := range pairs {
		if _, err := os.Stat(p.from); err != nil {
			continue
		}
		if err := copyFile(p.from, p.to); err != nil {
			return fmt.Errorf("copy %s: %w", p.from, err)
		}
		migrated++
	}
	if migrated > 0 {
		stamp := []byte("migrated_at=" + time.Now().UTC().Format(time.RFC3339) + "\nfrom=" + legacyDir + "\n")
		_ = os.WriteFile(filepath.Join(newDir, "MIGRATED_FROM_LLM_MEMORY"), stamp, 0o644)
		fmt.Fprintf(os.Stderr, "ginko: migrated %d file(s) from %s to %s. Old files preserved.\n", migrated, legacyDir, newDir)
	}
	return nil
}

func copyFile(from, to string) error {
	src, err := os.Open(from)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.OpenFile(to, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}
