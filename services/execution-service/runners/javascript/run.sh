#!/bin/sh
set -e
echo "$USER_CODE" > /tmp/solution.js
echo "$USER_INPUT" | node /tmp/solution.js
