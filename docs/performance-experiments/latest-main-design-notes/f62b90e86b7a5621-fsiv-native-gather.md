# CFG-local FSIV gather experiment

`JSGO_EXPERIMENT_FSIV_NATIVE_GATHER=1` enables a private builtin-only gather
backend. It is independent of the three backward-use backend flags and leaves
main's forward analysis, candidate filtering, cheap-check order, and lazy first
backward-use query unchanged. Public pass constructors and the original public
CFG callback keep the original traversal.

The existing fresh-main bundle allocation profile attributes 390.7 MiB to
`newNodeTraversal`; 68.2% is below FSIV, including 159.6 MiB below
`GatherCandidates.Visit`. These are cumulative allocated bytes, not live memory
or RSS. Removing gather constructors addresses part of this allocation target;
it does not imply a proportional CPU or wall-time improvement. No benchmark has
been run for this factor.

Each eligible function's gather establishes exact stock compiler, genuine fresh
CLI lease, normalized current-root, stock convention/diagnostic/callback, and
observer fences once. Loaded builtin runtime AST tables are allowed: this
callback never reads or invokes their suppliers. The proof is shared with the
current-main private backward-use admission helper, without depending on the
backward-use flags.

The concrete walker preserves the original CFG boundary test using live AST
parents, root callback parent `nil`, postorder, and captured-next child walk.
Function names, parameters and bodies retain their original traversal order.
This exact callback never observes scopes or traversal state; scope-stack work
is omitted. The default `NodeTraversal_Traverse` builder does not request special
execution ordering for destructuring/default nodes, so those keep ordinary child
order. Class, class-members, module and script special handlers remain original.

Candidate discovery reads AST, canonical forward Definition/Var objects, and
pure stock convention state. A reusable per-gather pointer slice stages fresh
candidate records, without publishing them to the pass's candidate set. A
complete walk publishes candidates in original insertion order. An unsupported
special node or any speculative panic discards staging and runs the original
traversal, which reproduces its own error context and partially published
candidate prefix. The fallback lazily creates one original NodeTraversal and
reuses it only in this private read-only function-gather window. No AST clone,
global cache, cross-pass state, or change to GOMAXPROCS is involved. Staged pointer
slots are cleared on completion or panic.

The unchanged public CFG callback is the oracle. Focused tests compare ordered
node/parent events and exact candidate name/Definition/use/CFG identities across
branches, loops, defaults, destructuring, optional chaining, catch/finally,
closures and class fallback. A malformed NAME after a real candidate checks
original exception formatting and the exact partial-candidate prefix. Complete
stock SIMPLE compilations assert positive post-colors activation independently
with full, projected, scalar, and original backward-use analyses; custom,
public, unleased, observer and copied-convention paths must fall back. Numeric
probes are caller-owned and nil in ordinary/timed runs.

The full-pass differential, race, actual workload activation and fair CLI
benchmarks remain mandatory before promotion. The public and fallback traversal
remain available as controls. ComputeDependence and Candidate.GetDefinition
traversals are deliberately unchanged in this initial experiment.
