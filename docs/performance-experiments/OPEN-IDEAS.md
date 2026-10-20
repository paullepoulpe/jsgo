# Open performance ideas

The [current round-4 report](ROUND4.md) supersedes historical profile weights
and source-only status below. Main `5de6e335` already fixes the large source-map
allocation issue; current `madvise`, futex and syscall CPU samples are small.
Explicit Java-style string hashing is under 0.5% in the current typed sample.
Those older targets should not drive another large rewrite without new evidence.

Current bounded candidates are a private Sass classifier returning a value and
a parser-node slab cap of 102. Both have focused tests and original-candidate
Java artifact gates, but no measured speed/RSS win yet. Actual assembly caught
a missed-inlining problem in the first slab implementation; a separately frozen
min/max refinement must pass its own gates before the five-variant screen.

The strongest open structural questions are AST locality and scope/variable
arena traffic. The naive per-store metadata COW design modeled more allocation
than its favorable node-slab saving; this does not reject a direct-pointer tagged
union. That diagnostic is still unbuilt. Scope/Var/slot sites allocate roughly
755 MiB on typed input, but numeric requested/reserved histories are still
unmeasured and do not establish live retention. Keyless variable slots require
original-key snapshots before public mutable escapes, and common early iteration
may erase the saving. Exact-type traversal no-op classification has a much
smaller sampled dispatch budget than inclusive traversal suggests and adds
branches to other traversals; broad specialization is not yet justified.

These are remaining experiments and source findings, not claimed gains or an
exhaustive list. The final benchmark report determines what was validated and
measured. All future comparisons should use the same latest-main baseline,
ordinary CLI inputs, shipped PGO, default runtime settings, exact Java artifact
gates, repeated interleaved elapsed measurements, and peak RSS/time ratios.

| Idea | Status and next useful discriminator | Main risks / compatibility boundary |
|---|---|---|
| **Remove remaining scope-free traversal objects** | Compact-stack variants and the concrete gather-walker prototype already target part of this family; the final report records their validation/selection. Next isolate `MustBeReachingVariableDef.ComputeDependence` and `Candidate.GetDefinition`, possibly the use-count callback. First count actual calls/allocations on the final selected binary, then port one callback and compare exact visit/mutation order. | Preserve final CFG cut predicate, FUNCTION/CLASS/default/computed-member order, linked-set insertion and matching Node identity. Preserve exception context without replaying partially applied callbacks. Private stock factory only; public/custom callbacks keep original traversal. |
| **Reuse lexical binding resolution within one inference operation** | Source proof identifies repeated same-context lookup on NAME reads and assignment updates. A private helper could carry the current canonical TypedVar/owner through that one operation, with no retained cache. Start with assignments, then reads; collect real eligible-operation/saved-lookup counters before timing. | Read type values at every original semantic use; keep flow versions, joins, diagnostics and transfer order. `DeclareEsModuleImports` can replace imported TypedVars during inference even after scope creation freezes. Therefore a function-wide binding cache is unsafe; implicit/qualified/custom names retain original helpers. |
| **Re-evaluate native liveness on current main correctly** | Historical native engines are not a fresh-main result. Current BitSet contains a self-referential inline buffer and lattices embed it. A clean port must initialize `lattices[i].liveBits` **in place**, set `liveSet=&liveBits`, and remove old separate sets/temporary struct copies. Compare against the improved current generic engine with positive post-colors activation. | Copying `*NewBitSetWithSize` into another value leaves its slice pointing at the temporary's inline buffer, retaining extra objects and unused buffers. Old owned-word slabs cannot remove backing allocations that current main already eliminated. Public annotation/query/clone/capacity semantics and fallback remain exact. |
| **Publish native liveness words directly** | After the correct in-place bridge, an independent two-word initializer could replace per-set-bit Set loops without extra arrays or pointer fields. Next measure actual published sets, popcounts and node/step-weighted variable counts. A one-word solver state for ≤64 variables is another separate candidate. | Direct initialization must preserve Size, Length, trimming, clone and failure publication. One-word state needs an unchanged two-word fallback, escaped-local semantics and worklist order. The old 128→actual scratch reservation saves no allocation on current main; it is not a useful memory optimization. |
| **Private indexed AST / regions** | A fresh stock CLI may retain function/program regions until compilation ends and still match outputs. A real experiment could use stable facades plus pointer-free intraregion links and typed payload storage. Before implementing, inventory cross-region links, surviving rows, detached payloads, synthesized/cloned nodes and lookup costs on actual workloads. | Public/custom paths need original Node identities and detached aliases. Whole-program owners can retain deleted nodes/types/properties; bounded regions reduce but do not eliminate coupling. Facade branches/index lookups can cost more than current direct inlined GetNext. A shadow index over unchanged Nodes adds anchors and does not remove original scan work. No process-global allocator/pool. |
| **Broaden type-constructor specialization carefully** | The native Array factors cover a specific measured constructor family. An untimed census by raw type, callsite, resolver state and map shape should identify the next material population before generalizing. | Fresh wrappers/maps/replacers and public mutable-slice alias rules must remain exact. Registry closure does not make all JSType/property graphs immutable. Do not memoize mutable type results or count the entire overlapping template-family profile as removable cost. |
| **Use more cores where dependencies permit** | Parsing, annotation, printing and other existing parallel work are already baseline features. New parallel phases require disjoint ownership or a read-only boundary and deterministic replay/diagnostics. First validate on a host with more than the current two available CPUs, with Java receiving the same CPU budget. | Lexical sibling scopes are not independent: inference can update descendant signatures, shared properties, imports and registry/type state. “Scope creator frozen” is not “types frozen.” Require race tests, exact ordered-output proof, activation, throughput/latency and peak RSS scaling. N bounds compiler work concurrency; it is not a global goroutine or OS-thread limit. |

## Avoid double-counting the traversal evidence

The original ff06 sample attributed 159.59 MiB of NodeTraversal construction to
gathering, 42.02 MiB to ComputeDependence and 26.51 MiB to GetDefinition. The first
site is already explored by the private gather prototype (`cbf0034b`, integrated
as `baf328cb` before final validation); this does not establish that every
workload dispatches or that the entire site disappears. The final report owns
its test/activation/timing status; all three originally used the much
larger 16-entry inline stack. **228.12MiB is not an untouched remaining saving.**
The latter two sum to 68.53 MiB in that old sample, and compact4 already reduces
per-constructor storage. Actual remaining costs need fresh final-binary evidence;
sampled allocations are neither peak RSS nor an elapsed-time forecast.

## Source variants that are not performance results

* **Stack8** is a separate source alternative to stack4/16; it needs a direct
  matched comparison before it can be selected as a speed/memory improvement.
* **Deeper lazy/eager and balanced reentrant growth tests** strengthen correctness
  coverage. The isolated obsolete-inline-reference cleanup also adds per-push
  checks; weak-lifetime success alone does not establish a macro gain. Any
  variant absent from the final measured batch remains unmeasured.
* **Current-main `-l=4` build alternative** requires matched repeated benchmarks
  and all artifact gates. Successful compilation, aggressive inlining or an
  older-base result is not a current-main gain.

## Supporting audits

* [Current AST/GC ownership and concrete-walker design](third-followup-records/round2/latest-main-ast-gc-design-audit.md)
* [Current inference and operation-local binding reuse](third-followup-records/round2/latest-main-type-inference-next-audit.md)
* [Current-main liveness publication compatibility](latest-main-design-notes/79d338a71d4f2522-native-live-newmain-materialization-audit.md)
* [Traversal inline retention and reentrancy](third-followup-records/round2/traversal-inline-retention-audit.md)
* [Fresh original-main profile evidence](third-followup-records/round2/latest-main-profiles-v2/evidence.md)

The GC profile identifies real CPU work, but it does not supply a fixed wall-time
ceiling for these designs. Allocation, traversal, type operations, memory
bandwidth and parallel GC overlap. Each proposal needs its own correctness and
activation proof, then measured time and RSS, before combining it with retained
changes.

## Frontend architecture ideas

These remain substantial, unproven designs. The experiments so far do not
establish that full parser/IR redesign is exhausted, and do not predict a win.
No new code, build, profile or benchmark accompanies this note.

### What was actually explored

The current clean-main owned-ASCII factor retains the complete original
Parser.ParseProgram → ProgramTree → IRFactory → Rhino AST pipeline. It changes
only source backing for eligible fresh types-only CLI inputs. Its current
screen was flat/slower; that does **not** measure eliminating the intermediate
parse tree or location graph. Its full legacy parsing and focused differential
passes establish backing equivalence, not equivalence of a future fused parser.

Earlier isolated direct-Rhino experiments initially supported a restricted
expression/statement subset and rejected files containing functions/comments.
Later versions covered substantially more real programs, reused original
annotation/validation helpers, and delegated unsupported syntax to original
parser fragments. Thus accepted-file coverage was not wholly direct parsing.
They still consumed the original Scanner's eagerly allocated tokens, positions
and ranges, and retained whole-file fallback. Some early results were only
component/microbenchmarks; later macro results were mixed. Parser allocation
batching, token-self compaction and last-position reuse tested narrower storage
changes. None was a complete shared-grammar value/event frontend.

### 1. One grammar, alternative parse-tree/event/IR sinks

A compiler-private sink could receive grammar production results as compact
value records/events, then construct ordinary Rhino nodes without a parallel
ProgramTree object graph. Keep the original grammar and parse-tree sink as the
public implementation/oracle. A private event sink is a plausible intermediate
step before attempting immediate Rhino construction; it can defer decisions
without retaining every generated Java-style tree object.

This is more than replacing constructors. Current grammar routines inspect
intermediate-tree kinds and payloads while resolving arrow/parenthesized cover
grammar, patterns, assignments, functions and context. Such inspection needs an
explicit representation/contract. Source comments are available to the original
IR transform after parsing the whole file; fileoverview, licenses, JSDoc claims,
closure-unaware regions and warning order cannot simply move earlier. Bounded
transactional statement/function event windows or a separate attachment stage
must preserve those dependencies. Existing validation walks may be fused only
after proving the same diagnostic chronology and resulting metadata, including
failure/KEEP_GOING paths. Retain fallback; count delegated work explicitly.

### 2. Private value tokens and lazy public positions/ranges

Source identity could live once per lexer/file, with immutable value spans
(offset/line/column endpoints) carried by private tokens/events. Public
SourcePosition/SourceRange objects would materialize only where the existing
public parser, retained comments, diagnostics or other exposed APIs require
them. Narrow coordinates need checked overflow and an original/full-width path.
Ordinary GC-tracked Rhino nodes remain unchanged; this is not an integer-handle
AST, global pool or cross-compilation cache.

Today public token/location fields are read directly, and native expressions
retain position pointers across lookahead and fragment rewinds. Merely changing
a getter, compacting nativeExpression, or recycling existing position cells
cannot remove the scanner's already-created location objects. A genuine shared
lexical kernel plus public/private adapters is required. Bridges must preserve
lookahead payloads, comment callback chronology, rewinds, UTF16 indices and source
identity. Retained comments need independent stable positions; temporary region
storage must not retain unrelated syntax, and identifier/comment strings must
not borrow large source buffers. Existing public APIs must keep their fields,
identity and lifetime behavior.

### Evidence and decision gates

Fresh profiles containing parser token/tree/position/range allocation sites
motivate investigating representation. Inclusive sites overlap; adding them
would double-count allocation, and neither cumulative allocated bytes nor a
simple object-size model predicts live memory, peak RSS or elapsed improvement.
Old corpus acceptance/component timings are not current full-CLI performance.

Any new implementation should stay behind the fresh CLI boundary and keep
public/custom/IDE paths intact. First compare complete AST properties, JSDoc,
features, UTF16 source ranges, diagnostics/order/exceptions, comments and source
maps against the original grammar across valid, ambiguous and malformed input.
Then prove actual ordinary CLI activation/optimized source bytes and remaining
fragment/fallback work, including N1/N2 and race gates. Only afterward compare
paired default-GC CLI time and peak RSS independently. No cached-worker baseline,
GC tuning, semantic shortcut or speed prediction is implied by these ideas.
