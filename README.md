# jsgo

A Go port of [Google Closure Compiler](https://github.com/google/closure-compiler), aiming to be a drop-in replacement with the same flags, byte-identical output, and identical diagnostics.

The port is a whole-file, 1:1 translation of the pinned Java release (**v20261028**; originally ported from v20260415 and upgraded through v20261015 and v20261021), built bottom-up and verified layer by layer against the Java compiler on large real-world corpora.

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

## Benchmarks (`jsgo` vs. Java Closure Compiler v20261028)

Gated on **100% byte-identical output** (JS, `.i.js` type summaries, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on an Apple M2 (macOS arm64, 24 GiB RAM, `--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `75bca77c`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 208.7 s | — | **80.1 s**<br>2.61×±0.07 | 12.31 GiB | **5.60 GiB**<br>2.20×±0.06 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.6 s | 40.6 s | **15.7 s**<br>3.16×±0.03 | 4.69 GiB | **1,023 MiB**<br>4.69×±0.57 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.9 s | 26.4 s | **9.08 s**<br>3.95×±0.06 | 3.71 GiB | **761 MiB**<br>4.99×±0.28 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 23.2 s | 15.8 s | **5.19 s**<br>4.47×±0.07 | 3.43 GiB | **464 MiB**<br>7.58×±0.37 |
| **`conformance-checks`**<br>*CHECKS+types+conformance · 1.89M nodes* | 9.49 s | 4.60 s | **1.72 s**<br>5.53×±0.06 | 1.48 GiB | **378 MiB**<br>4.02×±0.31 |
| **`generate-ijs`**<br>*GENERATE_IJS · 1.89M nodes* | 3.14 s | 0.57 s | **0.35 s**<br>9.05×±0.52 | 1.03 GiB | **275 MiB**<br>3.84×±0.20 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.94 s | 2.73 s | **1.11 s**<br>6.26×±0.10 | 1.35 GiB | **297 MiB**<br>4.67×±0.18 |
| **`incremental-ijs-check`**<br>*CHECKS+types · 1,743 `.i.js` + 3 `.js`* | 4.01 s | 1.26 s | **0.58 s**<br>6.92×±0.17 | 880 MiB | **207 MiB**<br>4.25×±0.07 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 13.1 s | 8.70 s | **2.77 s**<br>4.73×±0.06 | 1.60 GiB | **252 MiB**<br>6.51×±0.34 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.8 s | 6.80 s | **2.23 s**<br>5.72×±0.06 | 1.44 GiB | **249 MiB**<br>5.90×±0.28 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.65 s | 1.52 s | **0.63 s**<br>7.37×±0.11 | 608 MiB | **125 MiB**<br>4.88×±0.12 |

Across all whole-program, Stage-1 (`.i.js` / conformance), and staged workloads, `jsgo` is **2.6×–9.0× faster** and uses **2.2×–7.6× less memory** than cold Java (and 1.5×–3.1× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#75bca77c-on-jdk-27-apple-m2) for all 16 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
