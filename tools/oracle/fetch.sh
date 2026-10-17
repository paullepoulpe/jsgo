#!/usr/bin/env bash
# Fetches the pinned Closure Compiler release (binary jar + Java sources) into
# .cache/ and verifies the jar checksum. Idempotent.
#
#   .cache/closure-compiler.jar   the oracle
#   .cache/closure-src/           Java sources for porting (com/google/...)
set -euo pipefail
cd "$(dirname "$0")/../.."

VERSION=v20260415
JAR_SHA256=20313094c6e9cac8c685931e0b3fd67f82e1534c30448cdff404ee350460e391
PATH_PREFIX=com/google/javascript/closure-compiler/$VERSION/closure-compiler-$VERSION
# Maven Central, with Google's mirror as fallback (Central rate-limits).
MIRRORS=(https://repo1.maven.org/maven2 https://maven-central.storage-download.googleapis.com/maven2)

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

fetch .jar .cache/closure-compiler.jar
echo "$JAR_SHA256  .cache/closure-compiler.jar" | sha256sum -c - >/dev/null

if [[ ! -d .cache/closure-src/com ]]; then
  fetch -sources.jar .cache/closure-src.jar
  rm -rf .cache/closure-src && mkdir -p .cache/closure-src
  (cd .cache/closure-src && unzip -qo ../closure-src.jar)
fi
echo "closure-compiler $VERSION ready in .cache/"
