# Open performance ideas

The next integration is`42aaa853` on main`abcf69a0` (production-identical
to`ae9008f2`; two new test fixtures corrected). Full normal and three strict
Java CLI byte gates now pass; all28 new regression roots pass. All completed warm-JVM results below retain their
named earlier source. New source-only factors target this integration; none
has yet produced a new performance result. The old Scope186/Var341 factor is
deferred because new page-aligned large chunks remove its density rationale.

The prior fully validated frozen baseline was `1fe9acbb`: main `bb521c0f` plus four validated
compatibility corrections. It passes the full normal suite and three Java
CLI byte gates. An 18-row Go-only screen now measures the main update and
compatibility cost separately. Three completed fresh typed JVM cohorts give
fixed Go a Java/Go time ratio1.10789 and RSS ratio2.14739, with every byte
matching. These results include main optimizations, not the unrun factors.

The fresh18-row Go-only screen on42aaa measures3.41% less typed elapsed time
and3.22% less bundle elapsed against1fe, with flat typed RSS and7.97% less
bundle RSS (median same-round ratios; three samples/build/task, raw variation
retained). All18 match Java's bytes. A fresh warm-JVM comparison on42aaa is
still needed. The experimental Node/slab/container/traversal factors remain
unmeasured and unpromoted. Fresh42aaa CPU/alloc profiles now pass both Java gates and are decoded/stopped:
GetNext16.32% flattyped CPU, GCscan union9.77%, madvise+otherSyscalls0.318%.
Typed sampled allocation4182.08MiB vs prior5788.92 is diagnostic cumulative
traffic, not a prediction of elapsed/RSS. The older attribution below is historical.

Historical1fe fixed-main profiles attribute 13.40% of typed CPU to traversal methods,
8.99% to the inlined GetNext load and a 14.16% union to GC scan stacks. These
are attribution views, not additive removable budgets. Typed allocations
total 5788.92 MiB, including scope 382.12 MiB, string/pointer LinkedHashMap
346.69 MiB, Var 195.85 MiB and CFG nodes 190.09 MiB. Source-only pilots target
scope/Var allocator geometry, the common traversal path, large container
growth, a physical Node64 metadata union, and cap127 node slabs. Node80/cap127
is a separate control for Node64/cap127; Node64/default128 remains separate.
Their tests, actual allocation-path inlining and timings remain pending. No
GC tuning is part of these factors.

A native assembly audit of the exact stockabcf CLI rejects one broad premise:
string-key/pointer-value ordered maps already use direct fast-string runtime
lookups and direct string equality, including concrete Scope and namespace
lookup callers. Generic dictionary type metadata and explicit mapping callbacks
remain distinct costs. This does not justify copying a specialized string-map
implementation. A private callback/miss-path change requires fresh profiling;
no elapsed or RSS improvement was measured by this audit.

The older proposals below retain their named historical source/profile
context. Main now implements private small-variable staging, independent
lazy scope arenas and another concrete scope-free walker; adapt older census
hooks before using them on this baseline.

A concrete isolated indexed-AST model (`0ec2cc95`) now passes nine roots
under normal, race and package-only checkptr builds. It preserves stable
native/foreign/copied-node identity with 56-byte handles and GC-tracked tables.
It does not change Rhino or the CLI. Append-only node/type history, pointer
registries, fallback snapshots, and extra getter/registration work remain
explicit costs; the smaller header is not a compiler memory or speed result.

The [current round-4 report](ROUND4.md) supersedes historical profile weights
and source-only status below. Main `5de6e335` already fixes the large source-map
allocation issue; current `madvise`, futex and syscall CPU samples are small.
Explicit Java-style string hashing is 0.73% in the latest typed sample.
Those older targets should not drive another large rewrite without new evidence.

The five-variant Sass/slab screen on `5de6e335` is complete: 30 invocations,
all Java artifacts matching. The refined slab-only variant takes about 1.7%
less typed-build time against its same-binary off control; three samples per
variant and the control's difference from stock make this an exploratory lead,
not a selected improvement. Sass alone shows no reliable gain, and neither
factor establishes a consistent typed RSS saving. Actual assembly caught and
corrected a missed-inlining problem in the first slab implementation; both
source states and failed fixture attempts remain archived. Neither factor has
been measured on `bb521c0f`.

The strongest open structural questions are AST locality and scope/variable
arena traffic. The naive per-store metadata COW design modeled more allocation
than its favorable node-slab saving; this does not reject a direct-pointer tagged
union. The physical Node64 prototype at `99f9fc7b` now implements that approach
in Rhino, preserving hot pointer links and full metadata widths; it is still
unbuilt. Immutable tuple allocations and getter branches must be measured,
and simultaneous disjoint metadata writes to the same Node are a new shared-
representation boundary. No affected stock path was found in the bounded
source review; that is not whole-compiler race evidence. Scope/Var/slot sites allocate roughly
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

A bounded recheck also leaves a narrow SourcePosition-only slab factor open.
The old experiment changed all generic slab ramps and node allocation together,
had one macro run per task, and did not establish position-specific inlining or
retention. Current token-position slab samples total 168.75 MiB on typed input,
but only an unknown maximum-chunk portion could benefit from a 127 cap. More
frequent rollovers can reserve more memory for some file sizes. Comment
positions keep their independent lifetime; SourceRange/token/tree allocation
would remain unchanged. A growth/tail census, actual caller inlining and retained
position tests are needed before another small measured trial. This remains an
open follow-up, with no new implementation or performance result.

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


### Latest-main checkpoint: e5d0987a / 4159dc5b

The bare Node64/cap128 experiment is now **measured and rejected for promotion**:
mean paired elapsed ratios are 1.053216 typed and 1.080983 bundle; peak-RSS ratios
are 1.032740 and 1.017659. All sixteen rows match the Java bytes, four balanced
pairs per task, default runtime, actual 4 GiB/two-CPU caps. Complete codegen,
failed infrastructure attempts and stopped/OOM evidence are in records v9.
Header size alone did not predict compiler speed or RSS.

Open follow-ups remain source-only: specialize common cold-field setters and
hoist template decoding without changing value-copy/nonzero-init/panic behavior;
measure cap127 allocator geometry; rebase container-growth and traversal factors;
and run a metadata-only scope-pool census only if new profiles still justify it.
The private ArrayList pilot targets only about 0.24%/0.12% of old sampled bytes
and can increase partially filled-block retention, so it is lower priority.
No modeled allocation reduction is an elapsed/RSS gain until measured.

Two newer upstream rounds changed scopes, parser arenas, streaming output,
template maps, type caches and PGO. The f1 integration passes full normal tests;
the latest e5 integration is source-reviewed but runtime validation is pending
at this checkpoint. Validate and profile the latest correct compiler first.
Do not time the stock template shortcut that skips function-signature work.
Keep ordinary CLI versus freshly warmed Java, generated BigQuery/protobuf,
large many-chunk builds, memory caps, variance and RSS ratios as the measurement
priorities. Previous madvise/futex hypotheses no longer dominate the completed
42aaa853 profile and must not drive speculative redesigns without new evidence.
