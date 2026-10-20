# Immutable function-analysis prototype

Enable this experimental path with `JSGO_EXPERIMENT_FUNCTION_ANALYSIS=1`.
It uses the existing `--num_parallel_threads` budget and never changes
`GOMAXPROCS`. The original implementation remains the default and fallback.

The first real consumer is CoalesceVariableNames. Its private work covers CFG
construction, liveness, the variable interference graph, and deterministic graph
coloring together. It does not merely overlap two reaching-definition solvers.

## Ownership and ordering

1. On the compiler goroutine, discover and materialize canonical lexical scopes
   in stable traversal order using the pass's existing MemoizedScopeCreator.
   Store only a function scope and revision header for each possible job.
2. When a job's result is needed, prepare at most one workers-sized window of
   variable maps and escape sets serially. Escape discovery resolves nested
   names, including lazy implicit scope variables, before workers run.
3. Each worker owns its CFG nodes/annotations, liveness states/worklist, escaped
   set, interference graph and coloring. Its compiler facade exposes options
   and local exception formatting only; accidental use of another compiler API
   fails the job instead of reaching shared mutable compiler state.
4. Join every job in the window before returning to the original traversal.
   There are no AST writes while workers read. Consume summaries in the original
   function-entry order and use the existing rename/diagnostic/change-reporting
   code unchanged.
5. Before consumption, compare the function's parent, every enclosing function
   and script identity, and their change timestamps with the snapshot. Reparenting
   or an earlier outer rewrite invalidates a result. Discard it and execute the
   unchanged serial implementation at that original traversal position. Worker
   failures also fall back there, preserving exception/report order.

Only one window of pending CFG/lattice/coloring results is retained. Variable
maps transfer to their liveness result or are released after the window. The
normal pass's stack can additionally retain active enclosing analyses, as it
already did. Canonical scope metadata remains cached for the pass; materializing
it earlier can affect peak memory and must be measured.

## Shared-state audit

The CFG callbacks read AST and options while mutating their private graph only.
Liveness uses a read-only name-to-Var map and HasOwnSlot queries; its GetVar calls
are restricted to explicitly declared simple formal parameters for arguments
handling. Escape discovery and all ScopeCreator access remain serial. Graph
construction and coloring share canonical Var identities but mutate no Scope,
Var declaration, AST or type-registry data. Var's cached name hash already uses
atomic loads/stores; immutable names and scope roots determine it.

Upfront scope materialization is specific to this consumer: the original
CoalesceVariableNames liveness constructor already invokes ComputeEscaped over
nested functions before outer rewrites, using the same memoized scope creator.
Captured variables are escaped and excluded from coalescing. An ancestor that
fails the eligibility limit does not rewrite its function. This is not a claim
that the same scope timing is safe for arbitrary passes.

TypedAST compilers, compilers with existing errors, and nonstandard compiler
implementations keep the original path. The experiment does not parallelize
renaming, type inference, scope creation, diagnostics or compilation requests.

## Validation provided

`TestImmutableFunctionAnalysisMatchesOriginalCoalescing` compares the original
path with the experimental path at 1, 2 and 4 workers, with and without pseudo
names. Cases cover independent and nested functions, captured outer bindings,
arguments/two-parameter behavior, loops, catches, default/destructuring parameters,
named function expressions, initializer reparenting, classes and class-field
captures. It asserts that actual worker results are consumed, preventing an
all-fallback implementation from silently passing the comparison.

Additional tests verify ancestor/parent invalidation, repeated real compiler
change stamps, stale-result rejection, bounded pending results, and two genuinely
overlapping jobs behind a blocking wrapper around the real CFG/liveness/coloring
strategy. Existing CoalesceVariableNames tests, race checks and end-to-end
Java/output gates are still required before retaining this candidate.

## Expected ceiling and follow-on work

Existing bundle profiles attribute approximately 0.91 of 22.21 sampled CPU seconds
(4.1%) to this entire pass. That is an upper bound on the removable work, not a
predicted speedup. Two available CPUs imply only roughly a 2% ideal overall gain
if all of it parallelized perfectly. Serial preparation, invalidated jobs,
additional scope traversal, load imbalance, allocator/GC contention and memory
bandwidth lower that ceiling. Small functions can lose to scheduling overhead.
Measure the actual wall time and peak RSS before deciding whether to retain it.

The more general result is an ownership boundary: canonical immutable inputs,
private per-function analyses, bounded execution and revision-validated ordered
publication. The indexed CFG/data-flow experiment can run inside these jobs
because their graphs and solver states are private. Other consumers such as
FlowSensitiveInlineVariables need their own invalidation and canonical-symbol
proof before reusing this protocol. Function side-effect summaries followed by
ordered call-graph reduction remain a larger future target; concurrent whole-pass
mutation is not a substitute for that split.
