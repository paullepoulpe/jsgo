# Native Array<?> binding recipe with fresh results

`JSGO_EXPERIMENT_ARRAY_UNKNOWN_RECIPE=1` is a default-off, ordinary-CLI inference
experiment. It also needs `JSGO_EXPERIMENT_CLI_SOURCE_LOCK=1` for the existing
fresh-input ownership proof. It is independent of the earlier temporary-vector
experiment: leave `JSGO_EXPERIMENT_FRESH_TEMPLATE_VALUES` off to compare only this
recipe against the unchanged original constructor.

The public registry/TemplatizedType factories remain unchanged and fresh. The
private default inference window and its V2 admission/error-count token are
reused. Unknown compilers, callbacks, warning guards, RAI chains, public direct
inference and unowned sources retain their original factories. The optional
caller-owned probe counts real original seeds, fresh coallocated hits and native
mutation fallbacks; it is nil during ordinary/timed compilation.

## Evidence and scope

Existing post-GC checks-only allocation profiles attribute approximately
117 MiB to NewTemplatizedType and 108.55 MiB to CopyFilledWithValues, with
substantial overlap. Those totals include required work and more construction
sites than this prototype handles. The whole old constructor family has roughly
1.6% of sampled checks CPU. No new whole-build improvement is predicted from
those inclusive counts, and downstream equality/cache effects are not measured.

The pinned Java type-algebra oracle records the bare registry's native Array
map: eight keys/seven existing values/one unfilled Array T, with resolved values
`[?, ?, ?, number, ?, ?, ?, ?]`. That is not the full browser compilation shape.
The browser externs re-register IArrayLike/ReadonlyArray inheritance and Array's
implemented interfaces; their ordinary map extension concatenates keys.

The first untimed real-task activation run exposed this mismatch: the private
window opened, but the fixed 8/7 recipe had zero seeds/hits. The clean detailed
run identified only the `map_shape` rejection: typed-app made 2,559 exact native
Array<?> calls; checks-only made 21,686 (21,690 total constructor calls). Each
attempted Capture once. Both actual source maps had **31 keys/30 existing values**
and 31 resolved values, with native/resolver identity checks satisfied. Their
Java artifact checks matched. These old zero-hit runs do not measure the recipe.

Capture now accepts bounded vectors of **1–64 keys**, retaining exactly one
unfilled native Array T. It validates **every** key and source/resolved value,
including entries beyond the diagnostic snapshot's first sixteen. The bound
limits work and pass-owned snapshot retention; it does not recognize benchmark
names or substitute a hard-coded 31-entry type. Larger or semantically unsupported
maps use the unchanged original factory.

Only exact native Array plus exact native unchecked UNKNOWN arguments at the
existing builtin inference construction sites qualify. Promise, arbitrary rest
parameter values and map-based generic specialization remain original when the
independent temporary-vector flag is off.

## What is reused and what remains fresh

The first actual Array<?> construction uses the original factory and returns
that actual result to inference. It seeds an immutable recipe; no extra hidden
exemplar is constructed. The recipe retains only the same registry's native
Array/map/keys and native primitive bindings, not the exemplar or its helpers.

Each hit creates a fresh ordinary TemplatizedType, TemplateTypeMap and partial
TemplateTypeReplacer in one GC-traced header with stable interior pointers. A
second exactly sized GC-traced block owns `n` original-value slots, `n` resolved
slots and one argument. Full-slice caps fence the three regions, so appending to
one cannot overwrite another. All per-result mutable identities and exposed
value slices remain independent. Template keys/submap boundaries keep the same sharing with
the current raw native map as the original factory.

The hit skips only derivation of the validated pure native binding vectors:
temporary total replacement, its traversal/seen sets, source value copies and
the repeated native UNKNOWN-only specialization classification. Capture accepts
only concrete native UNKNOWN/NUMBER or template-key source values; keys have
exact native UNKNOWN bounds and no transformations. Fully resolved results must
be the same native UNKNOWN/NUMBER objects. These traversals have no reporter,
custom visitor, lazy named type or externally supplied bound callback to skip.

Proxy/Object/JSType initialization still runs in original order, including
AddUnresolved. The fresh partial replacer has the same options and initial state.
The original final ResolveIfClosed runs after every fresh type is initialized;
resolution and its capture-stack checks are not replaced with a cached resolved
type pointer. The raw Array and UNKNOWN are already resolved to themselves.

Before every hit, bounded concrete checks reject changed raw/native identities,
map/key/value shape or identities, template bounds/transformations, submap bits
or an open/busy resolver. A miss calls the original factory; a successful normal
construction can seed a new valid recipe. An unsupported seed remains on the
original factory rather than repeatedly attempting capture. Raw properties are
live: proxy property reads/writes still use the original shared raw Array, and
are not frozen or cached by this recipe.

## Why not a shared Array<?> singleton

A whole specialization has mutable resolution/template-check/recursion state
and an exposed mutable template-value slice. Sharing that whole object could
also change pointer-fast-path visitation decisions and public GetCompiler graph
aliasing. Those are not automatically all CLI compatibility contracts, but the
fresh-result design avoids needing that weaker proof.

Property writes already delegate to the same raw Array in the original design;
therefore raw property mutation alone is not a valid argument against sharing.
Likewise, TemplatizedType's partial replacer has visitProperties=false, so its
persistent record-property replacement cache is not active. The actual proof
burden is narrower than declaring all type metadata immutable or all pointer
identity sacred. This prototype nevertheless preserves those mutable per-result
identities while reusing only native binding derivation.

## Allocation layout and retention limits

Static amd64 fields, checked by a focused source test: TemplatizedType 192 bytes,
TemplateTypeMap 88, partial replacer 56. The composite header has 336 payload bytes
and uses the 352-byte allocation class. Two exactly sized allocations avoid a
large fixed-capacity block for every result.

| Binding count | Header payload/class | Data payload/class | Combined classes | Original retained classes |
| --- | ---: | ---: | ---: | ---: |
| Bare registry n=8 | 336/352 B | 272/288 B | 640 B | 624 B |
| Browser externs n=31 | 336/352 B | 1008/1024 B | 1376 B | 1392 B |

These Go 1.27 class calculations concern retained construction components and
exclude shared raw types, recipe snapshots, allocator/span metadata and later
helper allocations. They do not predict peak RSS or elapsed time. Avoiding
repeated binding derivation and reducing temporary objects remains the target.
The one bounded recipe has four exactly sized native-pointer vectors and no
per-result padding. Its maximum is 64 keys; it is dropped when inference exits.

A map/replacer survivor retains exactly the current 352-byte header and its
current value block. Compared with the original map alone, the extra header
coupling is 256 bytes (352 versus 96); compared with an original map+partial
replacer, it is 192 bytes (352 versus 160), plus differences in value-block class
rounding/argument storage. There is no neighboring slab or older wrapper.
Keeping only an exposed argument/value slice retains the data block, not the
header: data entries point only at native bindings and never back at their
wrapper. Thus the weak-pointer tests require that otherwise-dead header to
collect even while its argument slice survives. An argument-only survivor can
still retain both vectors (288 bytes at n=8 or 1024 at n=31, versus 16 originally).
Peak RSS and survivor distribution remain independent measurement gates.

The admitted block has no flow-scope, CFG, source-node or previous-result field.
Binding payloads point only at validated native keys/UNKNOWN/NUMBER from the same
registry, which the original map/value types already retain. The fresh proxy has
no AST source field; documentation/property writes delegate to raw Array as
before. The partial replacer's record cache is disabled. GuardAgainstCycles removes
its transient seen key with a defer, including on panic; the collection removes
and zeros that entry. CaseTemplateType's separate nondeferred replacement path
can only follow the validated native bindings here, so it has no custom bound
visitor to leave behind. No concrete admitted panic-path graph-retention issue
was found. Retaining a map/replacer also retains any later wrapper/replacer
state, not just its construction-time zero fields; this audit's native-only
claim covers construction and admitted native replacement. Consequently retaining
a map/replacer adds one native-result wrapper's lifetime, not a link to an older
flow/AST graph during admitted CLI.
There is no global pool and no cross-compilation recipe. A Go caller can mutate
an exported graph after DoRun in ways that add arbitrary references; this design
does not promise identical weak-pointer timing or retention under those new
mutations. Original public factories do not use coallocation.

## Source coverage and remaining gates

Rhino fixtures compare ordinary vs recipe strings, equality, hash, template
vectors, fresh helper identities and resolver capture/queue state; independent
public value-slice mutations; live raw property reads/writes; raw binding/key/
bound/submap/native identity changes and resolver-open fallbacks; and bare
native shape/layout. Added shape oracles cover 1/2/8/17/31/64 keys, a 65-key
fallback, every expanded metadata/value position, late-entry mutations, seed
slice independence and append-cap boundaries. Weak-pointer fixtures require the first exemplar and prior
results to collect, while intentionally allowing a current map to retain its
one coallocated header until that map is released. A surviving argument/data
slice must not retain the header.

Compiler fixtures require actual private seed/hit counters and compare full
typed AST/scope/source/diagnostic snapshots across loops, nested arrays, class
fields/methods, defaults, rest arguments, try/catch and modules under both CFG
engines. They cover public/custom fallback, malformed constructor argument
exception order, routing after raw map mutation, and actual fresh N=1/N=2 CLI
checks parity. A new stock-runner fixture loads the real default browser
extern set, asserts the actual 31/30 map and positive seeds/hits, and compares
complete typed AST/scope/source/diagnostic snapshots against recipe-off. The earlier vector mode-off fixtures explicitly disable this new
independent flag so their original assertions remain meaningful.

The first 8/7 implementation passed focused Rhino normal/race and compiler
oracles, but real benchmarks proved it inactive. The generalized revision is
unbuilt/untested at drafting. Correctness, race gates and
same-binary default-GC whole-build measurements must pass before any performance
claim or inclusion. No workers, preparse cache, GC tuning, global type intern
cache or weakened expected vectors are involved.
