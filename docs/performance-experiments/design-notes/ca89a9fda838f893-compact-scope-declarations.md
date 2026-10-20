# Compact scope declarations

This source candidate replaces only `AbstractScopeBase.vars`. Empty scopes still
hold a nil pointer. Declaration objects, declaration indices, implicit variables,
scope ancestry, and scope lifetimes stay unchanged. There is no cross-pass cache,
runtime switch, pool, slab, or chunk allocation.

## Representation and ownership

A nonempty map owns a dense vector of 24-byte entries on amd64. Each entry holds
the original declaration pointer, a pointer to its immutable name header, and two
32-bit insertion-order links. The map header is 48 bytes, compared with the old
32-byte linked-map header. A string-to-index Go map is created after eight live
entries, matching the old linked map's lookup threshold; eight is a hypothesis to
check against actual scope sizes, not a new measured optimum.

The private `getNameCell` method borrows the existing Var/TypedVar name field.
Source inspection finds that field written only during construction, and the
public Var name API exposes only a string value. A mismatching raw
`DeclareInternal(key, declaration)` key is preserved in its own string cell.
Replacing an entry changes its borrowed pointer so the old declaration is not
retained. Removal and clear zero both declaration and name pointers. Go's GC
tracks all pointers, including interior pointers; no unsafe or integer pointer
encoding is used. The construction invariant matters: reinitializing an existing
Var with `InitAbstractVarBase` would violate the borrowed-name assumption.

Stable integer links preserve live iteration across vector reallocation.
Removed entries retain their next link while active iterators skip them. Clear
disconnects the old chain from subsequent insertions. An atomic active-iterator
count prevents tombstone compaction while an iterator holds an index, and permits
concurrent read-only iterations. Mutations remain single-owner operations, as
with the original linked map. Compaction changes private storage indices only;
it never changes declaration objects or their declaration indices.

## Resource limits

Smaller retained storage does not guarantee lower cumulative allocation. For
ordinary sequential inserts, ignoring allocator metadata and using amd64 size
classes:

| Declarations | Original retained/cumulative bytes | Dense retained bytes | Dense cumulative bytes |
| --- | ---: | ---: | ---: |
| 1 | 80 | 72 | 72 |
| 2 | 128 | 96 | 120 |
| 3 | 176 | 144 | 216 |
| 4 | 224 | 144 | 216 |
| 5 | 272 | 240 | 408 |
| 8 | 416 | 240 | 408 |

The old entry has size 40 and is individually rounded to 48 bytes in tiny maps.
The dense table includes its 48-byte header and geometric backing-array growth.
Larger maps also have index storage; the original starts slab allocation after
eight entries, so extrapolating tiny-map savings to large maps is invalid.
Backing capacity and interior name-pointer loads may offset locality benefits.
Values iteration adds two atomic operations. No performance gain is claimed.

## Validation and distribution evidence

The standalone package tests compare actual linked-map behavior across random
histories, aliased keys, null values, live views, mutation during iteration,
reallocation, clear, nested iteration, tombstone compaction, concurrent readers,
and garbage collection of removed declarations. A compiler-package test checks
actual Scope declaration identity, declaration indices, and raw alias keys.

The optional `scopemapstats` **build tag** records newly-created untyped and typed
scope sizes, nil maps, indexed maps, vector length, and capacity in a
compiler-owned diagnostic state, then emits `JSGO_DECLARATION_STATS` JSON on
stderr. Ordinary builds contain no histogram storage or updates. Record this in
an untimed run of real CLI workloads, not a storage microbenchmark. Counts are
creation snapshots: later declarations/removals can change a scope, and a
subsequent pass may create a new scope for the same AST node. Histograms are
needed to assess the threshold and capacity slack before accepting the design.

At source freeze, tests, actual distribution counts, full compiler parity, and
ordinary CLI time/peak-RSS measurements are pending the central build slot.
