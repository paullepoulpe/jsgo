# Captured native binding validation session

This separate default-off factor builds on the clean-main recipe. Set
`JSGO_EXPERIMENT_ARRAY_CAPTURED_VALIDATION=1` in addition to
`JSGO_EXPERIMENT_ARRAY_UNKNOWN_RECIPE=1`. Shared-resolved storage stays independent.
It changes validation work, not the type construction or public data layout.

The ordinary `TryNew` and `TryNewSharedResolved` methods retain the full bounded
key/bound/value/submap scan. A captured session first performs that full scan,
then uses O(1) concrete native identities, resolver state/stack/queue, source-map
identity, vector lengths/backing identities and submap-object identity on hits.
It still executes the exact same fresh allocation/proxy/resolver sequence.

## Lifetime proof on ff06ac15

The proof applies only after the original `TypedScopeCreator.FinishAndFreeze`,
inside the private DefaultPassConfig inference window, with its existing stock
CLI lease/root/options/creator/registry/reporting/RAI/convention fences. Public
constructors, public inference, injected typed ASTs, custom callbacks and saved
compiler inputs do not enter that window. Each captured hit also requires the
creator's actual FROZEN stage. The first actual original Array<?> result supplies
the full snapshot; no hidden exemplar is constructed.

The source-write audit is narrower than a claim that all JSTypes are immutable:

| Captured dependency | Production writers | Why absent during the window |
| --- | --- | --- |
| Native registry entries | JSTypeRegistry initialization/RegisterNativeType | Inference uses the existing closed registry; no TI/RAI/TTL call reaches registration |
| Raw Array map pointer | initialization and JSTypeBase.MergeSupertypeTemplateTypes | Runtime merge callers are FunctionType.SetImplementedInterfaces/SetExtendedInterfaces, invoked by registry initialization or jscomp FunctionTypeBuilder during typed scope building |
| Template bounds / referenced bound | TemplateType constructor/SetBound | Only jscomp caller is FunctionTypeBuilder; the frozen creator cannot create an unvisited scope |
| Template key/value vector elements | constructor stores of supplied arrays | No production assignment through GetTemplateKeys/GetTemplateValues exists; matching/equality/type inference only read these vectors |
| Submap bits | TemplateTypeMap constructor/copy helpers | Copy helpers allocate a new bitset before writes; none mutate an existing map's bitset |
| Native Unknown/Number and raw Array resolution | native initialization / ordinary resolve | Full capture requires their existing resolveResult already equal themselves; later Resolve returns those cached selves |

The post-freeze cast path evaluates JSDoc via JSTypeRegistry's comment-node
interpreter. It creates fresh function/record/templatized/NamedType objects and
may report; it does not use jscomp FunctionTypeBuilder, set native interfaces or
set template bounds. The type-transformation evaluator constructs fresh types
and reads bindings; it has no map/key/bound writer. Invocation template matching
and inference replace values in fresh temporary binding maps and new function
types, not the source native map. Stock RAI changes flow scopes, not native
binding definitions. Property inference can write native properties: these are
intentionally outside the snapshot and remain live through ordinary proxy
lookup, including current-main PropertyMap/prototype caches.

Closed per-result resolver registration remains unchanged. Native recipe values
resolve only canonical Unknown/Number/template bindings; source raw resolution
is cached-self. Fresh cast/TTL resolution cannot reopen registry definition or
run FunctionTypeBuilder. Both OPEN and CLOSING states and nonempty resolver stacks
are rejected by the session header check.

## Explicit exclusions and invalidation limits

An O(1) session deliberately cannot detect an in-place write to an existing key,
value, bound or bitset. Such a write violates its lifetime precondition. This is
not a drop-in replacement for the public mutable TemplateTypeMap API. Public
recipe methods keep full validation and reject these mutations. Header changes
still return nil and the compiler falls back to the original factory and full
capture. No mutation version counter or persistent cache layout is added.

A private diagnostic constructor observer retires the captured lifetime proof
permanently before it runs. Its mutations use ordinary full recipe validation;
reentry from that observer invokes the original registry factory. Removing the
observer later does not rearm captured validation. The full per-scope admission
oracle likewise uses full recipe validation. A nested private inference window
cannot acquire captured validation. Numeric activation probes do not supply any
callback and may count actual admitted sessions/hits/fallbacks without changing
admission.

Sessions are owned only by the inference mode and are retired with it. They hold
the already retained recipe/source vectors, no result wrapper, flow/AST state or
cross-invocation owner. Public result slices and maps are identical to the recipe
control. There is one small session allocation per successful private recipe.
The expected saving is repeated 31-key validation work, not per-result payload
bytes. Speed and peak RSS must be measured independently on current main;
historical recipe timings provide no result for this factor.
