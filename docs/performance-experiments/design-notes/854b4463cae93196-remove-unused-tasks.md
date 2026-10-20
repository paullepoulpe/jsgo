# Invocation-owned RemoveUnusedCode tasks

`JSGO_EXPERIMENT_REMOVE_UNUSED_TASKS=1` selects an experimental continuation
representation only for a fresh private default-pass invocation on the current
compiler-owned root. It additionally requires the existing private
`CLI_SOURCE_LOCK` ownership proof from the stock fresh CLI runner. It remains
off by default. There are no measured performance
or memory improvements yet; source inspection attributes approximately 10.3% of
one bundle CPU profile to the complete RemoveUnusedCode pass, not these
allocations alone.

The original representation allocates a 24-byte ArrayList header for each
RemovableBuilder, a pointer backing array for nonempty lists, and a 24-byte
Continuation object for each deferred subtree. Continuation is already concrete;
this experiment removes ownership/allocation layers around it. Native task rows
hold real GC-visible Node and Scope pointers plus a list link (24 bytes on amd64).
An 8-byte list range is shared by all removables built from that builder. Its
uint32 identity occupies existing builder/base padding: the builder remains
32 bytes and RemovableBase remains 56 bytes. VarInfo and Removable objects,
canonical identity relationships and their original operations remain intact.

A concrete FIFO contains 16-byte Node/Scope values. It classifies the node as a
function at application time, exactly as Continuation.Apply does. Transferring
a list zeroes task-row pointers; consuming the FIFO zeroes its pointer cells
before traversal can enqueue more work. Freed row IDs are reused within this
invocation, bounding task storage by peak pending tasks rather than cumulative
applied tasks. Shared list ranges clear in place, preserving builder reuse,
multi-removable aliasing, separate applied bits and FIFO order. Task/list chunks
hold 64 rows, respectively 1,536/512 bytes plus slice-directory overhead. This
introduces rounding costs for small phases; the full compiler must decide
whether reduced allocation count outweighs them. The FIFO doubles in capacity,
matching the original queue growth policy, but stores values instead of pointers.

Public Build/NewRemoveUnusedCode_RemovableBuilder/AddContinuation and direct
TraverseAndRemoveUnusedReferences calls retain the original representation.
Only Process on the private fresh builtin pass may open the native window;
custom compilers, conventions, redeclaration handlers, change observers, error
handlers, warning guards, logging/tracing and custom pass schedules fall back.
The private CLI loader witness excludes public/default-factory manual compiler
invocations, saved/typed input and custom static source implementations; stock
callbacks cannot inject a public legacy continuation queue during native drain. PhaseOptimizer
constructs and invokes each NamedPass instance immediately and retains no pass
instance (phase_optimizer.go); stock callbacks cannot capture these private
builders, VarInfos or removables. The arena is released by a Process defer on
success or panic. There are no global pools, cross-phase caches, Node/Scope owner
fields, integer AST handles or unsafe pointer conversions.

The uint32 list/task index limit is checked before incrementing. Reaching it
would require at least 32 GiB of list rows or 96 GiB of simultaneously pending
task rows, respectively, before directory/object overhead; this experimental
representation reports an invariant failure rather than silently wrapping.

Source fixtures compare complete real-pass AST/output/diagnostics/change stamps
and VarInfo sharing against the original path, including long transitive function
dependencies, dead functions, classes/properties, accessor definitions, closures,
loops and catch scopes. Separate fixtures exercise shared-list clear semantics,
FIFO append during draining, growth/wrap, immediate pointer release/recycling,
public identity/fallback and exceptional phase cleanup. A private optional
counter object proves actual native construction/enqueue/application in untimed
tests; ordinary compilation allocates no counter object. Isolated compiler
fixtures explicitly install the package-private CLI proof; a separate fixture
removes it and checks original dispatch.

Validation status: source formatted and whitespace checked only. The candidate
requires compiler tests, race checks and whole-CLI timing/RSS measurements before
retention. A storage microbenchmark cannot establish an end-to-end gain.
