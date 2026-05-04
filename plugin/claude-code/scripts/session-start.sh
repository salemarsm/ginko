#!/usr/bin/env bash
# ginko SessionStart hook
#
# Fires when a Claude Code session starts (or resumes). Stdout from this
# script is injected as context that Claude can see in the session
# transcript. Stderr goes to the debug log.
#
# We call `ginko context` to fetch the most relevant memories for the
# current project and print them. The agent reads this and starts the
# session already aware of prior decisions, preferences, and facts.
#
# This script is non-blocking: any failure exits 0 with a stderr note,
# so the user's session is never blocked by ginko being unavailable.

set -uo pipefail

# Read hook input JSON from stdin (we don't use it currently, but it must
# be drained or the parent may block).
INPUT=$(cat || true)

# If ginko is not on PATH, exit silently. The user installed Claude Code
# but not ginko — that's fine, just no memory context.
if ! command -v ginko >/dev/null 2>&1; then
  echo "[ginko] binary not found on PATH — skipping context injection." >&2
  exit 0
fi

# Detect project name from cwd (basename) as a fallback subject.
# When ginko gains git-aware project detection (roadmap v0.3), this script
# can be simplified to just `ginko context`.
PROJECT="$(basename "$(pwd)" 2>/dev/null || echo "")"

# Fetch context. Bound the budget so we don't dump megabytes into the
# agent's context window. Cap to ~600 tokens of memory.
if ! OUTPUT=$(ginko context "$PROJECT" --limit-tokens 600 2>/dev/null); then
  # Older ginko versions may not support the flags. Fall back to default.
  OUTPUT=$(ginko context 2>/dev/null || true)
fi

# Empty output is OK — first session, no memories yet. Print nothing.
if [ -z "${OUTPUT:-}" ]; then
  exit 0
fi

# Print as a clearly delimited context block. The agent treats stdout
# from SessionStart as additional context.
cat <<EOF
=== ginko: prior context for "$PROJECT" ===
$OUTPUT
=== end ginko context ===
EOF

exit 0
