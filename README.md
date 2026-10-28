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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `440a0571`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 205.1 s | — | **117.7 s**<br>1.74×±0.02 | 12.43 GiB | **5.75 GiB**<br>2.17×±0.19 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.4 s | 40.4 s | **23.1 s**<br>2.14×±0.03 | 4.99 GiB | **889 MiB**<br>5.74×±0.58 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.0 s | 26.6 s | **15.6 s**<br>2.24×±0.03 | 3.72 GiB | **722 MiB**<br>5.28×±0.17 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 23.0 s | 15.5 s | **9.18 s**<br>2.51×±0.02 | 3.44 GiB | **450 MiB**<br>7.82×±0.30 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.73 s | 2.80 s | **1.87 s**<br>3.60×±0.06 | 1.32 GiB | **269 MiB**<br>5.02×±0.21 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 12.9 s | 8.59 s | **4.13 s**<br>3.13×±0.03 | 1.63 GiB | **242 MiB**<br>6.87×±0.30 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.7 s | 7.07 s | **3.79 s**<br>3.34×±0.11 | 1.44 GiB | **216 MiB**<br>6.84×±0.14 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.46 s | 1.47 s | **0.88 s**<br>5.06×±0.11 | 600 MiB | **127 MiB**<br>4.74×±0.20 |

Across all whole-program and staged workloads, `jsgo` is **1.7×–5.8× faster** and uses **2.2×–7.9× less memory** than cold Java (and 1.2×–2.1× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#440a0571-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
