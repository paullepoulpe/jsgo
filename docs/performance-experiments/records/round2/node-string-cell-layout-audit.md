# Node layout and shared string-cell audit

Source-only audit, 2026-10-19. No build, test, new profile capture, profile decoding
or timing was run. Node80 source: `/workspace/jsgo-structural-next` at
`f248613d2ebb4df5596f503d576d4127bb247145`; Node64 source: `3594959c`.
Allocation figures below are copied from existing `structural-allocation-audit.json`
and `frontend-allocation-audit.json`, not newly measured runtime inventories.

## Shared current string headers are already implemented

Both representations use `rhinoStringPoolCell`: a mutex-protected
`map[string]*string` that creates one immutable 16-byte string header plus a
cloned backing string on a pool miss. Repeated StringNodes of the same spelling
point at that same header. Node80 stores it directly in `str`; unextended
Node64 stores the same pointer directly in `payload`. There is no separate
per-StringNode payload allocation in either case. The global pool clears at
1,048,576 names, so sharing is by pool generation, not permanent pointer identity.
String equality remains value equality.

The parser-owned NodeAllocator.NewString and ordinary IRFactory string paths
both call SetString and therefore already use this pool. Public SetString
replaces the node's pointer; it does not mutate an existing shared cell. String
cloning copies the current cell pointer without reinterning or cloning its
backing string. Cells are GC tracked and remain alive through the nodes after
pool clearing. Parser-local interning would remove some global mutex operations
but create separate copies of names repeated across input files; it does not
remove per-node headers, since those were removed previously.

Original names differ. SetOriginalName/SetOriginalNameFromName create a new
escaped string-header cell, and GetOriginalName returns that stored mutable
`*string`. Clones deliberately share its pointer. The payload tests verify that
modifying this original-name header does not change the current name. Aliasing
it with the shared current-name cell would violate that property. Sharing
original-name cells across unrelated nodes would also expose new cross-node
mutation through the existing public API. A new immutable StringCell is not a
free substitute for this public pointer contract.

## Exact amd64 layout constraint

The current Node80 payload is:

| Field group | Bytes |
| --- | ---: |
| token int32 + linenoCharno int32 | 8 |
| next, first, parent, previous pointers | 32 |
| length int32 + kind byte + type/color bool + padding | 8 |
| typeOrColor pointer | 8 |
| propListHead pointer | 8 |
| shared string-cell pointer | 8 |
| per-node extension pointer | 8 |
| Total payload / individual allocation class | 80 / 80 |

Node64 keeps the four links, type pointer, property pointer and one payload
pointer (56 bytes), leaving eight bytes for packed metadata and source position.
It packs token/kind/type tag/ordinary length and uses cold overflow metadata for
arbitrary int32 token and length values. Ordinary string reads access the shared
cell; extended strings first read a private 24-byte `{coldMeta; value}` payload.
The common header makes GetToken smaller, but cannot eliminate this additional
extended-string read without changing where its per-node mutable metadata lives.

Keeping full direct token32 + position32 and the same seven pointer slots already
uses exactly 64 bytes. Length32, subclass kind (independent of SetToken) and the
type/color discriminator still need storage. Keeping these direct while merging
str/ext into one payload gives 72 bytes, which the actual Go toolchain rounds to
the 80-byte class (`internal/runtime/gc/sizeclasses.go`). Parser-owned []Node
blocks could save raw 72-versus-80 stride bytes, but individual nodes would not;
this also introduces the same payload discrimination that current Node80 avoids.
A uniform 64-byte version must move/compress another field or remove a link.
It cannot preserve all direct accesses merely by sharing string cells again.

Numeric values cannot be placed as raw float bits in a pointer slot: Go's GC
would interpret the bits as a pointer. Current Node64 uses an ordinary typed
`*float64` allocation (8-byte payload; tiny-allocation rules apply), or a typed
24-byte promoted numeric payload. BigInt uses its existing GC-tracked pointer.
There is no uintptr payload or dynamically changing GC bitmap. Adding a separate
float64 scalar alongside a genuine pointer again grows the node to at least72
bytes. Removing previous/first or moving them behind a descriptor changes common
traversal/edit operations and requires a new, independently justified design.

## Allocation evidence and limits

Existing Node80 sampled allocation summary:

| Workload | Whole alloc bytes | Individual NewNode bytes | Parser NodeAllocator bytes |
| --- | ---: | ---: | ---: |
| bundle | 5,071,370,244 | 99,622,320 | 115,473,010 |
| typed-app | 1,315,261,615 | 20,973,120 | 24,331,937 |
| checks-only | 1,832,776,077 | 11,010,888 | 124,429,266 |

NodeAllocator object counts are chunk allocations, not AST node counts; slab
bytes include unused capacity. Combining Node80 allocations for bundle gives
about215.1 MB (4.24% of cumulative allocation). Shrinking just that payload from80
to64 saves about43.0 MB, or0.85% of cumulative allocation, before different cold
payload costs. This is not an elapsed-time ceiling: live nodes can be a larger
share of resident memory and locality can matter disproportionately. The prior
matched Node64 screens already directly tested the broad locality hypothesis.

The existing frontend-only leaf summaries attribute pool allocation bytes of
3,752,524 (bundle), 1,614,114 (typed), 5,849,708 (checks), and2,152,152 (whitespace)
to rhinoStringPoolCell. These include map growth and unique headers, not a
per-node string payload. They are subsets of total pool work because the report
is scoped to frontend ancestry. The sampled object counts cannot be treated as
exact unique-name counts (large map allocations and small headers have different
sampling behavior). Typed whole-pass summaries additionally show4,718,736 bytes
in node_getExt and1,572,888 in SetOriginalNameFromName; these are not complete
counts of promoted strings or all original-name setters.

No existing exact inventory establishes N total nodes, S StringNodes, unique
current-cell pointers C, distinct spellings U, or E extended strings. Such an
untimed diagnostic would need to walk real parsed and optimized ASTs and record
those values by kind/promotion plus source/clone provenance. `16*(S-C)` header
bytes are already avoided. Parser-local pools could increase C relative to the
current cross-file sharing. A new payload proposal should first show how it saves
objects/bytes beyond that baseline, with retained graph and pointer semantics
accounted for, rather than interpreting slab allocations as node counts.

## Conclusion

There is no new shared-StringCell implementation to pursue: it is already in
both Node80 and Node64. A direct-token 64-byte GC-safe node still needs another
real storage tradeoff, while a 72-byte descriptor variant does not improve
individual allocation class. No production prototype is proposed from this
audit. Further layout work needs actual promotion/name-sharing inventory and a
specific hot-field benefit sufficient to justify new indirection or API costs.
