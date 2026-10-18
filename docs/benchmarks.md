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
