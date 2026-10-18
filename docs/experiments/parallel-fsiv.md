# Parallel flow-sensitive inlining experiment

`JSGO_EXPERIMENT_PARALLEL_FSIV=1` enables an opt-in ordinary-CLI prototype.
The existing `--num_parallel_threads` setting bounds concurrently active script
jobs; the pass never changes process-wide `GOMAXPROCS`. Serial execution remains
the default and the fallback. There is no worker process or extern cache.

## Implemented ownership boundary

The compiler goroutine creates one canonical syntactic global Scope from the
actual combined externs/JavaScript ROOT. Every script traversal receives that
same Scope, preserving global Var identities and declaration-node aliases.
Creating separate global parents for different scripts would change binding
identity. FSIV only rewrites function-local expressions, so global declarations
remain intact. A ROOT cannot lazily create implicit arguments/this/super/exports
variables, and Var's cached name hash already uses atomic loads/stores.

Each job owns a separate FSIV instance, scope creator, descendant scopes,
CFG/lattices, candidates, and analyzer settings. It runs the original traversal,
analysis, and rewriting order throughout its script. An outer function therefore
finishes its rewrites before traversal subsequently processes its body functions.
No ASTs are cloned, and no function analysis is speculatively reused after a
mutation. The independent native-transfer-facts experiment can operate inside
these private analyses with its existing function-generation lifetime.

A small compiler facade exposes audited read operations and buffers diagnostics
and scope-change reports. Enclosing change scopes are resolved when reported:
one FSIV rewrite reports a parent before detaching it, so deferred ancestor lookup
would be incorrect. After joining jobs, publication follows original script order
and preserves report multiplicity, change stamps, timeline order, and the earliest
input-order exception. On an unexpected internal exception, later jobs may already
have changed their ASTs; serial partial in-memory AST equivalence after failure is
not guaranteed. This prototype is not a general concurrent compiler API.

Custom compiler implementations, coding conventions, observers/error handlers,
debug/tracing configurations, unsupported roots/inputs, and insufficient work
retain serial execution. Raw JSType nodes are rejected when type-based local
optimization would inspect them, avoiding shared type-resolution caches.
AstAnalyzer's color comparison is against the primitive STRING color: its empty
sets prevent equality from calling the otherwise lazily cached Color.HashCode.
Accessor summaries and input maps are read-only during this pass.

## Limits and validation

Whole-script scheduling cannot balance one dominant input or a single large
bundle. Concurrent jobs also increase live analysis scratch space, while buffered
events remain retained until publication. Both elapsed time and peak RSS must be
measured before retaining the design.

The older bundle CPU profile attributes 2.07 of 22.21 sampled seconds (9.3%) to
FSIV. Perfect two-CPU parallelization of that fraction would save about 4.7% of
total time; serial scope/type preparation, load imbalance, scheduling and GC reduce
the practical ceiling. This is a profile-based bound, not a measured improvement
or a Java speed ratio.

New tests compare serial and parallel output, AST structure, change stamps and
timeline order at two/four workers with native facts off/on. Fixtures include
cross-script globals, nested functions/arrows, loops, arguments, catches,
destructuring, classes, modules, and initializer moves. Additional tests check
shared canonical global Vars, fallback, primitive colors, before-detach reporting,
and ordered exception-prefix publication. These tests, race checks, and CLI/Java
output gates are required; their current results belong in the benchmark report.

## Further function-cohort feasibility: not implemented

A giant IIFE can defeat both script and top-level-statement scheduling. A possible
finer boundary is a cohort of non-overlapping complete function subtrees inside
an enclosing function, using actual traversal cut points:

1. Process functions nested in default-parameter expressions before the enclosing
   body's FSIV work. NodeTraversal visits PARAM_LIST before entering BLOCK; a
   blanket parent-before-descendants depth schedule changes the observed AST.
2. Finish the enclosing body's complete FSIV analysis and rewriting before any
   body-descendant job starts. Escape discovery visits descendant references;
   RHS safety checks deliberately inspect names inside function-expression bodies;
   ancestor inlining can detach or move those function roots. These are concrete
   ancestor-read/descendant-write and ancestor-write/descendant-read conflicts.
3. Discover child roots after those rewrites, retaining their actual canonical
   enclosing Scope chains. Before sharing these scopes, materialize required
   implicit slots serially: function/module/field scopes, unlike ROOT, can lazily
   allocate arguments/this/super/exports Vars when sibling closures resolve names.
4. Run independent sibling subtrees with their internal traversal order intact,
   publishing events in original depth-first order. Use one bounded execution
   budget; recursively creating an N-worker pool per function oversubscribes it.

FSIV already skips analysis for functions with more than 100 local variables, so
the serial enclosing-IIFE barrier might be inexpensive. That does not remove the
scope-ownership and traversal-order requirements or establish an actual speedup.
