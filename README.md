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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `29f91ce6`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 205.3 s | — | **119.6 s**<br>1.72×±0.03 | 12.71 GiB | **6.09 GiB**<br>2.09×±0.17 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.3 s | 40.18 s | **23.6 s**<br>2.09×±0.02 | 4.74 GiB | **927 MiB**<br>5.21×±0.80 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.2 s | 26.45 s | **16.1 s**<br>2.19×±0.02 | 3.71 GiB | **755 MiB**<br>5.04×±0.09 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 23.1 s | 15.57 s | **9.45 s**<br>2.44×±0.05 | 3.42 GiB | **440 MiB**<br>7.96×±0.32 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.74 s | 2.77 s | **1.95 s**<br>3.46×±0.13 | 1.32 GiB | **262 MiB**<br>5.15×±0.38 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 12.9 s | 8.57 s | **4.32 s**<br>2.99×±0.08 | 1.61 GiB | **237 MiB**<br>6.97×±0.44 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.5 s | 6.77 s | **3.85 s**<br>3.24×±0.09 | 1.43 GiB | **213 MiB**<br>6.84×±0.24 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.44 s | 1.48 s | **0.90 s**<br>4.95×±0.16 | 603 MiB | **128 MiB**<br>4.70×±0.17 |

Across all whole-program and staged workloads, `jsgo` is **1.7×–5.5× faster** and uses **2.1×–8.4× less memory** than cold Java (and 1.2×–2.0× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#29f91ce6-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
