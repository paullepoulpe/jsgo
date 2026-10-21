PASS independent review; one cohort, five measured pairs; NONE_PROMOTED. No confidence interval or established speed gain.
All 92 requests matched all six oracle digests; all driver exits zero; no new OOM, resource/provenance/pressure-proxy or cleanup violation; PSI unavailable.
No telemetry anomalies. Approximate lifetime HWM includes warmup/JIT; current RSS is a postrequest snapshot; allocation is VM-wide request delta, not retained heap.

| Metric | CONTROL | source-position-last | Candidate / CONTROL |
| --- | ---: | ---: | ---: |
| Elapsed median s | 4.192934694 | 4.228053542 | 1.008375720 |
| Elapsed mean s | 4.133574304 | 4.203515707 | 1.016920321 |
| CPU median s | 4.800000000 | 4.840000000 | 1.008333333 |
| CPU mean s | 4.790000000 | 4.910000000 | 1.025052192 |
| Allocated median bytes | 2097608688.000000000 | 2067180904.000000000 | 0.985494061 |
| Current RSS median MiB | 957.703125000 | 626.097656250 | 0.653749205 |
| Approx. lifetime HWM MiB | 1454.605468750 | 1186.257812500 | 0.815518598 |

CONTROL: 41 matched warmups; drift -1.295726%, warm JIT/wall 1.679998%, spread 1.049530816.
Measured JIT deltas [172, 64, 20, 30, 8554] ms; sum 8840 ms / 20.667871520 s = 42.771700%; MEASURED_JIT_ABOVE_2_PERCENT; invalid intervals 0.

source-position-last: 41 matched warmups; drift -0.315955%, warm JIT/wall 1.422669%, spread 1.070350859.
Measured JIT deltas [132, 81, 23, 54, 100] ms; sum 390 ms / 21.017578533 s = 1.855590%; MEASURED_JIT_AT_OR_BELOW_2_PERCENT; invalid intervals 0.

One-cohort CONTROL screen only. Elapsed and CPU medians increased while allocated bytes fell 1.4506%. CONTROL recorded an 8554 ms cumulative compilation-time counter delta on its fifth measured request, putting its batch JIT fraction at 42.7717%, versus 1.8556% for the candidate. CompilationMXBean measures approximate accumulated compilation elapsed time, potentially summed across threads, not CPU%. A compilation spanning a paused container may contribute paused duration when a later counter update occurs; that mechanism was not diagnosed for this jump. Treat the large ratio as a potentially contaminated elapsed-time proxy, not proof of JIT CPU share or slowdown causality. This single cohort does not establish a settled benefit; lower allocation and approximate RSS remain descriptive. No independent-cohort interval or speed gain is established; none promoted. All rows and measured JIT diagnostics are retained, with no postmeasurement rejection.

The independent reviewer checked manifest/runner pins, matched source-family/compiler flags, loaded class locations/digests, the exact compiler command and all six recorded outputs on every request. Owned JVM exits and OOM deltas were zero, and both containers were removed. No full-JAR hashing or JVM/build execution occurred in this review. Source/probe correctness review is separate.

Evidence: review.json, rows.csv, cohorts.csv, ratios.csv, measured-jit.csv and compact-summary.json. Original results are unchanged.

| Source identity (SHA-256) | CONTROL | source-position-last |
| --- | --- | --- |
| original_sha256 | 3dfbc67cea91c7f39700865cac659622089f22297f2af2616e79d5fac82f30c0 | 3dfbc67cea91c7f39700865cac659622089f22297f2af2616e79d5fac82f30c0 |
| effective_sha256 | 3dfbc67cea91c7f39700865cac659622089f22297f2af2616e79d5fac82f30c0 | b88bbd45a5e8ad51625efbfac70ac697dd38da584f9d663e8e1d23e713e727b5 |
| relocated_sha256 | 69e158c0f08962220b95394439fb8e77d4e66a0462a43e5615b95f04516ce865 | ca325c90edbaec2fa7ced875754763da4061c560819ea23409f78d5aeb2eda6d |

Small source/receipt identity checks passed against the matching manifest. The rebuilt family is LineNumberScanner only, with matching CONTROL/candidate compiler flags and independent loaded-origin checks. The change keeps one scanner-local immutable SourcePosition for an identical file/offset/line/column result; no cross-request compiler cache is introduced. Semantic/probe review is separate from this measurement review.
