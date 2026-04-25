// Port of com/google/javascript/rhino/jstype/NamedType.java @ closure-compiler v20260415
//
// NamedType is a forward-reference placeholder: a type that carries a
// string name but whose resolved JSType is filled in later, once the
// full type graph has been built. Until resolveViaRegistry succeeds,
// the type behaves like Unknown so analysis stays conservative.
//
// Subset ported: the fields, ResolutionKind, getReferenceName,
// resolveViaRegistry, and the TypeName / IsSubtypeOf contract. The
// full resolveInternal (TYPEOF, StaticTypedScope lookup, cycle checks)
// is stubbed — they are exercised only by Tier-5+ tests; TODO markers
// call the specific Java methods being deferred.

package jstype

import "strings"

// ResolutionKind mirrors NamedType.ResolutionKind in Java.
// NONE means the referenced type is already known (no lookup needed).
// TYPE_NAME means the name should be looked up in the registry.
// TYPEOF means the name should be resolved as a typeof expression.
type ResolutionKind int

const (
	ResolutionKindNone     ResolutionKind = iota // Java: NONE
	ResolutionKindTypeName                        // Java: TYPE_NAME
	ResolutionKindTypeof                          // Java: TYPEOF
)

// NamedType is a named reference to some other type.  This provides a
// convenient mechanism for implementing forward references to types; a
// NamedType can be used as a placeholder until its reference is
// resolved.  It is also useful for representing type names in JsDoc
// type annotations, some of which may never be resolved (e.g. types
// in host systems not yet supported by JSCompiler).
//
// Port of com/google/javascript/rhino/jstype/NamedType.java.
type NamedType struct {
	// reference is the textual name of the referenced type, e.g. "Foo"
	// or "goog.Foo". Mirrors NamedType.reference in Java.
	reference string

	// sourceName / lineno / charno carry the location of the JSDoc
	// annotation that produced this type, for error reporting.
	// Mirrors NamedType.{sourceName,lineno,charno}.
	sourceName string
	lineno     int
	charno     int

	// resolutionKind controls how the type is resolved.
	// Mirrors NamedType.resolutionKind.
	resolutionKind ResolutionKind

	// referencedType is nil until resolveViaRegistry succeeds.  After
	// resolution it holds the looked-up JSType.  Mirrors the
	// referencedType field on Java's ProxyObjectType superclass.
	referencedType JSType

	// restrictByNull mirrors NamedType.restrictByNull — true when the
	// "!" operator (getBangType) has been applied before resolution.
	restrictByNull bool

	// templateTypes mirrors NamedType.templateTypes.  Stored but ignored
	// during resolution for now (Tier-5+; TODO: port the template-args
	// substitution from NamedType.resolveInternal).
	templateTypes []JSType
}

// NamedTypeBuilder mirrors the Builder inner class in Java.  Call
// NamedTypeBuilder.Build() to produce the NamedType.
type NamedTypeBuilder struct {
	registry       *Registry
	referenceName  string
	resolutionKind ResolutionKind
	sourceName     string
	lineno         int
	charno         int
	referencedType JSType
	restrictByNull bool
	templateTypes  []JSType
}

// NewNamedTypeBuilder is the factory for NamedTypeBuilder.  Mirrors
// NamedType.builder(registry, reference) in Java.
// Go-only: language-adapter
func NewNamedTypeBuilder(registry *Registry, referenceName string) *NamedTypeBuilder {
	return &NamedTypeBuilder{
		registry:       registry,
		referenceName:  referenceName,
		resolutionKind: ResolutionKindTypeName,
	}
}

// SetResolutionKind mirrors Builder.setResolutionKind.
func (b *NamedTypeBuilder) SetResolutionKind(k ResolutionKind) *NamedTypeBuilder {
	b.resolutionKind = k
	return b
}

// SetErrorReportingLocation mirrors Builder.setErrorReportingLocation.
func (b *NamedTypeBuilder) SetErrorReportingLocation(sourceName string, lineno, charno int) *NamedTypeBuilder {
	b.sourceName = sourceName
	b.lineno = lineno
	b.charno = charno
	return b
}

// SetReferencedType mirrors Builder.setReferencedType.  Used when
// ResolutionKind is NONE (the type is already known at construction).
func (b *NamedTypeBuilder) SetReferencedType(t JSType) *NamedTypeBuilder {
	b.referencedType = t
	return b
}

// SetRestrictByNull mirrors Builder.setRestrictByNull (private in Java;
// exposed here for getBangType / toBuilder).
func (b *NamedTypeBuilder) SetRestrictByNull(v bool) *NamedTypeBuilder {
	b.restrictByNull = v
	return b
}

// SetTemplateTypes mirrors Builder.setTemplateTypes.
func (b *NamedTypeBuilder) SetTemplateTypes(tpls []JSType) *NamedTypeBuilder {
	b.templateTypes = tpls
	return b
}

// Build mirrors Builder.build().  When resolutionKind is not NONE and
// the registry has already seen a resolved type for this name, the
// NamedType is immediately resolved.  Otherwise it is placed on the
// registry's late-resolve worklist.
func (b *NamedTypeBuilder) Build() *NamedType {
	nt := &NamedType{
		reference:      b.referenceName,
		sourceName:     b.sourceName,
		lineno:         b.lineno,
		charno:         b.charno,
		resolutionKind: b.resolutionKind,
		referencedType: b.referencedType,
		restrictByNull: b.restrictByNull,
		templateTypes:  b.templateTypes,
	}
	if b.registry != nil {
		// Eagerly resolve when the registry already knows the type; if
		// not, enqueue for late resolution. Mirrors the
		// registry.getResolver().resolveIfClosed(this, TYPE_CLASS) call
		// in Java's NamedType constructor.
		b.registry.enqueueOrResolve(nt)
	}
	return nt
}

// GetReferenceName mirrors NamedType.getReferenceName() in Java.
// Returns the textual name of the referenced type.
func (n *NamedType) GetReferenceName() string {
	return n.reference
}

// IsResolved mirrors JSType.isResolved() — true once resolveViaRegistry
// has placed a concrete type in referencedType.
func (n *NamedType) IsResolved() bool {
	return n.referencedType != nil
}

// GetReferencedType mirrors NamedType.getReferencedType() / the Java
// ProxyObjectType.getReferencedTypeInternal() accessor.  Returns Unknown
// when still unresolved.
func (n *NamedType) GetReferencedType() JSType {
	if n.referencedType == nil {
		return Unknown
	}
	return n.referencedType
}

// setReferencedType mirrors ProxyObjectType.setReferencedType.
func (n *NamedType) setReferencedType(t JSType) {
	n.referencedType = t
}

// TypeName mirrors JSType.toString() / NamedType.appendTo().  Follows
// Java: when resolved, delegate; when unresolved, return the reference
// string.
// Go-only: language-adapter
func (n *NamedType) TypeName() string {
	if n.referencedType != nil {
		return n.referencedType.TypeName()
	}
	return n.reference
}

// IsSubtypeOf mirrors NamedType → ProxyObjectType → JSType.isSubtypeOf.
// An unresolved NamedType acts like Unknown (conservative: any other type
// is accepted). Once resolved it delegates to the concrete type.
// Go-only: language-adapter
func (n *NamedType) IsSubtypeOf(other JSType) bool {
	if n.referencedType != nil {
		return n.referencedType.IsSubtypeOf(other)
	}
	// Unresolved: behave like Unknown — a subtype of everything so that
	// forward-referenced types don't spuriously fail checks before the
	// type graph is complete.
	return true
}

// GetBangType mirrors NamedType.getBangType().  Returns a non-null
// version of this type by setting restrictByNull and rebuilding, or
// (if already resolved) by delegating to the referenced type.
//
// TODO: wire through the full restrictByNotNullOrUndefined semantics
// from JSType.java when Tier-5 resolution lands.
func (n *NamedType) GetBangType() JSType {
	if n.restrictByNull {
		return n
	}
	if n.IsResolved() {
		return n.referencedType
	}
	// Return a copy with restrictByNull = true.
	cp := *n
	cp.restrictByNull = true
	return &cp
}

// GetTemplateTypes mirrors NamedType.getTemplateTypes().
func (n *NamedType) GetTemplateTypes() []JSType {
	return n.templateTypes
}

// resolveViaClosureNamespace mirrors the resolveViaClosureNamespace branch
// within NamedType.resolveInternal() in Java. Attempts to resolve a dotted
// reference (e.g. "goog.module.Foo") through the closure-namespace registry.
// Returns true when the lookup succeeds and the NamedType is now resolved.
// Only called for TYPE_NAME resolution when the reference contains a dot.
// Go-only: language-adapter
func (n *NamedType) resolveViaClosureNamespace(reg *Registry) bool {
	if n.resolutionKind != ResolutionKindTypeName {
		return false
	}
	if !strings.Contains(n.reference, ".") {
		return false
	}
	t := reg.ResolveViaClosureNamespace(n.reference)
	if t == nil {
		return false
	}
	if n.restrictByNull {
		// TODO(Tier-5): call t.restrictByNotNullOrUndefined().
	}
	n.referencedType = t
	return true
}

// resolveViaRegistry mirrors NamedType.resolveViaRegistry().  Looks up
// the reference name in the registry's type table.  Returns true when
// the lookup succeeds and the NamedType is now resolved.
//
// When the lookup fails, handleUnresolvedType is called: forward-declared
// names (goog.forwardDeclare) silently receive a NoResolvedType so they
// are assignment-compatible everywhere, while truly unknown names stay
// unresolved (Unknown) so callers can emit JSC_UNRECOGNIZED_TYPE_ERROR.
// Mirrors NamedType.resolveViaRegistry / handleUnresolvedType in Java.
func (n *NamedType) resolveViaRegistry(reg *Registry) bool {
	if n.resolutionKind == ResolutionKindNone {
		// Already has a referenced type; nothing to look up.
		return n.referencedType != nil
	}
	t := reg.lookupNamedType(n.reference)
	if t == nil {
		// handleUnresolvedType: if forward-declared, store NoResolvedType
		// so the name is silently accepted; otherwise leave unresolved.
		// Mirrors NamedType.handleUnresolvedType in Java.
		if reg.IsForwardDeclaredType(n.reference) {
			n.referencedType = NewNoResolvedType(n.reference, n.templateTypes)
			return true
		}
		return false
	}
	if n.restrictByNull {
		// TODO(Tier-5): call t.restrictByNotNullOrUndefined() here.
	}
	n.referencedType = t
	return true
}

// resolveTypeof mirrors NamedType.resolveTypeof().  TYPEOF resolution
// requires scope lookup and is deferred to Tier-5.
//
// TODO(Tier-5): implement typeof resolution via StaticTypedScope.
func (n *NamedType) resolveTypeof(_ *Registry) bool {
	// TYPEOF resolution is not yet implemented.
	return false
}

// warning mirrors NamedType.warning().  Emits a diagnostic at the
// location captured when this NamedType was constructed.
//
// TODO(Tier-5): hook into the project's ErrorReporter abstraction.
func (n *NamedType) warning(_ string) {
	// No-op until the ErrorReporter port lands.
}

// Reference / SourceName / Lineno / Charno expose the worklist fields.
func (n *NamedType) Reference() string  { return n.reference }
func (n *NamedType) SourceName() string { return n.sourceName }
func (n *NamedType) Lineno() int        { return n.lineno }
func (n *NamedType) Charno() int        { return n.charno }
