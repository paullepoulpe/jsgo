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

| Benchmark | Mode | `jsgo` (`t=1` / `t=4`) | Java Cold | Java Warm | `jsgo` vs Cold Java (`t=1` / `t=4`) | `jsgo` vs Warm Java (`t=1` / `t=4`) | Peak RSS: `jsgo` vs Cold / Warm Java |
|---|---|---:|---:|---:|---:|---:|---:|
| **`typed-app`** *(Closure Library + app)* | `ADVANCED` + types + map | **1.57 s** / **1.44 s** | 6.14 s | 2.18 s | **3.91×** / **4.26× faster** | **1.39×** / **1.51× faster** | **256–260 MiB** vs 758 / 1,871 MiB (**2.9× / 7.2× lower**) |
| **`bundle`** *(12.6 MB: three + lodash + react-dom + TS)* | `SIMPLE` | **4.89 s** / **4.82 s** | 16.40 s | 8.68 s | **3.35×** / **3.40× faster** | **1.78×** / **1.80× faster** | **577 MiB** vs 2,151 / 2,343 MiB (**3.7× / 4.1× lower**) |
| **`downlevel-whitespace`** *(Angular + RxJS + core-js → ES5)* | `WHITESPACE_ONLY` | **0.58 s** / **0.48 s** | 2.55 s | 0.60 s | **4.40×** / **5.31× faster** | **1.03×** / **1.25× faster** | **165–187 MiB** vs 553 / 1,437 MiB (**3.0× / 8.7× lower**) |
| **`downlevel-simple`** *(Angular + RxJS + core-js → ES5)* | `SIMPLE` | **1.21 s** / **1.12 s** | 4.68 s | 1.60 s | **3.87×** / **4.18× faster** | **1.32×** / **1.43× faster** | **199–206 MiB** vs 613 / 1,489 MiB (**3.0× / 7.5× lower**) |
| **`checks-only`** *(Closure Library, 971 files)* | `--checks_only` + types | **2.73 s** / **2.33 s** | 6.55 s | 3.20 s | **2.40×** / **2.81× faster** | **1.17×** / **1.37× faster** | **611–640 MiB** vs 1,591 / 3,562 MiB (**2.5× / 5.6× lower**) |

### Large & Mega Real-World `ADVANCED` Workloads

| Benchmark | Input Size & Workload | `jsgo` (`t=1` / `t=4`) | Java Cold | Java Warm | `jsgo` vs Cold Java (`t=1` / `t=4`) | `jsgo` vs Warm Java (`t=1` / `t=4`) | Peak RSS: `jsgo` vs Cold / Warm Java |
|---|---|---:|---:|---:|---:|---:|---:|
| **`typed-mega`** | **44.6 MB, 4,385 files** *(10 Closure `ADVANCED` codebases: Closure Library + NGA OpenSphere 23.0.1 + OpenLayers 4.6.5 + Shaka Player 5.3.1 + Blockly 8.0.5 + Camptocamp ngeo 2.2.10 + YouTube SPFjs 2.4.0 + ol-cesium 1.37.0 + google-protobuf 3.21.2 + app, 94.1% typed → 14.33 MB JS + 22.29 MB map)* | **20.96 s** / **19.83 s** | 38.18 s | 29.53 s | **1.82×** / **1.93× faster** | **1.41×** / **1.49× faster** *(−9.70 s)* | **1,804–1,908 MiB** vs 6,608 / 6,881 MiB (**3.5× / 3.6× lower**) |
| **`bundle-mega`** | **36.1 MB, 26 UMD libraries** *(plotly.js, handsontable, echarts, ag-grid-community, pdfjs-dist, three.js, react-dom, pixi.js, fabric, mathjs, konva, moment, d3, google-libphonenumber, lodash, lovefield, vue, jquery, rxjs, google-protobuf, handlebars, immutable, chart.js, zone.js, underscore, 27,597 warnings → 13.09 MB JS + 33.99 MB map)* | **26.51 s** / **25.70 s** | 66.71 s | 53.60 s | **2.52×** / **2.60× faster** | **2.02×** / **2.09× faster** *(−27.90 s)* | **2,387–2,444 MiB** vs 4,421 / 6,775 MiB (**1.8× / 2.8× lower**) |
| **`typed-full`** | **28.5 MB, 1,818 files** *(Closure Library + Google Shaka Player 5.3.1 + Google Blockly 8.0.5 + app, 97.3% typed → 6.75 MB JS + 13.31 MB map)* | **12.13 s** / **11.48 s** | 31.60 s | 20.21 s | **2.61×** / **2.75× faster** | **1.67×** / **1.76× faster** *(−8.73 s)* | **1,036–1,057 MiB** vs 1,603 / 4,724 MiB (**1.5× / 4.5× lower**) |
| **`bundle-advanced-ts`** | **12.6 MB, 4 libraries incl. 8.9 MB `typescript.js`** *(`ADVANCED` + 10,751 warnings → 3.25 MB JS + 9.39 MB map)* | **7.54 s** / **7.41 s** | 20.78 s | 12.41 s | **2.76×** / **2.80× faster** | **1.65×** / **1.67× faster** *(−5.00 s)* | **633–677 MiB** vs 2,422 / 2,850+ MiB (**3.6× / 4.2× lower**) |
| **`shaka-blockly`** | **8.6 MB, 847 files** *(Google Shaka Player 5.3.1 + Google Blockly 8.0.5, 93.9% typed → 1.86 MB JS + 4.94 MB map)* | **5.08 s** / **4.84 s** | 15.26 s | 8.67 s | **3.00×** / **3.15× faster** | **1.71×** / **1.79× faster** *(−3.83 s)* | **513–528 MiB** vs 1,341 / 2,852 MiB (**2.5× / 5.4× lower**) |
| **`bundle-advanced`** | **5.4 MB, 11 UMD libraries** *(three, lodash, react-dom, rxjs, d3, moment, vue, jquery, handlebars, immutable, underscore → 1.88 MB JS + 4.35 MB map)* | **3.83 s** / **3.66 s** | 11.30 s | 6.79 s | **2.95×** / **3.09× faster** | **1.77×** / **1.86× faster** *(−3.13 s)* | **395–404 MiB** vs 1,600 / 3,795 MiB (**4.0× / 9.4× lower**) |

Across every benchmark—from sub-second down-levelling to 44.6 MB / 4,385-file and 36.1 MB / 26-library whole-program `ADVANCED` compilations—`jsgo` is **1.93×–5.31× faster than cold Java**, **1.25×–2.09× faster than a warm JIT-compiled JVM**, and uses **2.8×–9.4× less memory** than warm Java.
