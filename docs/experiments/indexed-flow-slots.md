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
