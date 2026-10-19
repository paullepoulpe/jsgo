# jsgo vs Java Closure Compiler: end-to-end benchmark

> **Preliminary.** Measured on a shared machine with other jobs running (see the load averages), under `nice -n 10`. A quiet, pinned re-run will replace these numbers.

Spec: PLAN.md §6. Harness: `tools/bench/` (README.md there). Values are medians, with the min–max range over the runs in parentheses. Wall and CPU are seconds; CPU is user+system of the process (for Java warm: the JVM's process CPU during the iteration, JIT and GC threads included). RSS is peak resident memory (Java warm: the JVM's peak over its whole life, warm-up included). Speed-ups are Java time ÷ jsgo time (above 1× means jsgo is faster) and only appear when jsgo's output, source map and diagnostics match Java's byte for byte. *load* is the 1-minute load average at the start of each Java cold run.

## Overall

| Task | Java cold wall | Java cold CPU | Java cold RSS | Java warm wall | Java warm CPU | Java warm RSS | jsgo | jsgo wall | jsgo CPU | jsgo RSS | jsgo speed-up vs cold | vs warm | load (1 min) |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| Typed app build | 12.8 s (12.7–13.4) | 25.2 s (25.0–25.7) | 572 MB | 3.83 s (3.62–4.19) | 4.93 s (4.32–5.48) | 919 MB | MATCH | 5.98 s (5.90–5.99) | 7.48 s (7.35–7.53) | 272 MB | 2.15× | 0.64× | 1.08 (0.82–1.32) |
| Large plain-JS bundle | 27.9 s (26.8–28.0) | 51.7 s (49.9–52.2) | 1491 MB | 16.7 s (16.5–17.0) | 20.9 s (20.0–22.3) | 1906 MB | MATCH | 23.3 s (23.3–24.3) | 28.7 s (28.7–29.9) | 678 MB | 1.20× | 0.72× | 1.61 (1.54–1.70) |
| Down-levelling, WHITESPACE_ONLY | 4.38 s (4.35–4.69) | 8.54 s (8.22–9.19) | 334 MB | 0.85 s (0.81–0.94) | 1.34 s (1.17–1.64) | 890 MB | MATCH | 1.49 s (1.49–1.50) | 1.88 s (1.88–1.88) | 168 MB | 2.94× | 0.57× | 1.45 (1.39–1.49) |
| Down-levelling, SIMPLE | 9.18 s (8.78–9.21) | 18.1 s (17.2–18.1) | 378 MB | 2.68 s (2.53–3.02) | 3.24 s (3.00–3.89) | 1092 MB | MATCH | 4.71 s (4.65–4.93) | 5.83 s (5.80–6.06) | 220 MB | 1.95× | 0.57× | 1.86 (1.82–1.96) |
| Checks only | 14.0 s (13.7–14.5) | 27.6 s (26.8–28.4) | 1078 MB | 5.95 s (5.81–6.06) | 7.75 s (7.23–8.25) | 1421 MB | MATCH | 9.04 s (8.83–9.14) | 11.3 s (10.9–11.4) | 710 MB | 1.55× | 0.66× | 1.65 (1.61–1.68) |

## Per task and per phase

Per-phase times come from separate runs with `--tracer_mode=TIMING_ONLY` (milliseconds, median over the trace runs; for multi-invocation tasks, summed over the invocations of a run). Passes are rolled up with `tools/bench/phases.tsv`. *startup* is process start to the tracker's start; *code printing and source maps* is the tracker's report to process exit (the report is printed before the output is generated, so this row holds toSource, the source map, writing the output and exit). For Java warm, the iteration's start and end stand for process start and exit; with several invocations per iteration (the single-file task) startup and printing cannot be told apart inside the JVM, so all time outside the trackers is on the printing row and startup is 0. *untraced* is the tracker's wall time not attributed to any pass. Speed-ups (Java ÷ jsgo) only for MATCH tasks.

### Typed app build (`typed-app`)

A Closure Library-based app (`tools/bench/app/`: a notes editor written for this benchmark, three goog.module/goog.provide files on top of goog.editor, goog.ui, goog.net, goog.storage and the library's own editor demo plugins `goog.demos.editor.*`) compiled with the whole of Closure Library v20230802 (`closure/goog` and `third_party/closure/goog`, tests excluded, 971 files) as candidate inputs; `--dependency_mode=PRUNE` keeps the app's transitive goog.module/goog.provide dependencies.

Mode: ADVANCED, type checking on (`--jscomp_warning=checkTypes`), source map. Invocations per run: 1. Java cold runs: 3; Java warm: 12 warm-up + 10 measured iteration(s) in one JVM; task file `tools/bench/tasks/typed-app.json`.

- Java exit status: [0].
- Java determinism: all cold runs byte-identical.
- Warm driver fidelity: its last iteration reproduces Java cold run 1 byte for byte.
- Java warm iterations (wall s): 12.10*, 5.60*, 4.79*, 4.08*, 3.85*, 4.11*, 4.32*, 3.82*, 4.10*, 3.87*, 4.14*, 3.88*, 4.00, 3.62, 3.62, 3.83, 3.81, 4.19, 3.94, 3.76, 3.88, 3.82 (* = warm-up).
- jsgo: MATCH (stdout, stderr, exit status and output files byte-identical).

(no per-phase data)

### Large plain-JS bundle (`bundle`)

Four plain scripts compiled into one bundle (12.6 MB): three.js r159 `build/three.js` (UMD, the last release with a script build), lodash 4.17.21 `lodash.js`, react-dom 18.3.1 `umd/react-dom.development.js` and TypeScript 5.6.3 `lib/typescript.js` (8.9 MB).

Mode: SIMPLE. Invocations per run: 1. Java cold runs: 3; Java warm: 12 warm-up + 10 measured iteration(s) in one JVM; task file `tools/bench/tasks/bundle.json`.

- Java exit status: [0].
- Java determinism: all cold runs byte-identical.
- Warm driver fidelity: its last iteration reproduces Java cold run 1 byte for byte.
- Java warm iterations (wall s): 26.79*, 18.23*, 16.35*, 17.02*, 16.67*, 16.64*, 16.03*, 16.84*, 17.01*, 17.22*, 16.92*, 17.00*, 16.52, 16.77, 16.50, 16.46, 16.85, 16.65, 16.72, 16.50, 16.75, 17.02 (* = warm-up).
- jsgo: MATCH (stdout, stderr, exit status and output files byte-identical).

(no per-phase data)

### Down-levelling, WHITESPACE_ONLY (`downlevel-whitespace`)

ES2015+ library code down-levelled to ES5: Angular 16.2.12 `fesm2022/core.mjs` (ES2022 module; 16.x is the last major whose core imports nothing but rxjs), the rxjs 7.8.2 ES2015 module tree (`dist/esm`), tslib 2.8.1, and core-js 3.50.0 (CommonJS modules, entry `stable/index.js`). Inputs are passed through a `node_modules/` layout of symlinks with their package.json files; NODE module resolution, CommonJS rewriting, `--dependency_mode=PRUNE` from the two entry points (2,094 candidate files).

Mode: WHITESPACE_ONLY, `--language_out=ECMASCRIPT5`. Invocations per run: 1. Java cold runs: 3; Java warm: 12 warm-up + 10 measured iteration(s) in one JVM; task file `tools/bench/tasks/downlevel-whitespace.json`.

- Java exit status: [0].
- Java determinism: all cold runs byte-identical.
- Warm driver fidelity: its last iteration reproduces Java cold run 1 byte for byte.
- Java warm iterations (wall s): 3.93*, 1.92*, 1.62*, 1.39*, 1.18*, 1.11*, 0.99*, 0.94*, 0.88*, 0.88*, 0.92*, 0.92*, 0.91, 0.85, 0.94, 0.87, 0.85, 0.84, 0.81, 0.84, 0.82, 0.85 (* = warm-up).
- jsgo: MATCH (stdout, stderr, exit status and output files byte-identical).

(no per-phase data)

### Down-levelling, SIMPLE (`downlevel-simple`)

ES2015+ library code down-levelled to ES5: Angular 16.2.12 `fesm2022/core.mjs` (ES2022 module; 16.x is the last major whose core imports nothing but rxjs), the rxjs 7.8.2 ES2015 module tree (`dist/esm`), tslib 2.8.1, and core-js 3.50.0 (CommonJS modules, entry `stable/index.js`). Inputs are passed through a `node_modules/` layout of symlinks with their package.json files; NODE module resolution, CommonJS rewriting, `--dependency_mode=PRUNE` from the two entry points (2,094 candidate files).

Mode: SIMPLE, `--language_out=ECMASCRIPT5`. Invocations per run: 1. Java cold runs: 3; Java warm: 12 warm-up + 10 measured iteration(s) in one JVM; task file `tools/bench/tasks/downlevel-simple.json`.

- Java exit status: [0].
- Java determinism: all cold runs byte-identical.
- Warm driver fidelity: its last iteration reproduces Java cold run 1 byte for byte.
- Java warm iterations (wall s): 8.69*, 5.21*, 3.78*, 3.45*, 3.14*, 3.20*, 3.02*, 3.41*, 3.09*, 2.62*, 2.52*, 2.58*, 2.85, 2.57, 2.60, 2.53, 2.62, 2.75, 2.57, 2.77, 3.02, 2.74 (* = warm-up).
- jsgo: MATCH (stdout, stderr, exit status and output files byte-identical).

(no per-phase data)

### Checks only (`checks-only`)

All of Closure Library v20230802 (`closure/goog` and `third_party/closure/goog`, tests excluded, 971 files), dependency-sorted, checked and type-checked without output.

Mode: `--checks_only`, `--jscomp_warning=checkTypes`. Invocations per run: 1. Java cold runs: 3; Java warm: 12 warm-up + 10 measured iteration(s) in one JVM; task file `tools/bench/tasks/checks-only.json`.

- Java exit status: [0].
- Java determinism: all cold runs byte-identical.
- Warm driver fidelity: its last iteration reproduces Java cold run 1 byte for byte.
- Java warm iterations (wall s): 13.80*, 8.21*, 6.88*, 6.52*, 6.63*, 6.30*, 5.99*, 5.97*, 5.85*, 6.09*, 5.95*, 6.20*, 6.06, 5.98, 5.96, 6.03, 5.91, 5.81, 5.96, 5.82, 5.92, 5.94 (* = warm-up).
- jsgo: MATCH (stdout, stderr, exit status and output files byte-identical).

(no per-phase data)

## Machine

- date: 2026-10-18 12:19 UTC
- cpu: AMD EPYC 9V74 80-Core Processor
- cpus: 3
- memory: 9.7 GiB
- kernel: 6.18.55
- os: Debian GNU/Linux 13 (trixie)
- go: go version go1.27.2 linux/amd64
- jar_sha256: d6864ebf70bca5ba6e237898de6611a389aaf8b933c2549036ced4a684b664b6
- nice: 10
- locale: C.UTF-8
- loadavg_start: [0.82, 0.54, 0.22]
- loadavg_end: [1.68, 1.65, 1.22]
- bench_commit: c9872522ef49054007c7ded25ee33b634c0ea989
- java: openjdk version "21.0.12.1" 2026-08-18 / OpenJDK Runtime Environment (build 21.0.12.1+1-1-deb13u1-Debian) / OpenJDK 64-Bit Server VM (build 21.0.12.1+1-1-deb13u1-Debian, mixed mode, sharing)
- jsgo: source binary /workspace/jsgo/bin/jsgo-go124, sha256 132709913b97b95e5442085161b8fc8ad48e9f00d33ec3717d4325010e45b684
- environment: JAVA_TOOL_OPTIONS, _JAVA_OPTIONS and JDK_JAVA_OPTIONS removed for Java (default JVM settings); GOGC, GOMEMLIMIT, GOMAXPROCS and GODEBUG removed for jsgo (default Go GC).
