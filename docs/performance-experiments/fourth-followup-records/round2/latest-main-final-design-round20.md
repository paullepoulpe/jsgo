# Final clean-main candidate screen (round20)

Frozen source `baf328cba812a7ebadd2eb3f9912107382805e8a`, descended from main `ff06ac15`. This adds the inline-scope retention repair, captured Array validation and private CFG-local gathering to the earlier compact4 integration. Ordinary fresh CLI, Go 1.27.2 standard inlining, shipped main PGO, default runtime GC/GOMAXPROCS, CPUs 0,1 and compiler N=1. No worker or extern reuse. Build-only resource controls are recorded separately and were cleared for measurement. No builds, tests, profiles or archiving overlapped timings.

Five variants × five tasks × three observations = 75 timed runs. `control_next` disables all optional factors; `gather` enables only the private gather walker; `recipe_shared` enables the Array recipe/shared-vector group; `recipe_captured` additionally enables captured validation; `compound` combines captured recipe/sharing, gather, and cardinality/projected/scalar use analysis. Parser and CommonJS factors are off. All variants use the same binary; the unconditional compact-stack/retention changes are in their common control. This is NOT a paired stock-main comparison or an additional whole-main speed claim.

|Task|Variant|Go seconds [min,max]|Change vs control|Peak RSS MiB [min,max]|Java/Go time|Java/Go RSS|
|---|---|---:|---:|---:|---:|---:|
|typed-app|control_next|5.142 [5.086,5.448]|+0.00%|251.2 [251.0,254.6]|0.727|4.34|
|typed-app|gather|5.209 [5.127,5.225]|+1.30%|254.3 [253.1,255.8]|0.718|4.29|
|typed-app|recipe_shared|5.162 [5.099,5.235]|+0.39%|253.0 [252.5,267.6]|0.724|4.31|
|typed-app|recipe_captured|5.508 [5.240,5.563]|+7.12%|252.8 [250.1,254.6]|0.679|4.31|
|typed-app|compound|5.314 [5.173,5.757]|+3.34%|254.3 [252.8,258.0]|0.704|4.29|
|bundle|control_next|18.270 [18.103,18.675]|+0.00%|495.6 [493.7,511.8]|0.909|3.42|
|bundle|gather|18.696 [18.444,18.918]|+2.34%|600.9 [490.9,685.5]|0.888|2.82|
|bundle|recipe_shared|18.313 [18.058,18.419]|+0.24%|582.3 [438.9,585.8]|0.906|2.91|
|bundle|recipe_captured|18.334 [18.185,18.805]|+0.35%|587.0 [499.3,598.6]|0.905|2.89|
|bundle|compound|18.230 [18.127,18.766]|-0.22%|541.2 [508.1,552.7]|0.911|3.13|
|downlevel-whitespace|control_next|1.340 [1.261,1.363]|+0.00%|149.9 [147.7,159.1]|0.700|6.21|
|downlevel-whitespace|gather|1.354 [1.348,1.357]|+1.10%|155.6 [155.1,157.6]|0.693|5.98|
|downlevel-whitespace|recipe_shared|1.335 [1.273,1.337]|-0.31%|159.8 [148.0,160.1]|0.703|5.82|
|downlevel-whitespace|recipe_captured|1.327 [1.255,1.363]|-0.94%|157.1 [153.8,160.7]|0.707|5.92|
|downlevel-whitespace|compound|1.325 [1.287,1.383]|-1.11%|155.2 [151.8,157.0]|0.708|5.99|
|downlevel-simple|control_next|3.899 [3.794,4.004]|+0.00%|196.7 [195.4,199.1]|0.719|5.69|
|downlevel-simple|gather|3.918 [3.845,3.970]|+0.49%|195.1 [193.8,200.9]|0.716|5.73|
|downlevel-simple|recipe_shared|3.879 [3.762,3.928]|-0.51%|197.1 [194.7,199.8]|0.723|5.67|
|downlevel-simple|recipe_captured|3.922 [3.906,4.008]|+0.59%|195.3 [192.1,200.7]|0.715|5.73|
|downlevel-simple|compound|3.907 [3.886,3.970]|+0.19%|198.9 [196.3,202.8]|0.718|5.62|
|checks-only|control_next|8.174 [8.025,8.187]|+0.00%|662.2 [645.2,662.3]|0.739|2.05|
|checks-only|gather|8.307 [8.040,8.576]|+1.63%|664.7 [660.5,671.8]|0.727|2.05|
|checks-only|recipe_shared|7.969 [7.864,8.012]|-2.51%|625.0 [618.8,651.0]|0.758|2.18|
|checks-only|recipe_captured|8.147 [8.070,8.165]|-0.33%|641.2 [629.2,653.1]|0.741|2.12|
|checks-only|compound|8.035 [7.878,8.179]|-1.70%|632.3 [629.5,650.5]|0.751|2.15|

Negative time changes mean less elapsed time. All 75 measured exits/artifact digests match pinned Java. Ratios are Java/Go; above 1 favors Go. Java uses the stored warm-JVM reference with fresh compiler and reparsed externs per iteration. These are not fresh paired Java confidence data or settled-JIT proof. Go RSS is the median fresh-process lifetime peak; Java RSS is the warm-JVM lifetime peak including warmup. Repetition ranges and inactive-workload variation limit conclusions; modeled buffer savings and cumulative allocation profiles are not peak RSS.

Compatibility evidence: the plain focused race run passes 205 test/subtest records with one optional Angular fixture skipped because its input environment was absent. Running only that fixture with the pinned input passes separately, bringing exercised passing records to 206 across the two runs. No failures/data races. All 15 CLI gates match at N=1/2/4. Real-input diagnostic runs match all stdout/stderr/exit/JS/map digests and prove actual gather/captured execution; see `latest-main-next-activation-report.md`. These cover exercised behavior, not the entire port.

The frozen source bundle published before these tests correctly labels this integration unvalidated at its capture time; these later records supersede that status without rewriting the published archive. The source commit is unchanged. Raw source/build/overlay/test/activation/timing commands and identities are retained.

No broad default promotion, exhaustive search or fixed theoretical maximum is inferred. The report and open-idea notes distinguish a measured negative from an unimplemented design. Any future attempt to improve the new backends should profile admission/staging/dispatch and allocation costs on the current binary rather than assume object elimination guarantees speed. Fresh paired stock-main/Java confirmation is required before claiming parity or a broad win.

The private gather walker is 1.30% slower on typed-app and 2.34% slower on bundle versus its same-binary control; bundle median RSS is 600.9 vs 495.6 MiB. It does not produce a measured win here despite avoiding 300,227 bundle CFG-local traversal constructions. This is a three-sample default-runtime result, not proof that every refined implementation must lose.

Captured validation is 6.71% slower on typed-app and 2.23% slower on checks-only than the proper same-binary recipe_shared control. Verified positive execution therefore does not establish a useful faster factory. The compound is 3.34% slower on typed, 0.22% faster on bundle, 1.11% faster on whitespace, 0.19% slower on SIMPLE and 1.70% faster on checks. All five Java/Go time ratios remain below 1. No broad promotion or Java victory is supported. Some factors are inactive on whitespace/checks; their variation remains a control/noise observation.
