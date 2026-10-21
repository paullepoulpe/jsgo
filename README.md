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
| **`typed-app`** *(Closure Library + app)* | `ADVANCED` + types + map | **1.35 s** / **1.23 s** *(`1.22 s` `--fast_mode`)* | 6.24 s `[6.17, 6.29]` | 1.99 s `[1.97, 2.00]` | **4.62×** / **5.07×** *(`5.10×`)* | **1.47×** / **1.62×** *(`1.63×`)* | **264 MiB** vs 749 / 1,212 MiB *(peak 1,875)* (**2.83× / 4.58× lower**, **7.09×** vs JVM peak) |
| **`bundle`** *(12.6 MB: three + lodash + react-dom + TS)* | `SIMPLE` | **4.03 s** / **3.75 s** | 12.54 s `[12.47, 12.72]` | 8.97 s `[8.94, 9.03]` | **3.11×** / **3.34×** | **2.23×** / **2.39×** | **559 MiB** vs 2,143 / 2,900 MiB *(peak 3,633)* (**3.83× / 5.19× lower**, **6.50×** vs JVM peak) |
| **`downlevel-whitespace`** *(Angular + RxJS + core-js → ES5)* | `WHITESPACE_ONLY` | **0.52 s** / **0.43 s** | 2.41 s `[2.38, 2.45]` | 0.56 s `[0.55, 0.58]` | **4.63×** / **5.60×** | **1.08×** / **1.30×** | **155 MiB** vs 555 / 548 MiB *(peak 1,829)* (**3.58× / 3.54× lower**, **11.81×** vs JVM peak) |
| **`downlevel-simple`** *(Angular + RxJS + core-js → ES5)* | `SIMPLE` | **1.01 s** / **0.90 s** | 4.65 s `[4.55, 4.71]` | 1.53 s `[1.53, 1.54]` | **4.60×** / **5.17×** | **1.51×** / **1.70×** | **204 MiB** vs 619 / 912 MiB *(peak 1,889)* (**3.04× / 4.48× lower**, **9.28×** vs JVM peak) |
| **`checks-only`** *(Closure Library, 971 files)* | `--checks_only` + types | **2.39 s** / **2.00 s** | 6.92 s `[6.87, 7.01]` | 3.07 s `[3.06, 3.09]` | **2.90×** / **3.46×** | **1.28×** / **1.53×** | **519 MiB** vs 1,622 / 1,322 MiB *(peak 4,385)* (**3.12× / 2.55× lower**, **8.44×** vs JVM peak) |

### Large & Mega Real-World `ADVANCED` Workloads

| Benchmark | Input Size & Workload | `jsgo` (`t=1` / `t=4`) | Java Cold | Java Warm | `jsgo` vs Cold Java (`t=1` / `t=4`) | `jsgo` vs Warm Java (`t=1` / `t=4`) | Peak RSS: `jsgo` (`t=4`) vs Cold / Warm Java |
|---|---|---:|---:|---:|---:|---:|---:|
| **`typed-mega`** | **44.6 MB, 4,385 files** *(10 Closure `ADVANCED` codebases: Closure Library + NGA OpenSphere 23.0.1 + OpenLayers 4.6.5 + Shaka Player 5.3.1 + Blockly 8.0.5 + Camptocamp ngeo 2.2.10 + YouTube SPFjs 2.4.0 + ol-cesium 1.37.0 + google-protobuf 3.21.2 + app, 94.1% typed → 14.33 MB JS + 22.29 MB map)* | **17.96 s** / **16.72 s** *(`16.13 s` `--fast_mode`)* | 35.43 s `[35.04, 35.55]` | 27.45 s `[27.30, 27.59]` | **1.97×** / **2.12×** *(`2.20×`)* | **1.53×** / **1.64×** *(`1.70×`)* | **1,773 MiB** vs 6,581 / 5,957 MiB *(peak 7,113)* (**3.71× / 3.36× lower**, **4.01×** vs JVM peak) |
| **`bundle-mega`** | **36.1 MB, 26 UMD libraries** *(plotly.js, handsontable, echarts, ag-grid-community, pdfjs-dist, three.js, react-dom, pixi.js, fabric, mathjs, konva, moment, d3, google-libphonenumber, lodash, lovefield, vue, jquery, rxjs, google-protobuf, handlebars, immutable, chart.js, zone.js, underscore, 27,597 warnings → 13.09 MB JS + 33.99 MB map)* | **21.67 s** / **20.27 s** *(`20.22 s` `--fast_mode`)* | 48.26 s `[48.12, 48.52]` | 41.20 s `[41.01, 42.88]` | **2.23×** / **2.38×** *(`2.39×`)* | **1.90×** / **2.03×** *(`2.04×`)* | **1,746 MiB** vs 6,666 / 6,333 MiB *(peak 7,485)* (**3.82× / 3.63× lower**, **4.29×** vs JVM peak) |
| **`typed-full`** | **28.5 MB, 1,818 files** *(Closure Library + Google Shaka Player 5.3.1 + Google Blockly 8.0.5 + app, 97.3% typed → 6.75 MB JS + 13.31 MB map)* | **10.25 s** / **9.50 s** *(`9.23 s` `--fast_mode`)* | 23.02 s `[22.91, 23.11]` | 16.10 s `[16.00, 16.17]` | **2.25×** / **2.42×** *(`2.49×`)* | **1.57×** / **1.69×** *(`1.74×`)* | **1,113 MiB** vs 4,068 / 3,325 MiB *(peak 4,625)* (**3.66× / 2.99× lower**, **4.16×** vs JVM peak) |
| **`bundle-advanced-ts`** | **12.6 MB, 4 libraries incl. 8.9 MB `typescript.js`** *(`ADVANCED` + 10,751 warnings → 3.25 MB JS + 9.39 MB map)* | **6.16 s** / **5.86 s** *(`5.85 s` `--fast_mode`)* | 17.70 s `[17.62, 17.82]` | 11.61 s `[11.55, 11.75]` | **2.87×** / **3.02×** *(`3.03×`)* | **1.88×** / **1.98×** *(`1.99×`)* | **587 MiB** vs 2,358 / 3,278 MiB *(peak 4,821)* (**4.02× / 5.59× lower**, **8.22×** vs JVM peak) |
| **`shaka-blockly`** | **8.6 MB, 847 files** *(Google Shaka Player 5.3.1 + Google Blockly 8.0.5, 93.9% typed → 1.86 MB JS + 4.94 MB map)* | **4.18 s** / **3.88 s** *(`3.75 s` `--fast_mode`)* | 12.50 s `[12.36, 12.56]` | 6.97 s `[6.95, 7.05]` | **2.99×** / **3.22×** *(`3.33×`)* | **1.67×** / **1.80×** *(`1.86×`)* | **517 MiB** vs 1,809 / 2,861 MiB *(peak 4,586)* (**3.50× / 5.53× lower**, **8.87×** vs JVM peak) |
| **`bundle-advanced`** | **5.4 MB, 11 UMD libraries** *(three, lodash, react-dom, rxjs, d3, moment, vue, jquery, handlebars, immutable, underscore → 1.88 MB JS + 4.35 MB map)* | **3.14 s** / **2.91 s** | 11.35 s `[11.25, 11.45]` | 5.79 s `[5.75, 5.82]` | **3.61×** / **3.90×** | **1.84×** / **1.99×** | **312 MiB** vs 2,291 / 2,037 MiB *(peak 2,872)* (**7.33× / 6.52× lower**, **9.19×** vs JVM peak) |

Across every benchmark—from sub-second down-levelling to 44.6 MB / 4,385-file and 36.1 MB / 26-library whole-program `ADVANCED` compilations—`jsgo` (`t=4`) is **2.12×–5.60× faster than cold Java** (**2.20×–5.60×** with `--fast_mode`), **1.30×–2.39× faster than a 5-warmup JIT-compiled JVM**, and uses **2.55×–7.33× less per-iteration RSS** (**4.01×–11.81× less** than JVM peak RSS).
