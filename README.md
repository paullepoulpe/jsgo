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
| **`typed-app`** *(Closure Library + app)* | `ADVANCED` + types + map | **1.55 s** / **1.43 s** | 6.14 s | 2.18 s | **3.96×** / **4.29× faster** | **1.41×** / **1.52× faster** | **242–246 MiB** vs 758 / 1,871 MiB (**3.1× / 7.7× lower**) |
| **`bundle`** *(12.6 MB: three + lodash + react-dom + TS)* | `SIMPLE` | **4.81 s** / **4.70 s** | 16.40 s | 8.68 s | **3.41×** / **3.49× faster** | **1.80×** / **1.85× faster** | **558–597 MiB** vs 2,151 / 2,343 MiB (**3.6× / 4.2× lower**) |
| **`downlevel-whitespace`** *(Angular + RxJS + core-js → ES5)* | `WHITESPACE_ONLY` | **0.56 s** / **0.47 s** | 2.55 s | 0.60 s | **4.55×** / **5.43× faster** | **1.07×** / **1.28× faster** | **162–171 MiB** vs 553 / 1,437 MiB (**3.2× / 8.9× lower**) |
| **`downlevel-simple`** *(Angular + RxJS + core-js → ES5)* | `SIMPLE` | **1.20 s** / **1.10 s** | 4.68 s | 1.60 s | **3.90×** / **4.25× faster** | **1.33×** / **1.45× faster** | **205–211 MiB** vs 613 / 1,489 MiB (**2.9× / 7.3× lower**) |
| **`checks-only`** *(Closure Library, 971 files)* | `--checks_only` + types | **2.72 s** / **2.33 s** | 6.55 s | 3.20 s | **2.41×** / **2.81× faster** | **1.18×** / **1.37× faster** | **623–664 MiB** vs 1,591 / 3,562 MiB (**2.4× / 5.7× lower**) |

### Large & Mega Real-World `ADVANCED` Workloads

| Benchmark | Input Size & Workload | `jsgo` (`t=1` / `t=4`) | Java Cold | Java Warm | `jsgo` vs Cold Java (`t=1` / `t=4`) | `jsgo` vs Warm Java (`t=1` / `t=4`) | Peak RSS: `jsgo` vs Cold / Warm Java |
|---|---|---:|---:|---:|---:|---:|---:|
| **`typed-mega`** | **44.6 MB, 4,385 files** *(10 Closure `ADVANCED` codebases: Closure Library + NGA OpenSphere 23.0.1 + OpenLayers 4.6.5 + Shaka Player 5.3.1 + Blockly 8.0.5 + Camptocamp ngeo 2.2.10 + YouTube SPFjs 2.4.0 + ol-cesium 1.37.0 + google-protobuf 3.21.2 + app, 94.1% typed → 14.33 MB JS + 22.29 MB map)* | **20.74 s** / **19.61 s** | 38.18 s | 29.53 s | **1.84×** / **1.95× faster** | **1.42×** / **1.51× faster** *(−9.92 s)* | **1,821–1,864 MiB** vs 6,608 / 6,881 MiB (**3.5× / 3.8× lower**) |
| **`bundle-mega`** | **36.1 MB, 26 UMD libraries** *(plotly.js, handsontable, echarts, ag-grid-community, pdfjs-dist, three.js, react-dom, pixi.js, fabric, mathjs, konva, moment, d3, google-libphonenumber, lodash, lovefield, vue, jquery, rxjs, google-protobuf, handlebars, immutable, chart.js, zone.js, underscore, 27,597 warnings → 13.09 MB JS + 33.99 MB map)* | **26.15 s** / **25.39 s** | 66.71 s | 53.60 s | **2.55×** / **2.63× faster** | **2.05×** / **2.11× faster** *(−28.21 s)* | **2,288–2,410 MiB** vs 4,421 / 6,775 MiB (**1.8× / 3.0× lower**) |
| **`typed-full`** | **28.5 MB, 1,818 files** *(Closure Library + Google Shaka Player 5.3.1 + Google Blockly 8.0.5 + app, 97.3% typed → 6.75 MB JS + 13.31 MB map)* | **11.96 s** / **11.35 s** | 31.60 s | 20.21 s | **2.64×** / **2.78× faster** | **1.69×** / **1.78× faster** *(−8.86 s)* | **1,098–1,113 MiB** vs 1,603 / 4,724 MiB (**1.4× / 4.3× lower**) |
| **`bundle-advanced-ts`** | **12.6 MB, 4 libraries incl. 8.9 MB `typescript.js`** *(`ADVANCED` + 10,751 warnings → 3.25 MB JS + 9.39 MB map)* | **7.40 s** / **7.28 s** | 20.78 s | 12.41 s | **2.81×** / **2.85× faster** | **1.68×** / **1.70× faster** *(−5.13 s)* | **627–673 MiB** vs 2,422 / 2,850+ MiB (**3.6× / 4.5× lower**) |
| **`shaka-blockly`** | **8.6 MB, 847 files** *(Google Shaka Player 5.3.1 + Google Blockly 8.0.5, 93.9% typed → 1.86 MB JS + 4.94 MB map)* | **4.98 s** / **4.78 s** | 15.26 s | 8.67 s | **3.06×** / **3.19× faster** | **1.74×** / **1.81× faster** *(−3.89 s)* | **547–551 MiB** vs 1,341 / 2,852 MiB (**2.4× / 5.2× lower**) |
| **`bundle-advanced`** | **5.4 MB, 11 UMD libraries** *(three, lodash, react-dom, rxjs, d3, moment, vue, jquery, handlebars, immutable, underscore → 1.88 MB JS + 4.35 MB map)* | **3.77 s** / **3.60 s** | 11.30 s | 6.79 s | **3.00×** / **3.14× faster** | **1.80×** / **1.89× faster** *(−3.19 s)* | **389–395 MiB** vs 1,600 / 3,795 MiB (**4.1× / 9.8× lower**) |

Across every benchmark—from sub-second down-levelling to 44.6 MB / 4,385-file and 36.1 MB / 26-library whole-program `ADVANCED` compilations—`jsgo` is **1.95×–5.43× faster than cold Java**, **1.28×–2.11× faster than a warm JIT-compiled JVM**, and uses **3.0×–9.8× less memory** than warm Java.
