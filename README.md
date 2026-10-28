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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `a94169fa`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 207.5 s | — | **120.2 s**<br>1.73×±0.04 | 12.04 GiB | **5.95 GiB**<br>2.02×±0.08 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.7 s | 40.2 s | **23.6 s**<br>2.11×±0.04 | 4.99 GiB | **898 MiB**<br>5.69×±0.37 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.4 s | 26.0 s | **16.1 s**<br>2.20×±0.02 | 3.71 GiB | **725 MiB**<br>5.24×±0.13 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 23.1 s | 15.5 s | **9.45 s**<br>2.45×±0.05 | 3.44 GiB | **441 MiB**<br>7.99×±0.26 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.64 s | 2.76 s | **1.96 s**<br>3.38×±0.05 | 1.30 GiB | **274 MiB**<br>4.86×±0.26 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 12.9 s | 8.54 s | **4.35 s**<br>2.96×±0.03 | 1.62 GiB | **241 MiB**<br>6.91×±0.70 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.4 s | 6.69 s | **3.81 s**<br>3.26×±0.07 | 1.43 GiB | **218 MiB**<br>6.71×±0.27 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.50 s | 1.49 s | **0.91 s**<br>4.96×±0.19 | 602 MiB | **125 MiB**<br>4.84×±0.25 |

Across all whole-program and staged workloads, `jsgo` is **1.7×–5.7× faster** and uses **2.0×–8.5× less memory** than cold Java (and 1.2×–2.0× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#a94169fa-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
