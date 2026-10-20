# Round 16: native Array recipes and shared resolved vectors

This screen uses historical prototype source `7b801677`. It does **not** measure an incremental improvement over latest main `ff06ac15`.

All variants use the same binary and the earlier native9 plus CLI ownership configuration, Go 1.27.2, shipped historical PGO and the historical `-l4` build flags. Runtime GC and GOMAXPROCS remain at their defaults. CPU affinity is 0,1 and the compiler thread count is one. Three rotating repetitions of six variants on two typed tasks produced 36 measured runs; every output matched Java. Validation passed 1,312 focused race cases with two skips, no failures and no detected races, plus all 15 artifact gates at compiler thread counts 1, 2 and 4.

| Task | Variant | Median seconds | Range | RSS MiB | Time change vs control | Java/Go time | Java/Go RSS |
|---|---|---:|---:|---:|---:|---:|---:|
| typed-app | control | 5.330 | 4.963–5.534 | 247.0 | +0.00% | 0.702 | 4.41 |
| typed-app | fresh | 5.157 | 5.117–5.214 | 249.6 | -3.23% | 0.725 | 4.37 |
| typed-app | recipe | 5.264 | 5.245–5.268 | 276.0 | -1.24% | 0.710 | 3.95 |
| typed-app | shared | 5.279 | 5.162–5.286 | 276.1 | -0.95% | 0.708 | 3.95 |
| typed-app | fresh_shared | 5.077 | 5.013–5.372 | 253.3 | -4.74% | 0.737 | 4.30 |
| typed-app | compound_shared | 5.082 | 5.044–5.217 | 277.4 | -4.65% | 0.736 | 3.93 |
| checks-only | control | 7.324 | 7.309–7.631 | 660.9 | +0.00% | 0.824 | 2.06 |
| checks-only | fresh | 7.448 | 7.350–7.507 | 668.7 | +1.70% | 0.811 | 2.03 |
| checks-only | recipe | 7.414 | 7.386–7.599 | 654.5 | +1.23% | 0.814 | 2.08 |
| checks-only | shared | 7.511 | 7.388–7.928 | 638.6 | +2.57% | 0.804 | 2.13 |
| checks-only | fresh_shared | 7.569 | 7.391–7.607 | 638.5 | +3.36% | 0.798 | 2.13 |
| checks-only | compound_shared | 7.503 | 7.481–7.740 | 645.4 | +2.45% | 0.805 | 2.11 |

Negative time changes mean less elapsed time than the same-binary control. The factors are:

- `fresh`: the earlier `FRESH_TEMPLATE_VALUES` factor.
- `recipe`: generalized `Array<?>` construction recipes.
- `shared`: recipes plus private immutable resolved-vector sharing.
- `fresh_shared`: `fresh` and `shared` together.
- `compound_shared`: `fresh_shared` plus the corrected eight-backend compound configuration.

Public wrappers, maps and slices remain fresh; sharing applies only to the private immutable resolved vectors.

The full DoRun diagnostic confirmed activation: 2,558 recipe/shared hits for typed-app and 21,685 for checks-only, with one original seed per task and zero native-mutation fallbacks. Both diagnostic outputs matched Java. An earlier diagnostic omitted the required flags and recorded zero activation. It is retained as a control, separately from the flagged command and positive log; it is not evidence that sharing ran.

Results are mixed. `fresh_shared` reduced typed-app median time by about 4.75% but increased checks-only time by 3.36%. Recipes alone changed these times by −1.24% and +1.23%, respectively. Adding vector sharing provides no clear CPU benefit. Requested-byte reductions did not predict the observed RSS changes. Validation of the 31-key domain on every hit may offset construction savings; an immutable-session proof remains a research option.

There is no broad Java win or default promotion. These are three-run screens with overlapping cohort ranges and no confidence estimate. Do not add their percentages to other cohorts or to latest main. The old substrate includes native-parser and BitSet layout experiments that have not been ported wholesale.

Ratios are Java divided by Go; values above one favor Go. The stored warm Java runner creates a fresh Compiler and reparses externs on every iteration. Go RSS is the median lifetime peak of fresh CLI processes; Java RSS is the warm JVM lifetime peak, including warmup. These are different memory windows. A current-main control must be measured separately.
