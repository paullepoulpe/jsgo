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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on an 8-vCPU AMD EPYC-Genoa VPS (Linux, `--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 27.0.1), at commit `1652e273`; ratios (`Java / jsgo`, against cold Java) report geometric means with 95% CIs (`±`). The warm-JVM column is Java's time in one JVM after warm-up runs (`CliBatch`), from the same measurement; `googleapis` has no warm run.

| Benchmark | Java Time (cold) | Java Time (warm JVM) | `jsgo` Time | Java RSS | `jsgo` RSS |
|---|---:|---:|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 418.8 s | — | **204.3 s**<br>2.05×±0.05 | 16.77 GiB | **5.59 GiB**<br>3.00×±0.24 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 93.2 s | 78.2 s | **32.2 s**<br>2.89×±0.13 | 3.71 GiB | **935 MiB**<br>4.06×±0.22 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 72.7 s | 59.9 s | **25.6 s**<br>2.84×±0.23 | 2.05 GiB | **709 MiB**<br>2.96×±0.18 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 45.5 s | 36.5 s | **14.4 s**<br>3.16×±0.09 | 1.53 GiB | **427 MiB**<br>3.68×±0.11 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 12.0 s | 6.30 s | **3.12 s**<br>3.85×±0.34 | 1.32 GiB | **272 MiB**<br>4.96×±0.34 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 22.7 s | 16.6 s | **6.07 s**<br>3.73×±0.10 | 1.62 GiB | **228 MiB**<br>7.28×±0.49 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 23.4 s | 15.2 s | **6.20 s**<br>3.77×±0.11 | 1.44 GiB | **221 MiB**<br>6.68×±0.11 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 7.52 s | 3.43 s | **1.33 s**<br>5.64×±0.35 | 615 MiB | **119 MiB**<br>5.18×±0.25 |

Across all whole-program and staged workloads, `jsgo` is **2.0×–6.1× faster** and uses **3.0×–8.5× less memory** than cold Java (and 1.7×–2.7× faster than a warm JVM). See [docs/benchmarks.md](docs/benchmarks.md#1652e273-on-jdk-27-epyc-vps) for all 13 workloads, warm JVM numbers, the per-stage `googleapis` breakdown, `t=1` results, the comparison with the previous measurement, and full `[lo, hi]` 95% confidence intervals.
