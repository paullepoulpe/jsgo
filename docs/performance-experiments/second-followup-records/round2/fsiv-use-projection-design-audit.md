# FSIV quotient projection and successor designs

Source-only during the compound timing reservation; no builds/tests/profiles.
Current corrected source9d190a2b; standalone projection prototype derives from
corrected quotient fcc17e2c in /workspace/jsgo-fsiv-use-cardinality. Kernel audit
uses f79f7da8. Existing diagnostic JSON is read, not newly captured.

## Existing actual activation and size

The corrected all-variable quotient really runs on the pinned inputs:

| Task | Runs | Canonical unescaped columns summed | Actual queries | Matrix bytes cumulative | Largest matrix pair |
|---|---:|---:|---:|---:|---:|
| typed-app |1352|5878|1393|1,257,888|111,280|
| bundle |16494|88147|17475|27,397,952|676,368|
| downlevel-simple |3697|16378|2618|3,005,584|76,176|

All attempted analyses admitted in this capture. The source is actual corrected
post-colors state with CLI ownership and parallel FSIV enabled. Data files:
`array-recipe-diagnostic-overlay/runs/<task>/job-0/post-colors-activation.json`.
These diagnostic timings are not performance results.

Average columns/function are about4.35 typed,5.34 bundle,4.43 simple. Bundle's
node-weighted width is10.71 (matrix bytes/(8*CFG nodes)), so larger functions have
more columns. Actual distinct queried Vars are at most the query counts: at most
23.7%,19.8%,16.0% respectively of the summed existing columns. This is only a
perfect-knowledge count bound. It does NOT measure the selected width from all
initial candidates, does not weight by CFG size, and cannot predict runtime or
matrix savings. Many initial candidates fail a pre-query check; running those
checks early would change the original mutation-dependent schedule.

## Exact first factor: complete-candidate projection

The implementation snapshots every candidate's original allVarsInFn Var after
GatherCandidates. Later code only removes candidates, never adds or renames
records. All future query variables are therefore contained in the snapshot.
Keep every allVars name alias to each selected Var, not just candidate spellings.
Keep original escaped exclusion and canonical pointer aliasing. All backward
operations are componentwise gen/whole-kill/union; no omitted-variable dependency
exists. Forward MustDef and Definition.depends remain full, preserving candidate
pruning after earlier rewrites. Parameter/loop/side-effect early checks retain
original execution order; no speculative CanInline pass occurs.

Full-domain validation must precede projection. Distinct omitted Vars with equal
HAMT hashes can affect original collision behavior, so full collision rejection
remains. Full AST validation remains even when all columns are omitted. Most
critically, the original step-cap failure proof must retain ALL generated facts;
excluded variables can determine original convergence and may not be removed
from that safety bound. The prototype compiles/proves full programs first, then
compacts actions/columns in place before allocating state. Nil selection means
original full product; nonnil empty selection means a validated zero product.

## Sparse storage: reject a blind hash-map replacement

Current average width is small, and int32 matrices are contiguous/pointer-free.
Per-row Go maps would add headers/buckets/pointers at widths where a whole dense
row is often20–44 bytes. Fixed sparse-pair structs with overflow pointers can also
exceed those rows before holding values. Persistent sparse rows reintroduce
allocation/versioning that the quotient just removed. Current largest matrix pair
is only676KB, so cumulative27.4MB is not a27.4MB retained-RSS opportunity. Measure
selected widths and nonzero-state density before choosing a representation.

## Larger exact successor if selected widths remain costly

The product is independent by Var. A scalar solver can process one selected Var
at a time with two CFG-sized int32 rows, then retain only values at that Var's
candidate definition nodes. All candidate (canonical Var,defCFG) pairs are known
before rewrites and never grow. Current candidate decisions read exactly those
OUT values; other CFG state has no consumer. This could reduce state fromO(N*S)
toO(N+queried pairs) and avoid coupled requeues sweeping unrelated variables.

It requires one full original validation/cap proof first, a preindexed per-Var
kill/gen program, exact original singleton CFG IDs/ON_EX behavior, alias handling
and graph-membership failures. The admitted transfer is callback-free and the
original capped failure is proved impossible, permitting independent scalar
convergence without changing reporting/exception behavior. Do not drop the full
proof merely because scalar convergence is shorter. Repeated heap setup per Var
could outweigh reduced coupled iterations; it needs a separate factor and counters
for pairwise steps/heap operations. Dense projection is the simpler first step.

The original full FSIV label was6.85% CPU, not wall time, and backward use is only
part of it. Neither the quotient nor projection establishes a10–25% whole-program
win. Wider ambitions should combine independent measured phase reductions rather
than count the whole pass as removable work.

## Composition with inline candidate kernel f79f7da8

That kernel specializes InlineFunctions.FindCandidateFunctions, preserving its
live ordered traversal/context and optional reference itinerary. It does not
replace FSIV's GatherCandidates or dataflow. Its itinerary lives only through
the same InlineFunctions invocation; quotient/projection state lives only in the
current FSIV EnterScope and is retired before return. They have no shared cache,
queue, scope IDs or transfer state. PhaseOptimizer runs the whole passes in the
original sequence; FSIV builds fresh canonical scope facts/CFGs from the current
AST after any prior inlining. Forward and rewritten AST semantics remain real.
No source-level composition blocker found. Combined artifact/race gates and an
isolated factor screen remain necessary; this source audit does not prove gains.
