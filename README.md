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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on an 8-vCPU AMD EPYC-Genoa VPS (Linux, `--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `cd72a29e`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 370.1 s | — | **142.8 s**<br>2.59×±0.01 | 16.57 GiB | **5.52 GiB**<br>3.00×±0.17 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 81.4 s | 70.2 s | **25.4 s**<br>3.20×±0.04 | 3.71 GiB | **970 MiB**<br>3.91×±0.26 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 63.4 s | 51.7 s | **16.1 s**<br>3.94×±0.24 | 1.99 GiB | **733 MiB**<br>2.78×±0.18 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 41.2 s | 30.8 s | **9.20 s**<br>4.47×±0.30 | 1.56 GiB | **464 MiB**<br>3.45×±0.22 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 10.5 s | 5.17 s | **2.03 s**<br>5.18×±0.27 | 1.31 GiB | **291 MiB**<br>4.58×±0.42 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 19.4 s | 13.8 s | **4.39 s**<br>4.41×±0.12 | 1.61 GiB | **234 MiB**<br>7.09×±0.43 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 20.5 s | 12.1 s | **3.85 s**<br>5.33×±0.04 | 1.43 GiB | **235 MiB**<br>6.24×±0.25 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 6.68 s | 2.41 s | **1.03 s**<br>6.51×±0.33 | 617 MiB | **120 MiB**<br>5.14×±0.19 |

Across all whole-program and staged workloads, `jsgo` is **2.6×–7.1× faster** and uses **2.8×–7.6× less memory** than cold Java (and 1.4×–3.4× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#cd72a29e-on-jdk-27-epyc-vps) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
