# Structural follow-up screen: default runtime GC

The user stopped GC tuning. This screen compares code/design changes with GOGC, GOMEMLIMIT, GOMAXPROCS and GODEBUG removed from the benchmark environment. Each Go sample is a fresh CLI compilation; no worker or pre-parsed extern AST is used.

Production source: `5906ebbb`; follow-up test fixes: `b585e18a`. The prior binary is production `5b720ead`. Go 1.27.2, shipped PGO, same four package-specific `-l=4` build flags, affinity CPUs 0 and 1, two-CPU quota, compiler threads 1. Seven variants, two alternating/rotated runs each, five unchanged workloads: all 70 output/diagnostic/source-map gates matched.

This is an exploratory screen, not a confidence-backed speed claim. Java ratios use the stored stock warm JVM cohort, not a fresh paired cohort. Java creates a new compiler/runner and reparses externs each iteration. Go RSS is median fresh-process peak; Java RSS is JVM-lifetime peak including warmup. Ratios above one favor Go.

| Task | Control s / MiB | Combined s / MiB | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|
| typed-app | 4.925 / 280.7 | 5.043 / 283.0 | 0.742 | 3.85 |
| bundle | 17.972 / 439.1 | 18.000 / 413.3 | 0.922 | 4.10 |
| downlevel-whitespace | 1.168 / 167.7 | 1.167 / 167.1 | 0.804 | 5.57 |
| downlevel-simple | 3.635 / 206.4 | 3.752 / 206.7 | 0.748 | 5.41 |
| checks-only | 7.491 / 666.7 | 7.268 / 680.3 | 0.831 | 2.00 |

The combined configuration adds inference-scope scheduling, contiguous typed declarations, a direct-only final inlining worklist, and chunked reference-discovery events to the first structural batch. Speed changes are mixed: no broad win, and Go still trails the stored warmed Java times. Bundle combined RSS is about 6% below the same-binary control in these two samples, with virtually identical time; this needs repetition. No new feature is promoted to default on this evidence.

Individual factors (seconds / MiB):

| Task | Scope schedule | Typed tables | Direct worklist | Chunked tape |
|---|---:|---:|---:|---:|
| typed-app | 4.860 / 259.0 | 4.982 / 271.1 | 5.146 / 266.2 | 4.986 / 271.3 |
| bundle | 18.490 / 436.1 | 18.180 / 427.4 | 18.486 / 445.2 | 18.136 / 437.6 |
| downlevel-whitespace | 1.133 / 164.2 | 1.158 / 161.0 | 1.153 / 169.7 | 1.234 / 172.1 |
| downlevel-simple | 3.705 / 198.0 | 3.651 / 202.7 | 3.719 / 208.2 | 3.613 / 210.9 |
| checks-only | 7.277 / 666.8 | 7.295 / 656.0 | 7.164 / 683.1 | 7.368 / 656.5 |

The checks-only same-binary control itself regresses relative to the previous binary. Therefore, the combined improvement against that control does not establish an improvement over the prior implementation. Non-inference whitespace changes from scope/typed switches are also useful noise controls, not attributed wins.

Validation: new focused suite passes after correcting two test-only expectations; the broad pre-fix test run has 4,874 passing test/subtest records and 20 skips, with only the known exception-context fixture failure. All five CLI workloads match at compiler thread counts one and two. Race validation for the new follow-up designs is pending.

A separate diagnostic overlay confirms CLI-owned source files activate on all five real tasks; silent fallback does not explain the earlier neutral source-lock screen. Its pass labels and atomic counters are absent from all timed binaries. Profile capture durations are not reported as performance.

Raw results, full variant environments, hashes and machine metadata are in `followup-screen/`; binary provenance is `build-followup-provenance.json`. The following source batch separately explores access-control facts, unused-code task storage, mixed direct/block inlining replay and ordered parallel peepholes; those changes are not part of these measurements.
