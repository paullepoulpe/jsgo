Future confirmation protocol, recorded 2026-10-21 08:50:50 UTC, before any v4 result. Source preparation only; root selects the candidate/workload and owns execution. This is a prospective response to residual measured-phase JIT activity and wide cohort uncertainty in v3. It does not change v3 sources, rows, historical PASS statuses, or the interpretation of unfavorable results.

| Item | Fixed v4 rule |
| --- | --- |
| Comparison | Selected candidate directly against unchanged published STOCK; screens are separate |
| Independent units | Five fresh independent pairs of JVMs, one workload per run |
| Warmup | Matched paired warmups; minimum 60, maximum 150 requests per variant |
| Settling window | Last eight requests versus preceding eight, after at least 60 |
| Absolute median-wall drift | <=3% |
| Last-eight summed JIT / summed wall | <=1%, with valid nondecreasing JIT counters |
| Last-eight maximum/minimum elapsed | <=1.15 |
| Stopping warmup | First paired prefix where both variants meet every condition; if none by 150, preserve the incomplete failure |
| Measurements | Exactly 20 alternating pairs per cohort, fixed before measurement; no adaptive extension or early success stop |
| Resource/work checks | Unchanged v3 4 GiB/no swap/two CPUs, default collector/-Xmx3g, fresh compiler/extern parsing, six digests on every request, class provenance, pressure/headroom/disk/timeout guards and owned teardown |
| Primary estimator | Geometric mean of five within-cohort candidate/STOCK median-wall ratios |
| Primary interval | Two-sided 95% Student log-t, df=4, critical value 2.7764451051977987 |
| Sensitivities | All cohort effects; means and paired ratios; whole-cohort bootstrap using 10,000 resamples; CPU/allocation/approximate RSS descriptive |

Retain every warmup and measured row, all ordered JIT before/after counters and deltas, and all failed attempts. Do not pool screens, v3 runs, individual requests across JVMs, or workloads into the five-cohort inference. The 200 measured requests are not 200 independent JVM comparisons. More requests reduce within-JVM sampling variability; five independent cohorts still cannot promise detection of a 2% effect or establish the normal-log-effect assumption. A bootstrap based on five observed cohorts remains limited; it is not the sole promotion criterion.

Measured JIT is diagnostic, not a new rejection or exclusion rule. Report each whole 20-request batch and four fixed consecutive five-request blocks: every delta, sum JIT, sum wall, fraction, strict >1% flag and strict >2% flag. The >2% flag preserves comparability with the prior declared review; >1% adds the tighter warmup target. Unavailable/decreasing counters are explicit, not zero. Also show first-ten/last-ten median-wall drift as a diagnostic without any filtering rule. A later JIT wave remains possible after any warmup gate; inspect the ordered evidence before describing a difference as a settled gain. No row/cohort replacement, result-dependent retry, threshold relaxation, or favorable-subset analysis is authorized.

The unchanged six-output/resource/provenance/cleanup conditions still determine protocol admission. Measured JIT annotations do not turn a prior protocol PASS into a failure. The existing warmup-only VmHWM exception (later strictly exceeded, exact retained proof) carries forward prospectively; measured VmHWM and any cgroup-peak decrease remain fatal. Approximate VmHWM is not exact physical peak memory, and unavailable PSI is not zero.

Before launching, root records the selected candidate, workload, source/build/spec/runner identities, seed and timeout bounds in the normal runtime coordination receipt; this is not a new approval workflow. The default seed stays 20261020, request timeout 300 seconds and total timeout 14,400 seconds. Root may predeclare a different total timeout in that receipt before starting, without changing samples/thresholds after seeing results. The executable protocol itself fixes STOCK, five cohorts, 20 pairs and 60–150 warmups. Run v4 into new result/review directories. The isolated launch-v4.py supplies this receipt, defaults to BigQuery, and never automatically launches bundle. Existing v3 launcher/reviewer behavior stays unchanged; the publisher dispatches separately on the recorded v3/v4 protocol.

This note fixes the sampling protocol, not an open-ended search family. For a single selected comparison, a median-wall log-t interval wholly below 1 is the proposed statistical evidence of lower elapsed cost, subject to work parity, practical size and residual-JIT interpretation. If several v4 candidates/workloads will support one selection claim, declare the comparison family and multiplicity treatment before their results; these per-comparison 95% intervals alone do not establish familywise control. Preserve all attempts and report decisions to stop pursuing a candidate, including unfavorable outcomes. Bundle generalization is a separate workload, not extra BigQuery cohorts.

Cost: five cohorts × two JVMs × (60–150 warm +20 measured) means 800–1,700 compilations. At roughly 4.2 s/request this is 56–119 minutes of BigQuery compile time plus startup/artifact/management overhead; around 11 s/request gives 2.4–5.2 hours for bundle. These are planning estimates, not performance results. The existing four-hour default may time out a difficult bundle run; that partial outcome remains evidence. Use this more expensive protocol for a plausible selected candidate, not as a way to rerun until a screen looks favorable.
