#!/usr/bin/env bash
# Runs the problem-content gate: every function-mode problem's reference
# solution must be accepted, and its generated starter must not be.
#
# This is the check that would have caught the whole original audit — the
# unparseable multi-line inputs, the JSON-quoted string returns, the Python
# drivers calling .trim(), and the 361 starters that shipped the answer.
#
# Requires the runner images (./scripts/build-runners.sh) and a reachable
# Postgres. Intended for CI on every change to codegen, the judge, or seed data.
set -euo pipefail
cd "$(dirname "$0")/.."

LANGS="${LANGS:-javascript,python,java,cpp,go}"
exec go run ./services/execution-service/cmd/validate-problems -langs "$LANGS" "$@"
