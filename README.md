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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `8238e94e`; ratios (`Java / jsgo`) report geometric means with 95% CIs (`±`).

| Benchmark | Java Time | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 203.8 s | **121.8 s**<br>1.67×±0.02 | 12.70 GiB | **5.73 GiB**<br>2.22×±0.09 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.2 s | **23.6 s**<br>2.08×±0.01 | 4.74 GiB | **912 MiB**<br>5.32×±0.26 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.6 s | **16.1 s**<br>2.21×±0.03 | 3.70 GiB | **728 MiB**<br>5.20×±0.11 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 23.1 s | **9.52 s**<br>2.42×±0.02 | 3.44 GiB | **439 MiB**<br>8.03×±0.27 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.77 s | **1.94 s**<br>3.48×±0.06 | 1.35 GiB | **277 MiB**<br>4.98×±0.20 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 12.8 s | **4.32 s**<br>2.96×±0.05 | 1.61 GiB | **234 MiB**<br>7.05×±0.35 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.4 s | **3.86 s**<br>3.22×±0.06 | 1.44 GiB | **225 MiB**<br>6.53×±0.17 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.50 s | **0.91 s**<br>4.96×±0.20 | 607 MiB | **132 MiB**<br>4.61×±0.17 |

Across all whole-program and staged workloads, `jsgo` is **1.7×–5.5× faster** and uses **2.2×–8.0× less memory** than cold Java (and 1.2×–2.0× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#8238e94e-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
