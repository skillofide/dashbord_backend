#!/bin/sh
set -e
# printf, not echo: sh's echo interprets backslash escapes, so `\n` or a
# regex like `\d` in submitted code would be corrupted before it runs.
printf '%s\n' "$USER_CODE" > /tmp/solution.go
cd /tmp
# Build first so compile errors are reported without the program running.
go build -o /tmp/solution /tmp/solution.go
printf '%s\n' "$USER_INPUT" | /tmp/solution
