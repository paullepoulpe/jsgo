# jsgo

A Go port of [Google Closure Compiler](https://github.com/google/closure-compiler), aiming to be a drop-in replacement with the same flags, byte-identical output, and identical diagnostics.

The port is a whole-file, 1:1 translation of the pinned Java release (**v20261015**; originally ported from v20260415 and upgraded), built bottom-up and verified layer by layer against the Java compiler on large real-world corpora.

## Quick Start

```sh
make                # build, vet, and run unit tests
make oracle-tools   # fetch the pinned Java jar/sources and build oracle tools (requires JDK 21+)
```

## Documentation & Layout

- [PLAN.md](PLAN.md) — feasibility assessment, architecture, and milestone roadmap.
- [PORTING.md](PORTING.md) — the 1:1 translation contract every ported file follows.
- [docs/benchmarks.md](docs/benchmarks.md) — benchmark methodology and full results.
- [docs/performance-roadmap.md](docs/performance-roadmap.md) — round-by-round performance engineering log.
- [docs/memory-research/](docs/memory-research/) — memory round 1 (peak live heap −34 to −62 %, peak RSS −26 to −53 %, byte-identical; wall time +9 to +19 %), research on other compilers, and the round-2 design.
- [tools/oracle/](tools/oracle/) — tools that extract ground-truth behavior from the Java compiler ([pin.env](tools/oracle/pin.env)).
- [testdata/](testdata/) — Closure type-checker acceptance fixtures and Java oracles.

---

## Benchmarks (`jsgo` vs. Java Closure Compiler v20261015)

All benchmarks are gated on **100% byte-identical output** (compiled JS, source maps, stdout, stderr diagnostics, and exit code; `--fast_mode` enables additional multi-core passes that relax arbitrary hash-iteration ordering while preserving semantics) against `closure-compiler-v20261015.jar`. Measured on Apple M2 (macOS arm64, 24 GiB RAM; Go 1.27.3 with PGO enabled via [`cmd/jsgo/default.pgo`](cmd/jsgo/default.pgo), OpenJDK 26.0.2.1), with both `jsgo` (`t=4`) and Java running `--num_parallel_threads=4`.

**Statistical methodology ($N = 5$):** Every benchmark runs **$N = 5$ counterbalanced (`AB / BA`) interleaved blocks** to cancel thermal and background drift. Wall times and Peak RSS report sample means; speedups and memory ratios (`Java / jsgo`) report geometric means with **two-sided 95% confidence intervals** (`[lo, hi]`): paired Student-$t$ ($df = 4$) on log-ratios for Cold, Batch, Staged, and Peak RSS comparisons, and two-sample Welch-$t$ on log-times for Warm JVM comparisons (5 measured iterations after 5 in-JVM JIT warmup iterations).

### 1. Production-Scale Staged Compilation (`googleapis-bench`)

Compiles **43,917 JS files (~35.1M AST nodes)** across **7,460 chunks** (Closure Library + `google-protobuf` runtime + Closure-style `googleapis` protobufs) in `ADVANCED` mode with type checking (`--num_parallel_threads=4`), staged across `CHECKS → OPTIMIZATIONS → FINALIZATIONS`:

| Stage | `jsgo` (`t=4`) | `--fast_mode` | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (95% CI) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (95% CI) |
|---|---:|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **50.08 s** | **48.53 s** | 98.76 s | **1.97×** `[1.93, 2.02]` *(2.03× `[1.96, 2.11]`)* | **6,610 MiB** *(6,291)* | 13,272 MiB | **2.01×** `[1.93, 2.09]` *(−6.51 GiB)* |
| **`OPTIMIZATIONS`** | **42.03 s** | **41.46 s** | 88.90 s | **2.12×** `[2.09, 2.14]` *(2.14× `[2.10, 2.19]`)* | **4,540 MiB** *(4,518)* | 14,424 MiB | **3.18×** `[3.00, 3.37]` *(−9.65 GiB)* |
| **`FINALIZATIONS`** | **6.27 s** | 6.32 s | 20.22 s | **3.22×** `[3.18, 3.26]` *(3.20× `[3.17, 3.22]`)* | **1,325 MiB** *(1,349)* | 6,884 MiB | **5.20×** `[4.94, 5.48]` *(−5.43 GiB)* |
| **Total** | **98.38 s** | **96.31 s** | **207.88 s** | **2.11×** `[2.08, 2.14]` *(2.16× `[2.12, 2.20]`)* | **6,610 MiB** *(6,291)* | **14,424 MiB** | **2.18×** `[2.11, 2.25]` *(−7.63 GiB)* |

### 2. Large & Mega Whole-Program `ADVANCED` Workloads

All tasks run `--compilation_level=ADVANCED` with source maps enabled. Peak RSS ratios compare single-compilation processes (`Java Cold / jsgo`; the long-lived Warm JVM process peaks higher at 3.1–8.7 GiB, or 5.4×–19.0× `jsgo`'s RSS):

| Benchmark | Workload | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (95% CI) | Speedup vs. Warm (95% CI) | Peak RSS & Memory Ratio vs. Cold (95% CI) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-mega`** | 44.6 MB, 4,385 files *(10 typed Closure apps)* | **17.97 s** / **15.41 s** *(`14.81 s` fast)* | 34.94 s / 26.74 s | **2.27×** `[2.24, 2.29]` *(2.36× `[2.33, 2.38]`)* | **1.73×** `[1.71, 1.76]` *(1.81× `[1.78, 1.83]`)* | **862 MiB** vs. 3,800 MiB — **4.41×** `[4.32, 4.50]` |
| **`bundle-mega`** | 36.1 MB, 26 UMD libraries *(27.6k warnings)* | **23.08 s** / **20.59 s** *(`20.17 s` fast)* | 49.91 s / 40.82 s | **2.42×** `[2.40, 2.45]` *(2.47× `[2.44, 2.51]`)* | **1.98×** `[1.97, 1.99]` *(2.02× `[2.02, 2.03]`)* | **888 MiB** vs. 5,232 MiB — **5.89×** `[5.53, 6.27]` |
| **`typed-full`** | 28.5 MB, 1,818 files *(Closure + Shaka + Blockly)* | **10.96 s** / **9.48 s** *(`9.06 s` fast)* | 24.45 s / 15.73 s | **2.58×** `[2.44, 2.73]` *(2.70× `[2.56, 2.85]`)* | **1.66×** `[1.55, 1.78]` *(1.74× `[1.63, 1.85]`)* | **540 MiB** vs. 3,523 MiB — **6.53×** `[6.15, 6.92]` |
| **`bundle-advanced-ts`** | 12.6 MB, 4 libraries *(incl. 8.9 MB `typescript.js`)* | **6.70 s** / **6.57 s** | 17.89 s / 13.57 s | **2.73×** `[2.52, 2.96]` | **2.07×** `[1.87, 2.29]` | **306 MiB** vs. 1,907 MiB — **6.22×** `[5.73, 6.75]` |
| **`shaka-blockly`** | 8.6 MB, 847 files *(Shaka Player + Blockly)* | **4.57 s** / **4.03 s** *(`3.78 s` fast)* | 14.73 s / 9.50 s | **3.66×** `[3.58, 3.73]` *(3.90× `[3.87, 3.93]`)* | **2.34×** `[1.97, 2.77]` *(2.49× `[2.10, 2.96]`)* | **267 MiB** vs. 1,627 MiB — **6.06×** `[5.34, 6.86]` |
| **`bundle-advanced`** | 5.4 MB, 11 UMD libraries *(React, D3, Three, etc.)* | **3.33 s** / **2.94 s** | 13.14 s / 8.62 s | **4.42×** `[3.82, 5.13]` | **2.93×** `[2.61, 3.29]` | **184 MiB** vs. 1,597 MiB — **8.65×** `[7.67, 9.76]` |

### 3. Standard Benchmark Suite ([`tools/bench/tasks/`](tools/bench/tasks/))

| Benchmark | Mode & Workload | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (95% CI) | Speedup vs. Warm (95% CI) | Peak RSS & Memory Ratio vs. Cold (95% CI) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | `ADVANCED` + types *(2.8 MB, Closure Library + app)* | **1.40 s** / **1.26 s** *(`1.23 s` fast)* | 7.42 s / 2.21 s | **5.91×** `[5.88, 5.94]` *(6.06× `[5.89, 6.23]`)* | **1.76×** `[1.70, 1.83]` *(1.80× `[1.74, 1.87]`)* | **182 MiB** vs. 893 MiB — **4.91×** `[4.66, 5.18]` |
| **`bundle`** | `SIMPLE` *(12.6 MB, Three + Lodash + React + TS)* | **4.17 s** / **3.81 s** | 12.77 s / 8.68 s | **3.35×** `[3.30, 3.40]` | **2.28×** `[2.26, 2.29]` | **293 MiB** vs. 2,180 MiB — **7.22×** `[5.29, 9.86]` |
| **`downlevel-whitespace`** | `WHITESPACE_ONLY` *(2.1 MB, Angular + RxJS → ES5)* | **0.47 s** / **0.40 s** | 2.34 s / 0.50 s | **5.88×** `[5.79, 5.96]` | **1.26×** `[1.18, 1.36]` | **106 MiB** vs. 559 MiB — **5.27×** `[4.98, 5.58]` |
| **`downlevel-simple`** | `SIMPLE` *(2.1 MB, Angular + RxJS → ES5)* | **0.96 s** / **0.83 s** | 4.67 s / 1.50 s | **5.64×** `[5.47, 5.82]` | **1.81×** `[1.74, 1.89]` | **132 MiB** vs. 629 MiB — **4.78×** `[4.66, 4.89]` |
| **`checks-only`** | `--checks_only` + types *(11.3 MB, 971 Closure files)* | **2.22 s** / **1.72 s** | 7.01 s / 2.86 s | **4.07×** `[3.95, 4.19]` | **1.66×** `[1.59, 1.73]` | **336 MiB** vs. 1,402 MiB — **4.18×** `[3.91, 4.47]` |
| **`single-files`** | 1,000 standalone CLI invocations (`SIMPLE`) | **98.54 s** | — / 112.88 s *(batch)* | — | **1.15×** `[1.13, 1.16]` *(vs. batch)* | **72.5 MiB** vs. 1,569 MiB — **21.5×** `[17.8, 25.8]` |

Across all workloads, `jsgo` is **2.1×–5.9× faster than cold Java** (`2.2×–6.1×` with `--fast_mode`), **1.3×–2.9× faster than a JIT-warmed JVM** (**1.15× faster** even when spawning 1,000 fresh `jsgo` CLI processes against a single in-JVM batch runner), and uses **2.2×–21.5× less memory**. See [docs/benchmarks.md](docs/benchmarks.md) and [docs/performance-roadmap.md](docs/performance-roadmap.md) for historical benchmarks, pass-level breakdowns, and optimization details.
