# Native liveness owned-word publication screen

All 60 fresh CLI runs (four variants × five unchanged tasks × three runs) matched pinned Java output, diagnostics, exit status and source maps. Each task/variant gate reports MATCH; all 60 raw process exits are zero. No worker, parsed extern cache or GC tuning. CPU affinity 0,1 under the two-CPU quota; compiler jobs=1; runtime GOGC/GOMEMLIMIT/GODEBUG/GOMAXPROCS unset. Timing ran without concurrent builds/tests/profiling.

Candidate production source **92fa5325**, documentation-only **f2c17056**; prior production source **b2e407f6**. Go 1.27.2, shipped PGO and the same four package `-l=4` flags. Raw environments, binary hashes, execution order, per-run bytes/digests and measurements: [native-live-words-screen/](native-live-words-screen/). This is a small randomized screen, with three observations per cell; it provides no independent-cohort confidence interval and changes no default.

V2 replaces the native liveness solver’s separate tiny published BitSet word allocations with one invocation-owned array, partitioned into exact-capacity windows. It preserves the original BitSet representation and API. Each window is capped before trimming, so later growth cannot overwrite a neighboring lattice. The native solver, lazy escape-sensitive gen/kill, fixed-point order and standard annotation publication remain the same. The separate generic inline-BitSet representation is not included.

Variants: `prior_base` and `prior_live` use the prior binary with native liveness off/on; `base` and `words` use the V2 binary with native liveness off/on. Every variant retains the same nine preceding native/inline facts and CLI_SOURCE_LOCK ownership witness. Other round-8 compact/scope/output factors and scope recipes remain off. Same-binary off/on comparisons isolate enabling the complete native solver; `words` versus `prior_live` compares the two enabled implementations across binaries, including code-generation/layout differences.

No broad improvement or warm-Java win is established, and V2 is not promoted. The clearest followup is simple transpilation: `words` is 4.6% faster than its same-binary base and 4.8% faster than the prior enabled solver, with essentially unchanged median RSS. Typed compilation is 1.3% slower than its same-binary base; bundle improves only 0.8%. Whitespace and checks-only do not run these optimization passes, so their apparent elapsed/RSS changes are inactive controls, not native-liveness gains. Prior enabled liveness is slower than prior base in every task in this screen, including inactive controls; that variation cautions against treating across-binary changes as attributable wins.

| Task | `words` elapsed vs `base` | `words` elapsed vs `prior_live` | `words` RSS vs `base` |
|---|---:|---:|---:|
| typed-app | +1.3% | -1.8% | -2.6% |
| bundle | -0.8% | -3.1% | -0.7% |
| downlevel-whitespace | -0.2% | -6.8% | +6.9% |
| downlevel-simple | -4.6% | -4.8% | -0.1% |
| checks-only | -1.2% | -3.3% | +1.8% |

Negative deltas mean lower time/RSS. Simple-transpilation elapsed ranges are 3.624–3.770 s for `words`, 3.795–3.863 s for `base` and 3.704–3.893 s for `prior_live`; the prior-enabled range overlaps. Other active-task elapsed ranges also overlap. This warrants a larger paired confirmation rather than extrapolating a broad gain. RSS is mixed: typed falls 2.6% versus base, bundle falls 0.7%, simple is flat, while inactive whitespace rises 6.9% and checks rises 1.8%.

Java/Go ratios above 1 favor Go. **These use older stored warmed Java measurements from 2026-10-18 22:08 UTC against current Go measurements from 2026-10-19 11:05 UTC**, on the same recorded host/toolchain; they are not a fresh paired Java confidence estimate. Stock Java creates a fresh compiler and reparses externs each iteration. Go elapsed is the median fresh CLI wall time. Go RSS is median fresh-process peak RSS; Java RSS is the stored warm JVM’s whole-lifetime peak, including warmup. The RSS ratio describes those different windows, not equivalent per-compilation live memory. Java remains faster in every row.

| Task | Variant | Go s | Peak MiB | Java/Go time | Java/Go RSS |
|---|---|---:|---:|---:|---:|
| typed-app | prior_base | 5.036 | 280.7 | 0.743 | 3.88 |
| typed-app | prior_live | 5.087 | 276.0 | 0.735 | 3.95 |
| typed-app | base | 4.930 | 283.7 | 0.759 | 3.84 |
| typed-app | words | 4.993 | 276.3 | 0.749 | 3.95 |
| bundle | prior_base | 18.671 | 435.2 | 0.889 | 3.90 |
| bundle | prior_live | 19.031 | 408.0 | 0.872 | 4.16 |
| bundle | base | 18.581 | 431.6 | 0.893 | 3.93 |
| bundle | words | 18.438 | 428.6 | 0.900 | 3.96 |
| downlevel-whitespace | prior_base | 1.216 | 171.0 | 0.771 | 5.44 |
| downlevel-whitespace | prior_live | 1.257 | 168.4 | 0.747 | 5.52 |
| downlevel-whitespace | base | 1.174 | 158.0 | 0.799 | 5.89 |
| downlevel-whitespace | words | 1.172 | 168.8 | 0.801 | 5.51 |
| downlevel-simple | prior_base | 3.786 | 213.1 | 0.741 | 5.25 |
| downlevel-simple | prior_live | 3.861 | 206.5 | 0.727 | 5.41 |
| downlevel-simple | base | 3.851 | 212.2 | 0.728 | 5.27 |
| downlevel-simple | words | 3.674 | 212.0 | 0.763 | 5.28 |
| checks-only | prior_base | 7.698 | 658.7 | 0.784 | 2.06 |
| checks-only | prior_live | 7.759 | 654.9 | 0.778 | 2.08 |
| checks-only | base | 7.594 | 666.4 | 0.795 | 2.04 |
| checks-only | words | 7.503 | 678.7 | 0.805 | 2.00 |

Validation: focused V2 race suite **931 PASS, 2 SKIP**; all five CLI tasks match pinned stock Java at compiler job counts **1, 2 and 4**. Coverage includes exact-capacity adoption, zero/one/two-word boundaries, trimming and retained Size, hash/clone/growth behavior, disjoint mutation, native solver state/step/escaped order and failure publication, plus consumer/fallback tests. These correctness checks establish compatibility for the exercised paths, not a performance win.
