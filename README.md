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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `f60f193b`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 206.2 s | — | **117.0 s**<br>1.76×±0.02 | 12.44 GiB | **5.70 GiB**<br>2.18×±0.15 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.3 s | 40.2 s | **23.0 s**<br>2.14×±0.02 | 4.72 GiB | **896 MiB**<br>5.39×±0.41 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.0 s | 26.3 s | **15.4 s**<br>2.26×±0.03 | 3.70 GiB | **730 MiB**<br>5.19×±0.13 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 22.9 s | 15.5 s | **9.02 s**<br>2.54×±0.03 | 3.42 GiB | **442 MiB**<br>7.94×±0.19 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.68 s | 2.86 s | **1.88 s**<br>3.56×±0.09 | 1.29 GiB | **270 MiB**<br>4.89×±0.47 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 12.9 s | 8.60 s | **4.07 s**<br>3.16×±0.06 | 1.61 GiB | **240 MiB**<br>6.87×±0.47 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.4 s | 6.79 s | **3.70 s**<br>3.36×±0.06 | 1.42 GiB | **215 MiB**<br>6.78×±0.34 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.50 s | 1.51 s | **0.88 s**<br>5.12×±0.13 | 602 MiB | **126 MiB**<br>4.79×±0.16 |

Across all whole-program and staged workloads, `jsgo` is **1.8×–5.9× faster** and uses **2.2×–8.2× less memory** than cold Java (and 1.2×–2.1× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#f60f193b-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
