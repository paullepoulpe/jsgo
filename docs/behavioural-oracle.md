# The behavioural oracle (PLAN.md §8, phase 1)

`jsgo --fast_mode` may produce output that differs from compat mode's, which
is byte-identical to Java's. The behavioural oracle checks that a fast build
*does the same thing* as the compat build: each test is compiled with compat
and with fast (same flags otherwise), the uncompiled test and both builds run
under node, and their behaviour is compared, compat being the reference.
Phase 1 is the harness, test262's runtime suite, Closure's runtime tests and
the seeded-fault check; it does not yet replace the alpha-isomorphism layer of
`tools/bench/fast_mode_gate.py` (phase 4).

Code: [tools/behave/](../tools/behave/) (`behave.py` the command line,
`core.py` the engine, `t262.py` and `closure_rt.py` the suites, `runner.js`
the node runner); seeded faults: `internal/jscomp/seeded_faults*.go`.

## Running it

```sh
go build -o /tmp/jsgo ./cmd/jsgo
bash tools/behave/fetch.sh            # the suites, at the pins of tools/behave/pins.json
python3 tools/behave/behave.py run --jsgo /tmp/jsgo --out /tmp/behave-pc               # per-change subset
python3 tools/behave/behave.py run --jsgo /tmp/jsgo --out /tmp/behave-full --subset full --java
go build -tags jsgo_faults -o /tmp/jsgo-faults ./cmd/jsgo
python3 tools/behave/behave.py faults --jsgo /tmp/jsgo-faults --baseline /tmp/behave-pc --out /tmp/behave-faults
python3 tools/behave/behave.py explain --jsgo /tmp/jsgo-faults --baseline /tmp/behave-pc --out /tmp/behave-faults
```

- `run` writes `results.jsonl` (one line per test and configuration, with the
  three outcomes), `summary.md` and `summary.json`, the compat builds
  (`builds/`, gzipped), and for every difference a directory under `diffs/`
  with both builds, both outcomes and the command lines. It exits 1 when a
  test that passes on the compat build behaves differently on the fast build
  and is not listed in `tools/behave/expect/compat-fast.txt`. `--resume`
  continues an interrupted run; `--suites`, `--configs`, `--filter` and
  `--limit` narrow it; `--jobs` defaults to the number of CPUs.
- `--java` adds the uncompiled-vs-compat check (below): `findings.jsonl`,
  and `findings/` for port differences.
- `--rerun-identical` also runs fast builds that are byte-identical to the
  compat build (by default their run is reused), which checks that the runner
  is deterministic: identical bytes giving different behaviour is reported as
  `nondeterministic`.
- `faults` rebuilds the fast builds of a finished run once per seeded fault
  and writes `faults.md` and `faults.json`; `explain` then tells, for each
  change no test caught, whether the changed code ran (below).
- `replay` runs a finished run's compat builds again and checks that each
  behaves as recorded (the runner's determinism, or the engine's after a node
  upgrade).
- `probe` and `sample` regenerate the test262 sample (below); after a new
  sample, `behave.py verify --update` records its digest.

Requirements: node 22 (`--node`, default the `node` on `PATH`), Python 3,
and for `--java` a JDK and the pinned jar (`tools/oracle/fetch.sh`).

## What is compared

For each test, configuration (`simple-es5`, `simple-es2015`,
`advanced-es5`, `advanced-es2015`: `-O SIMPLE|ADVANCED --language_in
ECMASCRIPT_NEXT --language_out ECMASCRIPT5|ECMASCRIPT_2015`, both modes at
`--num_parallel_threads 1`, since fast mode's output is the same at every
thread count) and build, the oracle records:

- whether it compiled, and the error codes when it did not;
- whether the run finished within 30 s;
- the exit status, stdout and stderr;
- the runner's record: whether an exception escaped (from the script, a timer
  callback or a microtask) and its error type (the `name` of an Error,
  `Test262Error` for test262's error class, else the constructor's name) and
  message; the unhandled promise rejections at the end; the number of timer
  callbacks run and left; for Closure's tests, goog.testing's per-test results.

The suite's rule turns a run into pass or fail: test262's INTERPRETING.md
(a negative runtime test must throw the named error; an async test must print
`Test262:AsyncTestComplete`; any other must complete), and for Closure's
tests goog.testing's verdict with at least one test run.

Compat versus fast:

| result | meaning | counts against fast mode |
|---|---|---|
| same | the records are equal | |
| same-renamed | equal once the identifiers that a build has and the sources do not (names the compiler made up) are replaced by numbered placeholders | |
| regression | compat passes, fast fails | yes |
| behaviour | both pass, but the records differ | yes |
| fast-only-pass | compat fails, fast passes | reported only |
| weak | both fail, differently | reported only |

A test that fails on the compat build carries no signal (PLAN.md §8): the
summary counts compat passes per configuration, which is how many tests check
something.

Allowed to differ: stack traces (the runner records none, and the V8 frame
lines a program prints are dropped before comparing), `Function.prototype.toString`
(only seen if a test prints a function, which no phase-1 test does), timing
(the clock is virtual), renamed identifiers in messages (normalised as above).
The allowance for the enumeration order of renamed properties is not needed
by phase 1, whose outputs enumerate nothing but goog.testing's results (sorted
by name); it comes with the library suites of phase 2.

**The deterministic prelude** (`runner.js`). The program runs as a classic
script in node's main realm (`vm.runInThisContext`), after: a fixed-seed
`Math.random`; virtual timers (`setTimeout`, `setInterval`, `setImmediate`
and their `clear` functions): after the script, the runner fires them in
(due time, creation order), jumping the clock, and drains the promise jobs
between two timers; `performance.now()` on the same virtual clock; `print`
and test262's `$262` (`global`, `evalScript`,
`detachArrayBuffer`); optionally a frozen `Date` that reads the virtual
clock (on for Closure's tests, off for test262, whose Date tests check the
constructor itself), and a minimal DOM stub (`window`, `document`,
`navigator`, `location`, `window.onload`) for goog.testing. The process runs
with `TZ=UTC`, and the Go and node variables that change behaviour are
cleared.

The first seeded-fault run showed why the prelude has to be complete: two
runs of the same goog.testing build differed in the "N tests run in X ms"
line, because `performance.now()` still read the real clock. Since it reads
the virtual one, `replay` ran the per-change compat builds a second time
(5,447 test262 builds and 829 Closure builds, hours apart) with no
difference. Running in the main realm has one known weakness: a test that
puts properties on built-in prototypes (such as a getter for
`Object.prototype[1]`) can break node's own event loop, and the run hangs;
56 of the 58 test262 tests that time out uncompiled modify built-in
prototypes, and they are out of the sample.

## The suites

### test262

test262 at commit `91c70759` (2026-10-29), pinned with a digest of the
harness and the sampled tests in `tools/behave/pins.json`.

**Compilation.** The harness (`assert.js`, `sta.js`, `doneprintHandle.js` for
async tests, then the test's `includes`) and the test are compiled together,
as separate `--js` inputs, so that ADVANCED renames them consistently, with
`tools/behave/test262/externs.js` declaring the host's `print` and `$262`. A
test without flags runs in strict mode (Closure assumes strict input): the
uncompiled program is `"use strict";` followed by the inputs, and the builds
get `--emit_use_strict`. `noStrict` and `raw` tests run sloppy: the builds
get `--strict_mode_input=false`. `raw` tests get no harness.

**Exclusions** (`tools/behave/test262/exclusions.json`, generated by
`behave.py sample` from a probe). Of the 54,045 files under `test/`:

| excluded | tests | why |
|---|---:|---|
| `negative:parse` | 4,669 | the expected outcome is an early error: a parser test with no runtime behaviour (the parser is checked byte for byte against Java, and test262-parser-tests is in the corpus) |
| `intl402/` | 3,382 | ECMA-402 is a library outside the compiler's scope, and its results depend on ICU data |
| `staging/` | 1,487 | proposals and engine-contributed tests not yet in the main suite |
| `module` flag | 598 | ES module tests need a loader and their `_FIXTURE` imports; Closure rewrites modules into one script, so `import.meta`, live bindings across files, top-level await and namespace objects do not carry over |
| host features | 339 | `cross-realm`, `IsHTMLDDA`, `host-gc-required`, or `atomicsHelper.js` (`$262.agent`): hooks the runner does not provide |
| `_FIXTURE` files | 294 | imported by module tests, not tests |
| `harness/` | 116 | tests of test262's own includes |
| `negative:resolution` | 34 | module resolution errors |
| `CanBlock` flags | 9 | need `$262.agent` worker threads |
| node 22 lacks the feature | 6,408 | 33 features where at least half the tests fail uncompiled: Temporal (4,613), dynamic-import (a loader), explicit-resource-management, ShadowRealm, decorators, source-phase-imports, import-defer, iterator-chunking/-includes/-sequencing, joint-iteration, await-dictionary, upsert, uint8array-base64, Float16Array, RegExp.escape, Math.sumPrecise, Error.isError, tail-call-optimization, … |
| Closure rejects the feature | 4,126 | 9 features where at least 75% of the probed tests do not compile at `simple-es2015`: private class members (5 features, `JSC_UNSUPPORTED_LANGUAGE_FEATURE`), regexp lookbehind and unicode property escapes (`JSC_UNTRANSPILABLE` to ES2015), the regexp `v` flag and import attributes (`JSC_PARSE_ERROR`) |
| fails uncompiled | 500 | the test fails on node 22 without compilation, so it cannot check anything |

That leaves a pool of 32,083 tests. The sample is 5,000 of them
(`tools/behave/test262/sample.txt`), stratified by the first three path
components (493 strata; each gets its proportional share, at least 2, the
lowest sha256(seed, path) first); the per-change subset is a stratified
1,500 of the sample (`per-change.txt`). Features that Closure supports only
partly (classes, BigInt, async iteration, …) stay in: their compat failures
are counted as tests without signal.

### Closure Compiler's runtime tests

`test/com/google/javascript/jscomp/runtime_tests` (Apache 2.0): 163
goog.testing test files covering language features, ES modules, goog.module
and polyfills, written to run compiled. GitHub's mirror of Closure Compiler
carries them only up to tag **v20231112** (they are absent from v20240317 to
v20261021), so the oracle uses that tag, with Closure Library 20230802.0.0
from the corpus for goog.testing. Each test is compiled as the directory's
`utils/build.sh` does (`--dependency_mode PRUNE` with the test as entry
point, `--process_common_js_modules --module_resolution NODE`), with the
files it needs found by a namespace and import scan, and runs with the DOM
stub and a frozen clock; goog.testing's runner starts from `window.onload`.

- ADVANCED renames the test functions that goog.testing looks up by name.
  Closure's own test builds use `CompilerOptions.exportTestFunctions`, which
  the command line does not expose, so each case gets two generated inputs
  (`.cache/behave/closure-gen/`, the same in every configuration): an externs
  file whose object literal lists the test names (property renaming keeps
  them), and for jsunit-style tests an extra entry point calling
  `goog.exportSymbol` on each test function.
- The 97 polyfill tests also run as `#forced` variants (ES5 output only): the
  directory's `force_polyfills.js` deletes the engine's built-ins before the
  build runs, so that the build's polyfills are the code under test.
- Skipped: the two `polyfill_isolation_*` tests, which need their HTML page.
  That leaves 161 tests and 97 forced variants: 258 cases, 838 builds per
  mode over the four configurations.
- There is no uncompiled form (goog.module and ES modules need the compiler).

## Seeded faults

The faults are deliberate fast-mode bugs that exist only in a test binary
built with `go build -tags jsgo_faults -o /tmp/jsgo-faults ./cmd/jsgo`
(`internal/jscomp/seeded_faults*.go`). In that binary
`JSGO_FAST_FAULT=<name>` turns one on, in fast mode only. In every other
build the fault switch is the constant "none": each check is false at
compile time, and the environment is never read. `behave.py faults` and
`explain` refuse a binary on which `fold-negzero` does not fire on
`tools/behave/faults_toy.js`.

| fault | the bug |
|---|---|
| `inline-args` | flipped inlining condition: FunctionArgumentInjector ignores an argument's own side effects, so the argument goes into the inlined body as it is: dropped when the parameter is unused, evaluated twice when it is used twice |
| `drop-unused-init` | dropped side effect: RemoveUnusedCode deletes `var x = f()` without keeping `f()` |
| `swap-args` | swapped evaluation order: FunctionToBlockMutator declares the temporaries of an inlined call's arguments in reverse order |
| `fold-negzero` | PeepholeFoldConstants folds an arithmetic result of `-0` to `0` |
| `hook-negation` | PeepholeMinimizeConditions drops the leading `!` of a `?:` condition without swapping the branches |
| `collapse-unsafe` | CollapseProperties' fast-mode-only early flattening skips its safety check |
| `inline-vars` | InlineVariables moves a variable's value to its only use without checking that nothing in between has side effects |

Two earlier faults were replaced. A first `inline-args` skipped only the
check that the callee's side effects can change an argument; it changed 507
builds, no test caught it, and every change inspected was equivalent on its
input. `swap-operands`, which let PeepholeSubstituteAlternateSyntax rotate
`a * (b / c)` into `b / c * a` without its side-effect check, never changed
a single one of the 6,838 builds: no test has that shape after inlining.

### Results

The faults run over the per-change subset's 6,838 builds (1,500 test262
tests and 258 Closure cases, four configurations), against the per-change
run's compat builds. A build is *changed* when its output differs from the
fault-free build (the compat build: the two are byte-identical, see
"Current numbers");
the others are byte for byte the fault-free build, so nothing could catch
the fault there. *Signal* builds are those whose compat build passes.
`behave.py explain` then rebuilds every changed signal build that no test
caught, locates the changed tokens (renamed identifiers ignored), and runs
the build under V8 block coverage to see whether the changed code ran.

| fault | output identical | output changed | changed, signal | caught | caught / changed signal | uncaught, changed code ran | uncaught, never ran |
|---|---:|---:|---:|---:|---:|---:|---:|
| `inline-args` | 6,377 | 461 | 409 | 0 | 0% | 409 | 0 |
| `drop-unused-init` | 6,270 | 568 | 492 | 37 | 8% | 426 | 29 |
| `swap-args` | 5,505 | 1,333 | 1,119 | 0 | 0% | 1,105 (+14 renaming only) | 0 |
| `fold-negzero` | 6,820 | 18 | 18 | 4 | 22% | 8 | 6 |
| `hook-negation` | 5,621 | 1,217 | 1,116 | 348 | 31% | 505 | 263 |
| `collapse-unsafe` | 6,368 | 470 | 427 | 427 | 100% | 0 | 0 |
| `inline-vars` | 4,033 | 2,805 | 2,365 | 290 | 12% | 1,315 | 760 |
| all | | 6,872 | 5,946 | 1,106 | 19% | | |

**Mutation score: 5 of 7 faults caught** (`inline-args` and `swap-args`
are not), and 1,106 of the 5,946 changed signal builds. 400 of the catches
are `collapse-unsafe` crashing the compiler on Closure Library code, which
the existing strict layer catches too; without them the oracle catches 706
of 5,546 (13%). Caught tests: 17, 1, 122, 267 and 96 of the 1,589 signal
tests for the faults caught. Each fault run takes 18 to 32 minutes.

Of the 4,840 changed signal builds that no test caught, 3,768 ran the
changed code, 1,058 never ran it, and 14 differ only by renaming. A changed
build that no test catches is one of:
- **benign on its input:** the changed code ran, or would have, and did
  the same thing;
- **a blind spot:** the change does alter behaviour, in code that no test
  ran, or that a test ran without checking the result.

Coverage separates "ran" from "never ran", but whether a change that ran
was harmless takes reading it. The split is approximate: a deletion is
placed at the token after it, so it counts as run when that token ran; and
identifiers of up to three characters are treated as renamed, so a change
that only swaps two short names is not seen. Hand checks of samples (hunks under
`<out>/<fault>/explain/`):

- `inline-args`: benign. The arguments that Closure's side-effect analysis
  calls impure but that only build fresh objects are dropped or moved:
  `function f({}) {...}; f(Object.defineProperty({}, 'attr', {get: ...}))`
  in `language/statements/function/dstr/obj-ptrn-empty.js` loses the
  `defineProperty` call on an object nobody keeps; a `new` of a local
  iterator class in goog.Uri moves into its only use. A real double or
  dropped evaluation (`once(g())` loses `g()` in the fault's own toy test)
  never occurs in an inlined call of these suites.
- `drop-unused-init`: benign where it ran. The dropped initializers are calls
  with no actual effect: `var TypedArray = Object.getPrototypeOf(Int8Array)`
  in test262's `testTypedArray.js`, Closure Library's unused module constants
  (`new SafeUrl("about:invalid#zClosurez", token)`). The catches are real
  effects: a `for (let [...] = ...)` initializer whose calls fill the array
  the test expects (`array_pattern_test`), destructurings in test262 whose
  iterator or getter must run, a rejected promise that disappears
  (`async_function_test`).
- `swap-args`: benign. The reversed temporaries are pure in practice: the
  message literal of an inlined `assert.sameValue(actual, expected, message)`
  now comes before `actual`; two field reads in goog.Promise swap. No inlined
  call in these suites has two arguments whose side effects interfere.
- `fold-negzero`: benign where it ran: `-0` as a typed-array index (ToIndex
  makes it `0`) in the Atomics tests, and `1/-Infinity` compared with `===`
  in `assert._isSameValue`. The changes that never ran are in a harness
  helper the test does not call (`testWithAtomicsInBoundsIndices`).
- `hook-negation`: both kinds. Benign: `makeArray` in `testTypedArray.js`
  returns `4` instead of `["0","0","0","0"]`, and `new TA(4)` is the same
  array. Blind spot, ran but unchecked: goog.testing's `safeSetUpHelper_`
  now skips `setUp`, which only the tests with a real `setUp` notice (those
  are among the catches); every other Closure test passes the same.
  Blind spot, never ran: 263 test262 builds, in harness code paths that only
  failing tests take or helpers that the test does not call.
- `inline-vars`: both kinds. Benign: values moved past code that does not
  touch them (`Object.getPrototypeOf(...)` chains, propertyHelper's option
  reads, goog.log creating a parent logger after registering the child).
  Blind spot, never ran: the Promise polyfill's `then` in an ES5 build
  creates `new Promise(function(res, rej) {k = res; m = rej})` after passing
  `k` and `m` on, so `then` is broken, but node's own Promise is used and
  the polyfill never runs (the `#forced` variants, where it does, catch it).
- `collapse-unsafe`: everything is caught, 400 builds by a compiler crash and
  27 test262 builds by a TypeError.

**What this means.** The oracle is sharp where code runs and its result is
checked: a fault that changes such code is caught. But on test262 and
Closure's tests, most of the code a fault changes either does the same
thing on these inputs or is never run (harness failure paths, unused helpers,
library code behind a feature test), and some of it runs without anyone
checking (goog.testing's own setUp path). So the oracle cannot replace the
alpha-isomorphism layer on its own: a fast-mode change in code that no test
runs would pass unseen. It can replace it under PLAN.md §8's coverage rule,
where a change outside alpha-isomorphism passes only if the changed code ran
and behaved the same, and a change in code that did not run is a blind spot
that keeps the structural check. `explain` already computes the "ran" half
for any two builds; phase 4 needs jsgo to log where fast-only
transformations applied. Phase 2 has to bring tests that run library code
broadly and assert on its results, so that the changed code is both run and
checked; the blind spots above are the measure of what it must cover.

## Current numbers

Measured on 2026-11-01 with jsgo `ca96b0e0` (before the fix below), node
22.22.0, on the 4-CPU machine with 4 jobs. Both builds run at
`--num_parallel_threads 1`: fast mode's output is the same at every thread
count (the strict layer of `fast_mode_gate.py`).

Full run (`--subset full`):

| suite / configuration | tests | pass uncompiled | compat compiles | compat passes (signal) | fast passes | fast bytes = compat | differences |
|---|---:|---:|---:|---:|---:|---:|---:|
| test262 / simple-es5 | 5,000 | 5,000 | 4,670 | 4,224 | 4,224 | 4,670 of 4,670 | 0 |
| test262 / simple-es2015 | 5,000 | 5,000 | 4,742 | 4,328 | 4,328 | 4,742 of 4,742 | 0 |
| test262 / advanced-es5 | 5,000 | 5,000 | 4,307 | 3,351 | 3,351 | 4,307 of 4,307 | 0 |
| test262 / advanced-es2015 | 5,000 | 5,000 | 4,280 | 3,290 | 3,290 | 4,280 of 4,280 | 0 |
| closure / simple-es5 | 258 | n/a | 256 | 249 | 249 | 256 of 256 | 0 |
| closure / simple-es2015 | 161 | n/a | 159 | 153 | 153 | 159 of 159 | 0 |
| closure / advanced-es5 | 258 | n/a | 256 | 247 | 247 | 256 of 256 | 0 |
| closure / advanced-es2015 | 161 | n/a | 158 | 153 | 153 | 158 of 158 | 0 |

That is 15,193 test262 builds and 802 Closure builds that check something.
The per-change subset has 1,299, 1,326, 1,080 and 1,075 signal builds of
1,500 tests in the four configurations, plus the whole Closure suite.

**Compat vs fast: no difference.** Every fast build that compiled is byte
for byte the compat build: fast mode's shortcuts today (InlineObjectLiterals'
per-function temporaries, CollapseProperties' early flattening, work moved to
workers) leave no trace on programs of this size. The oracle's power on fast
mode is therefore measured by the seeded faults, not by these runs.

Closure's compat failures (no signal): two module tests import files that
GitHub's mirror lacks (`module_test_resources/+subdir/`, `a.js`),
`error_cause_test` and the non-forced `unhandled_rejection` tests skip
themselves under the stub's user agent and events, and a few tests fail on
Closure's own output (async generator delegation, `class extends Error` and
static `super` at ES2015, forced `TypedArray`/`startsWith`/`endsWith`
polyfills at ES5).

## Uncompiled passes, compat fails

With `--java`, every build whose uncompiled test passes and whose compat build
fails is compiled again by the pinned Java compiler, with the same command
line (`tools/oracle/java/CliBatch.java`, one warm JVM per job slot), and the
outputs are compared byte for byte (or the error codes, when neither
compiles).

Full run: **4,807 builds** (776, 672, 1,649 and 1,710 in the four
configurations) on 1,801 tests. **4,712 are upstream behaviour** (Java's
output is the same) and **95 were a jsgo port bug**, now fixed.

**The port bug** (95 builds, all `advanced-es2015`): jsgo stopped with an
internal compiler error (a nil dereference reported as Java's
NullPointerException) where Java compiles; Java's build passes the test in
58 of them. Minimal repro:

```js
var C = class { m() { window['n'] = arguments.length; } };
C.prototype.m(1, 2);
```

`jsgo -O ADVANCED --language_out ECMASCRIPT_2015 --js repro.js` failed;
Java prints `(class{h(){window.n=arguments.length}}).prototype.h(1,2);`.
Once the class expression is inlined into `C.prototype`, it has no name, and
`AnalyzePrototypeProperties.processMemberDef` checks
`classNameNode != null && classNameNode.isName()`; the port had dropped the
null check. The fix restores it (with a regression test in
`cross_chunk_method_motion_test.go`); after it all 95 builds are byte-identical
to Java's. Outputs that did not crash are unchanged.

**Upstream behaviour** (4,712 builds), by the compat build's failure:

| failure | simple-es5 | simple-es2015 | advanced-es5 | advanced-es2015 |
|---|---:|---:|---:|---:|
| throws Test262Error (a wrong value) | 323 | 275 | 645 | 660 |
| `JSC_UNTRANSPILABLE` (cannot lower the feature to the output language) | 184 | 127 | 182 | 125 |
| `JSC_UNDEFINED_VARIABLE` (globals created by `eval`, `this.x` or the harness at run time) | 0 | 0 | 269 | 269 |
| async test fails (transpiled async generators and for-await) | 89 | 121 | 93 | 124 |
| throws TypeError | 19 | 6 | 155 | 144 |
| `JSC_PARSE_ERROR` (sloppy-mode and annex B syntax, HTML comments, …) | 63 | 63 | 63 | 63 |
| `JSC_VAR_MULTIPLY_DECLARED_ERROR` (a test redeclares a harness name) | 0 | 0 | 98 | 98 |
| `JSC_DECLARATION_NOT_DIRECTLY_IN_BLOCK` (annex B function declarations) | 31 | 31 | 31 | 31 |
| throws ReferenceError | 5 | 5 | 49 | 51 |
| internal compiler error in Java too | 4 | 1 | 4 | 4 |
| other compile errors (12 codes) | 48 | 36 | 46 | 35 |
| other runtime errors (strings, SyntaxError, RangeError) | 10 | 7 | 14 | 11 |

What they are, from samples:
- SIMPLE: Closure assumes code does not depend on the errors that dead or
  unused code would throw, so it removes destructuring that should throw
  (`null` patterns, unresolvable defaults, iterator steps and closes),
  temporal-dead-zone reads and assignments to unresolvable references; it
  loses the `name` that function-name inference gives arrows and functions in
  destructuring defaults; and its transpiled async generators and
  `yield*` deviate from the spec on abrupt completions.
- ADVANCED adds renaming: properties read by string (`verifyProperty`,
  `hasOwnProperty`, descriptors), built-ins missing from the default externs
  (resizable ArrayBuffer, `Array.fromAsync`, …), function names, and globals
  reached through `this`, `eval` or `globalThis` (`$DONE` for `asyncHelpers.js`).
- Both: syntax Closure rejects or cannot lower to ES5 or ES2015 (BigInt
  literals, regexp flags and groups, sloppy-mode-only forms).

Everything in this table is Closure's behaviour, the same in Java; it costs
signal (13–16% of the SIMPLE builds, 32–33% of the ADVANCED ones), not
correctness.

## Wall time

On the 4-CPU machine, `--jobs 4`:

| run | what | wall time |
|---|---|---:|
| per-change (`run`, the default) | 1,500 test262 tests and the 258 Closure cases, 4 configurations | **49 min** (test262 36.1, Closure 13.1) |
| full (`run --subset full --java`) | 5,000 test262 tests and the Closure cases, then the Java check | **2 h 33 min** (test262 124.4, Closure 13.1, Java check of 4,807 builds 15.2) |
| `faults` | each fault rebuilds the per-change fast builds | 18–32 min per fault |
| `probe` (when the test262 pin moves) | 43,117 tests uncompiled, 1,324 compiles | 21 min |

Compiling dominates: a test262 compile takes about 0.27 s at SIMPLE and
0.7 s at ADVANCED, most of it parsing and type-inferring the default externs
(80,000 lines; `inferTypes` alone is 0.22 s of an ADVANCED compile), against
0.1 s for a node run. A fast build that is byte-identical to the compat build
is not run again (its behaviour is the compat run's); `--rerun-identical`
runs it anyway. The full run is a nightly job; the per-change subset fits the
30–60 minute budget of PLAN.md §8 on its own, which leaves no room for the
later phases without a cheaper compile (see PLAN.md §8, "Phase 1").

## Adding a suite (phases 2–4)

A suite is a function returning `core.Case` objects: the directory the
compiler runs in, the `--js` inputs in order, extra flags, the uncompiled
program (or none), runner options, and a verdict function from the runner's
record and stdout to pass/fail. The engine does the rest (compiling, running,
comparing, the Java check, faults, reports). A library suite (phase 2) is a
Case per test file, compiled with the library in ADVANCED, with a verdict
read from the test framework's output, as `closure_rt.py` does with
goog.testing; a workload check (phase 3) is a Case whose program prints its
results, compared through stdout. The coverage rule (phase 4) needs two
additions: jsgo logging where fast-only transformations applied, and the
runner's `NODE_V8_COVERAGE` output mapped through the builds' source maps.

## Moving the pins

Change the revision in `tools/behave/pins.json`, run `fetch.sh` (it refetches
a checkout whose commit differs), regenerate the sample if test262 moved
(`behave.py probe`, then `behave.py sample --probe DIR`), and record the new
digests with `behave.py verify --update`.
