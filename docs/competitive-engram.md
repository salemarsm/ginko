# Competitive analysis: Engram

Engram is the closest known direct comparator for `llm-memory`: Go, SQLite, FTS5, MCP, HTTP API, CLI/TUI, and a strong focus on coding-agent memory.

This note extracts practices worth adopting while keeping `llm-memory` differentiated as a canonical, auditable memory database rather than only a coding-agent observation log.

## What Engram does well

### 1. Crisp distribution story

Engram's onboarding promise is simple: one binary, one SQLite file, no runtime dependencies.

Practices to adopt:

- Ship one primary binary for normal users.
- Add GitHub Actions and GoReleaser early.
- Publish release archives for Linux/macOS/Windows with checksums.
- Document `go install` as the trust-preserving path for technical users.
- Add Homebrew-ready packaging later.

`llm-memory` angle:

> One local canonical memory database. SQLite source of truth. HTTP/MCP/CLI over the same lifecycle model.

### 2. Agent setup is product, not just docs

Engram provides per-agent setup docs and `engram setup <agent>` commands. It also explains compaction survival and project auto-detection.

Practices to adopt:

- Implement `llm-memory integrate <agent>` or `llm-memory setup <agent>` for at least OpenClaw, Claude Code, Codex-like CLI, and generic MCP.
- Keep manual JSON examples, but make the happy path one command where safe.
- Add bootstrap prompts / agent instructions that encode memory usage policy.
- Add smoke tests for generated config snippets where possible.

`llm-memory` differentiation:

- Setup should emphasize transparent canonical memory flow:
  1. context before answering,
  2. suggest after meaningful work,
  3. approve/store/supersede with audit trail.

### 3. Progressive disclosure retrieval

Engram documents a useful three-layer pattern: search compact results, then timeline, then full observation.

Practices to adopt:

- Keep `/api/context` compact and token-budgeted.
- Add explicit drill-down endpoints/workflows:
  - compact memory search,
  - memory detail,
  - event/supersession timeline,
  - evidence/document chunks.
- Make MCP tools mirror that flow.

`llm-memory` stronger version:

```txt
memory_context        -> prompt-ready compact canonical memory
memory_search         -> compact candidates with IDs
memory_get            -> full memory with provenance/supersession
memory_timeline       -> lifecycle/audit trail
memory_evidence       -> supporting document chunks/events
```

### 4. Memory hygiene primitives

Engram has pragmatic hygiene features: topic keys, revision counts, duplicate counts, soft delete, dedupe hashes, and project/scope filters.

Practices to adopt/adapt:

- Add a stable `topic_key` concept for evolving subjects/decisions.
- Add duplicate detection metadata instead of blindly inserting repeated memories.
- Add revision counters for updates/supersession chains.
- Keep soft-delete semantics visible in API/GUI.
- Provide a helper to suggest canonical topic keys.

`llm-memory` adaptation:

- Do not collapse everything into mutable observations. Preserve the stronger canonical lifecycle:
  - new memory,
  - superseded memory,
  - replacement memory,
  - event log,
  - provenance/evidence.

### 5. Project detection and ambiguity handling

Engram treats project identity as a first-class operational issue. It detects from cwd/git/config and refuses to guess when ambiguous.

Practices to adopt:

- Add project/subject detection helpers for agent integrations.
- Return structured ambiguity errors instead of silently writing to the wrong subject/project.
- Support repo-local config, e.g. `.llm-memory/config.json`, for canonical project identity.
- Add merge/consolidate tools for similar project names.

`llm-memory` angle:

- Generalize this beyond code projects: `subject` identity should be explicit, validated, and auditable.
- For coding agents, project should likely be a subject dimension or metadata field, not a hidden global.

### 6. Diagnostics and repair flows

Engram has a visible `doctor` posture, diagnostic packages, repair flows, and tests around migrations/legacy schemas.

Practices to adopt:

- Expand `llm-memory doctor` into read-only checks plus explicit repair commands.
- Diagnose DB schema version, migrations, FTS health, config, server reachability, auth exposure, and MCP config.
- Never auto-repair destructive issues without confirmation.
- Add migration/legacy-schema tests.

### 7. Conflict surfacing

Engram's newer work includes memory relations and conflict/judge flows.

Practices to adopt:

- Add explicit conflict candidates to `memory_suggest` results.
- Represent conflicts/supersession/reinforcement as first-class relations.
- Keep judge decisions auditable: model/actor, reason, evidence, confidence.

`llm-memory` stronger version:

- Treat conflict surfacing as governance over canonical memory, not just search-time warnings.

### 8. Sync model is local-first replication

Engram's cloud/git sync message is carefully framed: local SQLite remains authoritative; cloud is opt-in replication.

Practices to adopt later:

- If sync is added, keep local canonical DB primary.
- Make replication opt-in and project/subject-scoped.
- Prefer append/export chunks or migration-safe bundles over raw DB sync.
- Include conflict handling before multi-writer sync.

## What not to copy blindly

- Do not become only a coding-agent observation tracker. `llm-memory` should stay broader: personal AI infrastructure, RAG, assistants, and coding agents.
- Do not over-expand MCP tools before the lifecycle is stable. Fewer canonical tools are easier to integrate correctly.
- Do not introduce cloud before local safety/auth/governance are solid.
- Do not make raw prompt/session capture the default source of truth. It can be evidence, not canonical memory.

## Differentiation to reinforce

Engram's strongest message is "persistent memory for AI coding agents." `llm-memory` should avoid competing as a clone.

Recommended positioning:

> `llm-memory` is a local-first canonical memory database for AI agents: SQLite source of truth, HTTP/MCP interface, auditable lifecycle, token-budgeted context, and evidence-aware RAG bridge.

Short contrast:

```txt
Engram      = coding-agent observations and session memory
llm-memory  = canonical operational memory with evidence, audit, supersession, and agent-agnostic retrieval
```

Core differentiators to keep visible:

- canonical memory vs raw observation log,
- provenance/evidence links,
- supersession lifecycle,
- append-only events,
- RAG thesis: evidence vs conclusion,
- LLM as client, not database,
- embeddings as replaceable indexes,
- local-first governance before cloud.

## Implementation priorities derived from this review

1. **Packaging/CI**: GitHub Actions + GoReleaser + version command.
2. **Agent setup UX**: one-command setup for key agents plus tested manual snippets.
3. **Memory hygiene**: `topic_key`, duplicate/revision metadata, better soft-delete docs.
4. **Progressive disclosure**: add detail/timeline/evidence MCP/API flow.
5. **Project/subject identity**: detection, ambiguity errors, consolidation tools.
6. **Doctor/repair**: actionable diagnostics and explicit safe repair commands.
7. **Conflict governance**: relation table and suggestion-time conflict surfacing.

## Near-term recommended sequence

```txt
v0.1 polish: CI, version, release config, README install path
v0.2: auth + doctor + setup skeleton
v0.3: agent integrations + progressive disclosure MCP tools
v0.4: topic keys + duplicate/revision + conflict relations + approval queue
v0.5+: evidence/RAG ingestion and citation-aware context
```
