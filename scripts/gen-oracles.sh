#!/usr/bin/env bash
# Generate oracle files for fixtures in parallel.

set -u
repo_root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

JAR="${JAR:-.cache/closure-compiler.jar}"
WORKERS="${WORKERS:-8}"

one() {
    local fix="$1"
    local rel="${fix#testdata/fixtures/}"
    local out="testdata/oracle/${rel%.js}.txt"
    mkdir -p "$(dirname "$out")"
    local tmp
    tmp="$(mktemp)"
    java -jar "$JAR" \
        --checks_only \
        --jscomp_error=checkTypes \
        --warning_level=VERBOSE \
        --language_in=ECMASCRIPT_NEXT \
        --js "$fix" > "$tmp" 2>&1
    local ec=$?
    { cat "$tmp"; echo "---EXIT---"; echo "$ec"; } > "$out"
    rm -f "$tmp"
}
export -f one
export JAR

find testdata/fixtures -name '*.js' | sort | \
    xargs -I{} -P "$WORKERS" bash -c 'one "$@"' _ {}

echo "done"
