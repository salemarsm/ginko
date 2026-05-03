# LLM-Agnostic Memory Store

Implementação em Go + SQLite de uma camada de memória independente de LLM, com API HTTP e GUI web local.

## Ideia

O LLM não é o dono da memória. Ele só chama uma API sobre uma fonte canônica externa.

Componentes:

- `events`: log append-only de acontecimentos brutos.
- `memories`: memória canônica estruturada.
- `memories_fts`: índice SQLite FTS5 para busca textual.
- `embedding_refs_json`: ponteiros opcionais para índices vetoriais externos.
- `server`: API HTTP + GUI web embutida.
- `config`: configuração de DB, servidor, LLM e embeddings.

Embeddings não são a fonte da verdade. Podem ser apagados e regenerados.

## Rodar GUI/API

Requer Go 1.22+.

```bash
go test ./...
go run ./cmd/memserver -config ./config.example.json
```

Abrir:

```txt
http://127.0.0.1:8787
```

Gerar config default:

```bash
go run ./cmd/memserver -write-config ./config.local.json
```

## Configuração de LLM

A config existe para declarar qual LLM/agente vai consumir a memória. Ela não muda o formato canônico da memória.

```json
{
  "server": { "addr": "127.0.0.1:8787" },
  "database": { "path": "./memory.db" },
  "llm": {
    "provider": "openai",
    "model": "gpt-5.5",
    "api_key_env": "OPENAI_API_KEY"
  },
  "embedding": {
    "provider": "openai",
    "model": "text-embedding-3-small",
    "index": "sqlite-fts",
    "api_key_env": "OPENAI_API_KEY"
  }
}
```

Nota: no estado atual, `llm` e `embedding` são metadados/configuração para integração. O servidor ainda não chama APIs externas de LLM; isso preserva o núcleo LLM-agnóstico.

## API HTTP

### Listar/buscar memórias

```bash
curl 'http://127.0.0.1:8787/api/memories?q=respostas&subject=botmaster'
```

### Busca POST

```bash
curl -X POST http://127.0.0.1:8787/api/search \
  -H 'content-type: application/json' \
  -d '{"text":"respostas diretas","subject":"botmaster","limit":10}'
```

### Criar/atualizar memória

```bash
curl -X POST http://127.0.0.1:8787/api/memories \
  -H 'content-type: application/json' \
  -d '{
    "type":"preference",
    "subject":"botmaster",
    "content":"Prefere respostas diretas, técnicas e sem enrolação.",
    "source":{"kind":"api","ref":"manual"},
    "scope":"global",
    "confidence":0.95,
    "tags":["style","preference"],
    "embedding_refs":{}
  }'
```

### Buscar por ID

```bash
curl http://127.0.0.1:8787/api/memories/mem_123
```

### Supersede

```bash
curl -X POST http://127.0.0.1:8787/api/supersede/mem_antiga \
  -H 'content-type: application/json' \
  -d '{"type":"preference","subject":"botmaster","content":"Nova preferência.","source":{"kind":"api","ref":"manual"},"scope":"global","confidence":0.9}'
```

### Forget/delete

```bash
curl -X DELETE http://127.0.0.1:8787/api/memories/mem_123
```

### Eventos

```bash
curl http://127.0.0.1:8787/api/events?limit=50
```

### Config efetiva

```bash
curl http://127.0.0.1:8787/api/config
```

## API Go

```go
store, _ := memory.Open("memory.db")
defer store.Close()

m, _ := store.UpsertMemory(ctx, memory.Memory{
    Type:       memory.TypePreference,
    Subject:    "botmaster",
    Content:    "Prefere respostas diretas, técnicas e sem enrolação.",
    Source:     memory.Source{Kind: "conversation", Ref: "msg-123"},
    Scope:      memory.ScopeGlobal,
    Confidence: 0.95,
    Tags:       []string{"style"},
})

items, _ := store.Search(ctx, memory.Query{
    Text:    "respostas diretas",
    Subject: "botmaster",
    Limit:   10,
})

_ = m
_ = items
```

## Schema conceitual

```json
{
  "id": "mem_...",
  "type": "preference | fact | decision | task | note | relationship",
  "subject": "botmaster",
  "content": "Prefere respostas diretas, técnicas e sem enrolação.",
  "source": { "kind": "conversation", "ref": "session/message" },
  "scope": "global | project | session | private",
  "confidence": 0.95,
  "created_at": "...",
  "updated_at": "...",
  "valid_from": null,
  "valid_until": null,
  "supersedes_id": null,
  "superseded_by": null,
  "tags": ["style"],
  "embedding_refs": { "default": "vec_..." }
}
```

## Operações implementadas

- `Open(path)` cria/conecta e roda migrações.
- `AppendEvent(ctx, event)` grava histórico bruto.
- `ListEvents(ctx, limit)` lista eventos recentes.
- `UpsertMemory(ctx, memory)` cria/atualiza memória canônica.
- `GetMemory(ctx, id)` busca por ID.
- `Search(ctx, query)` busca por texto, tipo, escopo, assunto e tags.
- `Supersede(ctx, oldID, newer)` substitui memória antiga sem apagar histórico.
- `Forget(ctx, id)` remove memória canônica e índice FTS.

## Decisões de design

- SQLite local, simples, auditável.
- `modernc.org/sqlite` para evitar CGO.
- FTS5 para busca textual sem depender de embeddings.
- JSON em `tags_json` e `embedding_refs_json` para manter flexibilidade.
- `superseded_by IS NULL` por padrão em `Search`, então resultados ativos não trazem memórias obsoletas.
- GUI sem framework pesado: HTML/CSS/JS embutido no binário Go.
- Configuração de LLM desacoplada do store para manter portabilidade entre modelos.

## Próximos passos naturais

1. Autenticação local/API token.
2. Ranking híbrido: FTS + recência + confiança + embedding opcional.
3. Compactador: eventos brutos → memórias consolidadas.
4. Cliente SDK para agentes/LLMs.
5. Soft-delete auditável separado de hard-delete para privacidade.
