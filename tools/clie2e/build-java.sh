#!/usr/bin/env bash
# Compiles the Java side of tools/clie2e into DIR (default .cache/clie2e-build):
# CliE2E and CorpusExclusions, plus XtbMessageBundle recompiled from the pinned
# sources. The pinned closure-compiler.jar cannot load its own XtbMessageBundle
# (its jarjar rules relocate javax.xml.parsers without bundling those classes,
# so --translations_file throws NoClassDefFoundError; see
# tools/oracle/xtb_oracle.sh). The projects suite puts DIR ahead of the jar on
# the classpath, so the recompiled class, identical to the jar's but for the
# Guava package of its imports, replaces it.
#
#   tools/clie2e/build-java.sh [DIR]
#   go run ./tools/clie2e -suite projects -javabuild DIR ...
set -euo pipefail
cd "$(dirname "$0")/../.."
CACHE=${CACHE:-.cache}
[ -d "$CACHE/closure-src" ] || CACHE=/home/user/jsgo/.cache
out=${1:-$CACHE/clie2e-build}
mkdir -p "$out"
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
javac -nowarn -cp "$CACHE/closure-compiler.jar" -d "$out" \
  tools/oracle/java/CliE2E.java tools/oracle/java/CorpusExclusions.java
sed 's/import com\.google\.common\./import com.google.javascript.jscomp.jarjar.com.google.common./;
     s/import static com\.google\.common\./import static com.google.javascript.jscomp.jarjar.com.google.common./' \
  "$CACHE/closure-src/com/google/javascript/jscomp/XtbMessageBundle.java" > "$work/XtbMessageBundle.java"
javac -nowarn -cp "$CACHE/closure-compiler.jar" -d "$out" "$work/XtbMessageBundle.java"
