# RAG pipeline

> RAG is evidence. Memory is conclusion.

Recommended pipeline:

```txt
PDF/DOCX/HTML
  -> Docling
  -> Markdown/JSON
  -> chunks
  -> document evidence
  -> extracted canonical memory candidates
```

## Separation

- `documents`: source metadata and file hash
- `chunks`: retrievable evidence with headings/pages/token counts
- `memories`: durable conclusions extracted from evidence
- `embedding_refs`: optional vector index pointers

## Example

Document chunk:

```txt
In meeting notes, the team decided to use SQLite for local-first storage.
```

Canonical memory:

```txt
The llm-memory project uses SQLite as the canonical local-first store.
```

Evidence link:

```txt
doc_id=..., chunk_id=...
```
