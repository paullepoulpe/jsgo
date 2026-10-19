# Next fundamental pass-local designs

Source audit, 2026-10-19. Active source: validated `c7998d45` (the native9 runtime
configuration is independent of commit identity). Read existing structural pass
labels/allocation JSON and experiment sources. No builds, tests, fresh captures,
pprof decoding, GC tuning, workers or preparsed externs.

## Prioritized: FSIV backward-use quotient, not another solver wrapper

The strongest new source finding is an unnecessarily rich result domain.
FlowSensitiveInlineVariables (FSIV) computes full MaybeReachingVariableUse sets
for all local variables and CFG points, but its sole private GetUses consumer is:

    Candidate.HasExactlyOne(reachingUses.GetUses(varName, definitionCFGNode))

That helper counts at most two elements. It does not inspect iteration order or
compare the singleton with the candidate's actual use node. Same-CFG repeated
NAME reads collapse to one set element; numUsesWithinCfgNode is checked by a
separate original candidate traversal.

The entire backward lattice supports only:

* Add one CFG-node use to one variable.
* Kill ALL uses of one variable on an unconditional write.
* Union corresponding variable sets at control-flow joins.

It has no individual-element subtraction, intersection, use-dependent transfer,
or inter-variable dependence. Therefore abstract each variable's set to empty,
one specific CFG-node ID, or many. Union/add preserve a singleton only when its
exact node identity agrees; many is absorbing until a whole-variable kill. This
is an exact homomorphic quotient for the only FSIV query, not an approximation.
The public MaybeReachingVariableUse API and clients needing actual sets remain
unchanged.

A further algorithmic simplification follows: every add within a CFG node adds
that SAME node. Compile the original reverse traversal once into per-variable
kill/gen actions. A kill clears earlier gen, a later add sets gen. This removes
AST rescanning, name-to-Var lookup, inner Node HAMTs, outer Var HAMTs and immutable
full-set copies from every repeated backward worklist transfer. The existing
native CFG/priority queue can be reused; changing only its generic dispatch
would miss this opportunity. Full-variable representation is the first factor;
projection to candidate variables is a separate, unnecessary second step.

### Bounds and evidence

Existing exact CPU pass labels (CPU time, not wall time):

| Pass | Bundle CPU /whole | Typed CPU /whole | SIMPLE CPU /whole |
| --- | ---: | ---: | ---: |
| InlineVariables | 1.53s /7.22% | 0.44s /7.42% | 0.19s /4.04% |
| RemoveUnusedCode | 1.90s /8.97% | 0.32s /5.40% | 0.31s /6.60% |
| FSIV | 1.45s /6.85% | 0.11s /1.85% | 0.26s /5.53% |

The bundle totals21.18 CPU seconds. Initial early-IV is separate where labeled.
The rows are exclusive pass labels, but source work inside them includes scope
creation, traversal, must-def analysis and rewrites. The backward-use quotient
cannot save the whole6.85% FSIV share. Parallel cohorts, GC and kernel scheduling
mean a CPU percentage cannot be called a wall speedup. Even removing these
three entire bundle passes would cap the direct old CPU reduction near23%; no
real prototype approaches that elimination. Closing a10–25% wall gap requires
several material wins or a larger shared cause, not optimistic addition of
nested profiles.

Existing bundle cumulative allocation:

| Family | Bytes | Whole allocation |
| --- | ---: | ---: |
| Entire FSIV | 674,756,466 | 13.31% |
| Entire IV | 382,352,567 | 7.54% |
| Entire RUC | 326,206,203 | 6.43% |
| FSIV Var-to-use-set HAMT nodes | 69,210,240 | 1.36% |
| FSIV Node-to-Node HAMT nodes | 57,150,880 | 1.13% |

The two HAMT shapes correspond directly to the backward-use representation;
additional child slices, joins and lattice wrappers exist. The corresponding
nearest-caller Put68.69MB and Join39.32MB rows overlap these allocation leaves
and must not be added again. Dense quotient matrices/programs still allocate,
so126.36MB is eliminated-family evidence, not a predicted net saving or speedup.
A dense all-variable implementation must bound its N*V memory and diagnose
admission by graph weight, not quietly benchmark fallback.

### Exactness and admission risks

Canonical Vars come from the same allVarsInFn map and escaped set as the original
analysis, including its existing shadowing decisions. Same-Var aliases need the
same dense column; distinct equal-hash keys conservatively fall back to preserve
archived HAMT collision behavior. ON_EX edges suppress kills exactly as before.
No candidate filtering or reordering is required. Candidate rewrites still use
the pre-rewrite backward result just as the original pass does.

The quotient can converge earlier. The original per-node cap is20,000 steps,
not a millisecond limit. An admitted stock, immutable-AST transfer has no custom
callbacks or diagnostics, but preserving failure matters: conservatively prove
the original step bound before dispatch. Original IN sets grow monotonically
from empty and contain at most F distinct generated (Var,CFG-node) facts. Each
successor changes at most F times, so a node is processed at most1+outdegree*F
times; strict backward-priority graphs need one pass. Otherwise use the original
backend. This may reject useful large graphs and needs weighted eligibility
counters before judging the design. Unsupported/malformed preparation must defer
to the original solver's error order. Public/custom/debug/step-observer paths
stay original. Dense arrays/programs are invocation-owned and released before
function processing returns; no global fact cache or scope identity reuse.

A separate source prototype is now underway on `perf/fsiv-use-cardinality` from
`c7998d45`, with all-state oracle and real-pass/CLI activation fixtures. Its
source is not yet built or timed. Dense_analysis independently confirmed the
consumer/transfer restriction and reviewed the lattice proof; its canonical-Var
alias finding was corrected. Parallel_passes owns CommonJS and dense_analysis
owns Array specialization; these tracks do not overlap.

## Second: native IV symbol-state engine

IV still models each variable with VarExpert, InitialValueAnalysis, interface
analysis outcomes, closure-based inline actions, delayed-negative objects,
ObjectLinkedHashMap handled analyses and SetMultimap alias retry handlers.
The seven expert memo cells are ALREADY embedded; reference storage is ALREADY
one native scoped-name index with ordered entries and mutable-map escape. A
second collector index or another callback trampoline would repeat old work.

A private engine could instead use compact typed symbol records, tagged outcomes
(no-inline/aliases-safe/alias-wait/inline-kind), direct inline operations and
ordered alias adjacency. Current-scope handled state can use canonical Var
indices with root/name validation; cross-scope alias retries need retained record
IDs and exact scoped-name equality, not raw Var pointer identity. This replaces
an entire pass-local decision protocol rather than pooling its objects. Original
variable and retry order, lazy negative decisions and mutation-before-next-var
behavior must remain exact. Eagerly precomputing reference facts for the whole
scope is unsafe: earlier inlines change AST parents/assigned values that later
Reference methods read. Keep facts lazy at original decision points first.

Old bundle leaves include expert objects25.17MB and handled-analysis map Put /
chain17.30MB /7.96MB, plus outcome/closure/multimap costs not fully inventoried.
These are narrower than IV's382MB/7.22% CPU; source alone does not show enough
remaining dynamic-dispatch CPU to promise a large gain. A closed typed decision
engine is feasible, but the quotient's semantic reduction has stronger evidence
and a much smaller validation surface. Proposed only after its result is known.

## Third: native RUC reachability-record kernel

RUC's continuation FIFO experiments changed a small boundary. VarInfo remains an
ObjectLinkedHashMap[*Var, interface], and each RealVarInfo, Removable subtype,
builder and list is still an object/interface graph. A private typed mark/sweep
kernel could own symbol IDs, compact flags, removable-kind records and ordered
adjacency lists, while keeping exact AST traversal decisions, first-encounter
order, FIFO activation, property pinning and final removal order. The current
public algorithm remains an oracle. This is substantially more than changing a
queue or bump-allocating continuation records.

Existing bundle leaves: RealVarInfo40.37MB; varInfoMap chain28.81MB; Assign18.88MB;
IndirectAssign10.49MB; NameDeclarationStatement9.96MB. The nearest GetVarInfo
54.70MB row overlaps its map allocations. The whole phase is8.97% old CPU, not a
measured dispatch-only ceiling. A staged kernel can replace VarInfo lookup/state
first, then closed removable variants; actual symbol/edge distributions and
method CPU are still unknown. Its validation surface is much larger than FSIV.

Do not replace it with an eager ordinary static DCE graph without more proof.
The existing algorithm models conjunction-like dependencies, nested lexical
scopes, nonlocal assignment effects and lazy dead-subtree traversal. Eagerly
querying side effects or resolving types in originally skipped subtrees can
change diagnostics/lazy mutation. Preserve the original traversal/continuation
schedule while changing record representation first.

## Lower priority and excluded repeats

Cross-pass IV→RUC declaration recipes were actually tested and mixed/negative;
whole Scope reuse had mutation/parent/lifetime problems. Unchanged-function FSIV
result caching lacks dependency invalidation for descendants, ancestors and
shared side-effect/type summaries. Native reference storage already eliminates
the duplicate index; reference batches, pass tapes and intrusive FIFO have
mixed screens. Shared string cells and119 direct Is* predicates already exist.
Node64, generic inference solver rewrites and scope pooling are not new plans.
The archived stable-ID type-flow pages grew sparse versions and bundle RSS.
None should be presented as an unexplored large win.

No audit proves all drop-in-compatible designs exhausted. The immediate next
experiment is an exact query-driven analysis reduction; subsequent whole-pass
native kernels remain plausible but need fresh targeted evidence and strict
artifact/oracle validation. Until matched timings exist, report hypotheses,
activation populations, time and memory separately, without claiming Java parity.
