# jsgo

A Go port of Google Closure Compiler's type checker, inspired by
tsgo (TypeScript → Go). jsgo aims to produce the same diagnostics
as `closure-compiler --checks_only` on the same source — modulo
message strings; only `(path, line, col, code)` tuples are part of
the conformance check.

Currently at **722/1444 = 50.00%** on an extracted benchmark of
`TypeCheckTest.java` cases; Tier-0..Tier-3 hand-authored corpora
are all at 100%.

## Documents

- [SPEC.md](SPEC.md) — the MVP contract, tiered corpora, and
  enforced guardrails.
- [ARCHITECTURE.md](ARCHITECTURE.md) — port philosophy, code layout
  mapping to Java, pass pipeline, harness.
- [TODO.md](TODO.md) — ranked list of remaining Java ports needed
  to reach 80%+ on the extracted benchmark.

## Usage

```sh
go build -o bin/jsgo ./cmd/jsgo
./bin/jsgo typecheck testdata/fixtures/tier0/001_type_mismatch_var.js
```

## Development

```sh
make all                    # build + go test + ci-guards
make oracle                 # regenerate cached Java oracles (needs JDK)
go run ./tools/extract_fixtures <TypeCheckTest.java> testdata/fixtures/extracted
```

Guardrails are enforced by `scripts/ci-guards.sh` — pure-Go runtime,
no oracle leakage into source, every `internal/` Go file traces to
a real closure-compiler Java source. See [SPEC.md](SPEC.md) for the
full list.
