# Round 15: matched-thread CommonJS rewriting

Do not promote this parallel CommonJS prototype. At N=2, enabling it increased the whitespace task's elapsed median by **7.97%** and CPU median by **10.85%**. At N=4, elapsed median increased **1.43%**. At N=1, admission rejects parallel work; its apparent 2.43% improvement is control variation, not a parallel speedup. All 30 fresh Go invocations matched the stored Java artifacts.

## Matched Java references

Each row uses the pinned stock Java cohort with the same `--num_parallel_threads=N`: 12 warmup compilations followed by 10 measured compilations in one warm JVM, each with a fresh compiler and normal extern parsing. These cohorts were stored earlier, not rerun concurrently with this on/off experiment. Their elapsed medians and lifetime peak RSS are the denominators for that N only.

| Compiler N | Java warm median, s | Java JVM lifetime peak RSS, MiB |
| ---: | ---: | ---: |
| 1 | 0.8962 | 989.98 |
| 2 | 0.8818 | 917.80 |
| 4 | 0.8475 | 869.32 |

## CommonJS off versus on

Go cells show the median followed by `[minimum, maximum]` across five fresh CLI processes. RSS cells are per-process peak RSS, not live heap. Both ratio columns are **Java / Go; above 1 favors Go**. Java remains faster at every N/setting, while Go uses substantially less peak resident memory under the stated measurement windows.

| N | CommonJS | Go elapsed s [min, max] | Go peak RSS MiB [min, max] | Java/Go time | Java/Go RSS |
| ---: | --- | --- | --- | ---: | ---: |
| 1 | off | 1.3047 [1.1813, 1.3765] | 162.86 [161.23, 165.11] | 0.6869 | 6.0786 |
| 1 | on | 1.2730 [1.2032, 1.3930] | 166.94 [164.62, 172.38] | 0.7040 | 5.9301 |
| 2 | off | 1.1763 [1.1443, 1.2243] | 170.52 [165.41, 171.89] | 0.7496 | 5.3824 |
| 2 | on | 1.2701 [1.1856, 1.3613] | 168.75 [164.06, 174.21] | 0.6943 | 5.4389 |
| 4 | off | 1.1705 [1.1487, 1.2420] | 163.11 [157.05, 174.67] | 0.7241 | 5.3298 |
| 4 | on | 1.1872 [1.1737, 1.2212] | 164.05 [159.55, 169.48] | 0.7139 | 5.2991 |

Positive deltas below mean the CommonJS experiment consumed more time or RSS. Ratios of medians differ from the median of paired per-round ratios; both raw round sequences are retained in the aggregate JSON.

| N | Elapsed on−off, s | Elapsed delta | Peak RSS on−off, MiB | RSS delta | CPU median off→on, s | CPU delta |
| ---: | ---: | ---: | ---: | ---: | --- | ---: |
| 1 | -0.0317 | -2.43% | +4.08 | +2.50% | 1.6546→1.5660 | -5.36% |
| 2 | +0.0938 | +7.97% | -1.77 | -1.04% | 1.5220→1.6871 | +10.85% |
| 4 | +0.0167 | +1.43% | +0.95 | +0.58% | 1.5251→1.5908 | +4.31% |

Five interleaved on/off runs at each N provide a small screen, not a confidence interval. On was slower in four of five same-numbered rounds at N=2 and four of five at N=4. At N=1 it was slower in two of five despite sharing the same serial compilation path. The N=2 RSS reduction (1.04%) and N=4 RSS increase (0.58%) are small relative to the per-run RSS ranges; no memory improvement is established. The wider N=2 on-time range also argues against treating the median regression as a precise universal constant.

## Runtime, source and artifact provenance

- CPU affinity was 0,1, with quota `200000 100000` (two CPU equivalents); memory limit was 8 GiB. `--num_parallel_threads` limits eligible compiler work items. N=4 does not create four CPUs or set GOMAXPROCS; N=1 does not make GC/runtime activity single-threaded.
- GOGC, GOMEMLIMIT, GOMAXPROCS and GODEBUG were unset for Go. This is the binary's default runtime configuration, including its embedded build defaults, with no GC tuning.
- Both variants use exactly `/workspace/jsgo-results/round2/bin/jsgo-post-colors-kernel-l4`, SHA-256 `084c5d89af17cd9629ba034862b47a8daf016f210913da04983c68a980fce05d`. Its source checkout is `/workspace/jsgo-post-colors-candidate-kernel`, head `22358ec4da55bab6ffb5c2ec732a1ef7e71485f5`. The recorded runner `bench_commit=4bbbb8a9` identifies the runner checkout, not this binary's source.
- The binary uses Go 1.27.2, parser `-l=4`, and its existing PGO profile. The only environment difference between same-N variants is `JSGO_EXPERIMENT_PARALLEL_COMMONJS=1`. The shared environment is the nine existing native switches plus CLI source ownership; newer kernel/storage switches are not implicitly enabled. Exact flags are in `commonjs-matched-thread-commands.json` and each screen's `machine.json`.
- Java is the pinned stock Closure jar SHA-256 `d6864ebf70bca5ba6e237898de6611a389aaf8b933c2549036ced4a684b664b6`, running OpenJDK 21.0.12.1. Same-N warm JSON/jobs originate from `default-gc-white-traversal-thread-sweep/threads-N/downlevel-whitespace`.
- `commonjs-matched-thread-reference/threads-N` copies the sweep's machine metadata with explicit N and per-file provenance/SHA-256. Its `java-cold.json` and raw cold directory come from `baseline/downlevel-whitespace` solely for artifact digests and byte comparisons; no cold timing or memory is used as a thread-matched denominator.
- Memory ratios compare a fresh Go process peak against the same-N Java JVM lifetime peak including warmups. They quantify those windows, not equal steady-state heap occupancy. Stored Java ratios are not a fresh paired-Java confidence estimate.
- All 30 measured Go processes exited successfully and matched stdout, stderr, JS and configured artifacts. This is compatibility evidence for the fixed whitespace task; it is not a complete Java/drop-in parity claim.

## What activation does and does not show

The separate untimed N=2 activation capture recorded 842 dispatched scripts in 421 two-script cohorts for whitespace. It rejected 10 scripts through the conservative UMD-condition guard and left 10 eligible scripts unpaired. Angular core.mjs was one of the serial rejections. Eligible/dispatched node coverage is not CPU coverage: rejected files have no complete recorded weight, individual node kinds have different costs, and the task's expensive serial work can remain outside those cohorts. The earlier approximately 7.9% CommonJS CPU fraction sets only a broad upper bound; it does not measure this prototype's parallel fraction.

## Source audit of the regression

The source provides concrete overhead mechanisms, but no new profile was captured for this screen. Their relative contribution remains unmeasured.

1. **A fresh worker pool for each tiny cohort.** `runParallelCommonJSJobs` invokes `parallelPrebuild` once per group. At the observed N=2 activation this means 421 channels/joins/failure arrays and 842 worker goroutine starts, for 842 files. Median dispatched file weight was only 59.5 AST nodes (mean 117.1). These static weights identify fine scheduling granularity, not service time.
2. **Additional serial admission and per-file allocation.** Each file gets a preflight AST walk, function-ownership map, job/facade, copied loader/resolver/path headers and buffered effects. Each worker also creates an outer pass/traversal for its SCRIPT; the original path shares its outer pass/root traversal across files. Stable sorting and execution-order slices occur even for two jobs.
3. **A barrier every N scripts.** At N=2 an idle worker cannot take the next independent file until its paired file finishes and events publish. Unequal adjacent file costs can leave one CPU idle. N=4 can amortize barriers and expose more ready work to two CPUs, although at most two can execute CPU work simultaneously. The smaller N=4 regression means oversubscription alone cannot explain why N=2 was worst.
4. **Extra CPU/GC/cache work remains plausible.** N=2 median CPU rose 0.1651 s (10.85%); pure sleeping/channel wait does not itself explain extra CPU consumption. Allocation, goroutine/runtime work, cache behavior and different concurrent GC/assist timing need direct attribution.
5. **No obvious global resolver lock was found.** Workers copy resolver/loader headers and read shared lookup tables. The shared source-line mutex is used by exception formatting, not ordinary successful CommonJS work. Warning guards are admitted read-only and error manager publication is serial. There is no evidence here for a contended pass-wide mutex; do not describe the regression as measured lock contention.

## Next measurement before another implementation

First capture a phase-labeled CPU/allocation profile and untimed per-cohort service/idle spans on this same binary/source and fixed inputs, with CommonJS off/on. Attribute admission/setup, actual finder/rewrite work, event replay, channel/wait/goroutine activity, allocation and GC assistance separately. Record service-time distributions and serial-file phase cost; eligible node counts are insufficient.

Then use a diagnostic-only serial-facade factor: retain admission, private headers and ordered effects, but execute admitted cohort jobs inline. Compare original serial, inline facade and the existing parallel route. This distinguishes abstraction/preflight costs from worker scheduling without relaxing correctness guards or changing global-name publication.

If worker scheduling is material, prototype one pass-owned persistent bounded pool while retaining the original cohort barriers and ordered event publication, then measure it separately. Larger/batched independent windows could reduce barriers but require another ownership/diagnostic/global-name audit. If admission dominates, investigate reusing actual finder discovery rather than adding another complete scan. None of these designs is justified for promotion without the same output/source-map/diagnostic/stamp/race gates and a separate same-binary default-runtime screen. Do not select thresholds from benchmark filenames.

This round retains the prototype as an experiment. It establishes neither a speed win nor comfortable warm-Java parity.
