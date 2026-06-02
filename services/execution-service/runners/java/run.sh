#!/bin/sh
set -e
echo "$USER_CODE" > /tmp/Solution.java
cd /tmp
javac -J-XX:TieredStopAtLevel=1 -J-XX:+UseSerialGC -J-Xms8m -J-Xmx128m Solution.java
echo "$USER_INPUT" | java -XX:TieredStopAtLevel=1 -XX:+UseSerialGC -Xms8m -Xmx128m Solution

