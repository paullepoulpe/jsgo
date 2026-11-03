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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on an 8-vCPU AMD EPYC-Genoa VPS (Linux, `--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `a822a48e`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 421.8 s | — | **193.2 s**<br>2.18×±0.06 | 15.90 GiB | **5.39 GiB**<br>2.95×±0.12 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 92.5 s | 79.2 s | **31.3 s**<br>2.96×±0.08 | 3.56 GiB | **919 MiB**<br>3.96×±0.20 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 71.3 s | 62.1 s | **23.3 s**<br>3.05×±0.04 | 2.06 GiB | **733 MiB**<br>2.87×±0.16 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 46.6 s | 36.2 s | **13.9 s**<br>3.35×±0.09 | 1.55 GiB | **435 MiB**<br>3.66×±0.13 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 12.3 s | 6.53 s | **3.18 s**<br>3.86×±0.20 | 1.37 GiB | **277 MiB**<br>5.04×±0.74 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 22.8 s | 18.4 s | **6.03 s**<br>3.77×±0.10 | 1.62 GiB | **225 MiB**<br>7.36×±0.64 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 24.6 s | 16.1 s | **6.28 s**<br>3.92×±0.36 | 1.46 GiB | **225 MiB**<br>6.65×±0.45 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 7.53 s | 3.46 s | **1.29 s**<br>5.82×±0.35 | 620 MiB | **119 MiB**<br>5.21×±0.16 |

Across all whole-program and staged workloads, `jsgo` is **2.2×–6.1× faster** and uses **2.9×–8.6× less memory** than cold Java (and 1.9×–3.0× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#a822a48e-on-jdk-27-epyc-vps) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
