---
description: Use Ginko persistent memory when starting work, recalling project decisions, saving durable coding knowledge, or replacing outdated memories.
---

# Ginko Memory Protocol

Ginko gives Claude Code persistent memory across sessions through llm-memory MCP tools.

Use memory as durable project knowledge, not as raw chat history.

## Recall

- At the start of meaningful work, use the provided session context first.
- If the user asks “remember”, “what did we decide”, or “how did we solve this before”, call `memory_search` or `memory_context`.
- Before repeating architecture/debug work, search for prior decisions and bugfixes.

## Save proactively

Save a memory after:

- architecture decisions
- root-cause discoveries
- non-obvious bugfixes
- recurring user/project preferences
- constraints that future sessions must respect

Prefer:

- `type=decision` for what/why architectural choices
- `type=fact` for durable project facts and bugfix lessons
- `type=preference` for user/project style preferences

Skip routine progress, transient todos, secrets, and anything already obvious in repository docs.

## Good memory shape

Bad:

```txt
fixed auth bug
```

Good:

```txt
auth/middleware.go: token refresh failed when iat was slightly in the future; added 30s clock-skew tolerance after reproducing with CI clock drift.
```

A future session should be able to act without rediscovering the whole context.

## Supersession

If an old memory is obsolete, supersede it instead of duplicating it. Keep history auditable.

## Privacy

Never save content wrapped in `<private>...</private>`. Ask before saving sensitive personal or credential-like data.
