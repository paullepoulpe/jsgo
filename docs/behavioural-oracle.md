# The behavioural oracle (PLAN.md §8, phases 1 and 2)

`jsgo --fast_mode` may produce output that differs from compat mode's, which
is byte-identical to Java's. The behavioural oracle checks that a fast build
*does the same thing* as the compat build: each test is compiled with compat
and with fast (same flags otherwise), the uncompiled test and both builds run
under node, and their behaviour is compared, compat being the reference.

- Phase 1: the harness, test262's runtime suite, Closure's runtime tests and
  the seeded-fault check.
- Phase 2: the coverage rule (a fast build that is not alpha-isomorphic to
  compat passes only where its changes ran and behaved the same, and the
  changes that never ran are reported as blind spots); library suites that run
  library code broadly and check its results (Closure Library's unit tests,
  lodash's test suite, marked's spec tests); one bench-scale workload brought
  forward from phase 3 (the TypeScript compiler compiling a fixed program);
  benign seeded variants and the opt-in loop cap as "different but correct"
  inputs; a compat cache and a smaller test262 subset for the per-change
  budget.

The oracle does not yet replace the alpha-isomorphism layer of
`tools/bench/fast_mode_gate.py`; that is phase 4. Today fast output is byte
for byte compat's on every build of every suite (and on every task of the
fast-mode gate), so the oracle's power is measured with seeded faults, benign
variants and the loop cap, not with fast mode itself.

**What does not work yet** (details below):
- Almost no build is vouched outright: nearly every build carries some
  changed code its test never runs (polyfills behind feature tests, error
  paths, unused helpers). Under the loop cap 440 of 3,098 changed builds are
  vouched, while 57% of the distinct changes are. The gate rule that would
  act on distinct changes is phase 4's; it is not written.
- "Ran" is per process, not per test, and "vouched" means the change ran and
  nothing the tests check differed. The seeded faults `inline-args` and
  `swap-args` are still never caught; their instances are harmless where
  they occur, and the rule vouches for them as for correct changes.
- The TypeScript workload runs a 150-line program; it runs 35% of the code
  the loop cap changes in the compiler.
- The faults were re-run under the coverage rule on the library suites only,
  in two of the four configurations; the phase-1 suites' fault numbers are
  phase 1's.
- The enumeration-order allowance for renamed properties is still not
  implemented (no suite needs it), and failure messages that print compiled
  code make some builds differ inside tests that already fail (they stay
  blind, never fail).

Code: [tools/behave/](../tools/behave/) (`behave.py` the command line,
`core.py` the engine, `rule.py` the coverage rule, `t262.py`, `closure_rt.py`,
`closure_lib.py`, `libs.py` and `workloads.py` the suites, `runner.js` the
node runner); the structural diff:
[tools/bench/ast_isomorphism/diff.go](../tools/bench/ast_isomorphism/diff.go);
seeded faults: `internal/jscomp/seeded_faults*.go`.

## Running it

```sh
go build -o /tmp/jsgo ./cmd/jsgo
bash tools/behave/fetch.sh            # the suites, at the pins of tools/behave/pins.json
bash tools/bench/fetch-inputs.sh      # TypeScript 5.6.3 for the tsc workload (bench inputs)
python3 tools/behave/behave.py run --jsgo /tmp/jsgo --out /tmp/behave-pc               # per-change subset
python3 tools/behave/behave.py run --jsgo /tmp/jsgo --out /tmp/behave-full --subset full --java
python3 tools/behave/behave.py run --jsgo /tmp/jsgo --out /tmp/behave-cap \
    --fast-flags=--max_optimization_loop_iterations=1                                 # a "different but correct" input
go build -tags jsgo_faults -o /tmp/jsgo-faults ./cmd/jsgo
python3 tools/behave/behave.py faults --jsgo /tmp/jsgo-faults --baseline /tmp/behave-pc --out /tmp/behave-faults
python3 tools/behave/behave.py faults --jsgo /tmp/jsgo-faults --baseline /tmp/behave-pc --out /tmp/behave-benign \
    --faults benign-rename,benign-no-inline
```

- `run` writes `results.jsonl` (one line per test and configuration, with the
  three outcomes and the coverage rule's verdict), `summary.md` and
  `summary.json`, the compat builds (`builds/`, gzipped), for every behaviour
  difference a directory under `diffs/` with both builds, both outcomes and
  the command lines, and for every build that is not byte-identical the
  located changes and their verdicts under `rule/`. It exits 1 when a test
  that passes on the compat build behaves differently on the fast build and is
  not listed in `tools/behave/expect/compat-fast.txt`. `--resume` continues an
  interrupted run; `--suites` (test262, closure, closure-lib, lodash, marked,
  tsc), `--configs`, `--filter` and `--limit` narrow it; `--jobs` defaults to
  the number of CPUs.
- `--fast-flags` adds compiler flags to the fast builds only. It exists to
  feed the oracle output changes that keep behaviour, such as the opt-in loop
  cap; it never changes jsgo's defaults.
- `--compat-cache DIR` (default `<cache>/behave/compat-cache`; `none` turns
  it off): compat builds, their compile results and their runs are kept
  across runs (below, "Wall time").
- `--java` adds the uncompiled-vs-compat check (below): `findings.jsonl`,
  and `findings/` for port differences.
- `--rerun-identical` also runs fast builds that are byte-identical to the
  compat build (by default their run is reused), which checks that the runner
  is deterministic: identical bytes giving different behaviour is reported as
  `nondeterministic`.
- `faults` rebuilds the fast builds of a finished run once per seeded fault
  (or benign variant) and writes `faults.md` and `faults.json`, with the
  coverage rule's verdicts; `--suites` and `--configs` narrow it to part of
  the baseline. `explain` (phase 1) tells, for each change no test caught,
  whether the changed code ran; the coverage rule now does this for every
  differing build.
- `replay` runs a finished run's compat builds again and checks that each
  behaves as recorded (the runner's determinism, or the engine's after a node
  upgrade).
- `probe`, `sample` and `subset` regenerate the test262 sample and per-change
  subset (below); after a new sample or a moved pin, `behave.py verify
  --update` records the digests.

Requirements: node 22 (`--node`, default the `node` on `PATH`), Python 3, Go
(the coverage rule builds `tools/bench/ast_isomorphism` into the cache), and
for `--java` a JDK and the pinned jar (`tools/oracle/fetch.sh`).

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
  callbacks run and left; for goog.testing and the QUnit suites, the per-test
  results (each test's failure messages).

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
something. Where the runner reports per-test results (goog.testing, QUnit),
signal is per test: a test that passes on the compat build and not on the
fast build is a `regression` even when another test of the same file fails on
both, and a build "checks something" when at least one of its tests passes
on compat. (Phase 1 counted whole files: a Closure test file with one failing
test carried no signal at all.)

Allowed to differ: stack traces (the runner records none, and the V8 frame
lines a program prints are dropped before comparing), `Function.prototype.toString`
(only seen if a test prints a function, which no phase-1 test does), timing
(the clock is virtual), renamed identifiers in messages (normalised as above;
since phase 2 every identifier of one or two characters counts as renamed,
because compiled code makes them even when a source happens to use the same
name). The allowance for the enumeration order of renamed properties is still
not implemented: the library suites' outputs are test results keyed by test
name, digests of strings, and assertions whose failure messages name the
assertion, not the values compared. A program that prints a name the compiler
chose cannot be compared through a digest: the tsc workload first printed its
source maps' digest, and ADVANCED had renamed the key "mappings" differently
in the compat and the loop-capped builds (`FG`, `DG`), which the digest hid
from the normalisation; it now digests the maps' values.

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
`navigator`, `location`, `window.onload`) for goog.testing. Phase 2 adds:
`--env qunit`, a QUnit subset (`QUnit.module`, `test`, `skip`, `config`,
`load`, `start`, `done` and the assertions of QUnit 1 and 2) that lives in
the runner, uncompiled, so that the checking code is not part of the build
under test (`tools/behave/libs/qunit_externs.js` declares it to the
compiler); tests run one after the other on the virtual clock, and an
exception in a timer while a test runs fails that test; and `--date-step N`,
which moves the virtual clock on by N ms at each read of the frozen `Date`,
so that code which waits for the clock in a loop (lodash's throttle tests)
ends. The process runs with `TZ=UTC` and `--no-warnings` (node's process
warnings carry the process id), and the Go and node variables that change
behaviour are cleared.

The first seeded-fault run showed why the prelude has to be complete: two
runs of the same goog.testing build differed in the "N tests run in X ms"
line, because `performance.now()` still read the real clock. Since it reads
the virtual one, `replay` ran the per-change compat builds a second time
(5,447 test262 builds and 829 Closure builds, hours apart) with no
difference. Running in the main realm has one known weakness: a test that
puts properties on built-in prototypes (such as a getter for
`Object.prototype[1]`) can break node's own event loop, and the run hangs;
56 of the 58 test262 tests that time out uncompiled modify built-in
prototypes, and they are out of the sample. Phase 1 suggested a fresh `vm`
context per test file for the library suites; it was not needed: every case
runs in its own node process, and the library tests chosen do not patch
built-ins in a way that reaches node (lodash's test file patches them only in
a section it skips without `require`).

## The coverage rule

A fast build that is not byte-identical to its compat build passes only where
the oracle can vouch for it (PLAN.md §8). `run` and `faults` apply the rule
to every such build ([tools/behave/rule.py](../tools/behave/rule.py)):

1. **Locate.** `tools/bench/ast_isomorphism -diff compat.js fast.js` (the
   gate's isomorphism checker, extended in
   [diff.go](../tools/bench/ast_isomorphism/diff.go)) parses both builds with
   jsgo's parser, walks the two ASTs as the layer-2 check does, but aligns
   them instead of stopping at the first difference: where two matched nodes
   have children that differ in number or shape, the children are aligned by
   a longest common subsequence of name-insensitive subtree hashes, and the
   ones left between two aligned children are paired by token (and compared
   further down) or reported as inserted or deleted. Names are compared after
   the walk, up to one renaming: a variable of the compat build corresponds to
   the variable of the fast build its references meet most often (matched
   declarations break ties), both declared in matched scopes; an undeclared
   name (`window`) only to itself; with ADVANCED, an unquoted property name to
   the name it meets most often, one to one, except the externs' property
   names, which only correspond to themselves; at SIMPLE every property name
   only to itself. A reference whose names do not correspond is a change.
   Every node of the fast build is then either matched to a node of the
   compat build with the same token, flags and literal, or inside a reported
   change; a function with no change in it (outside its nested functions) is
   still alpha-isomorphic to its counterpart. Each change comes with its
   offsets in both builds and its unit: the innermost function around it, or
   the top-level statement outside every function. Node columns stop at 4,095
   in jsgo's AST and compiled lines are longer, so the tool parses a copy with
   newlines after some commas and semicolons and uses it only if its AST is
   the same (it was, on every build so far).
2. **Run.** Both builds run under V8 block coverage (`NODE_V8_COVERAGE`). A
   change ran when a point of its range in the fast build has a non-zero
   count (the count of the innermost coverage range around it); a deletion,
   when its range in the compat build does.
3. **Judge**, per change:

   | verdict | when |
   |---|---|
   | vouched | the change ran, the build checks something (its compat build passes, or passes at least one of its tests) and the fast build behaved the same (`same` or `same-renamed`) |
   | blind | the change never ran; or the build checks nothing; or the behaviour differs only where nothing is checked (`weak`, `fast-only-pass`). Reported, never passed |
   | fail | the fast build behaved differently on something checked (`regression`, `behaviour`) and the change ran |

   A build is *identical*, *isomorphic* (no change), *vouched* (every change
   vouched), *blind* (some change blind, none failing) or *fail* (some change
   failing, or a behaviour difference on something checked even when no
   change shows as run: a run killed by the time limit writes no coverage). A build
   that does not compile on one side is *fail* (fast fails where compat
   compiles, or other errors) or *blind* (fast compiles where compat does
   not); the same errors on both sides are no difference.

`run` writes the verdicts into `results.jsonl` (`rule`), the changes of each
differing build into `rule/<suite>/<config>/<id>.json` with the fast build,
and the table and the distinct changes into `summary.md`. A library is
compiled into many builds, so the summary also counts *distinct changes*:
the same function changed the same way in several builds (same normalised
text of its unit in both builds, same kind, same changed text) counts once,
vouched if some build vouched it and none failed it.

**Why a structural diff, not a log from jsgo.** PLAN.md proposed that jsgo
log where each fast-only transformation applied. The diff of the outputs
does the same job, and better:
- it sees the knock-on effects: a different decision early (one more
  inlining, one optimisation round fewer) changes what later shared passes
  do, far from the site a log would name. The loop cap is the extreme case:
  it applies nowhere in particular and changes code everywhere;
- it does not trust the code under test to report its own changes; a log
  that misses a site would pass it unseen;
- its positions are those of the output that runs, which the coverage reads,
  whereas a log's positions would have to be carried through later passes;
- it needs no change to jsgo and works on any two builds.
It costs a parse of both builds per differing build: 0.1 to 0.3 s for test
and library builds, a few seconds for a bundle.

**Limits.**
- Coverage is per run, not per test: a change that ran in a test whose result
  nobody checks counts as vouched (goog.testing's own setUp path in phase 1
  is the known example). The library suites assert on almost every result,
  which is why they matter; the per-test results and the printed outputs
  narrow it further.
- ADVANCED's property ambiguation can make a pure renaming non-bijective (two
  properties that share a name in one build have two names in the other).
  The diff reports those as changes, which the tests then vouch for or leave
  blind; it never calls them isomorphic. The `benign-rename` variant measures
  this (below).
- A change's range is tested at up to 17 points; a change inside a larger
  expression counts as run when any of them ran.

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
lowest sha256(seed, path) first); the per-change subset (`per-change.txt`)
was a stratified 1,500 of the sample in phase 1 and is 600 since phase 2
(`behave.py subset`: the same stratification, at least one test per stratum,
446 strata), to make room for the library suites. Features that Closure supports only
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


### Closure Library's unit tests (phase 2)

Closure Library is the library written for ADVANCED, and its goog.testing
unit tests assert on the results of nearly every module. The npm package in
the corpus leaves the tests out, so `fetch.sh` checks out GitHub's
google/closure-library at the tag of the corpus version, v20230802
(`closure/goog` only; its sources are the corpus package's, file for file).
Each `*_test.js` is a case compiled as Closure's runtime tests are
(dependency pruning from the test, the generated test-name externs and
exports for ADVANCED) and run under the DOM stub with a frozen clock.

The pinned list (`tools/behave/closure_lib/tests.txt`) is the 302 test files
at least one of whose tests passes on the compat build in all four
configurations, out of 443 probed (`run --suites closure-lib --subset
probe`, which leaves out `ui/`, `editor/`, `demos/`, `graphics/`, `fx/`,
`positioning/`, `style/`, `a11y/`, `cssom/`, `history/` and `window/`: they
need a real DOM in nearly every test). They cover `labs`, `testing`, `math`,
`vec`, `structs`, `net`, `html`, `events`, `dom`, `i18n`, `async`, `crypt`,
`debug`, `storage`, `proto2`, `string` and 25 smaller directories: 4,307
passing tests at SIMPLE and 3,834 at ADVANCED (of 4,651 and 4,401). All 302
run per change. 46 more files have signal in some configurations only, 44
of them at SIMPLE and not at ADVANCED (ADVANCED renames a name that their
test or goog.testing reaches by string, or drops code they need); the 88
ADVANCED builds of those were checked against Java: 82 are Java's bytes and 4 fail with Java's error (`JSC_MSG_ORPHANED_NODE`), upstream
behaviour; 2 were a port bug, now fixed (below, "Uncompiled passes, compat
fails"). `net/channelrequest_test.js`, one of the files the bug kept out,
compiles now and can join the list at its next regeneration.

### lodash (phase 2)

lodash 4.17.21 (MIT; GitHub's tag, since npm leaves `test/` out):
`lodash.js` and `test/test.js` in one case, 3,095 QUnit tests over 307
modules (the file generates tests in loops), of which 3,084 pass on the
compat build in every configuration, as uncompiled (the other 11 need node's
`Buffer` and `vm` realms the setup does not provide). The runner provides
QUnit. `test.js` reads a second, "stable" lodash for its helpers (lodash's own
setup loads it from `node_modules`): the runner's `--pre` script loads the
same `lodash.js` uncompiled for that, so the helpers are not the code under
test, fixes `process.argv` (the test file prints it) and sets the `ui` object
the test file would otherwise work out with `require`. The clock is frozen
with `--date-step 1`.

For ADVANCED, lodash and its tests reach methods by name (`_[methodName]`,
lists of method names), and the test calls lodash's API on the uncompiled
stable copy too: a generated externs file
(`.cache/behave/lib-gen/lodash/string_names_externs.js`) lists every
identifier-like string literal of the library and its tests, every property
they read from the global object, lodash's public API, and the keys of object
literals read back with `Object.entries`/`keys`, so that property renaming
leaves those names alone. Everything else ADVANCED does (inlining, dead code
removal, devirtualisation, renaming of the internals) still applies.

### marked (phase 2)

marked 15.0.12 (MIT): the corpus package's `lib/marked.umd.js` compiled with
a driver (`tools/behave/libs/marked_driver.js`) and the spec tests of
`test/specs/new`, `original` and `redos` at the tag v15.0.12: 146 specs,
markdown in and HTML expected, each a QUnit test that compares the output
with the expected HTML (whitespace next to tags ignored; marked's own runner
uses an HTML differ) and prints a digest of every output, so that the builds'
outputs are compared in full, the failing specs' included. 127 pass at
SIMPLE, as uncompiled; 123 at ADVANCED, whose builds are Java's bytes
(upstream behaviour). The CommonMark and GFM example files are left out:
they are the specifications' text, under CC-BY-SA, not the libraries'
licences.

### The TypeScript compiler (phase 3, brought forward)

TypeScript 5.6.3's `lib/typescript.js` (Apache 2.0; 8.9 MB, the bench
bundles' input, `tools/bench/inputs.txt`) compiled with a driver
(`tools/behave/workloads/tsc/driver.js`) that type-checks and emits a small
fixed program in memory (`workloads/tsc/*.ts`, three files written for the
oracle, with four deliberate type errors) against TypeScript's own lib
`.d.ts` files, twice (ES2015 with declarations and source maps; ES5 CommonJS),
and prints every diagnostic and a digest of every emitted file. The scanner,
parser, binder, checker, the declaration emitter and the ES5 transforms all
run. The verdict is that the run completes; the comparison is the full
output. It runs in all four configurations (ADVANCED with generated name
externs as for lodash, and the bench task's `--jscomp_warning=checkVars`); a
build takes 8 s (SIMPLE) to 13 s (ADVANCED) to compile and 1.5 s to run. The
ADVANCED build is Java's bytes. This is the bench task `bundle-advanced-ts`
narrowed to TypeScript, since the bundle's other inputs have no driver.

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
| `benign-rename` (phase 2) | not a bug: RenameVars and RenameProperties skip the first name they would assign, so every renamed name shifts (a pure renaming, which the diff must call alpha-isomorphic) |
| `benign-no-inline` (phase 2) | not a bug: InlineFunctions inlines nothing (a different, correct decision at every call it would have inlined) |

Two earlier faults were replaced. A first `inline-args` skipped only the
check that the callee's side effects can change an argument; it changed 507
builds, no test caught it, and every change inspected was equivalent on its
input. `swap-operands`, which let PeepholeSubstituteAlternateSyntax rotate
`a * (b / c)` into `b / c * a` without its side-effect check, never changed
a single one of the 6,838 builds: no test has that shape after inlining.

### Phase 1 results (test262 and Closure's runtime tests)

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

**What this meant for phase 2.** The oracle is sharp where code runs and its result is
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


### Phase 2: the library suites, under the coverage rule

The faults and the benign variants rebuilt the fast builds of the library
suites (the 302 Closure Library tests, lodash, marked and tsc) in two
configurations, SIMPLE ES5 and ADVANCED ES2015 (610 builds; `faults --suites
closure-lib,lodash,marked,tsc --configs simple-es5,advanced-es2015`), against
the per-change run's compat builds. Every one of these builds checks
something (per-test signal), so "changed" builds are all signal builds.

| fault | changed builds | caught builds | caught test files | rule, builds: vouched / blind / fail | changes vouched / blind / fail | distinct changes vouched / blind / fail | min |
|---|---:|---:|---:|---|---|---|---:|
| `inline-args` | 306 | 0 | 0 | 2 / 304 / 0 | 2,228 / 1,783 / 0 | 892 / 975 / 0 | 6.0 |
| `drop-unused-init` | 313 | 9 | 6 | 23 / 281 / 9 | 4,257 / 2,601 / 236 | 305 / 161 / 177 | 6.2 |
| `swap-args` | 315 | 0 | 0 | 17 / 298 / 0 | 4,042 / 6,680 / 0 | 1,034 / 1,879 / 0 | 5.4 |
| `fold-negzero` | 0 | 0 | 0 | | | | 5.2 |
| `hook-negation` | 610 | 380 | 207 | 0 / 230 / 380 | 474 / 26,436 / 1,358 | 3 / 2,189 / 224 | 8.6 |
| `collapse-unsafe` | 305 | 305 | 305 | 0 / 0 / 305 | (302 compiler crashes) | | 3.4 |
| `inline-vars` | 610 | 143 | 74 | 0 / 467 / 143 | 16,316 / 174,276 / 16,412 | 4,904 / 23,596 / 10,879 | 6.1 |
| `benign-rename` | 610 | 0 | 0 | 607 isomorphic; 1 / 2 / 0 | 194 / 181 / 0 | 194 / 181 / 0 | 4.8 |
| `benign-no-inline` | 610 | 0 | 0 | 0 / 610 / 0 | 100,562 / 108,432 / 0 | 19,432 / 16,748 / 0 | 5.4 |

**The library suites catch what changes results.** `hook-negation` is
caught in 380 of its 610 changed builds (62%; 31% on phase 1's suites),
`inline-vars` in 143 of 610 (23%; 12%), `collapse-unsafe` in all 305 (302 by
crashing the compiler on Closure Library code, as in phase 1, and now also at
run time in lodash, marked and tsc), `drop-unused-init` in 9 of 313. lodash,
marked and tsc catch `hook-negation`, `inline-vars` and `collapse-unsafe` in
every build those change, and none of `inline-args`, `swap-args` or the one
tsc build `drop-unused-init` changes. `fold-negzero` changes none of these builds (no arithmetic that
folds to -0).

**`inline-args` and `swap-args` are still never caught**, now on 621 more
changed builds. The vouched instances read by hand are harmless where they
are: the faulty inlining moves an argument past code it does not interact
with. In lodash's `LazyWrapper.prototype.value`, `swap-args` reads
`this.__dir__` before `this.__wrapped__.value()` instead of after, and
`value()` does not touch `__dir__`; in `_.update`, `inline-args` evaluates
`castFunction(updater)` after `baseGet(object, path)`, both pure; in
goog.testing's TestRunner, `inline-args` moves a `goog.bind(...)` call past a
property read. Inlining here only fires on small functions whose arguments
are rarely both effectful and interfering, and the tests never meet a case
where the moved evaluation changes a value.

**What the coverage rule does with them.** A fault build passes the rule
only if every change is vouched: 42 of the 2,154 changed builds of the
faults that do not crash the compiler (2%; 23 `drop-unused-init`, 17
`swap-args`, 2 `inline-args`, all harmless on their input as far as read),
against 440 of the loop cap's 3,098 (14%). By distinct change, the faults
that change results are mostly failed or blind (`hook-negation`: 3 vouched,
224 failed, 2,189 blind; `inline-vars`: 12% vouched, 28% failed), while
`inline-args`, `swap-args` and `drop-unused-init` are vouched about as often
as the loop cap (48%, 35% and 47%, against 57%): where a fault's change runs
without changing anything a test sees, the oracle cannot tell it from a
correct change, and it does not try to. "Vouched" means "ran and did the same
thing on these inputs", not "correct in general".

**The benign variants are never failed.** `benign-rename` (a pure renaming)
leaves 607 of its 610 builds alpha-isomorphic under the diff; the other three
are a renamed property name that Closure writes into a string literal
(`goog.reflect.objectProperty`, a real output difference, vouched) and two
ADVANCED ES2015 builds (marked, tsc) where property ambiguation makes the
renaming non-bijective and the diff reports changes (blind, never failed).
`benign-no-inline` changes all 610 builds and is vouched for 54% of its
distinct changes, failed for none: the same profile as the loop cap.

### Different but correct: the loop cap (phase 2)

`--max_optimization_loop_iterations=N` (merged opt-in in `8f5b28b9`; Java's
`CompilerOptions.setMaxOptimizationLoopIterations`) stops each code-removal
loop after N rounds. It changes the output and keeps behaviour, so it is the
ready-made "different but correct" input: `run --fast-flags=--max_optimization_loop_iterations=1`
adds it to the fast builds only. (It is a test input; nothing turns it on by
default.) With N = 1, on the per-change set (`run-cap`, 4,458 builds):

| suite | builds changed | rule: vouched / blind / fail builds | changes vouched / blind / fail | changes vouched | distinct changes vouched / blind / fail |
|---|---:|---|---|---:|---|
| test262 (600 × 4) | 1,048 of 2,400 | 438 / 610 / 0 | 3,341 / 3,828 / 0 | 47% | 511 / 357 / 0 |
| Closure runtime tests | 830 of 838 | 0 / 830 / 0 | 67,371 / 73,266 / 0 | 48% | 1,290 / 226 / 0 |
| Closure Library tests | 1,208 of 1,208 | 0 / 1,208 / 0 | 112,298 / 115,467 / 0 | 49% | 8,410 / 3,095 / 0 |
| lodash | 4 of 4 | 0 / 4 / 0 | 784 / 21 / 0 | 97% | 434 / 11 / 0 |
| marked | 4 of 4 | 2 / 2 / 0 | 95 / 25 / 0 | 79% | 24 / 0 / 0 |
| tsc | 4 of 4 | 0 / 4 / 0 | 4,159 / 7,682 / 0 | 35% | 3,627 / 7,124 / 0 |
| all | 3,098 | 440 / 2,658 / 0 | 188,048 / 200,289 / 0 | 48% | 14,296 / 10,813 / 0 |

(Distinct changes: the same function changed the same way counts once, and
is vouched if some build ran it and behaved the same; the per-suite split
attributes a change to the first build it was seen in.)

- **Nothing fails.** No build behaves differently on anything a test checks.
  38 builds (31 Closure Library, 7 test262) differ only inside tests that
  already fail on the compat build, in failure messages that print compiled
  code: goog.testing's PropertyReplacer quotes a function's source text, and
  renamed receivers are numbered in a different order. The rule keeps their
  changes blind (`differs-unchecked`); they are not failures. One false
  alarm was found and fixed on the way: the tsc workload digested its source
  maps, whose renamed key ADVANCED had named differently in the two builds
  (above).
- **Where tests run the code, they vouch for it:** 97% of the changes in
  lodash, 79% in marked, and per distinct change 73% of the Closure Library
  changes and 85% of those in Closure's runtime tests.
- **Per build, almost nothing is vouched outright**: 440 of 3,098 changed
  builds, because nearly every build carries some changed code that its test
  never runs: polyfills behind feature tests, error paths, helpers the test
  does not call (`$jscomp.getGlobal`'s throw when no global object exists,
  lodash's buffer paths, the polyfill branch for array-likes without
  `Symbol.iterator`). A per-build gate would refuse every non-isomorphic
  change; a gate per distinct change (phase 4) vouches for 57% of them here
  and lists the rest.
- **The bench-scale workload is where the blind spots are**: the TypeScript
  compiler checking a 150-line program runs 35% of the changed code. A
  checker this size needs a larger corpus of programs to vouch for most of
  its changes (phase 3).

### Mutation score, before and after

| | phase 1 (test262 and Closure's runtime tests, 4 configurations) | phase 2 (library suites, 2 configurations) |
|---|---|---|
| faults caught (of 7) | 5 (all but `inline-args`, `swap-args`) | 4 of the 6 that change a build (`fold-negzero` changes none) |
| changed signal builds caught | 1,106 of 5,946 (19%) | 837 of 2,459 (34%) |
| without `collapse-unsafe` (compiler crashes) | 706 of 5,546 (13%) | 532 of 2,154 (25%) |
| `hook-negation` / `inline-vars` / `drop-unused-init` | 31% / 12% / 8% | 62% / 23% / 3% |
| benign changes failed | (none tried) | 0 (loop cap, `benign-rename`, `benign-no-inline`) |

Together the suites catch the same five faults as phase 1: the library
suites catch the result-changing faults about twice as often per changed
build, and a correct change has not been failed once; the two faults whose
instances are harmless on every input the suites contain stay uncaught, and
the coverage rule vouches for those instances as it vouches for correct
changes.

## Current numbers

### Phase 2: the per-change set

Measured on 2026-11-01 with jsgo `66cdfea8` (the release build; the phase-2
changes to jsgo are seeded variants behind the `jsgo_faults` tag), node
22.22.0, 4 jobs (`run-pc`):

| suite | tests (cases × configurations) | builds that check something (compat) | tests passing on compat, per configuration (SIMPLE ES5 / ES2015 / ADVANCED ES5 / ES2015) | fast bytes = compat | differences |
|---|---:|---:|---|---:|---:|
| test262 (per-change 600) | 2,400 | 1,952 | 527 / 542 / 438 / 445 (whole tests) | 2,205 of 2,205 compiled | 0 |
| Closure runtime tests | 838 | 814 | 1,098 / 722 / 1,096 / 722 | 830 of 830 | 0 |
| Closure Library tests | 1,208 | 1,208 | 4,307 / 4,307 / 3,834 / 3,830 | 1,208 of 1,208 | 0 |
| lodash | 4 | 4 | 3,084 in each | 4 of 4 | 0 |
| marked | 4 | 4 | 127 / 127 / 123 / 123 | 4 of 4 | 0 |
| tsc | 4 | 4 | the whole run, in each | 4 of 4 | 0 |

Per-test signal changes what counts: the Closure Library builds whose every
test passes are 254 (SIMPLE) and 210 (ADVANCED) of 302, but all 1,208 builds
have passing tests, 16,278 in all; per file (phase 1's rule) 280 of them,
nearly a quarter, would carry no signal. The library suites add 29,114
passing tests per run (16,278 + 12,336 + 500) to test262's 1,952 and the
3,638 of Closure's runtime tests (their goog.testing results, counted the
same way).

**Compat vs fast: still no difference**, on every build of every suite,
including the 8.9 MB TypeScript compiler.

### Phase 1: test262 and Closure's runtime tests

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

### Phase 1

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

### Phase 2

The library suites have no uncompiled form for Closure Library (goog.module
needs the compiler), and lodash, marked and tsc compile as one case, so the
check was done by hand on the builds that lose tests to compilation; none is
a port bug:
- lodash: every build fails the same 11 tests as uncompiled.
- marked: the ADVANCED builds fail 4 specs that pass uncompiled
  (`codespan_newline`, `em_list_links`, `unicode_punctuation`,
  `link_redos`); both builds are byte-identical to Java's.
- tsc: the ADVANCED builds print different source maps than uncompiled
  (TypeScript serialises an object literal whose key `mappings` ADVANCED
  renames); the ADVANCED ES5 build is byte-identical to Java's.
- Closure Library: 44 test files pass some tests at SIMPLE and none at
  ADVANCED in some configuration (88 builds), checked against the pinned
  Java compiler with the same command lines: 82 builds are Java's bytes
  (ADVANCED renames a name the test or goog.testing reaches by string, or
  removes code it needs), 4 fail to compile with Java's error
  (`JSC_MSG_ORPHANED_NODE`), and **2 were a port bug**: jsgo stopped with
  `java.lang.IllegalStateException: Only use on nodes with colors present` on
  `net/channelrequest_test.js` (both ADVANCED configurations), where Java
  compiles. ConvertTypesToColors has a Go-only walk over the subtrees the
  checks detached but the compiler still references (so that their JSTypes
  can be collected); Java never visits those nodes. Their types were never
  gathered, so a cast among them got no color, and marking it as colored
  from a cast threw. The walk now leaves such a node uncolored and unmarked
  (`9445b09a`, with a regression test that throws without the fix); both
  builds are now Java's bytes, and the bench Java gate (`typed-app`,
  `typed-full`, `downlevel-simple`, `bundle-advanced`,
  `bundle-advanced-ts`, one run each) matches. Only nodes that crashed take
  the new path, so no other output changes.

## Wall time

On the 4-CPU machine, `--jobs 4`.

**Per change (phase 2).** The per-change set is test262's 600-test subset,
Closure's runtime tests, the 302 Closure Library tests, lodash, marked and
tsc, in four configurations: 4,458 builds per mode.

| run | wall time |
|---|---:|
| phase 1's per-change run (1,500 test262 tests and Closure's runtime tests, no cache) | 49 min |
| phase 2, compat cache cold (every compat build compiled and run: the first run after a pin moves) | about 40 min (measured 32.3 with the Closure Library's compat builds already cached by its probe, whose compiles take 7.5 min of wall time) |
| phase 2, compat cache warm, every fast build different from compat (the loop cap: every changed build is diffed and both builds run under coverage) | **22.1 min** (test262 5.4, Closure runtime 5.8, Closure Library 8.5, tsc 2.0, lodash 0.3, marked 0.1) |
| phase 2, compat cache warm, no fast build different (fast mode today) | about 19 min (the fast compiles alone: 73 CPU-minutes) |

The two measures phase 1 proposed are both in, and both were needed:
- **The compat cache.** A compat build is Java's output, so it depends only
  on its inputs, its flags and the pinned Closure release, not on the jsgo
  under test (the parity suites check jsgo against Java separately): the key
  is a digest of those (and of every file a flag names), and the cache keeps
  the compat build, its compile result and its runs (keyed further by the
  node version, the runner and the run options). A change then compiles and
  runs only its fast builds, and runs those only when they differ from
  compat. This halves the compiles. The key leaves the jsgo binary out on
  purpose; `--compat-cache none` rebuilds everything, and the nightly full
  run should use it, so that a stale cache cannot hide a compat change.
- **A smaller test262 subset**, 600 instead of 1,500 (`behave.py subset`).
  test262's small programs rarely reach what fast mode changes (under the loop
  cap 1,352 of its 2,400 builds stay byte-identical), while the library
  suites are where changed code runs and is checked.

What is left of the 60 minutes (about 38 in the worst case) is room for phase
3's workloads and for running each change's differing builds a second time
if a gate wants it.

**Other runs.**

| run | what | wall time |
|---|---|---:|
| full (`run --subset full --java`, phase 1) | 5,000 test262 tests and the Closure cases, then the Java check | 2 h 33 min |
| `faults`, phase 1 suites | each fault rebuilds the per-change fast builds | 18–32 min per fault |
| `faults`, library suites (2 configurations) | each fault or benign variant | 3.4–8.6 min |
| Closure Library probe (443 files, 4 configurations, both modes) | | 21.4 min |
| `probe` (when the test262 pin moves) | 43,117 tests uncompiled, 1,324 compiles | 21 min |

Compiling dominates: a test262 compile takes about 0.27 s at SIMPLE and 0.7 s
at ADVANCED (mostly the default externs), a Closure Library test 1 to 2 s, a
tsc build 8 to 13 s; a node run 0.1 to 1.5 s, and 1.5 to 2 times that under
coverage. The coverage rule's diff takes 0.1 to 0.3 s per test build and
about 4 s for a tsc build (3 MB).

## Adding a suite

A suite is a function returning `core.Case` objects: the directory the
compiler runs in, the `--js` inputs in order, extra flags, the uncompiled
program (or none), runner options, and a verdict function from the runner's
record and stdout to pass/fail; register it in `behave.py` (`SUITES`,
`load_cases`) and pin its sources in `pins.json` and `fetch.sh`. The engine
does the rest (compat cache, compiling, running, comparing, the coverage rule,
the Java check, faults, reports). What phase 2 learnt:
- Compile the library with its tests, so that ADVANCED renames both, and
  generate externs for the names reached by string (`libs.string_names`:
  identifier-like string literals, properties read from the global object,
  keys read back with `Object.entries`); declare host names the code reads
  (`libs/qunit_externs.js`).
- Report per-test results in the record (`rec.suite.tests`, as goog.testing
  and the QUnit shim do), so that each test carries its own signal; keep
  failure messages free of values the compiler names.
- When a program prints its results, print them in full or digest only what
  the program itself chose: a digest hides a renamed name from the
  normalisation (the tsc source maps).
- Anything the test harness needs from outside the build (a second copy of
  the library, a test framework) runs uncompiled, in the runner or a `--pre`
  script.

## Moving the pins

Change the revision in `tools/behave/pins.json`, run `fetch.sh` (it refetches
a checkout whose commit differs), regenerate the sample if test262 moved
(`behave.py probe`, then `behave.py sample --probe DIR`), and record the new
digests with `behave.py verify --update`.
