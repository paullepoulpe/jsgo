# Intrinsic leaf traversal prototype

This default-off experiment removes a recursive `TraverseBranch` call, token
dispatch, virtual preorder callback, and scope classification for childless
intrinsic leaves. Set `JSGO_EXPERIMENT_TRAVERSAL_LEAVES=1` before constructing a
traversal. The normal callback `Visit` still runs. It introduces no public option,
parallel work, AST cache, alternate graph, or runtime thread/GC setting.

The motivating profiles have persistent `TraverseBranch` and sibling traversal
cost across many passes. `GetNext` itself is already a field read; this prototype
does not remove necessary sibling reads. There is no measured speedup or memory
claim yet. Added dispatch/code size and a larger traversal header can outweigh
the saved frames, and the default-off branch also needs a same-binary control.

Eligibility uses an exact concrete callback whitelist rather than a base marker.
The selected abstract scoped/postorder bases return true directly and contain no
virtual `self` pointer. Other admitted preorder implementations have effects
only on non-leaf tokens: catch bookkeeping on BLOCK, function pruning on
FUNCTION, and script/function itinerary recording. The changed-scope wrapper is
allowed only with a recursively admitted delegate and a FUNCTION/SCRIPT target.
No custom accessor is called to determine eligibility. Copied callback outer
fields can change `Visit` behavior; that method executes on the original object.

Unknown callbacks and subclasses always retain the original loop. In particular,
access-control checks can establish context on leaf class-extends targets, type
checks can start logging on annotated leaves, and scope builders have meaningful
preorder work. They are excluded. Scope/special-context tokens are excluded even
when childless. Every decision checks the live token and children, so a NAME with
an initializer or a malformed STRINGLIT with children takes the original route.

Each fast child visit caches its next sibling before invoking `Visit` and sets
`currentNode` at the original point. No other traversal context is changed.
Postorder callbacks may add children, change tokens, detach/replace the current
node, or insert/detach siblings. The fast loop preserves the original handling,
including visiting a cached next sibling that has since been detached, using the
original parent argument. Function/default/class/module traversal remains in the
existing special handlers, including execution-order and lazy scope behavior.
The duplicated generic loop avoids adding another recursive helper frame.

Source-only differential tests cover actual positive leaf dispatch with a
caller-owned nil-by-default counter, current-node/lazy-scope/input/chunk contexts,
both default/destructuring execution orders, postorder mutations, custom
preorder overrides/pruning, exception source contexts, early/late fixed-point
peephole changes and feature/stamp order, and complete SIMPLE/ADVANCED compiler
output/source maps/AST/diagnostics. No build, tests, or benchmarks have been run
by the author; the central validation batch must establish correctness before
any performance result is used.
