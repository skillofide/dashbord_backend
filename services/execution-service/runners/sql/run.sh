#!/bin/sh
# Boots the pre-initialised PostgreSQL cluster, then hands over to the harness,
# which builds the fixture tables, runs USER_CODE, and prints the result rows
# as JSON.
#
# The cluster listens only on a unix socket in /tmp: the sandbox gives the
# container no network, and nothing outside it should be able to reach this.
set -e

pg_ctl -D "$PGDATA" -w -t 30 start >/dev/null 2>&1 || {
  echo "Failed to start the database sandbox." >&2
  exit 1
}

exec python3 /harness.py
