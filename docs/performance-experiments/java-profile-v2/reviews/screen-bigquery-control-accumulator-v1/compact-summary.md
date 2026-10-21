Independent review passed for the BigQuery CONTROL/accumulator screen. This is one cohort with five measured requests per variant; it provides no confidence interval or significance claim. The accumulator is not selected as a speed winner.

| Metric | CONTROL | accumulator | Candidate / CONTROL |
|---|---:|---:|---:|
| Elapsed median, s | 4.163091146 | 4.237157120 | 1.017791101 |
| CPU median, s | 4.880000 | 4.880000 | 1.000000000 |
| Allocated median, bytes | 2,079,177,752 | 2,021,800,312 | 0.972403783 |
| Lifetime HWM, MiB; includes warmup | 1334.335938 | 1318.394531 | 0.988052929 |
| Current RSS median, MiB; postrequest snapshot | 702.003906 | 776.714844 | 1.106425245 |

The candidate elapsed median was 1.7791% higher, CPU median was equal, allocation was 2.7596% lower, and observed lifetime HWM was 1.1947% lower. Current RSS median was 10.6425% higher; its postrequest snapshot differs from the lifetime HWM window. These are screen observations, not established general effects.

| Variant | Measured JIT deltas, ms, in request order | Sum JIT, ms | Sum elapsed, s | Batch JIT / elapsed | Above 2%? |
|---|---|---:|---:|---:|---|
| CONTROL | 33, 122, 80, 115, 60 | 410 | 21.066121810 | 1.946253% | No |
| accumulator | 36, 69, 58, 108, 107 | 378 | 21.205412861 | 1.782564% | No |

Both variants completed 54 matched warmup requests. Every warm and measured request matched all six oracle digests. Resource, pressure, provenance, and teardown review passed; driver exits and OOM deltas were zero. Neither measured batch exceeded the predeclared JIT flag threshold. No postmeasurement drift rule or new rejection was added; no outcome was discarded.

Retained evidence: [independent review](review.json), [all request rows](rows.csv), [cohort medians](cohorts.csv), [ratios](ratios.csv), and [compact machine-readable summary](compact-summary.json).
