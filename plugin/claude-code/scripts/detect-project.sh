#!/usr/bin/env sh
set -eu
if [ -f .llm-memory/config.json ]; then
  python3 - <<'PY' 2>/dev/null || true
import json
with open('.llm-memory/config.json') as f:
    data=json.load(f)
print(data.get('project') or data.get('subject') or '')
PY
fi
remote="$(git config --get remote.origin.url 2>/dev/null || true)"
if [ -n "$remote" ]; then
  printf '%s\n' "$remote" | sed -E 's#^git@[^:]+:##; s#^https?://[^/]+/##; s#\.git$##; s#[^A-Za-z0-9._/-]+#-#g; s#/#-#g' | tr '[:upper:]' '[:lower:]'
  exit 0
fi
basename "$(pwd)" | tr '[:upper:]' '[:lower:]' | sed -E 's#[^a-z0-9._-]+#-#g'
