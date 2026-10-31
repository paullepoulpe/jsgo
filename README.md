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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `8f5b28b9`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 208.2 s | — | **101.9 s**<br>2.04×±0.04 | 12.12 GiB | **5.54 GiB**<br>2.19×±0.07 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.2 s | 40.3 s | **20.1 s**<br>2.44×±0.02 | 4.69 GiB | **904 MiB**<br>5.31×±0.39 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.2 s | 26.4 s | **13.6 s**<br>2.58×±0.05 | 3.70 GiB | **707 MiB**<br>5.36×±0.06 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 23.0 s | 15.4 s | **7.89 s**<br>2.91×±0.03 | 3.42 GiB | **419 MiB**<br>8.35×±0.24 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.66 s | 2.82 s | **1.65 s**<br>4.02×±0.07 | 1.34 GiB | **276 MiB**<br>4.99×±0.17 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 13.0 s | 8.88 s | **3.55 s**<br>3.67×±0.05 | 1.62 GiB | **238 MiB**<br>6.97×±0.52 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.4 s | 6.83 s | **3.24 s**<br>3.82×±0.07 | 1.43 GiB | **218 MiB**<br>6.72×±0.22 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.47 s | 1.49 s | **0.77 s**<br>5.78×±0.06 | 606 MiB | **121 MiB**<br>5.01×±0.12 |

Across all whole-program and staged workloads, `jsgo` is **2.0×–6.6× faster** and uses **2.2×–8.6× less memory** than cold Java (and 1.4×–2.5× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#8f5b28b9-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
