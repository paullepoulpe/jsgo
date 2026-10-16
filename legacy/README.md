# legacy/ — the first jsgo attempt (frozen)

This directory holds the original jsgo code base: a Go type checker for
Closure-annotated JavaScript developed by iterating against diagnostic
fixtures extracted from `TypeCheckTest.java`. It is kept, building and
tested, for reference only. It is not being developed further; see
[`../PLAN.md`](../PLAN.md) for why and for what replaces it.

Measured state when frozen (commit `c05fc09`, harness in `harness/`):

| Corpus | Default | `strictCheckTypes` |
|---|---|---|
| `testdata/fixtures/extracted` (1444 fixtures, the development set) | 93.70% | 90.17% |

The documents in this directory (`SPEC.md`, `ARCHITECTURE.md`, `TODO.md`,
`docs/`) describe that attempt. Scripts under `scripts/` and `tools/` assume
the old layout and are not maintained; the Java oracles they produced remain
in `../testdata/oracle*/` and are still valid.
