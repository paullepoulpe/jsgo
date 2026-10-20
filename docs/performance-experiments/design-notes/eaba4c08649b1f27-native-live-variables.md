# Dense live-variable fixed-point experiment

`JSGO_EXPERIMENT_NATIVE_LIVE_VARIABLES=1` selects an invocation-owned liveness
solver for the fresh builtin CoalesceVariableNames and DeadAssignmentsElimination
passes. Public constructors and public `LiveVariablesAnalysis.Analyze` retain the
original engine. The backend additionally requires exact stock CFG/edge objects,
stock compiler observers/conventions/reporters, a private fresh-CLI source
ownership witness, and at most 128 variable indices. Existing consumers already
limit normal analysis to about 100 locals. Unsupported cases fall back before
effects or annotations change. There is no GC tuning or persistent cache.

The current engine allocates two temporary BitSets, a copied output lattice and
its BitSet at every transfer; multi-edge joins allocate additional lattices. It
also recomputes the same AST gen/kill sets on every revisit. This prototype keeps
in/out/gen/kill as two-word values in one per-function slice, tracks queue
membership in those rows, and uses a concrete integer heap with the original
priority comparator, insertion order and tie rules. It preserves incoming and
outgoing edge order, duplicate-edge enqueue behavior, zero-successor joins,
implicit return, step counts and the original divergence limit.

`ComputeGenKill` is deliberately reused rather than reimplemented. In particular,
reading `arguments` calls `MarkAllParametersEscaped` during the transfer. Eagerly
computing all gen/kill sets or premarking parameters would change the original
sequence and intermediate states. Instead each node computes gen/kill at its
original first transfer using two reused scratch BitSets. The cache is versioned
by escaped-set size; a transfer during which that set grows is not cacheable.
Later transfers recompute after any growth. Thus even names preceding an
`arguments` read in the same expression observe the original escaped set.

At return, including exception unwinding, the backend publishes ordinary
`LinearFlowState<LiveVariableLattice>` annotations with ordinary BitSets. Existing
coalescing/interference and dead-assignment consumers remain unchanged. The
private owned interval excludes callbacks that could observe delayed intermediate
publication. Failure publication includes the step increment and joined out
state preceding a failed transfer. Public annotation object identity and public
custom graph/callback schedules remain on the original path.

Parallel coalescing inherits eligibility only from its guarded coordinator and
creates an independent solver inside each owned function job. No queue, cache,
mask, or mutable scratch BitSet is shared. The optional test probe records
coordinator-side runs only, avoiding shared worker counters.

Source-only differential tests compare every node's final in/out values, hash,
step count, and escaped insertion order, including loops, conditional writes,
try/catch/finally, optional chains, closures, default/destructured parameters,
self updates, arguments-induced escape and the 64-bit word boundary. Full
consumer fixtures compare generated source, complete AST, diagnostics and change
stamps. A malformed VAR after an arguments read checks original failure text and
partial annotation publication. Public/unowned/custom-observer fallbacks require
zero native dispatch. Actual transfer reuse/escape transitions/consumer dispatch
must occur; merely comparing equal output is insufficient.

The V2 publication experiment additionally allocates one flat uint64 backing for
all result BitSets. Each set adopts an exact three-index window whose capacity is
ceil(variable-count/64), then trims zero words without losing reserved capacity.
Direct word copies replace per-bit Set loops. This removes two tiny word-array
allocations per CFG node without changing BitSet size/hash/empty/growth behavior:
growing beyond a window must allocate private storage and cannot overwrite the
next state. Public headers/lattices already share invocation-owned arrays, so
the flat word backing adds no new cross-function lifetime. Dedicated zero/one/
two-word and disjoint mutation/clone/growth tests supplement the full state and
consumer oracles. The ordinary BitSet constructors and representation remain
unchanged; the new helper transfers explicit caller ownership.

This is an unvalidated prototype. It retains the original generic CFG and
publishes standard result objects, so it does not eliminate CFG allocation or
the interference graph. Its additional per-function rows/maps may outweigh
saved iterations for tiny functions. Correctness suites, race checks and
isolated default-GC CLI time/RSS/allocation measurements are required before
adoption or performance claims.
