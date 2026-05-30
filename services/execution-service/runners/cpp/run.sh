#!/bin/sh
set -e
echo "$USER_CODE" > /tmp/solution.cpp
g++ -O2 -o /tmp/solution /tmp/solution.cpp
echo "$USER_INPUT" | /tmp/solution
