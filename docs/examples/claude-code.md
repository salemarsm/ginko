# Claude Code setup

Ginko is the Claude Code-facing plugin/setup name for `ginko`.

Dry-run first:

```bash
bin/ginko setup claude-code --dry-run
```

Apply:

```bash
bin/ginko setup claude-code
```

What it does:

- finds Claude Code settings (`.claude/settings.json` in the current project if present, otherwise `~/.claude/settings.json`)
- merges a `ginko` MCP server entry without overwriting unrelated settings
- writes a `.bak` backup before modifying an existing file
- points Claude Code at the local `memmcp` binary

Legacy snippet-only flow remains available:

```bash
bin/ginko install-mcp claude-code
```

See `plugin/claude-code/` for the Ginko plugin skeleton: MCP declaration, hooks, skill, and slash-command prompts.
