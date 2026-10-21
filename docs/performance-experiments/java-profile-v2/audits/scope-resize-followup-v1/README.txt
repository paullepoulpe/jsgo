Scope declaration and map-resize follow-up

Source-only outcome: no stronger low-complexity candidate justified by the
existing evidence. No production source was changed, and no Java, build,
benchmark, archive operation, or full-jar hash was run for this audit.

Reassessment of the closed scope-declare-once candidate

The patch removes one hasOwnSlot(name) call for each new declaration reaching
AbstractScope.declareInternal. That call repeats an explicit map lookup and,
on a miss, the implicit-name predicate. Existing-name replacements and all
canDeclare calls retain their original behavior. The source proof remains valid
for actual Scope/TypedScope implementations, including their lazy implicit vars.

The 8.9149% createScope, 8.3746% populate, and 2.9716% declareVar inclusive CPU
shares include substantial other work. They do not establish the cost of the
removed query. The candidate removes no Var, scope, or explicit variable-map
allocation. It is worth evaluating as a small redundant-work removal, but this
review found no new reason to expect a large elapsed-time change.

Existing resize attribution

Corrected profile-stock-bigquery-steady-v2 has 57 HashMap.resize leaf samples
among 2,221 CPU samples (2.5664%). Its hot-callers.tsv shows 20 caller prefixes
covering 48 samples; 9 samples from 9 distinct prefixes are unshown. See the
compact resize-callers.tsv and attribution.json in this directory.

Largest identified callers:
  ReferenceCollector.addReference: 12 samples, 0.5403% of all CPU samples.
  AbstractScope.declareInternal via Scope.declare: 9 samples, 0.4052%.
  TypedScopeCreator.createScope: 3 samples, 0.1351%.
  Every other shown prefix: at most 2 samples, 0.0900% each.

These are identified shares, not exhaustive per-caller ceilings, because nine
prefixes are unshown. Function-only scope sizing would affect only some scope
declarations, and resize also includes initial table allocation. It cannot be
credited with the aggregate 2.5664% resize share.

Why no map-sizing patch

ReferenceCollector's LinkedHashMap starts empty and spans variables from nested
scopes, with optional filtering and special hoisted-function traversal order.
getAllSymbols exposes its insertion order. Before traversal there is no cheap
complete count of the variables that will actually receive references. Global
scope variable count omits locals; AST size is an overestimate requiring another
walk; a fixed larger capacity is a workload guess. Replacing this map with an
unordered identity map would break its encounter-order contract.

Function-scope parameter count is available, but obtaining it requires a child
walk, gives only a lower bound for destructuring, and excludes the bulk of block
scope declarations. A bounded hint computed at first map allocation can retain
laziness, yet the currently attributed opportunity is smaller than the broader
resize total. No evidence here supports prioritizing that complexity over the
already prepared SourcePosition reuse candidate.

Other examined leads

The FUNCTION arm of ScopeScanner.populate could skip declareLHS(scope, args)
when the already-validated PARAM_LIST has no children. The old path creates a
capturing consumer, repeats validation, and visits an empty list; it invokes no
declaration or redeclaration callbacks. Bleeding-name declaration would continue
unchanged. This is a valid tiny source lead, but its empty-parameter fraction and
surviving allocation cost are unknown. The 26 declareLHS leaf samples and 91
inclusive samples include real declarations. The 27 inclusive lambda-consumer
samples cannot support an empty-list guard, since empty lists never invoke that
consumer. No patch was prepared on this evidence.

Generalizing declareLHS's existing one-NAME fast path to multiple declarations
must preserve sibling traversal when custom redeclaration callbacks mutate ASTs;
an up-front classification walk also adds work. Broader scope memoization needs
AST invalidation and stable parent scope identities; NodeTraversal already creates
its per-traversal scopes lazily. These are not clear small improvements here.

Recommended next screen

Use the already closed source-position-last candidate rather than adding another
weak map guard. Its target class accounts for 5.3498% sampled allocation weight,
but qualifying adjacent reuse frequency and net CPU effect remain unknown. Its
semantic probes and compiler byte gates must pass before timing. This is a test
priority recommendation, not a claim that it will improve performance.
