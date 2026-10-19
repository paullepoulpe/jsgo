# Serial FSIV read-only scope window

Base: `ef9657d7` (native80). This is a separate default-off experiment selected by
`JSGO_EXPERIMENT_FSIV_SCOPE_SHARING=1`. It does not add Scope or Var fields, a
compiler-wide cache, change stamps, declaration snapshots or mutation hooks.

Serial FSIV already has a function-body scope when it enters its analysis. It
then calls `NodeUtil_GetAllVarsDeclaredInFunction` and
`DataFlowAnalysis_ComputeEscaped` consecutively before building dataflow results
or gathering/replacing candidates. The first helper reads declarations into its
own maps/list. The second reads names and stock coding-convention export answers
into its escaped set. Their callbacks do not edit AST nodes or declarations.
Lazy implicit-variable lookup keeps its normal existing function-scope owner.
The default redeclaration handler is a no-op; custom handlers retain fresh scans.

The private creator seeds the exact existing body Scope and matches root and
parent pointers exactly. It records lexical scopes encountered in the first
helper and looks them up in the second. New misses in the second helper, such as
nested/default-parameter functions skipped by the first, delegate to the original
creator without recording. The private creator is released with defer before
returning, including panic exits. Its delegate, body pointer and map are cleared
before dataflow/candidate/rewrite work. Only the ordinary returned Var/scope facts
retain their normal analysis lifetime; no creator is installed on the outer
traversal or stored on a compiler.

Eligibility requires the exact stock Compiler in normalized unobfuscated state,
stock default/Closure/Google convention chains, the stock observer/error-manager
guard, and an unchanged default syntactic creator/redeclaration handler. Custom
compiler facades, conventions, declaration handlers, error/change observers,
debug/tracing modes and public helper calls keep their original schedule.
Parallel FSIV cohorts are not opted in. This guard is checked for each serial
pass; creator identity is checked for each window.

A private observer receives numeric counters only after release: requests/hits
in both helper phases, seed hits, lexical scopes recorded, uncached second-phase
misses, completion and release status. It provides actual-dispatch evidence in
tests without tracing or adding a user-facing reporting mode.
Counter storage is allocated only when that private observer is installed (or a
unit test explicitly requests instrumentation); ordinary experiment runs keep a
nil stats pointer and perform no counter updates. The amd64 window has four
pointers plus a recording boolean (40-byte layout, 48-byte allocation class).

Tests cover declaration order/index/input/name-node/root relationships, canonical
existing body Vars, escaped equality/encounter order, nested functions, defaults,
arrows/arguments, lexical/loop/catch/class-field scopes, exact parent identity,
non-recording second-phase misses, mutations between windows, exception cleanup,
stock/custom guards, public freshness and custom callback schedules. Actual FSIV
dispatch tests compare output, diagnostics and ordered change stamps with the
original path and require nonzero sharing hits.

Source inspection establishes the read-only boundary for these stock helpers;
it does not establish a performance gain. A window object, optional map and hit
accounting have costs. Root currently reserves builds and benchmarks; at source
freeze, tests, CLI Java-output gates and ordinary time/peak-RSS/total-allocation
measurements are pending. Reject the experiment if realistic weighted savings
do not exceed these costs.
