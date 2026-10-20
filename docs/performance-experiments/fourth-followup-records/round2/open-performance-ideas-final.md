# Open performance ideas

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

* [Current AST/GC ownership and concrete-walker design](latest-main-ast-gc-design-audit.md)
* [Current inference and operation-local binding reuse](latest-main-type-inference-next-audit.md)
* [Current-main liveness publication compatibility](native-live-newmain-materialization-audit.md)
* [Traversal inline retention and reentrancy](traversal-inline-retention-audit.md)
* [Fresh original-main profile evidence](latest-main-profiles-v2/evidence.md)

The GC profile identifies real CPU work, but it does not supply a fixed wall-time
ceiling for these designs. Allocation, traversal, type operations, memory
bandwidth and parallel GC overlap. Each proposal needs its own correctness and
activation proof, then measured time and RSS, before combining it with retained
changes.
