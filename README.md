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

## Benchmarks (`jsgo` vs. Java Closure Compiler v20261021)

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on an 8-vCPU AMD EPYC-Genoa VPS (Linux, `--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `a167d9c7`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 431.8 s | — | **205.3 s**<br>2.10×±0.02 | 16.73 GiB | **5.50 GiB**<br>3.05×±0.38 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 93.7 s | 81.0 s | **32.7 s**<br>2.87×±0.08 | 3.63 GiB | **938 MiB**<br>3.96×±0.12 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 73.4 s | 61.8 s | **25.3 s**<br>2.90×±0.05 | 2.07 GiB | **705 MiB**<br>3.00×±0.06 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 47.9 s | 38.1 s | **14.9 s**<br>3.21×±0.09 | 1.53 GiB | **439 MiB**<br>3.56×±0.20 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 11.8 s | 6.53 s | **3.12 s**<br>3.76×±0.30 | 1.30 GiB | **273 MiB**<br>4.87×±0.34 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 23.3 s | 17.1 s | **6.31 s**<br>3.69×±0.15 | 1.61 GiB | **226 MiB**<br>7.31×±0.30 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 23.9 s | 16.1 s | **6.44 s**<br>3.71×±0.11 | 1.43 GiB | **219 MiB**<br>6.71×±0.23 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 7.50 s | 3.61 s | **1.35 s**<br>5.57×±0.22 | 609 MiB | **115 MiB**<br>5.30×±0.19 |

Across all whole-program and staged workloads, `jsgo` is **2.1×–6.9× faster** and uses **3.0×–8.4× less memory** than cold Java (and 1.9×–2.7× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#a167d9c7-on-jdk-27-epyc-vps) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
