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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `0223e8e5`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 200.2 s | — | **99.8 s**<br>2.01×±0.03 | 14.29 GiB | **5.48 GiB**<br>2.61×±0.05 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.2 s | 39.9 s | **18.9 s**<br>2.60×±0.02 | 4.87 GiB | **901 MiB**<br>5.53×±0.40 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.1 s | 26.6 s | **12.8 s**<br>2.75×±0.01 | 3.70 GiB | **719 MiB**<br>5.27×±0.30 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 22.9 s | 15.5 s | **7.39 s**<br>3.09×±0.01 | 3.42 GiB | **424 MiB**<br>8.27×±0.14 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.69 s | 2.76 s | **1.64 s**<br>4.07×±0.09 | 1.31 GiB | **269 MiB**<br>4.99×±0.83 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 12.9 s | 8.60 s | **3.52 s**<br>3.67×±0.03 | 1.61 GiB | **246 MiB**<br>6.70×±0.34 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.4 s | 6.69 s | **3.08 s**<br>4.02×±0.03 | 1.43 GiB | **214 MiB**<br>6.86×±0.21 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.52 s | 1.48 s | **0.75 s**<br>6.01×±0.13 | 607 MiB | **120 MiB**<br>5.06×±0.09 |

Across all whole-program and staged workloads, `jsgo` is **2.0×–6.7× faster** and uses **2.6×–8.3× less memory** than cold Java (and 1.4×–2.4× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#0223e8e5-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
