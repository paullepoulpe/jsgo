# Private native typed declarations

`JSGO_EXPERIMENT_NATIVE_TYPED_DECLARATIONS=1` selects a compact ordered declaration
table only for fresh scopes made by a stock `TypedScopeCreator`. Public
`TypedScope` constructors and custom compiler/convention/observer configurations
keep the original generic storage. This is an isolated, unmeasured prototype.

Current label-filtered evidence is in
`jsgo-results/round2/next-typed-scope-storage-profile.json`. In checks-only,
typed-scope creation accounts for approximately 253 MiB of 1772 MiB sampled
allocation. `TypedScope.Declare` has 38.9 MiB of allocation context; typed lookup
has 0.20 sampled CPU seconds of 7.93 total. String linked-map entries and their
slabs account for 17.3 MiB and 9.0 MiB inside scope creation, but these include
other string-keyed maps and therefore bound eligible bytes from above. This
does not imply that the entire scope-creation cost disappears.

The flow lattice remains `LinkedFlowScope` with its original HAMTs and lazy type
operations. This candidate changes declaration storage used during type building
and subsequent typed lookups, not convergence, transfer functions, or type joins.
`NATIVE_LOCAL_FLOW` independently specializes MustDef/MaybeUse for FSIV.

## Storage and lifetime

On the benchmark's 64-bit platform:

| Object | Previous | Candidate |
| --- | ---: | ---: |
| TypedScope struct | 72 bytes, 80-byte allocation class | 80 bytes, same allocation class |
| TypedVar struct | 88 bytes, 96-byte allocation class | unchanged |
| Declaration record | 48-byte linked entry | 24-byte contiguous key/value |
| Initialized table | 32-byte generic map, excluding entries | 72-byte table, including first record |

Empty scopes point only to an immutable empty sentinel and allocate no table.
The first declaration is fused into its table allocation. Up to eight entries
use a contiguous linear lookup; larger tables add a native name-to-ordinal index.
Variables retain their original pointer identity, declaration order, indices,
type fields, reserved-name checks, implicit-slot creation, and qualified shadowing.
ASCII/valid UTF-8 hits use exact bytes; noncanonical WTF-8/malformed keys normalize
to the same Java UTF-16 units before insertion or after a failed lookup.

There is no global owner, registry, pool, recycling, or separate variable graph.
Appending may copy a private backing array; iterators read the current table by
index and no declaration-record address escapes. Superseded arrays are not kept.
The fused inline record is cleared after growth so it cannot retain a replaced
variable. Slice growth can increase cumulative copied bytes or unused capacity;
fewer allocation objects and less pointer scanning do not guarantee lower RSS.

Removal and clear permanently materialize the original linked map. Active
iterators transfer at their existing position before mutation, preserving
remove-current, remove-future, nested iteration, append, and clear/reinsert
behavior. Native entries/index/inline copy are then released. The remaining
72-byte table wrapper forwards existing views to that original mutable map.

## API boundary and validation

Normal TypedScope calls and explicit embedded-base method calls use the same
store. A nil-vars bridge checks canonical private ownership; ordinary initialized
generic maps return before type inspection. Untyped `Scope.GetVar` keeps its
existing concrete direct-map path. Other inherited base methods gain a branch,
whose cost must be included in bundle/whitespace negative-control measurements.

Go value copies of an opted-in compiler-managed TypedScope or just its embedded
AbstractScopeBase are outside this prototype's private ownership contract. Such
copies can lose the storage association or retain the old `self` pointer. Public
constructors/default mode retain their old copy behavior. Explicit reinitializing
of the canonical base retires native ownership and restores generic storage.
If this candidate proves useful, fence it to the private CLI lifecycle and
materialize scopes before publishing them, or solve that broader ownership API
before enabling it by default.

The table is owned by one compiler traversal at a time. Concurrent independent
compilers are supported; concurrent iterators on the same privately mutable
table are not a new supported contract, because iterator migration tracks active
cursors. Future parallel typed readers need an explicit frozen/read-only
iteration path or materialization first.

Source tests cover layout, stock dispatch, public/custom fallback, explicit base
methods, captured empty/initialized views, Java-equivalent UTF-16 spellings,
105 live mutation histories against the unchanged linked-map oracle, and full
typed-inference AST/symbol/diagnostic/change-order comparisons with both CFG
engines. Central normal/race tests, Java CLI byte gates, and paired elapsed/RSS/
allocation measurements are still required. No improvement is claimed yet.
