Independent review passed for the BigQuery CONTROL/traversal-name-switch screen. This is one cohort with five measured requests per variant, without a confidence interval or significance claim. **Both measured batches exceed the predeclared 2% JIT flag threshold; the observed elapsed reduction is preliminary and does not establish a settled speed gain.**

| Metric | CONTROL | traversal-name-switch | Candidate / CONTROL |
|---|---:|---:|---:|
| Elapsed median, s | 4.261608284 | 4.174620770 | 0.979588102 |
| CPU median, s | 4.910000 | 4.820000 | 0.981670061 |
| Allocated median, bytes | 2,087,058,360 | 2,071,695,408 | 0.992638945 |
| Lifetime HWM, MiB; includes warmup | 1303.441406 | 1243.777344 | 0.954225743 |
| Current RSS median, MiB; postrequest snapshot | 666.953125 | 613.468750 | 0.919807895 |

Observed elapsed median was 2.0412% lower, CPU median 1.8330% lower, allocation 0.7361% lower, and lifetime HWM 4.5774% lower. Current RSS median changed by -8.0192%; its window differs from lifetime HWM. These are retained screen observations, not established general effects.

| Variant | Measured JIT deltas, ms, in request order | Sum JIT, ms | Sum elapsed, s | Batch JIT / elapsed | Above 2%? |
|---|---|---:|---:|---:|---|
| CONTROL | 318, 170, 20, 124, 65 | 697 | 21.287470096 | 3.274227% | Yes |
| traversal-name-switch | 270, 146, 11, 88, 185 | 700 | 20.828052225 | 3.360852% | Yes |

Both variants completed 45 matched warmup requests. Every warm and measured request matched all six oracle digests. Resource, pressure, provenance, and teardown review passed; driver exits and OOM deltas were zero. The JIT flags are annotations under the predeclared note, not new rejection rules. The original PASS is preserved, and no outcome was discarded or retried.

Retained evidence: [independent review](review.json), [all request rows](rows.csv), [cohort medians](cohorts.csv), [ratios](ratios.csv), and [compact machine-readable summary](compact-summary.json).
