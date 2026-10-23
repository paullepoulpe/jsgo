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
- **Speed** ([benchmarks.md](benchmarks.md#end-to-end-v20261015) and Rounds 3–22 below):
  after twenty-two rounds of systematic substrate, data-structure, pass-level,
  inter-pass arena, and multi-threaded middle-end architectural optimizations,
  `jsgo` (`t=4`) is **2.12×–5.60× faster than cold Java (`java -jar`)** (**2.20×–5.60×** with `--fast_mode`)
  and **1.30×–2.39× faster than a 5-warmup JIT-compiled JVM** across all 11
  whole-program benchmarks (including 4,385-file / 44.6 MB `typed-mega` in **16.72 s** `compat` / **16.13 s** `fast` and
  26-library / 36.1 MB `bundle-mega` in **20.27 s** `compat` / **20.22 s** `fast`), while using
  **2.55×–7.33× less per-iteration RSS** (**4.01×–11.81× less** than JVM peak RSS).
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

## Round 6: Optimization-Loop Allocation & Hash-Collision Elimination (`−3.79 GB` / `−22.1%` `alloc_space` on `bundle-mega`) (2026-10-20)

Profiling Round 5 (`738e04ca`) on `bundle-mega` (`17.17 GB` `alloc_space`) and `typed-mega` (`11.61 GB` `alloc_space`) uncovered five remaining bottlenecks in the optimization loop (`RemoveUnusedCode`, `InlineVariables`, `ReferenceCollector`, `SyntacticScopeCreator`/`Scope`, `FlowSensitiveInlineVariables`/`DataFlowAnalysis`, `MinimizedCondition`, and `InlineFunctions`):

1. **`RemoveUnusedCode` Hash-Collision Fix, Empty-Block Scope Avoidance, and Pooled Arena (`Subagent 6A`, `remove_unused_code.go`):**
   - **Fixed $O(N^2)$ `GetVarInfo` Bucket Chain Collision:** Previously, `getCanonicalVarSlot` hashed `varKey{name, scopeRoot}` using ONLY `v.nameHash` (`StringHashCode(v.name)`). In `bundle-mega` (where thousands of inner functions declare identically named local variables like `e`, `t`, `n`, `r`, `i`), all `Var`s with the same name collided into a single hash bucket chain, spending `550 ms` of flat CPU scanning 500+-element linked lists in `GetVarInfo`. Mixing `uintptr(unsafe.Pointer(scopeRoot))` into `hashVarKey(name, scopeRoot)` restored $O(1)$ lookup.
   - **Empty Block-Scope Elision (`hasBlockScopedDeclarations`):** Added a shallow child check before `t.GetScope()` on `BLOCK`, `SWITCH_BODY`, `FOR`, `FOR_IN`/`FOR_OF`/`FOR_AWAIT_OF`, and `CLASS` nodes, skipping `CreateScopeOfScope` and `NodeUtil_GetInputId` parent walks on blocks with no lexical declarations.
   - **Pooled `removeUnusedCodeArena` & Intrusive Property Lists:** Replaced `removablesForPropertyNames` (`*jutil.HashMultimap`) and per-instance slabs with a `sync.Pool`-backed `removeUnusedCodeArena` (`rucSlab[T]`) and intrusive FIFO singly-linked property lists (`propRemovables`).
2. **`InlineVariables` Lifetime-Separated `StandardVarExpert` Slabs & Flat `ReferenceCollector` Table (`Subagent 6B`, `inline_variables.go`, `reference_collector.go`, `inline_object_literals.go`):**
   - **Eliminated `360.70 MB` `expertSlab` Leak in `InlineVariables`:** Previously, `DoInlinesForScope` only reset `b.expertSlab` when `len(b.aliasRetryHandlers) == 0`. A single unresolved cross-scope/hoisted alias in `b.aliasRetryHandlers` prevented `b.expertSlab` from resetting for the rest of the AST traversal, allocating `360.70 MB` of `StandardVarExpert` slabs on `bundle-mega`. Evaluating `expert.Analyze()` in a reusable `b.scratchExpert` and splitting slabs into `scopeExpertSlab` (reset unconditionally at the start of every `DoInlinesForScope`) and `retryExpertSlab` (allocated only when `ShouldWaitForAliasedVar()` waits on an unresolved alias) eliminated 99%+ of `StandardVarExpert` allocations.
   - **Open-Addressed `ReferenceCollector` Table & Pooled Chunks:** Replaced `c.refByVar` (`map[*Var]*ReferenceCollection`) with a flat power-of-2 open-addressed hash table `c.table []rcTableEntry` plus `lastVar`/`lastCollection` cache and pooled `referenceCollectorArena` (`sync.Pool`), clearing `col.refsBuf[:]` on release so grown `ReferenceCollection`s never retain stale `*Reference` / `*Scope` pointers across passes.
3. **Lazy `Scope` `CompilerInput` Resolution & Exact-Sized `ScanVars` Arena Buffers (`Subagent 6C`, `scope.go`, `syntactic_scope_creator.go`, `abstract_scope.go`, `node_traversal.go`):**
   - **Lazy `getInput()` & `inheritedInput()`:** Removed eager `NodeUtil_GetInputId(n)` parent walks and `compiler.GetInput(id)` hash lookups from `NewScope`/`CreateScopeOfScope`, resolving `CompilerInput` lazily on first variable declaration (`s.getInput()`) and inheriting `p.input` from parent scopes inside the same script.
   - **Stack-Buffered `smallVars [8]*Var` in `SyntacticScopeCreator_ScopeScanner`:** Buffered the first 8 local `Var`s on `SyntacticScopeCreator_ScopeScanner.smallVars [8]*Var` during `ScanVars` and flushed them once at the end of `Generate`, allocating the exact slot count (`1..8`) from `scopeArena` and pre-sizing `sc.vars` (`max(16, count*2)`) for scopes with $\ge 9$ variables without intermediate heap slice doublings (`−736 MB` in `Scope.declareNewVar`).
   - **Concrete `NodeTraversal.GetScope()` Fast Paths:** Embedded `defaultScopeCreator` by value in `NodeTraversal` and added type-switched fast paths for `*SyntacticScopeCreator` and `*Scope`.
4. **Reusable `FlowSensitiveInlineVariables` Dataflow Engines & Fast `ComputeEscaped` (`Subagent 6D`, `flow_sensitive_inline_variables.go`, `must_be_reaching_variable_def.go`, `maybe_reaching_variable_use.go`, `data_flow_analysis.go`, `control_flow_analysis.go`, `graph_linked_directed_graph.go`):**
   - Reused embedded `MustBeReachingVariableDef`, `MaybeReachingVariableUse`, `DataFlowAnalysisBase`, `ControlFlowAnalysis`, and `LinkedDirectedGraph` slabs across `FlowSensitiveInlineVariables.EnterScope` invocations (`−17.8%` pass time on `typed-mega`, `−12.6%` on `bundle-mega`).
   - Added `dfaHasNestedFunctionOrClass` and `isInsideEscapingScope` fast paths in `DataFlowAnalysis_ComputeEscaped` so functions without nested functions/classes avoid creating inner block scopes.
5. **Deferred AST Materialization in `PeepholeMinimizeConditions` & Fused `InlineFunctions` Analysis (`Subagent 6E`, `minimized_condition.go`, `peephole_minimize_conditions.go`, `inline_functions.go`, `optimize_calls.go`):**
   - Added reusable `scratchMC MinimizedCondition` (`[32]MinimizedCondition_MeasuredNode` inline + reusable `overflow` slice) and `virtualToken` deferred AST node materialization in `PeepholeMinimizeConditions`, cloning/constructing AST nodes only when the minimized representation actually wins (`−18.9%` `latePeepholeOptimizations` time on `bundle-mega`).
   - Fused 3 separate AST passes (`NodeUtil_AllArgsUnescapedLocal` + `CheckNameUsage` + `FindCalledFunctions`) into a single `analyzeFunctionBody` walk in `InlineFunctions`, added slab allocators and 2-slot inline reference maps on `FunctionState`, and added a zero-allocation fast path in `OptimizeCalls_ReferenceMap_GetFunctionNodes` (`−10.0%` `optimizeCalls` time on `typed-mega`).

Total `alloc_space` on `bundle-mega` dropped by **−3.79 GB (−22.1%)** from `17.17 GB` (Round 5) to `13.38 GB` (Round 6).

### Final measured results — Round 6 (`t=1` and `t=4`, interleaved, N=3, PGO enabled)

#### Standard 5 Benchmarks

| Benchmark | Original (`c9872522`) | Round 5 (`738e04ca` `t=1`/`t=4`) | **Round 6 (`t=1` / `t=4`)** | Round 6 vs Round 5 (`t=1`/`t=4`) | Total vs Original (`t=1`/`t=4`) | Cold / Warm Java (`t=4`) | **`jsgo` vs Cold / Warm Java (`t=4`)** | Peak RSS: `jsgo` vs Cold / Warm Java |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| **`typed-app`** | 2.86 s | 1.79 s / 1.67 s | **1.72 s** / **1.60 s** | **−4.0%** / **−4.5%** | **−39.9%** / **−44.1%** | 6.14 s / 2.18 s | **3.84×** / **1.36×** | **240 MiB** vs 758 / 1,871 MiB (**3.2× / 7.8× lower**) |
| **`bundle`** (12.6 MB) | 11.82 s | 6.04 s / 5.88 s | **5.59 s** / **5.46 s** | **−7.4%** / **−7.2%** | **−52.7%** / **−53.8%** | 16.40 s / 8.68 s | **3.00×** / **1.59×** | **526–556 MiB** vs 2,151 / 2,343 MiB (**3.9× / 4.2× lower**) |
| **`downlevel-whitespace`** | 0.86 s | 0.60 s / 0.50 s | **0.58 s** / **0.49 s** | **−1.7%** / **−2.0%** | **−32.6%** / **−43.0%** | 2.55 s / 0.60 s | **5.20×** / **1.22×** | **149–151 MiB** vs 553 / 1,437 MiB (**3.7× / 9.6× lower**) |
| **`downlevel-simple`** | 2.39 s | 1.42 s / 1.33 s | **1.35 s** / **1.26 s** | **−5.3%** / **−5.2%** | **−43.5%** / **−47.3%** | 4.68 s / 1.60 s | **3.71×** / **1.27×** | **207–209 MiB** vs 613 / 1,489 MiB (**2.9× / 7.1× lower**) |
| **`checks-only`** | 4.40 s | 3.02 s / 2.63 s | **2.98 s** / **2.59 s** | **−1.6%** / **−1.3%** | **−32.3%** / **−41.1%** | 6.55 s / 3.20 s | **2.53×** / **1.24×** | **629–649 MiB** vs 1,591 / 3,562 MiB (**2.5× / 5.5× lower**) |

#### Large & Mega Real-World `ADVANCED` Benchmarks

| Benchmark | Workload (`--compilation_level=ADVANCED`) | Round 5 (`738e04ca` `t=1`/`t=4`) | **Round 6 (`t=1` / `t=4`)** | Round 6 vs Round 5 (`t=1`/`t=4`) | Cold Java | Warm Java | **`jsgo` vs Cold Java (`t=1`/`t=4`)** | **`jsgo` vs Warm Java (`t=1`/`t=4`)** | Peak RSS: `jsgo` vs Cold / Warm Java |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| **`typed-mega`** | **44.6 MB, 4,385 files** (10 Closure `ADVANCED` codebases, 94.1% typed → 14.33 MB JS + 22.29 MB map) | 25.21 s / 24.04 s | **24.07 s** / **23.06 s** | **−4.5%** / **−4.1%** | 38.18 s | 29.53 s | **1.59×** / **1.66× faster** | **1.23×** / **1.28× faster** (−6.47 s) | **1,834–1,938 MiB** vs 6,608 / 6,881 MiB (**3.4× / 3.6× lower**, −62 MiB vs R5) |
| **`bundle-mega`** | **36.1 MB, 26 UMD libraries** (27,597 warnings → 13.09 MB JS + 33.99 MB map) | 33.50 s / 32.66 s | **31.24 s** / **30.33 s** | **−6.7%** / **−7.1%** | 66.71 s | 53.60 s | **2.14×** / **2.20× faster** | **1.72×** / **1.77× faster** (−23.27 s) | **2,224–2,417 MiB** vs 4,421 / 6,775 MiB (**1.8× / 2.8× lower**, −281 MiB vs R5) |
| **`typed-full`** | **28.5 MB, 1,818 files** (Closure Library + Shaka Player 5.3.1 + Blockly 8.0.5 + app, 97.3% typed → 6.75 MB JS + 13.31 MB map) | 14.65 s / 13.94 s | **14.00 s** / **13.37 s** | **−4.5%** / **−4.1%** | 31.60 s | 20.21 s | **2.26×** / **2.36× faster** | **1.44×** / **1.51× faster** (−6.84 s) | **1,103–1,182 MiB** vs 1,603 / 4,724 MiB (**1.4× / 4.0× lower**) |
| **`bundle-advanced-ts`** | **12.6 MB, 4 libraries incl. 8.9 MB `typescript.js`** (10,751 warnings → 3.25 MB JS + 9.39 MB map) | 9.08 s / 8.96 s | **8.53 s** / **8.41 s** | **−6.1%** / **−6.1%** | 20.78 s | 12.41 s | **2.44×** / **2.47× faster** | **1.45×** / **1.48× faster** (−4.00 s) | **648–649 MiB** vs 2,422 / 2,850+ MiB (**3.7× / 4.4× lower**) |
| **`shaka-blockly`** | **8.6 MB, 847 files** (Shaka Player 5.3.1 + Blockly 8.0.5, 93.9% typed → 1.86 MB JS + 4.94 MB map) | 6.22 s / 5.97 s | **5.88 s** / **5.66 s** | **−5.4%** / **−5.2%** | 15.26 s | 8.67 s | **2.60×** / **2.70× faster** | **1.47×** / **1.53× faster** (−3.01 s) | **495–501 MiB** vs 1,341 / 2,852 MiB (**2.7× / 5.7× lower**) |
| **`bundle-advanced`** | **5.4 MB, 11 UMD libraries** (three, lodash, react-dom, rxjs, d3, moment, vue, jquery, handlebars, immutable, underscore → 1.88 MB JS + 4.35 MB map) | 4.69 s / 4.49 s | **4.39 s** / **4.23 s** | **−6.5%** / **−5.9%** | 11.30 s | 6.79 s | **2.57×** / **2.67× faster** | **1.55×** / **1.61× faster** (−2.56 s) | **368–408 MiB** vs 1,600 / 3,795 MiB (**3.9× / 9.3× lower**) |

## Round 7: Inter-Pass Arena Retention, Reusable CFG/Liveness State, and Pass-Level Slabs (`−4.35 GB` / `−29.2%` `alloc_space` on `bundle-mega`, `−3.12 GB` / `−30.0%` on `typed-mega`) (2026-10-20)

Profiling Round 6 (`d68a146b`) on `bundle-mega` (`14.92 GB` `alloc_space`) and `typed-mega` (`10.40 GB` `alloc_space`) showed that `sync.Pool`-backed arenas (`scopeArena`, `referenceCollectorArena`, `removeUnusedCodeArena`) were being evicted every 2 GC cycles (`runtime.poolCleanup`) between compiler passes, and that `ControlFlowAnalysis`, `CoalesceVariableNames`, `DeadAssignmentsElimination`, `CrossChunkCodeMotion`, `CodePrinter`, `SourceMapGeneratorV3`, `ReplaceCssNames`, `GlobalNamespace`, and `DisambiguateProperties` were still allocating per-invocation state:

1. **GC-Proof 1-Slot `atomic.Pointer` Inter-Pass Arena Cache & Scaled Chunk Limits (`scope.go`, `reference_collector.go`, `cross_chunk_reference_collector.go`, `remove_unused_code.go`, `memoized_scope_creator.go`):**
   - **Why `sync.Pool` Alone Fails Across Multi-Pass Compilations:** Between two invocations of `InlineVariables` or `RemoveUnusedCode`, 5–10 other passes execute and trigger 1–2 Go GC cycles (`30–34` GC cycles across `179–204` passes in `bundle-mega` / `typed-mega`). Go's `runtime.poolCleanup` moves `local` to `victim` on GC #1 and drops `victim` on GC #2. Placing a 1-slot `atomic.Pointer[Arena]` cache (`cachedScopeArena`, `cachedRCArena`, `cachedCCAArena`, `cachedRUCArena`) in front of `sync.Pool` retains a single cleared, zero-pointer arena across GC cycles while falling back to `sync.Pool` under concurrent compilations.
   - **Scaled `scopeArena` & `rucSlab` Retention Limits:** Raised `scopeArena` chunk retention (`maxScopeChunkSize = 1024, maxScopeChunks = 80`, `maxVarChunkSize = 2048, maxVarChunks = 80`, `maxSlotChunkSize = 2048, maxSlotChunks = 80`) and `rucSlab` `maxRetainedChunks` (`8 → 40`), and wired `ReleaseArena()` into `MemoizedScopeCreator`, `CoalesceVariableNames`, `FlowSensitiveInlineVariables`, `DeadAssignmentsElimination`, `InferConsts`, `Normalize`, and `VarCheck`.
2. **Reusable `ControlFlowAnalysis` / `AstControlFlowGraph` & `LiveVariablesAnalysis` State (`Subagent 7C`, `node_traversal.go`, `type_inference_pass.go`, `coalesce_variable_names.go`, `dead_assignments_elimination.go`, `live_variables_analysis.go`, `check_unreachable_code.go`, `check_missing_return.go`):**
   - Replaced per-scope `ComputeCfg()` allocations in `NodeTraversal.GetControlFlowGraph` (`1,284 MB` on `bundle-mega`) and `TypeInferencePass.ComputeCfg` with `t.cfaStorage.ComputeReusableCfg(..., &t.cfgStorage)`, cutting `ControlFlowAnalysis` allocations by **−97.7%** (`1,284 MB → 30 MB`).
   - Replaced per-function `LiveVariablesAnalysis`, `CombinedLiveVariablesAnalysis`, and `UndirectedGraph` allocations in `CoalesceVariableNames` and `DeadAssignmentsElimination` with per-depth reusable `coalesceScopeState` (`−22.1%` `coalesceVariableNames` pass time on `bundle-mega`, `−21.3%` on `typed-mega`).
3. **Pooled `CrossChunkReferenceCollector` / `CrossChunkCodeMotion` Arenas (`Subagent 7B`, `reference_collector.go`, `cross_chunk_reference_collector.go`, `cross_chunk_code_motion.go`):**
   - Added `crossChunkArena` slabbing `TopLevelStatement`, `GlobalSymbol`, `DeclarationStatementGroup`, and `GlobalSymbolCycle` and indexing `symbolsByVarIndex` by `Var.index` (`−41.2%` `crossChunkCodeMotion` pass time on `typed-mega`, `−32.3%` on `bundle-mega`).
4. **Pooled `CodePrinter` Buffers, `SourceMapGeneratorV3` Mapping Slabs, and Direct AST Walk in `ReplaceCssNames` (`Subagent 7D`, `code_printer.go`, `source_map_generator_v3.go`, `replace_css_names.go`, `node_util.go`):**
   - Pooled `CodePrinter` buffers (`codePrinterBufferPool`) and added 4,096-entry `Mapping` slabs in `SourceMapGeneratorV3`, cutting `CodePrinter` allocations by **−44.5%** (`1,040 MB → 577 MB`).
   - Replaced full `NodeTraversal` in `ReplaceCssNames` with a direct recursive AST walker (`−55.9%` pass time on `typed-mega`, `−55.8%` on `bundle-mega`).
5. **`GlobalNamespace` Ref Slabs & Pre-Sized `DisambiguateProperties` (`Subagent 7E`, `global_namespace.go`, `disambiguate_color_find_property_references.go`, `disambiguate_property_clustering.go`, `memoized_scope_creator.go`, `optimize_calls.go`):**
   - Replaced per-ref heap allocations in `GlobalNamespace` with a `jutil.Slab[GlobalNamespace_Ref]`, pre-sized `PropClustering` maps in `DisambiguateProperties` (`−14.8%` pass time on `typed-mega`, `−13.1%` on `bundle-mega`), and added a `[8]memoizedScopeEntry` inline buffer + `lastNode`/`lastScope` cache in `MemoizedScopeCreator`.

Total `alloc_space` dropped by **−4.35 GB (−29.2%)** on `bundle-mega` (`14.92 GB → 10.57 GB`) and **−3.12 GB (−30.0%)** on `typed-mega` (`10.40 GB → 7.28 GB`).

### Final measured results — Round 7 (`t=1` and `t=4`, interleaved, N=3, PGO enabled)

#### Standard 5 Benchmarks

| Benchmark | Original (`c9872522`) | Round 6 (`d68a146b` `t=1`/`t=4`) | **Round 7 (`t=1` / `t=4`)** | Round 7 vs Round 6 (`t=1`/`t=4`) | Total vs Original (`t=1`/`t=4`) | Cold / Warm Java (`t=4`) | **`jsgo` vs Cold / Warm Java (`t=4`)** | Peak RSS: `jsgo` vs Cold / Warm Java |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| **`typed-app`** | 2.86 s | 1.73 s / 1.61 s | **1.69 s** / **1.57 s** | **−2.5%** / **−2.5%** | **−40.9%** / **−45.1%** | 6.14 s / 2.18 s | **3.91×** / **1.39×** | **260–265 MiB** vs 758 / 1,871 MiB (**2.9× / 7.1× lower**) |
| **`bundle`** (12.6 MB) | 11.82 s | 5.59 s / 5.48 s | **5.47 s** / **5.32 s** | **−2.2%** / **−2.9%** | **−53.7%** / **−55.0%** | 16.40 s / 8.68 s | **3.08×** / **1.63×** | **597–622 MiB** vs 2,151 / 2,343 MiB (**3.5× / 3.9× lower**) |
| **`downlevel-whitespace`** | 0.86 s | 0.59 s / 0.49 s | **0.59 s** / **0.49 s** | **−0.2%** / **+0.4%** | **−31.4%** / **−43.0%** | 2.55 s / 0.60 s | **5.20×** / **1.22×** | **165–166 MiB** vs 553 / 1,437 MiB (**3.3× / 8.7× lower**) |
| **`downlevel-simple`** | 2.39 s | 1.35 s / 1.26 s | **1.30 s** / **1.20 s** | **−4.0%** / **−4.6%** | **−45.6%** / **−49.8%** | 4.68 s / 1.60 s | **3.90×** / **1.33×** | **207–215 MiB** vs 613 / 1,489 MiB (**2.9× / 7.2× lower**) |
| **`checks-only`** | 4.40 s | 2.98 s / 2.61 s | **2.98 s** / **2.59 s** | **−0.1%** / **−0.8%** | **−32.3%** / **−41.1%** | 6.55 s / 3.20 s | **2.53×** / **1.24×** | **650–687 MiB** vs 1,591 / 3,562 MiB (**2.4× / 5.5× lower**) |

#### Large & Mega Real-World `ADVANCED` Benchmarks

| Benchmark | Workload (`--compilation_level=ADVANCED`) | Round 6 (`d68a146b` `t=1`/`t=4`) | **Round 7 (`t=1` / `t=4`)** | Round 7 vs Round 6 (`t=1`/`t=4`) | Cold Java | Warm Java | **`jsgo` vs Cold Java (`t=1`/`t=4`)** | **`jsgo` vs Warm Java (`t=1`/`t=4`)** | Peak RSS: `jsgo` vs Cold / Warm Java |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| **`typed-mega`** | **44.6 MB, 4,385 files** (10 Closure `ADVANCED` codebases, 94.1% typed → 14.33 MB JS + 22.29 MB map) | 24.08 s / 22.92 s | **23.63 s** / **22.35 s** | **−1.9%** / **−2.5%** | 38.18 s | 29.53 s | **1.62×** / **1.71× faster** | **1.25×** / **1.32× faster** (−7.18 s) | **1,846–1,932 MiB** vs 6,608 / 6,881 MiB (**3.4× / 3.6× lower**) |
| **`bundle-mega`** | **36.1 MB, 26 UMD libraries** (27,597 warnings → 13.09 MB JS + 33.99 MB map) | 31.29 s / 30.39 s | **30.63 s** / **29.82 s** | **−2.1%** / **−1.9%** | 66.71 s | 53.60 s | **2.18×** / **2.24× faster** | **1.75×** / **1.80× faster** (−23.78 s) | **2,283–2,365 MiB** vs 4,421 / 6,775 MiB (**1.9× / 2.9× lower**, −62 MiB vs R6) |
| **`typed-full`** | **28.5 MB, 1,818 files** (Closure Library + Shaka Player 5.3.1 + Blockly 8.0.5 + app, 97.3% typed → 6.75 MB JS + 13.31 MB map) | 14.08 s / 13.40 s | **13.64 s** / **12.97 s** | **−3.2%** / **−3.2%** | 31.60 s | 20.21 s | **2.32×** / **2.44× faster** | **1.48×** / **1.56× faster** (−7.24 s) | **1,114–1,120 MiB** vs 1,603 / 4,724 MiB (**1.4× / 4.2× lower**, −61 MiB vs R6) |
| **`bundle-advanced-ts`** | **12.6 MB, 4 libraries incl. 8.9 MB `typescript.js`** (10,751 warnings → 3.25 MB JS + 9.39 MB map) | 8.51 s / 8.38 s | **8.31 s** / **8.18 s** | **−2.3%** / **−2.4%** | 20.78 s | 12.41 s | **2.50×** / **2.54× faster** | **1.49×** / **1.52× faster** (−4.23 s) | **602 MiB** vs 2,422 / 2,850+ MiB (**4.0× / 4.7× lower**, −41 MiB vs R6) |
| **`shaka-blockly`** | **8.6 MB, 847 files** (Shaka Player 5.3.1 + Blockly 8.0.5, 93.9% typed → 1.86 MB JS + 4.94 MB map) | 5.89 s / 5.64 s | **5.69 s** / **5.47 s** | **−3.5%** / **−3.0%** | 15.26 s | 8.67 s | **2.68×** / **2.79× faster** | **1.52×** / **1.59× faster** (−3.20 s) | **486–500 MiB** vs 1,341 / 2,852 MiB (**2.7× / 5.7× lower**, −21 MiB vs R6) |
| **`bundle-advanced`** | **5.4 MB, 11 UMD libraries** (three, lodash, react-dom, rxjs, d3, moment, vue, jquery, handlebars, immutable, underscore → 1.88 MB JS + 4.35 MB map) | 4.37 s / 4.22 s | **4.27 s** / **4.09 s** | **−2.4%** / **−3.1%** | 11.30 s | 6.79 s | **2.65×** / **2.76× faster** | **1.59×** / **1.66× faster** (−2.70 s) | **370–371 MiB** vs 1,600 / 3,795 MiB (**4.3× / 10.2× lower**) |

## Round 8: Mega-Scale `ReferenceCollector`/`Scope` Arena Retention, Token-Masked `PeepholeOptimizationsPass`, `HamtPMap` / `TemplateTypeMap` Type-System Fast Paths, and Fused Inlining (`−2.18 GB` / `−20.1%` `alloc_space` on `bundle-mega`, `−6.3%` to `−8.9%` Wall Time) (2026-10-20)

Profiling Round 7 (`abcf69a0`) on `bundle-mega` (`10.83 GB` `alloc_space`) and `typed-mega` (`7.45 GB` `alloc_space`) exposed five major CPU and allocation bottlenecks across mega-scale `ReferenceCollector` / `Scope` arena caps, `RemoveUnusedCode` / `PolyfillUsageFinder` guarded-AST callbacks, `PeepholeOptimizationsPass` per-node virtual dispatch, `InlineFunctions` / `FlowSensitiveInlineVariables` AST walks, and `HamtPMap` / `TemplateTypeMap` / `TypedScope` in `inferTypes`:

1. **Mega-Scale `ReferenceCollector` & `scopeArena` Retention + Subtree Scope Isolation (`Subagent 8A`, `reference_collector.go`, `cross_chunk_reference_collector.go`, `scope.go`, `syntactic_scope_creator.go`):**
   - **Eliminated `1.50 GB` `ReferenceCollector` Arena Drop on `bundle-mega`:** On `bundle-mega` (~300,000 `Var`s and ~600,000 `Reference`s per pass), `c.table` grew to `1,048,576` entries (`16 MB`), `c.vars`/`c.cols` grew to `524,288` (`4 MB`), and `refChunks` reached ~240 chunks (`~44 MB`). Round 7's `release()` discarded `c.table` when `> 524,288`, `c.vars`/`c.cols` when `> 262,144`, and truncated `refChunks` to `80`—throwing away `~65 MB` of arena buffers on every single one of the 15 `ReferenceCollector` passes (`1.50 GB` total) and re-growing `c.table` from `1,024` via 10 doublings (`366 MB` in `growTable`). Raising retention caps (`len(c.table) <= 2,097,152`, `cap(c.vars) <= 1,048,576`, `maxRetainedRefChunks = 256`, `maxRetainedColChunks = 192`, `maxRetainedBlkChunks = 224`) reduced `ReferenceCollector.getOrCreateReferences` allocations on `bundle-mega` by **−89.6%** (`930.60 MB → <97.15 MB`)!
   - **Prevented Subtree Traversals from Stealing Full-AST `scopeArena`s:** `SyntacticScopeCreator.CreateScopeOfScope` previously called `acquireScopeArena()` even for small subtree traversals that never called `ReleaseArena()`, stealing the warm 30 MB full-AST `scopeArena` from `cachedScopeArenas` and leaking it to the GC. Restricting `acquireScopeArena()` to full-AST traversals (`n.IsRoot()` / top-level `SCRIPT`) or explicit `c.poolable` creators and raising `maxScopeChunkSize = 2048, maxScopeChunks = 96`, `maxVarChunkSize = 4096, maxVarChunks = 96` cut `scopeArena` allocations on `bundle-mega` by **−32.8%** (`1,636 MB → 1,107 MB`).
2. **`PolyfillUsageFinder` Bitmask Rejections & `RemoveUnusedCode` Fast Paths (`Subagent 8B`, `polyfill_usage_finder.go`, `guarded_callback.go`, `remove_unused_code.go`):**
   - **Early `canPossiblyBeGuarded()` Bailout & `[256]uint64` First-Byte/Length Bitmasks:** In `RemoveUnusedCode`, `PolyfillUsageFinder.TraverseOnlyGuarded` visited every `NAME` and `GETPROP` node across the entire AST only to discard 99.99% when `!IsGuarded()`. Checking `canPossiblyBeGuarded()` at the top of `VisitGuarded` and precomputing `[256]uint64` `(firstByte, length)` bitmasks + native Go maps on `PolyfillUsageFinder_Polyfills` eliminated ~95% of `PolyfillUsageFinder` CPU time.
   - **Scaled `RemoveUnusedCode.varBuckets` (`16,384 → 262,144`) + `lastVar` Cache:** Replaced `hashVarKey(name, scopeRoot)` with precomputed `var_.HashCode()` mixing, added a 1-entry `(lastVar, lastVarInfo)` cache on `GetVarInfo`, and deferred `codingConvention.IsExportedWithLocal` until first `VarInfo` creation.
3. **Token-Indexed `[256]uint16` Dispatch Table in `PeepholeOptimizationsPass` (`Subagent 8C`, `peephole_optimizations_pass.go`, `peephole_substitute_alternate_syntax.go`, `peephole_remove_dead_code.go`, `peephole_fold_constants.go`, `peephole_replace_known_methods.go`):**
   - Precomputed a `[256]uint16` bitmask table `activeMaskByToken` mapping each `rhino.Token` to the exact subset of `peepholeOptimizations` that handle that token, iterating only active bits via `bits.TrailingZeros16(mask)` and recomputing higher bits only when `OptimizeSubtree` changes `currentNode.GetToken()`.
   - Added token-specific fast paths in `PeepholeRemoveDeadCode.TryFoldExpr` / `TryOptimizeBlock`, `PeepholeSubstituteAlternateSyntax.TryReduceCall`, `PeepholeFoldConstants.TryFoldGetProp`, and `PeepholeReplaceKnownMethods`.
4. **Fused `FunctionInjector` / `InlineFunctions` & Bitmask `MustBeReachingVariableDef` / `FlowSensitiveInlineVariables` (`Subagent 8D`, `function_injector.go`, `inline_functions.go`, `flow_sensitive_inline_variables.go`, `must_be_reaching_variable_def.go`, `maybe_reaching_variable_use.go`, `data_flow_analysis.go`):**
   - Fused `DoesFunctionMeetMinimumRequirements` (`referencesArguments` + `blocksInjection`), `ReferencesEval`, `CallMeetsBlockInliningRequirements`, and `InlineCostDelta` in `FunctionInjector` into allocation-free recursive AST walks, and shared a pooled `SyntacticScopeCreator` across `InlineFunctions.Process`.
   - Added `depMask uint64`, `escapedMask uint64`, `paramMask uint64`, and `catchMask uint64` bitmasks in `MustBeReachingVariableDef`, `EqualsLattice` in `DataFlowAnalysis` (bypassing `jutil.ObjectsEquals` interface boxing on every CFG flow step), and direct recursive AST walks in `FlowSensitiveInlineVariables`.
5. **Zero-Alloc `HamtPMap.CheckInvariants`, Concrete `LinkedFlowScope`, `TemplateTypeMap` Fast Paths, and `typedScopeArena` (`Subagent 8E`, `hamt_p_map.go`, `jstype_template_type_map.go`, `linked_flow_scope.go`, `typed_scope.go`, `typed_scope_creator.go`, `type_inference.go`, `type_inference_pass.go`):**
   - Disabled debug `HamtPMap.CheckInvariants()` recursive tree walks on every `Plus`/`Minus`/`Reconcile` in production builds (`hamtCheckInvariants = false`), added concrete `*HamtPMap` methods in `LinkedFlowScope`, and cached `TypedScope.HashCode()` / pointer-identity `Equals` (`−34%` `inferTypes` CPU time).
   - Replaced `ArrayList` + `BitSet` allocations in `NewTemplateTypeMap` / `CopyFilledWithValues` / `CopyWithExtensionOfKeysWithSubMapStarts` with `resolveSimpleTemplateValue` fast paths, and slabbed `TypedScope` / `TypedVar` / `reservedNames` in `typedScopeArena`.

Total `alloc_space` dropped by another **−2.18 GB (−20.1%)** on `bundle-mega` (`10.83 GB → 8.65 GB`, **−49.6% vs Round 5**) and **−780 MB (−10.5%)** on `typed-mega` (`7.45 GB → 6.67 GB`, **−42.5% vs Round 5**).

### Final measured results — Round 8 (`t=1` and `t=4`, interleaved, N=3, PGO enabled)

#### Standard 5 Benchmarks

| Benchmark | Original (`c9872522`) | Round 7 (`abcf69a0` `t=1`/`t=4`) | **Round 8 (`t=1` / `t=4`)** | Round 8 vs Round 7 (`t=1`/`t=4`) | Total vs Original (`t=1`/`t=4`) | Cold / Warm Java (`t=4`) | **`jsgo` vs Cold / Warm Java (`t=4`)** | Peak RSS: `jsgo` vs Cold / Warm Java |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| **`typed-app`** | 2.86 s | 1.68 s / 1.57 s | **1.62 s** / **1.50 s** | **−3.8%** / **−4.5%** | **−43.4%** / **−47.6%** | 6.14 s / 2.18 s | **4.09×** / **1.45×** | **233–245 MiB** vs 758 / 1,871 MiB (**3.1× / 7.6× lower**) |
| **`bundle`** (12.6 MB) | 11.82 s | 5.46 s / 5.33 s | **5.11 s** / **4.97 s** | **−6.4%** / **−6.7%** | **−56.8%** / **−58.0%** | 16.40 s / 8.68 s | **3.30×** / **1.75×** | **603–610 MiB** vs 2,151 / 2,343 MiB (**3.5× / 3.9× lower**) |
| **`downlevel-whitespace`** | 0.86 s | 0.58 s / 0.49 s | **0.58 s** / **0.49 s** | **−1.0%** / **−1.3%** | **−32.6%** / **−43.0%** | 2.55 s / 0.60 s | **5.20×** / **1.22×** | **158–162 MiB** vs 553 / 1,437 MiB (**3.4× / 9.1× lower**, −15.5% RSS) |
| **`downlevel-simple`** | 2.39 s | 1.29 s / 1.21 s | **1.23 s** / **1.14 s** | **−4.9%** / **−5.6%** | **−48.5%** / **−52.3%** | 4.68 s / 1.60 s | **4.11×** / **1.40×** | **209–210 MiB** vs 613 / 1,489 MiB (**2.9× / 7.1× lower**) |
| **`checks-only`** | 4.40 s | 2.98 s / 2.60 s | **2.80 s** / **2.42 s** | **−6.3%** / **−7.1%** | **−36.4%** / **−45.0%** | 6.55 s / 3.20 s | **2.71×** / **1.32×** | **605–629 MiB** vs 1,591 / 3,562 MiB (**2.5× / 5.9× lower**, −11.7% RSS) |

#### Large & Mega Real-World `ADVANCED` Benchmarks

| Benchmark | Workload (`--compilation_level=ADVANCED`) | Round 7 (`abcf69a0` `t=1`/`t=4`) | **Round 8 (`t=1` / `t=4`)** | Round 8 vs Round 7 (`t=1`/`t=4`) | Cold Java | Warm Java | **`jsgo` vs Cold Java (`t=1`/`t=4`)** | **`jsgo` vs Warm Java (`t=1`/`t=4`)** | Peak RSS: `jsgo` vs Cold / Warm Java |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| **`typed-mega`** | **44.6 MB, 4,385 files** (10 Closure `ADVANCED` codebases, 94.1% typed → 14.33 MB JS + 22.29 MB map) | 23.47 s / 22.37 s | **21.99 s** / **20.88 s** | **−6.3%** / **−6.7%** | 38.18 s | 29.53 s | **1.74×** / **1.83× faster** | **1.34×** / **1.41× faster** (−8.65 s) | **1,769–1,807 MiB** vs 6,608 / 6,881 MiB (**3.7× / 3.9× lower**, −149 MiB vs R7) |
| **`bundle-mega`** | **36.1 MB, 26 UMD libraries** (27,597 warnings → 13.09 MB JS + 33.99 MB map) | 30.96 s / 29.76 s | **28.20 s** / **27.41 s** | **−8.9%** / **−7.9%** | 66.71 s | 53.60 s | **2.37×** / **2.43× faster** | **1.90×** / **1.96× faster** (−26.19 s) | **2,357–2,393 MiB** vs 4,421 / 6,775 MiB (**1.8× / 2.8× lower**) |
| **`typed-full`** | **28.5 MB, 1,818 files** (Closure Library + Shaka Player 5.3.1 + Blockly 8.0.5 + app, 97.3% typed → 6.75 MB JS + 13.31 MB map) | 13.61 s / 12.97 s | **12.71 s** / **12.07 s** | **−6.6%** / **−6.9%** | 31.60 s | 20.21 s | **2.49×** / **2.62× faster** | **1.59×** / **1.67× faster** (−8.14 s) | **1,008–1,011 MiB** vs 1,603 / 4,724 MiB (**1.6× / 4.7× lower**, −109 MiB vs R7) |
| **`bundle-advanced-ts`** | **12.6 MB, 4 libraries incl. 8.9 MB `typescript.js`** (10,751 warnings → 3.25 MB JS + 9.39 MB map) | 8.32 s / 8.18 s | **7.81 s** / **7.67 s** | **−6.2%** / **−6.3%** | 20.78 s | 12.41 s | **2.66×** / **2.71× faster** | **1.59×** / **1.62× faster** (−4.74 s) | **608–609 MiB** vs 2,422 / 2,850+ MiB (**4.0× / 4.7× lower**) |
| **`shaka-blockly`** | **8.6 MB, 847 files** (Shaka Player 5.3.1 + Blockly 8.0.5, 93.9% typed → 1.86 MB JS + 4.94 MB map) | 5.70 s / 5.46 s | **5.31 s** / **5.07 s** | **−6.9%** / **−7.0%** | 15.26 s | 8.67 s | **2.87×** / **3.01× faster** | **1.63×** / **1.71× faster** (−3.60 s) | **503–506 MiB** vs 1,341 / 2,852 MiB (**2.7× / 5.6× lower**) |
| **`bundle-advanced`** | **5.4 MB, 11 UMD libraries** (three, lodash, react-dom, rxjs, d3, moment, vue, jquery, handlebars, immutable, underscore → 1.88 MB JS + 4.35 MB map) | 4.27 s / 4.07 s | **3.99 s** / **3.82 s** | **−6.6%** / **−6.3%** | 11.30 s | 6.79 s | **2.83×** / **2.96× faster** | **1.70×** / **1.78× faster** (−2.97 s) | **367–376 MiB** vs 1,600 / 3,795 MiB (**4.3× / 10.1× lower**) |

## Round 9: Arena-Backed `Scope` Hash Tables, Pooled `ParseArena`, Zero-Alloc Streaming IO / Source Maps, and `Var` Pointer Hashing (`−2.01 GB` / `−23.2%` `alloc_space` on `bundle-mega`, `−1.46 GB` / `−21.9%` on `typed-mega`) (2026-10-20)

Profiling Round 8 (`f1c0cd2e`) on `bundle-mega` (`8.65 GB` `alloc_space`) and `typed-mega` (`6.67 GB` `alloc_space`) uncovered five remaining allocation and CPU bottlenecks across `Scope` / `scopeArena`, `ParserRunner_Parse`, `BufferedWriter` / `OutputStreamWriter` / `SourceMap`, `ReferenceCollector` / `RemoveUnusedCode`, and `ConvertTypesToColors` / `SubtypeChecker`:

1. **Arena-Allocated `Scope` `varList []*Var` + Open-Addressed `varHash []int32` & `scopeArena` Pool Leak Fix (`Subagent 9A`, `scope.go`, `syntactic_scope_creator.go`, `abstract_scope.go`, `memoized_scope_creator.go`, and 12 pass files):**
   - **Stopped `NodeTraversal.defaultScopeCreator` from Stealing and Leaking `cachedScopeArenas`:** In Round 8, `SyntacticScopeCreator.CreateScopeOfScope` acquired a warm `scopeArena` whenever `parent == nil && n.IsRoot()`, which matched `t.defaultScopeCreator` in passes that never called `ReleaseArena()`—stealing warm 50 MB `scopeArena`s and dropping them when `releaseNodeTraversal` cleared `*t`. Restricting `acquireScopeArena()` to `c.poolable` creators and wiring `EnableArenaPool()` + `ReleaseArena()` across `MemoizedScopeCreator`, `CoalesceVariableNames`, `DeadAssignmentsElimination`, `FlowSensitiveInlineVariables`, `InferConsts`, `InlineFunctions`, `Normalize`, `PolyfillUsageFinder`, `VarCheck`, `VariableReferenceCheck`, `PureFunctionIdentifier`, `MakeDeclaredNamesUnique`, and `RenameVars` eliminated `~500 MB` of leaked `scopeArena` slabs.
   - **Replaced `LinkedHashMap` on `Scope` with Arena-Allocated `varList []*Var` and `varHash []int32`:** Replaced `Scope`'s `LinkedHashMap[string, *Var]` (which allocated `291 MB` in `buildIndex` and `175 MB` in `addSlow` on large UMD scopes) with `varList []*Var` (8 bytes/var instead of 32 bytes/slot) and an open-addressed `varHash []int32` index allocated from `scopeArena.hashChunk`, plus lazy single-hash ancestor-chain lookup in `Scope.GetVar` and `Scope.HasSlot`.
2. **Pooled `parser.ParseArena` & Inline `EnumSet` (`Subagent 9B`, `parse_arena.go`, `parser.go`, `scanner.go`, `line_number_scanner.go`, `source_file.go`, `parser_runner.go`, `feature_set.go`, `enumset.go`):**
   - Created a pooled `parser.ParseArena` acquired per file in `ParserRunner_Parse`, reusing `SourceFile.ContentsUTF16` (`[]uint16`), `SourcePosition`, `SourceRange`, `TokenBase`, `IdentifierToken`, `LiteralTokenBase`, `StringLiteralToken`, and 35 high-frequency `ParseTree` node slabs, and cloning `topLevelStatementRanges` when `ShouldParseDescriptions()` is enabled.
   - Added an inline `buf [2]uint64` to `jutil.EnumSet` and eliminated temporary `EnumSet` allocations in `FeatureSet.ContainsFeature`.
3. **Zero-Allocation `BufferedWriter` / `OutputStreamWriter`, `JsonReader` Buffer Reuse, and Value `SourceMap` Mappings (`Subagent 9C`, `buffered_writer.go`, `output_stream_writer.go`, `gson_json_reader.go`, `source_map_consumer_v3.go`, `source_map.go`, `code_printer.go`, `compiler.go`):**
   - Eliminated the 5 heap allocations per 8 KB chunk (`~350 MB` on `bundle-mega`) in `BufferedWriter.Write`, `BufferedWriter.FlushBuffer`, and `OutputStreamWriter.Write` by streaming UTF-16 code units directly into a reusable `byteBuf []byte` while preserving exact JDK 21 surrogate splitting and replacement semantics.
   - Added a reusable `builder []uint16` buffer in `JsonReader.nextQuotedValue` (`−120 MB`), incremental ASCII decoding in `StringReader.Read` (`−51 MB`), ASCII VLQ fast paths in `SourceMapConsumerV3_StringCharIterator` (`−22 MB`), value-returning `GetMappingForLineValue` / `getSourceMappingValue` in `SourceMap.AddMappingOfNodePositions` (`−163 MB`), `//go:noinline` panic extraction in `codePrinter_MappedCodePrinter_adjustEndPositionVal` (`−49 MB`), and `toSourceWithoutMappings` in `Compiler.ToSourceOfNode`.
4. **48-Byte `Var` Pointer Hash De-Clustering (`rcHashVarPtr`) & Direct-Mapped `RemoveUnusedCode` Cache (`Subagent 9D`, `remove_unused_code.go`, `reference_collector.go`, `cross_chunk_reference_collector.go`, `inline_variables.go`, `optimize_calls.go`):**
   - **Fixed 16-Byte Alignment Clustering in `ReferenceCollector.table`:** Because `sizeof(Var) == 48` (`3 × 16`), slab-allocated `Var`s in `scopeArena` have 16-byte-aligned addresses where the low 4 bits of `uintptr(unsafe.Pointer(v))` are always `0`, leaving 15/16 of `c.table` buckets empty under `IdentityHashCodeOfPointer(v)`. Replacing with `rcHashVarPtr(v)` (`uintptr >> 4` + MurmurHash3 64-bit finalizer) cut `ReferenceCollector.getOrCreateReferences` CPU time by **−71%** (`0.65s → 0.19s`).
   - Added a 2048-entry direct-mapped `*Var -> RemoveUnusedCode_VarInfo` pointer cache (`varPtrCache`) + unremovable `lastVarInfo` caching in `RemoveUnusedCode.GetVarInfo`, cutting `GetVarInfo` CPU time by **−56%** (`0.75s → 0.33s`).
5. **Direct AST Walks in `SerializeTypesToPointers`, Single-Pass `TakeTypesAndJSDocForColorConversion`, and `SubtypeChecker` Fast Paths (`Subagent 9E`, `serialization_serialize_types_to_pointers.go`, `serialization_convert_types_to_colors.go`, `node.go`, `jstype_equality_checker.go`, `jstype_subtype_checker.go`, `type_check.go`):**
   - Replaced both `NodeTraversal_Traverse` passes in `SerializeTypesToPointers.GatherTypesOnAst` with cached direct recursive AST walks, added single-pass `Node.TakeTypesAndJSDocForColorConversion()` for `ConvertTypesToColors`, and replaced `jutil.ObjectsEquals` interface boxing in `SubtypeChecker` with `equalityChecker_AreIdentityEqual`.

Total `alloc_space` dropped by another **−2.01 GB (−23.2%)** on `bundle-mega` (`8.65 GB → 6.64 GB`, **−61.3% vs Round 5**) and **−1.46 GB (−21.9%)** on `typed-mega` (`6.67 GB → 5.21 GB`, **−55.1% vs Round 5**).

### Final measured results — Round 9 (`t=1` and `t=4`, interleaved, N=3, PGO enabled)

#### Standard 5 Benchmarks

| Benchmark | Original (`c9872522`) | Round 8 (`f1c0cd2e` `t=1`/`t=4`) | **Round 9 (`t=1` / `t=4`)** | Round 9 vs Round 8 (`t=1`/`t=4`) | Total vs Original (`t=1`/`t=4`) | Cold / Warm Java (`t=4`) | **`jsgo` vs Cold / Warm Java (`t=4`)** | Peak RSS: `jsgo` vs Cold / Warm Java |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| **`typed-app`** | 2.86 s | 1.62 s / 1.50 s | **1.58 s** / **1.46 s** | **−2.5%** / **−2.9%** | **−44.8%** / **−49.0%** | 6.14 s / 2.18 s | **4.21×** / **1.49×** | **286–295 MiB** vs 758 / 1,871 MiB (**2.6× / 6.3× lower**) |
| **`bundle`** (12.6 MB) | 11.82 s | 5.12 s / 4.96 s | **4.99 s** / **4.90 s** | **−2.6%** / **−1.2%** | **−57.8%** / **−58.5%** | 16.40 s / 8.68 s | **3.35×** / **1.77×** | **720–798 MiB** vs 2,151 / 2,343 MiB (**2.7× / 2.9× lower**) |
| **`downlevel-whitespace`** | 0.86 s | 0.59 s / 0.48 s | **0.58 s** / **0.49 s** | **−1.4%** / **+0.5%** | **−32.6%** / **−43.0%** | 2.55 s / 0.60 s | **5.20×** / **1.22×** | **205–220 MiB** vs 553 / 1,437 MiB (**2.5× / 6.5× lower**) |
| **`downlevel-simple`** | 2.39 s | 1.23 s / 1.14 s | **1.23 s** / **1.14 s** | **0.0%** / **0.0%** | **−48.5%** / **−52.3%** | 4.68 s / 1.60 s | **4.11×** / **1.40×** | **251–261 MiB** vs 613 / 1,489 MiB (**2.3× / 5.7× lower**) |
| **`checks-only`** | 4.40 s | 2.80 s / 2.41 s | **2.72 s** / **2.33 s** | **−2.8%** / **−3.2%** | **−38.2%** / **−47.0%** | 6.55 s / 3.20 s | **2.81×** / **1.37×** | **672–695 MiB** vs 1,591 / 3,562 MiB (**2.3× / 5.1× lower**) |

#### Large & Mega Real-World `ADVANCED` Benchmarks

| Benchmark | Workload (`--compilation_level=ADVANCED`) | Round 8 (`f1c0cd2e` `t=1`/`t=4`) | **Round 9 (`t=1` / `t=4`)** | Round 9 vs Round 8 (`t=1`/`t=4`) | Cold Java | Warm Java | **`jsgo` vs Cold Java (`t=1`/`t=4`)** | **`jsgo` vs Warm Java (`t=1`/`t=4`)** | Peak RSS: `jsgo` vs Cold / Warm Java |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| **`typed-mega`** | **44.6 MB, 4,385 files** (10 Closure `ADVANCED` codebases, 94.1% typed → 14.33 MB JS + 22.29 MB map) | 21.97 s / 20.79 s | **21.37 s** / **20.13 s** | **−2.7%** / **−3.2%** | 38.18 s | 29.53 s | **1.79×** / **1.90× faster** | **1.38×** / **1.47× faster** (−9.40 s) | **1,883–1,903 MiB** vs 6,608 / 6,881 MiB (**3.5× / 3.6× lower**) |
| **`bundle-mega`** | **36.1 MB, 26 UMD libraries** (27,597 warnings → 13.09 MB JS + 33.99 MB map) | 28.27 s / 27.41 s | **27.55 s** / **26.70 s** | **−2.5%** / **−2.6%** | 66.71 s | 53.60 s | **2.42×** / **2.50× faster** | **1.95×** / **2.01× faster** (−26.90 s) | **2,574–2,804 MiB** vs 4,421 / 6,775 MiB (**1.6× / 2.4× lower**) |
| **`typed-full`** | **28.5 MB, 1,818 files** (Closure Library + Shaka Player 5.3.1 + Blockly 8.0.5 + app, 97.3% typed → 6.75 MB JS + 13.31 MB map) | 12.74 s / 12.06 s | **12.36 s** / **11.73 s** | **−3.0%** / **−2.7%** | 31.60 s | 20.21 s | **2.56×** / **2.69× faster** | **1.64×** / **1.72× faster** (−8.48 s) | **1,142–1,185 MiB** vs 1,603 / 4,724 MiB (**1.4× / 4.0× lower**) |
| **`bundle-advanced-ts`** | **12.6 MB, 4 libraries incl. 8.9 MB `typescript.js`** (10,751 warnings → 3.25 MB JS + 9.39 MB map) | 7.79 s / 7.67 s | **7.64 s** / **7.49 s** | **−1.9%** / **−2.3%** | 20.78 s | 12.41 s | **2.72×** / **2.77× faster** | **1.62×** / **1.66× faster** (−4.92 s) | **799–906 MiB** vs 2,422 / 2,850+ MiB (**2.7× / 3.1× lower**) |
| **`shaka-blockly`** | **8.6 MB, 847 files** (Shaka Player 5.3.1 + Blockly 8.0.5, 93.9% typed → 1.86 MB JS + 4.94 MB map) | 5.30 s / 5.08 s | **5.20 s** / **4.96 s** | **−2.1%** / **−2.4%** | 15.26 s | 8.67 s | **2.93×** / **3.08× faster** | **1.67×** / **1.75× faster** (−3.71 s) | **547–579 MiB** vs 1,341 / 2,852 MiB (**2.3× / 4.9× lower**) |
| **`bundle-advanced`** | **5.4 MB, 11 UMD libraries** (three, lodash, react-dom, rxjs, d3, moment, vue, jquery, handlebars, immutable, underscore → 1.88 MB JS + 4.35 MB map) | 4.00 s / 3.83 s | **3.86 s** / **3.68 s** | **−3.3%** / **−3.8%** | 11.30 s | 6.79 s | **2.93×** / **3.07× faster** | **1.76×** / **1.85× faster** (−3.11 s) | **454–531 MiB** vs 1,600 / 3,795 MiB (**3.0× / 7.1× lower**) |

## Round 10: `scopeArena` Chunk Retention (`−60%` `allocScope`), `DrainParseArenaPool` (`−637 MiB` / `−21%` Peak RSS), `RemoveUnusedCode` Polyfill Bitmask, Zero-Alloc `GuardedCallbackBase`, and Early Dataflow Bailouts (2026-10-20)

Profiling Round 9 (`e5d0987a`) on `bundle-mega` (`6.64 GB` `alloc_space`) and `typed-mega` (`5.21 GB` `alloc_space`) uncovered five remaining bottlenecks:

1. **`scopeArena` Chunk Truncation Fix & Pass Pooling (`Subagent 10A`, `scope.go`, `typed_scope_creator.go`, and 9 pass files):**
   - **Eliminated `~372 MB` of `scopeArena.allocScope` Churn on `bundle-mega`:** In Round 9, `maxScopeChunkSize = 2048, maxScopeChunks = 96` only retained `196,608` `Scope`s, whereas `bundle-mega` creates ~350,000 `Scope`s per pass (`FUNCTION` + `BLOCK` scopes plus inner block scopes). On every `release()`, `a.scopeChunks[:96]` truncated ~80 chunks (`~12 MB`) of `[]Scope` slabs, forcing the next pass to allocate and zero new OS pages. Raising `maxScopeChunkSize = 4096, maxScopeChunks = 128` (`524,288` `Scope`s) and `maxVarChunks = maxSlotChunks = maxHashChunks = 128`, plus wiring `EnableArenaPool()` in `TypedScopeCreator.CreateInitialScope`, `FindModuleDependencies`, `CheckAccessControls`, `CheckTypeImportCodeReferences`, `CheckClosureImports`, `ClosureRewriteModule`, `ConstCheck`, `Es6CheckModule`, `CheckSideEffects`, and `Denormalize` cut `scopeArena.allocScope` by **−60.4%** on `bundle-mega` (`615.84 MB → 244 MB`) and **−61.2%** on `typed-mega` (`487 MB → 189 MB`).
2. **Unpinning `ParseArena` After `Compiler.ParseInputs` (`Subagent 10B`, `parse_arena.go`, `compiler.go`):**
   - Added `parser.DrainParseArenaPool()` and called `defer parser.DrainParseArenaPool()` in `Compiler.ParseInputs()`, allowing the 4 `ParseArena`s to be shared across worker goroutines during parallel input parsing and then immediately reclaimed by the GC before `TypeCheck` and `PhaseOptimizer` begin—reducing Peak RSS by **−636.5 MiB (−21.1%)** on `bundle-mega`, **−250.4 MiB (−30.3%)** on `bundle`, **−219.4 MiB (−24.5%)** on `bundle-advanced-ts`, and **−122.8 MiB (−6.1%)** on `typed-mega`.
3. **Polyfill Bitmask Filter in `RemoveUnusedCode` & Zero-Alloc `GuardedCallbackBase` (`Subagent 10C`, `remove_unused_code.go`, `guarded_callback.go`, `polyfill_usage_finder.go`):**
   - Added `polyfillLenMask [256]uint64` and `polyfillEndMask [256]uint64` `(firstByte/lastByte, length)` bitmasks on `RemoveUnusedCode` so `TraverseNameNode` and `TraverseNormalOrOptChainGetProp` skip `r.polyfills.ContainsKey(name)` string hashing (`jutil.StringHashCode`) on 98%+ of `NAME` and `GETPROP` nodes.
   - Replaced `SetMultimap`/`HashMultiset`/`ListMultimap` in `GuardedCallbackBase` with lazily allocated native Go maps, added `emptyDepth int` and `[256]uint8` `guardedCallback_tokenFlags` to bypass `descendInto` and `NodeUtil_IsStatement` when at `EMPTY` context, and added `sc.EnableArenaPool()` in `PolyfillUsageFinder.TraverseExcludingGuarded` / `TraverseIncludingGuarded`.
4. **`DataFlowAnalysis_ComputeEscaped`, `FlowSensitiveInlineVariables`, `CoalesceVariableNames`, and `NodeUtil_GetAllVarsDeclaredInFunction` Fast Paths (`Subagent 10D`, `data_flow_analysis.go`, `flow_sensitive_inline_variables.go`, `coalesce_variable_names.go`, `node_util.go`):**
   - Added early bailouts in `DataFlowAnalysis_ComputeEscaped` (`dataFlowAnalysis_escapedFinder`), `FlowSensitiveInlineVariables.EnterScope`, and `CoalesceVariableNames.EnterScopeWithCfg` when all local variables in a function are already escaped or when fewer than 2 unescaped coalescable candidates exist (skipping CFG construction, reaching-definition / liveness dataflow analysis, and interference graph coloring), and pruned non-declaring statement subtrees in `NodeUtil_GetAllVarsDeclaredInFunction`.
5. **`LiveVariablesAnalysis` & `DeadAssignmentsElimination` Fast Paths (`Subagent 10E`, `live_variables_analysis.go`, `dead_assignments_elimination.go`):**
   - Replaced `l.GetCfg().GetOutEdges(node)` map lookups in `LiveVariablesAnalysis.FlowThrough` with `l.curDiGraphNode.GetOutEdges()`, added a linear scan fast path (`useSmallVarScan`) for functions with $\le 16$ variables without shadowed names, and reordered `DeadAssignmentsElimination.EliminateDeadAssignments` to skip CFG construction and liveness analysis when a function has 0 unescaped local variables.

Total `alloc_space` dropped to **6.21 GB** on `bundle-mega` (**−63.8% vs Round 5's `17.16 GB`**) and **5.05 GB** on `typed-mega` (**−56.5% vs Round 5's `11.61 GB`**), and `typed-mega` (`t=4`) broke the **20-second barrier** (`19.83 s`, **1.93× faster than cold Java and 1.49× faster than warm Java**), while `bundle-mega` (`t=4`) reached **25.70 s** (**2.60× faster than cold Java and 2.09× faster than warm Java**).

### Final measured results — Round 10 (`t=1` and `t=4`, interleaved, N=3, PGO enabled)

#### Standard 5 Benchmarks

| Benchmark | Original (`c9872522`) | Round 9 (`e5d0987a` `t=1`/`t=4`) | **Round 10 (`t=1` / `t=4`)** | Round 10 vs Round 9 (`t=1`/`t=4`) | Total vs Original (`t=1`/`t=4`) | Cold / Warm Java (`t=4`) | **`jsgo` vs Cold / Warm Java (`t=4`)** | Peak RSS: `jsgo` vs Cold / Warm Java |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| **`typed-app`** | 2.86 s | 1.59 s / 1.46 s | **1.57 s** / **1.44 s** | **−1.2%** / **−0.9%** | **−45.1%** / **−49.7%** | 6.14 s / 2.18 s | **4.26×** / **1.51×** | **256–260 MiB** vs 758 / 1,871 MiB (**2.9× / 7.2× lower**, −11.3% RSS) |
| **`bundle`** (12.6 MB) | 11.82 s | 4.99 s / 4.88 s | **4.89 s** / **4.82 s** | **−2.0%** / **−1.3%** | **−58.6%** / **−59.2%** | 16.40 s / 8.68 s | **3.40×** / **1.80×** | **577 MiB** vs 2,151 / 2,343 MiB (**3.7× / 4.1× lower**, −30.3% RSS) |
| **`downlevel-whitespace`** | 0.86 s | 0.58 s / 0.49 s | **0.58 s** / **0.48 s** | **−1.2%** / **−1.7%** | **−32.6%** / **−44.2%** | 2.55 s / 0.60 s | **5.31×** / **1.25×** | **165–187 MiB** vs 553 / 1,437 MiB (**3.0× / 8.7× lower**, −25.0% RSS) |
| **`downlevel-simple`** | 2.39 s | 1.26 s / 1.16 s | **1.21 s** / **1.12 s** | **−3.8%** / **−3.9%** | **−49.4%** / **−53.1%** | 4.68 s / 1.60 s | **4.18×** / **1.43×** | **199–206 MiB** vs 613 / 1,489 MiB (**3.0× / 7.5× lower**, −22.9% RSS) |
| **`checks-only`** | 4.40 s | 2.72 s / 2.33 s | **2.73 s** / **2.33 s** | **+0.1%** / **+0.1%** | **−38.0%** / **−47.0%** | 6.55 s / 3.20 s | **2.81×** / **1.37×** | **611–640 MiB** vs 1,591 / 3,562 MiB (**2.5× / 5.6× lower**, −12.0% RSS) |

#### Large & Mega Real-World `ADVANCED` Benchmarks

| Benchmark | Workload (`--compilation_level=ADVANCED`) | Round 9 (`e5d0987a` `t=1`/`t=4`) | **Round 10 (`t=1` / `t=4`)** | Round 10 vs Round 9 (`t=1`/`t=4`) | Cold Java | Warm Java | **`jsgo` vs Cold Java (`t=1`/`t=4`)** | **`jsgo` vs Warm Java (`t=1`/`t=4`)** | Peak RSS: `jsgo` vs Cold / Warm Java |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| **`typed-mega`** | **44.6 MB, 4,385 files** (10 Closure `ADVANCED` codebases, 94.1% typed → 14.33 MB JS + 22.29 MB map) | 21.40 s / 20.29 s | **20.96 s** / **19.83 s** | **−2.0%** / **−2.3%** | 38.18 s | 29.53 s | **1.82×** / **1.93× faster** | **1.41×** / **1.49× faster** (−9.70 s) | **1,804–1,908 MiB** vs 6,608 / 6,881 MiB (**3.5× / 3.6× lower**, −123 MiB vs R9) |
| **`bundle-mega`** | **36.1 MB, 26 UMD libraries** (27,597 warnings → 13.09 MB JS + 33.99 MB map) | 27.57 s / 26.73 s | **26.51 s** / **25.70 s** | **−3.9%** / **−3.8%** | 66.71 s | 53.60 s | **2.52×** / **2.60× faster** | **2.02×** / **2.09× faster** (−27.90 s) | **2,387–2,444 MiB** vs 4,421 / 6,775 MiB (**1.8× / 2.8× lower**, −637 MiB / −21.1% vs R9) |
| **`typed-full`** | **28.5 MB, 1,818 files** (Closure Library + Shaka Player 5.3.1 + Blockly 8.0.5 + app, 97.3% typed → 6.75 MB JS + 13.31 MB map) | 12.37 s / 11.70 s | **12.13 s** / **11.48 s** | **−2.0%** / **−1.9%** | 31.60 s | 20.21 s | **2.61×** / **2.75× faster** | **1.67×** / **1.76× faster** (−8.73 s) | **1,036–1,057 MiB** vs 1,603 / 4,724 MiB (**1.5× / 4.5× lower**, −117 MiB vs R9) |
| **`bundle-advanced-ts`** | **12.6 MB, 4 libraries incl. 8.9 MB `typescript.js`** (10,751 warnings → 3.25 MB JS + 9.39 MB map) | 7.62 s / 7.49 s | **7.54 s** / **7.41 s** | **−1.0%** / **−1.0%** | 20.78 s | 12.41 s | **2.76×** / **2.80× faster** | **1.65×** / **1.67× faster** (−5.00 s) | **633–677 MiB** vs 2,422 / 2,850+ MiB (**3.6× / 4.2× lower**, −219 MiB / −24.5% vs R9) |
| **`shaka-blockly`** | **8.6 MB, 847 files** (Shaka Player 5.3.1 + Blockly 8.0.5, 93.9% typed → 1.86 MB JS + 4.94 MB map) | 5.19 s / 4.95 s | **5.08 s** / **4.84 s** | **−2.1%** / **−2.2%** | 15.26 s | 8.67 s | **3.00×** / **3.15× faster** | **1.71×** / **1.79× faster** (−3.83 s) | **513–528 MiB** vs 1,341 / 2,852 MiB (**2.5× / 5.4× lower**, −47 MiB vs R9) |
| **`bundle-advanced`** | **5.4 MB, 11 UMD libraries** (three, lodash, react-dom, rxjs, d3, moment, vue, jquery, handlebars, immutable, underscore → 1.88 MB JS + 4.35 MB map) | 3.86 s / 3.69 s | **3.83 s** / **3.66 s** | **−0.8%** / **−0.9%** | 11.30 s | 6.79 s | **2.95×** / **3.09× faster** | **1.77×** / **1.86× faster** (−3.13 s) | **395–404 MiB** vs 1,600 / 3,795 MiB (**4.0× / 9.4× lower**, −123 MiB / −23.3% vs R9) |

## Round 11: 64-Bit `Scope` Bloom Filter, Lazy `GlobalNamespace` Scopes (`−55.5%` `allocScope`), `RemoveUnusedCode` Slab Retention (`−166 MB`), Pooled `ReferenceCollector` Pointer Slices (`−80 MB`), and `OptimizeCalls` / `InlineFunctions` / `InlineSimpleMethods` / `ColorFindPropertyReferences` Pooling (2026-10-20)

Profiling Round 10 (`9c895b43`) on `bundle-mega` (`6.21 GB` `alloc_space`) and `typed-mega` (`5.05 GB` `alloc_space`) uncovered five remaining bottlenecks:

1. **$O(1)$ 64-Bit 2-Bit-per-Name Bloom Filter (`varBloom uint64`) on `Scope` (`Subagent 11A`, `scope.go`, `syntactic_scope_creator.go`):**
   - Added `varBloom uint64` on `Scope` and an $O(1)$ 3-byte hash `scopeBloomMask(s string) uint64` (using `(s[0], s[n-1], s[n>>1], len(s))` with zero loops), gating every ancestor scope in `Scope.GetVar` and `Scope.HasSlot` on `(s.varBloom & bloom) == bloom` while keeping `scopeNameHash(name)` lazy (`if s.varHash != nil`). This rejects $>98\%$ of local ancestor scopes in a single bitwise instruction without touching `s.varList` or `s.varHash` and without running `scopeNameHash` on local scopes (`s.varHash == nil`).
   - Deduplicated duplicate checks in `Scope.Declare`, `Scope.DeclareInternal`, `Scope.CanDeclare`, and `SyntacticScopeCreator_ScopeScanner.DeclareVar` via `alreadyMaybePresent := (s.varBloom & bloom) == bloom`.
2. **Lazy `t.GetScope()` / `t.GetChunk()` in `GlobalNamespace` & `Scope` Arena Pooling (`Subagent 11B`, `global_namespace.go`, `process_defines.go`, `inline_and_collapse_properties.go`, `es6_rewrite_class.go`, `analyze_prototype_properties.go`, `cross_chunk_method_motion.go`, `es6_rewrite_block_scoped_declaration.go`, `gather_module_metadata.go`):**
   - Replaced eager `b.Collect(t.GetScope(), t.GetChunk(), n)` in `GlobalNamespace_BuildGlobalNamespace.ShouldTraverse` with `b.collect(nil, nil, t, n)`, resolving `t.GetScope()` only when a candidate global token (`NAME`, `GETPROP`, `STRING_KEY`, `MEMBER_*`) passes initial filters and resolving `t.GetChunk()` only when a reference is actually recorded.
   - Pooled `SyntacticScopeCreator` arenas across `GlobalNamespace` (with `Release()` wired in `ProcessDefines`, `InlineAndCollapseProperties`, and `Es6RewriteClass`), `AnalyzePrototypeProperties` (released in `CrossChunkMethodMotion`), `Es6RewriteBlockScopedDeclaration` (also skipping `LoopClosureTransformer` when `letConsts.IsEmpty()`), and `GatherModuleMetadata`.
   - Together, `scopeArena.allocScope` dropped another **−55.5%** in Round 11 (`244 MB → 108.73 MB` on `bundle-mega`, **−82.3% vs Round 9's `615.84 MB`**).
3. **`RemoveUnusedCode` Slab Retention (`maxRetainedChunks = 160`), De-Virtualized `*RemoveUnusedCode_RealVarInfo`, and Traversal Fast Paths (`Subagent 11C`, `remove_unused_code.go`):**
   - Raised `maxRetainedChunks` in `rucSlab.reset()` from `40` (`143,296` items) to `160` (`655,360` items), completely eliminating the `144.05 MB` of `realVarInfoSlab.New` and `21.70 MB` of `assignSlab.New` chunk truncation across the 10 `RemoveUnusedCode` passes on `bundle-mega`.
   - De-virtualized `RemoveUnusedCode_VarInfo` to concrete `*RemoveUnusedCode_RealVarInfo`, converted `varPtrCache` (`[4096]`), `scopeVarCache` (`[512]`), and `pinnedPropCache` (`[256]`) to pointer-backed fixed-size arrays on `removeUnusedCodeArena` (eliminating 48 KB struct-by-value copies per pass), eliminated the duplicate `NodeUtil_EvaluatesToLocalValue` subtree walk in `RealVarInfo.AddRemovable`, and added fast paths in `TraverseCall`, `ConsiderForAccessorSideEffects`, and `MayHaveSideEffects`.
4. **Pooled `ReferenceCollector` Pointer Slices (`refPtrChunks`), Open-Addressed `hoistedTable`, and `InlineVariables` Compact Delayed Analysis (`Subagent 11D`, `reference_collector.go`, `inline_variables.go`):**
   - Eliminated the `80.20 MB` (`bundle-mega`) / `63.48 MB` (`typed-mega`) `ReferenceCollection.Add` -> `ArrayList.Add` heap allocations by sub-allocating grown `[]*Reference` buffers from pooled `refPtrChunks` on `referenceCollectorArena`.
   - Replaced `collectedHoistedMap map[*rhino.Node]struct{}` with a pooled open-addressed pointer table (`hoistedTable`) + `lastHoisted` MRU cache, eliminated `c.blockStack` heap allocations in `maybeJumpToHoistedFunctionSlow`, and replaced 200-byte `scopeExpertSlab` copies in `InlineVariables` with a 16-byte `inlineVariables_delayedNegativeAnalysis` slab.
5. **Pooled `OptimizeCalls_ReferenceMap` (`−123 MB`), `InlineSimpleMethods`, `InlineFunctions`, and Memoized `ColorFindPropertyReferences` (`Subagent 11E`, `optimize_calls.go`, `inline_simple_methods.go`, `inline_functions.go`, `disambiguate_color_find_property_references.go`):**
   - Pooled `OptimizeCalls_ReferenceMap` and its multi-chunk `optimizeCalls_refListSlot` arena across `OptimizeCalls.Process` runs, eliminating `123 MB` of `LinkedHashMap.PutNew` and `ArrayList.Add` allocations on `bundle-mega`.
   - Pooled `ismPoolState` (`props map[string]*ismPropSlot` + `slotChunks`) in `InlineSimpleMethods` (`−31 MB`), pooled `inlineFunctionsPoolState` and cached `bodyAstSize` / `lastContainingFnSize` in `InlineFunctions` (`−20 MB`), and memoized `createNode(owner *colors.Color)` + eliminated per-node closures in `ColorFindPropertyReferences` (`−56 MB`).

Total `alloc_space` dropped another **−410 MB (−6.6%)** on `bundle-mega` (`6.21 GB → 5.80 GB`, **−66.2% vs Round 5's `17.16 GB`**) and **−270 MB (−5.3%)** on `typed-mega` (`5.05 GB → 4.78 GB`, **−58.8% vs Round 5's `11.61 GB`**), while speeding up 21 of 22 benchmark configurations.

### Final measured results — Round 11 (`t=1` and `t=4`, interleaved, N=3, PGO enabled)

#### Standard 5 Benchmarks

| Benchmark | Original (`c9872522`) | Round 10 (`9c895b43` `t=1`/`t=4`) | **Round 11 (`t=1` / `t=4`)** | Round 11 vs Round 10 (`t=1`/`t=4`) | Total vs Original (`t=1`/`t=4`) | Cold / Warm Java (`t=4`) | **`jsgo` vs Cold / Warm Java (`t=4`)** | Peak RSS: `jsgo` vs Cold / Warm Java |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| **`typed-app`** | 2.86 s | 1.57 s / 1.46 s | **1.55 s** / **1.43 s** | **−0.8%** / **−1.8%** | **−45.8%** / **−50.0%** | 6.14 s / 2.18 s | **4.29×** / **1.52×** | **242–246 MiB** vs 758 / 1,871 MiB (**3.1× / 7.7× lower**, −14.4% RSS) |
| **`bundle`** (12.6 MB) | 11.82 s | 4.90 s / 4.79 s | **4.81 s** / **4.70 s** | **−1.8%** / **−2.0%** | **−59.3%** / **−60.2%** | 16.40 s / 8.68 s | **3.49×** / **1.85×** | **558–597 MiB** vs 2,151 / 2,343 MiB (**3.6× / 4.2× lower**) |
| **`downlevel-whitespace`** | 0.86 s | 0.58 s / 0.48 s | **0.56 s** / **0.47 s** | **−2.0%** / **−1.8%** | **−34.9%** / **−45.3%** | 2.55 s / 0.60 s | **5.43×** / **1.28×** | **162–171 MiB** vs 553 / 1,437 MiB (**3.2× / 8.9× lower**) |
| **`downlevel-simple`** | 2.39 s | 1.21 s / 1.12 s | **1.20 s** / **1.10 s** | **−1.5%** / **−1.4%** | **−49.8%** / **−54.0%** | 4.68 s / 1.60 s | **4.25×** / **1.45×** | **205–211 MiB** vs 613 / 1,489 MiB (**2.9× / 7.3× lower**) |
| **`checks-only`** | 4.40 s | 2.74 s / 2.32 s | **2.72 s** / **2.33 s** | **−0.5%** / **+0.3%** | **−38.2%** / **−47.0%** | 6.55 s / 3.20 s | **2.81×** / **1.37×** | **623–664 MiB** vs 1,591 / 3,562 MiB (**2.4× / 5.7× lower**) |

#### Large & Mega Real-World `ADVANCED` Benchmarks

| Benchmark | Workload (`--compilation_level=ADVANCED`) | Round 10 (`9c895b43` `t=1`/`t=4`) | **Round 11 (`t=1` / `t=4`)** | Round 11 vs Round 10 (`t=1`/`t=4`) | Cold Java | Warm Java | **`jsgo` vs Cold Java (`t=1`/`t=4`)** | **`jsgo` vs Warm Java (`t=1`/`t=4`)** | Peak RSS: `jsgo` vs Cold / Warm Java |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| **`typed-mega`** | **44.6 MB, 4,385 files** (10 Closure `ADVANCED` codebases, 94.1% typed → 14.33 MB JS + 22.29 MB map) | 20.95 s / 19.80 s | **20.74 s** / **19.61 s** | **−1.0%** / **−1.0%** | 38.18 s | 29.53 s | **1.84×** / **1.95× faster** | **1.42×** / **1.51× faster** (−9.92 s) | **1,821–1,864 MiB** vs 6,608 / 6,881 MiB (**3.5× / 3.8× lower**, −84 MiB vs R10) |
| **`bundle-mega`** | **36.1 MB, 26 UMD libraries** (27,597 warnings → 13.09 MB JS + 33.99 MB map) | 26.44 s / 25.80 s | **26.15 s** / **25.39 s** | **−1.1%** / **−1.6%** | 66.71 s | 53.60 s | **2.55×** / **2.63× faster** | **2.05×** / **2.11× faster** (−28.21 s) | **2,288–2,410 MiB** vs 4,421 / 6,775 MiB (**1.8× / 3.0× lower**, −82 MiB vs R10) |
| **`typed-full`** | **28.5 MB, 1,818 files** (Closure Library + Shaka Player 5.3.1 + Blockly 8.0.5 + app, 97.3% typed → 6.75 MB JS + 13.31 MB map) | 12.16 s / 11.47 s | **11.96 s** / **11.35 s** | **−1.6%** / **−1.1%** | 31.60 s | 20.21 s | **2.64×** / **2.78× faster** | **1.69×** / **1.78× faster** (−8.86 s) | **1,098–1,113 MiB** vs 1,603 / 4,724 MiB (**1.4× / 4.3× lower**) |
| **`bundle-advanced-ts`** | **12.6 MB, 4 libraries incl. 8.9 MB `typescript.js`** (10,751 warnings → 3.25 MB JS + 9.39 MB map) | 7.53 s / 7.41 s | **7.40 s** / **7.28 s** | **−1.8%** / **−1.8%** | 20.78 s | 12.41 s | **2.81×** / **2.85× faster** | **1.68×** / **1.70× faster** (−5.13 s) | **627–673 MiB** vs 2,422 / 2,850+ MiB (**3.6× / 4.5× lower**) |
| **`shaka-blockly`** | **8.6 MB, 847 files** (Shaka Player 5.3.1 + Blockly 8.0.5, 93.9% typed → 1.86 MB JS + 4.94 MB map) | 5.08 s / 4.84 s | **4.98 s** / **4.78 s** | **−2.1%** / **−1.3%** | 15.26 s | 8.67 s | **3.06×** / **3.19× faster** | **1.74×** / **1.81× faster** (−3.89 s) | **547–551 MiB** vs 1,341 / 2,852 MiB (**2.4× / 5.2× lower**) |
| **`bundle-advanced`** | **5.4 MB, 11 UMD libraries** (three, lodash, react-dom, rxjs, d3, moment, vue, jquery, handlebars, immutable, underscore → 1.88 MB JS + 4.35 MB map) | 3.83 s / 3.65 s | **3.77 s** / **3.60 s** | **−1.6%** / **−1.4%** | 11.30 s | 6.79 s | **3.00×** / **3.14× faster** | **1.80×** / **1.89× faster** (−3.19 s) | **389–395 MiB** vs 1,600 / 3,795 MiB (**4.1× / 9.8× lower**) |

## Round 12: `ReferenceCollector` Full Mega-Chunk Retention (`−190 MB`), `AnalyzePrototypeProperties` Edge Deduplication (`−93 MB`), `InferConsts` Per-Function State Flushing (`−132 MB`), `RenameVars` Chunked Slabs (`−83 MB`), Lazy `GetScope` Guards, `OptimizeCalls` Multimap Elimination, and `TypedScope` / `TypedScopeCreator` Optimizations (2026-10-20)

Profiling Round 11 (`88070d50`) on `bundle-mega` (`5.80 GB` `alloc_space`) and `typed-mega` (`4.89 GB` `alloc_space`) identified five disjoint optimization targets:

1. **`ReferenceCollector` Full Mega-Chunk Retention & `AnalyzePrototypeProperties` Duplicate Edge Elimination (`Subagent 12A`, `reference_collector.go`, `analyze_prototype_properties.go`, `graph_linked_directed_graph.go`):**
   - Raised `allocReference` chunk growth cap to `8192` (`maxRetainedRefChunks = 256` = ~2M `Reference`s) and `allocBasicBlock` / `allocReferenceCollection` chunk growth caps to `4096` (`maxRetainedBlkChunks = 128`, `maxRetainedColChunks = 128`, `maxRetainedRefPtrChunks = 256`), preventing `VariableReferenceCheck`'s ~1.35M references and ~350k basic blocks from truncating tail chunks before `InlineVariables` and `InlineObjectLiterals`.
   - Added chunked slab allocation (`nameInfoChunks`) and inline `declBuf [2]AnalyzePrototypeProperties_Symbol` on `AnalyzePrototypeProperties_NameInfo`, cached `graphNode` directly on `NameInfo`, and deduplicated `(globalNode, chunk, dest)`, `(externNode, chunk, dest)`, and repeated `(lastSrc, lastChunk, dest)` edges via `connectSymbolNodes` before calling `a.symbolGraph.ConnectNodes`.
2. **`InferConsts` Per-Outer-Function State Flushing & `RenameVars` Chunked Slabs (`Subagent 12B`, `infer_consts.go`, `rename_vars.go`):**
   - Split `inferConstsCollector` state into `localStates` (positive `v.nameHash`) for function-local variables and `globalStates` (negative `v.nameHash`), and immediately flushed and reset `localStates` via `flushLocalStates()` whenever `funcDepth == 0 && jumpDepth == 0` upon exiting a top-level function, preventing `states` from doubling to ~650,000 entries (`132.36 MB`) on `bundle-mega`.
   - Replaced flat `localRefs` and `globalRefs` slices in `RenameVars` with chunked slabs (`localRefChunks` growing up to `16384` and `globalRefChunks` growing up to `8192`) and indexed `varAssignments` once per unique `RenameVars_Assignment` inside `allocAssignment`.
3. **Lazy `t.GetScope()` Guards & `FunctionToBlockMutator` Local Arena Reuse (`Subagent 12C`, `pure_function_identifier.go`, `function_to_block_mutator.go`, `check_type_import_code_references.go`, `es6_check_module.go`, `find_module_dependencies.go`):**
   - Replaced `!t.GetScope().IsFunctionScope()` in `PureFunctionIdentifier_FunctionBodyAnalyzer.EnterScope` with `!t.GetScopeRoot().IsFunction()` and deferred `t.GetScope()` in `ExitScope` until after `functionEntry.root == nil` and `taintedVars.IsEmpty()` checks.
   - Embedded a dedicated reusable `localArena scopeArena` and `localScopeCreator SyntacticScopeCreator` on `FunctionToBlockMutator` for `MakeLocalNamesUnique` during `InlineFunctions`.
   - Added a `collectRequireTypeNames` pre-scan in `CheckTypeImportCodeReferences.Process`, skipped non-module `SCRIPT` subtrees in `Es6CheckModule.ShouldTraverse`, and guarded `FindModuleDependencies.ExitScope` with `f.dynamicImportScope != nil`.
4. **`OptimizeCalls` Backing Slice Retention & `GetFunctionNodes` Multimap Elimination (`Subagent 12D`, `optimize_calls.go`, `optimize_parameters.go`, `optimize_returns.go`, `devirtualize_methods.go`):**
   - Raised `ocRefSlotChunkSize` to `4096`, retained `slot.list` backing arrays up to capacity `65536` in `releaseOptimizeCallsRefMap`, and pre-sized `props` in `acquireOptimizeCallsRefMap` and `EnterScope`.
   - Replaced `OptimizeCalls_ReferenceMap_GetFunctionNodes` multimap allocations in `OptimizeParameters`, `OptimizeReturns`, and `DevirtualizeMethods` with zero-allocation `OptimizeCalls_ReferenceMap_AppendFunctionNodes` / `optimizeCalls_getSingleFunctionNode` helpers and reusable scratch slices.
5. **`TypedScope` Closure/Struct Optimization & `FirstOrderFunctionAnalyzer` In-Place `var_.nameHash` Tracking (`Subagent 12E`, `typed_scope.go`, `typed_scope_creator.go`, `type_inference.go`):**
   - Eliminated `GetDeclarativelyUnboundVarsWithoutTypes` heap closure allocations in `TypeInference.InferDeclarativelyUnboundVarsWithoutTypes`, removed the unused 96-byte embedded `varSlots [4]` array from `TypedScope` (shrinking `TypedScope` from 216 B to 120 B), and increased `typedScopeArena` slab caps.
   - Tracked local variable assignment count and escaped status directly on syntactic `var_.nameHash` during `FirstOrderFunctionAnalyzer.Visit`, populating a single `varFlags` map in `ExitScope` only when `h != 1`, and pre-sized `c.memoized` at the end of `CreateInitialScope` using `c.scopeCount`.

Total `alloc_space` dropped another **−552 MB (−9.5%)** on `bundle-mega` (`5.80 GB → 5.25 GB`, **−69.4% vs Round 5's `17.16 GB`**) and **−292 MB (−6.0%)** on `typed-mega` (`4.89 GB → 4.60 GB`, **−60.4% vs Round 5's `11.61 GB`**).

### Final measured results — Round 12 (`t=1` and `t=4`, 5 warmups + 10 Latin-square-interleaved runs, PGO enabled)

See the full `Median [95% CI]` and `Mean ± SD` tables in [benchmarks.md](benchmarks.md#rounds-212-and-large--mega-real-world-advanced-benchmarks) and [README.md](../README.md).

## Round 13: `GlobalNamespace` Object-Literal LValue Caching (`−154 MB`), `TemplateTypeReplacer` / `TemplateTypeMap` Inline Buffers (`−237 MB`), `PropertyUseSites` Arena (`−46 MB`), `MustDef` 512-Var Wide Bitmasks & Multi-Chunk Arenas (`−48 MB`), `FirstOrderFunctionAnalyzer` Single-Pass Traversal (`−61 MB`), and `CodePrinter` Pooling (`−65 MB`) (2026-10-21)

Profiling Round 12 (`b44f03c1`) on `typed-mega` (`4.60 GB` `alloc_space`) and `bundle-mega` (`5.25 GB` `alloc_space`) identified five disjoint optimization targets:

1. **`GlobalNamespace` Object-Literal LValue Caching, Top-Variable Scope Pre-Filtering, and `refSliceBuf` Slab (`Subagent 13A`, `global_namespace.go`):**
   - Added an 8-entry MRU cache (`[8]globalNamespace_objLitCacheEntry`) and `resolveObjectLit` on `GlobalNamespace_BuildGlobalNamespace` to memoize the LValue / `Object.defineProperties` resolution (`prefix`, `root`, `scope`, `state`) per `OBJECTLIT` node across all sibling keys and nested object literals.
   - Added `globalNamespace_topNameOfQName` to extract the top variable name of an LValue or qualified name node without allocating qualified name strings, checking `getRootNodeOfTopName` and `IsTopLevelScopeRoot(root)` **before** calling `GetQualifiedNameString()` or `nodeUtilConcatDot` (`−123.19 MB` on `typed-mega`).
   - Added `confirmExisting` on `GlobalNamespace_BuildGlobalNamespace`, made `n.refsIndex` lazily constructed via `ensureRefsIndex()`, and added a bump-allocated slice buffer (`refSliceBuf` / `allocRefSlice`) on `GlobalNamespace` (`−30.70 MB` on `typed-mega`).
2. **`TemplateTypeReplacer` & `TemplateTypeMap` Inline Buffers, Direct Dispatch, and Slice Sharing (`Subagent 13B`, `internal/rhino/jstype_template_type_replacer.go`, `internal/rhino/jstype_template_type_map.go`):**
   - Replaced `TemplateTypeReplacer.seenTypes` (`*jutil.LinkedHashSet[JSType]`, `148.08 MB` on `typed-mega`) with an 8-element inline LIFO buffer (`seenTypesBuf [8]JSType`) + overflow slice, inlined `seenAdd`/`seenRemove` cycle guards into `CaseFunctionType`, `CaseObjectType`, `CaseTemplatizedType`, `CaseUnionType`, and `CaseProxyObjectType`, and added direct `visitJSType` type-switch dispatch.
   - Replaced `visitedObjectTypes` (`*jutil.LinkedHashMap`) with `visitedObjBuf [4]visitedObjEntry` + lazy overflow map, and lazily initialized `results`/`builder`/`paramBuilder`/`RecordTypeBuilder` only when a type argument or property changes.
   - Stack-allocated `TemplateTypeReplacer` inside `initTemplateTypeMap`, reused `m.templateValues` directly for `m.resolvedTemplateValues` when all values resolve to themselves, embedded `valuesBuf [4]JSType` in `TemplateTypeMap`, and added shared immutable `BitSet`s (`sharedSubMapStartsEmpty`, `sharedSubMapStartsZero`).
3. **`Normalize` Direct `varList` Iteration, `StrictModeCheck` Bitmask Member Check, and `PropertyUseSites` Arena (`Subagent 13C`, `normalize.go`, `strict_mode_check.go`, `disambiguate_property_clustering.go`, `disambiguate_color_find_property_references.go`):**
   - Replaced `GetVarIterable()` range-over-func loops in `Normalize_NormalizeStatements.fastEnterScope` with direct `varList` iteration, deferred/pre-sized `s.nameUsage`, and avoided inserting generated `name + "$jscomp$" + id` strings into `s.nameUsage` unless an input variable containing `"$jscomp$"` is encountered (`−34.75 MB` on `bundle-mega`).
   - Replaced the 4 per-node `LinkedHashSet` allocations in `StrictModeCheck.CheckObjectLiteralOrClass` with a 0/1-child fast path, a 16-entry stack buffer `seenBuf [16]strictModeMemberEntry`, and a reusable `memberFlags map[string]uint8` (`−27.01 MB` on `typed-mega`, `−17.54 MB` on `bundle-mega`).
   - Replaced `useSites jutil.LinkedHashMap[*rhino.Node, *ColorGraphNode]` on `PropertyClustering` with `PropertyUseSites` (56 B inline with `[2]PropertyUseSite` + geometric unrolled linked chunks carved from a contiguous `[]PropertyUseSite` arena on `ColorFindPropertyReferences`) and added `ufNodeSlab` (`−46.19 MB` on `bundle-mega`, `−45.81 MB` on `typed-mega`).
4. **`MustBeReachingVariableDef` 512-Var Wide Bitmasks & Multi-Chunk Arenas + `InlineVariables` Scope-Indexed Slice (`Subagent 13D`, `must_def_map.go`, `must_be_reaching_variable_def.go`, `inline_variables.go`):**
   - Extended `mustDefLayout` and `mustDefMap` from 64 to 512 variables (`[8]uint64` wide bitmasks + `uint16` IDs) while keeping the single-word `uint64` fast path for $\le 64$ variables, and converted `mustDefSlab` and `defSlab` to multi-chunk arenas (`mustDefChunks`, `defChunks`) retained across `Reset()` (`−48.61 MB` on `bundle-mega`).
   - Replaced `handledVarsMap` in `InlineVariables_InliningBehavior` with a direct scope-indexed `handledScopeVars []InlineVariables_InlineVarAnalysis` slice indexed by `v.GetIndex()`, converted `delayedNegSlab`, `retryExpertSlab`, and `retryHandlerSlab` to multi-chunk arenas, and pooled `InlineVariables_InliningBehavior` (`−32.93 MB` on `typed-mega`).
5. **`PureFunctionIdentifier` Lazy `GetScope()`, `TypedScopeCreator` Single-Pass Unmemoized `FirstOrderFunctionAnalyzer`, and `CodePrinter` Pooling (`Subagent 13E`, `pure_function_identifier.go`, `typed_scope_creator.go`, `code_printer.go`):**
   - Pre-populated empty `Scope` instances on `t.scopes` outside all functions in `PureFunctionIdentifier`, deferred `traversal.GetScope()` in `visitLhsNode` and `Token_NAME` skiplisting until needed, and added `propSummariesByName` (`−31.06 MB` on `bundle-mega`).
   - Removed the `MemoizedScopeCreator` wrapper around `SyntacticScopeCreator` in `TypedScopeCreator.CreateInitialScope`, merged `IdentifyEnumsAndTypedefsAsNonNullable` into `FirstOrderFunctionAnalyzer.Visit` (eliminating a full AST traversal), and deferred `t.GetScope()` to `FirstOrderFunctionAnalyzer.ExitScope` (`−61.57 MB` on `bundle-mega`, `−29.54 MB` on `typed-mega`).
   - Added `sync.Pool`s (`codePrinter_compactPool` and `codePrinter_prettyPool`) for `CodePrinter_CompactCodePrinter` and `CodePrinter_PrettyCodePrinter` (`−65.69 MB` on `typed-mega`, `−21.06 MB` on `bundle-mega`).

Total `alloc_space` dropped another **−430 MB (−9.3%)** on `typed-mega` (`4.60 GB → 4.17 GB`, **−64.1% vs Round 5's `11.61 GB`**) and **−300 MB (−5.7%)** on `bundle-mega` (`5.25 GB → 4.95 GB`, **−71.2% vs Round 5's `17.16 GB`**), speeding up **all 22 of 22 benchmark × thread count configurations**!

### Final measured results — Round 13 (`t=1` and `t=4`, 1 warmup + 5 interleaved runs, PGO enabled)

| Benchmark | Round 12 (`b44f03c1` `t=1` / `t=4`) | **Round 13 (`t=1` / `t=4`)** | Round 13 vs Round 12 (`t=1` / `t=4`) | Peak RSS (`t=1` / `t=4`) | Parity |
|---|---:|---:|---:|---:|:---:|
| **`typed-mega`** (44.6 MB, 4,385 files) | 20.55 s / 19.37 s | **19.87 s** / **18.80 s** | **−3.27%** / **−2.93%** | 1,798 / 1,893 MiB (**−1.24%** `t=4`) | `MATCH` |
| **`bundle-mega`** (36.1 MB, 26 UMD libs) | 25.91 s / 25.00 s | **25.26 s** / **24.38 s** | **−2.51%** / **−2.46%** | 2,358 / 2,456 MiB | `MATCH` |
| **`typed-full`** (28.5 MB, 1,818 files) | 11.80 s / 11.19 s | **11.46 s** / **10.81 s** | **−2.93%** / **−3.40%** | 1,010 / 1,061 MiB (**−3.69%** `t=4`) | `MATCH` |
| **`bundle-advanced-ts`** (12.6 MB) | 7.31 s / 7.21 s | **7.14 s** / **7.01 s** | **−2.33%** / **−2.68%** | 600 / 629 MiB (**−4.92%** `t=1`) | `MATCH` |
| **`bundle`** (12.6 MB) | 4.76 s / 4.64 s | **4.67 s** / **4.53 s** | **−1.88%** / **−2.37%** | 555 / 547 MiB (**−2.97%** `t=4`) | `MATCH` |
| **`shaka-blockly`** (8.6 MB) | 4.91 s / 4.70 s | **4.76 s** / **4.53 s** | **−3.01%** / **−3.57%** | 489 / 511 MiB (**−3.58%** `t=4`) | `MATCH` |
| **`bundle-advanced`** (5.4 MB) | 3.72 s / 3.55 s | **3.63 s** / **3.47 s** | **−2.58%** / **−2.39%** | 371 / 393 MiB (**−2.31%** `t=1`) | `MATCH` |
| **`checks-only`** (11.3 MB) | 2.68 s / 2.29 s | **2.58 s** / **2.19 s** | **−3.67%** / **−4.45%** | 596 / 554 MiB (**−7.18%** `t=4`) | `MATCH` |
| **`typed-app`** (2.8 MB) | 1.52 s / 1.40 s | **1.49 s** / **1.37 s** | **−2.11%** / **−2.14%** | 264 / 270 MiB (**−2.41%** `t=1`) | `MATCH` |
| **`downlevel-simple`** (2.1 MB) | 1.18 s / 1.09 s | **1.17 s** / **1.09 s** | **−0.84%** / **−0.79%** | 210 / 205 MiB | `MATCH` |
| **`downlevel-whitespace`** (2.1 MB) | 0.57 s / 0.48 s | **0.57 s** / **0.47 s** | **−0.95%** / **−0.80%** | 170 / 169 MiB | `MATCH` |

## Round 14: `GlobalNamespace` Chunked Arena Pooling & 4096-Entry QName Cache (`−312 MB`), Batched Parallel `ToSource` & Pointer-Free `SourceMapGeneratorV3_Mapping` (`−135 MB`), `JSDocSerializer` / `JSTypeReconserializer` / `ColorPool_Builder` Direct Construction (`−160 MB`), `LinkedDirectedGraph` / `FunctionTypeBuilder` Pooling (`−105 MB`), `LinkedFlowScope` Arena & Scratch `FunctionArgumentInjector` (`−140 MB`), and Pass Arena Eviction (2026-10-21)

Profiling Round 13 (`fc40a451`) on `typed-mega` (`4.17 GB` `alloc_space`) and `bundle-mega` (`4.95 GB` `alloc_space`) identified five disjoint optimization targets:

1. **`GlobalNamespace` Chunked Arena Pooling, 4096-Entry Qualified-Name Cache, and `SetShouldTraverseScriptPredicate` Private Arena (`Subagent 14A`, `global_namespace.go`, `process_defines.go`, `inline_and_collapse_properties.go`):**
   - Extracted `GlobalNamespace` allocation state into a pooled `globalNamespaceArena` using chunked fixed-size array pointers (`nameChunks []*[256]GlobalNamespace_Name`, `refChunks []*[256]GlobalNamespace_Ref`, `propsChunks []*[256]globalNamespace_propsSlot`, `refSliceChunks []*[1024]*GlobalNamespace_Ref`) plus reusable `nameMap`, `globalNames`, and `spreadSiblingCache`.
   - Added a 4096-entry direct-mapped `concatCache` (`concatDotCached`, `concatDot3Cached`, `concatDot4Cached`, `qnameStringCached`) and `globalNamespace_hasPrototype` early bailout across `GlobalNamespace` and `InlineAndCollapseProperties`.
   - Prevented `Compiler.GetTranspilationNamespace()` (`SetShouldTraverseScriptPredicate`) from taking the full-AST `globalNamespaceArena` from `ProcessDefines`, allowing `ProcessDefines` (`CHECK` and `OPTIMIZE`), `Es6RewriteClass`, and `InlineAndCollapseProperties` to share a single `globalNamespaceArena` (`−312 MB` on `typed-mega`) and evicting it to `sync.Pool` immediately after `InlineAndCollapseProperties` completes.
2. **`CodePrinter` Pool Reuse With `gatherSourceMap == true`, 64-Script Batched Parallel `ToSource`, and 32-Byte Pointer-Free `SourceMapGeneratorV3_Mapping` (`Subagent 14B`, `code_printer.go`, `compiler.go`, `internal/sourcemap/source_map_generator_v3.go`):**
   - Detached source mappings into `CodePrinter_SourceAndMappings` (`takeMappingsFrom`, `iterMappings`, `Release`) and retained `m.codeASCII` scratch buffers across pool reuses so `pooledCompact` and `pooledPretty` are unconditionally recycled even when `gatherSourceMap == true`.
   - Changed `Compiler.ToSource` to print scripts in parallel batches of 64 scripts and immediately release each batch's `*[512]SourceMap_Mapping` slabs (`−113.8 MB` of `codePrinter_mappingSlabPool` allocations on `typed-mega`), and evicted cached pass arenas (`evictPassArenaCaches`) before `ToSource`.
   - Replaced `originalName *string` on `SourceMapGeneratorV3_Mapping` with `originalNameIdx int32` into `g.addedSymbolNames`, shrinking `SourceMapGeneratorV3_Mapping` from 40 B to 32 B with zero pointers (`noscan`, `−22.1 MB` on `typed-mega`) and pooling `*[16384]SourceMapGeneratorV3_Mapping` slabs.
3. **Direct `JSDocSerializer_ConvertJSDocInfoForOptimizations`, Single-Map `JSTypeReconserializer` Slabs, and Inline `ColorPool_Builder` Entries (`Subagent 14C`, `serialization_js_doc_serializer.go`, `serialization_serialize_types_to_pointers.go`, `serialization_js_type_reconserializer.go`, `serialization_color_pool.go`, `serialization_convert_types_to_colors.go`):**
   - Rewrote `JSDocSerializer_ConvertJSDocInfoForOptimizations` to construct the optimization `*rhino.JSDocInfo` directly via a stack-allocated `rhino.JSDocInfo_Builder` (returning `nil` with 0 allocations when no optimization properties exist) and replaced `TreeSet` allocations for `modifies` and `suppressions` with alphabetical boolean flags (`−41.0 MB` on `typed-mega`).
   - Replaced `JSTypeReconserializer.typeToRecordCache` with `map[rhino.JSType]*JSTypeReconserializer_SeenTypeRecord`, cached forwarded/axiomatic types directly so `SerializeTypesToPointers` does not populate a duplicate `typePointersByJstype` map during `GatherTypesOnAst`, and added contiguous slab allocations for `TypeProto`, `ObjectTypeProto`, `UnionTypeProto`, `SubtypingEdge`, and `SeenTypeRecord`.
   - Replaced `ColorPool_Builder.idToProto` (`*jutil.HashBasedTable`, ~30K inner `LinkedHashMap`s) with a single `LinkedHashMap[colors.ColorId, colorPoolIdEntry]` storing the first `(shard, proto)` inline, shared empty sets in `reconcileObjectProtos`, and eliminated `[]*colors.ColorId` pointer slices in `colorPool_createTrimmedOffsetToId`.
4. **`LinkedDirectedGraph` Multi-Chunk Slab Retention, Root CFG Map Pooling, and `FunctionTypeBuilder` Pooling (`Subagent 14D`, `graph_linked_directed_graph.go`, `control_flow_analysis.go`, `function_type_builder.go`, `typed_scope_creator.go`):**
   - Retained allocated node and edge slab chunks (`nodeChunks`, `edgeChunks`) and `savedNodes` across `ResetLinkedDirectedGraph`, pooled root CFG `nodes` and `astPosition` maps in `ControlFlowAnalysis`, and replaced `c.finallyMap.Clear()` with `cfaEmptyFinallyMap`.
   - Pooled `FunctionTypeBuilder` via `functionTypeBuilderPool` with embedded `defaultContents FunctionTypeBuilder_AstFunctionContents`, avoided duplicate `LinkedHashSet` / `ToSlice()` allocations in `InferParameterTypesWithParamsParent` and `MaybeUseNativeClassTemplateNames`, and added zero-allocation `StaticTypedScope_LookupQualifiedNameString` lookups in `TypedScopeCreator`.
5. **Chunked `linkedFlowScopeArena`, Batched Unbound-Var Inference, `TypeInference` QName/Template Caches, and Scratch `FunctionArgumentInjector` (`Subagent 14E`, `linked_flow_scope.go`, `type_inference.go`, `function_injector.go`, `function_argument_injector.go`):**
   - Added a reusable chunked `linkedFlowScopeArena` on `TypeInference` (`[256]LinkedFlowScope`, `[256]LinkedFlowScope_OverlayScope`, `[256]LinkedFlowScope_OverlaySlot` chunks), batched `InferDeclarativelyUnboundVarsWithoutTypes` onto a single `OverlayScope`, added a 512-entry `qnameCache`, and reused `matchInvocationTemplateTypes` across calls.
   - Added `GetFunctionCallParameterMapScratch` (using preallocated sentinel `undefined`/`this` nodes and a scratch `LinkedHashMap` + pair slab), `HasCallArgumentsNeedingTemps`, and a reusable `ExpressionDecomposer` on `FunctionInjector` / `FunctionArgumentInjector` (`−68.6 MB` on `bundle-mega`).

Total `alloc_space` dropped another **−669.5 MB (−15.7%)** on `typed-mega` (`4.17 GB → 3.51 GB`, **−69.7% vs Round 5's `11.61 GB`**) and **−179.4 MB (−3.5%)** on `bundle-mega` (`4.95 GB → 4.78 GB`, **−72.2% vs Round 5's `17.16 GB`**), while end-of-compilation `inuse_space` dropped by **−161.2 MB** on `typed-mega` (`1,161.8 MB → 1,000.6 MB`) and **−312.8 MB (−22.1%)** on `bundle-mega` (`1,413.9 MB → 1,101.2 MB`).

### Final measured results — Round 14 (`t=1` and `t=4`, 1 warmup + 5 interleaved runs, PGO enabled)

| Benchmark | Round 13 (`fc40a451` `t=1` / `t=4`) | **Round 14 (`t=1` / `t=4`)** | Round 14 vs Round 13 (`t=1` / `t=4`) | Peak RSS (`t=1` / `t=4`) | Parity |
|---|---:|---:|---:|---:|:---:|
| **`typed-mega`** (44.6 MB, 4,385 files) | 19.84 s / 18.68 s | **19.77 s** / **18.59 s** | **−0.33%** / **−0.53%** | 1,862 / 1,838 MiB (**−0.17%** `t=4`, **−15.7% `alloc_space`**) | `MATCH` |
| **`bundle-mega`** (36.1 MB, 26 UMD libs) | 25.05 s / 24.24 s | **25.04 s** / **24.26 s** | **−0.05%** / **+0.08%** | 2,340 / 2,494 MiB (**−1.65%** `t=1`, **−22.1% `inuse_space`**) | `MATCH` |
| **`typed-full`** (28.5 MB, 1,818 files) | 11.43 s / 10.80 s | **11.36 s** / **10.74 s** | **−0.58%** / **−0.52%** | 1,089 / 1,115 MiB | `MATCH` |
| **`bundle-advanced-ts`** (12.6 MB) | 7.08 s / 6.95 s | **7.06 s** / **6.94 s** | **−0.31%** / **−0.28%** | 635 / 631 MiB (**−0.88%** `t=1`) | `MATCH` |
| **`bundle`** (12.6 MB) | 4.64 s / 4.52 s | **4.59 s** / **4.48 s** | **−1.23%** / **−0.94%** | 534 / 552 MiB (**−2.60%** `t=1`) | `MATCH` |
| **`shaka-blockly`** (8.6 MB) | 4.75 s / 4.55 s | **4.72 s** / **4.49 s** | **−0.74%** / **−1.39%** | 511 / 506 MiB (**−1.07%** `t=4`) | `MATCH` |
| **`bundle-advanced`** (5.4 MB) | 3.62 s / 3.44 s | **3.60 s** / **3.43 s** | **−0.43%** / **−0.05%** | 375 / 384 MiB (**−2.55%** `t=4`) | `MATCH` |
| **`checks-only`** (11.3 MB) | 2.57 s / 2.20 s | **2.59 s** / **2.21 s** | +0.73% / +0.52% | 580 / 561 MiB (**−8.24%** `t=4`) | `MATCH` |
| **`typed-app`** (2.8 MB) | 1.49 s / 1.38 s | **1.48 s** / **1.37 s** | **−0.87%** / **−0.43%** | 277 / 262 MiB | `MATCH` |
| **`downlevel-simple`** (2.1 MB) | 1.17 s / 1.08 s | **1.16 s** / **1.08 s** | **−0.48%** / **−0.09%** | 214 / 216 MiB | `MATCH` |
| **`downlevel-whitespace`** (2.1 MB) | 0.56 s / 0.47 s | **0.56 s** / **0.47 s** | +0.05% / +0.29% | 169 / 169 MiB | `MATCH` |

## Round 15: Zero-Clone Inlined Call Arguments (`−307 MB`), `SourceMapGeneratorV3` Symbol Deduplication & Adaptive `ToSource` Batching (`−285 MB`), Zero-Copy File & JSON Decoding (`−187 MB`), Stack `TemplateTypeReplacer` & `JSTypeRegistry` Slabs (`−95 MB`), and `OptimizeCalls` / `PureFunctionIdentifier` / `ColorPool` Reuse (`−145 MB`) (2026-10-21)

Profiling Round 14 (`6daef8ee`) on `typed-mega` (`3.51 GB` `alloc_space`) and `bundle-mega` (`4.81 GB` `alloc_space`) identified five disjoint optimization targets:

1. **Zero-Copy UTF-8 File Decoding, `JsonReader` Single-Allocation ASCII Fast Path, `SourceMapConsumerV3` Slice Trim, `PrintStream`, and `LightweightMessageFormatter` (`Subagent 15A`, `internal/jutil/files_read.go`, `internal/jutil/input_stream_reader.go`, `internal/jutil/gson_json_reader.go`, `internal/jutil/print_stream.go`, `internal/sourcemap/source_map_consumer_v3.go`, `lightweight_message_formatter.go`, `source_file.go`):**
   - Returned `unsafe.String(unsafe.SliceData(b), len(b))` directly when `utf8.Valid(b)` holds on the freshly allocated `[]byte` from `os.ReadFile` / `io.ReadAll` in `Files_ReadString` and `InputStreamReader.decodeAll` (`−80.25 MB` on `bundle-mega`).
   - Added a single-allocation ASCII fast path in `JsonReader.nextQuotedValue` when reading multi-MB unescaped strings (such as inline sourcemap `"mappings"`) from a `*StringReader`, plus `2x` capacity growth on `builder` (`−107.11 MB` on `bundle-mega`).
   - Replaced `sourceMapConsumerV3_copyOfInt32s` / `sourceMapConsumerV3_copyOfInts` in `SourceMapConsumerV3_MappingBuilder.Trim` with in-place reslicing (`−23.77 MB` on `bundle-mega`), avoided `[]byte` allocations in `PrintStream.write` / `Println` (`−44.55 MB` on `bundle-mega`), and eliminated intermediate builders, UTF-16 slice conversions, and line-split slices in `LightweightMessageFormatter` and `SourceFile.FindLineOffsets` (`−65+ MB` on `bundle-mega`).
2. **Zero-Clone Call Argument Transfer in `FunctionToBlockMutator`, Static `staticNumberExts` Table, and Slab-Backed `CloneTreeWithTypeExprs` (`Subagent 15B`, `internal/rhino/node.go`, `function_to_block_mutator.go`, `function_injector.go`, `inline_functions.go`):**
   - Added `staticNumberExts [1024]node_ext` in `internal/rhino/node.go` for non-negative integer number nodes in `[0, 1024)` (`−23.50 MB`) and slab-backed `CloneTreeWithTypeExprs` (`nodeCloneSlab`, capped at 512 nodes per chunk and sized to exact subtree node count with zero tail waste).
   - Added `takeOrCloneCallArg` and `rebindMovedArgFunctions` in `FunctionToBlockMutator.AliasAndInlineArguments` (plus `movedArgs` tracking in `InlineFunctions.RemoveInlinedFunctions`) to move call arguments directly out of `callNode` instead of cloning 100,000-node UMD factory arguments (`−307.02 MB` on `bundle-mega`).
3. **`SourceMapGeneratorV3` Symbol & Source Deduplication and Adaptive `Compiler.ToSource` Batch Size (`Subagent 15C`, `internal/sourcemap/source_map_generator_v3.go`, `compiler.go`):**
   - Deduplicated `g.addedSymbolNames` and `g.addedSourceFiles` in `SourceMapGeneratorV3.AddMappingPositions` (`−165.45 MB` on `bundle-mega`, `−109.99 MB` on `typed-mega`) and added dense `symbolIdxToNameId` and `sourceIdxToFileId` O(1) lookup slices in `WriteEntry`.
   - Made `batchSize` in `Compiler.ToSource` adaptive (`max(numThreads, 4)` when `len(scripts) <= 64`) so `bundle-mega` (26 huge scripts) prints 4 scripts per batch and recycles `codePrinter_mappingSlabPool` slabs across batches (`−120 MB` on `bundle-mega`).
4. **Stack-Allocated `TemplateTypeReplacer` and `JSTypeRegistry` Slabs for `TemplateTypeMap` & `TemplatizedType` (`Subagent 15D`, `internal/rhino/jstype_template_type_replacer.go`, `internal/rhino/jstype_template_type_map.go`, `internal/rhino/jstype_templatized_type.go`, `internal/rhino/jstype_js_type.go`, `internal/rhino/jstype_js_type_registry.go`, `internal/rhino/jstype_function_type.go`):**
   - Added value constructors `MakeTemplateTypeReplacerForPartialReplacement` and `TemplateTypeReplacer_MakeForInference` plus direct `VisitJSType`, and removed the eager per-instance `replacer *TemplateTypeReplacer` field from `TemplatizedType` (`−38.01 MB` on `typed-mega`).
   - Added `templateTypeMapSlab`, `templatizedTypeSlab`, and `jsTypeSliceSlab` on `JSTypeRegistry`, removed the 64-byte `valuesBuf [4]JSType` from `TemplateTypeMap`, reused unchanged `*FunctionType_Parameter` instances in `CaseFunctionTypeUnguarded`, and used a stack `FunctionType_Builder` (`−55+ MB` on `typed-mega`).
5. **`OptimizeCalls_ReferenceMap` Slice Retention, Direct `PureFunctionIdentifier` Iteration, and Single-Shard `ColorPool_Builder` Fast Path (`Subagent 15E`, `optimize_calls.go`, `pure_function_identifier.go`, `serialization_js_type_reconserializer.go`, `serialization_color_pool.go`):**
   - Increased `optimizeCalls_refListSlot.buf` to `[8]*rhino.Node` and retained grown `ArrayList` backing slices (`cap <= 4096`) across `releaseOptimizeCallsRefMap` pool reuse (`−60+ MB` on `typed-mega`).
   - Eliminated `JavaHashMap` and `ArrayList` copies in `PureFunctionIdentifier.PopulateDatastructuresForAnalysisTraversal`, replaced `summariesForAllNamesOfFunctionByNode` (`ListMultimap`) with `map[*rhino.Node][]*PureFunctionIdentifier_AmbiguatedFunctionSummary` using `summary.selfSlice[:1:1]`, and added a single-shard direct-indexed fast path and set slabs in `ColorPool_Builder.Build` (`−75+ MB` on `typed-mega`).

Total `alloc_space` dropped by **−827.3 MB (−16.8%)** on `bundle-mega` (`4.81 GB → 4.00 GB`, **−76.7% vs Round 5's `17.16 GB`**) and **−138.8 MB (−3.9%)** on `typed-mega` (`3.51 GB → 3.38 GB`, **−70.9% vs Round 5's `11.61 GB`**), cutting peak RSS by **−642.2 MiB (−25.86%)** on `bundle-mega (t=4)` (`2,483.8 MiB → 1,841.6 MiB`) and **−396.5 MiB (−17.55%)** on `bundle-mega (t=1)` (`2,258.9 MiB → 1,862.4 MiB`), with RSS reductions across all 11 benchmarks!

### Final measured results — Round 15 (`t=1` and `t=4`, 1 warmup + 5 interleaved runs, PGO enabled)

| Benchmark | Round 14 (`6daef8ee` `t=1` / `t=4`) | **Round 15 (`t=1` / `t=4`)** | Round 15 Peak RSS (`t=1` / `t=4`) | Peak RSS vs Round 14 (`t=1` / `t=4`) | Parity |
|---|---:|---:|---:|---:|:---:|
| **`bundle-mega`** (36.1 MB, 26 UMD libs) | 25.02 s / 24.22 s | **25.11 s** / **24.32 s** | **1,862** / **1,842 MiB** | **−17.55%** / **−25.86%** (**−827.3 MB `alloc_space`**) | `MATCH` |
| **`typed-mega`** (44.6 MB, 4,385 files) | 19.75 s / 18.59 s | **19.80 s** / **18.59 s** | **1,791** / **1,786 MiB** | **−8.02%** / **−2.36%** (**−138.8 MB `alloc_space`**) | `MATCH` |
| **`typed-full`** (28.5 MB, 1,818 files) | 11.34 s / 10.76 s | **11.41 s** / **10.73 s** | **1,082** / **1,078 MiB** | **−1.58%** / **−2.53%** | `MATCH` |
| **`bundle-advanced-ts`** (12.6 MB) | 7.06 s / 6.94 s | **7.07 s** / **6.96 s** | **567** / **579 MiB** | **−11.73%** / **−7.72%** | `MATCH` |
| **`bundle`** (12.6 MB) | 4.59 s / 4.47 s | **4.65 s** / **4.55 s** | **538** / **523 MiB** | **−0.05%** / **−6.59%** | `MATCH` |
| **`shaka-blockly`** (8.6 MB) | 4.71 s / 4.50 s | **4.71 s** / **4.49 s** | **487** / **500 MiB** | **−3.85%** / **−6.84%** | `MATCH` |
| **`bundle-advanced`** (5.4 MB) | 3.60 s / 3.42 s | **3.61 s** / **3.46 s** | **334** / **331 MiB** | **−9.12%** / **−11.21%** | `MATCH` |
| **`checks-only`** (11.3 MB) | 2.60 s / 2.20 s | **2.62 s** / **2.22 s** | **573** / **540 MiB** | **−0.23%** / **−7.91%** | `MATCH` |
| **`typed-app`** (2.8 MB) | 1.48 s / 1.36 s | **1.48 s** / **1.37 s** | **266** / **271 MiB** | **−1.86%** / +1.29% | `MATCH` |
| **`downlevel-simple`** (2.1 MB) | 1.17 s / 1.08 s | **1.18 s** / **1.09 s** | **207** / **205 MiB** | **−5.03%** / **−5.62%** | `MATCH` |
| **`downlevel-whitespace`** (2.1 MB) | 0.56 s / 0.48 s | **0.57 s** / **0.48 s** | **165** / **169 MiB** | **−3.59%** / **−1.01%** | `MATCH` |

## Round 16: `HamtPMapArena` & `UnionType_Builder` Stack (`−187 MB`), `ReferenceCollector` / `RemoveUnusedCode` Direct `Var.nameHash` Indexing & Largest-First `scopeArena` (`−110 MB`), Byte-Backed `Compiler_CodeBuilder` & `SourceMapGeneratorV3` Ordered Slices (`−105 MB`), `ParseArena` Chunk Retention & `LineNumberScanner` Fast Path (`−55 MB`), and `GlobalNamespace` 64K Single-String Concat Cache (`−60 MB`) (2026-10-21)

Profiling Round 15 (`80199f92`) on `typed-mega` (`3,459.97 MB` `alloc_space`) and `bundle-mega` (`4,094.09 MB` `alloc_space`) identified five disjoint optimization targets across all 11 benchmarks:

1. **`ParseArena` Chunk Retention, Doubling UTF-16 Buffer Growth, and `LineNumberScanner` Same-Line Fast Path (`Subagent 16A`, `internal/parser/parse_arena.go`, `internal/parser/line_number_scanner.go`):**
   - Increased `arenaSlabMaxChunk` from `4096` to `8192` and `arenaSlabMaxChunks` from `64` to `256` so 2–3 MB scripts with 400K–800K `SourcePosition`s, `SourceRange`s, and `Token`s in `PrebuildAst` no longer truncate and reallocate chunks `64..195` after every script (`−53.78 MB` on `bundle-mega`), and grew `ParseArena.decodeUTF16`'s backing buffer by doubling capacity.
   - Guarded `jutil.CheckArgumentMsg2` in `LineNumberScanner.getSourcePosition` to avoid boxing `offset` and `l.lastLineStart` to `any` on every token, and added a same-line inline fast path in `getTokenSourceRange`.
2. **Direct `Var.nameHash` Pass Indexing in `ReferenceCollector`, `CrossChunkReferenceCollector`, & `RemoveUnusedCode`, and Largest-First `cachedScopeArenas` Selection (`Subagent 16B`, `reference_collector.go`, `cross_chunk_reference_collector.go`, `remove_unused_code.go`, `scope.go`, `node_util.go`):**
   - Replaced hash table probing (`c.table []rcTableEntry`, `growTable()`) in `ReferenceCollector` and `CrossChunkReferenceCollector` with O(1) direct 1-based indexing via `v.nameHash` (preserving `sizeof(Var) == 48 B`), explicit `2x` `growVarsCols()` growth, and arena pooling in `NodeUtil_CollectExternVariableNames` (`−61.25 MB` and `−290 ms` CPU on `bundle-mega`).
   - Added a `var_.nameHash` O(1) fast path at the top of `RemoveUnusedCode.GetVarInfo` and `TraverseVar`.
   - Updated `acquireScopeArena()` and `scopeArena.release()` to select and retain the largest pre-warmed `scopeArena` across the 4-slot `cachedScopeArenas` cache (`−17.94 MB` on `bundle-mega`, `−14.39 MB` on `typed-mega`).
3. **`HamtPMapArena[K, V]` in `linkedFlowScopeArena` & `MaybeReachingVariableUse`, and `UnionType_Builder` Stack (`Subagent 16C`, `internal/rhino/hamt_p_map.go`, `internal/rhino/jstype_union_type.go`, `internal/rhino/jstype_js_type_registry.go`, `internal/rhino/jstype_template_type_replacer.go`, `linked_flow_scope.go`, `type_inference.go`, `maybe_reaching_variable_use.go`):**
   - Added `HamtPMapArena[K, V]` (`PlusMapWithArena`, `MinusMapWithArena`, `ReconcileMapWithArena`) and unboxed `jutil.ObjectsEqualsOf` value comparisons in `hamt_p_map.go`, wiring `scopeMapArena` and `slotMapArena` into `linkedFlowScopeArena` (`TypeInference`) and concrete `*HamtPMap` + `varMapArena`/`nodeMapArena` into `MaybeReachingVariableUse` (`−158.50 MB` on `bundle-mega`, `−55.50 MB` on `typed-mega`).
   - Added `unionBuilderStack [8]UnionType_Builder` (with `inlineAlternates [32]JSType`) and `unionTypeSlab` on `JSTypeRegistry`, eliminating `UnionType_Builder` heap escapes caused by the mutual recursion SCC between `CreateUnionTypeOfList`, `AddAlternate`, `GetAlternates`, and `IsSubtypeWithoutStructuralTyping` (`−53.01 MB` on `bundle-mega`, `−15.50 MB` on `typed-mega`).
4. **Byte-Backed `Compiler_CodeBuilder` and Ordered-Slice `SourceMapGeneratorV3` (`Subagent 16D`, `compiler.go`, `internal/sourcemap/source_map_generator_v3.go`):**
   - Replaced `sb *jutil.StringBuilder` (`[]uint16` UTF-16 backing slice) in `Compiler_CodeBuilder` with a UTF-8 `buf []byte`, single-pass `utf16Len`/`lineCount`/`colCount` tracking, and zero-copy `String()` (`−62.33 MB` on `bundle-mega`, `−89.25 MB` on `typed-mega`).
   - Replaced `sourceFileMap` and `originalNameMap` (`*jutil.LinkedHashMap[string, int]`) in `SourceMapGeneratorV3` with flat `orderedSourceFiles` and `orderedSymbolNames` slices populated directly in `WriteEntry` (`−22.90 MB` on `typed-mega`, `−15.88 MB` on `bundle-mega`).
5. **`GlobalNamespace` 64K Single-String Concat Cache and Zero-Allocation Qualified-Name Checks (`Subagent 16E`, `global_namespace.go`, `coding_convention.go`, `js_message_visitor.go`, `check_side_effects.go`, `replace_id_generators.go`, `process_closure_primitives.go`, `closure_check_module.go`, `analyze_prototype_properties.go`):**
   - Shrunk `GlobalNamespace` `concatCache` entries from `struct{a, b, res string}` (48 B, 3 pointers) to `string` (16 B, 1 pointer) verified via `gnConcatMatch(res, a, b)`, expanded `gnConcatCacheSize` from `4,096` to `65,536`, and added recursive caching for 5+ segment `GETPROP` chains (`−30.15 MB` on `typed-mega`).
   - Eliminated `GetQualifiedName()` / `NodeUtil_GetBestLValueName()` allocations in `CodingConvention_AssertionFunctionLookup.LookupByCallee`, `JsMessageVisitorBase.isDeclareIcuTemplateCallee`, `CheckSideEffects.Visit`, `ReplaceIdGenerators_ReplaceGenerators.Visit`, `ProcessClosurePrimitives.Visit`, `ClosureCheckModule.VisitWithModule`, and `AnalyzePrototypeProperties` (`−30+ MB` across `typed-mega` and `bundle-mega`).

Total `alloc_space` dropped by **−417.52 MB (−10.20%)** on `bundle-mega` (`4,094.09 MB → 3,676.57 MB`, **−78.6% vs Round 5's `17.16 GB`**) and **−268.32 MB (−7.75%)** on `typed-mega` (`3,459.97 MB → 3,191.65 MB`, **−72.5% vs Round 5's `11.61 GB`**), while end-of-compilation `inuse_space` dropped to **`985.48 MB` (−57.70 MB, breaking below 1 GB)** on `typed-mega` and **`1,194.87 MB` (−23.36 MB)** on `bundle-mega`, and wall time improved across **all 22/22 (`t=1` and `t=4`) benchmark configurations** (`−4.51%` / `−1.10 s` on `bundle-mega`, `−3.48%` / `−0.69 s` on `typed-mega`)!

### Final measured results — Round 16 (`t=1` and `t=4`, 1 warmup + 5 interleaved runs, PGO enabled)

| Benchmark | Round 15 (`80199f92` `t=1` / `t=4`) | **Round 16 (`t=1` / `t=4`)** | Round 16 vs Round 15 (`t=1` / `t=4`) | Peak RSS (`t=1` / `t=4`) | Parity |
|---|---:|---:|---:|---:|:---:|
| **`bundle-mega`** (36.1 MB, 26 UMD libs) | 25.07 s / 24.32 s | **23.96 s** / **23.22 s** | **−4.44%** / **−4.51%** | **1,784** / **1,771 MiB** (**−2.97%** / **−3.48%**, **−417.5 MB `alloc_space`**) | `MATCH` |
| **`typed-mega`** (44.6 MB, 4,385 files) | 19.83 s / 18.58 s | **19.14 s** / **18.02 s** | **−3.48%** / **−3.06%** | 1,899 / 1,861 MiB (**−268.3 MB `alloc_space`**, **985.5 MB `inuse_space`**) | `MATCH` |
| **`typed-full`** (28.5 MB, 1,818 files) | 11.37 s / 10.71 s | **11.05 s** / **10.46 s** | **−2.77%** / **−2.40%** | 1,113 / 1,156 MiB | `MATCH` |
| **`bundle-advanced-ts`** (12.6 MB) | 7.07 s / 6.96 s | **6.82 s** / **6.70 s** | **−3.58%** / **−3.81%** | **554** / 616 MiB (**−1.58%** `t=1`) | `MATCH` |
| **`bundle`** (12.6 MB) | 4.66 s / 4.53 s | **4.48 s** / **4.37 s** | **−3.79%** / **−3.66%** | 545 / **593 MiB** (**−1.88%** `t=4`) | `MATCH` |
| **`shaka-blockly`** (8.6 MB) | 4.70 s / 4.48 s | **4.57 s** / **4.35 s** | **−2.71%** / **−3.06%** | **516** / 519 MiB | `MATCH` |
| **`bundle-advanced`** (5.4 MB) | 3.61 s / 3.46 s | **3.49 s** / **3.33 s** | **−3.33%** / **−3.72%** | **317** / **306 MiB** (**−4.21%** / **−6.77%**) | `MATCH` |
| **`checks-only`** (11.3 MB) | 2.63 s / 2.21 s | **2.55 s** / **2.14 s** | **−2.97%** / **−2.82%** | 579 / 557 MiB | `MATCH` |
| **`typed-app`** (2.8 MB) | 1.48 s / 1.36 s | **1.44 s** / **1.32 s** | **−2.55%** / **−3.01%** | 262 / **259 MiB** (**−4.86%** `t=4`) | `MATCH` |
| **`downlevel-simple`** (2.1 MB) | 1.18 s / 1.09 s | **1.15 s** / **1.06 s** | **−2.66%** / **−2.59%** | **199** / **197 MiB** (**−7.34%** / **−3.19%**) | `MATCH` |
| **`downlevel-whitespace`** (2.1 MB) | 0.57 s / 0.48 s | **0.56 s** / **0.47 s** | **−2.07%** / **−1.75%** | **157** / **152 MiB** (**−6.47%** / **−9.78%**) | `MATCH` |

## Round 17: `OptimizeCalls_ReferenceMap` Arena Pooling (`−78 MB`), `ControlFlowAnalysis` Contiguous `DiGraph` Slices (`−55 MB`), `DisambiguateProperties` / `AmbiguateProperties` / `RenameProperties` Pooling & Flat Arrays (`−35 MB`), `JSTypeRegistry.namedToType` Pre-Sizing (`−17 MB`), and Zero-Copy `BufferedWriter.WriteString` (`−11 MB`) (2026-10-21)

Profiling Round 16 (`080e9677`) on `typed-mega` (`3,191.65 MB` `alloc_space`) and `bundle-mega` (`3,676.57 MB` `alloc_space`) targeted Stage 2 peak-heap allocations and property renaming passes:

1. **`OptimizeCalls_ReferenceMap` Arena Pooling (`internal/jscomp/optimize_calls.go`):**
   - Retained `nameSlots`, `propSlots`, `refSlots`, and `nodeSlices` backing slices (capped at 64 chunks) in `releaseOptimizeCallsRefMap` so all 4–5 iterations of `OptimizeCalls` in the Stage 2 optimization loop reuse pre-warmed slab chunks (`−78 MB` `alloc_space` on `typed-mega` and `bundle-mega`).
2. **Contiguous Slice Backing for `ControlFlowAnalysis` `DiGraphNode` & `DiGraphEdge` (`internal/jscomp/control_flow_analysis.go`, `internal/jscomp/node_traversal.go`):**
   - Added `cfaInitialNodesCap` / `cfaInitialEdgesCap` and retained `nodeChunk` / `edgeChunk` buffers in `cfaArena.reset` (up to `32,768` nodes and `65,536` edges), plus stack-buffered `[16]*rhino.Node` roots in `NodeTraversal.TraverseRoots` (`−55 MB` `alloc_space`).
3. **`DisambiguateProperties`, `AmbiguateProperties`, and `RenameProperties` Pooling & Contiguous Arrays (`internal/jscomp/disambiguate_ambiguate_properties.go`, `internal/jscomp/rename_properties.go`):**
   - Pooled `disambiguateProperties_Property` slabs and `nodeListPool` across pass invocations (`−20 MB`), replaced per-node `fixedBitSet` slice headers in `AmbiguateProperties` with a flat contiguous `bitsArena []uint64`, and replaced `string` keys in `RenameProperties` `Node` maps with `*RenameProperties_Property` pointer keys and a contiguous `RenameProperties_Property` slab (`−15 MB`).
4. **`JSTypeRegistry.namedToType` Pre-Sizing and Zero-Copy `BufferedWriter.WriteString` / `OutputStreamWriterCharset.WriteString` (`internal/rhino/jstype_js_type_registry.go`, `internal/jutil/buffered_writer.go`, `internal/jutil/output_stream_writer_charset.go`):**
   - Pre-sized `JSTypeRegistry.namedToType` (`NewLinkedHashMapWithCapacity(256)`) to avoid early map rehashing (`−17.02 MB` `inuse_space`, `−8.51 MB` `alloc_space`), and added zero-copy `WriteString(string)` fast paths to `BufferedWriter` and `OutputStreamWriterCharset` (`−11 MB` `alloc_space`).

### Final measured results — Round 17 (`t=4`, interleaved runs, PGO enabled)

| Benchmark | Round 16 (`t=4`) | **Round 17 (`t=4`)** | Wall Time Delta | Round 16 Peak RSS | **Round 17 Peak RSS** | Peak RSS Delta | Parity |
|---|---:|---:|---:|---:|---:|---:|:---:|
| **`typed-mega`** (44.6 MB, 4,385 files) | 18.05 s | **17.98 s** | **−0.39%** | 1,842 MiB | **1,660 MiB** | **−9.85%** (**−181 MiB**) | `MATCH` |
| **`bundle-mega`** (36.1 MB, 26 UMD libs) | 23.24 s | **23.16 s** | **−0.34%** | 1,744 MiB | **1,750 MiB** | +0.31% | `MATCH` |
| **`shaka-blockly`** (8.6 MB) | 4.36 s | **4.35 s** | **−0.26%** | 520 MiB | **453 MiB** | **−12.91%** (**−67 MiB**) | `MATCH` |
| **`checks-only`** (11.3 MB) | 2.16 s | **2.16 s** | +0.07% | 535 MiB | **564 MiB** | +5.49% | `MATCH` |
| **`downlevel-simple`** (2.1 MB) | 1.07 s | **1.06 s** | **−0.72%** | 205 MiB | **190 MiB** | **−7.29%** (**−15 MiB**) | `MATCH` |

## Round 18: `node_sourceFilePropSet` (`−63.5 MB`), `functionTypeSlab` (`−53.5 MB`), Iterative `NodeUtil_CountAstSize` & Loop Guards, Pooled `inferConstsCollector` (`−30 MB`), 1-Slot Loop-Retained `OptimizeCalls_ReferenceMap` & 4x `stringNodesToRename` Handoff (`−115 MB`), and Shared `ReusableCfg` Pool (`−45 MB`) (2026-10-21)

Profiling Round 17 (`b97b6be8`) identified five disjoint targets:

1. **`node_sourceFilePropSet` Shared `IntPropListItem`s & `functionTypeSlab` (`Subagent 18A`, `internal/rhino/node.go`, `internal/rhino/jstype_function_type.go`, `internal/rhino/jstype_js_type_registry.go`):**
   - Embedded the 16 most common `(Node_Prop, int32)` property items directly alongside each file's `Node_Prop_SOURCE_FILE` item in `node_sourceFilePropSet` (`base` + `commons [16]Node_PropListItem` with `shared: true`), eliminating `63.5 MB` of `NewNode_IntPropListItem` allocations on `typed-mega` and `bundle-mega`.
   - Added `functionTypeSlab []FunctionType` (`allocFunctionType()`) on `JSTypeRegistry` (`−53.5 MB` `NewFunctionType`).
2. **Guarded & Iterative `NodeUtil_CountAstSize` and Closure-Free `NodeUtil_CountAstSizeUpToLimit` (`Subagent 18B`, `internal/jscomp/phase_optimizer.go`, `internal/jscomp/node_util.go`):**
   - Guarded `NodeUtil_CountAstSize(root)` in `PhaseOptimizer_Loop.Process` on `trackAstSize` (`p.useSizeHeuristicToStopOptimizationLoop && l.isCodeRemovalLoop`) and `lastIterMadeChanges`, and replaced recursive `NodeUtil_CountAstSize` with a 2-level leaf-unrolled iterative walk.
3. **Pooled `inferConstsCollector` & Chunked `BasicBlock` Slabs (`Subagent 18C`, `internal/jscomp/infer_consts.go`):**
   - Pooled `inferConstsCollector` (`inferConstsCollectorPool`), replaced `blockSlab` with multi-chunk `blockChunks`, and pooled `blockStack` slices across `maybeJumpToHoistedFunction`.
4. **Loop-Retained `OptimizeCalls_ReferenceMap` & 4x-Grown `stringNodesToRename` Handoff (`Subagent 18D`, `internal/jscomp/optimize_calls.go`, `internal/jscomp/disambiguate_ambiguate_properties.go`, `internal/jscomp/rename_properties.go`):**
   - Added 1-slot `cachedOptimizeCallsRefMap` and `cachedOptimizeCallsCallback` (evicted via `evictOptimizeCallsLoopCache()` when `PhaseOptimizer_Loop.Process` exits) so `OptimizeCalls` retains its pre-grown reference map across GC cycles between loop iterations (`−79 MB` `alloc_space`).
   - Added `addPropNodeToRename` (4x slice growth below 256K elements instead of Go's 1.25x `growslice` factor) and a 1-slot handoff `cachedPropNodesToRename` from `AmbiguateProperties` to `RenameProperties` (`−55 MB` `alloc_space`).
5. **Shared `ReusableCfg` Pool Across `TypeInferencePass`, `FlowSensitiveInlineVariables`, and `NodeTraversal_AbstractCfgCallbackBase` (`Subagent 18E`, `internal/jscomp/control_flow_analysis.go`, `internal/jscomp/graph_linked_directed_graph.go`, `internal/jscomp/node_traversal.go`, `internal/jscomp/type_inference_pass.go`, `internal/jscomp/flow_sensitive_inline_variables.go`):**
   - Added `ClearSlabsForPool()` on `LinkedDirectedGraph` and a 4-slot `cachedReusableCfgs` + `reusableCfgPool` (evicted via `evictCachedReusableCfgs()` when `PhaseOptimizer.Process` exits), sharing pre-warmed CFG node/edge chunks across all CFG-consuming passes.

### Final measured results — Round 18 (`t=1` and `t=4`, 1 warmup + 5 interleaved runs, PGO enabled)

| Benchmark | Round 17 (`b97b6be8` `t=1` / `t=4`) | **Round 18 (`t=1` / `t=4`)** | Round 18 vs Round 17 (`t=1` / `t=4`) | Peak RSS (`t=1` / `t=4`) | Parity |
|---|---:|---:|---:|---:|:---:|
| **`bundle-mega`** (36.1 MB, 26 UMD libs) | 24.04 s / 23.24 s | **24.02 s** / **23.11 s** | **−0.10%** / **−0.56%** | 1,742 / 1,743 MiB | `MATCH` |
| **`typed-mega`** (44.6 MB, 4,385 files) | 19.27 s / 18.13 s | **19.28 s** / **18.12 s** | +0.03% / **−0.01%** | 1,728 / **1,714 MiB** (**−2.88%** `t=4`) | `MATCH` |
| **`typed-full`** (28.5 MB, 1,818 files) | 11.08 s / 10.50 s | **11.06 s** / **10.47 s** | **−0.18%** / **−0.35%** | 1,096 / 1,093 MiB | `MATCH` |
| **`bundle-advanced-ts`** (12.6 MB) | 6.86 s / 6.73 s | **6.82 s** / **6.70 s** | **−0.56%** / **−0.45%** | 554 / **547 MiB** (**−27.27%** / **−205 MiB** `t=4`) | `MATCH` |
| **`bundle`** (12.6 MB) | 4.47 s / 4.36 s | **4.46 s** / **4.36 s** | **−0.23%** / +0.21% | **530** / 736 MiB (**−3.70%** `t=1`) | `MATCH` |
| **`shaka-blockly`** (8.6 MB) | 4.60 s / 4.37 s | **4.58 s** / **4.34 s** | **−0.50%** / **−0.47%** | 490 / 474 MiB | `MATCH` |
| **`bundle-advanced`** (5.4 MB) | 3.50 s / 3.36 s | **3.49 s** / **3.34 s** | **−0.32%** / **−0.56%** | **323** / 309 MiB | `MATCH` |
| **`checks-only`** (11.3 MB) | 2.57 s / 2.16 s | **2.56 s** / **2.16 s** | **−0.56%** / **−0.23%** | **522** / **557 MiB** (**−5.50%** / **−2.63%**) | `MATCH` |
| **`typed-app`** (2.8 MB) | 1.44 s / 1.32 s | **1.43 s** / **1.32 s** | **−0.54%** / **−0.33%** | 266 / 257 MiB | `MATCH` |
| **`downlevel-simple`** (2.1 MB) | 1.14 s / 1.05 s | **1.14 s** / **1.06 s** | +0.09% / +0.28% | 196 / 203 MiB | `MATCH` |
| **`downlevel-whitespace`** (2.1 MB) | 0.55 s / 0.46 s | **0.55 s** / **0.46 s** | **−0.47%** / +0.24% | 156 / 157 MiB | `MATCH` |

## Round 19 (Phase 1 / Bucket A): Inline `rhino.Node` `boolFlags uint16` Bitmask Mirror, `NodeTraversal` Leaf-Child Unrolling & `noShouldTraverse`/`noVisit`, Empty-Block-Bypassing `t.LookupVar` & `Scope.lookupParent`, `scopeNameEq`, `TypedScope` Arena Pre-Sizing & `varBloom`, 2-Way Set-Associative `GlobalNamespace` Cache, and O(1) Ring-Buffer `FixedPointGraphTraversal` (2026-10-21)

Phase 1 (Bucket A — 100% byte-identical shared architectural upgrades) delivered five structural improvements across the AST, scope resolution, traversal, and fixed-point graph engines:

1. **Inline `boolFlags uint16` Mirror on `rhino.Node` & Slab-Backed `RhinoStringPool` (`Subagent 19A`, `internal/rhino/node.go`, `internal/rhino/rhino_string_pool.go`):**
   - Packed `boolFlags uint16` into the 2 bytes of previously unused struct padding at offset 46..48 of `rhino.Node` (keeping `sizeof(rhino.Node) == 80`), mirroring the 14 hottest boolean properties (`IS_PARENTHESIZED`, `INCRDECR`, `QUOTED`, `SYNTHETIC`, `ADDED_BLOCK`, `IS_CONSTANT_NAME`, `IS_NAMESPACE`, `FREE_CALL`, `STATIC_MEMBER`, `ARROW_FN`, `ASYNC_FN`, `GENERATOR_FN`, `IS_SHORTHAND_PROPERTY`), 2-bit `CONSTANT_VAR_FLAGS`, and `nodeFlagNoBlockDecls` (`bit 15`).
   - Turned hot property reads (`GetBooleanProp`, `GetIntProp`, `IsSyntheticBlock`, `IsArrowFunction`, `IsConstantVar`, `SetConstantVarFlag`, etc.) into single-instruction bitmask tests and short-circuited no-op writes before walking or allocating `propListHead`.
   - Added per-shard `cellSlab []string` and `byteSlab []byte` to `rhinoStringPoolShard` and bumped `rhinoStringPool_MAX_SHARD_SIZE` from `16,384` to `65,536`.
2. **`NodeTraversal` Leaf-Child Unrolling, `noShouldTraverse`/`noVisit` Dispatch Bypass, `t.LookupVar(name)` Empty-Block Scope Bypass, `Scope.lookupParent` + `scopeNameEq`, and `TypedScope` Arena Pre-Sizing & `varBloom` (`Subagent 19B`, `internal/jscomp/node_traversal.go`, `internal/jscomp/es6_check_module.go`, `internal/jscomp/syntactic_scope_creator.go`, `internal/jscomp/scope.go`, `internal/jscomp/typed_scope.go`):**
   - Removed `NodeTraversal_AbstractPostOrderCallbackBase` embedding from `Es6CheckModule`, enabling marker interfaces `nodeTraversal_noShouldTraverse` and `nodeTraversal_noVisit` to skip no-op `ShouldTraverse` and `Visit` interface dispatches across `TraverseBranch` and `TraverseChildren`.
   - Unrolled leaf children (`child.GetFirstChild() == nil && nodeTraversalTokenKind[uint8(child.GetToken())] == 0`) inside `TraverseBranch` and `TraverseChildren`, eliminating ~60% of recursive `TraverseBranch` calls.
   - Added `scopeNameEq(a, b string)` (`len(a) == len(b) && (unsafe.StringData(a) == unsafe.StringData(b) || (len(a) > 0 && a[0] == b[0] && a == b))`), eliminating out-of-line `runtime.memequal` calls on same-length variable name mismatches.
   - Added `lookupParent *Scope` and `hasImplicitVars bool` on `Scope` plus `t.LookupVar(name)` (`lookupScope *Scope` on `scopeStackEntry`) and `syntacticScopeHasBlockDecls` (`nodeFlagNoBlockDecls`) so declaration-free inner `BLOCK` and `FOR*` scopes are never instantiated or walked during variable lookups.
   - Eliminated the duplicate 256-node `countScopeVars` AST walk in `largeScopeVarSlotCount` when `sc.arena != nil`.
   - Pre-sized `TypedScope.vars` from `typedScopeArena.allocVarSlots` for function parameter scopes and `LET`/`CONST`/`CLASS` block scopes, and added `varBloom uint64` on `TypedScope`.
3. **`ReferenceCollector` `t.LookupVar` & `OptimizeCalls` 4x Growth + `cap <= 65536` Retention (`Subagent 19C`, `internal/jscomp/reference_collector.go`, `internal/jscomp/optimize_calls.go`):**
   - Wired `t.LookupVar(name)` into `ReferenceCollector_CollectorCallback.Visit`, deferring `t.GetScope()` until `v != nil` and `varFilter` passes.
   - Added `optimizeCallsAddNode` with 4x geometric slice growth and raised `releaseOptimizeCallsRefMap` capacity retention from `4,096` to `65,536`.
4. **Word-Sampling String Hash & 2-Way Set-Associative `concatCache` in `GlobalNamespace`, plus `cachedInferConstsCollector` & `t.LookupVar` in `InferConsts` (`Subagent 19D`, `internal/jscomp/global_namespace.go`, `internal/jscomp/infer_consts.go`):**
   - Replaced 5-byte `gnConcatHash` with 64-bit word-sampling `gnHashStr` across the full string (`0`, `n-8`, `(n-8)>>1`, `n>>2`, `n-16`, `(n*3)>>2`) and a 2-way set-associative `concatCache [65536][2]string`, eliminating `.prototype` hash collisions, and wired `t.LookupVar(topName)` into `getRootNodeOfTopNameWithTraversal`.
   - Added 1-slot `cachedInferConstsCollector`, 4x state slice growth, and `t.LookupVar(name)` in `InferConsts`.
5. **Zero-Allocation O(1) Ring-Buffer `FixedPointGraphTraversal` for `LinkedDirectedGraph` (`Subagent 19E`, `internal/jscomp/graph_fixed_point_graph_traversal.go`, `internal/jscomp/graph_linked_directed_graph.go`):**
   - Added `inWorkSet bool` in the unused struct padding of `LinkedDirectedGraph_LinkedDiGraphNode[N, E]` and implemented `computeFixedPointLinked` using a power-of-2 ring buffer, eliminating `31.75 MB` of `LinkedHashSet` allocations and `O(Iterations^2)` tombstone scans in `PureFunctionIdentifier`, `AnalyzePrototypeProperties`, `AmbiguateProperties`, and `DisambiguateProperties`.

### Final measured results — Round 19 (`t=1` and `t=4`, 1 warmup + 5 interleaved runs, PGO enabled)

| Benchmark | Round 18 (`bf0be5bd` `t=1` / `t=4`) | **Round 19 (`t=1` / `t=4`)** | Round 19 vs Round 18 (`t=1` / `t=4`) | Peak RSS (`t=1` / `t=4`) | Parity |
|---|---:|---:|---:|---:|:---:|
| **`bundle-mega`** (36.1 MB, 26 UMD libs) | 24.12 s / 23.27 s | **23.33 s** / **22.54 s** | **−3.28%** / **−3.15%** | **1,675** / **1,657 MiB** (**−3.61%** / **−5.23%**, **−91.5 MiB**) | `MATCH` |
| **`typed-mega`** (44.6 MB, 4,385 files) | 19.32 s / 18.18 s | **18.74 s** / **17.70 s** | **−3.00%** / **−2.66%** | **1,738** / **1,656 MiB** (**−3.86%** / **−1.67%**, **−69.7 MiB**) | `MATCH` |
| **`typed-full`** (28.5 MB, 1,818 files) | 11.10 s / 10.37 s | **10.71 s** / **10.18 s** | **−3.53%** / **−1.82%** | **1,040** / **1,071 MiB** (**−1.00%** / **−1.77%**) | `MATCH` |
| **`bundle-advanced-ts`** (12.6 MB) | 6.88 s / 6.75 s | **6.68 s** / **6.55 s** | **−3.02%** / **−2.95%** | 546 / 754 MiB | `MATCH` |
| **`bundle`** (12.6 MB) | 4.49 s / 4.39 s | **4.36 s** / **4.21 s** | **−2.78%** / **−3.94%** | 667 / **518 MiB** (**−4.17%** `t=4`) | `MATCH` |
| **`shaka-blockly`** (8.6 MB) | 4.57 s / 4.35 s | **4.45 s** / **4.23 s** | **−2.71%** / **−2.76%** | 507 / 491 MiB | `MATCH` |
| **`bundle-advanced`** (5.4 MB) | 3.48 s / 3.32 s | **3.36 s** / **3.21 s** | **−3.54%** / **−3.39%** | 329 / 308 MiB | `MATCH` |
| **`checks-only`** (11.3 MB) | 2.54 s / 2.15 s | **2.50 s** / **2.11 s** | **−1.47%** / **−1.90%** | 525 / **520 MiB** (**−3.98%** `t=4`) | `MATCH` |
| **`typed-app`** (2.8 MB) | 1.43 s / 1.32 s | **1.41 s** / **1.30 s** | **−1.61%** / **−1.85%** | 264 / **250 MiB** (**−2.78%** `t=4`) | `MATCH` |
| **`downlevel-simple`** (2.1 MB) | 1.16 s / 1.05 s | **1.11 s** / **1.02 s** | **−4.23%** / **−3.61%** | **193** / **190 MiB** (**−3.19%** / **−5.35%**) | `MATCH` |
| **`downlevel-whitespace`** (2.1 MB) | 0.55 s / 0.46 s | **0.54 s** / **0.46 s** | **−2.34%** / **−1.79%** | 161 / 164 MiB | `MATCH` |

## Round 20 (Phase 2 / Bucket B): Dual-Mode `compat` vs. `fast` + 4-Layer Equivalence Gate, O(1) Generation-Stamped `ControlFlowAnalysis`, Concrete `DataFlowAnalysis` Worklist, and Parallel Function-Local Middle-End Passes (`DeadAssignmentsElimination`, `CoalesceVariableNames`, `FlowSensitiveInlineVariables`, `InlineObjectLiterals`) (2026-10-21)

Phase 2 (Bucket B — dual-mode `compat` vs. `fast`, 4-layer equivalence gate, and parallel function-local middle-end passes) delivered five structural improvements:

1. **Dual-Mode `compat` vs. `fast` (`--fast_mode` / `JSGO_FAST_MODE=1`) & 4-Layer Equivalence Gate (`Subagent 20A`, `internal/jscomp/compiler_options.go`, `internal/jscomp/command_line_runner.go`, `tools/bench/fast_mode_gate.py`):**
   - Added `CompilerOptions.FastMode()` / `SetFastMode(bool)`, CLI flag `--fast_mode`, and `JSGO_FAST_MODE=1` environment override. Default `compat` mode remains 100% byte-for-byte identical to Java (`t=1 == t=2 == t=4`), while `fast` mode unlocks optimizations whose output differs only in deterministic temporary/pseudo-name numbering while preserving 100% Go-to-Go multi-thread determinism (`t=1 == t=2 == t=4`).
   - Built [`tools/bench/fast_mode_gate.py`](../tools/bench/fast_mode_gate.py), verifying 21 benchmark + CLI e2e workloads across 4 layers:
     - **Layer 1 (Go-to-Go Multi-Thread Determinism):** Identical `(rc, stdout, stderr, js, map)` SHA-256 across `t=1`, `t=2`, and `t=4` under `-race`.
     - **Layer 2 (Canonical Alpha-Normalized Equivalence):** Identical SHA-256 after canonicalizing compiler-generated internal temporaries (`$jscomp$inline_*`, `JSCompiler_inline_label_*`, `JSCompiler_object_inline_*`).
     - **Layer 3 (Raw & Gzip Byte-Size Non-Regression):** `fast_raw_bytes <= compat_raw_bytes` and `fast_gzip_bytes <= compat_gzip_bytes` (`0.00%` regression across all 21 tasks).
     - **Layer 4 (Structural AST & Diagnostic Census):** Exact match on exit code, warning/error lines, function count, call count, statement count, and identifier count.
2. **O(1) Generation-Stamped `posTable` in `ControlFlowAnalysis` & Concrete `LinkedDiGraphNode` Worklist in `DataFlowAnalysisBase` (`Subagent 20B`, `internal/jscomp/control_flow_analysis.go`, `internal/jscomp/data_flow_analysis.go`, `internal/jscomp/live_variables_analysis.go`):**
   - Replaced `c.astPosition *jutil.JavaHashMap[*rhino.Node, int]` (`O(N)` allocations/probing + `node.HashCode()`) in `ControlFlowAnalysis` with a pre-sized open-addressed power-of-2 `posTable []cfaPosEntry` using 32-bit `posGen uint32` generation stamps (`O(1)` reset with zero clearing) and `c.posCounter` node priority.
   - Replaced interface-boxed `DiGraphNode` maps/sets in `DataFlowAnalysisBase` (`orderedWorkSet` and `stepCountMap`) with a concrete `[]*LinkedDirectedGraph_LinkedDiGraphNode[*rhino.Node, ControlFlowGraph_Branch]` slice sorted by `node.GetPriority()`, a flat `inWorkSet []bool` bitmap, and a flat `stepCounts []int32` array indexed by `node.GetNodeIndex()`.
   - Precomputed `l.varHasSameName []bool` once in `LiveVariablesAnalysis` so `MarkAllParametersEscaped` and `isRearrangeableAssign` skip `scope.Parent()` walks when a parameter/variable is unshadowed.
3. **Pre-Filtered & Parallelized `DeadAssignmentsElimination` and `CoalesceVariableNames` (`Subagent 20C`, `internal/jscomp/dead_assignments_elimination.go`, `internal/jscomp/coalesce_variable_names.go`, `internal/rhino/node.go`):**
   - Added a single pre-pass `daeCandidateScanner` and parallel worker pool in `DeadAssignmentsElimination.Process` (`448 ms → 168 ms` on `bundle-mega`, **2.7× faster**), skipping non-candidate functions without constructing scopes or CFGs and executing CFG + `LiveVariablesAnalysis` + dead assignment elimination in parallel across `NumParallelThreads()`, followed by deterministic pre-order `ReportCodeChange` replay.
   - Added `cvnCandidateScanner` and parallel worker pool in `CoalesceVariableNames.Process` (`642 ms → 290 ms` on `bundle-mega`, **2.2× faster**, and `235 ms → 68 ms` on `typed-mega`, **3.5× faster**) in both `compat` and `fast` modes (`usePseudoNames == false` produces 100% byte-identical AST mutations in parallel; `usePseudoNames == true` uses deterministic per-function local sorting in `fast` mode).
   - Fixed a redundant write (`tail.shared = true` on already-shared `SOURCE_FILE` property lists) in `rhino.Node.SetStaticSourceFileFrom` exposed under `-race`.
4. **Parallelized `FlowSensitiveInlineVariables` (`compat` + `fast`) (`Subagent 20D`, `internal/jscomp/flow_sensitive_inline_variables.go`):**
   - Parallelized `FlowSensitiveInlineVariables.Process` across `NumParallelThreads()` (`1,062 ms → 249 ms` in `compat` mode and `242 ms` in `fast` mode on `bundle-mega`, **4.4× faster**!):
     - Phase 1 runs `ControlFlowAnalysis`, `SyntacticScopeCreator`, `MustBeReachingVariableDef`, `MaybeReachingVariableUse`, and ` gatherCands` in parallel for each candidate function using per-worker pooled arenas (`fsivWorkerState`).
     - In `compat` mode, Phase 2 replays `cand.inlineVariable` sequentially in deterministic AST pre-order using the exact shared `UniqueIdSupplier`, preserving 100% byte-identical output to Java; in `fast` mode, Phase 2 inlines directly inside the parallel workers using deterministic function-seeded local ID suppliers.
5. **Candidate-Filtered & Parallel `InlineObjectLiterals` + `CheckUnreachableCode` / `CheckMissingReturn` / `CombinedCompilerPass` Fast Paths (`Subagent 20E`, `internal/jscomp/inline_object_literals.go`, `internal/jscomp/check_unreachable_code.go`, `internal/jscomp/check_missing_return.go`, `internal/jscomp/combined_compiler_pass.go`):**
   - Added `iolCandidateScanner` to `InlineObjectLiterals.Process` (`983 ms → 694 ms` in `fast` mode on `bundle-mega`), skipping functions with no local `OBJECTLIT` assignments/declarations and running function-local `ReferenceCollector` + object-literal splitting/inlining in parallel in `fast` mode.
   - Added `hasUnreachableStack []bool` in `CheckUnreachableCode` to skip `CheckMissingReturn` CFG traversals (`checkControlFlow`: `298 ms → 180 ms` on `bundle-mega`, **1.7× faster**), `cmrMightHaveEmptyReturn` pre-filter in `CheckMissingReturn`, and `noVisit` / `noShouldTraverse` single-callback specializations in `CombinedCompilerPass`.

### Final measured results — Round 20 (`t=1` and `t=4`, 1 warmup + 5 interleaved runs, `compat` mode, PGO enabled)

All **22 of 22 (`11 benchmarks × {t=1, t=4}`)** configurations improved in Round 20, with `bundle-mega` (`t=4`) dropping by **−1.50 s (−6.72%)** in a single round (`22.35 s → 20.85 s` in `compat` mode, and **`20.55 s`** in `--fast_mode`, down **−2.95 s / −12.4%** since Round 17) and `typed-mega` (`t=4`) dropping by **−0.63 s (−3.57%)** (`17.64 s → 17.01 s`, down **−1.57 s / −8.5%** since Round 17):

| Benchmark | Round 19 (`68cad1f2` `t=1` / `t=4`) | **Round 20 (`t=1` / `t=4`)** | Round 20 vs Round 19 (`t=1` / `t=4`) | Peak RSS (`t=1` / `t=4`) | Parity (`compat` & `fast` L1–L4) |
|---|---:|---:|---:|---:|:---:|
| **`bundle-mega`** (36.1 MB, 26 UMD libs) | 23.16 s / 22.35 s | **22.37 s** / **20.85 s** (`20.55 s` `fast`) | **−3.43%** / **−6.72%** (**−1.50 s**) | **1,700** / 1,797 MiB (**−1.77%** `t=1`) | `MATCH` |
| **`typed-mega`** (44.6 MB, 4,385 files) | 18.67 s / 17.64 s | **18.50 s** / **17.01 s** | **−0.95%** / **−3.57%** (**−0.63 s**) | **1,703** / **1,707 MiB** (**−1.04%** / **−4.72%**) | `MATCH` |
| **`typed-full`** (28.5 MB, 1,818 files) | 10.70 s / 10.11 s | **10.58 s** / **9.80 s** | **−1.16%** / **−3.05%** | 1,042 / 1,119 MiB | `MATCH` |
| **`bundle-advanced-ts`** (12.6 MB) | 6.60 s / 6.49 s | **6.40 s** / **6.16 s** | **−3.01%** / **−5.09%** | **576** / 743 MiB (**−19.67%** `t=1`) | `MATCH` |
| **`bundle`** (12.6 MB) | 4.33 s / 4.21 s | **4.14 s** / **3.87 s** | **−4.41%** / **−8.18%** | 538 / 565 MiB | `MATCH` |
| **`shaka-blockly`** (8.6 MB) | 4.45 s / 4.23 s | **4.35 s** / **4.03 s** | **−2.25%** / **−4.65%** | 481 / 514 MiB | `MATCH` |
| **`bundle-advanced`** (5.4 MB) | 3.36 s / 3.21 s | **3.26 s** / **3.01 s** | **−3.04%** / **−6.20%** | **323** / 321 MiB (**−3.26%** `t=1`) | `MATCH` |
| **`checks-only`** (11.3 MB) | 2.51 s / 2.11 s | **2.44 s** / **2.05 s** | **−2.61%** / **−2.74%** | 531 / 541 MiB | `MATCH` |
| **`typed-app`** (2.8 MB) | 1.41 s / 1.29 s | **1.38 s** / **1.25 s** | **−2.49%** / **−3.08%** | **256** / 267 MiB (**−3.59%** `t=1`) | `MATCH` |
| **`downlevel-simple`** (2.1 MB) | 1.10 s / 1.02 s | **1.09 s** / **0.97 s** | **−1.17%** / **−4.46%** | 198 / 204 MiB | `MATCH` |
| **`downlevel-whitespace`** (2.1 MB) | 0.54 s / 0.45 s | **0.53 s** / **0.45 s** | **−0.15%** / **−0.21%** | **150** / **158 MiB** (**−5.58%** / **−0.96%**) | `MATCH` |

## Round 21 (Phase 3 / Post-Architecture Turn 1): `NodeTraversal` Leaf-Bypass Markers (`noVisitLeaves` / `onlyNameLeaves`, `−0.78 s` CPU), `RemoveUnusedCode` & `GuardedCallback` Leaf Fast-Exits, Lazy `ReferenceCollector` & `InferConsts` `BasicBlock`s, Direct `OptimizeCalls` Walk, `Normalize` String Slabs (`−12.9 MB`), and 2-Way `TypeInference.concatQName` Cache (`−20 MB`) (2026-10-21)

Profiling Round 20 (`5c782211`) across `bundle-mega` (`3.59 GB` `alloc_space`) and `typed-mega` (`2.95 GB` `alloc_space`) revealed newly exposed traversal and per-pass bottlenecks after the Phase 1–2 architectural upgrades:

1. **`NodeTraversal` Leaf-Node Dispatch Bypass (`noVisitLeaves` & `onlyNameLeaves`) + `InlineFunctions` / `FunctionInjector` Leaf Unrolling (`internal/jscomp/node_traversal.go`, `internal/jscomp/inline_functions.go`, `internal/jscomp/function_injector.go`):**
   - Added internal marker interfaces `nodeTraversal_noVisitLeaves` (`noVisitLeaves()`) and `nodeTraversal_onlyNameLeaves` (`onlyNameLeaves()`) to `NodeTraversal`, skipping `Visit` interface dispatch (and `t.currentNode = child` when `noShouldTraverse` is also set) on normal leaf AST nodes across `TraverseBranch` and `TraverseChildren` (`−0.78 s` cumulative CPU in `TraverseBranch` on `bundle-mega`).
   - Wired `noVisitLeaves()` on `InlineFunctions_FindCandidateFunctions`, `InlineFunctions_CallVisitorBase`, and `InlineFunctions_Inline`, and `onlyNameLeaves()` on `InlineFunctions_FindCandidatesReferences`. Precomputed a 64-bit `paramBloom` mask in `inlineFunctions_analyzeBodyNode`, and unrolled leaf children (`c.GetFirstChild() == nil`) inline without recursion in `inlineFunctions_analyzeBodyNode`, `inlineFunctions_hasLetOrConst`, `inlineFunctions_findCalledFunctions`, `functionInjector_blocksMinimumRequirements`, `functionInjector_referencesEval`, and `functionInjector_findInnerFunctions`.
2. **`RemoveUnusedCode` Leaf Fast-Path (`traverseNameReference`) & `GuardedCallback` / `PolyfillUsageFinder` Leaf Bypass (`internal/jscomp/remove_unused_code.go`, `internal/jscomp/guarded_callback.go`, `internal/jscomp/polyfill_usage_finder.go`):**
   - Extracted `traverseNameReference(n, scope)` and added leaf-node (`GetFirstChild() == nil`) fast-exits at the top of `RemoveUnusedCode.TraverseNode` and inside `TraverseChildren`, skipping `TraverseNode` switch dispatch on all non-`NAME` leaf AST nodes, and enlarged `scopeVarCache` (`[512] → [2048]`) and `pinnedPropCache` (`[256] → [1024]`) with `scopeNameEq` (`−0.19 s` cumulative CPU in `TraverseNameNode`, `−0.11 s` flat in `TraverseNode`).
   - Added leaf-node bypass in `GuardedCallbackBase.ShouldTraverse` for `polyfillUsageFinder_Traverser` (skipping context stack push/pop and `VisitGuarded` on non-`NAME` leaves and non-static-polyfill `NAME` leaves) plus end-byte bitmasks (`staticsEndMask`, `suffixesEndMask`, `methodsEndMask`, `anyPropMask`, `anyPropEndMask`) in `PolyfillUsageFinder_Polyfills`.
3. **Lazy `BasicBlock` Materialization & Leaf Skipping in `ReferenceCollector` and `InlineObjectLiterals` (`internal/jscomp/reference_collector.go`, `internal/jscomp/inline_object_literals.go`):**
   - Replaced eager `allocBlock` allocations on every control-flow boundary in `ReferenceCollector_CollectorCallback.ShouldTraverse` with lazy `[]basicBlockEntry` stack entries (`getCurrentBlock()` / `materializeBlock(idx)`), allocating `BasicBlock` structs only when a tracked variable reference is actually recorded inside that block, and skipped leaf nodes other than `NAME`, `IMPORT_STAR`, `BLOCK`, and `SWITCH_BODY`.
   - Unrolled leaf children in `iolCandidateScanner.scan`.
4. **Direct `OptimizeCalls` AST Walk, `InlineSimpleMethods` / `InferConsts` Leaf Bypass, and Pooled `Normalize` String Slabs (`internal/jscomp/optimize_calls.go`, `internal/jscomp/inline_simple_methods.go`, `internal/jscomp/infer_consts.go`, `internal/jscomp/normalize.go`):**
   - Replaced `NodeTraversal_TraverseRoots` in `OptimizeCalls.Process` with a direct `cb.traverseNode` recursive AST walk with leaf unrolling (`−70 ms` in `OptimizeCalls`).
   - Added `noVisitLeaves()` / `noShouldTraverse()` on `InlineSimpleMethods_GatherSignatures` and `InlineSimpleMethods_InlineTrivialAccessors`, leaf fast-exit and lazy `BasicBlock`s in `InferConsts`, `onlyNameLeaves()` / `noShouldTraverse()` on `Normalize_PropagateConstantPropertyOverVars` and `ValidityCheck_VerifyConstants`, stricter `normalizeHasJscompSeparator`, and pooled `normalizeFastRenamerPoolState` + 32 KB `strSlab` in `Normalize` (`−12.9 MB` `alloc_space` on `typed-mega`, `−8.9 MB` on `bundle-mega`).
5. **2-Way `TypeInference.concatQName` Cache (`−20 MB`), `CheckAccessControls` Token Table & `Var` Cache, and `PeepholeOptimizationsPass` Token Pre-Filters (`internal/jscomp/type_inference.go`, `internal/jscomp/check_access_controls.go`, `internal/jscomp/peephole_optimizations_pass.go`):**
   - Upgraded `TypeInference.concatQName` to a 2-way set-associative `[2048][2]` cache backed by `nodeUtilConcatDot` (`−20 MB` `alloc_space` on `typed-mega`), cached `(lastScopeRoot, lastModule, lastSyntacticBlockScope)` in `FlowThrough`, and fast-pathed `TightenTypeAfterDereference`.
   - Added a `[256]bool` `checkAccessControlsActiveToken` table, a 2-way `[512][2]cacVarCacheEntry` `(scope, name) -> *Var` cache, and deferred `NodeUtil_GetBestJSDocInfo` in `CheckAccessControls`; added token pre-filters in `PeepholeOptimizationsPass_PeepCallback.Visit`.

Total `alloc_space` on `typed-mega` dropped another **−80 MB** down to **`2.87 GB`** (**−75.3% vs Round 5's `11.61 GB`**) and `inuse_space` dropped **−207.8 MB**.

### Final measured results — Round 21 (`t=1` and `t=4`, 1 warmup + 5 interleaved runs, `compat` mode, PGO enabled)

| Benchmark | Round 20 (`5c782211` `t=1` / `t=4`) | **Round 21 (`t=1` / `t=4`)** | Round 21 vs Round 20 (`t=1` / `t=4`) | Peak RSS (`t=1` / `t=4`) | Parity (`JAVA_MATCH`) |
|---|---:|---:|---:|---:|:---:|
| **`bundle-mega`** (36.1 MB, 26 UMD libs) | 22.29 s / 20.95 s | **22.26 s** / **20.85 s** | **−0.11%** / **−0.46%** | 1,741 / **1,701 MiB** (**−4.82%** / **−86 MiB** `t=4`) | `MATCH` |
| **`typed-mega`** (44.6 MB, 4,385 files) | 18.44 s / 17.09 s | **18.37 s** / **17.02 s** | **−0.41%** / **−0.43%** | 1,746 / 1,843 MiB | `MATCH` |
| **`typed-full`** (28.5 MB, 1,818 files) | 10.55 s / 9.80 s | **10.54 s** / **9.77 s** | **−0.07%** / **−0.28%** | 1,046 / 1,072 MiB | `MATCH` |
| **`bundle-advanced-ts`** (12.6 MB) | 6.40 s / 6.14 s | **6.38 s** / **6.10 s** | **−0.30%** / **−0.56%** | 671 / **563 MiB** (**−1.91%** `t=4`) | `MATCH` |
| **`bundle`** (12.6 MB) | 4.13 s / 3.85 s | **4.15 s** / **3.86 s** | +0.55% / +0.25% | **536** / 545 MiB (**−0.63%** `t=1`) | `MATCH` |
| **`shaka-blockly`** (8.6 MB) | 4.34 s / 4.01 s | **4.33 s** / **4.00 s** | **−0.37%** / **−0.28%** | **491** / **509 MiB** (**−1.35%** / **−1.38%**) | `MATCH` |
| **`bundle-advanced`** (5.4 MB) | 3.26 s / 3.00 s | **3.24 s** / **2.99 s** | **−0.69%** / **−0.20%** | **324** / 326 MiB (**−0.46%** `t=1`) | `MATCH` |
| **`checks-only`** (11.3 MB) | 2.45 s / 2.05 s | **2.43 s** / **2.04 s** | **−0.81%** / **−0.48%** | 538 / **522 MiB** (**−2.00%** `t=4`) | `MATCH` |
| **`typed-app`** (2.8 MB) | 1.38 s / 1.26 s | **1.39 s** / **1.26 s** | +0.48% / **−0.19%** | 265 / 271 MiB | `MATCH` |
| **`downlevel-simple`** (2.1 MB) | 1.09 s / 0.97 s | **1.09 s** / **0.97 s** | +0.12% / +0.37% | **197** / **204 MiB** (**−0.96%** / **−0.37%**) | `MATCH` |
| **`downlevel-whitespace`** (2.1 MB) | 0.54 s / 0.45 s | **0.54 s** / **0.45 s** | **−0.09%** / +1.20% | **149** / **159 MiB** (**−0.35%** / **−2.71%**) | `MATCH` |

## Round 22 (Phase 3 / Post-Architecture Turn 2): `nodeFlagHasJSDoc` Bitmask (`−0.79 s` CPU), `Scope` / `TypedScope` / `SyntacticScopeCreator` Fast Paths (`−1.86 s` CPU), `CoalesceVariableNames` / `LiveVariablesAnalysis` / `DataFlowAnalysis` Early Bailouts (`−0.79 s` CPU), and Compact 24-Byte `SourceMap_Mapping` (`−40%` Memory) (2026-10-21)

Profiling Round 21 (`f3eb8b4c`) on `bundle-mega` and `typed-mega` identified four remaining post-architecture hotspots across `rhino.Node` property lists, scope lookup/scanning, data-flow analysis, and source-map generation:

1. **Inline `nodeFlagHasJSDoc uint16 = 1 << 4` Bitmask on `rhino.Node` (`internal/rhino/node.go`, `internal/jscomp/node_util.go`):**
   - Repurposed bit 4 of `Node.boolFlags` as `nodeFlagHasJSDoc` (keeping `sizeof(rhino.Node) == 80`). Since ~98% of AST nodes have no `JSDOC_INFO` property, `n.GetJSDocInfo()` and `n.RemoveProp(Node_Prop_JSDOC_INFO)` now return immediately after a 2-instruction bitmask test without walking `n.propListHead` (`−0.34 s` flat in `Node.LookupProperty`, `−0.26 s` in `Node.RemoveProp` across `MarkNoSideEffectCalls` / `PureFunctionIdentifier`, and `−0.19 s` in `NodeUtil_GetBestJSDocInfo`).
2. **`Scope`, `TypedScope`, and `SyntacticScopeCreator` Lookup & Scanner Fast Paths (`internal/jscomp/scope.go`, `internal/jscomp/typed_scope.go`, `internal/jscomp/syntactic_scope_creator.go`, `internal/jscomp/abstract_scope.go`, `internal/jscomp/variable_reference_check.go`, `internal/jscomp/rename_vars.go`, `internal/jscomp/denormalize.go`, `internal/jscomp/make_declared_names_unique.go`):**
   - Eliminated per-ancestor `hasImplicitVars` branches inside the `Scope.GetVar` and `Scope.HasSlot` upward walk loops by handling `arguments` once on entry via `GetArgumentsVar()`, made `scopeBloomMask` bounds-check-free and inlineable (`109 ms → 0 ms`), replaced byte-by-byte FNV-1a in `scopeNameHash` with O(1) word-sampling, and added `TypedScope.lookupRootVar` for `TypeCheck` / `TypeValidator` (`−1.18 s` cumulative CPU across `Scope.GetVar` and `TypedScope.GetVar`).
   - Added `parentBloom` caching, inlined `getInput()`, and unrolled `VAR`/`LET`/`CONST`/`FUNCTION`/`CLASS` child declarations in `SyntacticScopeCreator_ScopeScanner.scanVars` (`−0.68 s` cumulative CPU).
   - Replaced closure-allocating `GetOwnSlotMap().Values()` iterations with direct `scope.varList` slice loops in `VariableReferenceCheck`, `RenameVars`, `Denormalize`, and `MakeDeclaredNamesUnique`.
3. **`CoalesceVariableNames`, `LiveVariablesAnalysis`, and `DataFlowAnalysis` Early Bailouts & `onlyNameLeaves()` (`internal/jscomp/coalesce_variable_names.go`, `internal/jscomp/live_variables_analysis.go`, `internal/jscomp/data_flow_analysis.go`, `internal/jscomp/flow_sensitive_inline_variables.go`, `internal/jscomp/inline_object_literals.go`):**
   - Added `onlyNameLeaves()` to `CoalesceVariableNames_Callback`, hoisted `lvaAllVarsFinder.collectWithLimit(fnScope, LiveVariablesAnalysis_MAX_VARIABLES_TO_ANALYZE)` before `ComputeEscaped` and `UndirectedGraph` construction in `CoalesceVariableNames.EnterScope` and `FlowSensitiveInlineVariables`, and added `onlyNameLeaves()` + early `nonEscapingDepth > 0` skip to `dataFlowAnalysis_escapedFinder` (`−0.79 s` cumulative CPU on `bundle-mega`).
   - Reused `behavior` per worker with `behavior.staleVars.Clear()` in `InlineObjectLiterals.Process` (`--fast_mode`), fixing a recycled `*Var` arena pointer bug between `candFns` in single-threaded `--fast_mode`.
4. **Compact 8-Byte `sourcemap.FilePosition` & 24-Byte `SourceMap_Mapping` (`−40%` Memory) (`internal/sourcemap/file_position.go`, `internal/sourcemap/source_map_generator_v3.go`, `internal/jscomp/source_map.go`, `internal/jscomp/code_printer.go`, `internal/jscomp/compiler.go`):**
   - Changed `FilePosition.line` and `column` from `int` (16 B) to `int32` (8 B) and reordered `SourceMap_Mapping` fields to eliminate padding, shrinking `SourceMap_Mapping` from 40 B to **24 B (−40%)**.
   - Raised `codePrinter_mappingSlabSize` from `512` to `4096`, iterated `mappingChunks`/`mappingSlab` by value without constructing an intermediate `allMappings []*SourceMap_Mapping` pointer slice in `Compiler.appendSourceAndMappingsToCodeBuilder`, and added pointer-first string comparison in `SourceMapGeneratorV3.AddMappingPositions`.

### Final measured results — Round 22 (`t=1` and `t=4`, 1 warmup + 5 interleaved runs, `compat` & `fast` modes, PGO enabled)

All **22 of 22 (`11 benchmarks × {t=1, t=4}`)** configurations improved in Round 22 over Round 21 in `compat` mode (100% byte-identical to Java), with `--fast_mode` (`t=4`) reaching **16.13 s** on `typed-mega` (**2.20× faster than Cold Java**, **1.70× faster than Warm Java**) and **20.22 s** on `bundle-mega` (**2.39× faster than Cold Java**, **2.04× faster than Warm Java**):

| Benchmark | Round 21 (`f3eb8b4c` `t=1` / `t=4`) | **Round 22 `compat` (`t=1` / `t=4`)** | **`--fast_mode` (`t=4`)** | Round 22 vs Round 21 (`t=1` / `t=4`) | Cumulative vs Round 18 (`t=4` `compat` / `fast`) | Peak RSS (`t=1` / `t=4`) | Parity |
|---|---:|---:|---:|---:|---:|---:|:---:|
| **`bundle-mega`** (36.1 MB, 26 UMD libs) | 22.22 s / 20.98 s | **21.67 s** / **20.27 s** | **20.22 s** | **−2.48%** / **−3.38%** | **−12.89%** (**−3.00 s**) / **−13.11%** (**−3.05 s**) | 1,782 / 1,746 MiB | `MATCH` |
| **`typed-mega`** (44.6 MB, 4,385 files) | 18.43 s / 17.08 s | **17.96 s** / **16.72 s** | **16.13 s** | **−2.57%** / **−2.08%** | **−8.03%** (**−1.46 s**) / **−11.28%** (**−2.05 s**) | **1,727** / **1,773 MiB** (**−1.22%** / **−2.89%**) | `MATCH` |
| **`typed-full`** (28.5 MB, 1,818 files) | 10.55 s / 9.76 s | **10.25 s** / **9.50 s** | **9.23 s** | **−2.85%** / **−2.67%** | **−8.39%** (**−0.87 s**) / **−10.99%** (**−1.14 s**) | **1,038** / 1,113 MiB (**−5.11%** `t=1`) | `MATCH` |
| **`bundle-advanced-ts`** (12.6 MB) | 6.40 s / 6.12 s | **6.16 s** / **5.86 s** | **5.85 s** | **−3.73%** / **−4.11%** | **−13.19%** (**−0.89 s**) / **−13.33%** (**−0.90 s**) | 546 / 587 MiB | `MATCH` |
| **`bundle`** (12.6 MB) | 4.14 s / 3.86 s | **4.03 s** / **3.75 s** | — | **−2.74%** / **−2.79%** | **−14.58%** (**−0.64 s**) | **522** / **559 MiB** (**−4.31%** / **−2.84%**) | `MATCH` |
| **`shaka-blockly`** (8.6 MB) | 4.31 s / 3.98 s | **4.18 s** / **3.88 s** | **3.75 s** | **−2.95%** / **−2.65%** | **−10.80%** (**−0.47 s**) / **−13.79%** (**−0.60 s**) | **464** / 517 MiB (**−1.93%** `t=1`) | `MATCH` |
| **`bundle-advanced`** (5.4 MB) | 3.24 s / 2.99 s | **3.14 s** / **2.91 s** | — | **−2.82%** / **−2.93%** | **−12.35%** (**−0.41 s**) | **319** / 312 MiB (**−1.70%** `t=1`) | `MATCH` |
| **`checks-only`** (11.3 MB) | 2.43 s / 2.03 s | **2.39 s** / **2.00 s** | — | **−1.56%** / **−1.75%** | **−6.98%** (**−0.15 s**) | **545** / **519 MiB** (**−2.46%** / **−1.50%**) | `MATCH` |
| **`typed-app`** (2.8 MB) | 1.39 s / 1.25 s | **1.35 s** / **1.23 s** | **1.22 s** | **−2.65%** / **−1.89%** | **−6.82%** (**−0.09 s**) / **−7.58%** (**−0.10 s**) | **256** / **265 MiB** (**−1.37%** / **−1.23%**) | `MATCH` |
| **`downlevel-simple`** (2.1 MB) | 1.09 s / 0.98 s | **1.01 s** / **0.90 s** | — | **−7.41%** / **−8.16%** | **−14.29%** (**−0.15 s**) | **196** / 204 MiB (**−0.74%** `t=1`) | `MATCH` |
| **`downlevel-whitespace`** (2.1 MB) | 0.54 s / 0.45 s | **0.52 s** / **0.43 s** | — | **−3.34%** / **−3.69%** | **−6.52%** (**−0.03 s**) | **148** / **155 MiB** (**−8.54%** / **−2.33%**) | `MATCH` |

## Memory round 1 (2026-10-23): peak heap and pointers instead of wall time

Slabs, arenas, pools and cross-pass caches from the earlier rounds were the
main source of excess heap: removing them (plus fixing several retention
leaks and shrinking the type structs) cut the peak live heap by 34–62 % on
all 11 benchmarks, byte-identical with Java, at the cost of more allocation
(wall time +8–13 % on the large tasks with Go 1.27.3). Details, numbers and
the round-2 plan: [memory-research/round1-results.md](memory-research/round1-results.md).

## Smaller open items

- `tools/clie2e/sweep.go` runs its JVMs under its own `nice -n 10`; under
  an outer `nice` they end at 19 and run ~5× slower when other work
  competes. Drop the inner `nice` or document it.
- The released jar still cannot load its own XTB reader
  (`NoClassDefFoundError …jarjar/javax/xml/parsers/ParserConfigurationException`);
  `tools/clie2e/build-java.sh` works around it for the oracles.
