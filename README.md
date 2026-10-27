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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `80cb7956`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 206.6 s | — | **119.4 s**<br>1.73×±0.03 | 12.79 GiB | **5.52 GiB**<br>2.32×±0.13 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.0 s | 40.2 s | **23.6 s**<br>2.08×±0.02 | 4.66 GiB | **926 MiB**<br>5.14×±0.51 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.1 s | 26.5 s | **16.1 s**<br>2.18×±0.01 | 3.72 GiB | **741 MiB**<br>5.13×±0.18 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 22.9 s | 15.4 s | **9.45 s**<br>2.42×±0.03 | 3.42 GiB | **441 MiB**<br>7.95×±0.19 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.88 s | 2.83 s | **1.95 s**<br>3.53×±0.14 | 1.32 GiB | **270 MiB**<br>5.01×±0.38 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 12.9 s | 8.55 s | **4.32 s**<br>2.98×±0.05 | 1.60 GiB | **237 MiB**<br>6.90×±0.27 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.6 s | 6.75 s | **3.84 s**<br>3.27×±0.08 | 1.44 GiB | **220 MiB**<br>6.69×±0.26 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.51 s | 1.49 s | **0.90 s**<br>5.03×±0.14 | 603 MiB | **129 MiB**<br>4.69×±0.12 |

Across all whole-program and staged workloads, `jsgo` is **1.7×–5.5× faster** and uses **2.3×–8.2× less memory** than cold Java (and 1.2×–2.0× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#80cb7956-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
