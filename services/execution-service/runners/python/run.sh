#!/bin/sh
# Writes USER_CODE to /tmp/solution.py and runs it with USER_INPUT on stdin.
set -e
echo "$USER_CODE" > /tmp/solution.py
echo "$USER_INPUT" | python3 /tmp/solution.py
