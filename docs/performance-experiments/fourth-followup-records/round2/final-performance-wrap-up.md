# Performance investigation wrap-up

The investigation is paused at the user-requested quota/time wrap-up. The 30-minute deadline was 16:21:55 UTC on 2026-10-19; final measurements completed at 16:07:57 UTC, with publication/verification afterward. No further prototypes or benchmarks are planned in this handoff.

Latest main was fetched again and remains `ff06ac15304f14c756bce9cd36ad2205212f51da`. PR23 includes that main; compiler/library code and shipped PGO match it exactly. Experimental implementations are preserved as recoverable sources, not applied production changes. Main already integrates several older ideas; historical percentages are not additional gains on current main.

## Final result

No comfortable or broad warm-Java win was achieved. The final independently validated designs remain small/mixed or slower at default runtime settings. A lack of gains from these implementations is not proof that all drop-in designs are exhausted.

One complete final configuration, compound in round20, has the following three-observation medians versus its same-binary optional-factors-off control. These are not paired stock-main speedup measurements.

|Task|Control → compound s|Time change|Control → compound peak RSS MiB|Java/Go time|Java/Go RSS|
|---|---:|---:|---:|---:|---:|
|typed-app|5.142 → 5.314|+3.34%|251.2 → 254.3|0.704|4.29|
|bundle|18.270 → 18.230|-0.22%|495.6 → 541.2|0.911|3.13|
|downlevel-whitespace|1.340 → 1.325|-1.11%|149.9 → 155.2|0.708|5.99|
|downlevel-simple|3.899 → 3.907|+0.19%|196.7 → 198.9|0.718|5.62|
|checks-only|8.174 → 8.035|-1.70%|662.2 → 632.3|0.751|2.15|

Java/Go ratios above 1 favor Go. Java is the documented stored warm JVM with fresh compiler and reparsed externs. No worker reuse or runtime GC tuning was explored in this final work. Go RSS is median fresh-process lifetime peak; Java RSS is warm-JVM lifetime peak including warmup. Stored Java cohorts, GC-dependent RSS variation and small repetition counts limit conclusions. A fresh paired Java/stock-main rerun remains required before a parity claim.

## Completed evidence

- Round18: 40 matched timed layout/compound runs; round19: 50 matched factor runs. Round20 adds 75 matched final-design runs; round21 adds 30 matched thread runs. All 195 prototype runs match Java artifacts/exits. The separate new-main baseline has 30 matched runs. These counts are distinct from historical measurements in earlier archives.
- Final integration `baf328cb` compiles normally and with race instrumentation. The focused run passes 205 test/subtest records with one optional input-dependent skip; the skipped Angular fixture passes in a separate pinned-input invocation. No failures/data races. Fifteen artifact gates pass at compiler N=1/2/4.
- Five real DoRun diagnostics match full stdout/stderr/exits/JS/maps and prove active paths: gather handles 300,227 bundle CFG nodes with 4 fallback nodes; captured validation hits 2,558 typed and 21,685 checks constructors. Counts are not CPU coverage or RSS savings.
- The small-stack retention fix passes isolated weak-lifetime and balanced reentrant-creator checks. It is present in every final control/candidate; no isolated macro win is claimed for it.
- The new CommonJS run pool is 0.18% slower at N=1, 2.03% slower at N=2,and 0.76% faster at N=4 versus same-N controls. Each N uses its own matched warm Java reference. The host has two CPUs; N=4 is oversubscribed. CompilerN controls phase jobs, not all goroutines or runtime threads.

## What remains open

`open-performance-ideas-final.md` records remaining concrete walkers, operation-local resolved bindings, a correct clean-main native-liveness bridge/publication, private AST regions, constructor census, unmeasured stack8/build alternatives and larger-core parallelism. `open-parser-ir-ideas-final.md` separates incomplete historical directIR work from still-unimplemented shared-grammar event/IR sinks and private value tokens/lazy public positions. Each includes compatibility hazards and the next discriminator.

The strongest untested architectural possibilities require changing representation or removing whole layers of repeated work; component counters alone cannot forecast their benefit. Before refining a negative backend, profile its current guard/staging/dispatch costs, prove actual activation, and repeat matched time/RSS tests. Do not silently loosen alias, exception, source-map, diagnostic-order, mutation or public/custom fallback requirements.

Current-main profile data, failed first capture, corrected activation attribution, positive/negative/inconclusive results, source-only variants, failures/fixes and exact commands remain in the original archive plus immutable successive deltas. The clean-main bundle preserves 12 refs/42 commits (including final `baf328cb`); its capture-time unvalidated labels are historical and superseded by later validation records. Published sources and measurements use pinned commits/toolchains/flags rather than whichever branch is current.

PR23 should remain a draft evidence/runners review. No compiler prototype is promoted from this final screen, and no claim of exhaustive search or a theoretical maximum is supported.
