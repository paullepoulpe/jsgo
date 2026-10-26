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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `fbdac472`; ratios (`Java / jsgo`) report geometric means with 95% CIs (`±`).

| Benchmark | Java Time | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 202.6 s | **126.4 s**<br>1.60×±0.02 | 13.61 GiB | **5.22 GiB**<br>2.61×±0.03 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.3 s | **25.6 s**<br>1.93×±0.02 | 4.73 GiB | **873 MiB**<br>5.55×±0.23 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.2 s | **17.7 s**<br>1.99×±0.02 | 3.71 GiB | **700 MiB**<br>5.42×±0.17 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 23.3 s | **10.4 s**<br>2.23×±0.05 | 3.44 GiB | **425 MiB**<br>8.28×±0.25 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.55 s | **2.02 s**<br>3.24×±0.09 | 1.34 GiB | **269 MiB**<br>5.09×±0.30 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 13.0 s | **4.48 s**<br>2.90×±0.04 | 1.63 GiB | **235 MiB**<br>7.13×±0.32 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.4 s | **4.27 s**<br>2.89×±0.06 | 1.43 GiB | **203 MiB**<br>7.24×±0.22 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.48 s | **0.95 s**<br>4.73×±0.16 | 606 MiB | **124 MiB**<br>4.89×±0.24 |

Across all whole-program and staged workloads, `jsgo` is **1.6×–5.4× faster** and uses **2.6×–8.3× less memory** than cold Java (and 1.1×–2.0× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#fbdac472-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
