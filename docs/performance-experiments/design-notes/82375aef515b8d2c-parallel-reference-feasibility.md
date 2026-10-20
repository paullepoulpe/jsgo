# ReferenceCollector parallelization feasibility

Source/profile investigation only; this branch does not replace collection or
claim a performance gain.

## The apparent hotspot includes serial continuation work

In the existing bundle profile, CollectorCallback.Visit accounts for 2.58 of
22.21 sampled CPU seconds (11.6%). This is cumulative, not independently movable
work. MaybeJumpToHoistedFunction accounts for 1.53 seconds beneath that callback:
it recursively traverses a hoisted function when its name is first referenced.
That traversal determines reference order and BasicBlock parentage.

AddReference accounts for approximately 0.91 seconds (4.1%), including about
0.37 seconds in ordered object-map lookup. Lookup/find values are nested, not
additive. The unconditionally pure InferConsts caller accounts for only about
0.23 seconds (1%); restricting parallel collection to DO_NOTHING behavior does
not expose the apparent 11.6% opportunity. PureFunctionIdentifier has no samples
in this particular bundle profile and is not an evidenced replacement target.

## Existing behavior callbacks are part of the traversal

InlineVariables, InlineObjectLiterals and Denormalize transform the AST at scope
exit. VariableReferenceCheck emits ordered diagnostics there. Custom filters,
behaviors and scope creators can have additional effects.

The collector itself is not just a per-function AST walk:

- A first reference may jump to a later hoisted function before source-order
  traversal reaches its declaration. The declaration's NAME reference is added
  explicitly in the containing scope, then the function is traversed at scope.
- That continuation temporarily substitutes a BasicBlock stack appropriate for
  the hoist position. ReferenceCollection uses BasicBlock identity and ancestry
  in its semantic decisions; creating merely equivalent blocks is insufficient.
- Var objects are not stable across separate traversals. Equality uses the pair
  of declaration scope root and immutable name, while Reference retains its
  actual lexical Scope. Preserving variable equality alone does not establish
  that all consumer identity checks remain unchanged.
- Scope resolution may create implicit variables. Materializing all names and
  canonical scopes must happen serially before workers query them.
- Each scope-exit callback must see the same published prefix of references as
  in the original implementation, at the same traversal context and AST state.
  Publishing complete future-reference lists early requires an additional proof
  for each behavior and is not equivalent for arbitrary callbacks.

## Candidate event-tape architecture and unresolved proof

A possible larger design is to run the original continuation/control traversal
as a preparation phase. Record canonical reference contexts and scope-exit
contexts in exact event order. Workers could build private Reference headers and
symbol lists, followed by ordered prefix publication and serial behavior calls.
The saved traversal context must copy its scope stack, preserve exact Scope and
BasicBlock pointers, and materialize lazy state in the caller. No workers may
read AST or scopes while replayed callbacks mutate them.

For InlineVariables and InlineObjectLiterals, a promising property is that own
lexical-variable references live within an already visited subtree when its
scope exits. InlineVariables defers aliases of an enclosing variable until that
enclosing variable is handled, after the enclosing scope finishes. Object
splitting inserts replacement declarations in the exited scope/hoist context and
marks affected earlier references stale.

This source inspection does not prove that all future traversal decisions and
scope creation remain unchanged after those rewrites. Required proof cases
include hoisted recursive continuations, named function expressions, default
parameter scopes, block/catch/class scopes, delayed outer-alias retries, function
initializer reparenting and generated declarations. Unknown callbacks, predicates,
scope creators and coding-convention callbacks must retain the original path.
If a rewrite changes a future decision, a real continuation/restart mechanism is
needed; blindly replaying a stale tape is not a fallback.

Even if the semantic proof succeeds, a tape that leaves all canonical traversal
and lookup serial only moves part of the 0.91-second AddReference cost. Private
maps followed by a serial deterministic merge can duplicate the original map
work and add buffers, sorting, allocations and peak RSS. A low-risk prototype
must first establish that the removable cost exceeds those additions.

## More direct first experiment

Var equality within the ReferenceCollector key type is scoped-name equality:
`(scopeRoot *Node, name string)`. A collector-local native lookup index could
avoid repeated generic ObjectLinkedHashMap hashing/equality on every reference,
while retaining the original ordered map for first-symbol publication and its
existing API wrappers. New symbols would enter the ordered map at their original
first occurrence, and traversal, continuations, filter calls and behavior timing
would remain unchanged.

This requires exact equality and nil-root exception tests, handling distinct Var
objects that compare equal, preserving map mutation semantics for custom users,
and comparison across all real consumers. It is a representation experiment,
not parallelism. A measured serial reduction of collection costs would also
provide a cleaner baseline for deciding whether event-tape parallelism is worth
its larger semantic and memory cost.
