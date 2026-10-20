# Compact Node common-header access

This isolated followup to the 64-byte Node prototype removes an avoidable
kind-dispatch switch from cold metadata reads. It changes no field sizes,
allocation ownership, values, pointer identity or public API.

The existing Node64 CLI screen reduced peak RSS but slowed typed, bundle and
SIMPLE workloads. Inspection of the already-built matched balanced-PGO binaries
found an obvious source-level contributor: Node80 GetToken is a three-byte
load/return, while Node64 GetToken occupied about 90 bytes. Most extra code was
the five-way coldMeta kind switch, inlined through packedToken's overflow path.
The same switch can expand other frequently inlined metadata accessors. This
is evidence of code expansion, not a measurement of its runtime share.

Every promoted concrete payload starts with node_coldMeta at offset zero. The
header accessor now casts the existing GC-tracked payload pointer directly to
that first field. Concrete allocation types and their GC pointer bitmaps remain
unchanged. There are no integer pointers, low-bit tags, arenas, address arithmetic
or new lifetime owners. Existing arbitrary int32 token/length overflow behavior
and per-kind payload promotion remain intact.

GetString was also inspected in both existing binaries. The compiler already
eliminates the impossible template-kind branch after its StringNode cast. The
remaining extended-string branch/pointer indirection is inherent in the current
compact representation. Type/color fields and property lookup paths showed no
new hash lookup or property scan caused by the layout.

The existing payload differential tests exercise every kind, overflow tokens,
length transitions, originals, clone independence and GC roots. Their layout
assertion now additionally requires each promoted header to remain at offset
zero. This source-only candidate still needs normal/race/checkptr correctness
gates and a new matched CLI screen; no speed improvement is claimed yet.
