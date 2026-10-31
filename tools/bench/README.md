# Benchmarks

Two benchmarks live here:

- **End-to-end** (this file): jsgo against the pinned Java compiler on real
  CLI invocations, as specified in [PLAN.md §6](../../PLAN.md). Overall and
  per-phase numbers, behind a byte-for-byte correctness gate.
- **Parser** (milestone M1): `java/ParseBench.java`, `parsebench/` and
  `measure.py`; see [docs/benchmarks.md](../../docs/benchmarks.md).

## Quick start

```sh
tools/bench/fetch-inputs.sh                      # pinned inputs -> .cache/bench-inputs/
python3 -I tools/bench/bench.py run --out /tmp/bench --report /tmp/bench/report.md
```

That runs every task on Java cold, Java warm and jsgo, with N=5 runs each.
It takes several hours on 4 cores, most of it in the 1,000-invocation task
(see [Cost](#cost)). `--tasks a,b` selects tasks, and `bench.py report --out
/tmp/bench --md report.md` re-renders the report from stored results.

Everything uses the `.cache/` of the main checkout, shared by all worktrees;
`$JSGO_CACHE` overrides it. It holds the jar (`CLOSURE_JAR` of
[tools/oracle/pin.env](../oracle/pin.env), the one place that names the
pinned release, v20261021, and its jar), the corpus, `bench-inputs/` and
`bench-build/`.

## Files

| File | What |
|---|---|
| `inputs.txt` | Third-party inputs, pinned by version and checksum: npm sha512 integrity, or git commit + tree + a sha256 of the file listing |
| `fetch-inputs.sh` | Fetches and verifies them into `.cache/bench-inputs/<name>@<version>/`. Never commit them |
| `app/` | The typed sample app: original code, three files |
| `tasks/*.json` | One file per task: the exact flags, inputs and outputs, identical for both compilers |
| `tasks/single-files.txt` | The 1,000 files of the single-file task (`bench.py sample [--check]`) |
| `java/CliBatch.java` | The warm Java driver: repeats the same `CommandLineRunner` invocations in one JVM |
| `java/MaxLoopRunner.java`, `maxloop_oracle.py` | Java with `CompilerOptions.setMaxOptimizationLoopIterations`, read from jsgo's hidden `--max_optimization_loop_iterations` flag, and the byte-for-byte comparison of both on the tasks (fast mode round 3 in [docs/performance-roadmap.md](../../docs/performance-roadmap.md)) |
| `phases.tsv` | Pass name → phase, used for the per-phase roll-up |
| `bench.py` | Runner, correctness gate, tracer-report parser and Markdown report |
| `pgo.py` | Regenerates `cmd/jsgo/default.pgo`, the profile for profile-guided optimisation (see [PGO profile](#pgo-profile)) |

## Inputs

| Input | Version | Used by |
|---|---|---|
| Closure Library | v20230802, the last release; GitHub tag `ccf3f1dc` | typed app, checks only, shaka-blockly, typed-full, typed-mega |
| Shaka Player | 5.3.1 (`lib/`, `ui/`, `externs/`, `third_party/`, `dist/locales.js`) | shaka-blockly, typed-full, typed-mega |
| Blockly | 8.0.5 (`core/`, `blocks/`, `generators/`, last pure-Closure release) | shaka-blockly, typed-full, typed-mega |
| OpenSphere | release-23.0.1 (1,648 Closure `goog.provide`/`goog.module` files + externs) | typed-mega |
| OpenLayers | 4.6.5 (`src/ol/`, `externs/`, 539 files) | typed-mega |
| Camptocamp ngeo | 2.2.10 (`src/`, `contribs/`, `externs/`, 203 files) | typed-mega |
| YouTube SPFjs | 2.4.0 (`src/`, 36 files) | typed-mega |
| ol-cesium | 1.37.0 (`src/olcs/`, 40 files) | typed-mega |
| google-protobuf | 3.21.2 (Closure JS runtime + UMD bundle) | typed-mega, bundle-mega |
| three.js | 0.159.0 (`build/three.js`) | bundle, bundle-advanced, bundle-advanced-ts, bundle-mega |
| lodash | 4.17.21 (`lodash.js`) | bundle, bundle-advanced, bundle-advanced-ts, bundle-mega |
| react-dom | 18.3.1 (`umd/react-dom.development.js`) | bundle, bundle-advanced, bundle-advanced-ts, bundle-mega |
| TypeScript | 5.6.3 (`lib/typescript.js`, 8.9 MB) | bundle, bundle-advanced-ts |
| Large UMD suite | plotly.js 2.27.0, handsontable 6.2.2, echarts 5.4.3, ag-grid-community 30.2.1, pdfjs-dist 2.16.105, pixi.js 6.5.10, fabric 5.3.0, mathjs 11.11.2, konva 8.4.3, google-libphonenumber 3.2.33, lovefield 2.1.12, chart.js 2.9.4, zone.js 0.14.2 | bundle-mega |
| @angular/core | 16.2.12 (`fesm2022/core.mjs`) | down-levelling |
| rxjs | 7.8.2 (`dist/esm`, ES2015 modules; `dist/bundles/rxjs.umd.js`) | down-levelling, bundle-advanced, bundle-mega |
| tslib | 2.8.1 (`tslib.es6.js`) | down-levelling |
| core-js | 3.50.0 (CommonJS, from `stable/index.js`) | down-levelling |
| oracle corpus | `tools/oracle/corpus.txt` | single files, bundle-advanced, bundle-mega |

Choices:

- **Typed app.** Closure Library's demos are HTML pages with inline scripts,
  so they cannot be compiled as they are. `app/` is a small notes editor
  written for this benchmark instead, in the style of the library's editor
  demo:
  - `notestore.js` uses `goog.provide`; `notesview.js` and `main.js` use
    `goog.module`.
  - It uses goog.editor with ten plugins and the toolbar, goog.ui (TabBar,
    DatePicker, Button), goog.net.XhrIo, goog.storage and goog.i18n.
  - It also uses the library's own demo plugins,
    `goog.demos.editor.HelloWorld` and `HelloWorldDialogPlugin`.
  - The whole library (971 files) is passed. `--dependency_mode=PRUNE` from
    `goog:benchapp.main` keeps about 280 of them, a mix of `goog.module` and
    `goog.provide` files.
  - It compiles with 0 errors and 98% typed.
- **Bundle.** The four libraries are compiled as plain scripts into one
  output.
  - three.js r159 is the last release with a script (UMD) build. Later
    releases ship `three.cjs`, whose top-level `class Audio`/`const Cache`
    collide with the browser externs, and ES modules.
- **Down-levelling.**
  - Angular 16.x is the last major whose `core.mjs` imports nothing but
    `rxjs`. From 17 on it imports `@angular/core/primitives/*`, which only
    package.json `exports` maps resolve. Closure's NODE resolver does not
    read those maps, and a `--browser_resolver_prefix_replacements` value
    cannot start with `@` (args4j reads it as an @-file).
  - The task's working directory gets a `node_modules/` layout of symlinks
    into `.cache/bench-inputs`, and the package.json files are passed as
    inputs. `--module_resolution=NODE`, `--process_common_js_modules` and
    PRUNE from the two entry points (Angular core, core-js `stable`) then
    pick the files.
  - The same flags and inputs run under WHITESPACE_ONLY and SIMPLE, with
    `--language_out=ECMASCRIPT5`.
- **Single files.** A sample of 1,000 files from `.cache/corpus-files.txt`:
  - drawn with `random.Random(20260415).sample`;
  - excluding test262-parser-tests `fail/` and `early/`, which are invalid on
    purpose;
  - excluding files over 256 KiB, so that each invocation stays small.

  Every file is compiled alone, with its output on stdout.

## Tasks

| Task file | Mode | Outputs compared |
|---|---|---|
| `typed-app` | `ADVANCED`, `--jscomp_warning=checkTypes`, source map | stdout, stderr, exit, `out.js`, `out.js.map` |
| `bundle` | `SIMPLE` | stdout, stderr, exit, `out.js` |
| `downlevel-whitespace` | `WHITESPACE_ONLY`, `--language_out=ECMASCRIPT5` | stdout, stderr, exit, `out.js` |
| `downlevel-simple` | `SIMPLE`, `--language_out=ECMASCRIPT5` | stdout, stderr, exit, `out.js` |
| `checks-only` | `--checks_only`, `--jscomp_warning=checkTypes` | stdout, stderr, exit |
| `single-files` | `SIMPLE`, 1,000 processes per run | stdout, stderr, exit, per invocation |
| `shaka-blockly` | `ADVANCED`, `--jscomp_warning=checkTypes`, source map (847 files, 8.6 MB: Shaka Player 5.3.1 + Blockly 8.0.5) | stdout, stderr, exit, `out.js`, `out.js.map` |
| `typed-full` | `ADVANCED`, `--jscomp_warning=checkTypes`, source map (1,818 files, 28.5 MB: Closure Library + Shaka Player 5.3.1 + Blockly 8.0.5 + app) | stdout, stderr, exit, `out.js`, `out.js.map` |
| `bundle-advanced` | `ADVANCED`, `--jscomp_warning=checkVars`, source map (11 UMD libraries, 5.4 MB) | stdout, stderr, exit, `out.js`, `out.js.map` |
| `bundle-advanced-ts` | `ADVANCED`, `--jscomp_warning=checkVars`, source map (12.6 MB incl. `typescript.js`) | stdout, stderr, exit, `out.js`, `out.js.map` |
| `typed-mega` | `ADVANCED`, `--jscomp_warning=checkTypes`, source map (4,385 files, 44.6 MB: 10 Closure `ADVANCED` codebases) | stdout, stderr, exit, `out.js`, `out.js.map` |
| `bundle-mega` | `ADVANCED`, `--jscomp_warning=checkVars`, source map (26 UMD libraries, 36.1 MB) | stdout, stderr, exit, `out.js`, `out.js.map` |

Input paths are always relative: `inputs/…`, `app/…`, `corpus/…` and
`node_modules/…`. Each compile runs in its own directory holding those
symlinks, and outputs (`out.js`, `out.js.map`) are written there. Both
compilers therefore see the same argv, and source maps name the same paths.
`<results>/<task>/jobs.json` records the fully expanded argv; globs are
expanded in sorted order.

## What is measured

- **Java cold:** `java -jar closure-compiler.jar <args>` (the cache's
  `CLOSURE_JAR`), N runs (default 5).
- **Java warm:** `CliBatch`, one JVM per task.
  - It runs W warm-up iterations (default 2; 1 for the bundle and the
    single files), then N measured ones.
  - Each iteration runs the task's invocations exactly as
    `CommandLineRunner.main` would.
  - It reports wall time and process CPU per iteration (JIT and GC threads
    included), and the JVM's peak RSS over its whole life.
- **jsgo:** `jsgo <args>`, N runs, interleaved with the Java cold runs so
  that changes in load hit both alike.

For each run the report gives wall time, CPU time (user + system) and peak
RSS from `wait4`, as the median over the runs with its min–max range.

All measured processes run under `nice -n 10` (`--nice`). The 1-minute load
average is recorded before every run.

The environment is fixed:
- `LC_ALL=C.UTF-8` (`--locale`).
- For Java, `JAVA_TOOL_OPTIONS`, `_JAVA_OPTIONS` and `JDK_JAVA_OPTIONS` are
  removed, so the JVM runs with default settings.
- For jsgo, `GOGC`, `GOMEMLIMIT`, `GOMAXPROCS` and `GODEBUG` are removed, so
  it runs with the default Go GC.

**Correctness gate.** jsgo's first run is compared byte for byte with Java
cold run 1: stdout, stderr (the diagnostics), exit status and the task's
output files. The result is one of:

- **MATCH:** everything is identical. Only MATCH tasks report speed-ups
  (Java time ÷ jsgo time).
- **PENDING:** jsgo reached a member that is not ported yet (a `*PendingPort`
  panic, `pending port: <Java signature>` on stderr). The report names the
  signature, and jsgo is not timed further.
- **DIFF:** anything else. The report shows the first differing byte of
  each differing artifact, and the times without a ratio. DIFFs are sorted
  into three kinds:
  - **crash:** jsgo died with a Go panic or fatal error. A jsgo that crashes
    on every invocation is not timed.
  - **java-frames:** Java threw an internal exception, and only its stack
    frames differ. jsgo prints the same exception, message, exit status and
    output, but not the JVM's `\tat …` frames.
  - **output:** any other difference.

For the single-file task, every invocation is gated, and the task gets the
counts of each result. Examples of each DIFF kind are listed by input path;
each is a one-file reproducer with the task's flags. The task's speed-up is
then computed over the invocations that match.

`bench.py regate --out <dir>` recomputes `gate.json` from the stored first
runs without compiling anything, for example after a change to the gate.

Java is also checked against itself:
- every cold run must be byte-identical to run 1;
- the warm driver's last iteration must reproduce cold run 1.

The report says when either check fails.

**Per phase.** Separate runs append `--tracer_mode=TIMING_ONLY` (3 per
compiler by default, `--trace-runs`). `bench.py` reads the `Summary:`
section of each report (pass, runtime, runs) and the tracker's start and
end time, then:

- lines the passes up by name, and checks the names and run counts against
  each other: Java cold vs jsgo, and Java cold vs warm;
- rolls the passes up into phases with `phases.tsv`:
  - parse;
  - checks;
  - type checking;
  - transpilation (module rewriting, down-levelling, polyfills);
  - optimization loop;
  - renaming;
  - other passes.

A pass missing from `phases.tsv` is counted as "other passes" and listed as
unmapped.

`bench.py report` rolls the stored pass times up again with the current
`phases.tsv`, so mapping a new pass name also corrects the reports of
results measured before.

Three rows do not come from passes. `PerformanceTracker` prints its report
before the output is generated, so the timestamps cover what the passes do
not:

- **startup:** process start (or, warm, iteration start) to the tracker's
  start.
- **untraced (inside tracker):** the tracker's wall time not spent in any
  pass.
- **code printing and source maps:** the tracker's end to process exit (or
  iteration end): toSource, the source map, writing the output, and exit.
  - With several invocations per warm iteration, startup and printing cannot
    be told apart, so their sum appears on this row.

## Running jsgo later

The jsgo binary comes from one of:
- `--jsgo PATH`: an existing binary;
- `--jsgo-ref REF`: `git archive REF` is built in `<out>/jsgo-src`;
- by default, `./cmd/jsgo` of the current checkout.

If none of these exists, jsgo is reported as unavailable.

Java results stay valid while the jar is unchanged (a new pin in
`tools/oracle/pin.env` invalidates them). Once more of the port
lands, re-measure only jsgo against the stored Java reference outputs:

```sh
python3 -I tools/bench/bench.py run --out /tmp/bench --skip-java --jsgo-ref <branch> \
    --report /tmp/bench/report.md
```

These jsgo runs are then not interleaved with Java runs. For final numbers,
re-run everything on a quiet machine and render with `bench.py report
--final`.

## Interleaved suite

`interleaved.py` is the lighter runner behind the README and the
"v20261021 on JDK 27" section of [docs/benchmarks.md](../../docs/benchmarks.md):
cold Java against one or more jsgo binaries (each with its own
`--num_parallel_threads`), $N$ blocks in alternating order, every run gated
byte for byte against Java's first run, plus optional warm Java
(`java/CliBatch.java`) and the staged `googleapis` build:

```sh
python3 -I tools/bench/interleaved.py --out /tmp/suite --java /opt/homebrew/opt/openjdk/bin/java \
    --blocks 5 --warm --jsgo final=./jsgo@4 --jsgo final1=./jsgo@1 --tasks typed-full,bundle
python3 -I tools/bench/interleaved_stats.py /tmp/suite/results.jsonl
```

`--no-java --java-ref DIR` compares jsgo builds alone against the Java
references an earlier run left in `DIR/<task>/ref`. `interleaved_stats.py`
prints means and geometric-mean ratios with 95% CIs (paired $t$ against cold
Java, Welch's $t$ against warm Java). Unlike `bench.py`, it records no
per-phase traces.

## PGO profile

`go build ./cmd/jsgo` (and so this benchmark) applies profile-guided
optimisation with `cmd/jsgo/default.pgo`, a CPU profile of jsgo on these
tasks. The compiler uses it to inline hot calls and to call the usual
implementation of hot interface methods directly; the program's behaviour
does not depend on it. A stale profile still works (functions it no longer
matches are just not optimised), but regenerate it after large changes:

```sh
python3 -I tools/bench/pgo.py      # writes cmd/jsgo/default.pgo
```

`pgo.py` builds jsgo without the old profile and with a CPU-profiling hook
(through `-overlay`; no source changes), sets up the task workspaces as
`bench.py` does, runs the first invocation of typed-app,
downlevel-whitespace, downlevel-simple and checks-only and eight evenly
spaced single-file invocations, and merges the profiles. It takes a few
minutes. `go build -pgo=off` builds jsgo without the profile.

## Results layout

```
<out>/machine.json                  machine, toolchains, jar and jsgo identity, load averages
<out>/<task>/jobs.json              the expanded invocations
<out>/<task>/java-cold.json         per-run wall/CPU/RSS/load/exit and output digests
<out>/<task>/java-warm.json         per-iteration driver report, fidelity check
<out>/<task>/jsgo.json, gate.json   jsgo runs; the correctness gate and its first differences
<out>/<task>/trace.json             per-run phase and pass times for each compiler
<out>/<task>/java-cold/run-1/       Java's reference outputs (later runs keep only small files)
<out>/<task>/jsgo/run-1/            jsgo's outputs
```

## Cost

Wall time per run on the 4-CPU machine with nothing else running, from the
v20261015 run (results in
[docs/benchmarks.md](../../docs/benchmarks.md#end-to-end-v20261015)):

| Task | Java cold, per run | Java warm, per iteration | jsgo, per run |
|---|---|---|---|
| typed-app | 11–12 s | 4.5–5.1 s | 7.0–7.4 s |
| bundle | 26–27 s (12.6 MB in, 1.7 GB RSS) | 18–20 s (1.7 GB RSS) | 28–29 s (0.55 GB RSS) |
| downlevel-whitespace | 3.8–4.2 s | 1.0–1.5 s | 1.7–1.8 s |
| downlevel-simple | 7.7–8.7 s | 3.1–4.0 s | 5.3–6.3 s |
| checks-only | 13–14 s | 6.7–6.9 s | 9.8–10.2 s |
| single-files | 1,000 JVM launches, 1.1–2.1 s each (p10–p90): 30 min | 4.6 min | 7 min |

A full `bench.py run` took 4 h 15 min:
- the first five tasks together took 23 minutes (N=5, three trace runs);
- single-files took 3 h 50 min, two thirds of it in the 5,000 cold
  JVM launches.

On a loaded machine everything takes longer, and Java's times grow the
most: a cold JVM keeps about three cores busy. At load 5–17, the first
baseline measured typed-app at 17–37 s cold and the single-file runs at
40–75 min.

`--runs 1 --tasks single-files` gives one run, and the report then shows
per-invocation median, p10 and p90 from its 1,000 samples. Its cold
per-phase runs are off (`"trace_cold": false` in the task file).
