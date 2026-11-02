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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on an 8-vCPU AMD EPYC-Genoa VPS (Linux, `--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `effd02e2`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 448.8 s | — | **216.3 s**<br>2.08×±0.03 | 16.45 GiB | **5.31 GiB**<br>3.10×±0.16 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 98.8 s | 84.9 s | **33.8 s**<br>2.92×±0.04 | 3.65 GiB | **906 MiB**<br>4.13×±0.14 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 78.6 s | 65.5 s | **26.4 s**<br>2.98×±0.08 | 2.35 GiB | **717 MiB**<br>3.26×±1.20 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 50.3 s | 40.4 s | **15.6 s**<br>3.23×±0.03 | 1.53 GiB | **447 MiB**<br>3.51×±0.02 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 13.1 s | 6.92 s | **3.17 s**<br>4.12×±0.14 | 1.33 GiB | **270 MiB**<br>5.05×±0.49 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 24.1 s | 18.4 s | **6.35 s**<br>3.81×±0.10 | 1.62 GiB | **233 MiB**<br>7.12×±0.54 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 26.1 s | 17.3 s | **6.64 s**<br>3.93×±0.11 | 1.43 GiB | **223 MiB**<br>6.57×±0.42 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 8.43 s | 3.75 s | **1.47 s**<br>5.72×±0.28 | 617 MiB | **115 MiB**<br>5.37×±0.26 |

Across all whole-program and staged workloads, `jsgo` is **2.1×–6.7× faster** and uses **3.1×–8.8× less memory** than cold Java (and 1.8×–2.9× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#effd02e2-on-jdk-27-epyc-vps) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
