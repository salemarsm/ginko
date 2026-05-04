#!/usr/bin/env sh
set -eu
project="$(${CLAUDE_PLUGIN_ROOT:-$(dirname "$0")/..}/scripts/detect-project.sh | tail -n 1)"
if command -v memctl >/dev/null 2>&1; then
  memctl -subject "$project" -scope project -max-tokens 600 context "Claude Code session start: recall relevant project decisions, facts, preferences, and recent constraints." 2>/dev/null || true
fi
