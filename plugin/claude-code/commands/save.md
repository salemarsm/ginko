---
description: Save a structured memory about the current work via the ginko MCP server.
---

Save a durable memory using the ginko MCP server (`memory_remember` tool).

If the user did not specify the details, infer them from the recent conversation:

- **type**: one of `preference`, `fact`, `decision`, `task`, `note`. Pick the type that best matches the content.
- **subject**: the project (auto-detected — leave blank unless the user specifies).
- **content**: a self-contained sentence or two that a future session could act on without re-investigating. Lead with file path or scope when applicable.
- **confidence**: 0.9 default. Lower (0.6–0.8) if you are uncertain.

Then call `memory_remember` with those values and report the new memory ID.

If the work just completed is trivial (routine task, pure refactor, content already in CLAUDE.md), tell the user it is not worth saving and do not call the tool. Do not save by reflex.

Apply the rules from the Ginko Memory Protocol skill if available.
