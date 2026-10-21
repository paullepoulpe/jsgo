Independent review passed for the three-cohort BigQuery traversal-name-switch / published STOCK confirmation. **No demonstrated speed gain:** elapsed cohort ratios were 1.025077180, 0.938362416 and 1.004951924; their geometric mean was 0.988759840 with 95% log-cohort t(df=2) interval [0.881161793, 1.109496609]. The interval includes useful gains and regressions. Two of three cohorts had higher candidate elapsed medians. Five of six measured batches exceeded the predeclared JIT annotation; a settled gain is not established.

| Cohort | Variant | Warmups | Wall median s | Wall mean s | CPU median s | Allocated median bytes | Approx. lifetime HWM MiB | Current RSS median MiB |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| 1 | STOCK | 41 | 4.177404776 | 4.191072802 | 4.950000 | 2,088,974,152 | 1258.187500 | 655.152344 |
| 1 | traversal-name-switch | 41 | 4.282162307 | 4.311533827 | 4.980000 | 2,073,739,080 | 1320.765625 | 634.878906 |
| 2 | STOCK | 41 | 4.514256396 | 4.451430792 | 4.990000 | 2,082,348,008 | 1248.027344 | 762.781250 |
| 2 | traversal-name-switch | 41 | 4.236008536 | 4.270363684 | 4.870000 | 2,094,657,080 | 1313.074219 | 856.589844 |
| 3 | STOCK | 40 | 4.141000490 | 4.202893460 | 4.760000 | 2,087,996,152 | 1210.722656 | 945.894531 |
| 3 | traversal-name-switch | 40 | 4.161506408 | 4.159055219 | 4.890000 | 2,079,497,640 | 1286.218750 | 611.292969 |

| Cohort | Elapsed ratio | CPU ratio | Allocated ratio | Approx. lifetime HWM ratio | Current RSS median ratio | Mean-wall ratio |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| 1 | 1.025077180 | 1.006060606 | 0.992706912 | 1.049736724 | 0.969055384 | 1.028742289 |
| 2 | 0.938362416 | 0.975951904 | 1.005911150 | 1.052119751 | 1.122982302 | 0.959323841 |
| 3 | 1.004951924 | 1.027310924 | 0.995929824 | 1.062356224 | 0.646259121 | 0.989569509 |

| Metric | Geometric mean candidate/STOCK | 95% log-t(df=2) interval | Discrete whole-cohort bootstrap 95% |
| --- | ---: | --- | --- |
| elapsed | 0.988759840 | [0.881161793, 1.109496609] | [0.938362416, 1.025077180] |
| CPU | 1.002885815 | [0.940652297, 1.069236700] | [0.975951904, 1.027310924] |
| RSS_lifetimeHWM | 1.054723390 | [1.038224727, 1.071484236] | [1.049736724, 1.062356224] |
| allocated_bytes | 0.998166834 | [0.981236216, 1.015389579] | [0.992706912, 1.005911150] |

Mean-wall sensitivity: geometric mean 0.992139435, log-t95(df=2) [0.909512830, 1.082272426]. Approximate lifetime HWM was 4.97–6.24% higher in every cohort; geometric mean 5.47% higher. Its descriptive interval is also above 1, but this secondary telemetry is approximate and not an exact physical-memory effect. Allocation and CPU show no established improvement. Current RSS snapshots differ from the lifetime window.

| Cohort | Variant | All measured JIT deltas ms | Sum JIT ms | Sum elapsed s | JIT / elapsed | >2% flag |
| --- | --- | --- | ---: | ---: | ---: | --- |
| 1 | STOCK | 149, 59, 49, 39, 99 | 395 | 20.955364012 | 1.884959% | No |
| 1 | traversal-name-switch | 102, 52, 12, 56, 234 | 456 | 21.557669134 | 2.115257% | Yes |
| 2 | STOCK | 122, 44, 41, 82, 267 | 556 | 22.257153960 | 2.498073% | Yes |
| 2 | traversal-name-switch | 124, 103, 96, 183, 83 | 589 | 21.351818422 | 2.758547% | Yes |
| 3 | STOCK | 157, 134, 73, 61, 43 | 468 | 21.014467302 | 2.227037% | Yes |
| 3 | traversal-name-switch | 115, 110, 100, 37, 70 | 432 | 20.795276094 | 2.077395% | Yes |

Every warm and measured request passed all six oracle digests: 244 warm requests and 30 measured requests across six JVMs. Resource, pressure, loaded-origin and stopped/zero-new-OOM review passed. One existing-policy telemetry exception was admitted: STOCK cohort 1 VmHWM fell 168 KiB between warmup requests 6 and 7 (1,100,692 to 1,100,524 KiB), then rose to 1,288,384 KiB at request 8. The raw rows and later maximum are unchanged. No measured VmHWM dip was admitted, and no new rule was introduced.

The three fresh independent pair-of-JVM cohorts are the inference units; the 30 measured requests are not 30 independent JVM comparisons. The log-t interval uses exp(mean(log ratios) ± 4.302652729911275 × sample-SD(log ratios)/sqrt(3)); it assumes approximately normal independent log effects, which three observations cannot establish. The bootstrap is discrete and confined to these three observed cohorts, so it cannot be the sole inference. Screens are not pooled. Secondary intervals are descriptive, without a multiplicity-adjusted memory claim.

JIT flags preserve the original protocol PASS and all outcomes; no row or cohort was discarded, replaced, or retried. The JIT counter sums compiler work and is not elapsed compile time. Allocated bytes are VM-wide ThreadMXBean request deltas, not retained heap. RSS telemetry is approximate; lifetime HWM includes warmup/JIT, current RSS is a postrequest snapshot, and missing PSI remains unavailable rather than zero.

Retained evidence: [independent review](review.json), [all request rows](rows.csv), [cohort medians](cohorts.csv), [ratios and intervals](ratios.csv), [measured JIT](measured-jit.csv), [compact exact summary](compact-summary.json), and [report transform](summarize-confirmation.py). Bundle generalization remains a separate workload with no result in this review.
