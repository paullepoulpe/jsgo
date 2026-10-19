# jsgo

A Go port of [Google Closure Compiler](https://github.com/google/closure-compiler),
aiming to be a drop-in replacement: same flags, same output, same diagnostics.

The port is a whole-file, 1:1 translation of the pinned Java release
(**v20261015**; ported from v20260415, then upgraded), built bottom-up and
verified layer by layer against the Java compiler itself on large corpora.
The pin is named in one place, [tools/oracle/pin.env](tools/oracle/pin.env).

- [PLAN.md](PLAN.md) — feasibility assessment, strategy and roadmap.
- [PORTING.md](PORTING.md) — the translation contract every ported file follows.
- [docs/benchmarks.md](docs/benchmarks.md) & [docs/performance-roadmap.md](docs/performance-roadmap.md) — end-to-end benchmarks and performance engineering history.
- [tools/oracle/](tools/oracle/) — tools that extract ground truth from the Java compiler.
- [testdata/](testdata/) — Closure type-check fixtures and their Java oracles (acceptance set for the type checker, M7).

```sh
make            # build, vet, test
make oracle-tools   # fetch the pinned jar/sources and build the Java oracle tools (needs a JDK)
```

## Benchmarks (`jsgo` vs Java Closure Compiler v20261015)

Every benchmark is gated on **100% byte-identical output** (compiled JS, source maps, stdout, stderr diagnostics, and exit status) against `closure-compiler-v20261015.jar`.

Measured on macOS arm64 (Go 1.27.1 with PGO enabled via [`cmd/jsgo/default.pgo`](cmd/jsgo/default.pgo), OpenJDK 21.0.10, interleaved runs). Speed ratios above `1.00×` mean `jsgo` is faster than Java:

### Standard Benchmark Suite ([`tools/bench/tasks/`](tools/bench/tasks/))

| Benchmark | Mode | `jsgo` (`t=1` / `t=4`) | Java Cold | Java Warm | `jsgo` vs Cold Java | `jsgo` vs Warm Java | Peak RSS: `jsgo` vs Cold / Warm Java |
|---|---|---:|---:|---:|---:|---:|---:|
| **`typed-app`** *(Closure Library + app)* | `ADVANCED` + types + map | **2.26 s** / **2.11 s** | 6.14 s | 2.18 s | **2.92× faster** | **1.04× faster** | **243 MiB** vs 758 / 1,871 MiB (**3.1× / 7.7× lower**) |
| **`bundle`** *(12.6 MB: three + lodash + react-dom + TS)* | `SIMPLE` | **7.90 s** / **7.76 s** | 16.40 s | 8.68 s | **2.11× faster** | **1.12× faster** | **522 MiB** vs 2,151 / 2,343 MiB (**4.1× / 4.5× lower**) |
| **`downlevel-whitespace`** *(Angular + RxJS + core-js → ES5)* | `WHITESPACE_ONLY` | **0.66 s** / **0.55 s** | 2.55 s | 0.60 s | **4.64× faster** | **1.09× faster** | **149 MiB** vs 553 / 1,437 MiB (**3.7× / 9.6× lower**) |
| **`downlevel-simple`** *(Angular + RxJS + core-js → ES5)* | `SIMPLE` | **1.71 s** / **1.60 s** | 4.68 s | 1.60 s | **2.93× faster** | **1.00× (tied)** | **197 MiB** vs 613 / 1,489 MiB (**3.1× / 7.6× lower**) |
| **`checks-only`** *(Closure Library, 971 files)* | `--checks_only` + types | **3.47 s** / **3.08 s** | 6.55 s | 3.20 s | **2.12× faster** | **1.04× faster** | **643 MiB** vs 1,591 / 3,562 MiB (**2.5× / 5.5× lower**) |

### Large Real-World `ADVANCED` Workloads

| Benchmark | Input Size & Workload | `jsgo` (`t=4`) | Java Cold | Java Warm | `jsgo` vs Cold Java | `jsgo` vs Warm Java | Peak RSS: `jsgo` vs Cold / Warm Java |
|---|---|---:|---:|---:|---:|---:|---:|
| **`typed-full`** | **28.5 MB, 1,818 files** *(Closure Library + Google Shaka Player 5.3.1 + Google Blockly 8.0.5 + app, 97.3% typed → 6.75 MB JS + 13.31 MB map)* | **17.17 s** | 31.60 s | 20.21 s | **1.84× faster** | **1.18× faster** *(−3.04 s)* | **1,189 MiB** vs 1,603 / 4,724 MiB (**1.3× / 4.0× lower**) |
| **`shaka-blockly`** | **8.6 MB, 847 files** *(Google Shaka Player 5.3.1 + Google Blockly 8.0.5, 93.9% typed → 1.86 MB JS + 4.94 MB map)* | **7.44 s** | 15.26 s | 8.67 s | **2.05× faster** | **1.17× faster** *(−1.23 s)* | **552 MiB** vs 1,341 / 2,852 MiB (**2.4× / 5.2× lower**) |
| **`bundle-advanced`** | **5.4 MB, 11 UMD libraries** *(three, lodash, react-dom, rxjs, d3, moment, vue, jquery, handlebars, immutable, underscore → 1.88 MB JS + 4.35 MB map)* | **6.05 s** | 11.30 s | 6.79 s | **1.87× faster** | **1.12× faster** *(−0.74 s)* | **417 MiB** vs 1,600 / 3,795 MiB (**3.8× / 9.1× lower**) |
| **`bundle-advanced-ts`** | **12.6 MB, 4 libraries incl. 8.9 MB `typescript.js`** *(`ADVANCED` + 10,751 warnings → 3.25 MB JS + 9.39 MB map)* | **12.09 s** | 20.78 s | 12.41 s | **1.72× faster** | **1.03× faster** *(−0.32 s)* | **670 MiB** vs 2,422 / 2,850+ MiB (**3.6× / 4.3× lower**) |

Across every benchmark—from sub-second down-levelling to 28.5 MB / 1,818-file whole-program `ADVANCED` compilations—`jsgo` is **1.72×–4.64× faster than cold Java**, **1.00×–1.18× faster than a warm JIT-compiled JVM**, and uses **4.0×–9.6× less memory** than warm Java.
