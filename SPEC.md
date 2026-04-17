# jsgo — MVP Specification

This file is the contract for the jsgo MVP. It is intentionally small and
strict so that "done" is objectively verifiable and hard to game.

## Why

Explore whether Closure Compiler's type checker (Java) can be ported to Go,
in the spirit of tsgo (TypeScript compiler → Go). The MVP's purpose is to
prove the port approach is viable by getting a minimal end-to-end slice
working under guardrails that prevent shortcuts from claiming victory.

## MVP Goal

`jsgo typecheck <file.js>` emits diagnostics that match a pinned release of
the Java Closure Compiler
(`java -jar closure-compiler.jar --checks_only --jscomp_error=checkTypes
--warning_level=VERBOSE --language_in=ECMASCRIPT_NEXT --js <file.js>`).

"Match" means equal sets of `(path, line, col, error_code)` tuples plus an
equal exit-code class (zero vs non-zero). Message strings are NOT part of
the conformance check, but jsgo should still print Java-shaped output so
that dogfooding feels right.

## Pinned Oracle

- Closure Compiler version: **v20260415** (Maven Central)
- Jar SHA256: `20313094c6e9cac8c685931e0b3fd67f82e1534c30448cdff404ee350460e391`
- Corpus flags: `--checks_only --jscomp_error=checkTypes
  --warning_level=VERBOSE --language_in=ECMASCRIPT_NEXT`
- Cached oracle outputs live in `testdata/oracle/` and are checked in, so
  routine `go test` does not require Java.

## Hard Guardrails

These are enforced by `scripts/ci-guards.sh` (invoked from CI), not honor
system.

1. **Pure Go at runtime**: `cmd/jsgo` and `internal/` may not use `cgo`,
   `os/exec`, or the `net` package at runtime. `tools/` and `harness/` (test
   support) are exempt.
2. **No hardcoded oracle data**: no `map[input]→expected_output` lookup
   tables in `internal/`. The guard greps fixture paths/bodies against
   non-testdata source.
3. **Source traceability**: every non-trivial file under `internal/`
   begins with a header comment
   `// Port of <relative/path/File.java> @ closure-compiler <tag>` naming a
   real file in `github.com/google/closure-compiler`. The guard lints for
   this header.
4. **Tuple-level equivalence**, not byte equivalence. Prevents both
   message-string tarpits and "print whatever matches."
5. **Hold-out corpus** (populated at corpus freeze, empty at scaffold
   time): 20% of fixtures live in `testdata/holdout/` and are not looked
   at during development. MVP "done" requires hold-out at same pass-rate
   as visible corpus.
6. **Novel-input sanity**: jsgo must not panic on arbitrary JS input even
   when diagnostics are wrong.
7. **Performance window**: Tier-1 corpus wall-clock within `[0.1×, 2×]`
   Java baseline. Outside that window, something's wrong (either not doing
   the work, or pathological).

## Tiered Corpus

Fixtures are extracted from closure-compiler's own tests
(`com/google/javascript/jscomp/TypeCheckTest.java` and related), then
curated.

| Tier | Count | Scope |
|------|-------|-------|
| 0    | ~10   | primitives, `@type`, simple assignment mismatches |
| 1    | ~30   | `@param`/`@return`, calls, `@constructor`, nullability |
| 2    | ~50   | `@template`, `@record`, unions, subtype rules |
| Stretch | — | `closure-library/closure/goog/base.js` warning-count match |

## Exit Criteria for "MVP-complete"

- Tier 0 + Tier 1: **100%** tuple match on visible and hold-out.
- Tier 2: **≥80%** tuple match on visible and hold-out.
- All 7 guardrails green.
- Performance in `[0.1×, 2×]` on Tier-1.

The current scaffold milestone is much smaller: **one Tier-0 fixture
passing end-to-end** (`testdata/fixtures/tier0/001_type_mismatch_var.js`).

## Port Style

- Java `class` → Go `struct` + methods on pointer receivers.
- Java abstract class → Go interface + shared embedded struct.
- `extends` → struct embedding; `super.foo()` → explicit `base.foo(v)`.
- `null` → Go pointer / typed nil, matching Java null-checks 1:1.
- Java `enum` → Go `type X int` + constants.
- Closure's `JSType` subclass hierarchy → Go interface + type switch.
- Visitor (`NodeTraversal.Callback`) → Go interface with same methods.
- Do **not** introduce channels/goroutines/generics until parity is proven.

## Out of Scope

- Tiers 1 and 2 implementation.
- Full Rhino parser port (current parser is tracer-bullet sized).
- All non-type-checking `jscomp` passes.
- Source maps, bundling, emit.
- IDE / language-server integration.
