# Latest-main AST/GC redesign audit

Source-only audit, 2026-10-19. Main substrate: ff06ac15; combined source inspected
at a3a2f8f3, compact4 at9044aff9. No Go builds, tests, timing, new profiles or new
profile decoding were performed for this audit. Existing fresh ff06 capture and
its decoded allocation-owner stacks are reused. Runtime GC tuning is excluded.

## Finding

A private indexed AST can preserve the public signatures, but a stable public
`*Node` facade cannot simply replace four strong graph links with untracked
integers for free. Either an owner strongly retains the whole facade/row arena,
or live indexed links must somehow keep their target facade alive. Whole-arena
ownership pins detached nodes and their type/property payloads. A second indexed
traversal representation over the unchanged Node tree does not remove its scan
work and usually adds pointer anchors. The prior small-node/slab experiments do
not establish that these more invasive ownership costs will pay off.

The next concrete low-retention experiment should be **private concrete walkers
for the existing scope-free CFG-expression callbacks**, rather than another
parser slab or a purported pointer-free shadow of the same tree. This removes
actual short-lived traversal objects and their scanned state without changing
Node identity, mutation, detachment, property aliases or lifetime. A true bounded
indexed-node region remains an ambitious second-stage experiment only after its
cross-region edge/retention/lookup cost is supported by an untimed inventory.

## Current representation and what earlier work already covers

`internal/rhino/node.go` has an80-byte Node on linux/amd64:16 scalar bytes and eight
GC-tracked pointer slots. The four topology pointers are next/first/parent/
previous; the other four hold type-or-color, properties, a shared current-string
cell, and private extension data. GetNext/GetFirstChild/GetParent are direct
loads and already inline. The stale source comment mentioning96bytes is not the
current layout. The type-or-color unsafe.Pointer is a real GC-tracked pointer,
not a uintptr handle; numeric bits cannot be hidden in it safely.

The earlier64-byte Node packed metadata and merged cold payload storage while
retaining real topology pointers. It therefore changed layout/indirection and
one payload pointer, not the fundamental GC-visible graph. It already tested a
broad small-node/locality hypothesis; it is not evidence that repeating the
same layout will yield a new win. Its old-base results are not current-main
measurements.

The prior parser-owned NodeAllocator emitted ordinary typed []Node blocks,
growing8→16→32→64→128. Nodes did not retain the allocator, chunks were not pooled
or globally published, and copies kept GC-typed fields. This reduced allocation
objects and gave bounded chunk coupling. It did **not** make Node links noscan.
A live interior pointer keeps the entire backing chunk alive. Reintroducing it
under an indexed/arena label without actually removing pointer fields would
repeat the same experiment. The current clean full-parser ASCII factor changes
source backing and uses the full existing grammar/IRFactory; it is independent
of that old native parser and does not supply an indexed AST.

Shared current-string cells are already present, now through the main sharded
pool. Original-name cells have a different public mutable-pointer alias
contract. Interning every original name or creating another per-parser string
pool cannot be counted as a new per-Node header saving. No global mutex-backed
batch allocator or process-global AST owner is proposed.

## Why a facade/arena design is harder than its row size

A genuinely pointer-reduced row could store four32-bit link indices (16bytes),
with scalar metadata, and place GC-bearing type/property/string/extension values
in typed side storage. Public GetNext must still return the exact same *Node,
GetParent/previous semantics must remain live, and map keys/source definitions/
Scope.Var declarations must continue identifying that same facade.

Consider one64-byte facade containing metadata, four payload pointers, owner
pointer and index; a16-byte topology row; and an8-byte strong facade-table entry.
That is already88bytes/node before owner headers/capacity slack, versus80today.
Its six pointer slots/node (five in facade, one in table) reduce pointer slots
from eight by25%, while GetNext changes from one load into owner→row→index→table
lookup plus legacy/arena discrimination. The table/owner cycle retains every
facade and every detached payload in the arena. This is an illustrative model,
not a proposed measured layout.

Contiguous facade values permit `&owner.nodes[id]` instead of a pointer table:
64+16=80bytes/node and five pointer slots instead of eight (37.5% fewer). But any
interior facade pointer plus its owner retains the entire facade array/arena.
A private child detached from a long-lived source tree can no longer die on its
own. Clearing a removed row is incorrect if an external alias, diagnostic,
Reference/Scope or transformation candidate still holds that detached Node;
Detach explicitly preserves its subtree. CloneNode must keep its detached,
shallow metadata alias rules, while CloneTree creates distinct Node identities.
Moving/adopting a node between files cannot silently replace its facade.

A weak facade table is not a solution: an integer link is not a GC root, so a
child reachable only through a live parent's index could be collected. Recreating
a facade from the row changes pointer identity. Raw uintptr/pointer tagging,
noscan storage containing real Go pointers, manual off-heap allocation and a
custom collector are excluded here. Default/public NewNode also lacks a compiler
owner parameter; routing all synthesized stock nodes into a private region
requires explicit factory plumbing or adoption, not a hidden global allocator.

## Public lifetime contract versus private CLI ownership

The ownership issue is not a proof that indexing is impossible or cannot be a
mostly drop-in CLI replacement. For public/custom Node paths, callers may keep
arbitrary detached aliases and expect the ordinary live graph; those paths need
the original representation or fully equivalent stable facades. A fresh stock
CLI lease has stronger boundaries: no public callback receives its AST, and
compilation ends at output/diagnostics. It may intentionally retain a complete
function or even whole-program arena until compilation ends and still match all
Java artifacts. A program-sized strong row/facade owner is therefore an allowed
**risky measurable private option**, not a global API impossibility.

The tradeoff is deleted-node/type/property retention during subsequent passes,
region capacity slack and higher RSS/GC work despite fewer topology pointers.
Retaining the full arena until process exit also prevents reclamation between
optimization iterations. Function-sized regions reduce that risk but ordinary
rewrites move nodes between functions; stable identities and cross-region links
must still work. If the owner is to release earlier while surviving Nodes stay
in the output AST, migration/promotion must preserve all aliases, not simply
reset a phase table. A uniform indexed alternate binary or an explicitly private
NodeRef pipeline can remove facade branching, but requires much broader pass and
factory plumbing; an extra mode branch in every GetNext can instead consume the
locality benefit. Contiguous facades can additionally reduce GC object count,
not just pointer slots. These effects are not determined by a row-size model or
by old Node64/slab outcomes.

A legitimate experiment could compare compiler-owned program versus bounded
function regions using original-public fallback, exact fresh-CLI admission and
all default-runtime macro memory/CPU/elapsed gates. It needs a clear implementation
budget and real retention inventory first. No source-only percentage in this
audit fixes its eventual ceiling; the recommendation below selects a smaller
well-isolated next experiment while existing factors are being measured.

## A bounded indexed-region experiment, if revisited

A technically feasible shape is a private parse-owned region capped at a small
number of nodes, with contiguous stable facades and pointer-free intraregion
links. Cross-region/legacy edges use ordinary GC-tracked overflow pointers;
mutation updates/clears the exact edge storage. Regions never point to a global
registry or all other regions. A live facade pins at most its region directly,
with other regions retained only by actual cross edges. Detach removes its
external links but keeps its child graph; clone allocates a distinct facade;
reparenting does not change identity. Payload storage stays GC-typed and mutable
metadata keeps its original alias boundaries.

This differs substantively from the previous slab only if its actual GC bitmap
omits most intraregion topology pointer slots. Merely filling original pointers
with zeros does not change the bitmap or eliminate scans. The implementation
would require an alternate Node layout and mechanical routing of every link
getter/mutator, plus a fallback for public/mixed nodes. There must be no second
set of strong links shadowing the indices. Cross-edge pointer tables, owner
backpointers, empty fallback descriptors, capacity slack and detached-region
retention all belong in the byte/pointer budget. A small region can lose its
pointer advantage if most edges cross region boundaries; retaining per-node
legacy link records as well can make it worse than the control.

Before implementation, a diagnostic should count topology edges internal versus
external to hypothetical64/128-node regions, actual allocation order versus
preorder order, live rows per surviving region after mutation-heavy passes,
new synthesized/clone nodes, and node payload promotion/type/property occupancy.
Use real bundle/checks/typed inputs and all Node subgraphs (closure-unaware
shadows and JSDoc type-expression roots included). A preorder-only locality
microbenchmark or byte-sizeof test cannot answer this. Capturing maps that
strongly retain discarded nodes must not then be used to claim GC/RSS behavior.
This audit does not implement such a diagnostic or another massive AST prototype.

## Concrete next experiment: exact private CFG-expression traversal

Fresh ff06 allocation attribution shows bundle newNodeTraversal390.71sampledMiB.
Of that, nearest callbacks are GatherCandidates.Visit159.59MiB,
MustBeReachingVariableDef.ComputeDependence42.02MiB and
Candidate.GetDefinition26.51MiB:228.12MiB of direct constructor allocation in
three small scope-free traversal uses. Candidate.GetNumUseInUseCfgNode is another
similar site; inspect its independent contribution before adding it.

These callbacks ignore the NodeTraversal argument. They use tokens, real
Node links/strings, private reaching-definition state or a private candidate;
they neither request a lexical scope nor mutate tree topology. Their final
ShouldTraverse implementation is exactly `parent==nil ||
!ControlFlowGraph_IsEnteringNewCfgNode(n)`. Dependence collection uses the
existing Var map and linked identity-set; definition lookup stores the exact
original matching Node pointer. GatherCandidates also calls the coding
convention, so only an exact admitted stock convention may use the private path.

An invocation-local concrete walker can preserve this pre/post order and cut
predicate while keeping only existing Node/parent values on the Go stack. It
needs no owner arena, pointer table, permanent Node field or Node copy. Its
callback is a concrete private operation, not an arbitrary public callback.
Private compiler/factory admission must reject custom conventions, compiler,
observer/error hooks and public FSIV constructions just as the existing owned
experiments do. The original NodeTraversal path remains the oracle/fallback.
This removes actual heap allocation and scanned transient traversal state; it
does not claim to make the surviving AST graph noscan.

Important implementation requirements:

* Preserve the generic traversal's exact special ordering for FUNCTION, CLASS,
  computed members/default values and any nested CFG boundaries. Ordinary raw
  child recursion is not automatically equivalent. The current shouldTraverse
  predicate is final, but NodeTraversal token dispatch can still affect order.
* Use each callback's existing exact mutation/query order. Dependence-set
  insertion order and the last matching definition matter; identity/order must
  be compared against the original, not only final generated JS.
* Preserve malformed-AST/internal exception context and wrapping. Do not rerun
  partially applied candidate-gather callbacks without rollback: Candidate
  objects have identity, so duplicate callbacks can change the set. A rare
  error may construct the original context wrapper lazily, but only after its
  source/current-node/parent state is shown identical. Unknown or unsupported
  cases use the original path before visible effects.
* Prove complete original/private CLI output+diagnostics+source maps and actual
  dispatch, function/class/default/destructuring/exception fixtures, detached
  subtrees and custom fallback. Weak lifetime tests should show no extra Node
  retention after the invocation; no buffer survives into another function.

The lazy cardinality quotient does not remove these forward-definition,
gathering and definition-lookup traversals, so this is complementary. Compact4
already reduces their per-object footprint: a simple control576→256-byte class
model scales228.12MiB to roughly101.4MiB before growth costs. That is only about
2.2% of the original bundle sampled allocation total, not a guaranteed new
saving, and no current4 capture measures it yet. Eliminating callback objects,
repeated traversal dispatch and scope bookkeeping may affect CPU beyond those
bytes, but this must be measured. Start with one callback family and expand only
when exact-order/exception parity and a real macro gain justify it.

## CPU and memory bounds

Fresh checks scanSpan~18.9% cumulative CPU is not18.9% elapsed time, not all AST
work, and not a target that an AST layout can remove completely. Current profiles
lack a fresh forced-live per-type scan inventory; historical GC-lab Node shares
and3–7% wall tax are useful caution, not fresh-main parameters or permission to
change GC settings. Removing all18.9% of aggregate CPU would bound a purely CPU
model at1/(1-.189)=1.233×, while actual elapsed impact depends on parallel GC,
assists, memory bandwidth and allocation. A layout removing only some link
pointer slots is far below that assumption and can spend its benefit on extra
GetNext loads/branches (already4.9% flat CPU) or retained regions.

Measure ordinary default-runtime CLI elapsed, peak RSS and Java ratios separately
from CPU/sample allocation/GC scan counters. Retain a factor only after combined
correctness gates and repeated interleaved macro measurements. There is no
credible source-only route here to a2–3× wall gain, and no new win is claimed.
