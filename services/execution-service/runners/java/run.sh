#!/bin/sh
set -e
# printf, not echo: sh's echo interprets backslash escapes, so `\n` or a
# regex like `\d` in submitted code would be corrupted before it runs.
printf '%s\n' "$USER_CODE" > /tmp/Solution.java
cd /tmp
javac -J-XX:TieredStopAtLevel=1 -J-XX:+UseSerialGC -J-Xms8m -J-Xmx128m Solution.java
printf '%s\n' "$USER_INPUT" | java -XX:TieredStopAtLevel=1 -XX:+UseSerialGC -Xms8m -Xmx128m Solution

