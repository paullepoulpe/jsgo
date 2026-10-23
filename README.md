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
- [tools/oracle/](tools/oracle/) — tools that extract ground-truth behavior from the Java compiler ([pin.env](tools/oracle/pin.env)).
- [testdata/](testdata/) — Closure type-checker acceptance fixtures and Java oracles.

---

## Benchmarks (`jsgo` vs. Java Closure Compiler v20261015)

All benchmarks are gated on **100% byte-identical output** (compiled JS, source maps, stdout, stderr diagnostics, and exit code) against `closure-compiler-v20261015.jar`. Measured on macOS arm64 (Go 1.27.1 with PGO enabled via [`cmd/jsgo/default.pgo`](cmd/jsgo/default.pgo), OpenJDK 21.0.10).

### 1. Production-Scale Staged Compilation (`googleapis-bench`)

Compiles **43,917 JS files (~35.1M AST nodes)** across **7,460 chunks** (Closure Library + `google-protobuf` runtime + Closure-style `googleapis` protobufs) in `ADVANCED` mode with type checking (`--num_parallel_threads=4`), staged across `CHECKS → OPTIMIZATIONS → FINALIZATIONS`:

| Stage | `jsgo` | `--fast_mode` | Java (`-Xmx20g`) | Speedup vs. Java | `jsgo` Peak RSS | Java Peak RSS | Memory Saved |
|---|---:|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **49.36 s** | **48.71 s** | 150.70 s | **3.05×** *(3.09×)* | **6,084 MiB** | 10,516 MiB | **−4.33 GiB** |
| **`OPTIMIZATIONS`** | **39.53 s** | **39.30 s** | 102.77 s | **2.60×** *(2.62×)* | **5,525 MiB** | 11,892 MiB | **−6.22 GiB** |
| **`FINALIZATIONS`** | **6.52 s** | 6.56 s | 21.41 s | **3.28×** *(3.26×)* | **1,622 MiB** | 6,750 MiB | **−5.01 GiB** |
| **Total** | **95.41 s** | **94.57 s** | **274.88 s** | **2.88×** *(2.91×)* | **6,084 MiB** | **11,892 MiB** | **−5.67 GiB (1.95× lower)** |

### 2. Large & Mega Whole-Program `ADVANCED` Workloads

All tasks run `--compilation_level=ADVANCED` with source maps enabled. Java Warm is measured after 5 in-process JIT warmup iterations:

| Benchmark | Workload | `jsgo` (`t=1` / `t=4`) | Java (Cold / Warm) | Speedup vs. Cold / Warm (`t=4`) | Peak RSS (`jsgo` vs. Java) |
|---|---|---:|---:|---:|---:|
| **`typed-mega`** | 44.6 MB, 4,385 files *(10 typed Closure apps)* | **15.84 s** / **14.21 s** *(`13.71 s` fast)* | 35.43 s / 27.45 s | **2.49×** / **1.93×** *(`2.58×` / `2.00×`)* | **1,225 MiB** vs. 5,957 MiB (**4.9× lower**) |
| **`bundle-mega`** | 36.1 MB, 26 UMD libraries *(27.6k warnings)* | **20.76 s** / **19.34 s** *(`19.02 s` fast)* | 48.26 s / 41.20 s | **2.50×** / **2.13×** *(`2.54×` / `2.17×`)* | **1,614 MiB** vs. 6,333 MiB (**3.9× lower**) |
| **`typed-full`** | 28.5 MB, 1,818 files *(Closure + Shaka + Blockly)* | **9.21 s** / **8.30 s** *(`7.98 s` fast)* | 23.02 s / 16.10 s | **2.78×** / **1.94×** *(`2.89×` / `2.02×`)* | **869 MiB** vs. 3,325 MiB (**3.8× lower**) |
| **`bundle-advanced-ts`** | 12.6 MB, 4 libraries *(incl. 8.9 MB `typescript.js`)* | **5.77 s** / **5.52 s** | 17.70 s / 11.61 s | **3.21×** / **2.10×** | **615 MiB** vs. 2,358 MiB (**3.8× lower**) |
| **`shaka-blockly`** | 8.6 MB, 847 files *(Shaka Player + Blockly)* | **3.78 s** / **3.40 s** *(`3.34 s` fast)* | 12.50 s / 6.97 s | **3.67×** / **2.05×** *(`3.75×` / `2.09×`)* | **412 MiB** vs. 1,809 MiB (**4.4× lower**) |
| **`bundle-advanced`** | 5.4 MB, 11 UMD libraries *(React, D3, Three, etc.)* | **3.03 s** / **2.75 s** | 11.35 s / 5.79 s | **4.12×** / **2.10×** | **328 MiB** vs. 2,037 MiB (**6.2× lower**) |

### 3. Standard Benchmark Suite ([`tools/bench/tasks/`](tools/bench/tasks/))

| Benchmark | Mode & Workload | `jsgo` (`t=1` / `t=4`) | Java (Cold / Warm) | Speedup vs. Cold / Warm (`t=4`) | Peak RSS (`jsgo` vs. Java) |
|---|---|---:|---:|---:|---:|
| **`typed-app`** | `ADVANCED` + types *(2.8 MB, Closure Library + app)* | **1.21 s** / **1.08 s** *(`1.07 s` fast)* | 6.24 s / 1.99 s | **5.78×** / **1.84×** *(`5.83×` / `1.86×`)* | **225 MiB** vs. 749 MiB (**3.3× lower**) |
| **`bundle`** | `SIMPLE` *(12.6 MB, Three + Lodash + React + TS)* | **3.90 s** / **3.61 s** | 12.54 s / 8.97 s | **3.48×** / **2.49×** | **535 MiB** vs. 2,143 MiB (**4.0× lower**) |
| **`downlevel-whitespace`** | `WHITESPACE_ONLY` *(2.1 MB, Angular + RxJS → ES5)* | **0.45 s** / **0.38 s** | 2.41 s / 0.56 s | **6.33×** / **1.47×** | **138 MiB** vs. 548 MiB (**4.0× lower**) |
| **`downlevel-simple`** | `SIMPLE` *(2.1 MB, Angular + RxJS → ES5)* | **0.89 s** / **0.80 s** | 4.65 s / 1.53 s | **5.84×** / **1.92×** | **181 MiB** vs. 619 MiB (**3.4× lower**) |
| **`checks-only`** | `--checks_only` + types *(11.3 MB, 971 Closure files)* | **2.18 s** / **1.77 s** | 6.92 s / 3.07 s | **3.91×** / **1.73×** | **408 MiB** vs. 1,322 MiB (**3.2× lower**) |
| **`single-files`** | 1,000 standalone CLI invocations (`SIMPLE`) | **103.84 s** | 114.69 s *(batch)* | **1.10×** *(vs. in-JVM batch)* | **98 MiB** vs. 1,577 MiB (**16.1× lower**) |

Across all workloads, `jsgo` is **2.5×–6.3× faster than cold Java**, **1.5×–2.5× faster than a JIT-warmed JVM**, and uses **2.0×–16.1× less memory**. See [docs/benchmarks.md](docs/benchmarks.md) and [docs/performance-roadmap.md](docs/performance-roadmap.md) for full methodology, confidence intervals, and optimization details.
