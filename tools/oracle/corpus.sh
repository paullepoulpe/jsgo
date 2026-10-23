#!/usr/bin/env bash
# Fetches a pinned real-world JS corpus into a cache directory. Idempotent;
# verifies npm sha512 integrity.
#
#   tools/oracle/corpus.sh [manifest] [outdir]
#   default: tools/oracle/corpus.txt -> .cache/corpus
#   held out: tools/oracle/corpus-holdout.txt -> .cache/corpus-holdout
set -euo pipefail
cd "$(dirname "$0")/../.."
manifest=${1:-tools/oracle/corpus.txt}
out=${2:-.cache/corpus}
mkdir -p "$out" .cache/corpus-tgz
grep -vE '^\s*(#|$)' "$manifest" | while read -r name version; do
  dest="$out/$name@$version"
  [[ -d "$dest" ]] && continue
  meta=$(curl -fsS "https://registry.npmjs.org/${name/\//%2F}/$version")
  tarball=$(printf '%s' "$meta" | python3 -I -c 'import json,sys; print(json.load(sys.stdin)["dist"]["tarball"])')
  integrity=$(printf '%s' "$meta" | python3 -I -c 'import json,sys; print(json.load(sys.stdin)["dist"]["integrity"])')
  tgz=".cache/corpus-tgz/$(echo "$name@$version" | tr '/' '_').tgz"
  curl -fsS -o "$tgz" "$tarball"
  want=${integrity#sha512-}
  got=$(openssl dgst -sha512 -binary "$tgz" | base64 -w0)
  if [[ "$want" != "$got" ]]; then echo "integrity mismatch for $name@$version" >&2; exit 1; fi
  mkdir -p "$dest.tmp"
  # One pattern per call: GNU tar needs --wildcards for patterns and BSD tar
  # (macOS) rejects that option, but both match a lone pattern argument and
  # both fail when a pattern matches nothing.
  for pat in '*.js' '*.mjs' '*.cjs'; do
    tar -xzf "$tgz" -C "$dest.tmp" --no-same-owner --no-same-permissions "$pat" 2>/dev/null || true
  done
  mv "$dest.tmp" "$dest"
  echo "$name@$version: $(find "$dest" -type f | wc -l) files"
done
