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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `0a1a72c1`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 208.1 s | — | **101.8 s**<br>2.04×±0.06 | 12.33 GiB | **5.57 GiB**<br>2.22×±0.16 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.4 s | 40.0 s | **20.3 s**<br>2.44×±0.02 | 4.77 GiB | **881 MiB**<br>5.55×±0.33 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.3 s | 26.4 s | **13.7 s**<br>2.58×±0.03 | 3.73 GiB | **715 MiB**<br>5.34×±0.15 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 23.1 s | 15.6 s | **7.91 s**<br>2.92×±0.05 | 3.44 GiB | **421 MiB**<br>8.38×±0.35 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.89 s | 2.84 s | **1.66 s**<br>4.16×±0.06 | 1.32 GiB | **268 MiB**<br>5.06×±0.38 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 12.9 s | 8.56 s | **3.57 s**<br>3.62×±0.04 | 1.61 GiB | **246 MiB**<br>6.70×±0.21 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.4 s | 6.78 s | **3.24 s**<br>3.82×±0.09 | 1.44 GiB | **210 MiB**<br>7.02×±0.27 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.46 s | 1.50 s | **0.78 s**<br>5.71×±0.16 | 604 MiB | **120 MiB**<br>5.05×±0.12 |

Across all whole-program and staged workloads, `jsgo` is **2.0×–6.6× faster** and uses **2.2×–8.7× less memory** than cold Java (and 1.4×–2.4× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#0a1a72c1-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
