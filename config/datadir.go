package config

import (
	"os"
	"path/filepath"
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

// MaybeMigrateLegacyDataDir is retained as a no-op compatibility hook for
// older binaries that called it during startup.
func MaybeMigrateLegacyDataDir() error { return nil }
