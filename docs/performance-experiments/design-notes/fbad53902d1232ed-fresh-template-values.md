# Fresh template specialization value vectors

`JSGO_EXPERIMENT_FRESH_TEMPLATE_VALUES=1` is a default-off private stock CLI
inference experiment. It also needs `JSGO_EXPERIMENT_CLI_SOURCE_LOCK=1` for the
existing fresh-input ownership proof. Public inference constructors and the
original Rhino specialization factories retain their existing paths.

Existing default-GC checks heap512 profiles attribute 108.55 MiB / 1.239 million
sampled objects to TemplateTypeMap.CopyFilledWithValues, including 45.04 MiB of
temporary ArrayList growth and 12.51 MiB of the subsequent immutable value copy.
Typed-app's temporary ArrayList growth is 7.51 MiB. These are sampled cumulative
allocations, not RSS. The entire constructor family has 0.14/8.57s sampled CPU
on checks, including required lazy type resolution; it is not all removable
work and is not evidence of a broad predicted speedup.

The alternative fills one exactly sized fresh vector. It copies existing and
incoming values before calling any padding bound, captures the same padding
key slice, then calls DefaultValueType in original key order. DefaultValueType
can call mutable/custom GetBound().IsUnknownType(), so it is not hoisted or
cached. Null-element validation remains after every padding callback and keeps
the original index/message. Constructor registry/key snapshots occur at the
same point before validation. NewTemplateTypeMap still computes a distinct
resolved vector and performs its original replacer/type visits.

The explicit alternative specialization factory preserves the original incoming
value validation, base null check, ProxyObjectType initialization, template
resolution/equality loop, partial replacer creation and resolver close behavior.
Each call creates fresh specialization, template map and replacer identities.
No Array<?> singleton, interning, prior flow result or semantic type cache is
introduced. Only three inference construction sites select this path (rest
parameters, array literals and dynamic import Promise types). Map-based generic
call inference and scope building remain unchanged.

The private default inference factory opens a temporary mode window. Selection
requires exact compiler/registry/roots, stock reverse interpreters, conventions,
reporters, observers and warning guards, no custom passes/debug/tracing/type-AST
injection, and the fresh owned-CLI loader. Public InferAllScopes remains original.
The V2 guard captures an invocation-owned compiler/root/creator/registry/error
manager token. It repeats full admission once after FinishAndFreeze, before the
inference traversal/replay. Each scope then checks those exact identities and
the synchronized current error count. Errors from scope building or an earlier
inference scope therefore disable later-scope fresh construction at the same
boundary as before. Generic and native CFG engines receive the same mode. A
private callback-bearing initial-analysis probe keeps the original full guard
per scope; a test-only oracle also selects that original admission sequence.
The window restores on return or panic. Optional caller-owned admission and
constructor counters/observations are nil in ordinary timed compilation.
Scope/Var layouts do not change; each analysis holds only a mode pointer.

Component tests compare input snapshots, post-padding key mutation, callback
order, padding exceptions before null validation, overfill/null failures,
reentrant construction and output aliasing against the unchanged original.
Specialization tests require distinct identities and equivalent native types,
template values, resolver state and validation behavior. Full inference fixtures
compare typed ASTs, scopes, emitted source and diagnostics with both CFG engines,
require actual dispatch, and cover public/custom fallback, exceptional cleanup
and fresh N=1/N=2 CLI checks parity.

The error-count token contains a bounded, constructor-owned plan of exact
ThreadSafeDelegatingErrorManager wrappers ending in a concrete sorting counter.
Nil/custom/cyclic/deep chains are rejected before acquiring any private lock and
retain the original guard. Each scope tries the existing monitor mutexes in the
original outer-to-inner order. On success the concrete leaf sums its counters,
with the same exclusion against concurrent public Report; no unsynchronized
delegate field read is introduced. On any failed try, every acquired private
mutex is released before calling the original public GetErrorCount. This keeps
same-goroutine public reentry (including an already-held inner wrapper) and
cross-goroutine contention correct. A defer releases acquisitions on panic.
The additive Go-only Monitor.TryLockNonReentrant/UnlockNonReentrant API changes
no existing Lock, Unlock, monitor owner, or hold-count behavior. Its successful
critical section invokes no user callbacks, formatting, resolver, or generator.

This removes repeated convention/observer/warning/input admission and avoids
the ordinary synchronized counter getter's two GoroutineID/runtime.Stack walks
in the uncontended scope loop. It does not remove error checks or synchronization.
The full entry and phase-boundary guards retain their original synchronized
queries. Added state is one bounded plan and token per inference invocation;
Scope/Var/node layouts are unchanged, and no global cache is introduced.

The pre-V2 Rhino selection passed 33 normal and the same 33 race tests. Central
pre-V2 focused testing found a missing class-closing brace in one fixture; that
test-only typo was corrected without reducing coverage. V2 source adds admission
oracle differentials for clean, real scope-building-error, and inference-error
programs under both CFG engines; private callback and direct/public retirement
cases; and normal/private monitor contention, nested reentry, malformed-plan
rejection and panic-unwind coverage. These new V2 tests have not been built or
executed. No attribution of the observed pre-V2 macro regression to guard work
has been measured. Factory CPU attribution is not a bound on added guard cost.
Reduced allocation does not establish lower elapsed time or RSS; same-binary
default-GC correctness gates and measurements remain required.
