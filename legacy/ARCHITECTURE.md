# jsgo Architecture

This doc explains _how_ the Go port is shaped and _why_ it's shaped that
way. For the verifiable "done" contract, see [SPEC.md](SPEC.md); for
the ranked port gap list, see [TODO.md](TODO.md).

## Why a port, not a rewrite

`jsgo` mirrors Google Closure Compiler's type-checker architecture as
closely as Go idioms allow. The guiding rule (also called out in
SPEC.md) is:

> **Port, don't redesign.** Keep Java class names, method names,
> file layout, and pass structure recognisable to someone familiar
> with closure-compiler, even when that means a less idiomatic Go
> design.

Reasons:

1. The Java checker is the oracle. Matching its structure means bugs
   surface as behavioural divergence at the same spot in the Go code
   as in the Java. A rewrite risks re-inventing bugs.
2. The type system is subtle (structural + nominal + unions +
   templates + dataflow narrowing). A rewrite would bikeshed the
   lattice forever; a port lets us adopt Java's answers wholesale.
3. Every `// Port of` traceability header is a one-click jump to the
   Java source we're matching, which is invaluable when behaviour
   drifts.

Shortcuts are allowed when the scope is clearly narrower than Java —
but the port header must say so with a qualifier like `(subset)` /
`(stub)` / `(shim)`. `scripts/ci-guards.sh` enforces the header.

## Code layout

```
jsgo/
├── SPEC.md                              # MVP contract + guardrails
├── ARCHITECTURE.md                      # this file
├── TODO.md                              # ranked remaining port work
│
├── cmd/jsgo/main.go                     # CLI: `jsgo typecheck <file>`
│
├── internal/
│   ├── rhino/                           # com.google.javascript.rhino
│   │   ├── node.go                      #   Node.java (subset)
│   │   ├── token.go                     #   Token.java (subset)
│   │   ├── jsdoc.go                     #   JSDocInfo.java (subset)
│   │   ├── jstype_expr.go               #   JSTypeExpression.java
│   │   └── jstype/                      #   .rhino.jstype.*
│   │       ├── jstype.go                #     JSType.java (+ primitives)
│   │       ├── function_type.go         #     FunctionType.java (subset)
│   │       ├── object_type.go           #     Prototype/InstanceObjectType.java
│   │       ├── union_type.go            #     UnionType.java
│   │       └── registry.go              #     JSTypeRegistry.java (subset)
│   │
│   ├── jscomp/                          # com.google.javascript.jscomp
│   │   ├── compiler.go                  #   Compiler.java (shim)
│   │   ├── scope.go                     #   TypedScopeCreator.java (subset)
│   │   ├── typecheck.go                 #   TypeCheck.java (subset)
│   │   ├── inference.go                 #   TypeInference.java (stub)
│   │   ├── nodetraversal.go             #   NodeTraversal.java (subset)
│   │   └── parsing/                     #   .jscomp.parsing.*
│   │       ├── jsdoctoken.go            #     JsDocToken.java
│   │       ├── jsdoctokenstream.go      #     JsDocTokenStream.java
│   │       └── jsdocinfoparser.go       #     JsDocInfoParser.java
│   │
│   ├── parser/                          # .jscomp.parsing.parser
│   │   ├── scanner.go                   #   parser/Scanner.java (subset)
│   │   └── parser.go                    #   parser/Parser.java (narrow)
│   │
│   └── diag/
│       └── diag.go                      # JSError.java + DiagnosticType.java
│
├── tools/extract_fixtures/              # TypeCheckTest.java → fixtures/
├── harness/harness_test.go              # diff jsgo tuples vs Java oracle
├── scripts/
│   ├── gen-oracles.sh                   # parallel oracle regeneration
│   └── ci-guards.sh                     # enforces SPEC guardrails 1-3
│
└── testdata/
    ├── fixtures/{tier0,tier1,tier2,tier3,extracted}/
    └── oracle/{tier0,tier1,tier2,tier3,extracted}/
```

Each `internal/` Go file that claims to be a port carries a
`// Port of <java-file> @ closure-compiler v20260415 (qualifier)`
header. The qualifier flags intentional scope narrowing — `(subset)`,
`(stub)`, `(shim)`, `(narrow; …)`.

## Pass pipeline

The compiler runs three ordered passes (plus parse):

```
  source (.js text)
    │
    │   parser.Parse  ← internal/parser (Scanner + Parser)
    ▼
  rhino.Node tree (SCRIPT)
    │
    │   TypedScopeCreator.CreateScope    ← internal/jscomp/scope.go
    ▼
  TypedScope tree + FunctionType / ObjectType registry
    │
    │   TypeCheck.Check via NodeTraversal ← internal/jscomp/{typecheck,nodetraversal}.go
    ▼
  []diag.Diagnostic
```

`internal/jscomp/compiler.go` is a ~30-line shim that wires these
together for the CLI and test harness. A proper port of Java's
`Compiler.java` (PhaseOptimizer, pass registry, Result, module graph)
lands when we have more than these three passes.

## Oracle + test harness

- The Java Closure Compiler at a pinned SHA (currently `v20260415`)
  is the ground truth.
- `scripts/gen-oracles.sh` runs Java over every fixture and caches
  the stderr output in `testdata/oracle/`. Oracle files are checked
  in so the normal `go test` path doesn't require a JVM.
- `harness/harness_test.go` walks every fixture, runs `jscomp.CheckSource`
  on it, extracts `(path, line, col, code)` diagnostic tuples from
  both the cached oracle and jsgo's output, and asserts they match.
- Tier-0..Tier-3 corpora are hand-authored and mandatory (100%
  required). The `extracted/` corpus is 1,444 fixtures pulled from
  `com/google/javascript/jscomp/TypeCheckTest.java` by
  `tools/extract_fixtures`; it's the benchmark reported in
  SPEC.md#current-status.

## Guardrails (enforced; see SPEC.md)

1. Pure-Go runtime — no cgo, `os/exec`, `net/http`, etc.
2. No hardcoded oracle data in `internal/`.
3. Every `internal/` Go file carries a `// Port of …` header.
4. Diagnostics compared as `(path, line, col, code)` tuples.
5. Hold-out corpus slot reserved (to be populated later).
6. Novel-input sanity: the parser never panics.
7. Performance window: jsgo within [0.1×, 2×] of Java wall-clock.

`scripts/ci-guards.sh` covers #1, #2, #3 today.

## Current gaps vs. a full port

See [TODO.md](TODO.md) for the ranked list and effort estimates.
Briefly, the biggest three are `TypeInference.java` (real dataflow
narrowing), `Parser.java` (token-stream grammar, more JS constructs),
and a proper externs-loader with `DefaultExterns.js` content.
