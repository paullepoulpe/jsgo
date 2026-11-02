# jsgo

A Go port of [Google Closure Compiler](https://github.com/google/closure-compiler), aiming to be a drop-in replacement with the same flags, byte-identical output, and identical diagnostics.

The port is a whole-file, 1:1 translation of the pinned Java release (**v20261021**; originally ported from v20260415 and upgraded through v20261015), built bottom-up and verified layer by layer against the Java compiler on large real-world corpora.

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

## Benchmarks (`jsgo` vs. Java Closure Compiler v20261021)

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `58e97f06`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 200.7 s | — | **100.0 s**<br>2.01×±0.02 | 14.37 GiB | **5.43 GiB**<br>2.64×±0.14 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.3 s | 40.2 s | **19.0 s**<br>2.59×±0.03 | 4.61 GiB | **902 MiB**<br>5.22×±0.57 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.2 s | 26.2 s | **12.9 s**<br>2.74×±0.02 | 3.71 GiB | **725 MiB**<br>5.24×±0.09 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 23.0 s | 15.6 s | **7.45 s**<br>3.08×±0.03 | 3.44 GiB | **420 MiB**<br>8.38×±0.07 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.59 s | 2.79 s | **1.65 s**<br>3.99×±0.09 | 1.30 GiB | **253 MiB**<br>5.27×±0.28 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 13.0 s | 8.64 s | **3.53 s**<br>3.68×±0.03 | 1.61 GiB | **243 MiB**<br>6.80×±0.23 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.3 s | 6.70 s | **3.11 s**<br>3.95×±0.07 | 1.43 GiB | **212 MiB**<br>6.90×±0.17 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.50 s | 1.47 s | **0.76 s**<br>5.94×±0.27 | 612 MiB | **119 MiB**<br>5.13×±0.09 |

Across all whole-program and staged workloads, `jsgo` is **2.0×–6.7× faster** and uses **2.6×–8.4× less memory** than cold Java (and 1.4×–2.5× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#58e97f06-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
