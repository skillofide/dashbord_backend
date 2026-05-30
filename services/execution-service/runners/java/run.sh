#!/bin/sh
set -e
echo "$USER_CODE" > /tmp/Solution.java
cd /tmp
javac Solution.java
echo "$USER_INPUT" | java Solution
