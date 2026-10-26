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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `64c098f3`; ratios (`Java / jsgo`) report geometric means with 95% CIs (`±`).

| Benchmark | Java Time | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 205.0 s | **123.7 s**<br>1.66×±0.02 | 12.96 GiB | **5.31 GiB**<br>2.44×±0.09 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.2 s | **25.0 s**<br>1.97×±0.02 | 4.86 GiB | **885 MiB**<br>5.63×±0.17 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.7 s | **17.3 s**<br>2.07×±0.01 | 3.72 GiB | **737 MiB**<br>5.17×±0.16 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 23.5 s | **10.2 s**<br>2.30×±0.06 | 3.44 GiB | **431 MiB**<br>8.18×±0.14 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.86 s | **1.97 s**<br>3.48×±0.07 | 1.33 GiB | **261 MiB**<br>5.21×±0.22 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 13.3 s | **4.45 s**<br>2.98×±0.09 | 1.62 GiB | **236 MiB**<br>7.03×±0.34 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.6 s | **4.19 s**<br>3.00×±0.09 | 1.45 GiB | **210 MiB**<br>7.06×±0.40 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.43 s | **0.95 s**<br>4.68×±0.14 | 605 MiB | **126 MiB**<br>4.80×±0.12 |

Across all whole-program and staged workloads, `jsgo` is **1.7×–5.5× faster** and uses **2.4×–8.5× less memory** than cold Java (and 1.2×–2.0× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#64c098f3-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
