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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `7afbb2fe`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 206.7 s | — | **112.5 s**<br>1.84×±0.01 | 12.11 GiB | **5.54 GiB**<br>2.19×±0.18 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.3 s | 40.9 s | **21.8 s**<br>2.27×±0.03 | 4.92 GiB | **901 MiB**<br>5.58×±0.70 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.4 s | 25.9 s | **14.4 s**<br>2.45×±0.05 | 3.71 GiB | **715 MiB**<br>5.32×±0.21 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 23.1 s | 15.6 s | **8.37 s**<br>2.75×±0.07 | 3.43 GiB | **428 MiB**<br>8.21×±0.39 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.86 s | 2.75 s | **1.78 s**<br>3.84×±0.09 | 1.36 GiB | **276 MiB**<br>4.99×±0.89 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 12.9 s | 8.59 s | **3.94 s**<br>3.28×±0.07 | 1.61 GiB | **236 MiB**<br>7.02×±0.37 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.4 s | 6.60 s | **3.44 s**<br>3.62×±0.10 | 1.42 GiB | **212 MiB**<br>6.87×±0.20 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.48 s | 1.48 s | **0.85 s**<br>5.26×±0.15 | 605 MiB | **128 MiB**<br>4.71×±0.16 |

Across all whole-program and staged workloads, `jsgo` is **1.8×–6.2× faster** and uses **2.2×–8.2× less memory** than cold Java (and 1.3×–2.2× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#7afbb2fe-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
