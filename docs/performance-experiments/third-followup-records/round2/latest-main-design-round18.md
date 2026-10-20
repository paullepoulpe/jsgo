# Clean latest-main prototype screen (round18)

Pinned main: `ff06ac15304f14c756bce9cd36ad2205212f51da`. Ordinary fresh Go CLI, standard Go1.27.2 inlining, main’s shipped bundle-trained PGO, default runtime GC/GOMAXPROCS, affinity CPUs0,1 (two-CPU quota). No worker, extern reuse or GC tuning. All tasks use compiler N=1. This is an initial two-observation screen, not a confidence estimate or a demonstrated overall improvement.

Sources: stock main `ff06ac15`; integrated control16 `a3a2f8f3`; compact4 and compound4 `9044aff9`. Relative to a3 the compact source changes only the private scopes buffer16→4 in production; remaining changes are tests. Control16 and compact4 have all optional factors off. Compound4 enables the recipe/shared resolved vector, owned ASCII parser, and cardinality/projected/scalar lazy-use backend. It also enables the CommonJS guard/run-pool flags, which are inactive at N=1. Exact flags and binary hashes are in the command and machine records.

|Task|Variant|Go median s [min,max]|Go peak RSS MiB|Java/Go time|Java/Go RSS|
|---|---|---:|---:|---:|---:|
|typed-app|stock_main|5.467 [5.015,5.919]|257.0|0.684|4.24|
|typed-app|control16|5.085 [5.068,5.103]|259.0|0.735|4.21|
|typed-app|compact4|5.070 [5.039,5.100]|253.0|0.738|4.31|
|typed-app|compound4|4.987 [4.935,5.039]|255.3|0.750|4.27|
|bundle|stock_main|17.900 [17.761,18.038]|574.5|0.927|2.95|
|bundle|control16|18.281 [17.952,18.611]|603.0|0.908|2.81|
|bundle|compact4|18.124 [17.953,18.295]|612.5|0.916|2.77|
|bundle|compound4|18.216 [17.939,18.493]|595.0|0.911|2.85|
|downlevel-whitespace|stock_main|1.270 [1.223,1.317]|154.9|0.739|6.00|
|downlevel-whitespace|control16|1.274 [1.216,1.333]|157.0|0.736|5.92|
|downlevel-whitespace|compact4|1.264 [1.250,1.278]|158.1|0.742|5.88|
|downlevel-whitespace|compound4|1.284 [1.263,1.306]|156.7|0.731|5.94|
|downlevel-simple|stock_main|4.040 [3.915,4.165]|201.3|0.694|5.56|
|downlevel-simple|control16|3.975 [3.922,4.028]|194.3|0.706|5.76|
|downlevel-simple|compact4|3.856 [3.842,3.869]|195.6|0.728|5.72|
|downlevel-simple|compound4|3.837 [3.806,3.869]|195.0|0.731|5.73|
|checks-only|stock_main|8.118 [7.969,8.267]|657.0|0.744|2.07|
|checks-only|control16|8.028 [7.976,8.079]|678.1|0.752|2.01|
|checks-only|compact4|7.977 [7.932,8.022]|663.8|0.757|2.05|
|checks-only|compound4|7.937 [7.815,8.060]|658.4|0.761|2.07|

Ratios are Java/Go, above1 favors Go. Java references are the same stored warm-JVM cohorts as the baseline, with a fresh compiler and reparsed externs per iteration. They are not a fresh paired Java rerun or settled-JIT proof. Go RSS is the median fresh-process lifetime peak; Java RSS is the warm-JVM lifetime peak including warmup. RSS is not retained heap or allocated bytes.

|Task|Compact4 time vs control16|Compound4 time vs control16|Compound4 time vs compact4|Compound4 RSS vs control16|
|---|---:|---:|---:|---:|
|typed-app|-0.31%|-1.93%|-1.63%|-1.44%|
|bundle|-0.86%|-0.36%|+0.51%|-1.34%|
|downlevel-whitespace|-0.84%|+0.77%|+1.62%|-0.20%|
|downlevel-simple|-3.00%|-3.46%|-0.47%|+0.39%|
|checks-only|-0.63%|-1.12%|-0.49%|-2.91%|

Negative time changes mean less elapsed time. The compound is about1.9% faster on typed and3.5% on SIMPLE versus control16, but bundle is flat and whitespace slower. Its additional benefit over compact4 alone is below2% in every task. The stock typed control has a5.015–5.919s range; its unusually slow observation must not be turned into an8.8% code-win headline. Bundle RSS ranges vary by hundreds of MiB even for the same binary; two-sample peak medians cannot establish a reliable memory improvement. No broad promotion or warm-Java victory is supported.

All40 timed runs match every pinned Java artifact digest and successful exit. Both clean source variants separately passed15 CLI artifact gates (five tasks × N1/2/4). The compact4 focused race binary passes146 tests,0skips,0failures and0data races, including fixed real-Angular admission and deep growth fixtures. The original a3 race run’s Angular fixture failure is retained; its retry correction changes the test fixture, not production admission. This is exercised compatibility coverage, not whole-port parity.

Untimed actual DoRun counters confirm the ASCII and FSIV paths execute on real inputs; they do not prove speedup, allocation or RSS savings. See `latest-main-composed-activation-report.md` and the frozen diagnostic overlay/commands.

A source audit also found stale embedded scope entries after a stack backing array grows. Popping clears the current backing but leaves the original embedded entries reachable. Smaller capacity creates additional overflow cases. A separate unbuilt followup clears obsolete entries and snapshots nodes across reentrant custom creator callbacks; it is excluded from these numbers and requires validation. This prevents treating smaller modeled struct size as proven lower RSS.

Next discriminator: a same-binary factor screen separates parser, recipe/sharing and use-analysis changes. Additional source-only native gather, captured validation and retention prototypes are not measured improvements.

Raw evidence: `latest-main-design-screen/summary.json`, `machine.json`, per-variant `jsgo.json`, `latest-main-design-screen-command.json`, `latest-main-design-screen.log`, central gate logs, focused race logs and source bundle manifests.
