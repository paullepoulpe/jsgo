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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `79777d0a`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 205.4 s | — | **110.6 s**<br>1.86×±0.04 | 12.52 GiB | **5.61 GiB**<br>2.23×±0.15 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.4 s | 40.3 s | **21.3 s**<br>2.31×±0.03 | 4.84 GiB | **879 MiB**<br>5.64×±0.29 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.0 s | 26.4 s | **14.2 s**<br>2.48×±0.02 | 3.70 GiB | **718 MiB**<br>5.27×±0.15 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 23.0 s | 15.4 s | **8.21 s**<br>2.80×±0.06 | 3.43 GiB | **422 MiB**<br>8.33×±0.35 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.88 s | 2.79 s | **1.75 s**<br>3.93×±0.20 | 1.35 GiB | **268 MiB**<br>5.17×±0.42 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 12.8 s | 8.54 s | **3.82 s**<br>3.36×±0.09 | 1.62 GiB | **241 MiB**<br>6.89×±0.29 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.6 s | 6.87 s | **3.38 s**<br>3.72×±0.06 | 1.44 GiB | **208 MiB**<br>7.05×±0.20 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.46 s | 1.48 s | **0.83 s**<br>5.36×±0.18 | 601 MiB | **126 MiB**<br>4.79×±0.18 |

Across all whole-program and staged workloads, `jsgo` is **1.9×–6.2× faster** and uses **2.2×–8.8× less memory** than cold Java (and 1.3×–2.2× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#79777d0a-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
