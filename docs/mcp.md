# MCP integration

`memmcp` exposes:

- `memory_context`
- `memory_suggest`
- `memory_remember`
- `memory_search`

## Build

```bash
go build -o bin/memmcp ./cmd/memmcp
```

## Config snippet

```json
{
  "command": "/path/to/llm-memory/bin/memmcp",
  "args": ["-db", "/path/to/memory.db"]
}
```

## Agent bootstrap

```txt
Before answering, silently call memory_context with the user request, subject, relevant scopes, and max_tokens <= 1200.
Do not mention memory unless asked.
After answering, call memory_suggest with the user prompt, assistant response, and a concise LLM inference about durable learnings.
Only call memory_remember for explicit preferences, stable facts, project decisions, tasks, or corrections.
Ask before storing sensitive, private, or uncertain information.
Prefer compact memories over raw document chunks.
```
