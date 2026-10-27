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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `f99b1259`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 204.6 s | — | **119.6 s**<br>1.71×±0.06 | 12.61 GiB | **6.37 GiB**<br>1.98×±0.07 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.4 s | 40.5 s | **23.4 s**<br>2.11×±0.03 | 4.54 GiB | **950 MiB**<br>4.89×±0.41 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 35.2 s | 26.3 s | **15.9 s**<br>2.22×±0.03 | 3.72 GiB | **747 MiB**<br>5.11×±0.22 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 23.4 s | 15.4 s | **9.37 s**<br>2.50×±0.06 | 3.45 GiB | **446 MiB**<br>7.94×±0.24 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 6.97 s | 2.75 s | **1.96 s**<br>3.56×±0.12 | 1.31 GiB | **281 MiB**<br>4.78×±0.35 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 13.2 s | 8.69 s | **4.37 s**<br>3.03×±0.08 | 1.60 GiB | **233 MiB**<br>7.05×±0.40 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 12.8 s | 6.74 s | **3.81 s**<br>3.35×±0.13 | 1.43 GiB | **220 MiB**<br>6.67×±0.16 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.58 s | 1.50 s | **0.90 s**<br>5.10×±0.18 | 602 MiB | **125 MiB**<br>4.80×±0.20 |

Across all whole-program and staged workloads, `jsgo` is **1.7×–5.5× faster** and uses **2.0×–8.4× less memory** than cold Java (and 1.2×–2.0× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#f99b1259-on-jdk-27-apple-m2) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
