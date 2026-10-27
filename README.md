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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `9b7b72ed`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 209.6 s | — | **120.9 s**<br>1.73×±0.05 | 12.10 GiB | **6.15 GiB**<br>1.97×±0.22 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.4 s | 40.4 s | **23.7 s**<br>2.09×±0.01 | 4.61 GiB | **931 MiB**<br>5.07×±0.47 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.2 s | 26.3 s | **16.2 s**<br>2.17×±0.04 | 3.71 GiB | **729 MiB**<br>5.22×±0.10 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 23.0 s | 15.8 s | **9.51 s**<br>2.42×±0.04 | 3.44 GiB | **444 MiB**<br>7.93×±0.17 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.83 s | 2.87 s | **1.96 s**<br>3.49×±0.21 | 1.35 GiB | **269 MiB**<br>5.14×±0.24 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 13.1 s | 8.77 s | **4.34 s**<br>3.01×±0.03 | 1.64 GiB | **229 MiB**<br>7.32×±0.28 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.4 s | 6.82 s | **3.87 s**<br>3.21×±0.03 | 1.43 GiB | **219 MiB**<br>6.66×±0.03 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.54 s | 1.52 s | **0.91 s**<br>4.97×±0.25 | 600 MiB | **129 MiB**<br>4.65×±0.12 |

Across all whole-program and staged workloads, `jsgo` is **1.7×–5.6× faster** and uses **2.0×–8.8× less memory** than cold Java (and 1.2×–2.0× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#9b7b72ed-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
