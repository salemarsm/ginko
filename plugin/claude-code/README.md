# Ginko for Claude Code

Ginko is the Claude Code plugin surface for `llm-memory`: persistent canonical memory backed by SQLite and exposed through MCP.

Current contents:

- MCP server declaration for `memmcp`
- SessionStart / PreCompact / Stop hook skeletons
- Ginko memory protocol skill
- `/save` and `/recall` command prompts

The plugin name is intentionally short and memorable; the implementation remains in the `llm-memory` repository.

## Local setup today

Until marketplace packaging is finalized, use:

```sh
llm-memory setup claude-code --dry-run
llm-memory setup claude-code
```

The setup command merges a `ginko` MCP server into Claude Code settings without overwriting unrelated settings and writes a backup first.
