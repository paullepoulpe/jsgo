# Conservative whole-function Coalesce cohorts

This isolated prototype starts at ef9657d7. The opt-in internal environment flag
is JSGO_EXPERIMENT_PARALLEL_COALESCE=1; the only worker budget is the existing
num_parallel_threads option. It never modifies GOMAXPROCS or global runtime
settings. The original serial pass is the default and remains the oracle.

The bundle profile attributes roughly 3.8% CPU to CoalesceVariableNames. Even
perfect two-CPU scaling of that whole fraction would save only about 1.9% CPU,
before coordination, preparation and imbalance. Actual covered work is smaller.
This can compose with other work, but is not a route to beating Java by itself.
No new macro timing or memory result is claimed by this source prototype.

## Strong ownership boundary

Only dispatch when the coordinator has no active Coalesce coloring, including
all ancestors. Then its own NAME visits cannot rewrite local bindings. Eligible
function parents and their descendants stay together and run in original serial
order within one job. This avoids relying on speculative parent rewrites or on
an assumption that all captured names remain stable after arbitrary changes.

The collector chooses non-overlapping whole-function subtrees. A function whose
explicit root and immediate body scopes already have 100 or more variables
cannot be eligible; leave it for the coordinator to enter serially and discover
smaller children. Parameters live in the root scope and body declarations in a
separate block scope; their sum is a lower bound on the original pass's complete
per-function list, which also includes nested lexical scopes. This
can split large UMD/library wrappers instead of giving a whole bundle to one
worker. The test probe records actual job dispatch and wraps real job execution;
a successful entry into the experimental path does not itself prove parallel
work happened.

Each job owns its pass, memoized descendant scopes, CFG, liveness lattices,
interference graph, graph-coloring state and copied AstFactory header. Canonical
ancestor scopes are read only. All their implicit arguments/this/super slots are
materialized serially before workers start. Var's lazy name hash uses the
existing atomic implementation. There are no global fresh-name suppliers: this
pass chooses existing local names with the original variable-order tie breaker.

## Publication and shared-state guards

Coalesce mutates during postorder NAME visits. FSIV's scope-entry publication
pattern cannot simply be copied: it would move effects before surrounding serial
work. Workers buffer diagnostics and resolved change-scope events. After the
join, jobs remain pending until the coordinator reaches each original function
root. It then publishes that job's events and skips its already processed
subtree. The coordinator's remaining traversal and final lifecycle transition
stay in their original order.

The compiler facade leaves unaudited methods nil, shares only immutable settings
and input maps, serializes exception source-line formatting, and asserts every
change belongs to the job subtree. Original report handlers and diagnostics run
only in the coordinator. Work execution may be sorted by subtree weight; report
publication follows original traversal order. At most N analyses/rewrite jobs
run concurrently using the existing bounded pool. Job metadata and buffered
effects remain until their original roots are encountered.

Conservative fallback applies to custom compiler implementations/conventions,
custom observer/error hooks, debug/tracing/determinism modes, custom scope
creators/redeclaration handlers, pseudo names, and raw-JSType AstFactory mode.
Concrete convention proxy types must also have self pointing to themselves;
copied proxies are rejected, matching the later 59ad22c1 safety fence without
changing this branch's FSIV implementation. NONE/COLOR factory paths only create
private void 0 nodes with fixed native colors; they do not resolve shared raw
types or use lazy Color hashing. Raw-JSType mode remains serial.

## Proof and validation plan

Source fixtures compare the original serial pass against N=2/N=4 for actual output,
ASTs, change-scope order/stamps, warnings and lifecycle stage. They include
oversized wrappers, default-parameter functions, sibling arrow captures of
canonical arguments, catch/finally, destructuring, modules, class fields/methods,
IE8's two-parameter restriction and the color-producing undefined reset path.
An eligible-parent fixture must dispatch no descendant jobs. Guard fixtures
cover N=1, pseudo names, raw types and copied convention proxies.

The overlap probe blocks real jobs before running them, proves the existing
N=2 limit, then executes the real Coalesce analysis/rewrite. Publication checks
verify worker changes remain buffered after join and compare the final ordered
snapshot with the serial oracle. An oversized wrapper with default/body
functions followed by an eligible sibling tests deferred publication around
serial coordinator work.

A joint fixture runs FSIV followed by Coalesce on the same AST, matching their
relative order in DefaultPassConfig. It requires actual worker rewrites in both
phases at N=2/N=4 and compares the serial AST, output, ordered reports/stamps and
lifecycle after each phase. For the parser-allocation race configuration, set
JSGO_EXPERIMENT_NATIVE_PARSER=1 and
JSGO_EXPERIMENT_NATIVE_PARSER_ALLOCATOR=1 before starting the test binary; the
parser reads these flags during package initialization. This additionally checks
the same passes on independently mutated nodes within parse-owned chunks.

No builds/tests were run during the initial source-only reservation. Before
macro consideration: compile and run focused Coalesce/cohort tests, run the
existing Coalesce oracle suite, exercise nested/shared-ancestor fixtures under
the race detector, and run full SIMPLE/ADVANCED CLI parity at N=1/N=2. Then record
untimed dispatch coverage and matched baseline/experimental N=1/N=2 time and peak
RSS. The small Amdahl ceiling and added discovery/scope overhead can erase the
benefit; retain only if the complete CLI workload supports it. Partial AST
identity after an internal worker failure is outside this CLI prototype; the
earliest original-root failure/report prefix is preserved, as in FSIV.
