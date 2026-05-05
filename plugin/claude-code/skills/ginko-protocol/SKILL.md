---
description: Use this skill whenever working in a project that has a ginko (ginko) MCP server available — to recall prior context at the start of work, save durable memories after meaningful work, and supersede outdated memories. Trigger when the user mentions persistent memory, recall, decisions, preferences, or starts a task that may have prior context in this codebase.
---

# Ginko Memory Protocol

You have persistent memory across Claude Code sessions through the ginko MCP server. Memory is shared across all sessions in the same project.

## Available tools

- `memory_context` — recall relevant memories for the current task
- `memory_remember` — save a durable memory
- `memory_search` — explicit keyword search
- `memory_supersede` — replace an outdated memory with a new one (when supported by your installed version)

The `memory_context` tool is also called automatically by the SessionStart hook, so prior context is usually already injected when a session begins.

## When to save (proactive — don't wait to be asked)

Save after any of:

- **Bugfix.** Type `fact`. Content: root cause + how it was fixed + relevant file paths.
- **Architecture decision.** Type `decision`. Content: what was decided, why, what alternatives were rejected and why.
- **User preference.** Type `preference`. Content: stated or strongly inferred preference. Do not save preferences inferred from a single ambiguous turn.
- **Discovery.** Type `fact`. Content: non-obvious project knowledge that surprised you and would surprise the next agent.
- **Convention.** Type `preference` or `decision`. Content: project convention discovered (naming, layout, testing approach).

Skip saving for:

- Routine task completions
- Pure refactors with no learning
- Information already in `CLAUDE.md` or git history
- Trivial syntactic facts

## How to save well

A good memory is **actionable by a future session without further investigation**.

Bad: `memory_remember(content="fixed bug in auth")`

Good: `memory_remember(type="fact", subject="<project>", content="auth/middleware.ts: token refresh failed when iat was in the past — added 30s clock-skew tolerance via jose options.clockTolerance. Tests in test/auth.test.ts.", confidence=0.9)`

Rules:

- Lead with the file path or scope.
- State the problem and the fix in one sentence.
- Include test or verification location when applicable.
- Confidence: 0.9 default. Lower (0.6–0.8) when uncertain. Higher (0.95+) only for things you verified with multiple signals.

## When to recall

- **Beginning of meaningful work.** Already done by the SessionStart hook; you usually do not need to call `memory_context` again.
- **When the user references something you do not have in the current context.** "Remember when we...", "the way we did before", "our convention for X". Call `memory_search` with the relevant noun.
- **Before making a decision that may have precedent.** If you are about to choose between options, search for prior decisions on similar choices first.

Do not call `memory_search` on every turn. The token cost adds up. Search when there is a real reason.

## Supersession over duplication

If you find an existing memory that is now wrong or outdated, prefer `memory_supersede(old_id, new_content)` over creating a new memory. This keeps the audit trail intact.

If supersession is not available in your installed version, save a new memory and explicitly note `Supersedes: <old memory ID>` in the content.

## Privacy

Wrap sensitive content in `<private>` tags before passing it to memory tools. The server strips them before persistence:

```
memory_remember(content="API uses <private>sk-abc123</private> as the key")
```

Becomes stored as:

```
API uses [REDACTED] as the key
```

Never store live secrets, even briefly.

## Project scope

Memories are scoped per project. The project is auto-detected from the working directory (git remote URL or directory name). You do not need to pass a `project` field; if you do, it will be normalized.

## What memory is NOT

Memory is not a place to dump conversation logs, raw tool outputs, or evidence. It is for **conclusions** — durable claims worth re-applying in future sessions. Evidence (large file contents, document chunks) belongs to the RAG layer (when enabled) and is referenced from memories, not stored as memories.

A useful test before saving: *"if this exact memory surfaced in three months, would it still be useful and still be true?"* If no, do not save it.
