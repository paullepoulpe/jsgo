Both BigQuery JFR recordings are retained as diagnostic evidence. Neither supplies accepted timing results or a measured speedup. The [selection plan](selection-plan.json), recorded before any v2 candidate performance result, governs screening and confirmation.

The [first recording](results/profile-stock-bigquery-v1/work/profile.jfr) started JFR after 46 warmup requests had passed the existing stability gate. JFR activation then triggered a large compilation wave: the eight capture requests took 16.314124460, 11.134029267, 6.252886348, 5.451278774, 5.154809828, 4.730449773, 4.595196748, and 4.940759519 seconds; their JIT counter deltas were 19,224, 11,058, 5,614, 5,390, 4,920, 4,631, 4,456, and 2,636 ms. The JIT counter aggregates compiler work and can exceed elapsed time. Its retained `summary.json` records successful capture and byte checks, but does not establish a settled profile. Treat [its aggregates](results/profile-stock-bigquery-v1/summary.tsv) as disturbed leads only. [Per-request evidence](results/profile-stock-bigquery-v1/profiled-requests.json) and [cleanup evidence](results/profile-stock-bigquery-v1/stopped-proof.json) remain intact.

The [corrected recording](results/profile-stock-bigquery-steady-v2/work/profile.jfr) kept an identical priming recording active before READY and throughout warmup and capture. A second recording began after 51 warmup requests passed the same gate. Both recordings used the built-in profile configuration, 10 ms Java execution sampling, allocation sampling throttled at 150/s, and 128-frame stacks. This avoided the large first-capture JIT wave, but the final strict stability check still failed. Its [failure](results/profile-stock-bigquery-steady-v2/failure.json) is retained; no thresholds were relaxed.

| Corrected capture request | Elapsed seconds, diagnostic only | JIT counter delta, ms |
|---|---:|---:|
| 1 | 4.254345773 | 40 |
| 2 | 4.329946654 | 108 |
| 3 | 4.277252552 | 58 |
| 4 | 4.431505542 | 33 |
| 5 | 4.419336180 | 81 |
| 6 | 4.256371119 | 28 |
| 7 | 4.548891248 | 115 |
| 8 | 4.388034690 | 146 |

The [final gate](results/profile-stock-bigquery-steady-v2/profile-stability.json) reported last-four versus prior-four median drift **2.3256306634620527%** (limit 5%), last-four JIT/wall **2.100764803429376%** (limit 2%), and last-four maximum/minimum **1.068725240544515** (limit 1.20). Only the JIT criterion failed. This gate evaluates the final windows, not every capture request individually; all eight rows above are disclosed. All six compiler-output digests matched the stock oracle on every warmup and capture request. Both BigQuery JVMs exited zero and their outer OOM deltas were zero. The separate bundle v1 attempt was interrupted during warmup and is not a completed profile.

The corrected [targeted aggregates](results/profile-stock-bigquery-steady-v2/targets.tsv) contain **2,221 ExecutionSample events** and **4,951 ObjectAllocationSample events**, with **16,087,851,504 sampled allocation-weight bytes**. CPU stacks: zero missing, three truncated. Allocation stacks: one missing, carrying 640 weight bytes, and 11 truncated. No `jdk.DataLoss` events were recorded. These counts describe recorded evidence; they do not establish exact CPU time or exact allocation totals.

| Target | CPU leaf / inclusive samples | Inclusive share of CPU samples | Inclusive allocation-weight bytes | Share of all allocation weight |
|---|---:|---:|---:|---:|
| `NodeTraversal.traverseBranch` | 586 / 1,612 | 72.5799% | 10,387,047,592 | 64.5645% |
| `SyntacticScopeCreator` family | 65 / 198 | 8.9149% | 1,120,212,152 | 6.9631% |
| `FunctionArgumentInjector` family | 2 / 24 | 1.0806% | 709,441,888 | 4.4098% |
| `FunctionArgumentInjector.findModifiedParameters` | 1 / 18 | 0.8104% | 603,604,432 | 3.7519% |
| `GlobalNamespace$Name.getFullName` | 0 / 0 | 0.0000% | 37,859,296 | 0.2353% |
| `PropertyMap.getAllKeys` | 0 / 1 | 0.0450% | 8,383,880 | 0.0521% |
| `PropertyMap.collectAllAncestors` | 0 / 1 | 0.0450% | 7,335,904 | 0.0456% |
| `PropertyMap.hashCode` | 0 / 0 | 0.0000% | 1,048,000 | 0.0065% |
| `$AutoValue_Color.hashCode` | 0 / 1 | 0.0450% | 0 | 0.0000% |
| `AutoValue_Color.hashCode` | 0 / 0 | 0.0000% | 0 | 0.0000% |
| `CodePrinter$MappedCodePrinter.getSourceMappings` | 0 / 0 | 0.0000% | 3,143,984 | 0.0195% |
| `SourceMapGeneratorV3.addSourcesContentMap` | 0 / 0 | 0.0000% | 0 | 0.0000% |
| `NodeTraversal.createsBlockScope` | 0 / 0 | 0.0000% | 0 | 0.0000% |

Inclusive rows overlap: recursive calls and each method are counted once per sampled stack, while family totals deduplicate the whole family per stack. Do not add wrapper, recursive-helper, and family percentages. The large `traverseBranch` inclusive share includes descendant work; its **26.3845% leaf share** also includes work beyond any proposed callback guard. Zero means no matching observed samples, not proof that a method did not execute. Baseline `AutoValue_Color` inherits `hashCode` from `$AutoValue_Color`, so the declaring owner matters. Stack truncation, sparse sampling, inlining, container pauses, and profiler overhead constrain attribution. Observed JFR timestamp spans are not compile durations.

The targeted file also retains the full `traverseBranch` leaf line/immediate-caller histogram and top 20 allocation stacks for each selected object class. `LinkedHashMap`, `LinkedHashMap$Entry`, and `HashMap$Node[]` account for 5.0552%, 7.3649%, and 4.0492% of allocation weight respectively; these classes serve many callers and cannot all be attributed to scopes. Exact `Scope` objects account for 1.4832%, and separate `TypedScope` objects for 0.0065%. The corrected broad `JfrSummary` output was not yet produced when these notes were prepared; the table above comes directly from `targets.tsv`.

Candidate priorities follow the recorded plan: first the `findModifiedParameters` accumulator, then bypasses of final, pure no-op traversal callbacks. The accumulator preserves preorder traversal and closure-flag isolation while removing intermediate sets and merges; independent source review found no introduced semantic defect. Its probe compares ordered output, repeated calls, AST preservation, and unsupported-formal exceptions against the pinned recursive algorithm. The separate empty-parameter shortcut changes malformed-body exception behavior and is not part of the accumulator-only candidate. PropertyMap identity membership, full-name caching, color-hash memoization, and source-map copying remain lower priorities because their observed target shares are small. The subsequent `JSType.differsFrom` identity shortcut is a source-supported additional lead; this target file does not measure its avoidable cost.

An illustrative Amdahl calculation gives scale, not a prediction: if the helper's 0.8104% inclusive CPU-sample share represented the affected elapsed fraction, halving the entire helper would reduce total time about 0.405%; removing it entirely would reduce time about 0.810% (roughly 1.0082× speedup). The accumulator removes only part of that work. Allocation weight is not time share, and changes in GC or cache behavior can affect work outside the sampled helper. Eighteen inclusive CPU samples and a diagnostically retained recording do not establish a hard performance bound.

The one-cohort CONTROL/candidate screen selects candidates; it carries no significance or confidence-interval claim. Confirmation compares against published STOCK in three fresh independent paired-JVM cohorts, with five alternating measured requests per variant after matched settled warmup, followed by separate bundle generalization. Performance runs use unchanged `CliFactors` without JFR. Every warm and measured compilation must pass the six-digest oracle, resource/pressure checks, and teardown review. Full upstream test-suite coverage is not implied by the focused probes and workload checks.

The primary confirmation estimate is the geometric mean of cohort ratios of within-variant medians. Persistent requests are not independent JVM replicates; pooled request statistics are descriptive. With three cohorts the percentile bootstrap is highly discrete and its endpoints essentially reduce to the observed range: all-minimum and all-maximum resamples each have probability 1/27, already above 2.5%. Bootstrap exclusion of one alone is weak evidence. The separate log-Student interval uses df=2 and assumes approximately normal independent log-cohort effects. Report all cohort effects, both uncertainty calculations, allocation separately from runtime, and any losing screens. A narrow interval, faster sample, or allocation reduction is not by itself a general speed claim.

Raw recordings, request rows, source variants, failures, and cleanup evidence remain retained. No timing inference from either profile replaces the independent benchmark protocol, and frozen `transfer-v1` inputs were not changed.
