# Concepts

`ginko` is built around one distinction:

> RAG is evidence. Memory is conclusion.

## Layers

```txt
raw events  -> audit trail
documents   -> evidence
chunks      -> retrievable evidence
memories    -> canonical conclusions
context     -> compact prompt-ready projection
LLM         -> client
SQLite      -> source of truth
```

## Canonical memory

A canonical memory is a durable, structured record with provenance, confidence, scope, and lifecycle metadata.

It is not a vector. It is not a chat message. It is not an opaque prompt fragment.

## Evidence

Documents and chunks preserve where information came from. They are used for citations, verification, and RAG retrieval.

## Context projection

`/api/context` converts relevant memories into compact prompt-ready text under a token budget. Agents should prefer it over raw search.
