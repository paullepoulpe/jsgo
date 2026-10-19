# Declaration recipes in an IV-to-RUC optimization window

`JSGO_EXPERIMENT_SCOPE_RECIPES=1` enables a default-off, ordinary CLI
experiment. It records syntactic declaration operations while a private builtin
InlineVariables (IV) constructs scopes, and can replay them into fresh scopes
in the next private builtin RemoveUnusedCode (RUC). It caches no parsed externs,
resolved Scope, Var, CompilerInput, declaration map, or binding lookup result.

The measured broad subsystem attribution was about 8.9% CPU and 517 MiB sampled
allocation below syntactic declaration scans. This prototype only avoids the
statement-selection walk for eligible reused scopes; declaration expansion,
lookup, legality checks, Scope/Var/map allocation, globals and parameters remain.
The broad numbers are overlapping ceilings, not predicted savings. Recording
adds operation storage and a map, and can lose on edited or small scopes.

## Actual scheduling and lifetime

The stock SIMPLE configuration used by the bundle workload has these adjacent
sequences in DefaultPassConfig:

* Early loop: IV, InlineObjectLiterals, RUC, peephole optimizations.
* Main loop: InlineFunctions, IV, DeadAssignmentsElimination,
  InlineObjectLiterals, RUC, peephole optimizations.

The lease opens at IV and permits only the exact DefaultPassConfig factories
for DeadAssignmentsElimination and InlineObjectLiterals in between. These two
stock rewrites report changed enclosing owners; their custom public constructors
or suppliers do not qualify. The next exact builtin RUC consumes and retires
the lease. A different pass, skipped loop pass, new IV, optimizer/loop exit,
lifecycle change, or exception retires the current window. Nothing persists
between compilations or across an unrelated pass. Pass construction, BeforePass,
AfterPass and execution remain at their original scheduling points.

Ownership starts before Init in fresh exact stock CommandLineRunner.DoRun and
ends by defer on success, early return, or failure. Initialized/public compilers,
custom runners, typed-AST restore/save, custom passes/observers/error managers/
conventions, debug logging and validity checking use the original scope scan.
Public IV/RUC/scope creator constructors retain their original behavior. The
extra state is one private Compiler pointer and a creator flag in existing
padding; Scope and Var layouts are unchanged.

## Validity and fresh declarations

Eligible roots are function-owned BLOCK, FOR/FOR_IN/FOR_OF/FOR_AWAIT_OF and
SWITCH_BODY scopes. Function parameter scopes, class roots, field scopes,
SCRIPT/ROOT, module bodies (including bundled goog.module blocks), global and
legacy namespace scans stay generic. A class static BLOCK can qualify when its
nearest change owner is a function; its original hoist/block role is recorded.

Each entry is keyed by lexical root identity, nearest FUNCTION change owner and
its stamp, the owner's immediate parent, input ID, immediate lexical-root parent,
current parent-scope root identity, root token and block/synthetic roles.
Moved/cloned roots or changed owner/context miss. Each hit constructs a new
Scope with the current parent, initializes an original scanner with the current
input ID, and replays original DeclareLHS/DeclareVar operations. Current input
resolution, names, duplicate order, implicit slots, bleeding names, parameter
shadowing, declaration indices and linked-map semantics remain original.
Changing the CompilerInput associated with the same input ID is safe because
no resolved input is cached.

A nested function body's edit need not invalidate its outer owner: scanning the
outer body cuts at the nested FUNCTION. Function/class declaration operations
read their current name and expression/declaration role at replay, including
initially empty names. This covers child-name edits whose reported change stamp
does not propagate to the outer owner. Actual enclosing statement/declaration
changes rely on the audited stock passes' change reports; arbitrary public AST
mutation cannot enter this private lifecycle. The recorder uses the generic
scanner's exact statement cuts, catch checks, hoist/block target selection and
order. Unsupported nested input/module boundaries use the generic scanner and
do not publish a recipe.

## Evidence and measurement

Private aggregate counters record opened/consumed windows, recorded/hit/missed
scopes, owner/context invalidations, unsupported/ineligible scopes, original
statement-selection visits, saved visits and replayed declaration operations.
A private optional release observer can receive the pointer-free aggregate
snapshot; normal compilation has no output callback or diagnostic output.
`PeakRecipeBytesEstimate` counts recipe headers and allocated operation slice
capacities. It excludes map buckets, allocator rounding and AST nodes retained
by recipe pointers. It is not heap usage or RSS; matched CLI measurements must
report actual peak RSS separately. Window cleanup drops every recipe pointer.

`TestScopeRecipes*` covers ordered bindings/lookup/node/input relationships,
new Scope/Var identities, duplicate/hoist/catch/destructuring/arguments/module/
extern cases, owner edits, live nested names, nested-body-only edits, moved and
cloned roots, current input resolution, actual SIMPLE pass lists and full
compiler output/AST/change epochs/diagnostics, custom fallback, skipped phases,
exceptions and lifecycle cleanup. A weak-reference GC fixture checks that live
recipes do not retain the discarded Scope/Var graph. Positive replay counters
are required, so fallback-only tests cannot establish coverage.

This source prototype has not yet been built, run or timed. Normal/race tests,
Java output oracle checks, actual workload activation counts, retained-memory
measurement and a matched default-GC CLI screen are required before promotion.
