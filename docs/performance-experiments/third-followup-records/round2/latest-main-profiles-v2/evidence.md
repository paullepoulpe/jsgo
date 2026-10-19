# Fresh ff06 CPU and allocation evidence

Source `ff06ac15304f14c756bce9cd36ad2205212f51da`; binary `6d6540c1653875257b54ec802ac43028b97d9862d42b7e60a845a9e589b93877`. All five pinned Java artifact gates passed (exit, stdout, stderr, JS, source map). One untimed stock-main invocation per task, affinity 0,1, all experiments off, default runtime settings and shipped PGO. Only the main diagnostic adapter changes: it finalizes profiles immediately before the original successful exit receiver. The first failed capture is preserved separately.

CPU seconds below are sampled aggregate CPU, not elapsed runtime or Java speed ratios. Allocation bytes/objects are default-sampling cumulative estimates, not peak RSS or MemStats. No forced GC; the final uncompleted profiling cycle can be absent. No phase labels were injected. Cumulative stack entries and named families overlap and may be truncated.

| Task | Sampled CPU seconds | Sampled allocated MiB | Sampled allocated objects |
|---|---:|---:|---:|
| typed-app | 6.78 | 1,262.4 | 18,981,652 |
| bundle | 26.35 | 4,594.7 | 65,235,480 |
| downlevel-whitespace | 1.78 | 318.8 | 4,059,294 |
| downlevel-simple | 5.46 | 1,045.3 | 14,447,458 |
| checks-only | 10.67 | 1,535.4 | 18,522,070 |

## Priorities supported by current evidence

* Whole-tree traversal and repeated lexical scopes remain larger than small type-flow maps. Bundle: newNodeTraversal allocates 390.7 MiB directly (8.50%); SyntacticScopeCreator.CreateScopeOfScope cumulatively allocates 694.9 MiB (15.12%) and accounts for 2.45 sampled CPU seconds (9.30%). These figures overlap and must not be added.
* NodeTraversal now embeds 16 scopeStackEntry records (384 bytes) regardless of actual traversal depth. Source layout is 544 bytes before allocator rounding (expected class 576). Investigate a max-live-depth/overflow-growth histogram before choosing a smaller inline capacity; pointer scanning and retained backing arrays also matter. No representation change has been made.
* Bundle RemoveUnusedCode.Process: 2.33 CPU seconds (8.84%), 393.7 sampled MiB (8.57%). CoalesceVariableNames.Process: 0.98 CPU seconds (3.72%), 403.8 MiB (8.79%). Native liveness admission and fresh-main publication-layout compatibility still require their own activated experiment.
* BitSet is now only about 0.1% of direct sampled allocation in bundle. Existing embedded liveness sets and CloneInto invalidate old per-set backing-allocation estimates; old native publication code must not be copied unchanged.
* Checks: TypeInferencePass.InferAllScopes is 2.72 CPU seconds (25.49%) and 655.4 MiB (42.69%); TypedScopeCreator.CreateScope is 1.10 seconds (10.31%) and 248.7 MiB (16.20%). Typed inference and scopes remain substantial. This is cumulative ancestry attribution, not phase-label isolation.
* The clean full-parser ASCII prototype targets a real allocation source: parser ancestry is about 43.6% of checks sampled allocation, including SourcePosition slabs (8.34% flat), token slabs (4.06%), SourceRange slabs (3.98%) and StringToUTF16 (3.44%). Only actual gated measurement can establish its benefit.
* Qualified-name construction still makes many small objects: GetQualifiedName plus GetQualifiedNameForGetProp are 9.55% of checks sampled object count, 6.84% typed. Byte/CPU ceilings are lower than object-count percentages.

## Reproducibility

* overlay/manifest.json records original/replacement main, overlay, commit and shipped PGO hashes.
* build-provenance.json records the binary and action trace: only main was compiled, then one link; library package cache was reused without library edits.
* summary.json records artifact gates and profile hashes. run.py clears inherited experiments/GC/scheduler overrides before invoking exact baseline args.
* hotspots.json has flat/cumulative function and non-additive family summaries. Each task contains Go1.27.2 pprof CPU/alloc_space/alloc_objects flat and cumulative text views.
* Default allocs sampling is intentionally retained; do not compare these one-shot sampled totals to old post-GC or instrumented MemStats totals as if identical measurements.
