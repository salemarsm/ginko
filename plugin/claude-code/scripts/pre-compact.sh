#!/usr/bin/env sh
set -eu
project="$(${CLAUDE_PLUGIN_ROOT:-$(dirname "$0")/..}/scripts/detect-project.sh | tail -n 1)"
if command -v memctl >/dev/null 2>&1; then
  printf 'Claude Code is compacting this session. Save durable decisions, bugfixes, preferences, and constraints explicitly before compaction if they are not already stored.\n'
  memctl -subject "$project" -scope project -max-tokens 400 context "Pre-compact checkpoint: prior durable project context." 2>/dev/null || true
fi
