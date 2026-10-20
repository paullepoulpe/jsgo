# Callee facts within one read-only reference-finding phase

`JSGO_EXPERIMENT_READONLY_CALLEE_FACTS=1` selects this default-off experiment in
the stock InlineFunctions pass factory. It is independent of the existing
invocation-local `JSGO_EXPERIMENT_INLINE_ARGUMENT_FACTS` toggle.

The existing fused walk still computes the same callee's local names, reference
counts and parameter/side-effect facts for each callsite. During
FindCandidatesReferences, candidate bodies have not yet been rewritten: this
phase records call references and pass decisions. Expression decomposition,
safe-function cloning and actual inlining happen later. This narrower interval
can therefore reuse immutable callee facts without a cross-pass stamp cache.

An opt-in phase sidecar keys each lazily created summary by the exact
FunctionState pointer, preserving ordinary FunctionState allocation size. The injector sees it only
while checking that state's callsite. Cache keys compare every parameter name
and its order exactly, including `this` and synthetic excess-argument names.
Those names and their count affect the original side-effect collector's early
stopping rule and ordered sets. Four shapes per callee are retained; additional
shapes use invocation-owned facts without retention. DIRECT-to-BLOCK retries
can reuse an exact shape too.

Only callee facts are shared. Argument effects, mutable-state checks, export
decisions, temporary-name decisions and unique-name supplier calls execute at
their original callsites and in their original order. The immutable result is
never exposed to callers for mutation. The discovery phase's deferred cleanup
clears the entire per-state map and the injector context before returning,
including an internal exception. Nothing survives into cost trimming, conflict
resolution, decomposition, inlining or another discovery phase.

Public NewInlineFunctions accepts an arbitrary name-supplier callback. Such a
callback can mutate the callee during construction of an argument map. That
constructor therefore retains ordinary recomputation even with the flag set.
A private stock factory supplies the exact Compiler's own unique-name closure.
Custom compiler implementations, coding conventions, observers/error handlers,
debugging/tracing and raw-type local optimization also retain recomputation.
No Go function identity comparison or global AST-pointer cache is introduced.

Focused tests compare complete pass output, AST, diagnostics, change stamps and
timeline order against the public recomputing pass. They assert actual summary
hits/misses, DIRECT-to-BLOCK reuse, bounded exact parameter-shape handling,
cleanup before a subsequently mutated callee, and custom-supplier mutation
schedule/fallback. Optional caller-owned test probes are absent in normal CLI
timings. Existing public-constructor InlineFunctions fixtures remain an oracle;
they do not by themselves prove this private production factory executed.

The ceiling is the repeated callee walks and their temporary collection
allocation during this one phase. Caller work, function cloning, actual AST
rewrites and later phases are unchanged. Many block callees never require
GatherCallArgumentsNeedingTemps, and single-call functions provide no reuse.
Retained summaries can raise peak memory until the phase ends. Actual macro
timing, peak RSS, and eligible-hit coverage must establish whether this is useful;
no performance improvement is claimed by the source prototype.
