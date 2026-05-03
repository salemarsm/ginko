package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	Server    ServerConfig    `json:"server"`
	Database  DatabaseConfig  `json:"database"`
	LLM       LLMConfig       `json:"llm"`
	Embedding EmbeddingConfig `json:"embedding"`
}

type ServerConfig struct {
	Addr string `json:"addr"`
}

type DatabaseConfig struct {
	Path string `json:"path"`
}

type LLMConfig struct {
	Provider  string `json:"provider"`
	Model     string `json:"model"`
	APIKeyEnv string `json:"api_key_env"`
}

type EmbeddingConfig struct {
	Provider  string `json:"provider"`
	Model     string `json:"model"`
	Index     string `json:"index"`
	APIKeyEnv string `json:"api_key_env"`
}

func Default() Config {
	return Config{
		Server:   ServerConfig{Addr: "127.0.0.1:8787"},
		Database: DatabaseConfig{Path: "./memory.db"},
		LLM: LLMConfig{
			Provider:  "none",
			Model:     "",
			APIKeyEnv: "",
		},
		Embedding: EmbeddingConfig{
			Provider:  "none",
			Model:     "",
			Index:     "sqlite-fts",
			APIKeyEnv: "",
		},
	}
}

func Load(path string) (Config, error) {
	cfg := Default()
	if path == "" {
		return cfg, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	if err := json.Unmarshal(b, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, cfg.Validate()
}

func (c Config) Validate() error {
	if c.Server.Addr == "" {
		return errors.New("server.addr is required")
	}
	if c.Database.Path == "" {
		return errors.New("database.path is required")
	}
	return nil
}

func WriteDefault(path string) error {
	cfg := Default()
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return os.WriteFile(path, b, 0o644)
}
