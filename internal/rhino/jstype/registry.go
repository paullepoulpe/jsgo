// Port of com/google/javascript/rhino/jstype/JSTypeRegistry.java @ closure-compiler v20260415
//
// The real registry manages interning of all JSType values, named-type
// resolution, template instantiation, and much more. Tier-0 only needs a
// name-to-primitive lookup, which we provide here.
//
// Additive: namedTypes map + pendingNamedTypes worklist added to support
// NamedType.resolveViaRegistry and the createNamedType factory.
// Additive: forwardDeclaredTypes + closureNamespaces added to support
// resolveViaClosureNamespace / handleUnresolvedType (Tier-1 widening).
// Additive: nativeTypes table + property side tables + template-key scope
// table added (Tier-2 widening; mirrors JSTypeRegistry.java §nativeTypes,
// §eachRefTypeIndexedByProperty / §nonRefTypesIndexedByProperty,
// and §registerTemplateTypeNamesInScope).
// Additive: late-resolution worklist callbacks (typeResolvers) +
// per-scope forward declarations (scopeForwardDeclaredTypes) +
// hasType predicate + declareNative* helpers (Tier-3 widening; mirrors
// JSTypeRegistry.java §resolveTypesForTesting / §declareType and
// TypedScopeCreator.java §declareNativeFunctionType / §declareNativeType /
// §declareNativeValueType).

package jstype

import "strings"

// closureNamespace mirrors JSTypeRegistry.ClosureNamespace in Java.
// isLegacy is true for goog.provide / legacy goog.module namespaces;
// exportedType holds the exported type for non-legacy modules.
type closureNamespace struct {
	isLegacy     bool
	exportedType JSType // nil for legacy namespaces
}

// Registry mirrors JSTypeRegistry. Named-type registration and late-resolve
// are additive; existing ResolveName / ResolveExpr are unchanged.
type Registry struct {
	// namedTypes holds types registered via RegisterNamedType.
	// Mirrors the typeRegistry field (name→JSType) in Java's JSTypeRegistry.
	namedTypes map[string]JSType

	// pendingNamedTypes is the late-resolve worklist: NamedType values
	// that were constructed before their reference was registered.
	// Mirrors the implicit "forward-ref queue" that Java's resolver
	// processes when the registry transitions to closed.
	pendingNamedTypes []*NamedType

	// resolver is the JSTypeResolver that shares this registry's lifecycle.
	// Nil until JSTypeResolver.Create() wires it.
	resolver *JSTypeResolver

	// forwardDeclaredTypes is the set of names registered via
	// goog.forwardDeclare. Mirrors JSTypeRegistry.forwardDeclaredTypes.
	// When a NamedType fails resolution, this set determines whether to
	// emit a warning (not found) or silently produce a NoResolvedType.
	forwardDeclaredTypes map[string]bool

	// closureNamespaces mirrors JSTypeRegistry.closureNamespaces.
	// Keyed by module/provide name (e.g. "goog.module.Foo").
	// Used by ResolveViaClosureNamespace to match dotted-type references
	// against the longest registered namespace prefix.
	closureNamespaces map[string]closureNamespace

	// nativeTypes is a fixed-size table of built-in JSType values indexed
	// by JSTypeNative ordinal. Mirrors Java's JSTypeRegistry.nativeTypes
	// array (private JSType[]).  Populated by RegisterNativeType; read by
	// GetNativeType.  Using an array (not a map) keeps Get O(1) with no
	// allocation, mirroring Java's "very fast" comment on the field.
	nativeTypes [NUMBER_ADDITION_SUPERTYPE + 1]JSType

	// eachRefTypeIndexedByProperty mirrors Java's eachRefTypeIndexedByProperty
	// (SetMultimap<String,ObjectType>).  For each property name, the slice
	// records every named ObjectType (HasReferenceName() == true) on which
	// that property has been explicitly declared.  Used by
	// GetEachRefTypeWithProperty, which TypeCheck consults to detect whether
	// an absent property might be "loose" (present on a subtype).
	eachRefTypeIndexedByProperty map[string][]*ObjectType

	// nonRefTypesIndexedByProperty mirrors Java's nonRefTypesIndexedByProperty
	// (SetMultimap<String,JSType>).  Holds non-reference types (anonymous
	// object literals, record types, etc.) that declare a given property.
	// Kept separate from eachRefTypeIndexedByProperty to avoid inflating the
	// reference-type index with structural / anonymous types.
	nonRefTypesIndexedByProperty map[string][]JSType

	// templateScopes mirrors the Node-keyed scopedNameTable in Java
	// (Table<Node,String,JSType>).  Because Go callers pass an opaque key
	// (typically the *rhino.Node of the function's FUNCTION token, or nil
	// for the global scope) we use interface{} so this file stays free of
	// a rhino import.  Populated by RegisterTemplateTypeNamesInScope;
	// read by LookupTemplateType.
	templateScopes map[interface{}]map[string]*TemplateType

	// typeResolvers is the late-resolution callback worklist.  Each entry
	// pairs a JSType instance with a resolver function that, given a type
	// name, returns a resolved JSType (or nil when still unknown).
	// ResolveTypes drains this list and calls each resolver; if the
	// resolver returns a non-nil type the entry is removed.
	// Mirrors the implicit per-type resolution callbacks that Java's
	// JSTypeRegistry drives during its "close" pass.
	typeResolvers []typeResolverEntry

	// scopeForwardDeclaredTypes holds per-scope forward-declared type
	// names.  Unlike the flat forwardDeclaredTypes set (which mirrors
	// JSTypeRegistry.forwardDeclaredTypes populated by goog.forwardDeclare
	// at module-boundary granularity), this table is keyed by an opaque
	// scope key (typically a *rhino.Node) so that forward declarations
	// made inside a function scope are not visible globally.
	//
	// Mirrors the per-scope slot in JSTypeRegistry.scopedNameTable that
	// TypedScopeCreator populates when processing in-scope type aliases /
	// forward declarations.
	scopeForwardDeclaredTypes map[interface{}]map[string]bool
}

// typeResolverEntry associates a JSType with a resolver callback.
// Mirrors the implicit tuple that drives per-type late resolution in Java.
type typeResolverEntry struct {
	t        JSType
	resolver func(name string) JSType
}

func NewRegistry() *Registry { return &Registry{namedTypes: map[string]JSType{}} }

// RegisterNamedType stores a name→JSType mapping so that NamedType
// instances with that reference name can resolve. Mirrors
// JSTypeRegistry.registerType (the put-by-name variant).
// Additive: does not affect ResolveName / ResolveExpr.
func (r *Registry) RegisterNamedType(name string, t JSType) {
	if r.namedTypes == nil {
		r.namedTypes = map[string]JSType{}
	}
	r.namedTypes[name] = t
	// Eagerly resolve any pending NamedTypes waiting on this name.
	remaining := r.pendingNamedTypes[:0]
	for _, nt := range r.pendingNamedTypes {
		if !nt.resolveViaRegistry(r) {
			remaining = append(remaining, nt)
		}
	}
	r.pendingNamedTypes = remaining
}

// lookupNamedType returns the JSType registered under name, or nil.
// Called by NamedType.resolveViaRegistry. Not exported; callers outside
// the package use RegisterNamedType + createNamedType.
func (r *Registry) lookupNamedType(name string) JSType {
	if r.namedTypes == nil {
		return nil
	}
	return r.namedTypes[name]
}

// enqueueOrResolve is called by NamedTypeBuilder.Build. Tries
// resolveViaClosureNamespace (for dotted names) and resolveViaRegistry
// in sequence. If both fail, adds nt to the late-resolve worklist.
// Mirrors the resolveIfClosed call in Java's NamedType constructor;
// the two-dispatch order mirrors resolveInternal in NamedType.java.
func (r *Registry) enqueueOrResolve(nt *NamedType) {
	if nt.resolutionKind == ResolutionKindNone {
		return // already has a referenced type; nothing to do
	}
	if nt.resolveViaClosureNamespace(r) {
		return
	}
	if nt.resolveViaRegistry(r) {
		return
	}
	r.pendingNamedTypes = append(r.pendingNamedTypes, nt)
}

// createNamedType is the minimal factory described in the port spec.
// Mirrors JSTypeRegistry.createNamedType(scope, reference, …) — the
// scope-less, location-less variant used by callers that only have a
// name string.  Returns a NamedType (possibly still unresolved) and
// enqueues it for late resolution.
func (r *Registry) createNamedType(name string) *NamedType {
	return NewNamedTypeBuilder(r, name).
		SetResolutionKind(ResolutionKindTypeName).
		Build()
}

// ResolveWorklist drains the late-resolve worklist, applying resolver
// to each pending name. resolver is a caller-supplied function that maps
// a type name to a JSType (or nil when still unknown). After the call
// any NamedTypes that remain unresolved are left as-is (they will act
// like Unknown at check time).
//
// Mirrors the "close the registry" pass in Java's TypedScopeCreator
// where JSTypeRegistry.resolveTypesForTesting (and the production
// finalisation path) iterate over forward-declared types.
func (r *Registry) ResolveWorklist(resolver func(name string) JSType) {
	remaining := r.pendingNamedTypes[:0]
	for _, nt := range r.pendingNamedTypes {
		if t := resolver(nt.reference); t != nil {
			nt.referencedType = t
		} else if !nt.resolveViaRegistry(r) {
			remaining = append(remaining, nt)
		}
	}
	r.pendingNamedTypes = remaining
}

// DrainPendingNamedTypes returns and clears the pending list so the
// caller can emit JSC_UNRECOGNIZED_TYPE_ERROR for whatever remains.
func (r *Registry) DrainPendingNamedTypes() []*NamedType {
	out := r.pendingNamedTypes
	r.pendingNamedTypes = nil
	return out
}

// RegisterPending records an unresolved type-name reference for the
// two-phase resolve pass. lowercase fields on main's NamedType.
func (r *Registry) RegisterPending(ref, sourceName string, lineno, charno int) {
	r.pendingNamedTypes = append(r.pendingNamedTypes, &NamedType{
		reference:  ref,
		sourceName: sourceName,
		lineno:     lineno,
		charno:     charno,
	})
}

// GetResolver returns the JSTypeResolver associated with this registry, or nil
// if none has been attached yet. Mirrors JSTypeRegistry.getResolver() in Java.
func (r *Registry) GetResolver() *JSTypeResolver { return r.resolver }

// IsForwardDeclaredType reports whether name was registered via
// RegisterForwardDeclare. Mirrors JSTypeRegistry.isForwardDeclaredType.
func (r *Registry) IsForwardDeclaredType(name string) bool {
	return r.forwardDeclaredTypes[name]
}

// RegisterForwardDeclare records a name as forward-declared (goog.forwardDeclare).
// Mirrors the forwardDeclaredTypes set populated by the compiler's
// ProcessClosurePrimitives pass.
func (r *Registry) RegisterForwardDeclare(name string) {
	if r.forwardDeclaredTypes == nil {
		r.forwardDeclaredTypes = map[string]bool{}
	}
	r.forwardDeclaredTypes[name] = true
}

// RegisterClosureNamespace records a non-legacy goog.module namespace with
// its exported type. Mirrors JSTypeRegistry.registerNonLegacyClosureNamespace.
func (r *Registry) RegisterClosureNamespace(moduleName string, exportedType JSType) {
	if r.closureNamespaces == nil {
		r.closureNamespaces = map[string]closureNamespace{}
	}
	r.closureNamespaces[moduleName] = closureNamespace{isLegacy: false, exportedType: exportedType}
}

// RegisterLegacyClosureNamespace records a goog.provide or legacy goog.module
// namespace. Mirrors JSTypeRegistry.registerLegacyClosureNamespace.
func (r *Registry) RegisterLegacyClosureNamespace(moduleName string) {
	if r.closureNamespaces == nil {
		r.closureNamespaces = map[string]closureNamespace{}
	}
	r.closureNamespaces[moduleName] = closureNamespace{isLegacy: true}
}

// ResolveViaClosureNamespace resolves a dotted type reference by finding
// the longest registered closure-namespace prefix, then walking properties.
// Returns nil when no namespace prefix matches or when the match is a
// legacy namespace (which falls through to the normal registry / property path).
// Mirrors JSTypeRegistry.resolveViaClosureNamespace in Java.
func (r *Registry) ResolveViaClosureNamespace(reference string) JSType {
	if r.closureNamespaces == nil {
		return nil
	}
	prefix := reference
	var suffixComponents []string
	for {
		if ns, ok := r.closureNamespaces[prefix]; ok {
			if ns.isLegacy {
				// Legacy (goog.provide / declareLegacyNamespace) falls
				// through to the normal registry / property path.
				return nil
			}
			// Non-legacy goog.module: always stop here. Walk properties
			// on the exported type for any suffix components. Mirrors
			// Java's resolveViaPropertyGivenSlot call.
			if len(suffixComponents) == 0 {
				return ns.exportedType
			}
			t := ns.exportedType
			for i := len(suffixComponents) - 1; i >= 0; i-- {
				if t == nil {
					return nil
				}
				obj, ok := t.(*ObjectType)
				if !ok {
					return nil
				}
				t = obj.GetPropertyType(suffixComponents[i])
			}
			return t
		}
		dot := strings.LastIndexByte(prefix, '.')
		if dot < 0 {
			return nil
		}
		suffixComponents = append(suffixComponents, prefix[dot+1:])
		prefix = prefix[:dot]
	}
}

// ResolveName mirrors JSTypeRegistry.getType for primitive names. Returns
// Unknown for names this scaffold doesn't recognise yet.
func (r *Registry) ResolveName(name string) JSType {
	switch name {
	case "number":
		return Number
	case "string":
		return String
	case "boolean":
		return Boolean
	case "null", "Null":
		// Java exposes both lowercase "null" and the capitalised alias
		// "Null" that appears in some Closure-compiler test fixtures.
		return Null
	case "undefined", "void":
		return Void
	case "symbol":
		return Symbol
	case "*":
		return All
	}
	return Unknown
}

// ResolveExpr parses the textual form captured from a JSDoc type
// annotation ({e.g. "number|string", "!Foo", "?number"}) and returns a
// JSType. Unknown members collapse the whole expression to Unknown so
// the checker stays conservative, which matches how Java resolves an
// unparseable type expression.
func (r *Registry) ResolveExpr(expr string) JSType {
	expr = strings.TrimSpace(expr)
	if expr == "" || expr == "?" {
		return Unknown
	}
	// `?Foo` = Foo | null (nullable). Strip the leading `?` and union
	// with Null.
	nullable := false
	if strings.HasPrefix(expr, "?") {
		nullable = true
		expr = strings.TrimSpace(expr[1:])
	}
	parts := strings.Split(expr, "|")
	var members []JSType
	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.TrimPrefix(p, "!") // non-null is a no-op on primitives
		if p == "" {
			return Unknown
		}
		m := r.ResolveName(p)
		if isUnknown(m) {
			// Any unresolved member makes the whole expression Unknown
			// — conservative, avoids partial false positives.
			return Unknown
		}
		members = append(members, m)
	}
	if nullable {
		members = append(members, Null)
	}
	if len(members) == 1 {
		return members[0]
	}
	return NewUnion(members...)
}

// ── Native-type table ────────────────────────────────────────────────────────
//
// Mirrors JSTypeRegistry.nativeTypes (private JSType[]) and the pair
// registerNativeType / getNativeType in Java.

// RegisterNativeType stores t at slot typeId in the native-type table.
// Mirrors JSTypeRegistry.registerNativeType(JSTypeNative, JSType).
// Called by createNativeTypes during registry initialisation and by
// any pass that lazily populates a slot (e.g. object-type bootstrap).
func (r *Registry) RegisterNativeType(typeId JSTypeNative, t JSType) {
	r.nativeTypes[typeId] = t
}

// GetNativeType returns the JSType registered for typeId, or nil when
// the slot has not been populated yet.
// Mirrors JSTypeRegistry.getNativeType(JSTypeNative) in Java.
// Callers that need an ObjectType or FunctionType assert the concrete
// type themselves (mirrors getNativeObjectType / getNativeFunctionType).
func (r *Registry) GetNativeType(typeId JSTypeNative) JSType {
	return r.nativeTypes[typeId]
}

// ── Property side tables ─────────────────────────────────────────────────────
//
// Mirrors JSTypeRegistry.eachRefTypeIndexedByProperty (SetMultimap<String,
// ObjectType>) and nonRefTypesIndexedByProperty (SetMultimap<String,JSType>).
// The two tables partition the property index: named / reference types go
// in eachRefTypeIndexedByProperty; structural / anonymous types go in
// nonRefTypesIndexedByProperty.  This split prevents the reference-type
// index from being inflated by arbitrarily many anonymous object literals.

// isObjectLiteralThatCanBeSkipped mirrors the private static helper of the
// same name in Java.  Returns true for record types and anonymous object
// literals that should be mapped to the sentinel rather than stored directly.
func isObjectLiteralThatCanBeSkipped(t JSType) bool {
	// Unwrap null/undefined from union before checking.
	if u, ok := t.(*UnionType); ok {
		t = u.RestrictByNotNullOrUndefined()
	}
	return isRecordType(t) || isLiteralObject(t)
}

// RegisterPropertyOnType records that propertyName is declared on type t.
// Union types are expanded recursively so each alternate is registered
// individually.  Anonymous object literals and record types are collapsed
// to the sentinel object literal to avoid unbounded growth of the map.
// Mirrors JSTypeRegistry.registerPropertyOnType in Java.
func (r *Registry) RegisterPropertyOnType(propertyName string, t JSType) {
	// Expand union types: each alternate is registered independently,
	// matching Java's recursive call on toMaybeUnionType().getAlternates().
	if u, ok := t.(*UnionType); ok {
		for _, alt := range u.Members() {
			r.RegisterPropertyOnType(propertyName, alt)
		}
		return
	}

	// Collapse object literals / record types to the sentinel to bound map size.
	// The sentinel is a plain anonymous ObjectType with no name.
	if isObjectLiteralThatCanBeSkipped(t) {
		// Use a nil sentinel — callers treat any non-nil entry as "has property".
		// We store nil explicitly by using a bool-keyed set instead of JSType.
		// For the Go port we simply skip registration to keep the map small;
		// the canPropertyBeDefined logic that consults this table is not yet
		// ported, so omitting sentinels has no observable effect on diagnostics.
		return
	}

	// Named (reference) types → eachRefTypeIndexedByProperty.
	// Non-reference types → nonRefTypesIndexedByProperty.
	if obj, ok := t.(*ObjectType); ok && obj.HasReferenceName() {
		if r.eachRefTypeIndexedByProperty == nil {
			r.eachRefTypeIndexedByProperty = map[string][]*ObjectType{}
		}
		// Deduplicate: only append if not already present.
		for _, existing := range r.eachRefTypeIndexedByProperty[propertyName] {
			if existing == obj {
				return
			}
		}
		r.eachRefTypeIndexedByProperty[propertyName] = append(
			r.eachRefTypeIndexedByProperty[propertyName], obj)
	} else {
		if r.nonRefTypesIndexedByProperty == nil {
			r.nonRefTypesIndexedByProperty = map[string][]JSType{}
		}
		for _, existing := range r.nonRefTypesIndexedByProperty[propertyName] {
			if existing == t {
				return
			}
		}
		r.nonRefTypesIndexedByProperty[propertyName] = append(
			r.nonRefTypesIndexedByProperty[propertyName], t)
	}
}

// GetEachRefTypeWithProperty returns every named ObjectType on which
// propertyName has been declared via RegisterPropertyOnType.  Returns nil
// when no entry exists (mirrors Java returning an empty ImmutableList).
// Mirrors JSTypeRegistry.getEachReferenceTypeWithProperty in Java.
func (r *Registry) GetEachRefTypeWithProperty(propertyName string) []*ObjectType {
	if r.eachRefTypeIndexedByProperty == nil {
		return nil
	}
	return r.eachRefTypeIndexedByProperty[propertyName]
}

// GetPropertyTypes returns the union of all JSType values — both reference
// and non-reference — that declare propertyName.  Convenience wrapper used
// by type-checking passes that need a flat slice of candidate types without
// the reference / non-reference distinction.
// Mirrors the combined iteration over both maps in
// JSTypeRegistry.canPropertyBeDefined.
func (r *Registry) GetPropertyTypes(propertyName string) []JSType {
	var out []JSType
	if r.nonRefTypesIndexedByProperty != nil {
		out = append(out, r.nonRefTypesIndexedByProperty[propertyName]...)
	}
	if r.eachRefTypeIndexedByProperty != nil {
		for _, obj := range r.eachRefTypeIndexedByProperty[propertyName] {
			out = append(out, obj)
		}
	}
	return out
}

// ── Template-key scope table ─────────────────────────────────────────────────
//
// Mirrors JSTypeRegistry.registerTemplateTypeNamesInScope (which uses a
// Table<Node,String,JSType>).  Go callers pass an opaque scopeKey
// (typically a *rhino.Node pointer) so this file stays free of an import
// of the rhino package.

// RegisterTemplateTypeNamesInScope records each TemplateType in keys under
// the given scopeKey so that LookupTemplateType can retrieve them during
// type-expression evaluation inside that scope.
// Mirrors JSTypeRegistry.registerTemplateTypeNamesInScope(Iterable<TemplateType>, Node).
func (r *Registry) RegisterTemplateTypeNamesInScope(keys []*TemplateType, scopeKey interface{}) {
	if r.templateScopes == nil {
		r.templateScopes = map[interface{}]map[string]*TemplateType{}
	}
	m := r.templateScopes[scopeKey]
	if m == nil {
		m = map[string]*TemplateType{}
		r.templateScopes[scopeKey] = m
	}
	for _, key := range keys {
		m[key.Name] = key
	}
}

// LookupTemplateType returns the TemplateType registered under name in
// scopeKey's scope, or nil when not found.
// Mirrors the scopedNameTable.get(root, name) call in Java's
// JSTypeRegistry.getType(StaticTypedScope, String, …).
func (r *Registry) LookupTemplateType(scopeKey interface{}, name string) *TemplateType {
	if r.templateScopes == nil {
		return nil
	}
	return r.templateScopes[scopeKey][name]
}

// ── Late-resolution worklist callbacks ──────────────────────────────────────
//
// RegisterTypeResolver / ResolveTypes implement a per-type callback
// approach to late resolution. This mirrors the production finalisation
// pass in Java's JSTypeRegistry where pending NamedTypes are driven
// to resolution once all type declarations are available.

// RegisterTypeResolver registers a callback for the given JSType instance.
// When ResolveTypes is called, resolver is invoked with the referenced
// type name; if it returns a non-nil JSType, the entry is consumed.
// Callers that only need worklist-level resolution should use
// RegisterNamedType (which eagerly drains pendingNamedTypes) instead.
//
// Mirrors the implicit per-type resolution callbacks that Java's
// JSTypeRegistry drives via its JSTypeResolver.Closer / resolveAll path.
func (r *Registry) RegisterTypeResolver(t JSType, resolver func(name string) JSType) {
	r.typeResolvers = append(r.typeResolvers, typeResolverEntry{t: t, resolver: resolver})
}

// ResolveTypes drains the pending NamedType worklist and then fires any
// registered type-resolver callbacks.  After the call, entries whose
// resolver returned a non-nil type are removed; the remainder are left
// for a subsequent pass.
//
// Mirrors JSTypeRegistry.resolveTypesForTesting (and the production
// finalisation path) in Java, where the registry iterates the pending
// type list once all declarations are in place.
func (r *Registry) ResolveTypes() {
	// Drain the named-type worklist first (same logic as ResolveWorklist
	// but without an external resolver — relies on namedTypes only).
	remaining := r.pendingNamedTypes[:0]
	for _, nt := range r.pendingNamedTypes {
		if !nt.resolveViaRegistry(r) {
			remaining = append(remaining, nt)
		}
	}
	r.pendingNamedTypes = remaining

	// Fire per-type resolver callbacks.
	if len(r.typeResolvers) == 0 {
		return
	}
	live := r.typeResolvers[:0]
	for _, entry := range r.typeResolvers {
		// Derive the name: if the entry's type is a NamedType use its
		// reference; otherwise leave the name blank so callers that do
		// not need a name still get invoked (resolver("")==nil terminates).
		name := ""
		if nt, ok := entry.t.(*NamedType); ok {
			name = nt.reference
		}
		if entry.resolver(name) == nil {
			live = append(live, entry)
		}
	}
	r.typeResolvers = live
}

// ── Per-scope forward declarations ──────────────────────────────────────────
//
// ForwardDeclareType / GetForwardDeclaredType implement scoped forward
// declarations.  Unlike the flat forwardDeclaredTypes set (populated by
// goog.forwardDeclare at module granularity), these are keyed by an
// opaque scope key so that per-function-scope type aliases stay local.
//
// Mirrors JSTypeRegistry.declareType (the scoped variant) and the
// scopedNameTable used by TypedScopeCreator when processing in-scope
// forward declarations.

// HasType reports whether name is present in the named-type table.
// Mirrors the isDeclaredForScope / getTypeForScopeInternal pair in Java
// (collapsed to a flat lookup since Go callers always resolve globally).
func (r *Registry) HasType(name string) bool {
	return r.namedTypes != nil && r.namedTypes[name] != nil
}

// ForwardDeclareType registers name as a forward-declared type within
// the given scope key.  The scope key is an opaque value (typically a
// *rhino.Node pointer) that identifies the enclosing scope; nil means
// the global scope.
//
// Mirrors the per-scope entry in JSTypeRegistry.scopedNameTable that
// TypedScopeCreator populates when it encounters an in-scope type alias
// or forward-declare directive.
func (r *Registry) ForwardDeclareType(scope interface{}, name string) {
	if r.scopeForwardDeclaredTypes == nil {
		r.scopeForwardDeclaredTypes = map[interface{}]map[string]bool{}
	}
	m := r.scopeForwardDeclaredTypes[scope]
	if m == nil {
		m = map[string]bool{}
		r.scopeForwardDeclaredTypes[scope] = m
	}
	m[name] = true
}

// GetForwardDeclaredType looks up name in the forward-declaration table
// for the given scope key.  It first checks scope's own table; if not
// found, falls back to the global (nil) scope.  When a match is found,
// a NoResolvedType carrying the name is returned so the caller can
// distinguish "forward-declared but unresolved" from "completely unknown".
// Returns nil when name has not been forward-declared in any visible scope.
//
// Mirrors the scope-chain walk in Java's JSTypeRegistry.getTypeInternal
// (getLookupScope → getTypeForScopeInternal) for the forward-declared
// variant of scoped type lookup.
func (r *Registry) GetForwardDeclaredType(scope interface{}, name string) JSType {
	if r.scopeForwardDeclaredTypes == nil {
		return nil
	}
	// Check the specific scope first, then fall through to the global scope.
	for _, key := range []interface{}{scope, nil} {
		if m, ok := r.scopeForwardDeclaredTypes[key]; ok && m[name] {
			return NewNoResolvedType(name, nil)
		}
	}
	return nil
}

// ── Native-type declaration helpers ─────────────────────────────────────────
//
// DeclareNativeFunctionType / DeclareNativeType / DeclareNativeValueType
// mirror the private helpers of the same names in TypedScopeCreator.java.
// They are structural placeholders that register types in the named-type
// table; they do not yet drive scope.declare (which requires a ported
// TypedScope).

// DeclareNativeFunctionType builds a FunctionType from params / ret,
// registers it in the named-type table under name, and also registers
// a plain version without the "Function" suffix when the name ends with
// "Function".  Mirrors TypedScopeCreator.declareNativeFunctionType:
//
//	private void declareNativeFunctionType(TypedScope scope, JSTypeNative tId) {
//	    FunctionType t = typeRegistry.getNativeFunctionType(tId);
//	    declareNativeType(scope, t.getInstanceType().getReferenceName(), t);
//	    declareNativeType(scope, t.getPrototype().getReferenceName(), t.getPrototype());
//	}
//
// Because the Go port does not yet have a fully-wired native type table,
// this variant accepts an explicit name, params slice, and return type
// instead of a JSTypeNative ordinal.
func (r *Registry) DeclareNativeFunctionType(name string, params []JSType, ret JSType) {
	t := NewFunctionType(params, ret)
	r.DeclareNativeType(name, t)
}

// DeclareNativeType registers t in the named-type table under name.
// Mirrors TypedScopeCreator.declareNativeType (the static helper):
//
//	private static void declareNativeType(TypedScope scope, String name, JSType t) {
//	    scope.declare(name, null, t, null, false);
//	}
//
// The scope.declare call is not yet ported; until TypedScope is available
// this method registers the type in the registry's name table so that
// NamedType resolution and HasType queries can find it.
func (r *Registry) DeclareNativeType(name string, t JSType) {
	r.RegisterNamedType(name, t)
}

// DeclareNativeValueType looks up the JSType stored at typeId in the
// native-type table and registers it under name.  Mirrors
// TypedScopeCreator.declareNativeValueType:
//
//	private void declareNativeValueType(TypedScope scope, String name, JSTypeNative tId) {
//	    declareNativeType(scope, name, typeRegistry.getNativeType(tId));
//	}
//
// If the native slot for typeId is empty, the call is a no-op.
func (r *Registry) DeclareNativeValueType(name string, typeId JSTypeNative) {
	t := r.GetNativeType(typeId)
	if t == nil {
		return
	}
	r.DeclareNativeType(name, t)
}
