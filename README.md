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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `50da0981`; ratios (`Java / jsgo`) report geometric means with 95% CIs (`±`).

| Benchmark | Java Time | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 207.5 s | **129.1 s**<br>1.61×±0.04 | 13.37 GiB | **5.05 GiB**<br>2.65×±0.12 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.3 s | **25.6 s**<br>1.93×±0.02 | 4.66 GiB | **880 MiB**<br>5.42×±0.27 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.2 s | **17.7 s**<br>1.99×±0.02 | 3.71 GiB | **690 MiB**<br>5.50×±0.16 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 23.0 s | **10.4 s**<br>2.21×±0.02 | 3.44 GiB | **423 MiB**<br>8.34×±0.14 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.87 s | **2.02 s**<br>3.40×±0.11 | 1.29 GiB | **272 MiB**<br>4.85×±0.36 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 12.9 s | **4.46 s**<br>2.90×±0.03 | 1.61 GiB | **234 MiB**<br>7.03×±0.43 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.5 s | **4.26 s**<br>2.92×±0.03 | 1.43 GiB | **205 MiB**<br>7.13×±0.28 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.55 s | **0.95 s**<br>4.79×±0.12 | 603 MiB | **122 MiB**<br>4.94×±0.26 |

Across all whole-program and staged workloads, `jsgo` is **1.6×–5.5× faster** and uses **2.7×–8.3× less memory** than cold Java (and 1.2×–2.0× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#50da0981-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
