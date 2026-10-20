# Inline the two common BitSet words

This isolated binary alternative replaces the internal jutil BitSet
representation with 32 bytes: two inline uint64 words, a pointer to a larger
slice holder, and one byte each for logical used length and reserved capacity.
It adds no runtime flag, build tag, global cache or unsafe operation. Mainline
remains unchanged until validation and a matched whole-CLI comparison select a
representation. Do not combine this experiment with the separate native-liveness
V2 owned-word publication helper while attributing its effect.

## Capacity and ownership audit

The port's existing NewBitSet/zero value has Size()==0, despite the old Java
comment saying 64. NewBitSetWithSize preserves exact requested word capacity,
including 129 bits -> three reserved words even while empty. Clone copies only
used words and adopts Go append's capacity, so an empty reserved clone loses its
reservation. These established port behaviors are retained.

Small storage tracks logical capacity 0/1/2 separately from its physical two-word
array. ClearAll/trim retain that capacity. Extending used length clears newly
exposed words, matching the original append of zero words and preventing old
cleared bits from resurfacing. Promotion retains Go's exact old small-slice
capacity behavior: the growth target is max(newLen,2*oldCap), then a real uint64
append from nil obtains the runtime's rounded allocation capacity. A plain make
would incorrectly expose a smaller Size for some lengths. Larger storage uses
the original append/trim engine and is never demoted merely by clearing bits;
Clone may create a new small value because it already discarded unused capacity.
Large constructor capacities remain int; overflow errors are not truncated into
the small-capacity byte.

All production uses outside jutil hold *BitSet pointers, except one initializer
in nativeLiveVariablesSolver.publish: each array cell receives the value of a
unique NewBitSetWithSize temporary whose source pointer is discarded. There is
no surviving source alias there. Repository searches found no shared-backing
struct-copy dependence, external word-field access, BitSet literal initializer
outside jutil, or BitSet-specific reflection/unsafe-size dependency. Java-style
pointer copies still alias; Clone remains independent. Go copying the struct
by value is not a supported replacement for Clone: its incidental old slice
aliasing is intentionally not retained, as approved for this internal package.

## Method audit

| API/method family | Required behavior retained |
| --- | --- |
| zero/NewBitSet/NewBitSetWithSize/Size | logical capacity, negative/overflow constructor errors |
| ensure/trim/ClearAll | zero extension; used length versus retained capacity |
| Set/SetTo/SetRange | index/range validation and empty-range nil-receiver short circuit |
| Clear/Get | bounds behavior without growth; trim after clear |
| Flip/FlipRange | zero extension, complement, trim and empty-range behavior |
| Length/IsEmpty/Cardinality | original used-word invariant and bit count |
| NextSetBit/PreviousSetBit/NextClearBit | masks, word-boundary results, -1/negative cases |
| Or/And/AndNot/Intersects | original operand and self handling; empty receivers must not dereference nil operands for And/AndNot/Intersects |
| Clone/Equals | deep ownership, used-word capacity clone behavior, typed-nil and wrong-type equality |
| HashCode/String/ToByteArray | original signed Java hash, encounter order and independent little-endian bytes |

The production New/Clone and all 24 exported instance-method declarations were
reviewed together; the representation adds only private word-access helpers.

## Allocation tradeoff, not a performance claim

On amd64 the current header is 24 bytes. The larger holder adds a 24-byte slice
header only when capacity exceeds two words.

| Reservation/storage | Current objects / bytes excluding rounding beyond two words | Inline alternative |
| --- | --- | --- |
| capacity 0 | 1 / 24 B | 1 / 32 B (+8 B) |
| capacity 1 | 2 / 32 B | 1 / 32 B |
| capacity 2 | 2 / 40 B | 1 / 32 B (-8 B) |
| capacity >2 | 2 / (24 B + backing) | 3 / (56 B + same backing), +32 B |

Large allocation and zero-capacity clone frequency matter. Stock liveness has a
100-variable analysis limit, so its reservations fit two words; other consumers
include subtype/interference graphs, instrumentation lines, transitive chunk
sets, and template metadata and can grow much larger. Native liveness publishes
arrays of BitSet headers, so replacing per-cell backing allocations has a
different object model from separately escaped constructor calls. Aggregate
method-call counts cannot establish actual malloc bytes or RSS. Measure real
TotalAlloc/Mallocs, peak RSS and elapsed time under default GC, without tuning.

## Validation plan and current status

Source-only: formatting/whitespace checks pass; no Go build/test or runtime
measurement has run yet. Tests include frozen pre-redesign slice code from
7e23d53b as an independent oracle for Go-specific capacities, randomized mixed
operations, direct/sequential promotions, constructor overflow errors, trim and
stale-capacity extension, Clone/reference/array-cell ownership, nil short
circuits and allocation-object assertions. Existing JDK bitset/hash/flip/
previous-bit vectors remain required, plus complete compiler tests and CLI byte
gates. No speed or memory improvement is claimed from this source prototype.
