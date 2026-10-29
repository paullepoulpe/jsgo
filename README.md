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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `ef4d1df4`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 205.7 s | — | **115.4 s**<br>1.78×±0.03 | 12.39 GiB | **5.57 GiB**<br>2.23×±0.07 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.5 s | 40.3 s | **22.3 s**<br>2.22×±0.01 | 4.62 GiB | **887 MiB**<br>5.33×±0.68 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.2 s | 26.3 s | **14.8 s**<br>2.37×±0.02 | 3.70 GiB | **735 MiB**<br>5.16×±0.19 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 23.0 s | 15.4 s | **8.59 s**<br>2.68×±0.06 | 3.45 GiB | **427 MiB**<br>8.28×±0.29 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.80 s | 2.77 s | **1.82 s**<br>3.73×±0.05 | 1.32 GiB | **278 MiB**<br>4.88×±0.56 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 12.9 s | 8.59 s | **4.03 s**<br>3.20×±0.06 | 1.61 GiB | **231 MiB**<br>7.12×±0.16 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.5 s | 6.69 s | **3.56 s**<br>3.51×±0.09 | 1.44 GiB | **213 MiB**<br>6.95×±0.10 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.50 s | 1.50 s | **0.86 s**<br>5.23×±0.19 | 603 MiB | **124 MiB**<br>4.87×±0.17 |

Across all whole-program and staged workloads, `jsgo` is **1.8×–5.9× faster** and uses **2.2×–8.4× less memory** than cold Java (and 1.2×–2.1× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#ef4d1df4-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
