# Memory model

## Types

| Type | Use for | Example |
|---|---|---|
| `preference` | stable user/project preferences | `User prefers Go examples over Python examples.` |
| `fact` | durable factual statements | `The project uses SQLite as the canonical store.` |
| `decision` | architecture/project decisions | `Embeddings are optional indexes, not canonical truth.` |
| `task` | long-lived pending actions | `Add API token support.` |
| `note` | low-structure observations | `The GUI is currently local-only.` |
| `relationship` | links between entities | `Project X belongs to client Y.` |

## Schema

```json
{
  "id": "mem_...",
  "type": "preference",
  "subject": "botmaster",
  "content": "User prefers direct technical answers.",
  "source": { "kind": "conversation", "ref": "session/message" },
  "scope": "global",
  "confidence": 0.95,
  "tags": ["style"],
  "embedding_refs": {}
}
```

## Supersession

Use supersession to replace stale memory without pretending history never existed.

## Retrieval metadata is derived

Canonical memories are the durable records in the `memories` table. Retrieval scores are derived projections over those records and must not become the source of truth.

Planned ranking metadata may include:

- `lexical_score`
- `semantic_score`
- `recency_score`
- `confidence_score`
- `provenance_score`
- `final_score`
- `rank_reason`

These values explain a retrieval decision at a point in time. They may be logged in evaluation tables, but they should not replace canonical memory fields.

## Embeddings and canonical memory

Embeddings are optional indexes. They are useful for semantic candidate generation, but they are not memory.

Guidelines:

- Keep SQLite as the canonical source of truth.
- Keep `memories` as the canonical conclusion table.
- Preserve `embedding_refs`/`embedding_refs_json` as adapter bridge metadata when useful.
- Add auxiliary tables such as `memory_embeddings` only as derived indexes.
- Do not require a vector database for normal operation.
- Do not treat documents or chunks as canonical memories.
