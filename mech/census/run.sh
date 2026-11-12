#!/usr/bin/env bash
# Runs the external-API census over the pinned Closure production sources.
# Needs a JDK 21+ and `make oracle-tools` (for .cache/closure-src). Resolves
# types against the npm build of the same release, whose jar bundles Guava,
# Gson, protobuf and args4j unshaded (the Maven jar relocates them).
set -euo pipefail
cd "$(dirname "$0")/../.."
source tools/oracle/pin.env
ver=${CLOSURE_VERSION#v}.0.0
jar=.cache/closure-compiler-$CLOSURE_VERSION-npm.jar
if [[ ! -s $jar ]]; then
  curl -fsS "https://registry.npmjs.org/google-closure-compiler-java/-/google-closure-compiler-java-$ver.tgz" |
    tar -xzO package/compiler.jar >"$jar.tmp" && mv "$jar.tmp" "$jar"
fi
find .cache/closure-src -name '*.java' -not -path '*/ant/*' | LC_ALL=C sort >.cache/census-files.txt
java --add-modules jdk.compiler -Xmx6g mech/census/Census.java "$jar" mech/census/census.tsv @.cache/census-files.txt
