# Audited untyped scope reuse experiment

`JSGO_EXPERIMENT_SCOPE_REUSE=1` selects a private creator only for the builtin
InlineVariables and serial FlowSensitiveInlineVariables traversals. Public
SyntacticScopeCreator constructors preserve fresh Scope/Var identity. Custom
compiler implementations, error handlers, debug logging, redeclaration handlers,
module-body roots and scans involving implicit goog namespaces retain fresh scopes.
Parallel FSIV jobs retain their existing uncached creators.

Each stock compiler owns its cache. Keys contain the AST root and the exact parent
Scope pointer. A changed ChangeTracker stamp retires the whole cache. Direct
Declare, Undeclare, ClearVarsInternal and implicit namespace operations retire it
too; initialization resets ownership for the next compilation.

A stamp is insufficient to prove unchanged declarations. Each hit therefore runs
the original ScopeScanner with a private validation sink and compares every
declaration's node, spelling, input identity, target scope and encounter order.
It also checks root/parent context and constructor invariants. A mismatch retires
the whole cache and performs an ordinary fresh scan. Provide encounters always
reject caching, including a provide inserted without a change report. This avoids
retaining stale implicit namespace strength or accumulated provide declarations.

Cached Scope/Var identity is confined to the two audited internal consumers; no
public ReferenceCollector behavior callback receives these scopes. Declaration
nodes remain live, so Var metadata continues to follow the AST. This is not a
result cache and does not alter AST ownership or thread budgets.

The validation scan retains AST-walk costs. The intended saving is fresh Scope,
Var and ordered-map allocation and declaration hashing. Costs include one owner
pointer per Scope, declaration-event snapshots and the compiler cache itself.
No improvement is claimed until correctness gates, hit counts, retained memory
and ordinary CLI measurements establish a benefit.
