# Shared-cohort factor screen

`paired_factors.py` keeps the pinned Java comparator and complete task unchanged.
Each cohort starts one JVM, performs adaptive stock-compiler warmup, and measures
all Go configurations against that JVM. Every measured Go invocation has its own
fresh stock Java compilation; Java results are not shared between variants. Java
constructs a fresh compiler and parses ordinary externs for every request.

Provide a JSON array of objects with `name`, absolute `binary`, and `go_env`.
Optional `sha256` verifies the expected binary. Only `JSGO_EXPERIMENT_*` strings
may vary in the environment. Configurations cannot override arguments, inputs,
externs, chunks, locale, CPU/GC settings, Java flags or output requirements.
Configure stock main, factor off, each audited individual factor, and the compound
explicitly; the runner does not invent flag names or assume they activate.
Separate dispatch evidence is still required to attribute a factor's effect.

```sh
python3 tools/bench/paired_factors.py --spec /abs/frozen-task.json \
  --variants /abs/variants.json --control factor-off \
  --jar /abs/pinned.jar --driver /abs/compiled-driver \
  --java-home /abs/installed-jdk --image 'image@sha256:recorded-digest' \
  --out /abs/new-evidence --cohorts 3 --repetitions 5
```

There are five pairs **per variant per cohort** by default. Warmup is 12 Java
requests, extending to at most 30 under the existing declared stability policy,
and three byte-gated unmeasured fresh Go invocations per variant. Raw warmup,
wall/CPU/GC/JIT samples, every digest, process RSS spread, failure states, and row
mapping remain available. The controller records measurement drift too. A JVM
that fails to settle is marked inconclusive; failure/OOM is a failed experiment.

Variant order uses seeded randomized Latin rotations. Across each complete block
of V repetitions, every variant visits every position. Partial-block position
counts differ by at most one. Each variant's Java/Go first direction alternates
across repetitions and cohorts. `schedule.json` records the plan before running;
`row-mapping.json` adds exact protocol request IDs after adaptive warmup.

Every compiler uses the existing kernel-enforced 4 GiB / zero swap / two-CPU /
affinity 0,1 budget. Java remains at `-Xmx3g`; Go runtime defaults remain unchanged.
All JVM threads are paused during every Go process. Paused resident memory remains
charged and outer OOM invalidates the cohort. A common final JVM lifetime VmHWM
is used for every variant's memory ratio, avoiding order-dependent Java high-water
marks. Each Go row retains its fresh-process wait4 peak. These are different
memory windows, not per-compilation live-heap measurements.

Elapsed and memory confidence intervals resample independent restarted JVM
**cohort median ratios**, never correlated requests within a JVM. Direct Go
control/variant ratios use the same cohort pairing. Individual samples and sample
standard deviations remain in raw/cohort summaries. Three cohorts provide limited
precision; the script preserves uncertainty and does not promise a win threshold.

The script pins helper source hashes and copies exact helper bytes plus driver
class into the new evidence directory before use. Container helpers execute from
that read-only snapshot. Its Python import compiles checked bytes directly,
avoiding stale `.pyc` caches. The source manifest records originals, snapshots and
SHA-256. Source inputs, binary identities, stock jar and external JDK configuration
are checked between cohorts; the same complete immutable spec applies throughout.
No compiler build, image pull, fixture generation, worker or extern cache is used.

`paired_factors_test.py` exercises the schedule, row/statistics model, immutable
helper import and strict configuration/resource admission without launching a
compiler, JVM or Docker container. Actual protocol/performance validation is a
separate coordinated run.
