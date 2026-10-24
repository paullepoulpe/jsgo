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
