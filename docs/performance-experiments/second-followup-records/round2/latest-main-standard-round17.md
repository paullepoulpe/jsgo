# Round 17: latest-main standard-build baseline

Source is `ff06ac15304f14c756bce9cd36ad2205212f51da`. This uses standard compiler inlining, the new shipped PGO and Go 1.27.2, with no experiment flags, workers, caches or historical `-l4` flags. Runtime GC and GOMAXPROCS remain at their defaults. CPU affinity is 0,1 under a two-CPU quota.

The five-task screen has three fresh Go runs per task at compiler thread count one. A separate balanced whitespace screen has five runs at each of thread counts 1, 2 and 4. All 30 measured outputs matched Java. Comparisons use stored warm Java references; these are screening results, not fresh paired cohorts or confidence estimates.

| Task N1 | Go s | Range s | Go RSS MiB | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|---:|
| typed-app | 5.340 | 5.209–5.448 | 253.5 | 0.700 | 4.30 |
| bundle | 18.607 | 18.541–18.969 | 543.5 | 0.892 | 3.12 |
| downlevel-whitespace | 1.333 | 1.321–1.342 | 154.8 | 0.704 | 6.01 |
| downlevel-simple | 3.989 | 3.944–4.045 | 199.1 | 0.703 | 5.62 |
| checks-only | 8.183 | 8.180–8.187 | 675.7 | 0.738 | 2.01 |

## Balanced whitespace thread screen

Each row below uses its own Java reference with the **same compiler thread count**, from the earlier thread cohort with 12 warmups and 10 measured iterations. Its N1 reference differs from the five-task baseline reference above; keep those denominators separate. Each row has five fresh Go runs.

| N | Go s | Range s | Go RSS MiB | Java s | Java RSS MiB | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|---:|---:|---:|
| 1 | 1.249 | 1.235–1.352 | 156.3 | 0.896 | 990.0 | 0.717 | 6.33 |
| 2 | 1.273 | 1.214–1.283 | 160.0 | 0.882 | 917.8 | 0.693 | 5.74 |
| 4 | 1.240 | 1.215–1.287 | 164.1 | 0.848 | 869.3 | 0.683 | 5.30 |

The thread parameter bounds concurrent compiler work items. It does not set GOMAXPROCS or cap the total number of goroutines or OS threads. This host has a two-CPU quota, so these observations do not establish scaling on a machine with four available cores. In this small cohort N2 is slower than N1, while N4 is slightly faster. The ranges overlap; there is no confident CPU-scaling claim.

All ratios are Java divided by Go; values above one favor Go. Go RSS is the median fresh CLI lifetime peak, while Java RSS is the whole warm JVM lifetime peak including warmup. These memory windows differ.

Standard latest main still trails the stored warm Java references. No exhaustive performance ceiling is claimed. Historical prototype percentages are not incremental gains on this new main, and previously reported Mac or other-VM measurements are not combined with this machine. A build using the historical `-l4` configuration is being prepared separately to distinguish build effects from source changes; it is not completed evidence in this report.
