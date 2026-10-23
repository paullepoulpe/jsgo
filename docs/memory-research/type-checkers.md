# Research cluster: type-checkers — compact types, symbols, scopes and per-node results

Scope: Sorbet, TypeScript (JS checker and typescript-go), Flow (now a Rust port), Hack/hackrs, Pyright, mypy, Closure's own TypedAST/Colors (plus a note on rustc/Swift/go/types interning). Then a precise mapping onto jsgo's `internal/rhino/jstype_*.go`, `internal/colors`, `internal/jscomp/typed_scope*.go`/`typed_var.go`, the TypedAST serializers, and `rhino.Node.typeOrColor`.

Sources were read in shallow sparse clones under `scratchpad/research/src/{sorbet,typescript-go,flow,hhvm,pyright,mypy}` (none were built or run). Two numbers come from read-only greps of the typed-full inputs in `/home/user/jsgo/.cache/bench-inputs` and the bundled `externs.zip`. No jsgo code was built, run or changed.

---

## 0. Headline findings

1. **A global pins the whole JSType graph after `ConvertTypesToColors` (likely leak, near-zero-effort fix).** `rhino/jstype_property.go:46` declares a package-level `propertySymbolList []*KnownSymbolType`. Every symbol-keyed property (`[Symbol.iterator]` and the like, which the default externs declare, e.g. `Iterator.prototype[Symbol.iterator]` in es6.js and 20+ such members in html5/w3c_dom1/svg; they go through `DeclarePropertyIfNamespaceType` → `DefineDeclaredPropertyByKey` → `NewPropertyByKey`) appends its `KnownSymbolType` to it through `registerPropertySymbol`, and nothing ever clears it.
   - Each `KnownSymbolType` embeds `JSTypeBase.registry`. So after `Compiler.ClearJSTypeRegistry` the registry stays reachable from a global, and with it `nativeTypes`, `scopedNameTable`, `eachRefTypeIndexedByProperty`, `nonRefTypesIndexedByProperty` and `objectDocInfo`.
   - Through those tables stay alive all declared class, interface and namespace types, their `PropertyMap`s, every `Property` (its `docInfo` is the *full* declaration JSDoc, with its type-expression Node trees) and every method `FunctionType`. Through `FunctionType.source` and `Property.propertyNode`, AST subtrees that later passes deleted stay alive too.
   - It is consistent with `docs/gclab-2026-10-18.md` §3: "types & colors (JSType, registry, TemplateTypeMap)" is **18 % of the typed-app heap mid-loop vs 21 % at the post-check peak**. If the types really died at `convertTypesToColors`, that row would be colors only, a few percent.
   - It also leaks across compilations in a long-lived process (tests, batch drivers).
   - Fix: keep the symbol list per registry, or clear it in `ClearJSTypeRegistry`. Then verify with a `weak.Pointer` check (§5).
2. **The AST→type link should be a 32-bit handle, not `unsafe.Pointer`.** This is what Sorbet's `SymbolRef`/`TypePtr`, TS's numeric type and node ids, Flow's `Index<T>` and Closure's own TypedAST `TypePointer` all do. It is the missing piece for a pointer-free (noscan) `rhino.Node` and costs no pass changes.
3. **Colors throw away the compactness that Closure's TypePool already has.** The TypePool proto stores `int32` type pointers and `int32` string indices. jsgo's `colors.Color` turns that back into four `LinkedHashSet`s per color with Go-map indexes. In addition, `Color_CreateUnion` (called from `AstFactory` during optimizations) allocates a fresh union color with *copied* property-name sets each time. A TypePool-shaped Color (ID slices, interned unions) is a local change.
4. **Closure is the only checker studied that does not hash-cons structural types.** Sorbet, TS, Hack, rustc, Swift and go/types all do. Closure creates a new `UnionType` for every `?Foo`, `Foo` (implicitly nullable) and `T=` annotation, a new `TemplatizedType` + `TemplateTypeMap` per `Array<string>`, and a new `FunctionType` per call of a generic function (`type_inference.go:2746`).
   - Typed-full has 52.6 K source annotations plus 19.1 K in the default externs, but only about 6–8 K distinct annotation strings. Roughly 24 K source annotations allocate a union, templatized, function or record type.
   - Exact-key interning of *resolved* unions and templatized types is a credible 2–4 % cut in peak heap, with medium risk. Measure the duplication first (§5).
5. Peak-versus-GC realism. Typed-full peaks *during* type checking, so (1) and (3) do not lower the 552 MB peak; they lower the post-check heap. Under GOGC=100, total mark work is roughly allocation × (scannable bytes / live bytes), so freeing retained pointerful data mostly reduces memory, not GC CPU. GC CPU falls only when data becomes *noscan* (item 2 plus the int32-AST work) or allocation falls (items 3 and 4). The peak itself falls only through fewer or smaller type and scope objects (items 4 and 6 and the in-flight struct shrinking).

---

## 1. Per project: the techniques that matter

### 1.1 Sorbet (C++; `core/`)

- **Everything global is a 32-bit index into a flat table owned by `GlobalState`.**
  - The tables are values, not pointers: `std::vector<ClassOrModule> classAndModules; std::vector<Method> methods; std::vector<Field> fields; std::vector<TypeParameter> typeMembers, typeParameters` (`core/GlobalState.h` ~L795).
  - Sizes are asserted: `CheckSize(ClassOrModule, 128, 8)`, `Method` 136 B, `Field` 56 B, `TypeParameter` 56 B (`core/Symbols.h`).
  - Handles are 4 B: `ClassOrModuleRef`, `MethodRef`, `FieldRef`, `TypeMemberRef`, `TypeParameterRef`, and a generic `SymbolRef` with `KIND_BITS = 3` and 29 bits of id (`core/SymbolRef.h` L381–410).
  - Class members are `UnorderedMap<NameRef, SymbolRef>`, 8-byte pointer-free entries.
  - Method parameters are `InlinedVector<ParamInfo, EXPECTED_METHOD_PARAMS_COUNT>`, stored inline by value. Locations are `InlinedVector<Loc,2>` with `Loc` = 12 B.
- **Names are interned into three typed tables behind one 32-bit `NameRef` with 2 kind bits** (`core/NameRef.h` L88–113):
  - `UTF8Name` = 16 B `string_view` into 4 KiB pages of `StableStringStorage`;
  - `ConstantName` = 4 B (the `NameRef` of the original name);
  - `UniqueName` = 12 B (original, counter, kind).
  - Lookup goes through an open-addressing `NameHash` of `Bucket{uint32 hash; uint32 rawId}` (8 B, pointer-free), sized `2*nextPow2(n)` (`core/GlobalState.h` L45–100).
- **`TypePtr` is a 64-bit tagged word that inlines the commonest types** (`core/TypePtr.h`). The low 8 bits hold the tag and bit 8 means "not inlined"; bits 16–63 hold either a refcounted pointer or an inline payload.
  - Inlined, never allocated: `ClassType` (just a `ClassOrModuleRef`, `CheckSize(ClassType,8,8)`), `NamedLiteralType` (a `NameRef` plus a kind byte), `SelfTypeParam`, `AliasType`, `SelfType`, `BlamedUntyped` (`core/Types.h` L383, L623, `make_type<ClassType>` = `TypePtr(Tag::ClassType, ref.id())`).
  - Heap types are small, refcounted and immutable: `OrType` is 24 B (refcount plus two `TypePtr`s, a binary tree rather than a list); `AppliedType` is 32 B.
- **Per-method results are not kept.** Inference runs per method on a per-method CFG, and locals are `LocalVariable{NameRef; uint32 unique}` = 8 B (`core/LocalVariable.h`), not symbols. The CFG, environments and per-expression types are discarded after the method is checked.
  - `GlobalState` is frozen (`freezeSymbolTable`/`freezeNameTable`) and shared read-only by all workers. Only the symbol table stays alive, so memory scales with *declarations*, not with expressions.
- Transferable to Go: index handles and value tables (the GC does not scan a `[]Method` of pointer-free structs), 8-byte pointer-free hash buckets, and "inline the common case in the handle". Not transferable: tagging bits into a pointer. Go's GC must see real pointers, so tags only work in ID handles (`uint32`), not in pointer or interface words.

### 1.2 TypeScript: JS checker and typescript-go (`internal/checker`, `internal/core`)

- **Types carry a dense id and flags.**
  - JS: `{flags, id, symbol, ...}` with `typeCount++`.
  - Go: `type Type struct { flags TypeFlags; objectFlags ObjectFlags; id TypeId; symbol *ast.Symbol; alias *TypeAlias; checker *Checker; data TypeData }` (`types.go` L681). The concrete data struct embeds `Type`, and `data` is the self-interface, the same pattern as jsgo's `JSTypeBase.self`.
- **Structural types are interned by keys built from the ids of their parts.**
  - JS: `unionTypes: Map<string, UnionType>` keyed by `getTypeListId(types)`, ids run-length encoded as `"12:4,19"`. Per-target `instantiations: Map<string, Type>` use the same id-list keys. Relation caches are `Map<string, RelationComparisonResult>` keyed `"srcId,tgtId,..."`.
  - Constituents are *sorted by id* before interning, so the key is canonical.
  - **tsgo replaces the string keys with `CacheHashKey = xxh3.Uint128`.** `keyBuilder` serializes the ids into a 192-byte inline buffer and hashes them (`checker.go` L17471–17606). Then `unionTypes`, `intersectionTypes`, `tupleTypes`, `indexedAccessTypes`, `subtypeReductionCache`, and `Relation.results map[CacheHashKey]RelationComparisonResult` (`relater.go` L98) all use **pointer-free keys**.
  - Derived-type memo: `cachedTypes map[CachedTypeKey]*Type` with `CachedTypeKey{kind, typeId}` covers apparent, widened, awaited, permissive-instantiation and similar variants (`checker.go` L117–149). Again the key is pointer-free.
  - tsgo does **not** compare on a hash hit; it trusts the 128-bit hash.
- **Per-node and per-symbol analysis results live in side tables indexed by lazily assigned dense ids** (`links`).
  - JS: `nodeLinks: NodeLinks[]` indexed by `getNodeId(node)`, and `symbolLinks[]`.
  - tsgo: `core.LinkStore[K,V]` = `map[K]*V` + `Arena[V]`, and `PagedLinkStore[V]` = pages of 256 values indexed by `id>>8`, with a slice page table for low ids and a map above that (`internal/core/linkstore.go`).
  - `GetNodeId` assigns ids with an atomic counter only for nodes that need links (`ast/utilities.go` L22).
  - There are about 25 link kinds, each with its own small struct (`SymbolNodeLinks{resolvedSymbol}`, `TypeNodeLinks{resolvedType, outerTypeParameters}`). Cold data never bloats the node.
  - Types are *not* stored per expression node. They are recomputed or cached only where needed.
- `Arena[T]` (`internal/core/arena.go`) is a growable `[]T` that hands out `&data[i]` and doubles by size class. This is the same as jsgo's registry slabs.

### 1.3 Flow (now a Rust port: `rust_port/crates/flow_type_sig`, `flow_heap`)

- **"Compact tables" via a copying-collector-like pass** (`compact_table.rs` header comment): build nodes, mark the ones reachable from the exports, compact (assign 0-based indices in insertion order) and copy into `Table<T>` (a boxed slice). References become `Index<T>` (a `usize` newtype).
- A file signature is `Module { local_defs: Table<PackedDef<Index<Loc>>>, remote_refs: Table<...>, pattern_defs, patterns, module_refs }` (`packed_type_sig.rs`). Even *locations* are indices into a per-file loc table (`PackedALocTable`).
- The heap keeps type signatures as **serialized bytes** and calls `deserialize_type_sig(file, bytes)` on demand (`flow_heap/src/parsing_heaps.rs` L325). Most per-file data is held as `Arc<...>` blobs.
- Closure's `SerializeTypesToPointers` → `TypePool` → `ColorPool` already follows this "compact into index tables, insertion-ordered" model (§1.6).

### 1.4 Hack / hackrs

- **Hash-consing:** `hcons::Hc<T>(Arc<T>)` with a global `Conser<T> { table: DashMap<u64, Weak<T>> }`. Equality and hash are pointer-based, and `Conser::gc()` drops dead weak entries.
  - On a 64-bit FNV collision it only does `debug_assert!(x == *rc)`, marked "TODO: handle collisions" (`hphp/hack/src/hcons/lib.rs`).
- **Provenance is kept out of the shared part.** `pub struct Ty<R: Reason>(R, Hc<Ty_<R>>)` (`hackrs/ty/decl/ty.rs` L212): the reason/position sits *outside* the hash-consed node, so structurally equal types share one node even when they come from different places.
- `hh_shared.c` keeps decls in a shared-memory hash table as serialized values ("data with the same key are considered equivalent ... if data is missing it can be recomputed"), outside the OCaml GC heap. Folded decls are computed lazily.

### 1.5 Pyright and mypy

- Pyright, `TypeBase { category; flags; props; cached; shared; priv }` (`analyzer/types.ts` L201):
  - `shared` holds class details shared by all specializations (`ClassDetailsShared`);
  - `priv` holds the per-specialization part (`typeArgs`);
  - `props` holds rare optional fields;
  - `cached` holds memoized data that is never cloned.
  - A specialization copies only a small header and `priv`.
- Pyright keeps per-node results in side tables keyed by node id: `typeCache = new Map<number, TypeCacheEntry>()` (`typeEvaluator.ts` L707/802). Under heap pressure it **drops them** and recomputes lazily ("Emptying type cache to avoid heap overflow", `cacheManager.ts`; `program.ts` limits to about 750 K entries). This only works because the evaluator is demand-driven.
- mypy keeps per-expression types in a side table, `dict[Expression, Type]` on the TypeChecker. After each module is checked it calls `free_state()` and `free_tree(self.tree)`, dropping function bodies and the type map (`mypy/build.py` L3587–3591).
  - Every node and type class has `__slots__` and is compiled to a C struct by mypyc.
  - The cache format is a fixed binary format in which the most common `Instance` types are a single tag (`mypy/cache.py` header): the "inline the common case" idea again.

### 1.6 Closure's own designs (in jsgo already)

- **The TypedAST `TypePool` is an ID-based type graph** (`jscomp/proto_types.go`):
  - `TypeProto{object|union}`, `UnionTypeProto.unionMember []int32`;
  - `ObjectTypeProto.prototype/instanceType []int32` (type pointers), `ownProperty []int32` (string-pool indices);
  - `SubtypingEdge` pairs.
  - Type pointers reserve the first 10 offsets for axiomatic colors (`TypePointers_OFFSET_TO_AXIOMATIC_COLOR`): the "inline the natives" trick.
- jsgo already stores the serializer's pointer *inline* (`JSTypeBase.SerialGen/SerialPointer`) instead of an `IdentityHashMap`, interns property names (`Property.nameID rhinoStringID`, `PropertyMap.entries` scanned by ID, `ext.byID map[rhinoStringID]*Property` past 24 entries), tags type kinds in `objFlags`, and allocates types from per-kind registry slabs.
- `colors.Color` is the post-check summary. Inside a `ColorPool`, colors are canonical per `ColorId`, have value equality (AutoValue) and memoize their hash code.
- Elsewhere: rustc interns every `Ty<'tcx>` (an arena-allocated `WithCachedTypeInfo<TyKind>` with cached flags) and every generic-argument list. Swift uniques canonical types in `ASTContext`. go/types' `Context` dedupes identical generic instances by `typeMap map[string][]ctxtEntry` (hash → candidates with an exact `Identical` check). Filliâtre & Conchon's "type-safe modular hash-consing" adds a unique tag per value for O(1) hashing and ID-keyed maps. The common lesson is to **hash-cons structural types, give them dense IDs, key every cache by IDs, and keep provenance (positions) outside the shared node.**

---

## 2. jsgo today (what the techniques map onto)

| jsgo structure (file) | Layout now | Technique that applies |
|---|---|---|
| `JSTypeBase` (`rhino/jstype_js_type.go:248`) | 64 B: `self` interface 16, `resolveResult` interface 16, `templateTypeMap*`, `registry*`, `SerialGen`/`SerialPointer`, flags | dense `id uint32` replaces `SerialGen/SerialPointer`; struct shrinking is in flight |
| `UnionType` (`jstype_union_type.go:60`) | 128 B with `inlineAlternates [2]`; built by `UnionType_Builder`, **never interned**; alternates in insertion order; `AppendTo` sorts names; `HashCode` is a product, so order-independent | exact-key hash-consing of resolved unions; per-type nullable/optional cache |
| `TemplatizedType` + `TemplateTypeMap` (`jstype_templatized_type.go:42`, `jstype_template_type_map.go:44`) | new per `CreateTemplatizedTypeOfList` (no interning); a lossy 4096-entry `templateTypeMapCache` exists | intern by (base id, arg ids) |
| `FunctionTypeBase` (`jstype_function_type.go:133`) | 224 B + `functionTypeExt` + `[]*FunctionType_Parameter` (24 B objects plus pointer slots) | per-call instantiation memo (go/types `Context`); parameters by value (Sorbet `InlinedVector<ParamInfo>`) |
| `NamedType` (`jstype_named_type.go:69`) | 208 B per *unresolved reference occurrence*; Java comment "TODO: Each instance should support named type creation using interning" | keep provenance outside (Hack `Ty(R, Hc<Ty_>)`): fields used only for resolution go into an ext freed after resolution (struct shrinking) |
| `Property` (`jstype_property.go:71`), `PropertyMap` | 40 B, slab-allocated, keyed by `rhinoStringID`; `ext.byID` is a Go map with pointer values | already Sorbet-style; optional int32 open-addressing index |
| registry property indexes (`jstype_js_type_registry.go`) | `map[string]*ObjectLinkedHashSet[JSType]` ×2 | key by `rhinoStringID`, values `[]uint32` type IDs (noscan, insertion order kept) |
| `propertySymbolList` (`jstype_property.go:46`) | **global, never cleared, pins the registry** | per-registry table (finding 1) |
| `Node.typeOrColor unsafe.Pointer` (`rhino/node.go:597`) | `*JSTypeBase`, `*Color`, `*JSDocInfo` or `*nodeExtBox` | a 32-bit handle (Sorbet `TypePtr`/`SymbolRef`, TS `TypeId`, Flow `Index`) |
| `colors.Color` (`colors/color.go:36`) | 80 B + 4 sets (`ObjectLinkedHashSet[*Color]` ×3, `LinkedHashSet[string]`); unions copy the elements' property sets | TypePool-shaped: `[]uint32` string IDs and color indices; interned unions |
| `TypedScope`/`TypedVar` (`jscomp/typed_scope.go:144`, `typed_var.go:39`) | `vars LinkedHashMap[string,*TypedVar]`, `TypedVar.type_ rhino.JSType` (16 B interface) | ID handle for the type; ID/hash index (in flight) |

---

## 3. Recommendations, ranked by impact on peak / pointer count (and simplicity)

### R1. Free the type graph for real after colors: fix `propertySymbolList` and add a leak guard
- **Files:** `rhino/jstype_property.go` (`propertySymbolList`, `registerPropertySymbol`, `propertySymbolFromID`); `jscomp/compiler.go:2160` (`ClearJSTypeRegistry`).
- **Mechanism:** move the list onto `JSTypeRegistry` (`r.propertySymbols []*KnownSymbolType`). Decode through `p.type_.GetJSTypeBase().registry`; every `Property` has a non-nil type. As a stopgap for the CLI, nil the global in `ClearJSTypeRegistry`; that is safe only if no other compilation runs in the process.
- **Guard:** in memgate or a test, take `weak.Make(registry)` (Go 1.24 is in `go.mod`). After `convertTypesToColors`, run `runtime.GC()` and assert `.Value() == nil`. That catches this and any future retention path, such as a pooled builder or a cached `AstAnalyzer` (`compiler.go:1695` holds `typeRegistry` until the next `GetAstAnalyzer`).
- **Gain:** the post-check (mid-loop) live heap drops by the retained types, plus the declaration JSDoc reachable from `Property.docInfo`/`objectDocInfo`, plus dead AST reachable from `FunctionType.source`/`Property.propertyNode`. Measured typed-app share: 18 % of mid-loop live. For typed-full I estimate 50–90 MB of pointerful mid-loop heap. The 552 MB peak is unchanged, because it occurs during checks. GC CPU changes little under GOGC=100 (more, smaller cycles) but RSS headroom and cache footprint improve. It also fixes the cross-compilation leak.
- **API, risk, effort:** passes untouched. No byte-identity risk; symbol IDs stay stable within a compilation. About 30 lines. Default mode.

### R2. Make the AST↔type/color link a 32-bit handle, with dense type and color IDs (enables a noscan `rhino.Node`)
- **Files:** `rhino/node.go` (`typeOrColor` and accessors: `GetJSType`/`GetJSTypeBase`/`SetJSType`/`GetColor`/`SetColor`/`TakeTypesAndJSDocForColorConversion`, clone paths); `rhino/jstype_js_type_registry.go` (the `alloc*` slab functions assign `id`); `jscomp/serialization_color_pool.go` (a ShardView's `offsetToColor` is already the color table); `colors/color.go` (add `index uint32`); `jscomp/ast_factory.go` and the TypedAST deserializer (register new colors).
- **Mechanism:**
  - At allocation, give every `JSTypeBase` an `id uint32` and append to `registry.types []*JSTypeBase`. The registry is created first, so natives get IDs 1..N in ordinal order. This is Sorbet's and TS's dense id, and Closure's axiomatic-offset trick.
  - Replace `SerialGen`/`SerialPointer` with that id (the serializer keeps a `[]int32` indexed by type id), which frees 4 B per type.
  - Node stores `typeID uint32` while checks run and `colorIdx uint32` afterwards, with a kind bit already in `kindFlags`. Store JSDoc as a separate `uint32` handle (coordinate with the JSDoc-compaction team). Then a node carrying both JSDoc and a type no longer needs `nodeExtBox`, and the node holds 8 B of handles and no pointer.
  - Resolve handles through a per-compilation table. jsgo already runs process-global pools (the string pool), so either a global `atomic.Pointer[typeTable]` swapped per compilation, or a table reached through the node slab header.
- **Gain:** by itself 0 bytes (the 8 B slot becomes 4+4). Its value is that it removes the last semantic pointer that blocks the in-flight "pointer-free int32-linked AST": the AST is 32 % of typed-full's peak (160 MB) and 45–85 % of GC scan work (gclab). Also, `JSTypeReconserializer.typeToRecordCache map[JSType]…`, `standardColorByPtr`, `objectDocInfo map[*ObjectTypeBase]`, the identity maps in `InvocationTemplateTypeMatcher`, and `DisambiguateProperties`' `*Color`-keyed maps can all become `[]T` indexed by id (noscan keys, simpler code).
- **API:** pass APIs are unchanged; `GetJSType() JSType` still returns the interface.
- **Risk:** none to output, provided ids are never used for ordering, only identity. Identity hash codes stay pointer-based, as now.
- **Effort:** medium, all in the substrate. Default mode.

### R3. TypePool-shaped Colors, with interned union colors
- **Files:** `colors/color.go`, `colors/color_registry.go`, `jscomp/serialization_color_pool.go`, `jscomp/ast_factory.go` (6 `Color_CreateUnion` calls), `jscomp/disambiguate_color_find_property_references.go`:184/195, `disambiguate_disambiguate_properties.go`:320/339, `serialization_color_serializer.go`:209/346 (the only `GetOwnProperties` consumers).
- **Mechanism:**
  - `ownProperties` becomes an insertion-ordered `[]uint32` of interned string IDs. For more than ~16 entries, add a small sorted copy or an int32 open-addressing index for `Contains`. Callers already hold the property name's string ID in `Node.payload`.
  - `prototypes`, `instanceColors` and `unionElements` become `[]uint32` color indices, or `[]*Color` with nil for the very common empty case and a shared one-element backing.
  - `ColorRegistry.colorToDisambiguationSupertypeGraph` becomes adjacency lists of `uint32`.
  - Intern union colors in a per-compilation table keyed by the **ordered** element-index list, with an exact compare. Both `ColorPool` reconciliation and `AstFactory` (`createOr`/`createHook`/…) then return the same object instead of copying every element's property names.
- **Gain:** per object color, ~40 B of header per non-empty set plus ~45 B per property name with a map index, down to 4 B per name. I estimate colors at 15–20 MB for typed-full (about 30 K colors at ~500 B), shrinking to 3–5 MB and mostly noscan. That is roughly −3 % of the mid-loop heap, plus the `AstFactory` union churn. The peak is unchanged.
- **Byte identity:** keep insertion order exactly. `Equals`/`HashCode` stay set-semantic (sort both sides to compare; the hash is already memoized). Interning by *ordered* elements makes it a pure dedup of value-equal objects that Java's value-equality maps would also merge.
- **API:** `GetOwnProperties` returns a small `colors.PropSet` with `Contains`/`All`/`Size`, so ~10 call sites change. Passes otherwise stay 1:1. Effort: low–medium. Default mode.

### R4. Hash-cons resolved structural JSTypes (measure first; start with the nullable/optional wrappers)
- **Files:** `rhino/jstype_js_type_registry.go` (`CreateNullableType` L2312, `CreateOptionalType`, `CreateUnionTypeOfList`, `CreateTemplatizedTypeOfList` L2630), `rhino/jstype_union_type.go` (`UnionType_Builder.Build` L1273), `jscomp/type_inference.go:2746` (per-call `TemplateTypeReplacer` output).
- **Mechanism, in order of risk:**
  - **(a) Derived-type memo** (tsgo `cachedTypes{kind,typeId}`): `registry.nullableOf[id]`, `optionalOf[id]`, `nullOrUndefOf[id]`, only for resolved, non-proxy `T`. This covers about 15 K of the ~17 K annotation unions in typed-full sources (implicitly nullable names 8.6 K, `?X` 3.0 K, `T=` 3.1 K) plus the externs.
  - **(b) General union interning** when `alternatesResolvedBeforeBuild` holds, so the union can never be rebuilt. Key = the ordered alternate id list (TS-style bytes, or a 64-bit hash → candidate list as in go/types `Context`), **with an exact compare**. tsgo and Hack trust the hash; byte identity cannot.
  - **(c) `TemplatizedType` interning** by (base id, arg ids) for resolved args. This also dedupes its `TemplateTypeMap`.
  - **(d) Per-call generic instantiation memo** keyed by (generic `FunctionType` id, ordered inferred binding ids), only when the replacer touched no unresolved or proxy types.
  - Drop all tables in `ClearJSTypeRegistry`.
- **Gain (estimate; verify with §5):** unions ~100 K × ~140 B ≈ 14 MB with ~70–90 % duplicates; templatized types ~20–40 K × ~230 B ≈ 5–9 MB at ~70 %; per-call instantiations a few MB at ~50 %. Total **~15–20 MB at the typed-full peak (≈3–4 %)**, plus a matching share of the "types & colors 5 %" allocation churn. That is less than the AST work, but it is one of the few levers that lowers the *checks-phase* peak.
- **Byte-identity risk: medium.**
  - Identity checks exist: `base.JSCompObjects_Identical(fnSlotType, slot.GetType())` in `LinkedFlowScope` join (`linked_flow_scope.go` ~L774/794), and the `LinkedIdentityHashSet[JSType]` cycle guard in `invocation_template_type_matcher.go:42`. Interning makes more pairs identical. Outcomes should be structurally the same (the union builder dedupes equal alternates), but that cannot be proven by reading.
  - Mutable per-instance state (`hashCodeInProgress`, `visited`, `templateCheckResult`, `objectDocInfo` on object types) is reentrancy or memo state. Sharing it is safe single-threaded; it is unsafe if inference ever runs in parallel.
  - Structural `HashCode` and `InvalidatingTypes` (keyed by `equals`) are unaffected.
  - Keep alternate order exact, which `AppendTo` sorting already hides in messages.
- **API:** registry-internal; passes 1:1. Effort: (a) low, (b)+(c) medium, (d) medium.
- **Mode:** default if the memgate byte-identity corpus passes; otherwise behind `--fast_mode`.

### R5. Pointer-free, ID-keyed side tables in the registry and checker caches
- **Files:** `jstype_js_type_registry.go` (`eachRefTypeIndexedByProperty`, `nonRefTypesIndexedByProperty`, `objectDocInfo`, `leastSupertypeCache`, `templateTypeMapCache`), `jstype_property_map.go` (`propMapExt.byID`).
- **Mechanism:** after R2, key by `rhinoStringID` or type id; store values as `[]uint32` id lists in insertion order (dedupe through a structural-hash index that still compares with `equals`); make caches `[N]struct{a, b, res uint32}`; give `byID` an int32 open-addressing index into `entries`.
- **Gain:** small (≈2–4 MB, mostly pointer words), plus simpler code. No identity risk if dedupe semantics are kept. Low effort. Default mode.

### R6. Scope-side alignment (feeds the in-flight scope work)
- `TypedVar.type_` as a type id (4 B) or `*JSTypeBase` (8 B, interface rebuilt from `objFlags` as `Node.GetJSType` already does): 16 → 4–8 B per var.
- Key `TypedScope.vars` by the existing `Var.nameID` hash: no `string` in the slots, an int32 index.
- Gain ~3–5 MB at the typed-full peak (scopes ≈ 28 MB). Byte-identical (insertion order kept). Passes 1:1.

---

## 4. Fast-mode-only ideas, and rejected ideas

**`--fast_mode` candidates (deterministic, not Java-identical):**
- Canonicalize union alternate order (sort by type id, as TS does) before interning. More hits, but `GetAlternates()` order changes iteration-dependent results such as warning text order and color union element order.
- Intern `NamedType` by (scope root, reference) and keep only the first occurrence's location. The Java TODO points at this; Hack keeps provenance outside the shared node. It changes "unknown type" diagnostics and locations.
- Fuse `TypeCheck` into `TypeInference` per function and drop each function's local `TypedScope`/flow state right afterwards (Sorbet/mypy style). This could cut a large part of the scope share at the peak, but warning order and dedup change unless errors are buffered and re-sorted exactly.

**Rejected:**
- *Fully ID-based JSType graph* (Sorbet/Flow style: every internal type-to-type edge a `uint32`). It means rewriting ~34 K lines of `rhino/jstype_*` and every pass that walks types, and types still contain strings, slices, maps and funcs, so they never become noscan. At most about half of the 15 % types share is recoverable, for the largest possible rewrite.
- *Tagged pointers or inline types in Go interface or pointer words* (Sorbet `TypePtr`). Go's GC requires real pointers in pointer slots, and unsafe tagging crashes it. The idea survives only inside `uint32` handles (R2), where native types can be small constant ids.
- *Droppable, recompute-on-demand type caches* (Pyright) and *serialized off-heap decls* (Hack/Flow). Closure's inference is eager and its results are the AST annotations; nothing can be recomputed lazily. The TypedAST already covers the "serialized blob" niche for runtime libraries.
- *Trusting 64/128-bit hashes as intern keys* (tsgo, Hack `Conser`). A collision silently merges types, which is unacceptable under byte identity. Use the hash only to find candidates and compare exactly.
- *A separate per-node type side array* (tsgo links, mypy's type map). jsgo already has a dense inline slot in every node, so a side array adds an indirection without saving bytes. The right move is to shrink the inline slot to an id (R2) and keep rare data (`typeBeforeCast`) in the ext box.
- *Go's `unique` package for interning:* `Handle[T]` needs `comparable` fixed-size keys (not variable-length alternate lists), and its weak semantics are unnecessary because the tables die with the registry.

---

## 5. Measurement to do before R3 and R4 (cheap, via the existing `JSGO_PEAK_HEAP_PROFILE` hook)

1. **Per-kind bytes at the typed-full peak.** The registry slabs make inuse_space attribute per kind: `allocUnionType`, `allocTemplatizedType`, `allocTemplateTypeMap`, `allocFunctionType`/`Ext`/`Param*`, `allocProperty`, `allocNamedType`, `allocPrototypeObjectType`, `allocInstanceObjectType`, `jsTypeSliceSlab`, and the `TypedScope`/`TypedVar` arenas.
2. **Duplication ratios.** Behind a debug env var, count `UnionType_Builder.Build` results versus distinct exact keys (ordered alternate base pointers, resolved only), and the same for `NewTemplatizedType` and the `type_inference.go:2746` instantiations. R4's gain is these counts times the sizes.
3. **Retention check.** After `convertTypesToColors`, `runtime.GC()`, then confirm `weak.Make(registry).Value() == nil` and measure inuse by `rhino.*` alloc sites in the next pass. This confirms R1 and gives its exact size. Repeat after `clearTypedScopeCreatorPass` for scopes.
4. **Color inventory.** Count colors per pool and the `ownProperties` size histogram. Count `Color_CreateUnion` calls from `AstFactory` over a run versus distinct ordered element lists. This sizes R3.
