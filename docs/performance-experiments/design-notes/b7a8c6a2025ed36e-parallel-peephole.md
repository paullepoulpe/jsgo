# Parallel stock peephole script cohorts

Source prototype from 9aee69c4. Opt in with the internal environment flag
`JSGO_EXPERIMENT_PARALLEL_PEEPHOLE=1`. The worker budget is the existing
num_parallel_threads option. No process-wide runtime or GC setting changes.

Early and late peephole labels account for about 4.4% of bundle CPU in the
existing labeled profiles. Perfect two-CPU scaling of all of that work saves
at most roughly 2.2% before discovery, publication, imbalance and uncovered
work. This is an experiment to compose with other improvements; there is no
performance or memory result yet.

## Ownership and exact global rounds

The coordinator retains the original BeginTraversal/EndTraversal, every
GetChangedScopeNodesForPass call, retraverse option and global fixed-point
round. Jobs own entire SCRIPTs within a round. Changed roots belonging to one
file run serially in their original order, including when their slots are
interleaved with another file's roots. Workers buffer effects separately for
each original root slot. The coordinator publishes those slots in original
order before asking the tracker for the next global round. There are no
independent per-file fixed points or reordered name suppliers.

The three exact supported default recipes contain the stock six early
optimizers, seven normalized early optimizers (adding MinimizeExitPoints), or
six late optimizers (including StatementFusion). Each worker has its own
optimizer objects, feature sets, AstAnalyzer header, scope memo and event
buffers. The global scope is prepared serially over all files, with original
initial versus changed-root traversal scope roots; workers do not create an
isolated per-file global environment. Default recipes do not query lexical
scopes, but the canonical environment remains correct if a query occurs.

Only concrete stock compiler/conventions/optimizer selves and known change
handlers qualify. Custom recipes/callbacks, copied dispatch selves, raw JSTypes,
TypedAST loading, tracing/debug/determinism observers, fewer than two workers or
files, and an allowable feature set excluding LET use the original serial path.
An already-known changed set containing only one file also remains serial.
Reachable Color hashes are prepared serially to avoid their lazy cache writes
during typed equality. The facade leaves unaudited compiler methods nil.

## Metadata is part of correctness

Buffered changes do not touch scope stamps or the global tracker in workers.
Function deletion flags are also deferred: setting DELETED before publication
would cause an earlier buffered change to be ignored, losing its original
stamp. New detached bound-function clones register private function ownership
before attachment. A private NodeUtil_MarkNewScopesChanged interface hook records
those events without sharing the concrete ChangeTracker; ordinary compilers
keep their original route.

Some other CloneTree paths do not mark function clones new. If a child function
changes before a parent copies it, serial cloning copies its new ChangeTime;
simple buffered reporting would copy the old stamp. The first proof rejects
function-containing arithmetic/compound-assignment clone sources and potential
bind callees, including callee forms exposed by child constant folding. These
guards can reject a whole input set. Dispatch coverage must be measured; a
declined experiment is not parallel work. No negative stamps or whole-tree
metadata remapping are introduced.

SCRIPT is not a lexical scope. Original newFeatures can therefore remain
pending after a script and flush at another script's first scope exit or at
the single initial global exit. Creating a global exit for each worker would
change feature metadata. Workers record scope-exit checkpoints and unflushed
per-optimizer carry instead; the coordinator replays that accumulation with
the original optimizer order and global exit. The only feature introduced by
these stock recipes is LET, while their only script-feature query checks
SPREAD_EXPRESSIONS. Delaying LET publication therefore does not change worker
optimization decisions. Subsequent changed-root traversals keep their original
absence of an extra final global-scope feature flush.

## Evidence required before retention

Source tests require real worker changes at N=2/N=4 and multiple global rounds,
exact JS/AST, lifecycle, script/allowable feature sets and ordered scope stamps.
They exercise nonadjacent roots from the same file, cross-file global bindings
and shadows, feature carry both to a later inner scope and the global root
exit, changes before function deletion, and each conservative clone guard.
A blocking probe wraps actual job execution and proves the shared N=2 bound.
Public flag tests compare N=1 and N=2 to the original path.

Source-only peer review found no definite semantic counterexample. No tests,
builds, race runs, activation profiles, timing or RSS measurements were run by
the author during the reserved source slot. The central gate must run focused
cohort tests plus existing peephole oracles, race tests and all CLI output gates.
Then record actual recipe/file/round dispatch coverage and matched fresh-CLI
time and peak RSS. Preparation, global-scope construction, event buffering,
strict clone fallback and large-file imbalance can erase this small ceiling.
Partial in-memory AST state after an internal worker failure is outside the
prototype; original-root report prefix and earliest failure are preserved.
