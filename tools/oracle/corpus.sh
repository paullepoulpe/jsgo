#!/usr/bin/env bash
# Fetches the pinned real-world JS corpus listed in corpus.txt into
# .cache/corpus/. Idempotent; verifies npm sha512 integrity.
set -euo pipefail
cd "$(dirname "$0")/../.."
out=.cache/corpus
mkdir -p "$out" .cache/corpus-tgz
grep -vE '^\s*(#|$)' tools/oracle/corpus.txt | while read -r name version; do
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
  tar -xzf "$tgz" -C "$dest.tmp" --no-same-owner --no-same-permissions --wildcards '*.js' '*.mjs' '*.cjs' 2>/dev/null || true
  mv "$dest.tmp" "$dest"
  echo "$name@$version: $(find "$dest" -type f | wc -l) files"
done
