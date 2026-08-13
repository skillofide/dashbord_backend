#!/bin/sh
set -e
# printf, not echo: sh's echo interprets backslash escapes, so `\n` or a
# regex like `\d` in submitted code would be corrupted before it runs.
printf '%s\n' "$USER_CODE" > /tmp/solution.js
printf '%s\n' "$USER_INPUT" | node /tmp/solution.js
