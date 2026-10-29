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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `019139af`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 206.2 s | — | **104.7 s**<br>1.97×±0.01 | 12.37 GiB | **6.08 GiB**<br>2.04×±0.15 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.3 s | 40.1 s | **20.9 s**<br>2.36×±0.02 | 4.75 GiB | **894 MiB**<br>5.41×±0.93 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.2 s | 25.9 s | **14.0 s**<br>2.52×±0.02 | 3.71 GiB | **722 MiB**<br>5.27×±0.17 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 23.1 s | 15.5 s | **8.11 s**<br>2.85×±0.03 | 3.45 GiB | **426 MiB**<br>8.30×±0.25 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.85 s | 2.84 s | **1.71 s**<br>4.00×±0.10 | 1.36 GiB | **260 MiB**<br>5.37×±0.19 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 13.0 s | 8.69 s | **3.74 s**<br>3.47×±0.11 | 1.61 GiB | **229 MiB**<br>7.23×±0.29 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.6 s | 6.65 s | **3.32 s**<br>3.78×±0.10 | 1.43 GiB | **218 MiB**<br>6.73×±0.30 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.49 s | 1.46 s | **0.80 s**<br>5.58×±0.20 | 604 MiB | **123 MiB**<br>4.93×±0.23 |

Across all whole-program and staged workloads, `jsgo` is **2.0×–6.4× faster** and uses **2.0×–8.3× less memory** than cold Java (and 1.4×–2.3× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#019139af-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
