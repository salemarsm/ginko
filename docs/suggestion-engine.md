# Suggestion engine

Status: **implemented as conservative heuristics, not an LLM judge**.

The suggestion engine is deliberately small in v0.x. It proposes memory candidates from:

- user prompts
- assistant responses
- concise LLM-provided inference text

It does **not** automatically persist suggestions. It returns candidates with reasons and `requires_confirmation=true`.

## Flow

```txt
User prompt / assistant response / LLM inference
  ↓
/api/suggest or MCP memory_suggest
  ↓
heuristic candidate extraction
  ↓
candidate list with type/confidence/reason
  ↓
human or policy approval
  ↓
/api/memories or MCP memory_remember
```

## Current heuristics

The current implementation scans sentence-like fragments for durable signals.

Examples:

| Signal | Candidate type | Example |
|---|---|---|
| `I prefer`, `prefiro`, `me chame`, `quero que` | `preference` | `I prefer Go examples.` |
| `we decided`, `decidimos`, `vamos usar` | `decision` | `We decided to use SQLite.` |
| `need to`, `preciso`, `todo`, `next step` | `task` | `Need to add API tokens.` |
| project/system/user factual statements | `fact` | `The project uses SQLite.` |
| LLM inference says `should remember` / `aprendizado` | `note` | `Remember that the user values concise answers.` |

## Example request

```bash
curl -X POST http://127.0.0.1:8787/api/suggest \
  -H 'content-type: application/json' \
  -d '{
    "subject": "botmaster",
    "scope": "global",
    "user_prompt": "I prefer direct technical answers. We decided to use SQLite as canonical storage.",
    "max_candidates": 5
  }'
```

## Example response

```json
{
  "candidates": [
    {
      "memory": {
        "type": "preference",
        "subject": "botmaster",
        "content": "I prefer direct technical answers",
        "source": { "kind": "suggestion", "ref": "prompt" },
        "scope": "global",
        "confidence": 0.9,
        "tags": ["suggested"],
        "embedding_refs": {}
      },
      "reason": "explicit preference/style instruction",
      "requires_confirmation": true
    },
    {
      "memory": {
        "type": "decision",
        "subject": "botmaster",
        "content": "We decided to use SQLite as canonical storage",
        "source": { "kind": "suggestion", "ref": "prompt" },
        "scope": "global",
        "confidence": 0.82,
        "tags": ["suggested"],
        "embedding_refs": {}
      },
      "reason": "project decision or implementation choice",
      "requires_confirmation": true
    }
  ]
}
```

## Limitations

- It is heuristic and language-pattern based.
- It does not resolve contradictions yet.
- It does not classify sensitive data yet.
- It does not perform semantic dedupe beyond simple exact-ish candidate keys.
- It should be treated as a candidate generator, not an authority.

## Planned improvements

- sensitive-data detection
- contradiction checks against existing memories
- supersession recommendations
- optional LLM judge mode
- per-subject write policy
- explanation traces for rejected candidates
