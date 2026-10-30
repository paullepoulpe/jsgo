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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `d1c35b97`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 209.3 s | — | **102.0 s**<br>2.05×±0.02 | 12.27 GiB | **5.88 GiB**<br>2.09×±0.21 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.8 s | 40.1 s | **20.7 s**<br>2.41×±0.03 | 4.84 GiB | **911 MiB**<br>5.44×±0.38 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.4 s | 26.9 s | **13.7 s**<br>2.58×±0.03 | 3.71 GiB | **731 MiB**<br>5.20×±0.13 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 23.0 s | 15.3 s | **7.97 s**<br>2.89×±0.06 | 3.43 GiB | **425 MiB**<br>8.27×±0.18 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.77 s | 2.84 s | **1.71 s**<br>3.97×±0.16 | 1.28 GiB | **278 MiB**<br>4.69×±0.33 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 12.8 s | 8.76 s | **3.69 s**<br>3.48×±0.05 | 1.61 GiB | **236 MiB**<br>7.01×±0.43 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.5 s | 6.73 s | **3.29 s**<br>3.79×±0.08 | 1.43 GiB | **220 MiB**<br>6.66×±0.30 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.56 s | 1.49 s | **0.81 s**<br>5.66×±0.25 | 607 MiB | **122 MiB**<br>4.95×±0.15 |

Across all whole-program and staged workloads, `jsgo` is **2.1×–6.5× faster** and uses **2.1×–8.6× less memory** than cold Java (and 1.4×–2.4× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#d1c35b97-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
