# Release obsolete traversal inline scope entries

This separate followup is based on compact4 source9044aff9. It neither changes
the chosen inline capacity nor adds a runtime flag, pool, field or reuse cache.
No builds or tests have run while the initial macro comparison is active.

At the first append beyond the inline stack, Go copies the entries into new
backing storage. Pops subsequently clear only the new backing. The original
embedded entries would retain popped root/scope/declaration graphs until the
traversal itself becomes unreachable. Both push methods now recognize that
first inline-backed full stack, perform the unchanged append, then clear the
obsolete inline array before invoking EnterScope. A failed append therefore
never destroys the original active entries. Later heap growth and stack reuse
do not clear active entries. The same code works for capacity4/8/16.

GetAbstractScopeAt keeps the logical node value before recursively materializing
parents. A balanced custom creator may push/pop a temporary scope while ancestor
queries are suspended, causing the first growth and obsolete-buffer clearing.
Keeping the node value prevents those suspended queries from reading cleared
inline slots. The live stack still stores the exact originally returned Scope
and parent pointers. Callback order, public depth/root queries, quiet behavior,
materialization order and exact scope/Var identities remain unchanged for normal
or balanced reentrant creation schedules. Permanent replacement of a queried
stack entry during its creation already has storage-capacity-dependent behavior
in the original code; this change does not promise a new mutation protocol.

Weak-lifetime tests construct AST and scopes manually without installing them
in compiler. currentNode is a separate sentinel and callback/creator retain no
AST/scope references. They retain NodeTraversal while checking popped Scope and
root Node collection across no overflow, first overflow and later heap growth,
for lazy roots, materialized roots, prebuilt scopes and quiet scopes. This
isolates obsolete inline storage rather than mistaking the compiler's ordinary
AST ownership for a retention failure. A separate custom creator deliberately
causes first growth during recursive parent creation and verifies exact
node/scope/parent identities and subsequent grown-stack reuse.

Actual memory or timing benefit has not been measured. The added fullness and
backing-identity check runs on each push; the clear runs only at first overflow.
For usual traversals compiler/currentNode already retain the AST, so unique
reclaimed memory is principally scope/Var/declaration state when such entries
were materialized before first overflow. For traversal objects that die promptly,
there may be little practical retention benefit. Measure this followup separately
from compact capacity selection.
