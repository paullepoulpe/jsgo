# Fresh-template guard V2 screen

All **20 fresh CLI runs matched pinned stock Java**, with zero exit codes, identical output/diagnostic/exit/source-map digests across all samples and variants per task, and all four task/variant gates MATCH. Timing completed successfully at **2026-10-19 11:51:31 UTC**. The same-binary V2 comparison shows small net gains; **Java remains faster and no promotion is made**.

Both variants use production source **c7998d45**, Go 1.27.2, shipped PGO and identical four-package `-l=4` flags. Both retain the same nine earlier native factors, CLI_SOURCE_LOCK ownership witness and two-inline-word BitSet; native liveness is off. `control` disables fresh template values; `fresh` enables them. Only the two tasks that run type inference are measured: typed ADVANCED+checkTypes and checks-only+checkTypes. There are five fresh CLI processes per variant/task, compiler jobs=1, CPU affinity 0,1 under the two-CPU quota. Runtime GOGC/GOMEMLIMIT/GODEBUG/GOMAXPROCS are unset; runtime GC tuning is dropped. No builds, tests or profiling overlapped timing. Raw data/environments/hashes/order: [template-guard-v2-screen/](template-guard-v2-screen/) and [template-guard-v2-screen.log](template-guard-v2-screen.log).

V2 validates and caches the private invocation guard at inference boundaries, retaining original admission checks for callback-bearing paths. Its exact stock error-counter chain reads the current error count under the **same exclusion mutexes** using TryLockNonReentrant, avoiding new goroutine-identifying runtime.Stack calls at each uncontended scope. It still observes newly reported errors. If any mutex cannot be acquired, it releases acquired outer locks and calls the original reentrant getter; unknown counter implementations keep the original guard. This preserves error exclusion and reentry instead of reading an unsynchronized counter.

| Task | Fresh elapsed vs control | Fresh RSS vs control |
|---|---:|---:|
| typed-app | -3.7% | -8.4% |
| checks-only | -1.3% | +0.8% |

Typed elapsed falls from 5.123 to 4.934 s (3.7%), with median peak RSS falling from 275.7 to 252.6 MiB (8.4%). Checks-only falls from 7.504 to 7.405 s (1.3%), while RSS rises from 653.4 to 658.8 MiB (0.8%). Both elapsed ranges overlap. These are five-sample screens, not confidence intervals or independent JVM/Go cohorts.

The runner alternates AB/BA pair order after a seeded starting order: typed starts fresh first, checks starts control first. With five rounds, one variant runs first three times and the other twice. Fresh is faster in four of five within-round pairs for each task, but round pairing does not remove drift or establish significance; medians below summarize all five processes, not a paired confidence estimate.

Java/Go ratios above 1 favor Go. Java is the **older stored stock warm JVM cohort from 2026-10-18 22:08 UTC** (12 warmup/10 measured compilations), compared with **current Go measurements from 2026-10-19 11:49 UTC** on the same recorded host/toolchain. Java creates a fresh compiler and reparses externs each iteration. These are not fresh paired Java comparisons. Go RSS is the median fresh-process lifetime peak; Java RSS is the entire stored warm-JVM lifetime peak including warmup. Memory ratios describe those different windows, not equivalent per-compilation live memory.

| Task | Variant | Go median s [five-run min–max] | Peak MiB median [min–max] | Java/Go time | Java/Go RSS |
|---|---|---:|---:|---:|---:|
| typed-app | control | 5.123 [4.957–5.382] | 275.7 [275.3–282.9] | 0.730 | 3.95 |
| typed-app | fresh | 4.934 [4.797–5.161] | 252.6 [248.4–284.6] | 0.758 | 4.32 |
| checks-only | control | 7.504 [7.254–7.588] | 653.4 [641.1–668.2] | 0.805 | 2.08 |
| checks-only | fresh | 7.405 [7.330–7.469] | 658.8 [638.3–667.9] | 0.815 | 2.06 |

V1 enabled medians were 5.846 s typed and 9.127 s checks-only; V2 enabled medians are 4.934 and 7.405 s. **That is an across-cohort, across-binary comparison, not an isolated causal benchmark of the guard change.** The appropriate net-effect evidence here is the much smaller same-binary control/fresh difference. The large V1 regression is absent in this screen, but this does not establish a broad Java win or justify promoting template vectors or the unconditional BitSet representation.

Validation: compiler race suite **964 PASS, 2 SKIP**, no races; Monitor focused tests **5 normal PASS + 5 race PASS**; all five CLI tasks at compiler jobs **1, 2 and 4** match pinned stock Java (15 gates). All 20 timed raw zero exits/digests and four gates were independently checked for this report. No builds/tests, profiling, timing or source/git edits were performed to prepare the artifact.
