# Code/design screen at default GC

All 100 fresh CLI runs (10 variants × five unchanged tasks × two runs) matched the pinned Java outputs, diagnostics and source maps. No worker or pre-parsed extern AST is used. Runtime GOGC/GOMEMLIMIT/GOMAXPROCS/GODEBUG were removed. CPU affinity is 0,1 under a two-CPU quota. Compiler jobs=1 for the factor screen.

Candidate CLI was frozen at `0fafa262`; later `690e1dea` corrects a test fixture only. Go1.27.2, shipped PGO, `-l=4` for jscomp/rhino/jutil/parsing. Previous CLI source `5906ebbb`. Hashes/environments/order/machine/raw samples are in `code-design-screen/` and `build-traversal-provenance.json`.

No broad speed improvement is established. The complete compound and all-code configurations fail to beat the previous implementation across the suite, and still trail the stored warm Java times. No new default is recommended. Several apparent gains against the current same-binary control disappear against the previous binary; two-run noise and code-generation/layout changes matter.

The all-code configuration includes old pass-elision/reference batching/coalesce switches in addition to the six new factors. Compound includes only the new factors. Every environment is recorded verbatim; no minimum-per-task configuration is presented as a single result.

Ratios are Java / Go; values above one favor Go. Go RSS is median fresh CLI peak; Java RSS is the stored JVM-lifetime peak including warmup. The factor screen uses the older stored Java cohort, so these are reference ratios, not fresh paired confidence estimates.

## typed-app

| Variant | Time s | Peak MiB | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|
| prior | 5.116 | 271.3 | 0.731 | 4.02 |
| base | 4.710 | 266.9 | 0.794 | 4.08 |
| mixed_inlining | 4.834 | 280.5 | 0.773 | 3.89 |
| access_facts | 4.941 | 249.7 | 0.757 | 4.37 |
| tasks | 5.054 | 270.9 | 0.740 | 4.02 |
| intrusive | 4.780 | 251.9 | 0.782 | 4.33 |
| leaf_visits | 5.070 | 264.8 | 0.738 | 4.12 |
| scope_insertion | 4.978 | 267.0 | 0.751 | 4.08 |
| compound | 4.829 | 263.8 | 0.774 | 4.13 |
| all_code | 4.995 | 252.7 | 0.749 | 4.31 |

## bundle

| Variant | Time s | Peak MiB | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|
| prior | 18.303 | 435.9 | 0.907 | 3.89 |
| base | 19.162 | 423.6 | 0.866 | 4.01 |
| mixed_inlining | 18.524 | 425.4 | 0.896 | 3.99 |
| access_facts | 18.179 | 442.2 | 0.913 | 3.84 |
| tasks | 19.394 | 439.6 | 0.856 | 3.86 |
| intrusive | 18.796 | 431.2 | 0.883 | 3.93 |
| leaf_visits | 18.430 | 434.8 | 0.901 | 3.90 |
| scope_insertion | 18.691 | 433.7 | 0.888 | 3.91 |
| compound | 18.819 | 436.4 | 0.882 | 3.89 |
| all_code | 18.960 | 436.8 | 0.876 | 3.88 |

## downlevel-whitespace

| Variant | Time s | Peak MiB | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|
| prior | 1.180 | 174.5 | 0.795 | 5.33 |
| base | 1.162 | 168.4 | 0.807 | 5.52 |
| mixed_inlining | 1.258 | 166.3 | 0.746 | 5.59 |
| access_facts | 1.158 | 166.6 | 0.810 | 5.58 |
| tasks | 1.173 | 165.2 | 0.800 | 5.63 |
| intrusive | 1.188 | 170.2 | 0.790 | 5.46 |
| leaf_visits | 1.170 | 165.2 | 0.802 | 5.63 |
| scope_insertion | 1.201 | 162.6 | 0.781 | 5.72 |
| compound | 1.176 | 164.6 | 0.798 | 5.65 |
| all_code | 1.196 | 171.8 | 0.785 | 5.41 |

## downlevel-simple

| Variant | Time s | Peak MiB | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|
| prior | 3.674 | 213.0 | 0.763 | 5.25 |
| base | 3.841 | 208.7 | 0.730 | 5.36 |
| mixed_inlining | 3.751 | 213.5 | 0.748 | 5.24 |
| access_facts | 3.747 | 210.0 | 0.749 | 5.33 |
| tasks | 3.723 | 208.4 | 0.754 | 5.37 |
| intrusive | 3.675 | 200.5 | 0.763 | 5.58 |
| leaf_visits | 3.710 | 207.9 | 0.756 | 5.38 |
| scope_insertion | 3.565 | 211.6 | 0.787 | 5.29 |
| compound | 3.680 | 211.9 | 0.762 | 5.28 |
| all_code | 3.732 | 208.5 | 0.752 | 5.36 |

## checks-only

| Variant | Time s | Peak MiB | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|
| prior | 7.369 | 665.1 | 0.819 | 2.04 |
| base | 7.454 | 672.1 | 0.810 | 2.02 |
| mixed_inlining | 7.815 | 684.6 | 0.773 | 1.99 |
| access_facts | 7.550 | 659.2 | 0.800 | 2.06 |
| tasks | 7.661 | 667.0 | 0.788 | 2.04 |
| intrusive | 7.567 | 659.4 | 0.798 | 2.06 |
| leaf_visits | 7.395 | 663.2 | 0.816 | 2.05 |
| scope_insertion | 7.570 | 646.4 | 0.798 | 2.10 |
| compound | 7.904 | 663.9 | 0.764 | 2.05 |
| all_code | 7.417 | 683.0 | 0.814 | 1.99 |

## Fresh whitespace thread sweep

Same frozen candidate, all-code environment. Five Go CLI processes per count; pinned stock Java, twelve warmups and ten measurements per fresh JVM. Each Java iteration creates a fresh runner/compiler and reparses externs. No worker or cached AST setup. Both languages receive the same compiler job count. All artifacts match.

| Jobs | Go s | Java s | Go MiB | Java MiB | Java/Go time | Java/Go RSS |
|---:|---:|---:|---:|---:|---:|---:|
| 1 | 1.212 | 0.896 | 164.1 | 990.0 | 0.740 | 6.03 |
| 2 | 1.133 | 0.882 | 164.8 | 917.8 | 0.778 | 5.57 |
| 4 | 1.174 | 0.848 | 166.8 | 869.3 | 0.722 | 5.21 |

Two jobs produced the best observed Go median on this two-CPU machine. Four jobs measures oversubscription, not four physical-core scaling. Java remains faster. This is one JVM cohort per count, not an independently replicated confidence claim. Measured Java half-to-half drift percentages: N1=-3.36%, N2=-0.02%, N4=5.06%.

The job parameter bounds independent compiler work; it does not set GOMAXPROCS or cap runtime/GC goroutines or OS threads. Inherited GOMAXPROCS was removed and normal runtime CPU detection remained active. Raw fresh cohorts are in `default-gc-white-traversal-thread-sweep/`.

## Correctness and diagnostics

Corrected focused normal and race suites each have 760 passing checks and two skips, with no races. The broader pre-leaf/scope-insertion run has 6,270 passes, 51 skips and no failures. Latest CLI checks match all five tasks at N1/N2/N4. These are the reported suites, not full Closure parity.

Fresh builtin extern inventory passes exact AST/JSDoc/source/property/feature/diagnostic/exception comparisons with ordinary and private parser allocators. Native parsing accepts127/128 BROWSER entries, 96.12% of bytes, with zero original-parser fragments in accepted files. Sole whole-file rejection is es3.js at var arguments. Follow-up statement-fragment parsing and declaration recipes are separate, unmeasured prototypes.

A source audit of the compact Node64 layout confirms all119 simple Is* predicates already compare packed token constants without a cold GetToken fallback. That suspected missing optimization was already present; no new production change was made.

## Next work

Heap allocation attribution uses the post-forced-GC live profile alloc_space/alloc_objects columns. Exact MemStats deltas are distinct from sampled allocation shares and from peak RSS. Contextual renaming uses about314MiB of sampled bundle allocations (6.5%); pass-local compact shadow/member storage is being explored. Nothing in this report attributes gains to GC tuning.
