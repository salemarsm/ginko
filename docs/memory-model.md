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
