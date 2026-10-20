# Current-main CLI measurement protocol

This round starts from `da20b44f` and measures ordinary Go CLI compilations
against the pinned, unmodified stock Closure jar on a warmed JVM. Historical
profiles and measurements from the earlier research branch identify candidates;
they do not establish current-main costs or improvements. No parsed-extern
reuse or persistent Go process is part of this comparison.

## Decision criteria

The working target for a **comfortable win** is Go elapsed time at most 80% of
warm Java elapsed time, using identical compiler flags, input bytes, output
requirements and CPU budget. Equivalently, Java/Go elapsed time is at least
1.25x. A broad claim requires this on every task in the original five-task suite:
`typed-app`, `bundle`, `downlevel-whitespace`, `downlevel-simple`, and
`checks-only`. Report task-specific wins and regressions individually. Additional
single-file or held-out workloads supplement that suite; they do not replace or
dilute a slow task in an aggregate.

For a final confidence claim, require the upper 95% interval for Go/Java time to
be at most 0.80 across independent, paired measurement cohorts. A point estimate
or one favorable iteration is insufficient. Treat this criterion as a research
target, not an existing result or a predicted language-level maximum.

## Screening and final comparison

Build all binaries, the Java driver, and profiling tools before timed work.
Run one compilation process at a time on an otherwise quiet machine. Preserve
binary/source/PGO/toolchain hashes, expanded argv, corpus integrity, default
runtime settings and cgroup limits. Every Go run must match stock Java stdout,
stderr, exit, JavaScript and source maps under the established compatibility
rules. Reject failed Java references and incomplete measurement windows.

For screening, `experiments.py` interleaves Go variants by round and checks every
run. Its stored warm-Java median is context, not a newly paired Java confidence
estimate. Use at least ten paired Go rounds when deciding between small changes;
retain raw paired observations and bootstrap the paired log-time ratios with a
fixed seed. Report the interval and sample count alongside the median. Reject a
candidate if any material task regresses without an explicitly justified tradeoff.

For the final Java comparison, restart the stock JVM for each independent cohort
and warm it using fresh compiler instances and normal extern parsing. Predeclare
warm-up and measured counts; inspect all samples, early/late-half drift and tail
behavior. If compilation is still warming up, repeat the comparison with a
larger declared window rather than selecting favorable samples. Go remains a
fresh CLI process for each observation. Bracket/counterbalance JVM cohorts with
Go CLI runs to expose host drift. Do not run Go alongside a background JVM.

Use several independently restarted cohorts, preferably five or more, and form
paired Go/JVM cohort ratios. Bootstrap cohorts, not individual iterations from
one JVM: JIT, GC and thermal effects make those iterations correlated. Publish
per-cohort results when the available cohort count is too small for a useful
interval. A failed stability check must remain visible.

## Threads and cores

Report compiler job counts 1, 2 and 4 separately from CPU availability.
`--num_parallel_threads` bounds work in compiler phases; it does not set
`GOMAXPROCS`, cap all goroutines, or cap OS threads. Keep Go's runtime CPU
detection unchanged for the primary comparison and give Java the same affinity
and quota. On this environment's two-CPU quota, four jobs is an oversubscription
experiment, not a four-core scaling result. A genuine four-core claim requires
four CPUs of quota as well as suitable affinity. Report the one-/two-core
resource cohorts independently if both are tested.

## Memory and profiling

Both displayed ratios divide Java by Go: above 1 favors Go for elapsed time or
resident memory. Peak RSS is resident process memory, not live heap or allocated
bytes. Go CLI RSS is the distribution of fresh process peaks; Java warm RSS is
the JVM lifetime peak including startup and warm-up. State those windows, and
report medians/ranges or intervals rather than conflating them with live heap.
Record CPU time, cgroup memory/CPU limits, affinity, and wall time as well.

Profile separate instrumented builds, then time uninstrumented builds. CPU
`pprof`, default-sampled allocation profiles, and exit-time `runtime.MemStats`
need no compiler tracing or debug flags. Report allocation bytes/objects per
compile, collections and heap snapshots separately from RSS. An exit-time
`HeapAlloc` is a snapshot, not peak live heap; a forced-GC heap profile measures
retained state after that GC and must be labeled accordingly.

For pass attribution, an overlay can wrap the existing named-pass call in
`pprof` labels without enabling `--tracer_mode`, debug logging or validity
checking. Those modes can select different compiler paths. Optional dispatch
coverage can count callback types at traversal creation, and use traversal-owned
non-atomic counters published at completion. Avoid a global atomic increment per
AST node. Include scope-change checks/skips, native/parser acceptance and
fallback counts, and CFG/data-flow runs to explain *why* a candidate helps. The
instrumentation's coverage evidence is separate from final timing evidence.
