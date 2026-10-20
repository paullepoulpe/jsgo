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
- **Speed** ([benchmarks.md](benchmarks.md#end-to-end-v20261015) and Rounds 3–5 below):
  after five rounds of systematic substrate, data-structure, and pass-level
  optimizations, `jsgo` is **1.59×–5.10× faster than cold Java (`java -jar`)**
  and **1.20×–1.64× faster than a warm JIT-compiled JVM** across all 11
  whole-program benchmarks (including 4,385-file / 44.6 MB `typed-mega` and
  26-library / 36.1 MB `bundle-mega` `ADVANCED` builds), while using
  **2.7×–10.0× less peak RSS** than warm Java. Even single-threaded (`t=1`),
  `jsgo` matches or beats 4-thread warm Java (`1.00×–1.60×`) across all 11
  benchmarks.
- **Historical starting point (initial v20261015 port, before Rounds 1–3)** (per-phase tables in
  [benchmarks-e2e-v20261015.md](benchmarks-e2e-v20261015.md)): on the initial
  1:1 port, `jsgo` was 0.65–0.74× of warm Java's speed, trailing in the
  optimization loop (`flowSensitiveInlineVariables` 2.4–2.8× Java's time,
  `inlineVariables` ~2×, `removeUnusedCode` 1.5–1.9×,
  `collapseObjectLiterals` and `coalesceVariableNames` ~2×), parsing
  (155 ms per single-file compile against 61 ms, mostly the default
  externs), `inferTypes` (1.4–1.8×) and `checkAccessControls` (~2.1×).
  The sections below trace how each of those bottlenecks was measured and
  eliminated.

## Why: the CPU profile (initial port)

The initial profile behind `cmd/jsgo/default.pgo` (typed-app and friends, v20261015
port) attributed:

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

## Round 1: Allocation, compact layout, and `--num_parallel_threads` (2026-10-18)

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

## Round 2: Idiomatic concurrency, hot-path passes, and PR #23 integration (2026-10-19)

Building on `123b46e8`, a second round of systematic concurrency, data-structure,
and pass-level optimizations was implemented and verified step by step:

1. **Idiomatic concurrency & lock elimination (`step1` & `step7`):**
   - Unified `prebuild_ast_parallel.go` and `prebuild_dependency_info_parallel.go`
     into a generic `parallelPrebuild[T, R]` worker pool using lock-free
     `atomic.Int64` work-claiming and capped at `min(numParallelThreads, runtime.GOMAXPROCS(0), len(inputs))`.
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

## Round 3: Exhaustive substrate, pass, scanner, slab, and large-scale `ADVANCED` optimizations (2026-10-19)

Building on Round 2 (`ff06ac15`), Round 3 (`perf/round3-exhaustive`, Steps 10–18)
addressed the remaining CPU, allocation, and cache-locality bottlenecks across
the AST, scanner, collection substrate, CFG/dataflow engines, type system,
optimization passes, and source-map/error-reporting pipeline—and expanded the
benchmark suite with four large real-world `ADVANCED` compilation workloads:

1. **AST node accessors, traversal, and compact `BitSet` (`Step 10`):**
   - Outlined cold mutation paths and made `rhino.Node` child/property/flag
     accessors inlineable (`FirstChild`, `GetNext`, `GetSecondChild`, `GetLastChild`,
     `HasOneChild`, `HasTwoChildren`, `GetBooleanProp`, `GetProp`, `lookupProperty`).
   - Replaced `BitSet` dynamic backing slice on the small fast path (`len <= 128`
     bits) with an inline `buf [2]uint64` + `words []uint64` (shrinking `BitSet`
     from 40 B → 24 B and eliminating slice header indirection for $\le 128$ bits).
   - Fast-pathed `NodeTraversal` and `CheckSideEffects` AST walks.
2. **CodePrinter, `RemoveUnusedCode`, `GatherModuleMetadata`, and `JSTypeRegistry` (`Step 11`):**
   - Eliminated transient buffers/closures in `CodePrinter`, `RemoveUnusedCode`,
     `GatherModuleMetadata`, and `JSTypeRegistry` property/type lookup paths.
3. **`InlineFunctions`, `DefinitionUseSiteFinder`, `Peephole`, and `DevirtualizeMethods` (`Step 12`):**
   - Eliminated per-call-site map/slice churn in `InlineFunctions` and
     `FunctionInjector`; optimized `DefinitionUseSiteFinder` and
     `DevirtualizeMethods` collection lookups and peephole folding hot paths.
4. **`TypeInference`, `QualifiedName`, and `RenameVars` allocation elimination (`Step 13`):**
   - Eliminated transient `BooleanOutcomePair` and flow-scope wrapper allocations
     in `TypeInference`, optimized `QualifiedName.Matches`, and streamlined
     `RenameVars` assignment sorting and lookup.
5. **Scanner ASCII fast-path, keyword switch, and FNV-1a open-addressing `RhinoStringPool` (`Step 14`):**
   - Added pure-ASCII batch identifier scanning in `Scanner` (`scanIdentifierOrKeyword`)
     bypassing per-character `PeekChar`/`Advance` and `strings.Builder`.
   - Replaced the 47-case string `switch` in `stringToKeyword` with a
     length + first-byte switch.
   - Replaced Go maps inside `RhinoStringPool` with a custom open-addressing
     `[]*string` hash table using inline FNV-1a hashing (`fnvHash`) and
     `sync.RWMutex` read-lock fast paths.
6. **AST `Node` slab allocation in parse order & flat JSDoc comment slice (`Step 15`):**
   - Added 256-node contiguous `nodeSlab` bump allocation in `IRFactory` and
     `ParserRunner` (`ir_factory.go`), placing sibling and child AST nodes coded
     adjacent in source into contiguous cache lines in parse/pre-order sequence.
   - Replaced `TreeMap[int, *Comment]` in `IRFactory` with a flat sorted
     `[]*Comment` slice and binary search.
7. **Flat-slot `LinkedHashMap` / `ObjectLinkedHashMap` & compact `Var` / `Scope` (`Step 16`):**
   - Re-architected `jutil.LinkedHashMap` and `jutil.ObjectLinkedHashMap` to
     store entries in a single contiguous `entries []flatEntry` slice with
     `int32` hash-chain and insertion-order links plus $O(1)$ `first` live-index
     tracking and inlineable value iterators (`KeyIterator`, `ValueIterator`,
     `EntryIterator`).
   - Shrunk `Var` (64 B → 48 B), `TypedVar` (96 B → 64 B), `Scope` (96 B → 80 B),
     and `TypedScope` (112 B → 96 B) by removing redundant parent pointers,
     packing indices (`index int32`), and replacing `ImplicitGoogNamespace`
     per-scope hash maps with a 6-bit mask.
8. **Cached `HashCode` on hot `ObjectLinkedHashMap` keys & inlined `jutil` preconditions (`Step 17a`):**
   - Added cached `hash int32` (computed via `jutil.StringHashCode`) to
     `rhino.InputId`, `AbstractVarBase`, `ScopedName_Simple` (shrunk 40 B → 32 B),
     `Module`, `Binding`, and `Export`, eliminating ~96% of `ReflectHashCode`
     `reflect.ValueOf` allocations in `ObjectLinkedHashMap`.
   - Rewrote `jutil.IsNil` to inspect the 2-word `eface` header directly without
     `reflect.ValueOf`, and split `CheckNotNull`, `CheckElementIndex`,
     `CheckPositionIndex`, `CheckState`, `CheckArgument`, and `ObjectsEqualsOf`
     into inlineable fast paths with `@noinline` panic helpers.
9. **Type equality/subtyping fast paths, coding conventions, and zero-alloc `NodeTraversal` callbacks (`Step 17b`):**
   - Added `if t.self == that { return false }` identity check to
     `JSTypeBase.DiffersFrom`, stack-allocated `EqualityChecker` /
     `SubtypeChecker` / `UnionType_Builder` (`inlineAlts [8]JSType`), and
     short-circuited 2-element `CreateUnionType`.
   - Eliminated per-traversal `abstractPostOrderCallbackAdapter` heap allocation
     via `NodeTraversal_AbstractPostOrderCallbackInterface`.
   - Fixed `GoogleCodingConvention.IsClassFactoryCall` (which previously called
     `callTarget.String()`—allocating a debug `strings.Builder`—on every single
     `CALL` node in the compiler) and cached string pointers in
     `ClosureCodingConvention`.
10. **CFG & dataflow graph/lattice slab allocation (`Step 17c`):**
    - Added slab-allocated graph nodes/edges (`LinkedDirectedGraph`,
      `LinkedUndirectedGraph`), stack-backed worklists/heaps in
      `GraphReachability`, `ControlFlowAnalysis`, and `DataFlowAnalysis`, and
      per-pass slab allocators + copy-on-write lattices in
      `LiveVariablesAnalysis`, `MustBeReachingVariableDef`,
      `MaybeReachingVariableUse`, `CoalesceVariableNames`,
      `DeadAssignmentsElimination`, and `FlowSensitiveInlineVariables`.
11. **`SourceMapInput`, `Monitor`, and `SourceMapGeneratorV3` (`bundle-advanced-ts` discovery):**
    - Adding the 12.6 MB `bundle-advanced-ts` benchmark (`typescript.js` in
      `ADVANCED` mode, which emits 10,751 warnings and 641,000 source-map
      mappings) exposed two major bottlenecks where `jsgo` initially took 43.3 s
      vs Warm Java's 12.4 s:
      1. Every one of the 10,751 warnings called `Compiler.GetSourceMapping` →
         `SourceMapInput.GetSourceMap`, which re-attempted to read the missing
         `lib.d.ts.map` file from disk **32,253 times** and formatted a
         `JSError` warning (`SourceMap.Format` via `regexp` + `fmt.Sprintf`)
         every single time before `ComposeWarningsGuard` discarded it. Adding an
         `atomic.Bool` `ready` fast path to `SourceMapInput.GetSourceMap` and a
         `c.inputSourceMaps.IsEmpty()` fast path in `Compiler.GetSourceMapping`
         reduced disk attempts from 32,253 to 1.
      2. `ThreadSafeDelegatingErrorManager` acquired `jutil.Monitor.Enter()` on
         every warning, which called `goroutineID()` → `runtime.Stack(buf[:], false)`,
         unwinding the call stack 10,751 times (~20 seconds of CPU). Adding
         non-reentrant `LockFast()` / `UnlockFast()` to `jutil.Monitor` for
         leaf critical sections eliminated `runtime.Stack` entirely.
      3. Converting `FilePosition` from pointer fields to value `line, column int32`
         and slab-allocating `SourceMapGeneratorV3_Mapping` in 256-entry slabs
         eliminated ~1.9M small allocations during source-map generation.
    - Together, these fixes cut `bundle-advanced-ts` wall time from **43.29 s to
      12.09 s** (**3.58× faster**, beating Warm Java's 12.41 s).

### Final measured results — Standard 5 Benchmarks (interleaved, N=3, PGO enabled)

| Benchmark | Original (`c9872522`) | Round 2 (`ff06ac15` `t=1`/`t=4`) | **Round 3 (`f0303218` `t=1`/`t=4`)** | Round 3 vs Round 2 (`t=1`/`t=4`) | Total vs Original (`t=1`/`t=4`) | Cold / Warm Java (`t=4`) | **`jsgo` vs Cold / Warm Java (`t=4`)** | Peak RSS: `jsgo` vs Cold / Warm Java |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| **`typed-app`** | 2.86 s | 2.55 s / 2.41 s | **2.26 s** / **2.11 s** | **−11.2%** / **−12.6%** | **−21.0%** / **−26.4%** | 6.14 s / 2.18 s | **2.92×** / **1.04×** | **243 MiB** vs 758 / 1,871 MiB (**3.1× / 7.7× lower**) |
| **`bundle`** (12.6 MB) | 11.82 s | 9.29 s / 9.07 s | **7.90 s** / **7.76 s** | **−15.0%** / **−14.4%** | **−33.2%** / **−34.3%** | 16.40 s / 8.68 s | **2.11×** / **1.12×** | **522 MiB** vs 2,151 / 2,343 MiB (**4.1× / 4.5× lower**) |
| **`downlevel-whitespace`** | 0.86 s | 0.74 s / 0.64 s | **0.66 s** / **0.55 s** | **−11.3%** / **−13.6%** | **−23.3%** / **−36.0%** | 2.55 s / 0.60 s | **4.64×** / **1.09×** | **149 MiB** vs 553 / 1,437 MiB (**3.7× / 9.6× lower**) |
| **`downlevel-simple`** | 2.39 s | 2.03 s / 1.90 s | **1.71 s** / **1.60 s** | **−15.8%** / **−15.9%** | **−28.6%** / **−33.1%** | 4.68 s / 1.60 s | **2.93×** / **1.00×** | **197 MiB** vs 613 / 1,489 MiB (**3.1× / 7.6× lower**) |
| **`checks-only`** | 4.40 s | 3.95 s / 3.46 s | **3.47 s** / **3.08 s** | **−12.1%** / **−10.9%** | **−21.1%** / **−29.9%** | 6.55 s / 3.20 s | **2.12×** / **1.04×** | **643 MiB** vs 1,591 / 3,562 MiB (**2.5× / 5.5× lower**) |

### Final measured results — Large Real-World `ADVANCED` Benchmarks (interleaved, N=3, `t=4`)

To test `jsgo` on large real-world `ADVANCED` compilations, four new benchmarks
were added to [`tools/bench/tasks/`](../tools/bench/tasks/) using **Google Shaka
Player 5.3.1** (522 files), **Google Blockly 8.0.5** (324 files), **Closure
Library v20230802** (971 files), and large multi-library UMD bundles. Every task
produces **100% byte-identical `out.js`, `out.js.map`, stdout, stderr, and exit
code** to Java v20261015:

| Benchmark | Workload (`--compilation_level=ADVANCED`) | Round 2 (`ff06ac15` `t=1`/`t=4`) | **Round 3 (`f0303218` `t=1`/`t=4`)** | Round 3 Speedup (`t=4`) | Cold Java | Warm Java | **`jsgo` vs Cold Java (`t=1`/`t=4`)** | **`jsgo` vs Warm Java (`t=1`/`t=4`)** | Peak RSS: `jsgo` vs Cold / Warm Java |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| **`typed-full`** | **28.5 MB, 1,818 files** (Closure Library + Shaka Player 5.3.1 + Blockly 8.0.5 + app, 97.3% typed → 6.75 MB JS + 13.31 MB map) | 21.89 s / 19.90 s | **17.90 s** / **17.17 s** | **−13.7%** | 31.60 s | 20.21 s | **1.77×** / **1.84× faster** | **1.13×** / **1.18× faster** (−3.04 s) | **1,158–1,189 MiB** vs 1,603 / 4,724 MiB (**1.3× / 4.0× lower**) |
| **`shaka-blockly`** | **8.6 MB, 847 files** (Shaka Player 5.3.1 + Blockly 8.0.5, 93.9% typed → 1.86 MB JS + 4.94 MB map) | 9.16 s / 8.61 s | **7.57 s** / **7.44 s** | **−13.6%** | 15.26 s | 8.67 s | **2.02×** / **2.05× faster** | **1.14×** / **1.17× faster** (−1.23 s) | **513–552 MiB** vs 1,341 / 2,852 MiB (**2.4× / 5.2× lower**) |
| **`bundle-advanced`** | **5.4 MB, 11 UMD libraries** (three, lodash, react-dom, rxjs, d3, moment, vue, jquery, handlebars, immutable, underscore → 1.88 MB JS + 4.35 MB map) | 9.42 s / 8.67 s | **6.20 s** / **6.05 s** | **−30.2%** | 11.30 s | 6.79 s | **1.82×** / **1.87× faster** | **1.09×** / **1.12× faster** (−0.74 s) | **412–417 MiB** vs 1,600 / 3,795 MiB (**3.8× / 9.1× lower**) |
| **`bundle-advanced-ts`** | **12.6 MB, 4 libraries incl. 8.9 MB `typescript.js`** (10,751 warnings → 3.25 MB JS + 9.39 MB map) | — / 43.29 s | **12.20 s** / **12.09 s** | **−72.1% (3.58×)** | 20.78 s | 12.41 s | **1.70×** / **1.72× faster** | **1.02×** / **1.03× faster** (−0.32 s) | **668–670 MiB** vs 2,422 / 2,850+ MiB (**3.6× / 4.3× lower**) |

## Round 4: Mega-Benchmarks (`typed-mega` 44.6 MB & `bundle-mega` 36.1 MB) and Deep Pass/SourceMap/Slab Optimizations (2026-10-20)

To stress-test `jsgo` at even larger real-world scale and uncover additional
bottlenecks, Round 4 added two mega-benchmarks to [`tools/bench/tasks/`](../tools/bench/tasks/)
and executed a 6-agent parallel optimization swarm across the optimization loop,
type inference, reference collection, AST/scope hot paths, and source-map/diagnostic
formatting:

1. **Two New Mega-Benchmarks (`typed-mega` and `bundle-mega`):**
   - **`typed-mega`** (`tools/bench/tasks/typed-mega.json`): **4,385 files, 44.62 MB (42.56 MiB), 94.1% typed** across 10 real-world Closure `ADVANCED` codebases: Closure Library `v20230802` (971 files), NGA OpenSphere `release-23.0.1` (1,648 files), OpenLayers `4.6.5` (539 files), Google Shaka Player `5.3.1` (520 files), Google Blockly `8.0.5` (446 files), Camptocamp `ngeo@2.2.10` (203 files), YouTube SPFjs `2.4.0`, `ol-cesium@1.37.0`, `google-protobuf@3.21.2`, and `app/` → produces **14.33 MB `out.js` + 22.29 MB `out.js.map` + 3.15 MB `stderr`**.
   - **`bundle-mega`** (`tools/bench/tasks/bundle-mega.json`): **26 UMD/browser libraries, 36.14 MB (34.47 MiB)** (`plotly.js`, `handsontable`, `echarts`, `ag-grid-community`, `pdfjs-dist`, `three.js`, `react-dom`, `pixi.js`, `fabric`, `mathjs`, `konva`, `moment`, `d3`, `google-libphonenumber`, `lodash`, `lovefield`, `vue`, `jquery`, `rxjs`, `google-protobuf`, `handlebars`, `immutable`, `chart.js`, `zone.js`, `underscore`) in `ADVANCED` mode → produces **13.09 MB `out.js` + 33.99 MB `out.js.map` + 14.01 MB `stderr` (27,597 warnings)**.
2. **`RemoveUnusedCode`, `GlobalNamespace`, `ProcessDefines`, and `InlineAndCollapseProperties` (~970 MB allocation reduction):**
   - Replaced pointer `*VarInfo` polymorphism in `RemoveUnusedCode` (`remove_unused_code.go`) with a single concrete 48-byte `VarInfo` struct allocated from 256-entry slabs and stored directly in a native `varInfoMap map[*Var]*VarInfo`, plus native `worklistSet map[*VarInfo]struct{}` and `removableSet map[*rhino.Node]struct{}`.
   - Added 256-entry `refSlab` / `nameSlab` bump allocators in `GlobalNamespace` (`global_namespace.go`), replaced `GlobalNamespace_Name.AliasingGets` with `int32` and `docInfo` with a bitflag, pre-sized `globalNames`, and reused a single `NodeFilter` closure across `Ref_CreateRefForTesting`.
3. **`ReferenceCollector`, `BasicBlock`, `InlineVariables`, and `InlineObjectLiterals`:**
   - Added 256-entry `Reference` and `BasicBlock` slab allocators in `ReferenceCollectingCallback` (`reference_collector.go`) with an inline `blockStackBuf [16]*BasicBlock` stack, inline `provable` slice on `BasicBlock` (`basic_block.go`), and `NativeVarMap` in `ReferenceMap`.
   - Replaced `LinkedIdentityHashMap` / `LinkedIdentityHashSet` in `InlineVariables` (`inline_variables.go`) and `InlineObjectLiterals` (`inline_object_literals.go`) with native Go maps where iteration order is not observable.
4. **Type System & Property Disambiguation (`TypeInference`, `LinkedFlowScope`, `TypeCheck`, `TypedScopeCreator`, `DisambiguateProperties`):**
   - Added 256-entry `LinkedFlowSlot` slab allocation (`flowSlotSlab`), pooled `FlatFlowScopeCache`, and embedded `LinkedHashMap` by value in `FlatFlowScopeCache` (`linked_flow_scope.go`), eliminating ~400 MB of flow-scope churn.
   - Pre-bound `SyntacticScopeCreator` callback in `TypeInferencePass` (`type_inference_pass.go`), added zero-alloc `GetJSTypeBeforeCast` and fast paths in `TypeCheck` (`type_check.go`) and `TypedScopeCreator` (`typed_scope_creator.go`), and added slab-allocated union-find elements + pooled `Proprec` nodes in `DisambiguateProperties` (`disambiguate_ambiguate_properties.go`).
5. **Inlining, Peephole, and Call Optimization Passes (`InlineSimpleMethods`, `InlineProperties`, `Peephole`, `InlineFunctions`, `OptimizeCalls`, `PureFunctionIdentifier`):**
   - Replaced `ArrayListMultimap` in `InlineSimpleMethods` with a flat `map[string][]*rhino.Node`, replaced `Objects_Equals` reflection in `PeepholeRemoveDeadCode` and `PeepholeMinimizeConditions` with direct pointer comparisons, eliminated per-call-site slice/map allocations in `InlineFunctions` / `FunctionInjector` / `FunctionArgumentInjector`, and optimized `OptimizeCalls` and `PureFunctionIdentifier`.
6. **AST Traversal, `Scope`, `GuardedCallback`, and SourceMap / Diagnostic Formatting:**
   - Added 64-entry `Var` and `TypedVar` slabs on `Scope` and `TypedScope` (`scope.go`), a zero-alloc `nativeSeenSet` in `SyntacticScopeCreator` (`syntactic_scope_creator.go`), a native `map[string]bool` in `GuardedCallback` (`guarded_callback.go`), a scratch node pool in `MinimizedCondition` (`minimized_condition.go`), and direct `TraverseWithScope` walks across `NodeUtil` (`node_util.go`).
   - Replaced per-character `Append([]byte{...})` in `Base64VLQ_Encode` (`base64_vlq.go`) and `SourceMapGeneratorV3_AddNameMap` (`source_map_generator_v3.go`) with a 4 KB stack buffer and an ASCII fast path in `Util_EscapeString` (`util.go`), cached `lastSymbolName` in `GetNameId`, increased mapping slabs to 4,096 entries, and added an ASCII fast path to `LightweightMessageFormatter.PadLine` (`lightweight_message_formatter.go`).
   - Regenerated `cmd/jsgo/default.pgo` from 15 profiles including `typed-full` and `bundle-advanced-ts`.

### Final measured results — Round 4 (`t=1` and `t=4`, interleaved, N=3, PGO enabled)

#### Standard 5 Benchmarks

| Benchmark | Original (`c9872522`) | Round 3 (`f0303218` `t=1`/`t=4`) | **Round 4 (`t=1` / `t=4`)** | Round 4 vs Round 3 (`t=1`/`t=4`) | Total vs Original (`t=1`/`t=4`) | Cold / Warm Java (`t=4`) | **`jsgo` vs Cold / Warm Java (`t=4`)** | Peak RSS: `jsgo` vs Cold / Warm Java |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| **`typed-app`** | 2.86 s | 2.22 s / 2.10 s | **1.93 s** / **1.82 s** | **−13.3%** / **−13.5%** | **−32.5%** / **−36.4%** | 6.14 s / 2.18 s | **3.37×** / **1.20×** | **253–269 MiB** vs 758 / 1,871 MiB (**2.8× / 7.0× lower**) |
| **`bundle`** (12.6 MB) | 11.82 s | 7.89 s / 7.77 s | **6.74 s** / **6.60 s** | **−14.6%** / **−15.0%** | **−43.0%** / **−44.2%** | 16.40 s / 8.68 s | **2.48×** / **1.32×** | **542 MiB** vs 2,151 / 2,343 MiB (**4.0× / 4.3× lower**) |
| **`downlevel-whitespace`** | 0.86 s | 0.64 s / 0.55 s | **0.61 s** / **0.52 s** | **−4.7%** / **−5.9%** | **−29.1%** / **−39.5%** | 2.55 s / 0.60 s | **4.90×** / **1.15×** | **133–149 MiB** vs 553 / 1,437 MiB (**3.7× / 9.6× lower**) |
| **`downlevel-simple`** | 2.39 s | 1.69 s / 1.60 s | **1.52 s** / **1.43 s** | **−9.7%** / **−10.6%** | **−36.4%** / **−40.2%** | 4.68 s / 1.60 s | **3.27×** / **1.12×** | **199–201 MiB** vs 613 / 1,489 MiB (**3.1× / 7.4× lower**) |
| **`checks-only`** | 4.40 s | 3.48 s / 3.07 s | **3.21 s** / **2.82 s** | **−7.6%** / **−8.1%** | **−27.0%** / **−35.9%** | 6.55 s / 3.20 s | **2.32×** / **1.13×** | **642–685 MiB** vs 1,591 / 3,562 MiB (**2.3× / 5.2× lower**) |

#### Large & Mega Real-World `ADVANCED` Benchmarks

| Benchmark | Workload (`--compilation_level=ADVANCED`) | Round 3 (`f0303218` `t=1`/`t=4`) | **Round 4 (`t=1` / `t=4`)** | Round 4 vs Round 3 (`t=1`/`t=4`) | Cold Java | Warm Java | **`jsgo` vs Cold Java (`t=1`/`t=4`)** | **`jsgo` vs Warm Java (`t=1`/`t=4`)** | Peak RSS: `jsgo` vs Cold / Warm Java |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| **`typed-mega`** | **44.6 MB, 4,385 files** (10 Closure `ADVANCED` codebases, 94.1% typed → 14.33 MB JS + 22.29 MB map) | 30.61 s / 29.38 s | **27.55 s** / **26.14 s** | **−10.0%** / **−11.1%** | 38.18 s | 29.53 s | **1.39×** / **1.46× faster** | **1.07×** / **1.13× faster** (−3.39 s) | **2,091–2,563 MiB** vs 6,608 / 6,881 MiB (**2.6× / 2.7× lower**) |
| **`bundle-mega`** | **36.1 MB, 26 UMD libraries** (27,597 warnings → 13.09 MB JS + 33.99 MB map) | 42.59 s / 41.59 s | **37.68 s** / **36.35 s** | **−11.5%** / **−12.6%** | 66.71 s | 53.60 s | **1.77×** / **1.84× faster** | **1.42×** / **1.47× faster** (−17.25 s) | **2,442–2,465 MiB** vs 4,421 / 6,775 MiB (**1.8× / 2.75× lower**, −622 MiB vs R3) |
| **`typed-full`** | **28.5 MB, 1,818 files** (Closure Library + Shaka Player 5.3.1 + Blockly 8.0.5 + app, 97.3% typed → 6.75 MB JS + 13.31 MB map) | 17.97 s / 17.19 s | **15.91 s** / **15.21 s** | **−11.5%** / **−11.5%** | 31.60 s | 20.21 s | **1.99×** / **2.08× faster** | **1.27×** / **1.33× faster** (−5.00 s) | **1,411–1,418 MiB** vs 1,603 / 4,724 MiB (**1.1× / 3.3× lower**) |
| **`bundle-advanced-ts`** | **12.6 MB, 4 libraries incl. 8.9 MB `typescript.js`** (10,751 warnings → 3.25 MB JS + 9.39 MB map) | 12.21 s / 12.11 s | **10.22 s** / **10.07 s** | **−16.3%** / **−16.8%** | 20.78 s | 12.41 s | **2.03×** / **2.06× faster** | **1.21×** / **1.23× faster** (−2.34 s) | **716–717 MiB** vs 2,422 / 2,850+ MiB (**3.4× / 4.0× lower**) |
| **`shaka-blockly`** | **8.6 MB, 847 files** (Shaka Player 5.3.1 + Blockly 8.0.5, 93.9% typed → 1.86 MB JS + 4.94 MB map) | 7.66 s / 7.44 s | **6.76 s** / **6.52 s** | **−11.8%** / **−12.4%** | 15.26 s | 8.67 s | **2.26×** / **2.34× faster** | **1.28×** / **1.33× faster** (−2.15 s) | **497–569 MiB** vs 1,341 / 2,852 MiB (**2.4× / 5.0× lower**) |
| **`bundle-advanced`** | **5.4 MB, 11 UMD libraries** (three, lodash, react-dom, rxjs, d3, moment, vue, jquery, handlebars, immutable, underscore → 1.88 MB JS + 4.35 MB map) | 6.25 s / 6.07 s | **5.21 s** / **5.03 s** | **−16.7%** / **−17.2%** | 11.30 s | 6.79 s | **2.17×** / **2.25× faster** | **1.30×** / **1.35× faster** (−1.76 s) | **363–389 MiB** vs 1,600 / 3,795 MiB (**4.1× / 9.8× lower**) |

## Round 5: `typed-mega` & `bundle-mega` Profile-Driven Optimizations (2026-10-20)

Profiling Round 4 on the new 44.6 MB `typed-mega` (4,385 files) and 36.1 MB
`bundle-mega` (26 UMD libraries) benchmarks exposed five major allocation and
CPU hotspots across code printing, reference collection, scope slot allocation,
CFG/interference graphs, type/color serialization, and variable normalization/renaming:

1. **`CodePrinter`, `CodeGenerator`, and `SourceMap` (`Subagent 5A`):**
   - **Progressive Mapping Slabs & Inline Buffers (`code_printer.go`):** `ToSourceArray` creates a separate `CodePrinter` for each of the 4,385 `SCRIPT` nodes in `typed-mega`; previously, `allocMapping` allocated a 4,096-entry `SourceMap_Mapping` slab (`98 KB`) on the first `SCRIPT` node even for the ~3,500 empty/pruned scripts (`3.20 GB` / `20.5%` of total allocations on `typed-mega`). Adding inline `mappingBuf [16]SourceMap_Mapping` + progressive slab growth (`64 -> 256 -> 1024 -> 4096`) and an inline `codeASCIIBuf [128]byte` / `codeASCII []byte` fast path cut `typed-mega` peak RSS at `t=4` by **662 MiB (−25.3%)** (`2,612 MiB → 1,950 MiB`).
   - **Zero-Escape `SourceMap.AddMappingOfNodePositions` (`source_map.go`):** Outlined the non-V3 interface fallback into a `//go:noinline` helper and cached `generatorV3`, `lastSourceFile`, and `lastFixedSourceFileName`, eliminating 4 heap escapes per mapped AST node (`234 MB` flat on `bundle-mega`).
   - **$O(1)$ Bounded `childCount` in `CodeGenerator` (`code_generator.go`):** Replaced $O(N)$ `node.GetChildCount()` walks on every `SCRIPT`, `BLOCK`, `OBJECTLIT`, and `ARRAYLIT` node in `AddNodeWithContextAndComments` with a bounded $O(1)$ check (`0..5`).
2. **`ReferenceCollector`, `Scope` Slot Sizing, `ScanVars`, and `NodeTraversal` (`Subagent 5B`):**
   - **Pointer-Keyed `ReferenceCollector` (`reference_collector.go`):** Replaced eager 24-byte `scopedVarKey{scopeRoot, name}` map hashing (`771 MB` on `bundle-mega`) and heap `cols` slice growth (`228 MB`) with pointer-keyed `refByVar map[*Var]*ReferenceCollection` and $O(1)$ `argumentsRefs` / active-traversal scope lookup.
   - **Accurate `Scope.declareNewVar` Arena Slot Sizing (`scope.go`):** Fixed `declareNewVar` so global scopes don't waste a 6-slot arena buffer, function parameter scopes (`s.rootNode.IsFunction()`) size slots by parameter count, and function block scopes (`s.IsFunctionBlockScope()`, where local `var`/`let`/`const` declarations live) allocate up to 8 arena slots instead of `2`, eliminating `LinkedHashMap.addSlow` heap reallocations on functions with 3–8 local variables.
   - **256-Entry Token Tables in `ScanVars` & `NodeTraversal` (`syntactic_scope_creator.go`, `node_traversal.go`):** Added `scanVarsTokenKind`, `nodeTraversalBlockScopeToken`, and `nodeTraversalSpecialToken` lookup tables so `CreatesBlockScope` inlines into `TraverseBranch` and `ScanVars` skips recursive calls on expression/leaf subtrees.
3. **`ControlFlowAnalysis`, `LinkedDirectedGraph`, `LinkedUndirectedGraph`, `CoalesceVariableNames`, and `DAE` (`Subagent 5C`):**
   - **Pooled `astPosition` & Direct AST Walker in `ControlFlowAnalysis` (`control_flow_analysis.go`):** Pooled `astPosition` (`map[*rhino.Node]int32`), stack-allocated `ControlFlowAnalysis` in `ComputeCfg()`, stored `pos int32` inside `cfaWorkItem` (eliminating map lookups in `cfaSiftUp`/`cfaSiftDown`), and replaced `NodeTraversal` in `computeCfg` with a direct scope-free AST walker.
   - **44% Smaller `AstControlFlowGraph` & Zero-Alloc `linkedUndirectedSubGraph` (`graph_linked_directed_graph.go`, `graph_linked_undirected_graph.go`):** Shrunk `ControlFlowAnalysis_AstControlFlowGraph` from `1,904 B` to `1,104 B` and implemented `linkedUndirectedSubGraph` for `LinkedUndirectedGraph.NewSubGraph()`, eliminating 100% of `neighborList()` slice allocations (`192 MB` on `bundle-mega`) during `GreedyGraphColoring.Color`.
   - **Word-Level `BitSet.Or` in `CoalesceVariableNames` (`coalesce_variable_names.go`):** Replaced $O(V)$ per-node escaped-var loops and nested `NextSetBit` loops with word-level `BitSet.Or` and $O(1)$ token/parent checks; skipped candidate traversal in `DeadPropertyAssignmentElimination` when a function has `< 2` property writes.
4. **`JSTypeRegistry.CanPropertyBeDefined`, `UnionType`, `CheckAccessControls`, and `ConvertTypesToColors` (`Subagent 5D`):**
   - **Zero-Alloc `CanPropertyBeDefined` (`jstype_js_type_registry.go`):** Eliminated the `associatedTypes := []JSType{}` slice allocation (`305 MB` on `bundle-mega`) and `.ToSlice()` copies by iterating `nonRefTypesIndexedByProperty` and `eachRefTypeIndexedByProperty` directly and short-circuiting on the first match, and checking `droppedPropertiesOfUnions.Contains(propertyName)` before scanning `rec.GetPropertyMap().entries`.
   - **Inline `UnionType` Alternates (`jstype_union_type.go`), `CheckAccessControls` (`check_access_controls.go`), and `ConvertTypesToColors` (`serialization_convert_types_to_colors.go`):** Added `inlineAlternates [4]JSType` to `UnionType` / `UnionType_Builder`, reused `scratchPropRef` with early token filtering in `CheckAccessControls`, and pre-populated `colorsByPointer []*colors.Color` with a direct AST walk in `ConvertTypesToColors`.
5. **`InferConsts`, `Normalize`, `Denormalize`, `RenameVars`, and `Es6RewriteClass` (`Subagent 5E`):**
   - **Specialized `inferConstsCollector` & `denormalizeVarInliner` (`infer_consts.go`, `denormalize.go`):** Replaced full `ReferenceCollector` runs in `InferConsts` and `Denormalize` Pass 2 with lightweight single-pass collectors indexing state per `Var` via `v.nameHash`.
   - **Fast Undo-Stack Contextual Renamer in `Normalize` (`normalize.go`):** Replaced 7 per-node helper calls in `DoStatementNormalizations` with a single token switch and implemented `useFastContextualRenamer` using a flat `nameUsage` + `activeReplacements` + undo stack instead of per-scope `ContextualRenamer` maps and $O(\text{depth})$ stack walks.
   - **$O(1)$ Indexed `RenameVars` (`rename_vars.go`) & `Es6RewriteClass` (`es6_rewrite_class.go`):** Cached per-`Var` `*RenameVars_Assignment` resolution in `var_.nameHash`, added `localByIndex []*RenameVars_Assignment` (eliminating `"L " + strconv.Itoa(idx)` string formatting and hashing), eliminated `originalNameByNode`, and reused `metadataBuf` in `Es6RewriteClass`.

### Final measured results — Round 5 (`t=1` and `t=4`, interleaved, N=3, PGO enabled)

#### Standard 5 Benchmarks

| Benchmark | Original (`c9872522`) | Round 4 (`289058de` `t=1`/`t=4`) | **Round 5 (`t=1` / `t=4`)** | Round 5 vs Round 4 (`t=1`/`t=4`) | Total vs Original (`t=1`/`t=4`) | Cold / Warm Java (`t=4`) | **`jsgo` vs Cold / Warm Java (`t=4`)** | Peak RSS: `jsgo` vs Cold / Warm Java |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| **`typed-app`** | 2.86 s | 1.95 s / 1.82 s | **1.80 s** / **1.68 s** | **−7.5%** / **−8.0%** | **−37.1%** / **−41.3%** | 6.14 s / 2.18 s | **3.65×** / **1.30×** | **238–245 MiB** vs 758 / 1,871 MiB (**3.1× / 7.6× lower**) |
| **`bundle`** (12.6 MB) | 11.82 s | 6.75 s / 6.62 s | **6.05 s** / **5.92 s** | **−10.3%** / **−10.6%** | **−48.8%** / **−49.9%** | 16.40 s / 8.68 s | **2.77×** / **1.47×** | **571–593 MiB** vs 2,151 / 2,343 MiB (**3.6× / 4.1× lower**) |
| **`downlevel-whitespace`** | 0.86 s | 0.61 s / 0.52 s | **0.60 s** / **0.50 s** | **−2.0%** / **−3.4%** | **−30.2%** / **−41.9%** | 2.55 s / 0.60 s | **5.10×** / **1.20×** | **146–152 MiB** vs 553 / 1,437 MiB (**3.6× / 9.8× lower**) |
| **`downlevel-simple`** | 2.39 s | 1.52 s / 1.43 s | **1.42 s** / **1.32 s** | **−6.6%** / **−7.3%** | **−40.6%** / **−44.8%** | 4.68 s / 1.60 s | **3.55×** / **1.21×** | **192–197 MiB** vs 613 / 1,489 MiB (**3.1× / 7.8× lower**) |
| **`checks-only`** | 4.40 s | 3.23 s / 2.82 s | **3.03 s** / **2.64 s** | **−6.0%** / **−6.5%** | **−31.1%** / **−40.0%** | 6.55 s / 3.20 s | **2.48×** / **1.21×** | **671–682 MiB** vs 1,591 / 3,562 MiB (**2.3× / 5.3× lower**) |

#### Large & Mega Real-World `ADVANCED` Benchmarks

| Benchmark | Workload (`--compilation_level=ADVANCED`) | Round 4 (`289058de` `t=1`/`t=4`) | **Round 5 (`t=1` / `t=4`)** | Round 5 vs Round 4 (`t=1`/`t=4`) | Cold Java | Warm Java | **`jsgo` vs Cold Java (`t=1`/`t=4`)** | **`jsgo` vs Warm Java (`t=1`/`t=4`)** | Peak RSS: `jsgo` vs Cold / Warm Java |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| **`typed-mega`** | **44.6 MB, 4,385 files** (10 Closure `ADVANCED` codebases, 94.1% typed → 14.33 MB JS + 22.29 MB map) | 27.47 s / 26.09 s | **25.23 s** / **24.00 s** | **−8.2%** / **−8.0%** | 38.18 s | 29.53 s | **1.51×** / **1.59× faster** | **1.17×** / **1.23× faster** (−5.53 s) | **1,899–1,950 MiB** vs 6,608 / 6,881 MiB (**3.4× / 3.5× lower**, −662 MiB vs R4) |
| **`bundle-mega`** | **36.1 MB, 26 UMD libraries** (27,597 warnings → 13.09 MB JS + 33.99 MB map) | 37.52 s / 36.59 s | **33.56 s** / **32.70 s** | **−10.6%** / **−10.6%** | 66.71 s | 53.60 s | **1.99×** / **2.04× faster** | **1.60×** / **1.64× faster** (−20.90 s) | **2,407–2,477 MiB** vs 4,421 / 6,775 MiB (**1.8× / 2.7× lower**) |
| **`typed-full`** | **28.5 MB, 1,818 files** (Closure Library + Shaka Player 5.3.1 + Blockly 8.0.5 + app, 97.3% typed → 6.75 MB JS + 13.31 MB map) | 15.91 s / 15.27 s | **14.66 s** / **14.00 s** | **−7.8%** / **−8.3%** | 31.60 s | 20.21 s | **2.16×** / **2.26× faster** | **1.38×** / **1.44× faster** (−6.21 s) | **1,098–1,179 MiB** vs 1,603 / 4,724 MiB (**1.4× / 4.0× lower**, −288 MiB vs R4) |
| **`bundle-advanced-ts`** | **12.6 MB, 4 libraries incl. 8.9 MB `typescript.js`** (10,751 warnings → 3.25 MB JS + 9.39 MB map) | 10.20 s / 10.03 s | **9.08 s** / **8.94 s** | **−11.0%** / **−10.9%** | 20.78 s | 12.41 s | **2.29×** / **2.32× faster** | **1.37×** / **1.39× faster** (−3.47 s) | **643–667 MiB** vs 2,422 / 2,850+ MiB (**3.6× / 4.4× lower**) |
| **`shaka-blockly`** | **8.6 MB, 847 files** (Shaka Player 5.3.1 + Blockly 8.0.5, 93.9% typed → 1.86 MB JS + 4.94 MB map) | 6.76 s / 6.53 s | **6.21 s** / **5.97 s** | **−8.2%** / **−8.5%** | 15.26 s | 8.67 s | **2.46×** / **2.56× faster** | **1.40×** / **1.45× faster** (−2.70 s) | **495–529 MiB** vs 1,341 / 2,852 MiB (**2.5× / 5.4× lower**) |
| **`bundle-advanced`** | **5.4 MB, 11 UMD libraries** (three, lodash, react-dom, rxjs, d3, moment, vue, jquery, handlebars, immutable, underscore → 1.88 MB JS + 4.35 MB map) | 5.18 s / 4.98 s | **4.65 s** / **4.47 s** | **−10.2%** / **−10.3%** | 11.30 s | 6.79 s | **2.43×** / **2.53× faster** | **1.46×** / **1.52× faster** (−2.32 s) | **363–381 MiB** vs 1,600 / 3,795 MiB (**4.2× / 10.0× lower**) |

## Smaller open items

- `tools/clie2e/sweep.go` runs its JVMs under its own `nice -n 10`; under
  an outer `nice` they end at 19 and run ~5× slower when other work
  competes. Drop the inner `nice` or document it.
- The released jar still cannot load its own XTB reader
  (`NoClassDefFoundError …jarjar/javax/xml/parsers/ParserConfigurationException`);
  `tools/clie2e/build-java.sh` works around it for the oracles.
