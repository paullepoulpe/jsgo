# Type inference and FlowScope representation audit

Source-only audit on 2026-10-19. Read Node80 structural source at `f248613d`,
existing `structural-allocation-audit.json`, and the archived indexed-flow-slot
experiment documentation. No build, test, capture, profile decoding or timing.
The sampled summaries are evidence for prioritization, not current Java/Go ratios.

## Already optimized

The enabled native CFG path already builds direct topology and runs a concrete
TypeInference worklist. It keeps in/out states in one array and branch states
in another, shares the immutable bottom lattice, passes a single predecessor's
state directly, uses stack join reducers and one stack branch-condition cache
per transfer, and compares outputs in the original short-circuit order. It
preserves worklist/edge ordering and the original FlowThrough/type operations.
Reimplementing generic DataFlowAnalysis again would duplicate active work.

LinkedFlowScope already avoids whole-state copies: both the outer TypedScope
map and inner inferred-name maps are persistent HAMTs. Updates path-copy changed
paths while sharing unaffected subtrees and slots. OverlayScope.Infer returns
itself for an identical type pointer; outer map Plus then returns the existing
map, and InferSlotType returns the original flow state. WithSyntacticScope
returns the input when scope identity is unchanged. JoinFlow skips identical
scope-map/function-scope pairs, and HAMT Reconcile/Equivalent short-circuit
identical subtrees. InferDeclarativelyUnboundVarsWithoutTypes scans a scope only
once per inference using its existing inferredUnboundVars set. Entry-parameter
inference is not an all-nodes repeated scan.

Immutable versions are necessary: branch outcomes, short-circuit expressions,
reverse-interpretation results and CFG predecessor states retain older inputs.
Type operations and qualified-name declarations also mutate JSType/TypedVar
state outside the lattice. A syntactic expression being unchanged does not
make its inferred output reusable across fixpoint iterations.

## Allocation attribution, checks-only

Existing whole cumulative allocation: 1,832,776,077 bytes /22,972,844 objects.
The following inclusive rows overlap and cannot be added:

| Existing row | Bytes | Whole bytes |
| --- | ---: | ---: |
| InferAllScopes (includes initial scope construction) | 615,194,728 | 33.57% |
| InferScope | 325,121,408 | 17.74% |
| runNativeTypeInference | 236,445,371 | 12.90% |
| native compute | 223,180,637 | 12.18% |
| FlowThrough | 189,864,156 | 10.36% |
| TraverseObjectLiteral | 113,307,268 | 6.18% |
| TraverseArrayLiteral | 89,694,758 | 4.89% |
| CreateTemplatizedTypeOfList | 122,779,120 | 6.70% |

Concrete allocation leaves show a substantially narrower flow-storage target:

| Leaf | Objects | Bytes | Whole bytes |
| --- | ---: | ---: | ---: |
| NewLinkedFlowScope (48-byte state) | 163,845 | 7,864,680 | 0.429% |
| New OverlayScopeWithSlots (24-byte overlay) | 152,915 | 3,670,100 | 0.200% |
| String-key pointer-value HAMT nodes | 114,688 | 7,340,480 | 0.401% |
| Pointer-key pointer-value HAMT nodes | 139,266 | 8,913,440 | 0.486% |
| runNativeTypeInference direct arrays/objects | 63,730 | 11,167,502 | 0.609% |

These HAMT generic shapes include other consumers; the top-N summaries do not
supply a complete FlowScope allocation union, all child-slice allocations,
overlay-slot counts or flow update/join population. They are not all savings.
A prior mixed profile attributed about6.75% CPU and70.5 MiB to type equality,
unions and flow/HAMT together, but it predates current changes and overlaps
semantic type construction. It cannot justify a6.75% solver speedup forecast.

The larger concrete constructor target is repeated semantic type construction.
TraverseArrayLiteral unconditionally calls CreateTemplatizedType(Array, unknown)
after traversing children, even on repeated CFG visits. The source and the
89.69MB/122.78MB overlapping rows support the parent's current bounded template
candidate more directly than another flow-map backend. Those rows are allocation
ceilings for that family, not a proof that result sharing is legal or all bytes
can be removed: identity, mutability, resolver lifecycle and callback order still
need its existing candidate's gates.

## One feasible bounded flow successor

A narrowly actionable opportunity is to preserve immutable overlay identity on
no-change joins. LinkedFlowScope_Join currently always constructs a new
OverlayScope after reconciling its slots, even when the returned slot map is
exactly an input map and bestScope is that input overlay's same TypedScope.
A private stock join could perform every original ordered reconciliation/type
operation, then return that existing overlay when both identities match. The
outer HAMT can consequently retain an input path; a flow wrapper can be reused
only if map, syntactic scope and function scope all also match. This does not
introduce a cache or mutable state and does not rely on an AST change stamp.

This is not a large demonstrated win. The3.67MB overlay leaf bounds all such
wrappers, not only join wrappers; the8.91MB pointer-key HAMT shape bounds a broader
family. Even eliminating both whole leaves would be only0.69% of cumulative
allocation, before all remaining costs. Actual hits and saved outer paths need
an untimed diagnostic. It should retain the public/custom join as the oracle and
compare ordered type operations/diagnostics, bottom joins, shadowing, nil types,
three-way joins, qualified names, hash collisions and historical versions.
Future identity shortcuts must not suppress custom or lazy type callbacks.

The generic HAMT also has an adjacent asymmetry: one-sided ReconcileTrees
branches clone each node/child slice even when all callbacks return the input
values. Both-nonnull reconciliation already tries to reuse unchanged roots.
A private pointer-valued specialization could lazily copy children and retain
the input after executing callbacks in the original order. The generic function
cannot gain an extra ObjectsEquals call on arbitrary V: that can invoke custom
user equality absent in the original one-sided branch. This is a separate,
unmeasured, bounded structural-sharing option, not permission to alter the
public reconciler's callback behavior globally.

## Designs rejected or deferred

* The archived stable-ID paged inferred-slot backend already explored the
  obvious indexed representation. Its four-slot pages/flat directory can cost
  592 bytes for a sparse one-slot version versus64 for a HAMT node. Equality
  fallback reconstructs HAMTs, joins allocate discarded directories, IDs/maps
  remain attached to TypedScope, and collision/overflow wrappers add overhead.
  Its documented screen had no broad speed benefit and increased bundle RSS.
  Repeating it without occupancy/copy/invalidation evidence is not a new design.
* Coallocating FlowScope, OverlayScope and OverlaySlot as one version object is
  unsafe as a general memory improvement: future persistent versions share old
  slots. An old shared slot would then retain its old flow wrapper and map,
  potentially keeping every prior state graph alive. Unlike first Var+map
  coallocation, the current flow objects are not already mutually reachable.
* A mutable transfer-local map/builder is not automatically safe. FlowScope
  APIs retain old aliases across short-circuit and reverse-analysis branches;
  every fork/publication would need explicit freeze/copy ownership. Ordinary
  visitor methods do not establish unique ownership. The required redesign is
  substantially larger than changing a map implementation.
* Freezing symbol IDs/declaration defaults at function entry is insufficient:
  InferQualifiedSlot can lazily declare qualified TypedVars and widen their
  declaration types during transfer. Join fallback reads those current types.
  Shadowed root names and changes in syntactic scope also affect ownership.
* GetSlot repeats a typed lookup when an overlay lacks a requested inferred slot
  despite the initial syntactic GetVar already resolving it. A stock-only helper
  could reuse that current Var, but the existing summaries give no dedicated
  CPU ceiling. It is a smaller lookup optimization, not a fundamental state win.

No new production implementation is proposed this cycle. Current concrete
solver work plus the larger template-construction candidate should be validated
and measured first. If a flow followup is pursued, instrument no-change overlay
join reuse and one-sided reconciliation before adding storage or global caches.
