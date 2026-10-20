# Clean latest-main factor screen (round19)

Follow-up to round18 on the same compact4 binary, source9044aff9, pinned ff06ac15 baseline. Five variants × five tasks × two observations =50timed runs, all pinned Java artifact digests and exits MATCH. Same ordinary fresh CLI, standard Go1.27.2/new shipped PGO/default runtime GC andGOMAXPROCS, CPUs0,1, compilerN=1. No builds/tests/profiling overlapped measurement.

Factors: control4 alloff; parser4 onlyownedASCII; recipe_shared4 onlyArrayunknownrecipe +sharedresolvedvector; use_scalar4 cardinality+projection+scalar; compound4 allthreegroups. CommonJS flags are off in this cohort; its N1 workerpath is inactive in round18. Exact environments, source/binary hashes and commands remain in raw metadata. Captured-validation and nativeGather factors are absent from this binary.

|Task|Variant|Go median s [min,max]|Time vs control|Go peak RSS MiB|Java/Go time|Java/Go RSS|
|---|---|---:|---:|---:|---:|---:|
|typed-app|control4|5.244 [5.235,5.252]|+0.00%|252.5|0.713|4.32|
|typed-app|parser4|5.233 [5.226,5.241]|-0.19%|253.1|0.715|4.31|
|typed-app|recipe_shared4|5.274 [5.073,5.476]|+0.58%|253.4|0.709|4.30|
|typed-app|use_scalar4|5.117 [5.067,5.166]|-2.42%|258.8|0.731|4.21|
|typed-app|compound4|5.159 [5.043,5.274]|-1.61%|251.7|0.725|4.33|
|bundle|control4|18.222 [17.851,18.594]|+0.00%|622.0|0.911|2.73|
|bundle|parser4|18.438 [18.327,18.549]|+1.19%|557.2|0.900|3.04|
|bundle|recipe_shared4|18.228 [18.061,18.395]|+0.03%|614.3|0.911|2.76|
|bundle|use_scalar4|18.228 [18.071,18.385]|+0.03%|516.7|0.911|3.28|
|bundle|compound4|17.927 [17.869,17.985]|-1.62%|650.6|0.926|2.61|
|downlevel-whitespace|control4|1.292 [1.262,1.322]|+0.00%|152.4|0.726|6.10|
|downlevel-whitespace|parser4|1.302 [1.283,1.320]|+0.77%|158.9|0.721|5.85|
|downlevel-whitespace|recipe_shared4|1.311 [1.294,1.327]|+1.45%|153.5|0.716|6.06|
|downlevel-whitespace|use_scalar4|1.312 [1.223,1.401]|+1.58%|157.3|0.715|5.91|
|downlevel-whitespace|compound4|1.281 [1.231,1.331]|-0.83%|156.2|0.732|5.95|
|downlevel-simple|control4|3.823 [3.756,3.889]|+0.00%|200.8|0.734|5.57|
|downlevel-simple|parser4|3.895 [3.806,3.984]|+1.90%|196.3|0.720|5.70|
|downlevel-simple|recipe_shared4|3.946 [3.899,3.993]|+3.23%|196.1|0.711|5.70|
|downlevel-simple|use_scalar4|3.920 [3.887,3.953]|+2.55%|200.6|0.716|5.57|
|downlevel-simple|compound4|3.750 [3.709,3.791]|-1.89%|196.9|0.748|5.68|
|checks-only|control4|7.979 [7.951,8.008]|+0.00%|654.9|0.757|2.08|
|checks-only|parser4|8.166 [8.155,8.178]|+2.35%|664.9|0.739|2.05|
|checks-only|recipe_shared4|7.978 [7.965,7.990]|-0.02%|639.4|0.757|2.13|
|checks-only|use_scalar4|8.046 [8.015,8.078]|+0.84%|673.0|0.750|2.02|
|checks-only|compound4|8.028 [7.967,8.090]|+0.62%|628.8|0.752|2.16|

Negative time changes mean less elapsed time. Ratios are Java/Go; above1 favors Go. ReferenceJava is the stored warm-JVM/freshcompiler/reparsedexterns cohort, not a fresh paired rerun. GoRSS is median fresh-process lifetime peak; JavaRSS is warm-JVM lifetime peak including warmup.

The use-analysis factor is2.42% faster on typed-app but flat on bundle and2.55% slower on SIMPLE. Parser alone is essentially flat on typed and slower elsewhere. Recipe/sharing is mixed and3.23% slower on SIMPLE. The compound is1.6–1.9% faster on typed/bundle/SIMPLE in this screen, with checks0.62% slower. These small effects are comparable to inactive-path/control variation and cross-cohort drift; two observations do not confirm a reproducible speedup or interaction. They do not support enabling allfactors by default.

In particular, FSIV does not execute on whitespace/checks-only (verified by actualDoRun diagnostics), yet that factor’s times differ +1.58%/+0.84% fromcontrol. Those rows are noise/control observations, not FSIV algorithm regressions. Do not add round18/19percentages or choose different winning binaries per task as one candidate.

Bundle RSS is especially variable: use_scalar median516.7MiB versus622.0control, whilecompound650.6MiB. These are two lifetime-peak samples, not a reliable retained-memory conclusion. Full raw samples and ranges are preserved. No model-byte counts are substituted for RSS.

No variant beats the stored warmJava reference across these tasks. The most useful next independent target is to eliminate scope-free temporary gather traversals; captured validation also separates fullper-call proof cost from the recipe operation. Both new factors and the stack-retention fix are in a separately frozen integration and require new validation/activation/timing.

Compatibility/race support for this source is described in round18:15Java gates atN1/2/4 and146focusedracePASS/0skip. These50fresh timedmatches are additional exercisedartifactcoverage, not fullportparity.

Raw evidence: `latest-main-factor-screen/summary.json`, `machine.json`, per-variant `jsgo.json`, `latest-main-factor-screen-command.json`, `latest-main-factor-screen.log`.
