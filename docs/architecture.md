# Architecture

```txt
┌──────────────────────────┐
│ OpenClaw / Claude / Codex │
└────────────┬─────────────┘
             │ MCP / HTTP / CLI
┌────────────▼─────────────┐
│      Ginko API       │
├──────────────────────────┤
│ Context builder           │
│ Suggestion engine         │
│ Write policy              │
├──────────────────────────┤
│ Canonical memories        │
│ Append-only events        │
│ Documents / chunks        │
├──────────────────────────┤
│ SQLite + FTS5             │
│ Optional vector indexes   │
└──────────────────────────┘
```

## Operational flow

```txt
User prompt
  ↓
Agent calls memory_context
  ↓
ginko returns compact relevant memories
  ↓
Agent answers
  ↓
Agent calls memory_suggest
  ↓
Human/policy approves
  ↓
memory_remember stores or supersedes
  ↓
Event log records changes
```
