# Further compiler parallelism: prerequisites and experiments

This note describes design constraints found while exploring additional uses of
`--num_parallel_threads`. It is not a claim that these changes already improve
end-to-end performance.

## What the parameter controls

The parameter bounds simultaneously active compiler work items in phases that
support parallel execution. Parsing and dependency prebuilding are examples.
A goroutine is a scheduled unit of Go work, not an operating-system thread.
`GOMAXPROCS` independently limits how many goroutines can execute Go code at once
across the process. Runtime and GC goroutines share that CPU budget. Neither
setting is an exact bound on operating-system threads: blocked system calls and
the runtime can require additional threads.

Changing `GOMAXPROCS` inside a compiler request would change the CPU budget of
other compilers, request handling, GC, and embedding applications. Keep it as a
process-level deployment/benchmark setting. Bound the compiler's own parallel
work separately. A useful experiment matrix varies both settings and the JVM's
processor budget; comparing only worker counts would conflate concurrency with
available CPU.

## Code-printing prototype

Each script has a private CodePrinter, text buffer, relative mappings, and
license-event recorder. The calling compiler goroutine then replays events in
input order, applies the existing license de-duplication and mapping offsets,
and concatenates text using the existing rules. Built-in license trackers can
be recorded by eligible source-file transitions. Custom trackers retain the
sequential callback path, and typed extern generation has a serial fallback.
No AST rewriting, type-registry updates, source-map accumulation, or shared
license-tracker calls run concurrently.

Tests compare 1, 2, and 4 workers for binary, chunk, source-array and chunk-aware
license output, source maps, delimiters, strict prefixes, pretty printing,
duplicate licenses and empty files. Actual compiler output gates and race tests
are still required before retaining the prototype.

The existing profiles impose a low upper bound on its broad effect: printing
accounts for approximately 0.24 seconds of 22.21 sampled CPU seconds for bundle
and 0.05 seconds of 6.19 for typed-app. Even ideal printing parallelism cannot
close the overall Java gap. Large WHITESPACE outputs may have a larger fraction;
they require direct measurement. Do not retain the extra machinery merely
because it executes concurrently.

## Why wrapping a whole optimization pass in goroutines is unsafe

The current pass order contains dependencies both between passes and within a
single traversal:

- Lexical scope creation is lazy. SyntacticScopeCreator populates declarations,
  reports redeclarations and creates Var objects; MemoizedScopeCreator owns
  mutable caches. Sharing either creator without an ownership design is unsafe.
- TypedScopeCreator and the type registry update global declarations, named
  types, properties and inference state. Separate source files are not separate
  type universes. Imports, externs, inheritance and recursive functions connect
  them.
- Rewriting passes update the shared AST, scope change tracking, deleted-function
  information, reference graphs, and sometimes globally allocated names.
  Disjoint function bodies can still refer to the same symbols and metadata.
- The existing FlowSensitiveInlineVariables traversal analyzes and rewrites an
  outer function before subsequently visiting nested functions. Precomputing an
  inner function's analysis before those outer rewrites can change which AST and
  declaration state it observes. A cache keyed only by function node identity
  does not establish that the cached facts are still valid.
- Diagnostics and ordered collection iteration feed observable output. Locks
  can prevent data races without preserving the original decision or reporting
  order. Correctness requires deterministic publication, not merely race-free
  mutation.

## A practical design for function-level parallelism

Introduce an explicit analysis barrier rather than concurrent rewriting:

1. In the original traversal order, materialize lexical scopes, canonical Var
   identities, input metadata and required global type facts. Publish an
   immutable snapshot with a revision number for each relevant scope/AST region.
   This phase must retain existing diagnostic order.
2. Collect independent jobs in stable source/function order. A job receives
   immutable metadata and owns its CFG, data-flow states, worklist and result.
   It must not call compiler mutation/reporting APIs, fill a shared lazy cache,
   or allocate observable global names.
3. Run pure analysis jobs with at most `num_parallel_threads` active jobs. Keep
   results and diagnostics in job-indexed slots, not completion order.
4. Commit results on the compiler goroutine in the original order. Validate the
   revisions of all facts each result depended on. If an earlier rewrite
   invalidates a later result, recompute that result against the new state.
5. Establish compatibility with exact JS, maps, diagnostics and status gates,
   repeated determinism tests, then race checks. Measure both wall time and peak
   RSS: concurrently live CFGs and lattices increase memory even if total
   allocations stay constant.

Starting with leaf functions can reduce invalidation dependencies, but selecting
only leaves must not change the original order or semantics of transformations.
Cloning whole functions and splicing them back is not a shortcut: node identity,
source information, type/color metadata and compiler change reporting all need
an explicit mapping and validation scheme.

## Most promising next split: function summaries, then graph reduction

PureFunctionIdentifier already separates function-body observations from a
global reference/call-graph analysis conceptually. A larger design experiment
could compute private per-function side-effect summaries concurrently after
lexical/reference discovery, then merge summaries and solve the global graph in
stable order. The same idea applies to reference gathering and selected escape
analyses. These phases do not need concurrent AST rewrites, which makes their
ownership boundaries easier to specify.

The prerequisite is to extract summaries that contain canonical symbol IDs and
ordered facts rather than shared mutable builder objects. Recursive call-graph
components need a deterministic fixed-point reduction. Diagnostics and any
annotation writes remain in the serial commit phase. This is more invasive than
parallel prebuild and should be evaluated as a separate branch before combining
it with other representation changes.

Per-function CFG construction is also a pure candidate after scope/input state
is frozen. However, building every possible graph eagerly can do work that the
serial pass skips and retains multiple graphs simultaneously. Build only the
jobs the pass will actually consume, and discard stale results after rewrites.

Whole-program type inference is a later project: establish declarations first,
partition the dependency graph into strongly connected components, process
independent components concurrently, and merge type/property updates with the
same deterministic semantics. Scheduling the current inference callbacks
concurrently is not equivalent to that redesign.

## Performance limits

For a parallelizable wall-time fraction `p`, ideal speedup on `n` CPUs is bounded
by `1 / ((1 - p) + p/n)`, before queueing, synchronization, extra allocations,
GC contention and load imbalance. The fraction must be measured in the relevant
configuration; cumulative CPU percentages and nested traversal times are not
independent amounts that can be added together. The current two-CPU quota also
limits what a larger compiler work queue can demonstrate. Report processor
budget, compiler workers, wall time and memory separately for both runtimes.
