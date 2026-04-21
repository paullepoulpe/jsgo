// Port of com/google/javascript/rhino/jstype/NamedType.java
// @ closure-compiler v20260415 (narrow subset).
//
// In Java, NamedType is a ProxyObjectType subclass: a placeholder for
// a named reference that hasn't been resolved yet. The registry hands
// one out at JSDoc-parse time; a late-resolve pass (post-
// TypedScopeCreator) walks the registry and either substitutes the
// ReferencedType slot with a resolved JSType or emits the diagnostic
// JSC_UNRECOGNIZED_TYPE_ERROR at the NamedType's recorded source
// position.
//
// This Go port captures the two state machines Tier-3 actually needs:
//   - an unresolved NamedType behaves like Unknown in the subtype
//     lattice (no false positives while the scope is still being built,
//     mirroring ProxyObjectType.isSubtype's permissive default);
//   - once ReferencedType is set, IsSubtypeOf delegates to the proxied
//     type so downstream checks see the real shape.
//
// The richer Java NamedType also tracks template arguments, property
// continuations, and a ResolutionScope; those land with Tier-4+ generic
// instantiation. The current port is enough to thread forward
// references through scope.ResolveExpr and report unresolved names
// with JSDoc-position accuracy.

package jstype

// NamedType is an unresolved named reference. `Reference` is the
// source-level identifier (possibly qualified: `foo.bar.Baz`).
// `Lineno` and `Charno` are the absolute coordinates of the reference
// in the source — matching Node.getLineno (1-indexed) and
// Node.getCharno (0-indexed) semantics.
//
// `ReferencedType` starts nil; the late-resolve pass fills it with the
// resolved JSType when it finds one. `RestrictByNull` mirrors the
// leading `!` modifier: once resolved, the non-null arms of the type
// survive.
//
// Mirrors NamedType.java's reference / sourceName / lineno / charno
// / referencedType / restrictByNull fields.
type NamedType struct {
	Reference      string
	SourceName     string
	Lineno         int
	Charno         int
	ReferencedType JSType
	RestrictByNull bool
}

// NewNamedType constructs a NamedType with the given reference name
// and source position. Mirrors NamedType.Builder.build() for the
// forward-reference case (ResolutionKind.TYPE_NAME in Java).
func NewNamedType(reference, sourceName string, lineno, charno int) *NamedType {
	return &NamedType{
		Reference:  reference,
		SourceName: sourceName,
		Lineno:     lineno,
		Charno:     charno,
	}
}

// GetReferenceName mirrors NamedType.getReferenceName().
func (n *NamedType) GetReferenceName() string { return n.Reference }

// GetReferencedType mirrors NamedType.getReferencedType(): returns
// the resolved type, or Unknown while unresolved. Matches Java's
// ProxyObjectType.getReferencedType fallback.
func (n *NamedType) GetReferencedType() JSType {
	if n.ReferencedType == nil {
		return Unknown
	}
	return n.ReferencedType
}

// IsResolved reports whether the NamedType has been bound to a real
// JSType. Mirrors Java's isResolved() check on JSType.
func (n *NamedType) IsResolved() bool { return n.ReferencedType != nil }

// Resolve records the resolved referent. Mirrors
// NamedType.setReferencedAndResolvedType's main assignment.
func (n *NamedType) Resolve(t JSType) { n.ReferencedType = t }

// TypeName returns the reference name while unresolved, and delegates
// to the resolved type once bound. Mirrors NamedType.appendTo which
// prints the reference name until resolution, then the resolved
// TypeStringBuilder form.
func (n *NamedType) TypeName() string {
	if n.ReferencedType != nil {
		return n.ReferencedType.TypeName()
	}
	return n.Reference
}

// IsSubtypeOf: an unresolved NamedType is permissive (treated like
// Unknown) so downstream assignment checks don't false-positive while
// the scope is still being populated. Once resolved, delegates to the
// referent. Mirrors ProxyObjectType.isSubtypeOf in Java: before
// resolution, the proxy returns true unconditionally; after, it
// forwards to the concrete type.
func (n *NamedType) IsSubtypeOf(o JSType) bool {
	if n.ReferencedType == nil {
		return true
	}
	return n.ReferencedType.IsSubtypeOf(o)
}
