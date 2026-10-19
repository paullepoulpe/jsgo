# Inline BitSet and fresh-template results

Runtime GC tuning is excluded. All measured Go runs use default runtime settings; GOGC, GOMEMLIMIT, GOMAXPROCS and GODEBUG are unset. Fresh non-worker CLI only, no cached parsed extern ASTs. These candidates remain experimental; no broad warm-Java win or default promotion is established.

## Two inline BitSet words

Source a6afc95f versus prior b2e407f6; each retains the same nine earlier native factors and CLI_SOURCE_LOCK. Three fresh runs per cell, compiler jobs=1. The complete six-variant/five-workload screen had 90 matching outputs. The unconditional representation change is compared across separate binaries, including code-generation/layout differences.

| Workload | Prior Go s | Inline Go s | Prior peak MiB | Inline peak MiB |
|---|---:|---:|---:|---:|
| typed-app | 5.081 | 4.954 | 290.0 | 250.7 |
| bundle | 18.079 | 18.449 | 417.0 | 432.5 |
| downlevel-whitespace | 1.310 | 1.199 | 169.9 | 170.3 |
| downlevel-simple | 3.917 | 3.749 | 216.3 | 212.6 |
| checks-only | 7.396 | 7.336 | 668.7 | 664.4 |

Small-set allocation savings do not translate into a broad win: typed-app and simple transpilation improve in this screen, but bundle time and RSS regress. Empty sets grow from 24 to 32 bytes; larger sets add a 24-byte holder in addition to the larger header. A one-word 24-byte alternative needs an exact capacity/Clone distribution before implementation; sampled heap objects do not supply those counts.

## Fresh template vectors with corrected guards

Production source c7998d45, same binary with FRESH_TEMPLATE_VALUES off/on. Both variants have the same two-inline-word BitSet representation; native liveness is off. Five fresh runs per cell on the two workloads that execute inference. The untyped SIMPLE/whitespace tasks are inactive controls for this factor. All 20 measured compilations match pinned Java.

| Workload | Control Go s | Fresh Go s | Control peak MiB | Fresh peak MiB | Java/Go time (fresh) | Java/Go RSS (fresh) |
|---|---:|---:|---:|---:|---:|---:|
| typed-app | 5.123 | 4.934 | 275.7 | 252.6 | 0.758 | 4.32 |
| checks-only | 7.504 | 7.405 | 653.4 | 658.8 | 0.815 | 2.06 |

Typed-app time falls 3.7% and median RSS 8.4%; checks-only time falls 1.3% while RSS rises 0.8%. Sample ranges overlap; these are screening results, not independent-cohort confidence estimates. Java remains faster. Ratios are Java/Go, so values above 1 favor Go.

The initial vector experiment used a full per-scope admission check and regressed typed-app 15.8% and checks-only 22.6% against its own control. GetErrorCount entered/exited a reentrant monitor, adding two runtime.Stack walks per scope. The corrected version validates the private ownership/callback window at entry and after scope freezing, captures compiler/root/creator/registry/error-manager identities, and retains error-count checks through the same underlying mutex. It try-locks known wrappers and reads the concrete callback-free counter; on contention/reentry it releases all private locks before the original top-level getter. Public monitor methods remain unchanged. Internal callback-bearing probes retain full checks. Success/panic restoration and stop-after-error admission are tested.

The V1-to-V2 time recovery compares separate binaries/cohorts and does not isolate the guard change alone. The same-binary V2 comparison above establishes the observed net candidate effect. No semantic type/AST cache or specialization interning is used; every specialization and required resolver/replacer work stays fresh.

## Validation and measurement limits

* V2 compiler-focused race suite: 964 passed, 2 skipped, no race reports. Monitor selection: 5 normal and 5 race tests passed. Prior inline-only compiler race suite: 923 passed/2 skipped; initial combined constructor suite: 942 passed/2 skipped.
* Each candidate matched all five pinned Java artifacts/diagnostics/source maps at compiler job counts 1, 2, 4. These are correctness gates, not new thread-scaling performance measurements. The 2-CPU quota makes 4 jobs oversubscribed; the parameter bounds compiler work items rather than runtime goroutines or OS threads.
* An initial malformed class test fixture omitted a closing brace. The test-only fix 63dc7116 was followed by a complete clean race rerun. Original failed logs are retained; production behavior was not altered to weaken comparisons.
* Go 1.27.2, shipped PGO and identical -l=4 flags for jscomp/rhino/jutil/parsing. CPU affinity 0,1; no builds, tests or profiling overlapped timed measurements. Build-only memory controls were never benchmark runtime settings.
* Java ratios use the stored stock JVM cohort from 2026-10-18 22:08 UTC, 12 warmups plus 10 measured iterations, against Go measurements on2026-10-19. Each Java iteration constructs a fresh compiler and reparses externs. These are older-reference ratios, not fresh paired Java confidence estimates.
* Go RSS is median fresh-process lifetime peak. Java RSS is the warm JVM lifetime peak including warmup; the ratios describe these distinct windows, not equivalent per-compilation live memory.

Raw evidence is retained in the execution workspace under /workspace/jsgo-results/round2: inline-bitset-template-screen/, template-guard-v2-screen/, their command/provenance JSON, race logs and N1/N2/N4 gate logs. See fresh-template-values.md and inline-bitset.md for design details. Further design audits found shared string headers and concrete inference solver arrays already implemented; parser-cache locality and one-word BitSet occupancy require diagnostics before another representation experiment.
