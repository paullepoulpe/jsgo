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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on an 8-vCPU AMD EPYC-Genoa VPS (Linux, `--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `65c3133d`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 471.1 s | — | **201.3 s**<br>2.34×±0.02 | 16.60 GiB | **5.54 GiB**<br>3.00×±0.15 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 100.5 s | 86.5 s | **32.6 s**<br>3.08×±0.03 | 3.56 GiB | **914 MiB**<br>3.98×±0.21 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 80.0 s | 66.8 s | **25.3 s**<br>3.17×±0.04 | 1.98 GiB | **728 MiB**<br>2.78×±0.16 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 51.6 s | 41.9 s | **14.7 s**<br>3.51×±0.08 | 1.55 GiB | **450 MiB**<br>3.52×±0.10 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 12.9 s | 7.12 s | **3.37 s**<br>3.83×±0.21 | 1.32 GiB | **270 MiB**<br>5.00×±0.49 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 24.7 s | 18.4 s | **6.05 s**<br>4.08×±0.15 | 1.86 GiB | **226 MiB**<br>8.17×±2.67 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 25.1 s | 17.2 s | **6.57 s**<br>3.82×±0.06 | 1.43 GiB | **234 MiB**<br>6.28×±0.34 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 8.04 s | 4.01 s | **1.45 s**<br>5.54×±0.29 | 621 MiB | **117 MiB**<br>5.31×±0.19 |

Across all whole-program and staged workloads, `jsgo` is **2.3×–6.1× faster** and uses **2.8×–8.2× less memory** than cold Java (and 1.9×–3.0× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#65c3133d-on-jdk-27-epyc-vps) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
