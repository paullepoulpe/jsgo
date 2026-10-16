#!/usr/bin/env bash
# CI guards for jsgo. Each guard corresponds to a numbered entry in
# SPEC.md § "Hard Guardrails". Exits non-zero on first violation.

set -euo pipefail

repo_root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

fail=0
note() { echo "GUARD: $*"; }
violate() { echo "VIOLATION: $*" >&2; fail=1; }

# -- Guard #1: pure Go at runtime ---------------------------------------
# Forbidden packages under cmd/ and internal/. tools/ and harness/ (which
# only runs under `go test`) are exempt by design.
note "#1 pure-Go runtime (no cgo / os/exec / net)"
banned_imports=(
    '"os/exec"'
    '"net"'
    '"net/http"'
    '"net/rpc"'
)
for pat in "${banned_imports[@]}"; do
    if grep -R --include='*.go' -nF "$pat" cmd internal 2>/dev/null; then
        violate "banned import $pat found in cmd/ or internal/"
    fi
done
if grep -R --include='*.go' -nE '^import "C"|"C"[[:space:]]*$' cmd internal 2>/dev/null; then
    violate "cgo import found in cmd/ or internal/"
fi

# -- Guard #2: no hardcoded oracle data ---------------------------------
# For every fixture, take the filename stem and grep for it outside
# testdata/. Any match implies the fixture is being keyed on by source.
note "#2 no fixture names/paths leaked into cmd/ or internal/"
while IFS= read -r fix; do
    stem="$(basename "$fix" .js)"
    if grep -R --include='*.go' -nF "$stem" cmd internal 2>/dev/null; then
        violate "fixture stem '$stem' referenced from runtime source"
    fi
done < <(find testdata/fixtures -name '*.js' 2>/dev/null || true)

# Also scan for the oracle-file exit marker — it should never appear
# in runtime code (harness is a test support file and may parse it).
if grep -R --include='*.go' -nF -- '---EXIT---' cmd internal >/dev/null 2>&1; then
    violate "'---EXIT---' oracle marker leaked into runtime source"
fi

# -- Guard #3: traceability header --------------------------------------
# Every .go file under internal/ must open with "// Port of " in the
# first non-empty line. (cmd/ is exempt since it's not a port.)
# A leading `//go:build` constraint is allowed to precede the header
# so a file can be gated behind a build tag.
note "#3 traceability header on every file under internal/"
while IFS= read -r f; do
    first="$(awk 'NF && $0 !~ /^\/\/go:build/ { print; exit }' "$f")"
    case "$first" in
        "// Port of "*) ;;
        *) violate "missing '// Port of …' header: $f" ;;
    esac
done < <(find internal -name '*.go' -not -path 'internal/jutil/*' 2>/dev/null)  # jutil: Java-semantics runtime, not a port

# -- Report -------------------------------------------------------------
if [ "$fail" -ne 0 ]; then
    echo "ci-guards FAILED" >&2
    exit 1
fi
echo "ci-guards OK"
