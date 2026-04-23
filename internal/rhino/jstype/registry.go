// Port of com/google/javascript/rhino/jstype/JSTypeRegistry.java @ closure-compiler v20260415
//
// The real registry manages interning of all JSType values, named-type
// resolution, template instantiation, and much more. Tier-0 only needs a
// name-to-primitive lookup, which we provide here.
//
// Additive: namedTypes map + pendingNamedTypes worklist added to support
// NamedType.resolveViaRegistry and the createNamedType factory.

package jstype

import "strings"

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

// enqueueOrResolve is called by NamedTypeBuilder.Build. If the
// registry already knows the name, the NamedType is resolved
// immediately; otherwise it is added to the worklist.
// Mirrors the resolveIfClosed call in Java's NamedType constructor.
func (r *Registry) enqueueOrResolve(nt *NamedType) {
	if nt.resolutionKind == ResolutionKindNone {
		return // already has a referenced type; nothing to do
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

// GetResolver returns the JSTypeResolver associated with this registry, or nil
// if none has been attached yet. Mirrors JSTypeRegistry.getResolver() in Java.
func (r *Registry) GetResolver() *JSTypeResolver { return r.resolver }

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
