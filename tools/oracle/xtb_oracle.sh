#!/usr/bin/env bash
# Regenerates internal/jscomp/testdata/xtb/expected.txt (MessagesTypesMiscDump xtb mode).
#
# The pinned (Maven, shaded) closure-compiler.jar cannot run XtbMessageBundle: its jarjar rules
# relocate javax.xml.parsers to com.google.javascript.jscomp.jarjar.javax.xml.parsers without
# bundling those classes, so `new XtbMessageBundle(...)` throws NoClassDefFoundError. This script
# compiles XtbMessageBundle.java from the pinned sources (only its Guava imports rewritten to the
# jar's shaded package) and puts it ahead of the jar on the classpath.
# The jar is .cache/$CLOSURE_JAR (tools/oracle/pin.env).
set -euo pipefail
cd "$(dirname "$0")/../.."
source tools/oracle/pin.env
CACHE=${CACHE:-.cache}
[ -d "$CACHE/closure-src" ] || CACHE=/home/user/jsgo/.cache
JAR=$CACHE/$CLOSURE_JAR
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
mkdir -p "$work/src" "$work/xtb" "$work/build"
sed 's/import com\.google\.common\./import com.google.javascript.jscomp.jarjar.com.google.common./;
     s/import static com\.google\.common\./import static com.google.javascript.jscomp.jarjar.com.google.common./' \
  "$CACHE/closure-src/com/google/javascript/jscomp/XtbMessageBundle.java" > "$work/src/XtbMessageBundle.java"
javac -nowarn -cp "$JAR" -d "$work/xtb" "$work/src/XtbMessageBundle.java"
javac -nowarn -cp "$JAR" -d "$work/build" tools/oracle/java/MessagesTypesMiscDump.java \
  tools/oracle/java/CorpusExclusions.java
java -cp "$work/xtb:$JAR:$work/build" \
  com.google.javascript.jscomp.MessagesTypesMiscDump xtb - internal/jscomp/testdata/xtb \
  2>/dev/null > internal/jscomp/testdata/xtb/expected.txt
