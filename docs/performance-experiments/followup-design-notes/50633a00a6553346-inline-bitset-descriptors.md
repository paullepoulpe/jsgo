# Two inline words with a GC-safe capacity descriptor,24-byte BitSet

Isolated binary alternative based on `a6afc95f`, separate from frozen one-word
`7e0ba5fb` and native scratch `204015d1`. No build tag, runtime flag, unsafe
storage, mutable global cache or GC tuning is introduced. The amd64 header is
two uint64 words plus one genuine typed pointer:24 bytes.

## Representation and invariants

The holder pointer has four states: nil means inline capacity0; addresses of two
immutable package-owned typed objects mean inline capacity1/2; any per-set holder
means the ordinary large slice engine. These are normal GC-tracked pointers, not
tagged address bits or uintptr payloads. The two descriptor objects have nonzero
size, distinct addresses and nil words. No method writes them or exposes them as
a mutable slice holder. They retain no AST, registry, compiler or caller data.

No per-set used/capacity bytes remain. Small used length is derived from the two
word values: high nonzero→2, otherwise low nonzero→1, otherwise0. Logical reserved
capacity survives independently in the descriptor, so Size still reports its
old reservation after ClearAll or clearing the last bit.

The original public API supports this invariant by induction:

* Zero/NewBitSet/NewBitSetWithSize initialize **zero used words**, even when the
  backing reservation is positive. In particular NewBitSetWithSize128 does not
  expose a positive len of zero words; IsEmpty is true and Length is0.
* Set and nonempty SetRange make the last used word nonzero.
* Clear/Flip/FlipRange/And/AndNot trim trailing zero words; ClearAll uses len0.
* Or preserves a nonzero last word because both operands satisfy the invariant;
  an empty operand extends nothing. Clone copies only used trimmed words.
* Private ensure can temporarily expose zero words, but every public caller
  immediately sets/ORs/flips and trims when cancellation can clear the last word.
  No callback observes the transient state. NewBitSetFromOwnedWords, in its
  separate alternative branch, explicitly trims too; it is absent here.

Therefore deriving used length preserves Length/IsEmpty/query/clone behavior for
all states reachable through the public API. Deliberately malformed internal Go
struct literals or a private ensure call inspected before mutation were not
public states and have no repository production dependency.

ClearAll must physically zero both inline words because no separate used length
can hide stale bits. Trim is a no-op on small storage; the bits already imply
the trimmed length. Ensure returns its temporary extended slice directly,
rather than calling usedWords before its caller writes the new bit. Promotion
retains original append-based capacity doubling/rounding; original reserved
small capacity0/1/2 remains independent of physical two-word storage. Large
storage is never demoted merely by trim/clear. Clone of at most two used words
chooses a descriptor from that used count, preserving original clone capacity.

Java pointer copies and Clone retain their contracts. The previously approved
internal Go struct-copy limitation remains; repository's only production copy
initializes an array cell from a unique discarded constructor temporary.

## Costs and evidence

| Reservation | Original slice | Two-word32B | One-word24B | Descriptor two-word24B |
| --- | ---: | ---: | ---: | ---: |
| 0 | 1obj/24 B | 1/32 B | 1/24 B | 1/24 B |
| 1 | 2/32 B | 1/32 B | 1/24 B | 1/24 B |
| 2 | 2/40 B | 1/32 B | 3/64 B | 1/24 B |
| >2 | 2/(24+backing) | 3/(56+backing) | 3/(48+backing) | 3/(48+backing) |

These are derived individually escaped allocation models, not heap/RSS results.
Two immutable24-byte descriptors are shared process-wide. The main CPU tradeoff
is extra pointer-state classification and deriving small used length rather
than reading a byte. Fewer objects/bytes need not imply faster complete builds.

The original-slice inventory at
`/workspace/jsgo-results/round2/bitset-distribution-original/report.md` reproduces
all five Java artifact gates. It shows predominantly one-word result capacities,
but39,779 two-word constructor events in bundle and small→large promotions in
typed. The descriptor proposal avoids the one-word candidate's holder/backing
for those two-word states; it does not eliminate large fallback costs.

## Validation status

Standalone targeted jutil normal/race tests PASS with Go1.27.2 and default GC
(GOMAXPROCS=1/-p=1 constrain test resources only). Normal execution0.019 seconds,
race execution1.105 seconds. Exact commands, source and log paths are saved in
`/workspace/jsgo-results/round2/inline-descriptor-jutil-provenance.json`.
Gofmt and diff checks pass. Frozen original oracle, randomized mixed operations,
exact append-capacity/constructor/clone/nil tests and JDK vectors pass. New tests
explicitly exercise zero-used reservations,
OR with empty reserved sets, set-clear/double-flip/range/self-AndNot cancellation,
ClearAll followed by low-word extension, independent clones and concurrent
descriptor users, promotion capacity and never mutating shared descriptors.
Allocation assertions confirm0/1/64/65/128/129/1025-bit constructors including the
large fallback penalty, and prove the intended24-byte header on amd64.

Independent java_worker source review found no blocker in the public invariant,
temporary ensure slice, capacity/growth, descriptor mutability or GC pointers.
Compiler/CLI differential gates and matched default-GC elapsed/TotalAlloc/Mallocs/
RSS comparisons remain required. No speed, memory-ratio or Java-parity
improvement is claimed.
