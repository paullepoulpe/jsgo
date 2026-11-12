# Post-Flip Go Optimizations: Necessity, Interactions & Empirical ROI Report

Companion report to [`migration_plan.md`](migration_plan.md) (§8). All numbers in this document were measured sequentially (uncontended) on Apple M2 (`Mac14,2`, 24 GiB RAM, OpenJDK `27.0.1` vs. Go `1.27.3`) with **`100%` byte-for-byte output & sourcemap `MATCH`** across all `7,460` `googleapis-bench` chunks (`SHA-256: 7968e63b5cd96980...`) and all 9 standard benchmark tasks.

---

## 1. Executive Summary & Necessity Verdict

> [!IMPORTANT]
> **Zero (`0` of `23`) post-flip optimizations are required for the cutover.**
> Right at the flip (`L0_postflip_true_t4` — Non-Moving [`NodeStore`](internal/rhino/node_store.go) + upstream Java's existing `PrebuildAst` / `PrebuildDependencyInfo` only, with **all** Go-only pass parallelism, cross-pass caches, traversal bitmasks, slabs, `FastGzip`, and AST compaction disabled):
> - **`googleapis-bench` (`-t4`, `43,917` files, `35.6M` AST nodes, `7,460` chunks)**: Go is already **`-38.3%` faster** (`128.71s` vs. `208.69s` Java `-t4`) and uses **`-48.2%` less Peak RAM** (`6,536 MB` vs. `12,610 MB`).
> - **9-Task Suite (`-t4`)**: Go is already **`-65.0%` faster than Cold Java CLI** (`28.74s` vs. `82.17s` sum, winning `9/9` tasks) and **`-34.6%` faster than JIT-warmed Java** (`28.74s` vs. `43.97s` sum).

Between the post-flip baseline (`L0_postflip_true_t4`) and current `HEAD` (`v1104`), `jsgo` contains **23 distinct Go-only optimizations (`~29,955` gross Go-only LOC)** that collectively reduce post-flip wall time by **`-36.1%` on `googleapis-bench` (`128.71s → 82.20s`)** and **`-45.9%` on the 9-task suite (`28.74s → 15.54s`)**:

| Tier | Ranks | Go-Only LOC | `googleapis-bench` Effect vs. Post-Flip (`L0`) | 9-Task Suite Effect vs. Post-Flip (`L0`) | Share of Total `googleapis` Gain | Historical Bug Commits | Verdict |
|---|---|---:|---:|---:|---:|---:|---|
| **Tier 1 — Recommended** | `#1`–`#12` | `14,835` gross (`14,107` net) | **`-31.9%`** (`-41.07s`) | **`-34.6%`** (`-9.94s`) | **`88.3%`** | `13` *(mostly localized)* | **Implement selectively post-flip** |
| **Tier 2 — Optional** | `#14`–`#17` | `5,989` | **`-3.1%`** (`-4.00s`) | **`-6.6%`** (`-1.91s`) | `8.6%` | `5` *(core parser/scope/12-pass hooks)* | **Optional; invasive or niche** |
| **Tier 3 — Skip** | `#13`, `#18`–`#23` | `9,131` | **`+0.3%` net slower** *(`-1.6%` w/o `#23`)* | **`-3.9%`** (`-1.18s`) | `<3.5%` *(negative w/ `#23`)* | **`15`** *(races, OOMs, stale caches)* | **Do not re-implement** |

---

## 2. Quick-Scan Ranking: One Bullet Per Optimization (Percentages First)

All percentages below are **wall-time change relative to the post-flip Go baseline (`L0`: `128.71s` `googleapis`, `28.74s` 9-task)** so every bullet shares a single denominator and sums directly. Ordered strictly from highest ROI to lowest ROI:

### Tier 1 — Recommended (`-31.9%` `googleapis` / `-34.6%` 9-task; `88%` of speedup for `<50%` of LOC)

- **`-0.6%` to `-1.8%` `googleapis` / `-0.3%` 9-task (`0` LOC)** — **#1 Go PGO** ([`cmd/jsgo/default.pgo`](cmd/jsgo/default.pgo)): Profile-guided build flag applied automatically by `go build`; zero Go code or bug surface.
- **`-0.04%` `googleapis` / `-0.1%` 9-task (`-728` LOC)** — **#2 Ryū `dtoa`** ([`internal/dtoa/d_to_a.go`](internal/dtoa/d_to_a.go)): Replaces `1,053`-LOC port of Rhino `DToA.java` with a self-contained `325`-LOC leaf package (`65,675/65,675` Rhino vectors).
- **`-3.5%` `googleapis` / `-2.7%` 9-task (`418` LOC)** — **#3 Parallel `TypeCheck`** ([`type_check_parallel.go`](internal/jscomp/type_check_parallel.go), [`typed_scope_creator_subbuilders.go`](internal/jscomp/typed_scope_creator_subbuilders.go)): Checks disjoint top-level scopes concurrently (`-60.8%` on `checkTypes` itself).
- **`-6.0%` `googleapis` / `-7.0%` 9-task (`1,217` LOC)** — **#4 Local-Scope Mid-End Parallelism** ([`inline_variables_parallel.go`](internal/jscomp/inline_variables_parallel.go), [`dead_assignments_elimination.go`](internal/jscomp/dead_assignments_elimination.go), [`coalesce_variable_names.go`](internal/jscomp/coalesce_variable_names.go), [`flow_sensitive_inline_variables.go`](internal/jscomp/flow_sensitive_inline_variables.go), Stage-2 helpers): Embarrassingly parallel per-function dataflow/scope passes.
- **`-5.1%` `googleapis` / `-5.0%` 9-task (`1,210` LOC)** — **#5 `NodeTraversal` & `CombinedCompilerPass` Bitmasks** ([`node_traversal.go`](internal/jscomp/node_traversal.go), [`combined_compiler_pass.go`](internal/jscomp/combined_compiler_pass.go)): Token/leaf/JSDoc pre-filters skipping unneeded visitor callbacks across every stage (`0` historical bugs).
- **`-0.6%` `googleapis` / `0.0%` 9-task (`189` LOC)** — **#6 `FastGzip` Multi-Member Writer** ([`serialization_fast_gzip_output_stream.go`](internal/jscomp/serialization_fast_gzip_output_stream.go)): Parallel RFC 1952 concatenated gzip framing for `--filename_to_save_to` `.state` files (`0` bugs).
- **`-0.9%` `googleapis` / `-1.7%` 9-task (`510` LOC)** — **#7 `Peephole` Function-Wave Parallelism** ([`peephole_optimizations_pass.go`](internal/jscomp/peephole_optimizations_pass.go)): Single-file bottom-up AST folding across disjoint function depth waves.
- **`-1.3%` `googleapis` / `-2.6%` 9-task (`845` LOC)** — **#8 `ProcessDefines` Bloom + Single-Strong-Chunk Skip** ([`process_defines.go`](internal/jscomp/process_defines.go), [`js_chunk_graph.go`](internal/jscomp/js_chunk_graph.go)): Pre-filters candidate scripts and bypasses chunk-graph reachability on single-chunk builds (`0` bugs).
- **`-7.1%` `googleapis` / `-8.9%` 9-task (`4,930` LOC)** — **#9 Stage-1 Check/Color Parallelism + `RecordingCompiler`** ([`parallel_pass_helpers.go`](internal/jscomp/parallel_pass_helpers.go) + `32` passes): Largest single wall-time drop (`-17.6%` of `CHECKS`); includes the `1,295`-LOC `RecordingCompiler` harness shared by `#3`, `#4`, and `#7`.
- **`-1.5%` `googleapis` / `-3.4%` 9-task (`1,295` LOC)** — **#10 `InlineObjectLiterals` Fast Scanner** ([`inline_object_literals.go`](internal/jscomp/inline_object_literals.go)): Single-file Bloom + incremental scanner replacing full-AST `ReferenceCollector` (`-76.0%` on `collapseObjectLiterals`).
- **`-3.8%` `googleapis` / `0.0%` 9-task (`2,556` LOC)** — **#11 Parallel `TypedAst` & `RestoreState` Serde** ([`serialization_typed_ast_serializer.go`](internal/jscomp/serialization_typed_ast_serializer.go), [`serialization_typed_ast_deserializer.go`](internal/jscomp/serialization_typed_ast_deserializer.go), [`compiler_state.go`](internal/jscomp/compiler_state.go)): Wire-compatible parallel shard serialization/deserialization across multi-stage boundaries.
- **`-1.5%` `googleapis` / `-2.9%` 9-task (`1,665` LOC)** — **#12 Incremental `RemoveUnusedCode`** ([`remove_unused_code_incremental.go`](internal/jscomp/remove_unused_code_incremental.go), [`remove_unused_code_summary.go`](internal/jscomp/remove_unused_code_summary.go)): Per-function side-effect summaries & change stamps (`-32.5%` on `removeUnusedCode`).

### Tier 2 — Optional (`-3.1%` `googleapis` / `-6.6%` 9-task; Invasive or Workload-Specific)

- **`-1.3%` `googleapis` / `-1.9%` 9-task (`~1,850` LOC)** — **#14 UTF-8 `Scanner`, WTF-8 Checkpoints & `Parser` Slabs** ([`scanner.go`](internal/parser/scanner.go), [`source_file.go`](internal/jscomp/source_file.go), [`parser_slabs.go`](internal/parser/parser_slabs.go)): Cuts `parseInputs` by `-29.1%` (`-5.2%` total CPU), but rewrites core parser & UTF-16 offset internals.
- **`-0.8%` `googleapis` / `-0.7%` 9-task (`~1,150` LOC)** — **#15 Flat `Scope` / `TypedScope` Tables & Slabs** ([`scope.go`](internal/jscomp/scope.go), [`typed_scope.go`](internal/jscomp/typed_scope.go)): Adds `scopeVarIndex`, `varBloom`, and `lookupParent` skip pointers across core scope types.
- **`-1.0%` `googleapis` / `-2.9%` 9-task (`2,444` LOC)** — **#16 `PropUseIndex` Cross-Pass Property Index** ([`prop_use_index.go`](internal/jscomp/prop_use_index.go), [`prop_use_index_oc.go`](internal/jscomp/prop_use_index_oc.go)): Accelerates `OptimizeCalls` / `InlineProperties`, but couples invalidation hooks across 12 passes (`+11.9%` RSS on `typed-full`).
- **`0.0%` `googleapis` / `-1.0%` 9-task (`545` LOC)** — **#17 `CheckConformance` Static Token Bitmask** ([`check_conformance.go`](internal/jscomp/check_conformance.go), [`check_conformance_parallel.go`](internal/jscomp/check_conformance_parallel.go)): Useful only on conformance-heavy builds (`-8.6%` on `conformance-checks` alone).

### Tier 3 — Skip Permanently (`9,131` LOC for `+0.3%` net slowdown on `googleapis` & `15` bug-fix commits)

- **`-0.9%` `googleapis` / `-2.2%` 9-task (`1,419` LOC)** — **#13 `InlineFunctions` Dependency-Sharded Parallelism** ([`inline_functions_shard.go`](internal/jscomp/inline_functions_shard.go), [`inline_functions_parallel.go`](internal/jscomp/inline_functions_parallel.go)): Connected-component call-graph sharding with `3` race/ordering bug commits for `<1%` `googleapis` gain.
- **`-0.3%` `googleapis` / `-0.2%` 9-task (`~1,100` LOC)** — **#18 Embedded Flat `PropertyMap`, `HamtPMapArena` & `Color` Slabs** ([`jstype_property_map.go`](internal/rhino/jstype_property_map.go), [`hamt_p_map.go`](internal/rhino/hamt_p_map.go), [`color.go`](internal/colors/color.go)): Custom persistent HAMT + slab arena for `0.3%` wall gain.
- **`-0.4%` `googleapis` / `-0.7%` 9-task (`1,798` LOC)** — **#19 Parallel Multi-Chunk / `ToSource` & `SourceMap` Codegen** ([`compiler_codegen.go`](internal/jscomp/compiler_codegen.go), [`compiler_chunk_output.go`](internal/jscomp/compiler_chunk_output.go), [`source_map_generator_v3_parallel.go`](internal/sourcemap/source_map_generator_v3_parallel.go)): Serial Go codegen already emits `7,460` chunks in `~1.2s`; caused **`5` historical OOM/race/ordering bugs**.
- **`-0.1%` `googleapis` / `-0.4%` 9-task (`710` LOC)** — **#20 `CoalesceVariableNames` 128-Var Bitset** ([`small_bitset.go`](internal/jscomp/small_bitset.go), [`coalesce_variable_names.go`](internal/jscomp/coalesce_variable_names.go)): Redundant (`<0.04%`) once `#4` parallelizes `CVN` per function.
- **`+0.1%` slower `googleapis` / `-0.3%` 9-task (`1,365` LOC)** — **#21 `InlineFunctions` Subtree Change-Stamp Cache** ([`inline_functions_state.go`](internal/jscomp/inline_functions_state.go)): Bookkeeping overhead cancels skip savings at `-t4` (`2` stale-cache bugs).
- **`+0.1%` slower `googleapis` (`+3.4%` slower at `-t1`) / `-0.4%` 9-task (`1,545` LOC)** — **#22 Moving Pre-Order AST Compaction** ([`node_compact.go`](internal/rhino/node_compact.go), [`compiler_compact.go`](internal/jscomp/compiler_compact.go)): Zero `-t4` gain over Non-Moving `NodeStore`, `+6.3%` higher `OPT` RSS, and wipes `propIndex` + `inlineFnCache` (`3` dangling-`Node` bugs; keep only `13`-LOC non-moving `PruneDeadAux`).
- **`+2.0%` slower `googleapis` (`+3.0%` vs. `HEAD`) / `0.0%` 9-task (`1,194` LOC)** — **#23 Custom Binary Save-State** ([`compiler_state_binary.go`](internal/jscomp/compiler_state_binary.go)): Slower (`84.14s → 86.70s`) and bloats `checks.state` by **`+59.7%` (`69.57 MB → 111.13 MB`)**.

---

## 3. Benchmark Methodology & Ablation Architecture

Measuring each optimization solely in isolation against a cold single-threaded baseline would **double-count** overlapping optimizations (for example, a serial traversal filter and 4-worker pass parallelism speeding up the exact same pass). Conversely, measuring only knock-outs from `HEAD` can hide prerequisite couplings. We therefore combined **four uncontended sequential measurement passes**:

1. **Cumulative Bottom-Up `-t4` Ladder (`L0 → L1 → … → L10 → HEAD`)**:
   - Built `jsgo-mincore-switches` (`mincore` substrate revert + 14 startup-evaluated `JSGO_NO_*` kill-switches with zero `os.Getenv` hot-loop overhead, verified with PGO enabled and untouched upstream `PrebuildAst`/`PrebuildDependencyInfo`).
   - Enabled optimizations **one layer at a time** at `--num_parallel_threads=4`, recording wall time, CPU time, Peak RSS, `.state` sizes, and `JSGO_PASS_LOG` per-pass timings across every `googleapis-bench` stage (`CHECKS`, `OPTIMIZATIONS`, `FINALIZATIONS`) and all 9 standard tasks.
2. **Standalone Single-Cache `-t1` Runs (Phase 1)**:
   - Measured `JSGO_PROPINDEX=1`, `JSGO_RUC_INCR=1`, `JSGO_INLINE_FN_CACHE=1`, and `JSGO_AST_COMPACT=1` individually at `--num_parallel_threads=1` to compare isolated serial cache behavior against `-t4` stacked behavior.
3. **Single-Feature Knock-Out Binaries off `HEAD`**:
   - Built dedicated worktree binaries knocking out PGO (`-pgo=off`), `S2+S3` (`Scanner`/WTF-8/Parser slabs), `S4` (`Scope`/`TypedScope` flat tables), and `S5` (`PropertyMap`/`HamtPMap`/`Color` slabs) individually from `HEAD`, plus `JSGO_BINARY_SAVE_STATE=1`.
   - Verified additivity: the independent knock-out deltas of `S2+S3 (-1.68s) + S4 (-1.05s) + S5 (-0.40s) = -3.13s` matched the combined `L10 → HEAD` substrate step (`-3.21s`) to within **`97.5%`**.
4. **End-to-End Recommended Subset Verification (`R_recommended_subset`)**:
   - Measured Tier 1 enabled **together** in a single run (`L7 + JSGO_RUC_INCR=1`, disabling `>15,000` LOC of `#7` and `#13`–`#23`) to confirm zero hidden negative interactions when Tier 2 & Tier 3 are omitted.

---

## 4. Cross-Optimization Interactions (Synergies, Overlaps & Interference)

### 4.1 Positive Synergies (Optimizations That Amplify Each Other)

| Interacting Pair | Effect Alone | Effect Together `(measured)` | Mechanism Verified in `JSGO_PASS_LOG` |
|---|---:|---:|---|
| **`FastGzip` (`#6`) + Parallel `TypedAst` Serde (`#11`)** | `#6` alone: **`-0.6%`** (`-0.81s` `googleapis`) | **`-4.4%` `googleapis`** (`-5.73s`: `118.39s → 112.67s`) | Without `FastGzip`, stdlib single-threaded `gzip.Writer` serializes output from 4 parallel `TypedAst` shard encoders. Together, `CHECKS:serializeTypedAst` drops **`-28.6%`** (`8.43s → 6.02s`), `OPT:deserializeCompilerState` drops **`-48.4%`** (`4.77s → 2.46s`), and `OPT:serializeTypedAst` drops **`-48.6%`** (`1.75s → 0.90s`). |
| **`-t4` Worker Parallelism + `PropUseIndex` (`#16`)** | `#16` at `-t1`: **`-0.1%` `googleapis` `OPT`** (`51.61s → 51.49s`) | **`-1.0%` total / `-3.2%` `OPT` at `-t4`** (`37.85s → 36.63s` `OPT`; `-2.9%` 9-task) | Serially at `-t1`, walking `35.6M` AST nodes to build `propUseIndex` consumes almost all of the `-0.64s` saved in `optimizeCalls`. At `-t4`, `propUseIndex` builds in parallel across 4 workers, unlocking net wall-time gains across `optimizeCalls` (`-40.7%`), `inlineProperties` (`-83.9%`), and `DPAE` (`-43.5%`). |
| **Shared `RecordingCompiler` (`1,295` LOC in `#9`) + Pass Parallelism (`#3`, `#4`, `#7`, `#9`)** | `0.0%` *(infrastructure alone)* | **`-17.5%` `googleapis` (`-22.50s`) / `-20.4%` 9-task (`-5.85s`)** | `RecordingCompiler` + read-only `JSTypeRegistry` lock modes (`parallel_pass_helpers.go`) amortize across all 4 Tier-1 parallel pass groups (`6,565` LOC combined = `0.43s / 100 LOC`). |

### 4.2 Sub-Additive Overlaps (Pass Parallelism Shrinks Serial Pass Tweaks)

When `#4` and `#9` parallelize a pass by `~2.5–3×` across 4 workers, any serial optimization inside that same pass has proportionally less wall time left to save:

| Serial Optimization | Isolated / Pre-Parallelism Effect `(measured)` | Stacked *After* `-t4` Pass Parallelism `(measured)` | Practical Implication |
|---|---:|---:|---|
| **Incremental `RemoveUnusedCode` (`#12`)** | **`-1.9%` `googleapis`** (`-2.48s` at `-t1`; `-4.7%` 9-task) | **`-1.5%` `googleapis`** (`-1.92s` at `L8 → L8a`; `-2.9%` 9-task) | Retains ~`77%` of its wall-time win even at `-t4` because `RemoveUnusedCode`'s main reachability fixed-point remains partially serial (`6.19s → 4.18s`). **Keep (Tier 1).** |
| **`NodeTraversal` Bitmasks (`#5`)** | **`-5.1%` `googleapis`** (`-6.60s` at `L0 → L1`; `-5.0%` 9-task) | Est. **`~-2.0%` to `-2.5%`** wall (`-4.3%` CPU) after `#4`+`#9` | Even where parallelized, skipping callback/token dispatch cuts memory-bus pressure across all 4 workers and speeds up serial passes (`inferTypes`, `inlineAndCollapseProperties`, `closureCodeRemoval`). **Keep (Tier 1).** |
| **`CVN` 128-Var Bitset (`#20`)** | **`-0.12%` `googleapis`** (`-0.15s`: `0.63s → 0.48s` serial `CVN`) | **`<0.04%` `googleapis`** (`<0.05s` once `#4` drops `CVN` to `0.15s`) | Per-function `CVN` parallelism (`#4`) makes custom 128-var bitset coloring (`small_bitset.go`) redundant. **Skip (Tier 3).** |
| **`InlineFunctions` Subtree Cache (`#21`)** | **`-0.16%` `googleapis`** (`-0.22s` at `-t1`; `-1.9%` 9-task) | **`+0.08%` slower `googleapis`** (`+0.10s` at `L8b → L9`; `-0.3%` 9-task) | At `-t4`, subtree stamp maintenance overhead matches or exceeds skipped candidate scans. **Skip (Tier 3).** |

### 4.3 Destructive Interference (Optimizations That Actively Hurt Others)

- **Moving Pre-Order AST Compaction (`#22`) invalidates `PropUseIndex` (`#16`) and `InlineFunctions` Cache (`#21`)**:
  - Moving compaction allocates a fresh `NodeStore`, copies live AST nodes in pre-order, and remaps `Node` handles across 15 compiler structures. Because remapping every `[]Node` slice inside cross-pass caches was too brittle, [`compiler_compact.go:161-162`](internal/jscomp/compiler_compact.go) explicitly executes `c.inlineFnCache = nil` and `c.propIndex = nil` on every compaction.
  - Empirically, enabling moving compaction on top of all caches (`L9 → L10`) was **`+0.1%` slower (`85.32s → 85.41s`)** on `googleapis-bench` while raising `OPTIMIZATIONS` Peak RSS by **`+6.3%` (`5,452.8 MB → 5,798.5 MB`)** during dual-arena copying; at `-t1` it regressed `googleapis-bench` by **`+3.4%` (`133.98s → 138.57s`)**.
  - **Resolution**: Keep Non-Moving `NodeStore` + the **`13`-LOC** non-moving [`(*NodeStore).PruneDeadAux`](internal/rhino/node_compact.go) cleanup, and permanently skip moving compaction.

### 4.4 End-to-End Verification of the Recommended Tier 1 Bundle (`R_recommended_subset`)

To verify that skipping all `15,120` LOC of Tier 2 + Tier 3 causes no unexpected slowdown, we benchmarked **`R_recommended_subset`** (`L7 + JSGO_RUC_INCR=1`, i.e., Tier 1 minus `#7`, with `#13`–`#23` disabled) end-to-end:

| Metric | Post-Flip Baseline (`L0`) | Sum of Individual Ladder Steps (`L0` + `#1..#6, #8..#12`) | **Measured Together (`R_recommended_subset`)** | Adding Single-File `#7` (`Peephole` Wave) | Full `HEAD` (`All 23`) |
|---|---:|---:|---:|---:|---:|
| **`googleapis-bench` Total (`-t4`)** | `128.71s` (`0.0%`) | `88.92s` (**`-30.9%`**) | **`89.20s` (`-30.7%`)** | **`~87.99s` (`-31.6%`)** | `82.20s` (`-36.1%`) |
| **9-Task Suite Sum (`-t4`)** | `28.74s` (`0.0%`) | `18.80s` (**`-34.6%`**) | **`18.58s` (`-35.4%`)** | **`~18.08s` (`-37.1%`)** | `15.54s` (`-45.9%`) |
| **9-Task `typed-full` (`-t4`)** | `10.64s` (`0.0%`) | `6.56s` (**`-38.3%`**) | **`6.52s` (`-38.7%`)** | **`~6.34s` (`-40.4%`)** | `5.20s` (`-51.1%`) |

The end-to-end bundle matches the sum of individual steps within **`0.2%–0.8%`**, confirming that Tier 1 stands cleanly on its own without any Tier 2 or Tier 3 code.

---

## 5. Master Measurement & ROI Table (`#1` → `#23`)

*Both **% vs. Post-Flip (`L0`)** (fixed denominator: `128.71s` `googleapis`, `28.74s` 9-task) and **Step %** (relative to the immediately preceding rung) are shown below.*

| Rank | Optimization | Go-Only LOC | `googleapis` Δ vs `L0` *(Step %)* | `googleapis` Wall Δ | 9-Task Sum Δ vs `L0` *(Step %)* | `typed-full` Δ vs `L0` | Combined ROI (`s / 100 LOC`) | Bug Fixes | Verdict |
|---:|---|---:|---:|---:|---:|---:|---:|---:|---|
| **1** | Go PGO (`cmd/jsgo/default.pgo`) | `0` | **`-0.6%`** *(`-0.9%`; up to `-2.7%`)* | `-0.74s` *(to `-2.35s`)* | **`-0.3%`** *(`-0.6%`)* | `-0.5%` (`-0.06s`) | **`∞`** | `0` | **Tier 1** |
| **2** | Ryū `dtoa` (`internal/dtoa/d_to_a.go`) | `-728` | **`-0.04%`** *(`-0.1%`)* | `~-0.05s` | **`-0.1%`** *(`-0.1%`)* | `-0.1%` (`-0.01s`) | **`∞`** | `0` | **Tier 1** |
| **3** | Parallel `TypeCheck` | `418` | **`-3.5%`** *(`-4.0%`)* | `-4.46s` | **`-2.7%`** *(`-3.2%`)* | `-2.3%` (`-0.25s`) | **`1.26s`** | `1` | **Tier 1** |
| **4** | Local-Scope Mid-End Parallelism | `1,217` | **`-6.0%`** *(`-7.9%`)* | `-7.75s` | **`-7.0%`** *(`-9.3%`)* | `-8.8%` (`-0.93s`) | **`0.80s`** | `2` | **Tier 1** |
| **5** | `NodeTraversal` / `CombinedPass` Bitmasks | `1,210` | **`-5.1%`** *(`-5.1%`)* | `-6.60s` | **`-5.0%`** *(`-5.0%`)* | `-5.5%` (`-0.59s`) | **`0.66s`** | `0` | **Tier 1** |
| **6** | `FastGzip` Multi-Member Writer | `189` | **`-0.6%`** *(`-0.7%`)* | `-0.81s` | **`0.0%`** *(`0.0%`)* | `0.0%` (`0.00s`) | **`0.43s`** | `0` | **Tier 1** |
| **7** | `Peephole` Function-Wave Parallelism | `510` | **`-0.9%`** *(`-1.3%`)* | `-1.21s` | **`-1.7%`** *(`-2.5%`)* | `-1.7%` (`-0.18s`) | **`0.34s`** | `1` | **Tier 1** |
| **8** | `ProcessDefines` Bloom + Single-Chunk Skip | `845` | **`-1.3%`** *(`-1.4%`)* | `-1.65s` | **`-2.6%`** *(`-2.8%`)* | `-2.9%` (`-0.30s`) | **`0.28s`** | `0` | **Tier 1** |
| **9** | Stage-1 Check/Color Parallelism + `RecordingCompiler` | `4,930` | **`-7.1%`** *(`-8.4%`)* | `-9.08s` | **`-8.9%`** *(`-10.5%`)* | `-6.1%` (`-0.65s`) | **`0.24s`** | `3` | **Tier 1** |
| **10** | `InlineObjectLiterals` Fast Scanner | `1,295` | **`-1.5%`** *(`-1.5%`)* | `-1.88s` | **`-3.4%`** *(`-3.6%`)* | `-5.7%` (`-0.61s`) | **`0.22s`** | `1` | **Tier 1** |
| **11** | Parallel `TypedAst` & `RestoreState` Serde | `2,556` | **`-3.8%`** *(`-4.2%`)* | `-4.92s` | **`0.0%`** *(`0.0%`)* | `-0.1%` (`-0.01s`) | **`0.19s`** | `3` | **Tier 1** |
| **12** | Incremental `RemoveUnusedCode` | `1,665` | **`-1.5%`** *(`-2.2%`)* | `-1.92s` | **`-2.9%`** *(`-4.4%`)* | `-5.8%` (`-0.61s`) | **`0.16s`** | `2` | **Tier 1** |
| **13** | `InlineFunctions` Dependency Sharding | `1,419` | `-0.9%` *(`-1.4%`)* | `-1.22s` | `-2.2%` *(`-3.3%`)* | `-2.7%` (`-0.29s`) | `0.13s` | **`3`** | **Tier 3 — Skip** |
| **14** | Substrate S2+S3: UTF-8 `Scanner`, WTF-8 & Parser Slabs | `~1,850` | `-1.3%` *(`-2.0%`)* | `-1.68s` | `-1.9%` *(`-3.5%`)* | `-1.1%` (`-0.12s`) | `0.12s` | `2` | **Tier 2 — Opt** |
| **15** | Substrate S4: Flat `Scope`/`TypedScope` & Slabs | `~1,150` | `-0.8%` *(`-1.3%`)* | `-1.05s` | `-0.7%` *(`-1.3%`)* | `-0.7%` (`-0.08s`) | `0.11s` | `1` | **Tier 2 — Opt** |
| **16** | `PropUseIndex` Cross-Pass Property Index | `2,444` | `-1.0%` *(`-1.5%`)* | `-1.27s` | `-2.9%` *(`-4.8%`)* | `-4.4%` (`-0.47s`) | `0.09s` | `2` | **Tier 2 — Opt** |
| **17** | `CheckConformance` Static Token Bitmask | `545` | `0.0%` *(`0.0%`)* | `0.00s` | `-1.0%` *(`-1.2%`)* | `-0.1%` (`-0.01s`) | `0.06s` | `0` | **Tier 2 — Opt** |
| **18** | Substrate S5: Flat `PropertyMap`, `HamtPMap` & `Color` Slabs | `~1,100` | `-0.3%` *(`-0.5%`)* | `-0.40s` | `-0.2%` *(`-0.3%`)* | `-0.5%` (`-0.05s`) | `0.04s` | `1` | **Tier 3 — Skip** |
| **19** | Parallel Multi-Chunk / `ToSource` & `SourceMap` Codegen | `1,798` | `-0.4%` *(`-0.5%`)* | `-0.54s` | `-0.7%` *(`-0.8%`)* | `-0.7%` (`-0.08s`) | `0.04s` | **`5`** | **Tier 3 — Skip** |
| **20** | `CoalesceVariableNames` 128-Var Bitset | `710` | `-0.1%` *(`-0.1%`)* | `-0.15s` | `-0.4%` *(`-0.5%`)* | `-0.3%` (`-0.03s`) | `0.04s` | `0` | **Tier 3 — Skip** |
| **21** | `InlineFunctions` Subtree Change-Stamp Cache | `1,365` | **`+0.1%`** *(`+0.1%`)* | `+0.10s` | `-0.3%` *(`-0.5%`)* | `-0.6%` (`-0.06s`) | **`0.00s`** | `2` | **Tier 3 — Skip** |
| **22** | Moving Pre-Order AST Compaction | `1,545` | **`+0.1%`** *(`+0.1%`; `+3.4%` `-t1`)* | `+0.09s` | `-0.4%` *(`-0.7%`)* | `-1.0%` (`-0.10s`) | **`0.00s`** | **`3`** | **Tier 3 — Skip** |
| **23** | Custom Binary Save-State (`JSGO_BINARY_SAVE_STATE=1`) | `1,194` | **`+2.0%`** *(`+3.0%`)* | **`+2.56s`** | `0.0%` *(`0.0%`)* | `0.0%` (`0.00s`) | **`-0.21s`** | `1` | **Tier 3 — Skip** |

---

## 6. Complete Raw Benchmark Matrices & Per-Pass Attributions

### 6.1 `googleapis-bench` Complete Matrix (`43,917` Files, `35.6M` AST Nodes, `7,460` Chunks)

All runs verified `chunk_count = 7,460` and `chunks_sha256 = 7968e63b5cd96980ce477dea033cb4f3fbaef3268e622595dca3198a16fc1467`:

| Configuration | Threads | `CHECKS` Wall | `OPT` Wall | `FIN` Wall | **Total Wall** | **vs. `L0` (`%`)** | Total CPU | Peak RSS | `.state` (`checks` / `opt`) |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| **Stock Java `v20261104` (`-Xmx20g`)** | `4` | `97.55s` | `91.08s` | `20.06s` | **`208.69s`** | `+62.1%` | — | `12,610M` | — |
| **Go Heap `*Node` *(no `NodeStore`)* (`puremech`)** | `1` | `135.98s` | `74.23s` | `11.10s` | **`221.31s`** | `+71.9%` | `367.6s` | `11,112M` | — |
| **Pre-Flip `-t1` (`NodeStore` only, no compact)** | `1` | `69.84s` | `50.08s` | `9.17s` | **`129.10s`** | `+0.3%` | `139.8s` | `5,884M` | `68.1M` / `26.3M` |
| **`L0_postflip_true_t4` (Post-Flip `-t4` Baseline)** | `4` | `63.66s` | `55.36s` | `9.70s` | **`128.71s`** | **`0.0%`** | `154.2s` | `6,536M` | `68.1M` / `26.3M` |
| `L1_plus_traversal_filters` (`+#5`) | `4` | `60.25s` | `52.62s` | `9.26s` | **`122.12s`** | `-5.1%` | `148.6s` | `6,592M` | `68.1M` / `26.3M` |
| `L2_plus_serial_fastpaths` (`+#8, #10, #17, #20`) | `4` | `58.63s` | `50.61s` | `9.15s` | **`118.39s`** | `-8.0%` | `144.4s` | `6,284M` | `68.1M` / `26.3M` |
| `L2b_plus_fast_gzip` (`+#6`) | `4` | `58.34s` | `50.10s` | `9.15s` | **`117.59s`** | `-8.6%` | `142.6s` | `6,013M` | `68.1M` / `26.3M` |
| `L3_plus_state_io_par` (`+#11`) | `4` | `56.40s` | `47.55s` | `8.72s` | **`112.67s`** | `-12.5%` | `143.5s` | `6,057M` | `69.6M` / `26.5M` |
| `L4_plus_codegen_par` (`+#19`) | `4` | `56.41s` | `47.60s` | `8.11s` | **`112.13s`** | `-12.9%` | `143.6s` | `6,035M` | `69.6M` / `26.5M` |
| `L5_plus_typecheck_par` (`+#3`) | `4` | `52.07s` | `47.50s` | `8.09s` | **`107.66s`** | `-16.4%` | `146.1s` | `6,287M` | `69.6M` / `26.5M` |
| `L6_plus_stage1_par` (`+#9`) | `4` | `42.93s` | `47.52s` | `8.13s` | **`98.58s`** | `-23.4%` | `147.2s` | `6,077M` | `69.6M` / `26.5M` |
| `L7_plus_simple_midend_par` (`+#4`) | `4` | `42.02s` | `41.59s` | `7.23s` | **`90.84s`** | `-29.4%` | `150.0s` | `6,274M` | `69.6M` / `26.5M` |
| **`R_recommended_subset` (`L7 + #12`; skips `>15K` LOC)** | `4` | **`42.27s`** | **`39.70s`** | **`7.23s`** | **`89.20s`** | **`-30.7%`** | `148.0s` | `6,217M` | `69.6M` / `26.5M` |
| `L7b_plus_peephole_par` (`L7 + #7`) | `4` | `42.00s` | `40.78s` | `6.84s` | **`89.62s`** | `-30.4%` | `150.7s` | `6,102M` | `69.6M` / `26.5M` |
| `L8_plus_inline_fn_par` (`L7b + #13`) | `4` | `42.06s` | `39.69s` | `6.65s` | **`88.40s`** | `-31.3%` | `154.2s` | `5,920M` | `69.6M` / `26.5M` |
| `L8a_plus_ruc_incr` (`L8 + #12`) | `4` | `41.98s` | `37.85s` | `6.65s` | **`86.49s`** | `-32.8%` | `151.1s` | `6,177M` | `69.6M` / `26.5M` |
| `L8b_plus_propindex` (`L8a + #16`) | `4` | `41.98s` | `36.63s` | `6.61s` | **`85.22s`** | `-33.8%` | `150.4s` | `6,144M` | `69.6M` / `26.5M` |
| `L9_all_caches_nocompact` (`L8b + #21`) | `4` | `42.19s` | `36.48s` | `6.65s` | **`85.32s`** | `-33.7%` | `151.8s` | `6,285M` | `69.6M` / `26.5M` |
| `L10_mincore_full_compact` (`L9 + #22`) | `4` | `42.03s` | `36.71s` | `6.66s` | **`85.41s`** | `-33.6%` | `152.3s` | `5,950M` | `69.6M` / `26.5M` |
| **`v1104_full_anchor` (`HEAD` Full)** | `4` | **`39.67s`** | **`36.07s`** | **`6.46s`** | **`82.20s`** | **`-36.1%`** | `139.3s` | `5,707M` | `69.6M` / `26.5M` |
| *Knock-out*: `HEAD` without PGO (`#1`) | `4` | `39.94s` | `36.45s` | `6.56s` | `82.94s` | *(`+0.9%` vs `HEAD`)* | `141.3s` | `5,849M` | `69.6M` / `26.5M` |
| *Knock-out*: `HEAD` without `S2+S3` Parser (`#14`) | `4` | `41.23s` | `36.20s` | `6.45s` | `83.88s` | *(`+2.0%` vs `HEAD`)* | `147.3s` | `6,037M` | `69.6M` / `26.5M` |
| *Knock-out*: `HEAD` without `S4` Scope (`#15`) | `4` | `39.86s` | `36.77s` | `6.63s` | `83.25s` | *(`+1.3%` vs `HEAD`)* | `141.4s` | `5,870M` | `69.6M` / `26.5M` |
| *Knock-out*: `HEAD` without `S5` Types (`#18`) | `4` | `40.08s` | `36.08s` | `6.43s` | `82.60s` | *(`+0.5%` vs `HEAD`)* | `140.5s` | `6,035M` | `69.6M` / `26.5M` |
| *Phase 1*: `HEAD` + Binary Save-State (`#23`) *(vs `84.14s`)* | `4` | `43.50s` | `36.75s` | `6.45s` | `86.70s` | *(`+3.0%` vs `HEAD`)* | `155.1s` | `5,683M` | **`111.1M`** / `31.7M` |

### 6.2 Standard 9-Task Suite Complete Per-Task Matrix

All Go runs verified `100%` `MATCH` on JS output, sourcemaps, and diagnostics:

| Configuration | `typed-app` | `downlevel` | `checks-only` | `gen-ijs` | `inc-ijs` | `conformance` | `shaka` | `bundle` | `typed-full` | **9-Task Sum** | **vs. `L0` (`%`)** | Max RSS |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| **Cold Java CLI (`-t4`)** | `6.49s` | `4.62s` | `6.64s` | `3.04s` | `3.95s` | `9.12s` | `12.43s` | `12.84s` | `23.04s` | **`82.17s`** | `+185.9%` | `3,497M` |
| **Warm Java JVM (`-t4`)** | `2.00s` | `1.52s` | `2.73s` | `0.57s` | `1.26s` | `4.60s` | `6.80s` | `8.70s` | `15.79s` | **`43.97s`** | `+53.0%` | `3,584M` |
| `t1_base_nocompact` (`-t1`) | `1.40s` | `0.98s` | `2.49s` | `1.39s` | `1.04s` | `3.77s` | `4.32s` | `4.39s` | `10.28s` | **`30.05s`** | `+4.6%` | `542M` |
| `t1_only_propindex` (`-t1`) | `1.38s` | `0.98s` | `2.58s` | `1.43s` | `1.05s` | `3.70s` | `4.05s` | `4.37s` | `9.66s` | **`29.20s`** | `+1.6%` | `609M` |
| `t1_only_ruc_incr` (`-t1`) | `1.39s` | `0.96s` | `2.51s` | `1.37s` | `1.02s` | `3.69s` | `3.99s` | `4.24s` | `9.47s` | **`28.63s`** | `-0.4%` | `562M` |
| `t1_only_inline_fn_cache` (`-t1`) | `1.38s` | `0.99s` | `2.49s` | `1.39s` | `1.05s` | `3.73s` | `4.16s` | `4.24s` | `10.04s` | **`29.48s`** | `+2.6%` | `568M` |
| **`L0_postflip_true_t4` (`-t4` Base)** | **`1.35s`** | **`1.01s`** | **`2.06s`** | **`0.60s`** | **`0.79s`** | **`3.39s`** | **`4.44s`** | **`4.47s`** | **`10.64s`** | **`28.74s`** | **`0.0%`** | **`547M`** |
| `L1_plus_traversal_filters` (`+#5`) | `1.29s` | `0.97s` | `1.92s` | `0.60s` | `0.77s` | `3.20s` | `4.25s` | `4.28s` | `10.05s` | **`27.32s`** | `-5.0%` | `549M` |
| `L1b_plus_iol_fast` (`+#10`) | `1.27s` | `0.91s` | `1.92s` | `0.62s` | `0.77s` | `3.20s` | `3.97s` | `4.22s` | `9.44s` | **`26.33s`** | `-8.4%` | `535M` |
| `L1c_plus_defines_chunk_fast` (`+#8`) | `1.25s` | `0.89s` | `1.89s` | `0.60s` | `0.73s` | `3.12s` | `3.85s` | `4.13s` | `9.14s` | **`25.59s`** | `-11.0%` | `547M` |
| `L1d_plus_conformance_fast` (`+#17`) | `1.25s` | `0.89s` | `1.89s` | `0.60s` | `0.73s` | `2.83s` | `3.85s` | `4.13s` | `9.12s` | **`25.29s`** | `-12.0%` | `545M` |
| `L2_plus_dfa_fast` (`+#20`) | `1.25s` | `0.88s` | `1.89s` | `0.60s` | `0.74s` | `2.82s` | `3.82s` | `4.09s` | `9.09s` | **`25.17s`** | `-12.4%` | `551M` |
| `L3_plus_state_io_par` (`+#6, #11`) | `1.25s` | `0.88s` | `1.88s` | `0.60s` | `0.74s` | `2.84s` | `3.82s` | `4.08s` | `9.08s` | **`25.17s`** | `-12.4%` | `546M` |
| `L4_plus_codegen_par` (`+#19`) | `1.25s` | `0.88s` | `1.88s` | `0.55s` | `0.73s` | `2.84s` | `3.77s` | `4.08s` | `9.01s` | **`24.97s`** | `-13.1%` | `542M` |
| `L5_plus_typecheck_par` (`+#3`) | `1.20s` | `0.88s` | `1.72s` | `0.55s` | `0.73s` | `2.60s` | `3.69s` | `4.08s` | `8.76s` | **`24.19s`** | `-15.8%` | `531M` |
| `L6_plus_stage1_par` (`+#9`) | `1.08s` | `0.83s` | `1.24s` | `0.44s` | `0.64s` | `1.88s` | `3.45s` | `3.97s` | `8.11s` | **`21.63s`** | `-24.7%` | `532M` |
| `L7_plus_simple_midend_par` (`+#4`) | `0.97s` | `0.74s` | `1.24s` | `0.44s` | `0.64s` | `1.88s` | `2.99s` | `3.55s` | `7.17s` | **`19.63s`** | `-31.7%` | `555M` |
| **`R_recommended_subset` (`L7 + #12`)** | **`0.96s`** | **`0.72s`** | **`1.24s`** | **`0.44s`** | **`0.64s`** | **`1.89s`** | **`2.76s`** | **`3.41s`** | **`6.52s`** | **`18.58s`** | **`-35.4%`** | **`544M`** |
| `L7b_plus_peephole_par` (`L7 + #7`) | `0.95s` | `0.73s` | `1.23s` | `0.44s` | `0.64s` | `1.89s` | `2.85s` | `3.40s` | `7.00s` | **`19.12s`** | `-33.5%` | `550M` |
| `L8_plus_inline_fn_par` (`L7b + #13`) | `0.94s` | `0.72s` | `1.23s` | `0.44s` | `0.64s` | `1.89s` | `2.76s` | `3.16s` | `6.71s` | **`18.49s`** | `-35.7%` | `522M` |
| `L8a_plus_ruc_incr` (`L8 + #12`) | `0.99s` | `0.71s` | `1.24s` | `0.48s` | `0.70s` | `1.89s` | `2.55s` | `3.03s` | `6.10s` | **`17.67s`** | `-38.5%` | `557M` |
| `L8b_plus_propindex` (`L8a + #16`) | `0.89s` | `0.70s` | `1.24s` | `0.44s` | `0.64s` | `1.89s` | `2.37s` | `3.03s` | `5.63s` | **`16.82s`** | `-41.5%` | `623M` |
| `L9_all_caches_nocompact` (`L8b + #21`) | `0.89s` | `0.70s` | `1.28s` | `0.44s` | `0.65s` | `1.88s` | `2.36s` | `2.98s` | `5.56s` | **`16.74s`** | `-41.7%` | `644M` |
| `L10_mincore_full_compact` (`L9 + #22`) | `0.88s` | `0.70s` | `1.24s` | `0.44s` | `0.65s` | `1.89s` | `2.37s` | `3.02s` | `5.46s` | **`16.63s`** | `-42.1%` | `504M` |
| **`v1104_full_recheck` (`HEAD` Full)** | **`0.84s`** | **`0.63s`** | **`1.11s`** | **`0.39s`** | **`0.62s`** | **`1.71s`** | **`2.27s`** | **`2.78s`** | **`5.20s`** | **`15.54s`** | **`-45.9%`** | **`484M`** |
| *Knock-out*: `HEAD` without PGO (`#1`) | `0.83s` | `0.65s` | `1.13s` | `0.35s` | `0.58s` | `1.74s` | `2.28s` | `2.82s` | `5.26s` | `15.62s` | *(`+0.6%` vs `HEAD`)* | `463M` |
| *Knock-out*: `HEAD` without `S2+S3` (`#14`) | `0.86s` | `0.68s` | `1.18s` | `0.44s` | `0.61s` | `1.81s` | `2.29s` | `2.91s` | `5.32s` | `16.10s` | *(`+3.6%` vs `HEAD`)* | `466M` |
| *Knock-out*: `HEAD` without `S4` (`#15`) | `0.83s` | `0.65s` | `1.13s` | `0.35s` | `0.59s` | `1.73s` | `2.30s` | `2.90s` | `5.28s` | `15.74s` | *(`+1.3%` vs `HEAD`)* | `473M` |
| *Knock-out*: `HEAD` without `S5` (`#18`) | `0.83s` | `0.63s` | `1.14s` | `0.35s` | `0.60s` | `1.74s` | `2.27s` | `2.77s` | `5.25s` | `15.59s` | *(`+0.3%` vs `HEAD`)* | `516M` |

### 6.3 Top `JSGO_PASS_LOG` Pass/Tracer Drops Verified on `googleapis-bench`

- **`L0 → L1` (`NodeTraversal` Bitmasks `#5`, `-5.1%` total)**: `CHECKS:suspiciousCode` `-48.8%` (`1.60s → 0.82s`), `CHECKS:checkControlFlow` `-38.9%` (`1.85s → 1.13s`), `OPT:closureCodeRemoval` `-67.0%` (`1.03s → 0.34s`), `CHECKS:checkJsDocAndEs6Modules` `-41.5%` (`1.35s → 0.79s`), `OPT:inlineAndCollapseProperties` `-5.5%` (`6.72s → 6.35s`).
- **`L1 → L2` (Pass-Local Fast Paths `#8, #10, #20`, `-2.9%` total)**: `OPT:collapseObjectLiterals` **`-76.0%`** (`2.46s → 0.59s`), `CHECKS:processDefines_CHECK` **`-38.4%`** (`2.63s → 1.62s`), `OPT:processDefines_OPTIMIZE` **`-37.1%`** (`1.75s → 1.10s`), `FIN:coalesceVariableNames` **`-23.8%`** (`0.63s → 0.48s`).
- **`L2 → L3` (`FastGzip` `#6` + Parallel `TypedAst` Serde `#11`, `-4.4%` total)**: `CHECKS:serializeTypedAst` **`-28.6%`** (`8.43s → 6.02s`), `OPT:deserializeCompilerState` **`-48.4%`** (`4.77s → 2.46s`), `OPT:serializeTypedAst` **`-48.6%`** (`1.75s → 0.90s`), `FIN:deserializeCompilerState` **`-38.1%`** (`1.26s → 0.78s`).
- **`L4 → L5` (Parallel `TypeCheck` `#3`, `-3.5%` total)**: `CHECKS:checkTypes` **`-60.8%`** (`7.11s → 2.79s`).
- **`L5 → L6` (Stage-1 Check/Color Parallelism `#9`, `-7.1%` total)**: `checkVariableReferences` `-66.9%` (`1.39s → 0.46s`), `checkControlFlow` `-72.6%` (`1.13s → 0.31s`), `processDefines_CHECK` `-43.8%` (`1.62s → 0.91s`), `inferConsts` `-62.6%` (`1.07s → 0.40s`), `suspiciousCode` `-73.2%` (`0.82s → 0.22s`), `checkSideEffects` `-72.5%` (`0.80s → 0.22s`), `checkJsDocAndEs6Modules` `-70.9%` (`0.79s → 0.23s`), `checkStrictMode` `-73.8%` (`0.65s → 0.17s`), `checkVars` `-62.0%` (`0.71s → 0.27s`), `closureCheckModule` `-66.2%` (`0.65s → 0.22s`).
- **`L6 → L7` (Local-Scope Mid-End Parallelism `#4`, `-6.0%` total)**: `OPT:inlineVariables` **`-47.6%`** (`4.52s → 2.37s`), `OPT:earlyInlineVariables` **`-44.4%`** (`1.87s → 1.04s`), `OPT:deadAssignmentsElimination` **`-55.6%`** (`0.99s → 0.44s`), `OPT:inferConsts` **`-53.5%`** (`1.01s → 0.47s`), `FIN:coalesceVariableNames` **`-68.8%`** (`0.48s → 0.15s`), `OPT:normalize` **`-22.5%`** (`2.49s → 1.93s`).
- **`L7 → L7b` (`Peephole` Wave `#7`, `-0.9%` total)**: `OPT:peepholeOptimizations` `-56.3%` (`1.03s → 0.45s`), `FIN:peepholeOptimizations` `-56.4%` (`0.39s → 0.17s`), `FIN:latePeepholeOptimizations` `-59.4%` (`0.32s → 0.13s`).
- **`L7b → L8` (`InlineFunctions` Sharding `#13`, `-0.9%` total)**: `OPT:inlineFunctions` `-26.5%` (`4.49s → 3.30s`), `FIN:inlineFunctions` `-39.6%` (`0.48s → 0.29s`).
- **`L8 → L8a` (Incremental `RemoveUnusedCode` `#12`, `-1.5%` total)**: `OPT:removeUnusedCode` **`-32.5%`** (`6.19s → 4.18s`).
- **`L8a → L8b` (`PropUseIndex` `#16`, `-1.0%` total)**: `OPT:optimizeCalls` **`-40.7%`** (`1.62s → 0.96s`), `OPT:inlineProperties` **`-83.9%`** (`0.31s → 0.05s`), `OPT:deadPropertyAssignmentElimination` **`-43.5%`** (`0.23s → 0.13s`).
- **`L10 → HEAD` (Substrate `S2–S5`, `-2.5%` total)**: `CHECKS:parseInputs` **`-28.9%`** (`5.19s → 3.69s`, from `S2+S3`), `CHECKS:inferTypes` **`-3.3%`** (`20.96s → 20.27s`, from `S4+S5`).
