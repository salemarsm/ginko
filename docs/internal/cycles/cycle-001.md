# Ciclo 1 — Orientação e observações

- **Data:** 2026-10-14
- **Atividade:** onboarding + primeiro ciclo OODA no projeto
- **Decisão resultante:** [`../decisions/decision-001-repo-truth.md`](../decisions/decision-001-repo-truth.md)

## ORIENT — onde o projeto estava

- `main` em `90f7b49`, tags `v0.1.0-alpha.1`, `v0.1.0-alpha.2`, `v0.6`, `v0.7`, `v1.0`
- Último push em 2026-05-06; retomada em 2026-10-14 (5 meses)
- Backlog de 40+ itens em `docs/backlog.md`, **zero issues no GitHub**
- O último commit é `feat(v1.0): stable API contract, scope ranking, full eval fixtures, docs`

## OBSERVE — o que as medições disseram

Tudo abaixo foi medido com comando, não lido de relatório.

### Saudável

```
go build ./...      OK (2,5s)
go vet ./...        limpo
go test ./...       7 pacotes ok, 2 falhas de ambiente Windows
CI (Linux)          "Test and build" = success
```

As duas falhas locais — `TestWriteConfigModeAndRoundTrip` (espera `0600`, Windows
devolve `666`) e `TestIngestDoclingFileWithCLI` (Docling CLI ausente) — não são
defeitos do produto; o CI verde no Linux confirma. Rastreadas na issue #9.

### Defeitos encontrados

1. **O whitepaper descrevia o v0.5 como último milestone**, com v0.6 e v0.7
   marcados como *Planned* — três releases atrás do código.
2. **O PDF commitado era mais antigo que sua fonte**: PDF em `b3b3962`, `.tex`
   alterado depois em `c73af39`.
3. **O workflow `Compile Whitepaper` falhava em todo push desde o v0.6.**
4. **Nove arquivos Go realmente fora do `gofmt`** (alinhamento de comentários,
   uma linha em branco final).
5. **Deriva de documentação** em README, `docs/roadmap.md`, `docs/install.md` e
   CHANGELOG (contagem de tools MCP: 13 declarado, 12 real).

## Duas hipóteses erradas, registradas aqui

Ficam no registro porque o ciclo só produziu resultado confiável depois de
derrubá-las.

### Causa raiz errada do CI

**Hipótese:** faltava `permissions: contents: write` no `whitepaper.yml`, e o
`git-auto-commit-action` tomava 403. Sustentada por evidência circunstancial — o
`release.yml` declara `contents: write`, o `whitepaper.yml` não declarava nada, e
a `main` não é protegida.

**Teste:** adicionada a permissão e feito o push. **O workflow falhou de novo.**

**Causa real**, lida no log fresco:

```
! LaTeX Error: File `lmodern.sty' not found.
l.8 \usepackage{microtype}
!  ==> Fatal error occurred, no output PDF file produced!
```

Falha no passo *Compile PDF*, não no commit: o documento exige
`\usepackage{lmodern}` e o `apt-get` do workflow não instalava esse pacote
(separado dos `texlive-latex-*` em Debian/Ubuntu).

**Lição:** a hipótese era plausível e bem sustentada, e ainda assim errada. O log
original tinha expirado (HTTP 410 — o GitHub purga após 90 dias), então a única
forma de testar era provocar uma execução nova. Quando a causa raiz de uma falha
não é observável, **provocar o erro de novo vale mais do que inferir**.

### Medição errada de fim de linha

**Hipótese:** o repositório tinha CRLF gravado no índice git, e por isso
`gofmt -l` acusava 61 arquivos. Baseada em `grep -c $'\r'`, que devolveu
"353 linhas com CR".

**Medição correta**, em bytes:

```
cmd/ginko/main.go   working tree: 8749 bytes, 0 CRLF
                    blob HEAD:    8749 bytes, 0 CRLF   (mesmo hash de blob)
```

Os blobs sempre foram LF. O `grep` casou a **letra "r"**, não o byte CR — o
arquivo não tinha nenhum CR. O CRLF existia apenas no working tree, gerado pelo
`core.autocrlf=true` no checkout do Windows. O commit que eu descrevi como
"normalizar 115 arquivos" não tocou um único arquivo de código.

**Lição:** medir bytes com `grep`/`awk` neste shell não é confiável. Contar bytes
com Python é. E uma medição que confirma uma hipótese não substitui uma segunda
medição independente.

## Resultado

PR #10 na branch `chore/repo-truth-v1`. O workflow do whitepaper passou
(run `38596773277`) — verde pela primeira vez desde o v0.6. Cinco issues criadas
(#5–#9). Zero regressão no build ou nos testes.
