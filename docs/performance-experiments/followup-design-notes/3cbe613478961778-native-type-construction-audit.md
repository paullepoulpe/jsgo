# Broader native type-construction opportunities

Source-only audit; no implementation, build, profiling, timing or speed claim.
Source examined: structural-bitset-template current source, plus the generalized
Array recipe draft in jsgo-array-specialization-reuse. Full direct callsite list
is saved alongside this report as native-type-construction-callsites.txt.

## First distinguish constructor counts from per-binding work

The existing default-GC template-specialization-audit.json attributes checks
CopyFilledWithValues 113,822,790 sampled bytes / 1,238,945 sampled objects. Of
that caller family TraverseArrayLiteral accounts for 84,451,046 bytes (74.2%).
InferTemplatedTypesForCall accounts for 7,342,416; annotation declaration and
FunctionTypeBuilder parameter handling about 6.3MB each. These are sampled
allocation events, not constructor invocation counts or peak/live RSS. One real
Array<?> specialization with a 31-entry inherited map performs many per-entry
visitor/set allocations. Therefore 21,686 actual calls can explain a much larger
object count; a million sampled objects does not imply a million unhandled
specialization calls. The generalized recipe is aimed at that dominant known
family. Re-measure its real activation and output before extrapolating beyond it.

## Actual construction routes

The three private inference calls currently routed through createTemplatizedType
are rest-parameter Array construction, TraverseArrayLiteral and dynamic-import
Promise construction. Additional routes retain the public factory:

* TypeInference.TraverseInstantiation calls CreateTemplatizedTypeOfMap with
  inferred argument types, often arbitrary templates/unions or user classes.
* Promises_CreateAsyncReturnableType always constructs native IThenable<?> and
  may also construct IThenable<asyncTemplateUnion>; Promises_WrapInIThenable
  builds IThenable<unwrapped>. JsIterables_CreateIterableTypeOf builds native
  Iterable<elementType>. These public helpers run in multiple owners, not just
  the current private inference window.
* DestructuredTarget, TypedScopeCreator rest-parameter declarations,
  FunctionTypeBuilder return/yield helpers and AstFactory construct Array,
  Iterable/Generator/Promise types from both canonical and arbitrary inputs.
* Rhino TemplateTypeReplacer.CaseTemplatizedTypeUnguarded constructs a new
  specialization when its base or arguments change. JSType raw-type matching
  and UnionType alternate building construct unknown specializations from an
  empty argument list. NamedType.ResolveInternal constructs from resolved
  annotation arguments, and registry parsing constructs explicit generic types.
* Registry initialization constructs native interface inheritance maps while the
  resolver is open. Those calls cannot use a closed-registry inference recipe.
* TemplateTypeMap constructors also serve CopyWithExtension and CopyWithoutKeys;
  these are a different family from specializations and need separate counts.

Reachable larger private windows could cover InferTemplatedTypesForCall and the
stock Promises/JsIterables helper chains, but inputs must be classified by actual
canonical native identity, shape and resolver state before choosing a recipe.
Do not redirect public helpers globally or cache user-class specializations.
A bounded per-inference recipe table keyed by exact raw native identity plus
canonical primitive/unknown argument identity is plausible. Seed each recipe
from its first actual original result and keep every wrapper/map/replacer fresh.
Empty argument lists and explicit unknown arguments must remain distinct until
default-bound equivalence is proved; multi-parameter order matters. Resolver
OPEN/CLOSING, template transformations, Named/Proxy/union/custom inputs retain
the complete original factory. Inventory first: no coverage claim yet for these
additional families.

## Simpler next step before a broad lazy-map redesign

TemplateTypeMap.resolvedTemplateValues is private to rhino and only assigned by
NewTemplateTypeMap in ordinary production. Its read sites are inside that file;
no public method exposes the resolved vector. GetResolvedTemplateType returns
one JSType, while GetTemplateValues exposes only the separate unresolved vector.
Current copy operations copy/share unresolved values, not the resolved vector.
The generalized native recipe's resolved vector contains only canonical Unknown
or Number, owns its snapshot and is never mutated by stock methods.

Consequently a separate private recipe factor could SHARE ONLY that immutable
resolved vector across recipe-created maps, while keeping every map/header,
unresolved values and argument vector fresh. This adds no map fields or hot
lookup branches and removes one n-element copy plus one n-element backing region
per result. Original public factories remain unchanged. It requires explicitly
accepting that package-internal tests that directly mutate the unexported
resolvedTemplateValues slice are outside this new private implementation
contract; all exported Closure-like operations retain original behavior.
For n31 the vector is 496 payload bytes (normally class512); at 21,686 calls the
allocation opportunity is about 10.6MiB before revised coallocation rounding.
This is a model, not RSS or speed. A map/values survivor retains no previous
wrapper and only the immutable native snapshot; registry/native ownership was
already present. Public GetTemplateValues mutation must remain per-result and
must NOT mutate the private resolved snapshot (matching the original eager map,
whose resolved values do not recompute after such a mutation).

## Larger lazy unresolved-value view

A private recipe-backed map can defer its fresh unresolved vector while scalar
lookups and HasAnyTemplateTypes read the immutable descriptor. Before
GetTemplateValues exposes a slice, or CopyWithoutKeys preserves its original
backing-sharing relationship, materialize a fresh per-map vector. Every escaped
slice and every new wrapper/map/replacer remains independent. Recipe values
must be construction-time snapshots, never live views of subsequently modified
raw maps. Original constructor resolution effects must still occur for all
unknown/custom inputs; only proven pure-native derivation can be deferred.

There are 32 textual non-test templateValues/resolvedTemplateValues mentions in
current rhino, concentrated in TemplateTypeMap and fresh-vector helper. Copy and
constructor sites need review, especially CopyWithoutKeys' intentional unresolved
slice sharing, String order, the nil/empty distinctions, and callback timing of
CopyFilled/extension. A descriptor pointer would grow TemplateTypeMap88 to96
bytes, currently the same allocation class, but introduces branches and snapshot
retention on every result. Resolved-only sharing is much less invasive and
should be measured first. Native immutable views should remain private initially;
public generic factories need not change.

## Untimed constructor/exposure inventory design

Before routing more callsites or implementing a lazy unresolved view, add a
separate diagnostic SOURCE OVERLAY to a frozen ordinary two-inline control:

1. Count NewTemplatizedType once (not once per wrapper factory), grouping exact
   canonical native base enum or OTHER; supplied argument tuple of canonical
   native identities / TemplateType / OTHER; map key/value counts; resolver
   OPEN/CLOSING/CLOSED. Observe concrete fields/pointer identity only, never
   GetReferenceName/String/Equals/Resolve or arbitrary JSType methods.
2. Independently count NewTemplateTypeMap by CopyFilled, extension, removal,
   registry-init or OTHER route. Record keys/values/resolved lengths and pure
   native eligibility. Use bounded runtime.CallersFrames only in this untimed
   diagnostic to identify the nearest non-factory constructor/caller owner;
   retain numeric counters and function-name strings, not types/AST/maps.
3. To assess laziness, overlay one numeric group ID and exposure-bit field on
   TemplateTypeMap (no production state). Count each map's first
   GetTemplateValues exposure, each CopyWithoutKeys sharing boundary and copy
   route. Do not retain object pointers in diagnostic maps; no artificial
   lifetime extension. Getter instrumentation must not invoke additional type
   queries or materialize values itself.
4. Preserve actual native9+CLI_SOURCE_LOCK and default runtime GC. Keep all new
   constructor recipe factors off in the baseline census. Do not enable compiler
   tracing, change warning/convention/observer guards, or use diagnostic elapsed
   time/heap as a benchmark. Gate every artifact against Java and save frozen
   overlay source hashes. Ordinary method-event counts are not malloc counts.
5. Rank eligible original constructor bytes and exposure fractions by raw type,
   argument shape, owner and resolver state. Only then choose an additional
   private caller or a lazy map factor. A high unknown-specialization count
   during OPEN/CLOSING does not justify enabling the closed pure-native recipe.
