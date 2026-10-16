#!/usr/bin/env bash
# Generate oracle files for fixtures in parallel.
#
# Produces two oracle trees from every fixture:
#   testdata/oracle/<rel>.txt         — default `--jscomp_error=checkTypes`
#   testdata/oracle-strict/<rel>.txt  — +`--jscomp_warning=strictCheckTypes`
#
# Mirrors how closure-compiler's TypeCheckTest.java covers both the
# default `TypeCheck` diagnostic surface and the stricter
# `TypeValidator` operand checks gated behind STRICT_CHECK_TYPES.

set -u
repo_root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

JAR="${JAR:-.cache/closure-compiler.jar}"
WORKERS="${WORKERS:-8}"

one() {
    local fix="$1"
    local rel="${fix#testdata/fixtures/}"
    local out="testdata/oracle/${rel%.js}.txt"
    local outStrict="testdata/oracle-strict/${rel%.js}.txt"
    # Skip fixtures whose oracles are already fully populated so
    # gen-oracles.sh stays idempotent — re-runs only regenerate what's
    # missing, matching how Java's test infra caches checker output.
    # Set FORCE=1 to force regeneration.
    if [[ "${FORCE:-0}" == "0" && -s "$out" && -s "$outStrict" ]]; then
        return
    fi
    mkdir -p "$(dirname "$out")" "$(dirname "$outStrict")"
    local tmp
    tmp="$(mktemp)"
    # Default oracle.
    if [[ "${FORCE:-0}" == "1" || ! -s "$out" ]]; then
        java -jar "$JAR" \
            --checks_only \
            --jscomp_error=checkTypes \
            --warning_level=VERBOSE \
            --language_in=ECMASCRIPT_NEXT \
            --js "$fix" > "$tmp" 2>&1
        local ec=$?
        { cat "$tmp"; echo "---EXIT---"; echo "$ec"; } > "$out"
    fi
    # Strict oracle.
    if [[ "${FORCE:-0}" == "1" || ! -s "$outStrict" ]]; then
        java -jar "$JAR" \
            --checks_only \
            --jscomp_error=checkTypes \
            --jscomp_warning=strictCheckTypes \
            --warning_level=VERBOSE \
            --language_in=ECMASCRIPT_NEXT \
            --js "$fix" > "$tmp" 2>&1
        local ec=$?
        { cat "$tmp"; echo "---EXIT---"; echo "$ec"; } > "$outStrict"
    fi
    rm -f "$tmp"
}
export -f one
export JAR

find testdata/fixtures -name '*.js' | sort | \
    xargs -I{} -P "$WORKERS" bash -c 'one "$@"' _ {}

echo "done"
