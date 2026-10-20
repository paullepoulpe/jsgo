# Exact singleton-use analysis for private FSIV

`JSGO_EXPERIMENT_FSIV_USE_CARDINALITY=1` opts into a new backward-analysis domain
inside stock fresh CLI FlowSensitiveInlineVariables (FSIV). It requires the
native local-flow topology and fresh CLI source ownership. Public
MaybeReachingVariableUse, its full-set solver/API and custom/debug/step-observer
paths remain unchanged. No workers, extern cache, global state or GC tuning.

FSIV's only GetUses consumer asks whether the set has exactly one CFG node. Its
separate check still counts individual NAME reads within that CFG node. The new
domain stores0 for empty, positive CFG-id+1 for a singleton, and-1 for many.
Union/add agree with the original set operation under this abstraction; a whole
variable kill returns0. There is no partial removal/intersection or cross-variable
transfer. Thus this is an exact quotient for the requested decision.

The original reverse AST scan is compiled once, including destructuring/default
values, optional chains, logical short circuits and ON_EX conditional kills.
All adds within one CFG node add that same node, so each variable's sequence
reduces to kill/gen bits: kill resets earlier gen; a subsequent add sets gen.
The concrete native queue evaluates pointer-free int32 IN/OUT matrices with the
same CFG/priority infrastructure, then answers the original candidate decision.
All unescaped variables are represented; candidate-only projection is deferred.
Same canonical Var aliases share one column. The original allVarsInFn/escaped
lookup semantics, candidate order, side-effect checks and AST rewrites remain.

The admitted graph is immutable during analysis. Fresh stock CLI ownership,
stock warnings/reporters/conventions and no custom observers exclude callback
or injected-type behavior. Preparation that encounters unsupported AST falls
back without publishing state; the original solver then produces its original
failure/order. Var and generated CFG-node hash collisions also fall back rather
than assuming away the reference HAMT's collision-sensitive behavior.

Earlier convergence must not suppress the original20,000-step-per-node failure.
An original IN set grows monotonically from empty, with at most F generated
(Var,CFG-node) facts. Each successor changes at most F times, giving an upper
bound1+outdegree*F on processing a node. Strict backward-priority edges prove
one-pass execution even for large acyclic graphs. If neither proof fits the
original limit, retain the full-set solver. Two dense matrices are limited to
32MiB per function; this is an explicit eligibility bound, not a measured optimal
dense/sparse threshold. All matrices/programs are private to the current function
and cleared before EnterScope returns, including an exception.

Optional private caller-owned counters record attempts/admissions and fallback
reasons, attempted/admitted CFG-node weight, variable counts, generated facts,
compiled actions, actual worklist steps/queries/singletons, cumulative matrix
bytes and maximum matrix bytes. They emit nothing and are absent normally.
Matrix bytes exclude topology, program arrays, allocation rounding and RSS.
Parallel cohort jobs have private pass instances; no shared global counter is
introduced. Actual CLI diagnostics need to aggregate each job's private probe.

`TestFSIVUseCardinality*` compares all IN/OUT abstract states and decisions with
the original full-set lattice, plus real FSIV output/change stamps/diagnostics
across repetitions. Fixtures include same-node duplicate reads, loops, kill-after-
many, ON_EX/finally, switch, optional/default/destructuring, escaped/arguments,
canonical alias columns, hash collision, memory/step/malformed fallback, custom
observers and real fresh CLI SIMPLE admission with positive query counters.

Existing bundle CPU labels put the WHOLE FSIV phase at6.85% CPU, not wall time;
the backward part is smaller. Its two full-use HAMT allocation leaves sum to
126.36MB of5.071GB cumulative allocation, before child slices/wrappers. Matrices
and preparation still cost work. No speedup, memory reduction or broad parity
is claimed. Source has been formatted/checked only; central normal/race/oracle
and matched time/RSS measurements are required.

## Admission correction after initial validation

The first central race run passed every full-set state/decision oracle and the
isolated private-pass tests, but the complete SIMPLE CLI admission assertion
failed. Real-job diagnostics also reported zero admissions (typed1352,
bundle16494 and simple3697 attempts), all rejected by the stock fence. Those
results exercise the original solver and cannot measure the quotient.

The inherited scope-allocation guard rejected any loaded runtime-library table.
`ConvertTypesToColors.Process` initializes this table before normalization in
both checked and unchecked ordinary optimization pipelines. The quotient now
uses its own proof: stock callbacks/reporters/warnings, no injected typed AST
filesystem, and the original private fresh-CLI source lifetime remain required.
An already-loaded runtime-library table is admitted because the transfer does
not read or invoke it. No global allocator guard was weakened. Regression tests
include an unused throwing runtime supplier, preserved injected-filesystem/custom
warning fallbacks, and an explicit actual CLI post-typesToColors table assertion.
Central validation and measured activation must be repeated for this correction.
