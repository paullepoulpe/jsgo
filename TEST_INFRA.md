# End-to-End Testing Infrastructure (TEST_INFRA)

## 1. Test Philosophy

The Closure Compiler Go migration (`jsgo`) follows an **opaque-box, requirement-driven differential testing philosophy**:
- **Ground Truth Specification**: Pinned Closure Compiler Java release **v20261104** (Maven Central jar `closure-compiler-v20261104-maven.jar`) is the absolute authoritative oracle.
- **Differential Parity**: The translated `jsgo` CLI binary is tested against Closure Java ground truth with identical argument vectors, stdin streams, and input file hierarchies.
- **Strict Byte-for-Byte Verification**: Standard output, standard error (with JVM runtime warning headers and internal Java stack frames stripped), process exit codes, and all generated filesystem artifacts (source maps, chunks, renaming reports) must match ground truth.
- **Clean-Room Compliance**: Strict clean-room rules per `PORTING.md` are enforced: zero inspection of OpenJDK source code, zero `javap` inspection of JDK bytecode, and copyleft-clean throughout. All Go test harness and oracle code is built from public specification contracts.
- **Progressive Testability & Isolation**: Every test case sets up its own isolated temporary workspace, operates independently without execution order dependencies, and cleans up upon completion. Tests run reliably in offline mode using pre-recorded verified oracles in `testdata/e2e/manifest.json`, and can also be re-verified differentially online against Java via `--differential`.

---

## 2. 4-Tier Test Design

The test suite is structured into 4 comprehensive tiers covering isolated features, extreme edge conditions, cross-module feature interactions, and large-scale real-world projects.

```
+-------------------------------------------------------------------------------+
|                        4-TIER E2E TEST SUITE ARCHITECTURE                      |
+-------------------------------------------------------------------------------+
| Tier 1: Feature Coverage (45 cases, >=5 per core CLI feature)                 |
|   - Compilation Levels (5)      - Formatting & Output (5) - Source Maps (5)   |
|   - Language In/Out (5)         - Warning Levels (5)      - Type Checking (5) |
|   - Code Motion & Chunking (5)  - Defines (5)             - CLI Meta (5)      |
+-------------------------------------------------------------------------------+
| Tier 2: Boundary & Corner Cases (25 cases, 5 per boundary category)           |
|   - Empty & Trivial Inputs (5)  - Syntax & Malformed (5)  - UTF-16 & WTF-8 (5)|
|   - Deep ASTs & Nesting (5)     - Flag Combinations & Argument Errors (5)     |
+-------------------------------------------------------------------------------+
| Tier 3: Cross-Feature Combinations (15 cases)                                 |
|   - Pairwise & Multi-pass interactions across checking, optimization,         |
|     transpilation, polyfills, module bundling, chunking, and source maps      |
+-------------------------------------------------------------------------------+
| Tier 4: Real-World Application Scenarios (18 cases)                           |
|   - Multi-chunk Closure apps, ES modules (Preact, Three, Moment, Underscore), |
|     CommonJS (Lodash, RxJS, Handlebars), and Benchmark tasks (typed-app, etc.)|
+-------------------------------------------------------------------------------+
```

### Tier 1: Feature Coverage (45 test cases)
Requires $\ge 5$ test cases per feature across all core CLI capabilities:

1. **Compilation Levels (`--compilation_level`)** [5 cases]:
   - `tier1_c01_whitespace`: `WHITESPACE_ONLY` strips comments and spacing without AST modification.
   - `tier1_c02_simple`: `SIMPLE` performs constant folding, local variable renaming, and dead code elimination.
   - `tier1_c03_advanced`: `ADVANCED` performs aggressive cross-module renaming, global symbol mangling, and function inlining.
   - `tier1_c04_bundle`: `BUNDLE` concatenates and packages input scripts without minification.
   - `tier1_c05_transpile_only`: `TRANSPILE_ONLY` downlevels modern syntax to ES5 without symbol renaming.

2. **Formatting & Output Options** [5 cases]:
   - `tier1_f01_pretty_print`: `--formatting PRETTY_PRINT` produces indented, multi-line readable output.
   - `tier1_f02_print_input_delimiter`: `--formatting PRINT_INPUT_DELIMITER` emits delimiter comments between input files.
   - `tier1_f03_single_quotes`: `--formatting SINGLE_QUOTES` forces single quotes for emitted string literals.
   - `tier1_f04_output_wrapper`: `--output_wrapper "(function(){%output%})();"` wraps generated JS in an IIFE.
   - `tier1_f05_output_file`: `--js_output_file dist/app.js` writes generated output directly to specified path.

3. **Source Maps (`--create_source_map`)** [5 cases]:
   - `tier1_s01_basic_map`: Basic V3 JSON source map with valid version, file, sources, and VLQ mappings.
   - `tier1_s02_multi_input_map`: Source map spanning multiple input files accurately mapping source origins.
   - `tier1_s03_pretty_map`: Source map generated with `--formatting PRETTY_PRINT`.
   - `tier1_s04_wrapped_map`: Source map with `--output_wrapper` correctly applying line offsets to mappings.
   - `tier1_s05_advanced_map`: Source map generated under `ADVANCED` optimizations with renamed identifier mappings.

4. **Language In & Language Out** [5 cases]:
   - `tier1_l01_esnext_to_es5`: Transpilation of ES_NEXT (arrow functions, classes, templates) down to ECMASCRIPT5.
   - `tier1_l02_esnext_to_es2015`: Transpilation down to ECMASCRIPT_2015 preserving classes while downleveling newer syntax.
   - `tier1_l03_no_transpile`: `--language_out NO_TRANSPILE` preserving ES_NEXT modern constructs unmodified.
   - `tier1_l04_es5_strict`: `--language_in ECMASCRIPT5_STRICT --language_out ECMASCRIPT5` with `'use strict';` directive.
   - `tier1_l05_optional_chaining`: Transpilation of ECMASCRIPT_2020 optional chaining (`a?.b?.c?.()`) to ES5 conditionals.

5. **Warning Levels & Diagnostic Groups** [5 cases]:
   - `tier1_w01_quiet`: `--warning_level QUIET` suppresses all non-fatal warning diagnostics.
   - `tier1_w02_verbose`: `--warning_level VERBOSE` produces comprehensive diagnostic reports.
   - `tier1_w03_error_group`: `--checks_only --jscomp_error=checkTypes` turns type warnings into fatal exit code 1 errors.
   - `tier1_w04_warning_group`: `--jscomp_warning=checkVars` emits diagnostic warnings for undeclared variable references.
   - `tier1_w05_off_group`: `--jscomp_off=uselessCode` suppresses unreachable statement warnings.

6. **Static Checking & Type Inference** [5 cases]:
   - `tier1_t01_var_mismatch`: Detects assignment of string literal to `@type {number}` variable.
   - `tier1_t02_param_mismatch`: Detects function invocation passing invalid parameter type.
   - `tier1_t03_return_mismatch`: Detects return statement violating declared `@return {number}` contract.
   - `tier1_t04_interface_compliance`: Verifies structural type compliance against declared `@record` interface.
   - `tier1_t05_checks_only_clean`: `--checks_only` with clean typed code succeeds with code 0 and emits zero JS output.

7. **Code Motion & Chunking (`--chunk`)** [5 cases]:
   - `tier1_k01_two_chunks_linear`: 2-chunk linear graph (`base:1` and `feat:1:base`).
   - `tier1_k02_three_chunks_diamond`: 3-chunk diamond graph (`root:1`, `left:1:root`, `right:1:root`).
   - `tier1_k03_chunk_wrapper`: `--chunk_output_wrapper "/*chunk*/%output%"` applied per emitted chunk file.
   - `tier1_k04_code_motion`: Helper function used across leaf chunks automatically moved to common ancestor chunk.
   - `tier1_k05_chunk_cycle_error`: Cyclic chunk dependency (`a:1:b` and `b:1:a`) detected and reported with error exit.

8. **Defines & Compile-Time Constants (`--define`, `-D`)** [5 cases]:
   - `tier1_d01_bool_define`: Boolean `@define` overridden with `--define=DEBUG_FLAG=false`.
   - `tier1_d02_str_define`: String `@define` overridden with `-D BUILD_ENV='production'`.
   - `tier1_d03_num_define`: Numeric `@define` overridden with `--define=MAX_RETRIES=5`.
   - `tier1_d04_multi_define`: Multiple simultaneous defines across boolean, string, and number types.
   - `tier1_d05_implicit_bool`: Boolean define flag specified without explicit value (`-D ENABLE_METRICS`).

9. **CLI Meta Flags & Information** [5 cases]:
   - `tier1_m01_version`: `--version` prints Closure Compiler version string (`v20261104`) and exits 0.
   - `tier1_m02_help`: `--help` prints CLI option usage synopsis and exits 0.
   - `tier1_m03_help_markdown`: `--help_markdown` emits markdown formatted documentation table and exits 0.
   - `tier1_m04_unknown_flag`: Unknown CLI flag produces invalid option diagnostic and exits with code 255.
   - `tier1_m05_missing_arg`: Flag requiring parameter passed without argument produces error and exits non-zero.

---

### Tier 2: Boundary & Corner Cases (25 test cases)
Covers boundary conditions across 5 specific risk areas:

1. **Empty & Trivial Inputs** [5 cases]:
   - `tier2_b01_empty_file`: 0-byte input file compiles cleanly to empty output, exit code 0.
   - `tier2_b02_whitespace_only`: Input containing only whitespace (spaces, tabs, newlines) compiles cleanly, exit code 0.
   - `tier2_b03_comments_only`: Input containing only single-line and multiline comments compiles cleanly, exit code 0.
   - `tier2_b04_empty_with_output_file`: 0-byte input with `--js_output_file` produces 0-byte output file.
   - `tier2_b05_multiple_empty`: Multiple concatenated 0-byte input files compile cleanly, exit code 0.

2. **Syntax Errors & Malformed Inputs** [5 cases]:
   - `tier2_b06_unclosed_paren`: Unclosed parenthesis / expression fails compilation with parse error, exit code 1.
   - `tier2_b07_illegal_token`: Illegal stray non-ASCII / operator characters fail with parse error, exit code 1.
   - `tier2_b08_reserved_word_strict`: Using reserved keyword (`let`) as variable name in strict mode fails, exit code 1.
   - `tier2_b09_duplicate_param_strict`: Duplicate parameter declaration in strict mode function fails, exit code 1.
   - `tier2_b10_unclosed_string`: Unclosed quote in string literal fails compilation with parse error, exit code 1.

3. **UTF-16 Surrogate Pairs & Unicode** [5 cases]:
   - `tier2_b11_surrogate_pair_string`: Supplementary plane emoji (`😀🎉🚀`, surrogate pairs) preserved in string literals.
   - `tier2_b12_lone_high_surrogate`: Lone high surrogate (`\uD800`) handled cleanly without panic or WTF-8 corruption.
   - `tier2_b13_lone_low_surrogate`: Lone low surrogate (`\uDC00`) handled cleanly without panic or WTF-8 corruption.
   - `tier2_b14_unicode_identifiers`: Unicode identifier names (Japanese `変数`, Greek `π`, Cyrillic `результат`) parsed and compiled.
   - `tier2_b15_unicode_regex`: Regular expressions with Unicode range escapes (`/[\u00A0-\uD7FF]/gu`) compiled correctly.

4. **Deep ASTs & Extreme Structures** [5 cases]:
   - `tier2_b16_deep_binary_expr`: 250 levels of nested binary addition `(1 + 1 + ...)` compiled without stack overflow.
   - `tier2_b17_deep_parentheses`: 250 levels of nested parentheses `((((... 42 ...))))` compiled without stack overflow.
   - `tier2_b18_deep_blocks`: 60 levels of nested block scopes `{{{ ... }}}` compiled without stack overflow.
   - `tier2_b19_deep_function_chain`: 80 chained method invocations `builder.step().step()...` compiled without stack overflow.
   - `tier2_b20_deep_ternary`: 50 nested ternary conditional expressions `a ? 1 : (b ? 2 : ...)` compiled without stack overflow.

5. **Flag Combinations & Argument Errors** [5 cases]:
   - `tier2_b21_nonexistent_input`: Missing input file reports clear file not found diagnostic and exits non-zero.
   - `tier2_b22_nonexistent_flagfile`: Missing `@flagfile` reports error and exits non-zero.
   - `tier2_b23_invalid_compilation_level`: Unknown compilation level name reports error and exits non-zero.
   - `tier2_b24_invalid_charset`: Unrecognized charset name reports error and exits non-zero.
   - `tier2_b25_conflicting_flags`: Contradictory flag combination (`--isolation_mode IIFE` with `--output_wrapper`) handled with standard error exit.

---

### Tier 3: Cross-Feature Combinations (15 test cases)
Covers pairwise and multi-stage interactions:
- `tier3_x01_typecheck_and_advanced`: Type declarations (`@constructor`, `@type`) guide property renaming under `ADVANCED`.
- `tier3_x02_typecheck_halting_advanced`: Type error with `--jscomp_error=checkTypes` halts compilation before optimization phases execute.
- `tier3_x03_chunking_with_sourcemaps`: Multi-chunk compilation generating separate valid source maps per chunk file.
- `tier3_x04_chunking_codemotion_sourcemaps`: Cross-chunk code motion maintains mapping accuracy when moving functions to parent chunks.
- `tier3_x05_chunking_with_wrappers`: Chunk output wrappers applied across multi-chunk compilation outputs.
- `tier3_x06_transpile_polyfills_advanced`: Modern ES classes/promises transpiled to ES5 with polyfill injection under `ADVANCED`.
- `tier3_x07_transpile_isolate_polyfills`: Transpilation with isolated polyfills and output wrapper preventing global scope leakage.
- `tier3_x08_esm_bundling_sourcemaps`: Multiple ES modules (`import`/`export`) bundled into single output with source maps.
- `tier3_x09_cjs_to_esm_advanced`: CommonJS (`module.exports`/`require`) module transpilation and minification under `ADVANCED`.
- `tier3_x10_renaming_reports_pretty_wrapper`: Variable & property renaming reports emitted alongside pretty-printed wrapped output.
- `tier3_x11_strict_checks_and_defines`: Strict type checking (`--jscomp_warning=strictCheckTypes`) combined with compile-time defines.
- `tier3_x12_dynamic_imports_chunks`: Dynamic `import()` syntax driving chunk boundaries.
- `tier3_x13_language_out_formatting_wrapper`: Transpilation to ES5 combined with `SINGLE_QUOTES`, `PRETTY_PRINT`, and output wrapper.
- `tier3_x14_defines_inlining_deadcode`: False boolean define causes dead branch and its helper functions to be completely removed.
- `tier3_x15_externs_typecheck_advanced`: Custom externs file type-checked and preserved from renaming under `ADVANCED`.

---

### Tier 4: Real-World Application Scenarios (18 test cases)
Covers real projects from `tools/clie2e/projects/` and benchmark tasks from `tools/bench/tasks/`:
- `tier4_r01_closure_cli_globs`: Multi-file build using compiler-expanded globs (`closure-cli.json`).
- `tier4_r02_closure_cli_flagfile`: Compilation with options loaded from `@flagfile` (`closure-cli.json`).
- `tier4_r03_closure_cli_translations`: Compilation with `--translations_file` and XML translation bundles (`closure-cli.json`).
- `tier4_r04_closure_app_simple`: Closure Library application with `goog.provide`/`goog.require` under `SIMPLE` (`closure-app.json`).
- `tier4_r05_closure_app_advanced`: Closure Library application with namespace collapsing under `ADVANCED` (`closure-app.json`).
- `tier4_r06_chunks_closure`: Multi-chunk Closure Library application build (`chunks-closure.json`).
- `tier4_r07_chunks_modules`: Multi-chunk ES module application graph build (`chunks-modules.json`).
- `tier4_r08_esm_preact`: Preact ES module component scenario (`esm-preact.json`).
- `tier4_r09_esm_moment`: Moment.js date helper module scenario (`esm-moment.json`).
- `tier4_r10_esm_three`: Three.js 3D vector math module scenario (`esm-three.json`).
- `tier4_r11_esm_underscore`: Underscore collection utility pipeline scenario (`esm-underscore.json`).
- `tier4_r12_cjs_lodash`: CommonJS Lodash clamp utility module scenario (`cjs-lodash.json`).
- `tier4_r13_cjs_rxjs`: CommonJS RxJS observable pattern scenario (`cjs-rxjs.json`).
- `tier4_r14_cjs_handlebars`: CommonJS Handlebars template rendering scenario (`cjs-handlebars.json`).
- `tier4_r15_bench_typed_app`: Benchmark task `typed-app` (`tools/bench/tasks/typed-app.json`).
- `tier4_r16_bench_checks_only`: Benchmark task `checks-only` (`tools/bench/tasks/checks-only.json`).
- `tier4_r17_bench_downlevel_simple`: Benchmark task `downlevel-simple` (`tools/bench/tasks/downlevel-simple.json`).
- `tier4_r18_typecheck_acceptance`: Representative acceptance fixture from `testdata/fixtures/extracted/`.

---

## 3. Test Runner Harness Architecture

The test harness is implemented in `tools/clie2e/` with two native execution pathways:

1. **CLI Test Runner (`tools/clie2e`)**:
   ```bash
   go run ./tools/clie2e -suite tiers [-tier 1,2,3,4|all] [-case <id>] [-tier-jobs N] [-jsgo <bin>] [-differential]
   ```
2. **Go Test Framework Integration (`go test`)**:
   ```bash
   go test -v ./tools/clie2e -run TestTiers
   go test -v ./tools/clie2e -run TestTier1FeatureCoverage
   go test -v ./tools/clie2e -run TestTier2BoundaryCornerCases
   go test -v ./tools/clie2e -run TestTier3CrossFeatureCombinations
   go test -v ./tools/clie2e -run TestTier4RealWorldScenarios
   ```

### Binary Selection Hierarchy
When executing against the translated compiler, the harness resolves the target `jsgo` binary in the following priority order:
1. Explicit CLI argument: `-jsgo <path_to_binary>`
2. Environment variable: `JSGO_BIN=<path_to_binary>`
3. Mechanical migration binary: `transpiled/cmd/jsgo` (built into workspace)
4. Existing baseline binary: `./cmd/jsgo` (fallback)

### Pass / Fail Semantics & Exit Codes
- **Accurate Exit Code Handling**:
  - A test passes if and only if **all** expected attributes match: exit status, standard output, standard error, and output files.
  - If any test case diverges, the test runner prints the exact diff and terminates with exit code **1**.
  - Exit code **0** is returned if and only if 100% of selected test cases pass.
  - Exit code **2** is reserved for CLI argument or harness configuration errors.

---

## 4. Verification and Parity Matrix

All 103 test cases have been verified against Closure Compiler Java v20261104 ground truth:

| Tier | Category | Cases | Status | Ground Truth Source |
|------|----------|-------|--------|---------------------|
| Tier 1 | Compilation Levels | 5 | PASS (5/5) | Java v20261104 CLI |
| Tier 1 | Formatting & Output | 5 | PASS (5/5) | Java v20261104 CLI |
| Tier 1 | Source Maps | 5 | PASS (5/5) | Java v20261104 CLI |
| Tier 1 | Language In / Out | 5 | PASS (5/5) | Java v20261104 CLI |
| Tier 1 | Warning Levels | 5 | PASS (5/5) | Java v20261104 CLI |
| Tier 1 | Type Checking | 5 | PASS (5/5) | Java v20261104 CLI |
| Tier 1 | Chunking | 5 | PASS (5/5) | Java v20261104 CLI |
| Tier 1 | Defines | 5 | PASS (5/5) | Java v20261104 CLI |
| Tier 1 | CLI Meta Flags | 5 | PASS (5/5) | Java v20261104 CLI |
| **Tier 1 Total** | **Feature Coverage** | **45** | **PASS (45/45)** | **100.0% Parity** |
| Tier 2 | Empty & Trivial | 5 | PASS (5/5) | Java v20261104 CLI |
| Tier 2 | Syntax Errors | 5 | PASS (5/5) | Java v20261104 CLI |
| Tier 2 | UTF-16 & Unicode | 5 | PASS (5/5) | Java v20261104 CLI |
| Tier 2 | Deep ASTs | 5 | PASS (5/5) | Java v20261104 CLI |
| Tier 2 | Flags & Options | 5 | PASS (5/5) | Java v20261104 CLI |
| **Tier 2 Total** | **Boundary & Corner Cases** | **25** | **PASS (25/25)** | **100.0% Parity** |
| **Tier 3 Total** | **Cross-Feature Combinations**| **15** | **PASS (15/15)** | **100.0% Parity** |
| **Tier 4 Total** | **Real-World Scenarios** | **18** | **PASS (18/18)** | **100.0% Parity** |
| **GRAND TOTAL** | **Full E2E Suite** | **103** | **PASS (103/103)** | **100.0% Parity** |
