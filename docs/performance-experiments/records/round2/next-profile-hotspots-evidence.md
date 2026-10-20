# Labeled eecbcf44 diagnostic profiles and balanced PGO

Captured ordinary fresh CLI jobs with all candidate flags enabled, owned parser allocation enabled, GOGC200, CPU64-PC stacks and heap depth512. Byte gates passed for all five tasks and additional repetitions. These capture CPU samples are diagnostic, not warmed Java timing results or RSS measurements. Pass totals use exclusive `jsgo_pass` sample labels rather than retained stack-family unions.

## Base captures: exact CPU-label mass

| Task | Total sampled CPU | Labeled | Unlabeled | Largest pass labels |
|---|---:|---:|---:|---|
| typed-app | 6.68s | 75.0% | 1.67s | inferTypes11.4%, inlineVariables6.0%, removeUnusedCode5.8%, checkTypes5.1% |
| bundle | 19.33s | 84.1% | 3.08s | inlineFunctions10.9%, removeUnusedCode10.3%, inlineVariables7.6%, FSIV7.1% |
| downlevel-whitespace | 1.48s | 45.9% | 0.80s | normalize6.8%, rewriteCommonJsModules6.8% |
| downlevel-simple | 4.30s | 69.8% | 1.30s | removeUnusedCode6.0%, FSIV5.6% |
| checks-only | 7.93s | 76.0% | 1.90s | inferTypes28.0%, checkTypes9.0%, checkAccessControls7.8% |

Unlabeled CPU includes parsing/input handling/output and background runtime/GC work. It is not missing sample mass. Cumulative function ancestry within any label remains subject to64-PC truncation; flat innermost symbols are actual PCs. Labels survive lost owning-pass root frames and inherit into newly started goroutines.

## Larger measured targets

1. **SourceFile synchronization:** across eight whitespace repetitions (11.12s CPU), `jutil.GoroutineID` accounts for0.79s, of which0.78s lies under SourceFile.GetCode and0.76s has runtime.Stack frames. Owner breakdown: dependency scanning0.72s, AST parsing0.06s, other0.01s. GetCode total1.11s. This is an actual ~7% direct CPU target. Compiler-owned source synchronization can specialize audited CLI paths while preserving public/custom reentrant-monitor behavior. No extern/AST preparse cache is involved.
2. **Typed-scope construction:** checks inferTypes owns28% CPU, but TypedScopeCreator family alone retains1.05/7.93s (13.2% whole); native solver compute retains0.85s (10.7%). Typed-app scope family0.70/11.92s (5.9%, two repetitions). CreateScope is already memoized, so another scope-result cache is not the immediate answer. Investigate construction/annotation visits and binding/storage costs. LinkedFlowScope alone is a smaller fraction:0.15s checks,0.17s typed. Its post-GC sampled allocation is32MiB checks/9MiB typed; outer-HAMT leaf allocations5.5MiB/.5MiB and inner-HAMT6MiB/2MiB. Replacing tiny HAMT roots alone lacks evidence for a large win.
3. **Property lookup representation:** PropertyMap point lookups retain0.34s checks (4.3%) and0.25s typed (2.1%); full family0.45s checks (5.7%). The current backing TreeMap binary-searches sorted keys, invoking UTF16 StringCompareTo. That comparator already has an early ASCII-difference optimization. A private stock ASCII-key native index could preserve sorted iteration and non-ASCII UTF16 equality fallback while reducing repeated point-lookup comparisons. Keep insertion/JSDoc replacement, immutable map exceptions, custom backing-map behavior and noncanonical WTF8 equivalent-key semantics intact. This remains a bounded unimplemented opportunity, not a predicted win.
4. **Optimization discovery and storage:** bundle's largest labels are inlineFunctions, removeUnusedCode, inlineVariables and FSIV. Existing prototype work on discovery fusion and private collection storage addresses these measured groups. Static dispatch rewrites or emitter work without broader visit/allocation elimination have weaker evidence.

## Allocation publication warning

At GOGC200, the `.allocs.pprof` emitted before final forced GC can omit recent allocations. Checks total is1088.7MiB there and its LinkedFlowScope entries have zero allocation weights. The existing `.live.pprof` emitted after forced GC also contains cumulative `alloc_space`; it reports1771.6MiB and32MiB under LinkedFlowScope. Use its alloc_space column for these captures, and use MemStats deltas for exact total allocated bytes. Sampled allocation is not peak RSS. No new workload capture was required to expose this publication lag.

## Balanced training profile

`next-balanced-pgo/balanced.pgo`

SHA256 `bf41a4d79739dc07555dbeb4300672be5edb0cc5dcbf7365f37507d3f018e942`

All15 CPU profiles included: five base jobs plus typed1, whitespace7, simple2 extra repetitions. Repetitions merge within task first. Actual first sample type is verified as samples/count, then proportional largest-remainder integer rounding gives **exactly1,000,000 counts per task**,5,000,000 total. Every positive sample remains positive; all stack/location and label identities are preserved. This avoids unequal task weighting and avoids scalar-rounding drift.

Raw task counts after merging: typed1192, bundle1933, whitespace1112, simple1212, checks793. Each normalized CPU-nanosecond column totals10^13. The Go pprof merger accepted the normalized profiles. Four lightweight Python tool tests pass, covering exact mass, sample/label/non-value preservation, nonuniform weights, and rejection of malformed/wrong-type/downscaled/negative samples.

Training executable/source is eecbcf44, binarySHA and full environment/overlay/input/weight/tool hashes are in `next-balanced-pgo/manifest.json`. The source worktree subsequently advanced to9aee69c4. The helper required an explicit head-difference option and verified original overlay sources from pinned eec git blobs; it records current HEAD separately. Do not relabel this profile as training at a later commit. Parent-coordinated later changes were N2 parallel handling and tests; the capture core budget was1.

Executables: `balanced_pgo.py`, `pass_label_summary.py`, `label_hotspots.py`. Raw attribution: `next-pass-label-summary.json`, `next-label-hotspots.json`, `next-infer-scope-live-alloc-profile.json`, `sourcefile-monitor-stacks.json`.
