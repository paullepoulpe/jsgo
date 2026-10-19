# CLI toolchain and compact-node screening

Completed 2026-10-19, one-shot Go CLI versus stock warmed Java. No workers or pre-parsed extern ASTs. All six variants use the same five original tasks, CPU affinity 0,1, compiler N=1, GOGC=200, and private per-parser allocation blocks. Runtime GOMAXPROCS is unset. Two fresh processes per variant/task in interleaved order; no builds, tests, profiling, or parallel analysis during measurement. All 60 measured full CLI outputs/diagnostics/exits match the Java artifact reference.

Java ratios use the stored cohort from the original baseline (12 warmups and 10 measurements), not a fresh paired cohort or confidence interval. Go memory is the median fresh-process peak RSS; Java memory is the stored JVM lifetime peak including warmup. MiB throughout. Ratios are Java/Go; above 1 favors Go.

Variants: prior200 is the previous 4f417579 candidate with package-scoped -l=4; next200 adds module path helpers, empty-pass proofs, reference batches, and parallel Coalesce (the last is inactive at N=1); next_noproofs200 disables the four new opt-in switches while retaining automatic module helpers; pgo200 uses balanced fresh PGO with -l=4; newinline200 uses the same balanced profile and GOEXPERIMENT=newinliner instead of -l=4, including matching rebuilt standard library; node64_200 isolates compact Node storage plus current private parser allocation blocks against pgo200, excluding the old tiny-flow-map and qualified-name experiments.

Balanced training merged 15 byte-gated captures over all five tasks, normalizing the actual first samples/count column to exactly one million counts per task. The diagnostic overlay adds pass labels without enabling CompilerTracer. It trains on eecbcf44; the later coalesce wrapper splitting fix only changes N>=2, while all training/timing here uses N=1. Raw source/binary/profile/overlay hashes are retained in build manifests.

## typed-app

Java: 3.739s, 1090.2 MiB.

| Variant | Go seconds | Go RSS MiB | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|
| prior200 | 4.575 | 368.8 | 0.817 | 2.956 |
| next200 | 4.657 | 340.6 | 0.803 | 3.201 |
| next_noproofs200 | 4.666 | 339.4 | 0.801 | 3.213 |
| pgo200 | 4.693 | 348.4 | 0.797 | 3.130 |
| newinline200 | 4.828 | 339.9 | 0.775 | 3.207 |
| node64_200 | 4.868 | 325.3 | 0.768 | 3.352 |

## bundle

Java: 16.600s, 1696.6 MiB.

| Variant | Go seconds | Go RSS MiB | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|
| prior200 | 17.184 | 640.6 | 0.966 | 2.648 |
| next200 | 17.515 | 617.4 | 0.948 | 2.748 |
| next_noproofs200 | 17.277 | 665.9 | 0.961 | 2.548 |
| pgo200 | 17.432 | 619.8 | 0.952 | 2.737 |
| newinline200 | 17.518 | 632.6 | 0.948 | 2.682 |
| node64_200 | 17.751 | 610.5 | 0.935 | 2.779 |

## downlevel-whitespace

Java: 0.938s, 930.2 MiB.

| Variant | Go seconds | Go RSS MiB | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|
| prior200 | 1.370 | 193.5 | 0.685 | 4.808 |
| next200 | 1.234 | 193.8 | 0.761 | 4.799 |
| next_noproofs200 | 1.273 | 190.5 | 0.737 | 4.884 |
| pgo200 | 1.351 | 203.6 | 0.694 | 4.570 |
| newinline200 | 1.248 | 195.2 | 0.752 | 4.765 |
| node64_200 | 1.334 | 193.7 | 0.703 | 4.803 |

## downlevel-simple

Java: 2.805s, 1118.4 MiB.

| Variant | Go seconds | Go RSS MiB | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|
| prior200 | 3.620 | 271.1 | 0.775 | 4.125 |
| next200 | 3.595 | 271.9 | 0.780 | 4.114 |
| next_noproofs200 | 3.498 | 269.5 | 0.802 | 4.149 |
| pgo200 | 3.526 | 273.4 | 0.796 | 4.091 |
| newinline200 | 3.560 | 256.8 | 0.788 | 4.355 |
| node64_200 | 3.679 | 256.3 | 0.763 | 4.363 |

## checks-only

Java: 6.038s, 1360.0 MiB.

| Variant | Go seconds | Go RSS MiB | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|
| prior200 | 7.005 | 860.0 | 0.862 | 1.581 |
| next200 | 7.099 | 921.7 | 0.851 | 1.476 |
| next_noproofs200 | 7.086 | 857.4 | 0.852 | 1.586 |
| pgo200 | 7.373 | 883.3 | 0.819 | 1.540 |
| newinline200 | 7.686 | 972.8 | 0.786 | 1.398 |
| node64_200 | 7.312 | 814.9 | 0.826 | 1.669 |

## Interpretation

None beats Java. Next200 whitespace is 9.9% faster than the interleaved prior200 control, at similar RSS; other next200 tasks vary from slightly faster to 1.8% slower. Two runs are a screen, not evidence of a stable universal gain. Disabling the new optional switches is faster on SIMPLE and bundle in this sample; they are not approved as defaults.

Fresh balanced PGO and the alternate inliner show no broad gain; neither is retained as a performance improvement. Isolated Node64 reduces typed/checks RSS but again loses speed across this screen. This closes the previous confound from combining compact storage with tiny-flow/QName changes, and confirms that current private parser allocation blocks do not automatically rescue the compact layout's speed. It remains an experimental memory tradeoff.

Correctness completed before timing: all new focused normal and race tests pass, including real N2/N4 worker overlap and exact FSIV→Coalesce intermediate phase snapshots; reference allocator GC ownership; full compact-node parsing suite; and all five full CLI byte gates at N1/N2 for fresh-PGO, alternate-inliner, and compact-node builds. Optional external parser-oracle tests requiring separate environment input skipped; do not claim that oracle was run. Experimental compound still needs broad regression validation before shipping.

The next independent structural batch targets initial source-load monitor costs, redundant inlining/typed-analysis walks, property lookup, and dependency scanner adapters. No performance claim is made for those untested prototypes.

Evidence: next-screen/{machine,summary}.json, next-screen-command.json, build-next-{l4,balanced,node64}-provenance.json, next-balanced-pgo/manifest.json, next-pass-label-summary.json and corresponding gates/tests logs.
