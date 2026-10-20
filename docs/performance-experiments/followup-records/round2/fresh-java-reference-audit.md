# Fresh stock warm-Java reference audit

Prepared 2026-10-19; **no JVM, compiler build, or benchmark was run for this audit**. The ready command is in `fresh-java-n1-command.json`; its warm-only wrapper is `run_fresh_java_n1.py`. Expected duration is 12–15 minutes for all five original tasks, 12 warmups plus 10 measured iterations each. Run only in a coordinator-reserved quiet window. The output directory must not already exist.

## Equivalent work

The source audit uses `jsgo-structural-bitset-template/tools/bench/{bench.py,java/CliBatch.java}` and the pinned Apache Closure source. `CliBatch` creates a new stock `CommandLineRunner` for every job on every iteration; it overrides exit handling, not compiler construction or extern handling. Stock `AbstractCommandLineRunner.doRun` calls `createCompiler`, `createOptions`, and `createExterns` each time. Stock `CommandLineRunner.createCompiler` constructs a new `Compiler`; `createExterns` loads the normal builtin externs each time. There are no Java class overlays, parsed-extern caches, compiler-result caches, or worker helpers in this classpath.

The pinned jar SHA-256 is `d6864ebf70bca5ba6e237898de6611a389aaf8b933c2549036ced4a684b664b6` (v20261015). The wrapper verifies the installed jar against both the pin and the stored baseline metadata. It validates the stored cold oracle's complete artifact digests before starting. Each new JVM's final measured output must match all oracle stdout, stderr, exit, and task artifacts byte for byte, including missing-file state.

The old baseline omitted the thread flag; the pinned compiler's default is exactly one. The prepared command explicitly replaces any existing flag with one `--num_parallel_threads=1` on every job. All other arguments come from the stored baseline, including ADVANCED/checkTypes for typed-app, checks_only/checkTypes for checks-only, SIMPLE for bundle/downlevel-simple, and WHITESPACE_ONLY for downlevel-whitespace.

The helper compiles the unchanged benchmark driver into a new output-local directory before timing, checks javac success, and records source/class hashes and classpath. This avoids the general runner's mtime-only helper cache. Each task receives a fresh JVM. The measured iteration interval includes fresh compiler construction, normal extern parsing, compilation, and output handling; the JVM's classes and JIT remain warm. JVM injection environment variables are removed by the existing runner. No custom Java tuning flags are added.

## Is twelve warmups enough?

The old reference was captured on 2026-10-18 around 22:08 UTC on the same recorded host, with JDK 21.0.12.1 and the same pinned jar. Its metadata does not explicitly record affinity, so identical old affinity cannot be established retrospectively. The new command pins CPU affinity to 0,1 and records actual affinity, `cpu.max`, cpuset and memory limit. Expected quota is two CPUs (`200000 100000`).

All old warm runs had 12 warmups, 10 measured iterations, zero process exits and no cold-oracle mismatches. Twelve iterations remove the large initial ramp, but neither these measurements nor elapsed-time flatness alone proves that JIT compilation has finished. The runner does not record compilation-MXBean activity.

| Task | Warm first/last 3 median (s) | Measured first/last 5 median (s) | Wall drift | CPU drift | Measured range (s) |
|---|---:|---:|---:|---:|---:|
| typed-app | 5.792 / 3.911 | 3.713 / 3.778 | +1.8% | −6.3% | 3.480–18.427 |
| bundle | 18.007 / 17.068 | 16.485 / 16.759 | +1.7% | +6.0% | 16.402–17.202 |
| downlevel-whitespace | 2.060 / 0.995 | 0.907 / 0.947 | +4.4% | −8.9% | 0.824–1.018 |
| downlevel-simple | 4.643 / 2.840 | 2.761 / 2.814 | +1.9% | −2.0% | 2.679–2.874 |
| checks-only | 8.379 / 5.993 | 6.068 / 6.008 | −1.0% | −3.7% | 5.897–6.270 |

Drift is second-half/first-half median minus one, including process CPU from JIT/GC threads. The typed run's last two measured wall times were 6.098 and 18.427 s, versus process CPU 4.740 and 5.110 s. This is substantial wall-time variation; these records do not establish its cause. Retain the outliers and raw GC/CPU/wall rows. Medians resist them, but the old cohort should not be represented as a demonstrated quiet, fully settled reference.

Use the unchanged 12+10 protocol for the next fresh comparison, report half-drift and full ranges, and do not adapt warmup or drop iterations after seeing which runtime wins. If a final result needs tighter evidence, predeclare a longer fixed warmup (for example 20+10, approximately 16–20 minutes) or run independent paired cohorts using `cohort_compare.py`. Confidence intervals belong over independent restarted cohorts, not the ten correlated iterations inside one JVM. Longer warmup alone does not resolve host noise.

## Memory and comparison limits

Java RSS is the fresh JVM's **lifetime peak**, including startup, JIT, warmup and measured work. Go CLI RSS is the fresh Go process's lifetime peak. These are useful process-residency comparisons with different lifetime windows, not live heap per compilation or incremental request memory. Keep the window explicit for every Java/Go memory ratio. A new Java-only cohort does not create a paired confidence interval with older Go timings; final confirmation should rerun both in alternating order under the same quota and affinity.

The previous five JVM lifetimes totalled about 739 seconds (12.3 minutes), giving the expected duration above. Javac happens once beforehand and is excluded from compiler timings. The wrapper is source-prepared and syntax-parsed only; it has not been executed or used to replace any reference.
