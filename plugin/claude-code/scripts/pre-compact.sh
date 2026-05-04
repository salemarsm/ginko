#!/usr/bin/env bash
# ginko PreCompact hook
#
# Fires when Claude Code is about to compact (summarize) a long
# conversation. This is the danger zone: if the agent didn't save
# memories during the session, that knowledge is about to be lost.
#
# We auto-checkpoint by saving a short note-type memory marking that
# a compaction occurred. The actual summary should ideally be written
# by the agent itself via `memory_remember` before this fires; this
# script is the safety net.
#
# Non-blocking on failure.

set -uo pipefail

INPUT=$(cat || true)

if ! command -v ginko >/dev/null 2>&1; then
  exit 0
fi

PROJECT="$(basename "$(pwd)" 2>/dev/null || echo "")"
TIMESTAMP="$(date -u +%Y-%m-%dT%H:%M:%SZ)"

# Save a checkpoint note. If the agent already saved a richer summary,
# this is just a low-confidence breadcrumb and won't hurt retrieval.
ginko save \
  --type note \
  --subject "$PROJECT" \
  --confidence 0.5 \
  --tag ginko-checkpoint \
  "Auto-checkpoint at $TIMESTAMP: session compaction occurred. Earlier turns were summarized; check prior memories for canonical decisions." \
  >/dev/null 2>&1 || true

exit 0
