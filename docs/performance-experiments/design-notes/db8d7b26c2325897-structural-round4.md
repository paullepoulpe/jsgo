# Structural CLI screening with default GC

2026-10-19. Runtime GC uses its defaults, following the user's request to stop GC tuning. Inherited GOGC, GOMEMLIMIT, GODEBUG and GOMAXPROCS settings are cleared. Build-tool memory settings are separate. Each measurement runs a fresh CLI process, with no worker or parsed-extern cache. Two interleaved runs per variant and task use CPU affinity 0,1 and compiler thread count N=1. Builds, tests, profiling and parallel analysis stopped during measurement. All 70 outputs, diagnostics and exits match stock Java.

Ratios are Java/Go; above 1 favors Go. Go RSS is the median fresh-process peak; Java RSS is the stored JVM lifetime peak including warmups. Java's cohort used 12 warmups and 10 measured iterations. It is a historical reference, not a fresh paired confidence estimate. Two Go runs provide screening evidence only.

`control` is the preceding candidate with 80-byte nodes. `struct_base` is the new binary with all five new switches off. `source` enables only the private CLI SourceFile lock. `compound` enables that lock, sparse inlining discovery, property indexing, fused initial type analyses and the dependency scanner. These four use the shipped PGO profile and package-scoped `-l=4`. The separate `pgo80`, `compact` and `compact_header` variants share a balanced PGO profile and runtime flags, isolating compact nodes and their corrected metadata accessor. They exclude the old tiny-flow-map and qualified-name experiments.

## typed-app

Java 3.739s / 1090.2MiB.

| Variant | Go seconds | Go peak MiB | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|
| control | 5.184 | 271.1 | 0.721 | 4.021 |
| struct_base | 4.831 | 249.9 | 0.774 | 4.363 |
| source | 4.934 | 248.2 | 0.758 | 4.393 |
| compound | 5.115 | 280.8 | 0.731 | 3.882 |
| pgo80 | 5.312 | 272.1 | 0.704 | 4.006 |
| compact | 5.346 | 254.7 | 0.700 | 4.280 |
| compact_header | 5.476 | 257.7 | 0.683 | 4.231 |

## bundle

Java 16.600s / 1696.6MiB.

| Variant | Go seconds | Go peak MiB | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|
| control | 18.093 | 416.7 | 0.918 | 4.072 |
| struct_base | 18.469 | 419.0 | 0.899 | 4.049 |
| source | 18.265 | 429.8 | 0.909 | 3.947 |
| compound | 18.305 | 432.2 | 0.907 | 3.925 |
| pgo80 | 18.573 | 424.1 | 0.894 | 4.000 |
| compact | 18.791 | 392.4 | 0.883 | 4.324 |
| compact_header | 18.716 | 407.8 | 0.887 | 4.160 |

## downlevel-whitespace

Java 0.938s / 930.2MiB.

| Variant | Go seconds | Go peak MiB | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|
| control | 1.380 | 160.9 | 0.680 | 5.783 |
| struct_base | 1.354 | 171.3 | 0.693 | 5.430 |
| source | 1.356 | 166.0 | 0.692 | 5.602 |
| compound | 1.252 | 160.3 | 0.749 | 5.804 |
| pgo80 | 1.356 | 165.2 | 0.692 | 5.629 |
| compact | 1.405 | 165.7 | 0.668 | 5.614 |
| compact_header | 1.344 | 165.1 | 0.698 | 5.633 |

## downlevel-simple

Java 2.805s / 1118.4MiB.

| Variant | Go seconds | Go peak MiB | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|
| control | 3.934 | 202.1 | 0.713 | 5.533 |
| struct_base | 3.856 | 209.4 | 0.728 | 5.342 |
| source | 3.723 | 198.1 | 0.754 | 5.647 |
| compound | 4.059 | 201.6 | 0.691 | 5.547 |
| pgo80 | 3.824 | 204.5 | 0.734 | 5.469 |
| compact | 4.088 | 202.2 | 0.686 | 5.531 |
| compact_header | 4.191 | 205.6 | 0.669 | 5.439 |

## checks-only

Java 6.038s / 1360.0MiB.

| Variant | Go seconds | Go peak MiB | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|
| control | 7.502 | 673.5 | 0.805 | 2.019 |
| struct_base | 7.528 | 661.2 | 0.802 | 2.057 |
| source | 7.380 | 650.8 | 0.818 | 2.090 |
| compound | 7.328 | 683.1 | 0.824 | 1.991 |
| pgo80 | 7.532 | 653.3 | 0.802 | 2.082 |
| compact | 7.825 | 650.0 | 0.772 | 2.092 |
| compact_header | 7.859 | 647.5 | 0.768 | 2.101 |

## Findings and limits

None beats Java. Compound whitespace improves 9.2% against the control at essentially the same RSS. Checks improve 2.3% with 1.4% higher RSS; SIMPLE regresses 3.2%. Other differences are small. Source locking alone is neutral on whitespace, so the predicted lock benefit remains unproven and actual activation will be checked. The new binary with its switches off is itself 6.8% faster on typed-app than the older control. Compiler-generated code and allocation effects need to be separated from claims about individual features. The full compound is not approved as a default.

The compact-node metadata accessor shrank from 89 to 22 machine-code bytes; ordinary 80-byte nodes use 3 bytes. Payload, header and race tests pass. This fixes a real inefficiency, but does not produce a broad compiler speedup. Against the original compact variant, it is 4.3% faster on whitespace and slower on several others, with small RSS differences.

Before timing, 113 new behavior cases and 4,743 broader regression cases passed. Nineteen cases were skipped, including existing disabled cases and optional external coverage oracles. Fifteen PropertyMap cases and the focused structural suite passed under race detection. All five complete CLI tasks match Java at N=1 and N=2. A followup guard for wrapped inference compilers is queued on a separate branch; stock CLI measurements are unaffected. Public constructors retain their original paths. The new prototypes remain opt-in and unshipped.

Next, isolate feature effects and validate compact event containers, the final inlining worklist, recorded inference schedules and compact typed declarations. Thread scaling needs fresh Java cohorts at matching thread counts. No GC tuning variants are planned.

Raw evidence in the results directory: `structural-screen/{machine,summary}.json`, `structural-screen-command.json`, `build-structural-provenance.json`, `build-node64-header-provenance.json`, `node64-token-code-size.json`, `tests-structural-summary.json`, `tests-structural-focused-race.log`, and `gates-structural-l4.log`.
