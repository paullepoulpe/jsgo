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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on an 8-vCPU AMD EPYC-Genoa VPS (Linux, `--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `f0a308f8`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 414.7 s | — | **186.9 s**<br>2.22×±0.08 | 16.40 GiB | **5.24 GiB**<br>3.12×±0.27 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 91.5 s | 77.8 s | **30.8 s**<br>2.97×±0.05 | 3.75 GiB | **909 MiB**<br>4.22×±0.37 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 70.8 s | 59.2 s | **23.3 s**<br>3.04×±0.05 | 2.01 GiB | **715 MiB**<br>2.88×±0.21 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 44.7 s | 35.8 s | **13.4 s**<br>3.32×±0.07 | 1.53 GiB | **439 MiB**<br>3.57×±0.11 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 11.6 s | 6.17 s | **3.07 s**<br>3.78×±0.22 | 1.30 GiB | **280 MiB**<br>4.74×±0.39 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 22.6 s | 16.7 s | **5.87 s**<br>3.85×±0.17 | 1.87 GiB | **231 MiB**<br>8.05×±2.34 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 24.0 s | 15.2 s | **6.22 s**<br>3.86×±0.16 | 1.43 GiB | **228 MiB**<br>6.43×±0.17 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 7.27 s | 3.45 s | **1.28 s**<br>5.67×±0.20 | 615 MiB | **114 MiB**<br>5.38×±0.16 |

Across all whole-program and staged workloads, `jsgo` is **2.2×–6.1× faster** and uses **2.9×–8.4× less memory** than cold Java (and 1.7×–2.8× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#f0a308f8-on-jdk-27-epyc-vps) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
