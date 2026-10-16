# jsgo: assessment and plan

**Goal:** a drop-in replacement for Google Closure Compiler written in Go —
same command-line flags, same JavaScript output, same source maps, same
diagnostics — so that existing Closure users can switch binaries.

**Verdict:** feasible, but not with the approach the repository had been
following. The first attempt was a fixture-tuned type checker that does
not generalise, and it has been removed. It remains in git history at
commit `c05fc09`. Work has restarted as a bottom-up, whole-file, 1:1 port of the
pinned Java release, verified layer by layer against the Java compiler
itself. The rest of this document gives the evidence for both halves of that
verdict and the roadmap.

---

## 1. What "drop-in replacement" means in size

Closure Compiler v20260415 is **813 Java files, 273k lines, ~177k lines of
code** (excluding comments and blank lines). Rough breakdown by area:

| Area | Java SLOC | Share |
|---|---:|---:|
| Front end: AST (`rhino`), parser, IRFactory, JSDoc parser | 22.8k | 13% |
| Compiler core: driver, options, CLI, diagnostics, scopes, traversal, NodeUtil | 44.5k | 25% |
| Code printer + source maps | 5.4k | 3% |
| Checks and lint (untyped) | 7.8k | 4% |
| Transpilation, polyfills, module systems | 24.3k | 14% |
| Simple optimizations (peephole, inlining, dead code) | 19.8k | 11% |
| Type system and type checking | 25.4k | 14% |
| Advanced optimizations (property renaming/collapsing, devirtualization, disambiguation) | 20.5k | 12% |
| Serialization, instrumentation, refactoring, Ant | 6.3k | 4% |

A faithful Go port will be roughly 230–260k lines. That is large but
bounded, and of the same order as Microsoft's TypeScript-to-Go port.

The first attempt targeted only the type-checking diagnostics, about 14% of
the compiler. It never produced JavaScript output, and it had already grown
to 61k lines of Go for a slice whose Java is ~25k lines.

## 2. Why the previous approach did not work

### Measured

The legacy checker reports **93.7%** on its 1,444-fixture development set,
which was extracted from `TypeCheckTest.java`. The README's 83% figure is
stale. Note that 40% of those fixtures expect *no* diagnostics, so a no-op
checker already scores 40%.

To test generalisation, I extracted **544 held-out fixtures** from the sibling
test files written by the same Closure team for the same checker
(`TypeCheckNoTranspileTest`, `TypeCheckTemplatizedTest`,
`TypeCheckCovarianceTest`). I generated their oracles with the same pinned
jar and flags.

| | Development set | Held-out set |
|---|---:|---:|
| Fixtures with expected diagnostics: exact match | 90.4% | **26.4%** |
| Fixtures expecting no diagnostics: no false positives | 98.6% | **68.9%** |
| Overall | 93.7% | **41.7%** |
| Process crashes (unrecoverable stack overflow in `IsSubtypeOf`) | 0 | 2 |

Even the held-out file closest in style to the development set
(`TypeCheckTemplatizedTest`, ES5-style generics) drops to 48% on fixtures
with expected diagnostics. Most misses are plain `JSC_TYPE_MISMATCH`es in
ordinary modern code: destructuring, `async`/`await`, class fields and arrow
functions.

Other observations on new inputs:

- `function g(xs){ return xs.map(x => x*2); }` gets a false
  `JSC_UNDEFINED_VARIABLE`. Java reports nothing.
- Output order is nondeterministic: one input produced 3 different
  diagnostic orders in 12 runs, because properties live in Go maps. The
  tuple-sorting harness could not detect this.
- **Real-world code** (162 files: 150 from Closure Library plus lodash,
  jQuery, d3, three.js, React, Vue, moment, acorn, esprima):
  - Java parses all 162; jsgo parses 134. Its lexer rejects any regex with
    an escape (`var a = /\d/;` is a "parse error") and any template literal
    with `${}`, so 11 of the 12 plain libraries fail to parse.
  - `for (var i = 0, n = a.length; i < n; i++)` is silently misparsed into
    four bogus diagnostics.
  - Two more stack-overflow crashes are reachable from ordinary code
    (optional chaining `o?.a`; self-referential class properties).
  - Where both parse, jsgo's diagnostics agree with Java's on only 2–12% of
    tuples.
  - One quadratic pass makes jsgo 5–6× *slower* than the JVM on large files
    such as lodash.
- The repository's own guard script (`ci-guards.sh`) was failing on `main`.

### Root causes

An independent read of the code against the Java source found the following
(details in the session report; examples are reproducible at
`c05fc09`):

1. **Development was driven by fixtures, top-down.** Each change made some
   fixtures pass: 55 of the 96 visible commits cite fixture names, pass-count
   deltas or false positives. The result is heuristics that Java does not
   have. Examples:
   - suppressions keyed on what the oracle shows;
   - a raw-text scan of `/** … */` comments that emits
     `JSC_UNRECOGNIZED_TYPE_ERROR`;
   - checking for `=` on the source line to decide whether a variable has an
     initializer;
   - name-keyed de-duplication maps for diagnostics;
   - hard-coded builtin name lists to avoid "false-positive storms".

   Java's 15-line `visitGetProp` became 566 lines.

2. **The foundation was never ported.**
   - The parser emits AST nodes directly, with no parse-tree or IRFactory
     layer; the alternative pipeline silently drops 25 of 85 tree kinds.
   - JSDoc types are kept as strings and re-parsed by splitting on `|`.
   - The type registry is a switch over primitive names.
   - Scopes are `map[string]JSType`.
   - Type inference is re-run on demand, without dataflow.

   Every later Closure pass depends on these layers. This is why it cannot
   be fixed in place: the object model differs exactly where everything else
   depends on it.

3. **The metrics were gamed rather than met.**
   - 1,532 functions carry a `// Go-only: language-adapter` tag, which
     exempted them from the fidelity score; the tag was routinely applied to
     false-positive guards.
   - Java-named methods were added with empty bodies to raise
     method-overlap scores.
   - About 29% of functions are unreachable.

The lesson is not that porting Java to Go is hard. It is that the
verification signal — a fixed set of end-to-end diagnostic fixtures — was
too narrow and too easy to overfit, and that building top-down on an
unported foundation made every fix a local patch.

## 3. Why the project is feasible

- **Closure is deterministic by design.** The Java code constructs
  `LinkedHashMap`/`LinkedHashSet` 509 times and unordered
  `HashMap`/`HashSet` only 10 times. With ordered collections in Go
  (`internal/jutil`), byte-identical output is a realistic target, and that
  turns verification into a diff.
- **The Java→Go mapping is mechanical for nearly everything.** The hazards
  are enumerable and handled once, in `internal/jutil` (see
  [PORTING.md](PORTING.md) §6):
  - UTF-16 string semantics;
  - iteration order;
  - Java's saturating casts and masked shifts;
  - `Double.toString` (verified against JDK 21 on 27k vectors);
  - exceptions as panics with catch-by-subclass;
  - 52 `java.util.regex` sites to audit for RE2 compatibility.
- **The packages layer cleanly.** Treating generated protobuf classes as
  leaves, the Java package graph condenses to a DAG:
  - `sourcemap`, `base`, `parser` → `rhino` (AST + types) → `parsing` →
    `jscomp`;
  - the front end (23k SLOC) can be ported and verified before anything
    else exists.
- **There is an oracle for every layer, on unlimited input.**
  [`tools/oracle/`](tools/oracle/) runs the pinned jar to produce ground
  truth for any JavaScript file:
  - the exact AST, with positions, properties and JSDoc;
  - CLI output and diagnostics.

  More oracles will follow per milestone: printed code, per-pass ASTs and
  source maps. The corpus is not limited to hand-picked fixtures: it is
  14,690 real-world files (Closure Library, test262 parser tests, three.js,
  core-js, TypeScript's compiler, React, Angular, …), pinned and
  integrity-checked ([`tools/oracle/corpus.txt`](tools/oracle/corpus.txt)),
  plus Closure's own unit tests, which can be ported. For example,
  `ParserTest.java` is 8.8k lines and `JsDocInfoParserTest.java` is 6.5k
  lines.
- **A 1:1 port can follow upstream.** Each new Closure release is a Java
  diff that maps file-for-file onto the Go code. A redesign cannot do that.

### Risks

- **Volume.** About 800 files have to be translated consistently.
  Mitigation: strict mechanical conventions, so that independently ported
  files fit together (PORTING.md); whole-file ports in dependency order; and
  per-layer oracles that catch inconsistency immediately.
- **Hidden Java semantics.** Examples are `HashMap` iteration order at the
  ~10 unordered sites, identity hash codes, `String.format` details, regex
  dialect and `char` arithmetic. Mitigation: centralise them in `jutil`
  with JVM-generated test vectors, and treat any oracle diff as a port bug
  to root-cause, never to special-case.
- **Performance.** A line-for-line port inherits Java's allocation
  patterns. Go should still beat the JVM on wall-clock for typical builds,
  because there is no JIT warm-up or startup cost. Optimise after parity,
  behind the same oracles.
- **Reading discipline.** The previous attempt shows how easily
  "make the test pass" displaces "port the Java". The process below makes
  that structurally hard.

## 4. Strategy

1. **Bottom-up, in dependency order.** No layer is built on an unported one.
2. **Whole files, faithfully.** One Java file → one Go file, every member,
   Java order, no stubs, no subsets. The only permitted deviations are
   language adapters, and they are marked.
3. **Layer-wise differential verification against Java**, on the full
   corpus:
   - AST dumps for the front end;
   - printed code for the printer;
   - per-pass ASTs for transforms;
   - CLI output for the driver.

   Ported Java unit tests come first. A layer is done when its diff is
   empty, not when a percentage crosses a line.
4. **Never debug toward a fixture.** A diff means a mistranslation: find the
   Java logic and translate it.
5. **Byte-identical output is the target.** It is the cheapest definition of
   "drop-in" that can be checked mechanically.

## 5. Roadmap

Each milestone produces something verifiable. User-visible value starts at
M2.

| | Milestone | Java SLOC | Exit criterion |
|---|---|---:|---|
| M0 | Runtime (`jutil`), conventions, oracle tooling, corpus | — | done (this session) |
| M1 | Front end: `rhino` core (Node, Token, IR, JSDocInfo, …), `parser` (Scanner, Parser, 86 parse trees, FeatureSet), `parsing` (IRFactory, JsDocInfoParser, Config, ParserRunner) | 23k | ported unit tests pass; Go AST dump identical to `AstDump` on 100% of the corpus |
| M2 | Code printer + source maps (`CodePrinter`, `CodeGenerator`, `CodeConsumer`, `debugging/sourcemap`) | 5k | parse→print byte-identical to Java (compact and pretty) on the corpus; identical source maps. **First drop-in use: whitespace-only mode.** |
| M3 | Compiler core: `Compiler`, `CompilerOptions`, `CommandLineRunner` flag parsing, `PhaseOptimizer`, `DefaultPassConfig`, `NodeTraversal`, `NodeUtil`, scopes, error managers/formatters, `DiagnosticGroups`, warnings guards | 45k | CLI accepts every Closure flag; output format and exit codes identical |
| M4 | Untyped checks and lint | 8k | diagnostics identical on corpus |
| M5 | Transpilation (ES2015+ → ES5 …), polyfill injection, module rewriting (goog.module/provide, ES modules, CommonJS), chunks | 24k | identical output across `--language_out` levels |
| M6 | SIMPLE_OPTIMIZATIONS passes | 20k | **byte-identical `SIMPLE` output** on corpus |
| M7 | Type system and type checking (`rhino/jstype`, TypedScopeCreator, TypeInference, TypeCheck, conformance) | 25k | diagnostics identical on corpus and on all extracted Closure type-check tests (including the 1,444 + 544 fixtures already in this repo) |
| M8 | ADVANCED_OPTIMIZATIONS (collapse/rename/disambiguate properties, devirtualization, cross-chunk motion) | 21k | **byte-identical `ADVANCED` output** on corpus |
| M9 | Remainder: serialization (TypedAST), instrumentation, remaining flags; then performance work | 6k+ | full CLI parity; benchmark report vs Java (§6), overall and per phase |

The legacy fixtures and oracles in `testdata/` remain valid and become the
M7 acceptance set.

### Porting order: dependency waves, run in parallel

The milestones above are verification checkpoints. The *porting* order
follows the dependency graph, which is much wider than the milestone list
suggests. The 556 files of the Go `jscomp` package condense into 460
strongly connected components, from 137.7k non-blank, non-comment lines of
Java. `tools/oracle/jscompdeps.py` computes them, and
[docs/port-waves.md](docs/port-waves.md) lists every file by level.

The structure:

- **One core SCC:** 37 files, 15.7k lines. It contains `NodeUtil`,
  `CompilerOptions`, `AstFactory`, `NodeTraversal`, the scope classes,
  `AstAnalyzer`, `ControlFlowAnalysis`, `TypeValidator`, `SourceMap`, …
- **Below it:** ~180 small files, 17k lines (error managers, warnings
  guards, graphs, deps/modules, serialization).
- **Above it:** a very wide layer of ~300 files and ~85k lines (levels
  11–18). This holds the printer, every check, transpilation pass and
  optimization.

Two consequences:

- **The type system moves early.** The core references `TypedScope`,
  `TypedVar`, `TypeValidator` and `JSTypeRegistry`, so `rhino/jstype` (10k
  lines) is ported alongside `parsing`, right after rhino core. Type
  *checking* (`TypeCheck`, `TypeInference`, `TypedScopeCreator`) stays at
  its natural level.
- **The printer needs the core.** `CodePrinter` (level 13) depends on
  `NodeUtil`, `CompilerOptions` and `SourceMap`. M2 is therefore verified
  after the core lands.

| Wave | Contents | Parallelism |
|---|---|---|
| 1 | `parser` (done); `rhino` core, `base`, `dtoa`, `colors`; `sourcemap`; `regex`, `diagnostic`, `resources` | one agent per package |
| 2 | `parsing` (IRFactory, JsDocInfoParser, …); `rhino/jstype`; jscomp levels 0–9 | ~6 agents |
| 3 | jscomp core SCC | ~4 agents coding against each other's mechanically predictable names, then one integration pass |
| 4 | jscomp levels 11–18 | 10–20 agents, by pass family |
| 5 | `Compiler`, `DefaultPassConfig`, `AbstractCommandLineRunner`, `CommandLineRunner` | 1–2 agents |

Each wave starts as soon as the files it depends on are merged, not when a
milestone is declared done. Verification (unit tests, oracles) follows each
merge.

## 6. Benchmarking against the Java compiler

Once the port is a drop-in replacement, benchmark it against the pinned
Java compiler on a representative set of tasks. Report overall numbers and
numbers per compilation phase.

**Tasks.** Each one is a real CLI invocation, run with identical flags on
both compilers:

| Task | Inputs | Mode |
|---|---|---|
| Typed app build | Closure Library-based app with deps (`goog.module`/`goog.provide`) | `ADVANCED`, type checking on, source maps |
| Large plain-JS bundle | three.js, lodash, react-dom, TypeScript's `typescript.js` (9 MB) | `SIMPLE` |
| Modern code down-levelling | rxjs, Angular, core-js | `--language_out=ECMASCRIPT5`, `WHITESPACE_ONLY` and `SIMPLE` |
| Checks only | Closure Library | `--checks_only` with `checkTypes` |
| Many small invocations | 1,000 single-file compiles | `SIMPLE` (startup-dominated; models per-file build steps) |

**Correctness gate.** A run only counts if both compilers produce
byte-identical output, source maps and diagnostics. Otherwise the speed
comparison is meaningless.

**Measurements.** For each run, record:

- wall-clock time;
- CPU time;
- peak RSS.

Each figure is the median of N runs plus its spread, on a quiet, pinned
machine.

The Java compiler gets two rows:

- **cold:** `java -jar`, the way most users run it;
- **warm:** the same invocation repeated in one JVM (`CliBatch`), which
  models build-system persistent workers and isolates JIT warm-up from
  steady-state speed.

Both use default JVM and Go GC settings.

**Per phase.** Closure's `--tracer_mode=TIMING_ONLY` makes
`PerformanceTracker` report time per pass. Because the Go port is 1:1,
including `PerformanceTracker` and `PhaseOptimizer`, it runs the same passes
under the same names and emits the same report. The per-pass tables
therefore line up exactly. Passes are also rolled up into phases:

- parse;
- checks;
- type checking;
- transpilation;
- optimization loop;
- renaming;
- code printing and source maps.

**Earlier, per milestone.** A line-for-line port inherits Java's allocation
patterns, and the first attempt showed that a single quadratic pass can
make Go slower than the JVM. So each layer is benchmarked as it lands:

- **M1:** parse the 14,690-file corpus;
- **M2:** parse and print it;
- **later milestones:** the passes they add.

This way a slow foundation (e.g. node allocation, ordered maps, UTF-16
conversion) is caught before everything else is built on it. The full
benchmark above becomes the M9 acceptance report.

## 7. Status

Done:

- **M0.**
  - [PORTING.md](PORTING.md): the translation contract.
  - `internal/jutil`: Java-semantics runtime, tested against JVM-generated
    vectors.
  - `tools/oracle`:
    - pinned jar and source fetch with checksum;
    - `AstDump`, the exact Java AST for any file;
    - `CliBatch`, the CLI oracle without per-file JVM startup,
      byte-identical to the per-process oracles;
    - the package-SCC generator;
    - the pinned 14,690-file corpus.
  - The first attempt has been deleted; it is available at `c05fc09`.
    Its fixtures and oracles stay in `testdata/`.

- **M1, part 1: `internal/parser`.**
  - All 104 files of `jscomp.parsing.parser` (+ `trees`, `util`) are
    ported: 11.1k lines of Go for 10.1k of Java, with 22 marked language
    adapters.
  - The Java parse tree, diagnostics and comment list are reproduced byte
    for byte (`tools/oracle/java/ParseTreeDump.java`), with zero diffs on:
    - the 14,690-file corpus;
    - the 1,488 fixtures;
    - 29,892 randomly broken variants;

    under the ES_NEXT, ES5 and ES3 configurations.
  - Independently re-checked on a fresh corpus sample, the 544 held-out
    fixtures and 59 new adversarial inputs.

- **M1, part 2: `internal/rhino` core, `internal/base`, `internal/dtoa`,
  `internal/colors`.**
  - Node, JSDocInfo, IR and the rest, ported with the Java unit tests
    (IRTest, NodeTest, JSDocInfoTest, QualifiedNameTest, HamtPMapTest, …).
  - Formatting matches Java on a differential oracle.
  - DToA matches Java on 65,675 numbers.
- **First benchmark** ([docs/benchmarks.md](docs/benchmarks.md)): the
  ported parser against Java's.
  - After allocation fixes that are now conventions, Go takes 6.6 s per
    corpus parse, against Java's 5.9 s cold and 2.2 s warm.
  - Go has no warm-up cost and uses ~4× less memory.

- **Wave 1 complete:**
  - **`internal/sourcemap`:** all 19 classes.
    - Differential oracle on the generator and consumer API: 50,572
      results identical.
    - Independently re-run with a new seed: 84,012 identical.
  - **`internal/regex`:**
    - 253,105 regex patterns identical to Java, including every literal in
      the corpus.
    - Independently re-checked on 1,196 new pattern/flag combinations.
  - **`internal/diagnostic`, `internal/resources`:** the jar's JS
    resources are embedded.
  - **`internal/jutil`:** consolidated across branches (one helper per
    Java class, all tested against JDK-generated vectors).

- **M1 front end complete: `internal/parsing`.**
  - Ported: IRFactory, ParserRunner, JsDocInfoParser, JsDocTokenStream,
    TypeTransformationParser, Config, FeatureCollector, ….
  - `ParserRunner.parse` output is byte-identical to Java's
    (`ParsingDump`). That covers the AST with JSDoc, features, comments and
    every error and warning, on:
    - the 14,690-file corpus;
    - the 1,488 fixtures and 544 held-out fixtures;
    - 59 adversarial inputs;
    - all JsDocParsing modes, ES3, ES5 and ES2015;
    - KEEP_GOING;
    - 29,892 broken mutants.
  - Independently re-verified on a **fresh held-out corpus**
    (`tools/oracle/corpus-holdout.txt`): 10,145 real-world files from 16
    npm packages never used during porting. 10,145 of 10,145 identical.
  - Closure's own `ParserTest` (649 tests) and `JsDocInfoParserTest` (778)
    are ported. The only skips need `rhino/jstype`, plus one Java
    stack-overflow case that Go cannot reproduce.

- **jscomp levels 0–9, group B** (80 files: graph library, messages,
  name generators, variable maps, tracer, code consumer, …).
  - Merged with 26 ported Java test classes.
  - Zero diffs on 4,042 seeded differential cases. Re-run with a new
    seed: zero diffs.

In progress, in parallel (wave 2, PLAN.md §5):
- `rhino/jstype`, against a type-algebra oracle.
- jscomp levels 0–9, group A (89 files: diagnostics, error managers, deps/modules, serialization basics).
