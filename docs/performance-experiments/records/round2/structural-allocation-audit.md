# Default-GC remaining allocation redesign audit

2026-10-19. This reads existing `structural-profiles/*/profile.live.pprof` after
forced GC, using cumulative `alloc_space` and `alloc_objects`; it captures no
new workload. The pre-final-GC allocation profile can lag publication. Heap
profiles retain 512 PCs: observed maximum location depth is 425 for bundle and
308 for typed/checks, with no allocation mass at depth 512. CPU pass fractions
come from the existing exact `jsgo_pass` label summary, not heap-stack guesses.

| Task | Exact MemStats allocated MiB | Exact mallocs | Sampled heap allocated MiB |
|---|---:|---:|---:|
| bundle | 4,812.485 | 80,364,505 | 4,836.436 |
| typed-app | 1,220.378 | 21,721,475 | 1,254.331 |
| checks-only | 1,771.006 | 25,415,993 | 1,747.871 |

Sampled allocation estimates differ from exact MemStats totals, especially for
small objects. The following figures are sampled estimates, not exact savings
or RSS. Cumulative rows overlap; exclusive leaf/nearest-caller rows do not.
Heap owner attribution uses identifiable outermost pass Process roots. Anonymous
factory closures such as inference remain explicitly unattributed at that level;
their full function ancestry is still available. PhaseOptimizer wrappers and
ReferenceCollector helper Process calls are excluded as pass owners.

## Strongest isolated new candidate

Contextual name restoration owns 313.56 MiB / 1.88M sampled objects on the
bundle: 6.48% of sampled allocated bytes. Its labeled CPU is 0.76/21.18s (3.59%).

| Operation | MiB | Sampled objects |
|---|---:|---:|
| AddPotentialShadowVariable | 135.88 | 274,195 |
| AttemptRename temporary disallowed-name sets | 125.15 | 360,586 |
| Whole contextual inverter pass | 313.56 | 1,880,855 |

The first two operations build generic linked sets. Shadow records require
identity deduplication and insertion-order name reads. Temporary name sets are
never iterated after construction; they serve membership tests only. A private
stock-CLI path can replace shadow links with contiguous identity IDs and name
membership with native temporary storage, retaining original public/custom paths.
Names must still be read once per distinct shadow, in original order, at the
rename decision: earlier renames change the current names. UTF-16-equivalent and
malformed strings need a deliberate original-collection fallback.

This has a smaller correctness surface than replacing all Scope storage. The
6.48% allocation and 3.59% direct CPU fractions are upper bounds; the proposed
replacement still allocates and may retain records longer within a local-scope
cohort. No elapsed/RSS benefit has yet been measured. Parent authorized this
prototype separately, default off.

## Larger but less isolated remaining opportunities

Untyped scope declaration storage is larger: nearest-caller allocations are
296.22 MiB in DeclareInternal, 238.01 MiB in NewVar and 140.51 MiB in child Scope
construction (674.74 MiB / 13.95% of bundle bytes combined). These are separate
allocation sites. The typed-declaration-table candidate does not cover untyped
Scope. Blindly adding an 8-byte native-table field to Scope would grow its
48-byte allocation class to 64 bytes across roughly three million child scope
allocations, offsetting about 47 MiB before replacing any map. Direct GetVar,
canonical Var identity, implicit variables, live public iterators, mutation,
reentrant callbacks and compiler ownership all need an explicit representation
contract. Reusing old scopes after AST mutation is a different, harder proof.

CoalesceVariableNames and DeadAssignmentsElimination own 343.31 and 238.37 MiB
respectively on bundle (581.68 MiB / 12.03% combined); these entire pass totals
are not all removable CFG overhead. The existing native Must/Maybe backends do
not cover their LiveVariablesAnalysis graphs. Among exclusive allocation leafs,
the two passes spend 54 MiB / 2.61M objects on NewBitSetWithSize and 47 MiB on
linked CFG node allocation alone, with additional edges, priority records,
queues, lattice wrappers and graph collections. Their existing CPU labels total
1.68/21.18s (7.93%). Since shipped liveness limits analysis to 100 variables, a
direct native CFG plus two-word live-state kernel is a substantive next path.
It must preserve parameter escape side effects, transfer order, exceptional
edges and pass query APIs; it is more work than compact inverter records.

Peephole condition minimization allocates 111.51 MiB / 2.44M MeasuredNodes on
bundle (plus synthetic AST nodes). A private compact condition plan could avoid
temporary object trees and delay AST construction, but public mutable
MeasuredNode APIs and all expression evaluation/tie-break orders need a separate
proof. Its pass CPU mass overlaps multiple peephole cohorts, so no direct
111-MiB-to-time projection is justified.

Checks-only has a concrete inference construction site: TraverseArrayLiteral
allocates 85.04 MiB / 904k objects by repeatedly constructing Array<?> types;
TraverseObjectLiteral allocates 26.52 MiB / 525k objects. A shared global Array<?>
would conflate mutable identities. A per-literal or frozen type representation
requires an audit of template resolution, native registry mutation, contextual
node type replacement and published type consumers. The native_types owner has
been informed; no new type cache was implemented here.

## Existing tracks excluded from the recommendation

Reference constructors alone account for 151.50 MiB / 4.96M bundle objects;
ReferenceCollection growth adds 92.05 MiB and native reference-map appends
221.03 MiB. Existing reference batching/native storage experiments already target
these, so these totals cannot justify a second overlapping arena. The sparse
inlining discovery tape still allocates 279.46 MiB cumulatively; chunked tape and
direct-reference worklist candidates address that separate track. RUC task
arenas, compact typed declarations and typed scope schedules are also under
central evaluation, not newly credited here.

Raw machine-readable audit: `structural-allocation-audit.json`. Reader:
`structural_alloc_audit.py`, based on the existing public profile.proto helper.
