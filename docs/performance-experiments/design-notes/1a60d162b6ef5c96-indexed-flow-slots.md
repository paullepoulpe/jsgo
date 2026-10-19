# Indexed inferred-type slots (experimental)

`JSGO_EXPERIMENT_INDEXED_FLOW_SLOTS=1` selects an alternative persistent storage
backend for `LinkedFlowScope` inferred slots. The default remains the original
HAMT. This is a prototype, with no measured compiler performance improvement yet.

The type-inference rules, canonical `TypedScope`/`TypedVar` instances, declarations,
qualified-name handling, scope trimming and type-equality operations stay in the
existing compiler. Only each scope's inferred-name map changes representation:

* A scope owns stable symbol IDs, shared across its immutable flow versions.
* Four slot pointers occupy each compact page. An update copies one page and its
  directory, preserving the other pages and every original slot pointer.
* Joins address entries by integer ID and invoke type joins in canonical HAMT hash
  order. Equality likewise preserves the order of actual changed-type comparisons;
  differing key sets delegate to the original HAMT's early-exit traversal.
* More than 256 distinct symbols, cross-owner maps, or a Java string-hash collision
  use the original HAMT behind the same immutable wrapper. This bounds copying
  and preserves collision-sensitive behavior of the reference implementation.

The owner lives with its typed scope; it is not a global cache. Symbol interning
is synchronized. Owner initialization remains in the existing serial typed-scope
construction phase. This prototype does not parallelize type inference.

`internal/flowstore` can be tested without compiling the large `jscomp` test binary.
Its differential oracle is the unchanged Closure HAMT, including historical
versions, branches, deletions, joins, callback order, hash collisions, mixed
owners, overflow, and concurrent symbol interning. The compiler integration tests
also repeat the original flow-lattice suite and compare typed ASTs, source, and
ordered diagnostics with the reference backend for forward aliases, recursive
object types, unknowns, block shadowing, collisions, and large scopes.

A changed map representation is not automatically a performance win. Stable-ID
lookup introduces a scope-local Go map and synchronization, and directory copying
can lose against HAMT path copying for some densities. The disabled adapter also
adds an owner pointer to `TypedScope`. Benchmark enabled versus disabled builds
and the unchanged compound binary, measuring time and RSS, before retaining it.
Profile context: in the existing `current-profile/checks-only` capture, stacks
containing type equality, unions or flow-scope/HAMT operations accounted for about
6.75% of sampled CPU and 70.5 MiB (4.48%) of total allocation. These overlapping
stack fractions are prioritization evidence, not additive savings or a speedup
prediction. This change alone should not be expected to close the full Java gap.

## Screening decision and source allocation audit

Reject this representation for retention. The single-run ordinary-CLI screen
matched Java's output, but bundle time was 18.743 s / 578.99 MiB with indexed
slots versus 18.668 s / 507.16 MiB with the backend disabled: no speed benefit
and 14.2% more peak RSS. Typed-app was 5.298 s / 270.04 MiB versus 5.395 s /
273.04 MiB. One sample does not establish a repeatable gain. Raw evidence is
`ordinary-cli-screen/summary.json`; these are screening observations, not final
benchmark claims.

The following sizes follow the 64-bit field layouts (before allocator rounding):

| Storage object | Bytes |
| --- | ---: |
| `Store[*OverlaySlot]` header | 48 |
| Four-slot pointer page | 32 |
| String-keyed `HamtPMap` node | 64 |
| `OverlaySlot`, common to both backends | 32 |
| ID-owner header, excluding maps and backing arrays | 88 |

A sparse version that updates only ID 255 allocates a 512-byte directory plus a
32-byte page and 48-byte header: 592 bytes of storage, compared with a 64-byte
one-entry HAMT node. Including the common slot allocation gives 624 versus 96
bytes. This is a worst-case shape, not a measured occupancy claim.

Other concrete amplification paths are present in the implementation:

* `Equivalent` reconstructs two complete HAMTs whenever native key sets differ,
  discarding existing structural sharing for a common fixpoint comparison.
* `Reconcile` allocates a new directory before finding that its answer is an
  unchanged input. Those discarded directories still increase allocation.
* Overflow and collision fallbacks allocate an additional 48-byte wrapper for
  each changed HAMT version, providing no native storage benefit for that update.
* ID maps, names and sorted order remain attached to canonical typed scopes after
  a function's temporary flow states can die. Publishing a new immutable order
  for every inserted ID allocates 263,168 bytes cumulatively at 256 IDs, although
  only the final 2,048-byte order stays live. Two Go maps and the name array add
  further owner storage. Their aggregate contribution has not been profiled.

The `TypedScope` field itself grows its layout from 72 to 80 bytes; both fit the
same 80-byte Go allocation class on this platform. Its disabled-state field is
therefore not an explanation for the screen's RSS variation.

These mechanisms explain why compact pages do not guarantee compact compilation.
They do not prove which mechanism accounts for the observed RSS increase. No
occupancy counters, allocation profile, or repeated performance result has been
collected for this candidate. Keeping the original HAMT is the lower-complexity
choice; further production fixes are deferred rather than adding more controls.
