# Inline BitSet and fresh-template vector screen

All **90 fresh CLI runs matched pinned stock Java**, with all raw process exit codes equal to zero and 30 MATCH task/variant gates. The three raw samples in every cell have identical artifact/stdout/stderr/exit digests across all variants of each task. Timing completed successfully at 2026-10-19 11:42:03 UTC. **Results are mixed, Java remains faster, and neither experiment is promoted.**

Six variants × five unchanged tasks × three randomized/rotated fresh CLI runs per cell; default runtime GC, no workers or parsed extern cache, no concurrent builds/tests/profiling during timing. Compiler jobs=1, CPU affinity 0,1 under a two-CPU quota; GOGC/GOMEMLIMIT/GODEBUG/GOMAXPROCS unset. Go 1.27.2, shipped PGO and identical four-package `-l=4` flags. Each variant retains the same nine earlier native factors plus CLI_SOURCE_LOCK; newer unrelated factors remain off.

| Variant | Production source | BitSet | Native liveness | Fresh template values |
|---|---|---|---|---|
| prior | b2e407f6 | original slice | off | absent |
| inline | a6afc95f | two inline words | off | absent |
| inline_live | a6afc95f | two inline words | on | absent |
| template_control | 3ebe6449 | two inline words | off | off |
| template | 3ebe6449 | two inline words | off | on |
| compound | 3ebe6449 | two inline words | on | on |

Test-only followup source is 63dc7116. Provenance and commands: [inline-bitset-template-screen-protocol.md](inline-bitset-template-screen-protocol.md), [build-inline-bitset-provenance.json](build-inline-bitset-provenance.json), [build-bitset-template-provenance.json](build-bitset-template-provenance.json), and [inline-bitset-template-screen/](inline-bitset-template-screen/).

The inline representation embeds two common BitSet words, retaining logical capacity separately; empty sets grow by 8 bytes, one-word sets save one object without reducing allocated bytes, two-word sets save 8 bytes and one object, and larger sets add 32 bytes plus one object. The fresh-template prototype uses a fresh, exactly sized template-value vector for each specialization in the private inference path. Same-binary `inline_live/inline` isolates native liveness; `template/template_control` isolates enabling template vectors. `inline/prior` includes the unconditional representation and code-generation changes. `template_control/inline` includes binary/layout/noise differences and is not a vector-path benefit. The six variants do not isolate every interaction of a full three-factor matrix.

No broad Java win or promotion is supported. Inline BitSet is mixed: typed elapsed falls about 2.5% and RSS 13.6% versus prior, while bundle elapsed rises about 2.0% and RSS 3.7%. Typed ranges overlap (inline 4.878–5.072 s, prior 4.923–5.134 s), so the small time difference needs independent confirmation. Template vectors regress typed elapsed about 15.8% versus their same-binary control, and the compound retains most of that regression. The control and enabled typed ranges do not overlap in this three-sample screen. Checks-only also regresses: template is 22.6% slower than its control (9.127 vs 7.443 s), while compound is 25.4% slower.

A source audit identified repeated GetErrorCount monitor calls in the V1 template eligibility guard, adding two runtime.Stack operations per scope. A V2 guard change is being prepared, **not measured here**; it is not yet established how much of the observed regression that change would remove. No V2 result is mixed into these numbers.

| Task | inline/prior time | inline/prior RSS | inline_live/inline time | template/control time | compound/control time |
|---|---:|---:|---:|---:|---:|
| typed-app | -2.5% | -13.6% | +0.1% | +15.8% | +13.5% |
| bundle | +2.0% | +3.7% | +0.2% | +2.3% | +0.1% |
| downlevel-whitespace | -8.5% | +0.2% | +2.4% | -1.3% | -0.4% |
| downlevel-simple | -4.3% | -1.7% | -1.6% | +0.3% | +2.1% |
| checks-only | -0.8% | -0.7% | -0.7% | +22.6% | +25.4% |

Negative deltas mean lower time/RSS. Native liveness is active only in typed ADVANCED and the two SIMPLE tasks; it is inactive in whitespace/checks-only. Fresh-template inference is active only in typed ADVANCED+checkTypes and checks-only+checkTypes; it is inactive in bundle SIMPLE, downlevel SIMPLE and whitespace. Template on/off differences in all three inactive tasks are controls, not vector-path effects. Their flag-related fluctuations are not wins. Inline BitSet itself is unconditional and may still affect those workloads. Three observations per cell are a screen, not a confidence interval, and apparent factor improvements cannot simply be added.

Java/Go ratios above 1 favor Go. Java is the **older stored stock warm JVM cohort from 2026-10-18 22:08 UTC** (12 warmup/10 measured compilations), compared with **current Go runs from 2026-10-19 11:30 UTC** on the same recorded host/toolchain. This is not a fresh paired Java confidence estimate. Java creates a fresh compiler and reparses externs each iteration. Go RSS is median fresh CLI lifetime peak; Java RSS is the stored warm JVM’s entire lifetime peak, including warmup. These RSS ratios describe different measurement windows, not equivalent per-compilation live memory. Java remains faster in every row.

| Task | Variant | Go s | Peak MiB | Java/Go time | Java/Go RSS |
|---|---|---:|---:|---:|---:|
| typed-app | prior | 5.081 | 290.0 | 0.736 | 3.76 |
| typed-app | inline | 4.954 | 250.7 | 0.755 | 4.35 |
| typed-app | inline_live | 4.960 | 276.8 | 0.754 | 3.94 |
| typed-app | template_control | 5.049 | 275.8 | 0.741 | 3.95 |
| typed-app | template | 5.846 | 276.0 | 0.640 | 3.95 |
| typed-app | compound | 5.732 | 272.7 | 0.652 | 4.00 |
| bundle | prior | 18.079 | 417.0 | 0.918 | 4.07 |
| bundle | inline | 18.449 | 432.5 | 0.900 | 3.92 |
| bundle | inline_live | 18.478 | 437.9 | 0.898 | 3.87 |
| bundle | template_control | 18.130 | 429.0 | 0.916 | 3.95 |
| bundle | template | 18.541 | 425.9 | 0.895 | 3.98 |
| bundle | compound | 18.155 | 435.6 | 0.914 | 3.89 |
| downlevel-whitespace | prior | 1.310 | 169.9 | 0.716 | 5.47 |
| downlevel-whitespace | inline | 1.199 | 170.3 | 0.783 | 5.46 |
| downlevel-whitespace | inline_live | 1.228 | 167.4 | 0.764 | 5.56 |
| downlevel-whitespace | template_control | 1.285 | 170.1 | 0.730 | 5.47 |
| downlevel-whitespace | template | 1.269 | 162.7 | 0.739 | 5.72 |
| downlevel-whitespace | compound | 1.281 | 162.1 | 0.733 | 5.74 |
| downlevel-simple | prior | 3.917 | 216.3 | 0.716 | 5.17 |
| downlevel-simple | inline | 3.749 | 212.6 | 0.748 | 5.26 |
| downlevel-simple | inline_live | 3.691 | 208.3 | 0.760 | 5.37 |
| downlevel-simple | template_control | 3.752 | 213.9 | 0.748 | 5.23 |
| downlevel-simple | template | 3.763 | 211.0 | 0.745 | 5.30 |
| downlevel-simple | compound | 3.829 | 214.9 | 0.733 | 5.20 |
| checks-only | prior | 7.396 | 668.7 | 0.816 | 2.03 |
| checks-only | inline | 7.336 | 664.4 | 0.823 | 2.05 |
| checks-only | inline_live | 7.283 | 649.8 | 0.829 | 2.09 |
| checks-only | template_control | 7.443 | 665.3 | 0.811 | 2.04 |
| checks-only | template | 9.127 | 634.2 | 0.662 | 2.14 |
| checks-only | compound | 9.331 | 657.8 | 0.647 | 2.07 |

Correctness validation reported by the coordinator: inline-focused race suite **923 PASS, 2 SKIP**; combined suite **942 PASS, 2 SKIP**, with no race report; all five tasks at compiler jobs **1, 2 and 4** match pinned stock Java (15 gates). An initial new fixture omitted a class closing brace, causing two parsing failures plus its parent-suite failure. That fixture was fixed in test-only source; the complete rerun is clean and original failed logs are retained. Runtime candidate code was not changed to weaken the oracle.

Final summary and all 90 raw zero exits/digests plus all 30 gates were verified after the coordinator reported timing completion. No builds, tests, profiling or timing were performed to prepare this report. At this screen’s completion, template V2 had not been built or measured; production source for this screen remains 3ebe6449 with only the later test-fixture correction 63dc7116.
