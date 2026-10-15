package memory

import (
	"os"
	"path/filepath"
	"testing"
)

func escreverProjectConfig(t *testing.T, dir, conteudo string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, ".ginko"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".ginko", "config.json"), []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
}

// GINKO_PROJECT fixa a identidade de projeto explicitamente.
//
// Um servidor MCP herda o diretorio de trabalho do cliente, nao o do repositorio
// que ele serve. Sem um override explicito, a deteccao por CWD resolve para o
// basename da pasta do usuario (ex.: "salemarsm") e toda consulta devolve vazio
// em silencio -- sem erro, sem aviso.
func TestDetectProjectEnvOverride(t *testing.T) {
	t.Setenv("GINKO_PROJECT", "ginko")
	if got := DetectProject(t.TempDir()); got != "ginko" {
		t.Fatalf("GINKO_PROJECT deveria vencer a deteccao por CWD; got %q", got)
	}
}

func TestDetectProjectEnvOverrideIsNormalized(t *testing.T) {
	t.Setenv("GINKO_PROJECT", "  My Cool_Project  ")
	if got := DetectProject(t.TempDir()); got != "my-cool-project" {
		t.Fatalf("got %q, esperado my-cool-project", got)
	}
}

func TestDetectProjectEnvOverrideBeatsRepoConfig(t *testing.T) {
	dir := t.TempDir()
	escreverProjectConfig(t, dir, `{"project":"do-config"}`)
	t.Setenv("GINKO_PROJECT", "do-env")
	if got := DetectProject(dir); got != "do-env" {
		t.Fatalf("got %q, esperado do-env", got)
	}
}

func TestDetectProjectEmptyEnvKeepsDetection(t *testing.T) {
	t.Setenv("GINKO_PROJECT", "")
	dir := t.TempDir()
	escreverProjectConfig(t, dir, `{"project":"do-config"}`)
	if got := DetectProject(dir); got != "do-config" {
		t.Fatalf("got %q, esperado do-config: env vazio nao pode desligar a deteccao", got)
	}
}
