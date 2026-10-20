# JSType dispatch audit (source-only; no implementation)

> Attribution correction: retained stack-family totals are lower bounds, not hard performance ceilings. CPU stacks are limited to64 PCs; owning pass frames can be lost during generic recursion. See `profile-stack-depth-audit.md`.

Examined source 4f417579 and the already-captured enabled 11534351 CPU/allocation profiles. No compiler build, new profile capture, or benchmark was run for this audit.

## Measured opportunity

| Enabled profile | All sampled CPU | Is*Type / ToMaybe* / GetTypeClass stack union | GetTypeClass samples | Unknown-query samples |
|---|---:|---:|---:|---:|
| typed-app | 6.87s | 0.13s (1.9%) | 0 | 0.03s |
| checks-only | 10.20s | 0.21s (2.1%) | 0 | 0.03s |

The union is a conservative candidate-family measure of methods named `Is*Type`, `ToMaybe*`, and `GetTypeClass`; it includes useful semantic work, not just interface dispatch. It excludes unrelated algorithms such as `IsSubtypeOf`. Sampled method presence does not count every caller instruction or establish a rigorous cycle bound: very small indirect callees may have no sampled PC. Nevertheless this is insufficient evidence to justify a broad kind-tag redesign as the next material optimization.

`JSTypeBase` cumulative totals are 0.31s typed and 0.78s checks, but these contain `Equals`, `IsSubtypeOf`, and their algorithms. Treating those totals as predicate-dispatch savings would overstate the opportunity.

The premise that base predicates generally call `self.GetTypeClass` does not hold in this source. Most base scalar predicates return false, and most `ToMaybe` methods return nil. `IsUnionType`, `IsFunctionType`, `IsNamedType`, etc. call their virtual `self.ToMaybe...` methods; `IsObjectType` calls `self.IsObject`. The direct base `self.GetTypeClass` use is in `EagerlyResolveToSelf`, not the hot predicate family. The subtype dispatcher and enum `IsTypeOf` use `GetTypeClass` separately.

## Field layout and initialization

On amd64, source-derived `JSTypeBase` layout is 64 bytes: self interface at 0, resolveResult interface at 16, templateTypeMap pointer at 32, three bools at 40/41/42, five padding bytes at 43..47, templateCheckResult pointer at 48, registry pointer at 56. A uint8 stock-kind-plus-one field could fit at offset43 without increasing size; a four-byte enum could fit at44. This layout inference was not verified with a new executable during the build reservation.

The base self field is assigned only by public `InitJSTypeBaseWithTemplateTypeMap`; higher-level bases have distinct shadowing self fields. Any future tag must be reset on every public initialization, including reinitialization from previously tagged stock self to custom self. Zero must mean no proven stock identity. Initialization must use exact dynamic concrete-type recognition, not invoke custom `GetTypeClass` earlier than the existing constructor schedule. Arbitrary external wrappers can embed a stock base, override methods, and pass themselves as self; treating those as a stock type changes public virtual behavior.

## A runtime-class tag is not a semantic predicate tag

`NamedType`, `TemplateType`, and `TemplatizedType` extend `ProxyObjectTypeBase`. Their fixed runtime class does not determine the result of delegated predicates/downcasts. For example, `ProxyObjectTypeBase.ToMaybeFunctionType` and `ToMaybeUnionType` query the current referencedType, which `SetReferencedType` and named-type resolution can replace. `ObjectTypeBase.IsUnknownType` also follows prototype/interface state and mutates the existing unknown cache; `UnionType.IsUnknownType` checks alternates. A cached class cannot safely replace these operations.

An exact-stock type switch could eliminate selected second-hop dispatches for proven nonproxy scalar types, preserving fallback for custom wrappers and dynamic proxies. But an external helper would itself need `GetJSTypeBase` interface dispatch or a type switch, and the current measured candidate family is only around2%. No tag or broad overrides implemented.

## Equality/subtyping follow-up

| Enabled profile | EqualityChecker cumulative | SubtypeChecker cumulative | Sampled allocations under each |
|---|---:|---:|---:|
| typed-app | 0.06s | 0.19s | 0MiB / 4.5MiB |
| checks-only | 0.20s | 0.43s | 1.5MiB / 20.5MiB |

Totals overlap where subtyping invokes equality. Checker constructors have no positive allocation samples; introducing pooling or another global result cache lacks evidence.

Reflexive equality is already tested in `JSTypeBase.Equals` before checker construction, and again inside the recursive equality algorithm. `GetLeastSupertype` has an identity shortcut. SubtypeHelper also tests identity, but some subtype-class-specific dispatch occurs before it. Moving a new identity shortcut across that dispatch would need proof about exact stock behavior and custom callbacks; it is not a mechanical safe shortcut.

One real repeated-work observation: `EqualityChecker.AreTypeMapEqual` does not return early when its two concrete TemplateTypeMap pointers are identical. It performs a nested key scan and recursive value checks. However its entire cumulative cost is only0.01s typed /0.05s checks (at most0.5% checks CPU), and a blanket shortcut changes exceptional behavior on invalid/uninitialized maps. Do not prioritize this as a large win. There is no demonstrated >5% whole-program redundant equality/subtyping target in these captures.

Raw data: `jstype-predicate-profile.json`, `jstype-predicate-only-profile.json`, `jstype-equality-subtype-profile.json`.
