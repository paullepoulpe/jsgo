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

## GC lab (2026-10-18): measured, and it corrects the estimates above

Full write-up: [gclab-2026-10-18.md](gclab-2026-10-18.md); scripts and
the opt-in hook in [tools/gclab/](../tools/gclab/). All runs stayed
byte-identical.

- **The GC's wall-time cost is small: 3–7%** (typed-app 7%, downlevel-simple
  5%, bundle and checks-only 3%, single files 2%), i.e. 7–22% of the gap to
  warm Java. The GC takes 22–29% of CPU, but it runs mostly on otherwise
  idle cores: mark assists ~0.01 s, pauses under 10 ms. The profile's 27%
  is CPU, not wall time; the estimates in options 1 and 3 above were too
  high.
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

**Revised view.** No single substrate change is worth 2–3×; each is worth
single-digit percent. The gap to warm Java is spread over allocation
cost, memory layout and the passes' own algorithms. Reaching and passing
warm Java will take the combination (allocation-free collections,
arena-allocated AST in traversal order, cheaper per-pass data
structures, and per-pass work on the slowest passes), measured step by
step; the "2–3×" estimate in option 3 is not supported by the data.
Raising GOGC to 200–400 is a free 2–5% in the meantime, if the owner
wants it.

**Next measurements:** a prototype of the AST arena (noscan slab for
`rhino.Node`, via an overlay) over 10 paired runs; hand-mapping the
10–15% of allocations left as "other"; heap profiles for
downlevel-simple and checks-only.

## Allocation, compact layout, and `--num_parallel_threads` (2026-10-18)

Following the GC lab findings, a set of targeted allocation, struct-size,
data-flow lattice, and opt-in parallel prebuild changes were integrated:

- **Slab-batched collection entries (`internal/jutil/collections.go`):**
  `LinkedHashMap` and `ObjectLinkedHashMap` batch entry nodes in 64-entry
  slabs once a map grows past 8 entries, plus single-probe `LoadOrStore`.
- **Indexed 64-var bitmask lattice (`internal/jscomp/must_def_map.go`):**
  `MustBeReachingVariableDef` uses a dense indexed table and uint64 active
  bitset for functions with $\le 64$ tracked variables, falling back to
  `ObjectLinkedHashMap` only for larger functions.
- **Compact `rhino.Node` (96 B → 80 B) and `jscomp.Var` (80 B → 64 B):**
  pack `str` (`*string`) and `typeOrColor` (`unsafe.Pointer` with a 2-word
  tagged wrapper when both JSType and Color are present on the same node) in
  `rhino.Node`; pack boolean flags and lazy atomic `nameHash` on `AbstractVar`.
- **Scanner & rename-conflict deduplication:** reuse consecutive identical
  token `SourcePosition`s in `LineNumberScanner` and deduplicate conflict
  recording in `MakeDeclaredNamesUnique` for scopes with $\ge 8$ declarations.
- **Opt-in `--num_parallel_threads` prebuild:** bounded parallel AST parsing
  (`prebuild_ast_parallel.go`) and dependency scanning
  (`prebuild_dependency_info_parallel.go`) with deterministic pre-assigned
  diagnostics and synchronized ZIP entry caching (`js_comp_zip_file_cache.go`).
  Default sequential behavior (`--num_parallel_threads=1`) is unchanged.
- **Unicode 15.0 delta wrapper (`internal/jutil/unicode15.go`):** lightweight
  delta table over Go's `unicode` package so Go 1.27+ (Unicode 16+) remains
  100% byte-identical with JDK 21 (Unicode 15.0) without vendoring Go 1.24's
  `unicode` package.

### Measured results (interleaved, N=5, PGO enabled)

Measured on macOS arm64 (Go 1.27.1, `-pgo=cmd/jsgo/default.pgo`), interleaving
`main` (`c9872522`) and the updated branch (`N=5` cold Go runs each, `N=3`
cold Java, `N=5` warm Java, all 5 byte-equivalence gates `MATCH`, plus 48/48
`-race` parallel equivalence gate cases across `--num_parallel_threads={1,2,4}`):

| Task | Baseline (`main`) | Integrated (`threads=1`) | Speedup (`t=1`) | Peak RSS (`main` → new) | With `--num_parallel_threads=4` | vs Cold Java | vs Warm Java (`t=1` / `t=4`) |
|---|---:|---:|---:|---:|---:|---:|---:|
| `typed-app` | 2.86 s | **2.74 s** | **−4.0%** | 234 → 242 MiB | **2.64 s** (−7.7% vs `main`) | 2.20× (6.02 s) | 0.76× / 0.79× (2.07 s) |
| `bundle` (12.6 MB) | 11.82 s | **10.31 s** | **−12.8%** | 543 → **490 MiB** (−9.7%) | 10.31 s (single file) | **1.22×** (12.57 s) | **0.85×** (8.76 s) |
| `downlevel-whitespace` | 0.86 s | **0.84 s** | **−2.5%** | 168 → **145 MiB** (−13.7%) | **0.70 s** (−18.2% vs `main`) | 2.92× (2.44 s) | 0.74× / 0.88× (0.62 s) |
| `downlevel-simple` | 2.39 s | **2.25 s** | **−5.9%** | 203 → **190 MiB** (−6.6%) | **2.10 s** (−12.2% vs `main`) | 2.00× (4.51 s) | 0.72× / 0.77× (1.62 s) |
| `checks-only` | 4.40 s | **4.29 s** | **−2.4%** | 693 → **618 MiB** (−10.7%) | **3.82 s** (−13.1% vs `main`) | 1.57× (6.75 s) | 0.76× / 0.85× (3.25 s) |

Notably, on the 12.6 MB `bundle` task—previously the only whole-program task
where `jsgo` trailed cold Java (0.94×)—wall time dropped by **12.8%** and peak
RSS by **9.7%**, putting `jsgo` ahead of cold Java (**1.22×**) and closing the
warm-JVM ratio from 0.67× to **0.85×**.

## Smaller open items

- `tools/clie2e/sweep.go` runs its JVMs under its own `nice -n 10`; under
  an outer `nice` they end at 19 and run ~5× slower when other work
  competes. Drop the inner `nice` or document it.
- The released jar still cannot load its own XTB reader
  (`NoClassDefFoundError …jarjar/javax/xml/parsers/ParserConfigurationException`);
  `tools/clie2e/build-java.sh` works around it for the oracles.

