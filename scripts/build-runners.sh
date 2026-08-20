#!/usr/bin/env bash
# Builds the sandboxed language runner images the judge executes code inside.
#
# Without these, execution-service answers every run with
#   create container: No such image: skillofide/runner-<lang>:latest
# which surfaces to a candidate as a RuntimeError on correct code — and to a
# scholarship as a coding section everybody scores zero on. The images are built
# by CI for deploys and by start-dev.ps1 on Windows; this is the same list for
# anyone working on macOS or Linux.
#
#   ./scripts/build-runners.sh                 # all of them
#   ./scripts/build-runners.sh python javascript
set -euo pipefail

cd "$(dirname "$0")/.."

declare -a LANGS=(python javascript java cpp go sql)
if [ $# -gt 0 ]; then LANGS=("$@"); fi

for lang in "${LANGS[@]}"; do
  dir="services/execution-service/runners/$lang"
  [ -d "$dir" ] || { printf '  \033[31m✗\033[0m no runner directory for %s\n' "$lang"; exit 1; }
  printf '\n\033[1mbuilding skillofide/runner-%s:latest\033[0m\n' "$lang"
  docker build -q -t "skillofide/runner-$lang:latest" "$dir" \
    && printf '  \033[32m✓\033[0m %s\n' "$lang"
done

printf '\n\033[1mpresent now:\033[0m\n'
docker images --format '  {{.Repository}}:{{.Tag}}  {{.Size}}' | grep 'skillofide/runner-' | sort
