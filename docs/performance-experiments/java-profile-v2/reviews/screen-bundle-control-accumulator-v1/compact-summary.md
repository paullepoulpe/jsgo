PASS independent review; one cohort, five measured pairs; NONE_PROMOTED. No confidence interval or established speed gain.
All 54 requests matched all six oracle digests; all driver exits zero; no new OOM, resource/provenance/pressure-proxy or cleanup violation; PSI unavailable.
Telemetry anomalies admitted under the existing warmup-only policy: 0. Approximate lifetime HWM includes warmup/JIT; current RSS is a postrequest snapshot; allocation is VM-wide request delta, not retained heap.

| Metric | CONTROL | accumulator | Candidate / CONTROL |
| --- | ---: | ---: | ---: |
| Elapsed median s | 11.158163928 | 10.728804040 | 0.961520561 |
| Elapsed mean s | 11.043794191 | 10.860059846 | 0.983363114 |
| CPU median s | 12.880000000 | 12.360000000 | 0.959627329 |
| CPU mean s | 12.832000000 | 12.548000000 | 0.977867830 |
| Allocated median bytes | 4684739288.000000000 | 4367995056.000000000 | 0.932388077 |
| Current RSS median MiB | 633.699218750 | 637.539062500 | 1.006059411 |
| Approx. lifetime HWM MiB | 1303.429687500 | 1371.406250000 | 1.052152075 |

CONTROL: 22 matched warmups; drift 0.405565%, warm JIT/wall 1.571590%, spread 1.076015582.
Measured JIT deltas [421, 212, 163, 138, 110] ms; sum 1044 ms / 55.218970957 s = 1.890655%; MEASURED_JIT_AT_OR_BELOW_2_PERCENT; invalid intervals 0.

accumulator: 22 matched warmups; drift -1.197264%, warm JIT/wall 1.696773%, spread 1.044408491.
Measured JIT deltas [282, 145, 206, 148, 233] ms; sum 1014 ms / 54.300299231 s = 1.867393%; MEASURED_JIT_AT_OR_BELOW_2_PERCENT; invalid intervals 0.

One-cohort bundle CONTROL screen only. Median elapsed fell 3.8479%, but mean elapsed fell only 1.6637%; CPU and allocated-byte medians fell while approximate lifetime HWM increased. CompilationMXBean reports approximate accumulated compilation elapsed time, potentially summed across threads, not CPU%. Container pauses and later counter updates may affect that proxy; no causal attribution is inferred. No independent-cohort interval or confirmed speed gain is established; none promoted. Selected for the sole v4 bundle accumulator direct-STOCK holdout; this screen is excluded from its inference, and no holdout rows or results were considered here. All rows and measured JIT diagnostics are retained, with no postmeasurement rejection.

The independent reviewer checked manifest/runner pins, matched source-family/compiler flags, loaded class locations/digests, the exact compiler command and all six recorded outputs on every request. Owned JVM exits and OOM deltas were zero, and both containers were removed. No full-JAR hashing or JVM/build execution occurred in this review. Source/probe correctness review is separate.

Evidence: review.json, rows.csv, cohorts.csv, ratios.csv, measured-jit.csv and compact-summary.json. Original results are unchanged.
