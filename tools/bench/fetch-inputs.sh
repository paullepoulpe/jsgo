#!/usr/bin/env bash
# Fetches the benchmark inputs pinned in tools/bench/inputs.txt into
# .cache/bench-inputs/ and verifies their checksums. Idempotent: an input whose
# directory exists is not fetched again (delete it to re-fetch). Also makes
# sure the oracle corpus (.cache/corpus, the source of the single-file task) is
# present.
#
#   tools/bench/fetch-inputs.sh [manifest] [outdir]
#   default: tools/bench/inputs.txt -> $CACHE/bench-inputs
#
# $CACHE is $JSGO_CACHE if set, else the .cache/ of the main checkout (shared
# by all worktrees; it also holds the pinned jar and the corpus).
#
# The inputs are third-party code: they are data for the compilers and are
# never executed, and never committed.
set -euo pipefail
cd "$(dirname "$0")/../.."
cache=${JSGO_CACHE:-$(cd "$(git rev-parse --git-common-dir)/.." && pwd)/.cache}
manifest=${1:-tools/bench/inputs.txt}
out=${2:-$cache/bench-inputs}
mkdir -p "$out" "$cache/bench-tgz"

listing_sha256() { # <dir>: sha256 of the sorted per-file sha256 listing
  (cd "$1" && find . -type f -not -path './.git/*' -print0 | LC_ALL=C sort -z |
    xargs -0 sha256sum | sha256sum | cut -d' ' -f1)
}

grep -vE '^\s*(#|$)' "$manifest" | while read -r kind name version rest; do
  dest="$out/$name@$version"
  [[ -d "$dest" ]] && continue
  rm -rf "$dest.tmp"
  case "$kind" in
  npm)
    integrity=$rest
    tgz="$cache/bench-tgz/$(echo "$name@$version" | tr '/' '_').tgz"
    base=${name##*/}
    if [[ ! -s "$tgz" ]]; then
      curl -fsS -o "$tgz.tmp" "https://registry.npmjs.org/$name/-/$base-$version.tgz"
      mv "$tgz.tmp" "$tgz"
    fi
    got="sha512-$(openssl dgst -sha512 -binary "$tgz" | base64 -w0)"
    if [[ "$got" != "$integrity" ]]; then
      echo "integrity mismatch for $name@$version: got $got" >&2
      exit 1
    fi
    mkdir -p "$dest.tmp"
    tar -xzf "$tgz" -C "$dest.tmp" --strip-components=1 --no-same-owner --no-same-permissions
    ;;
  git)
    read -r url commit tree listing <<<"$rest"
    GIT_LFS_SKIP_SMUDGE=1 git -c advice.detachedHead=false clone -q --depth 1 --branch "$version" "$url" "$dest.tmp"
    got_commit=$(git -C "$dest.tmp" rev-parse HEAD)
    got_tree=$(git -C "$dest.tmp" rev-parse 'HEAD^{tree}')
    if [[ "$got_commit" != "$commit" || "$got_tree" != "$tree" ]]; then
      echo "$name $version: got commit $got_commit tree $got_tree, want $commit $tree" >&2
      exit 1
    fi
    rm -rf "$dest.tmp/.git"
    got_listing=$(listing_sha256 "$dest.tmp")
    if [[ "$got_listing" != "$listing" ]]; then
      echo "$name $version: file listing sha256 $got_listing, want $listing" >&2
      exit 1
    fi
    ;;
  *)
    echo "unknown input kind $kind" >&2
    exit 1
    ;;
  esac
  mkdir -p "$(dirname "$dest")"
  mv "$dest.tmp" "$dest"
  echo "$name@$version: $(find "$dest" -type f | wc -l) files"
done

# The single-file task draws from the oracle corpus.
if [[ ! -d "$cache/corpus" ]]; then
  bash tools/oracle/corpus.sh tools/oracle/corpus.txt "$cache/corpus" >/dev/null
fi
echo "bench inputs ready in $out/"
