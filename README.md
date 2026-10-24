# jsgo

A Go port of [Google Closure Compiler](https://github.com/google/closure-compiler), aiming to be a drop-in replacement with the same flags, byte-identical output, and identical diagnostics.

The port is a whole-file, 1:1 translation of the pinned Java release (**v20261015**; originally ported from v20260415 and upgraded), built bottom-up and verified layer by layer against the Java compiler on large real-world corpora.

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

## Benchmarks (`jsgo` vs. Java Closure Compiler v20261015)

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Cold CLI runs measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 26.0.2.1); ratios (`Java / jsgo`) report geometric means with 95% CIs (`±`).

| Benchmark | Java | `jsgo` Time | `jsgo` RSS |
|---|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 207.9 s<br>14.1 GiB | **98.4 s**<br>2.11×±0.03 | **6.46 GiB**<br>2.18×±0.07 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.9 s<br>5.11 GiB | **20.6 s**<br>2.42×±0.03 | **888 MiB**<br>5.89×±0.38 |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 34.9 s<br>3.71 GiB | **15.4 s**<br>2.27×±0.03 | **862 MiB**<br>4.41×±0.09 |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 24.5 s<br>3.44 GiB | **9.48 s**<br>2.58×±0.15 | **540 MiB**<br>6.53×±0.39 |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 7.01 s<br>1.37 GiB | **1.72 s**<br>4.07×±0.12 | **336 MiB**<br>4.18×±0.29 |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 12.8 s<br>2.13 GiB | **3.81 s**<br>3.35×±0.05 | **293 MiB**<br>7.22×±2.64 |
| **`shaka-blockly`**<br>*ADV+types · 618k nodes* | 14.7 s<br>1.59 GiB | **4.03 s**<br>3.66×±0.08 | **267 MiB**<br>6.06×±0.81 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.67 s<br>629 MiB | **0.83 s**<br>5.64×±0.17 | **132 MiB**<br>4.78×±0.11 |

Across all whole-program and staged workloads, `jsgo` is **2.1×–5.9× faster** and uses **2.2×–8.7× less memory** than Java. See [docs/benchmarks.md](docs/benchmarks.md) for all 13 workloads, warm JVM numbers, per-stage `googleapis` breakdowns, `t=1` / `--fast_mode` results, and full `[lo, hi]` 95% confidence intervals.
