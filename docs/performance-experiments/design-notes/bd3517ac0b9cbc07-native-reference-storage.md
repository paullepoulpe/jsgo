# Single native ReferenceCollector storage experiment

This candidate is based on current main (`da20b44f`). Enable with
`JSGO_EXPERIMENT_NATIVE_REFERENCES=1`. It changes ordinary CLI compilation only;
there is no worker, extern-preparse cache, vendored Unicode table or node-chunk
allocator dependency.

## Why the earlier lookup index was not retained

The previous candidate kept both a generic ordered map and a native side index.
A first reference allocated/stored an entry in both and performed both hash-map
updates. The ordered generic map remained live even when no consumer requested
its mutable API. This duplicated work and memory, and the experiment regressed.
It therefore did not establish that replacing the generic representation was a
bad idea. Generic lookup beneath AddReference was about 0.37 of 22.21 sampled
bundle CPU seconds in the earlier profile; nested lookup/find totals must not be
added together. The whole collector callback's cumulative 11.6% also includes
recursive hoisted-function traversal and cannot all be attributed to storage.

## Representation and unchanged semantics

The new candidate owns one Go map from `(scope-root pointer, immutable name)` to
an entry index, plus a compact insertion-order vector containing the first Var
object and its ReferenceCollection. The generic ObjectLinkedHashMap is absent
until a consumer requests the raw mutable map.

Var equality and hashing are final methods of the concrete `*Var` key type on
current main. An embedding type cannot override equality while being passed as
`*Var`. Different Var or Scope instances with equal names and the same root
address the same entry. The first actual Var object is retained as the ordered
key. Null keys, empty names and non-null Vars with null-root exceptions retain
the original behavior. No scope canonicalization is introduced.

The original traversal, hoisted first-reference continuations, scope lookup,
filter and behavior callbacks, reference allocation, InputId/Scope/BasicBlock
identities and AST writes remain at their original positions. Only collection
storage and lookup change.

## One-way mutable-map escape

A source audit found no production calls to GetRawReferenceMap; the method is a
public extension point. Collector-owned readonly wrappers delegate to their
owner, so retained wrappers see later additions without materializing a map.
Stringification formats the native entries directly in original order.

GetRawReferenceMap materializes one generic ordered map in first-reference order,
reusing the exact Var and ReferenceCollection objects, then permanently switches
the collector to it. The native index and vector are released. Raw replacement,
remove, clear, null values, and future mutations through a retained raw pointer
then have the original generic semantics; no index is rebuilt or synchronized.
Standalone public wrappers continue to wrap their supplied map directly.

A raw escape may occur inside an active GetAllSymbols iterator. Simply restarting
at a numeric index or taking a snapshot would change remove/clear/reinsert
behavior. The candidate instead positions a paused generic entry iterator at each
active native cursor before exposing the raw map. After the callback returns,
that same generic linked-entry iterator continues. This reproduces even the
existing map's behavior where Clear replaces its header while an active iterator
continues through old linked entries. Nested iterators bridge independently.
Early-stopped native cursors are removed immediately. No bridge coroutine exists
on the normal no-escape path.

## Diagnosis and validation

For an untimed diagnostic run, additionally set `JSGO_DIAG_REFERENCE_STORAGE=1`.
Each completed collector entry point writes a `reference_storage` JSON line to
stderr with native lookups/creates, generic lookups, raw escapes, materializations
and active-iterator bridges. Disable this instrumentation for timings and exact
output comparisons. No file names or source contents are logged.

Tests cover semantic duplicate Var instances, distinct roots, null keys and null
roots, first-key identity and symbol order, per-reference node/scope/input/block
identity, retained readonly/public wrappers, exact-once materialization and
replacement/removal/clear/null mutation behavior. Live-iterator tests compare the
original implementation for additions, removing current/next entries,
remove/reinsert, replacement, clear/addition, nested iterators, escape before
iteration starts, and early stopping. A real custom scope-exit callback and a
retained custom filter both mutate the raw map. Differential consumer tests cover
InferConsts, InlineVariables, InlineObjectLiterals, Denormalize and
VariableReferenceCheck, comparing AST, JavaScript and diagnostics.

Central builds/tests, an untimed dispatch diagnosis and paired CLI time/peak-RSS
measurements remain required before retention. A small storage hotspot does not
justify a broad Java-parity or speedup claim by itself.
