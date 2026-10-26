# Benchmarks: Go port vs Java Closure Compiler

Per-milestone performance checks (PLAN.md §6). Every comparison runs the
same layer of both implementations on the same inputs. Both outputs are
first verified identical by the oracles.

The end-to-end benchmark (whole CLI invocations, overall and per phase,
behind a byte-for-byte correctness gate) is `tools/bench/bench.py`; see
[tools/bench/README.md](../tools/bench/README.md).

## End-to-end, v20261015

The M9 acceptance run of PLAN.md §6: jsgo against the pinned Java
compiler, v20261015, on the six tasks of `tools/bench/tasks/`. The full
generated report, with every task's per-pass table, is
[benchmarks-e2e-v20261015.md](benchmarks-e2e-v20261015.md).

**Setup:**

- **Machine:** 4-vCPU VM (Intel Xeon @ 2.10 GHz, 15.7 GiB, Ubuntu 24.04,
  kernel 6.18), with nothing else running. The load averages in the
  report are the benchmark's own: a cold JVM keeps about three cores busy.
- **Java:** OpenJDK 21.0.12, default settings; the Maven Central jar of
  v20261015 (`tools/oracle/pin.env`).
- **jsgo:** 076ba05 (the work branch at 07b0f28 plus a regenerated
  `cmd/jsgo/default.pgo`), Go 1.24.7, built with PGO, default GC.
- **Method:** `bench.py run` with its defaults (2026-10-18, 01:43–05:58
  UTC):
  - 5 cold runs per compiler, Java and jsgo interleaved;
  - Java warm: 2 warm-up iterations (1 for the bundle and the single
    files), then 5 measured, in one JVM;
  - 3 `--tracer_mode=TIMING_ONLY` runs per compiler (1 for the single
    files, warm Java and jsgo only);
  - everything under `nice -n 10`.

**Correctness gate:**

- All five whole-program tasks **MATCH**: stdout, stderr, exit status,
  output and source map are byte-identical to Java's.
- Single files: 989 of 1,000 MATCH. The other 11 are inputs on which Java
  throws an internal exception. jsgo prints the same exception, message,
  exit status and output, but not the JVM's `\tat …` frames.
- Java's cold runs are byte-identical to each other, and the warm
  driver's last iteration reproduces cold run 1.

**Overall.** Wall time, median of the runs (min–max). Speed-ups are Java
÷ jsgo, so above 1× means jsgo is faster. CPU is user+system; for Java warm
it is the JVM's process CPU during an iteration, JIT and GC threads
included. Peak RSS for Java warm is the JVM's over its whole life.

| Task | Java cold | Java warm | jsgo | jsgo vs cold | jsgo vs warm | CPU s: cold / warm / jsgo | Peak RSS MB: cold / warm / jsgo |
|---|---|---|---|---|---|---|---|
| Typed app (ADVANCED, checkTypes, source map) | 11.6 s (11.3–12.1) | 4.90 s (4.47–5.14) | 7.27 s (7.04–7.40) | **1.60×** | 0.67× | 37.2 / 8.46 / 9.65 | 649 / 1133 / 252 |
| Bundle, 12.6 MB (SIMPLE) | 26.6 s (26.4–27.4) | 18.9 s (18.3–19.8) | 28.3 s (27.8–29.0) | **0.94×** | 0.67× | 63.3 / 22.3 / 37.7 | 1720 / 1741 / 552 |
| Down-levelling, WHITESPACE_ONLY | 3.91 s (3.78–4.18) | 1.26 s (1.04–1.48) | 1.72 s (1.68–1.83) | **2.28×** | 0.74× | 12.5 / 3.11 / 2.34 | 386 / 765 / 163 |
| Down-levelling, SIMPLE | 8.35 s (7.65–8.66) | 3.54 s (3.10–3.95) | 5.35 s (5.31–6.25) | **1.56×** | 0.66× | 27.0 / 5.65 / 7.18 | 447 / 1082 / 206 |
| Checks only (Closure Library) | 13.6 s (13.0–13.9) | 6.77 s (6.72–6.87) | 9.92 s (9.76–10.2) | **1.37×** | 0.68× | 41.4 / 9.93 / 13.4 | 1299 / 1545 / 690 |
| 1,000 single files, per run | 1805 s (1799–1816) | 278 s (275–278) | 426 s (421–436) | **4.24×** | 0.65× | 5255 / 283 / 587 | 373 / 922 / 83 |

Notes on the table:

- **Single files.** The speed-up against cold Java is over the 989
  matching compiles: 1.96 s per compile for Java, 0.48 s for jsgo. The
  ratio against warm Java is over all 1,000. The RSS is the largest single
  process, except for Java warm, which runs all compiles in one JVM.
- **Memory.** jsgo's peak RSS is 1.9–3.1× below cold Java's (4.5× per
  single-file process) and 2.2–5.3× below the warm JVM's.
- **CPU.** jsgo uses 1.3–1.4 CPU-seconds per second of wall time; the
  part above 1 is mostly the GC workers. Cold Java uses 2.4–3.2 (JIT and
  GC threads). On whole-program
  tasks jsgo needs 1.7–5.3× less CPU than cold Java, and 0.75–1.7× as
  much as warm Java.
- **Bundle.** It is the only task where jsgo is slower than cold Java. A
  27 s compile is long enough for the JIT to warm up during the run, so
  cold Java's optimization loop takes only 1.3× its warm time (14.5 s
  against 11.2 s). jsgo's 1.7× gap in that loop then shows in full.

**PGO.** Measured on a jsgo built with `go build -pgo=off`, interleaved
with the PGO build (N=5 each, in the control session below):

| Task | PGO (default.pgo) | `-pgo=off` | Difference |
|---|---|---|---|
| typed-app | 7.48 s (7.04–7.59), CPU 10.07 s | 7.30 s (7.26–7.52), CPU 9.95 s | none measurable (+2.5% wall, inside the spread) |
| downlevel-simple | 5.83 s (5.60–6.03), CPU 7.82 s | 5.98 s (5.83–6.15), CPU 8.06 s | −2.5% wall, −3.0% CPU |
| checks-only | 11.22 s (11.02–11.48), CPU 15.13 s | 11.49 s (10.92–11.79), CPU 15.64 s | −2.3% wall, −3.3% CPU |

The profile buys 0–3%. On the round-3 code it saved 2–7% CPU.

### Per phase

The tables below are milliseconds, medians of the trace runs.
`phases.tsv` rolls passes up into phases. *startup* is process start to
the tracker's start. *Code printing* is the tracker's end to process exit:
toSource, the source map, writing the output, and exit. Ratios are Java ÷
jsgo.

Typed app:

| Phase | Java cold | Java warm | jsgo | vs cold | vs warm |
|---|---|---|---|---|---|
| startup | 494 | 33 | 40 | 12.36× | 0.82× |
| parse | 1620 | 544 | 743 | 2.18× | 0.73× |
| checks | 1537 | 626 | 891 | 1.73× | 0.70× |
| type checking | 3171 | 1102 | 1586 | 2.00× | 0.69× |
| transpilation | 271 | 121 | 149 | 1.82× | 0.81× |
| optimization loop | 3647 | 2019 | 3516 | 1.04× | 0.57× |
| renaming | 125 | 71 | 121 | 1.03× | 0.59× |
| other passes | 70 | 28 | 30 | 2.33× | 0.93× |
| untraced (inside tracker) | 183 | 61 | 96 | 1.91× | 0.64× |
| code printing and source maps | 313 | 69 | 127 | 2.47× | 0.54× |
| **total** | **11432** | **4674** | **7299** | **1.57×** | **0.64×** |

Bundle:

| Phase | Java cold | Java warm | jsgo | vs cold | vs warm |
|---|---|---|---|---|---|
| startup | 444 | 16 | 30 | 14.97× | 0.54× |
| parse | 2547 | 864 | 1633 | 1.56× | 0.53× |
| checks | 4163 | 2665 | 3323 | 1.25× | 0.80× |
| type checking | 260 | 133 | 150 | 1.73× | 0.89× |
| transpilation | 742 | 522 | 495 | 1.50× | 1.05× |
| optimization loop | 14487 | 11215 | 19354 | 0.75× | 0.58× |
| renaming | 2401 | 1934 | 2828 | 0.85× | 0.68× |
| other passes | 273 | 162 | 207 | 1.32× | 0.78× |
| untraced (inside tracker) | 368 | 260 | 327 | 1.13× | 0.80× |
| code printing and source maps | 675 | 191 | 434 | 1.56× | 0.44× |
| **total** | **26360** | **17962** | **28781** | **0.92×** | **0.62×** |

Checks only:

| Phase | Java cold | Java warm | jsgo | vs cold | vs warm |
|---|---|---|---|---|---|
| startup | 514 | 19 | 38 | 13.60× | 0.50× |
| parse | 3237 | 888 | 1965 | 1.65× | 0.45× |
| checks | 4381 | 2613 | 3558 | 1.23× | 0.73× |
| type checking | 4803 | 2575 | 4351 | 1.10× | 0.59× |
| transpilation | 402 | 311 | 347 | 1.16× | 0.90× |
| untraced (inside tracker) | 159 | 54 | 78 | 2.04× | 0.69× |
| code printing and source maps | 108 | 5 | 53 | 2.03× | 0.09× |
| **total** | **13604** | **6465** | **10390** | **1.31×** | **0.62×** |

Single files, summed over the 989 compiles that print a tracer report.
Inside one JVM, startup and printing cannot be told apart, so warm Java's
time outside the trackers is all on the printing row:

| Phase | Java warm | jsgo |
|---|---|---|
| startup | 0 | 30616 |
| parse | 60049 | 153524 |
| checks | 85253 | 88985 |
| type checking | 18433 | 27846 |
| transpilation | 14482 | 12362 |
| optimization loop | 60012 | 85650 |
| renaming | 7230 | 12214 |
| other passes | 7320 | 7630 |
| untraced (inside tracker) | 3989 | 11069 |
| code printing and source maps | 16208 | 6398 |
| **total** | **272976** | **436293** |

**Where jsgo still loses to warm Java.** These are per pass, from the full
report:

- **The optimization loop has the largest absolute gap**, at 0.57–0.61× of
  warm Java's speed:
  - typed-app: 3.5 s against 2.0 s;
  - bundle: 19.4 s against 11.2 s, three quarters of the bundle's 10.8 s
    gap;
  - downlevel-simple: 2.6 s against 1.6 s.

  The worst passes, as jsgo's time over warm Java's:
  - `flowSensitiveInlineVariables`: 2.4–2.8× (bundle 3.55 s against
    1.28 s);
  - `inlineVariables`: about 2× (bundle 2.49 s against 1.25 s);
  - `removeUnusedCode`: 1.5–1.9×;
  - `collapseObjectLiterals` and `coalesceVariableNames`: about 2×.
- **Parse is 0.45–0.82× of warm Java's speed.**
  - checks-only: 1.97 s against 0.89 s.
  - bundle: 1.63 s against 0.86 s.
  - Single files: 155 ms per compile against 61 ms, mostly the default
    externs that every compile parses. This is the largest item of that
    task.
- **Type checking:**
  - `inferTypes` takes 1.4–1.8× warm Java's time (checks-only: 3.46 s
    against 1.95 s);
  - `checkAccessControls` about 2.1× (0.73 s against 0.34 s).
- **Renaming on the bundle:** `invertContextualRenaming` takes 2.14 s
  against 1.46 s.
- **Code printing and exit** run at 0.41–0.60× of warm Java's speed
  (checks-only: 53 ms against 5 ms), but take under 0.45 s everywhere.

jsgo is at or ahead of warm Java on:
- process startup: 30–46 ms, against 0.44–0.52 s for a cold JVM;
- the bundle's transpilation;
- several gathering passes (`gatherGettersAndSetters`,
  `gatherModuleMetadata`).

### Against the v20260415 round-3 numbers

These numbers differ a lot from the ones recorded after performance round
3 on v20260415:
- Java cold has about halved: typed-app went from 22.0 s to 11.6 s.
- Java warm has also dropped.
- jsgo is 18–27% faster.

Those earlier figures were measured on the same VM type, but next to
porting agents: load 5–15 for the Java baseline and 4–6 for jsgo. To
split the change, a same-machine control ran right after the main run
(06:00–06:20 UTC). In one session, interleaved, it timed:
- the v20260415 jar against the v20261015 jar, cold (N=3) and warm (2
  warm-up + 3 iterations);
- jsgo built from 2fb9080 (the v20260415 port with the round-3 code and
  its own `default.pgo`) against the current jsgo (N=5).

The old builds' outputs are not gated: they compile for v20260415.

| Wall time, median | typed-app | downlevel-simple | checks-only |
|---|---|---|---|
| Java cold, round 3 (v20260415, loaded) | 22.0 s | 14.6 s | 25.6 s |
| Java cold, v20260415 jar, quiet | 11.10 s | 8.53 s | 14.71 s |
| Java cold, v20261015 jar, quiet | 11.62 s (+4.7%) | 8.82 s (+3.4%) | 15.03 s (+2.2%) |
| Java warm, round 3 | 8.21 s | 3.99 s | 8.09 s |
| Java warm, v20260415 jar, quiet | 4.52 s | 3.87 s | 7.12 s |
| Java warm, v20261015 jar, quiet | 4.87 s (+7.7%) | 4.06 s (+4.9%) | 7.42 s (+4.2%) |
| jsgo, round 3 (loaded) | 8.84 s | 7.28 s | 12.87 s |
| jsgo 2fb9080 (v20260415 port), quiet | 7.10 s | 5.84 s | 11.14 s |
| jsgo 076ba05 (v20261015 port), quiet | 7.48 s (+5.4%) | 5.83 s (−0.2%) | 11.22 s (+0.7%) |

- **The machine accounts for nearly all of the change.** On the quiet
  machine, the same v20260415 jar runs 42–50% faster cold, and the
  round-3 jsgo runs 13–20% faster.
  - A cold JVM keeps about three cores busy with JIT and GC threads, so it
    lost the most when other jobs held those cores.
  - jsgo uses about 1.35 cores.
  - The round-3 warm baselines were themselves noisy: typed-app's warm
    iterations ranged from 6.7 to 14.8 s.
- **Upstream Java v20261015 is 2–8% slower** than v20260415 on these
  tasks.
- **The port's v20261015 code is 0–5% slower** than the round-3 code. The
  only clear case is typed-app's +5%, in the optimization loop, mostly
  `inlineVariables`: jsgo 544 → 659 ms. Java warm's `inlineVariables`
  slowed too (281 → 354 ms), so the port follows upstream there. Other
  per-phase differences are within the noise of three trace runs.
- **The pass names changed in two places:**
  - `exploitAssign` became `postNormalizePeephole`;
  - `es6RewriteRestAndSpread` split into `es6RewriteRestParameters` and
    `es6RewriteSpreadExpressions`.

  `phases.tsv` now maps all three.
- **The ratios against warm Java did not move.** In the same session,
  round-3 jsgo against the v20260415 warm JVM is 0.64× (typed-app),
  0.66× (downlevel-simple) and 0.64× (checks-only). v20261015 jsgo against
  the v20261015 warm JVM is 0.65×, 0.70× and 0.66×. Round 3's 0.93×
  against warm Java on typed-app came from comparing against a loaded
  JVM.
- **Absolute times drift on this VM from one session to the next.** The
  control measured the same jsgo and jar 0–13% slower than the main run
  had, four hours earlier: checks-only jsgo 11.2 s against 9.9 s, Java cold
  15.0 s against 13.6 s. Interleaved comparisons within a session are
  good to a few percent; absolute numbers across sessions are not.

To reproduce:

```sh
tools/bench/fetch-inputs.sh
python3 -I tools/bench/pgo.py      # cmd/jsgo/default.pgo
python3 -I tools/bench/bench.py run --out /tmp/bench --report /tmp/bench/report.md
python3 -I tools/bench/bench.py report --out /tmp/bench --md report.md --final
```

### Allocation, compact layout, and `--num_parallel_threads` pass

After the v20261015 baseline above, slab-batched `jutil` map entries, a 64-var
indexed bitmask lattice in `MustBeReachingVariableDef`, compact `rhino.Node`
(96 B → 80 B) and `jscomp.Var` (80 B → 64 B), scanner `SourcePosition` reuse,
rename-conflict deduplication, and opt-in `--num_parallel_threads` prebuild
were integrated (see [performance-roadmap.md](performance-roadmap.md)).

Interleaved comparison (`N=5` per Go build, `N=3` cold Java, `N=5` warm Java,
all byte-identical; macOS arm64, Go 1.27.1 with `default.pgo`):

| Task | Baseline (`c9872522`) | New (`t=1` / `t=4`) | Δ Wall (`t=1` / `t=4`) | vs Cold / Warm Java Time (`t=1` / `t=4`) | Peak RSS: Go (`base` → `new`) | Peak RSS: Cold / Warm Java | Java / Go RSS (Cold / Warm) |
|---|---:|---:|---:|---:|---:|---:|---:|
| `typed-app` | 2.86 s (2.84–3.48) | **2.74 s** / **2.64 s** | **−4.0%** / **−7.7%** | 2.20× / 0.76× (0.79× at `t=4`) | 234 → **242 MiB** | 758 / 1,871 MiB | **3.13× / 7.72×** |
| `bundle` (12.6 MB) | 11.82 s (11.80–11.98) | **10.31 s** / 10.31 s | **−12.8%** | **1.22×** / **0.85×** | 543 → **490 MiB** (−9.7%) | 2,151 / 2,343 MiB | **4.39× / 4.78×** |
| `downlevel-whitespace` | 0.86 s (0.85–0.87) | **0.84 s** / **0.70 s** | **−2.5%** / **−18.2%** | 2.92× / 0.74× (0.88× at `t=4`) | 168 → **145 MiB** (−13.7%) | 553 / 1,437 MiB | **3.82× / 9.92×** |
| `downlevel-simple` | 2.39 s (2.38–2.41) | **2.25 s** / **2.10 s** | **−5.9%** / **−12.2%** | 2.00× / 0.72× (0.77× at `t=4`) | 203 → **190 MiB** (−6.6%) | 613 / 1,489 MiB | **3.23× / 7.84×** |
| `checks-only` | 4.40 s (4.39–4.42) | **4.29 s** / **3.82 s** | **−2.4%** / **−13.1%** | 1.57× / 0.76× (0.85× at `t=4`) | 693 → **618 MiB** (−10.7%) | 1,591 / 3,562 MiB | **2.57× / 5.76×** |

On a 2-vCPU Linux VM (AMD EPYC 9V74, JDK 21.0.12.1, 12 warm-up + 10 measured
iterations), the same changes reduce Go elapsed time by 8.8–19.8% across the
five tasks (`bundle` 23.10 s → 18.53 s, −19.8%; RSS 675 → 504 MiB, −25.3%),
with warm Java / Go peak RSS ratios of **4.74×** (`typed-app`), **3.33×**
(`bundle`), **5.72×** (`downlevel-whitespace`), **7.13×** (`downlevel-simple`),
and **2.19×** (`checks-only`).

### Rounds 2–2S, `googleapis-bench` (35.1M AST Nodes, 7,460 Chunks), and Large / Mega Real-World `ADVANCED` Benchmarks

Following substrate, scanner, AST-slab, flat-map, CFG/dataflow-slab, type-system,
pass-level, inter-pass arena, peephole-dispatch, source-map, `Scope`
open-addressed hash table + 64-bit Bloom filter, `ParseArena` pool draining,
`GlobalNamespace` lazy scope resolution + O(1) tombstones, `RemoveUnusedCode`
slab retention, `ReferenceCollector` two-tier (`global` / `local`) chunk
retention, `AnalyzePrototypeProperties` duplicate-edge elimination, `InferConsts`
per-function local state flushing, `RenameVars` chunked reference slabs,
`OptimizeCalls` direct AST walk, parallel function-local middle-end passes
(`DeadAssignmentsElimination`, `CoalesceVariableNames`,
`FlowSensitiveInlineVariables`, `InlineObjectLiterals`), single-cache-line 64 B
`rhino.Node` (`5` pointers, inline leaf `payload uint64`, 64-bit `propFlags`),
pointer-free 64-shard `RhinoStringPool` with 32-bit `rhinoStringID` handles,
embedded `ArrowType` + 128 B `UnionType` (`[3]JSType` inline), parallel
`TypedAst` / `PrebuildAst` / Stage 3 `CodePrinter`, `≤ 32,760 B` small-span
arenas, and in-place AST node reuse in `CollapseProperties` &
`DevirtualizeMethods` (see [performance-roadmap.md](performance-roadmap.md)),
`jsgo` beats both cold Java and a 5-warmup JIT-compiled JVM across the standard
suite, 6 large/mega real-world `ADVANCED` compilation workloads up to 44.6 MB /
4,385 files, and the 35.1M-AST-node / 7,460-chunk `googleapis-bench` staged
compilation benchmark (macOS arm64, Go 1.27.1 with `default.pgo`, OpenJDK
21.0.10, all 100% byte-identical):

#### Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, ~35.1M AST Nodes, 7,460 Chunks, `--threads 4`)

| Stage | `jsgo` (`compat`) | `jsgo` (`--fast_mode`) | Java (`-Xmx20g`) | Speedup vs. Java (`compat` / `fast`) | `jsgo` Peak RSS | Java Peak RSS | Memory Saved |
|---|---:|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **49.36 s** | **48.71 s** | 150.70 s | **3.05×** / **3.09×** | **6,083.8 MiB** *(5,821.2 best)* | 10,515.8 MiB | **−4,432 MiB** |
| **`OPTIMIZATIONS`** | **39.53 s** | **39.30 s** | 102.77 s | **2.60×** / **2.62×** | **5,525.4 MiB** *(5,312.9 fast)* | 11,891.7 MiB | **−6,366 MiB** |
| **`FINALIZATIONS`** | **6.52 s** | 6.56 s | 21.41 s | **3.28×** / **3.26×** | **1,621.9 MiB** | 6,749.5 MiB | **−5,128 MiB** |
| **Total** | **95.41 s** *(93.90 best)* | **94.57 s** | **274.88 s** | **2.88×** / **2.91×** | **6,083.8 MiB** *(5,821.2 best)* | **11,891.7 MiB** | **−5,808 MiB (1.95× lower)** |

#### Standard Whole-Program & CLI Suite

| Task | Original (`c9872522`) | **Latest `jsgo` (`t=1` / `t=4`)** | Java Cold (`[95% CI]`) | Java Warm (`[95% CI]`) | **`jsgo` vs. Cold / Warm Java (`t=4`)** | Peak RSS: `jsgo` vs. Cold / Warm Java |
|---|---:|---:|---:|---:|---:|---:|
| **`typed-app`** | 2.86 s | **1.21 s** / **1.08 s** *(`1.07 s` fast)* | 6.24 s `[6.17, 6.29]` | 1.99 s `[1.97, 2.00]` | **5.78×** / **1.84×** *(`5.83×` / `1.86×`)* | **225 MiB** vs. 749 / 1,212 MiB *(peak 1,875)* (**3.33× / 5.39× lower**, **8.33×** vs. JVM peak) |
| **`bundle`** (12.6 MB) | 11.82 s | **3.90 s** / **3.61 s** | 12.54 s `[12.47, 12.72]` | 8.97 s `[8.94, 9.03]` | **3.48×** / **2.49×** | **535 MiB** vs. 2,143 / 2,900 MiB *(peak 3,633)* (**4.01× / 5.42× lower**, **6.79×** vs. JVM peak) |
| **`downlevel-whitespace`** | 0.86 s | **0.45 s** / **0.38 s** | 2.41 s `[2.38, 2.45]` | 0.56 s `[0.55, 0.58]` | **6.33×** / **1.47×** | **138 MiB** vs. 555 / 548 MiB *(peak 1,829)* (**4.02× / 3.97× lower**, **13.25×** vs. JVM peak) |
| **`downlevel-simple`** | 2.39 s | **0.89 s** / **0.80 s** | 4.65 s `[4.55, 4.71]` | 1.53 s `[1.53, 1.54]` | **5.84×** / **1.92×** | **181 MiB** vs. 619 / 912 MiB *(peak 1,889)* (**3.42× / 5.04× lower**, **10.44×** vs. JVM peak) |
| **`checks-only`** | 4.40 s | **2.18 s** / **1.77 s** | 6.92 s `[6.87, 7.01]` | 3.07 s `[3.06, 3.09]` | **3.91×** / **1.73×** | **408 MiB** vs. 1,622 / 1,322 MiB *(peak 4,385)* (**3.98× / 3.24× lower**, **10.75×** vs. JVM peak) |
| **`single-files`** (1,000 CLI runs) | 126.90 s | **103.84 s** | — | 114.69 s *(in-JVM batch)* | **1.10×** *(vs. in-JVM batch)* | **98 MiB** vs. 1,577 MiB (**16.09× lower**) |

#### Large & Mega Real-World `ADVANCED` Workloads

| Task | Input Size & Workload (`--compilation_level=ADVANCED`) | **Latest `jsgo` (`t=1` / `t=4`)** | Java Cold (`[95% CI]`) | Java Warm (`[95% CI]`) | **`jsgo` vs. Cold / Warm Java (`t=4`)** | Peak RSS: `jsgo` vs. Cold / Warm Java |
|---|---|---:|---:|---:|---:|---:|
| **`typed-mega`** | **44.6 MB, 4,385 files** (10 Closure `ADVANCED` codebases: Closure Library + NGA OpenSphere 23.0.1 + OpenLayers 4.6.5 + Shaka Player 5.3.1 + Blockly 8.0.5 + Camptocamp ngeo 2.2.10 + YouTube SPFjs 2.4.0 + ol-cesium 1.37.0 + google-protobuf 3.21.2 + app, 94.1% typed → 14.33 MB JS + 22.29 MB map) | **15.84 s** / **14.21 s** *(`13.71 s` fast)* | 35.43 s `[35.04, 35.55]` | 27.45 s `[27.30, 27.59]` | **2.49×** / **1.93×** *(`2.58×` / `2.00×`)* | **1,225 MiB** vs. 6,581 / 5,957 MiB *(peak 7,113)* (**5.37× / 4.86× lower**, **5.81×** vs. JVM peak) |
| **`bundle-mega`** | **36.1 MB, 26 UMD libraries** (plotly.js, handsontable, echarts, ag-grid-community, pdfjs-dist, three.js, react-dom, pixi.js, fabric, mathjs, konva, moment, d3, google-libphonenumber, lodash, lovefield, vue, jquery, rxjs, google-protobuf, handlebars, immutable, chart.js, zone.js, underscore, 27,597 warnings → 13.09 MB JS + 33.99 MB map) | **20.76 s** / **19.34 s** *(`19.02 s` fast)* | 48.26 s `[48.12, 48.52]` | 41.20 s `[41.01, 42.88]` | **2.50×** / **2.13×** *(`2.54×` / `2.17×`)* | **1,614 MiB** vs. 6,666 / 6,333 MiB *(peak 7,485)* (**4.13× / 3.92× lower**, **4.64×** vs. JVM peak) |
| **`typed-full`** | **28.5 MB, 1,818 files** (Closure Library + Shaka Player 5.3.1 + Blockly 8.0.5 + app, 97.3% typed → 6.75 MB JS + 13.31 MB map) | **9.21 s** / **8.30 s** *(`7.98 s` fast)* | 23.02 s `[22.91, 23.11]` | 16.10 s `[16.00, 16.17]` | **2.78×** / **1.94×** *(`2.89×` / `2.02×`)* | **869 MiB** vs. 4,068 / 3,325 MiB *(peak 4,625)* (**4.68× / 3.83× lower**, **5.32×** vs. JVM peak) |
| **`bundle-advanced-ts`** | **12.6 MB, 4 libraries incl. 8.9 MB `typescript.js`** (10,751 warnings → 3.25 MB JS + 9.39 MB map) | **5.77 s** / **5.52 s** | 17.70 s `[17.62, 17.82]` | 11.61 s `[11.55, 11.75]` | **3.21×** / **2.10×** | **615 MiB** vs. 2,358 / 3,278 MiB *(peak 4,821)* (**3.83× / 5.33× lower**, **7.84×** vs. JVM peak) |
| **`shaka-blockly`** | **8.6 MB, 847 files** (Shaka Player 5.3.1 + Blockly 8.0.5, 93.9% typed → 1.86 MB JS + 4.94 MB map) | **3.78 s** / **3.40 s** *(`3.34 s` fast)* | 12.50 s `[12.36, 12.56]` | 6.97 s `[6.95, 7.05]` | **3.67×** / **2.05×** *(`3.75×` / `2.09×`)* | **412 MiB** vs. 1,809 / 2,861 MiB *(peak 4,586)* (**4.39× / 6.94× lower**, **11.13×** vs. JVM peak) |
| **`bundle-advanced`** | **5.4 MB, 11 UMD libraries** (three, lodash, react-dom, rxjs, d3, moment, vue, jquery, handlebars, immutable, underscore → 1.88 MB JS + 4.35 MB map) | **3.03 s** / **2.75 s** | 11.35 s `[11.25, 11.45]` | 5.79 s `[5.75, 5.82]` | **4.12×** / **2.10×** | **328 MiB** vs. 2,291 / 2,037 MiB *(peak 2,872)* (**6.98× / 6.21× lower**, **8.76×** vs. JVM peak) |

### a95c4a93 on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `a95c4a93` (`50da0981` plus two refactors: the remaining `unsafe` imports and the type-erased interfaces of the type graph, CFG and threads eliminated, and the JSType `Visitor`/`RelationshipVisitor` parameterized). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `50da0981` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was quiet and on AC power: Java's cold times are within 2% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **63.27 s** | 97.74 s | **1.54× `[1.50, 1.59]`** | **5,169 MiB** | 11,921 MiB | **2.31× `[2.14, 2.48]`** |
| **`OPTIMIZATIONS`** | **55.91 s** | 91.31 s | **1.63× `[1.59, 1.68]`** | **4,492 MiB** | 12,496 MiB | **2.78× `[2.72, 2.85]`** |
| **`FINALIZATIONS`** | **7.98 s** | 20.09 s | **2.52× `[2.48, 2.55]`** | **1,136 MiB** | 7,332 MiB | **6.47× `[5.98, 6.99]`** |
| **Total** | **127.2 s** | **209.1 s** | **1.64× `[1.62, 1.67]`** | **5,169 MiB** | **12,496 MiB** | **2.42× `[2.28, 2.56]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **28.62 s** / **25.83 s** | 49.17 s / 39.94 s | 1.90× `[1.87, 1.93]` | 1.55× `[1.53, 1.57]` | **845 MiB** vs. 4,969 / 6,598 MiB — **5.87× `[5.36, 6.43]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **20.47 s** / **17.88 s** | 35.57 s / 26.84 s | 1.99× `[1.95, 2.02]` | 1.50× `[1.47, 1.53]` | **682 MiB** vs. 3,805 / 4,280 MiB — **5.58× `[5.34, 5.82]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **11.95 s** / **10.48 s** | 23.00 s / 15.43 s | 2.19× `[2.18, 2.21]` | 1.47× `[1.44, 1.51]` | **426 MiB** vs. 3,520 / 3,675 MiB — **8.26× `[8.07, 8.46]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.78 s** / **7.33 s** | 17.65 s / 11.47 s | 2.41× `[2.39, 2.42]` | 1.56× `[1.55, 1.57]` | **310 MiB** vs. 1,726 / 2,660 MiB — **5.56× `[5.33, 5.80]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.97 s** / **4.27 s** | 12.44 s / 6.86 s | 2.91× `[2.85, 2.97]` | 1.61× `[1.57, 1.65]` | **208 MiB** vs. 1,477 / 3,518 MiB — **7.12× `[6.65, 7.63]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.98 s** / **3.38 s** | 11.41 s / 5.53 s | 3.37× `[3.31, 3.43]` | 1.63× `[1.62, 1.64]` | **187 MiB** vs. 1,587 / 1,907 MiB — **8.46× `[7.75, 9.24]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.52 s** / **1.31 s** | 6.06 s / 1.84 s | 4.61× `[4.47, 4.75]` | 1.40× `[1.34, 1.46]` | **134 MiB** vs. 684 / 1,565 MiB — **5.08× `[4.11, 6.29]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.58 s** / **2.03 s** | 6.74 s / 2.79 s | 3.32× `[3.23, 3.41]` | 1.37× `[1.29, 1.46]` | **275 MiB** vs. 1,346 / 3,516 MiB — **4.88× `[4.37, 5.46]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.95 s** / **4.46 s** | 13.13 s / 8.58 s | 2.94× `[2.85, 3.04]` | 1.92× `[1.91, 1.94]` | **229 MiB** vs. 1,655 / 1,880 MiB — **7.22× `[7.01, 7.45]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **82.52 s** | — / 89.26 s | — | 1.08× `[1.06, 1.10]` | **67 MiB** vs. — / 1,731 MiB — **25.7×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.10 s** / **0.95 s** | 4.57 s / 1.50 s | 4.80× `[4.68, 4.92]` | 1.57× `[1.48, 1.67]` | **124 MiB** vs. 605 / 1,489 MiB — **4.86× `[4.77, 4.96]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.51 s** / **0.45 s** | 2.42 s / 0.51 s | 5.43× `[5.36, 5.50]` | 1.15× `[1.07, 1.23]` | **109 MiB** vs. 518 / 1,422 MiB — **4.75× `[4.62, 4.89]`** |

#### Against `50da0981`

Performance-neutral: measured in the same blocks (ratios `50da0981 / a95c4a93`), wall time is 0.99×–1.01× and CPU time 0.99×–1.01× on every task, and peak RSS is unchanged within noise except on `downlevel-whitespace`, where it is about 5% higher (0.95× `[0.92, 0.98]`). The lower `googleapis` memory ratio against Java (2.42×, against 2.65× in the previous section) comes from Java's side: its peak RSS was lower in this run, while `jsgo`'s is unchanged.

### 50da0981 on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `50da0981` (`692a3752` plus `dc2873d7`, which removes pass micro-caches, dead code and manual slice slabs from jscomp). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `692a3752` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was quiet and on AC power: Java's cold times are within 4% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **64.23 s** | 96.81 s | **1.51× `[1.45, 1.56]`** | **5,169 MiB** | 12,982 MiB | **2.51× `[2.43, 2.59]`** |
| **`OPTIMIZATIONS`** | **56.71 s** | 90.31 s | **1.59× `[1.53, 1.65]`** | **4,545 MiB** | 13,641 MiB | **3.01× `[2.73, 3.31]`** |
| **`FINALIZATIONS`** | **8.11 s** | 20.35 s | **2.51× `[2.48, 2.54]`** | **1,162 MiB** | 7,341 MiB | **6.32× `[6.03, 6.63]`** |
| **Total** | **129.1 s** | **207.5 s** | **1.61× `[1.57, 1.65]`** | **5,169 MiB** | **13,693 MiB** | **2.65× `[2.53, 2.77]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **28.43 s** / **25.58 s** | 49.33 s / 39.87 s | 1.93× `[1.90, 1.95]` | 1.56× `[1.55, 1.57]` | **880 MiB** vs. 4,770 / 6,599 MiB — **5.42× `[5.15, 5.69]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **20.26 s** / **17.72 s** | 35.23 s / 26.20 s | 1.99× `[1.97, 2.01]` | 1.48× `[1.45, 1.51]` | **690 MiB** vs. 3,796 / 3,994 MiB — **5.50× `[5.34, 5.67]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **11.89 s** / **10.41 s** | 23.00 s / 15.61 s | 2.21× `[2.19, 2.23]` | 1.50× `[1.45, 1.55]` | **423 MiB** vs. 3,525 / 3,572 MiB — **8.34× `[8.20, 8.48]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.78 s** / **7.31 s** | 17.86 s / 11.59 s | 2.44× `[2.41, 2.47]` | 1.58× `[1.57, 1.60]` | **312 MiB** vs. 1,743 / 3,788 MiB — **5.59× `[5.44, 5.76]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.95 s** / **4.26 s** | 12.46 s / 6.61 s | 2.92× `[2.89, 2.96]` | 1.55× `[1.50, 1.60]` | **205 MiB** vs. 1,462 / 3,532 MiB — **7.13× `[6.85, 7.42]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.96 s** / **3.36 s** | 11.44 s / 5.56 s | 3.40× `[3.32, 3.48]` | 1.65× `[1.62, 1.69]` | **189 MiB** vs. 1,574 / 3,374 MiB — **8.31× `[7.62, 9.06]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.52 s** / **1.31 s** | 6.14 s / 1.80 s | 4.68× `[4.59, 4.76]` | 1.37× `[1.32, 1.43]` | **141 MiB** vs. 629 / 1,566 MiB — **4.46× `[4.13, 4.81]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.56 s** / **2.02 s** | 6.87 s / 2.79 s | 3.40× `[3.29, 3.52]` | 1.38× `[1.33, 1.44]` | **272 MiB** vs. 1,319 / 3,500 MiB — **4.85× `[4.51, 5.23]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.95 s** / **4.46 s** | 12.94 s / 8.74 s | 2.90× `[2.87, 2.94]` | 1.96× `[1.95, 1.97]` | **234 MiB** vs. 1,644 / 1,888 MiB — **7.03× `[6.62, 7.47]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **82.15 s** | — / 87.78 s | — | 1.07× `[1.06, 1.07]` | **69 MiB** vs. — / 1,131 MiB — **16.5×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.11 s** / **0.95 s** | 4.55 s / 1.47 s | 4.79× `[4.67, 4.91]` | 1.55× `[1.48, 1.62]` | **122 MiB** vs. 603 / 1,487 MiB — **4.94× `[4.69, 5.20]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.52 s** / **0.44 s** | 2.42 s / 0.51 s | 5.48× `[5.32, 5.65]` | 1.15× `[1.07, 1.22]` | **107 MiB** vs. 519 / 1,426 MiB — **4.87× `[4.68, 5.08]`** |

#### Against `692a3752`

Measured in the same blocks (ratios `692a3752 / 50da0981`), the whole-program and CLI tasks are unchanged (wall and CPU time 0.99×–1.01×, peak RSS within noise). `googleapis` is 2% slower (0.98× `[0.97, 1.00]`; 0 to 3% in every block), all of it in the `CHECKS` stage (0.97×), with CPU time and peak RSS unchanged.

### 692a3752 on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `692a3752` (`fbdac472` plus two cleanups: slot-peeking fast paths and inline scratch buffers removed from jutil, rhino and jscomp, and the inlining walkers deduplicated with the bespoke startup and pass slabs removed). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `fbdac472` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was quiet and on AC power: Java's cold times are within 4.6% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **62.77 s** | 93.79 s | **1.49× `[1.46, 1.53]`** | **5,203 MiB** | 13,305 MiB | **2.56× `[2.38, 2.75]`** |
| **`OPTIMIZATIONS`** | **55.31 s** | 87.62 s | **1.58× `[1.58, 1.59]`** | **4,502 MiB** | 13,933 MiB | **3.10× `[2.99, 3.20]`** |
| **`FINALIZATIONS`** | **7.95 s** | 19.96 s | **2.51× `[2.49, 2.53]`** | **1,137 MiB** | 7,358 MiB | **6.48× `[5.84, 7.20]`** |
| **Total** | **126.0 s** | **201.4 s** | **1.60× `[1.58, 1.61]`** | **5,203 MiB** | **13,933 MiB** | **2.68× `[2.54, 2.83]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **28.48 s** / **25.61 s** | 49.15 s / 39.92 s | 1.92× `[1.90, 1.93]` | 1.56× `[1.55, 1.57]` | **854 MiB** vs. 5,117 / 6,587 MiB — **5.98× `[5.59, 6.40]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **20.25 s** / **17.72 s** | 35.07 s / 26.63 s | 1.98× `[1.97, 1.99]` | 1.50× `[1.48, 1.52]` | **711 MiB** vs. 3,805 / 5,053 MiB — **5.36× `[5.07, 5.66]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **11.90 s** / **10.40 s** | 23.02 s / 15.47 s | 2.21× `[2.18, 2.24]` | 1.49× `[1.44, 1.54]` | **419 MiB** vs. 3,535 / 3,573 MiB — **8.43× `[8.23, 8.64]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.78 s** / **7.32 s** | 17.80 s / 11.45 s | 2.43× `[2.40, 2.47]` | 1.56× `[1.55, 1.58]` | **312 MiB** vs. 1,731 / 2,342 MiB — **5.55× `[5.36, 5.75]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.96 s** / **4.27 s** | 12.42 s / 6.76 s | 2.91× `[2.86, 2.96]` | 1.58× `[1.57, 1.60]` | **205 MiB** vs. 1,474 / 3,471 MiB — **7.20× `[6.91, 7.51]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.97 s** / **3.38 s** | 11.52 s / 5.51 s | 3.41× `[3.37, 3.45]` | 1.63× `[1.61, 1.65]` | **189 MiB** vs. 1,559 / 1,931 MiB — **8.23× `[7.85, 8.63]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.52 s** / **1.31 s** | 6.06 s / 1.95 s | 4.62× `[4.49, 4.75]` | 1.48× `[1.44, 1.52]` | **141 MiB** vs. 675 / 1,563 MiB — **4.77× `[4.08, 5.57]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.56 s** / **2.03 s** | 6.85 s / 2.81 s | 3.37× `[3.24, 3.51]` | 1.38× `[1.32, 1.45]` | **271 MiB** vs. 1,380 / 3,554 MiB — **5.09× `[4.65, 5.56]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.94 s** / **4.46 s** | 12.93 s / 8.76 s | 2.90× `[2.84, 2.96]` | 1.97× `[1.96, 1.98]` | **230 MiB** vs. 1,667 / 2,337 MiB — **7.24× `[6.84, 7.65]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **82.00 s** | — / 87.98 s | — | 1.07× `[1.07, 1.08]` | **67 MiB** vs. — / 2,387 MiB — **35.6×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.11 s** / **0.96 s** | 4.47 s / 1.48 s | 4.65× `[4.57, 4.74]` | 1.54× `[1.48, 1.60]` | **123 MiB** vs. 607 / 1,492 MiB — **4.95× `[4.76, 5.16]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.52 s** / **0.44 s** | 2.33 s / 0.51 s | 5.30× `[5.17, 5.43]` | 1.15× `[1.07, 1.24]` | **107 MiB** vs. 515 / 1,425 MiB — **4.80× `[4.67, 4.93]`** |

#### Against `fbdac472`

Performance-neutral in time: measured in the same blocks (ratios `fbdac472 / 692a3752`), wall time is 0.99×–1.03× and CPU time 0.98×–1.02× on every task. Peak RSS is unchanged within noise except on `typed-app`, where it is 14% lower (1.16× `[1.06, 1.28]`, 141 MiB against about 165 MiB).

### fbdac472 on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `fbdac472` (`41b39832` plus one change: `AbstractVarBase.nameHash` is no longer reused as per-pass scratch space). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `41b39832` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was quiet and on AC power: Java's cold times are within 4.4% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **63.11 s** | 94.30 s | **1.49× `[1.46, 1.53]`** | **5,346 MiB** | 13,388 MiB | **2.50× `[2.42, 2.59]`** |
| **`OPTIMIZATIONS`** | **55.35 s** | 88.30 s | **1.60× `[1.58, 1.61]`** | **4,604 MiB** | 13,938 MiB | **3.03× `[2.84, 3.23]`** |
| **`FINALIZATIONS`** | **7.92 s** | 19.97 s | **2.52× `[2.50, 2.55]`** | **1,111 MiB** | 7,427 MiB | **6.69× `[6.17, 7.25]`** |
| **Total** | **126.4 s** | **202.6 s** | **1.60× `[1.58, 1.63]`** | **5,346 MiB** | **13,938 MiB** | **2.61× `[2.58, 2.63]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **28.44 s** / **25.55 s** | 49.32 s / 40.54 s | 1.93× `[1.91, 1.95]` | 1.59× `[1.58, 1.59]` | **873 MiB** vs. 4,844 / 6,584 MiB — **5.55× `[5.33, 5.79]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **20.27 s** / **17.73 s** | 35.24 s / 25.90 s | 1.99× `[1.97, 2.01]` | 1.46× `[1.45, 1.47]` | **700 MiB** vs. 3,794 / 3,778 MiB — **5.42× `[5.25, 5.59]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **12.02 s** / **10.45 s** | 23.27 s / 15.31 s | 2.23× `[2.18, 2.27]` | 1.46× `[1.43, 1.50]` | **425 MiB** vs. 3,519 / 3,578 MiB — **8.28× `[8.04, 8.54]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.80 s** / **7.33 s** | 17.66 s / 11.52 s | 2.41× `[2.38, 2.44]` | 1.57× `[1.55, 1.59]` | **311 MiB** vs. 1,741 / 3,292 MiB — **5.60× `[5.39, 5.82]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.98 s** / **4.27 s** | 12.36 s / 6.74 s | 2.89× `[2.84, 2.95]` | 1.58× `[1.55, 1.61]` | **203 MiB** vs. 1,466 / 3,537 MiB — **7.24× `[7.02, 7.47]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.95 s** / **3.36 s** | 11.28 s / 5.60 s | 3.36× `[3.30, 3.41]` | 1.67× `[1.63, 1.71]` | **190 MiB** vs. 1,533 / 3,357 MiB — **8.06× `[7.37, 8.82]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.51 s** / **1.31 s** | 6.00 s / 1.90 s | 4.57× `[4.51, 4.62]` | 1.44× `[1.41, 1.48]` | **167 MiB** vs. 628 / 1,571 MiB — **3.77× `[3.53, 4.03]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.56 s** / **2.02 s** | 6.55 s / 2.77 s | 3.24× `[3.16, 3.33]` | 1.37× `[1.30, 1.45]` | **269 MiB** vs. 1,368 / 3,523 MiB — **5.09× `[4.79, 5.40]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.97 s** / **4.48 s** | 13.00 s / 8.91 s | 2.90× `[2.87, 2.94]` | 1.99× `[1.96, 2.02]` | **235 MiB** vs. 1,672 / 1,873 MiB — **7.13× `[6.81, 7.46]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **81.78 s** | — / 86.49 s | — | 1.06× `[1.05, 1.06]` | **68 MiB** vs. — / 1,201 MiB — **17.7×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.10 s** / **0.95 s** | 4.48 s / 1.48 s | 4.73× `[4.58, 4.89]` | 1.56× `[1.50, 1.63]` | **124 MiB** vs. 606 / 1,481 MiB — **4.89× `[4.66, 5.14]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.51 s** / **0.44 s** | 2.41 s / 0.51 s | 5.44× `[5.19, 5.70]` | 1.14× `[1.05, 1.23]` | **107 MiB** vs. 520 / 1,427 MiB — **4.86× `[4.72, 5.01]`** |

#### Against `41b39832`

Performance-neutral: measured in the same blocks (ratios `41b39832 / fbdac472`), wall time is 0.98×–1.00× (`googleapis` 0.98× `[0.97, 1.00]`), CPU time 0.99×–1.00× and peak RSS 0.97×–1.03× on every task.

### 41b39832 on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `41b39832`: `a8c90bf9` plus 49 commits that port about twenty more passes to the typed AST view (among them `TypeCheck`, `ControlFlowAnalysis`, `SyntacticScopeCreator`, `ClosureCheckModule` and `ProcessDefines`), replace bespoke structures with plain Go ones (`PolyfillUsageFinder` bitmasks and Bloom filters, the `StringPool` id table, `RemoveUnusedCode`'s bitmask tables) and page the original-name column. Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`; `jsgo` is built from the tree as committed, and the previous section's build (`a8c90bf9`) ran in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was quiet and on AC power throughout: Java's cold times are within 6% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **62.86 s** | 94.96 s | **1.51× `[1.46, 1.56]`** | **5,304 MiB** | 13,182 MiB | **2.48× `[2.40, 2.58]`** |
| **`OPTIMIZATIONS`** | **55.03 s** | 88.00 s | **1.60× `[1.59, 1.61]`** | **4,749 MiB** | 13,864 MiB | **2.92× `[2.77, 3.07]`** |
| **`FINALIZATIONS`** | **7.79 s** | 19.91 s | **2.55× `[2.53, 2.58]`** | **1,154 MiB** | 7,308 MiB | **6.33× `[6.09, 6.58]`** |
| **Total** | **125.7 s** | **202.9 s** | **1.61× `[1.59, 1.64]`** | **5,304 MiB** | **13,873 MiB** | **2.62× `[2.54, 2.69]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **28.29 s** / **25.39 s** | 49.17 s / 40.12 s | 1.94× `[1.92, 1.96]` | 1.58× `[1.57, 1.59]` | **850 MiB** vs. 4,778 / 6,595 MiB — **5.62× `[5.39, 5.86]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **20.17 s** / **17.68 s** | 35.31 s / 25.90 s | 2.00× `[1.99, 2.00]` | 1.46× `[1.43, 1.50]` | **715 MiB** vs. 3,799 / 4,394 MiB — **5.32× `[5.03, 5.62]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **11.86 s** / **10.38 s** | 23.07 s / 15.30 s | 2.22× `[2.19, 2.25]` | 1.47× `[1.45, 1.50]` | **429 MiB** vs. 3,525 / 3,683 MiB — **8.22× `[7.97, 8.48]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.75 s** / **7.29 s** | 17.62 s / 11.46 s | 2.42× `[2.38, 2.45]` | 1.57× `[1.56, 1.58]` | **316 MiB** vs. 1,733 / 2,976 MiB — **5.49× `[5.22, 5.78]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.94 s** / **4.24 s** | 12.54 s / 6.67 s | 2.96× `[2.92, 3.00]` | 1.57× `[1.54, 1.60]` | **201 MiB** vs. 1,456 / 3,522 MiB — **7.25× `[7.02, 7.49]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.93 s** / **3.34 s** | 11.43 s / 5.58 s | 3.42× `[3.38, 3.46]` | 1.67× `[1.65, 1.69]` | **195 MiB** vs. 1,642 / 1,985 MiB — **8.41× `[7.72, 9.15]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.50 s** / **1.30 s** | 6.13 s / 1.85 s | 4.72× `[4.61, 4.82]` | 1.42× `[1.37, 1.47]` | **159 MiB** vs. 628 / 1,553 MiB — **3.95× `[3.74, 4.16]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.57 s** / **2.02 s** | 6.85 s / 2.84 s | 3.38× `[3.28, 3.49]` | 1.40× `[1.34, 1.46]` | **257 MiB** vs. 1,340 / 3,554 MiB — **5.20× `[4.72, 5.73]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.92 s** / **4.43 s** | 12.90 s / 8.51 s | 2.91× `[2.85, 2.97]` | 1.92× `[1.91, 1.94]` | **222 MiB** vs. 1,653 / 1,804 MiB — **7.45× `[7.09, 7.82]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **81.30 s** | — / 87.74 s | — | 1.08× `[1.07, 1.08]` | **67 MiB** vs. — / 1,282 MiB — **19.0×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.10 s** / **0.95 s** | 4.48 s / 1.49 s | 4.73× `[4.57, 4.89]` | 1.57× `[1.49, 1.66]` | **126 MiB** vs. 604 / 1,482 MiB — **4.77× `[4.66, 4.90]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.51 s** / **0.44 s** | 2.41 s / 0.50 s | 5.51× `[5.35, 5.66]` | 1.14× `[1.07, 1.22]` | **107 MiB** vs. 527 / 1,424 MiB — **4.91× `[4.73, 5.10]`** |

#### Against `a8c90bf9`

Measured in the same blocks (ratios `a8c90bf9 / 41b39832`, paired 95% CIs), those commits cost 0–3% of wall time (0.97×–1.00×; `googleapis` 0.97× `[0.96, 0.98]`, `shaka-blockly` and `typed-mega` 0.98×, the rest 0.99×–1.00×) and 0–2% of CPU time, with peak RSS unchanged within noise on every task (0.95×–1.04×).

### a8c90bf9 on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `a8c90bf9`: the node-id migration (`rhino.Node` is a 32-bit id into on-demand record blocks, 36-byte node records, in-place compaction, a typed AST view with InlineVariables, ConvertToDottedProperties and other passes ported to it, a regenerated `default.pgo`) and the cleanups that followed it. Same machine, method, Java (v20261021 on OpenJDK 27.0.1) and `--num_parallel_threads=4` as [the previous section](#v20261021-on-jdk-27-apple-m2); `jsgo` is built from the tree as committed. The previous section's build (`83d55bdc`) ran in the same blocks as a third contender.

**Gate.** As in the previous section: every `jsgo` run of every block was byte-identical to Java (stdout, stderr, exit status, `out.js`, `out.js.map`; for `googleapis` all three stages and the 7,483 chunk files), with the JVM's own `WARNING:` lines dropped from Java's stderr and, in `single-files`, the Java stack frames of the same 11 inputs on which both compilers crash with the same exception (see the previous section's note).

**Conditions.** The first pass ran partly on battery and while another session ran builds and tests and a Chrome tab used about two cores; Java's cold times then drifted up to 21% above the previous section's. `shaka-blockly`, `bundle`, `bundle-advanced-ts`, `typed-full`, `typed-mega`, `bundle-mega`, `googleapis` and `single-files` were rerun on AC power with nothing else running, and those reruns are what is reported; Java's cold times are now within 7% of the previous section's on every task (`googleapis` 204.3 s against 201.3 s).

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **62.31 s** | 95.84 s | **1.54× `[1.49, 1.59]`** | **5,162 MiB** | 12,893 MiB | **2.50× `[2.37, 2.63]`** |
| **`OPTIMIZATIONS`** | **53.48 s** | 88.49 s | **1.65× `[1.64, 1.67]`** | **4,577 MiB** | 13,632 MiB | **2.98× `[2.69, 3.31]`** |
| **`FINALIZATIONS`** | **7.68 s** | 20.00 s | **2.60× `[2.58, 2.63]`** | **1,140 MiB** | 7,164 MiB | **6.30× `[5.59, 7.10]`** |
| **Total** | **123.5 s** | **204.3 s** | **1.65× `[1.63, 1.68]`** | **5,162 MiB** | **13,632 MiB** | **2.64× `[2.49, 2.80]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **27.94 s** / **25.08 s** | 49.39 s / 40.30 s | 1.97× `[1.96, 1.98]` | 1.61× `[1.59, 1.62]` | **842 MiB** vs. 4,517 / 6,595 MiB — **5.35× `[4.89, 5.87]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **20.13 s** / **17.50 s** | 35.54 s / 26.34 s | 2.03× `[2.00, 2.06]` | 1.51× `[1.49, 1.52]` | **738 MiB** vs. 3,791 / 3,924 MiB — **5.14× `[4.87, 5.42]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **11.87 s** / **10.45 s** | 23.69 s / 15.57 s | 2.27× `[2.19, 2.35]` | 1.49× `[1.42, 1.56]` | **421 MiB** vs. 3,535 / 3,713 MiB — **8.39× `[8.08, 8.71]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.75 s** / **7.30 s** | 18.03 s / 11.55 s | 2.47× `[2.43, 2.51]` | 1.58× `[1.56, 1.60]` | **314 MiB** vs. 1,752 / 3,814 MiB — **5.58× `[5.26, 5.93]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.99 s** / **4.27 s** | 12.90 s / 6.97 s | 3.02× `[2.97, 3.07]` | 1.63× `[1.61, 1.65]` | **204 MiB** vs. 1,480 / 3,461 MiB — **7.26× `[7.21, 7.31]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **4.01 s** / **3.37 s** | 12.04 s / 5.62 s | 3.57× `[3.46, 3.68]` | 1.67× `[1.64, 1.70]` | **187 MiB** vs. 1,629 / 1,878 MiB — **8.68× `[7.89, 9.56]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.52 s** / **1.31 s** | 6.40 s / 1.90 s | 4.88× `[4.68, 5.08]` | 1.45× `[1.42, 1.48]` | **163 MiB** vs. 674 / 1,566 MiB — **4.10× `[3.58, 4.71]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.63 s** / **2.04 s** | 7.23 s / 2.93 s | 3.54× `[3.43, 3.66]` | 1.44× `[1.40, 1.47]` | **271 MiB** vs. 1,334 / 3,492 MiB — **4.93× `[4.51, 5.38]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **5.03 s** / **4.49 s** | 13.27 s / 8.79 s | 2.95× `[2.89, 3.01]` | 1.96× `[1.94, 1.97]` | **222 MiB** vs. 1,635 / 1,849 MiB — **7.38× `[6.79, 8.02]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **81.20 s** | — / 87.28 s | — | 1.07× `[1.07, 1.08]` | **68 MiB** vs. — / 2,321 MiB — **34.0×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.10 s** / **0.94 s** | 4.74 s / 1.47 s | 5.05× `[4.79, 5.33]` | 1.57× `[1.49, 1.65]` | **128 MiB** vs. 604 / 1,497 MiB — **4.73× `[4.66, 4.80]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.51 s** / **0.45 s** | 2.48 s / 0.55 s | 5.51× `[5.04, 6.01]` | 1.22× `[1.10, 1.35]` | **108 MiB** vs. 495 / 1,433 MiB — **4.57× `[3.97, 5.26]`** |

#### Against `83d55bdc`

Measured in the same blocks (ratios `83d55bdc / a8c90bf9`, paired 95% CIs), `a8c90bf9` is 3–11% slower on every task (wall time 0.90×–0.97×; `googleapis` 0.90× `[0.87, 0.93]`, 123.5 s against 111.1 s; `typed-full` 0.90× `[0.87, 0.94]`; `bundle-mega` 0.91× `[0.90, 0.91]`) and uses 5–11% more CPU, while its peak RSS is up to 31% lower on the larger workloads (`googleapis` 5.04 GiB against 6.33 GiB, 1.26× `[1.21, 1.31]`; `bundle` 1.44× `[1.26, 1.65]`; `typed-full` 1.19× `[1.17, 1.22]`) and about the same on the `downlevel` tasks, `bundle-advanced` and `bundle-advanced-ts`. Against Java, the speedup on `googleapis` moves from 1.74× to 1.65× and the memory ratio from 2.05× to 2.64×.

The regenerated `cmd/jsgo/default.pgo` of the node-id branch (`a93368f4`) is as good as Round 3A's profile: in an earlier, partly disturbed pass on `a8b05ce3`, a build of the same tree with the Round 3A profile was within 2% of it on seven of the eight smaller whole-program tasks and 1.05× `[0.97, 1.13]` on `shaka-blockly`.

### v20261021 on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, after the upgrade to v20261021 and the memory round 2 merge. Measured at commit `83d55bdc` (`jsgo` built from the tree as committed, so with the committed `cmd/jsgo/default.pgo`) against Maven Central's `closure-compiler-v20261021.jar`, on Apple M2 (macOS arm64, 24 GiB RAM; Go 1.27.3; OpenJDK 27.0.1). Java and `jsgo` both run `--num_parallel_threads=4`; the `t=1` column is `jsgo` at `--num_parallel_threads=1`.

**Method.** [`tools/bench/interleaved.py`](../tools/bench/interleaved.py) runs and gates the suite and [`interleaved_stats.py`](../tools/bench/interleaved_stats.py) computes the statistics. Every task runs $N = 5$ blocks; each block runs every contender once, in the order Java, `jsgo` builds in odd blocks and the reverse in even ones, which cancels linear thermal and background drift. Times and peak RSS (`wait4`'s `ru_maxrss`) are sample means. Ratios (`Java / jsgo`) are geometric means with two-sided 95% confidence intervals: paired Student-$t$ ($df = 4$) on per-block log-ratios against cold Java, and Welch's $t$ on log-times against warm Java. Warm Java is `tools/bench/java/CliBatch` in one JVM: 5 measured iterations after 5 warm-up iterations (2 for `single-files`); its peak RSS is the JVM's over its whole life. `single-files` (1,000 invocations) is compared with the warm batch only, as before. `googleapis` is the three-stage build of `tools/bench/run_googleapis_bench.py` (Java with `-Xmx20g`); its rows are per stage and their sum.

**Gate.** Every `jsgo` run of every block is compared byte for byte with Java's first cold run (for `single-files`, the warm driver's last iteration): stdout, stderr, exit status and the output files (`out.js`, `out.js.map`), and for `googleapis` each stage's stdout, stderr and exit status plus a digest of all 7,483 chunk files. Two things are set aside:
- the JVM's own `WARNING:` lines at the top of Java's stderr (the `sun.misc.Unsafe` deprecation notice that JDK 26 and 27 print for the bundled protobuf), which are not compiler output;
- in `single-files`, the Java stack frames of the 11 inputs on which both compilers crash with the same exception and exit status 254 (ten standalone Closure Library files such as `closure/goog/a11y/aria/roles.js`, where `RemoveUnusedCode.getVarForNameNode` throws a `NullPointerException` on the undeclared `goog`, and one test262 file, `pass-explicit/0d137e8a97ffe083.js`, an `IllegalStateException` for a `yield` outside a generator), bench.py's `java-frames` class. (This note first said one input: the gate listed only the first three differences.)

Everything else matched on every run.

**JDK 27.** The v20261015 jar gives byte-identical output on JDK 26.0.2.1 and 27.0.1 for ten of these tasks, so the JDK upgrade alone does not change Java's output. Java's cold times are within 6% of the Round 3A section's (JDK 26, v20261015), except `typed-app`, `shaka-blockly` and `bundle-advanced`, which are 13–17% faster now.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **55.88 s** | 93.65 s | **1.68× `[1.62, 1.73]`** | **6,904 MiB** | 14,135 MiB | **2.05× `[1.92, 2.19]`** |
| **`OPTIMIZATIONS`** | **52.47 s** | 87.84 s | **1.68× `[1.51, 1.86]`** | **5,741 MiB** | 14,580 MiB | **2.58× `[2.05, 3.26]`** |
| **`FINALIZATIONS`** | **7.50 s** | 19.83 s | **2.64× `[2.62, 2.67]`** | **1,393 MiB** | 7,370 MiB | **5.29× `[5.13, 5.46]`** |
| **Total** | **115.9 s** | **201.3 s** | **1.74× `[1.66, 1.82]`** | **7,131 MiB** | **14,580 MiB** | **2.05× `[1.88, 2.24]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **25.47 s** / **22.65 s** | 49.46 s / 40.72 s | 2.18× `[2.17, 2.20]` | 1.80× `[1.78, 1.81]` | **898 MiB** vs. 5,015 / 6,595 MiB — **5.59× `[5.29, 5.91]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **18.40 s** / **15.74 s** | 35.13 s / 25.95 s | 2.23× `[2.21, 2.26]` | 1.65× `[1.63, 1.67]` | **864 MiB** vs. 3,804 / 3,980 MiB — **4.40× `[4.30, 4.50]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **10.83 s** / **9.32 s** | 23.09 s / 15.54 s | 2.48× `[2.45, 2.50]` | 1.67× `[1.61, 1.72]` | **503 MiB** vs. 3,514 / 3,596 MiB — **6.99× `[6.88, 7.10]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.08 s** / **6.66 s** | 17.77 s / 11.33 s | 2.67× `[2.63, 2.71]` | 1.70× `[1.69, 1.72]` | **333 MiB** vs. 2,063 / 3,868 MiB — **5.95× `[4.22, 8.39]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.44 s** / **3.75 s** | 12.48 s / 6.69 s | 3.33× `[3.26, 3.39]` | 1.78× `[1.72, 1.85]` | **258 MiB** vs. 1,462 / 3,488 MiB — **5.69× `[4.97, 6.50]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.61 s** / **3.05 s** | 11.45 s / 5.49 s | 3.76× `[3.68, 3.84]` | 1.80× `[1.79, 1.81]` | **193 MiB** vs. 1,582 / 1,957 MiB — **8.21× `[7.60, 8.87]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.42 s** / **1.19 s** | 6.15 s / 1.94 s | 5.18× `[5.03, 5.32]` | 1.63× `[1.58, 1.69]` | **176 MiB** vs. 630 / 2,045 MiB — **3.59× `[3.47, 3.70]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.36 s** / **1.92 s** | 6.88 s / 2.75 s | 3.58× `[3.45, 3.71]` | 1.43× `[1.36, 1.49]` | **314 MiB** vs. 1,352 / 3,523 MiB — **4.31× `[4.09, 4.55]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.49 s** / **4.04 s** | 12.86 s / 8.64 s | 3.19× `[3.13, 3.24]` | 2.14× `[2.12, 2.16]` | **337 MiB** vs. 1,659 / 3,387 MiB — **4.92× `[4.67, 5.19]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **80.17 s** | — / 87.93 s | — | 1.10× `[1.09, 1.11]` | **73 MiB** vs. — / 2,329 MiB — **31.9×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.06 s** / **0.86 s** | 4.53 s / 1.46 s | 5.28× `[5.19, 5.37]` | 1.71× `[1.61, 1.81]` | **128 MiB** vs. 604 / 1,489 MiB — **4.73× `[4.54, 4.93]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.47 s** / **0.41 s** | 2.33 s / 0.50 s | 5.65× `[5.55, 5.75]` | 1.20× `[1.14, 1.27]` | **107 MiB** vs. 514 / 1,425 MiB — **4.81× `[4.66, 4.97]`** |

#### `googleapis` since Round 3A

At `83d55bdc`, `googleapis` takes 115.9 s against Round 3A's 98.4 s, and its speedup over Java falls from 2.11× to 1.74×. Two things contribute:

- **Code.** The `jsgo` builds of main's commits since Round 3A, measured alone (2 interleaved blocks, no Java, fast-mode runs, see below): 97.6 s before the merge (`1d2b3149`); 102.9 s after the memory round 2 merge (`540d836f`, +5%); 106.8 s after main's first refactors that replace `unsafe` tricks with standard library code (`e1a9dd5f`, `0ed5c4cc`, `81d319ca`; +4%, mostly in `CHECKS`); flat until `25688342` (`unsafe` removed from JSType and PropertyMap), from which it is 109.4–110.0 s (+3%). The v20261021 upgrade adds nothing measurable.
- **Variance.** The same binary's `OPTIMIZATIONS` stage takes either about 46 s or about 54 s. The slow mode appeared in 4 of the 5 blocks above, where each `jsgo` run follows or precedes a 200 s, 14 GiB Java run, in 6 of the 22 `jsgo`-only runs, and in none of three back-to-back runs. The machine is a fanless MacBook Air (`Mac14,2`), so thermal throttling under sustained load is the likely cause; the slow mode is not tied to a commit.

#### Effect of the memory round 2 merge (`jsgo` only)

The memory round 2 work ([docs/memory-research/](memory-research/): nodes in 32-bit-id chunks with no Go pointers, handles for types, colors and JSDoc, per-pass slabs, AST compaction) was merged into main at `540d836f`. Measured on the same machine, $N = 5$ interleaved blocks of the two builds alone (no Java), each gated against the v20261015 jar's output: `1d2b3149` (main before the merge) against `540d836f` (the merge), both with the committed `default.pgo`, both at `--num_parallel_threads=4`. Ratios are `before / after`: above 1 means the merge is faster or smaller.

| Task | Before the merge | After the merge | Wall time (`[95% CI]`) | CPU time (`[95% CI]`) | Peak RSS (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|
| `googleapis` | 97.72 s | 103.1 s | 0.95× `[0.93, 0.96]` | 1.14× `[1.13, 1.16]` | 1.11× `[1.02, 1.22]` |
| `bundle-mega` | 21.94 s | 22.38 s | 0.98× `[0.91, 1.04]` | 1.16× `[1.10, 1.22]` | 1.02× `[0.92, 1.12]` |
| `typed-mega` | 18.45 s | 17.42 s | 1.07× `[0.90, 1.26]` | 1.22× `[1.07, 1.40]` | 1.08× `[1.05, 1.11]` |
| `typed-full` | 11.37 s | 10.38 s | 1.08× `[0.99, 1.17]` | 1.24× `[1.20, 1.28]` | 1.15× `[1.08, 1.22]` |
| `bundle-advanced-ts` | 6.59 s | 6.74 s | 0.98× `[0.96, 0.99]` | 1.15× `[1.13, 1.17]` | 1.02× `[0.99, 1.06]` |
| `shaka-blockly` | 4.59 s | 4.48 s | 1.02× `[0.97, 1.08]` | 1.16× `[1.10, 1.22]` | 1.07× `[0.96, 1.20]` |
| `bundle-advanced` | 2.96 s | 3.11 s | 0.96× `[0.89, 1.03]` | 1.12× `[1.03, 1.21]` | 0.98× `[0.91, 1.05]` |
| `typed-app` | 2.01 s | 1.84 s | 1.13× `[0.90, 1.41]` | 1.22× `[1.03, 1.43]` | 1.11× `[1.04, 1.19]` |
| `checks-only` | 1.99 s | 2.04 s | 0.97× `[0.89, 1.07]` | 1.14× `[1.03, 1.25]` | 1.10× `[1.01, 1.21]` |
| `bundle` | 4.44 s | 4.54 s | 0.98× `[0.96, 1.00]` | 1.16× `[1.15, 1.18]` | 0.90× `[0.82, 0.98]` |
| `downlevel-simple` | 1.20 s | 1.27 s | 0.94× `[0.77, 1.14]` | 1.08× `[0.96, 1.21]` | 1.04× `[1.02, 1.06]` |
| `downlevel-whitespace` | 0.68 s | 0.67 s | 1.01× `[0.96, 1.07]` | 1.17× `[1.11, 1.23]` | 1.02× `[0.97, 1.06]` |

On this 8-core machine at 4 threads the merge cuts CPU time by 8–24% and peak RSS by up to 15% (except `bundle`, 10% higher), but wall time hardly moves: `googleapis` is 5% slower, `bundle` and `bundle-advanced-ts` about 2% slower, `typed-full` about 8% faster, the rest within noise. A likely explanation is that the garbage-collector work the merge removes ran on otherwise idle cores here, mostly off the critical path, while the id-based node accessors add a little to it. At `--num_parallel_threads=1` (measured only with a profile regenerated for the merged code, which is mostly a few percent slower than the committed one, see below) wall time is between 12% slower (`typed-app`) and 4% faster, mostly 0–5% slower, and CPU time is 9–19% lower (`typed-app` 2% higher). The branch's own measurement ([round3-flat-slices.md §7](memory-research/round3-flat-slices.md), −4% to −15% wall time) was taken on a 4-core Linux amd64 machine at `--num_parallel_threads=1`, where CPU time is the bottleneck; the gains do not carry over to this machine.

#### PGO profile

`tools/bench/pgo.py` regenerates `cmd/jsgo/default.pgo` from `t=1` runs of the bench tasks. A profile regenerated that way for `83d55bdc` was 3–4% slower than the committed profile (from Round 3A) on every task, with tight CIs (e.g. `bundle` 0.96× `[0.96, 0.96]`, `googleapis` 0.95× `[0.90, 1.00]`); one regenerated for the merge alone was mostly slower too (−7% to +6%). The committed profile is kept: regenerating it with `pgo.py` as it stands would be a regression on this machine.

### Round 3A + Memory/GC Optimizations ($N = 5$ Counterbalanced Interleaved Suite with 95% CIs)

Measured at commit `30771cf2` against `closure-compiler-v20261015.jar` on Apple M2 (macOS arm64, 24 GiB RAM; Go 1.27.3 with PGO enabled via `cmd/jsgo/default.pgo`, OpenJDK 26.0.2.1), with both `jsgo` (`t=4`) and Java running `--num_parallel_threads=4`, gated on 100% byte-identical output on every run.

**Statistical methodology ($N = 5$):** Every benchmark runs **$N = 5$ counterbalanced (`AB / BA`) interleaved blocks** to cancel linear thermal and background drift. Wall times and Peak RSS report sample means; speedups and memory ratios (`Java / jsgo`) report geometric means with **two-sided 95% confidence intervals** (`[lo, hi]`): paired Student-$t$ ($df = 4$) on log-ratios for Cold, Batch, Staged, and Peak RSS comparisons, and two-sample Welch-$t$ on log-times for Warm JVM comparisons (5 measured iterations after 5 in-JVM JIT warmup iterations).

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | `--fast_mode` | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **50.08 s** | **48.53 s** | 98.76 s | **1.97×** `[1.93, 2.02]` *(2.03× `[1.96, 2.11]`)* | **6,610 MiB** *(6,291)* | 13,272 MiB | **2.01×** `[1.93, 2.09]` *(−6.51 GiB)* |
| **`OPTIMIZATIONS`** | **42.03 s** | **41.46 s** | 88.90 s | **2.12×** `[2.09, 2.14]` *(2.14× `[2.10, 2.19]`)* | **4,540 MiB** *(4,518)* | 14,424 MiB | **3.18×** `[3.00, 3.37]` *(−9.65 GiB)* |
| **`FINALIZATIONS`** | **6.27 s** | 6.32 s | 20.22 s | **3.22×** `[3.18, 3.26]` *(3.20× `[3.17, 3.22]`)* | **1,325 MiB** *(1,349)* | 6,884 MiB | **5.20×** `[4.94, 5.48]` *(−5.43 GiB)* |
| **Total** | **98.38 s** | **96.31 s** | **207.88 s** | **2.11×** `[2.08, 2.14]` *(2.16× `[2.12, 2.20]`)* | **6,610 MiB** *(6,291)* | **14,424 MiB** | **2.18×** `[2.11, 2.25]` *(−7.63 GiB)* |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Parsed Nodes & Input | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | 4.19M nodes *(36.1 MB, 26 UMD libs)* | **23.08 s** / **20.59 s** *(`20.17 s` fast)* | 49.91 s / 40.82 s | **2.42×** `[2.40, 2.45]` *(2.47× `[2.44, 2.51]`)* | **1.98×** `[1.97, 1.99]` *(2.02× `[2.02, 2.03]`)* | **888 MiB** vs. 5,232 / 8,709 MiB — **5.89×** `[5.53, 6.27]` *(9.8× warm)* |
| **`typed-mega`** | 3.55M nodes *(44.6 MB, 4,385 files)* | **17.97 s** / **15.41 s** *(`14.81 s` fast)* | 34.94 s / 26.74 s | **2.27×** `[2.24, 2.29]` *(2.36× `[2.33, 2.38]`)* | **1.73×** `[1.71, 1.76]` *(1.81× `[1.78, 1.83]`)* | **862 MiB** vs. 3,800 / 4,634 MiB — **4.41×** `[4.32, 4.50]` *(5.4× warm)* |
| **`typed-full`** | 1.90M nodes *(25.9 MB, 1,820 files)* | **10.96 s** / **9.48 s** *(`9.06 s` fast)* | 24.45 s / 15.73 s | **2.58×** `[2.44, 2.73]` *(2.70× `[2.56, 2.85]`)* | **1.66×** `[1.55, 1.78]` *(1.74× `[1.63, 1.85]`)* | **540 MiB** vs. 3,523 / 3,620 MiB — **6.53×** `[6.15, 6.92]` *(6.7× warm)* |
| **`bundle-advanced-ts`** | 1.11M nodes *(11.9 MB, 4 libs)* | **6.70 s** / **6.57 s** | 17.89 s / 13.57 s | **2.73×** `[2.52, 2.96]` | **2.07×** `[1.87, 2.29]` | **306 MiB** vs. 1,907 / 3,091 MiB — **6.22×** `[5.73, 6.75]` *(10.1× warm)* |
| **`shaka-blockly`** | 618k nodes *(8.5 MB, 848 files)* | **4.57 s** / **4.03 s** *(`3.78 s` fast)* | 14.73 s / 9.50 s | **3.66×** `[3.58, 3.73]` *(3.90× `[3.87, 3.93]`)* | **2.34×** `[1.97, 2.77]` *(2.49× `[2.10, 2.96]`)* | **267 MiB** vs. 1,627 / 3,440 MiB — **6.06×** `[5.34, 6.86]` *(12.9× warm)* |
| **`bundle-advanced`** | 573k nodes *(5.7 MB, 11 UMD libs)* | **3.33 s** / **2.94 s** | 13.14 s / 8.62 s | **4.42×** `[3.82, 5.13]` | **2.93×** `[2.61, 3.29]` | **184 MiB** vs. 1,597 / 3,494 MiB — **8.65×** `[7.67, 9.76]` *(19.0× warm)* |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | `ADVANCED` + types *(1.29M nodes)* | **1.40 s** / **1.26 s** *(`1.23 s` fast)* | 7.42 s / 2.21 s | **5.91×** `[5.88, 5.94]` *(6.06× `[5.89, 6.23]`)* | **1.76×** `[1.70, 1.83]` *(1.80× `[1.74, 1.87]`)* | **182 MiB** vs. 893 / 3,368 MiB — **4.91×** `[4.66, 5.18]` *(18.5× warm)* |
| **`checks-only`** | `--checks_only` + types *(1.29M nodes)* | **2.22 s** / **1.72 s** | 7.01 s / 2.86 s | **4.07×** `[3.95, 4.19]` | **1.66×** `[1.59, 1.73]` | **336 MiB** vs. 1,402 / 3,507 MiB — **4.18×** `[3.91, 4.47]` *(10.4× warm)* |
| **`bundle`** | `SIMPLE` *(1.11M nodes)* | **4.17 s** / **3.81 s** | 12.77 s / 8.68 s | **3.35×** `[3.30, 3.40]` | **2.28×** `[2.26, 2.29]` | **293 MiB** vs. 2,180 / 3,826 MiB — **7.22×** `[5.29, 9.86]` *(13.1× warm)* |
| **`single-files`** | 1,000 CLI runs, `SIMPLE` *(306k nodes)* | **98.54 s** | — / 112.88 s *(batch)* | — | **1.15×** `[1.13, 1.16]` *(vs. batch)* | **72.5 MiB** vs. 1,569 MiB — **21.5×** `[17.8, 25.8]` |
| **`downlevel-simple`** | `SIMPLE` → ES5 *(232k nodes)* | **0.96 s** / **0.83 s** | 4.67 s / 1.50 s | **5.64×** `[5.47, 5.82]` | **1.81×** `[1.74, 1.89]` | **132 MiB** vs. 629 / 1,485 MiB — **4.78×** `[4.66, 4.89]` *(11.3× warm)* |
| **`downlevel-whitespace`** | `WHITESPACE_ONLY` → ES5 *(232k nodes)* | **0.47 s** / **0.40 s** | 2.34 s / 0.50 s | **5.88×** `[5.79, 5.96]` | **1.26×** `[1.18, 1.36]` | **106 MiB** vs. 559 / 1,431 MiB — **5.27×** `[4.98, 5.58]` *(13.5× warm)* |

## M1 — parser (`internal/parser` vs `parsing.parser.Parser.parseProgram`)

Harnesses:

- `tools/bench/java/ParseBench.java`;
- `tools/bench/parsebench` (Go);
- `tools/bench/measure.py` (wall time, CPU time, peak RSS via `wait4`).

**Setup:**

- **Inputs:** the full pinned corpus: 14,690 files, 84.8M UTF-16 chars
  (`tools/oracle/corpus.txt`). All files are read into memory before
  timing.
- **Runs:** 8 parses of the whole corpus per process, single-threaded, with
  the compiler's ES_NEXT parser configuration.
- **Settings:** default JVM (JDK 21, G1) and Go 1.24 (GOGC=100,
  GOMAXPROCS=4) settings.
- **Machine:** 4 vCPUs, *shared with several porting agents*. Absolute
  numbers are therefore noisy; ratios are indicative. Re-run on a quiet
  machine for the M9 report.
- **Sanity check:** both report the same 2,776 parse errors.

| Per full-corpus parse | Go (first port) | Go (after `jutil` fixes, 9dfa23e) | Java, 1st iteration (cold JIT) | Java, steady state (iters 7–8) |
|---|---:|---:|---:|---:|
| Wall time | 8.2–9.7 s | **6.6–6.8 s** | 5.9 s | **2.2 s** |
| Allocated | 3.11 GB | 2.63 GB | 2.45 GB | 2.45 GB |
| GCs / GC pause | 32 / ~5 ms | 28 / ~5 ms | 11 / 141 ms | 2 / 20 ms |
| Whole process: wall / peak RSS | 69 s / 540 MB (8 iterations) | 20.4–20.7 s / 530 MB (3 iterations) | 38 s / 2.25 GB (8 iterations) | (same process) |

**Findings:**

- **The first port was 1.4× slower than cold Java and ~4× slower than warm
  Java.** A CPU profile showed ~37% of the time in allocation.
  - The single biggest avoidable cost (8%) was `jutil.CheckArgumentMsg`
    boxing its varargs on every call, even though the check passes.
  - Two more costs were UTF-16 ↔ string conversions that allocated twice,
    and an `ArrayList.Get` that could not be inlined.
- **The fixes are now conventions (PORTING.md §5):**
  - fixed-arity Preconditions, as in Guava;
  - single-allocation string conversion;
  - out-of-line failure paths.

  Result: −22% time and −15% allocation.
- **GC policy is not the gap.** `GOGC=400` gains ~5%. `GOGC=off` is much
  slower (9.8 GB RSS, page-fault bound).
- **The remaining gap to a warm JVM (~3×) is structural.** Go's allocator
  is slower than HotSpot's thread-local bump allocation, and Go has no
  escape analysis across interface calls and inlines less. A line-for-line
  port allocates what Java allocates. Options, after parity, behind the
  same zero-diff oracles:
  - slab-allocating per-token and per-position objects;
  - reducing interface dispatch in the scanner's hot path;
  - parsing independent input files concurrently. The compiler already
    treats inputs independently at this stage, and the output stays
    deterministic.
- **Go already wins on process-level costs:**
  - no JIT warm-up: Go's first iteration is as fast as its last;
  - no JVM startup;
  - ~4× less memory.

  For single-shot CLI invocations, the common case, the comparison that
  matters is against cold Java.

To reproduce:

```sh
source tools/oracle/pin.env
go build -o .cache/bench-build/parsebench ./tools/bench/parsebench
javac -cp .cache/$CLOSURE_JAR -d .cache/bench-build tools/bench/java/ParseBench.java
python3 -I tools/bench/measure.py go.rusage .cache/bench-build/parsebench 8 .cache/corpus @.cache/corpus-files.txt go-perfile.tsv
python3 -I tools/bench/measure.py java.rusage java -Xss1g -cp .cache/$CLOSURE_JAR:.cache/bench-build ParseBench 8 .cache/corpus @.cache/corpus-files.txt java-perfile.tsv
PARSEBENCH_CPUPROFILE=parse.pprof .cache/bench-build/parsebench 2 .cache/corpus @.cache/corpus-files.txt /dev/null
```

## M1 — full front end (`internal/parsing.ParserRunner_Parse` vs `parsing.ParserRunner.parse`)

Parser + IRFactory + JsDocInfoParser, with the compiler's own
`parsing.Config` (ES_NEXT, TYPES_ONLY, STOP_AFTER_ERROR, STRICT).

- **Harness:** `parsebench --mode=runner` and `ParseBench --mode=runner`.
- **Setup:** same corpus and machine caveats as above. Three porting agents
  were running during the measurement.
- **Sanity check:** both report the same 1,157 errors.

| Per full-corpus pass | Go | Java, 1st iteration (cold JIT) | Java, steady state (iters 3–8) |
|---|---:|---:|---:|
| Wall time | 12.0–13.9 s | 9.2 s | 4.2–5.0 s |
| Allocated | 4.31 GB | 4.11 GB | 4.08 GB |
| GCs / GC pause | 44 / ~5 ms | 14 / 522 ms | 5–6 / 70–160 ms |
| Whole process, 8 iterations: wall / CPU / peak RSS | 103 s / 128 s / 689 MB | | 43 s / 64 s / 1.5 GB |

**Findings:**

- **Allocation volume now matches Java's (4.3 vs 4.1 GB).** The remaining
  gap, 1.3× cold Java and ~2.7× warm Java, is allocation *cost* and code
  generation, not wasted allocation.
- **The profile has no new avoidable hot spot of the boxing kind.** Two
  structural items are recorded below.

### Known hot spots to revisit after parity

- **`rhino.RhinoStringPool`: ~7% of front-end CPU.** It interns every
  identifier through Go's `unique.Make`, which is the faithful analogue of
  Guava's `Interners.newWeakInterner()`, but its trie and weak-pointer
  bookkeeping are expensive. Go compares strings by value, so a
  per-compilation strong interner would behave identically and could
  replace it.

- `rhino.TokenStream_IsKeyword` converts every NAME to UTF-16. That is
  ~5% of `ParserRunner.parse` allocations (found while porting IRFactory).
