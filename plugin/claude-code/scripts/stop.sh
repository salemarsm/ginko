#!/usr/bin/env bash
# ginko Stop hook
#
# Fires at the end of each turn. Currently a no-op placeholder.
# Reserved for future "refresh memory usage stats" or end-of-turn
# bookkeeping. Non-blocking.

set -uo pipefail
INPUT=$(cat || true)
exit 0
