package main

import (
	"path/filepath"
	"testing"

	"github.com/salemarsm/ginko/config"
)

func TestResolveDBPathPrefersExplicit(t *testing.T) {
	const explicito = "/tmp/explicito.db"
	if got := resolveDBPath(explicito); got != explicito {
		t.Fatalf("got %q, esperado %q", got, explicito)
	}
}

// Sem -db e sem GINKO_DB, o caminho nao pode ser relativo.
//
// Um servidor MCP herda o diretorio de trabalho do cliente, entao "./memory.db"
// aponta para um lugar arbitrario e diferente do banco canonico que a GUI e o
// memctl usam. O agente grava memoria que o usuario nunca ve, sem erro nenhum.
func TestResolveDBPathNeverRelative(t *testing.T) {
	t.Setenv("USERPROFILE", t.TempDir())
	t.Setenv("HOME", t.TempDir())

	got := resolveDBPath("")
	if got == "./memory.db" || got == "memory.db" {
		t.Fatalf("resolveu para o caminho relativo %q", got)
	}
	if !filepath.IsAbs(got) {
		t.Fatalf("esperado caminho absoluto, got %q", got)
	}
	if got != config.DefaultDBPath() {
		t.Fatalf("got %q, esperado o padrao %q", got, config.DefaultDBPath())
	}
}
