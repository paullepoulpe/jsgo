# One inline BitSet word: source-only alternative assessment

Status: held for a measured capacity distribution, neither rejected nor selected.
No implementation, diagnostic run, build, test, profile processing, GC change or
timing was performed for this assessment. Source reviewed: isolated two-word
BitSet (`/workspace/jsgo-inline-bitset`, initially 39c77a25), the ordinary
constructors/Clone call sites, native liveness publication/scratch construction,
and the existing diagnostic-overlay README.

## Why consider a different threshold

Root's matched screen reported the two-word 32-byte representation as mixed:

| Workload | Slice-control time / peak RSS | Two-word time / peak RSS |
| --- | --- | --- |
| typed-app | 5.081 s / 290 MiB | 4.954 s / 250.7 MiB |
| bundle | 18.079 s / 417 MiB | 18.449 s / 432.5 MiB |

Those observations do not identify empty/one/two/large-set frequencies or prove
the cause of the bundle result. A one-word representation could remove the
two-word design's extra eight header bytes while retaining inline storage for
the smallest nonempty sets. It also creates an important new two-word penalty.

## Feasible layout and allocation tradeoff

On amd64, `uint64 inline`, `*largeBitSetWords`, and byte-sized used/capacity fields
fit a 24-byte struct after padding. No union, unsafe pointer, tagged integer,
global pool or runtime switch is required. The large holder remains a 24-byte
ordinary slice header, owning the same uint64 backing/growth engine as the old
implementation.

The following counts describe separately escaping headers and reserved backing
at construction; they are not observed runtime counts or complete RSS estimates.

| Reserved word capacity | Original slice representation | Two inline words | One inline word |
| --- | --- | --- | --- |
| 0 | 24 B, one object | 32 B, one object | 24 B, one object |
| 1 | 24 + 8 B, two objects | 32 B, one object | 24 B, one object |
| 2 | 24 + 16 B, two objects | 32 B, one object | 24 + 24 + 16 = 64 B, three objects |
| >2 | 24 B + backing, two objects | 56 B + backing, three objects | 48 B + backing, three objects |

Relative to two inline words, one word saves eight header bytes for capacity
0/1/>2, but costs 32 bytes and two objects for capacity 2. In this simplified
separately escaping-header model it reduces requested bytes only when
`N0 + N1 + Nlarge > 4*N2`, equivalently capacity-2 instances are less than 20%
of the four categories. Relative to the original, one-word capacity-1 savings
must offset the 24-byte holder penalty on capacity-2 and larger instances:
`N1 > 3*(N2 + Nlarge)`.

These inequalities exclude allocator rounding, tiny-object packing, stack
allocation, promotion history, operations/inlining cost, GC work and liveness.
They are screening conditions, not performance predictions.

## What existing evidence establishes, and what it does not

The existing `next-infer-scope-live-alloc-profile.json` reports sampled allocated
MiB by stack family. It has no constructor argument, reserved-capacity, or Clone
used-length histogram. Sampled bytes or objects cannot be converted into the
required exact capacity frequencies. A BitSet allocation site services multiple
sizes, some header temporaries do not escape, and Clone can discard reservation.

The existing `/round2/bitset-distribution-overlay/README.md` describes exact
method-event counters and explicitly states that no diagnostic binary had been
built or run. No generated capacity-count result was found among the existing
distribution/count artifact filenames inspected. Actual frequency is therefore
unknown here; no numerical capacity distribution is asserted.

Source establishes these shapes and bounds, but not their workload frequency:

* `rhino.newTemplateTypeMapBitSet` reserves exactly 64 bits, including its empty
  sentinel. This is favorable to one word. Other template extensions reserve
  the number of keys; actual key-count distribution remains unknown. An empty
  extension can return the original map without retaining the constructor
  temporary, so constructor calls are not necessarily escaped headers.
* Stock liveness consumers normally cap analysis at about 100 variables. At
  65–100 variables, their initial/join reservations use two words. Transfer
  scratch capacity follows the input's observable Size, while Clone copies
  used words only; a two-word-reserved source may yield a zero- or one-word
  clone. Function-local count alone cannot classify all these instances.
* Native liveness currently reserves 128 bits for each of its two reused scratch
  BitSets, even in a function with at most 64 variables. A one-word layout would
  spill both. Changing that reservation to the variable count is a distinct
  possible optimization and must be controlled equally in any representation
  comparison rather than quietly bundled into one candidate.
* Native liveness V1 publishes an array of BitSet headers by copying unique
  constructor temporaries. Headers are array cells, not one heap allocation per
  constructor. V2's owned-word windows are another distinct ownership model.
  A general constructor histogram must label these cases before applying the
  simple object-count table. Array cells retain the same 24/32-byte stride cost.
* Coalescing starts interference rows with `NewBitSet()` and sets indices from
  variable order. A sparse set with only bit 95 still requires two words;
  cardinality is not a proxy for used/reserved word capacity.
* Subtype/color and transitive-chunk sets, line instrumentation, and parameter
  masks can grow beyond one/two words. Their size and promotion distributions
  are likewise not established by source. One word still pays an extra large
  holder compared with the original representation.

## Required semantic design if later measured

Retain the existing two-word implementation's separation of logical used length
from observable reserved capacity, reducing only the inline threshold to one.
Clear newly exposed inline storage after trim/ClearAll; promotion from capacity
0/1 must reproduce the old Go append capacity and size-class rounding, not use
an arbitrary exact-length allocation. Large storage stays large after clearing
because Size preserves reservation. Clone of a large source with zero/one used
words may produce an empty/inline value because original Clone already discards
unused reservation. Clone with two used words must spill and preserve the old
copy/growth capacity semantics. Preserve existing nil/short-circuit/overflow,
hash, byte-array and alias behavior and the unique-temporary native publication
copy. Existing independent old-port differential tests cover 65/128 boundaries
and should be reused with threshold-specific allocation expectations adjusted.

## Concrete next step, deferred to a future reserved diagnostic slot

Use the already prepared untimed original-slice diagnostic overlay on the same
five fresh CLI tasks and normal input/optimization/thread settings. Record:

1. Default/reserved/Clone method events by capacity 0/1/2/>2.
2. Clone source reservation **and** source used-word length (the existing matrix).
3. Actual capacity growth transitions, including capacity-1 to 2 and beyond.
4. Explicit identification of native publication's temporary-to-array copies
   and any borrowed/slab views; do not count their calls as escaping headers.

Keep that instrumented binary out of performance comparisons. Counters do not
retain BitSet pointers or change ownership and report package-initialization
events separately where feasible. They still measure method events, not exact
heap objects; bare zero values are unobserved. Pair the distribution with already
measured ordinary TotalAlloc/Mallocs/RSS rather than inferring these from sampled
pprof objects. If the two-word bucket is small enough to make the tradeoff
credible, build an isolated one-word binary and compare it against both original
slice and two-word controls, with default GC and the same correctness gates.

There is no basis yet for a blind one-word prototype or a negative performance
classification. Hold the idea pending that diagnostic evidence.
