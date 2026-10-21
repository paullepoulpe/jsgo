# Failed v4 confirmation: retained partial evidence

The fixed five-cohort confirmation failed in cohort 2 and remains rejected. No restart was selected; that decision closed before these partial performance summaries were inspected. These numbers describe retained observations only: no accepted benchmark summary, five-cohort interval, pooled estimate, or speed promotion.

The independent failure audit checked the runner/manifest/grant/spec/driver identities, loaded class origin and target-class bytes, compiler/resource configuration, zero-pressure status of all saved measured rows, metric conversions, and all 295 recorded six-digest sets against the pinned oracle. The frozen reviewer rejected the run with `failure evidence present: no accepted result`; its rejection is retained.

Cohort 1 contains 60 warmups per JVM and 20 complete measured pairs. Cohort 2 contains 60 warmups per JVM, seven complete measured pairs, one saved unpaired STOCK request 68, and failed accumulator request 68. Cohorts 3–5 never started. There are 296 protocol request directories, 295 complete metric/digest rows (240 warmup, 55 measured), and 27 complete measured pairs. No rows were dropped, replaced, smoothed, or imputed.

## Failure and cleanup

Accumulator request 68 was the second variant request within cohort 2, measured pair 8. The shared 8 GiB parent gained 387 local memory-max events and 14,879 direct-scan/direct-steal counts. The compiler returned exit 0, but the zero-pressure assertion ran before six-digest computation and row persistence. Its stdout/stderr/exit streams remain and independently hash to the oracle; timing, allocation, RSS and JIT metrics were not saved. The final output artifacts exist, but full six-digest parity for this request was not established. STOCK request 68 remains separately available below and in the CSV.

The snapshots show parent memory.current 7,709,859,840 → 8,051,613,696 bytes, below its 8,589,934,592-byte cap at those two instants. They do not rule out transient limit encounters. The shared parent covers paused JVMs and filesystem/kernel/cache state; its counters do not identify an exclusive causal process or prove an accumulator leak or slowdown. PSI was unavailable, not zero. All request and outer OOM deltas were zero.

All four owned JVM drivers exited 0; all four containers were removed without fallback. Cleanup errors, remaining owned containers/processes and active compiler lists were empty. No telemetry exceptions were needed. No JVM/build execution or large-JAR/bulk-artifact rehashing was performed by this audit. All stored six-digest sets were checked; only the failing request’s three streams were rehashed, so this is a scoped failed-run audit rather than a substitute PASS from the frozen reviewer.

## Descriptive comparisons

Each table uses the complete pairs within that cohort. The extra STOCK row is shown separately and never silently discarded. Lifetime HWM uses all retained rows, including warmup; cohort 2 has unequal and incomplete endpoint coverage.

Cohort 1: 20 complete pairs; COMPLETE_COHORT_OF_FAILED_RUN.

| Metric | STOCK | accumulator | Candidate / STOCK |
| --- | ---: | ---: | ---: |
| Elapsed median s | 10.724770394 | 10.873963507 | 1.013911078 |
| Elapsed mean s | 10.785541494 | 10.938052999 | 1.014140366 |
| CPU median s | 12.380000000 | 12.360000000 | 0.998384491 |
| CPU mean s | 12.420000000 | 12.444000000 | 1.001932367 |
| Allocated median bytes | 4756135828.000000000 | 4389663624.000000000 | 0.922947490 |
| Allocated mean bytes | 4756151076.800000191 | 4389709812.000000000 | 0.922954242 |
| Current RSS median MiB | 672.072265625 | 650.662109375 | 0.968143074 |
| Maximum retained lifetime HWM MiB | 1597.875000000 | 1526.824218750 | 0.955534206 |

STOCK: all 20 available measured JIT deltas [47, 97, 86, 81, 42, 148, 12, 16, 7, 202, 101, 2, 42, 13, 36, 35, 13, 69, 25, 36] ms; total 1110 ms / 215.710829872 s = 0.514578%; >1%: False; >2%: False; invalid intervals: 0.
Warm gate independently recomputed at request 60: drift 0.833763%, JIT/wall 0.899055%, spread 1.084540692.

Fixed five-request JIT blocks: 61–65: 0.655774%, 66–70: 0.721554%, 71–75: 0.353783%, 76–80: 0.331543%.

accumulator: all 20 available measured JIT deltas [46, 67, 214, 78, 4, 165, 11, 1, 6, 118, 318, 17, 11, 1, 21, 40, 18, 21, 16, 30] ms; total 1203 ms / 218.761059976 s = 0.549915%; >1%: False; >2%: False; invalid intervals: 0.
Warm gate independently recomputed at request 60: drift -0.452652%, JIT/wall 0.623119%, spread 1.063479721.

Fixed five-request JIT blocks: 61–65: 0.739587%, 66–70: 0.555900%, 71–75: 0.662982%, 76–80: 0.232313%.

Cohort 1 last10/first10 median elapsed drift: STOCK 1.019075%; accumulator -0.136200%. Diagnostic only; no rows excluded.

Cohort 2: 7 complete pairs; PARTIAL_COHORT_OF_FAILED_RUN.

| Metric | STOCK | accumulator | Candidate / STOCK |
| --- | ---: | ---: | ---: |
| Elapsed median s | 10.866691977 | 10.835871869 | 0.997163800 |
| Elapsed mean s | 10.865314758 | 10.889190525 | 1.002197430 |
| CPU median s | 12.500000000 | 12.620000000 | 1.009600000 |
| CPU mean s | 12.485714286 | 12.490000000 | 1.000343249 |
| Allocated median bytes | 4671851856.000000000 | 4378997456.000000000 | 0.937315136 |
| Allocated mean bytes | 4671910114.285714149 | 4378745294.857142448 | 0.937249474 |
| Current RSS median MiB | 660.292968750 | 654.027343750 | 0.990510841 |
| Maximum retained lifetime HWM MiB | 1491.132812500 | 1426.421875000 | 0.956602834 |

STOCK: all 8 available measured JIT deltas [46, 63, 94, 56, 65, 157, 17, 5] ms; total 503 ms / 86.795314804 s = 0.579524%; >1%: False; >2%: False; invalid intervals: 0.
Warm gate independently recomputed at request 60: drift -0.210647%, JIT/wall 0.712992%, spread 1.065354609.

Fixed five-request JIT blocks: 61–65: 0.592410%, 66–68: 0.557572% (partial).

accumulator: all 7 available measured JIT deltas [24, 60, 114, 76, 27, 132, 17] ms; total 450 ms / 76.224333677 s = 0.590363%; >1%: False; >2%: False; invalid intervals: 0.
Warm gate independently recomputed at request 60: drift 0.461209%, JIT/wall 0.459397%, spread 1.079830179.

Fixed five-request JIT blocks: 61–65: 0.551991%, 66–67: 0.686812% (partial).

Unpaired STOCK request 68: elapsed 10.738111495 s; CPU 12.240000000 s; allocated 4671884576 bytes; current RSS 666.574218750 MiB; lifetime HWM 1491.132812500 MiB; JIT delta 5 ms. Accumulator request 68 has no saved metric row.

CompilationMXBean reports approximate accumulated compilation elapsed time, may sum compiler threads, and may reflect compiles spanning container pauses. The JIT/wall annotations are not CPU percentages or proof of a performance cause. Allocation is a VM-wide request delta; RSS/HWM measures have the different scopes described above.

The screen findings remain exploratory, and this failed direct-STOCK attempt cannot establish a speed gain. Source correctness and separate upstream test results do not turn an incomplete performance run into confirmation.

Evidence: `failure-audit.json`, `retained-row-index.csv`, `partial-descriptive.json`, `all-available-rows.csv` (all 295 saved rows plus the explicit missing row), `measured-rows-with-JIT.csv`, and `frozen-reviewer-rejection/review.json`. Original results, runner, prespecification, root README and publication config were not changed.
