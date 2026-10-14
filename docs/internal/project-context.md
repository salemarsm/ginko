# Project Context — Ginko

> Documento vivo. Base de conhecimento que todo agente recebe para entender o
> projeto. Atualizar ao fim de cada ciclo.

## 1. Identificação

- Repo: `github.com/salemarsm/ginko` (público, MIT)
- Branch default: `main`
- Última análise: 2026-10-14
- Stack: Go (binários + CLI + MCP), SQLite, HTML/JS embutido (GUI), LaTeX (whitepaper)
- Estado: v1.0 tagueado; dormente de 2026-05-06 a 2026-10-14

## 2. Estrutura

```
ginko/
├── cmd/
│   ├── ginko/          CLI principal (setup, integrate, ingest, serve, upgrade)
│   ├── ginko-admin/    administração (tokens, config)
│   ├── memctl/         CLI de memória (remember, context, suggest, search, timeline)
│   ├── memmcp/         servidor MCP (stdio, 12 tools)
│   ├── memserver/      servidor HTTP + GUI embutida
│   └── memdemo/        demo (sem testes)
├── memory/             store, schema, migrations, ranking, ingestão, eval, signals
├── server/             handlers HTTP, middleware de auth, GUI
├── config/             config, datadir, validação
├── internal/version/   versão injetada por ldflags
├── plugin/claude-code/ plugin do marketplace do [CC]
├── docs/               18 documentos + whitepaper (LaTeX/PDF) + openapi.yaml
├── e2e_test.go         teste end-to-end
├── Makefile            build / test / check / install / snapshot / e2e
└── .github/workflows/  ci.yml, release.yml, whitepaper.yml
```

## 3. Stack e comandos

- Gerenciador: `go mod` (`go.mod` declara `go 1.22`)
- Dependências diretas: `github.com/google/uuid`, `modernc.org/sqlite`
- **Sem CGO** (`modernc.org/sqlite` é Go puro) — binário estático, cross-compile trivial
- Build: `go build ./...` · Testes: `go test ./... -count=1` · Vet: `go vet ./...`
- Formatação: `gofmt -l .` (não está no `make check`)
- CI local: `make check` (= `go test ./...` + `go build ./cmd/...`)
- Whitepaper: `pdflatex -interaction=nonstopmode ginko-whitepaper.tex`, **duas passadas**

## 4. Arquitetura

Padrão: núcleo (store + retrieval) cercado de adaptadores de entrada (HTTP, MCP,
CLI) — hexagonal leve. O núcleo vive em `memory/`.

Persistência: SQLite com FTS5. Tabelas canônicas `memories`, `memory_tags`,
`memories_fts`, `events` (append-only), mais `documents`/`chunks`/`ingestion_runs`,
`sessions`, `agent_signals` e as auxiliares de embeddings/eval.

Comunicação: HTTP REST (`/api/...` e `/api/v1/...`), MCP stdio, CLI.
A GUI é embutida via `//go:embed` — **editar HTML/JS não tem efeito sem rebuild**.

Filosofia não negociável: SQLite é a fonte canônica da verdade. Embeddings são
índices opcionais, nunca verdade. Chunks são evidência, memória é conclusão.
A LLM é cliente, não dona do estado.

## 5. Convenções

- Go idiomático, `gofmt`; sem linter configurado
- Testes: `*_test.go` no pacote; smoke tests por comando + `e2e_test.go` na raiz
- Commits: Conventional Commits (`feat(v0.6): …`, `fix: …`, `docs: …`)
- Versionamento: semver por milestone; ldflags injetam versão/commit/data
- Não introduzir dependência com CGO; não tornar embedding obrigatório

## 6. Estado atual

- Tags: `v0.1.0-alpha.1`, `v0.1.0-alpha.2`, `v0.6`, `v0.7`, `v1.0` (sem tags para v0.2–v0.5)
- Backlog: 40+ itens P0–P3 em `docs/backlog.md`
- Issues no GitHub: 5 abertas (#5–#9), criadas em 2026-10-14 a partir do backlog.
  Antes disso o repositório tinha **zero issues na história** (os números #1–#4
  são PRs já mergeados) — o backlog vivia apenas como arquivo

## 7. Áreas sensíveis (hotspots)

- `memory/` — store, ranking e retrieval. Mudança de scoring afeta `/api/context`
  e pode quebrar `retrieval_eval_test.go`. **Serializar alterações neste pacote.**
- `server/` — handlers HTTP; contrato público estável desde o v1.0
- `docs/openapi.yaml` — precisa acompanhar qualquer endpoint novo ou alterado
- Sem testes: `internal/version`, `cmd/memdemo`

## 8. Riscos de onboarding

- **A documentação estava três releases atrás do código** (corrigido no PR #10).
  O whitepaper dizia que o último milestone era o v0.5 com v0.6/v0.7 "Planned";
  o README se declarava "Experimental v0.x". Um agente que lesse os docs primeiro
  construiria sobre premissas erradas. Revalidar sempre contra o código.
- **Fim de linha:** os blobs do repositório sempre foram LF. O que existe é ruído
  local — no Windows com `core.autocrlf=true` o checkout gera CRLF no working
  tree e `gofmt -l` passa a acusar todos os arquivos Go. O `.gitattributes` com
  `eol=lf` resolve. **Não "conserte" isso de novo achando que o repo está sujo.**
- **Duas falhas de teste são de ambiente Windows, não defeito:**
  `TestWriteConfigModeAndRoundTrip` (espera `0600`, Windows devolve `666`) e
  `TestIngestDoclingFileWithCLI` (Docling CLI ausente). O CI no Linux passa.
  Não "consertar" esses testes às cegas — rastreado na issue #9.
- **O CLI tem aliases legados** (`install-mcp` mapeia para `ginko-admin`).
- **`gofmt -l` é a única verificção de formatação** e não está no `make check`.
- A janela de log do GitHub Actions é de 90 dias: para run antiga não há como ler
  a causa de uma falha. Reproduzir localmente é o único caminho.

## 9. Comandos de verificação

```bash
go build ./...                    # compila tudo
go vet ./...                      # estático
gofmt -l .                        # deve retornar vazio
go test ./... -count=1            # 7 pacotes ok + 2 falhas de ambiente Windows
make e2e                          # teste end-to-end
```
