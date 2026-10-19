# Current-main actual CLI activation evidence

All five fresh stock CLI invocations exited successfully. CommonJS, ASCII and FSIV were observed during the same DoRun per task. These are untimed normal-build diagnostics, not speed or peak-memory measurements.

Source `a3a2f8f3`; build/run provenance: `latest-main-composed-activation-build-command.json`, `latest-main-composed-activation-run-commands.json`, and `latest-main-focused-composed-activation-v3/manifest.json`. The diagnostic test binary uses PGO off. Affinity is CPUs 0,1 and compiler `--num_parallel_threads=2`; that job budget is separate from GOMAXPROCS. No runtime GC/GOMAXPROCS override is listed in the run command records. Build-only resource settings are not runtime settings; effective runtime env is not independently recorded in these activation JSONs.

## CommonJS dispatch

| Task | Scripts / eligible / dispatched | Maximal run sizes | Launches | Active job limit |
| --- | ---: | --- | ---: | ---: |
| typed-app | pass not executed | — | 0 | — |
| bundle | pass not executed | — | 0 | — |
| downlevel-whitespace | 862 / 862 / 862 | 862 | 1 | 2 |
| downlevel-simple | 863 / 863 / 863 | 863 | 1 | 2 |
| checks-only | pass not executed | — | 0 | — |

Both executed passes accepted the stock ownership/observer guards, used exact-TRUE admission and the run-pool factor, and completed. Every script was eligible and dispatched; there were no serial fallback gaps or singleton runs. A run size of 862/863 is the queued input count, not that many simultaneous workers. The private work-item limit remains two.

The Angular core script was actually dispatched in both downlevel tasks: **79,096 nodes and 1,932 functions**. Total dispatched nodes were 188,633 for whitespace and 188,634 for SIMPLE. Angular accounts for about 41.9% of those node counts; the next-largest script has 5,299 nodes. This establishes admission and potential scheduling skew, not parallel CPU coverage.

## ASCII backing

| Task | Backed / parse calls | Backed source bytes / all source bytes |
| --- | ---: | ---: |
| typed-app | 392 / 407 | 4,909,018 / 5,559,292 |
| bundle | 124 / 132 | 10,860,944 / 13,841,169 |
| downlevel-whitespace | 986 / 995 | 2,900,877 / 4,278,306 |
| downlevel-simple | 986 / 995 | 2,900,877 / 4,278,306 |
| checks-only | 1,055 / 1,099 | 12,060,140 / 19,476,443 |

The probe is attached to the actual runner before DoRun. `backed` confirms the guarded ASCII constructor route was selected; counts include externs. Byte counts are input payload sizes, not allocations or RSS.

## FSIV native use analysis

| Task | Completed pass invocations | Gathered function visits | Native admitted / attempted | Projected / scalar runs | Queries / singleton answers |
| --- | ---: | ---: | ---: | ---: | ---: |
| typed-app | 2 | 1,460 | 55 / 55 | 55 / 55 | 187 / 30 |
| bundle | 2 | 17,685 | 795 / 795 | 795 / 795 | 2,480 / 313 |
| downlevel-whitespace | not executed | — | — | — | — |
| downlevel-simple | 2 | 4,311 | 279 / 279 | 279 / 279 | 690 / 135 |
| checks-only | not executed | — | — | — | — |

Every attempted native program was admitted, projected and handled by the scalar backend; all recorded rejection counters are zero. This applies to the lazy attempted subset, not every gathered function/candidate or all compiler analysis. Function counters are visits across completed pass invocations, not distinct source functions.

| Task | Original dense model bytes | Projected dense model bytes | Scalar state + action accounting bytes |
| --- | ---: | ---: | ---: |
| typed-app | 367,904 | 128,792 | 47,374 |
| bundle | 6,842,200 | 2,150,096 | 791,072 |
| downlevel-simple | 495,584 | 166,016 | 136,668 |

These byte fields sum per-function modeled/accounted buffers over the run. They omit other objects/runtime overhead and do not describe simultaneously live memory, peak heap or RSS. `DenseBytes` and `MaxDenseBytes` were zero because all admitted runs used scalar storage; this is not a claim of zero analysis allocation. Actual RSS and Java/Go memory ratios belong in the separate timed benchmark report.

## Existing artifact checks

Read-only SHA-256 checks matched the pinned baseline for all **five emitted artifact files** (four JS outputs and the typed-app source map), all **five captured diagnostic streams**, and all **five exit statuses**. Checks-only emits no JS artifact. The reference was `baseline/TASK/java-cold.json`, used strictly for artifact/diagnostic/exit provenance. Activation stdout is not retained by this harness and was not independently checked here. Per-file paths, sizes and digests are in `latest-main-composed-activation-report.json`. No output was overwritten and no compiler was invoked for this report.

## Exact enabled factor flags

These match both the saved activation command env and `compound4` in `latest-main-design-screen-command.json`:

```json
{
  "JSGO_EXPERIMENT_ARRAY_UNKNOWN_RECIPE": "1",
  "JSGO_EXPERIMENT_ARRAY_SHARED_RESOLVED": "1",
  "JSGO_EXPERIMENT_PARSER_ASCII_SOURCE": "1",
  "JSGO_EXPERIMENT_PARALLEL_COMMONJS": "1",
  "JSGO_EXPERIMENT_COMMONJS_EXACT_TRUE_GUARD": "1",
  "JSGO_EXPERIMENT_COMMONJS_RUN_POOL": "1",
  "JSGO_EXPERIMENT_FSIV_USE_CARDINALITY": "1",
  "JSGO_EXPERIMENT_FSIV_USE_PROJECTION": "1",
  "JSGO_EXPERIMENT_FSIV_USE_SCALAR": "1"
}
```

The diagnostic binary is the normal a3 source with observation overlays; it is not the compact4 timed binary. Extra scans, probes, test durations and recorded phase elapsed fields are excluded from performance conclusions.
