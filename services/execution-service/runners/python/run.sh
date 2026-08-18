#!/bin/sh
# Writes USER_CODE to /tmp/solution.py and runs it with USER_INPUT on stdin.
set -e
# printf, not echo: sh's echo interprets backslash escapes, so `\n` or a
# regex like `\d` in submitted code would be corrupted before it runs.
printf '%s\n' "$USER_CODE" > /tmp/solution.py
printf '%s\n' "$USER_INPUT" | python3 /tmp/solution.py
