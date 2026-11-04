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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on an 8-vCPU AMD EPYC-Genoa VPS (Linux, `--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `1a98e8a3`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 461.9 s | — | **189.7 s**<br>2.44×±0.08 | 16.07 GiB | **5.43 GiB**<br>2.96×±0.11 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 100.5 s | 87.3 s | **31.4 s**<br>3.20×±0.05 | 3.54 GiB | **951 MiB**<br>3.82×±0.20 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 80.0 s | 67.3 s | **22.6 s**<br>3.54×±0.09 | 2.00 GiB | **754 MiB**<br>2.72×±0.20 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 51.5 s | 40.8 s | **13.2 s**<br>3.91×±0.07 | 1.55 GiB | **459 MiB**<br>3.46×±0.24 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 13.2 s | 7.27 s | **2.96 s**<br>4.44×±0.13 | 1.37 GiB | **283 MiB**<br>4.98×±0.12 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 24.7 s | 18.5 s | **6.07 s**<br>4.08×±0.10 | 1.61 GiB | **242 MiB**<br>6.82×±0.34 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 26.2 s | 17.3 s | **5.87 s**<br>4.46×±0.19 | 1.44 GiB | **231 MiB**<br>6.40×±0.26 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 8.45 s | 3.93 s | **1.35 s**<br>6.26×±0.28 | 619 MiB | **119 MiB**<br>5.22×±0.11 |

Across all whole-program and staged workloads, `jsgo` is **2.4×–6.4× faster** and uses **2.7×–8.2× less memory** than cold Java (and 1.8×–3.1× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#1a98e8a3-on-jdk-27-epyc-vps) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
