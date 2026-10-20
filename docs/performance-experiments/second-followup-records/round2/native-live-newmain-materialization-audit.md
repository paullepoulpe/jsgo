# Native liveness scratch / publication audit after ff06ac15

Source-only. No new prototype, build, profile decode or timed run. Old negative
nativeLVA/owned-word/scratch measurements did not establish execution before the
post-colors admission correction; they cannot be treated as performance results
for active backends. Current source choices must be compared against the new
main, not the former a6 inline2 utility or original slice-only BitSet.

## Source conservation and superseded experiments

Old projection728de924 still has two fixed128-bit scratch sets and per-set-bit
native publication. 204015d1 (scratch reservation bounded by orderedVars) and
90905e69 (flat owned uint64 publication windows) are absent. The latter assumed
BitSet was a 24-byte slice header; a6 then used a32-byte inline2 representation.
New main ff06ac15 instead has words[] +buf[2]uint64, 40-byte BitSet (48-byte class),
NewBitSet capacity64, InitWithSize and capacity-preserving CloneInto. Its stock
LiveVariableLattice embeds that BitSet, totaling48 bytes, with liveSet pointing
at its own liveBits. Generic FlowThrough now stack-allocates gen/kill with
InitWithSize; these improvements substantially change the old control.

Scratch128→actual<=100 saves ZERO allocated bytes/objects with both a6 and ff:
all reservations already use the same inline buffer. Only unobserved logical
capacity changes. It can add dynamic arithmetic and another flag with no useful
memory result. Do not port it as a new production factor on this baseline.

Owned-word windows would require an extra uint64 backing slab, while each new
embedded BitSet already contains two words. They cannot remove a backing object
that no longer exists; they add bytes, pointer targets and slice aliases. Do not
blindly port the old helper or its Size/Clone assumptions.

## New-main integration hazard: self-referential value copies

The retained old native publish allocates states[], lattices[] and sets[] then
executes sets[i]=*NewBitSetWithSize(variableCount). New-main words points into
the temporary BitSet's OWN buf; copying that struct cannot retarget its slice
pointer to sets[i].buf. The constructor object must survive through that slice
pointer, while its copied buf and the new lattice.liveBits remain unused.
This is valid GC-tracked memory but can waste a heap object for every published
set in addition to duplicated slabs. Never benchmark this naive merge as the
intended new-main native backend.

The compatibility publication bridge should initialize each existing
lattices[i].liveBits in place using InitWithSize, assign liveSet=&liveBits, then
perform the existing nativeBits.write Set loop. Drop sets[] entirely. Keep
original annotation/state objects, original reservation, failure publication,
step counts and unknown/public fallback behavior. This matches the new stock
lattice constructor rather than creating a new representation or exposing an
owned-word API. Existing native query/differential/failure tests must re-run
against the NEW generic stock engine after integration.

## Event model, not RSS or fresh-main measurements

Existing corrected actualDoRun telemetry source9d190a2b shows bundle28,981 native
runs /428,006 nodes, hence57,962 scratch constructors versus856,012 published
sets. Typed2,279 runs/34,992 nodes; SIMPLE6,063 runs/69,658 nodes; checks/whitespace
run no optimization liveness. These are method/phase events, not allocations or
current-main counts. New stock passes may change node/work distributions.
Machine-readable model with exact paths/counts: native-live-materialization-model.json.

Per node, excluding private solver state/topology and allocator span metadata:

| Publication layout | State | Two lattices | Separate sets | Retained temps | Total |
| --- | ---: | ---: | ---: | ---: | ---: |
| Old a6 inline2 |24|16|64|0|104 B |
| Naive ff merge |24|96|80|96|296 B |
| Embedded ff bridge |24|96|0|0|120 B |

The temporary column assumes normal escaping self-buffer constructor objects,
48-byte class each. Go compile escape output should confirm this before claiming
actual allocation costs. Array entries use payload size, not per-entry class
rounding. At the old bundle node count, the model is42.45MiB /120.82MiB /48.98MiB.
These are cumulative publication byte models, not live bytes or peak RSS. The
proper ff bridge can still cost more bytes than a6's smaller utility; fair new
stock control tests/measurements are essential.

## A useful separate publication factor

After the embedded bridge, a private in-place two-word initializer could write
the two known native words directly into lattice.liveBits.buf, set words length
from the high nonzero word and retain the exact original reservation capacity.
It removes the per-populated-bit Set loop without extra slabs, headers, GC
pointer fields or aliases. Its inputs are bounded by the existing native
orderedVars<=128 admission. Normal public constructors/CloneInto stay unchanged.
The new flag must be captured once per builtin consumer Process and propagated
to existing parallel-coalesce jobs, not Getenv per node or per transfer.

Before claiming material gain, collect private untimed PublishedSets,
PublishedSetBits, ReservedWordCount and DirectPublishedSets counters. Conservation:
PublishedSets=2*Nodes for completed native runs, DirectPublishedSets equals all
published sets only with the factor on, and original Set loop operations equal
the sum of popcounts of all in/out states. Counter collection must include actual
parallel workers and retain only numeric data. Compare query bitsets, Size,
Length, Clone/CloneInto, malformed failure publication and complete CLI artifacts
against the new main generic/control path. No extra pointer/GC scanning benefit
is expected: this factor saves CPU instruction/loop work, not object bytes.

## Potentially larger private-state factor, inventory first

Current nativeLiveVariableState stores four [2]uint64 vectors plus counters/
flags:88 bytes per node, pointer-free. If a function has<=64 orderedVars, an
exact one-word version can store four uint64 values and identical counters/
flags in56 bytes, saving32 bytes/node (36%) and halving hot word operations.
Two-word fallback handles larger functions; original transfer/escaped-local
side effects, queue order, equality and annotation types must remain unchanged.
This is a genuinely adaptive representation threshold, not a workload shortcut.

Need weighted <=64-variable native-node/step mass, not just function counts,
before choosing it. Add private diagnostic variable-count bins and node/step
weights at the existing solver boundary, without retaining scopes or altering
admission. At the old bundle node count the absolute maximum private-state byte
reduction would be13.06MiB if every node qualified, not a prediction. This may be
more material than changing a no-allocation scratch reservation, but is a
separate design/validation factor, not authorized source implementation yet.
