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

Measured on macOS arm64 (Go 1.27.1 with PGO enabled via [`cmd/jsgo/default.pgo`](cmd/jsgo/default.pgo), OpenJDK 21.0.10). Each benchmark runs **5 warmup iterations** followed by **10 Latin-square-interleaved measurement rounds** (`Java Cold`, `jsgo t=1`, `Java Warm`, `jsgo t=4`). Values below report the **median `[95% CI]`** across the 10 rounds; speed ratios above `1.00×` mean `jsgo` is faster than Java:

### Standard Benchmark Suite ([`tools/bench/tasks/`](tools/bench/tasks/))

| Benchmark | Mode | `jsgo` (`t=1` / `t=4`) | Java Cold | Java Warm | `jsgo` vs Cold Java (`t=1` / `t=4`) | `jsgo` vs Warm Java (`t=1` / `t=4`) | Peak RSS: `jsgo` (`t=4`) vs Cold / Warm Java |
|---|---|---:|---:|---:|---:|---:|---:|
| **`typed-app`** *(Closure Library + app)* | `ADVANCED` + types + map | **1.53 s** `[1.53, 1.54]` / **1.41 s** `[1.41, 1.41]` | 6.24 s `[6.17, 6.29]` | 1.99 s `[1.97, 2.00]` | **4.07×** / **4.43×** `[4.37, 4.45]` | **1.30×** / **1.41×** `[1.39, 1.42]` | **252 MiB** vs 749 / 1,212 MiB *(peak 1,875)* (**2.97× / 4.94× lower**, **7.44×** vs JVM peak) |
| **`bundle`** *(12.6 MB: three + lodash + react-dom + TS)* | `SIMPLE` | **4.75 s** `[4.74, 4.76]` / **4.62 s** `[4.61, 4.64]` | 12.54 s `[12.47, 12.72]` | 8.97 s `[8.94, 9.03]` | **2.64×** / **2.72×** `[2.70, 2.75]` | **1.89×** / **1.94×** `[1.93, 1.95]` | **563 MiB** vs 2,143 / 2,900 MiB *(peak 3,633)* (**3.80× / 5.17× lower**, **6.45×** vs JVM peak) |
| **`downlevel-whitespace`** *(Angular + RxJS + core-js → ES5)* | `WHITESPACE_ONLY` | **0.57 s** `[0.57, 0.57]` / **0.48 s** `[0.48, 0.48]` | 2.41 s `[2.38, 2.45]` | 0.56 s `[0.55, 0.58]` | **4.23×** / **5.03×** `[4.98, 5.13]` | **0.98×** / **1.18×** `[1.17, 1.20]` | **171 MiB** vs 555 / 548 MiB *(peak 1,829)* (**3.08× / 3.10× lower**, **10.70×** vs JVM peak) |
| **`downlevel-simple`** *(Angular + RxJS + core-js → ES5)* | `SIMPLE` | **1.18 s** `[1.18, 1.19]` / **1.09 s** `[1.09, 1.10]` | 4.65 s `[4.55, 4.71]` | 1.53 s `[1.53, 1.54]` | **3.93×** / **4.25×** `[4.16, 4.31]` | **1.30×** / **1.40×** `[1.39, 1.41]` | **205 MiB** vs 619 / 912 MiB *(peak 1,889)* (**3.05× / 4.46× lower**, **9.22×** vs JVM peak) |
| **`checks-only`** *(Closure Library, 971 files)* | `--checks_only` + types | **2.68 s** `[2.67, 2.70]` / **2.30 s** `[2.30, 2.30]` | 6.92 s `[6.87, 7.01]` | 3.07 s `[3.06, 3.09]` | **2.59×** / **3.01×** `[2.99, 3.04]` | **1.15×** / **1.34×** `[1.33, 1.34]` | **618 MiB** vs 1,622 / 1,322 MiB *(peak 4,385)* (**2.65× / 2.14× lower**, **7.09×** vs JVM peak) |

### Large & Mega Real-World `ADVANCED` Workloads

| Benchmark | Input Size & Workload | `jsgo` (`t=1` / `t=4`) | Java Cold | Java Warm | `jsgo` vs Cold Java (`t=1` / `t=4`) | `jsgo` vs Warm Java (`t=1` / `t=4`) | Peak RSS: `jsgo` (`t=4`) vs Cold / Warm Java |
|---|---|---:|---:|---:|---:|---:|---:|
| **`typed-mega`** | **44.6 MB, 4,385 files** *(10 Closure `ADVANCED` codebases: Closure Library + NGA OpenSphere 23.0.1 + OpenLayers 4.6.5 + Shaka Player 5.3.1 + Blockly 8.0.5 + Camptocamp ngeo 2.2.10 + YouTube SPFjs 2.4.0 + ol-cesium 1.37.0 + google-protobuf 3.21.2 + app, 94.1% typed → 14.33 MB JS + 22.29 MB map)* | **20.48 s** `[20.46, 20.52]` / **19.34 s** `[19.32, 19.38]` | 35.43 s `[35.04, 35.55]` | 27.45 s `[27.30, 27.59]` | **1.73×** / **1.83×** `[1.81, 1.84]` | **1.34×** / **1.42×** `[1.41, 1.43]` | **1,884 MiB** vs 6,581 / 5,957 MiB *(peak 7,113)* (**3.55× / 3.23× lower**, **3.78×** vs JVM peak) |
| **`bundle-mega`** | **36.1 MB, 26 UMD libraries** *(plotly.js, handsontable, echarts, ag-grid-community, pdfjs-dist, three.js, react-dom, pixi.js, fabric, mathjs, konva, moment, d3, google-libphonenumber, lodash, lovefield, vue, jquery, rxjs, google-protobuf, handlebars, immutable, chart.js, zone.js, underscore, 27,597 warnings → 13.09 MB JS + 33.99 MB map)* | **25.73 s** `[25.68, 25.84]` / **24.92 s** `[24.84, 24.98]` | 48.26 s `[48.12, 48.52]` | 41.20 s `[41.01, 42.88]` | **1.87×** / **1.93×** `[1.93, 1.94]` | **1.60×** / **1.65×** `[1.64, 1.72]` | **2,421 MiB** vs 6,666 / 6,333 MiB *(peak 7,485)* (**2.77× / 2.58× lower**, **3.09×** vs JVM peak) |
| **`typed-full`** | **28.5 MB, 1,818 files** *(Closure Library + Google Shaka Player 5.3.1 + Google Blockly 8.0.5 + app, 97.3% typed → 6.75 MB JS + 13.31 MB map)* | **11.80 s** `[11.73, 11.83]` / **11.17 s** `[11.15, 11.20]` | 23.02 s `[22.91, 23.11]` | 16.10 s `[16.00, 16.17]` | **1.95×** / **2.06×** `[2.05, 2.07]` | **1.36×** / **1.44×** `[1.43, 1.45]` | **1,106 MiB** vs 4,068 / 3,325 MiB *(peak 4,625)* (**3.65× / 3.13× lower**, **4.18×** vs JVM peak) |
| **`bundle-advanced-ts`** | **12.6 MB, 4 libraries incl. 8.9 MB `typescript.js`** *(`ADVANCED` + 10,751 warnings → 3.25 MB JS + 9.39 MB map)* | **7.29 s** `[7.28, 7.30]` / **7.17 s** `[7.16, 7.19]` | 17.70 s `[17.62, 17.82]` | 11.61 s `[11.55, 11.75]` | **2.42×** / **2.47×** `[2.46, 2.48]` | **1.59×** / **1.62×** `[1.61, 1.64]` | **623 MiB** vs 2,358 / 3,278 MiB *(peak 4,821)* (**3.77× / 5.30× lower**, **7.73×** vs JVM peak) |
| **`shaka-blockly`** | **8.6 MB, 847 files** *(Google Shaka Player 5.3.1 + Google Blockly 8.0.5, 93.9% typed → 1.86 MB JS + 4.94 MB map)* | **4.91 s** `[4.90, 4.92]` / **4.70 s** `[4.69, 4.73]` | 12.50 s `[12.36, 12.56]` | 6.97 s `[6.95, 7.05]` | **2.54×** / **2.64×** `[2.63, 2.68]` | **1.42×** / **1.48×** `[1.47, 1.50]` | **524 MiB** vs 1,809 / 2,861 MiB *(peak 4,586)* (**3.50× / 5.48× lower**, **8.75×** vs JVM peak) |
| **`bundle-advanced`** | **5.4 MB, 11 UMD libraries** *(three, lodash, react-dom, rxjs, d3, moment, vue, jquery, handlebars, immutable, underscore → 1.88 MB JS + 4.35 MB map)* | **3.71 s** `[3.70, 3.72]` / **3.54 s** `[3.53, 3.54]` | 11.35 s `[11.25, 11.45]` | 5.79 s `[5.75, 5.82]` | **3.07×** / **3.21×** `[3.18, 3.23]` | **1.56×** / **1.63×** `[1.63, 1.64]` | **394 MiB** vs 2,291 / 2,037 MiB *(peak 2,872)* (**5.75× / 5.21× lower**, **7.30×** vs JVM peak) |

Across every benchmark—from sub-second down-levelling to 44.6 MB / 4,385-file and 36.1 MB / 26-library whole-program `ADVANCED` compilations—`jsgo` (`t=4`) is **1.83×–5.03× faster than cold Java**, **1.18×–1.94× faster than a 5-warmup JIT-compiled JVM**, and uses **2.14×–5.48× less per-iteration RSS** (**3.09×–10.70× less** than JVM peak RSS).
