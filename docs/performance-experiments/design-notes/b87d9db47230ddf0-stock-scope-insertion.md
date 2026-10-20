# Prechecked insertion in private syntactic scope construction

`JSGO_EXPERIMENT_STOCK_SCOPE_INSERTION=1` selects a bounded source experiment
for DefaultPassConfig's early/main/constant InlineVariables passes and builtin
RemoveUnusedCode. Public scope creators, InlineVariables constructors and
RemoveUnusedCode builders retain the original construction path.

The existing bundle attribution reported approximately 517 MiB sampled
allocation and 8.9% CPU below syntactic declaration scans across all consumers.
Those are overlapping upper bounds for the broad subsystem, not savings
available from this narrower change. This experiment changes no allocation
representation and does not eliminate scope scans, Scope, Var or map objects.

An ordinary new declaration currently searches its own declaration map five
times: scanner getOwnSlotV; Scope.Declare's repeated getOwnSlotV assertion;
DeclareInternal's HasOwnSlot; CanDeclare's HasOwnSlot; LinkedHashMap.Put.
The scanner already proved the name absent before exact stock Compiler.GetInput
reads its stable input map. The private path removes the second absence check
and the redundant HasOwnSlot short-circuit, reducing these source-level own-map
searches from five to three. Parent parameter and shadowing checks remain;
the first declaration's nil-map checks should not be confused with hash work.

The original CanDeclare legality check is deliberately retained. Scanner-level
parameter-shadow handling alone does not prove all implicit-name and bleeding
function-name cases. Empty-name checks, NewVar construction/validation order,
Var indices, ordinary linked-map Put, duplicate handlers, goog namespace
replacement, canonical scope/variable identity and iteration behavior remain.
No borrowed string keys, dense-map iteration, snapshots or reused scope objects
are introduced.

Eligibility is established once when an audited private pass opens its creator
window: exact stock Compiler, actual compiler roots, default nonnil redeclaration
handler, no provide-as-redeclaration mode, no custom passes/observers/error handler,
and known conventions/error manager. No per-declaration environment lookup or
broad observer fence is added. Creator state is copied into the local scanner;
public scanner initialization clears it. The private pass closes the window
by defer, including an exception. Scope and Var layouts remain unchanged.

Optional caller-owned counters report scopes, original insertions and prechecked
insertions; they allocate nothing when absent. Differential fixtures compare
ordered declarations, exact binding/node/input relationships, lookup results,
fresh identities, duplicate/parameter/hoist/catch/destructuring/module/extern/
namespace handling, custom callback schedules, pass output/change stamps,
private-factory eligibility and exception cleanup. Normal/race/compiler-oracle
gates and a matched default-GC CLI screen are still required. No timing benefit
is claimed from search counts or historical attribution.

Broader cross-pass Scope reuse remains a separate unproven design. Exact-parent
identity causes global-declaration churn to invalidate descendants; node change
times do not propagate to ancestors, and direct public mutation requires replay.
The old cache additionally grew Scope's allocation class and retained mutable
Scope/Var graphs. A future immutable declaration-recipe cache could replay into
fresh scopes instead, but would save only scanner walks, needs explicit root/
input/parent invalidation, and must first demonstrate weighted reuse opportunity.
This insertion change avoids that dependency and lifetime problem entirely.
