# Bounded no-op proofs in required traversals

Base `4f417579`. Private stock factories select
`JSGO_EXPERIMENT_EMPTY_PASS_PROOFS=1`; public constructors retain their original
pass schedules. Custom compilers, roots, conventions, error/change observers,
custom passes and debug/tracing modes keep the original path. The pass remains
in the schedule, so PhaseOptimizer start/end reporting is unchanged. No change
stamp or cross-pass idempotence assumption is used.

## Block-scoped declarations

The existing first traversal collects LET/CONST declarations and performs its
usual checks/initialization. If the resulting set is empty, the second
LoopClosureTransformer traversal and declaration-rewrite phase cannot change
anything. Its reference callback rejects every Var which is neither LET nor
CONST before recording a reference or constructing loop objects; those predicates
test declaration token kinds, not inferred constness or JSDoc. Its other work
depends on the recorded maps. The traversal root creates its own global scope
from that root, so external scopes cannot introduce unseen LET/CONST targets.

The proof skips only after the original gather. Feature-removal bookkeeping still
runs, including on the empty path. Nonempty sets retain the complete original
transformer and declaration rewrite. Reusing a pass object with old nonempty
state can only retain extra original work, not cause an invalid skip.

The fresh parser-enabled CPU profile attributes about 30ms/0.57% of
downlevel-simple sampled CPU to this entire pass. Other sampled fixtures have
less visible work here. This is an upper bound on opportunity, not a measured
improvement or evidence of how often the empty path occurs. Existing tests plus
private probe tests must verify actual skips, initialized/closure-captured VAR
loops, real LET/CONST loop captures, output, diagnostics, ordered change reports,
feature bookkeeping and custom fallback.

At source freeze there are no new build/test/timing results; central validation
is pending the root scheduling slot.

## Denormalize declaration collection

The private wrapper performs the original mandatory first traversal, forwarding
the same callbacks before observing whether any visited NAME has a direct VAR
parent and no initializer children. This is precisely a necessary condition for
the second ReferenceCollector phase's `IsVarDeclaration` and
`!IsInitializingDeclaration` test. All such names count, including globals and
loop-header declarations that may not actually qualify, so uncertainty retains
the complete original collector. Destructuring bindings do not have a direct VAR
parent and do not satisfy the original reference predicate either.

The first traversal's transformations move existing VAR nodes into loop headers
or replace arithmetic/logical assignment operators. They create no VAR nodes and
remove no initializer from a VAR's NAME. A bare declaration therefore cannot
appear after its observation without having been visited earlier. There is no
additional tree/name scan. Candidate state resets on every Process invocation;
any observed candidate runs the original collector with the original callback.
Synthetic-block contexts retain the original path. Public Denormalize is unchanged.

Sampled whole-pass CPU is about 200ms/0.75% for bundle and 60ms/1.14% for
downlevel-simple. Only the collector portion can be saved; real skip frequency
and macro impact are unmeasured. Tests require actual skip/collect dispatch,
nested bare declarations, conservative loop-header handling, destructuring,
reuse of a pass object, custom fallback, and exact output/diagnostic/change-order
parity with the public original path.
