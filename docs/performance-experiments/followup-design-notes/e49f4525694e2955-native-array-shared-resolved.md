# Private immutable resolved-vector sharing

Status: isolated source factor atop generalized recipe 2a6fee38; unbuilt and
untested. Enable ARRAY_UNKNOWN_RECIPE plus ARRAY_SHARED_RESOLVED with the
`JSGO_EXPERIMENT_` prefix. The new flag alone does not activate a recipe or an
inference window. It is captured once when the private mode is constructed;
there is no environment lookup on constructor hits. Original recipe TryNew and
all ordinary public factories remain the control.

## What changes

Both modes still create a fresh ordinary TemplatizedType, TemplateTypeMap,
partial TemplateTypeReplacer and independent public unresolved/argument slices.
They perform the same native-snapshot validation, proxy initialization,
AddUnresolved capture and final ResolveIfClosed protocol. The sharing mode uses
the recipe's private immutable resolved native-value vector instead of copying
it into the current result's allocation. It preserves the independent ordinary
factory seed; later public factory calls remain entirely unchanged.

## Immutability and public/copy boundaries

A repository-wide non-test search for resolvedTemplateValues writes finds one
ordinary assignment: NewTemplateTypeMap assigns a freshly computed vector at
construction. Recipe initialization assigns its own immutable snapshot and
constructs map headers; no stock map or replacer method later writes resolved
vector cells. Its readers are confined to TemplateTypeMap. The public
GetResolvedTemplateType returns a single JSType, not the underlying slice.
GetTemplateValues exposes the distinct unresolved vector, which remains fresh
per result. The argument vector returned by GetTemplateTypes also remains fresh.
Canonical Unknown/Number objects are already shared by ordinary construction;
sharing their private interface-vector cells introduces no new mutable types.

Mutating a public unresolved-value slice never recomputes the original eager
resolved map. Sharing preserves this behavior, including mutations made before
future lookups/equality/hash operations. No exported operation exposes the
private resolved vector. Existing package-internal tests that directly overwrite
the unexported field are outside this new private implementation contract; the
original public and recipe-control constructors still own independent vectors.

Recipe maps have every key filled, so CopyWithoutKeys is the original no-op
returning the same map. CopyWithExtension always uses the untouched generic
constructor for actual extensions: it copies unresolved values and computes a
fresh resolved vector. If that extended map has an unfilled key subsequently
removed, CopyWithoutKeys preserves its original unresolved-backing sharing while
again computing a fresh resolved vector. No copy method propagates the private
recipe's shared resolved slice into a mutable output vector.

The immutable slice itself contains only canonical Unknown/Number references.
It points to no recipe object, wrapper, map, partial replacer, prior result, AST
or flow scope. A retained result can keep that native snapshot alive after the
inference window ends; those native types/registry were already retained by the
ordinary result. A retained vector alone must not keep the recipe descriptor or
result header alive; a weak-reference fixture checks both.

## Allocation model (amd64)

The existing 336-byte wrapper/map/replacer header uses class352 in both modes.
Control values use (2*n+1) JSType cells; shared values use (n+1), with full slice
caps separating unresolved values from the one public argument. The recipe's
resolved snapshot already exists in both modes, so sharing adds no extra
snapshot allocation.

| Keys | Control value block/class | Shared value block/class | Total control/shared |
| --- | ---: | ---: | ---: |
| 8 | 272/288 B | 144/144 B | 640/496 B |
| 31 | 1008/1024 B | 512/512 B | 1376/864 B |

This removes an n-element copy and backing region, not an allocation object:
both modes still allocate one header plus one current public-value block. At
21,686 31-key hits the byte model is about 10.6 MiB fewer allocation-class bytes;
real hits exclude the original seed and may change with workload. This is not
peak RSS, live memory, elapsed time or a claimed performance improvement.

## Validation prepared

Component fixtures cover n=1/8/31/64, actual resolved sharing and control
nonsharing, fresh object identities, equality/hash/format/resolution/protocol,
public mutation/append isolation, CopyWithExtension and CopyWithoutKeys alias
boundaries, and weak native-payload retention. Compiler fixtures compare full
typed snapshots under original/recipe/shared flags and both CFG solvers,
assert positive actual sharing dispatch and show that the flag is frozen at
private mode entry. All tests await a coordinated build slot. Runtime GC policy
is unchanged.
