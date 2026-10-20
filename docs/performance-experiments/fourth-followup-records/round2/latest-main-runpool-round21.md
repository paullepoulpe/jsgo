# Latest-main CommonJS run-pool thread screen (round21)

Same frozen `baf328cb` standard binary as round20. Whitespace task only, ordinary fresh CLI, default runtime GC/GOMAXPROCS, CPUs 0,1/two-CPU quota. Five alternating observations per variant at each compiler N=1/2/4, 30 timed runs total. `control` has all factors off; `run_pool` enables only parallel CommonJS, exact-TRUE admission and the persistent pool within one pass/run. It is not a persistent compiler worker. N=1 parallel dispatch is inactive. Other algorithm/parser/type factors remain off.

Each N uses its OWN matched stored warm-Java N cohort and artifact reference. Do not substitute the separate N=1 all-task denominator. Java is warm JVM with fresh compiler/extern parsing per iteration; memory windows differ as documented. This is not fresh paired Java confidence data.

|Compiler N|Variant|Go seconds [min,max]|Go peak RSS MiB [min,max]|Time vs same-N control|Java seconds / RSS MiB|Java/Go time|Java/Go RSS|
|---|---|---:|---:|---:|---:|---:|---:|
|1|control|1.294 [1.250,1.345]|155.7 [152.6,160.0]|+0.00%|0.896 / 990.0|0.692|6.36|
|1|run_pool|1.297 [1.256,1.372]|156.7 [150.8,162.4]|+0.18%|0.896 / 990.0|0.691|6.32|
|2|control|1.215 [1.149,1.248]|158.6 [154.7,164.6]|+0.00%|0.882 / 917.8|0.726|5.79|
|2|run_pool|1.240 [1.193,1.324]|159.9 [152.3,164.3]|+2.03%|0.882 / 917.8|0.711|5.74|
|4|control|1.237 [1.148,1.260]|162.2 [155.9,164.7]|+0.00%|0.848 / 869.3|0.685|5.36|
|4|run_pool|1.227 [1.123,1.246]|155.6 [154.0,165.5]|-0.76%|0.848 / 869.3|0.691|5.59|

All 30 measured exits/artifact digests match Java. Ratios are Java/Go; above 1 favors Go. Peak RSS is measured, not estimated. Source-only node admission is not parallel CPU coverage. Compiler N bounds concurrent phase jobs; it does not set GOMAXPROCS or cap all goroutines/OS threads. Runtime/GC goroutines continue normally. N=4 is oversubscribed on this two-CPU host, so it does not establish four-core scaling. More-core throughput/latency and memory results remain open.

The N=2 actual-DoRun diagnostics from round18 prove that the exact-TRUE/maximal-run design admits 862 whitespace scripts including Angular; stock-factory/lifecycle eligibility is retained in this source. Those diagnostics are not timed scaling data. The original round15 pool measurements used a different historical substrate and must not be mixed with these current-main results.

In this sweep the stock-phase control medians are 1.294 s/155.7 MiB at N=1,1.215 s/158.6 MiB at N=2,and 1.237 s/162.2 MiB at N=4. N=2 is 6.1% faster than N=1 in this cohort; N=4 does not improve further. This is descriptive on a two-CPU host, not evidence of four-core scaling. The pool factor itself is 0.18% slower at N=1 (inactive), 2.03% slower at N=2,and 0.76% faster at N=4 versus its same-N control. These effects do not establish a broadly useful faster phase, especially with overlapping sample ranges.
