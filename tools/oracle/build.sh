#!/usr/bin/env bash
# Compiles the Java oracle tools against the pinned jar into .cache/oracle-build.
# Run them with:  java -cp .cache/$CLOSURE_JAR:.cache/oracle-build <Tool> ...
# where CLOSURE_JAR comes from tools/oracle/pin.env (source it).
set -euo pipefail
cd "$(dirname "$0")/../.."
source tools/oracle/pin.env
bash tools/oracle/fetch.sh >/dev/null
mkdir -p .cache/oracle-build
javac -nowarn -cp ".cache/$CLOSURE_JAR" -d .cache/oracle-build tools/oracle/java/*.java
