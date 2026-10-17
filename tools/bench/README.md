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
`$JSGO_CACHE` overrides it. It holds the jar, the corpus, `bench-inputs/`
and `bench-build/`.

## Files

| File | What |
|---|---|
| `inputs.txt` | Third-party inputs, pinned by version and checksum: npm sha512 integrity, or git commit + tree + a sha256 of the file listing |
| `fetch-inputs.sh` | Fetches and verifies them into `.cache/bench-inputs/<name>@<version>/`. Never commit them |
| `app/` | The typed sample app: original code, three files |
| `tasks/*.json` | One file per task: the exact flags, inputs and outputs, identical for both compilers |
| `tasks/single-files.txt` | The 1,000 files of the single-file task (`bench.py sample [--check]`) |
| `java/CliBatch.java` | The warm Java driver: repeats the same `CommandLineRunner` invocations in one JVM |
| `phases.tsv` | Pass name → phase, used for the per-phase roll-up |
| `bench.py` | Runner, correctness gate, tracer-report parser and Markdown report |

## Inputs

| Input | Version | Used by |
|---|---|---|
| Closure Library | v20230802, the last release; GitHub tag `ccf3f1dc` | typed app, checks only |
| three.js | 0.159.0 (`build/three.js`) | bundle |
| lodash | 4.17.21 (`lodash.js`) | bundle |
| react-dom | 18.3.1 (`umd/react-dom.development.js`) | bundle |
| TypeScript | 5.6.3 (`lib/typescript.js`, 8.9 MB) | bundle |
| @angular/core | 16.2.12 (`fesm2022/core.mjs`) | down-levelling |
| rxjs | 7.8.2 (`dist/esm`, ES2015 modules) | down-levelling |
| tslib | 2.8.1 (`tslib.es6.js`) | down-levelling |
| core-js | 3.50.0 (CommonJS, from `stable/index.js`) | down-levelling |
| oracle corpus | `tools/oracle/corpus.txt` | single files |

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

Input paths are always relative: `inputs/…`, `app/…`, `corpus/…` and
`node_modules/…`. Each compile runs in its own directory holding those
symlinks, and outputs (`out.js`, `out.js.map`) are written there. Both
compilers therefore see the same argv, and source maps name the same paths.
`<results>/<task>/jobs.json` records the fully expanded argv; globs are
expanded in sorted order.

## What is measured

- **Java cold:** `java -jar closure-compiler.jar <args>`, N runs (default 5).
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
  each differing artifact, and the times without a ratio.

For the single-file task, every invocation is gated, and the task gets the
counts of each result.

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

Java results stay valid while the jar is unchanged. Once more of the port
lands, re-measure only jsgo against the stored Java reference outputs:

```sh
python3 -I tools/bench/bench.py run --out /tmp/bench --skip-java --jsgo-ref <branch> \
    --report /tmp/bench/report.md
```

These jsgo runs are then not interleaved with Java runs. For final numbers,
re-run everything on a quiet machine and render with `bench.py report
--final`.

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

Wall time per run on the shared 4-CPU development machine (load 5–17),
from the first baseline:

| Task | Java cold, per run | Java warm, per iteration |
|---|---|---|
| typed-app | 17–37 s | 7–15 s |
| bundle | 42–82 s (12.6 MB in, 1.6 GB RSS) | 20–29 s (3 GB RSS) |
| downlevel-whitespace | 7–11 s | 1.2–2.4 s |
| downlevel-simple | 14–20 s | ~4 s |
| checks-only | 18–40 s | 8–9 s |
| single-files | 1,000 JVM launches, 2.5–4.5 s each: 40–75 min | see the report |

All tasks except single-files take about 35 minutes together (N=5, three
trace runs). A full cold run of the single-file task is 1,000 JVM
startups, so 5 runs take several hours. `--runs 1 --tasks single-files`
gives one run, and the report then shows per-invocation median, p10 and p90
from its 1,000 samples. Its cold per-phase runs are off (`"trace_cold":
false` in the task file).
