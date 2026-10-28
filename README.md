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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `030244a9`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 209.1 s | — | **119.5 s**<br>1.75×±0.04 | 12.28 GiB | **5.64 GiB**<br>2.18×±0.09 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.3 s | 40.3 s | **23.3 s**<br>2.12×±0.02 | 4.79 GiB | **898 MiB**<br>5.46×±0.39 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.2 s | 26.2 s | **15.9 s**<br>2.21×±0.03 | 3.70 GiB | **736 MiB**<br>5.16×±0.14 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 22.9 s | 15.4 s | **9.31 s**<br>2.46×±0.02 | 3.44 GiB | **437 MiB**<br>8.07×±0.15 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.76 s | 2.77 s | **1.96 s**<br>3.46×±0.10 | 1.33 GiB | **280 MiB**<br>4.85×±0.34 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 13.1 s | 8.59 s | **4.34 s**<br>3.02×±0.06 | 1.61 GiB | **230 MiB**<br>7.17×±0.16 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.4 s | 6.71 s | **3.79 s**<br>3.28×±0.01 | 1.43 GiB | **219 MiB**<br>6.70×±0.11 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.48 s | 1.53 s | **0.90 s**<br>4.97×±0.10 | 601 MiB | **127 MiB**<br>4.73×±0.17 |

Across all whole-program and staged workloads, `jsgo` is **1.7×–5.6× faster** and uses **2.2×–8.1× less memory** than cold Java (and 1.2×–2.0× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#030244a9-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
