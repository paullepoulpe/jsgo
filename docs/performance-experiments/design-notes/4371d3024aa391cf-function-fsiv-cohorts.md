# Function-subtree FSIV cohorts

`JSGO_EXPERIMENT_PARALLEL_FSIV=1` enables this ordinary-CLI experiment. The
existing `--num_parallel_threads` budget bounds active jobs; `GOMAXPROCS` remains
independent. There are no worker processes, extern caches, or SourceFile changes.
The unchanged serial pass remains the default and unsupported-case fallback.

The coordinator retains the original NodeTraversal. When a function BLOCK is
entered, it first runs the parent's complete existing FSIV analysis and rewrites.
It then discovers the nearest descendant FUNCTION roots, pruning their subtrees,
and runs cohorts containing at least two roots. Each job runs the original pass
through its entire function subtree. The coordinator subsequently skips those
processed roots. Workers never construct coordinators or nested worker pools.

This exposes parallel work inside a single dominant IIFE instead of assigning
one whole input file to one worker. A subtree node count supplies a scheduling
weight: larger jobs start first, while change reports and exceptions publish in
original traversal order. The counters/blocking wrappers used by tests are
caller-owned optional probes, absent from ordinary CLI runs. Eligibility alone
does not prove that any cohort actually executed.

## Ordering and ownership

Default-parameter functions run before the enclosing body, exactly as in the
original traversal; static numeric function-depth scheduling would be wrong.
The parent body then finishes before any body-descendant job starts. Parent
escape discovery reads descendant references, RHS safety checks inspect nested
function bodies, and parent inlining can move function roots. These are concrete
ancestor/descendant conflicts, so discovering or analyzing descendants earlier
would need a different proof.

Discovery runs after parent rewrites using its actual canonical body Scope.
Collector-created block/loop/class scopes and function scopes retain that same
parent chain, preserving shared Var aliases. Before workers start, all shared
ancestor implicit slots are materialized serially. Sibling arrows can otherwise
race while creating the same enclosing `arguments` Var; function/module/field
scopes also differ from ROOT in their implicit `this`, `super`, and `exports`
behavior. Declaration maps remain intact: FSIV moves expressions and removes
assignment statements, rather than deleting local declaration names.

Each worker owns its scope creator, descendant scopes, CFG, lattices, candidates,
analyzer settings and buffered events. Its compiler facade exposes only audited
reads. Enclosing change scopes are resolved when reported, before a reported
assignment parent can detach. After joining, replay preserves notification
multiplicity, stamps and timeline order. Shared SourceFile exception-line
formatting is serialized through a cohort-owned mutex; no atomic-source-file
representation is introduced.

Custom compiler implementations, conventions, observers/error handlers and
debug/tracing configurations retain serial execution. Type-based local
optimization rejects raw JSType nodes that might access shared lazy type caches.
With type inspection disabled those types are unused. Primitive color comparison
against STRING uses empty comparison sets and does not call Color's lazy hash.

On an internal failure, all submitted jobs finish; only the original input-order
event prefix and earliest error publish. The coordinator stops remaining work.
Later jobs can already have changed their private AST subtrees: exact partial
in-memory AST equivalence after failure is outside this opt-in CLI prototype.

## Validation and limits

Tests compare serial output, AST structure, diagnostics, change stamps and
timeline order at two/four workers. Explicit fixtures cover parameter-before-body
ordering, post-rewrite root discovery, shared implicit arguments, genuine overlap
of real analyses, largest-first dispatch, the N budget, source formatting,
deterministic failure prefixes and conservative fallback.

Full-lifecycle tests use stock CLI compiler/options factories, builtin externs,
SIMPLE and ADVANCED with type warnings, async runtime injection, final JS and
source maps. A probe at the real pre-main-loop stage must execute user-source
jobs; runtime helpers or an all-fallback path cannot satisfy the assertion.
`TestFunctionFSIVRealInputCoverage` optionally accepts
`JSGO_FSIV_COVERAGE_INPUT=/path/to/bundle.js` for non-timed coverage of real input,
logging actual jobs and summed subtree node weights. Profiling/tracing runs
deliberately fall back and cannot establish parallel-job coverage.

The current cut is one cohort per serially reached function body; a worker keeps
all its descendants serial. A single dominant child can still limit balance.
Discovery, weight counting, scope materialization and event retention cost CPU
and memory, and N live analyses can increase peak RSS. The enclosing IIFE barrier
can be cheap when its more-than-100-local-variable limit skips FSIV analysis.
Only measured ordinary-CLI elapsed time and peak RSS can establish a benefit.
