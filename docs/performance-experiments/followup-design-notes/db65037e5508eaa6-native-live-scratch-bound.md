# Reserve private native-liveness scratch for actual variables

This is a separate factor from the one-word BitSet representation. The change
replaces the two unconditional128-bit gen/kill scratch reservations in the
private native solver with `len(analysis.orderedVars)`. It is independently
cherry-pickable onto the original slice, two-word or one-word representation.
It adds no flag or public API and does not change lattice publication.

## Source proof

The original public LiveVariablesAnalysis.FlowThrough remains unchanged. The
private solver is entered only with canonical analysis/CFG ownership and at
most128 variables. NewLiveVariablesAnalysis's AddScopeVariables assigns indices
by enumerating orderedVars into scopeVariables. AddToSetIfLocal is the only setter reached
from ComputeGenKill and uses those same indices. Neither ComputeGenKill nor its
recursive helpers reads scratch Size or capacity; they only Set known variable
bits. The private Compute loop clears the scratch sets between transfers and
copies their bit values into its native two-word state immediately afterward.
No scratch set is stored in the analysis, graph annotation, optional probe,
escaped collection or a callback. The short-lived destructuring visitor
closures synchronously use the same setter, without publishing the sets.

The `escaped.Size()` reads in compute refer to the escaped-variable ordered set,
not BitSet capacity. nativeLiveBitsFromSet only iterates NextSetBit. Published
input/output BitSets retain their preexisting actual-variable reservation and
normal annotation APIs; no externally observable Size changes. Public/custom
fallback and the native eligibility predicate are unchanged.

For ≤64 locals this reduces the original scratch backing from two words to one;
for0 locals it allocates no backing. For65–128 locals the reservation is unchanged.
The one-word representation additionally avoids large-holder construction for
small functions. There are only two scratch sets per native solver invocation,
so this is a bounded supplementary optimization, not a claim that it dominates
whole-compiler allocation. Existing128-bit reservation behavior is not being
removed from published public lattice states.

## Validation state

SOURCE ONLY; no builds/tests or measurements. The existing full native-versus-
generic fixedpoint differential suite now adds exact1/63/64/65/100/127/128-local
functions. Each highest variable is used across a loop backedge, so a bad bound
would change real liveness facts. The existing empty-function, arguments escape,
exceptions, nested closures and full consumer/parallel-worker differential tests
remain required. Constructor capacity distributions and macro measurements must
attribute this scratch factor separately from the representation factor.
