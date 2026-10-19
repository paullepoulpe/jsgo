# Tiny inferred-slot maps

This candidate replaces only inferred-slot maps containing at most four names.
It uses a copy-on-write slice of the existing immutable slot pointers: a 24-byte
header plus 8 bytes per slot on amd64. Larger maps and Java string-hash collisions
return the original HAMT directly. There is no per-scope ID owner, fallback
wrapper, runtime switch, or scope-size-dependent directory.

Unique-hash keys stay in canonical HAMT hash order. Joins first check whether
the union remains eligible, then invoke callbacks in that order. A collision or
larger union delegates to the original HAMT before any callback runs, preserving
its collision pivot behavior. Equality compares changed slots in the same order;
unchanged pointers take the same no-type-operation path as the original identity
check. Mixed tiny/HAMT operations normalize only the tiny input, retaining the
original large map and its structural sharing. Unchanged join results return the
existing input map. Slot pointers, types, names and scope identities are retained.

Four is a hypothesis to test against actual compiler states, not an established
optimal cutoff. The `flowmapstats` build tag adds compiler-owned histograms for
inference results, join inputs/results and transitions into the HAMT. Exact sizes
0..63 have separate bins; the final bin counts sizes >=64. Rows count operation
observations, not distinct immutable versions. At the normal compiler report,
the instrumented binary writes one `JSGO_FLOWMAP_STATS {json}` line to stderr.
Ordinary builds use zero-size diagnostic state and inline no-op observers.
Instrumented runs are untimed and must not be used for performance comparisons.

The histograms must show how often maps stay tiny and how often repeated growth
rebuilds a HAMT before choosing the cutoff or retaining the implementation.
Transitions with result sizes >4 are capacity fallbacks; smaller results indicate
collision or key/value-name fallback. Compare normal prebuilt binaries using the
existing ordinary-CLI Java output gates, elapsed time and peak RSS. No measured
performance improvement or state-size distribution is claimed yet.

Independent tests use the unchanged HAMT as the oracle across persistent branch
histories, collisions, deletions, mixed backend orientations, callback arguments
and order, changed-slot comparison order, pointer sharing and independent
concurrent versions. Compiler lattice and full pipeline parity remain required.
