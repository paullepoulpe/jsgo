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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `233d5db4`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 205.1 s | — | **123.7 s**<br>1.66×±0.01 | 12.86 GiB | **5.31 GiB**<br>2.42×±0.21 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.5 s | 40.8 s | **25.0 s**<br>1.98×±0.02 | 4.82 GiB | **921 MiB**<br>5.37×±0.39 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.4 s | 26.6 s | **17.0 s**<br>2.09×±0.02 | 3.71 GiB | **725 MiB**<br>5.24×±0.22 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 23.2 s | 15.4 s | **10.1 s**<br>2.31×±0.05 | 3.44 GiB | **442 MiB**<br>7.98×±0.44 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.87 s | 2.82 s | **1.96 s**<br>3.50×±0.12 | 1.29 GiB | **259 MiB**<br>5.11×±0.41 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 13.1 s | 8.86 s | **4.40 s**<br>2.97×±0.08 | 1.60 GiB | **234 MiB**<br>6.99×±0.19 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.5 s | 6.76 s | **4.14 s**<br>3.02×±0.09 | 1.43 GiB | **205 MiB**<br>7.15×±0.09 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.47 s | 1.54 s | **0.93 s**<br>4.78×±0.15 | 605 MiB | **127 MiB**<br>4.77×±0.12 |

Across all whole-program and staged workloads, `jsgo` is **1.7×–5.6× faster** and uses **2.4×–8.0× less memory** than cold Java (and 1.2×–2.0× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#233d5db4-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
