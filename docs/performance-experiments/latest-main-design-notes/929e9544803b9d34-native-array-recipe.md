# Private native Array binding recipe on current main

This default-off experiment ports only generalized native Array<?> binding
construction and an independent immutable resolved-vector sharing factor onto
ff06ac15. It requires the scoped stock CLI source provenance foundation. It does
not import the earlier native parser, CFG/data-flow solver, BitSet changes,
Unicode implementation, source lock wrappers, node slabs, typed-scope scheduling,
initial-analysis fusion, or temporary fresh-template-vector factor.

Flags:

- `JSGO_EXPERIMENT_ARRAY_UNKNOWN_RECIPE=1`: record the first actual ordinary
  closed-registry native Array<?> result, then construct fresh results from a
  bounded, completely validated snapshot.
- `JSGO_EXPERIMENT_ARRAY_SHARED_RESOLVED=1`: within recipe hits only, share the
  recipe's immutable resolved vector. Public unresolved/argument slices remain
  fresh and capacity-fenced. This flag alone does not enable the recipe.

Only the private DefaultPassConfig inference factory opens a construction mode.
Public inference/factories retain their original behavior. Entry and the
post-FinishAndFreeze boundary validate exact stock compiler/creator/registry,
current root, active CLI lease, reporting/convention/RAI chains and callback
exclusions. Per-scope admission compares captured identities and a synchronized
error count; the known counter plan uses the original monitor's same mutex,
falls back to its original reentrant getter after releasing all cold locks, and
never calls an unknown getter under a private lock. No full input scan or
runtime.Stack-based monitor walk is added to the uncontended per-scope path.

The original generic scope traversals, CFG, data-flow solver and transfer order
remain unchanged. Only three existing TypeInference factory call sites route
through a helper. Non-native Array/arguments, Promise and malformed cases call
the original registry factory. Both flags are captured once per private mode.
Mode state retires on normal and exceptional exit.

The snapshot accepts 1..64 keys, including the actual browser-extern-expanded
31-key/30-source-value map and the bare registry's 8/7 map. All keys, bindings,
native identities, submap bits, bounds and closed-resolver conditions are checked
at capture and on every hit. Raw source binding mutation falls back. Property
maps/prototype caches remain live and are delegated through the original proxy;
there is no dependency on old PropertyMap/PrototypeObjectType layouts.

Every hit owns fresh TemplatizedType, TemplateTypeMap and partial replacer headers.
Their coallocation ties the current header lifetimes together, not previous flow,
AST or result graphs. A separately allocated exact-sized GC-traced value block
keeps publicly exposed slices independent from the header and from one another.
The shared variant's resolved-only survivor retains canonical native primitives,
not the recipe/header. There is no process-global cache or singleton wrapper.

At the previously measured amd64 layouts, header payload336 rounds to352 bytes;
31-key control data1008 rounds to1024 (total1376), shared data512 (total864).
Bare8-key control rounds to640 versus ordinary624. These are allocation-class
models, not CPU or peak-RSS results. Old-base profiles and historical timings are
not new-main performance evidence. New-main end-to-end Java artifact gates and
independent default-GC speed/peak-RSS measurements are required before promotion.
