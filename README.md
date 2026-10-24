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

Gated on **100% byte-identical output** (JS, source maps, diagnostics, and exit code). Measured across **$N = 5$ counterbalanced interleaved blocks** on Apple M2 (`--num_parallel_threads=4`, Go 1.27.3 + [PGO](cmd/jsgo/default.pgo), OpenJDK 26.0.2.1); ratios (`Java / jsgo`) report geometric means with 95% CIs (`±`).

| Benchmark | Java *(Cold / Warm)* | `jsgo` Time | `jsgo` RSS |
|---|---:|---:|---:|
| **`googleapis`**<br>*staged ADV+types · 35.6M nodes* | 207.9 s *(cold)*<br>14.1 GiB | **98.4 s**<br>2.11×±0.03 | **6.46 GiB**<br>2.18×±0.07 |
| **`bundle-mega`**<br>*ADV · 4.19M nodes* | 49.9 / 40.8 s<br>5.11 / 8.50 GiB | **20.6 s**<br>2.42×±0.03 / 1.98×±0.01 | **888 MiB**<br>5.9×±0.4 / 9.8× |
| **`typed-mega`**<br>*ADV+types · 3.55M nodes* | 34.9 / 26.7 s<br>3.71 / 4.53 GiB | **15.4 s**<br>2.27×±0.03 / 1.73×±0.03 | **862 MiB**<br>4.4×±0.1 / 5.4× |
| **`typed-full`**<br>*ADV+types · 1.90M nodes* | 24.5 / 15.7 s<br>3.44 / 3.54 GiB | **9.48 s**<br>2.58×±0.15 / 1.66×±0.12 | **540 MiB**<br>6.5×±0.4 / 6.7× |
| **`checks-only`**<br>*CHECKS+types · 1.29M nodes* | 7.01 / 2.86 s<br>1.37 / 3.42 GiB | **1.72 s**<br>4.07×±0.12 / 1.66×±0.07 | **336 MiB**<br>4.2×±0.3 / 10.4× |
| **`bundle`**<br>*SIMPLE · 1.11M nodes* | 12.8 / 8.68 s<br>2.13 / 3.74 GiB | **3.81 s**<br>3.35×±0.05 / 2.28×±0.02 | **293 MiB**<br>7.2×±2.6 / 13.1× |
| **`single-files`**<br>*1k CLI runs · 306k nodes* | 112.9 s *(batch)*<br>1.53 GiB | **98.5 s**<br>1.15×±0.01 *(vs. batch)* | **72.5 MiB**<br>21.5×±4.4 |
| **`downlevel-simple`**<br>*SIMPLE→ES5 · 232k nodes* | 4.67 / 1.50 s<br>0.61 / 1.45 GiB | **0.83 s**<br>5.64×±0.17 / 1.81×±0.08 | **132 MiB**<br>4.8×±0.1 / 11.3× |

Across all 13 benchmarks, `jsgo` is **2.1×–5.9× faster than cold Java**, **1.3×–2.9× faster than a JIT-warmed JVM**, and uses **2.2×–21.5× less memory**. See [docs/benchmarks.md](docs/benchmarks.md) for all 13 workloads, per-stage `googleapis` numbers, `--fast_mode` / `t=1` results, and full `[lo, hi]` 95% confidence intervals.
