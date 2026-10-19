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

| Task | Baseline (`main`) | Integrated (`t=1` / `t=4`) | Speedup (`t=1` / `t=4`) | vs Cold / Warm Java Time (`t=1` / `t=4`) | Peak RSS: Go (`main` → new) | Peak RSS: Cold / Warm Java | Java / Go RSS (Cold / Warm) |
|---|---:|---:|---:|---:|---:|---:|---:|
| `typed-app` | 2.86 s | **2.74 s** / **2.64 s** | **−4.0%** / **−7.7%** | 2.20× / 0.76× (0.79× at `t=4`) | 234 → **242 MiB** | 758 / 1,871 MiB | **3.13× / 7.72×** |
| `bundle` (12.6 MB) | 11.82 s | **10.31 s** / 10.31 s | **−12.8%** | **1.22×** / **0.85×** | 543 → **490 MiB** (−9.7%) | 2,151 / 2,343 MiB | **4.39× / 4.78×** |
| `downlevel-whitespace` | 0.86 s | **0.84 s** / **0.70 s** | **−2.5%** / **−18.2%** | 2.92× / 0.74× (0.88× at `t=4`) | 168 → **145 MiB** (−13.7%) | 553 / 1,437 MiB | **3.82× / 9.92×** |
| `downlevel-simple` | 2.39 s | **2.25 s** / **2.10 s** | **−5.9%** / **−12.2%** | 2.00× / 0.72× (0.77× at `t=4`) | 203 → **190 MiB** (−6.6%) | 613 / 1,489 MiB | **3.23× / 7.84×** |
| `checks-only` | 4.40 s | **4.29 s** / **3.82 s** | **−2.4%** / **−13.1%** | 1.57× / 0.76× (0.85× at `t=4`) | 693 → **618 MiB** (−10.7%) | 1,591 / 3,562 MiB | **2.57× / 5.76×** |

On a 2-vCPU Linux VM (AMD EPYC 9V74, JDK 21.0.12.1, 12 warm-up + 10 measured
iterations), the same changes reduce elapsed time by 8.8–19.8% across the five
tasks (`bundle` 23.10 s → 18.53 s, −19.8%, RSS 675 → 504 MiB, −25.3%) with
warm Java / Go peak RSS ratios of **4.74×** (`typed-app`), **3.33×** (`bundle`),
**5.72×** (`downlevel-whitespace`), **7.13×** (`downlevel-simple`), and
**2.19×** (`checks-only`).

Notably, on the 12.6 MB `bundle` task—previously the only whole-program task
where `jsgo` trailed cold Java (0.94×)—wall time dropped by **12.8%** (19.8% on
Linux) and peak RSS by **9.7%** (25.3% on Linux), putting `jsgo` ahead of cold
Java (**1.22×**) and closing the warm-JVM speed ratio from 0.67× to **0.85–0.89×**
at **3.3–4.8× lower peak RSS**.

## Idiomatic concurrency, hot-path passes, and PR #23 integration (2026-10-19)

Building on `123b46e8`, a second round of systematic concurrency, data-structure,
and pass-level optimizations was implemented and verified step by step:

1. **Idiomatic concurrency & lock elimination (`step1` & `step7`):**
   - Unified `prebuild_ast_parallel.go` and `prebuild_dependency_info_parallel.go`
     into a generic `parallelPrebuild[T, R]` worker pool using lock-free
     `atomic.Int64` work-claiming and capped at `min(numParallelThreads, len(inputs))`.
     This bounds compiler work items; it does not change `GOMAXPROCS` or cap all
     runtime goroutines and OS threads.
   - Sharded `RhinoStringPool` (`internal/rhino/rhino_string_pool.go`) into 64
     `sync.RWMutex` shards keyed by `StringHashCode(s)` with an allocation-free
     read lock fast path.
   - Replaced `jutil.Monitor` with `sync.Once` in `ParserRunner` /
     `ParserConfiguration`, `sync.Mutex` + $O(1)$ `entries map[string]*zip.File`
     in `JSCompZipFileCache`, lock-free `atomic.Pointer` / `atomic.Uint64` on
     `colors.Color` (shrinking `Color` from 112 B → 80 B), `sync.Mutex` +
     `atomic.Pointer[string]` on `SourceFile` (`GetCodeNoCache`), and
     `sync.Mutex` on `DiagnosticGroup_ForType`.
2. **Parallel `SourceInformationAnnotator` and per-script `CodePrinter` (`step2`):**
   - Replaced `NodeTraversal` + `LinkedIdentityHashMap` scope stack in
     `SourceInformationAnnotator` with a direct recursive AST walk (`AnnotateTree`)
     and parallelized per-root annotation across worker goroutines.
   - Parallelized per-script `toSourceAndMappings` in `Compiler.ToSourceArray` /
     `ToSource` when `numParallelThreads > 1`, assembling source strings and
     source-map mappings deterministically in script index order, and replaced
     the UTF-16 `StringCharAt` suffix check with `strings.HasSuffix(code, ";")`.
3. **Value-embedded collections & concrete scope/traversal stacks (`step3`):**
   - Replaced `scopes []any` in `NodeTraversal` with concrete `[]scopeStackEntry`
     backed by an inline `[16]scopeStackEntry` array; lazily initialized
     `SyntacticScopeCreator` only when a traversal requests a scope.
   - Embedded `jutil.LinkedHashMap` by value in `AbstractScope`,
     `LinkedIdentityHashMap`, and `LinkedIdentityHashSet`; added inline
     `refsBuf [4]*Reference` in `ReferenceCollection` and concrete `GetVar` /
     `GetSlot` methods on `*Scope`.
4. **Hot optimization pass improvements (`step4`, `step5`, `step8`):**
   - **`FlowSensitiveInlineVariables`:** filtered function-parameter definitions
     during candidate gathering, skipped `MustBeReachingVariableDef` when a
     function has no candidates, deferred backward `MaybeReachingVariableUse`
     until at least one candidate passes cheap AST/CFG/side-effect checks, and
     embedded `dependsSet` inline in `MustBeReachingVariableDef_Definition`.
   - **`RemoveUnusedCode` & `InlineAndCollapseProperties`:** guarded discarded
     `DebugLog` closures (`IsLogging()`), checked `varInfo.IsRemovable()` before
     allocating closures in `TraverseNode`, converted `Continuation` to a value
     struct with inline `continuationsBuf [2]Continuation` and
     `removablesBuf [2]Removable`, and cached parsed `js/polyfills.txt` via
     `sync.OnceValue` across `PolyfillUsageFinder`, `RemoveUnusedCode`,
     `RewritePolyfills`, and `IsolatePolyfills`.
   - **`InlineVariables`, `ControlFlowAnalysis`, `LiveVariablesAnalysis`, `MakeDeclaredNamesUnique`:**
     embedded all 9 `InitiallyUnknown[T]` fields by value on `IdentifyConstants`;
     replaced `ControlFlowAnalysis.astPosition` with `map[*rhino.Node]int`;
     added inline `buf [2]uint64` (128 bits) to `jutil.BitSet` and stack-allocated
     `gen`/`kill` bitsets in `LiveVariablesAnalysis`; skipped `disallowedNames`
     map allocation when `preferredName` has no conflict in `AttemptRename` and
     used `map[string][]VariableInfo` in `VariableInfoStackMap`.
5. **`PropertyMap`, `TreeMap`, and PGO with `bundle` (`step7` & `step9`):**
   - Replaced `sort.Search` closure in `TreeMap.search` with an inline 3-way
     binary search (`c == 0` early return).
   - Replaced `TreeMap` inside `rhino.PropertyMap` (`jstype_property_map.go`,
     embedded by value in `PrototypeObjectTypeBase`) with a flat `entries` slice,
     $O(1)$ `byName` map for $> 8$ properties, lazy `StringCompareTo` sorting on
     key-ordered iteration, and allocation-free `findClosestPairByKey` across
     the prototype chain.
   - Added `bundle` to default PGO profile collection (`tools/bench/pgo.py`) and
     regenerated `cmd/jsgo/default.pgo`.

### Measured results (interleaved, N=3 per step + final paired control)

| Task | Original (`c9872522`) | Phase 1 (`123b46e8` `t=1`/`t=4`) | Phase 2 (`t=1` / `t=4`) | Phase 2 vs Phase 1 (`t=1` / `t=4`) | Total vs Original (`t=1` / `t=4`) | Warm Java (`t=4`) | Go / Warm Java Speed (`t=1` / `t=4`) | Go RSS vs Cold / Warm Java RSS |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `typed-app` | 2.86 s | 2.75 s / 2.61 s | **2.53 s** / **2.41 s** | **−8.1%** / **−7.5%** | **−11.7%** / **−15.7%** | 2.18 s | **0.86×** / **0.90×** | 230 MiB (**3.21× / 8.06×** lower) |
| `bundle` (12.6 MB) | 11.82 s | 10.31 s | **9.40–9.44 s** | **−8.7% to −10.5%** | **−20.1%** | 8.68 s | **0.92×** / **0.92×** | 487 MiB (**3.99× / 4.84×** lower) |
| `downlevel-whitespace` | 0.86 s | 0.84 s / 0.70 s | **0.74 s** / **0.63 s** | **−12.6%** / **−10.7%** | **−14.3%** / **−27.0%** | 0.60 s | **0.81×** / **0.95×** | 150 MiB (**3.24× / 12.08×** lower) |
| `downlevel-simple` | 2.39 s | 2.24 s / 2.10 s | **2.02 s** / **1.90 s** | **−10.0%** / **−9.2%** | **−15.5%** / **−20.4%** | 1.60 s | **0.79×** / **0.84×** | 183 MiB (**3.22× / 7.76×** lower) |
| `checks-only` | 4.40 s | 4.28 s / 3.83 s | **3.96 s** / **3.48 s** | **−7.5%** / **−9.0%** | **−10.0%** / **−20.9%** | 3.20 s | **0.81×** / **0.92×** | 637 MiB (**2.45× / 4.84×** lower) |

## Smaller open items

- `tools/clie2e/sweep.go` runs its JVMs under its own `nice -n 10`; under
  an outer `nice` they end at 19 and run ~5× slower when other work
  competes. Drop the inner `nice` or document it.
- The released jar still cannot load its own XTB reader
  (`NoClassDefFoundError …jarjar/javax/xml/parsers/ParserConfigurationException`);
  `tools/clie2e/build-java.sh` works around it for the oracles.


