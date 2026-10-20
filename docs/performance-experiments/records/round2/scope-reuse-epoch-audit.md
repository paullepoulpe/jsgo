# Re-audit: cross-pass untyped scope reuse without ownership hooks

Source-only audit, 2026-10-19, against `4f417579`; previous candidate
`56275ea4`. No implementation, build, test or new runtime profile was run.

## Assessment

Removing Scope/Var owner fields and replacing the global pass stamp with
function change times repairs two real defects in the old candidate. It does
not establish a useful hit rate or a low-memory cache. A bounded, private
optimizer cache is technically plausible. Do not revive the old compiler-wide
cache yet: first obtain weighted cross-pass hit opportunities, root declaration
churn and retained-size estimates with an untimed diagnostic.

The quoted ~517 MiB / 8.9% bundle costs cover all syntactic declaration scans,
not just InlineVariables and serial FSIV. They are an upper bound on the target,
not the saving available to this cache. A validation replay retains the scan
CPU; it primarily avoids Scope, Var and declaration-map allocation.

## What change times actually guarantee

`ChangeTracker.recordChange` sets exactly one SCRIPT/FUNCTION node's change
time to the current stamp, increments the stamp, and adds that node to the
timeline. It does not mark any ancestor. `ReportChangeToChangeScope` selects
the supplied SCRIPT/FUNCTION exactly. `ReportChangeToEnclosingScope` selects
the nearest change root starting at the supplied node's **parent**, except a
SCRIPT argument selects itself. Reporting on a FUNCTION therefore marks its
enclosing change scope; reporting on its NAME marks that FUNCTION.

Pass-start and change-query stamp increments do not update these stored node
times. Equality of a function's stored time can therefore survive unrelated
passes and changes to sibling functions, unlike the old global-stamp key.
Lexical blocks/classes have no separate tracked time; their conservative epoch
must come from their current enclosing FUNCTION/SCRIPT. A changed function
invalidates its parameter, body and lexical scopes together.

ROOT has no valid tracked epoch. TraverseRoots materializes one global Scope
for the common parent of externsRoot and jsRoot, including both trees.
`ReportFunctionDeleted` only sets deleted and removes the function from the
timeline; it does not increment time, notify observers, or mark a parent.
Function addition/removal/movement must also be reported by its caller against
the containing scope. Compiler synthetic-input insertion/removal does explicitly
report SCRIPT changes, but global validation must also account for changed root,
script identity/order, InputId and CompilerInput bindings, including externs.
RestoreState resets the tracker and reads zero-time ASTs; cache lifetime cannot
cross initialization, restore, separate compilation or unrelated public calls.

Function declaration names are a boundary case. The normal renamers and
InlineVariables explicitly report the containing scope in addition to changing
the function. ChangeVerifier independently compares declaration names in the
parent while excluding function bodies. Thus SCRIPT-only epochs are meaningful
under the stock change-reporting contract, but no propagation in ChangeTracker
enforces that contract. A raw NAME.SetString plus reporting only the NAME can
leave the SCRIPT time unchanged while changing a global declaration.

Neither an unchanged epoch nor a cached Var-map snapshot detects an unreported
AST-only declaration addition/removal. Arbitrary public/custom AST mutation
requires exact declaration replay; a local epoch-only fast path must be confined
to an audited stock optimizer schedule and explicitly rely on its reporting
contract. This assumption must not silently extend to public creators.

## Root and child proof

An exact global declaration rescan can preserve the same global Scope despite
body-only changes and despite SCRIPT epochs changing for non-declaration work.
Use the original scanner's statement traversal, declaration cuts and input
resolution. A hash alone is not an equality proof; collision-safe comparison
against recorded declarations or the expected resulting map is required.
Reject ModuleBody and provide/legacy-goog-module accumulation, as before.

If global declarations differ, create a new global Scope. Exact parent identity
then rejects every cached descendant, even functions whose own time is stable.
This is a major multiplicative hit-rate limit, not a minor root lookup expense.
Fresh module parents similarly block reuse below them. Do not update a cached
global scope in place merely to retain child hits: that weakens the exact-parent
and immutable-binding proof.

For a child hit check current root identity, exact parent Scope, constructor
invariants, Scope.self/root/parent/depth, root parent/token/input context,
arrow/function-expression/class-expression/static-block classification, current
enclosing change-root identity and its recorded time. A moved node or changed
scope chain must miss even when its time happens to match. Parent validation
must precede child lookup, because function-block shadowing depends on parent
parameter declarations.

## Detecting direct Scope/Var mutation without hooks

Map pointer and size are insufficient. Clear/remove/reinsert can retain the map
pointer and count while changing values or iteration order. Record ordered keys
and exact Var pointers, plus declaration node and CompilerInput pointers. Validate
every live entry, its order, and the Var's immutable constructor metadata before
publishing a hit. Derive expected self/scope from the cached Scope and Var pointer,
name from the recorded map key, and index from its original contiguous fresh-scan
position. Reject implicit namespace entries. This detects direct Declare,
Undeclare, ClearVarsInternal, replacement through DeclareInternal, and public
Scope/Var initialization when they alter the visible state. A mutation that
restores exactly the original visible state need not be rejected as history.

A minimal explicit-Var snapshot is approximately 40 bytes per declaration on
amd64: key string16, Var pointer8, name-node pointer8, input pointer8. Copy string
headers, not a borrowed pointer to Var.name: public Var initialization can change
that field. The snapshot must not merely keep a live Var pointer and then compare
the Var to itself. nameHash is a derived cache, not independent binding metadata.

Implicit arguments/this/super/exports require separate care: stock lookups lazily
populate implicitVars during a consumer, so literal comparison with the initial
nil array would reject legitimate reads. Validate their canonical constructor
metadata against the current root and parent, or snapshot only at a proven
consumer boundary. Reusing private implicit identity across passes has to remain
unobservable to any public callback. Scope reinitialization also requires the
Scope.self/root/parent/depth checks above, not just declaration-map checks.

## A smaller exact-replay option

The old cache records scanner events, approximately 40 bytes each, with geometric
slice growth. The target Scope is redundant in a per-scope entry; removing it
still leaves ~32 bytes/event, in addition to live-Var validation snapshots.

A possible simpler validator avoids event snapshots entirely. After validating
the cached ordered live-Var snapshot, replay declarations using a seen bitset
indexed by the original contiguous Var indices. On each first legal declaration,
verify the next declaration index, key, name node and resolved input against the
expected fresh map. Later duplicate declarations use the original default no-op
redeclaration rules; validate shadowing and arguments rules before treating an
event as an explicit binding. At end every expected binding must have been seen
and no extra binding may have appeared. Provide/module encounters reject reuse.
This proves the resulting fresh declaration map without hashing or retaining a
duplicate event list. Scratch costs roughly one bit/Var, with a single word for
<=64 declarations; larger scratch can be entry-owned and cleared per validation.

This is a source proposal, not a proven implementation. It needs differential
tests for duplicate declaration order, the same NAME scanned twice, destructuring,
parameter shadowing, arguments, malformed declarations, input changes, provides,
public reinitialization and ordered map mutation. It still walks every statement
tree on a hit. An epoch-only local shortcut reduces that walk only by accepting
the stock change-reporting contract described above.

## Concrete consumers and remaining lifetime limits

InlineVariables runs in early/main and optional post-localization optimization
loops. Repeated runs of this same stock consumer are the strongest actual
cross-pass path: unchanged functions can reuse exact parent chains while other
functions change. The main loop also runs InlineFunctions, RemoveUnusedCode and
other rewriting passes between requests. Global removals/renames can invalidate
the root and all children; changed functions invalidate all their local scopes.
Actual declaration-weighted survival is not known.

InlineVariables itself uses one ReferenceCollector traversal. Its traversal
already memoizes materialized scopes and deduplicates hoisted-function jumps, so
there is little obvious repeated same-pass creation. Its callback reads scopes
and Vars, stores them in private reference/alias collections, and rewrites ASTs;
it does not directly Declare/Undeclare/Clear or reinitialize scopes/Vars. Those
private collections must remain private and die with the pass. Preserve the
original stale-within-current-traversal map: invalidate future cache requests,
do not change the scope currently being iterated during rewrites.

FSIV runs outside the fixed-point loop before and after main optimization. Its
strongest repeat is the adjacent all-vars/escaped-variable read-only window,
already handled separately by the bounded experiment. Compiler-wide reuse adds
cross-pass benefit only through unchanged chains shared with InlineVariables or
the later FSIV invocation. Parallel cohort facades must remain uncached. Its
candidate/dataflow collections retain Vars during analysis and AST rewriting, so
future invalidation must not mutate those previously handed-out scopes.

Selection should occur only in private default pass factories, not by modifying
public InlineVariables/FSIV constructors or public SyntacticScopeCreator behavior.
Require exact stock Compiler, current roots and normalized stage, stock handler,
stock convention including copied-proxy self checks, known error/change observers,
no custom passes/debug/tracing/determinism callbacks. Recheck eligibility at each
selected pass entry. Merely checking it when a creator is first built does not
cover later observer/configuration changes. No custom callbacks receive cached
Scope/Var identity. Unsupported schedules must use the original creator.

Do not keep cache entries until compilation end by lazy invalidation. Bound the
cache to one private optimizer invocation/loop and release in a defer on success,
early return and panic. Sweep changed/deleted/reparented entries at selected-pass
boundaries, and drop descendants of replaced parent scopes; otherwise old Scope
parents keep invalid generations alive. A map keyed by (root, exact parent) that
simply adds new generations leaks those retained graphs until cache destruction.

## Evidence needed before another implementation

Record private, untimed counters only: requests and possible cross-pass hits by
consumer; declaration counts weighted by hits; global proof success/change; exact
parent misses; per-function epoch changes; module/provide exclusions; number and
size of scopes that must be retained just to anchor descendants; entries surviving
between pass invocations; expected snapshot and scratch capacities. Separate the
existing bounded FSIV-window savings from new cross-pass hits.

Snapshots and cache entries consume memory even on misses. A one-declaration
scope currently costs about 192 bytes (Scope48 + Var64 + linked header32 +
entry48); a 40-byte snapshot plus entry/context/key storage may consume most of
the first avoided reconstruction. Empty scopes cost only48 but may need caching
to keep parent identity for profitable descendants, so a naive per-scope
declaration-count threshold is unsound as an economic model. Evaluate whole
function/parent chains. Cache retention can raise peak live heap/RSS even when
total allocation falls. Do not claim benefit until ordinary CLI time, TotalAlloc,
objects and peak RSS improve together for representative selected workloads.
