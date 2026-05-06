# MCP integration

`memmcp` exposes the Ginko memory layer over the Model Context Protocol (JSON-RPC stdio, protocol 2024-11-05).

## Setup

```bash
ginko setup claude-code   # writes memmcp entry to ~/.claude/settings.json
```

Manual `settings.json` entry:

```json
{
  "mcpServers": {
    "ginko": {
      "command": "ginko",
      "args": ["mcp"]
    }
  }
}
```

Or directly:

```json
{
  "mcpServers": {
    "ginko": {
      "command": "/path/to/memmcp",
      "args": ["-db", "/path/to/ginko.db"]
    }
  }
}
```

---

## Tool reference

### `memory_context`

Retrieve token-budgeted prompt context relevant to the current task.

| Field | Type | Description |
|-------|------|-------------|
| `query` | string | Current task or user prompt |
| `subject` | string | Project/subject identifier |
| `scopes` | string[] | `["project","global"]` etc. |
| `max_tokens` | integer | Token budget (default 1200) |

Returns: `{ context_id, context, items[], rankings{}, conflicts[], estimated_tokens, truncated }`.

Call this silently at the start of every session.

---

### `memory_suggest`

Extract durable memory candidates from the current conversation without persisting.

| Field | Type | Description |
|-------|------|-------------|
| `subject` | string | Project/subject identifier |
| `scope` | string | `project` or `global` |
| `user_prompt` | string | User's message |
| `assistant_response` | string | Assistant's response |
| `llm_inference` | string | Concise inference about durable learnings |
| `max_candidates` | integer | Max suggestions (default 5) |

Returns: `{ candidates: [{ memory, reason, requires_confirmation }] }`.

---

### `memory_remember`

Persist a durable memory.

| Field | Type | Description |
|-------|------|-------------|
| `type` | string | `fact` \| `decision` \| `preference` \| `task` \| `note` |
| `subject` | string | Project/subject |
| `content` | string | Memory content (max 250 chars recommended) |
| `source` | object | `{ kind, ref }` — provenance |
| `scope` | string | `global` \| `project` \| `session` |
| `confidence` | number | 0.0–1.0 |
| `tags` | string[] | Optional labels |
| `topic_key` | string | If set, auto-supersedes previous memory with same (subject, topic_key) |
| `dry_run` | boolean | Preview without persisting |

Returns: `{ memory, conflicts[], duplicates[] }`. Auto-creates a `conflict` signal when conflicts are detected.

---

### `memory_search`

Search raw memories by text across subjects.

| Field | Type | Description |
|-------|------|-------------|
| `text` | string | Search query |
| `subject` | string | Optional subject filter |
| `scopes` | string[] | Optional scope filter |
| `limit` | integer | Max results (default 20) |

Returns: array of `Memory` objects with `lexical_score`.

---

### `memory_get`

Fetch a single memory record by ID.

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Memory ID (`mem_...`) |

Returns: full `Memory` object.

---

### `memory_timeline`

Inspect the full lifecycle of a memory — supersessions, deletions, approvals.

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Memory ID |
| `limit` | integer | Max events (default 20) |

Returns: array of `Event` objects ordered newest-first.

---

### `memory_session_start`

Start a project session. Recovers the last closed session summary.

| Field | Type | Description |
|-------|------|-------------|
| `project` | string | Project identifier |

Returns: `{ session_id, project, started_at, last_summary? }`.

---

### `memory_session_end`

End the current session with a summary for future continuity.

| Field | Type | Description |
|-------|------|-------------|
| `project` | string | Project identifier |
| `summary` | string | What was done; what to continue next session |

Returns: closed session record.

---

### `memory_session_summary`

Retrieve the summary of the most recent closed session.

| Field | Type | Description |
|-------|------|-------------|
| `project` | string | Project identifier |

Returns: `{ summary, ended_at }` or `null` if no closed session.

---

### `signal_create`

Create a coordination signal visible to all agents sharing this database.

| Field | Type | Description |
|-------|------|-------------|
| `project` | string | Project identifier |
| `kind` | string | `notice` \| `lease` \| `handoff` \| `conflict` \| `review_request` \| `blocker` |
| `owner_agent` | string | Name of the creating agent (e.g. `claude-code`) |
| `target_agent` | string | Optional target agent |
| `payload` | string | Human-readable description |
| `topic_key` | string | Optional stable topic identifier |
| `expires_at` | string | RFC3339 — **required for `lease` signals** |

Returns: `AgentSignal` record.

---

### `signal_list`

List coordination signals for the current project.

| Field | Type | Description |
|-------|------|-------------|
| `project` | string | Project identifier |
| `kind` | string | Optional kind filter |
| `status` | string | Default `active`; use `*` for all |
| `agent` | string | Optional agent filter (owner or target) |
| `limit` | integer | Max results (default 100) |

Returns: array of `AgentSignal`.

---

### `signal_update`

Transition a signal to a new status.

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Signal ID (`sig_...`) |
| `status` | string | `acknowledged` \| `resolved` \| `cancelled` |

Returns: updated `AgentSignal`.

---

## Agent bootstrap prompt

Add to `CLAUDE.md` or system prompt:

```
Before answering, silently call memory_context with the user request, subject,
relevant scopes, and max_tokens <= 1200. Do not mention memory unless asked.

After answering, call memory_suggest with the user prompt, assistant response,
and a concise inference about durable learnings. Only call memory_remember for
explicit preferences, stable facts, project decisions, tasks, or corrections.
Ask before storing sensitive, private, or uncertain information.
Prefer compact memories over raw document chunks.
```

---

## MCP contract stability

From v1.0, the following tools are stable:

- `memory_context`, `memory_suggest`, `memory_remember`, `memory_search` — stable input/output schemas.
- `memory_get`, `memory_timeline` — stable.
- `memory_session_start`, `memory_session_end`, `memory_session_summary` — stable.
- `signal_create`, `signal_list`, `signal_update` — stable.

New optional fields may be added without a major version bump. Removing or renaming fields requires a major version.

See also: [Suggestion engine](suggestion-engine.md), [Claude Code integration](agents/claude-code.md).
