# java2go vs javaport Evaluation

**Date:** 2026-04-24  
**java2go version:** v0.0.0-20251207015311-79857fbb7b4c  
**javaport:** tools/javaport (≈4547 LoC, regex/AST-based)

---

## 1. Maintenance Status

| Metric | Value |
|--------|-------|
| Last commit | 2025-12-07 (≈4.5 months ago) |
| Stars | 13 |
| Forks | 7 |
| Open issues | 3 |
| License | None (issue #8 requests one) |

The repository is technically recent (under 1 year), but has 13 stars, no license, and a
critical unfixed crash bug (see §3). Activity appears to be dependency-update-only rather
than functional improvements.

---

## 2. Summary Table

| File | Java LoC | javaport Go LoC | javaport Go/Java% | javaport gofmt errors | java2go Go LoC | java2go gofmt errors |
|------|----------|-----------------|-------------------|----------------------|----------------|----------------------|
| CheckSideEffects.java | 274 | 199 | 72% | 0 | — (crash) | — |
| CheckJSDoc.java | 799 | 533 | 66% | 5 | — (crash) | — |
| VarCheck.java | 583 | 322 | 55% | 0 | — (crash) | — |
| ControlFlowAnalysis.java | 1050 | 780 | 74% | 14 | — (crash) | — |
| TypeValidator.java | 1303 | 913 | 70% | 433 | — (crash) | — |
| NodeUtil.java | 5851 | 3912 | 66% | 2177 | — (crash) | — |

**java2go produced zero output on all 6 reference files.**

---

## 3. java2go: Critical Crash Bug

### Root cause

java2go uses tree-sitter to parse Java. The tree-sitter Java grammar (current versions) emits
two distinct node types for comments:

- `"line_comment"` — for `// single-line` comments
- `"block_comment"` — for `/* multi-line */` comments

java2go's dispatcher (tree_sitter.go line 275) only handles the legacy unified `"comment"` type
and falls through to `panic("Unknown node type: line_comment")` on any `//` comment in function
body scope. Separately, the `ParseDecls` dispatcher dereferences a nil `*FileScope` (issue #4)
when `--symbols` mode fails to produce a symbol table, causing a second class of panic.

### Reproduction

```
panic: Unknown node type: line_comment
goroutine 1 [running]:
main.ParseNode(...) tree_sitter.go:275
```

Every one of our 6 reference files (and any real Java file with `//` comments in method bodies)
triggers this panic. The `--symbols=false` flag only changes which panic fires first. No output
is written before the crash.

### Affected statements in java2go source

- `statement.go:151` — `if line.Type() == "comment"` (misses `line_comment`)
- `declaration.go:104` — `case "field_declaration", "comment":` (misses `line_comment`)
- `expression.go:24` — `case "comment":` (misses `line_comment`)
- `tree_sitter.go:275` — default `panic(...)` reached for unhandled node types

The fix is mechanical (~10 lines: replace `"comment"` with `"line_comment", "block_comment"` in
every `case`/`if` that mentions it), but it has not been applied in the repository as of the
evaluation date, despite issue #4 being open since ≥2022.

---

## 4. Feature Comparison

### 4.1 Closure-specific needs (javaport vs java2go)

| Feature | javaport | java2go |
|---------|----------|---------|
| `// Port of <path> @ closure-compiler v20260415` header | Yes | No (tool-agnostic) |
| Source-line pins (`// @ path:line`) | Yes | No |
| Import graph mapping (`com.google...jscomp` → `internal/jscomp`) | Yes | No |
| JsDoc preservation as Go `//` comments | Yes | No |
| AutoValue → Go struct + factory | Yes (partial) | No |
| Enum-with-methods → interface + per-constant struct | Yes | N/A (crashes) |
| Go generics for Java type parameters | Yes | N/A (crashes) |

### 4.2 Language constructs

| Construct | javaport | java2go |
|-----------|----------|---------|
| Lambda `(x) -> expr` | Transliterates as `x - > expr` (invalid Go, but preserved for manual fix) | Crashes |
| Switch expression `case X ->` | Transliterates (Go switch body preserved) | Crashes |
| Method refs `Foo::bar` | Drops `::`, emits just `Foo` or transliterates partially | Crashes |
| `new ArrayList<>()` | Preserves in comment, maps to `[]T{}` where type is known | Crashes |
| Guava `ImmutableList.of(...)` | Maps to `[]T{...}` where possible | Crashes |
| `Predicate<Node>` generics | Preserved as `Predicate<Node>` (invalid Go) for manual fix | Crashes |

### 4.3 javaport known issues (observed in run)

- **Lambdas** (`->` in expressions): emitted as `c - > expr` which breaks gofmt. Occurs in
  NodeUtil (406 Java arrows → 3 in output body scope, rest commented out) and ControlFlowAnalysis.
- **Inline block comments in `if` conditions**: Java allows `if (x /* comment */ || y)`;
  javaport preserves the comment inline, breaking Go syntax
  (see CheckJSDoc.go:309 — `if hasNonDescMsgTag / / Don't error...`).
- **Raw generics** in field/variable types emitted as `LinkedIdentityHashMap<Node, Integer>`
  (invalid Go), flagged with `// TODO: Java type ... not mapped`.
- **Constructor bodies** sometimes emit return statements for builder patterns without
  properly scoping the receiver (ControlFlowAnalysis, 14 gofmt errors).
- **Large files**: NodeUtil (5851 LoC Java) → 2177 gofmt parse errors, mostly from lambda
  transliteration and raw generic types cascading into argument-list confusion.

---

## 5. Specific Strengths and Weaknesses

### java2go

**Strengths (theoretical):**
- AST-based approach (tree-sitter) is structurally sounder than regex for expression parsing.
- Has a symbol table / scope resolution mode (`--symbols`, default on) that could, in principle,
  produce better type-aware output.
- Handles basic struct declarations, simple method bodies, and constructor calls on test files
  (its own `testfiles/`).

**Weaknesses (observed):**
- **Crashes on 6/6 real-world Closure Compiler Java files** due to `line_comment` node type.
- Does not strip or translate Java annotations, JsDoc, or import statements to Go equivalents.
- Produces no Closure-specific metadata (headers, line pins, package mapping).
- No license — legally ambiguous for incorporation.
- Effectively unmaintained (13 stars, no real bugfixes in 4.5 months).

### javaport

**Strengths:**
- Runs successfully on 5/6 files (CheckJSDoc.java and CheckSideEffects.java produce 0 gofmt
  parse errors; VarCheck.java produces 0).
- Produces exactly the headers, line pins, and import mappings required by ci-guards.sh.
- Handles Guava collection types, AutoValue, enums-with-methods, and JsDoc preservation.
- Integrated into the project build (Makefile, cache, CI).

**Weaknesses:**
- Lambda expressions (`->`) are transliterated as `x - > expr` — invalid Go, requires manual
  fixup per lambda.
- Inline comments inside Java `if`/`while` conditions bleed into the Go expression.
- Large files with many generics (NodeUtil, TypeValidator) generate thousands of cascading
  gofmt errors from raw `<T>` leakage.

---

## 6. Recommendation

**Decision: `keep-ours`**

java2go is not viable as a replacement or preprocessor for our javaport today. It panics on
every real Closure Compiler source file due to a trivial but unpatched bug in comment node
handling (`"comment"` vs `"line_comment"`/`"block_comment"`). Even if that bug were fixed
upstream, java2go produces no Closure-specific output: no `// Port of` headers, no source-line
pins, no `com.google.javascript.jscomp` → `internal/jscomp` import mapping — all of which are
required by ci-guards.sh and the project porting convention.

The one technique worth studying from java2go is its **tree-sitter-based AST traversal**:
javaport's regex approach mishandles inline comments inside expressions and complex generic
signatures. If javaport's lambda and inline-comment issues become blockers, the right path is
to replace javaport's ad-hoc regex with a proper tree-sitter (or ANTLR) Java parser —
taking inspiration from java2go's dispatcher pattern — rather than adopting java2go wholesale.

For now, the practical fixes that would most reduce javaport's gofmt error rate are:
1. Strip inline `// comments` that appear inside `if`/`while` condition spans (the `hasNonDescMsgTag` class of error).
2. Convert lambda `(x) -> expr` to a `// TODO: lambda: x -> expr` comment block rather than transliterating.

Neither fix warrants a commit to `tools/javaport/main.go` at this time without a broader test
plan; they are recorded here for prioritization.

---

## 7. Sample Outputs

Sample outputs are in `docs/java2go-samples/`:

- `javaport-<File>.go` — actual javaport output for each reference file
- `java2go-<File>.go` — crash report (no Go output produced)

---

## Appendix: Install Notes

```
go install github.com/NickyBoy89/java2go@latest
# Downloaded: v0.0.0-20251207015311-79857fbb7b4c
# Dependency: github.com/smacker/go-tree-sitter v0.0.0-20230113054119-af7e2ef5fed6
# Install time: ~5 seconds, no build errors.
# Runtime: panics immediately on any real Java file.
```

```
# java2go CLI flags:
#   -ast          Print Go AST instead of source
#   -exclude-annotations  Comma-separated list of annotations to exclude
#   -output       Output directory (default ".")
#   -q            Quiet mode
#   -symbols      Symbol-aware mode (default true) — crashes differently with =false
#   -sync         Parse files sequentially
#   -w            Write to disk instead of stdout
```
