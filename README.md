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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on an 8-vCPU AMD EPYC-Genoa VPS (Linux, `--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `645ecb68`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 438.4 s | — | **164.7 s**<br>2.66×±0.05 | 16.47 GiB | **5.54 GiB**<br>2.97×±0.08 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 97.0 s | 82.8 s | **29.4 s**<br>3.30×±0.03 | 3.61 GiB | **960 MiB**<br>3.85×±0.16 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 78.4 s | 64.6 s | **20.1 s**<br>3.90×±0.08 | 2.00 GiB | **781 MiB**<br>2.62×±0.24 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 50.7 s | 40.2 s | **11.7 s**<br>4.35×±0.08 | 1.55 GiB | **449 MiB**<br>3.53×±0.08 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 13.2 s | 7.36 s | **2.58 s**<br>5.14×±0.12 | 1.31 GiB | **289 MiB**<br>4.63×±0.19 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 24.8 s | 17.9 s | **5.28 s**<br>4.70×±0.07 | 1.62 GiB | **249 MiB**<br>6.69×±0.66 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 26.4 s | 17.6 s | **5.27 s**<br>5.02×±0.18 | 1.43 GiB | **239 MiB**<br>6.16×±0.42 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 8.46 s | 3.78 s | **1.26 s**<br>6.72×±0.33 | 625 MiB | **118 MiB**<br>5.29×±0.40 |

Across all whole-program and staged workloads, `jsgo` is **2.7×–6.7× faster** and uses **2.6×–7.8× less memory** than cold Java (and 2.0×–3.4× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#645ecb68-on-jdk-27-epyc-vps) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
