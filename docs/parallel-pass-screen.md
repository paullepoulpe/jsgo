# Parallel inlining and reference lookup screen

Two fresh CLI runs per variant, alternating order, with two compiler jobs,
the same binary, default GC and the same two-CPU affinity/quota. These are
incremental Go comparisons. Java artifact oracles check stdout, stderr, status,
JavaScript and source maps; no Java timing ratios are inferred from this screen.
The separate matched Java/Go sweep is in [performance.md](performance.md).

| Workload | Control time / peak RSS | Native reference index | Index + parallel inlining |
|---|---:|---:|---:|
| typed-app | 5.060 s / 279.8 MiB | 5.459 s / 268.8 MiB | 5.516 s / 271.7 MiB |
| bundle | 18.748 s / 589.9 MiB | 19.198 s / 509.9 MiB | 19.213 s / 571.6 MiB |

The reference index regressed in this screen, so a second alternating screen
removed it to isolate parallel inlining:

| Workload | Control time / peak RSS | Parallel inlining alone |
|---|---:|---:|
| typed-app | 5.167 s / 264.6 MiB | 5.200 s / 275.1 MiB |
| bundle | 18.668 s / 540.8 MiB | 18.483 s / 560.6 MiB |

Whole-script parallel inlining alone is essentially flat: +0.6% elapsed time
on typed-app and -1.0% on bundle. Two runs cannot establish that small a gain,
and RSS increases in both isolated rows. Neither experiment is promoted to the
draft PR. RSS is median fresh-process peak, not live heap or allocated bytes.

All measured artifacts match the pinned Java oracle. The enabled CLI candidate
also matches fresh Java at two jobs on SIMPLE, ADVANCED and dependency-pruned
fixtures. Source review finds no permanent CLI eligibility fallback; normal
compiler, observers, colors and synthetic inputs qualify. Separate tracer-mode
measurements intentionally use the serial path, so they cannot profile this
parallel implementation's timing.

Raw measurements and machine/binary metadata are in
[parallel-pass-screen.json](parallel-pass-screen.json). Ownership and fallback
constraints, including partial AST state after an unexpected internal exception,
are documented in [experiments/parallel-fsiv.md](experiments/parallel-fsiv.md).

The script-level workload can be imbalanced. Finer function cohorts remain a
source design, with canonical scope, implicit-slot and ancestor/descendant
ordering requirements; they have not been implemented or measured here.
