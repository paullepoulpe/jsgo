# Direct-allocation, per-kind AST layout experiment

This branch starts at current-main research base `da20b44f`. It changes only
Rhino's private representation and preserves the exported `*Node` API and object
identity. It does **not** restore the rejected global node batches, use an arena,
retain an allocation registry, or keep a second AST. The source is a prototype;
no new correctness or timing result is claimed until the central validation run.

## Why revisit it

The fresh bundle profile in `jsgo-results/round2/profiles/bundle` has 25.99 sampled
CPU seconds. `GetNext` accounts for 1.79 s flat (6.89%); it is already one direct
pointer load. Rewriting that accessor cannot remove the pointer chase. Layout
and fewer repeated walks are the relevant avenues. `runtime.scanSpan` is 5.49 s
cumulative (21.12%), but it scans all compiler objects, not just AST nodes.

`NewNode` accounts for approximately 181 MiB, 3.56% of 4.98 GiB sampled allocated
bytes. Cutting its 80-byte allocation to 64 bytes therefore saves about 36 MiB
of total allocated bytes, **about 0.71% of whole-compile allocation**, before
payload changes. It does not imply 20% lower compiler memory or runtime. The
post-forced-GC live profile has released the AST and cannot estimate peak AST
memory; measure peak RSS and allocation separately.

The previous packed64 experiment was timed with the obsolete global allocation
chunks. More importantly, it widened a shared cold payload to 48 bytes and used
it even for a number with no other annotations. It also treated equal-valued
current/original strings as interchangeable cells. The new candidate keeps the
caller's exact original-name pointer, including visible mutation of that cell.
It can be measured independently on today's direct allocations.

## Layout and payload ownership

The 64-byte Node contains a 32-bit packed scalar word, a 32-bit source position,
four ordinary `*Node` links, a tracked type/color pointer, a property-list pointer,
and a tracked payload pointer. All pointers remain real Go GC roots. Pointer bits
are never hidden in `uintptr`, metadata bits, numeric values, or integer IDs.

The scalar word contains an 8-bit common token, 3-bit folded subclass, type/color
bit, promotion bit, and 19-bit common source length. Token 255 and length 524287
are reserved sentinels. Negative/arbitrary int32 tokens and negative/large int32
lengths are retained exactly in promoted metadata. `SetLength(int)` still first
truncates to Java int32, as the original implementation does. Current enum token
predicates compare the low byte directly; a test requires every real enum token
to remain below the sentinel.

| Kind | Common payload | With original name or overflow |
| --- | --- | --- |
| Base | nil, or exact original-name string cell | 16-byte cold metadata |
| Number | 8-byte pointer-free float allocation | 24-byte metadata + float |
| BigInt | direct immutable `*big.Int` | 24-byte metadata + pointer |
| String | direct immutable interned string cell | 24-byte metadata + cell |
| Template substring | 32-byte metadata/raw/cooked record | same record |

The 8-byte float uses Go's normal tiny allocator, which may share a 16-byte tiny
allocation with another scalar. It introduces no compiler-owned lifetime group.
The baseline allocates a 32-byte cold record for a numeric value; this candidate
cuts that record and removes its pointer scanning in the common case. Promotion
is one-way for simplicity, so changing back to a common length/token does not
immediately reclaim the cold record. Exact heap allocation and GC cost still
need measurement, including promotion frequency.

Every mutable promoted record is owned by one node. Clone constructors allocate
independent numeric and metadata storage; strings and immutable BigInt values
retain the original sharing semantics. `Node_CopyBaseNodeFields` uses setters so
source and destination subclasses may differ. Constructors migrate arbitrary
NewNode token overflows when setting the folded subclass. Nodes continue to be
independent `&Node{}` allocations with ordinary GC lifetime.

## Why pointer-free adjacency is a larger compatibility boundary

Replacing four links with integer IDs requires an ID-to-`*Node` resolver. A
strong compiler-owned node table roots all nodes, including detached and dead
rewritten subtrees, until the table dies. A small escaped `*Node` would require
an owner pointer to keep the resolver alive and could retain the whole compiler
or AST. Global tables additionally introduce lifetime/concurrency problems.
Weak tables cannot guarantee that an adjacency target is alive when only an
integer references it. Reconstructing pointers from stored integers is not safe
under Go GC and checkptr. An extra `*Node` adjacency shadow would duplicate the
graph and defeat the stated constraint.

A pointer-free, per-compilation arena is attractive for an internal immutable
IR with explicit compiler lifetime and integer-handle API. It is not a drop-in
representation of this mutable, publicly exposed Node graph: passes detach,
reparent, replace, and retain arbitrary subtrees, and maps use Node identity.
Obtaining that design needs a deliberately separate pass boundary, explicit
conversion/lifetime semantics, and a measured payoff that exceeds conversion
and duplicate-live-IR costs. This prototype instead retains all four pointers.

## Validation required

Source fixtures exercise every folded subclass, cooked/uncookable templates,
full int32 metadata boundaries, transition back to compact metadata, constructor
promotion, exact original-name pointer identity, renaming, clone mutation
isolation, and GC reachability through tagged payloads after clearing the string
pool. The existing Rhino and parser/oracle suites remain the larger oracle.

Before retaining this candidate: run Rhino tests and checkptr, parser oracle and
compiler semantic gates, then compare exact JS/maps/diagnostics across ordinary
CLI tasks. Run paired baseline/candidate timing with the same PGO treatment,
thread count, GOMAXPROCS, input, and Java contract. Report wall time, peak RSS,
and Java/Go ratios separately. Include number-heavy and original-name-heavy
held-out inputs; reduced node size may lose to payload branches on small inputs.
A 64-byte node still requires all pointer chases, so neither 6.9% GetNext nor 21%
GC is a promised speedup or an additive ceiling for this change.

## Isolated private-parser allocation follow-up

The new follow-up branch starts from `eecbcf44`, preserving the current private
per-attempt parser allocator and reference-batch experiment. It reapplies only
this layout/payload change, not QName APIs, tiny flow maps, or global node
allocation chunks. The existing `NodeAllocator` uses the same packed-token and
per-kind setters as the ordinary constructors; it still retains only its current
8-to-128-node slab and does not pool or recycle published nodes.

This combination was not isolated in the previous mixed experiment. A full
128-node slab has 8192 bytes of node backing storage instead of 10240, before
payloads and allocator size classes. The maximum neighboring-node retention
remains 127; smaller backing storage does not remove those neighboring nodes'
outgoing links. String payloads retain their direct pointer; common numeric
payloads shrink from a 32-byte extension to an 8-byte pointer-free float.

Potential gains are fewer node bytes scanned and better cache/slab density.
Packed-field branches and cold promotion can offset them. This source branch
requires fresh normal/race/checkptr, exact CLI output gates, and paired whole
compile time/RSS/allocation measurements; no speed or memory win is established.
