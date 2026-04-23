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
