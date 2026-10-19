# Thread budget and compiler phase map

`--num_parallel_threads=N` bounds compiler jobs in phases that explicitly use
the option. Each pool starts at most `min(N, work items)` worker goroutines and
joins before publishing results or entering the next phase. A goroutine is a
scheduled unit of Go work, not an operating-system thread. `GOMAXPROCS` separately
limits simultaneous Go execution across the process; it is not set by this
compiler option. Runtime/GC goroutines and threads are not counted as compiler
jobs. Raising N cannot make a serial pass parallel or bypass a two-CPU quota.

## Present ordinary-CLI paths

| Phase | Work unit and publication | Status |
| --- | --- | --- |
| AST prebuild | Independently parse source files, then install roots, diagnostics, comments and source-map data in original input order. | Uses N when multiple independent inputs and supported ordinary parsing are available. Shared SourceFiles, custom compiler implementations and typed-AST loading retain the original path. |
| Dependency prebuild | Scan per-input dependency information with private resolver diagnostic buffers, then replay in input order. | Uses N for supported stock resolvers/input forms; final dependency ordering, pruning and graph updates remain serial. |
| Flow-sensitive variable inlining (FSIV) | Finish the parent function first; discover nearest sibling function subtrees, run private analyses/rewrites, join and publish ordered reports. Default-parameter-before-body order remains intact. | Default-off `JSGO_EXPERIMENT_PARALLEL_FSIV=1` experiment using the same N. Each worker keeps descendants serial; there are no nested worker pools. Custom observers/conventions, raw-type local optimization and debugging configurations fall back. |

The FSIV cohort cut exposes work inside a dominant file or IIFE. A whole-file
pool would leave a nine-megabyte input as one job. A single dominant child still
limits balance, and discovery, scope materialization, weights and buffered
events add costs. Correctness fixtures prove actual overlap and real user-source
jobs in complete SIMPLE/ADVANCED lifecycles; elapsed-time and peak-RSS measurements
must decide whether enabling it is useful.

## Remaining parallel opportunities and barriers

| Phase | Plausible independent unit | Required boundary before parallelizing |
| --- | --- | --- |
| CoalesceVariableNames | Function CFG, liveness, interference and coloring. | Canonical scopes/escaped vars established serially; private graph/results; apply rewrites in original order. An earlier prototype explored this narrow boundary; it is not an active N-controlled phase here. |
| Type inference | Private function CFG/solver work after enclosing facts are ready. | Typed-scope initialization, registry/lazy caches and shared flow updates need an explicit ownership/publication design. The native CFG/solver representation experiment changes storage, not concurrency. |
| Transpilation/normalization | Independent script or function subtrees for a whitelist of passes. | Unique-name suppliers, input-hash collisions, feature/type bookkeeping, global declarations and change/diagnostic callbacks must be isolated or published serially. Arrow/optional-chain-only work was too small in the profile to justify a broad facade. |
| Peephole local rewriting | Sibling function subtrees after enclosing work. | Each pass needs private mutable state plus the same default-parameter and ancestor/descendant ordering proof as FSIV; it is not safe to parallelize an arbitrary traversal callback. |
| Reference collection/variable inlining | Potentially read-only facts before ordered rewrites. | Hoisted-function jumps change collection order; scope-exit inlining changes AST and alias decisions. Arbitrary sibling collection/rewrite cannot be reordered without retaining those dependencies. |
| Output generation | Independent output chunk or script fragments. | Source-map offsets, output separators, license comments and final ordering need private fragment results followed by deterministic assembly. Current code generation does not use N. |

Whole-compiler pass order, global graph construction, global renaming and
cross-function alias/inlining decisions remain dependency-bound. Parallelizing
only proven independent work preserves the mostly drop-in behavior and makes
time/memory comparisons meaningful. More live analyses can increase peak RSS;
every N sweep should report both Java/Go elapsed-time and peak-memory ratios.
