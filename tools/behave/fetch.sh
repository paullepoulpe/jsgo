#!/usr/bin/env bash
# Fetches the behavioural oracle's suites at the revisions pinned in
# tools/behave/pins.json and verifies their digests (behave.py verify).
# Idempotent. Nothing fetched here is committed.
#
#   .cache/behave/test262          test262 (BSD), shallow, at the pinned commit
#   .cache/behave/closure-compiler Closure Compiler (Apache 2.0) at the pinned tag,
#                                  sparse: only runtime_tests/
#   .cache/corpus/google-closure-library@<version>
#                                  Closure Library (Apache 2.0) from the corpus
#                                  (tools/oracle/corpus.sh), for goog.testing
#
# The cache is the main checkout's .cache/ (shared by all worktrees), or
# $JSGO_CACHE.
set -euo pipefail
cd "$(dirname "$0")/../.."
CACHE=${JSGO_CACHE:-$(dirname "$(git rev-parse --path-format=absolute --git-common-dir)")/.cache}
pin() { python3 -c "import json,sys; d=json.load(open('tools/behave/pins.json')); print(d[sys.argv[1]][sys.argv[2]])" "$1" "$2"; }

# fetch_commit <dir> <repo> <commit> [sparse path]: a shallow checkout of one commit.
fetch_commit() {
  local dir="$1" repo="$2" commit="$3" sparse="${4:-}"
  if [[ "$(git -C "$dir" rev-parse HEAD 2>/dev/null)" == "$commit" ]]; then
    return 0
  fi
  rm -rf "$dir"
  mkdir -p "$dir"
  git -C "$dir" init -q
  git -C "$dir" remote add origin "$repo"
  if [[ -n "$sparse" ]]; then
    git -C "$dir" config core.sparseCheckout true
    git -C "$dir" sparse-checkout set "$sparse"
    GIT_LFS_SKIP_SMUDGE=1 git -C "$dir" fetch -q --depth 1 --filter=blob:none origin "$commit"
  else
    GIT_LFS_SKIP_SMUDGE=1 git -C "$dir" fetch -q --depth 1 origin "$commit"
  fi
  git -C "$dir" -c advice.detachedHead=false checkout -q FETCH_HEAD
}

fetch_commit "$CACHE/$(pin test262 dir)" "$(pin test262 repo)" "$(pin test262 commit)"
fetch_commit "$CACHE/$(pin closure dir)" "$(pin closure repo)" "$(pin closure commit)" "$(pin closure path)"
if [[ ! -d "$CACHE/$(pin closure_library dir)" ]]; then
  bash tools/oracle/corpus.sh tools/oracle/corpus.txt "$CACHE/corpus"
fi
python3 tools/behave/behave.py verify --cache "$CACHE"
