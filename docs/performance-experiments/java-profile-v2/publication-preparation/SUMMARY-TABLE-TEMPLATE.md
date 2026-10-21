Draft report data and remaining placeholders. Completed rows below come from independent reviews of eight BigQuery screens through reference-list-two, the bundle accumulator screen, and the finished three-cohort BigQuery confirmation. Keep the original PASS/FAIL outcomes and all diagnostic annotations. The running v4 bundle holdout and any later results are outside this cutoff.

Completed confirmation finding: **the identifier dispatch guard did not establish a speed gain against published STOCK**. Its geometric mean elapsed ratio was 0.988759840, with 95% log-cohort t(df=2) interval [0.881161793, 1.109496609]. Approximate lifetime HWM was higher in all three cohorts. Other candidate work and separate bundle generalization remain outside this completed result. This is a Java-versus-Java study of Closure v20261015; it does not establish a new Java/Go ratio.

Completed BigQuery screens: one fresh paired-JVM cohort, five measured requests per variant, candidate / same-family CONTROL; below 1 means lower cost. These are exploratory, without confidence intervals or significance claims.

| Candidate | Elapsed median ratio | CPU median ratio | Allocated median ratio | Approx. lifetime HWM ratio | Current RSS median ratio | Independent review | Measured JIT CONTROL / candidate |
| --- | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Shared parameter accumulator | 1.017791101 | 1.000000000 | 0.972403783 | 0.988052929 | 1.106425245 | PASS | 1.946253% / 1.782564%; neither flagged |
| Identifier dispatch guard (`traversal-name-switch`) | 0.979588102 | 0.981670061 | 0.992638945 | 0.954225743 | 0.919807895 | PASS | 3.274227% / 3.360852%; both flagged |
| Guaranteed-empty callback bypass (`traversal-noop`) | 0.986064334 | 1.022312373 | 1.001319747 | 1.265146397 | 0.979156294 | PASS | 3.512882% / 4.427500%; both flagged |
| Leaf identifier recursion bypass (`traversal-leaf-name`) | 0.994452727 | 0.997933884 | 1.010159249 | 1.057753552 | 0.959077630 | PASS | 1.534589% / 2.052817%; candidate flagged |
| Prototype/scalar equality shortcut (`equality-prototype-scalar`) | 1.013348675 | 0.966535433 | 0.993685628 | 0.997544935 | 1.132016383 | PASS | 1.7943% / 1.5774%; neither flagged |
| Empty-parameter shortcut (`empty-parameters`) | 1.007876370 | 1.048728814 | 1.000269872 | 0.925947483 | 1.102701509 | PASS | 1.795339% / 3.159669%; candidate flagged |
| Last source-position reuse (`source-position-last`) | 1.008375720 | 1.008333333 | 0.985494061 | 0.815518598 | 0.653749205 | PASS | 42.771700% / 1.855590%; CONTROL flagged |
| Reference-list initial capacity two (`reference-list-two`) | 1.037192071 | 1.012244898 | 1.001337697 | 1.129029255 | 0.995090237 | PASS | 2.044633% / 1.435168%; CONTROL flagged |

The accumulator reduced measured request allocation totals by 2.7596% in this screen but was not selected as a speed winner; its elapsed median increased 1.7791%. The identifier dispatch screen observed 2.0412% lower elapsed median, with both measured batches above the predeclared JIT threshold; it does not establish a settled gain. Both outcomes and unfavorable metrics remain in the report.

Reference-list-two completed 62 matched warmups and five measured pairs, all 134 requests matching six digests. Elapsed/CPU/allocation medians and approximate lifetime HWM were higher; no benefit was established. Its matching compact review preserves all absolute metrics, CPU means, ordered JIT counters and clean resource/provenance/cleanup checks. One existing-policy warmup-only VmHWM exception remains: candidate warmup 8→9 fell 24 KiB (1,209,936→1,209,912), then warmup 29 strictly exceeded it at 1,282,072 KiB, proven by exact retained rows. No measured exception, new rule, smoothing or data exclusion was added.

The later callback and leaf-identifier screens were also left unpromoted, with unfavorable allocation/CPU/RSS observations preserved. The equality candidate's elapsed median increased 1.3349%; its lower CPU median does not establish a benefit, and CPU means were nearly equal (candidate only 0.2003% lower). Exact per-variant medians, all JIT deltas and interpretations for these three screens are in their matching `reviews/*/compact-summary.txt` files and the main README. The detailed tables immediately below retain the first two screens; they are not the full result set. [Accumulator interpretation](ACCUMULATOR-ALLOCATION-INTERPRETATION.md) provides qualified allocation-focused wording without changing v4's elapsed primary endpoint.

Empty-parameters and source-position-last each completed 41 matched warmups and five measured pairs, with all 92 requests matching six digests and clean provenance/resource/cleanup review. Neither had an elapsed lead. Source-position CONTROL recorded an 8,554 ms cumulative compilation-time counter delta on its fifth measured request: its 42.7717% batch JIT fraction must remain visible beside the lower allocation/RSS observations. CompilationMXBean records approximate accumulated compilation elapsed time, potentially summed across threads, not CPU%. A compilation spanning a container pause may inflate a later counter update; this mechanism was not diagnosed for the observed jump. Treat the ratio as a potentially contaminated elapsed-time proxy without a JIT CPU-share or slowdown-causality claim; see ../JIT-COUNTER-INTERPRETATION.md. All rows and PASS outcomes remain; neither screen establishes a settled benefit. Their matching `reviews/*/compact-summary.md`, `compact-summary.json` and `measured-jit.csv` retain absolute medians, CPU means, source identities where checked, all ordered JIT counters and limitations.

| Screen | Variant | Matched warmups | Elapsed median s | CPU median s | Allocated median bytes | Approx. lifetime HWM MiB | Current RSS median MiB |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: |
| accumulator | CONTROL | 54 | 4.163091146 | 4.880000 | 2,079,177,752 | 1334.335938 | 702.003906 |
| accumulator | accumulator | 54 | 4.237157120 | 4.880000 | 2,021,800,312 | 1318.394531 | 776.714844 |
| identifier dispatch | CONTROL | 45 | 4.261608284 | 4.910000 | 2,087,058,360 | 1303.441406 | 666.953125 |
| identifier dispatch | traversal-name-switch | 45 | 4.174620770 | 4.820000 | 2,071,695,408 | 1243.777344 | 613.468750 |

| Screen / variant | All measured JIT deltas ms, in order | Sum JIT ms | Sum elapsed s | Batch JIT / elapsed | >2% flag |
| --- | --- | ---: | ---: | ---: | --- |
| accumulator / CONTROL | 33, 122, 80, 115, 60 | 410 | 21.066121810 | 1.946253% | No |
| accumulator / accumulator | 36, 69, 58, 108, 107 | 378 | 21.205412861 | 1.782564% | No |
| identifier dispatch / CONTROL | 318, 170, 20, 124, 65 | 697 | 21.287470096 | 3.274227% | Yes |
| identifier dispatch / traversal-name-switch | 270, 146, 11, 88, 185 | 700 | 20.828052225 | 3.360852% | Yes |

Both screens passed six oracle digests for every warm and measured request, resource/pressure/provenance review, zero new OOM, and stopped cleanup. The >2% flag is a diagnostic annotation, not a new rejection rule. JIT totals can sum work across compiler threads and are not elapsed compile time. Source: `reviews/screen-bigquery-control-accumulator-v1/` and `reviews/screen-bigquery-control-traversal-name-switch-v1/`, especially compact-summary.json, review.json, rows.csv, cohorts.csv and ratios.csv.

Completed direct-STOCK confirmation: three fresh paired-JVM cohorts, five measured requests per variant; screens excluded from inference. Matched warmup counts were 41, 41 and 40 per variant.

| Candidate / workload | Cohort | STOCK elapsed median s | Candidate elapsed median s | Candidate/STOCK elapsed | CPU ratio | Allocation ratio | Approx. lifetime HWM ratio | Current RSS ratio | STOCK / candidate measured JIT fractions and flags |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| traversal-name-switch / BigQuery | 1 | 4.177404776 | 4.282162307 | 1.025077180 | 1.006060606 | 0.992706912 | 1.049736724 | 0.969055384 | 1.884959% / 2.115257%; candidate flagged |
| traversal-name-switch / BigQuery | 2 | 4.514256396 | 4.236008536 | 0.938362416 | 0.975951904 | 1.005911150 | 1.052119751 | 1.122982302 | 2.498073% / 2.758547%; both flagged |
| traversal-name-switch / BigQuery | 3 | 4.141000490 | 4.161506408 | 1.004951924 | 1.027310924 | 0.995929824 | 1.062356224 | 0.646259121 | 2.227037% / 2.077395%; both flagged |

| Independent-cohort analysis | Raw cohort ratios | Geometric mean | 95% log-cohort t(df=2) interval | Whole-cohort bootstrap 95% | Finding / limitations |
| --- | --- | --- | --- | --- | --- |
| Median elapsed, primary | 1.025077180, 0.938362416, 1.004951924 | 0.988759840 | [0.881161793, 1.109496609] | [0.938362416, 1.025077180] | No established speed gain; two cohorts above 1 |
| Median CPU, descriptive | 1.006060606, 0.975951904, 1.027310924 | 1.002885815 | [0.940652297, 1.069236700] | [0.975951904, 1.027310924] | No established improvement |
| Median allocated bytes, descriptive | 0.992706912, 1.005911150, 0.995929824 | 0.998166834 | [0.981236216, 1.015389579] | [0.992706912, 1.005911150] | No established improvement |
| Approx. lifetime HWM, descriptive | 1.049736724, 1.052119751, 1.062356224 | 1.054723390 | [1.038224727, 1.071484236] | [1.049736724, 1.062356224] | Higher in all three cohorts; approximate secondary telemetry |

| Confirmation cohort / variant | All measured JIT deltas ms | Sum JIT ms | Sum elapsed s | JIT / elapsed | >2% flag |
| --- | --- | ---: | ---: | ---: | --- |
| 1 / STOCK | 149, 59, 49, 39, 99 | 395 | 20.955364012 | 1.884959% | No |
| 1 / traversal-name-switch | 102, 52, 12, 56, 234 | 456 | 21.557669134 | 2.115257% | Yes |
| 2 / STOCK | 122, 44, 41, 82, 267 | 556 | 22.257153960 | 2.498073% | Yes |
| 2 / traversal-name-switch | 124, 103, 96, 183, 83 | 589 | 21.351818422 | 2.758547% | Yes |
| 3 / STOCK | 157, 134, 73, 61, 43 | 468 | 21.014467302 | 2.227037% | Yes |
| 3 / traversal-name-switch | 115, 110, 100, 37, 70 | 432 | 20.795276094 | 2.077395% | Yes |

Mean-wall sensitivity gives geometric mean 0.992139435, log-t95(df=2) [0.909512830, 1.082272426]. All 244 warm and 30 measured requests passed six digests; resource/pressure/provenance/stopped/zero-new-OOM review passed. One existing-policy exception remains disclosed: cohort 1 STOCK VmHWM fell 168 KiB between warmup requests 6 and 7, then was strictly exceeded at request 8. Raw values and the retained maximum were unchanged; no measured dip or new rejection rule was admitted.

Five of six batches carry the predeclared JIT annotation; all outcomes and the original PASS are retained. Three independent cohort effects are the inference units, not 30 independent JVM comparisons. The log-t interval uses critical t=4.302652729911275 and assumes approximately normal independent log effects; three observations cannot establish that assumption. Three-cohort bootstrap intervals are discrete and weak. Secondary intervals are descriptive; do not infer a multiplicity-adjusted or exact physical-memory claim. Exact medians/SDs, current RSS, mean sensitivity and JIT data are retained in `reviews/confirm-bigquery-stock-traversal-name-switch-v1/{compact-summary.json,compact-summary.md,review.json,rows.csv,cohorts.csv,ratios.csv,measured-jit.csv}`.

Completed bundle CONTROL screen: one fresh paired cohort, 22 matched warmups and five measured pairs; all 54 requests matched six digests, and provenance/resource/pressure-proxy/cleanup review passed. No new OOM, driver exit failure or telemetry anomaly; PSI unavailable.

| Bundle accumulator metric | CONTROL | Candidate | Candidate / CONTROL |
| --- | ---: | ---: | ---: |
| Elapsed median s | 11.158163928 | 10.728804040 | 0.961520561 |
| Elapsed mean s | 11.043794191 | 10.860059846 | 0.983363114 |
| CPU median s | 12.880 | 12.360 | 0.959627329 |
| CPU mean s | 12.832 | 12.548 | 0.977867830 |
| Allocated median bytes | 4,684,739,288 | 4,367,995,056 | 0.932388077 |
| Approx. lifetime HWM MiB | 1303.429688 | 1371.406250 | 1.052152075 |
| Current RSS median MiB | 633.699219 | 637.539063 | 1.006059411 |

CONTROL JIT deltas `[421,212,163,138,110]` sum to 1,044 ms / 55.218970957 s = 1.890655%; candidate `[282,145,206,148,233]` sum to 1,014 ms / 54.300299231 s = 1.867393%. Neither carries the historical >2% flag; the elapsed-compilation proxy limitations still apply. Median elapsed was 3.8479% lower, but mean elapsed only 1.6637% lower; higher HWM/current RSS remain visible. This is an exploratory CONTROL screen without an independent-cohort confidence interval, not a confirmed speed benefit.

The sole v4 comparison selected before execution is accumulator versus published STOCK on bundle, with five fresh cohorts, 20 measured pairs, unchanged 60–150 warmups/thresholds and predeclared total timeout 21,600 seconds. See ../SELECTED-V4-CONFIRMATION.md. It is running; no partial rows or results are interpreted here. Neither this screen nor BigQuery is pooled into the new inference. Evidence: `reviews/screen-bundle-control-accumulator-v1/{compact-summary.md,compact-summary.json,review.json,rows.csv,cohorts.csv,ratios.csv,measured-jit.csv}`.

Memory labels: allocated bytes are VM-wide ThreadMXBean deltas around each request, not JFR sample weights or retained heap. Approximate lifetime HWM is the maximum observed process VmHWM including warmup/JIT; current RSS is a postrequest snapshot. These windows differ, and approximate RSS telemetry is not an exact physical peak. Missing PSI remains unavailable, not zero.

Profile status: the first BigQuery recording was disturbed by enabling JFR after warmup; the corrected primed recording narrowly failed the final JIT gate (2.1007648% versus 2%) despite 2.3256% drift. Both are diagnostic only. Do not substitute CPU sample shares for wall-time fractions or allocation weights for exact allocated bytes. Bounded audits left ancillary strings unresolved; raw JFR is omitted from publication and retained locally, while aggregates and audit receipts are included. See JFR-PUBLICATION-DISPOSITION.json.

Incomplete/prepared work: `identity-differs` passed builds/probes but its screen was intentionally stopped after two paired warmups and before any measurements; no ratio exists. Preserve identity-screen-stop.json and the interrupted result/cleanup evidence. Other proposals, diagnostics, successful builds/probes and active runs retain their documented statuses; none becomes a completed benchmark through this draft. Full upstream suites were not run.
