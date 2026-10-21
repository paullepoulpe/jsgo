# Profile-guided Java optimization experiments

This batch searches for the smallest worthwhile code optimization in stock Java
Closure v20261015. It is separate from the frozen `transfer-v1` study; no historical
result is replaced. No Go performance claim follows from these Java-versus-Java
experiments.

The strongest small candidate is the [shared parameter accumulator](upstream-patches/accumulator.patch):
one invocation-local set replaces recursive temporary sets and their merges,
while preserving every node visit, predicate and first-insertion order. In the
completed CONTROL screens it reduced allocated bytes by 2.76% on generated
BigQuery code and 6.76% on the eleven-library bundle (about 302 MiB per bundle
compile). The complete upstream FunctionArgumentInjectorTest class passed all
91 tests on STOCK, rebuilt CONTROL and candidate.

**A repeatable speed gain is not established.** The bundle screen observed 3.85%
lower median elapsed time but only 1.66% lower mean elapsed time, with 5.22% higher
approximate lifetime peak RSS. BigQuery did not show an elapsed-time gain. The
sole five-cohort direct-STOCK confirmation stopped during cohort 2 because the
parent host recorded measured memory pressure and direct reclaim. It remains
incomplete; no smaller confirmation or five-cohort confidence interval is claimed.
The small patch is an allocation-saving candidate, not a demonstrated general
speed improvement or a proven RSS reduction.

The separate completed three-cohort identifier-dispatch confirmation also did
not establish a speed gain: elapsed ratio 0.98876 with a 95% log-cohort t interval
of 0.88116–1.10950, and higher approximate lifetime peak RSS in all three cohorts.
Every negative result, failed attempt and remaining idea is retained below and
in the [open-ideas record](OPEN-IDEAS.md).

## Work and resource boundaries

Each request constructs a fresh CommandLineRunner/compiler and parses normal
externs. Inputs, flags, diagnostics, generated JavaScript, source maps, input
manifests and exit status must match the pinned stock run. Neither parsed externs
nor compilation results are cached. Runtime flags and the collector are unchanged.
Performance tests run without JFR, alternate the two implementations in the same
session, and enforce 4 GiB memory / no swap / two CPUs per JVM. Only one compiler
runs at a time; the other JVM is paused. Warmup requests are checked for parity as
well as measured requests.

The two workloads are a generated BigQuery/protobuf ADVANCED build (1,270 supplied
files, 251 retained inputs) and an ADVANCED bundle of eleven libraries. These do not
establish behavior for every Closure workload or a full enterprise chunk graph.

## Measurement

The historical `selection-plan.json`, recorded before candidate performance
results, used exploratory one-cohort screens and direct published-STOCK
confirmation with three fresh independent paired-JVM cohorts under v3. The
completed identifier-dispatch confirmation below follows that protocol.

The v4 confirmation was separately prespecified in the frozen
[v4 note](CONFIRMATION-V4-PRESPECIFIED.md): five fresh direct-STOCK cohorts,
20 fixed measured pairs per cohort, and 60–150 matched warmups per variant.
The final eight warmups must have at most 3% median drift versus the prior eight,
compilation elapsed-time proxy at most 1% of wall time, and max/min elapsed at
most 1.15. The sole selected accumulator/bundle holdout launched under the
[pre-execution selection note](SELECTED-V4-CONFIRMATION.md), with a predeclared
21,600-second total timeout and unchanged 300-second request timeout. It stopped
at its pressure guard during cohort 2. Its planned primary interval was the
log-cohort t interval with four degrees of freedom; that inference cannot be
computed from this incomplete run. All retained rows, failure records and
predeclared JIT annotations remain, without adaptive exclusion or replacement.

Per-request samples from one JVM are not independent cold starts. Results report
elapsed time, CPU, VM-wide allocated bytes, current RSS, and approximate lifetime
peak RSS including warmup/JIT. Allocation is not retained memory. Few-cohort
bootstrap intervals are limited; individual cohort effects and log-cohort t
intervals matter, and five cohorts do not promise detection of a 2% gain.
The JIT counter reports approximate accumulated compilation elapsed time,
potentially summed across threads and affected by container pauses or later
updates; it is not CPU%. See [counter limitations](JIT-COUNTER-INTERPRETATION.md).
Historical v3 admission outcomes and prospective v4 thresholds are unchanged.

`PROFILE-NOTES.md` documents diagnostic captures and their limitations. The first
capture was disturbed by enabling JFR after warmup; the second enabled it before
warmup but narrowly missed the strict final JIT threshold. Neither capture is a
performance comparison or a settled benchmark. All evidence is retained.

## Sources and validation

Each `candidates` directory preserves the proposed patch, pinned original source,
and focused probes where prepared. `build/<candidate>/manifest.json` identifies
same-family original-source CONTROL and candidate overlays; dependencies are
mechanically relocated to the published jar's packaging. Compiler flags are the
same on both sides. Performance runs check loaded class provenance and every output.
The complete pinned upstream `FunctionArgumentInjectorTest` class passed on
STOCK, CONTROL, and the shared parameter accumulator: 91 tests per variant, with
zero failures, ignores, or assumption failures. All test bodies and four upstream
helpers are retained; four compiler-facing Guava imports are mechanically relocated
for the published jar. Independent source and receipt reviews verified exact test
names, class origins, and CONTROL/candidate class digests. See
`validation/function-argument-injector-upstream-v1/COMPLETED-VALIDATION.txt` and its
linked results. This validates that specific upstream class; the full upstream
Java suite has not been run. Focused probes and workload output parity remain
separate checks.

Raw results, launch receipts, source/build identities, failure logs and teardown
proofs are retained in `results` and `launches`. Prepared but unexecuted candidates
are not claimed as improvements. Results below are updated only after completed runs and independent review.

## Incomplete accumulator confirmation

The frozen v4 controller rejected candidate request 68 in cohort 2 (the second
request of measured pair 8) after the parent cgroup recorded 387 new memory-limit
events and 14,879 pages of direct scan/reclaim. Cohort 1 completed 60 warmups per
variant and all 20 measured pairs. Cohort 2 completed 60 warmups per variant and
seven measured pairs, then the unpaired STOCK request 68. There are 295 retained
full metric/digest rows: 240 warmups and 55 measurements. Each matches all six
oracle digests. The failed candidate request has stdout/stderr/exit and pressure
evidence, but the pressure assertion ran before metric-row/digest persistence;
six-output parity is not established for that request.

The official v4 reviewer correctly rejects the attempt because failure evidence
is present. Four owned JVMs exited with status zero, all containers were removed,
no compiler or driver remained, and there were no new OOM events. The independent
[failure audit](reviews/confirm-bundle-stock-accumulator-v4-attempt-1/failure-audit.md)
retains the exact failure, row inventory and cleanup evidence. Partial descriptive
rows are not a successful confirmation, and no completed cohort is substituted
for the planned five-cohort comparison.

The retained complete pairs give the following **descriptive** candidate/STOCK
ratios. They support an allocation reduction but no elapsed-time gain. Cohort 2
is truncated; its HWM endpoints are incomplete and unequal. The extra STOCK row,
all JIT deltas and every available metric remain in the linked audit/CSV files.
No confidence interval or pooled estimate is calculated from these partial data.

| Retained cohort | Complete pairs | Elapsed median ratio | Elapsed mean ratio | CPU median ratio | Allocated-byte ratio | Approx. lifetime peak RSS ratio |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| 1, complete cohort of failed run | 20 | 1.01391 | 1.01414 | 0.99838 | 0.92295 | 0.95553 |
| 2, partial cohort of failed run | 7 | 0.99716 | 1.00220 | 1.00960 | 0.93732 | 0.95660 |

After shutdown, the 8 GiB parent still held about 5.64 GB, mostly file cache and
reclaimable kernel slab. Bounded task-file cache advice did not recover reliable
headroom for a clean restart. The no-restart decision was made before root saw
partial performance ratios; no JVM settings or benchmark admission rules changed.
See the [host investigation](environment-preflight/README.md). A future complete
confirmation needs sufficient host headroom as well as the original 4 GiB
per-JVM cap. The failed run remains evidence, not a speed claim.

## Completed identifier-dispatch confirmation

The identifier dispatch guard (`traversal-name-switch`) passed independent review
of three fresh BigQuery paired-JVM cohorts, with five measured requests per variant
after 41, 41 and 40 matched warmups. Screen observations were not pooled. Ratios
are candidate / published STOCK; below 1 means lower cost.

| Cohort | STOCK / candidate elapsed median s | Elapsed ratio | CPU ratio | Allocated-byte ratio | Approx. lifetime peak RSS ratio |
| --- | --- | ---: | ---: | ---: | ---: |
| 1 | 4.177405 / 4.282162 | 1.02508 | 1.00606 | 0.99271 | 1.04974 |
| 2 | 4.514256 / 4.236009 | 0.93836 | 0.97595 | 1.00591 | 1.05212 |
| 3 | 4.141000 / 4.161506 | 1.00495 | 1.02731 | 0.99593 | 1.06236 |

| Metric | Geometric mean ratio | 95% log-cohort t interval, df=2 |
| --- | ---: | --- |
| Elapsed median, primary | 0.98876 | 0.88116–1.10950 |
| CPU median, descriptive | 1.00289 | 0.94065–1.06924 |
| Allocated bytes, descriptive | 0.99817 | 0.98124–1.01539 |
| Approx. lifetime peak RSS, descriptive | 1.05472 | 1.03822–1.07148 |

Two of three elapsed ratios were above 1, and the interval includes useful gains
and regressions. Mean-wall sensitivity also includes 1: ratio 0.99214, interval
0.90951–1.08227. The three independent cohorts are the inference units; the 30
measured requests are not 30 independent JVM comparisons. Log-t assumes approximately
normal independent log effects, which three observations cannot establish. The
discrete whole-cohort bootstrap is reported in the detailed review and is not the
sole basis for inference.

Measured JIT fractions for STOCK/candidate were **1.885%/2.115%, 2.498%/2.759%,
and 2.227%/2.077%**. Five of six batches exceeded the predeclared 2% annotation;
a settled speed gain is not established. Every JIT delta is retained in the
[compact review](reviews/confirm-bigquery-stock-traversal-name-switch-v1/compact-summary.md)
and [JIT CSV](reviews/confirm-bigquery-stock-traversal-name-switch-v1/measured-jit.csv).
These flags preserve the original PASS; no row or cohort was discarded or retried.

Approximate lifetime peak RSS increased 4.97–6.24% in all three cohorts, a geometric
mean increase of 5.47%. This is unfavorable approximate secondary telemetry, including
warmup/JIT, rather than an exact physical-memory effect or adjusted confirmatory
memory claim. Current RSS median ratios were 0.96906, 1.12298 and 0.64626; these
postrequest snapshots have a different window. CPU and allocation did not establish
an improvement either.

All 244 warm and 30 measured requests matched six oracle digests; resource, pressure,
loaded-origin, stopped and zero-new-OOM review passed. The existing reviewer admitted
one documented warmup-only VmHWM dip: STOCK cohort 1 fell 168 KiB at request 7, then
was strictly exceeded at request 8. Raw values and the retained maximum were unchanged.
No measured dip or new rule was admitted. See the
[independent review](reviews/confirm-bigquery-stock-traversal-name-switch-v1/review.json)
for the exact anomaly and telemetry limits. No bundle confirmation of the
identifier-dispatch candidate was performed; the separate bundle screen below
evaluates the accumulator.

## Completed screens

Ratios are candidate / same-family CONTROL; below 1 means lower cost. These are
single independent paired-JVM cohorts, not confirmation or significance tests.

| BigQuery candidate | Time | CPU | Allocated bytes | Approx. lifetime peak RSS | Current RSS median |
| --- | ---: | ---: | ---: | ---: | ---: |
| Shared parameter accumulator | 1.01779 | 1.00000 | 0.97240 | 0.98805 | 1.10643 |
| Identifier dispatch guard | 0.97959 | 0.98167 | 0.99264 | 0.95423 | 0.91981 |
| Guaranteed-empty callback bypass | 0.98606 | 1.02231 | 1.00132 | 1.26515 | See matching raw review |
| Leaf identifier recursion bypass | 0.99445 | 0.99793 | 1.01016 | 1.05775 | 0.95908 |
| Prototype/scalar equality shortcut | 1.01335 | 0.96654 | 0.99369 | 0.99754 | 1.13202 |
| Empty-parameter shortcut | 1.00788 | 1.04873 | 1.00027 | 0.92595 | 1.10270 |
| Last source-position reuse | 1.00838 | 1.00833 | 0.98549 | 0.81552 | 0.65375 |
| Reference-list initial capacity two | 1.03719 | 1.01224 | 1.00134 | 1.12903 | 0.99509 |

Accumulator CONTROL/candidate medians were 4.16309/4.23716 seconds and
2,079,177,752/2,021,800,312 allocated bytes. Approximate lifetime peak RSS was
1,334.34/1,318.39 MiB; current RSS medians were 702.00/776.71 MiB. These memory
windows differ: reduced allocated bytes do not guarantee lower live RSS. Both
variants passed all six digests on 54 matched warmups and five measured requests.
Measured JIT fractions were 1.946%/1.783%; neither exceeded the predeclared 2% flag.
The screen reduced allocation 2.76% but did not nominate a speed winner. Full raw
results and independent review are under `results/screen-bigquery-control-accumulator-v1`
and the matching `reviews` directory.

Identifier-dispatch CONTROL/candidate medians were 4.26161/4.17462 seconds, with
approximate lifetime peak RSS 1,303.44/1,243.78 MiB and current RSS medians
666.95/613.47 MiB. Both variants matched all outputs for 45 warmups and five measured
requests. Measured JIT fractions were **3.274%/3.361%**, so both carry the predeclared
flag. The roughly 2.04% elapsed lead is preliminary, not an established settled gain.
This candidate was selected for a fresh three-cohort comparison directly against
published STOCK. The completed confirmation above did not establish a speed gain.

The three-line `identity-differs` patch compiled and passed focused oracle probes on
STOCK/CONTROL/candidate. Its timing screen was intentionally stopped during warmup
after more specific profile attribution became available. Two paired warmups had
completed; there were no measured requests. Most EqualityChecker allocation came
from `equals`, which already has an identity shortcut; the two shown `differsFrom`
prefixes account for only 0.0782% of total sampled allocation. Three unshown allocation
samples carry another 0.0195% total weight, so this is a prioritization observation,
not proof of zero remaining opportunity. `identity-screen-stop.json` preserves the
reason, the original interrupted result remains failed, and teardown confirms both
JVMs exited and containers were removed with no new OOM events. No speed ratio is
reported for this intentionally incomplete screen.

`traversal-noop` completed a one-cohort CONTROL screen: elapsed medians
4.28706/4.22732 seconds, CPU 4.93/5.04 seconds, allocated bytes
2,082,314,704/2,085,062,832, and approximate lifetime peak RSS
1,196.85/1,514.19 MiB. Forty-five matched warmups and all five measured pairs
passed all six outputs. Measured JIT was 3.513%/4.428%, both flagged. The small
elapsed lead alongside increased CPU and peak RSS does not nominate a useful
speed improvement. Full results and independent review are retained under
`screen-bigquery-control-traversal-noop-v1`.

`traversal-leaf-name` completed 53 matched warmups and five measured pairs,
with all outputs matching. CONTROL/candidate elapsed medians were
4.16362/4.14053 seconds; CPU 4.84/4.83 seconds; allocated bytes
2,075,729,072/2,096,816,920; lifetime peak RSS 1,211.84/1,281.83 MiB.
Measured JIT fractions were 1.535%/2.053%, flagging the candidate. This screen
establishes no useful speed improvement; no holdout was nominated.

`equality-prototype-scalar` passed independent review of 40 matched warmups and
five measured pairs, with all 90 requests matching all six oracle digests.
CONTROL/candidate elapsed medians were 4.32734/4.38510 seconds, a 1.335% increase;
allocated bytes were 2,074,036,912/2,060,940,672, a 0.631% decrease. CPU medians
were 5.08/4.91 seconds, but CPU means were nearly equal at 4.992/4.982 seconds
(only 0.200% lower), so the median CPU difference does not establish a repeatable
improvement. Approximate lifetime peak RSS was 1,179.00/1,176.11 MiB; current RSS
medians were 652.35/738.47 MiB and represent a different measurement window.
Measured JIT fractions were 1.794%/1.577%, both below the predeclared 2% flag.
Source/build receipts, loaded origins, pressure proxies and graceful cleanup
passed; PSI was unavailable. This single-cohort screen establishes no speed gain,
has no confidence interval, and leaves the candidate unpromoted. Every measured
JIT delta and metric is retained in the
[compact review](reviews/screen-bigquery-control-equality-prototype-scalar-v1/compact-summary.txt)
and [independent review](reviews/screen-bigquery-control-equality-prototype-scalar-v1/review.json).

`empty-parameters` passed independent review of 41 matched warmups and five
measured pairs, with all 92 requests matching all six oracle digests.
CONTROL/candidate elapsed medians were 4.15796/4.19071 seconds (0.788% higher);
CPU medians were 4.72/4.95 seconds (4.873% higher), and CPU means were
4.776/4.938 seconds (3.392% higher). Allocated-byte medians were
2,085,108,512/2,085,671,224 (0.027% higher). Approximate lifetime peak RSS was
1,295.27/1,199.35 MiB, while current RSS medians were 601.37/663.13 MiB;
these different windows do not establish a repeatable memory benefit.
Ordered measured JIT deltas were `[117, 98, 23, 30, 105]` ms for CONTROL and
`[150, 72, 22, 69, 345]` ms for the candidate. Their batch JIT fractions were
1.795%/3.160%, flagging the candidate under the predeclared 2% annotation.
Manifest/runner pins, loaded class provenance, resource/pressure proxies and
graceful cleanup passed; OOM deltas were zero, PSI was unavailable, and no
telemetry anomaly was admitted. Every row remains, with no postmeasurement
rejection. This one-cohort screen has no confidence interval, establishes no
speed gain, and leaves the candidate unpromoted. See the
[compact review](reviews/screen-bigquery-control-empty-parameters-v1/compact-summary.md)
and [independent review](reviews/screen-bigquery-control-empty-parameters-v1/review.json).

`source-position-last` passed independent review of 41 matched warmups and five
measured pairs, with all 92 requests matching all six oracle digests.
CONTROL/candidate elapsed medians were 4.19293/4.22805 seconds (0.838% higher);
CPU medians were 4.80/4.84 seconds, and CPU means were 4.790/4.910 seconds
(2.505% higher). Allocated-byte medians were 2,097,608,688/2,067,180,904
(1.451% lower). Approximate lifetime peak RSS was 1,454.61/1,186.26 MiB;
current RSS medians were 957.70/626.10 MiB.
CONTROL's ordered JIT deltas were `[172, 64, 20, 30, 8554]` ms, while the
candidate's were `[132, 81, 23, 54, 100]` ms. The 8,554 ms counter jump on
CONTROL's fifth measured request puts its batch JIT fraction at **42.772%**,
versus 1.856% for the candidate. This counter reports approximate accumulated
compilation elapsed time, potentially summed across threads, **not CPU%**.
A compilation spanning a container pause may contribute paused elapsed time
to a later counter update; that mechanism is not established for this jump.
Treat the ratio as a potentially contaminated elapsed-time proxy, without
attributing a specific slowdown to JIT. The single cohort establishes no
settled benefit; lower allocation and approximate RSS remain descriptive.
See [counter interpretation](JIT-COUNTER-INTERPRETATION.md).
Source/build receipt identity, loaded origins, resource/pressure proxies and
graceful cleanup passed, with zero new OOM and no telemetry anomalies; PSI was
unavailable. All rows and the original PASS remain. No speed gain is established,
no confidence interval is reported, and the candidate remains unpromoted. See
the [compact review](reviews/screen-bigquery-control-source-position-last-v1/compact-summary.md)
and [independent review](reviews/screen-bigquery-control-source-position-last-v1/review.json).

`reference-list-two` passed independent review of 62 matched warmups and five
measured pairs, with all 134 requests matching all six oracle digests.
CONTROL/candidate elapsed medians were 4.11444/4.26746 seconds (3.719% higher);
CPU medians were 4.90/4.96 seconds, and CPU means were 4.974/5.020 seconds
(0.925% higher). Allocated-byte medians were 2,079,609,448/2,082,391,336
(0.134% higher). Approximate lifetime peak RSS was 1,300.67/1,468.49 MiB
(12.903% higher), while current RSS medians were 608.64/605.65 MiB.
Ordered measured JIT deltas were `[181, 93, 18, 72, 65]` ms for CONTROL and
`[163, 54, 23, 46, 24]` ms for the candidate, giving batch fractions of
2.045%/1.435%; CONTROL carries the predeclared flag. The compilation elapsed-time
proxy has the [pause/accounting limitations](JIT-COUNTER-INTERPRETATION.md)
described above and does not establish JIT CPU share or slowdown causality.
Manifest/runner pins, loaded origins, caps, pressure proxies and graceful cleanup
passed, with zero new OOM; PSI was unavailable. One existing-policy telemetry
exception is retained: candidate VmHWM fell 24 KiB at warmup 8→9
(1,209,936→1,209,912 KiB), then warmup 29 strictly exceeded it at 1,282,072 KiB.
Exact retained rows support that exception; no measured decrease, new rule or
data exclusion was admitted. This single-cohort screen establishes no speed
gain and leaves the candidate unpromoted. See the
[compact review](reviews/screen-bigquery-control-reference-list-two-v1/compact-summary.md)
and [independent review](reviews/screen-bigquery-control-reference-list-two-v1/review.json).

`leaf-children-return`, `reference-function-first-child`, and `graph-edge-copy`
have successful matched builds and focused probes,
but no end-to-end timing result yet. Other prepared sources and rejected
source-level directions remain explicitly labeled in their own directories.
The [prepared-candidate source audit](source-audit/prepared-small-candidates-ranking.txt)
prioritizes the empty-child guard for its small change and simple proof; it makes
no speedup prediction or measured bundle comparison.

## Completed bundle screen

The shared parameter accumulator passed independent review of one fresh
same-family CONTROL cohort after 22 matched warmups and five measured pairs.
All 54 requests matched all six oracle digests. Ratios are candidate / CONTROL;
the screen is separate from both BigQuery and the selected direct-STOCK holdout.

| Metric | CONTROL | accumulator | Ratio |
| --- | ---: | ---: | ---: |
| Elapsed median s | 11.15816 | 10.72880 | 0.96152 |
| Elapsed mean s | 11.04379 | 10.86006 | 0.98336 |
| CPU median s | 12.880 | 12.360 | 0.95963 |
| CPU mean s | 12.832 | 12.548 | 0.97787 |
| Allocated median bytes | 4,684,739,288 | 4,367,995,056 | 0.93239 |
| Approx. lifetime peak RSS MiB | 1303.43 | 1371.41 | 1.05215 |
| Current RSS median MiB | 633.70 | 637.54 | 1.00606 |

The median elapsed lead was 3.848%, while the mean lead was only 1.664%.
These are observations from one persistent-JVM pair, with no independent-cohort
confidence interval or confirmed speed claim. The higher lifetime HWM and current
RSS remain alongside the lower allocation and elapsed observations.
Ordered measured JIT deltas were `[421, 212, 163, 138, 110]` ms for CONTROL and
`[282, 145, 206, 148, 233]` ms for the candidate: 1,044 ms / 55.21897 s = 1.891%
and 1,014 ms / 54.30030 s = 1.867%, respectively. Neither exceeds the historical
2% annotation; neither establishes complete settling. The
[counter limitations](JIT-COUNTER-INTERPRETATION.md) still apply.
Manifest/runner pins, loaded origins, caps, pressure proxies and graceful cleanup
passed, with zero new OOM, driver exit codes zero and no telemetry anomalies; PSI was
unavailable. See the [compact review](reviews/screen-bundle-control-accumulator-v1/compact-summary.md)
and [independent review](reviews/screen-bundle-control-accumulator-v1/review.json).

The selected v4 holdout planned five new direct-STOCK paired-JVM cohorts and 20
fixed measured pairs, with the frozen 60–150 warmups and all existing thresholds.
Its pressure failure and incomplete status are documented above. The screen,
BigQuery observations and partial holdout rows are not pooled into a confirmation.

## Additional diagnostics and prepared directions

A one-compilation instrumented `EqualityChecker.check` diagnostic passed all six
stock output digests and provenance checks. It recorded 1,863,747 entries across
324 class-pair keys; the printed top 128 covered all but 8,374 calls, with zero
counter overflow. Nonidentical InstanceObjectType pairs were 21.65% of entries.
The instrumented run is not a timing or allocation result. Source, probe, exact
counts, original stderr and teardown receipts are retained under
`diagnostics/equality-pairs` and `results/diagnostic-bigquery-equality-pairs-v1`.

The extra allocation-caller analysis in `analysis-driver-v1` uses the existing
corrected diagnostic JFR, with its previously disclosed unsettled timing limits.
It identifies temporary per-node arrays in graph-edge snapshot copying and the
ordinary-token SourcePosition allocation paths as further small candidates.
Prepared sources are proposals, not measured improvements.

Raw JFR recordings remain local. Bounded consumer audits completed but left
ancillary string values, including initial environment/property fields, uncleared
for publication. The public snapshot omits raw `.jfr` files; aggregate reports,
audit receipts, sources, benchmark rows and failed attempts remain included.
No specific secret was identified or exported. See
`publication-preparation/JFR-PUBLICATION-DISPOSITION.json`.
