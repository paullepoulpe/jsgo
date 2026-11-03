# jsgo

A Go port of [Google Closure Compiler](https://github.com/google/closure-compiler), aiming to be a drop-in replacement with the same flags, byte-identical output, and identical diagnostics.

The port is a whole-file, 1:1 translation of the pinned Java release (**v20261028**; originally ported from v20260415 and upgraded through v20261015 and v20261021), built bottom-up and verified layer by layer against the Java compiler on large real-world corpora.

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

## Benchmarks (`jsgo` vs. Java Closure Compiler v20261028)

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on an 8-vCPU AMD EPYC-Genoa VPS (Linux, `--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `54998951`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 423.2 s | — | **213.0 s**<br>1.99×±0.08 | 16.99 GiB | **5.45 GiB**<br>3.12×±0.15 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 91.7 s | 78.5 s | **32.5 s**<br>2.83×±0.06 | 3.61 GiB | **911 MiB**<br>4.06×±0.23 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 71.9 s | 60.3 s | **25.2 s**<br>2.85×±0.06 | 2.01 GiB | **721 MiB**<br>2.85×±0.14 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 46.6 s | 36.8 s | **14.7 s**<br>3.16×±0.18 | 1.56 GiB | **428 MiB**<br>3.73×±0.19 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 12.2 s | 6.40 s | **3.10 s**<br>3.94×±0.24 | 1.30 GiB | **277 MiB**<br>4.81×±0.45 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 22.8 s | 16.6 s | **6.22 s**<br>3.67×±0.15 | 1.61 GiB | **222 MiB**<br>7.44×±0.62 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 23.5 s | 15.4 s | **6.45 s**<br>3.65×±0.14 | 1.44 GiB | **212 MiB**<br>6.96×±0.33 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 7.65 s | 3.48 s | **1.34 s**<br>5.72×±0.53 | 616 MiB | **117 MiB**<br>5.26×±0.19 |

Across all whole-program and staged workloads, `jsgo` is **2.0×–6.4× faster** and uses **2.9×–8.6× less memory** than cold Java (and 1.8×–2.7× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#54998951-on-jdk-27-epyc-vps) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
