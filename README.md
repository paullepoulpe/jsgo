# jsgo

A Go port of [Google Closure Compiler](https://github.com/google/closure-compiler), aiming to be a drop-in replacement with the same flags, byte-identical output, and identical diagnostics.

The port is a whole-file, 1:1 translation of the pinned Java release (**v20261028**; originally ported from v20260415 and upgraded through v20261015 and v20261021), built bottom-up and verified layer by layer against the Java compiler on large real-world corpora.

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

## Benchmarks (`jsgo` vs. Java Closure Compiler v20261028)

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on an 8-vCPU AMD EPYC-Genoa VPS (Linux, `--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `d6afd440`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 468.5 s | — | **197.8 s**<br>2.37×±0.01 | 16.17 GiB | **5.38 GiB**<br>3.00×±0.12 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 101.7 s | 88.6 s | **32.3 s**<br>3.15×±0.03 | 3.77 GiB | **915 MiB**<br>4.21×±0.15 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 79.8 s | 67.7 s | **24.9 s**<br>3.20×±0.07 | 2.05 GiB | **732 MiB**<br>2.86×±0.26 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 51.5 s | 41.7 s | **14.6 s**<br>3.52×±0.06 | 1.56 GiB | **438 MiB**<br>3.65×±0.19 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 12.8 s | 7.36 s | **3.08 s**<br>4.15×±0.14 | 1.34 GiB | **286 MiB**<br>4.78×±0.17 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 25.1 s | 18.4 s | **6.15 s**<br>4.08×±0.10 | 1.62 GiB | **226 MiB**<br>7.36×±0.40 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 26.2 s | 17.1 s | **6.65 s**<br>3.94×±0.12 | 1.44 GiB | **227 MiB**<br>6.48×±0.26 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 8.51 s | 4.01 s | **1.45 s**<br>5.88×±0.30 | 621 MiB | **117 MiB**<br>5.30×±0.23 |

Across all whole-program and staged workloads, `jsgo` is **2.4×–6.4× faster** and uses **2.9×–8.2× less memory** than cold Java (and 1.9×–3.0× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#d6afd440-on-jdk-27-epyc-vps) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
