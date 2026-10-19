# Access-control property resolver experiment

This is an opt-in ordinary CLI experiment, not a persistent type cache. Enable
`JSGO_EXPERIMENT_ACCESS_PROPERTY_FACTS=1` together with
`JSGO_EXPERIMENT_CLI_SOURCE_LOCK=1`. No timing or memory gain is claimed yet.

`CheckPropertyVisibility` repeats defining-source and visibility-owner resolution
inside `GetEffectivePropertyVisibility`. Applicable references separately walk
implicit prototypes for deprecation and constness. The existing
`constPropertyInits` table tracks ordered writes; it does not cache these facts.

A cache belongs to one built-in `CheckAccessControls.Process` invocation. It
stores at most 8192 exact object-identity/property pairs. Facts are lazy and
independent: inherited definition node, ordinary visibility owner, overriding
visibility owner, deprecation result, and constness. Cached negative results have
explicit presence bits. At capacity, existing hits continue and new keys use the
original resolver. The deferred cleanup also runs when traversal throws. The
public constructor, static helpers, and public constness query keep their
original behavior.

The admission proof is deliberately narrow. Only DefaultPassConfig's private
constructor can enable it. The compiler must be exact stock, hold the private
fresh-CLI source-loader ownership marker, and have no typed/saved input or custom
passes. Known stock warning guards (including nested compose guards), compiler
error reporter/manager, reverse interpreter, convention, redeclaration handler,
and change observers are required. This excludes arbitrary callback insertion of
custom type graphs. The loader marker is installed only by exact stock DoRun
before Compiler.Init and removed when it returns. Public Compile calls cannot
acquire it. The experiment depends on this marker rather than claiming all
publicly mutable type objects are immutable. Proxy, Named, Templatized and custom
receiver implementations use their original queries.

Two distinctions must remain intact. The report owner treats a documented
declaration as overriding, but effective visibility uses `Override()` alone.
Defining-source lookup also searches secondary interfaces, whereas visibility,
deprecation and constness walk implicit prototypes. A missing definition falls
back to the current reference's source; that fallback is not cached by property.
Ordered constant initializations and all per-reference diagnostic decisions are
always recomputed.

The map and per-key records cost memory. On amd64 the source layout gives a
64-byte record, a 32-byte identity/string key, and a pointer map value, plus map
metadata/capacity overhead. The cap bounds record payload at 512 KiB; total cache
memory is larger. This tradeoff needs macro elapsed/RSS/allocation measurements
at the default GC setting. The previously labeled checks-only profile attributed
7.8% of CPU to the whole access-control pass; this cache replaces a subset of
that work and should not be presented as a larger measured win.

Focused tests cover real Process dispatch and cache hits against independent
stock compilers, repeated properties, private/protected/deprecated/constant
properties, interfaces, multiple files, exact diagnostics/output/AST/change
stamps, both override definitions, missing-source fallbacks, mutations between
windows, secondary-interface definitions, custom/public fallbacks, capacity,
public-method behavior and exceptional retirement. They have been prepared
source-only and await the coordinator's build slot.
