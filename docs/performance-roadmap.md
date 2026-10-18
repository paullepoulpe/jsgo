# Performance roadmap and open explorations

Status on 2026-10-18, after the upgrade to Closure Compiler v20261015. This
file records what we know about jsgo's performance, the experiments under
way and the options the repository owner has discussed, so that the work
can be resumed from the repository alone.

## Where things stand

- **Parity.** jsgo is byte-identical with Java v20261015 on every
  end-to-end suite: 40,005 fresh Java builds compared (projects 120/120,
  held-out sweep 7,077/7,077, development flag sweep 16,200/16,200,
  held-out flags 3,340/3,340, test262 7,947/7,947, development samples
  1,360/1,360, bench gate 5/5 tasks). The only differences are accepted
  JVM-only ones: Java's stack frames (11 of the 1,000 single files) and
  one fixture where Java catches a StackOverflowError (PLAN.md §7).
- **Speed** ([benchmarks.md](benchmarks.md#end-to-end-v20261015)):
  against `java -jar`, jsgo is 1.4–2.3× faster on whole-program builds
  and 4.2× faster on 1,000 single-file compiles; the 27 s bundle is the
  exception at 0.94×. Against a warm JVM it is at 0.65–0.74× of Java's
  speed, with 2–5× less memory. An interleaved control showed that the
  round-3 figure of 0.93× against warm Java was measured against a loaded
  JVM; on a quiet machine the round-3 code was also at about 0.65×.
- **Where warm Java wins** (per-phase tables in
  [benchmarks-e2e-v20261015.md](benchmarks-e2e-v20261015.md)): the
  optimization loop (`flowSensitiveInlineVariables` 2.4–2.8× Java's time,
  `inlineVariables` ~2×, `removeUnusedCode` 1.5–1.9×,
  `collapseObjectLiterals` and `coalesceVariableNames` ~2×), parsing
  (155 ms per single-file compile against 61 ms, mostly the default
  externs), `inferTypes` (1.4–1.8×) and `checkAccessControls` (~2.1×).

## Why: the CPU profile

The profile behind `cmd/jsgo/default.pgo` (typed-app and friends, v20261015
port) attributes:

- **~37% of CPU samples to memory management:** `runtime.gcDrain` 27%
  cumulative (`scanobject` 25%), `mallocgc` 11%. Go's collector is
  non-generational: every cycle re-traces the whole live heap (the AST,
  the type graph, scopes), while the translated passes allocate
  short-lived objects at Java's rate. HotSpot's generational collector
  makes exactly this pattern cheap; Go's makes it expensive.
- **9% flat in `rhino.(*Node).GetNext`**, a one-line accessor: memory
  latency from AST nodes scattered across the heap.
- **~8% in `NodeTraversal.GetAbstractScopeAt`:** scopes rebuilt per pass.
- **~5% in Java emulation for byte-identity:** `jutil.StringHashCode`,
  `StringCompareTo`, `StringFromUTF16`, Java-ordered hash maps.

## Options, in the order discussed with the owner

1. **GC setting** (`GOGC` / `GOMEMLIMIT`). jsgo uses 2–5× less memory than
   Java, so trading some of it for fewer GC cycles is cheap and
   reversible; estimated 10–25%. The owner chose Go's default for now
   (PLAN.md §7). The GC lab (below) measures it; its main use is to size
   option 3.
2. **Incremental, inside the 1:1 structure:** cut allocation in the slow
   passes named above; parse the default externs once (an embedded,
   pre-parsed snapshot) for the single-file case. Expected: about parity
   with warm Java, perhaps 1.2×.
3. **Substrate redesign, passes kept 1:1:** a pointer-free AST (nodes in
   arenas linked by int32 indices, so the GC skips them and traversal is
   cache-friendly), side tables keyed by node index instead of maps keyed
   by `*Node`, interned string ids instead of emulated Java strings,
   scopes reused across passes; then the same treatment for the other
   pointer-heavy structures (JSType graph, Var/Scope/Reference, CFG and
   data-flow lattices, jutil collections). Upstream rarely changes this
   substrate, so the ~1,200 ported files and upstream tracking stay
   cheap. Estimated 2–3× warm Java on the large tasks. Once the AST is
   pointer-free, option 1's gain shrinks to a few percent.
4. **Parallelism** (parsing, per-function work): the biggest lever for
   large builds; deferred by the owner (no goroutines for now).

**Byte-identical output stays.** Relaxing it would buy about 5% directly
(the Java emulation above) and would remove the oracle that makes
option 3 safe to do.

## Output quality: improving on Closure

Possible, as deliberate divergences after parity, behind an opt-in flag
so that the default stays byte-identical with Java and every parity
suite keeps working. Candidates: ES-module tree-shaking, modern-syntax
output, gzip-aware renaming, and fixing upstream crashes such as the
one below. Prerequisites: an output-quality benchmark (gzip/brotli size
and run time against terser, esbuild and Closure itself) and an
execution-based oracle (running compiled outputs against test suites),
since byte comparison with Java no longer applies to changed output.

**Upstream crash found by the parity re-run:** v20261015 throws a
`checkNotNull` NullPointerException in
`RemoveUnusedCode.getVarForNameNode` on ES-module files in the `*-lax`
configurations (588 of the 7,077 held-out sweep jobs: monaco-editor,
mathjs, svelte, axios, date-fns, prettier), plus 7
`IllegalStateException: AST should not contain Dynamic module import`.
v20260415 had none. jsgo reproduces Java's crash (PLAN.md §7).

## GC lab (2026-10-18, in progress)

An experiment measuring the GC tax (`GOGC` 100/200/400/off, and off with
Java's memory budget) on typed-app, bundle, downlevel-simple, checks-only
and 100 single files, plus the live heap and allocation churn by
component and a traversal-locality micro-benchmark. Results:
[gclab-2026-10-18.md](gclab-2026-10-18.md) (added when available; if it
is missing, the experiment was cut short and should be re-run as
described here).

## Smaller open items

- `tools/clie2e/sweep.go` runs its JVMs under its own `nice -n 10`; under
  an outer `nice` they end at 19 and run ~5× slower when other work
  competes. Drop the inner `nice` or document it.
- The released jar still cannot load its own XTB reader
  (`NoClassDefFoundError …jarjar/javax/xml/parsers/ParserConfigurationException`);
  `tools/clie2e/build-java.sh` works around it for the oracles.
