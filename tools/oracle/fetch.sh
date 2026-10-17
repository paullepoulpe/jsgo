#!/usr/bin/env bash
# Fetches the pinned Closure Compiler release (binary jar + Java sources) into
# .cache/ and verifies them. Idempotent. The release and the jar's name and
# checksums are in tools/oracle/pin.env.
#
#   .cache/$CLOSURE_JAR   the oracle jar (Maven Central's closure-compiler-$CLOSURE_VERSION.jar)
#   .cache/closure-src/   Java sources for porting (com/google/...)
set -euo pipefail
cd "$(dirname "$0")/../.."
source tools/oracle/pin.env

PATH_PREFIX=com/google/javascript/closure-compiler/$CLOSURE_VERSION/closure-compiler-$CLOSURE_VERSION
# Maven Central, with Google's mirror as fallback (Central rate-limits).
MIRRORS=(https://repo1.maven.org/maven2 https://maven-central.storage-download.googleapis.com/maven2)
JAR=.cache/$CLOSURE_JAR

mkdir -p .cache

fetch() { # <suffix> <dest>
  local suffix="$1" dest="$2"
  [[ -s "$dest" ]] && return 0
  for m in "${MIRRORS[@]}"; do
    for attempt in 1 2 3; do
      if curl -fsSL -o "$dest.tmp" "$m/$PATH_PREFIX$suffix"; then
        mv "$dest.tmp" "$dest"
        return 0
      fi
      sleep $((attempt * 2))
    done
  done
  echo "failed to fetch $suffix" >&2
  return 1
}

fetch .jar "$JAR"
sum=$(sha256sum "$JAR" | cut -d' ' -f1)
if [[ " $CLOSURE_JAR_SHA256 " != *" $sum "* ]]; then
  echo "$JAR (sha256 $sum) is not closure-compiler $CLOSURE_VERSION (tools/oracle/pin.env);" \
    "delete it to fetch the pinned jar" >&2
  exit 1
fi

# The sources jar holds the generated CommandLineRunnerVersion.java at its
# root, which names the release.
if [[ -d .cache/closure-src/com ]] &&
  ! grep -q "\"$CLOSURE_VERSION\"" .cache/closure-src/CommandLineRunnerVersion.java 2>/dev/null; then
  echo ".cache/closure-src holds another release than $CLOSURE_VERSION (tools/oracle/pin.env);" \
    "delete it and .cache/closure-src.jar to fetch the pinned sources" >&2
  exit 1
fi
if [[ ! -d .cache/closure-src/com ]]; then
  rm -f .cache/closure-src.jar
  fetch -sources.jar .cache/closure-src.jar
  rm -rf .cache/closure-src && mkdir -p .cache/closure-src
  (cd .cache/closure-src && unzip -qo ../closure-src.jar)
fi
echo "closure-compiler $CLOSURE_VERSION ready in .cache/"
