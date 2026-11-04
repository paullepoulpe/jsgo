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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on an 8-vCPU AMD EPYC-Genoa VPS (Linux, `--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `a6c7ea5d`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 462.1 s | — | **190.9 s**<br>2.42×±0.02 | 16.36 GiB | **5.42 GiB**<br>3.02×±0.12 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 100.7 s | 87.5 s | **30.5 s**<br>3.30×±0.05 | 3.68 GiB | **939 MiB**<br>4.01×±0.19 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 80.1 s | 68.6 s | **22.9 s**<br>3.49×±0.06 | 2.03 GiB | **771 MiB**<br>2.69×±0.21 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 51.2 s | 40.1 s | **13.1 s**<br>3.91×±0.05 | 1.54 GiB | **463 MiB**<br>3.40×±0.20 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 13.0 s | 7.13 s | **2.98 s**<br>4.35×±0.26 | 1.34 GiB | **279 MiB**<br>4.92×±0.33 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 24.9 s | 18.5 s | **6.01 s**<br>4.14×±0.10 | 1.61 GiB | **225 MiB**<br>7.34×±0.35 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 26.1 s | 17.1 s | **5.90 s**<br>4.42×±0.05 | 1.44 GiB | **225 MiB**<br>6.54×±0.24 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 8.40 s | 4.05 s | **1.33 s**<br>6.32×±0.13 | 620 MiB | **115 MiB**<br>5.40×±0.09 |

Across all whole-program and staged workloads, `jsgo` is **2.4×–6.3× faster** and uses **2.7×–8.1× less memory** than cold Java (and 1.9×–3.1× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#a6c7ea5d-on-jdk-27-epyc-vps) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
