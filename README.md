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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `33574cef`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 204.4 s | — | **115.2 s**<br>1.77×±0.02 | 12.52 GiB | **5.70 GiB**<br>2.20×±0.19 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.3 s | 40.2 s | **22.3 s**<br>2.22×±0.02 | 4.58 GiB | **879 MiB**<br>5.33×±0.37 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.1 s | 26.3 s | **14.8 s**<br>2.37×±0.05 | 3.71 GiB | **718 MiB**<br>5.29×±0.22 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 23.1 s | 15.6 s | **8.59 s**<br>2.68×±0.03 | 3.42 GiB | **430 MiB**<br>8.15×±0.31 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.96 s | 2.80 s | **1.83 s**<br>3.80×±0.12 | 1.36 GiB | **282 MiB**<br>4.93×±0.30 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 12.9 s | 8.51 s | **4.01 s**<br>3.22×±0.05 | 1.62 GiB | **231 MiB**<br>7.19×±0.60 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.5 s | 6.64 s | **3.54 s**<br>3.53×±0.08 | 1.43 GiB | **215 MiB**<br>6.84×±0.31 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.49 s | 1.51 s | **0.86 s**<br>5.20×±0.20 | 604 MiB | **126 MiB**<br>4.78×±0.05 |

Across all whole-program and staged workloads, `jsgo` is **1.8×–6.0× faster** and uses **2.2×–8.3× less memory** than cold Java (and 1.3×–2.1× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#33574cef-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
