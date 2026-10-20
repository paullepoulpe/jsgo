# Lazy second-word reservation: source design only

2026-10-19. No implementation, builds, tests, profiles or timing. The current
one-word source remains frozen at7e0ba5fb; scratch204015d1 is a separate factor.

## Bounded one-word extension

The one-word24B BitSet can represent logical capacity2 while storing only a
single used inline word. Keep the same existing fields and allow capacity0/1/2
with used0/1. `usedWords` must cap its physical slice at `min(capacity,1)`;
observable `Size()` still reads logical capacity and reports128 for cap2.
All public operations consume used words, not the private slice's capacity.

NewBitSetWithSize65..128 would create only the24B object, no holder/backing.
Set/Flip bits0..63 remain inline. Clone copies the used word count as before,
so an empty reserved128 set clones to cap0 and a low-only set clones to cap1.
ClearAll and trim retain logical capacity2. They must continue clearing any
newly exposed inline word after the used length was reset.

When a write needs two or more used words, promotion has two distinct cases:

* Requested length≤logical capacity: allocate a holder with `make(words,length,
  logicalCapacity)`. In particular, first high-word write to reservation2 must
  remain cap2, not double to4.
* Requested length>logical capacity: use the already-reviewed exact Go-append
  growth and allocator-rounding policy. Capacity2→length3 must become cap4;
  direct larger lengths must match the frozen original oracle too.

Never demote a promoted large set on trim, because its reserved Size remains
observable. All nil/index/range short-circuit order and internal Go copy-by-value
limitations remain the same as the one-word prototype.

Derived constructor layout for cap2 then becomes1object24B versus original
2objects40B and frozen one-word3objects64B. If the set later needs its high word,
it still ends up allocating a holder/backing: aggregate allocation then matches
the frozen one-word model, merely later. Thus zero/low-only reservation frequency
is the deciding evidence, not constructor capacity alone.

Existing original-slice inventory observes39,779 cap2 constructor results in
bundle, none in the other four tasks. Clone events from a cap2 source use0/1/2
words681/1,662/3,704 times in bundle. These are observations at clone sites, not
a census of all39,779 sets' lifetime maxima. They establish that low/empty cap2
states exist, but cannot quantify total avoided backing allocations. Whole-CLI
time and real allocation/RSS remain independent gates.

Required tests would compare frozen-slice semantics and allocation transitions:
reserved128 empty and low-only stay inline; Get/Size/Length before high writes;
first bit64/127 uses cap2 exactly; direct bit128 uses cap4; ClearAll/stale bits;
empty/low/high Clone capacity; OR/ranges crossing63/64; and subsequent large
growth. No runtime flag or additional production counters are required.

## Stronger separate alternative: two words in24B

There is a GC-safe way to store two inline words while retaining the original
24-byte header: `[2]uint64` plus the existing genuine holder pointer. That
pointer can select four representation states:

| Pointer value | Representation / capacity |
| --- | --- |
| nil | inline capacity0 |
| address of immutable descriptor1 | inline capacity1 |
| address of immutable descriptor2 | inline capacity2 |
| any real per-set holder | ordinary large slice engine |

Descriptors would be two immutable package-owned typed objects, not pointer
bit tags, uintptr payloads, mutable global caches or per-set allocations. All
references remain ordinary GC-tracked Go pointers. The inline used length is
derived from the two words: high!=0→2, otherwise low!=0→1, otherwise0. Consequently
ClearAll must zero both words, and trim needs no separate small used-length
mutation. Logical capacity survives through the descriptor selector.

This gives one24B object for every≤128-bit set, with zero/one/two-word logical
capacity preserved, while>128-bit storage uses the24B holder fallback. It removes
the lazy design's eventual high-word promotion cost as well. Clone chooses the
descriptor from its used-word count, preserving exact cap0/1/2 clone behavior.

The cost is extra pointer-state classification and deriving small used length
on reads. Those operations may offset allocation/locality gains in hot BitSet
loops, so this is a separate representation alternative requiring measurement.
Private methods must never treat an immutable descriptor as a mutable holder.
Raw internal Go struct value-copy alias behavior remains the same previously
approved limitation; Java pointer/reference copies and Clone retain contracts.

This design is proposed only, not authorized as a production implementation.
It would reuse the same frozen original oracle and full semantic/CLI gates,
and would need explicit tests proving descriptor immutability across independent
sets and concurrent owners. Do not combine it with V2 owned-word publication
until independent effects are known.
