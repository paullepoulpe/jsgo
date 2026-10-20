# One inline BitSet word at the original 24-byte size

This isolated binary representation alternative starts from the validated
two-word prototype `a6afc95f`. It stores one inline uint64, a pointer to an
optional slice holder, and logical used/capacity bytes. The amd64 BitSet is
24 bytes, matching the original slice header; no runtime flag, unsafe storage,
global state, GC tuning or public Java API change is introduced. The older
`inline-bitset.md` describes the separate two-word alternative, not this layout.

Only the inline threshold changes from two words to one. All word operations,
nil/index/range behavior, hash/equality, zero-extension after trim/ClearAll and
logical Size behavior retain the earlier implementation and frozen slice oracle.
Reservations above64 bits immediately use the original slice engine through the
holder. Promotion from capacity0/1 to two or more words uses the same runtime
append-based capacity rounding as the original. Large storage never demotes on
trim: observable Size retains its old reservation. Clone copies only used words,
so a large empty reservation clones to capacity0 and a one-used-word large set
clones to inline capacity1, exactly matching the old capacity behavior.

Java pointer copies still alias and Clone is independent. The previously
approved internal Go struct-copy limitation remains: repository production has
only one by-value initializer, copying a unique constructor temporary into a
native liveness publication array with no surviving source alias. This branch
does not contain the separate V2 owned-word publication helper.

## Allocation model, before measurements

Derived amd64 payload/class costs for individually escaped constructors:

| Word reservation | Original slice BitSet | Two-word alternative | One-word alternative |
| --- | ---: | ---: | ---: |
| 0 | 1 object /24 B | 1 /32 B | 1 /24 B |
| 1 (1–64 bits) | 2 /32 B | 1 /32 B | 1 /24 B |
| 2 (65–128 bits) | 2 /40 B | 1 /32 B | 3 /64 B |
| >2 | 2 /(24 B + backing) | 3 /(56 B + backing) | 3 /(48 B + backing) |

One word removes the empty-set growth penalty and saves8 bytes for common
single-word sets, but **adds24 bytes and one object for every two-or-more-word
set** versus the original. Compared with two inline words, two-word reservations
are32 bytes/two objects more expensive. Header arrays and escape elimination
have different object accounting, so constructor counts cannot establish actual
malloc bytes, live heap or RSS. Promotions also allocate a holder and replacement
backing while retaining the original header; their real distribution matters.

The original-implementation diagnostic is prepared separately at
`/workspace/jsgo-results/round2/bitset-distribution-overlay/native-analysis-f2c17056/overlay.json`.
It records exact constructor capacities, Clone source-capacity/used/result
triples and actual capacity-growth transitions. NewBitSetFromOwnedWords is a
separate row to detect accidental V2 publication. The diagnostic binary must not
be timed: its numeric maps/atomics add overhead, but retain no original sets or
compiler graph references. No inventory has run yet, so this proposal has no
demonstrated distribution or whole-CLI win.

## Validation state

SOURCE ONLY, unbuilt and untested. Gofmt and whitespace/diff checks pass.
Existing frozen-port differential/random-operation, capacity/promotion, nil,
ownership and JDK vectors are unchanged. Allocation assertions now cover empty,
1/64/65/128/129/1025-bit constructors and empty/one/two/three-used-word clones,
including the expected larger fallback cost. These assertions still require
execution under normal and race builds, then full compiler/CLI differential
gates and matched default-GC time, TotalAlloc/Mallocs and RSS measurements.

Native liveness currently reserves128 bits for its two private gen/kill scratch
sets regardless of actual variable count. Reducing that reservation may avoid
unnecessary fallback in this layout, but must remain a **separate commit and
factor**. It is not included in the initial one-word candidate.
