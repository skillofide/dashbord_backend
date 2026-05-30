#!/usr/bin/env bash
# scripts/migrate.sh
# Runs database migrations for all services using golang-migrate.
#
# Prerequisites: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
#
# Usage:
#   ./scripts/migrate.sh up      # apply all pending migrations
#   ./scripts/migrate.sh down    # rollback last migration
#   ./scripts/migrate.sh status  # show current migration status

set -euo pipefail

DIRECTION="${1:-up}"
DSN="${POSTGRES_DSN:-postgres://skillofide:password@localhost:5432/skillofide?sslmode=disable}"
ROOT="$(dirname "$0")/.."

SERVICES=(
  "problem-service"
  "submission-service"
  "progress-service"
)

for svc in "${SERVICES[@]}"; do
  MIGRATION_DIR="$ROOT/services/$svc/migrations"
  if [ ! -d "$MIGRATION_DIR" ]; then
    echo "  [skip] No migrations dir for $svc"
    continue
  fi

  echo "  → Running $DIRECTION migrations for $svc..."
  migrate \
    -path "$MIGRATION_DIR" \
    -database "$DSN" \
    "$DIRECTION"
done

echo "Migrations complete."
