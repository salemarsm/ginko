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

## Retrieval and embeddings

RAG retrieval and memory retrieval may both use lexical and semantic indexes, but they remain separate layers.

- Memory retrieval returns canonical conclusions from `memories`.
- Chunk retrieval returns evidence from `chunks`.
- Embeddings over memories are optional retrieval indexes.
- Embeddings over chunks are optional evidence indexes.
- Vector search may help find candidates, but it never becomes the canonical store.

Planned auxiliary tables:

- `memory_embeddings`: derived semantic index rows for canonical memories
- `chunk_embeddings`: derived semantic index rows for evidence chunks
- `retrieval_eval_runs`: evaluation run metadata
- `retrieval_eval_items`: per-query expected/actual retrieval results

A citation-aware context builder should prefer compact canonical memories, with evidence links or evidence snippets only when useful and within `max_tokens`.
