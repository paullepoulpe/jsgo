#!/usr/bin/env bash
# Compiles the Java oracle tools against the pinned jar into .cache/oracle-build.
# Run them with:  java -cp .cache/closure-compiler.jar:.cache/oracle-build <Tool> ...
set -euo pipefail
cd "$(dirname "$0")/../.."
bash tools/oracle/fetch.sh >/dev/null
mkdir -p .cache/oracle-build
javac -nowarn -cp .cache/closure-compiler.jar -d .cache/oracle-build tools/oracle/java/*.java
