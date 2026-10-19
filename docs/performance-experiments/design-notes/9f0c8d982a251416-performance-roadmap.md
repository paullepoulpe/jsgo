# Performance roadmap and open explorations

Updated 2026-10-18. The active goal is ordinary Go CLI performance against
stock warmed Java, with a fresh compiler and normal extern parsing for every
compilation. Worker processes, parsed-extern caches and embedded pre-parsed
extern snapshots are outside scope. Parallel compiler work is authorized using
the existing thread parameter; additional user-facing configurations are not a
goal.

[performance.md](performance.md) records the retained changes and measured CLI
results. The baseline figures and GC lab below are earlier snapshots from the
same date, before the current redesign experiments. Their measurements remain
historical evidence, not performance claims for unmeasured prototypes.

## Earlier v20261015 baseline

- **Parity at that baseline.** jsgo was byte-identical with Java v20261015 on every
  end-to-end suite: 40,005 fresh Java builds compared (projects 120/120,
  held-out sweep 7,077/7,077, development flag sweep 16,200/16,200,
  held-out flags 3,340/3,340, test262 7,947/7,947, development samples
  1,360/1,360, bench gate 5/5 tasks). The only differences are accepted
  JVM-only ones: Java's stack frames (11 of the 1,000 single files) and
  one fixture where Java catches a StackOverflowError (PLAN.md §7).
- **Baseline speed** ([benchmarks.md](benchmarks.md#end-to-end-v20261015)):
  against `java -jar`, jsgo was 1.4–2.3× faster on whole-program builds
  and 4.2× faster on 1,000 single-file compiles; the 27 s bundle is the
  exception at 0.94×. Against a warm JVM it was at 0.65–0.74× of Java's
  speed, with 2–5× less memory. An interleaved control showed that the
  round-3 figure of 0.93× against warm Java was measured against a loaded
  JVM; on a quiet machine the round-3 code was also at about 0.65×.
- **Where warm Java won at that baseline** (per-phase tables in
  [benchmarks-e2e-v20261015.md](benchmarks-e2e-v20261015.md)): the
  optimization loop (`flowSensitiveInlineVariables` 2.4–2.8× Java's time,
  `inlineVariables` ~2×, `removeUnusedCode` 1.5–1.9×,
  `collapseObjectLiterals` and `coalesceVariableNames` ~2×), parsing
  (155 ms per single-file compile against 61 ms, mostly the default
  externs), `inferTypes` (1.4–1.8×) and `checkAccessControls` (~2.1×).

## Earlier CPU profile

The earlier profile behind `cmd/jsgo/default.pgo` (typed-app and friends,
v20261015 port) attributed:

- **~37% of CPU samples to memory management** (CPU, not wall time:
  the GC lab below measured the wall-time cost at 3–7%): `runtime.gcDrain` 27%
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

## Active exploration priorities

1. **Reduce allocations in ordinary compilations.** Parse directly from scanner
   tokens into final Rhino nodes, avoiding intermediate parse trees where exact
   semantics are implemented. Keep the original parser as the oracle/fallback.
   Use compact native CFGs, indexed data-flow state and bounded persistent type
   maps to reduce temporary graph/collection objects. Externs still follow the
   normal per-compilation parsing path.
2. **Avoid repeated pass work.** Reuse analysis only while its canonical scope,
   AST position and enclosing change stamps remain valid. Specialize hot
   traversal callbacks without changing visits, mutation order or diagnostics.
   [Function analysis](experiments/function-analysis.md) and
   [concrete traversal](../tools/bench/concrete-traversal.md) describe prototypes;
   neither represents a measured gain until the CLI comparisons establish it.
3. **Parallelize independent phases within `--num_parallel_threads`.** Bounded
   parsing and dependency scanning are implemented; output and per-function
   analysis candidates need deterministic publication and race/byte gates.
   Analysis workers must finish reading before serial AST rewrites resume.
   [Parallel analysis](experiments/parallel-analysis.md) records the boundaries.
   The parameter controls compiler work items, not every goroutine or OS thread.
   `GOMAXPROCS` remains an independent runtime CPU limit.
4. **Improve storage/locality where profiles justify it.** AST traversal order,
   collection storage and type-graph representation remain research targets.
   A pointer-free indexed AST would require preserving pointer identity and
   mutation semantics at existing interfaces; it is not a promised speedup.
   [Indexed flow slots](experiments/indexed-flow-slots.md) and
   [source locks](experiments/native-source-locks.md) are narrower experiments.

Default GC behavior remains unchanged. The GC lab below bounds the measured
benefit of increasing heap headroom; allocation and algorithmic work remain
important even when collection costs are reduced.

**Byte-identical output stays.** Each candidate must preserve the existing
compatibility contract and pass differential checks before its timing matters.
Measure candidates separately and together; isolated percentages cannot be
added to predict compound improvements. There is no supported theoretical
maximum or current 2–3× target against warm Java.

## Separate future topic: output quality

Outside the current performance scope: deliberate divergences after parity,
behind an opt-in flag so that the default stays byte-identical with Java and every parity
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

## GC lab (2026-10-18): historical measurements

Full write-up: [gclab-2026-10-18.md](gclab-2026-10-18.md); scripts and
the opt-in hook in [tools/gclab/](../tools/gclab/). All runs stayed
byte-identical.

- **The GC's wall-time cost is small: 3–7%** (typed-app 7%, downlevel-simple
  5%, bundle and checks-only 3%, single files 2%), i.e. 7–22% of the gap to
  warm Java. The GC takes 22–29% of CPU, but it runs mostly on otherwise
  idle cores: mark assists ~0.01 s, pauses under 10 ms. The profile's 27%
  is CPU, not wall time; earlier GC/layout speedup estimates were too high.
- **GOGC=200** cuts GC CPU to 13–17%, **GOGC=400** to 7–10%, at 1.3–2.6×
  the RSS (still below Java's). The wall gain is bounded by the 3–7% tax.
  `GOGC=off` without a limit is slower here (page faults: bundle +86%).
- **Allocation itself costs 12–15% of CPU** on the main thread, and no GC
  setting removes it: allocation churn matters more than GC scanning.
- **Heap:** 80–96% of the live heap holds pointers. AST nodes are 37–66%
  of the live heap on typed-app and 76–89% on bundle; types are 21% at
  typed-app's peak. Churn: jutil collection entries 31–40%, pass-local
  scopes/references/graph nodes 17–19%, Var/Scope 9–11%;
  `flowSensitiveInlineVariables` alone is 21% of bundle's allocation.
- **Locality:** walking the real AST mid-optimization costs 26–33 ns per
  node against 9–15 ns for a contiguous pre-order copy (~2.5×).
- **Ranked candidates (estimates):** (1) AST in pointer-free arenas:
  −50 to −80% of GC mark work, wall −2 to −5%, up to −5% more with
  pre-order layout; (2) jutil hash maps/sets without per-entry objects:
  wall −3 to −6%, jutil only; (3) data-flow lattices as dense arrays;
  (4) arenas for scopes, vars and references; (5) types.

**What this established.** These GC/layout measurements did not support the
earlier 2–3× predictions. The remaining gap spans allocation, memory locality,
repeated analysis and pass algorithms. Larger changes require new measurements
rather than extrapolation from GC CPU percentages. Increasing GOGC also raises
memory use and remains an experiment, not the default recommendation.

**Current next checks:** finish differential and race tests for the CLI
prototypes, compare prebuilt candidates individually and in combination, and
measure stock warm Java versus Go at 1, 2 and 4 compiler workers under the same
CPU affinity/quota. Report elapsed-time and peak-RSS ratios together, with JVM
warm-up drift and measurement windows. The earlier plan to prototype AST arenas
via an overlay is historical; current retention decisions follow the measured
CLI report.

## Smaller open items

- `tools/clie2e/sweep.go` runs its JVMs under its own `nice -n 10`; under
  an outer `nice` they end at 19 and run ~5× slower when other work
  competes. Drop the inner `nice` or document it.
- The released jar still cannot load its own XTB reader
  (`NoClassDefFoundError …jarjar/javax/xml/parsers/ParserConfigurationException`);
  `tools/clie2e/build-java.sh` works around it for the oracles.
