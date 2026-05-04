---
description: Recall ginko memories relevant to the current task or to a specific query the user provides.
---

Use the ginko MCP server to recall relevant memories.

If the user provided a query as argument:

1. Call `memory_search` with that query and a limit of 10.
2. Summarize the top results, grouped by type (`decision`, `fact`, `preference`, `task`, `note`).
3. Highlight any memory with `confidence < 0.7` or that appears stale (older than 90 days) so the user can decide if it is still valid.

If the user did not provide a query:

1. Call `memory_context` with the current project subject and a token budget of about 800.
2. Print the summary in a compact form.
3. Ask the user what they want to do with that context (continue prior work, supersede, ignore).

Do not list memory IDs unless the user asks. Reference memories by short summary.
