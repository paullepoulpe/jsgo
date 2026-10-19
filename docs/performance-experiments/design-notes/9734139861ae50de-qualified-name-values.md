# Qualified names at hot callsites

Fresh round-two allocation profiles attribute 30.5 MiB (2.52%) in typed-app and
51 MiB (1.00%) in bundle to `Node.GetQualifiedName`. These are sampled allocations,
not peak RSS. In bundle, 39.5 MiB is the returned string-header allocation and
11.5 MiB is the constructed dotted-name backing storage. The existing pointer
wrapper is not consistently inlined, so consumers can allocate a header even
when they immediately dereference it as a map key.

This candidate exposes a Go value-result adapter to the exact existing qualified
name implementation, returning `(string, bool)` instead of `*string`. False has
the same meaning as the original null result. Global namespace collection, CSS
symbol counting, alias inlining, property type inference, and reverse-interpreter
property lookups use that adapter. No name is cached across AST mutations, no
string is interned, and map lookup order and lazy type operations are unchanged.
Backing storage still allocates where a qualified string is actually required.

Two callsites avoid useless work altogether: CSS symbol-count traversal is
skipped when no classes object was registered, and alias inlining traversal is
skipped when collection found no aliases. In those cases the original callbacks
cannot collect a name, report a classes-object diagnostic, or rewrite a node.
CSS counting accounts for 4 MiB of typed-app and 16 MiB of bundle qualified-name
allocations; alias inlining accounts for 10.5 MiB in bundle. Those numbers are
upper bounds on affected allocation, not expected measured savings: empty-state
frequency has not been measured. With real aliases, an unchanged node's name is
computed once instead of again for alias resolution.

The peephole fixed built-in receiver checks use existing `MatchesName` predicates
for Array, Math, and Number while retaining the preceding qualified-name guard.
They avoid dotted-name construction just to reject qualified receivers.

The value adapter differential test covers exact null/valid results, nonstandard
dotted and empty names, import-star exclusion, this/super, computed/call bases,
UTF-16 surrogate contents, malformed AST exceptions, and nil receivers. Existing
ReplaceCssNames, InlineAndCollapseProperties, GlobalNamespace, TypeInference,
reverse-interpreter, and PeepholeReplaceKnownMethods tests remain the principal
behavior gates. Full Java-output comparisons and time/peak-RSS benchmarks remain
pending central validation. No performance improvement is claimed from source
inspection or profile attribution alone.
