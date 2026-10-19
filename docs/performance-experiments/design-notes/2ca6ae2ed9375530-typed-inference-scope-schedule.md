# Private typed inference scope schedule

`JSGO_EXPERIMENT_TYPE_SCOPE_SCHEDULE=1` selects an independent private default
inference experiment. It depends on the initial-fusion branch only for the
shared private factory/eligibility helper; it works with initial fusion off or
on. Public TypeInferencePass.InferAllScopes retains its original second walk.

The required first typed-scope traversal already visits scopes in preorder and
materializes each scope. Record only entries satisfying the original second
callback's `IsCfgRootScope && !IsModuleScope` predicate, with its exact CurrentNode
and canonical TypedScope pointer, plus current SCRIPT for exception context. Complete weak imports, registry resolution and
FinishAndFreeze exactly as before. Replay InferScope sequentially in recorded
order. The original second callback's Visit and ExitScope are empty, so its full
AST traversal is unnecessary under these conditions. No type/property graph is
shared between concurrent writers; FROZEN is not assumed to imply immutability.

Stock TypedScopeCreator construction, weak-import resolution, FinishAndFreeze,
CFG computation and TypeInference do not rewrite AST child links/tokens/names.
They set node/type annotations and update mutable property/type graphs. Scope
classification and preorder therefore remain stable between record and replay,
and each previous InferScope sees the same graph before the next one executes.
The actual node is recorded rather than recomputed from the scope root because
field/default-parameter traversal context can distinguish them.

Custom compiler/reverse interpreter/convention/error observers, custom passes,
native analysis observers, tracing/debug/determinism and custom type reporters
use the original second traversal. Warning guards need a further explicit fence:
Compiler.Report invokes Level synchronously, and a public custom guard can edit
the AST while inference runs. Recursively inspect Compose's existing children,
admit only stock read-only Level implementations, require SuppressDoc's compiler
to be the exact owner, and reject unknown guards without invoking Level or
GetPriority. Validate both compiler and options chains, before recording and
again after freeze. Error-producing resolution can conservatively fall back to
the original second traversal. Public/custom constructors and callbacks retain
their schedules.

The recorder allocates a temporary slice of three typed pointers per eligible CFG
root (24 bytes/entry on amd64, geometric growth). The current SCRIPT preserves
TraverseWithScope's outer exception/context formatting if inference fails; the
normal path creates no traversal per recorded scope. Go's current Platform
interruption check is a constant false, so removed script visits do not remove an
active interruption check. It exists for one inference
invocation, clears before final Var normalization and via defer on panic, and is
not stored on scopes, Vars or the compiler. Private flag windows restore their
previous values on panic. No global state, cross-pass cache or mutation hooks
are added. Extra slice allocation must be measured against the removed walk;
no performance or memory improvement is claimed yet.

Private differential probes require actual replay dispatch and compare the exact
eligible node/scope/parent-root sequence with the original second traversal.
Fixtures compare typed AST, ordered TypedVars and flags, emitted source and
diagnostics across nested functions/default params/class fields/static blocks,
catch scopes, ES modules and goog.loadModule, with initial fusion off/on. Fallback
fixtures cover public calls, compiler/RAI wrappers, native observers, custom and
nested warning guards, and foreign/custom suppression compilers.

Source-only prototype: no build, test, capture or timing has run yet.
