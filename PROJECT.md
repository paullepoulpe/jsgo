# Project: Closure Java to Go Production Migration

## Architecture
The migration executes in a dedicated subdirectory `/Users/damien/code/jsgo/transpiled/` to isolate the mechanical translation pipeline, minimal shims, prep steps, and translated compiler binary from the existing codebase.

### Directory Layout
- `transpiled/prep/`: Deterministic mechanical Java-prep patches (PRs 1–13) applied to upstream Closure Java `.cache/closure-src/`.
- `transpiled/shims/`: Clean-room minimal Go shims (~7,075 LOC) including Non-Moving `NodeStore` and `strpool` (1,725 LOC) for high-performance AST representation without pointer-chasing GC overhead.
- `transpiled/translator/`: Simple, mechanical Java-to-Go translator (`j2go`) using OpenJDK 27 standard compiler API (`javax.tools.JavaCompiler`, `com.sun.source.*`) for 1:1 AST mapping, plus modular Go AST post-step passes (e.g., `ArrayList` -> Go slice lowering).
- `transpiled/gen/`: 1:1 mechanically translated Go code corresponding to all 813 non-Ant production Closure Java classes.
- `transpiled/cmd/jsgo/`: CLI driver producing the compiled `jsgo` binary.

## Feature Inventory
| # | Feature | Description | Milestone | Source |
|---|---------|-------------|-----------|--------|
| 1 | Mechanical Java-Prep Patches (PRs 1–13) | Deterministic mechanical patches to de-reflect flags/Gson, desugar streams, and replace complex Guava collections | M1 | Survey |
| 2 | Clean-Room Minimal Go Shims Foundation | High-performance shims based on public specifications, strictly adhering to clean-room rules (no OpenJDK source, no javap, copyleft-clean) | M1 | Survey |
| 3 | Non-Moving NodeStore & StrPool Integration | Integration of zero-GC non-moving node store and interned string pool for Node representation | M1 | Survey |
| 4 | Mechanical Java-to-Go Translator (j2go) | Simple Java AST-to-Go mechanical translator with minimal review surface | M2 | Survey |
| 5 | Modular ArrayList-to-Slice Pass | Independent post-step Go AST transformation lowering ArrayList usage to idiomatic Go slices | M2 | Survey |
| 6 | Full Upstream Codebase Transpilation | Transpilation of all 813 non-Ant production Java classes into transpiled/gen/ | M3 | Survey |
| 7 | CLI Entrypoint & Build Wiring | Standalone transpiled/cmd/jsgo driver with arguments matching Closure Compiler CLI | M3 | Survey |
| 8 | Differential Oracle Parity Verification | Byte-for-byte output parity verification against Closure Java v20261104 via tools/oracle/ and tools/clie2e/ | M4 | Survey |
| 9 | Fast Benchmark Validation (9-Task Suite) | Inner-loop verification using generate-ijs, incremental-ijs-check, downlevel-simple, typed-app, checks-only | M4 | Survey |
| 10 | Primary Performance Gate (googleapis-bench) | Uncontended execution beating Cold Java CLI execution time (208.7s) across all 7,460 chunks | M4 | Survey |
| 11 | Ground-Truth Migration Plan Update | Update migration_plan.md with real measured LOCs, parity percentages, and benchmark numbers | M5 | Survey |
| 12 | Git Delivery & Push to Main | Commit all work under transpiled/ and migration_plan.md, and push to origin/main | M5 | Survey |

## Milestones
| # | Name | Scope | Dependencies | Status |
|---|------|-------|-------------|--------|
| 1 | M1: Prep & Shims | Java-prep patches (PRs 1–13) & minimal clean-room Go shims (including Non-Moving NodeStore) | none | IN_PROGRESS |
| 2 | M2: Translator & Modular Passes | Simple mechanical j2go tool & independent ArrayList-to-slice Go AST pass | M1 | PLANNED |
| 3 | M3: Full Transpilation & CLI | Transpilation of all 813 Java files into transpiled/gen/ and transpiled/cmd/jsgo binary build | M1, M2 | PLANNED |
| 4 | M4: Parity & Performance Gate | Oracle parity verification (tools/clie2e/) and uncontended googleapis-bench beating Cold Java CLI | M3 | PLANNED |
| 5 | M5: Documentation & Git Delivery | Update migration_plan.md with ground-truth numbers and push to origin/main | M4 | PLANNED |

## Interface Contracts
### transpiled/prep ↔ transpiled/translator
- Prep patches modify `.cache/closure-src/` in-place or via staged directory.
- Output: Clean, desugared Java source files compiling with OpenJDK 27 with reflection and complex collections removed.

### transpiled/shims ↔ transpiled/gen
- `transpiled/shims/` exports standard runtime types matching Java APIs (`java.util`, `java.io`, `java.lang`, `com.google.common.collect`).
- `Node` implementation binds directly to Non-Moving `NodeStore` using integer handles.

### transpiled/gen ↔ transpiled/cmd/jsgo
- `transpiled/cmd/jsgo/main.go` invokes `transpiled/gen/com/google/javascript/jscomp.CommandLineRunner.Main(os.Args[1:])`.

## Code Layout
```
/Users/damien/code/jsgo/
├── transpiled/
│   ├── prep/
│   │   ├── patches/             # PRs 1-13 mechanical diffs
│   │   └── apply_prep.sh        # Deterministic patch application script
│   ├── shims/                   # Clean-room Go shims
│   │   ├── nodestore/           # Non-moving NodeStore & strpool
│   │   ├── javalang/            # java.lang basics (jstring, Object, Class)
│   │   ├── javautil/            # java.util collections
│   │   └── guava/               # com.google.common adapters
│   ├── translator/              # j2go tool
│   │   ├── src/                 # Java AST visitor (OpenJDK javax.tools/com.sun.source)
│   │   ├── passes/              # Go AST modular post-step passes (ArrayList -> slice)
│   │   └── run_j2go.sh          # Orchestration script for translator
│   ├── gen/                     # Generated Go packages (813 classes)
│   └── cmd/
│       └── jsgo/
│           └── main.go          # Standalone CLI binary
```
