# CLI

## llm-memory

```bash
bin/llm-memory init
bin/llm-memory doctor
bin/llm-memory paths
bin/llm-memory mcp-config
bin/llm-memory install-mcp claude-code
bin/llm-memory install-mcp codex
bin/llm-memory install-mcp openclaw
bin/llm-memory ui
```

## memctl

```bash
bin/memctl -subject botmaster search "direct answers"
bin/memctl -subject botmaster -max-tokens 800 context "How should I answer?"
bin/memctl -subject botmaster -type preference remember "User prefers direct answers."
bin/memctl -subject botmaster suggest "I prefer Go examples."
```
