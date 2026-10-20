# Native scoped-name lookup experiment

Set `JSGO_EXPERIMENT_NATIVE_REFERENCES=1` to enable a collector-local native Go
lookup index. The original ordered object map remains authoritative for symbol
publication, first-key identity, iteration and API wrappers. The default path is
unchanged.

ReferenceCollector's key type is the concrete `*Var`. Its final equality is the
immutable pair `(scopeRoot *Node, name string)`, so a comparable Go struct can
replace repeated generic object hashing/equality during lookup. Separate Var and
Scope objects with equal names and the same root still address one collection;
the original first Var object remains the ordered map's key. A Go type embedding
Var cannot override equality while still being passed as a `*Var` key. Null keys
retain their existing handling; a non-null Var with a null scope root still
throws the original NullPointerException.

Every new collection is inserted into both maps at its original first-reference
position. Existing collections are the same objects in both maps. No traversal,
hoisted-function continuation, lexical lookup, filter call, BasicBlock or Scope
identity, reference allocation, AST mutation, diagnostic or behavior callback is
moved or parallelized.

## Mutation and API compatibility

The production collector only writes map entries in AddReference. Its public
map escape is ReferenceMapWrapper.GetRawReferenceMap. Collector-created wrappers
carry their owner while an index is active. Before returning the raw map, the
wrapper permanently retires that collector's index. Replacement, removal, clear,
null values, and mutations through a retained raw map then use the original map
for all later lookups and additions. Previously retained wrappers also delegate
to that authoritative map after retirement.

GetReferences still returns the same mutable ReferenceCollection object, so
changes to a collection's references do not require index invalidation. The
public standalone wrapper constructor keeps its original raw-map behavior.
Unknown/custom behaviors and predicates retain their original callback timing;
if they request a mutable raw map, the same retirement guard applies before any
mutation can occur.

## Evidence required

Tests compare distinct-but-equal Var instances, different roots, empty names,
null keys, first-symbol identity/order, per-reference Node/Scope/InputId/BasicBlock
identity, null-root exceptions, replacement/removal/clear/null map values and
retained/public wrappers. Differential tests exercise InferConsts,
InlineVariables, InlineObjectLiterals, Denormalize and VariableReferenceCheck on
nested/hoisted functions, aliases, objects, arguments, defaults and catches,
comparing AST properties, printed JS and diagnostics with the original path.

The existing profile attributes approximately 0.37 seconds of 22.21 sampled CPU
seconds to object-map lookup underneath AddReference. Lookup and find samples
are nested and must not be added together. This small representation change has
an expected broad ceiling around 1–2%, with extra index memory and first-insertion
work potentially offsetting the saved lookup work. Central tests, race checks,
Java output gates and paired time/RSS measurements are required before retention.
