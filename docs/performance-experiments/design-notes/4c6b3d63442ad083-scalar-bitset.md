# Scalar-first BitSet experiment

This isolated candidate starts at current-main research base `da20b44f` and
changes no compiler options or process settings. It retains BitSet's 24-byte
header with one inline uint64, a GC-tracked pointer to a boxed large slice, and
an integer reserved inline capacity. There is no production unsafe code.

The fresh bundle allocation profile attributes approximately 54.5 MiB and
2,686,974 sampled allocation objects to NewBitSetWithSize: 1,769,469 24-byte
headers and 917,505 16-byte allocations/tiny blocks. All sampled backing entries
are the 16-byte bucket. That bucket combines one-word tiny allocations and real
two-word backing arrays; it does not prove the exact <=64-bit fraction.

The sampled callers are overwhelmingly LiveVariablesAnalysis lattice creation
and FlowThrough (gen/kill). Coalesce rejects functions with 100 or more local
variables, so its allocations are at most two words. DeadAssignments has a
function/body scope-count check before collecting all nested locals; it is not
a universal bound on all LVA sizes. A runtime request-size histogram remains a
useful untimed follow-up before claiming a macro improvement.

For zero/one-word capacity, construction allocates only the BitSet header.
Logical inline length follows whether the scalar is nonzero; all public methods
maintain the existing trimmed-high-word invariant. Set/Flip and scalar queries
operate directly on the word. Large sets promote once to a boxed slice and keep
that representation, even when cleared, because Size exposes reserved capacity.
The promotion uses the exact former append expression and the inline slice's
existing len/cap, retaining Go's capacity growth and allocation rounding. Large
sets incur an additional 24-byte slice-header object, so the candidate can lose
if two-or-more-word sets are frequent.

Public method behavior, logical words, equality/hash, iteration, byte output,
zero value, exceptions, and current capacity behavior are preserved. In
particular, the shipped NewBitSet().Size() is zero despite a comment mentioning
Java's default64. Clone drops reserved capacity exactly as the former
append(nil, words...) did; an empty clone has size0 and a one-word clone has
size64. ClearAll retains capacity and later growth clears newly logical words,
preventing stale backing values from resurfacing.

One explicit compatibility limit: arbitrary Go struct-value copies no longer
share one-word backing storage, and promoted copies share a boxed slice header.
Repository consumers use *BitSet and Clone exclusively, matching Java reference
semantics. Preserving arbitrary raw struct-copy aliasing would require the
backing allocation being removed here. This candidate should not ship if raw
struct copying is considered a supported external contract.

Validation includes the existing Java-derived BitSet TSV fixtures and all jutil
unit tests, plus a frozen copy of the original slice implementation as a test
oracle. Differential tests exercise 18,000 mixed operations across constructor,
scalar/promotion, clearing, clone, binary-operation, capacity-growth, query,
exception and nil boundaries. A separate allocation assertion proves one-word
constructor+Clone allocate only two headers, with no backing arrays, and checks
the header remains24 bytes. Full jutil tests pass; the race run and macro CLI
measurements are reported separately by the central runner.

No wall-time or peak-RSS improvement is claimed here. Measure the same ordinary
CLI workloads, output/maps/diagnostics, thread/runtime settings, and PGO treatment
as the baseline, and report both time and peak RSS. Keep this separate from the
frozen AST64 candidate until measuring the combination.

## Optional untimed histogram build

A separate diagnostic commit adds the `jsgo_bitset_diagnostic` build tag. The
ordinary build selects a compile-time false constant: recording branches and
exit output are dead code, with no counters, initialization, environment
lookups, or diagnostic imports. Do not use the tagged binary for timings: it
adds atomic event increments.

Build the same CLI with `go build -tags=jsgo_bitset_diagnostic ./cmd/jsgo` and run
an untimed workload. The normal System_Exit path writes one JSON object to stderr
before exiting; successful CLI runs use that path. A Go panic caught by main's
direct os.Exit does not promise a report. No new public CLI option is added.

The report buckets counts by words: 0,1,2,3-4,5-8,9-16,>16. It records sized
constructor requested capacity, default constructors, Clone logical word count,
first inline-word reservation after zero capacity, promotions by old/new
capacity, and actual large backing-capacity growth. These are allocation-boundary
events rather than live-size or per-object-lifetime counts. No per-object IDs or
ownership map is introduced.

`avoided_backing_allocation_events` sums one-word sized constructors, one-word
clones, and first zero-to-one inline reservations. `added_large_header_events`
sums large sized constructors, large clones, and promotions. These describe
allocation calls affected by the representation; tiny-allocation accounting and
bytes require separate interpretation. Neither event count measures runtime,
peak memory, or Java/Go ratios. Compare diagnostic results across the actual
SIMPLE/ADVANCED, coverage, and multi-chunk workloads before extrapolating.

## Current-main integration source review

The central `bd306adb` baseline has unchanged jutil BitSet source relative to
`da20b44f`, so the storage candidate has no source conflict. All production
consumers found use *BitSet or []*BitSet; value copies were not found. The
explicit Clone consumers are LVA lattice copying and parameter removal masks.

Default zero-capacity sets grow in parameter masks, Coalesce interference rows,
chunk dependency graphs, coverage line masks, and property/color relations.
Coalesce's local-variable bound usually keeps its interference rows short;
chunk indices can exceed63, coverage line numbers frequently do, and global
color/type IDs can be much larger. CrossChunkCodeMotion preallocates the number
of chunks and preserves it via Size. These are real large-set cases where the
extra boxed header and lookup indirection may offset LVA savings. Exact request
and mutation histograms are needed; allocation microtests are not macro proof.
