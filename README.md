# jsgo

Experimental Go port of Google Closure Compiler's type checker, inspired
by tsgo. See [SPEC.md](SPEC.md) for the MVP contract.

Status: **scaffold milestone** — one Tier-0 fixture passes end-to-end.

```
go build ./cmd/jsgo
./jsgo typecheck testdata/fixtures/tier0/001_type_mismatch_var.js
go test ./...
```
