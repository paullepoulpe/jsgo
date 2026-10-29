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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `89d80973`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 205.7 s | — | **112.4 s**<br>1.83×±0.06 | 12.39 GiB | **6.18 GiB**<br>2.01×±0.06 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.4 s | 40.3 s | **21.6 s**<br>2.28×±0.01 | 4.62 GiB | **905 MiB**<br>5.22×±0.21 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.3 s | 26.4 s | **14.3 s**<br>2.46×±0.04 | 3.71 GiB | **731 MiB**<br>5.20×±0.23 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 23.2 s | 15.5 s | **8.32 s**<br>2.78×±0.04 | 3.42 GiB | **434 MiB**<br>8.06×±0.29 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.81 s | 2.77 s | **1.77 s**<br>3.84×±0.11 | 1.33 GiB | **266 MiB**<br>5.13×±0.71 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 12.9 s | 8.76 s | **3.89 s**<br>3.33×±0.03 | 1.95 GiB | **231 MiB**<br>8.27×±3.16 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.6 s | 6.76 s | **3.42 s**<br>3.67×±0.07 | 1.43 GiB | **219 MiB**<br>6.67×±0.21 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.51 s | 1.48 s | **0.84 s**<br>5.36×±0.20 | 606 MiB | **126 MiB**<br>4.82×±0.13 |

Across all whole-program and staged workloads, `jsgo` is **1.8×–6.2× faster** and uses **2.0×–8.5× less memory** than cold Java (and 1.3×–2.2× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#89d80973-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
