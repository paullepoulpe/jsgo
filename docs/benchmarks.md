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

### 90bc249a on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `90bc249a`: `3e2e776c` plus `346ce473` (Refactor Round 67: no production `jutil.Try` left, a `TypedScope` skip chain, token and `strpool.ID` fast paths in the visitors, Stage 2 and `jutil` deduplication). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `3e2e776c` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power and kept awake; Java's cold times are within 1.9% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **53.02 s** | 95.68 s | **1.80× `[1.71, 1.90]`** | **5,642 MiB** | 11,807 MiB | **2.09× `[1.95, 2.25]`** |
| **`OPTIMIZATIONS`** | **50.11 s** | 89.34 s | **1.78× `[1.76, 1.81]`** | **6,191 MiB** | 12,596 MiB | **2.04× `[1.92, 2.16]`** |
| **`FINALIZATIONS`** | **6.70 s** | 19.86 s | **2.96× `[2.95, 2.98]`** | **1,110 MiB** | 7,249 MiB | **6.56× `[5.86, 7.33]`** |
| **Total** | **109.8 s** | **204.9 s** | **1.87× `[1.82, 1.91]`** | **6,191 MiB** | **12,596 MiB** | **2.04× `[1.92, 2.16]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **25.08 s** / **20.97 s** | 49.30 s / 40.42 s | 2.35× `[2.34, 2.37]` | 1.93× `[1.91, 1.94]` | **904 MiB** vs. 4,698 / 6,579 MiB — **5.19× `[4.77, 5.64]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **17.45 s** / **13.98 s** | 35.36 s / 26.25 s | 2.53× `[2.49, 2.57]` | 1.88× `[1.85, 1.91]` | **714 MiB** vs. 3,771 / 3,947 MiB — **5.29× `[5.03, 5.55]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **10.15 s** / **8.10 s** | 23.14 s / 15.62 s | 2.85× `[2.80, 2.91]` | 1.93× `[1.86, 2.00]` | **431 MiB** vs. 3,516 / 3,584 MiB — **8.15× `[7.87, 8.45]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **6.81 s** / **6.33 s** | 17.70 s / 11.37 s | 2.80× `[2.74, 2.86]` | 1.80× `[1.78, 1.82]` | **327 MiB** vs. 1,732 / 3,909 MiB — **5.30× `[5.18, 5.43]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.26 s** / **3.33 s** | 12.52 s / 6.75 s | 3.76× `[3.68, 3.84]` | 2.03× `[2.00, 2.06]` | **210 MiB** vs. 1,476 / 3,503 MiB — **7.04× `[6.94, 7.13]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.49 s** / **2.66 s** | 11.44 s / 5.51 s | 4.30× `[4.21, 4.38]` | 2.07× `[2.04, 2.10]` | **184 MiB** vs. 1,489 / 1,951 MiB — **8.11× `[7.65, 8.59]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.34 s** / **1.09 s** | 6.13 s / 1.97 s | 5.65× `[5.51, 5.79]` | 1.81× `[1.67, 1.96]` | **129 MiB** vs. 661 / 3,263 MiB — **5.11× `[4.29, 6.09]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.21 s** / **1.72 s** | 6.84 s / 2.76 s | 3.98× `[3.79, 4.17]` | 1.61× `[1.52, 1.70]` | **269 MiB** vs. 1,337 / 3,528 MiB — **4.97× `[4.52, 5.47]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.23 s** / **3.77 s** | 13.05 s / 8.60 s | 3.46× `[3.37, 3.55]` | 2.28× `[2.23, 2.33]` | **239 MiB** vs. 1,649 / 3,628 MiB — **6.90× `[6.64, 7.17]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **65.87 s** | — / 87.69 s | — | 1.33× `[1.33, 1.34]` | **71 MiB** vs. — / 2,245 MiB — **31.7×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **0.96 s** / **0.81 s** | 4.50 s / 1.48 s | 5.58× `[5.37, 5.79]` | 1.84× `[1.74, 1.94]` | **124 MiB** vs. 600 / 1,484 MiB — **4.85× `[4.68, 5.03]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.45 s** / **0.37 s** | 2.41 s / 0.50 s | 6.50× `[6.29, 6.72]` | 1.36× `[1.29, 1.44]` | **102 MiB** vs. 513 / 1,420 MiB — **5.04× `[4.79, 5.31]`** |

#### Against `3e2e776c`

Neutral to 1% faster: measured in the same blocks (ratios `3e2e776c / 90bc249a`), every whole-program task is 1.00×–1.01×, and only `bundle-mega` has an interval clearly above 1 (1.01× `[1.01, 1.01]`). The `t=1` times are 0.3–1.8% lower on every task. `googleapis` is 1.01× `[1.00, 1.01]` (109.8 s, 0.4–0.8% faster in every block). `single-files`, which starts the compiler 1,000 times, takes 65.9 s against 67.6 s (2.5% less, CPU time 123.7 s against 125.7 s). Peak RSS is unchanged within its [run-to-run spread](#peak-rss-between-runs); `googleapis` ran at the upper of its two usual levels in this suite (6.2 GiB), so its memory ratio reads 2.04×.

### 3e2e776c on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `3e2e776c`: `5e3cb80a` plus `efa96422` (Refactor Round 66: polymorphic peephole dispatch, an `ASCIICharSet` bitset, stage cleanups and file splits). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `5e3cb80a` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power and kept awake; Java's cold times are within 1.4% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **53.29 s** | 97.18 s | **1.82× `[1.78, 1.87]`** | **5,522 MiB** | 12,049 MiB | **2.18× `[2.04, 2.33]`** |
| **`OPTIMIZATIONS`** | **50.25 s** | 88.97 s | **1.77× `[1.74, 1.80]`** | **5,770 MiB** | 12,850 MiB | **2.23× `[2.01, 2.49]`** |
| **`FINALIZATIONS`** | **6.81 s** | 19.93 s | **2.93× `[2.90, 2.96]`** | **1,108 MiB** | 7,234 MiB | **6.58× `[5.61, 7.73]`** |
| **Total** | **110.3 s** | **206.1 s** | **1.87× `[1.84, 1.89]`** | **5,887 MiB** | **12,850 MiB** | **2.19× `[2.01, 2.39]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **25.23 s** / **21.23 s** | 49.15 s / 40.04 s | 2.32× `[2.30, 2.33]` | 1.89× `[1.88, 1.89]` | **876 MiB** vs. 4,948 / 6,665 MiB — **5.65× `[5.50, 5.80]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **17.55 s** / **14.07 s** | 35.16 s / 26.29 s | 2.50× `[2.48, 2.52]` | 1.87× `[1.84, 1.90]` | **728 MiB** vs. 3,788 / 3,952 MiB — **5.20× `[5.06, 5.35]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **10.20 s** / **8.15 s** | 23.01 s / 15.52 s | 2.82× `[2.78, 2.87]` | 1.91× `[1.89, 1.92]` | **433 MiB** vs. 3,513 / 3,576 MiB — **8.11× `[7.95, 8.27]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **6.84 s** / **6.35 s** | 17.55 s / 11.41 s | 2.76× `[2.70, 2.83]` | 1.80× `[1.78, 1.81]` | **314 MiB** vs. 1,737 / 3,646 MiB — **5.53× `[5.39, 5.67]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.29 s** / **3.36 s** | 12.62 s / 6.69 s | 3.75× `[3.64, 3.87]` | 1.99× `[1.95, 2.03]` | **211 MiB** vs. 1,464 / 3,533 MiB — **6.96× `[6.66, 7.26]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.52 s** / **2.69 s** | 11.55 s / 5.50 s | 4.30× `[4.20, 4.39]` | 2.05× `[2.03, 2.06]` | **185 MiB** vs. 1,627 / 1,953 MiB — **8.77× `[8.12, 9.48]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.34 s** / **1.09 s** | 6.17 s / 1.89 s | 5.65× `[5.55, 5.76]` | 1.73× `[1.68, 1.79]` | **130 MiB** vs. 624 / 1,567 MiB — **4.81× `[4.47, 5.18]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.25 s** / **1.74 s** | 6.93 s / 2.76 s | 3.98× `[3.93, 4.04]` | 1.58× `[1.47, 1.70]` | **272 MiB** vs. 1,370 / 3,513 MiB — **5.04× `[4.67, 5.44]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.26 s** / **3.81 s** | 13.00 s / 8.57 s | 3.41× `[3.36, 3.47]` | 2.25× `[2.22, 2.28]` | **232 MiB** vs. 1,641 / 1,850 MiB — **7.08× `[6.70, 7.47]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **67.55 s** | — / 87.40 s | — | 1.29× `[1.29, 1.30]` | **72 MiB** vs. — / 1,960 MiB — **27.3×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **0.98 s** / **0.82 s** | 4.58 s / 1.48 s | 5.61× `[5.40, 5.82]` | 1.81× `[1.72, 1.91]` | **121 MiB** vs. 606 / 1,487 MiB — **4.99× `[4.88, 5.11]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.45 s** / **0.38 s** | 2.38 s / 0.51 s | 6.23× `[5.97, 6.51]` | 1.33× `[1.24, 1.43]` | **101 MiB** vs. 515 / 1,432 MiB — **5.09× `[4.95, 5.24]`** |

#### Against `5e3cb80a`

Neutral: measured in the same blocks (ratios `5e3cb80a / 3e2e776c`), every task is 1.00× with intervals that include 1, and the `t=1` times are within 1%. `googleapis` is 1.00× `[1.00, 1.00]` (110.3 s), and `single-files` takes 67.6 s against 67.5 s. Peak RSS is unchanged within its [run-to-run spread](#peak-rss-between-runs).

### 5e3cb80a on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `5e3cb80a`: `79777d0a` plus `af4b509e` (flattened `SortingErrorManager` and `NoObjectType` hierarchies, simpler CLI, `jutil` and `TypeInference` code, 7 more file splits). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `79777d0a` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power and kept awake; Java's cold times are within 2.4% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **53.20 s** | 97.06 s | **1.82× `[1.78, 1.87]`** | **5,543 MiB** | 12,059 MiB | **2.17× `[2.07, 2.29]`** |
| **`OPTIMIZATIONS`** | **50.45 s** | 89.07 s | **1.77× `[1.73, 1.80]`** | **5,588 MiB** | 12,703 MiB | **2.27× `[2.19, 2.36]`** |
| **`FINALIZATIONS`** | **6.88 s** | 20.01 s | **2.91× `[2.88, 2.94]`** | **1,181 MiB** | 7,304 MiB | **6.25× `[5.08, 7.69]`** |
| **Total** | **110.5 s** | **206.1 s** | **1.86× `[1.84, 1.89]`** | **5,608 MiB** | **12,703 MiB** | **2.27× `[2.18, 2.35]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **25.32 s** / **21.29 s** | 49.22 s / 40.03 s | 2.31× `[2.30, 2.33]` | 1.88× `[1.86, 1.90]` | **876 MiB** vs. 4,954 / 6,578 MiB — **5.64× `[5.24, 6.08]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **17.61 s** / **14.11 s** | 35.17 s / 26.25 s | 2.49× `[2.46, 2.53]` | 1.86× `[1.84, 1.88]` | **725 MiB** vs. 3,804 / 3,961 MiB — **5.25× `[5.12, 5.38]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **10.25 s** / **8.17 s** | 22.94 s / 15.49 s | 2.81× `[2.75, 2.86]` | 1.89× `[1.84, 1.95]` | **432 MiB** vs. 3,518 / 3,614 MiB — **8.15× `[7.97, 8.34]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **6.86 s** / **6.38 s** | 17.60 s / 11.50 s | 2.76× `[2.70, 2.81]` | 1.80× `[1.79, 1.81]` | **323 MiB** vs. 1,731 / 3,278 MiB — **5.37× `[5.07, 5.68]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.30 s** / **3.37 s** | 12.45 s / 6.61 s | 3.70× `[3.61, 3.80]` | 1.96× `[1.89, 2.04]` | **210 MiB** vs. 1,469 / 3,500 MiB — **7.00× `[6.77, 7.25]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.53 s** / **2.69 s** | 11.57 s / 5.55 s | 4.30× `[4.19, 4.41]` | 2.07× `[2.03, 2.10]` | **186 MiB** vs. 1,623 / 3,458 MiB — **8.71× `[7.90, 9.61]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.35 s** / **1.09 s** | 6.18 s / 1.89 s | 5.66× `[5.46, 5.87]` | 1.73× `[1.68, 1.78]` | **130 MiB** vs. 655 / 1,570 MiB — **5.03× `[4.34, 5.82]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.24 s** / **1.74 s** | 6.84 s / 2.74 s | 3.94× `[3.81, 4.08]` | 1.58× `[1.49, 1.66]` | **275 MiB** vs. 1,356 / 3,498 MiB — **4.92× `[4.60, 5.27]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.26 s** / **3.80 s** | 12.91 s / 8.71 s | 3.39× `[3.36, 3.43]` | 2.29× `[2.25, 2.33]` | **229 MiB** vs. 1,654 / 3,628 MiB — **7.22× `[6.84, 7.63]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **67.47 s** | — / 87.84 s | — | 1.30× `[1.30, 1.31]` | **71 MiB** vs. — / 1,438 MiB — **20.3×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **0.98 s** / **0.82 s** | 4.56 s / 1.51 s | 5.55× `[5.43, 5.68]` | 1.84× `[1.76, 1.92]` | **124 MiB** vs. 604 / 1,482 MiB — **4.87× `[4.74, 5.00]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.45 s** / **0.38 s** | 2.36 s / 0.50 s | 6.21× `[6.10, 6.33]` | 1.32× `[1.24, 1.40]` | **100 MiB** vs. 518 / 1,422 MiB — **5.17× `[5.10, 5.24]`** |

#### Against `79777d0a`

Neutral: measured in the same blocks (ratios `79777d0a / 5e3cb80a`), every task is 1.00×–1.01× with intervals that include or touch 1, and the `t=1` times are within 1%. `googleapis` is 1.00× `[1.00, 1.01]` (110.5 s). `single-files` takes 67.5 s against 68.2 s (1.0% less). Peak RSS is unchanged within its [run-to-run spread](#peak-rss-between-runs); `googleapis` ran at the lower of its two usual levels in this suite (5.6 GiB), so its memory ratio reads 2.27×.

### 79777d0a on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `79777d0a`: `89d80973` plus `96e75119` (fewer `jutil` shims, shared `NodeTraversal`, dataflow and inliner tables, 11 packages simplified). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `89d80973` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power and kept awake; Java's cold times are within 1.0% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **53.30 s** | 96.08 s | **1.80× `[1.72, 1.89]`** | **5,437 MiB** | 12,070 MiB | **2.22× `[2.10, 2.34]`** |
| **`OPTIMIZATIONS`** | **50.43 s** | 89.31 s | **1.77× `[1.74, 1.80]`** | **5,746 MiB** | 12,824 MiB | **2.23× `[2.09, 2.39]`** |
| **`FINALIZATIONS`** | **6.84 s** | 20.03 s | **2.93× `[2.89, 2.97]`** | **1,020 MiB** | 7,295 MiB | **7.16× `[6.88, 7.45]`** |
| **Total** | **110.6 s** | **205.4 s** | **1.86× `[1.82, 1.90]`** | **5,746 MiB** | **12,824 MiB** | **2.23× `[2.09, 2.39]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **25.39 s** / **21.33 s** | 49.35 s / 40.26 s | 2.31× `[2.28, 2.35]` | 1.89× `[1.87, 1.91]` | **879 MiB** vs. 4,957 / 6,588 MiB — **5.64× `[5.35, 5.94]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **17.67 s** / **14.16 s** | 35.05 s / 26.38 s | 2.48× `[2.45, 2.50]` | 1.86× `[1.82, 1.91]` | **718 MiB** vs. 3,785 / 5,024 MiB — **5.27× `[5.12, 5.43]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **10.29 s** / **8.21 s** | 23.02 s / 15.41 s | 2.80× `[2.74, 2.86]` | 1.88× `[1.85, 1.90]` | **422 MiB** vs. 3,513 / 3,660 MiB — **8.33× `[7.99, 8.69]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **6.90 s** / **6.42 s** | 17.75 s / 11.31 s | 2.77× `[2.72, 2.81]` | 1.76× `[1.74, 1.78]` | **323 MiB** vs. 1,741 / 2,604 MiB — **5.40× `[5.25, 5.55]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.32 s** / **3.38 s** | 12.57 s / 6.87 s | 3.72× `[3.66, 3.78]` | 2.03× `[2.01, 2.06]` | **208 MiB** vs. 1,470 / 3,530 MiB — **7.05× `[6.85, 7.26]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.55 s** / **2.70 s** | 11.51 s / 5.53 s | 4.26× `[4.14, 4.38]` | 2.04× `[2.03, 2.06]` | **186 MiB** vs. 1,645 / 2,019 MiB — **8.83× `[8.18, 9.54]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.36 s** / **1.10 s** | 6.14 s / 1.90 s | 5.58× `[5.45, 5.72]` | 1.73× `[1.68, 1.78]` | **131 MiB** vs. 630 / 2,002 MiB — **4.83× `[4.49, 5.19]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.25 s** / **1.75 s** | 6.88 s / 2.79 s | 3.93× `[3.73, 4.13]` | 1.59× `[1.51, 1.67]` | **268 MiB** vs. 1,387 / 3,493 MiB — **5.17× `[4.77, 5.61]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.29 s** / **3.82 s** | 12.83 s / 8.54 s | 3.36× `[3.27, 3.46]` | 2.24× `[2.21, 2.26]` | **241 MiB** vs. 1,662 / 1,910 MiB — **6.89× `[6.60, 7.19]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **68.17 s** | — / 88.26 s | — | 1.29× `[1.28, 1.31]` | **70 MiB** vs. — / 2,406 MiB — **34.6×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **0.98 s** / **0.83 s** | 4.46 s / 1.48 s | 5.36× `[5.18, 5.55]` | 1.78× `[1.68, 1.88]` | **126 MiB** vs. 601 / 1,492 MiB — **4.79× `[4.61, 4.97]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.46 s** / **0.38 s** | 2.38 s / 0.50 s | 6.18× `[5.88, 6.50]` | 1.29× `[1.21, 1.37]` | **103 MiB** vs. 519 / 1,422 MiB — **5.03× `[4.78, 5.29]`** |

#### Against `89d80973`

Equal or faster on every task: measured in the same blocks (ratios `89d80973 / 79777d0a`), `bundle` and `shaka-blockly` are 1.02× `[1.01, 1.03]`, `bundle-advanced`, `bundle-advanced-ts` and `bundle-mega` 1.01× (intervals above 1), and the other whole-program tasks 1.00×–1.01× with intervals touching 1. The `t=1` times fall by 1–2.4% on ten tasks. `googleapis` is 1.02× `[1.01, 1.02]` (110.6 s, 1.3–1.8% faster in every block, each stage 1.0–1.7% faster). `single-files` takes 68.2 s against 69.0 s (1.2% less). Peak RSS is unchanged within its [run-to-run spread](#peak-rss-between-runs).

### 89d80973 on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `89d80973`: `7afbb2fe` plus `72d6be4a` (no more `jutil.StringBuilder` or production `DataFlowAnalysisBase`, a smaller `Property_Key`, 11 fewer `jutil.Try` blocks, every file over 1,520 lines split). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `7afbb2fe` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power and kept awake; Java's cold times are within 2.2% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **54.14 s** | 96.37 s | **1.78× `[1.68, 1.89]`** | **5,331 MiB** | 11,973 MiB | **2.25× `[2.11, 2.39]`** |
| **`OPTIMIZATIONS`** | **51.39 s** | 89.41 s | **1.74× `[1.71, 1.77]`** | **6,324 MiB** | 12,689 MiB | **2.01× `[1.95, 2.07]`** |
| **`FINALIZATIONS`** | **6.91 s** | 19.86 s | **2.88× `[2.85, 2.90]`** | **1,133 MiB** | 7,286 MiB | **6.45× `[5.70, 7.30]`** |
| **Total** | **112.4 s** | **205.7 s** | **1.83× `[1.77, 1.89]`** | **6,324 MiB** | **12,689 MiB** | **2.01× `[1.95, 2.07]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **25.79 s** / **21.61 s** | 49.37 s / 40.29 s | 2.28× `[2.27, 2.30]` | 1.86× `[1.85, 1.88]` | **905 MiB** vs. 4,727 / 6,574 MiB — **5.22× `[5.02, 5.43]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **17.96 s** / **14.34 s** | 35.27 s / 26.41 s | 2.46× `[2.42, 2.50]` | 1.84× `[1.83, 1.86]` | **731 MiB** vs. 3,800 / 4,697 MiB — **5.20× `[4.97, 5.44]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **10.45 s** / **8.32 s** | 23.17 s / 15.46 s | 2.78× `[2.74, 2.83]` | 1.86× `[1.79, 1.92]` | **434 MiB** vs. 3,500 / 3,590 MiB — **8.06× `[7.78, 8.36]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.02 s** / **6.50 s** | 17.68 s / 11.41 s | 2.72× `[2.68, 2.76]` | 1.76× `[1.74, 1.78]` | **326 MiB** vs. 1,748 / 3,818 MiB — **5.36× `[5.25, 5.48]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.37 s** / **3.42 s** | 12.55 s / 6.76 s | 3.67× `[3.60, 3.74]` | 1.98× `[1.93, 2.02]` | **219 MiB** vs. 1,460 / 3,536 MiB — **6.67× `[6.47, 6.89]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.61 s** / **2.74 s** | 11.48 s / 5.49 s | 4.18× `[4.11, 4.25]` | 2.00× `[1.99, 2.02]` | **187 MiB** vs. 1,588 / 1,926 MiB — **8.48× `[7.53, 9.55]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.37 s** / **1.12 s** | 6.18 s / 1.83 s | 5.53× `[5.29, 5.77]` | 1.64× `[1.58, 1.70]` | **142 MiB** vs. 724 / 1,563 MiB — **5.07× `[4.02, 6.40]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.28 s** / **1.77 s** | 6.81 s / 2.77 s | 3.84× `[3.74, 3.95]` | 1.56× `[1.48, 1.64]` | **266 MiB** vs. 1,367 / 3,521 MiB — **5.13× `[4.47, 5.89]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.39 s** / **3.89 s** | 12.95 s / 8.76 s | 3.33× `[3.30, 3.36]` | 2.25× `[2.23, 2.28]` | **231 MiB** vs. 1,996 / 3,665 MiB — **8.27× `[5.69, 12.01]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **69.03 s** | — / 87.81 s | — | 1.27× `[1.26, 1.28]` | **71 MiB** vs. — / 1,113 MiB — **15.6×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.00 s** / **0.84 s** | 4.51 s / 1.48 s | 5.36× `[5.16, 5.56]` | 1.76× `[1.67, 1.86]` | **126 MiB** vs. 606 / 1,489 MiB — **4.82× `[4.69, 4.95]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.46 s** / **0.39 s** | 2.39 s / 0.50 s | 6.17× `[5.88, 6.48]` | 1.30× `[1.20, 1.40]` | **102 MiB** vs. 517 / 1,424 MiB — **5.10× `[5.04, 5.16]`** |

#### Against `7afbb2fe`

Neutral: measured in the same blocks (ratios `7afbb2fe / 89d80973`), every task is 0.99×–1.01×, and only `bundle` has an interval above 1 (1.01× `[1.01, 1.02]`, CPU time 1.02×). The `t=1` times are within 1%. `googleapis` is 1.00× `[1.00, 1.01]` (112.4 s), and `single-files` takes 69.0 s against 69.2 s. Peak RSS is unchanged within its [run-to-run spread](#peak-rss-between-runs); `googleapis` ran at the upper of its two usual levels in this suite (6.3 GiB against 5.7 GiB), so its memory ratio reads 2.01× rather than 2.19×.

### 7afbb2fe on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `7afbb2fe`: `ef4d1df4` plus five refactoring rounds (`c5c0e721`, `72de78b1`, `f6f8dd0f`, `fdab3341` and `7afbb2fe`: value `TypeMismatch` and `AstFactory_Type`, `strpool.ID` in more passes, a concrete `ControlFlowGraph`, a pointer-free `LiveVariablesAnalysis`, monomorphized dataflow worklists, file splits and consolidations). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `ef4d1df4` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power and kept awake; Java's cold times are within 3.6% of the previous section's on every task.

The three `main` builds between `ef4d1df4` and `7afbb2fe` were measured but not published, because they failed the gate: from `c5c0e721` on, `typed-mega`'s output differed from Java's on every run and from run to run (with `--num_parallel_threads=1` as well; about half of the calls to the Cesium externs were dropped as if side-effect free), and `fdab3341` also crashed on `bundle-mega` (`LiveVariablesAnalysis supports at most 128 variables`, a function with 131). `7afbb2fe` fixes both.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **54.41 s** | 95.64 s | **1.76× `[1.71, 1.81]`** | **5,423 MiB** | 11,802 MiB | **2.18× `[2.05, 2.31]`** |
| **`OPTIMIZATIONS`** | **51.21 s** | 90.99 s | **1.78× `[1.71, 1.84]`** | **5,644 MiB** | 12,402 MiB | **2.20× `[2.02, 2.40]`** |
| **`FINALIZATIONS`** | **6.93 s** | 20.04 s | **2.89× `[2.88, 2.91]`** | **1,149 MiB** | 7,313 MiB | **6.39× `[5.63, 7.24]`** |
| **Total** | **112.5 s** | **206.7 s** | **1.84× `[1.83, 1.85]`** | **5,675 MiB** | **12,402 MiB** | **2.19× `[2.01, 2.38]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **25.93 s** / **21.77 s** | 49.32 s / 40.89 s | 2.27× `[2.23, 2.30]` | 1.88× `[1.86, 1.90]` | **901 MiB** vs. 5,043 / 6,611 MiB — **5.58× `[4.92, 6.33]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **18.01 s** / **14.41 s** | 35.36 s / 25.93 s | 2.45× `[2.40, 2.51]` | 1.80× `[1.78, 1.82]` | **715 MiB** vs. 3,803 / 4,074 MiB — **5.32× `[5.12, 5.53]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **10.50 s** / **8.37 s** | 23.06 s / 15.56 s | 2.75× `[2.68, 2.83]` | 1.86× `[1.79, 1.93]` | **428 MiB** vs. 3,514 / 3,609 MiB — **8.21× `[7.83, 8.61]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.05 s** / **6.54 s** | 17.78 s / 11.58 s | 2.72× `[2.66, 2.78]` | 1.77× `[1.75, 1.79]` | **324 MiB** vs. 1,729 / 3,740 MiB — **5.34× `[5.05, 5.65]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.40 s** / **3.44 s** | 12.45 s / 6.60 s | 3.62× `[3.52, 3.72]` | 1.92× `[1.85, 1.98]` | **212 MiB** vs. 1,459 / 3,536 MiB — **6.87× `[6.67, 7.08]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.62 s** / **2.76 s** | 11.51 s / 5.47 s | 4.17× `[4.08, 4.26]` | 1.98× `[1.95, 2.01]` | **190 MiB** vs. 1,552 / 1,966 MiB — **8.16× `[7.27, 9.16]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.38 s** / **1.12 s** | 6.11 s / 1.95 s | 5.48× `[5.29, 5.67]` | 1.74× `[1.72, 1.77]` | **136 MiB** vs. 632 / 1,564 MiB — **4.65× `[4.17, 5.17]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.28 s** / **1.78 s** | 6.86 s / 2.75 s | 3.84× `[3.75, 3.94]` | 1.54× `[1.45, 1.63]` | **276 MiB** vs. 1,392 / 3,498 MiB — **4.99× `[4.18, 5.97]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.44 s** / **3.94 s** | 12.92 s / 8.59 s | 3.28× `[3.21, 3.35]` | 2.18× `[2.15, 2.20]` | **236 MiB** vs. 1,653 / 1,853 MiB — **7.02× `[6.66, 7.40]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **69.21 s** | — / 87.86 s | — | 1.27× `[1.26, 1.27]` | **70 MiB** vs. — / 1,162 MiB — **16.5×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.01 s** / **0.85 s** | 4.48 s / 1.48 s | 5.26× `[5.11, 5.41]` | 1.74× `[1.67, 1.81]` | **128 MiB** vs. 605 / 1,468 MiB — **4.71× `[4.55, 4.88]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.46 s** / **0.39 s** | 2.44 s / 0.51 s | 6.23× `[5.89, 6.59]` | 1.29× `[1.20, 1.39]` | **102 MiB** vs. 518 / 1,428 MiB — **5.07× `[5.01, 5.14]`** |

#### Against `ef4d1df4`

Faster on every task, most on the typed and `ADVANCED` builds: measured in the same blocks (ratios `ef4d1df4 / 7afbb2fe`), `shaka-blockly` is 1.04× `[1.03, 1.04]`, `typed-mega`, `bundle-mega` and `bundle-advanced-ts` 1.03×, `typed-app`, `bundle-advanced`, `bundle` and `typed-full` 1.02× (intervals above 1), `checks-only` 1.01× (CPU time 1.03×), and the two small tasks 1.00×–1.01×. The `t=1` times fall by 2–3.5% on ten tasks. `googleapis` is 1.03× `[1.02, 1.03]` (112.5 s, 2.4–3.1% faster in every block; `CHECKS` and `OPTIMIZATIONS` 2.3–2.4% faster, `FINALIZATIONS` 7.9%). `single-files` takes 69.2 s against 70.1 s (1.3% less). Peak RSS is unchanged within its [run-to-run spread](#peak-rss-between-runs).

### ef4d1df4 on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `ef4d1df4`: `33574cef` plus `ec9611dd` (non-generic `Scope`/`Var`, consolidated args4j, singleton `PassFactory` values, module, peephole and regex cleanups). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `33574cef` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power, kept awake and quiet for the whole suite; Java's cold times are within 2.3% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **55.69 s** | 96.48 s | **1.73× `[1.66, 1.80]`** | **5,598 MiB** | 11,956 MiB | **2.13× `[2.02, 2.25]`** |
| **`OPTIMIZATIONS`** | **52.24 s** | 89.25 s | **1.71× `[1.68, 1.73]`** | **5,631 MiB** | 12,689 MiB | **2.25× `[2.18, 2.33]`** |
| **`FINALIZATIONS`** | **7.42 s** | 20.02 s | **2.70× `[2.68, 2.71]`** | **1,150 MiB** | 7,281 MiB | **6.34× `[6.13, 6.54]`** |
| **Total** | **115.4 s** | **205.7 s** | **1.78× `[1.75, 1.81]`** | **5,703 MiB** | **12,689 MiB** | **2.23× `[2.16, 2.29]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **26.62 s** / **22.29 s** | 49.52 s / 40.25 s | 2.22× `[2.21, 2.24]` | 1.81× `[1.79, 1.82]` | **887 MiB** vs. 4,729 / 6,588 MiB — **5.33× `[4.69, 6.05]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **18.60 s** / **14.83 s** | 35.20 s / 26.32 s | 2.37× `[2.35, 2.40]` | 1.77× `[1.74, 1.81]` | **735 MiB** vs. 3,793 / 4,065 MiB — **5.16× `[4.97, 5.35]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **10.80 s** / **8.59 s** | 23.03 s / 15.37 s | 2.68× `[2.62, 2.74]` | 1.79× `[1.74, 1.83]` | **427 MiB** vs. 3,537 / 3,582 MiB — **8.28× `[8.00, 8.57]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.28 s** / **6.74 s** | 17.70 s / 11.50 s | 2.63× `[2.58, 2.67]` | 1.71× `[1.69, 1.72]` | **321 MiB** vs. 2,050 / 3,695 MiB — **6.16× `[4.27, 8.88]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.55 s** / **3.56 s** | 12.49 s / 6.69 s | 3.51× `[3.42, 3.59]` | 1.88× `[1.82, 1.94]` | **213 MiB** vs. 1,479 / 3,522 MiB — **6.95× `[6.85, 7.05]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.71 s** / **2.83 s** | 11.52 s / 5.63 s | 4.07× `[3.97, 4.17]` | 1.99× `[1.95, 2.03]` | **187 MiB** vs. 1,573 / 3,294 MiB — **8.40× `[7.97, 8.85]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.41 s** / **1.14 s** | 6.11 s / 1.95 s | 5.37× `[5.25, 5.50]` | 1.71× `[1.68, 1.74]` | **141 MiB** vs. 630 / 1,565 MiB — **4.47× `[4.32, 4.62]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.36 s** / **1.82 s** | 6.80 s / 2.77 s | 3.73× `[3.68, 3.78]` | 1.52× `[1.42, 1.62]` | **278 MiB** vs. 1,357 / 3,513 MiB — **4.88× `[4.35, 5.47]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.54 s** / **4.03 s** | 12.89 s / 8.59 s | 3.20× `[3.14, 3.26]` | 2.13× `[2.12, 2.15]` | **231 MiB** vs. 1,648 / 1,868 MiB — **7.12× `[6.96, 7.29]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **70.11 s** | — / 88.75 s | — | 1.27× `[1.26, 1.27]` | **73 MiB** vs. — / 1,663 MiB — **22.8×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.03 s** / **0.86 s** | 4.50 s / 1.50 s | 5.23× `[5.04, 5.42]` | 1.74× `[1.66, 1.83]` | **124 MiB** vs. 603 / 1,489 MiB — **4.87× `[4.70, 5.05]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.46 s** / **0.40 s** | 2.36 s / 0.50 s | 5.89× `[5.74, 6.05]` | 1.24× `[1.17, 1.32]` | **104 MiB** vs. 513 / 1,426 MiB — **4.96× `[4.76, 5.16]`** |

#### Against `33574cef`

Neutral on the compilations themselves, faster to start: measured in the same blocks (ratios `33574cef / ef4d1df4`), every whole-program task is 0.99×–1.00× with intervals that include or touch 1 (`shaka-blockly` 0.99× `[0.99, 1.00]` is the largest difference), the `t=1` times are within 1.5%, and `googleapis` is 1.00× `[1.00, 1.01]` (115.4 s). `single-files`, which starts the compiler 1,000 times, takes 70.1 s against 71.9 s (2.5% less, CPU time 129.1 s against 131.2 s). Peak RSS is unchanged within its [run-to-run spread](#peak-rss-between-runs).

### 33574cef on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `33574cef`: `f60f193b` plus `d29c3f47` (direct `ColorPool` lowering, `strpool.ID` in `RemoveUnusedCode` and Rhino, pass flattening) and `07deaeb6` (compact `JSDocInfo` slices, `strpool.ID` renamers, a peephole prefilter, pass and cache cleanup). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `f60f193b` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power, kept awake and quiet for the whole suite; Java's cold times are within 1.3% of the previous section's on every task but `checks-only` (4.3% slower).

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **55.96 s** | 95.34 s | **1.70× `[1.66, 1.74]`** | **5,529 MiB** | 12,021 MiB | **2.17× `[2.07, 2.28]`** |
| **`OPTIMIZATIONS`** | **51.83 s** | 89.06 s | **1.72× `[1.70, 1.74]`** | **5,760 MiB** | 12,823 MiB | **2.23× `[2.01, 2.48]`** |
| **`FINALIZATIONS`** | **7.40 s** | 20.05 s | **2.71× `[2.68, 2.73]`** | **1,107 MiB** | 7,346 MiB | **6.65× `[6.08, 7.27]`** |
| **Total** | **115.2 s** | **204.4 s** | **1.77× `[1.75, 1.80]`** | **5,841 MiB** | **12,823 MiB** | **2.20× `[2.02, 2.39]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **26.54 s** / **22.27 s** | 49.33 s / 40.16 s | 2.22× `[2.20, 2.24]` | 1.80× `[1.78, 1.82]` | **879 MiB** vs. 4,691 / 6,581 MiB — **5.33× `[4.98, 5.72]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **18.59 s** / **14.83 s** | 35.14 s / 26.29 s | 2.37× `[2.32, 2.42]` | 1.77× `[1.74, 1.81]` | **718 MiB** vs. 3,795 / 4,141 MiB — **5.29× `[5.07, 5.51]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **10.81 s** / **8.59 s** | 23.05 s / 15.61 s | 2.68× `[2.65, 2.72]` | 1.82× `[1.78, 1.86]` | **430 MiB** vs. 3,505 / 3,616 MiB — **8.15× `[7.85, 8.47]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.25 s** / **6.72 s** | 17.76 s / 11.40 s | 2.64× `[2.62, 2.66]` | 1.70× `[1.68, 1.71]` | **318 MiB** vs. 1,734 / 3,524 MiB — **5.46× `[5.15, 5.79]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.54 s** / **3.54 s** | 12.51 s / 6.64 s | 3.53× `[3.45, 3.62]` | 1.87× `[1.81, 1.94]` | **215 MiB** vs. 1,469 / 3,502 MiB — **6.84× `[6.54, 7.16]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.71 s** / **2.81 s** | 11.44 s / 5.60 s | 4.06× `[3.93, 4.20]` | 1.99× `[1.97, 2.01]` | **189 MiB** vs. 1,572 / 1,947 MiB — **8.28× `[7.48, 9.16]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.41 s** / **1.14 s** | 6.12 s / 1.92 s | 5.35× `[5.21, 5.50]` | 1.68× `[1.62, 1.75]` | **142 MiB** vs. 628 / 1,569 MiB — **4.42× `[4.19, 4.67]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.38 s** / **1.83 s** | 6.96 s / 2.80 s | 3.80× `[3.69, 3.92]` | 1.53× `[1.46, 1.59]` | **282 MiB** vs. 1,390 / 3,510 MiB — **4.93× `[4.64, 5.24]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.56 s** / **4.01 s** | 12.92 s / 8.51 s | 3.22× `[3.17, 3.28]` | 2.12× `[2.10, 2.15]` | **231 MiB** vs. 1,660 / 1,798 MiB — **7.19× `[6.61, 7.82]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **71.94 s** | — / 89.56 s | — | 1.24× `[1.24, 1.25]` | **74 MiB** vs. — / 1,069 MiB — **14.5×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.04 s** / **0.86 s** | 4.49 s / 1.51 s | 5.20× `[5.01, 5.40]` | 1.74× `[1.66, 1.83]` | **126 MiB** vs. 604 / 1,485 MiB — **4.78× `[4.74, 4.83]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.47 s** / **0.40 s** | 2.40 s / 0.50 s | 6.03× `[5.82, 6.25]` | 1.26× `[1.17, 1.36]` | **104 MiB** vs. 518 / 1,424 MiB — **4.99× `[4.89, 5.10]`** |

#### Against `f60f193b`

Faster on every task, most on the typed `ADVANCED` builds: measured in the same blocks (ratios `f60f193b / 33574cef`), `typed-full` is 1.05× `[1.04, 1.06]`, `shaka-blockly` 1.04× `[1.04, 1.05]`, `typed-mega` 1.04× `[1.03, 1.05]`, `bundle-mega` 1.03× `[1.03, 1.03]`, and `bundle-advanced`, `checks-only`, `bundle` and `bundle-advanced-ts` 1.02× (intervals above 1); the three smallest tasks are 1.01×–1.02× with intervals touching 1. CPU time falls by the same amounts, and the `t=1` times by 1–4% on ten tasks (`typed-full` 3.9%, `shaka-blockly` and `typed-mega` 3.3%, `bundle-mega` 2.8%). `googleapis` is 1.02× `[1.02, 1.02]` (115.2 s, 1.6–1.9% faster in every block, each stage 1.5–1.9% faster, CPU time 2% lower). `single-files` takes 71.9 s against 73.3 s (1.9% less, CPU time 131.2 s against 135.8 s). Peak RSS is unchanged within its [run-to-run spread](#peak-rss-between-runs).

### f60f193b on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `f60f193b`: `440a0571` plus `e4d93288` (simpler inlining, purity and call passes, `JSType` visitors and `CodeGenerator`). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `440a0571` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power, kept awake and quiet for the whole suite; Java's cold times are within 2.1% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **56.81 s** | 97.02 s | **1.71× `[1.65, 1.77]`** | **5,640 MiB** | 12,201 MiB | **2.16× `[2.06, 2.27]`** |
| **`OPTIMIZATIONS`** | **52.71 s** | 89.18 s | **1.69× `[1.67, 1.71]`** | **5,709 MiB** | 12,731 MiB | **2.24× `[2.04, 2.45]`** |
| **`FINALIZATIONS`** | **7.52 s** | 19.97 s | **2.66× `[2.64, 2.67]`** | **1,107 MiB** | 7,316 MiB | **6.62× `[6.13, 7.15]`** |
| **Total** | **117.0 s** | **206.2 s** | **1.76× `[1.74, 1.78]`** | **5,838 MiB** | **12,738 MiB** | **2.18× `[2.04, 2.34]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **27.29 s** / **23.00 s** | 49.28 s / 40.19 s | 2.14× `[2.12, 2.17]` | 1.75× `[1.74, 1.75]` | **896 MiB** vs. 4,831 / 6,577 MiB — **5.39× `[4.99, 5.81]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **19.22 s** / **15.44 s** | 34.98 s / 26.30 s | 2.26× `[2.24, 2.29]` | 1.70× `[1.69, 1.72]` | **730 MiB** vs. 3,788 / 4,004 MiB — **5.19× `[5.06, 5.32]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **11.25 s** / **9.02 s** | 22.87 s / 15.53 s | 2.54× `[2.50, 2.57]` | 1.72× `[1.69, 1.76]` | **442 MiB** vs. 3,507 / 3,603 MiB — **7.94× `[7.75, 8.13]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.37 s** / **6.85 s** | 17.60 s / 11.48 s | 2.57× `[2.51, 2.63]` | 1.68× `[1.66, 1.69]` | **325 MiB** vs. 1,730 / 3,708 MiB — **5.33× `[5.03, 5.65]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.69 s** / **3.70 s** | 12.42 s / 6.79 s | 3.36× `[3.30, 3.41]` | 1.84× `[1.80, 1.87]` | **215 MiB** vs. 1,459 / 3,537 MiB — **6.78× `[6.45, 7.12]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.79 s** / **2.89 s** | 11.40 s / 5.53 s | 3.95× `[3.88, 4.02]` | 1.91× `[1.87, 1.96]` | **193 MiB** vs. 1,594 / 3,409 MiB — **8.23× `[7.37, 9.19]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.43 s** / **1.17 s** | 6.08 s / 1.87 s | 5.18× `[5.05, 5.31]` | 1.59× `[1.53, 1.65]` | **146 MiB** vs. 630 / 1,566 MiB — **4.31× `[4.16, 4.46]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.42 s** / **1.88 s** | 6.68 s / 2.86 s | 3.56× `[3.47, 3.65]` | 1.52× `[1.47, 1.58]` | **270 MiB** vs. 1,320 / 3,512 MiB — **4.89× `[4.44, 5.39]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.57 s** / **4.07 s** | 12.89 s / 8.60 s | 3.16× `[3.11, 3.22]` | 2.11× `[2.10, 2.12]` | **240 MiB** vs. 1,648 / 1,823 MiB — **6.87× `[6.42, 7.36]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **73.33 s** | — / 87.63 s | — | 1.19× `[1.19, 1.20]` | **75 MiB** vs. — / 2,184 MiB — **29.2×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.05 s** / **0.88 s** | 4.50 s / 1.51 s | 5.12× `[4.99, 5.26]` | 1.72× `[1.62, 1.82]` | **126 MiB** vs. 602 / 1,478 MiB — **4.79× `[4.63, 4.96]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.48 s** / **0.40 s** | 2.37 s / 0.50 s | 5.87× `[5.72, 6.03]` | 1.24× `[1.16, 1.34]` | **104 MiB** vs. 516 / 1,417 MiB — **4.95× `[4.74, 5.17]`** |

#### Against `440a0571`

Neutral to 2% faster: measured in the same blocks (ratios `440a0571 / f60f193b`), `bundle` is 1.02× `[1.01, 1.02]`, `shaka-blockly` and `bundle-advanced-ts` 1.01× `[1.01, 1.02]`, `bundle-advanced`, `typed-full` and `typed-mega` 1.01×, and the other whole-program tasks 0.99×–1.01× with intervals spanning 1. The `t=1` times fall most on the `ADVANCED` and output-heavy tasks (`shaka-blockly` 2.9%, `bundle` 1.6%, `typed-full` 1.4%, `typed-mega` 1.2%). `googleapis` is 1.01× `[1.01, 1.02]` (117.0 s, 0.7–1.6% faster in every block; `OPTIMIZATIONS` 1.2%). `single-files` is unchanged (73.3 s), as is peak RSS within its [run-to-run spread](#peak-rss-between-runs).

### 440a0571 on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `440a0571`: `1d3f88f7` plus `99ad2308` (Round 53: a value `Token`, a scanner without slabs, `IRFactory` fast paths, and denser graph and dataflow code). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `1d3f88f7` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power and kept awake. During the suite's `bundle` and `bundle-advanced-ts`, heavy terminal output (iTerm at up to 76% of a core) slowed Java by 6–8%, so both tasks were run again right after on a quiet machine; their rows come from that rerun. Java's cold times are then within 3% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **57.13 s** | 95.20 s | **1.67× `[1.62, 1.71]`** | **5,616 MiB** | 12,013 MiB | **2.14× `[2.06, 2.22]`** |
| **`OPTIMIZATIONS`** | **53.03 s** | 89.98 s | **1.70× `[1.66, 1.74]`** | **5,831 MiB** | 12,729 MiB | **2.19× `[1.98, 2.41]`** |
| **`FINALIZATIONS`** | **7.53 s** | 19.94 s | **2.65× `[2.62, 2.67]`** | **1,120 MiB** | 7,416 MiB | **6.64× `[5.97, 7.38]`** |
| **Total** | **117.7 s** | **205.1 s** | **1.74× `[1.72, 1.76]`** | **5,886 MiB** | **12,729 MiB** | **2.17× `[1.99, 2.36]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **27.35 s** / **23.07 s** | 49.42 s / 40.39 s | 2.14× `[2.12, 2.17]` | 1.75× `[1.74, 1.76]` | **889 MiB** vs. 5,107 / 6,599 MiB — **5.74× `[5.19, 6.35]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **19.45 s** / **15.64 s** | 35.02 s / 26.63 s | 2.24× `[2.20, 2.27]` | 1.70× `[1.68, 1.72]` | **722 MiB** vs. 3,811 / 4,136 MiB — **5.28× `[5.11, 5.45]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **11.41 s** / **9.18 s** | 23.01 s / 15.49 s | 2.51× `[2.49, 2.53]` | 1.69× `[1.64, 1.74]` | **450 MiB** vs. 3,520 / 3,747 MiB — **7.82× `[7.52, 8.13]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.44 s** / **6.92 s** | 17.76 s / 11.45 s | 2.57× `[2.54, 2.59]` | 1.65× `[1.64, 1.67]` | **333 MiB** vs. 1,740 / 3,790 MiB — **5.23× `[4.92, 5.56]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.83 s** / **3.79 s** | 12.68 s / 7.07 s | 3.34× `[3.23, 3.46]` | 1.87× `[1.81, 1.93]` | **216 MiB** vs. 1,477 / 3,469 MiB — **6.84× `[6.70, 6.99]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.81 s** / **2.91 s** | 11.34 s / 5.54 s | 3.89× `[3.83, 3.95]` | 1.90× `[1.88, 1.92]` | **195 MiB** vs. 1,537 / 1,999 MiB — **7.88× `[7.14, 8.70]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.44 s** / **1.19 s** | 5.97 s / 1.95 s | 5.03× `[4.91, 5.15]` | 1.64× `[1.60, 1.68]` | **140 MiB** vs. 628 / 1,565 MiB — **4.49× `[4.19, 4.80]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.43 s** / **1.87 s** | 6.73 s / 2.80 s | 3.60× `[3.54, 3.67]` | 1.50× `[1.44, 1.56]` | **269 MiB** vs. 1,348 / 3,509 MiB — **5.02× `[4.81, 5.23]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.65 s** / **4.13 s** | 12.93 s / 8.59 s | 3.13× `[3.10, 3.17]` | 2.08× `[2.07, 2.10]` | **242 MiB** vs. 1,664 / 1,814 MiB — **6.87× `[6.58, 7.17]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **73.43 s** | — / 88.77 s | — | 1.21× `[1.20, 1.21]` | **74 MiB** vs. — / 1,086 MiB — **14.6×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.04 s** / **0.88 s** | 4.46 s / 1.47 s | 5.06× `[4.96, 5.17]` | 1.66× `[1.56, 1.77]` | **127 MiB** vs. 600 / 1,485 MiB — **4.74× `[4.54, 4.94]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.48 s** / **0.41 s** | 2.35 s / 0.51 s | 5.79× `[5.65, 5.92]` | 1.24× `[1.15, 1.34]` | **104 MiB** vs. 515 / 1,421 MiB — **4.96× `[4.90, 5.02]`** |

#### Against `1d3f88f7`

Equal or faster on every task, mostly in CPU time: measured in the same blocks (ratios `1d3f88f7 / 440a0571`), wall time is 1.00×–1.02× and CPU time 1.00×–1.03× (`checks-only` CPU 1.03× `[1.03, 1.04]`), and the `t=1` times fall by 1–4% on ten tasks, most on the parse-heavy ones (`downlevel-whitespace` 3.7%, `downlevel-simple` 2.6%, `checks-only` 2.5%). `googleapis` is 1.01× `[1.00, 1.01]` (117.7 s, 0.4–0.7% faster in every block, CPU time 1.4% lower, most of it in `CHECKS`). `single-files` takes 73.4 s against 74.5 s (1.5% less, CPU time 136.0 s against 140.5 s). Peak RSS is unchanged within its [run-to-run spread](#peak-rss-between-runs).

### 1d3f88f7 on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `1d3f88f7`: `a94169fa` plus `1d3f88f7` itself (Round 52 readability and performance rewrites in `jscomp`, `rhino`, `colors`, `graph` and `jutil`). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `a94169fa` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power, kept awake and quiet; Java's cold times are within 1.5% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **57.77 s** | 96.25 s | **1.67× `[1.62, 1.71]`** | **5,700 MiB** | 12,056 MiB | **2.11× `[1.95, 2.30]`** |
| **`OPTIMIZATIONS`** | **53.28 s** | 89.31 s | **1.68× `[1.66, 1.69]`** | **5,791 MiB** | 12,814 MiB | **2.22× `[2.03, 2.43]`** |
| **`FINALIZATIONS`** | **7.58 s** | 20.12 s | **2.66× `[2.63, 2.69]`** | **1,173 MiB** | 7,300 MiB | **6.24× `[5.69, 6.85]`** |
| **Total** | **118.6 s** | **205.7 s** | **1.73× `[1.70, 1.76]`** | **5,951 MiB** | **12,814 MiB** | **2.16× `[2.01, 2.31]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **27.62 s** / **23.24 s** | 48.95 s / 40.84 s | 2.11× `[2.08, 2.14]` | 1.76× `[1.75, 1.77]` | **907 MiB** vs. 4,866 / 6,599 MiB — **5.36× `[4.92, 5.85]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **19.60 s** / **15.72 s** | 35.15 s / 26.39 s | 2.24× `[2.21, 2.27]` | 1.68× `[1.65, 1.71]` | **718 MiB** vs. 3,794 / 4,978 MiB — **5.28× `[5.15, 5.42]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **11.50 s** / **9.18 s** | 23.03 s / 15.45 s | 2.51× `[2.48, 2.54]` | 1.68× `[1.65, 1.72]` | **441 MiB** vs. 3,516 / 3,589 MiB — **7.97× `[7.72, 8.22]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.54 s** / **7.00 s** | 17.64 s / 11.43 s | 2.52× `[2.49, 2.55]` | 1.63× `[1.62, 1.64]` | **329 MiB** vs. 1,743 / 3,787 MiB — **5.30× `[5.02, 5.59]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.80 s** / **3.77 s** | 12.34 s / 6.78 s | 3.27× `[3.21, 3.34]` | 1.80× `[1.77, 1.83]` | **218 MiB** vs. 1,484 / 3,540 MiB — **6.80× `[6.60, 7.01]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.85 s** / **2.93 s** | 11.42 s / 5.58 s | 3.89× `[3.85, 3.93]` | 1.90× `[1.87, 1.93]` | **196 MiB** vs. 1,651 / 3,515 MiB — **8.42× `[8.04, 8.81]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.46 s** / **1.20 s** | 6.04 s / 1.89 s | 5.03× `[4.90, 5.18]` | 1.58× `[1.53, 1.63]` | **145 MiB** vs. 683 / 1,576 MiB — **4.70× `[4.06, 5.45]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.49 s** / **1.89 s** | 6.70 s / 2.77 s | 3.54× `[3.50, 3.58]` | 1.46× `[1.38, 1.54]` | **266 MiB** vs. 1,383 / 3,503 MiB — **5.21× `[4.72, 5.74]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.74 s** / **4.20 s** | 12.88 s / 8.59 s | 3.07× `[2.99, 3.15]` | 2.05× `[2.02, 2.07]` | **243 MiB** vs. 1,632 / 1,798 MiB — **6.72× `[6.17, 7.33]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **74.52 s** | — / 87.33 s | — | 1.17× `[1.17, 1.17]` | **74 MiB** vs. — / 1,726 MiB — **23.3×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.07 s** / **0.89 s** | 4.48 s / 1.50 s | 5.02× `[4.90, 5.14]` | 1.68× `[1.60, 1.76]` | **126 MiB** vs. 605 / 1,485 MiB — **4.80× `[4.62, 4.99]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.50 s** / **0.42 s** | 2.40 s / 0.51 s | 5.73× `[5.55, 5.91]` | 1.21× `[1.12, 1.31]` | **104 MiB** vs. 470 / 1,429 MiB — **4.45× `[3.41, 5.81]`** |

#### Against `a94169fa`

Faster on every task, with CPU time down by about as much: measured in the same blocks (ratios `a94169fa / 1d3f88f7`), `bundle` is 1.04× `[1.03, 1.04]`, `checks-only` and `downlevel-whitespace` 1.03×, `typed-full`, `typed-mega` and `downlevel-simple` 1.02×, and `bundle-advanced`, `shaka-blockly`, `bundle-advanced-ts`, `bundle-mega` and `typed-app` 1.01×, all but `downlevel-simple` with CIs above 1.00. The `t=1` times fall by as much (`bundle` 4.74 s against 4.89 s), so the gain is per operation, not more parallelism. `googleapis` is 1.02× `[1.01, 1.02]` (118.6 s), with each stage faster in every block (`CHECKS` 1.3%, `OPTIMIZATIONS` 1.8%, `FINALIZATIONS` 1.7%). `single-files` takes 74.5 s, against 75.6 s for `a94169fa` in its interleaved check (1.5% less, CPU time 140.5 s against 143.3 s). Peak RSS is unchanged within its [run-to-run spread](#peak-rss-between-runs).

### a94169fa on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `a94169fa`: `a5ea6874` plus `a94169fa` itself (a Stage 2 segment configuration, a parallel chunk render API and conformance state hooks, none of which the default build path uses, and readability and de-reflection cleanups across `cli`, `jscomp`, `diagnostic`, `graph`, `rhino` and `jutil`). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `a5ea6874` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power, kept awake and quiet; Java's cold times are within 3% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **58.37 s** | 96.87 s | **1.66× `[1.56, 1.77]`** | **5,481 MiB** | 11,441 MiB | **2.09× `[1.94, 2.25]`** |
| **`OPTIMIZATIONS`** | **54.16 s** | 90.70 s | **1.67× `[1.64, 1.71]`** | **6,097 MiB** | 12,328 MiB | **2.02× `[1.94, 2.11]`** |
| **`FINALIZATIONS`** | **7.69 s** | 19.89 s | **2.59× `[2.57, 2.60]`** | **1,160 MiB** | 7,327 MiB | **6.33× `[5.74, 6.97]`** |
| **Total** | **120.2 s** | **207.5 s** | **1.73× `[1.68, 1.77]`** | **6,097 MiB** | **12,328 MiB** | **2.02× `[1.94, 2.11]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **27.87 s** / **23.58 s** | 49.68 s / 40.24 s | 2.11× `[2.07, 2.14]` | 1.71× `[1.70, 1.72]` | **898 MiB** vs. 5,114 / 6,601 MiB — **5.69× `[5.33, 6.08]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **19.89 s** / **16.12 s** | 35.45 s / 25.95 s | 2.20× `[2.18, 2.22]` | 1.61× `[1.58, 1.64]` | **725 MiB** vs. 3,794 / 4,603 MiB — **5.24× `[5.11, 5.36]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **11.71 s** / **9.45 s** | 23.11 s / 15.52 s | 2.45× `[2.39, 2.50]` | 1.64× `[1.61, 1.68]` | **441 MiB** vs. 3,520 / 3,587 MiB — **7.99× `[7.73, 8.25]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.61 s** / **7.06 s** | 17.63 s / 11.41 s | 2.50× `[2.48, 2.52]` | 1.62× `[1.60, 1.64]` | **327 MiB** vs. 1,740 / 2,545 MiB — **5.32× `[5.16, 5.49]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.84 s** / **3.81 s** | 12.40 s / 6.69 s | 3.26× `[3.19, 3.32]` | 1.76× `[1.73, 1.78]` | **218 MiB** vs. 1,462 / 3,437 MiB — **6.71× `[6.44, 6.99]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.89 s** / **2.99 s** | 11.39 s / 5.54 s | 3.81× `[3.69, 3.94]` | 1.85× `[1.80, 1.91]` | **196 MiB** vs. 1,668 / 1,977 MiB — **8.51× `[7.89, 9.19]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.49 s** / **1.23 s** | 5.97 s / 1.96 s | 4.87× `[4.78, 4.96]` | 1.59× `[1.52, 1.67]` | **144 MiB** vs. 697 / 2,064 MiB — **4.82× `[4.33, 5.37]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.56 s** / **1.96 s** | 6.64 s / 2.76 s | 3.38× `[3.33, 3.44]` | 1.41× `[1.31, 1.51]` | **274 MiB** vs. 1,333 / 3,525 MiB — **4.86× `[4.61, 5.13]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.89 s** / **4.35 s** | 12.86 s / 8.54 s | 2.96× `[2.92, 2.99]` | 1.96× `[1.94, 1.98]` | **241 MiB** vs. 1,659 / 1,904 MiB — **6.91× `[6.25, 7.64]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **76.64 s** | — / 89.92 s | — | 1.17× `[1.15, 1.20]` | **75 MiB** vs. — / 1,359 MiB — **18.2×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.09 s** / **0.91 s** | 4.50 s / 1.49 s | 4.96× `[4.77, 5.16]` | 1.64× `[1.55, 1.74]` | **125 MiB** vs. 602 / 1,481 MiB — **4.84× `[4.59, 5.09]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.51 s** / **0.43 s** | 2.44 s / 0.50 s | 5.67× `[5.30, 6.06]` | 1.17× `[1.11, 1.24]` | **103 MiB** vs. 488 / 1,431 MiB — **4.71× `[3.63, 6.10]`** |

#### Against `a5ea6874`

Performance-neutral: measured in the same blocks (ratios `a5ea6874 / a94169fa`), wall and CPU time are 0.99×–1.01× on every whole-program task, with every CI including 1.00 except `typed-app` 0.99× `[0.99, 1.00]`, and `googleapis` within 0.3% in every block (1.00× `[1.00, 1.00]`, 120.2 s). `single-files` took 76.6 s in the suite, 2% more than in the previous section, but run against `a5ea6874` in five interleaved blocks right after (no Java), the two builds took 75.6 s and 75.7 s (1.00×, every block within 0.4%), so that difference was the machine. Peak RSS is unchanged within its [run-to-run spread](#peak-rss-between-runs).

### a5ea6874 on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `a5ea6874`: `030244a9` plus `4ffd6ae4` (the remaining monoliths split across all packages, 28 micro-files merged, Java idioms removed) and `e8311d87` (100 micro-files consolidated, `TypeInference` and the Polymer passes on `internal/ast`, the AutoBuilder boilerplate removed). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `030244a9` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power, kept awake and quiet (the run started only with WindowServer below 15% of a core); Java's cold times are within 3.5% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **58.40 s** | 95.32 s | **1.63× `[1.58, 1.69]`** | **5,498 MiB** | 11,986 MiB | **2.18× `[2.02, 2.35]`** |
| **`OPTIMIZATIONS`** | **54.07 s** | 90.04 s | **1.67× `[1.64, 1.69]`** | **6,018 MiB** | 12,726 MiB | **2.12× `[1.91, 2.35]`** |
| **`FINALIZATIONS`** | **7.54 s** | 19.98 s | **2.65× `[2.64, 2.66]`** | **1,208 MiB** | 7,340 MiB | **6.08× `[5.64, 6.56]`** |
| **Total** | **120.0 s** | **205.3 s** | **1.71× `[1.68, 1.75]`** | **6,018 MiB** | **12,729 MiB** | **2.12× `[1.91, 2.35]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **27.88 s** / **23.45 s** | 49.25 s / 40.74 s | 2.10× `[2.09, 2.11]` | 1.74× `[1.72, 1.76]` | **914 MiB** vs. 4,676 / 6,578 MiB — **5.12× `[4.83, 5.42]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **19.91 s** / **16.00 s** | 34.96 s / 26.17 s | 2.19× `[2.17, 2.20]` | 1.64× `[1.61, 1.67]` | **742 MiB** vs. 3,791 / 5,144 MiB — **5.12× `[4.83, 5.42]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **11.73 s** / **9.41 s** | 22.92 s / 15.48 s | 2.44× `[2.38, 2.49]` | 1.65× `[1.62, 1.67]` | **441 MiB** vs. 3,522 / 3,737 MiB — **7.99× `[7.94, 8.05]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.60 s** / **7.06 s** | 17.70 s / 11.26 s | 2.51× `[2.49, 2.52]` | 1.60× `[1.58, 1.61]` | **333 MiB** vs. 1,738 / 2,257 MiB — **5.22× `[5.04, 5.41]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.86 s** / **3.82 s** | 12.45 s / 6.67 s | 3.26× `[3.15, 3.37]` | 1.75× `[1.70, 1.79]` | **217 MiB** vs. 1,459 / 3,519 MiB — **6.73× `[6.53, 6.94]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.89 s** / **2.97 s** | 11.39 s / 5.54 s | 3.84× `[3.77, 3.91]` | 1.87× `[1.86, 1.88]` | **195 MiB** vs. 1,565 / 1,979 MiB — **8.03× `[7.35, 8.77]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.49 s** / **1.22 s** | 6.06 s / 1.93 s | 4.95× `[4.85, 5.06]` | 1.58× `[1.55, 1.61]` | **148 MiB** vs. 632 / 1,548 MiB — **4.26× `[3.93, 4.62]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.57 s** / **1.96 s** | 6.83 s / 2.80 s | 3.49× `[3.39, 3.59]` | 1.43× `[1.38, 1.48]` | **275 MiB** vs. 1,386 / 3,508 MiB — **5.04× `[4.58, 5.56]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.91 s** / **4.36 s** | 12.80 s / 8.65 s | 2.94× `[2.91, 2.97]` | 1.99× `[1.97, 2.00]` | **232 MiB** vs. 1,632 / 1,897 MiB — **7.05× `[6.79, 7.31]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **75.27 s** | — / 88.69 s | — | 1.18× `[1.17, 1.18]` | **75 MiB** vs. — / 2,549 MiB — **34.1×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.10 s** / **0.91 s** | 4.48 s / 1.47 s | 4.94× `[4.76, 5.11]` | 1.62× `[1.54, 1.70]` | **129 MiB** vs. 605 / 1,489 MiB — **4.70× `[4.58, 4.82]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.52 s** / **0.44 s** | 2.47 s / 0.50 s | 5.60× `[5.35, 5.87]` | 1.13× `[1.05, 1.23]` | **105 MiB** vs. 524 / 1,420 MiB — **4.97× `[4.83, 5.10]`** |

#### Against `030244a9`

Neutral to 1% slower: measured in the same blocks (ratios `030244a9 / a5ea6874`), the larger `ADVANCED` builds, typed or not, are about 1% slower at the edge of their CIs (`typed-app`, `shaka-blockly`, `typed-full`, `typed-mega` and `bundle-mega` 0.99×; `bundle-mega`'s CPU time 0.99× `[0.99, 0.99]`), while the other tasks are within 1% either way. `googleapis` is 0–0.4% slower in every block (1.00× `[1.00, 1.00]`, 120.0 s) and `single-files` takes 75.3 s, as before. Peak RSS is unchanged within its [run-to-run spread](#peak-rss-between-runs). As `bundle-mega` does not run type inference, the small cost is not specific to the `TypeInference` port; like the package split two sections below, it would also fit a smaller share of the code being covered by the PGO profile.

### 030244a9 on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `030244a9`: `0ddcfe05` plus `030244a9` itself (12 monoliths in `jscomp` and `rhino` split, Java OOP hierarchies flattened, `internal/ast` adopted further). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `0ddcfe05` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power and kept awake. During the suite's first pass, screen activity (WindowServer at 20–54% of a core, and a browser) slowed Java by up to 14% and `jsgo` by up to 8%, so the eleven whole-program tasks and `googleapis` were run again afterwards on a quiet machine (load below 2, WindowServer below 15%); their rows come from those reruns, and `single-files`'s from the first pass, which ran after the screen activity had stopped. Java's cold times are within 3% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **58.32 s** | 98.75 s | **1.69× `[1.61, 1.78]`** | **5,434 MiB** | 11,892 MiB | **2.19× `[2.00, 2.39]`** |
| **`OPTIMIZATIONS`** | **53.68 s** | 90.34 s | **1.68× `[1.64, 1.72]`** | **5,690 MiB** | 12,571 MiB | **2.21× `[2.08, 2.36]`** |
| **`FINALIZATIONS`** | **7.53 s** | 19.98 s | **2.65× `[2.62, 2.69]`** | **1,145 MiB** | 7,296 MiB | **6.38× `[5.89, 6.91]`** |
| **Total** | **119.5 s** | **209.1 s** | **1.75× `[1.71, 1.79]`** | **5,775 MiB** | **12,571 MiB** | **2.18× `[2.09, 2.27]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **27.72 s** / **23.29 s** | 49.31 s / 40.33 s | 2.12× `[2.10, 2.14]` | 1.73× `[1.72, 1.74]` | **898 MiB** vs. 4,904 / 6,599 MiB — **5.46× `[5.08, 5.87]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **19.85 s** / **15.89 s** | 35.15 s / 26.18 s | 2.21× `[2.18, 2.24]` | 1.65× `[1.62, 1.68]` | **736 MiB** vs. 3,794 / 4,333 MiB — **5.16× `[5.02, 5.30]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **11.66 s** / **9.31 s** | 22.92 s / 15.36 s | 2.46× `[2.44, 2.48]` | 1.65× `[1.62, 1.68]` | **437 MiB** vs. 3,526 / 3,597 MiB — **8.07× `[7.91, 8.22]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.60 s** / **7.03 s** | 17.81 s / 11.38 s | 2.53× `[2.52, 2.55]` | 1.62× `[1.60, 1.63]` | **333 MiB** vs. 1,751 / 2,367 MiB — **5.27× `[5.02, 5.53]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.82 s** / **3.79 s** | 12.44 s / 6.71 s | 3.28× `[3.27, 3.29]` | 1.77× `[1.76, 1.78]` | **219 MiB** vs. 1,465 / 3,454 MiB — **6.70× `[6.59, 6.82]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.87 s** / **2.96 s** | 11.46 s / 5.58 s | 3.88× `[3.83, 3.92]` | 1.89× `[1.88, 1.90]` | **194 MiB** vs. 1,558 / 1,866 MiB — **8.00× `[7.08, 9.03]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.48 s** / **1.22 s** | 6.11 s / 1.89 s | 5.01× `[4.92, 5.10]` | 1.55× `[1.50, 1.60]` | **146 MiB** vs. 630 / 1,520 MiB — **4.32× `[4.02, 4.64]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.56 s** / **1.96 s** | 6.76 s / 2.77 s | 3.46× `[3.36, 3.56]` | 1.41× `[1.35, 1.48]` | **280 MiB** vs. 1,357 / 3,519 MiB — **4.85× `[4.52, 5.19]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.89 s** / **4.34 s** | 13.09 s / 8.59 s | 3.02× `[2.95, 3.08]` | 1.98× `[1.96, 2.00]` | **230 MiB** vs. 1,649 / 1,902 MiB — **7.17× `[7.02, 7.33]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **75.27 s** | — / 87.11 s | — | 1.16× `[1.15, 1.16]` | **73 MiB** vs. — / 2,008 MiB — **27.4×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.08 s** / **0.90 s** | 4.48 s / 1.53 s | 4.97× `[4.87, 5.08]` | 1.70× `[1.63, 1.78]` | **127 MiB** vs. 601 / 1,480 MiB — **4.73× `[4.56, 4.90]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.51 s** / **0.43 s** | 2.39 s / 0.50 s | 5.56× `[5.26, 5.87]` | 1.15× `[1.08, 1.23]` | **106 MiB** vs. 523 / 1,423 MiB — **4.93× `[4.77, 5.10]`** |

#### Against `0ddcfe05`

Performance-neutral: measured in the same blocks (ratios `0ddcfe05 / 030244a9`), wall time is 0.99×–1.01× on every task, with every CI including 1.00, `googleapis` within 0.2% in every block (1.00× `[1.00, 1.00]`, 119.5 s), and `single-files` at 75.3 s, as before. Peak RSS is unchanged within its [run-to-run spread](#peak-rss-between-runs).

### 0ddcfe05 on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `0ddcfe05`: `f99b1259` plus `32e870a9` (the transpilation and peephole passes made idiomatic, three large passes split, `internal/ast` adopted). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `f99b1259` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power and kept awake; Java's cold times are within 8% of the previous section's on every task (2–8% faster on seven of them). During the suite's `single-files`, a browser used up to 70% of a core and both Java's warm run (91.3 s) and `jsgo` (77.1 s) were 3% slower than in the previous section, so `single-files` was run again right after on a quiet machine; its row comes from that rerun.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **57.98 s** | 95.06 s | **1.64× `[1.60, 1.68]`** | **5,442 MiB** | 12,350 MiB | **2.27× `[2.11, 2.44]`** |
| **`OPTIMIZATIONS`** | **54.15 s** | 88.02 s | **1.63× `[1.58, 1.67]`** | **6,386 MiB** | 13,113 MiB | **2.05× `[1.94, 2.17]`** |
| **`FINALIZATIONS`** | **7.57 s** | 19.90 s | **2.63× `[2.58, 2.68]`** | **1,223 MiB** | 7,344 MiB | **6.02× `[5.43, 6.67]`** |
| **Total** | **119.7 s** | **203.0 s** | **1.70× `[1.66, 1.73]`** | **6,386 MiB** | **13,113 MiB** | **2.05× `[1.94, 2.17]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **27.61 s** / **23.23 s** | 49.23 s / 40.21 s | 2.12× `[2.11, 2.13]` | 1.73× `[1.72, 1.74]` | **919 MiB** vs. 4,704 / 6,642 MiB — **5.11× `[4.54, 5.74]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **19.77 s** / **15.85 s** | 35.14 s / 26.28 s | 2.22× `[2.20, 2.23]` | 1.66× `[1.64, 1.68]` | **756 MiB** vs. 3,794 / 4,205 MiB — **5.02× `[4.84, 5.20]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **11.67 s** / **9.36 s** | 23.13 s / 15.47 s | 2.47× `[2.43, 2.51]` | 1.65× `[1.63, 1.68]` | **442 MiB** vs. 3,526 / 3,694 MiB — **7.98× `[7.85, 8.11]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.64 s** / **7.10 s** | 17.71 s / 11.37 s | 2.49× `[2.42, 2.57]` | 1.60× `[1.58, 1.63]` | **343 MiB** vs. 1,751 / 3,861 MiB — **5.11× `[4.88, 5.35]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.83 s** / **3.78 s** | 12.26 s / 6.65 s | 3.24× `[3.20, 3.28]` | 1.76× `[1.71, 1.81]` | **218 MiB** vs. 1,487 / 3,537 MiB — **6.84× `[6.65, 7.03]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.87 s** / **2.95 s** | 11.29 s / 5.57 s | 3.83× `[3.76, 3.90]` | 1.89× `[1.88, 1.91]` | **194 MiB** vs. 1,525 / 1,877 MiB — **7.86× `[7.38, 8.37]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.49 s** / **1.22 s** | 5.92 s / 1.83 s | 4.85× `[4.72, 4.98]` | 1.50× `[1.43, 1.56]` | **146 MiB** vs. 652 / 1,579 MiB — **4.46× `[4.05, 4.93]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.55 s** / **1.95 s** | 6.64 s / 2.76 s | 3.41× `[3.32, 3.49]` | 1.42× `[1.32, 1.52]` | **270 MiB** vs. 1,369 / 3,512 MiB — **5.08× `[4.45, 5.79]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.92 s** / **4.36 s** | 12.89 s / 8.58 s | 2.96× `[2.90, 3.01]` | 1.97× `[1.94, 1.99]` | **231 MiB** vs. 1,655 / 1,837 MiB — **7.18× `[6.87, 7.50]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **75.25 s** | — / 88.74 s | — | 1.18× `[1.18, 1.18]` | **73 MiB** vs. — / 1,700 MiB — **23.4×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.08 s** / **0.90 s** | 4.44 s / 1.50 s | 4.94× `[4.84, 5.05]` | 1.67× `[1.58, 1.76]` | **130 MiB** vs. 605 / 1,476 MiB — **4.64× `[4.52, 4.77]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.51 s** / **0.43 s** | 2.37 s / 0.50 s | 5.50× `[5.45, 5.54]` | 1.16× `[1.09, 1.24]` | **106 MiB** vs. 515 / 1,418 MiB — **4.88× `[4.71, 5.06]`** |

#### Against `f99b1259`

Performance-neutral: measured in the same blocks (ratios `f99b1259 / 0ddcfe05`), wall and CPU time are 0.99×–1.01× on every task, with every CI including 1.00 except `googleapis`, 1.01× `[1.00, 1.01]` (119.7 s), whose gain comes from one block (the others are within 0.4%). `single-files` takes 75.3 s against 75.0 s. Peak RSS is unchanged within its [run-to-run spread](#peak-rss-between-runs).

### f99b1259 on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `f99b1259`: `b7b8fb1c` plus `f856cdb8` (the `cli`, `transpile`, `bundle`, `refactoring`, `linter` and `testjunit` packages extracted from `jscomp`) and `e3c51e21` (core monoliths split, the typed AST adopted in `NodeUtil`). The committed PGO profile (`cmd/jsgo/default.pgo`, from Round 3A) predates the package split, so functions that moved to new packages are no longer matched by it. Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `b7b8fb1c` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power, kept awake and quiet; Java's cold times are within 4.6% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **58.18 s** | 95.75 s | **1.64× `[1.56, 1.74]`** | **5,470 MiB** | 12,057 MiB | **2.20× `[2.05, 2.37]`** |
| **`OPTIMIZATIONS`** | **53.94 s** | 88.97 s | **1.65× `[1.61, 1.69]`** | **6,522 MiB** | 12,913 MiB | **1.98× `[1.91, 2.05]`** |
| **`FINALIZATIONS`** | **7.52 s** | 19.87 s | **2.64× `[2.62, 2.67]`** | **1,180 MiB** | 7,338 MiB | **6.24× `[5.64, 6.90]`** |
| **Total** | **119.6 s** | **204.6 s** | **1.71× `[1.65, 1.77]`** | **6,522 MiB** | **12,913 MiB** | **1.98× `[1.91, 2.05]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **27.86 s** / **23.36 s** | 49.40 s / 40.46 s | 2.11× `[2.08, 2.15]` | 1.73× `[1.72, 1.74]` | **950 MiB** vs. 4,653 / 6,570 MiB — **4.89× `[4.49, 5.32]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **19.80 s** / **15.89 s** | 35.23 s / 26.33 s | 2.22× `[2.19, 2.25]` | 1.66× `[1.63, 1.68]` | **747 MiB** vs. 3,814 / 3,947 MiB — **5.11× `[4.89, 5.33]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **11.67 s** / **9.37 s** | 23.42 s / 15.37 s | 2.50× `[2.44, 2.56]` | 1.64× `[1.60, 1.68]` | **446 MiB** vs. 3,537 / 3,575 MiB — **7.94× `[7.70, 8.18]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.61 s** / **7.10 s** | 18.31 s / 11.51 s | 2.58× `[2.52, 2.64]` | 1.62× `[1.60, 1.64]` | **332 MiB** vs. 1,725 / 3,470 MiB — **5.20× `[4.88, 5.54]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.86 s** / **3.81 s** | 12.78 s / 6.74 s | 3.35× `[3.23, 3.48]` | 1.77× `[1.75, 1.79]` | **220 MiB** vs. 1,464 / 3,533 MiB — **6.67× `[6.51, 6.82]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.89 s** / **2.97 s** | 11.66 s / 5.54 s | 3.93× `[3.86, 3.99]` | 1.86× `[1.85, 1.88]` | **198 MiB** vs. 1,661 / 1,978 MiB — **8.39× `[7.88, 8.93]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.49 s** / **1.23 s** | 6.41 s / 1.97 s | 5.21× `[4.84, 5.60]` | 1.60× `[1.56, 1.64]` | **148 MiB** vs. 669 / 1,557 MiB — **4.49× `[3.90, 5.17]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.56 s** / **1.96 s** | 6.97 s / 2.75 s | 3.56× `[3.44, 3.68]` | 1.40× `[1.32, 1.49]` | **281 MiB** vs. 1,343 / 3,466 MiB — **4.78× `[4.45, 5.14]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.90 s** / **4.37 s** | 13.22 s / 8.69 s | 3.03× `[2.95, 3.11]` | 1.99× `[1.98, 2.00]` | **233 MiB** vs. 1,643 / 1,821 MiB — **7.05× `[6.67, 7.46]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **74.96 s** | — / 86.48 s | — | 1.15× `[1.15, 1.16]` | **74 MiB** vs. — / 1,039 MiB — **14.0×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.09 s** / **0.90 s** | 4.58 s / 1.50 s | 5.10× `[4.92, 5.28]` | 1.67× `[1.58, 1.76]` | **125 MiB** vs. 602 / 1,476 MiB — **4.80× `[4.60, 5.00]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.50 s** / **0.43 s** | 2.37 s / 0.50 s | 5.51× `[5.39, 5.63]` | 1.16× `[1.10, 1.22]` | **107 MiB** vs. 509 / 1,420 MiB — **4.77× `[4.63, 4.90]`** |

#### Against `b7b8fb1c`

Neutral to 1% slower: measured in the same blocks (ratios `b7b8fb1c / f99b1259`), wall and CPU time are 0.99×–1.01× on every task, with `shaka-blockly` and `typed-full` 0.99× `[0.99, 1.00]` at the edge of their CIs, and `googleapis` 0.3% slower in every block (0.996×–0.999×, 119.6 s). `single-files` takes 75.0 s against 75.4 s. Peak RSS is unchanged. A slightly smaller share of the code covered by the PGO profile would fit this pattern; the run does not separate it from the refactor itself.

### b7b8fb1c on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `b7b8fb1c`: `a8dc2699` plus `26cc3661` (Java inheritance hierarchies flattened, sentinel string pointers, `Invalidations`), `aefc27ab` (the side-effect, reference, fusion and inlining foundations on the typed AST) and `a8ed485b` (a larger typed AST, and the call and method optimization passes made idiomatic). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `a8dc2699` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power, kept awake and quiet; Java's cold times are within 3.1% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **58.34 s** | 99.04 s | **1.70× `[1.58, 1.82]`** | **5,467 MiB** | 11,616 MiB | **2.12× `[1.96, 2.30]`** |
| **`OPTIMIZATIONS`** | **54.57 s** | 90.33 s | **1.66× `[1.61, 1.71]`** | **5,927 MiB** | 12,503 MiB | **2.11× `[1.90, 2.35]`** |
| **`FINALIZATIONS`** | **7.52 s** | 20.03 s | **2.66× `[2.65, 2.68]`** | **1,122 MiB** | 7,279 MiB | **6.49× `[6.14, 6.86]`** |
| **Total** | **120.4 s** | **209.4 s** | **1.74× `[1.67, 1.81]`** | **5,927 MiB** | **12,503 MiB** | **2.11× `[1.90, 2.35]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **27.77 s** / **23.36 s** | 49.38 s / 40.38 s | 2.11× `[2.10, 2.13]` | 1.73× `[1.72, 1.74]` | **916 MiB** vs. 4,956 / 6,588 MiB — **5.41× `[5.30, 5.53]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **19.74 s** / **15.84 s** | 35.14 s / 25.97 s | 2.22× `[2.19, 2.24]` | 1.64× `[1.62, 1.66]` | **750 MiB** vs. 3,796 / 3,835 MiB — **5.06× `[4.95, 5.18]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **11.63 s** / **9.30 s** | 23.08 s / 15.54 s | 2.48× `[2.44, 2.52]` | 1.67× `[1.62, 1.72]` | **436 MiB** vs. 3,522 / 3,647 MiB — **8.08× `[7.94, 8.21]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.58 s** / **7.05 s** | 17.64 s / 11.36 s | 2.50× `[2.47, 2.54]` | 1.61× `[1.59, 1.63]` | **336 MiB** vs. 1,735 / 2,271 MiB — **5.17× `[4.92, 5.44]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.85 s** / **3.82 s** | 12.48 s / 6.90 s | 3.27× `[3.19, 3.35]` | 1.80× `[1.74, 1.87]` | **223 MiB** vs. 1,467 / 3,538 MiB — **6.58× `[6.44, 6.72]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.90 s** / **2.96 s** | 11.68 s / 5.58 s | 3.95× `[3.88, 4.02]` | 1.89× `[1.87, 1.90]` | **197 MiB** vs. 1,637 / 1,950 MiB — **8.29× `[7.92, 8.67]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.49 s** / **1.23 s** | 6.13 s / 1.94 s | 4.99× `[4.84, 5.16]` | 1.58× `[1.54, 1.62]` | **151 MiB** vs. 699 / 1,571 MiB — **4.59× `[3.70, 5.71]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.57 s** / **1.97 s** | 6.91 s / 2.80 s | 3.50× `[3.34, 3.68]` | 1.42× `[1.35, 1.49]` | **277 MiB** vs. 1,384 / 3,551 MiB — **5.01× `[4.76, 5.26]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.90 s** / **4.37 s** | 13.09 s / 8.78 s | 2.99× `[2.91, 3.08]` | 2.01× `[1.99, 2.03]` | **238 MiB** vs. 1,644 / 1,792 MiB — **6.91× `[6.39, 7.47]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **75.39 s** | — / 88.92 s | — | 1.18× `[1.17, 1.19]` | **73 MiB** vs. — / 2,296 MiB — **31.3×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.08 s** / **0.91 s** | 4.51 s / 1.54 s | 4.97× `[4.77, 5.18]` | 1.70× `[1.60, 1.81]` | **127 MiB** vs. 604 / 1,494 MiB — **4.74× `[4.54, 4.95]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.51 s** / **0.43 s** | 2.40 s / 0.51 s | 5.54× `[5.28, 5.82]` | 1.18× `[1.11, 1.25]` | **105 MiB** vs. 517 / 1,439 MiB — **4.93× `[4.77, 5.10]`** |

#### Against `a8dc2699`

Performance-neutral: measured in the same blocks (ratios `a8dc2699 / b7b8fb1c`), wall and CPU time are 0.99×–1.02× on every task, with every CI including 1.00; `googleapis` is 1.00× `[0.99, 1.02]` (120.4 s) and `single-files` takes 75.4 s against 75.7 s. Peak RSS is unchanged.

### a8dc2699 on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `a8dc2699`: `611eda44` plus `a49e0bab` (an `internal/graph` package, `jutil.ReuseSlab` and the `ast_isomorphism` tool) and `f5681cff` (typed AST helpers, and seven CFG and dataflow passes ported to the typed AST). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `611eda44` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power and kept awake; Java's cold times are within 2.4% of the previous section's on every task. The terminal (iTerm) used 35–40% of a core for part of the run; in block 4 of `checks-only` the head ran 11% slow (2.17 s against 1.94–1.98 s).

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **58.38 s** | 98.10 s | **1.68× `[1.61, 1.75]`** | **5,375 MiB** | 11,725 MiB | **2.18× `[2.08, 2.28]`** |
| **`OPTIMIZATIONS`** | **54.35 s** | 88.83 s | **1.63× `[1.63, 1.64]`** | **5,902 MiB** | 12,448 MiB | **2.11× `[1.94, 2.30]`** |
| **`FINALIZATIONS`** | **7.56 s** | 19.98 s | **2.64× `[2.63, 2.66]`** | **1,213 MiB** | 7,268 MiB | **6.00× `[5.59, 6.44]`** |
| **Total** | **120.3 s** | **206.9 s** | **1.72× `[1.69, 1.75]`** | **5,902 MiB** | **12,476 MiB** | **2.12× `[1.95, 2.30]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **27.78 s** / **23.40 s** | 49.74 s / 40.28 s | 2.13× `[2.09, 2.16]` | 1.72× `[1.71, 1.74]` | **919 MiB** vs. 4,911 / 6,582 MiB — **5.34× `[5.07, 5.62]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **19.99 s** / **16.01 s** | 35.57 s / 26.44 s | 2.22× `[2.19, 2.25]` | 1.65× `[1.62, 1.68]` | **746 MiB** vs. 3,795 / 3,974 MiB — **5.09× `[4.94, 5.25]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **11.74 s** / **9.43 s** | 23.24 s / 15.54 s | 2.46× `[2.43, 2.50]` | 1.65× `[1.62, 1.67]` | **442 MiB** vs. 3,523 / 3,583 MiB — **7.97× `[7.71, 8.25]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.64 s** / **7.11 s** | 17.87 s / 11.48 s | 2.52× `[2.44, 2.59]` | 1.62× `[1.58, 1.65]` | **325 MiB** vs. 1,740 / 2,360 MiB — **5.36× `[5.02, 5.72]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.85 s** / **3.80 s** | 12.45 s / 6.80 s | 3.27× `[3.24, 3.30]` | 1.79× `[1.76, 1.82]` | **219 MiB** vs. 1,470 / 3,536 MiB — **6.71× `[6.43, 7.00]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.90 s** / **2.98 s** | 11.33 s / 5.51 s | 3.80× `[3.68, 3.92]` | 1.85× `[1.82, 1.87]` | **196 MiB** vs. 1,624 / 1,931 MiB — **8.25× `[7.54, 9.03]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.50 s** / **1.22 s** | 6.24 s / 1.95 s | 5.11× `[5.01, 5.21]` | 1.60× `[1.55, 1.64]` | **150 MiB** vs. 669 / 1,577 MiB — **4.45× `[3.50, 5.66]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.56 s** / **2.01 s** | 6.83 s / 2.85 s | 3.41× `[3.21, 3.61]` | 1.42× `[1.33, 1.52]` | **286 MiB** vs. 1,411 / 3,494 MiB — **4.93× `[4.66, 5.21]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.90 s** / **4.36 s** | 13.09 s / 8.61 s | 3.00× `[2.94, 3.06]` | 1.97× `[1.95, 2.00]` | **235 MiB** vs. 1,662 / 3,131 MiB — **7.07× `[6.77, 7.38]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **75.70 s** | — / 88.85 s | — | 1.17× `[1.17, 1.18]` | **75 MiB** vs. — / 1,081 MiB — **14.4×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.08 s** / **0.91 s** | 4.52 s / 1.52 s | 4.97× `[4.89, 5.05]` | 1.68× `[1.59, 1.77]` | **128 MiB** vs. 601 / 1,477 MiB — **4.71× `[4.60, 4.83]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.51 s** / **0.43 s** | 2.38 s / 0.51 s | 5.51× `[5.34, 5.68]` | 1.17× `[1.11, 1.23]` | **105 MiB** vs. 516 / 1,426 MiB — **4.92× `[4.86, 4.98]`** |

#### Against `611eda44`

Neutral to 2% slower: measured in the same blocks (ratios `611eda44 / a8dc2699`), wall and CPU time are 0.98×–1.00× on every task, with `downlevel-simple` 0.98× `[0.97, 1.00]` and `bundle-advanced`, `bundle` and `typed-full` 0.99× at the edge of their CIs; `googleapis` is 1.00× `[0.99, 1.01]` (120.3 s) and `single-files` takes 75.7 s against 75.3 s. Peak RSS is unchanged.

#### Peak RSS between runs

The same binary's peak RSS can move by 5–10% from one run to the next, with every block of a run on the same side, which the per-run CIs (computed across blocks) do not capture. `611eda44`'s peak RSS, for example, was 132–143 MiB on `typed-app` in its own run and 130–156 MiB in this one, 332–351 MiB against 305–357 MiB on `bundle-advanced-ts`, and 5.6–6.5 GiB against 6.4–6.7 GiB on `googleapis` (where this run's 1.11× `[1.03, 1.19]` for `a8dc2699` against `611eda44` is therefore not a saving). Differences in peak RSS below about 10% between two builds are only reported below when another run confirms them; the previous section's two are corrected accordingly.

### 611eda44 on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `611eda44`: `9b7b72ed` plus `1ad51cb5`, which replaces `JSDocInfo`'s `ArrayList`s with Go slices, reflects the log keys and drops the source-map scratch unpacking. Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `9b7b72ed` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power, kept awake and quiet throughout; Java's cold times are within 4% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **58.21 s** | 96.43 s | **1.66× `[1.60, 1.72]`** | **5,641 MiB** | 11,824 MiB | **2.09× `[1.98, 2.21]`** |
| **`OPTIMIZATIONS`** | **54.34 s** | 89.64 s | **1.65× `[1.64, 1.66]`** | **5,714 MiB** | 12,641 MiB | **2.22× `[2.00, 2.45]`** |
| **`FINALIZATIONS`** | **7.53 s** | 20.00 s | **2.65× `[2.63, 2.68]`** | **1,136 MiB** | 7,356 MiB | **6.48× `[6.00, 7.01]`** |
| **Total** | **120.1 s** | **206.1 s** | **1.72× `[1.68, 1.75]`** | **5,842 MiB** | **12,641 MiB** | **2.17× `[2.00, 2.34]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **27.66 s** / **23.25 s** | 49.17 s / 40.79 s | 2.11× `[2.10, 2.13]` | 1.75× `[1.74, 1.77]` | **915 MiB** vs. 4,773 / 6,579 MiB — **5.22× `[4.85, 5.61]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **19.81 s** / **15.90 s** | 35.29 s / 26.53 s | 2.22× `[2.18, 2.26]` | 1.67× `[1.64, 1.70]` | **740 MiB** vs. 3,784 / 4,216 MiB — **5.11× `[5.03, 5.20]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **11.67 s** / **9.36 s** | 23.34 s / 15.40 s | 2.49× `[2.44, 2.55]` | 1.65× `[1.61, 1.68]` | **441 MiB** vs. 3,521 / 3,657 MiB — **7.99× `[7.85, 8.12]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.65 s** / **7.05 s** | 17.71 s / 11.67 s | 2.51× `[2.46, 2.56]` | 1.65× `[1.63, 1.68]` | **342 MiB** vs. 2,056 / 3,855 MiB — **5.79× `[4.07, 8.22]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.84 s** / **3.81 s** | 12.49 s / 6.86 s | 3.27× `[3.16, 3.39]` | 1.80× `[1.77, 1.83]` | **224 MiB** vs. 1,480 / 3,465 MiB — **6.62× `[6.26, 7.00]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.88 s** / **2.95 s** | 11.47 s / 5.56 s | 3.88× `[3.79, 3.97]` | 1.88× `[1.84, 1.93]` | **197 MiB** vs. 1,649 / 3,496 MiB — **8.37× `[7.67, 9.14]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.49 s** / **1.22 s** | 6.17 s / 1.94 s | 5.05× `[4.87, 5.25]` | 1.59× `[1.53, 1.65]` | **136 MiB** vs. 668 / 1,554 MiB — **4.90× `[4.09, 5.86]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.57 s** / **1.96 s** | 6.92 s / 2.78 s | 3.53× `[3.48, 3.59]` | 1.42× `[1.34, 1.50]` | **267 MiB** vs. 1,389 / 3,505 MiB — **5.22× `[4.92, 5.53]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.89 s** / **4.34 s** | 13.00 s / 8.59 s | 2.99× `[2.90, 3.09]` | 1.98× `[1.96, 2.00]` | **238 MiB** vs. 1,661 / 1,803 MiB — **6.99× `[6.79, 7.20]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **75.33 s** | — / 88.52 s | — | 1.18× `[1.17, 1.18]` | **74 MiB** vs. — / 1,321 MiB — **17.8×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.08 s** / **0.90 s** | 4.44 s / 1.49 s | 4.91× `[4.78, 5.05]` | 1.65× `[1.57, 1.73]` | **127 MiB** vs. 607 / 1,484 MiB — **4.79× `[4.55, 5.05]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.51 s** / **0.43 s** | 2.32 s / 0.51 s | 5.38× `[5.20, 5.57]` | 1.19× `[1.11, 1.27]` | **106 MiB** vs. 520 / 1,429 MiB — **4.90× `[4.68, 5.12]`** |

#### Against `9b7b72ed`

Measured in the same blocks (ratios `9b7b72ed / 611eda44`), the whole-program `ADVANCED` builds are 2–3% faster, with CPU time down by about as much: `bundle-advanced` 1.03× `[1.01, 1.04]`, `typed-full`, `typed-mega` and `bundle-mega` 1.02× (CIs from 1.01), `shaka-blockly` and `bundle-advanced-ts` 1.02× (CIs from 1.00). The `SIMPLE`, `WHITESPACE_ONLY` and `CHECKS` builds and `typed-app` are unchanged (0.99×–1.00×), as are `googleapis` (1.00× `[1.00, 1.01]`, 120.1 s; faster in every block by 0–0.7%) and `single-files` (75.3 s against 75.4 s). Peak RSS is unchanged. (This paragraph first reported `typed-app` 10% lower and `bundle-advanced-ts` 5% higher; the next run measured `611eda44` itself at 130–156 MiB and 305–357 MiB on those tasks, so both were [run-to-run noise](#peak-rss-between-runs).)

### 9b7b72ed on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `9b7b72ed`: `80cb7956` plus fast mode's first round (`JSGO_PASS_LOG` per-pass profiling, the relaxed fast-mode gate, and in fast mode the optimization loop's AST sizes counted and the source-map mappings encoded and resolved on the workers). The suite runs the default mode, which these commits leave as it was. Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `80cb7956` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power and kept awake; Java's cold times are within 1.6% of the previous section's on every task, but the terminal (iTerm, 35–45% of a core) and another agent (up to 26%) were busy at times. In blocks 0 and 1 of `downlevel-whitespace` the previous build ran 15–28% slow while the head did not, so that task was run again right after the suite; its row comes from that rerun.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **58.68 s** | 97.60 s | **1.66× `[1.58, 1.75]`** | **5,470 MiB** | 11,539 MiB | **2.11× `[1.99, 2.23]`** |
| **`OPTIMIZATIONS`** | **54.65 s** | 91.70 s | **1.68× `[1.60, 1.76]`** | **6,278 MiB** | 12,081 MiB | **1.93× `[1.64, 2.26]`** |
| **`FINALIZATIONS`** | **7.59 s** | 20.26 s | **2.67× `[2.62, 2.72]`** | **1,135 MiB** | 7,266 MiB | **6.42× `[5.79, 7.12]`** |
| **Total** | **120.9 s** | **209.6 s** | **1.73× `[1.68, 1.79]`** | **6,294 MiB** | **12,391 MiB** | **1.97× `[1.76, 2.21]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **28.07 s** / **23.67 s** | 49.43 s / 40.40 s | 2.09× `[2.08, 2.10]` | 1.71× `[1.70, 1.71]` | **931 MiB** vs. 4,723 / 6,590 MiB — **5.07× `[4.62, 5.56]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **20.07 s** / **16.20 s** | 35.22 s / 26.26 s | 2.17× `[2.14, 2.21]` | 1.62× `[1.58, 1.66]` | **729 MiB** vs. 3,804 / 4,195 MiB — **5.22× `[5.12, 5.32]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **11.85 s** / **9.51 s** | 23.00 s / 15.76 s | 2.42× `[2.38, 2.46]` | 1.66× `[1.62, 1.69]` | **444 MiB** vs. 3,518 / 3,589 MiB — **7.93× `[7.76, 8.10]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.77 s** / **7.22 s** | 18.02 s / 11.75 s | 2.50× `[2.46, 2.53]` | 1.63× `[1.60, 1.65]` | **342 MiB** vs. 2,092 / 3,600 MiB — **5.86× `[3.93, 8.73]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.94 s** / **3.87 s** | 12.42 s / 6.82 s | 3.21× `[3.18, 3.24]` | 1.76× `[1.71, 1.81]` | **219 MiB** vs. 1,460 / 3,489 MiB — **6.66× `[6.62, 6.69]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **4.00 s** / **3.04 s** | 11.42 s / 5.57 s | 3.75× `[3.59, 3.92]` | 1.83× `[1.76, 1.90]` | **195 MiB** vs. 1,706 / 1,892 MiB — **8.76× `[8.30, 9.25]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.51 s** / **1.23 s** | 6.07 s / 1.82 s | 4.92× `[4.76, 5.09]` | 1.48× `[1.42, 1.54]` | **145 MiB** vs. 657 / 1,521 MiB — **4.51× `[4.16, 4.90]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.59 s** / **1.96 s** | 6.83 s / 2.87 s | 3.49× `[3.29, 3.70]` | 1.47× `[1.40, 1.53]` | **269 MiB** vs. 1,382 / 3,518 MiB — **5.14× `[4.90, 5.38]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.89 s** / **4.34 s** | 13.06 s / 8.77 s | 3.01× `[2.97, 3.04]` | 2.02× `[1.99, 2.05]` | **229 MiB** vs. 1,675 / 1,794 MiB — **7.32× `[7.04, 7.61]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **75.41 s** | — / 88.12 s | — | 1.17× `[1.16, 1.17]` | **76 MiB** vs. — / 2,231 MiB — **29.5×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.12 s** / **0.91 s** | 4.54 s / 1.52 s | 4.97× `[4.73, 5.22]` | 1.66× `[1.56, 1.76]` | **129 MiB** vs. 600 / 1,480 MiB — **4.65× `[4.53, 4.78]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.50 s** / **0.43 s** | 2.41 s / 0.50 s | 5.62× `[5.52, 5.71]` | 1.16× `[1.10, 1.23]` | **107 MiB** vs. 519 / 1,427 MiB — **4.84× `[4.69, 4.99]`** |

#### Against `80cb7956`

Performance-neutral: measured in the same blocks (ratios `80cb7956 / 9b7b72ed`), wall and CPU time are 0.99×–1.00× on every task (`bundle-advanced-ts` and `typed-mega` 0.99× `[0.98, 1.00]`, at the edge of their CIs), `googleapis` is 1.00× `[0.99, 1.01]` (120.9 s, with the blocks between 0.985× and 1.004×), and `single-files` takes 75.4 s against 75.3 s. Peak RSS is unchanged within noise on every task. `googleapis`'s memory ratio against Java (1.97×) is within its block-to-block spread: `jsgo`'s peak RSS ranged over 5.4–6.4 GiB in this run and Java's was lower than usual (12,391 MiB).

### 80cb7956 on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `80cb7956`: `29f91ce6` plus the typed AST view's third round (`ast.RoleOf`, role aliases and `Must<X>()` accessors; `RemoveUnusedCode`, `InlineFunctions`, `FlowSensitiveInlineVariables`, `Es6NormalizeClasses`, `ClosureRewriteModule` and `ProcessCommonJSModules` reading the AST through it), the last non-test `jutil.ArrayList` uses removed, and the `maprangecheck` guard on Go map iteration. Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `29f91ce6` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power, kept awake and quiet throughout; Java's cold times are within 2.1% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **58.04 s** | 97.52 s | **1.68× `[1.62, 1.74]`** | **5,483 MiB** | 12,314 MiB | **2.25× `[2.15, 2.34]`** |
| **`OPTIMIZATIONS`** | **53.82 s** | 89.14 s | **1.66× `[1.64, 1.67]`** | **5,587 MiB** | 13,094 MiB | **2.35× `[2.19, 2.51]`** |
| **`FINALIZATIONS`** | **7.54 s** | 19.96 s | **2.65× `[2.63, 2.66]`** | **1,169 MiB** | 7,343 MiB | **6.29× `[5.88, 6.73]`** |
| **Total** | **119.4 s** | **206.6 s** | **1.73× `[1.70, 1.76]`** | **5,649 MiB** | **13,094 MiB** | **2.32× `[2.19, 2.45]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **27.97 s** / **23.59 s** | 48.96 s / 40.18 s | 2.08× `[2.06, 2.09]` | 1.70× `[1.70, 1.71]` | **926 MiB** vs. 4,768 / 6,591 MiB — **5.14× `[4.65, 5.68]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **20.01 s** / **16.11 s** | 35.10 s / 26.49 s | 2.18× `[2.17, 2.19]` | 1.64× `[1.61, 1.68]` | **741 MiB** vs. 3,806 / 4,142 MiB — **5.13× `[4.95, 5.32]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **11.77 s** / **9.45 s** | 22.92 s / 15.41 s | 2.42× `[2.40, 2.45]` | 1.63× `[1.60, 1.66]` | **441 MiB** vs. 3,502 / 3,621 MiB — **7.95× `[7.76, 8.14]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.68 s** / **7.12 s** | 17.75 s / 11.46 s | 2.49× `[2.44, 2.55]` | 1.61× `[1.59, 1.63]` | **335 MiB** vs. 1,746 / 3,772 MiB — **5.21× `[4.96, 5.47]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.88 s** / **3.84 s** | 12.55 s / 6.75 s | 3.27× `[3.19, 3.35]` | 1.76× `[1.74, 1.77]` | **220 MiB** vs. 1,471 / 3,473 MiB — **6.69× `[6.43, 6.96]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.91 s** / **2.99 s** | 11.43 s / 5.52 s | 3.83× `[3.76, 3.90]` | 1.85× `[1.82, 1.87]` | **196 MiB** vs. 1,615 / 3,506 MiB — **8.23× `[7.46, 9.09]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.49 s** / **1.23 s** | 6.13 s / 1.91 s | 5.00× `[4.83, 5.17]` | 1.56× `[1.49, 1.63]` | **141 MiB** vs. 629 / 1,561 MiB — **4.46× `[4.13, 4.83]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.55 s** / **1.95 s** | 6.88 s / 2.83 s | 3.53× `[3.39, 3.67]` | 1.45× `[1.39, 1.51]` | **270 MiB** vs. 1,349 / 3,519 MiB — **5.01× `[4.64, 5.40]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.87 s** / **4.32 s** | 12.88 s / 8.55 s | 2.98× `[2.93, 3.03]` | 1.98× `[1.96, 1.99]` | **237 MiB** vs. 1,637 / 1,880 MiB — **6.90× `[6.63, 7.18]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **75.31 s** | — / 87.68 s | — | 1.16× `[1.16, 1.17]` | **75 MiB** vs. — / 1,191 MiB — **15.9×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.08 s** / **0.90 s** | 4.51 s / 1.49 s | 5.03× `[4.90, 5.17]` | 1.66× `[1.58, 1.75]` | **129 MiB** vs. 603 / 1,488 MiB — **4.69× `[4.57, 4.81]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.50 s** / **0.43 s** | 2.41 s / 0.50 s | 5.55× `[5.31, 5.80]` | 1.16× `[1.08, 1.24]` | **107 MiB** vs. 511 / 1,428 MiB — **4.76× `[4.62, 4.90]`** |

#### Against `29f91ce6`

Performance-neutral: measured in the same blocks (ratios `29f91ce6 / 80cb7956`), wall time is 0.98×–1.00× and CPU time 0.98×–1.01× on every task, with `downlevel-simple` and `checks-only` 0.99× `[0.99, 1.00]` at the edge of their CIs; `googleapis` is 1.00× `[1.00, 1.01]` (119.4 s) and `single-files` takes 75.3 s against 75.2 s. Peak RSS is unchanged within noise on every task. `googleapis`'s peak RSS was 5,533–6,060 MiB for `80cb7956` against 5,602–6,634 MiB for `29f91ce6` in the same blocks, which puts its memory ratio against Java back at 2.32× (2.09× in the previous section); the difference between the two builds (1.07× `[0.98, 1.16]`) is within that block-to-block spread.

### 29f91ce6 on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `29f91ce6`: `8238e94e` plus three performance commits: `6ae7148b` (a flatter `GlobalNamespace`, without property-renaming graph roots), `3c4af5c1` (Stage 3 chunk source maps and file writes in parallel) and `6b9683cf` (inline source maps decoded lazily, `SaveState` in parallel). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `8238e94e` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power, kept awake and quiet throughout; Java's cold times are within 1.3% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **57.99 s** | 96.92 s | **1.67× `[1.62, 1.73]`** | **5,532 MiB** | 11,973 MiB | **2.16× `[1.97, 2.37]`** |
| **`OPTIMIZATIONS`** | **54.02 s** | 88.47 s | **1.64× `[1.64, 1.64]`** | **6,233 MiB** | 13,020 MiB | **2.09× `[1.93, 2.27]`** |
| **`FINALIZATIONS`** | **7.54 s** | 19.91 s | **2.64× `[2.63, 2.65]`** | **1,200 MiB** | 7,349 MiB | **6.13× `[5.78, 6.50]`** |
| **Total** | **119.6 s** | **205.3 s** | **1.72× `[1.69, 1.74]`** | **6,233 MiB** | **13,020 MiB** | **2.09× `[1.93, 2.27]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **28.00 s** / **23.61 s** | 49.30 s / 40.18 s | 2.09× `[2.07, 2.11]` | 1.70× `[1.70, 1.71]` | **927 MiB** vs. 4,850 / 6,580 MiB — **5.21× `[4.48, 6.07]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **19.96 s** / **16.06 s** | 35.20 s / 26.45 s | 2.19× `[2.18, 2.21]` | 1.65× `[1.62, 1.67]` | **755 MiB** vs. 3,803 / 3,909 MiB — **5.04× `[4.95, 5.14]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **11.77 s** / **9.45 s** | 23.06 s / 15.57 s | 2.44× `[2.39, 2.49]` | 1.65× `[1.60, 1.70]` | **440 MiB** vs. 3,501 / 3,579 MiB — **7.96× `[7.65, 8.29]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.66 s** / **7.11 s** | 17.74 s / 11.35 s | 2.49× `[2.47, 2.52]` | 1.60× `[1.58, 1.61]` | **332 MiB** vs. 1,737 / 3,739 MiB — **5.23× `[4.92, 5.56]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.87 s** / **3.85 s** | 12.49 s / 6.77 s | 3.24× `[3.15, 3.33]` | 1.76× `[1.74, 1.78]` | **213 MiB** vs. 1,460 / 3,451 MiB — **6.84× `[6.61, 7.09]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.90 s** / **3.00 s** | 11.46 s / 5.49 s | 3.82× `[3.76, 3.89]` | 1.83× `[1.81, 1.85]` | **195 MiB** vs. 1,635 / 1,984 MiB — **8.38× `[7.83, 8.98]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.49 s** / **1.23 s** | 6.06 s / 1.83 s | 4.95× `[4.93, 4.96]` | 1.49× `[1.43, 1.56]` | **144 MiB** vs. 627 / 1,574 MiB — **4.36× `[4.04, 4.71]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.55 s** / **1.95 s** | 6.74 s / 2.77 s | 3.46× `[3.33, 3.60]` | 1.42× `[1.35, 1.50]` | **262 MiB** vs. 1,351 / 3,502 MiB — **5.15× `[4.78, 5.54]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.88 s** / **4.32 s** | 12.94 s / 8.57 s | 2.99× `[2.91, 3.08]` | 1.98× `[1.97, 2.00]` | **237 MiB** vs. 1,649 / 1,795 MiB — **6.97× `[6.54, 7.41]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **75.22 s** | — / 87.47 s | — | 1.16× `[1.16, 1.17]` | **73 MiB** vs. — / 1,092 MiB — **14.9×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.08 s** / **0.90 s** | 4.44 s / 1.48 s | 4.95× `[4.79, 5.11]` | 1.65× `[1.55, 1.75]` | **128 MiB** vs. 603 / 1,492 MiB — **4.70× `[4.53, 4.88]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.50 s** / **0.43 s** | 2.36 s / 0.50 s | 5.48× `[5.24, 5.73]` | 1.15× `[1.08, 1.22]` | **107 MiB** vs. 514 / 1,431 MiB — **4.81× `[4.66, 4.95]`** |

#### Against `8238e94e`

Measured in the same blocks (ratios `8238e94e / 29f91ce6`), `googleapis` is 2% faster (1.02× `[1.01, 1.02]`, 119.6 s, faster in every block) with unchanged CPU time: its `CHECKS` stage is 2.3% faster in every block and `OPTIMIZATIONS` 1.4%, the two stages that end by saving their state, while `FINALIZATIONS` is unchanged. The whole-program tasks and `single-files` (75.2 s against 75.6 s) are unchanged within 1% (0.99×–1.00×), and so is their peak RSS.

`googleapis`'s peak RSS varies more from block to block since the previous section's parallel Stage 2 passes: 5.5–6.6 GiB for both builds in this run (6,233 MiB for `29f91ce6` and 6,206 MiB for `8238e94e` on average, against 5,866 MiB for `8238e94e` in its own run). The memory ratio against Java (2.09×, against 2.22× in the previous section) reflects that spread rather than a change between the two builds.

### 8238e94e on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `8238e94e`: `e6727a94` plus `81aa894d`, which parallelizes the function-local Stage 2 passes. Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `e6727a94` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power and kept awake throughout; Java's cold times are within 6% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **59.22 s** | 94.92 s | **1.60× `[1.55, 1.65]`** | **5,440 MiB** | 12,097 MiB | **2.22× `[2.07, 2.39]`** |
| **`OPTIMIZATIONS`** | **54.99 s** | 88.97 s | **1.62× `[1.59, 1.64]`** | **5,866 MiB** | 12,942 MiB | **2.21× `[2.10, 2.32]`** |
| **`FINALIZATIONS`** | **7.59 s** | 19.94 s | **2.63× `[2.61, 2.64]`** | **1,219 MiB** | 7,207 MiB | **5.92× `[5.41, 6.47]`** |
| **Total** | **121.8 s** | **203.8 s** | **1.67× `[1.65, 1.70]`** | **5,866 MiB** | **13,000 MiB** | **2.22× `[2.13, 2.31]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **28.00 s** / **23.60 s** | 49.18 s / 40.33 s | 2.08× `[2.07, 2.10]` | 1.71× `[1.70, 1.72]` | **912 MiB** vs. 4,852 / 6,670 MiB — **5.32× `[5.07, 5.59]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **20.08 s** / **16.14 s** | 35.64 s / 26.54 s | 2.21× `[2.18, 2.24]` | 1.64× `[1.62, 1.67]` | **728 MiB** vs. 3,791 / 4,010 MiB — **5.20× `[5.10, 5.31]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **11.85 s** / **9.52 s** | 23.06 s / 15.40 s | 2.42× `[2.40, 2.44]` | 1.62× `[1.58, 1.65]` | **439 MiB** vs. 3,525 / 3,571 MiB — **8.03× `[7.76, 8.31]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.70 s** / **7.15 s** | 17.82 s / 11.53 s | 2.49× `[2.45, 2.54]` | 1.61× `[1.60, 1.62]` | **336 MiB** vs. 1,746 / 3,679 MiB — **5.20× `[5.03, 5.38]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.91 s** / **3.86 s** | 12.43 s / 6.68 s | 3.22× `[3.16, 3.28]` | 1.73× `[1.69, 1.77]` | **225 MiB** vs. 1,471 / 3,480 MiB — **6.53× `[6.36, 6.70]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.91 s** / **2.99 s** | 11.45 s / 5.51 s | 3.83× `[3.74, 3.91]` | 1.84× `[1.82, 1.86]` | **198 MiB** vs. 1,564 / 1,919 MiB — **7.87× `[7.24, 8.56]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.50 s** / **1.23 s** | 6.00 s / 1.94 s | 4.88× `[4.70, 5.07]` | 1.58× `[1.55, 1.61]` | **147 MiB** vs. 628 / 1,560 MiB — **4.29× `[3.99, 4.60]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.55 s** / **1.94 s** | 6.77 s / 2.85 s | 3.48× `[3.42, 3.54]` | 1.46× `[1.40, 1.53]` | **277 MiB** vs. 1,379 / 3,515 MiB — **4.98× `[4.78, 5.18]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.86 s** / **4.32 s** | 12.81 s / 8.67 s | 2.96× `[2.92, 3.01]` | 2.01× `[1.99, 2.03]` | **234 MiB** vs. 1,650 / 1,857 MiB — **7.05× `[6.71, 7.41]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **75.56 s** | — / 88.71 s | — | 1.17× `[1.16, 1.19]` | **75 MiB** vs. — / 1,311 MiB — **17.4×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.08 s** / **0.91 s** | 4.50 s / 1.52 s | 4.96× `[4.76, 5.16]` | 1.67× `[1.61, 1.74]` | **132 MiB** vs. 607 / 1,478 MiB — **4.61× `[4.45, 4.78]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.51 s** / **0.43 s** | 2.35 s / 0.50 s | 5.45× `[5.28, 5.64]` | 1.15× `[1.10, 1.21]` | **108 MiB** vs. 516 / 1,426 MiB — **4.79× `[4.66, 4.93]`** |

#### Against `e6727a94`

Measured in the same blocks (ratios `e6727a94 / 8238e94e`), the parallel Stage 2 passes make the whole-program `ADVANCED` builds 5–9% faster: `bundle-advanced` 1.09× `[1.09, 1.10]` (now under 3 s), `shaka-blockly` 1.07×, `typed-full` and `typed-mega` 1.06×, `typed-app` and `bundle-mega` 1.05×. The `SIMPLE`, `WHITESPACE_ONLY` and `CHECKS` builds, which run few of these passes, and `bundle-advanced-ts`, whose time is mostly serial (`t=1` 7.70 s against `t=4` 7.15 s), are 1–2% faster. `single-files` takes 75.6 s, against 78.4 s in the previous section (Java's warm time 88.7 s against 89.5 s). CPU time is 2–6% higher on every task, the cost of the extra parallelism, and peak RSS is unchanged except on `shaka-blockly`, 7% higher (0.93× `[0.91, 0.95]`).

`googleapis` is the exception: 0.4% slower overall (1.00× `[0.99, 1.00]`, slower in every block), with 6% more CPU time and 7% more peak RSS (5,866 MiB, now reached in `OPTIMIZATIONS` rather than `CHECKS`). Its `OPTIMIZATIONS` stage is 2.9% slower in every block (0.97×) with 13% more CPU time, while `CHECKS` is 1% and `FINALIZATIONS` 4% faster. That stage already keeps the four threads busy, so the extra parallel work only adds overhead there.

### e6727a94 on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `e6727a94`: `233d5db4` plus two cleanup rounds (`10aad291`: 9 unused `jutil` helpers and 107 dead `*_values()` enums deleted, 18 passes on Go slices and maps; `aa435168`: `ListMultimap` and dead `jutil`/oracle files deleted, `ArrayList`s and `StringBuilder`s replaced across 32 files) and `788cf2cf`, which skips no-op AST compactions and parallelizes `Compactor.Slide`. Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `233d5db4` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was on AC power and kept awake (`caffeinate`) throughout; Java's cold times are within 3% of the previous section's on every task. In block 0 of `checks-only`, all three `jsgo` runs were 10–30% slower than in the other blocks while Java's was not, so `checks-only` was run again right after the suite; its row comes from that rerun.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **59.93 s** | 96.84 s | **1.61× `[1.54, 1.70]`** | **5,508 MiB** | 12,335 MiB | **2.24× `[2.19, 2.29]`** |
| **`OPTIMIZATIONS`** | **53.56 s** | 88.71 s | **1.66× `[1.64, 1.67]`** | **4,955 MiB** | 13,185 MiB | **2.66× `[2.42, 2.93]`** |
| **`FINALIZATIONS`** | **7.92 s** | 20.06 s | **2.53× `[2.49, 2.58]`** | **1,145 MiB** | 7,345 MiB | **6.42× `[5.98, 6.89]`** |
| **Total** | **121.4 s** | **205.6 s** | **1.69× `[1.65, 1.74]`** | **5,508 MiB** | **13,185 MiB** | **2.39× `[2.29, 2.50]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **28.42 s** / **25.18 s** | 50.43 s / 40.85 s | 2.00× `[1.96, 2.04]` | 1.62× `[1.61, 1.64]` | **957 MiB** vs. 5,141 / 6,576 MiB — **5.37× `[5.05, 5.71]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **20.19 s** / **17.30 s** | 36.07 s / 26.48 s | 2.08× `[2.05, 2.12]` | 1.53× `[1.49, 1.57]` | **722 MiB** vs. 3,795 / 5,050 MiB — **5.26× `[5.09, 5.43]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **12.11 s** / **10.32 s** | 23.86 s / 15.70 s | 2.31× `[2.24, 2.39]` | 1.52× `[1.50, 1.54]` | **443 MiB** vs. 3,506 / 3,710 MiB — **7.92× `[7.65, 8.20]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.71 s** / **7.25 s** | 17.94 s / 11.76 s | 2.47× `[2.44, 2.51]` | 1.62× `[1.60, 1.65]` | **331 MiB** vs. 1,730 / 3,783 MiB — **5.22× `[4.95, 5.51]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.97 s** / **4.21 s** | 12.88 s / 6.79 s | 3.06× `[2.98, 3.14]` | 1.61× `[1.58, 1.65]` | **214 MiB** vs. 1,468 / 3,518 MiB — **6.88× `[6.70, 7.05]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.97 s** / **3.32 s** | 11.52 s / 5.71 s | 3.47× `[3.36, 3.58]` | 1.72× `[1.67, 1.77]` | **199 MiB** vs. 1,547 / 3,180 MiB — **7.75× `[7.05, 8.52]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.53 s** / **1.32 s** | 6.37 s / 1.97 s | 4.84× `[4.71, 4.96]` | 1.49× `[1.37, 1.62]` | **147 MiB** vs. 630 / 3,287 MiB — **4.31× `[3.97, 4.67]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.56 s** / **1.95 s** | 6.66 s / 2.81 s | 3.41× `[3.29, 3.54]` | 1.44× `[1.37, 1.51]` | **269 MiB** vs. 1,307 / 3,495 MiB — **4.87× `[4.40, 5.38]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.89 s** / **4.37 s** | 12.93 s / 8.70 s | 2.96× `[2.89, 3.03]` | 1.99× `[1.97, 2.01]` | **235 MiB** vs. 1,652 / 1,787 MiB — **7.04× `[6.72, 7.38]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **78.36 s** | — / 89.48 s | — | 1.14× `[1.14, 1.15]` | **73 MiB** vs. — / 2,517 MiB — **34.5×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.11 s** / **0.96 s** | 4.57 s / 1.55 s | 4.77× `[4.50, 5.06]` | 1.62× `[1.54, 1.69]` | **130 MiB** vs. 606 / 1,479 MiB — **4.64× `[4.49, 4.81]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.52 s** / **0.44 s** | 2.47 s / 0.52 s | 5.57× `[5.41, 5.73]` | 1.17× `[1.08, 1.27]` | **108 MiB** vs. 517 / 1,421 MiB — **4.77× `[4.60, 4.94]`** |

#### Against `233d5db4`

Measured in the same blocks (ratios `233d5db4 / e6727a94`), `googleapis` is 1% faster (1.01× `[1.01, 1.02]`), all of it in its `OPTIMIZATIONS` stage, which is 2.3–2.9% faster in every block (53.6 s), where AST compaction runs; `CHECKS` and `FINALIZATIONS` are unchanged. The whole-program tasks, whose optimizations are much shorter, are within 1% (0.99×–1.02×, every CI including 1.00 except `bundle` 1.01× `[1.00, 1.01]` and `typed-mega` 0.99× `[0.99, 1.00]`). Peak RSS is unchanged within noise on every task. `single-files`, which is not run with the previous build, takes 78.4 s, against 80.8 s in the previous section's rerun, which ran under more background load (Java's warm time 89.5 s against 91.1 s).

### 233d5db4 on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `233d5db4`: `64c098f3` plus three commits that delete more `jutil` collections (`ArrayDeque`, `ImmutableList_Builder`, `LinkedHashMultiset`, `LinkedHashMultimap`, `ArrayListMultimap`, `EnumSet`, `TreeMultimap`, `HashMultimap`, `Sets` and `JavaHashSet`), store `FeatureSet` in a `uint64`, and convert about 25 passes to Go slices and maps (`22d6ff18`, `42a06d12`, `21f33b8b`). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `64c098f3` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). Java's cold times are within 5% of the previous section's on every task. The whole-program tasks and `googleapis` ran on AC power. The laptop then went on battery and slept for 106 s during `single-files`, which was therefore run again on AC power right after; its row comes from that rerun, under light background load (other applications using about two thirds of a core in all), with Java's warm time 4% slower than in the previous section (91.1 s against 87.2 s) and `jsgo`'s 2% slower (80.8 s against 79.1 s).

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **60.46 s** | 95.59 s | **1.58× `[1.57, 1.60]`** | **5,438 MiB** | 12,420 MiB | **2.28× `[2.10, 2.49]`** |
| **`OPTIMIZATIONS`** | **55.33 s** | 89.33 s | **1.61× `[1.59, 1.64]`** | **4,864 MiB** | 13,144 MiB | **2.70× `[2.41, 3.03]`** |
| **`FINALIZATIONS`** | **7.91 s** | 20.21 s | **2.55× `[2.49, 2.62]`** | **1,154 MiB** | 7,355 MiB | **6.38× `[6.05, 6.73]`** |
| **Total** | **123.7 s** | **205.1 s** | **1.66× `[1.64, 1.67]`** | **5,438 MiB** | **13,172 MiB** | **2.42× `[2.22, 2.64]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **28.14 s** / **24.96 s** | 49.51 s / 40.81 s | 1.98× `[1.96, 2.00]` | 1.63× `[1.61, 1.66]` | **921 MiB** vs. 4,940 / 6,560 MiB — **5.37× `[4.99, 5.77]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **19.89 s** / **16.97 s** | 35.39 s / 26.61 s | 2.09× `[2.06, 2.11]` | 1.57× `[1.55, 1.59]` | **725 MiB** vs. 3,799 / 3,977 MiB — **5.24× `[5.03, 5.46]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **11.75 s** / **10.05 s** | 23.22 s / 15.39 s | 2.31× `[2.26, 2.37]` | 1.53× `[1.50, 1.57]` | **442 MiB** vs. 3,523 / 3,614 MiB — **7.98× `[7.56, 8.43]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.70 s** / **7.27 s** | 18.20 s / 11.62 s | 2.50× `[2.41, 2.60]` | 1.60× `[1.58, 1.62]` | **327 MiB** vs. 1,712 / 3,694 MiB — **5.24× `[5.14, 5.34]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.90 s** / **4.14 s** | 12.52 s / 6.76 s | 3.02× `[2.94, 3.11]` | 1.63× `[1.60, 1.67]` | **205 MiB** vs. 1,469 / 3,533 MiB — **7.15× `[7.06, 7.25]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.92 s** / **3.27 s** | 11.51 s / 5.65 s | 3.52× `[3.42, 3.61]` | 1.73× `[1.69, 1.76]` | **197 MiB** vs. 1,588 / 3,369 MiB — **8.04× `[7.44, 8.70]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.50 s** / **1.31 s** | 6.33 s / 1.86 s | 4.83× `[4.65, 5.01]` | 1.42× `[1.35, 1.49]` | **150 MiB** vs. 660 / 1,504 MiB — **4.38× `[3.94, 4.88]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.54 s** / **1.96 s** | 6.87 s / 2.82 s | 3.50× `[3.38, 3.62]` | 1.43× `[1.38, 1.49]` | **259 MiB** vs. 1,323 / 3,516 MiB — **5.11× `[4.72, 5.53]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.88 s** / **4.40 s** | 13.05 s / 8.86 s | 2.97× `[2.89, 3.05]` | 2.02× `[1.99, 2.04]` | **234 MiB** vs. 1,634 / 1,886 MiB — **6.99× `[6.81, 7.19]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **80.76 s** | — / 91.05 s | — | 1.13× `[1.12, 1.14]` | **68 MiB** vs. — / 1,659 MiB — **24.4×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.08 s** / **0.93 s** | 4.47 s / 1.54 s | 4.78× `[4.63, 4.93]` | 1.65× `[1.55, 1.75]` | **127 MiB** vs. 605 / 1,490 MiB — **4.77× `[4.65, 4.89]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.51 s** / **0.43 s** | 2.42 s / 0.50 s | 5.56× `[5.36, 5.77]` | 1.15× `[1.08, 1.22]` | **108 MiB** vs. 518 / 1,422 MiB — **4.79× `[4.71, 4.88]`** |

#### Against `64c098f3`

Neutral to 1% faster: measured in the same blocks (ratios `64c098f3 / 233d5db4`), wall and CPU time are 1.00×–1.02× on every task, with `shaka-blockly`, `typed-full`, `typed-mega`, `bundle-mega` and `googleapis` 1.01× at the edge of their CIs (`googleapis` 1.01× `[1.00, 1.02]`, faster in four blocks of five, spread over its `OPTIMIZATIONS` and `FINALIZATIONS` stages). Peak RSS is unchanged within noise on every task (0.94×–1.05×).

### 64c098f3 on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `64c098f3`: `f18206ae` plus the parse work for 100K+ inputs ([docs/prebuild-ast-scaling.md](prebuild-ast-scaling.md): `0276ff52`, 32-bit source-file ids in a 36-byte node record, and `1e84e983`, per-worker node id ranges), two commits that delete unused `jutil` collections and convert passes to Go slices and maps (`4f621f6f`, `71c154ba`), and the `googleapis` input pre-read (`1a939f5c`). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `f18206ae` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was quiet and on AC power: Java's cold times are within 5% of the previous section's on every task.

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **60.07 s** | 96.06 s | **1.60× `[1.57, 1.63]`** | **5,438 MiB** | 12,522 MiB | **2.30× `[2.20, 2.41]`** |
| **`OPTIMIZATIONS`** | **55.64 s** | 88.85 s | **1.60× `[1.58, 1.61]`** | **4,812 MiB** | 13,270 MiB | **2.76× `[2.63, 2.90]`** |
| **`FINALIZATIONS`** | **7.95 s** | 20.07 s | **2.53× `[2.50, 2.55]`** | **1,175 MiB** | 7,248 MiB | **6.17× `[5.65, 6.74]`** |
| **Total** | **123.7 s** | **205.0 s** | **1.66× `[1.64, 1.68]`** | **5,438 MiB** | **13,270 MiB** | **2.44× `[2.35, 2.53]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **28.09 s** / **24.96 s** | 49.24 s / 40.14 s | 1.97× `[1.96, 1.99]` | 1.61× `[1.60, 1.62]` | **885 MiB** vs. 4,978 / 6,589 MiB — **5.63× `[5.46, 5.80]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **20.18 s** / **17.28 s** | 35.69 s / 26.17 s | 2.07× `[2.05, 2.08]` | 1.51× `[1.50, 1.53]` | **737 MiB** vs. 3,808 / 3,899 MiB — **5.17× `[5.01, 5.33]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **11.99 s** / **10.21 s** | 23.51 s / 15.69 s | 2.30× `[2.25, 2.36]` | 1.54× `[1.50, 1.57]` | **431 MiB** vs. 3,526 / 3,584 MiB — **8.18× `[8.04, 8.33]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.83 s** / **7.42 s** | 18.70 s / 11.74 s | 2.52× `[2.42, 2.62]` | 1.58× `[1.54, 1.62]` | **325 MiB** vs. 1,743 / 3,693 MiB — **5.37× `[5.12, 5.62]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **4.95 s** / **4.19 s** | 12.57 s / 6.75 s | 3.00× `[2.90, 3.09]` | 1.61× `[1.58, 1.64]` | **210 MiB** vs. 1,482 / 3,506 MiB — **7.06× `[6.68, 7.47]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **3.94 s** / **3.28 s** | 11.49 s / 5.69 s | 3.51× `[3.43, 3.59]` | 1.74× `[1.70, 1.77]` | **197 MiB** vs. 1,682 / 3,491 MiB — **8.54× `[7.63, 9.56]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.50 s** / **1.30 s** | 6.08 s / 1.88 s | 4.69× `[4.59, 4.80]` | 1.45× `[1.43, 1.47]` | **145 MiB** vs. 672 / 1,566 MiB — **4.62× `[3.59, 5.94]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.55 s** / **1.97 s** | 6.86 s / 2.81 s | 3.48× `[3.41, 3.56]` | 1.43× `[1.38, 1.47]` | **261 MiB** vs. 1,359 / 3,484 MiB — **5.21× `[4.99, 5.44]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **4.92 s** / **4.45 s** | 13.27 s / 8.78 s | 2.98× `[2.89, 3.07]` | 1.97× `[1.94, 2.00]` | **236 MiB** vs. 1,658 / 1,866 MiB — **7.03× `[6.70, 7.37]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **79.09 s** | — / 87.21 s | — | 1.10× `[1.09, 1.11]` | **69 MiB** vs. — / 1,202 MiB — **17.5×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.09 s** / **0.95 s** | 4.43 s / 1.49 s | 4.68× `[4.54, 4.82]` | 1.58× `[1.50, 1.66]` | **126 MiB** vs. 605 / 1,489 MiB — **4.80× `[4.68, 4.92]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.51 s** / **0.44 s** | 2.42 s / 0.50 s | 5.54× `[5.34, 5.75]` | 1.15× `[1.07, 1.23]` | **108 MiB** vs. 514 / 1,423 MiB — **4.77× `[4.56, 4.99]`** |

#### Against `f18206ae`

Measured in the same blocks (ratios `f18206ae / 64c098f3`), the typed and larger `ADVANCED` builds are 1–4% faster, with CPU time down by as much: `typed-mega` 1.04× `[1.03, 1.06]`; `checks-only`, `bundle-advanced` and `typed-full` 1.03×; `shaka-blockly` and `bundle-mega` 1.02×; `googleapis` 1.01× `[1.01, 1.01]`, all of it in its `CHECKS` stage, which is 2.1–2.7% faster in every block (the other two stages are unchanged). `bundle`, `bundle-advanced-ts`, `typed-app` and the two `downlevel` tasks are within 2%. `single-files`, which is not run with the previous build, takes 79.1 s against 82.3 s in the previous section (4% less). Peak RSS is higher: 7% on `typed-mega` (0.93× `[0.91, 0.97]`, 737 MiB against 689 MiB), about 4% on `googleapis` (5,438 against 5,217 MiB) and `bundle-mega` (885 against 853 MiB), and 4–6% within noise on most other tasks, in line with the about 5% more peak live heap that [prebuild-ast-scaling.md](prebuild-ast-scaling.md) gives for the larger node record. This run does not separate the four commits.

#### The input pre-read

This is the first run with [the `googleapis` input pre-read](#googleapis-and-the-file-cache). `jsgo`'s five `googleapis` runs are now within 0.9 s of each other (123.2–124.1 s, against 124.4–127.3 s for `f18206ae` in the previous section), whatever ran before them. Java's are within 3.9 s (202.8–206.7 s, against 199.5–208.7 s) with the same mean (205.0 s against 204.9 s): its `CHECKS` stage now takes 94.4–97.7 s in every block, where it took about 92 s after a `jsgo` run and about 99.5 s after a Java run. So the pre-read removes the order effect without making Java faster on average. The higher speedup against Java (1.66×, against 1.62×) comes from `jsgo`'s side.

### f18206ae on JDK 27 (Apple M2)

$N = 5$ counterbalanced interleaved suite with 95% CIs, at commit `f18206ae` (`a95c4a93` plus `5cd39f45` and `63a53705`, which replace lookup-only `jutil` collections with Go maps, slices and switches in 38 passes and delete `LinkedIdentityHash{Map,Set}`). Same machine, method, gate and Java (v20261021 on OpenJDK 27.0.1) as [the v20261021 section](#v20261021-on-jdk-27-apple-m2), at `--num_parallel_threads=4`, with `a95c4a93` in the same blocks. Every `jsgo` run was byte-identical to Java (in `single-files`, apart from the Java stack frames of the same 11 inputs on which both compilers crash). The machine was quiet and on AC power: Java's cold times are within 6% of the previous section's on every task. This is the last run without the input pre-read described [below](#googleapis-and-the-file-cache).

#### 1. Production-Scale Staged Compilation (`googleapis-bench`: 43,917 Files, 35.6M Parsed AST Nodes, 7,460 Chunks)

| Stage | `jsgo` (`t=4`) | Java (`t=4`, `-Xmx20g`) | Speedup vs. Java (`[95% CI]`) | `jsgo` Peak RSS | Java Peak RSS | Memory Ratio vs. Java (`[95% CI]`) |
|---|---:|---:|---:|---:|---:|---:|
| **`CHECKS`** | **62.85 s** | 96.39 s | **1.53× `[1.50, 1.57]`** | **5,229 MiB** | 12,187 MiB | **2.33× `[2.23, 2.44]`** |
| **`OPTIMIZATIONS`** | **55.31 s** | 88.57 s | **1.60× `[1.59, 1.62]`** | **4,714 MiB** | 13,074 MiB | **2.78× `[2.56, 3.02]`** |
| **`FINALIZATIONS`** | **7.93 s** | 19.90 s | **2.51× `[2.50, 2.52]`** | **1,167 MiB** | 7,332 MiB | **6.29× `[5.88, 6.72]`** |
| **Total** | **126.1 s** | **204.9 s** | **1.62× `[1.60, 1.64]`** | **5,229 MiB** | **13,074 MiB** | **2.50× `[2.37, 2.64]`** |

#### 2. Large & Mega Whole-Program `ADVANCED` Workloads

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`bundle-mega`** | ADV · 4.19M nodes | **28.39 s** / **25.54 s** | 49.28 s / 40.51 s | 1.93× `[1.91, 1.95]` | 1.59× `[1.56, 1.61]` | **846 MiB** vs. 4,812 / 6,594 MiB — **5.68× `[5.06, 6.36]`** |
| **`typed-mega`** | ADV+types · 3.55M nodes | **20.28 s** / **17.75 s** | 35.55 s / 26.62 s | 2.00× `[1.97, 2.03]` | 1.50× `[1.49, 1.51]` | **712 MiB** vs. 3,791 / 4,192 MiB — **5.33× `[5.08, 5.58]`** |
| **`typed-full`** | ADV+types · 1.90M nodes | **11.97 s** / **10.45 s** | 23.14 s / 15.59 s | 2.21× `[2.18, 2.25]` | 1.49× `[1.45, 1.53]` | **416 MiB** vs. 3,520 / 3,629 MiB — **8.46× `[8.19, 8.73]`** |
| **`bundle-advanced-ts`** | ADV · 1.11M nodes | **7.82 s** / **7.35 s** | 17.87 s / 11.51 s | 2.43× `[2.38, 2.48]` | 1.57× `[1.55, 1.58]` | **321 MiB** vs. 1,746 / 3,875 MiB — **5.44× `[5.27, 5.61]`** |
| **`shaka-blockly`** | ADV+types · 618k nodes | **5.03 s** / **4.33 s** | 13.04 s / 6.86 s | 3.01× `[2.91, 3.11]` | 1.58× `[1.54, 1.63]` | **203 MiB** vs. 1,461 / 3,490 MiB — **7.20× `[7.02, 7.39]`** |
| **`bundle-advanced`** | ADV · 573k nodes | **4.05 s** / **3.45 s** | 12.04 s / 5.66 s | 3.49× `[3.45, 3.53]` | 1.64× `[1.63, 1.66]` | **188 MiB** vs. 1,585 / 1,935 MiB — **8.41× `[7.56, 9.35]`** |

#### 3. Standard Whole-Program & CLI Suite

| Benchmark | Mode & Parsed Nodes | `jsgo` (`t=1` / `t=4`) | Java `t=4` (Cold / Warm) | Speedup vs. Cold (`[95% CI]`) | Speedup vs. Warm (`[95% CI]`) | Peak RSS (`jsgo` vs. Cold / Warm) & Cold Ratio (`[95% CI]`) |
|---|---|---:|---:|---:|---:|---:|
| **`typed-app`** | ADV+types · 1.29M nodes | **1.52 s** / **1.31 s** | 6.35 s / 1.94 s | 4.84× `[4.71, 4.97]` | 1.48× `[1.45, 1.52]` | **139 MiB** vs. 633 / 1,560 MiB — **4.56× `[4.26, 4.89]`** |
| **`checks-only`** | CHECKS+types · 1.29M nodes | **2.58 s** / **2.04 s** | 7.04 s / 2.82 s | 3.46× `[3.40, 3.51]` | 1.38× `[1.31, 1.46]` | **257 MiB** vs. 1,342 / 3,503 MiB — **5.21× `[4.79, 5.67]`** |
| **`bundle`** | SIMPLE · 1.11M nodes | **5.02 s** / **4.52 s** | 13.46 s / 8.74 s | 2.98× `[2.87, 3.09]` | 1.93× `[1.91, 1.96]` | **228 MiB** vs. 1,652 / 1,841 MiB — **7.25× `[6.98, 7.52]`** |
| **`single-files`** | 1,000 CLI runs, SIMPLE · 306k nodes | **82.29 s** | — / 90.25 s | — | 1.10× `[1.08, 1.11]` | **68 MiB** vs. — / 1,355 MiB — **19.9×** vs. batch |
| **`downlevel-simple`** | SIMPLE→ES5 · 232k nodes | **1.12 s** / **0.99 s** | 4.65 s / 1.51 s | 4.70× `[4.56, 4.85]` | 1.53× `[1.46, 1.61]` | **124 MiB** vs. 603 / 1,479 MiB — **4.86× `[4.71, 5.02]`** |
| **`downlevel-whitespace`** | WHITESPACE_ONLY→ES5 · 232k nodes | **0.52 s** / **0.45 s** | 2.45 s / 0.51 s | 5.44× `[5.29, 5.60]` | 1.13× `[1.06, 1.20]` | **107 MiB** vs. 522 / 1,424 MiB — **4.89× `[4.69, 5.10]`** |

#### Against `a95c4a93`

Neutral in time: measured in the same blocks (ratios `a95c4a93 / f18206ae`), wall and CPU time are 0.98×–1.00× on every task. Peak RSS is 5% lower on `typed-full` (1.05× `[1.03, 1.06]`, 416 MiB against 436 MiB) and unchanged within noise elsewhere (`bundle-mega` 1.04× `[1.00, 1.08]`, 846 MiB against 877 MiB). `googleapis` is unchanged: 0.99× `[0.98, 1.00]` over all five blocks, and 1.00× (0.998 and 1.000) in the two blocks where neither build ran right after Java, which are the ones to read (next section).

#### `googleapis` and the file cache

The `googleapis` comparisons between two `jsgo` builds in this section and the ones below carry a bias against the newer build. A Java `googleapis` run (14 GiB on this 24 GiB machine) evicts the 44,240 input files (485 MB) from the file cache: reading them all takes 8.4 s right after Java's `CHECKS` stage, against 1.2 s otherwise. Whichever run comes next reads its inputs from disk, and its `CHECKS` stage is slower. Run by hand, `jsgo`'s `CHECKS` stage took 64.7 s right after Java, against 61.7 s after a `jsgo` run and 61.8 s after Java with the inputs read again first. In the suite, Java's own `CHECKS` stage was 3–8% slower after another Java run.

The block order (Java, head, previous; then the reverse) puts the head build right after Java in blocks 0, 2 and 4, and the previous build never. In every section from `41b39832` on, the head build's `CHECKS` stage is 3–5% slower than the previous build's in those three blocks, and within 1% in blocks 1 and 3 (2.3% once). Read from blocks 1 and 3 only, the head-against-previous `googleapis` ratios (`previous / head`) are:

| Section | All five blocks (as first published) | Blocks 1 and 3 |
|---|---:|---:|
| `41b39832` against `a8c90bf9` | 0.97× | 0.98× (0.984, 0.981) |
| `fbdac472` against `41b39832` | 0.98× | 0.99× (0.995, 0.995) |
| `692a3752` against `fbdac472` | 0.99× | 1.00× (1.002, 1.001) |
| `50da0981` against `692a3752` | 0.98× | 0.99× (0.988, 0.998) |
| `a95c4a93` against `50da0981` | 0.99× | 1.00× (0.994, 1.001) |
| `f18206ae` against `a95c4a93` | 0.99× | 1.00× (0.998, 1.000) |

The comparison paragraphs below are corrected accordingly. The speedups against Java are hardly affected, because Java also ran right after a Java run in blocks 2 and 4, and in block 0 after the warm-JVM runs of the previous task. Since `1a939f5c`, [`interleaved.py`](../tools/bench/interleaved.py) reads the inputs, untimed, before every `googleapis` run, so no run pays for a cold cache. From the next section on, every `googleapis` run starts with its inputs cached ([effect](#the-input-pre-read)).

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

Measured in the same blocks (ratios `692a3752 / 50da0981`), the whole-program and CLI tasks are unchanged (wall and CPU time 0.99×–1.01×, peak RSS within noise). `googleapis` is unchanged too: 0.99× (0.988 and 0.998) in the two blocks where neither build ran right after Java, with CPU time and peak RSS unchanged. (This paragraph first said 2% slower, all of it in the `CHECKS` stage, from all five blocks; that was the file-cache bias described in [the f18206ae section](#googleapis-and-the-file-cache).)

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

Performance-neutral: measured in the same blocks (ratios `41b39832 / fbdac472`), wall time is 0.99×–1.00× (`googleapis` 0.99× in the two blocks where neither build ran right after Java; 0.98× over all five blocks, [biased](#googleapis-and-the-file-cache)), CPU time 0.99×–1.00× and peak RSS 0.97×–1.03× on every task.

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

Measured in the same blocks (ratios `a8c90bf9 / 41b39832`, paired 95% CIs), those commits cost 0–2% of wall time (0.98×–1.00×; `googleapis` 0.98× in the two blocks where neither build ran right after Java, against 0.97× `[0.96, 0.98]` over all five blocks, which is [biased](#googleapis-and-the-file-cache); `shaka-blockly` and `typed-mega` 0.98×, the rest 0.99×–1.00×) and 0–2% of CPU time, with peak RSS unchanged within noise on every task (0.95×–1.04×).

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

**Method.** [`tools/bench/interleaved.py`](../tools/bench/interleaved.py) runs and gates the suite and [`interleaved_stats.py`](../tools/bench/interleaved_stats.py) computes the statistics. Every task runs $N = 5$ blocks; each block runs every contender once, in the order Java, `jsgo` builds in odd blocks and the reverse in even ones, which cancels linear thermal and background drift. Times and peak RSS (`wait4`'s `ru_maxrss`) are sample means. Ratios (`Java / jsgo`) are geometric means with two-sided 95% confidence intervals: paired Student-$t$ ($df = 4$) on per-block log-ratios against cold Java, and Welch's $t$ on log-times against warm Java. Warm Java is `tools/bench/java/CliBatch` in one JVM: 5 measured iterations after 5 warm-up iterations (2 for `single-files`); its peak RSS is the JVM's over its whole life. `single-files` (1,000 invocations) is compared with the warm batch only, as before. `googleapis` is the three-stage build of `tools/bench/run_googleapis_bench.py` (Java with `-Xmx20g`); its rows are per stage and their sum. Since `1a939f5c`, the 44,240 `googleapis` input files are read, untimed, before every `googleapis` run: a Java run evicts them from the file cache, and the order above always put the head build right after Java ([details](#googleapis-and-the-file-cache)).

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
