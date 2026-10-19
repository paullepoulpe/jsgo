# Invocation-owned inline argument facts

`JSGO_EXPERIMENT_INLINE_ARGUMENT_FACTS=1` enables this opt-in prototype. The
FunctionArgumentInjector constructor snapshots the setting; the original
recomputation path remains the default and test oracle.

GatherCallArgumentsNeedingTemps previously walked the complete function body
separately for every argument's name-reference count. The experiment collects
all requested counts during its existing local-name discovery pass, reducing
that part from O(arguments × body nodes) to O(body nodes + arguments). It then
runs the original side-effect collector unchanged and consumes the counts in
the original argument loop.

The two fused predicates have different boundaries: local-name discovery stops
at nested functions, but name-reference counts include their declaration names,
parameters and bodies. The recursive walker preserves that distinction using a
local-scope flag. It preserves local-set insertion order and counts NAME tokens
exactly, including the original distinction between NAME "this" and token THIS.

The prunable side-effect traversal is deliberately separate. Its original
Apply/Visit order, stopping rules, loop state, side-effect checks and parameter-set
ordering remain intact. Collecting counts only on the nodes that survive that
traversal would miss later references once all relevant side effects were found.

Facts belong to one invocation and are discarded when it returns. There is no
AST-pointer cache, cross-call reuse, fixed-point reuse or ChangeTracker-based
invalidation. Body mutations before a subsequent call are observed by its new
facts. During the supported call, stock compiler reads and known convention
methods do not mutate the body. Custom compiler/convention implementations keep
per-argument recomputation; raw type optimization without optimization colors
also retains the original path. A regression convention mutates the body while
processing an earlier argument and verifies that this fallback still recomputes
the later argument's count.

Differential tests compare every count, local-name order, both side-effect sets,
and final temporary-name order with the original implementation. Cases include
early stopping, loops, nested functions, classes/field closures, destructuring,
catches, rest/excess arguments, optional chains, explicit receivers and Unicode
names. Existing FunctionArgumentInjector and InlineFunctions tests under the
flag, race checks and ordinary CLI output gates remain required.

No speedup is claimed before measurement. The eligible parameter count, map
allocation/lookup overhead, larger nested-function scans and live facts can
offset saved walks, especially for small callees or already-aliased parameters.
Measure elapsed CLI time and peak RSS; the whole InlineFunctions profile fraction
is only an upper bound on the work this narrower change can affect.
