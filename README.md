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
- [docs/memory-research/](docs/memory-research/) — memory round 1 (peak heap −34 to −62 %, byte-identical), research on other compilers, and the round-2 design.
- [tools/oracle/](tools/oracle/) — tools that extract ground-truth behavior from the Java compiler ([pin.env](tools/oracle/pin.env)).
- [testdata/](testdata/) — Closure type-checker acceptance fixtures and Java oracles.

---

## Benchmarks (`jsgo` vs. Java Closure Compiler v20261015)

All benchmarks are gated on **100% byte-identical output** (compiled JS, source maps, stdout, stderr diagnostics, and exit code) against `closure-compiler-v20261015.jar`. Measured on macOS arm64 (Go 1.27.1 with PGO enabled via [`cmd/jsgo/default.pgo`](cmd/jsgo/default.pgo), OpenJDK 21.0.10).

### 1. Production-Scale Staged Compilation (`googleapis-bench`)

Compiles **43,917 JS files (~35.6M AST nodes)** across **7,460 chunks** (Closure Library + `google-protobuf` runtime + Closure-style `googleapis` protobufs) in `ADVANCED` mode with type checking (`--num_parallel_threads=4`), staged across `CHECKS → OPTIMIZATIONS → FINALIZATIONS`:

| Stage | `jsgo` | `--fast_mode` | Java (`-Xmx20g`) | Speedup vs. Java | `jsgo` Peak RSS | Java Peak RSS | Memory Saved |
|---|---:|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **60.75 s** | 61.14 s | 150.70 s | **2.48×** | **8,697 MiB** | 10,516 MiB | **−1.78 GiB** |
| **`OPTIMIZATIONS`** | **48.94 s** | **48.05 s** | 102.77 s | **2.10×** *(2.14×)* | **8,572 MiB** | 11,892 MiB | **−3.24 GiB** |
| **`FINALIZATIONS`** | **8.68 s** | **8.65 s** | 21.41 s | **2.47×** *(2.48×)* | **2,490 MiB** | 6,750 MiB | **−4.16 GiB** |
| **Total** | **118.37 s** | **117.85 s** | **274.88 s** | **2.32×** *(2.33×)* | **8,697 MiB** | **11,892 MiB** | **−3.12 GiB (1.37× lower)** |

### 2. Large & Mega Whole-Program `ADVANCED` Workloads

All tasks run `--compilation_level=ADVANCED` with source maps enabled. Java Warm is measured after 5 in-process JIT warmup iterations:

| Benchmark | Workload | `jsgo` (`t=1` / `t=4`) | Java (Cold / Warm) | Speedup vs. Cold / Warm (`t=4`) | Peak RSS (`jsgo` vs. Java) |
|---|---|---:|---:|---:|---:|
| **`typed-mega`** | 44.6 MB, 4,385 files *(10 typed Closure apps)* | **17.31 s** / **16.72 s** *(`16.13 s` fast)* | 35.43 s / 27.45 s | **2.12×** / **1.64×** *(`2.20×` / `1.70×`)* | **1,581 MiB** vs. 5,957 MiB (**3.8× lower**) |
| **`bundle-mega`** | 36.1 MB, 26 UMD libraries *(27.6k warnings)* | **21.15 s** / **20.27 s** *(`20.22 s` fast)* | 48.26 s / 41.20 s | **2.38×** / **2.03×** *(`2.39×` / `2.04×`)* | **1,673 MiB** vs. 6,333 MiB (**3.8× lower**) |
| **`typed-full`** | 28.5 MB, 1,818 files *(Closure + Shaka + Blockly)* | **9.97 s** / **9.50 s** *(`9.23 s` fast)* | 23.02 s / 16.10 s | **2.42×** / **1.69×** *(`2.49×` / `1.74×`)* | **995 MiB** vs. 3,325 MiB (**3.3× lower**) |
| **`bundle-advanced-ts`** | 12.6 MB, 4 libraries *(incl. 8.9 MB `typescript.js`)* | **5.99 s** / **5.86 s** *(`5.85 s` fast)* | 17.70 s / 11.61 s | **3.02×** / **1.98×** *(`3.03×` / `1.99×`)* | **500 MiB** vs. 2,358 MiB (**4.7× lower**) |
| **`shaka-blockly`** | 8.6 MB, 847 files *(Shaka Player + Blockly)* | **4.07 s** / **3.88 s** *(`3.75 s` fast)* | 12.50 s / 6.97 s | **3.22×** / **1.80×** *(`3.33×` / `1.86×`)* | **458 MiB** vs. 1,809 MiB (**3.9× lower**) |
| **`bundle-advanced`** | 5.4 MB, 11 UMD libraries *(React, D3, Three, etc.)* | **3.06 s** / **2.91 s** | 11.35 s / 5.79 s | **3.90×** / **1.99×** | **312 MiB** vs. 2,037 MiB (**6.5× lower**) |

### 3. Standard Benchmark Suite ([`tools/bench/tasks/`](tools/bench/tasks/))

| Benchmark | Mode & Workload | `jsgo` (`t=1` / `t=4`) | Java (Cold / Warm) | Speedup vs. Cold / Warm (`t=4`) | Peak RSS (`jsgo` vs. Java) |
|---|---|---:|---:|---:|---:|
| **`typed-app`** | `ADVANCED` + types *(2.8 MB, Closure Library + app)* | **1.31 s** / **1.23 s** | 6.24 s / 1.99 s | **5.07×** / **1.62×** | **248 MiB** vs. 749 MiB (**3.0× lower**) |
| **`bundle`** | `SIMPLE` *(12.6 MB, Three + Lodash + React + TS)* | **3.93 s** / **3.75 s** | 12.54 s / 8.97 s | **3.34×** / **2.39×** | **537 MiB** vs. 2,143 MiB (**4.0× lower**) |
| **`downlevel-whitespace`** | `WHITESPACE_ONLY` *(2.1 MB, Angular + RxJS → ES5)* | **0.52 s** / **0.43 s** | 2.41 s / 0.56 s | **5.60×** / **1.30×** | **155 MiB** vs. 548 MiB (**3.5× lower**) |
| **`downlevel-simple`** | `SIMPLE` *(2.1 MB, Angular + RxJS → ES5)* | **0.98 s** / **0.90 s** | 4.65 s / 1.53 s | **5.17×** / **1.70×** | **193 MiB** vs. 619 MiB (**3.2× lower**) |
| **`checks-only`** | `--checks_only` + types *(11.3 MB, 971 Closure files)* | **2.24 s** / **2.00 s** | 6.92 s / 3.07 s | **3.46×** / **1.53×** | **495 MiB** vs. 1,322 MiB (**2.7× lower**) |
| **`single-files`** | 1,000 standalone CLI invocations (`SIMPLE`) | **103.84 s** | 114.69 s *(batch)* | **1.10×** *(vs. in-JVM batch)* | **98 MiB** vs. 1,577 MiB (**16.1× lower**) |

Across all workloads, `jsgo` is **2.1×–5.6× faster than cold Java**, **1.3×–2.4× faster than a JIT-warmed JVM**, and uses **1.4×–16.1× less memory**. See [docs/benchmarks.md](docs/benchmarks.md) and [docs/performance-roadmap.md](docs/performance-roadmap.md) for full methodology, confidence intervals, and optimization details.
