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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `a8c90bf9`; ratios (`Java / jsgo`) report geometric means with 95% CIs (`±`).

| Benchmark | Java Time | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 204.3 s | **123.5 s**<br>1.65×±0.02 | 13.31 GiB | **5.04 GiB**<br>2.64×±0.15 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.4 s | **25.1 s**<br>1.97×±0.01 | 4.41 GiB | **842 MiB**<br>5.35×±0.49 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.5 s | **17.5 s**<br>2.03×±0.03 | 3.70 GiB | **738 MiB**<br>5.14×±0.27 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 23.7 s | **10.45 s**<br>2.27×±0.08 | 3.45 GiB | **421 MiB**<br>8.39×±0.32 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 7.23 s | **2.04 s**<br>3.54×±0.11 | 1.30 GiB | **271 MiB**<br>4.93×±0.44 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 13.3 s | **4.49 s**<br>2.95×±0.06 | 1.60 GiB | **222 MiB**<br>7.38×±0.62 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.9 s | **4.27 s**<br>3.02×±0.05 | 1.45 GiB | **204 MiB**<br>7.26×±0.05 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.74 s | **0.94 s**<br>5.05×±0.27 | 604 MiB | **128 MiB**<br>4.73×±0.07 |

Across all whole-program and staged workloads, `jsgo` is **1.6×–5.5× faster** and uses **2.6×–8.7× less memory** than cold Java (and 1.2×–2.0× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#a8c90bf9-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
