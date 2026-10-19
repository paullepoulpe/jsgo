# Lazy exact FSIV use backends on ff06ac15

Clean source experiment based on latest main ff06ac15 plus the shared fresh CLI
lease foundation. This does not transplant the historical native9 framework or
restore eager backward analysis. No worker/cache/GC tuning or SourceFile lock
bypass. No build/test/performance claim yet.

Three independent default-off factors:

- FSIV_USE_CARDINALITY: exact empty/singleton CFG ID/many quotient for the sole
  private exactly-one query; full-domain dense state.
- FSIV_USE_PROJECTION: requires the quotient; reduce state/actions to canonical
  Vars of the current complete remaining candidate set, preserving all aliases.
- FSIV_USE_SCALAR: requires both above; independently solve selected Vars with
  reused O(CFG nodes) scratch and retain only original (Var,defCFG) OUT answers.

Every flag has the JSGO_EXPERIMENT_ prefix. Public NewFlowSensitiveInlineVariables
and MaybeReachingVariableUse retain original behavior. Only the private default
FSIV factory with a genuine stock CLI lease may choose the backend. Current root,
normalized phase, exact stock callback/convention/reporter/warning guards, no
injected typed AST filesystem and no debug/tracer/custom observers are checked.
Loaded built-in runtime AST suppliers are allowed because no transfer invokes
them. The lease changes no SourceFile synchronization or cache operation.

## Preserve latest-main lazy behavior

Main's candidate parameter filter, early loop/RHS/side-effect/use-count checks
and their original order remain unchanged. The existing exactly-one query is the
only new dispatch point. If all candidates fail earlier checks, no sidecar,
program, projection or backward solver is constructed. The first native attempt
is made once at the original first getReachingUses demand; rejection calls the
original lazy full-set analysis there.

No successful FSIV rewrite can precede that first demand: every successful
CanInline must pass exactly-one before InlineVariable runs. Earlier candidates
only failed and were removed. Snapshot the CURRENT remaining list, including the
triggering candidate. It only shrinks thereafter; captured names/Definition.node
never change. Do not pre-run CanInline or bring removed candidates back.

## Graph and algorithm ownership

Forward MustDef/escape/dependency work and canonical Scope/Var identities stay
original. At the first admitted demand, freeze a private immutable integer sidecar
from that exact existing generic CFA graph; do not rewalk AST or rebuild CFA.
Copy original creation-order node IDs/priorities and EACH incoming/outgoing edge
list separately, including ON_EX and implicit-return boundary. Unknown graph,
edge, priority or ownership shapes fall back. No public graph API is replaced.

Compile the full original reverse transfer program, preserving conditional kills,
logical/optional/default/destructuring order and whole-variable updates. Check
ALL original nonescaped Var/CFG-atom hash collisions and aliases, plus the full
original gen-fact convergence bound before any projection/scalar filtering. This
prevents earlier abstract convergence from hiding the original20,000-step failure.
A matrix budget bounds admitted products. Unsupported preparation publishes no
state, and original getReachingUses produces its original errors/order.

Projection retains every spelling aliasing a selected canonical Var, including
noncandidate spellings. Forward Definition.depends stays complete. Backward
transfer components only gen/whole-kill/union their own Var, making projection
exact. Empty selection still validates/proves the entire original domain first.

Scalar preparation validates every original candidate pair before discarding
dense state, builds a private CSR transpose of selected kill/gen actions and
retains exact OUT values keyed by canonical Var ID + original defCFG ID. Reset
IN/OUT/kinds/heap/queued/steps before every Var. All reads in one CFG node share
its atom, and many can become singleton only through whole-component kill. The
immutable program never consults the rewritten AST during result lookup. Scalar
lattice/result storage is O(nodes+querypairs); program/offset/map storage remains
additional. Repeated heap setup can outweigh shorter scalar joins and must be
measured independently from full/projected dense controls.

The native result/sidecar is retired before current EnterScope returns, including
panic. Stock CLI code cannot observe private forward-graph annotations during
backward transfer; public/custom/debug paths retain original annotation behavior.
There is no cross-function, cross-pass or cross-compilation cache.

## Validation and counters

Caller-owned numeric probes are nil normally. They distinguish gathered functions/
candidates, actual first-query attempts, fallbacks, original/selected columns and
actions, full gen-fact proof universe, scalar pairs/heap operations and state
bytes. Matrix/scratch counts exclude maps, capacities/rounding and RSS.

Tests compare the frozen sidecar's exact topology, every selected dense/scalar
IN/OUT state and OUT decision against the ORIGINAL full-set solver, singleton CFG
identity, aliases/escaped/arguments/loop/exception/default cases, excluded hash/
malformed validation, cap bound and self-loop reset hazards. Complete actual CLI
lifecycle exercises all8 flag combinations, positive original-vs-native dispatch,
public/unowned/custom-observer fallback and an actual gathered candidate that
fails main's early loop check without constructing any backward backend.

Central normal/race/Java artifact gates, real pinned-job admission/projection
counts and matched time/RSS screens are required on this current base. Historical
old-base measurements are not evidence for the latest-main implementation.
