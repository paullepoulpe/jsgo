# Latest-main inference: next design audit

Source-only, 2026-10-19. Inspected ff06ac15/current clean combined a3a2f8f3,
`TypeInferencePass`, `TypeInference`, `LinkedFlowScope`, `TypedScopeCreator`,
`TypedScope`, `ModuleImportResolver`, and generic data flow. Used already-decoded
`latest-main-profiles-v2/hotspots.json` and `allocation-owners.json`; no profile
capture/decoding, build, test or timing. No active source tree changed.

## Actual target, and limits

The fresh checks-only profile contains 10.67 sampled CPU seconds and
1,609,979,406 cumulative allocated bytes. Inclusive rows overlap:

| Current function/family | CPU / whole | Allocation bytes / whole |
| --- | ---: | ---: |
| InferAllScopes | 2.72 s /25.49% | 687,230,193 /42.69% |
| Initial scope-building phase | 1.23 s /11.53% | 263,375,704 /16.36% |
| TypedScopeCreator.CreateScopeInternal | 1.09 s /10.22% | 258,633,420 /16.06% |
| InferScope | 1.43 s /13.40% | 413,475,655 /25.68% |
| Generic data-flow Analyze, FlowScope shape | 1.22 s /11.43% | 297,935,207 /18.51% |
| Generic Flow, FlowScope shape | .99 s /9.28% | 218,688,304 /13.58% |
| TypeInference.FlowThrough | .83 s /7.78% | not separately in retained top byte rows |
| TypeInference.Traverse | .84 s /7.87% | not separately in retained top byte rows |
| Flow-map family, overlapping | .19 s /1.78% | 33,557,052 /2.08% |
| Template-type family, overlapping | .32 s /3.00% | 135,840,360 /8.44% |

These are CPU and cumulative sampled allocations, not elapsed time or RSS.
The whole inference row is not the removable cost of a solver or map rewrite.
Typed-app has only11.21% CPU/17.98% allocation below InferAllScopes, further
limiting a general end-to-end claim.

## Best new bounded experiment: resolved-binding flow operations

Fuse repeated lexical binding resolution **within one read or assignment**, then
pass the current canonical TypedVar/owner directly to private stock flow helpers.
Do not cache an inferred type, whole expression result or binding across visits.
Keep exactly the original worklist, transfers, branch order, diagnostics and
public/custom FlowScope path. This is a representation seam between lexical
bindings and immutable inferred-slot maps, not another lattice or solver.

Concrete repeated work in current source:

* `UpdateScopeForAssignmentWithUpdateNode` on NAME resolves a declared TypedVar,
  then `RedeclareSimpleVar` resolves the same declaration again, then
  `LinkedFlowScope.InferSlotType` resolves the root name again merely to find the
  overlay's owner. Common assignments therefore perform at least three lexical
  lookups before the persistent overlay update. All three are to the same current
  syntactic declaration scope, with no AST or binding-table write between them.
* `TraverseName` calls `scope.GetSlot`; concrete `LinkedFlowScope.GetSlot` first
  resolves the declared TypedVar. If an owner overlay exists but contains no
  inferred slot for that name, `OverlayScope.GetSlot` calls the typed scope lookup
  again. `TraverseName` then resolves the declared var again for `IsUnflowable`.
* `IsUnflowable` itself needs only the current canonical Var's local/escaped
  flags and owner/container scopes. Repeating name resolution supplies no new
  information in this exact read-only interval.

A minimal isolated implementation would add private helpers that return both
`(current slot, current declared Var)` and accept a resolved owner on an inferred
slot update. The update still reads the current persistent scopes map, creates
fresh OverlaySlot/OverlayScope/LinkedFlowScope versions when required, and uses
all original JSType operations. Original generic methods remain unchanged as
oracles. No new Scope/TypedVar fields, global dictionary, persistent vector,
retained AST index or dense node×symbol product is needed. Operand facts live on
the stack only until the current operation returns.

Restrict the first factor to an exact private TypeInference instance from the
fresh stock CLI lease, the standard linked FlowScope, stock reverse interpreter
and convention/reporting/observer path. The resolved typed declaration must be
non-nil, explicit and have matching name/owner. Check the selected overlay's
scope pointer matches the resolved owner; otherwise use the original operation.
Leave unknown/implicit names, qualified names and all custom/public paths on the
original route. Read `TypedVar.GetType()` at every original semantic use: cache
neither types nor inferred/declared decisions past an operation that can change
those values. Reverse-interpreter APIs continue to receive original FlowScopes.

The first testable factor is the assignment chain: reuse the already-resolved
Var from UpdateScopeForAssignment through the stock IsUnflowable decision and
InferSlotType owner selection. A second independent factor can fuse NAME reads.
This reduces three same-context lexical lookups to one on eligible assignments;
it does **not** eliminate joins, type construction, HAMT path copying or the
semantic type operation itself. Source provides no trustworthy whole-program
speed forecast. Even eliminating all FlowThrough work would cap direct sampled
CPU at7.78%, and this target is only a subset of that. Existing top rows do not
isolate the repeated lexical lookup union.

Validation should compare JS, every node's JSType, canonical TypedVar identities,
per-CFG IN/OUT/edge states, full ordered diagnostics, step/queue sequence and
original branch outcomes. Cover shadowed block bindings, declared/typeless const,
catch scopes, escaped closures, arguments/this/super, forward named references,
ES module import replacement, global qualified names, bottom scopes, and retained
older flow versions. Instrument eligible operations, original resolution count,
resolved-helper hits/fallback reasons and actual saved lookups in untimed real
CLI runs. Do not time a silent fallback. This audit adds no implementation.

## Why longer-lived binding/type memoization is unsafe without more proof

* The scope creator is FROZEN after its first phase, but declaration maps are not
  globally immutable. `DeclareEsModuleImports` calls `scope.Declare` during
  FlowThrough and can replace a canonical imported TypedVar. An entry-time cache
  of a simple imported binding can become stale. This is a definite blocker to
  a blanket per-function binding cache; operation-local fusion avoids it.
* `InferQualifiedSlot` may lazily declare qualified TypedVars or widen declaration
  types. InferScope order matters for modules, namespace exports, calls which
  back-infer function-argument types, shared properties and escaped variables.
* Implicit arguments/this/super TypedVars are materialized lazily; their initial
  type can depend on current enclosing function/extern type state. Eagerly
  warming or caching them would change creation timing.
* Equal input FlowScopes do not prove a repeated transfer is a no-op. Transfer
  methods read and mutate shared JSType/TypedVar/property graphs outside the
  lattice, and previous node annotations may have changed. Skipping a transfer
  would require comprehensive semantic dependency epochs, not pointer equality
  or an AST change stamp.
* Precompiling an expression's immutable token/operand shape is possible under
  stock ownership, but it still must execute original type operations in order.
  Its cost pays back only when enough CFG nodes are revisited; the default profile
  does not provide the required node-step population. A whole transfer tape is
  premature without an untimed step/size histogram, especially for mostly
  straight-line code. Operation-local fusion pays back on first visits too.

## Larger proposals checked against the archive

* **Native indexed CFG/concrete solver:** already implemented and oracle/race
  checked in historical experiments; the clean main currently uses the generic
  solver. The earlier macro screens did not show a broad win. Porting that
  engine again is a control/re-evaluation, not a new algorithm. Its pure overhead
  is narrower than Analyze minus FlowThrough because joins/branchers also execute
  semantic type operations. Here CreateFlowBrancher alone allocates19,399,448
  bytes (1.20% of whole); full states allocate8,388,992 (.52%). Useful bounded
  allocation targets, insufficient evidence for a fresh whole solver rewrite.
* **Persistent indexed FlowScope slots:** the earlier paged-ID backend already
  had sparse-version and fallback/reconstruction overhead, with mixed/negative
  macros. Existing HAMTs share unchanged roots and paths; no whole-state copies
  happen on each transfer. The fresh flow-map family is only2.08% of bytes.
* **Initial-analysis walk fusion / recorded second-scope schedule:** both already
  exist as separate archived source prototypes. They preserve known callback
  sequences but are not new designs. Recheck activation if retiming; do not cite
  them as unexplored fundamental opportunities.
* **Replace initial syntactic Scope/Var objects with a compact declaration graph:**
  potentially feasible because private first-order facts mostly need names,
  owners and assignment counts, while the registry's nonnullable-name table
  ultimately keys scope roots. However, source still needs exact declaration,
  alias, implicit, bleeding-name, catch and module semantics, and the enum pass
  uses real scope name resolution. Fresh attribution limits syntactic scope
  creation *inside this inference phase* to11,010,768 bytes (10.50 MiB), not the
  whole83MB syntactic family. The necessary scanner/registry adapter redesign
  is disproportionate to this local allocation evidence. Keep it below binding
  fusion in priority rather than claiming it removes the whole scope phase.
* **Parallel scopes:** parent inference can back-infer descendant function
  signatures; siblings can mutate shared variables, registry and property graphs.
  No disjoint ownership partition is established by lexical siblinghood. Exact
  diagnostic and alias order blocks blanket parallel InferScope.
* **Native Array specializations:** current dense-analysis factors already attack
  constructor/replacer work, with recipe/captured-validation/shared-resolved
  controls. Generic fresh-template construction was also tried historically.
  Generalizing it to another raw type needs actual per-type allocation population
  and template/resolve/metadata alias proofs; it is not permission to memoize all
  mutable JSType results. The whole8.44% template-family row is only an upper
  bound, shared with scope construction and other operations.

A larger remaining inference gain may require a broad JSType/type-builder
representation redesign, but these profiles do not yet identify a bounded,
semantically proved10–25% whole-program cut. The viable next isolated source
experiment is the resolved-binding operation seam; its actual hit population,
CPU allocation impact and elapsed-time benefit remain unknown.
