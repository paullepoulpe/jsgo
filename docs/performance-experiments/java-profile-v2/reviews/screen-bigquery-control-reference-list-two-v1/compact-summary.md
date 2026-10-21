PASS independent review; one cohort, five measured pairs; NONE_PROMOTED. No confidence interval or established speed gain.
All 134 requests matched all six oracle digests; all driver exits zero; no new OOM, resource/provenance/pressure-proxy or cleanup violation; PSI unavailable.
Telemetry anomalies admitted under the existing warmup-only policy: 1. Approximate lifetime HWM includes warmup/JIT; current RSS is a postrequest snapshot; allocation is VM-wide request delta, not retained heap.

| Metric | CONTROL | reference-list-two | Candidate / CONTROL |
| --- | ---: | ---: | ---: |
| Elapsed median s | 4.114438774 | 4.267463272 | 1.037192071 |
| Elapsed mean s | 4.196352601 | 4.320051266 | 1.029477662 |
| CPU median s | 4.900000000 | 4.960000000 | 1.012244898 |
| CPU mean s | 4.974000000 | 5.020000000 | 1.009248090 |
| Allocated median bytes | 2079609448.000000000 | 2082391336.000000000 | 1.001337697 |
| Current RSS median MiB | 608.640625000 | 605.652343750 | 0.995090237 |
| Approx. lifetime HWM MiB | 1300.667968750 | 1468.492187500 | 1.129029255 |

CONTROL: 62 matched warmups; drift -2.686385%, warm JIT/wall 1.998121%, spread 1.022550965.
Measured JIT deltas [181, 93, 18, 72, 65] ms; sum 429 ms / 20.981763006 s = 2.044633%; MEASURED_JIT_ABOVE_2_PERCENT; invalid intervals 0.

reference-list-two: 62 matched warmups; drift 4.446247%, warm JIT/wall 1.776481%, spread 1.089741568.
Measured JIT deltas [163, 54, 23, 46, 24] ms; sum 310 ms / 21.600256332 s = 1.435168%; MEASURED_JIT_AT_OR_BELOW_2_PERCENT; invalid intervals 0.

One-cohort CONTROL screen only. Elapsed, CPU and allocated-byte medians and approximate lifetime HWM increased. CompilationMXBean reports approximate accumulated compilation elapsed time, potentially summed across threads, not CPU%. Container pauses and later counter updates may affect that proxy; no causal attribution is inferred. No independent-cohort interval or speed gain is established; none promoted. All rows and measured JIT diagnostics are retained, with no postmeasurement rejection.

The independent reviewer checked manifest/runner pins, matched source-family/compiler flags, loaded class locations/digests, the exact compiler command and all six recorded outputs on every request. Owned JVM exits and OOM deltas were zero, and both containers were removed. No full-JAR hashing or JVM/build execution occurred in this review. Source/probe correctness review is separate.

Evidence: review.json, rows.csv, cohorts.csv, ratios.csv, measured-jit.csv and compact-summary.json. Original results are unchanged.

| Existing-policy telemetry exception | Prior KiB | Later KiB | Delta KiB | Retained superseding sample |
| --- | ---: | ---: | ---: | --- |
| reference-list-two, warmup 8 to 9 | 1209936 | 1209912 | -24 | warmup 29: 1282072 KiB |

This is the pre-existing warmup-only VmHWM amendment, originally adopted after an earlier v1 observation. Exact retained rows prove the later strictly larger sample. No new exception, tolerance, smoothing or row/cohort removal was introduced; raw values and maximum-sample RSS ratios remain unchanged. No measured VmHWM or cgroup-peak decrease was admitted.
