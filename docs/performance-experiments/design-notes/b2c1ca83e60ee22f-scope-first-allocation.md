# Coallocate the first private Scope declaration

`JSGO_EXPERIMENT_SCOPE_FIRST_ALLOCATION=1` is a default-off source experiment
for private stock InlineVariables and RemoveUnusedCode scope creators. It
requires the existing fresh-CLI source-loader ownership proof and audited stock
callbacks; public constructors, raw Scope.Declare and public mutable collection
constructors retain their original allocation.

The existing Scope record remains 48 bytes on amd64 and its existing vars field
remains exactly `*LinkedHashMap[string,*Var]`. No Scope/Var owner pointer, new
storage representation, key borrowing, side table, cache or unsafe conversion is
introduced. The original linked-map engine still performs every lookup, later
insertion, removal, clear and live iteration.

At the first ordinary declaration of an empty private scope, a generic jutil
constructor allocates `{value Var; map LinkedHashMap; entry lhmEntry}` together.
It returns ordinary GC-tracked interior pointers to the Var and map, with the
map's first entry referencing that exact Var under an owned string key. jscomp
runs the same pure Var initialization/token validation and legality checks as
the original declaration path before publishing the map into Scope.vars. The
public NewVar constructor is unchanged. The independently enabled stock-insertion
experiment still controls whether its two redundant map checks are omitted;
coallocation alone retains the original checks and exception order.

The current Var is **64 bytes**, including its self, scope and cached hash fields.
The amd64 allocation model is therefore:

| First declaration | Separate | Coallocated |
| --- | ---: | ---: |
| Var payload / class | 64 / 64 B | included |
| Linked map | 32 B | included |
| Entry payload / class | 40 / 48 B | included |
| Total allocated class bytes | 144 B | 136 payload / 144 B |
| Heap objects | 3 | 1 |

This is byte-neutral. It removes two malloc/GC objects per eligible nonempty
scope map, with no claimed RSS or elapsed improvement. Subsequent declarations
are unchanged. Empty scopes still allocate no map. Actual nonempty-scope counts
and whole-compile time/RSS/TotalAlloc/Mallocs are required to establish value.

Coallocation intentionally ties these three objects' lifetimes. This is sound
in the selected private consumers: their declaration maps are not removed or
cleared after population; RUC's synthetic global declaration only appends.
While the first binding remains present, map -> first entry -> first Var ->
Scope -> map already makes all three objects mutually reachable. Implicit goog
namespace replacement does not use this constructor: the implicit declaration
has already created a normal map, so its later replacement is not a nil-map first
insertion. The private fresh-CLI proof excludes public capture, custom source/
warning/reporting callbacks and arbitrary pass schedules. No scope graph remains
owned by this constructor beyond the original private pass lifetime.

Do not apply this constructor to unrestricted public mutable collections. A
removed embedded first value or its old next link can remain reachable with the
map; public remove/clear/reinitialization can otherwise change object lifetimes.
Functional mutation and live-iterator behavior still use the unchanged linked
engine, with differential fixtures covering current/next removal, replacement,
clear-and-append, early exit and promotion. The common-lifetime restriction is
explicit in the internal jutil constructor documentation.

Private untimed counters report actual first coallocations. Scope fixtures
compare ordered bindings, canonical Var/Scope identity relationships, declaration
indices, inputs, implicit/bleeding/parameter cases, constructor failure types/
messages and unpublished storage on failure. Real IV/RUC pass snapshots compare
AST/output/diagnostics/change epochs; fallback and partial-failure fixtures check
ownership fences and private flag cleanup. A standalone allocation assertion
checks one allocation versus three for equivalent escaped value/map/entry setups.

Status: formatted and whitespace checked only. No Go builds/tests or runtime
measurements were run during the root benchmark/build reservation. All benchmark
comparisons must use the user's default Go GC settings, without GOGC tuning.
