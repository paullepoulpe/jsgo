# Small Java candidates: source audit only

Pinned source: `/workspace/jsgo/.cache/closure-src` (Closure v20261015). No Java process, compile, benchmark, Docker, network, or frozen-v1 change was performed. These are profile targets, not speedup claims. The first two affect ordinary checking/optimization; the latter two require source-map output.

## 1. Use identity membership when collecting PropertyMap ancestors

`rhino/jstype/PropertyMap.java:213–270` collects ancestors in `LinkedHashSet<PropertyMap>` on **every** `getAllKeys()` call, before checking its cached key set. `PropertyMap` has identity equality but its `hashCode()` at lines 312–322 enumerates `properties.keySet()` and possibly known-symbol keys. Consequently each visited ancestor hashes all its own properties just to prevent duplicate visits. This is additional to the genuinely necessary inheritance traversal/cache validation. Consumers include `ObjectType.getPropertyNames/getAllKeys` and `SubtypeChecker`.

Small change: collect the existing DFS sequence in an `ArrayList<PropertyMap>` and use an `IdentityHashMap<PropertyMap, Boolean>` for visited membership. On entry, return when `seen.put(this, TRUE) != null`, otherwise append to the list and recurse in the original order. Preserve the list order because known-symbol key insertion order currently follows traversal order. Keep `PropertyMap.hashCode()` itself unchanged: JSType hashing requires its existing meaning.

This replaces property-count-dependent hashing (plus content-hash collisions) with identity lookup. HotSpot cannot generally remove those key-set walks through the current general collection API. Cost: an extra small list/map pair; shallow maps with few properties could lose. Validate diamond/cyclic inheritance, cache invalidation after property or parent changes, and known-symbol order. Profile targets: `PropertyMap.getAllKeys`, `collectAllAncestors`, `hashCode`, `AbstractSet.hashCode`, `TreeMap$KeyIterator.next`.

## 2. Reuse the namespace index's full-name string

`jscomp/GlobalNamespace.Name.getFullName()` at lines 1322–1323 recursively concatenates all parent prefixes on every call. Both `baseName` and `parent` are final (1231–1232); constructors only create a new child beneath an existing parent, and no reparenting/renaming API exists. `getOrCreateName` (1140–1163) already receives and retains the complete exact string as the namespace index key, including per-module keys.

[full-name-cache.patch](full-name-cache.patch) adds one nullable string field, seeds it from that existing key when the Name is created, and lazily computes it only for testing/direct `addProperty` paths. The [unshaded candidate source](full-name-cache/com/google/javascript/jscomp/GlobalNamespace.java) is ready for an isolated overlay. The production path retains no second full-name string. `getName()` inherits the improvement. This differs from v1's delayed scope/chunk access.

The avoided concatenation returns strings used in AST annotations and lookups, so it is real work that can survive JIT optimization. The field adds a reference to every Name and may increase aligned object size; do not assume allocation savings guarantee wall savings. Validate dotted globals, module-local names, direct test construction, and collapsed output/source-map original names. Profile targets: `GlobalNamespace$Name.getFullName`, `InlineAndCollapseProperties` flatten/alias methods, and string-concat allocation below them. A local single-call reuse in `rewriteAliasProps` (903, 906) is still smaller, but only removes its duplicate call.

## 3. Finalize source mappings in place

`jscomp/CodePrinter.MappedCodePrinter.getSourceMappings` (141–154) creates another list and another `SourceMap.Mapping` for every mapping solely to adjust `.end`. `toSource` creates a fresh printer (861), calls `endFile` (884), and calls `getSourceMappings` once (890). Mutating each existing mapping's `.end` and returning `allMappings` removes one escaping record allocation per mapping and the duplicate reference array.

This is a small allocation change, but ownership needs review: a custom `CodeGeneratorFactory` receives the consumer and can technically retain it. Check that returning existing records does not violate expected snapshot behavior. Validate mapping order, inserted line cuts, newline/end-position boundaries, and output bytes. Profile targets: `CodePrinter$MappedCodePrinter.getSourceMappings`, `SourceMap$Mapping` allocations. Defer if source-map output is a negligible profile share.

## 4. Skip the empty source-content scan

`com/google/debugging/sourcemap/SourceMapGeneratorV3.addSourcesContentMap` (486–503) allocates/fills an array-backed list and looks up every source before discovering there is no source content. Add `if (sourceFileContentMap.isEmpty()) { return; }` at entry. An empty map cannot produce `found == true`, so omission of `sourcesContent` stays identical. This avoids O(source files) work per output map and is the smallest low-risk patch here. Check actual map state, not the compiler option, because input maps may supply content independently. `includeSourcesContent` defaults false.

Expected whole-compile opportunity is limited unless there are many sources/output maps. Validate empty map, populated map, and partial content. Profile target: `SourceMapGeneratorV3.addSourcesContentMap`; allocations may show `Collections$CopiesList.toArray` or `ArrayList.addAll` below it.

## Screen-outs

- `Node.getQualifiedName()` already uses a capacity-reserved builder; do not replace it without a demonstrated caller-level redundancy.
- Several apparent `name.matches("...")` sites use `Property.Key.matches`, not Java regexes.
- `NodeUtil.getOriginalName(String)` does compile a regex via `replaceAll`, but there are no callers in this pinned source tree; it is a poor current-workload target.
- Source-map `Util.escapeString` explicitly documents choosing string append over ranged append for performance; changing it blindly would reverse an existing deliberate choice.

Use warmed JFR to prioritize candidates 1/2, then check 3/4 only if mappings are enabled and visible. Independent byte checks and end-to-end confirmation remain necessary before applying any candidate to a shipping artifact.
